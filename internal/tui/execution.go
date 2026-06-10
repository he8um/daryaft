package tui

import (
	"context"
	"errors"
	"fmt"

	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/downloader"

	tea "github.com/charmbracelet/bubbletea"
)

type itemRecord struct {
	Index          int
	URL            string
	Status         string
	Err            string
	ChecksumStatus string // "", "verified", "failed"
}

type executionState struct {
	Running         bool
	Done            bool
	ItemIndex       int
	ItemTotal       int
	CurrentURL      string
	TargetPath      string
	Status          string
	DownloadedBytes int64
	TotalBytes      int64
	Percent         float64
	Speed           float64
	Message         string
	Summary         executionSummary
	Items           []itemRecord
}

type executionSummary struct {
	Total                int
	Completed            int
	Failed               int
	Cancelled            int
	Skipped              int
	ChecksumVerified     int
	Failures             []executionFailure
	ItemChecksumStatuses []itemChecksumStatus
}

type itemChecksumStatus struct {
	Index  int
	Status string // "", "verified", "failed"
}

type executionFailure struct {
	URL   string
	Error string
}

type ExecutionRunner func(context.Context, download.Plan, downloader.BatchHandlers) downloader.BatchResult

func defaultExecutionRunner(ctx context.Context, plan download.Plan, handlers downloader.BatchHandlers) downloader.BatchResult {
	// Checksum verification (both single-target Plan.Checksum and per-target
	// Plan.TargetChecksums) is handled inside DownloadBatchContext, so the TUI
	// does not perform any checksum verification itself.
	return downloader.New().DownloadBatchContext(ctx, plan, handlers)
}

func newExecutionState(plan download.Plan) executionState {
	state := executionState{
		Running:   true,
		ItemTotal: len(plan.URLs),
		Status:    "Starting",
	}
	if len(plan.URLs) > 0 {
		state.ItemIndex = 1
		state.CurrentURL = plan.URLs[0]
	}
	return state
}

func runExecution(ctx context.Context, plan download.Plan, runner ExecutionRunner) <-chan tea.Msg {
	messages := make(chan tea.Msg, 64)
	if runner == nil {
		runner = defaultExecutionRunner
	}

	go func() {
		defer close(messages)

		result := runner(ctx, plan, downloader.BatchHandlers{
			ItemStarted: func(item downloader.BatchItem) {
				messages <- executionItemStartedMsg{Item: item}
			},
			Event: func(item downloader.BatchItem, event downloader.Event) {
				messages <- executionEventMsg{Item: item, Event: event}
			},
		})

		messages <- executionFinishedMsg{Summary: summaryFromBatch(result)}
	}()

	return messages
}

func waitForExecution(messages <-chan tea.Msg) tea.Cmd {
	if messages == nil {
		return nil
	}
	return func() tea.Msg {
		msg, ok := <-messages
		if !ok {
			return executionClosedMsg{}
		}
		return msg
	}
}

func (m *Model) applyExecutionItemStarted(msg executionItemStartedMsg) {
	m.execution.ItemIndex = msg.Item.Index
	m.execution.ItemTotal = msg.Item.Total
	m.execution.CurrentURL = msg.Item.URL
	m.execution.TargetPath = ""
	m.execution.Status = "Starting"
	m.execution.DownloadedBytes = 0
	m.execution.TotalBytes = 0
	m.execution.Percent = 0
	m.execution.Speed = 0
	m.execution.Message = ""
	m.execution.Items = append(m.execution.Items, itemRecord{
		Index:  msg.Item.Index,
		URL:    msg.Item.URL,
		Status: "Downloading",
	})
}

func (m *Model) applyExecutionEvent(msg executionEventMsg) {
	event := msg.Event
	if msg.Item.Index > 0 {
		m.execution.ItemIndex = msg.Item.Index
		m.execution.ItemTotal = msg.Item.Total
	}
	if event.URL != "" {
		m.execution.CurrentURL = event.URL
	}
	if event.TargetPath != "" {
		m.execution.TargetPath = event.TargetPath
	}

	m.execution.DownloadedBytes = event.DownloadedBytes
	m.execution.TotalBytes = event.TotalBytes
	m.execution.Percent = event.Percent
	m.execution.Speed = event.SpeedBytesPerSecond
	m.execution.Status = statusFromEvent(event.Type)
	m.execution.Message = messageFromEvent(event)
	if len(m.execution.Items) > 0 {
		last := &m.execution.Items[len(m.execution.Items)-1]
		switch event.Type {
		case downloader.EventCompleted:
			last.Status = "Completed"
		case downloader.EventFailed:
			last.Status = "Failed"
			if event.Error != nil {
				last.Err = event.Error.Error()
			} else if event.Message != "" {
				last.Err = event.Message
			}
		case downloader.EventCancelled:
			last.Status = "Cancelled"
		case downloader.EventResuming:
			last.Status = "Resuming"
		case downloader.EventRestarting:
			last.Status = "Restarting"
		}
	}
}

