package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/he8um/daryaft/internal/download"
)

const progressInterval = 200 * time.Millisecond

type Downloader struct {
	client  *http.Client
	sleeper Sleeper
}

func New() *Downloader {
	return &Downloader{
		client:  defaultHTTPClient(),
		sleeper: timerSleep,
	}
}

func NewWithClient(client *http.Client) *Downloader {
	if client == nil {
		client = defaultHTTPClient()
	}
	return &Downloader{
		client:  client,
		sleeper: timerSleep,
	}
}

func (d *Downloader) Download(plan download.Plan) (Result, error) {
	return d.DownloadContext(context.Background(), plan)
}

func (d *Downloader) DownloadWithEvents(plan download.Plan, handler EventHandler) (Result, error) {
	return d.DownloadWithEventsContext(context.Background(), plan, handler)
}

func (d *Downloader) DownloadContext(ctx context.Context, plan download.Plan) (Result, error) {
	return d.DownloadWithEventsContext(ctx, plan, nil)
}

func (d *Downloader) DownloadWithEventsContext(ctx context.Context, plan download.Plan, handler EventHandler) (Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(plan.URLs) != 1 {
		return Result{}, fmt.Errorf("single URL download requires exactly one URL")
	}

	rawURL := plan.URLs[0]
	policy := d.retryPolicy(plan.Retries)
	attempts := policy.MaxAttempts()

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			emitEvent(handler, Event{
				Type:        EventCancelled,
				URL:         rawURL,
				Message:     "Download cancelled. Partial file kept for resume.",
				Error:       ErrCancelled,
				Attempt:     attempt,
				MaxAttempts: attempts,
			})
			return Result{}, ErrCancelled
		}

		result, err := d.downloadAttempt(ctx, plan, handler, attempt, attempts)
		if err == nil {
			return result, nil
		}

		lastErr = err
		if isCancellationError(err) {
			return Result{}, ErrCancelled
		}
		if attempt >= attempts || !IsRetryableError(err) {
			emitEvent(handler, Event{
				Type:        EventFailed,
				URL:         rawURL,
				Error:       err,
				Attempt:     attempt,
				MaxAttempts: attempts,
			})
			return Result{}, err
		}

		nextAttempt := attempt + 1
		delay := BackoffDelay(attempt)
		emitEvent(handler, Event{
			Type:        EventRetrying,
			URL:         rawURL,
			Error:       err,
			Attempt:     nextAttempt,
			MaxAttempts: attempts,
			NextDelay:   delay,
		})
		if err := sleepWithContext(ctx, policy.Sleep, delay); err != nil {
			emitEvent(handler, Event{
				Type:        EventCancelled,
				URL:         rawURL,
				Message:     "Download cancelled. Partial file kept for resume.",
				Error:       ErrCancelled,
				Attempt:     attempt,
				MaxAttempts: attempts,
			})
			return Result{}, ErrCancelled
		}
	}

	return Result{}, lastErr
}

func (d *Downloader) retryPolicy(retries int) RetryPolicy {
	return RetryPolicy{
		Retries: retries,
		Sleep:   d.sleeper,
	}
}

func sleepWithContext(ctx context.Context, sleep Sleeper, delay time.Duration) error {
	if sleep == nil {
		sleep = timerSleep
	}
	return sleep(ctx, delay)
}

