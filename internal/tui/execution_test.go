package tui

import (
	"errors"
	"testing"

	"github.com/he8um/daryaft/internal/downloader"
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
