package tui

import tea "github.com/charmbracelet/bubbletea"

func isUpKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "up", "k":
		return true
	default:
		return false
	}
}

func isDownKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "down", "j":
		return true
	default:
		return false
	}
}

func isBackKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "esc", "backspace":
		return true
	default:
		return false
	}
}

func isQuitKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "q", "ctrl+c":
		return true
	default:
		return false
	}
}
