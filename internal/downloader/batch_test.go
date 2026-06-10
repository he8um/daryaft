package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/he8um/daryaft/internal/checksum"
	"github.com/he8um/daryaft/internal/download"
)

func TestDownloadBatchDownloadsMultipleURLs(t *testing.T) {
	var (
		mu    sync.Mutex
		paths []string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()

		switch r.URL.Path {
		case "/a.txt":
			_, _ = w.Write([]byte("a"))
		case "/b.txt":
			_, _ = w.Write([]byte("b"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	result := New().DownloadBatch(download.Plan{
		URLs: []string{
			server.URL + "/a.txt",
			server.URL + "/b.txt",
		},
		Output: dir,
	}, BatchHandlers{})

	if result.Err() != nil {
		t.Fatalf("result.Err() = %v", result.Err())
	}
	if result.Total() != 2 || result.Completed() != 2 || result.Failed() != 0 {
		t.Fatalf("counts = total %d completed %d failed %d", result.Total(), result.Completed(), result.Failed())
	}
	if !reflect.DeepEqual(paths, []string{"/a.txt", "/b.txt"}) {
		t.Fatalf("request order = %#v", paths)
	}

	assertFileContent(t, filepath.Join(dir, "a.txt"), "a")
	assertFileContent(t, filepath.Join(dir, "b.txt"), "b")
}

func TestDownloadBatchCancellationStopsRemainingItems(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte("body"))
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	result := New().DownloadBatchContext(ctx, download.Plan{
		URLs: []string{
			server.URL + "/a.txt",
			server.URL + "/b.txt",
		},
		Output: t.TempDir(),
	}, BatchHandlers{
		Event: func(_ BatchItem, event Event) {
			if event.Type == EventStarted {
				cancel()
			}
		},
	})

	if !errors.Is(result.Err(), ErrCancelled) {
		t.Fatalf("result.Err() = %v, want ErrCancelled", result.Err())
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
	if result.Total() != 2 || result.Completed() != 0 || result.Failed() != 0 || result.Cancelled() != 1 || result.Skipped() != 1 {
		t.Fatalf("counts = total %d completed %d failed %d cancelled %d skipped %d", result.Total(), result.Completed(), result.Failed(), result.Cancelled(), result.Skipped())
	}
}

func TestDownloadBatchContinuesAfterFailedItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fail.txt":
			http.Error(w, "nope", http.StatusTeapot)
		case "/ok.txt":
			_, _ = w.Write([]byte("ok"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	result := New().DownloadBatch(download.Plan{
		URLs: []string{
			server.URL + "/fail.txt",
			server.URL + "/ok.txt",
		},
		Output: dir,
	}, BatchHandlers{})

	if result.Err() == nil {
		t.Fatal("result.Err() = nil, want batch failure")
	}
	if result.Total() != 2 || result.Completed() != 1 || result.Failed() != 1 {
		t.Fatalf("counts = total %d completed %d failed %d", result.Total(), result.Completed(), result.Failed())
	}
	assertFileContent(t, filepath.Join(dir, "ok.txt"), "ok")

	failures := result.FailedItems()
	if len(failures) != 1 {
		t.Fatalf("len(failures) = %d, want 1", len(failures))
	}
	if failures[0].Item.URL != server.URL+"/fail.txt" {
		t.Fatalf("failed URL = %q", failures[0].Item.URL)
	}
	if !strings.Contains(failures[0].Err.Error(), "server returned 418") {
		t.Fatalf("failure error = %q", failures[0].Err)
	}
}

func TestDownloadBatchExistingTargetDoesNotStopFollowingDownloads(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/exists.txt":
			_, _ = w.Write([]byte("new"))
		case "/next.txt":
			_, _ = w.Write([]byte("next"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	existingPath := filepath.Join(dir, "exists.txt")
	if err := os.WriteFile(existingPath, []byte("old"), 0o600); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	result := New().DownloadBatch(download.Plan{
		URLs: []string{
			server.URL + "/exists.txt",
			server.URL + "/next.txt",
		},
		Output: dir,
	}, BatchHandlers{})

	if result.Total() != 2 || result.Completed() != 1 || result.Failed() != 1 {
		t.Fatalf("counts = total %d completed %d failed %d", result.Total(), result.Completed(), result.Failed())
	}
	assertFileContent(t, existingPath, "old")
	assertFileContent(t, filepath.Join(dir, "next.txt"), "next")

	failures := result.FailedItems()
	if len(failures) != 1 || !strings.Contains(failures[0].Err.Error(), "target file already exists") {
		t.Fatalf("failures = %#v", failures)
	}
}

func TestBatchResultSummaryString(t *testing.T) {
	result := BatchResult{
		Items: []BatchItemResult{
			{Item: BatchItem{Index: 1, Total: 2, URL: "https://example.com/a.txt"}},
			{Item: BatchItem{Index: 2, Total: 2, URL: "https://example.com/b.txt"}, Err: os.ErrPermission},
		},
	}

	summary := result.SummaryString()
	for _, want := range []string{
		"Daryaft batch summary",
		"Total: 2",
		"Completed: 1",
		"Failed: 1",
		"Failed downloads:",
		"- https://example.com/b.txt: permission denied",
	} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary missing %q in:\n%s", want, summary)
		}
	}
}

func TestBatchSummaryNotStartedLabel(t *testing.T) {
	result := BatchResult{
		Planned: 3,
		Items: []BatchItemResult{
			{Item: BatchItem{Index: 1, Total: 3, URL: "https://example.com/a.txt"}},
			{Item: BatchItem{Index: 2, Total: 3, URL: "https://example.com/b.txt"}, Err: ErrCancelled},
		},
	}

	if result.Skipped() == 0 {
		t.Fatal("expected Skipped() > 0 for this result")
	}

	summary := result.SummaryString()
	if !strings.Contains(summary, "Not started:") {
		t.Errorf("SummaryString() missing 'Not started:' label; got:\n%s", summary)
	}
	if strings.Contains(summary, "Skipped:") {
		t.Errorf("SummaryString() still contains old 'Skipped:' label; got:\n%s", summary)
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	if string(content) != want {
		t.Fatalf("%s content = %q, want %q", path, content, want)
	}
}

func sha256Hex(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func TestDownloadBatchContext_ChecksumVerified(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a.txt":
			_, _ = w.Write([]byte("alpha"))
		case "/b.txt":
			_, _ = w.Write([]byte("beta"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	urlA := server.URL + "/a.txt"
	urlB := server.URL + "/b.txt"
	result := New().DownloadBatchContext(context.Background(), download.Plan{
		URLs:   []string{urlA, urlB},
		Output: dir,
		TargetChecksums: map[string]checksum.Spec{
			urlA: {Algorithm: checksum.AlgorithmSHA256, Expected: sha256Hex("alpha")},
			urlB: {Algorithm: checksum.AlgorithmSHA256, Expected: sha256Hex("beta")},
		},
	}, BatchHandlers{})

	if result.Err() != nil {
		t.Fatalf("result.Err() = %v", result.Err())
	}
	if result.ChecksumVerified != 2 {
		t.Fatalf("ChecksumVerified = %d, want 2", result.ChecksumVerified)
	}
	for _, item := range result.Items {
		if item.ChecksumStatus != "verified" {
			t.Fatalf("item %s ChecksumStatus = %q, want verified", item.Item.URL, item.ChecksumStatus)
		}
	}
}

func TestDownloadBatchContext_ChecksumMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("alpha"))
	}))
	defer server.Close()

	dir := t.TempDir()
	urlA := server.URL + "/a.txt"
	result := New().DownloadBatchContext(context.Background(), download.Plan{
		URLs:   []string{urlA},
		Output: dir,
		TargetChecksums: map[string]checksum.Spec{
			urlA: {Algorithm: checksum.AlgorithmSHA256, Expected: sha256Hex("WRONG")},
		},
	}, BatchHandlers{})

	if result.Failed() != 1 {
		t.Fatalf("Failed() = %d, want 1", result.Failed())
	}
	if result.ChecksumVerified != 0 {
		t.Fatalf("ChecksumVerified = %d, want 0", result.ChecksumVerified)
	}
	if len(result.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(result.Items))
	}
	item := result.Items[0]
	if item.ChecksumStatus != "failed" {
		t.Fatalf("ChecksumStatus = %q, want failed", item.ChecksumStatus)
	}
	if item.Err == nil || !strings.Contains(item.Err.Error(), "checksum mismatch") {
		t.Fatalf("item.Err = %v, want checksum mismatch", item.Err)
	}
	// Downloaded file is left in place for inspection.
	assertFileContent(t, filepath.Join(dir, "a.txt"), "alpha")
}

func TestDownloadBatchContext_NoChecksum(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("alpha"))
	}))
	defer server.Close()

	dir := t.TempDir()
	result := New().DownloadBatchContext(context.Background(), download.Plan{
		URLs:   []string{server.URL + "/a.txt"},
		Output: dir,
	}, BatchHandlers{})

	if result.Err() != nil {
		t.Fatalf("result.Err() = %v", result.Err())
	}
	if result.ChecksumVerified != 0 {
		t.Fatalf("ChecksumVerified = %d, want 0", result.ChecksumVerified)
	}
	if result.Items[0].ChecksumStatus != "" {
		t.Fatalf("ChecksumStatus = %q, want empty", result.Items[0].ChecksumStatus)
	}
}

