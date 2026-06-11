package tui

import (
	"context"

	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/inspect"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.withWindowSize(msg.Width, msg.Height), nil
	case executionItemStartedMsg:
		m.applyExecutionItemStarted(msg)
		return m, waitForExecution(m.executionMessages)
	case executionEventMsg:
		m.applyExecutionEvent(msg)
		return m, waitForExecution(m.executionMessages)
	case executionFinishedMsg:
		m.applyExecutionFinished(msg)
		return m, nil
	case executionClosedMsg:
		return m, nil
	case inspectFinishedMsg:
		m.applyInspectFinished(msg)
		return m, nil
	case inspectClosedMsg:
		return m, nil
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		if m.isInputScreen() {
			return m.updateTextInput(msg)
		}
		return m, nil
	}

	if key.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if key.String() == "q" {
		if m.execution.Running {
			if m.executionCancel != nil {
				m.executionCancel()
			}
			m.execution.Status = "Cancelling"
			m.execution.Message = "Cancelling..."
			return m, nil
		}
		if m.inspect.Running {
			if m.inspectCancel != nil {
				m.inspectCancel()
			}
			m.inspect.Status = "Cancelling"
			m.inspect.Message = "Cancelling..."
			return m, nil
		}
		return m, tea.Quit
	}

	switch m.screen {
	case screenPlan:
		if key.String() == "enter" {
			return m.startExecution()
		}
		if isBackKey(key) {
			return m.back()
		}
		if key.String() == "h" {
			return m.home(), nil
		}
		return m, nil
	case screenURLInput, screenFileInput:
		if key.String() == "esc" || key.String() == "backspace" && m.input.Value() == "" {
			return m.back()
		}
		if key.String() == "enter" {
			return m.submitSourceInput()
		}
		return m.updateTextInput(msg)
	case screenInspectInput:
		if key.String() == "esc" || key.String() == "backspace" && m.input.Value() == "" {
			return m.back()
		}
		if key.String() == "enter" {
			return m.submitInspectInput()
		}
		return m.updateTextInput(msg)
	case screenOutputInput:
		if key.String() == "esc" || key.String() == "backspace" && m.input.Value() == "" {
			return m.back()
		}
		if key.String() == "enter" {
			return m.submitOutputInput()
		}
		return m.updateTextInput(msg)
	case screenFilenameInput:
		if key.String() == "esc" || key.String() == "backspace" && m.input.Value() == "" {
			return m.back()
		}
		if key.String() == "enter" {
			return m.submitFilenameInput()
		}
		return m.updateTextInput(msg)
	case screenChecksumInput:
		if key.String() == "esc" || key.String() == "backspace" && m.input.Value() == "" {
			return m.back()
		}
		if key.String() == "enter" {
			return m.submitChecksumInput()
		}
		return m.updateTextInput(msg)
	case screenHelp, screenVersion, screenSettings:
		if isBackKey(key) {
			return m.back()
		}
		return m, nil
	case screenExecution:
		if !m.execution.Running && (key.String() == "enter" || key.String() == "h") {
			return m.home(), nil
		}
		return m, nil
	case screenInspectExecution:
		return m, nil
	case screenInspectResult:
		if key.String() == "enter" || key.String() == "h" {
			return m.home(), nil
		}
		if isBackKey(key) {
			return m.back()
		}
		return m, nil
	case screenInspectError:
		if key.String() == "h" {
			return m.home(), nil
		}
		if isBackKey(key) {
			return m.back()
		}
		return m, nil
	}

	switch {
	case isUpKey(key):
		return m.moveUp(), nil
	case isDownKey(key):
		return m.moveDown(), nil
	case key.String() == "enter":
		return m.enter()
	case key.String() == "c" && m.screen == screenHome:
		m.screen = screenSettings
		return m, nil
	default:
		return m, nil
	}
}