func (d *Downloader) downloadAttempt(ctx context.Context, plan download.Plan, handler EventHandler, attempt, attempts int) (Result, error) {
	rawURL := plan.URLs[0]
	candidate, err := findResumeCandidate(plan.Output, rawURL, plan.Name)
	if err != nil {
		return Result{}, err
	}

	resumeOffset := int64(0)
	if plan.Resume && candidate.CanResume {
		resumeOffset = candidate.PartialSize
	}

	request, err := newRequestWithContext(ctx, rawURL)
	if err != nil {
		return Result{}, nonRetryableError{err: fmt.Errorf("create download request: %w", err)}
	}
	if resumeOffset > 0 {
		request.Header.Set("Range", fmt.Sprintf("bytes=%d-", resumeOffset))
	}

	response, err := d.client.Do(request)
	if err != nil {
		if isCancellationError(err) {
			emitCancelledEvent(handler, rawURL, candidate.Target.Final, candidate.Target.Partial, candidate.PartialSize, 0)
			return Result{}, ErrCancelled
		}
		return Result{}, fmt.Errorf("download %q: %w", rawURL, err)
	}

	prepared, err := d.prepareDownloadResponse(ctx, plan, rawURL, response, candidate, resumeOffset, handler)
	if err != nil {
		if isCancellationError(err) {
			emitCancelledEvent(handler, rawURL, candidate.Target.Final, candidate.Target.Partial, candidate.PartialSize, 0)
			return Result{}, ErrCancelled
		}
		return Result{}, err
	}
	response = prepared.response
	defer response.Body.Close()

	var partial *os.File
	if prepared.mode == downloadModeAppend {
		partial, err = os.OpenFile(prepared.target.Partial, os.O_WRONLY|os.O_APPEND, 0o600)
	} else {
		partial, err = os.Create(prepared.target.Partial)
	}
	if err != nil {
		return Result{}, fmt.Errorf("open partial file %q: %w", prepared.target.Partial, err)
	}

	metadata := metadataFromResponse(rawURL, prepared.target, response, prepared.totalBytes, prepared.startBytes, candidate)
	if err := savePartialMetadata(prepared.metadataPath, metadata); err != nil {
		_ = partial.Close()
		return Result{}, err
	}

	emitEvent(handler, Event{
		Type:            EventStarted,
		URL:             rawURL,
		TargetPath:      prepared.target.Final,
		PartialPath:     prepared.target.Partial,
		DownloadedBytes: prepared.startBytes,
		TotalBytes:      prepared.totalBytes,
		StatusCode:      response.StatusCode,
		Status:          response.Status,
		Attempt:         attempt,
		MaxAttempts:     attempts,
		Timestamp:       time.Now(),
	})

	downloadedBytes, copyErr := copyAndClose(partial, response.Body, copyEvents{
		URL:          rawURL,
		TargetPath:   prepared.target.Final,
		PartialPath:  prepared.target.Partial,
		InitialBytes: prepared.startBytes,
		TotalBytes:   prepared.totalBytes,
		Handler:      handler,
		Context:      ctx,
	})
	if copyErr != nil {
		metadata.DownloadedBytes = downloadedBytes
		_ = savePartialMetadata(prepared.metadataPath, metadata)
		if isCancellationError(copyErr) {
			emitCancelledEvent(handler, rawURL, prepared.target.Final, prepared.target.Partial, downloadedBytes, prepared.totalBytes)
			return Result{}, ErrCancelled
		}
		return Result{}, copyErr
	}
	metadata.DownloadedBytes = downloadedBytes
	if err := savePartialMetadata(prepared.metadataPath, metadata); err != nil {
		return Result{}, err
	}

	if err := os.Rename(prepared.target.Partial, prepared.target.Final); err != nil {
		return Result{}, fmt.Errorf("complete download %q: %w", prepared.target.Final, err)
	}
	if err := removePartialMetadata(prepared.metadataPath); err != nil {
		return Result{}, err
	}

	emitEvent(handler, Event{
		Type:            EventCompleted,
		URL:             rawURL,
		TargetPath:      prepared.target.Final,
		PartialPath:     prepared.target.Partial,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      prepared.totalBytes,
		Percent:         percent(downloadedBytes, prepared.totalBytes),
		Attempt:         attempt,
		MaxAttempts:     attempts,
		Timestamp:       time.Now(),
	})

	return Result{
		URL:  rawURL,
		Path: prepared.target.Final,
	}, nil
}

type downloadMode int

const (
	downloadModeRestart downloadMode = iota
	downloadModeAppend
)

type preparedDownload struct {
	response     *http.Response
	mode         downloadMode
	target       targetPaths
	metadataPath string
	startBytes   int64
	totalBytes   int64
}

