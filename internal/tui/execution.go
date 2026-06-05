package tui

import (
	"context"
	"errors"
	"fmt"

	"github.com/he8um/daryaft/internal/checksum"
	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/downloader"

	tea "github.com/charmbracelet/bubbletea"
)

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
}

type executionSummary struct {
	Total     int
	Completed int
	Failed    int
	Cancelled int
	Skipped   int
	Failures  []executionFailure
}

type executionFailure struct {
	URL   string
	Error string
}

type ExecutionRunner func(context.Context, download.Plan, downloader.BatchHandlers) downloader.BatchResult

func defaultExecutionRunner(ctx context.Context, plan download.Plan, handlers downloader.BatchHandlers) downloader.BatchResult {
	result := downloader.New().DownloadBatchContext(ctx, plan, handlers)
	return verifyCompletedChecksum(plan, result)
}

func verifyCompletedChecksum(plan download.Plan, result downloader.BatchResult) downloader.BatchResult {
	if plan.Checksum == nil || len(plan.URLs) != 1 || len(result.Items) != 1 || result.Items[0].Err != nil {
		return result
	}

	if _, err := checksum.VerifyFile(result.Items[0].Result.Path, *plan.Checksum); err != nil {
		result.Items[0].Err = err
	}
	return result
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
}

func (m *Model) applyExecutionFinished(msg executionFinishedMsg) {
	m.execution.Running = false
	m.execution.Done = true
	m.execution.Summary = msg.Summary
	m.executionCancel = nil
	m.executionMessages = nil
	if msg.Summary.Cancelled > 0 {
		m.execution.Status = "Cancelled"
		m.execution.Message = "Download cancelled. Partial file kept for resume."
		return
	}
	if msg.Summary.Failed > 0 {
		m.execution.Status = "Failed"
		return
	}
	m.execution.Status = "Completed"
}

func summaryFromBatch(result downloader.BatchResult) executionSummary {
	summary := executionSummary{
		Total:     result.Total(),
		Completed: result.Completed(),
		Failed:    result.Failed(),
		Cancelled: result.Cancelled(),
		Skipped:   result.Skipped(),
	}
	for _, failure := range result.FailedItems() {
		summary.Failures = append(summary.Failures, executionFailure{
			URL:   failure.Item.URL,
			Error: failure.Err.Error(),
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