func (m Model) submitInspectInput() (Model, tea.Cmd) {
	rawURL := m.input.Value()
	if err := inspect.ValidateURL(rawURL); err != nil {
		m.errorMessage = err.Error()
		return m, nil
	}

	m.inspectInput = rawURL
	ctx, cancel := context.WithCancel(context.Background())
	m.screen = screenInspectExecution
	m.errorMessage = ""
	m.inspect = newInspectState(rawURL)
	m.inspectCancel = cancel
	m.inspectMessages = runInspect(ctx, rawURL, m.inspectRunner)
	m.input.Blur()
	return m, waitForInspect(m.inspectMessages)
}

func (m Model) updateTextInput(msg tea.Msg) (Model, tea.Cmd) {
	before := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.input.Value() != before {
		m.errorMessage = ""
	}
	return m, cmd
}

func (m Model) submitSourceInput() (Model, tea.Cmd) {
	var (
		plan download.Plan
		err  error
	)

	switch m.screen {
	case screenURLInput:
		if err = tuiURLError(m.input.Value()); err != nil {
			m.errorMessage = err.Error()
			return m, nil
		}
		plan, err = planFromURL(m.input.Value(), "", "", "", m.retries, m.resume)
	case screenFileInput:
		if err = tuiFilePathError(m.input.Value()); err != nil {
			m.errorMessage = err.Error()
			return m, nil
		}
		plan, err = planFromFile(m.input.Value(), "", m.retries, m.resume)
	default:
		return m, nil
	}

	if err != nil {
		m.errorMessage = err.Error()
		return m, nil
	}

	m.plan = plan
	m.sourceInput = m.input.Value()
	m.sourceScreen = m.screen
	m.outputDirInput = m.defaultOutputDir
	m.filenameInput = ""
	m.checksumInput = ""
	m.screen = screenOutputInput
	m.errorMessage = ""
	m.input = m.newOutputInput(m.outputDirInput)
	return m, m.input.Focus()
}

func (m Model) submitOutputInput() (Model, tea.Cmd) {
	m.outputDirInput = outputDirValue(m.input.Value(), m.defaultOutputDir)
	if m.sourceScreen == screenURLInput {
		m.plan.Output = m.outputDirInput
		m.screen = screenFilenameInput
		m.errorMessage = ""
		m.input = m.newFilenameInput(m.filenameInput)
		return m, m.input.Focus()
	}

	m.plan.Output = m.outputDirInput
	m.screen = screenPlan
	m.errorMessage = ""
	m.input.Blur()
	return m, nil
}

func (m Model) submitFilenameInput() (Model, tea.Cmd) {
	filename, err := filenameValue(m.input.Value())
	if err != nil {
		m.errorMessage = err.Error()
		return m, nil
	}

	plan, err := planFromURL(m.sourceInput, m.outputDirInput, filename, m.checksumInput, m.retries, m.resume)
	if err != nil {
		m.errorMessage = err.Error()
		return m, nil
	}

	m.filenameInput = filename
	m.plan = plan
	m.screen = screenChecksumInput
	m.errorMessage = ""
	m.input = m.newChecksumInput(m.checksumInput)
	return m, m.input.Focus()
}

func (m Model) submitChecksumInput() (Model, tea.Cmd) {
	checksum := m.input.Value()
	plan, err := planFromURL(m.sourceInput, m.outputDirInput, m.filenameInput, checksum, m.retries, m.resume)
	if err != nil {
		m.errorMessage = err.Error()
		return m, nil
	}

	m.checksumInput = checksum
	m.plan = plan
	m.screen = screenPlan
	m.errorMessage = ""
	m.input.Blur()
	return m, nil
}

func (m Model) startExecution() (Model, tea.Cmd) {
	ctx, cancel := context.WithCancel(context.Background())
	m.screen = screenExecution
	m.errorMessage = ""
	m.execution = newExecutionState(m.plan)
	m.executionCancel = cancel
	m.executionMessages = runExecution(ctx, m.plan, m.executionRunner)
	return m, waitForExecution(m.executionMessages)
}
