package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
		client: defaultHTTPClient(),
	}
}

func NewWithClient(client *http.Client) *Downloader {
	if client == nil {
		client = defaultHTTPClient()
	}
	return &Downloader{
		client:  client,
		sleeper: time.Sleep,
	}
}

func (d *Downloader) Download(plan download.Plan) (Result, error) {
	return d.DownloadWithEvents(plan, nil)
}

func (d *Downloader) DownloadWithEvents(plan download.Plan, handler EventHandler) (Result, error) {
	if len(plan.URLs) != 1 {
		return Result{}, fmt.Errorf("single URL download requires exactly one URL")
	}

	rawURL := plan.URLs[0]
	policy := d.retryPolicy(plan.Retries)
	attempts := policy.MaxAttempts()

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		result, err := d.downloadAttempt(plan, handler, attempt, attempts)
		if err == nil {
			return result, nil
		}

		lastErr = err
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
		policy.Sleep(delay)
	}

	return Result{}, lastErr
}

func (d *Downloader) retryPolicy(retries int) RetryPolicy {
	sleep := d.sleeper
	if sleep == nil {
		sleep = time.Sleep
	}
	return RetryPolicy{
		Retries: retries,
		Sleep:   sleep,
	}
}

func (d *Downloader) downloadAttempt(plan download.Plan, handler EventHandler, attempt, attempts int) (Result, error) {
	rawURL := plan.URLs[0]
	totalBytes := int64(0)

	request, err := newRequest(rawURL)
	if err != nil {
		return Result{}, nonRetryableError{err: fmt.Errorf("create download request: %w", err)}
	}

	response, err := d.client.Do(request)
	if err != nil {
		return Result{}, fmt.Errorf("download %q: %w", rawURL, err)
	}
	defer response.Body.Close()

	if response.ContentLength > 0 {
		totalBytes = response.ContentLength
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		return Result{}, httpStatusError{
			StatusCode: response.StatusCode,
			Status:     response.Status,
		}
	}

	filename := FilenameFromResponse(rawURL, response.Header, plan.Name)
	target, err := prepareTarget(plan.Output, filename)
	if err != nil {
		return Result{}, err
	}

	// Resume is not implemented yet; restart any existing partial file for now.
	partial, err := os.Create(target.Partial)
	if err != nil {
		return Result{}, fmt.Errorf("create partial file %q: %w", target.Partial, err)
	}

	emitEvent(handler, Event{
		Type:        EventStarted,
		URL:         rawURL,
		TargetPath:  target.Final,
		TotalBytes:  totalBytes,
		Attempt:     attempt,
		MaxAttempts: attempts,
		Timestamp:   time.Now(),
	})

	downloadedBytes, copyErr := copyAndClose(partial, response.Body, copyEvents{
		URL:        rawURL,
		TargetPath: target.Final,
		TotalBytes: totalBytes,
		Handler:    handler,
	})
	if copyErr != nil {
		return Result{}, copyErr
	}

	if err := os.Rename(target.Partial, target.Final); err != nil {
		return Result{}, fmt.Errorf("complete download %q: %w", target.Final, err)
	}

	emitEvent(handler, Event{
		Type:            EventCompleted,
		URL:             rawURL,
		TargetPath:      target.Final,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      totalBytes,
		Percent:         percent(downloadedBytes, totalBytes),
		Attempt:         attempt,
		MaxAttempts:     attempts,
		Timestamp:       time.Now(),
	})

	return Result{
		URL:  rawURL,
		Path: target.Final,
	}, nil
}

type copyEvents struct {
	URL        string
	TargetPath string
	TotalBytes int64
	Handler    EventHandler
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
	buffer := make([]byte, 64*1024)
	startedAt := time.Now()
	lastProgressAt := time.Time{}
	var lastProgressBytes int64
	var downloadedBytes int64

	for {
		n, readErr := body.Read(buffer)
		if n > 0 {
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
				emitEvent(events.Handler, newProgressEvent(events.URL, events.TargetPath, downloadedBytes, events.TotalBytes, startedAt))
				lastProgressAt = now
				lastProgressBytes = downloadedBytes
			}
		}

		if readErr == io.EOF {
			if lastProgressAt.IsZero() || lastProgressBytes != downloadedBytes {
				emitEvent(events.Handler, newProgressEvent(events.URL, events.TargetPath, downloadedBytes, events.TotalBytes, startedAt))
			}
			return downloadedBytes, nil
		}
		if readErr != nil {
			return downloadedBytes, readErr
		}
	}
}

func percent(downloadedBytes, totalBytes int64) float64 {
	if totalBytes <= 0 {
		return 0
	}
	return float64(downloadedBytes) / float64(totalBytes) * 100
}
