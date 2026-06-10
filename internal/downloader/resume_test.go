package downloader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
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

func TestResumeAfterCancellation(t *testing.T) {
	body := bytes.Repeat([]byte("r"), 192*1024)
	requests := 0
	var resumeRange string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
			resumeRange = rangeHeader
			var start int
			if _, err := fmt.Sscanf(rangeHeader, "bytes=%d-", &start); err != nil {
				t.Fatalf("parse range %q: %v", rangeHeader, err)
			}
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, len(body)-1, len(body)))
			w.Header().Set("Content-Length", strconv.Itoa(len(body)-start))
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write(body[start:])
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		_, _ = w.Write(body)
	}))
	defer server.Close()

	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	_, err := New().DownloadWithEventsContext(ctx, download.Plan{
		URLs:   []string{server.URL + "/file.bin"},
		Output: dir,
		Resume: true,
	}, func(event Event) {
		if event.Type == EventProgress {
			cancel()
		}
	})
	if !errors.Is(err, ErrCancelled) {
		t.Fatalf("error = %v, want ErrCancelled", err)
	}

	partial := filepath.Join(dir, "file.bin.part")
	info, err := os.Stat(partial)
	if err != nil {
		t.Fatalf("partial missing after cancellation: %v", err)
	}
	if info.Size() <= 0 || info.Size() >= int64(len(body)) {
		t.Fatalf("partial size = %d, want between 0 and %d", info.Size(), len(body))
	}

	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.bin"},
		Output: dir,
		Resume: true,
	})
	if err != nil {
		t.Fatalf("resume Download returned error: %v", err)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
	if resumeRange != fmt.Sprintf("bytes=%d-", info.Size()) {
		t.Fatalf("resume Range = %q, want bytes=%d-", resumeRange, info.Size())
	}
	content, err := os.ReadFile(result.Path)
	if err != nil {
		t.Fatalf("read resumed file: %v", err)
	}
	if !bytes.Equal(content, body) {
		t.Fatalf("resumed content length = %d, want %d", len(content), len(body))
	}
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

func TestRequestedRangeNotSatisfiableRestartsSafely(t *testing.T) {
	requests := 0
	var ranges []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		ranges = append(ranges, r.Header.Get("Range"))
		if requests == 1 {
			w.Header().Set("Content-Range", "bytes */4")
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		_, _ = w.Write([]byte("fresh"))
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

	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
	if strings.Join(ranges, ",") != "bytes=5-," {
		t.Fatalf("ranges = %#v, want first range then full", ranges)
	}
	assertFileContent(t, result.Path, "fresh")
	restarting := findEvent(events, EventRestarting)
	if restarting.Message != resumeNotSupportedMessage {
		t.Fatalf("restart message = %q", restarting.Message)
	}
}

func TestRestartWithNewRequestClosesOldAndActiveResponsesOnce(t *testing.T) {
	oldCloses := 0
	newCloses := 0
	requests := 0
	var ranges []string

	client := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			requests++
			ranges = append(ranges, request.Header.Get("Range"))
			if requests == 1 {
				return &http.Response{
					StatusCode: http.StatusRequestedRangeNotSatisfiable,
					Status:     "416 Requested Range Not Satisfiable",
					Header:     make(http.Header),
					Body:       &countingReadCloser{reader: strings.NewReader("stale range"), closes: &oldCloses},
					Request:    request,
				}, nil
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Header:     http.Header{"Content-Length": []string{"5"}},
				Body:       &countingReadCloser{reader: strings.NewReader("fresh"), closes: &newCloses},
				Request:    request,
			}, nil
		}),
	}

	dir := t.TempDir()
	partial := filepath.Join(dir, "file.txt.part")
	if err := os.WriteFile(partial, []byte("stale"), 0o600); err != nil {
		t.Fatalf("write partial: %v", err)
	}

	result, err := NewWithClient(client).Download(download.Plan{
		URLs:   []string{"https://example.com/file.txt"},
		Output: dir,
		Resume: true,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
	if strings.Join(ranges, ",") != "bytes=5-," {
		t.Fatalf("ranges = %#v, want first range then full", ranges)
	}
	if oldCloses != 1 {
		t.Fatalf("old response closes = %d, want 1", oldCloses)
	}
	if newCloses != 1 {
		t.Fatalf("new response closes = %d, want 1", newCloses)
	}
	assertFileContent(t, result.Path, "fresh")
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

func TestChangedLastModifiedRestartsSafely(t *testing.T) {
	requests := 0
	var ranges []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		ranges = append(ranges, r.Header.Get("Range"))
		if requests == 1 {
			w.Header().Set("Last-Modified", "Mon, 01 Jun 2026 00:00:00 GMT")
			w.Header().Set("Content-Range", "bytes 5-7/8")
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write([]byte("bad"))
			return
		}

		w.Header().Set("Last-Modified", "Mon, 01 Jun 2026 00:00:00 GMT")
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
		LastModified:    "Sun, 31 May 2026 00:00:00 GMT",
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

func TestPartialLargerThanRemoteRestartsSafely(t *testing.T) {
	// Server returns a 206 whose Content-Range total is smaller than the local
	// partial offset on the first request, then a full body on restart. This
	// exercises the proactive partial-overflow guard: appending would write past
	// the real end of the remote file, so Daryaft must restart from byte 0.
	remote := "fresh"
	requests := 0
	var ranges []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		ranges = append(ranges, r.Header.Get("Range"))
		if requests == 1 {
			// Misbehaving partial response: total (5) <= requested offset (10).
			w.Header().Set("Content-Range", "bytes 10-10/5")
			w.WriteHeader(http.StatusPartialContent)
			return
		}
		_, _ = w.Write([]byte(remote))
	}))
	defer server.Close()

	dir := t.TempDir()
	partial := filepath.Join(dir, "file.txt.part")
	// Local partial (10 bytes) is larger than the remote file (5 bytes).
	if err := os.WriteFile(partial, []byte("0123456789"), 0o600); err != nil {
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

	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}
	if strings.Join(ranges, ",") != "bytes=10-," {
		t.Fatalf("ranges = %#v, want first range then full", ranges)
	}
	assertFileContent(t, result.Path, remote)
	restarting := findEvent(events, EventRestarting)
	if restarting.Message != partialLargerThanRemoteMessage {
		t.Fatalf("restart message = %q, want %q", restarting.Message, partialLargerThanRemoteMessage)
	}
}

