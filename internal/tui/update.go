package tui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if isQuitKey(key) {
		return m, tea.Quit
	}

	if m.screen != screenHome {
		if isBackKey(key) {
			return m.back(), nil
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
