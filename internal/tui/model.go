package tui

import (
	"github.com/he8um/daryaft/pkg/version"

	tea "github.com/charmbracelet/bubbletea"
)

type Options struct {
	NoColor bool
}

type Model struct {
	screen   screen
	selected int
	styles   styles
	version  version.Details
}

func NewModel(options Options) Model {
	return Model{
		screen:  screenHome,
		styles:  newStyles(options.NoColor),
		version: version.Info(),
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
	m.screen = item.screen
	return m, nil
}

func (m Model) back() Model {
	m.screen = screenHome
	return m
}
