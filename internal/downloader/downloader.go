package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/he8um/daryaft/internal/download"
)

type Result struct {
	URL  string
	Path string
}

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
	if len(plan.URLs) != 1 {
		return Result{}, fmt.Errorf("single URL download requires exactly one URL")
	}

	rawURL := plan.URLs[0]
	request, err := newRequest(rawURL)
	if err != nil {
		return Result{}, fmt.Errorf("create download request: %w", err)
	}

	response, err := d.client.Do(request)
	if err != nil {
		return Result{}, fmt.Errorf("download %q: %w", rawURL, err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		return Result{}, fmt.Errorf("download %q failed: server returned %s", rawURL, response.Status)
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

	copyErr := copyAndClose(partial, response.Body)
	if copyErr != nil {
		return Result{}, copyErr
	}

	if err := os.Rename(target.Partial, target.Final); err != nil {
		return Result{}, fmt.Errorf("complete download %q: %w", target.Final, err)
	}

	return Result{
		URL:  rawURL,
		Path: target.Final,
	}, nil
}

func copyAndClose(file *os.File, body io.Reader) error {
	_, copyErr := io.Copy(file, body)
	closeErr := file.Close()

	if copyErr != nil {
		return fmt.Errorf("write partial file %q: %w", file.Name(), copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close partial file %q: %w", file.Name(), closeErr)
	}

	return nil
}
