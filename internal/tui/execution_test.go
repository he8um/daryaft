package tui

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/downloader"
	"github.com/he8um/daryaft/internal/httpopts"
)

func TestExecutionStateAccumulatesItems(t *testing.T) {
	model := NewModel(Options{NoColor: true})

	model = updateWithMsg(t, model, executionItemStartedMsg{
		Item: downloader.BatchItem{Index: 1, URL: "https://a.com/f1.zip"},
	})
	model = updateWithMsg(t, model, executionEventMsg{
		Item: downloader.BatchItem{Index: 1},
		Event: downloader.Event{
			Type:            downloader.EventCompleted,
			DownloadedBytes: 100,
		},
	})
	model = updateWithMsg(t, model, executionItemStartedMsg{
		Item: downloader.BatchItem{Index: 2, URL: "https://a.com/f2.zip"},
	})
	model = updateWithMsg(t, model, executionEventMsg{
		Item: downloader.BatchItem{Index: 2},
		Event: downloader.Event{
			Type:  downloader.EventFailed,
			Error: errors.New("timeout"),
		},
	})
	model = updateWithMsg(t, model, executionItemStartedMsg{
		Item: downloader.BatchItem{Index: 3, URL: "https://a.com/f3.zip"},
	})

	exec := model.execution
	if len(exec.Items) != 3 {
		t.Fatalf("want 3 items, got %d", len(exec.Items))
	}
	if exec.Items[0].Status != "Completed" {
		t.Errorf("item 0 want Completed, got %s", exec.Items[0].Status)
	}
	if exec.Items[1].Status != "Failed" {
		t.Errorf("item 1 want Failed, got %s", exec.Items[1].Status)
	}
	if exec.Items[1].Err != "timeout" {
		t.Errorf("item 1 Err want timeout, got %s", exec.Items[1].Err)
	}
	if exec.Items[2].Status != "Downloading" {
		t.Errorf("item 2 want Downloading, got %s", exec.Items[2].Status)
	}
}

func TestStatusMarkerNoColor(t *testing.T) {
	cases := []struct{ status, want string }{
		{"Completed", "[ok]"},
		{"Failed", "[!]"},
		{"Cancelled", "[!]"},
		{"Downloading", "[>]"},
		{"Starting", "[>]"},
		{"Unknown", "[-]"},
	}
	for _, c := range cases {
		got := statusMarker(c.status, true)
		if got != c.want {
			t.Errorf("statusMarker(%q, true) = %q, want %q", c.status, got, c.want)
		}
	}
}

func TestTruncateURL(t *testing.T) {
	cases := []struct {
		url, want string
		max       int
	}{
		{"https://example.com/short.zip", "https://example.com/short.zip", 55},
		{"https://example.com/a/very/long/path/to/a/file/with/a/really/long/name.zip", "https://example.com/a/very/long/path/to/a/file/with/...", 55},
		{"abc", "abc", 3},
		{"abcd", "abc", 3},
		{"abc", "", 0},
	}
	for _, c := range cases {
		got := truncateURL(c.url, c.max)
		if got != c.want {
			t.Errorf("truncateURL(%q, %d) = %q, want %q", c.url, c.max, got, c.want)
		}
	}
}

func TestNewDefaultExecutionRunnerAppliesUserAgentAndTimeout(t *testing.T) {
	var receivedUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("content"))
	}))
	defer server.Close()

	runner := newDefaultExecutionRunner("10s")
	plan := download.Plan{
		URLs:        []string{server.URL},
		Output:      t.TempDir(),
		HTTPOptions: httpopts.Options{UserAgent: "DaryaftTUI/1.13"},
	}

	result := runner(context.Background(), plan, downloader.BatchHandlers{})
	if result.Completed() != 1 {
		t.Fatalf("expected 1 completed, got %d (failed: %d)", result.Completed(), result.Failed())
	}
	if receivedUA != "DaryaftTUI/1.13" {
		t.Fatalf("received User-Agent = %q, want DaryaftTUI/1.13", receivedUA)
	}
}

func TestNewDefaultExecutionRunnerInvalidTimeout(t *testing.T) {
	runner := newDefaultExecutionRunner("not-a-duration")
	plan := download.Plan{
		URLs:   []string{"http://127.0.0.1:9999/test"},
		Output: t.TempDir(),
	}
	result := runner(context.Background(), plan, downloader.BatchHandlers{})
	if result.Failed() != 1 {
		t.Fatalf("expected 1 failure for invalid timeout, got %d", result.Failed())
	}
	if len(result.Items) == 0 || !strings.Contains(result.Items[0].Err.Error(), "invalid timeout") {
		t.Fatalf("expected invalid timeout error, got %v", result.Items)
	}
}

func TestNewDefaultExecutionRunnerNonPositiveTimeout(t *testing.T) {
	for _, timeout := range []string{"0s", "-5s"} {
		runner := newDefaultExecutionRunner(timeout)
		plan := download.Plan{
			URLs:   []string{"http://127.0.0.1:9999/test"},
			Output: t.TempDir(),
		}
		result := runner(context.Background(), plan, downloader.BatchHandlers{})
		if result.Failed() != 1 {
			t.Fatalf("expected 1 failure for non-positive timeout %q, got %d", timeout, result.Failed())
		}
		if len(result.Items) == 0 || !strings.Contains(result.Items[0].Err.Error(), "timeout must be positive") {
			t.Fatalf("expected timeout must be positive error, got %v", result.Items)
		}
	}
}

func TestDefaultExecutionRunnerUsesEmptyTimeout(t *testing.T) {
	var receivedUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	plan := download.Plan{
		URLs:        []string{server.URL},
		Output:      t.TempDir(),
		HTTPOptions: httpopts.Options{UserAgent: "CustomAgent/1"},
	}

	result := defaultExecutionRunner(context.Background(), plan, downloader.BatchHandlers{})
	if result.Completed() != 1 {
		t.Fatalf("expected 1 completed, got %d", result.Completed())
	}
	if receivedUA != "CustomAgent/1" {
		t.Fatalf("received User-Agent = %q, want CustomAgent/1", receivedUA)
	}
}
