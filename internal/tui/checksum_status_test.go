package tui

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/checksum"
	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/downloader"
)

func tuiSHA256Hex(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func TestStatusMarker_ChecksumOK_Color(t *testing.T) {
	if got := statusMarker("Checksum OK", false); got != "✓" {
		t.Fatalf("statusMarker(Checksum OK, false) = %q, want ✓", got)
	}
}

func TestStatusMarker_ChecksumOK_NoColor(t *testing.T) {
	if got := statusMarker("Checksum OK", true); got != "[ok]" {
		t.Fatalf("statusMarker(Checksum OK, true) = %q, want [ok]", got)
	}
}

func TestStatusMarker_ChecksumFailed_Color(t *testing.T) {
	if got := statusMarker("Checksum Failed", false); got != "✗" {
		t.Fatalf("statusMarker(Checksum Failed, false) = %q, want ✗", got)
	}
}

func TestStatusMarker_ChecksumFailed_NoColor(t *testing.T) {
	if got := statusMarker("Checksum Failed", true); got != "[!]" {
		t.Fatalf("statusMarker(Checksum Failed, true) = %q, want [!]", got)
	}
}

func TestApplyExecutionFinished_SetsChecksumOK(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model = updateWithMsg(t, model, executionItemStartedMsg{
		Item: downloader.BatchItem{Index: 1, Total: 1, URL: "https://a.com/f1.zip"},
	})
	model = updateWithMsg(t, model, executionEventMsg{
		Item:  downloader.BatchItem{Index: 1, Total: 1},
		Event: downloader.Event{Type: downloader.EventCompleted, DownloadedBytes: 5},
	})

	model.applyExecutionFinished(executionFinishedMsg{
		Summary: executionSummary{
			Total:            1,
			Completed:        1,
			ChecksumVerified: 1,
			ItemChecksumStatuses: []itemChecksumStatus{
				{Index: 1, Status: "verified"},
			},
		},
	})

	if model.execution.Items[0].Status != "Checksum OK" {
		t.Fatalf("item status = %q, want Checksum OK", model.execution.Items[0].Status)
	}
	if model.execution.Items[0].ChecksumStatus != "verified" {
		t.Fatalf("ChecksumStatus = %q, want verified", model.execution.Items[0].ChecksumStatus)
	}
}

func TestApplyExecutionFinished_SetsChecksumFailed(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model = updateWithMsg(t, model, executionItemStartedMsg{
		Item: downloader.BatchItem{Index: 1, Total: 1, URL: "https://a.com/f1.zip"},
	})
	model = updateWithMsg(t, model, executionEventMsg{
		Item:  downloader.BatchItem{Index: 1, Total: 1},
		Event: downloader.Event{Type: downloader.EventCompleted, DownloadedBytes: 5},
	})

	model.applyExecutionFinished(executionFinishedMsg{
		Summary: executionSummary{
			Total:  1,
			Failed: 1,
			ItemChecksumStatuses: []itemChecksumStatus{
				{Index: 1, Status: "failed"},
			},
		},
	})

	if model.execution.Items[0].Status != "Checksum Failed" {
		t.Fatalf("item status = %q, want Checksum Failed", model.execution.Items[0].Status)
	}
	if model.execution.Items[0].ChecksumStatus != "failed" {
		t.Fatalf("ChecksumStatus = %q, want failed", model.execution.Items[0].ChecksumStatus)
	}
}

func TestSummaryView_ShowsChecksumVerifiedCount(t *testing.T) {
	out := summaryView(executionSummary{
		Total:            2,
		Completed:        2,
		ChecksumVerified: 2,
	})
	if !strings.Contains(out, "Checksum verified: 2") {
		t.Fatalf("summary view missing checksum verified count:\n%s", out)
	}
}

func TestSummaryView_NoChecksumLineWhenZero(t *testing.T) {
	out := summaryView(executionSummary{Total: 1, Completed: 1})
	if strings.Contains(out, "Checksum verified") {
		t.Fatalf("summary view unexpectedly contains checksum verified line:\n%s", out)
	}
}

func TestDefaultExecutionRunner_ChecksumVerified(t *testing.T) {
	content := "alpha"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(content))
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	plan := download.Plan{
		URLs:   []string{urlA},
		Output: t.TempDir(),
		TargetChecksums: map[string]checksum.Spec{
			urlA: {Algorithm: checksum.AlgorithmSHA256, Expected: tuiSHA256Hex(content)},
		},
		HasChecksumFile: true,
	}

	result := defaultExecutionRunner(context.Background(), plan, downloader.BatchHandlers{})
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

func TestDefaultExecutionRunner_ChecksumFailed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("alpha"))
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	plan := download.Plan{
		URLs:   []string{urlA},
		Output: t.TempDir(),
		TargetChecksums: map[string]checksum.Spec{
			urlA: {Algorithm: checksum.AlgorithmSHA256, Expected: tuiSHA256Hex("WRONG")},
		},
		HasChecksumFile: true,
	}

	result := defaultExecutionRunner(context.Background(), plan, downloader.BatchHandlers{})
	if result.Failed() != 1 {
		t.Fatalf("Failed() = %d, want 1", result.Failed())
	}
	if result.Items[0].ChecksumStatus != "failed" {
		t.Fatalf("ChecksumStatus = %q, want failed", result.Items[0].ChecksumStatus)
	}

	// Verify the summary correlates and the queue renders Checksum Failed.
	summary := summaryFromBatch(result)
	if len(summary.ItemChecksumStatuses) != 1 || summary.ItemChecksumStatuses[0].Status != "failed" {
		t.Fatalf("summary checksum statuses = %+v", summary.ItemChecksumStatuses)
	}
}