func (d *Downloader) prepareDownloadResponse(ctx context.Context, plan download.Plan, rawURL string, response *http.Response, candidate resumeCandidate, resumeOffset int64, handler EventHandler) (preparedDownload, error) {
	if resumeOffset > 0 {
		switch {
		case response.StatusCode == http.StatusPartialContent:
			if !contentRangeStartsAt(response, resumeOffset) {
				_ = response.Body.Close()
				return preparedDownload{}, fmt.Errorf("server returned unexpected content range %q", response.Header.Get("Content-Range"))
			}
			if candidate.HasMetadata && remoteChanged(candidate.Metadata, response.Header) {
				return d.restartWithNewRequest(ctx, plan, rawURL, response, candidate, handler, remoteChangedMessage)
			}

			totalBytes := responseTotalBytes(response, resumeOffset)
			emitEvent(handler, Event{
				Type:            EventResuming,
				URL:             rawURL,
				TargetPath:      candidate.Target.Final,
				PartialPath:     candidate.Target.Partial,
				DownloadedBytes: resumeOffset,
				TotalBytes:      totalBytes,
				Message:         fmt.Sprintf("Resuming from %d bytes", resumeOffset),
			})
			return preparedDownload{
				response:     response,
				mode:         downloadModeAppend,
				target:       candidate.Target,
				metadataPath: candidate.MetadataPath,
				startBytes:   resumeOffset,
				totalBytes:   totalBytes,
			}, nil

		case response.StatusCode == http.StatusOK:
			return d.restartWithResponse(plan, rawURL, response, candidate, handler, resumeNotSupportedMessage)
		case response.StatusCode == http.StatusRequestedRangeNotSatisfiable:
			return d.restartWithNewRequest(ctx, plan, rawURL, response, candidate, handler, resumeNotSupportedMessage)
		}
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		_ = response.Body.Close()
		return preparedDownload{}, httpStatusError{
			StatusCode: response.StatusCode,
			Status:     response.Status,
		}
	}

	target := candidate.Target
	metadataPath := candidate.MetadataPath
	if target.Final == "" {
		filename := FilenameFromResponse(rawURL, response.Header, plan.Name)
		var err error
		target, err = prepareTarget(plan.Output, filename)
		if err != nil {
			_ = response.Body.Close()
			return preparedDownload{}, err
		}
		metadataPath = metadataPathForPartial(target.Partial)
	} else if _, err := prepareTarget(plan.Output, filepath.Base(target.Final)); err != nil {
		_ = response.Body.Close()
		return preparedDownload{}, err
	}

	return preparedDownload{
		response:     response,
		mode:         downloadModeRestart,
		target:       target,
		metadataPath: metadataPath,
		totalBytes:   responseTotalBytes(response, 0),
	}, nil
}

func (d *Downloader) restartWithResponse(plan download.Plan, rawURL string, response *http.Response, candidate resumeCandidate, handler EventHandler, message string) (preparedDownload, error) {
	emitEvent(handler, Event{
		Type:            EventRestarting,
		URL:             rawURL,
		TargetPath:      candidate.Target.Final,
		PartialPath:     candidate.Target.Partial,
		DownloadedBytes: candidate.PartialSize,
		Message:         message,
	})

	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		_ = response.Body.Close()
		return preparedDownload{}, httpStatusError{
			StatusCode: response.StatusCode,
			Status:     response.Status,
		}
	}

	return preparedDownload{
		response:     response,
		mode:         downloadModeRestart,
		target:       candidate.Target,
		metadataPath: candidate.MetadataPath,
		totalBytes:   responseTotalBytes(response, 0),
	}, nil
}