func TestMissingSidecarMetadataRestartsSafely(t *testing.T) {
	// A .part file exists with no sidecar metadata. Daryaft should still produce
	// the correct final file without panicking or unsafely appending.
	body := "hello resumed world"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
			var start int
			if _, err := fmt.Sscanf(rangeHeader, "bytes=%d-", &start); err != nil {
				t.Fatalf("parse range %q: %v", rangeHeader, err)
			}
			if start < 0 || start > len(body) {
				w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", len(body)))
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
				return
			}
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, len(body)-1, len(body)))
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write([]byte(body[start:]))
			return
		}
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	dir := t.TempDir()
	partial := filepath.Join(dir, "file.txt.part")
	if err := os.WriteFile(partial, []byte(body[:6]), 0o600); err != nil {
		t.Fatalf("write partial: %v", err)
	}
	// Deliberately no sidecar metadata file written.
	assertMissing(t, metadataPathForPartial(partial))

	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
		Resume: true,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	assertFileContent(t, result.Path, body)
}

func TestCorruptSidecarMetadataRestartsSafely(t *testing.T) {
	// A .part file exists alongside a corrupt sidecar metadata file. Daryaft must
	// not panic and must produce the correct final file.
	body := "hello resumed world"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
			var start int
			if _, err := fmt.Sscanf(rangeHeader, "bytes=%d-", &start); err != nil {
				t.Fatalf("parse range %q: %v", rangeHeader, err)
			}
			if start < 0 || start > len(body) {
				w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", len(body)))
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
				return
			}
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, len(body)-1, len(body)))
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write([]byte(body[start:]))
			return
		}
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	dir := t.TempDir()
	partial := filepath.Join(dir, "file.txt.part")
	if err := os.WriteFile(partial, []byte(body[:6]), 0o600); err != nil {
		t.Fatalf("write partial: %v", err)
	}
	if err := os.WriteFile(metadataPathForPartial(partial), []byte("{ not valid json"), 0o600); err != nil {
		t.Fatalf("write corrupt metadata: %v", err)
	}

	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
		Resume: true,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	assertFileContent(t, result.Path, body)
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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

type countingReadCloser struct {
	reader *strings.Reader
	closes *int
}

func (body *countingReadCloser) Read(p []byte) (int, error) {
	return body.reader.Read(p)
}

func (body *countingReadCloser) Close() error {
	(*body.closes)++
	return nil
}
