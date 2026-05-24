package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(options Options) error {
	program := tea.NewProgram(NewModel(options), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("start TUI: %w", err)
	}
	return nil
}