func (d *Downloader) restartWithNewRequest(ctx context.Context, plan download.Plan, rawURL string, oldResponse *http.Response, candidate resumeCandidate, handler EventHandler, message string) (preparedDownload, error) {
	_ = oldResponse.Body.Close()

	emitEvent(handler, Event{
		Type:            EventRestarting,
		URL:             rawURL,
		TargetPath:      candidate.Target.Final,
		PartialPath:     candidate.Target.Partial,
		DownloadedBytes: candidate.PartialSize,
		Message:         message,
	})

	request, err := newRequestWithContext(ctx, rawURL)
	if err != nil {
		return preparedDownload{}, nonRetryableError{err: fmt.Errorf("create download request: %w", err)}
	}
	response, err := d.client.Do(request)
	if err != nil {
		if isCancellationError(err) {
			return preparedDownload{}, ErrCancelled
		}
		return preparedDownload{}, fmt.Errorf("download %q: %w", rawURL, err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		_ = response.Body.Close()
		return preparedDownload{}, httpStatusError{
			StatusCode: response.StatusCode,
			Status:     response.Status,
		}
	}

	return preparedDownload{
		response:     response,
		mode:         downloadModeRestart,
		target:       candidate.Target,
		metadataPath: candidate.MetadataPath,
		totalBytes:   responseTotalBytes(response, 0),
	}, nil
}

func metadataFromResponse(rawURL string, target targetPaths, response *http.Response, totalBytes, downloadedBytes int64, candidate resumeCandidate) partialMetadata {
	metadata := candidate.Metadata
	metadata.URL = rawURL
	metadata.TargetPath = target.Final
	metadata.PartialPath = target.Partial
	metadata.TotalBytes = totalBytes
	metadata.DownloadedBytes = downloadedBytes
	metadata.ETag = response.Header.Get("ETag")
	metadata.LastModified = response.Header.Get("Last-Modified")
	metadata.AcceptRanges = response.Header.Get("Accept-Ranges")
	return metadata
}

type copyEvents struct {
	URL          string
	TargetPath   string
	PartialPath  string
	InitialBytes int64
	TotalBytes   int64
	Handler      EventHandler
	Context      context.Context
}

func copyAndClose(file *os.File, body io.Reader, events copyEvents) (int64, error) {
	downloadedBytes, copyErr := copyWithProgress(file, body, events)
	closeErr := file.Close()

	if copyErr != nil {
		return downloadedBytes, fmt.Errorf("write partial file %q: %w", file.Name(), copyErr)
	}
	if closeErr != nil {
		return downloadedBytes, fmt.Errorf("close partial file %q: %w", file.Name(), closeErr)
	}

	return downloadedBytes, nil
}

func copyWithProgress(file *os.File, body io.Reader, events copyEvents) (int64, error) {
	ctx := events.Context
	if ctx == nil {
		ctx = context.Background()
	}
	buffer := make([]byte, 64*1024)
	startedAt := time.Now()
	lastProgressAt := time.Time{}
	var lastProgressBytes int64
	downloadedBytes := events.InitialBytes

	for {
		if err := ctx.Err(); err != nil {
			return downloadedBytes, ErrCancelled
		}

		n, readErr := body.Read(buffer)
		if isCancellationError(readErr) {
			return downloadedBytes, ErrCancelled
		}
		if n > 0 {
			if err := ctx.Err(); err != nil {
				return downloadedBytes, ErrCancelled
			}
			written, writeErr := file.Write(buffer[:n])
			downloadedBytes += int64(written)
			if writeErr != nil {
				return downloadedBytes, writeErr
			}
			if written != n {
				return downloadedBytes, io.ErrShortWrite
			}

			now := time.Now()
			if lastProgressAt.IsZero() || now.Sub(lastProgressAt) >= progressInterval {
				emitEvent(events.Handler, newProgressEvent(events.URL, events.TargetPath, events.PartialPath, downloadedBytes, downloadedBytes-events.InitialBytes, events.TotalBytes, startedAt))
				lastProgressAt = now
				lastProgressBytes = downloadedBytes
			}
		}

		if readErr == io.EOF {
			if lastProgressAt.IsZero() || lastProgressBytes != downloadedBytes {
				emitEvent(events.Handler, newProgressEvent(events.URL, events.TargetPath, events.PartialPath, downloadedBytes, downloadedBytes-events.InitialBytes, events.TotalBytes, startedAt))
			}
			return downloadedBytes, nil
		}
		if readErr != nil {
			return downloadedBytes, readErr
		}
	}
}

func emitCancelledEvent(handler EventHandler, rawURL, targetPath, partialPath string, downloadedBytes, totalBytes int64) {
	emitEvent(handler, Event{
		Type:            EventCancelled,
		URL:             rawURL,
		TargetPath:      targetPath,
		PartialPath:     partialPath,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      totalBytes,
		Message:         "Download cancelled. Partial file kept for resume.",
		Error:           ErrCancelled,
	})
}

func isCancellationError(err error) bool {
	return errors.Is(err, ErrCancelled) || errors.Is(err, context.Canceled)
}

func percent(downloadedBytes, totalBytes int64) float64 {
	if totalBytes <= 0 {
		return 0
	}
	return float64(downloadedBytes) / float64(totalBytes) * 100
}
