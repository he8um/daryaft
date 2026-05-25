package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	panel    lipgloss.Style
	title    lipgloss.Style
	subtitle lipgloss.Style
	item     lipgloss.Style
	selected lipgloss.Style
	muted    lipgloss.Style
	body     lipgloss.Style
	error    lipgloss.Style
}

func newStyles(noColor bool) styles {
	s := styles{
		panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 3).
			Margin(1, 2).
			Width(56),
		title: lipgloss.NewStyle().
			Bold(true),
		subtitle: lipgloss.NewStyle(),
		item: lipgloss.NewStyle().
			PaddingLeft(2),
		selected: lipgloss.NewStyle().
			Bold(true).
			PaddingLeft(2),
		muted: lipgloss.NewStyle(),
		body: lipgloss.NewStyle().
			MarginTop(1),
		error: lipgloss.NewStyle(),
	}

	if noColor {
		return s
	}

	s.panel = s.panel.BorderForeground(lipgloss.Color("63"))
	s.title = s.title.Foreground(lipgloss.Color("81"))
	s.subtitle = s.subtitle.Foreground(lipgloss.Color("245"))
	s.selected = s.selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("63"))
	s.muted = s.muted.Foreground(lipgloss.Color("242"))
	s.error = s.error.Foreground(lipgloss.Color("203"))
	return s
}
