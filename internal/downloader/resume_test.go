package downloader

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/he8um/daryaft/internal/download"
)

func TestFreshDownloadRemovesPartialAndMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "11")
		_, _ = w.Write([]byte("hello fresh"))
	}))
	defer server.Close()

	dir := t.TempDir()
	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
		Resume: true,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	assertFileContent(t, result.Path, "hello fresh")
	assertMissing(t, result.Path+".part")
	assertMissing(t, result.Path+".part.daryaft.json")
}

func TestPartialFileCanResumeWithRangeAndComplete(t *testing.T) {
	body := "hello resumed world"
	var gotRange string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRange = r.Header.Get("Range")
		if gotRange != "bytes=6-" {
			t.Fatalf("Range = %q, want bytes=6-", gotRange)
		}
		w.Header().Set("Content-Range", "bytes 6-18/19")
		w.Header().Set("Content-Length", "13")
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write([]byte(body[6:]))
	}))
	defer server.Close()

	dir := t.TempDir()
	partial := filepath.Join(dir, "file.txt.part")
	if err := os.WriteFile(partial, []byte(body[:6]), 0o600); err != nil {
		t.Fatalf("write partial: %v", err)
	}

	var events []Event
	result, err := New().DownloadWithEvents(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
		Resume: true,
	}, func(event Event) {
		events = append(events, event)
	})
	if err != nil {
		t.Fatalf("DownloadWithEvents returned error: %v", err)
	}

	assertFileContent(t, result.Path, body)
	if gotRange != "bytes=6-" {
		t.Fatalf("Range = %q, want bytes=6-", gotRange)
	}
	resuming := findEvent(events, EventResuming)
	if resuming.Type != EventResuming {
		t.Fatal("missing resuming event")
	}
	if resuming.DownloadedBytes != 6 {
		t.Fatalf("resuming.DownloadedBytes = %d, want 6", resuming.DownloadedBytes)
	}
}

func TestServerWithoutRangeSupportRestartsSafely(t *testing.T) {
	var gotRange string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRange = r.Header.Get("Range")
		_, _ = w.Write([]byte("fresh body"))
	}))
	defer server.Close()

	dir := t.TempDir()
	partial := filepath.Join(dir, "file.txt.part")
	if err := os.WriteFile(partial, []byte("stale"), 0o600); err != nil {
		t.Fatalf("write partial: %v", err)
	}

	var events []Event
	result, err := New().DownloadWithEvents(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
		Resume: true,
	}, func(event Event) {
		events = append(events, event)
	})
	if err != nil {
		t.Fatalf("DownloadWithEvents returned error: %v", err)
	}

	if gotRange != "bytes=5-" {
		t.Fatalf("Range = %q, want bytes=5-", gotRange)
	}
	assertFileContent(t, result.Path, "fresh body")
	restarting := findEvent(events, EventRestarting)
	if restarting.Message != resumeNotSupportedMessage {
		t.Fatalf("restart message = %q", restarting.Message)
	}
}

func TestNoResumeRestartsFromByteZero(t *testing.T) {
	var gotRange string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRange = r.Header.Get("Range")
		_, _ = w.Write([]byte("fresh"))
	}))
	defer server.Close()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "file.txt.part"), []byte("stale"), 0o600); err != nil {
		t.Fatalf("write partial: %v", err)
	}

	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
		Resume: false,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	if gotRange != "" {
		t.Fatalf("Range = %q, want empty", gotRange)
	}
	assertFileContent(t, result.Path, "fresh")
}

func TestExistingFinalTargetFailsBeforeTouchingPartial(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte("new"))
	}))
	defer server.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "file.txt")
	partial := target + ".part"
	if err := os.WriteFile(target, []byte("old"), 0o600); err != nil {
		t.Fatalf("write target: %v", err)
	}
	if err := os.WriteFile(partial, []byte("partial"), 0o600); err != nil {
		t.Fatalf("write partial: %v", err)
	}

	_, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
		Resume: true,
	})
	if err == nil {
		t.Fatal("Download returned nil error")
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
	assertFileContent(t, target, "old")
	assertFileContent(t, partial, "partial")
}

func TestPartialMetadataLoadSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt.part.daryaft.json")
	createdAt := time.Date(2026, 5, 24, 1, 2, 3, 0, time.UTC)
	want := partialMetadata{
		URL:             "https://example.com/file.txt",
		TargetPath:      "file.txt",
		PartialPath:     "file.txt.part",
		TotalBytes:      10,
		DownloadedBytes: 4,
		ETag:            `"abc"`,
		LastModified:    "Sun, 24 May 2026 00:00:00 GMT",
		AcceptRanges:    "bytes",
		CreatedAt:       createdAt,
	}

	if err := savePartialMetadata(path, want); err != nil {
		t.Fatalf("savePartialMetadata returned error: %v", err)
	}
	got, err := loadPartialMetadata(path)
	if err != nil {
		t.Fatalf("loadPartialMetadata returned error: %v", err)
	}

	if got.URL != want.URL || got.TargetPath != want.TargetPath || got.PartialPath != want.PartialPath {
		t.Fatalf("metadata identity = %#v, want %#v", got, want)
	}
	if got.TotalBytes != want.TotalBytes || got.DownloadedBytes != want.DownloadedBytes {
		t.Fatalf("metadata bytes = %d/%d, want %d/%d", got.DownloadedBytes, got.TotalBytes, want.DownloadedBytes, want.TotalBytes)
	}
	if got.ETag != want.ETag || got.LastModified != want.LastModified || got.AcceptRanges != want.AcceptRanges {
		t.Fatalf("metadata headers = %#v", got)
	}
	if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
		t.Fatalf("metadata timestamps not set: %#v", got)
	}
}

func TestRetryAfterPartialFailureResumesWhenEnabled(t *testing.T) {
	body := "hello retry resume"
	requests := 0
	var secondRange string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			w.Header().Set("Content-Length", "18")
			_, _ = w.Write([]byte(body[:6]))
			return
		}

		secondRange = r.Header.Get("Range")
		if secondRange != "bytes=6-" {
			t.Fatalf("second Range = %q, want bytes=6-", secondRange)
		}
		w.Header().Set("Content-Range", "bytes 6-17/18")
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write([]byte(body[6:]))
	}))
	defer server.Close()

	d := newTestDownloader(server)
	result, err := d.Download(download.Plan{
		URLs:    []string{server.URL + "/file.txt"},
		Output:  t.TempDir(),
		Resume:  true,
		Retries: 1,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
	if secondRange != "bytes=6-" {
		t.Fatalf("second Range = %q, want bytes=6-", secondRange)
	}
	assertFileContent(t, result.Path, body)
}

func TestChangedETagRestartsSafely(t *testing.T) {
	requests := 0
	var ranges []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		ranges = append(ranges, r.Header.Get("Range"))
		if requests == 1 {
			w.Header().Set("ETag", `"v2"`)
			w.Header().Set("Content-Range", "bytes 5-7/8")
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write([]byte("bad"))
			return
		}

		w.Header().Set("ETag", `"v2"`)
		_, _ = w.Write([]byte("new body"))
	}))
	defer server.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "file.txt")
	partial := target + ".part"
	if err := os.WriteFile(partial, []byte("hello"), 0o600); err != nil {
		t.Fatalf("write partial: %v", err)
	}
	if err := savePartialMetadata(metadataPathForPartial(partial), partialMetadata{
		URL:             server.URL + "/file.txt",
		TargetPath:      target,
		PartialPath:     partial,
		DownloadedBytes: 5,
		TotalBytes:      8,
		ETag:            `"v1"`,
	}); err != nil {
		t.Fatalf("save metadata: %v", err)
	}

	var events []Event
	result, err := New().DownloadWithEvents(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
		Resume: true,
	}, func(event Event) {
		events = append(events, event)
	})
	if err != nil {
		t.Fatalf("DownloadWithEvents returned error: %v", err)
	}

	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
	if strings.Join(ranges, ",") != "bytes=5-," {
		t.Fatalf("ranges = %#v, want first range then full", ranges)
	}
	assertFileContent(t, result.Path, "new body")
	restarting := findEvent(events, EventRestarting)
	if restarting.Message != remoteChangedMessage {
		t.Fatalf("restart message = %q", restarting.Message)
	}
}

func findEvent(events []Event, eventType EventType) Event {
	for _, event := range events {
		if event.Type == eventType {
			return event
		}
	}
	return Event{}
}

func assertMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("%s exists or stat failed: %v", path, err)
	}
}
