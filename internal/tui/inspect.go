package tui

import (
	"context"
	"errors"
	"fmt"

	"github.com/he8um/daryaft/internal/inspect"

	tea "github.com/charmbracelet/bubbletea"
)

type inspectState struct {
	Running bool
	Done    bool
	URL     string
	Status  string
	Message string
	Result  inspect.Result
	Error   string
}

type InspectRunner func(context.Context, string) (inspect.Result, error)

func defaultInspectRunner(ctx context.Context, rawURL string) (inspect.Result, error) {
	return inspect.URL(ctx, inspect.Options{URL: rawURL})
}

func newInspectState(rawURL string) inspectState {
	return inspectState{
		Running: true,
		URL:     rawURL,
		Status:  "Inspecting",
		Message: "Inspecting URL metadata...",
	}
}

func runInspect(ctx context.Context, rawURL string, runner InspectRunner) <-chan tea.Msg {
	messages := make(chan tea.Msg, 1)
	if runner == nil {
		runner = defaultInspectRunner
	}

	go func() {
		defer close(messages)
		result, err := runner(ctx, rawURL)
		messages <- inspectFinishedMsg{Result: result, Err: err}
	}()

	return messages
}

func waitForInspect(messages <-chan tea.Msg) tea.Cmd {
	if messages == nil {
		return nil
	}
	return func() tea.Msg {
		msg, ok := <-messages
		if !ok {
			return inspectClosedMsg{}
		}
		return msg
	}
}

func (m *Model) applyInspectFinished(msg inspectFinishedMsg) {
	m.inspect.Running = false
	m.inspect.Done = true
	m.inspectCancel = nil
	m.inspectMessages = nil
	if msg.Err != nil {
		m.screen = screenInspectError
		m.inspect.Status = "Failed"
		m.inspect.Error = inspectErrorMessage(msg.Err)
		m.inspect.Message = m.inspect.Error
		return
	}
	m.screen = screenInspectResult
	m.inspect.Status = "Completed"
	m.inspect.Message = ""
	m.inspect.Result = msg.Result
}

func inspectErrorMessage(err error) string {
	if errors.Is(err, context.Canceled) {
		return "Inspection cancelled."
	}
	return fmt.Sprintf("Inspect failed: %v", err)
}
