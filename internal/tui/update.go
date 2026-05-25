package tui

import (
	"github.com/he8um/daryaft/internal/download"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		if m.isInputScreen() {
			return m.updateTextInput(msg)
		}
		return m, nil
	}

	if isQuitKey(key) {
		return m, tea.Quit
	}

	switch m.screen {
	case screenPlan:
		if isBackKey(key) {
			return m.back()
		}
		if key.String() == "h" {
			return m.home(), nil
		}
		return m, nil
	case screenURLInput, screenFileInput:
		if isBackKey(key) {
			return m.back()
		}
		if key.String() == "enter" {
			return m.submitInput()
		}
		return m.updateTextInput(msg)
	case screenHelp, screenVersion:
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
	default:
		return m, nil
	}
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

func (m Model) submitInput() (Model, tea.Cmd) {
	var (
		plan download.Plan
		err  error
	)

	switch m.screen {
	case screenURLInput:
		plan, err = planFromURL(m.input.Value())
	case screenFileInput:
		plan, err = planFromFile(m.input.Value())
	default:
		return m, nil
	}

	if err != nil {
		m.errorMessage = err.Error()
		return m, nil
	}

	m.plan = plan
	m.planReturn = m.screen
	m.screen = screenPlan
	m.errorMessage = ""
	m.input.Blur()
	return m, nil
}
