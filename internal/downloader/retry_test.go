package downloader

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/he8um/daryaft/internal/download"
)

func TestDownloadRetriesZeroMakesOneAttempt(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "try later", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	var retryEvents []Event
	d := newTestDownloader(server)
	_, err := d.DownloadWithEvents(download.Plan{
		URLs:    []string{server.URL + "/file.txt"},
		Output:  t.TempDir(),
		Retries: 0,
	}, func(event Event) {
		if event.Type == EventRetrying {
			retryEvents = append(retryEvents, event)
		}
	})
	if err == nil {
		t.Fatal("DownloadWithEvents returned nil error")
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
	if len(retryEvents) != 0 {
		t.Fatalf("retry events = %d, want 0", len(retryEvents))
	}
}

func TestDownloadRetriesTransient503AndSucceeds(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			http.Error(w, "try later", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	var retryEvents []Event
	d := newTestDownloader(server)
	dir := t.TempDir()
	result, err := d.DownloadWithEvents(download.Plan{
		URLs:    []string{server.URL + "/file.txt"},
		Output:  dir,
		Retries: 3,
	}, func(event Event) {
		if event.Type == EventRetrying {
			retryEvents = append(retryEvents, event)
		}
	})
	if err != nil {
		t.Fatalf("DownloadWithEvents returned error: %v", err)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
	assertFileContent(t, result.Path, "ok")

	if len(retryEvents) != 1 {
		t.Fatalf("retry events = %d, want 1", len(retryEvents))
	}
	event := retryEvents[0]
	if event.Attempt != 2 || event.MaxAttempts != 4 {
		t.Fatalf("retry event attempt = %d/%d, want 2/4", event.Attempt, event.MaxAttempts)
	}
	if event.NextDelay != time.Second {
		t.Fatalf("retry event delay = %s, want 1s", event.NextDelay)
	}
	if event.Error == nil || !strings.Contains(event.Error.Error(), "temporary server error: 503") {
		t.Fatalf("retry event error = %v", event.Error)
	}
}

func TestDownloadCancelDuringRetryBackoffReturnsPromptly(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "try later", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	sleepStarted := make(chan struct{})
	d := NewWithClient(server.Client())
	d.sleeper = func(ctx context.Context, delay time.Duration) error {
		close(sleepStarted)
		<-ctx.Done()
		return ctx.Err()
	}

	errCh := make(chan error, 1)
	go func() {
		_, err := d.DownloadWithEventsContext(ctx, download.Plan{
			URLs:    []string{server.URL + "/file.txt"},
			Output:  t.TempDir(),
			Retries: 1,
		}, nil)
		errCh <- err
	}()

	select {
	case <-sleepStarted:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("retry sleep did not start")
	}
	cancel()

	select {
	case err := <-errCh:
		if !errors.Is(err, ErrCancelled) {
			t.Fatalf("error = %v, want ErrCancelled", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("DownloadWithEventsContext did not return promptly after cancellation")
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
}

func TestDownloadDoesNotRetry404(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.NotFound(w, r)
	}))
	defer server.Close()

	d := newTestDownloader(server)
	_, err := d.DownloadWithEvents(download.Plan{
		URLs:    []string{server.URL + "/missing.txt"},
		Output:  t.TempDir(),
		Retries: 3,
	}, nil)
	if err == nil {
		t.Fatal("DownloadWithEvents returned nil error")
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
	if !strings.Contains(err.Error(), "server returned 404") {
		t.Fatalf("error = %q, want 404 context", err)
	}
}

func TestDownloadDoesNotRetryInvalidURL(t *testing.T) {
	var retryEvents []Event
	d := newTestDownloader(nil)
	_, err := d.DownloadWithEvents(download.Plan{
		URLs:    []string{"://bad-url"},
		Output:  t.TempDir(),
		Retries: 3,
	}, func(event Event) {
		if event.Type == EventRetrying {
			retryEvents = append(retryEvents, event)
		}
	})
	if err == nil {
		t.Fatal("DownloadWithEvents returned nil error")
	}
	if len(retryEvents) != 0 {
		t.Fatalf("retry events = %d, want 0", len(retryEvents))
	}
}

func TestDownloadDoesNotRetryExistingTarget(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte("new"))
	}))
	defer server.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(target, []byte("old"), 0o600); err != nil {
		t.Fatalf("write existing target: %v", err)
	}

	d := newTestDownloader(server)
	_, err := d.DownloadWithEvents(download.Plan{
		URLs:    []string{server.URL + "/file.txt"},
		Output:  dir,
		Retries: 3,
	}, nil)
	if err == nil {
		t.Fatal("DownloadWithEvents returned nil error")
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
	if !strings.Contains(err.Error(), "target file already exists") {
		t.Fatalf("error = %q, want existing target context", err)
	}
	assertFileContent(t, target, "old")
}

func TestDownloadBatchRetriesEachItemIndependently(t *testing.T) {
	requestsByPath := make(map[string]int)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestsByPath[r.URL.Path]++
		switch r.URL.Path {
		case "/flaky.txt":
			if requestsByPath[r.URL.Path] == 1 {
				http.Error(w, "try later", http.StatusServiceUnavailable)
				return
			}
			_, _ = w.Write([]byte("flaky ok"))
		case "/ok.txt":
			_, _ = w.Write([]byte("ok"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	d := newTestDownloader(server)
	result := d.DownloadBatch(download.Plan{
		URLs: []string{
			server.URL + "/flaky.txt",
			server.URL + "/ok.txt",
		},
		Output:  t.TempDir(),
		Retries: 3,
	}, BatchHandlers{})
	if result.Err() != nil {
		t.Fatalf("result.Err() = %v", result.Err())
	}
	if result.Completed() != 2 || result.Failed() != 0 {
		t.Fatalf("counts = completed %d failed %d", result.Completed(), result.Failed())
	}
	if requestsByPath["/flaky.txt"] != 2 {
		t.Fatalf("flaky requests = %d, want 2", requestsByPath["/flaky.txt"])
	}
	if requestsByPath["/ok.txt"] != 1 {
		t.Fatalf("ok requests = %d, want 1", requestsByPath["/ok.txt"])
	}
}

func TestBackoffDelayCapsAtEightSeconds(t *testing.T) {
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{attempt: 1, want: time.Second},
		{attempt: 2, want: 2 * time.Second},
		{attempt: 3, want: 4 * time.Second},
		{attempt: 4, want: 8 * time.Second},
		{attempt: 5, want: 8 * time.Second},
	}

	for _, test := range tests {
		if got := BackoffDelay(test.attempt); got != test.want {
			t.Fatalf("BackoffDelay(%d) = %s, want %s", test.attempt, got, test.want)
		}
	}
}

func TestCancellationErrorsAreNotRetryable(t *testing.T) {
	if IsRetryableError(ErrCancelled) {
		t.Fatal("ErrCancelled classified as retryable")
	}
	if IsRetryableError(context.Canceled) {
		t.Fatal("context.Canceled classified as retryable")
	}
}

func newTestDownloader(server *httptest.Server) *Downloader {
	var client *http.Client
	if server != nil {
		client = server.Client()
	}
	d := NewWithClient(client)
	d.sleeper = func(context.Context, time.Duration) error { return nil }
	return d
}
