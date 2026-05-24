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
	client *http.Client
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
	return &Downloader{client: client}
}

func (d *Downloader) Download(plan download.Plan) (Result, error) {
	return d.DownloadWithEvents(plan, nil)
}

func (d *Downloader) DownloadWithEvents(plan download.Plan, handler EventHandler) (Result, error) {
	if len(plan.URLs) != 1 {
		return Result{}, fmt.Errorf("single URL download requires exactly one URL")
	}

	rawURL := plan.URLs[0]
	targetPath := ""
	totalBytes := int64(0)
	downloadedBytes := int64(0)

	fail := func(err error) (Result, error) {
		emitEvent(handler, Event{
			Type:            EventFailed,
			URL:             rawURL,
			TargetPath:      targetPath,
			DownloadedBytes: downloadedBytes,
			TotalBytes:      totalBytes,
			Error:           err,
		})
		return Result{}, err
	}

	request, err := newRequest(rawURL)
	if err != nil {
		return fail(fmt.Errorf("create download request: %w", err))
	}

	response, err := d.client.Do(request)
	if err != nil {
		return fail(fmt.Errorf("download %q: %w", rawURL, err))
	}
	defer response.Body.Close()

	if response.ContentLength > 0 {
		totalBytes = response.ContentLength
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		return fail(fmt.Errorf("download %q failed: server returned %s", rawURL, response.Status))
	}

	filename := FilenameFromResponse(rawURL, response.Header, plan.Name)
	target, err := prepareTarget(plan.Output, filename)
	if err != nil {
		return fail(err)
	}
	targetPath = target.Final

	// Resume is not implemented yet; restart any existing partial file for now.
	partial, err := os.Create(target.Partial)
	if err != nil {
		return fail(fmt.Errorf("create partial file %q: %w", target.Partial, err))
	}

	emitEvent(handler, Event{
		Type:       EventStarted,
		URL:        rawURL,
		TargetPath: target.Final,
		TotalBytes: totalBytes,
		Timestamp:  time.Now(),
	})

	downloadedBytes, copyErr := copyAndClose(partial, response.Body, copyEvents{
		URL:        rawURL,
		TargetPath: target.Final,
		TotalBytes: totalBytes,
		Handler:    handler,
	})
	if copyErr != nil {
		return fail(copyErr)
	}

	if err := os.Rename(target.Partial, target.Final); err != nil {
		return fail(fmt.Errorf("complete download %q: %w", target.Final, err))
	}

	emitEvent(handler, Event{
		Type:            EventCompleted,
		URL:             rawURL,
		TargetPath:      target.Final,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      totalBytes,
		Percent:         percent(downloadedBytes, totalBytes),
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