func TestDownloadBatchContext_ChecksumSkippedOnDownloadFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	dir := t.TempDir()
	urlA := server.URL + "/missing.txt"
	result := New().DownloadBatchContext(context.Background(), download.Plan{
		URLs:    []string{urlA},
		Output:  dir,
		Retries: 0,
		TargetChecksums: map[string]checksum.Spec{
			urlA: {Algorithm: checksum.AlgorithmSHA256, Expected: sha256Hex("alpha")},
		},
	}, BatchHandlers{})

	if result.Failed() != 1 {
		t.Fatalf("Failed() = %d, want 1", result.Failed())
	}
	item := result.Items[0]
	if item.ChecksumStatus != "" {
		t.Fatalf("ChecksumStatus = %q, want empty (download failed before verification)", item.ChecksumStatus)
	}
	if item.Err == nil || strings.Contains(item.Err.Error(), "checksum") {
		t.Fatalf("item.Err = %v, want download error not checksum error", item.Err)
	}
}

func TestDownloadBatchContext_MultipleMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a.txt":
			_, _ = w.Write([]byte("alpha"))
		case "/b.txt":
			_, _ = w.Write([]byte("beta"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	urlA := server.URL + "/a.txt"
	urlB := server.URL + "/b.txt"
	result := New().DownloadBatchContext(context.Background(), download.Plan{
		URLs:   []string{urlA, urlB},
		Output: dir,
		TargetChecksums: map[string]checksum.Spec{
			urlA: {Algorithm: checksum.AlgorithmSHA256, Expected: sha256Hex("WRONG-A")},
			urlB: {Algorithm: checksum.AlgorithmSHA256, Expected: sha256Hex("WRONG-B")},
		},
	}, BatchHandlers{})

	if result.Failed() != 2 {
		t.Fatalf("Failed() = %d, want 2", result.Failed())
	}
	if result.ChecksumVerified != 0 {
		t.Fatalf("ChecksumVerified = %d, want 0", result.ChecksumVerified)
	}
}

func TestDownloadBatchContext_SingleChecksumViaBatchPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("alpha"))
	}))
	defer server.Close()

	dir := t.TempDir()
	urlA := server.URL + "/a.txt"
	spec := checksum.Spec{Algorithm: checksum.AlgorithmSHA256, Expected: sha256Hex("alpha")}
	result := New().DownloadBatchContext(context.Background(), download.Plan{
		URLs:     []string{urlA},
		Output:   dir,
		Checksum: &spec,
	}, BatchHandlers{})

	if result.Err() != nil {
		t.Fatalf("result.Err() = %v", result.Err())
	}
	if result.ChecksumVerified != 1 {
		t.Fatalf("ChecksumVerified = %d, want 1", result.ChecksumVerified)
	}
	if result.Items[0].ChecksumStatus != "verified" {
		t.Fatalf("ChecksumStatus = %q, want verified", result.Items[0].ChecksumStatus)
	}
}

func TestBatchResult_SummaryString_WithChecksumVerified(t *testing.T) {
	result := BatchResult{
		Planned:          2,
		ChecksumVerified: 2,
		Items: []BatchItemResult{
			{Item: BatchItem{Index: 1, Total: 2, URL: "https://example.com/a.zip"}, ChecksumStatus: "verified"},
			{Item: BatchItem{Index: 2, Total: 2, URL: "https://example.com/b.zip"}, ChecksumStatus: "verified"},
		},
	}
	summary := result.SummaryString()
	if !strings.Contains(summary, "Checksum verified: 2") {
		t.Fatalf("summary missing checksum verified line:\n%s", summary)
	}
}

func TestBatchResult_SummaryString_NoChecksumLineWhenZero(t *testing.T) {
	result := BatchResult{
		Planned: 1,
		Items: []BatchItemResult{
			{Item: BatchItem{Index: 1, Total: 1, URL: "https://example.com/a.zip"}},
		},
	}
	summary := result.SummaryString()
	if strings.Contains(summary, "Checksum verified") {
		t.Fatalf("summary unexpectedly contains checksum verified line:\n%s", summary)
	}
}
