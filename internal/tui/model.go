package tui

import (
	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/pkg/version"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Options struct {
	NoColor bool
}

type Model struct {
	screen       screen
	selected     int
	styles       styles
	version      version.Details
	input        textinput.Model
	errorMessage string
	plan         download.Plan
	planReturn   screen
}

func NewModel(options Options) Model {
	styles := newStyles(options.NoColor)
	return Model{
		screen:     screenHome,
		styles:     styles,
		version:    version.Info(),
		input:      newTextInput(styles),
		planReturn: screenURLInput,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) moveUp() Model {
	if m.selected <= 0 {
		m.selected = len(homeMenu) - 1
		return m
	}
	m.selected--
	return m
}

func (m Model) moveDown() Model {
	if m.selected >= len(homeMenu)-1 {
		m.selected = 0
		return m
	}
	m.selected++
	return m
}

func (m Model) enter() (Model, tea.Cmd) {
	item := homeMenu[m.selected]
	if item.title == "Quit" {
		return m, tea.Quit
	}
	if item.screen == screenURLInput || item.screen == screenFileInput {
		return m.openInput(item.screen)
	}
	m.screen = item.screen
	return m, nil
}

func (m Model) openInput(next screen) (Model, tea.Cmd) {
	m.screen = next
	m.input = newTextInput(m.styles)
	m.errorMessage = ""
	m.plan = download.Plan{}
	m.planReturn = next
	return m, m.input.Focus()
}

func (m Model) back() (Model, tea.Cmd) {
	if m.screen == screenPlan {
		m.screen = m.planReturn
		return m, m.input.Focus()
	}
	m.screen = screenHome
	m.errorMessage = ""
	return m, nil
}

func (m Model) home() Model {
	m.screen = screenHome
	m.errorMessage = ""
	return m
}

func (m Model) isInputScreen() bool {
	return m.screen == screenURLInput || m.screen == screenFileInput
}
