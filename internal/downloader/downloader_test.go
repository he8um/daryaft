package downloader

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/pkg/version"
)

func TestDownloadSingleURL(t *testing.T) {
	var userAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.UserAgent()
		w.Header().Set("Content-Disposition", `attachment; filename="server-name.txt"`)
		_, _ = w.Write([]byte("hello daryaft"))
	}))
	defer server.Close()

	dir := t.TempDir()
	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/ignored-path.bin"},
		Output: dir,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	wantPath := filepath.Join(dir, "server-name.txt")
	if result.Path != wantPath {
		t.Fatalf("result.Path = %q, want %q", result.Path, wantPath)
	}

	content, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}
	if string(content) != "hello daryaft" {
		t.Fatalf("downloaded content = %q", content)
	}
	if _, err := os.Stat(wantPath + ".part"); !os.IsNotExist(err) {
		t.Fatalf("partial file still exists or stat failed: %v", err)
	}
	if userAgent != "Daryaft/"+version.Version {
		t.Fatalf("User-Agent = %q", userAgent)
	}
}

func TestDefaultHTTPClientUsesPhaseTimeoutsWithoutTotalTimeout(t *testing.T) {
	client := defaultHTTPClient()
	if client.Timeout != 0 {
		t.Fatalf("client.Timeout = %s, want no total download timeout", client.Timeout)
	}
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("client.Transport = %T, want *http.Transport", client.Transport)
	}
	if transport.ResponseHeaderTimeout != DefaultResponseHeaderTimeout {
		t.Fatalf("ResponseHeaderTimeout = %s, want %s", transport.ResponseHeaderTimeout, DefaultResponseHeaderTimeout)
	}
	if transport.TLSHandshakeTimeout != DefaultTLSHandshakeTimeout {
		t.Fatalf("TLSHandshakeTimeout = %s, want %s", transport.TLSHandshakeTimeout, DefaultTLSHandshakeTimeout)
	}
}

func TestResponseHeaderTimeoutFailsSlowHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)
		_, _ = w.Write([]byte("late"))
	}))
	defer server.Close()

	client := server.Client()
	client.Transport = &http.Transport{
		ResponseHeaderTimeout: 25 * time.Millisecond,
	}
	_, err := NewWithClient(client).Download(download.Plan{
		URLs:   []string{server.URL + "/slow-headers.txt"},
		Output: t.TempDir(),
	})
	if err == nil {
		t.Fatal("Download returned nil error, want header timeout")
	}
}

func TestSlowBodyStreamingIsNotKilledByHeaderTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "6")
		if flusher, ok := w.(http.Flusher); ok {
			_, _ = w.Write([]byte("abc"))
			flusher.Flush()
		}
		time.Sleep(75 * time.Millisecond)
		_, _ = w.Write([]byte("def"))
	}))
	defer server.Close()

	client := server.Client()
	client.Transport = &http.Transport{
		ResponseHeaderTimeout: 25 * time.Millisecond,
	}
	result, err := NewWithClient(client).Download(download.Plan{
		URLs:   []string{server.URL + "/slow-body.txt"},
		Output: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	assertFileContent(t, result.Path, "abcdef")
}

func TestDownloadCancellationLeavesPartialAndMetadata(t *testing.T) {
	body := bytes.Repeat([]byte("a"), 256*1024)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		_, _ = w.Write(body)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var events []Event
	dir := t.TempDir()
	_, err := New().DownloadWithEventsContext(ctx, download.Plan{
		URLs:    []string{server.URL + "/file.bin"},
		Output:  dir,
		Resume:  true,
		Retries: 3,
	}, func(event Event) {
		events = append(events, event)
		if event.Type == EventStarted {
			cancel()
		}
	})
	if !errors.Is(err, ErrCancelled) {
		t.Fatalf("error = %v, want ErrCancelled", err)
	}

	final := filepath.Join(dir, "file.bin")
	partial := final + ".part"
	metadata := partial + ".daryaft.json"
	assertMissing(t, final)
	if _, err := os.Stat(partial); err != nil {
		t.Fatalf("partial missing after cancellation: %v", err)
	}
	if _, err := os.Stat(metadata); err != nil {
		t.Fatalf("metadata missing after cancellation: %v", err)
	}

	cancelled := findEvent(events, EventCancelled)
	if cancelled.Type != EventCancelled {
		t.Fatal("missing cancelled event")
	}
	if cancelled.TargetPath != final {
		t.Fatalf("cancelled.TargetPath = %q, want %q", cancelled.TargetPath, final)
	}
	if cancelled.PartialPath != partial {
		t.Fatalf("cancelled.PartialPath = %q, want %q", cancelled.PartialPath, partial)
	}
	if cancelled.Message != "Download cancelled. Partial file kept for resume." {
		t.Fatalf("cancelled.Message = %q", cancelled.Message)
	}
}

func TestDownloadCancellationIsNotRetried(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write(bytes.Repeat([]byte("a"), 128*1024))
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var retryEvents []Event
	_, err := newTestDownloader(server).DownloadWithEventsContext(ctx, download.Plan{
		URLs:    []string{server.URL + "/file.bin"},
		Output:  t.TempDir(),
		Retries: 3,
	}, func(event Event) {
		if event.Type == EventStarted {
			cancel()
		}
		if event.Type == EventRetrying {
			retryEvents = append(retryEvents, event)
		}
	})
	if !errors.Is(err, ErrCancelled) {
		t.Fatalf("error = %v, want ErrCancelled", err)
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
	if len(retryEvents) != 0 {
		t.Fatalf("retry events = %d, want 0", len(retryEvents))
	}
}

func TestDownloadRejectsNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusTeapot)
	}))
	defer server.Close()

	_, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: t.TempDir(),
	})
	if err == nil {
		t.Fatal("Download returned nil error")
	}
	if !strings.Contains(err.Error(), "server returned 418") {
		t.Fatalf("error = %q, want status context", err)
	}
}

func TestDownloadCreatesOutputDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	dir := filepath.Join(t.TempDir(), "nested", "downloads")
	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	if result.Path != filepath.Join(dir, "file.txt") {
		t.Fatalf("result.Path = %q", result.Path)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("output directory was not created: %v", err)
	}
}

func TestDownloadRejectsExistingTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("new"))
	}))
	defer server.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(target, []byte("old"), 0o600); err != nil {
		t.Fatalf("write existing target: %v", err)
	}

	_, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
	})
	if err == nil {
		t.Fatal("Download returned nil error")
	}
	if !strings.Contains(err.Error(), "target file already exists") {
		t.Fatalf("error = %q, want existing target context", err)
	}
}

func TestDownloadEmitsStartedProgressCompletedEvents(t *testing.T) {
	body := bytes.Repeat([]byte("a"), 128*1024)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		_, _ = w.Write(body)
	}))
	defer server.Close()

	var events []Event
	dir := t.TempDir()
	result, err := New().DownloadWithEvents(download.Plan{
		URLs:   []string{server.URL + "/file.bin"},
		Output: dir,
	}, func(event Event) {
		events = append(events, event)
	})
	if err != nil {
		t.Fatalf("DownloadWithEvents returned error: %v", err)
	}

	if len(events) < 3 {
		t.Fatalf("got %d events, want at least started, progress, completed", len(events))
	}
	if events[0].Type != EventStarted {
		t.Fatalf("first event = %q, want %q", events[0].Type, EventStarted)
	}
	if events[len(events)-1].Type != EventCompleted {
		t.Fatalf("last event = %q, want %q", events[len(events)-1].Type, EventCompleted)
	}

	var progress Event
	for _, event := range events {
		if event.Type == EventProgress {
			progress = event
		}
		if event.URL != server.URL+"/file.bin" {
			t.Fatalf("event.URL = %q", event.URL)
		}
		if event.TargetPath != "" && event.TargetPath != result.Path {
			t.Fatalf("event.TargetPath = %q, want %q", event.TargetPath, result.Path)
		}
		if event.Timestamp.IsZero() {
			t.Fatalf("event %q has zero timestamp", event.Type)
		}
	}
	if progress.Type != EventProgress {
		t.Fatal("no progress event emitted")
	}
	if progress.DownloadedBytes != int64(len(body)) {
		t.Fatalf("progress.DownloadedBytes = %d, want %d", progress.DownloadedBytes, len(body))
	}
	if progress.TotalBytes != int64(len(body)) {
		t.Fatalf("progress.TotalBytes = %d, want %d", progress.TotalBytes, len(body))
	}
	if progress.Percent != 100 {
		t.Fatalf("progress.Percent = %f, want 100", progress.Percent)
	}
	if progress.SpeedBytesPerSecond <= 0 {
		t.Fatalf("progress.SpeedBytesPerSecond = %f, want > 0", progress.SpeedBytesPerSecond)
	}

	completed := events[len(events)-1]
	if completed.DownloadedBytes != int64(len(body)) {
		t.Fatalf("completed.DownloadedBytes = %d, want %d", completed.DownloadedBytes, len(body))
	}
	if completed.Percent != 100 {
		t.Fatalf("completed.Percent = %f, want 100", completed.Percent)
	}
}

func TestDownloadProgressUnknownTotal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		_, _ = w.Write([]byte(" daryaft"))
	}))
	defer server.Close()

	var progressEvents []Event
	_, err := New().DownloadWithEvents(download.Plan{
		URLs:   []string{server.URL + "/stream.txt"},
		Output: t.TempDir(),
	}, func(event Event) {
		if event.Type == EventProgress {
			progressEvents = append(progressEvents, event)
		}
	})
	if err != nil {
		t.Fatalf("DownloadWithEvents returned error: %v", err)
	}
	if len(progressEvents) == 0 {
		t.Fatal("no progress event emitted")
	}

	lastProgress := progressEvents[len(progressEvents)-1]
	if lastProgress.TotalBytes != 0 {
		t.Fatalf("lastProgress.TotalBytes = %d, want 0", lastProgress.TotalBytes)
	}
	if lastProgress.Percent != 0 {
		t.Fatalf("lastProgress.Percent = %f, want 0", lastProgress.Percent)
	}
	if lastProgress.DownloadedBytes != int64(len("hello daryaft")) {
		t.Fatalf("lastProgress.DownloadedBytes = %d", lastProgress.DownloadedBytes)
	}
}

func TestDownloadEmitsFailedEventForNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusTeapot)
	}))
	defer server.Close()

	var events []Event
	_, err := New().DownloadWithEvents(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: t.TempDir(),
	}, func(event Event) {
		events = append(events, event)
	})
	if err == nil {
		t.Fatal("DownloadWithEvents returned nil error")
	}
	if len(events) != 1 {
		t.Fatalf("got %d events, want one failed event", len(events))
	}
	if events[0].Type != EventFailed {
		t.Fatalf("event.Type = %q, want %q", events[0].Type, EventFailed)
	}
	if events[0].Error == nil {
		t.Fatal("failed event Error is nil")
	}
	if !strings.Contains(events[0].Error.Error(), "server returned 418") {
		t.Fatalf("failed event error = %q", events[0].Error)
	}
}