func (m *Model) applyExecutionFinished(msg executionFinishedMsg) {
	m.execution.Running = false
	m.execution.Done = true
	m.execution.Summary = msg.Summary
	m.executionCancel = nil
	m.executionMessages = nil
	m.applyChecksumStatuses(msg.Summary.ItemChecksumStatuses)
	if msg.Summary.Cancelled > 0 {
		m.execution.Status = "Cancelled"
		m.execution.Message = "Download cancelled. Partial file kept for resume."
		if len(m.execution.Items) > 0 {
			last := &m.execution.Items[len(m.execution.Items)-1]
			if last.Status == "" || last.Status == "Downloading" || last.Status == "Starting" || last.Status == "Resuming" || last.Status == "Restarting" {
				last.Status = m.execution.Status
			}
		}
		return
	}
	if msg.Summary.Failed > 0 {
		m.execution.Status = "Failed"
		if len(m.execution.Items) > 0 {
			last := &m.execution.Items[len(m.execution.Items)-1]
			if last.Status == "" || last.Status == "Downloading" || last.Status == "Starting" || last.Status == "Resuming" || last.Status == "Restarting" {
				last.Status = m.execution.Status
			}
		}
		return
	}
	m.execution.Status = "Completed"
	if len(m.execution.Items) > 0 {
		last := &m.execution.Items[len(m.execution.Items)-1]
		if last.Status == "" || last.Status == "Downloading" || last.Status == "Starting" || last.Status == "Resuming" || last.Status == "Restarting" {
			last.Status = m.execution.Status
		}
	}
}

// applyChecksumStatuses updates the queue item records with the final checksum
// verification result reported by the downloader. Items are correlated by their
// 1-based batch index. A checksum result only overrides a successfully completed
// download status; it never overrides a cancelled or download-failed item.
func (m *Model) applyChecksumStatuses(statuses []itemChecksumStatus) {
	for _, status := range statuses {
		if status.Status != "verified" && status.Status != "failed" {
			continue
		}
		for i := range m.execution.Items {
			item := &m.execution.Items[i]
			if item.Index != status.Index {
				continue
			}
			item.ChecksumStatus = status.Status
			switch item.Status {
			case "", "Downloading", "Starting", "Resuming", "Restarting", "Completed":
				if status.Status == "verified" {
					item.Status = "Checksum OK"
				} else {
					item.Status = "Checksum Failed"
				}
			}
		}
	}
}

func summaryFromBatch(result downloader.BatchResult) executionSummary {
	summary := executionSummary{
		Total:            result.Total(),
		Completed:        result.Completed(),
		Failed:           result.Failed(),
		Cancelled:        result.Cancelled(),
		Skipped:          result.Skipped(),
		ChecksumVerified: result.ChecksumVerified,
	}
	for _, failure := range result.FailedItems() {
		summary.Failures = append(summary.Failures, executionFailure{
			URL:   failure.Item.URL,
			Error: failure.Err.Error(),
		})
	}
	for _, item := range result.Items {
		if item.ChecksumStatus == "" {
			continue
		}
		summary.ItemChecksumStatuses = append(summary.ItemChecksumStatuses, itemChecksumStatus{
			Index:  item.Item.Index,
			Status: item.ChecksumStatus,
		})
	}
	return summary
}

func statusFromEvent(eventType downloader.EventType) string {
	switch eventType {
	case downloader.EventStarted, downloader.EventProgress:
		return "Downloading"
	case downloader.EventRetrying:
		return "Retrying"
	case downloader.EventResuming:
		return "Resuming"
	case downloader.EventRestarting:
		return "Restarting"
	case downloader.EventCompleted:
		return "Completed"
	case downloader.EventFailed:
		return "Failed"
	case downloader.EventCancelled:
		return "Cancelled"
	case downloader.EventWarning:
		return "Downloading"
	default:
		return "Starting"
	}
}

func messageFromEvent(event downloader.Event) string {
	if event.Message != "" {
		return event.Message
	}
	if event.Error != nil {
		switch event.Type {
		case downloader.EventRetrying:
			return fmt.Sprintf("Retrying %d/%d in %s: %v", event.Attempt, event.MaxAttempts, event.NextDelay, event.Error)
		case downloader.EventCancelled:
			if errors.Is(event.Error, downloader.ErrCancelled) {
				return "Download cancelled. Partial file kept for resume."
			}
			return event.Error.Error()
		default:
			return event.Error.Error()
		}
	}
	return ""
}
