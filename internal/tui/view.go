package tui

import (
	"fmt"
	"strings"

	"github.com/he8um/daryaft/internal/config"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var content string
	switch m.screen {
	case screenHome:
		content = m.homeView()
	case screenURLPlanned, screenFilePlanned:
		content = m.plannedView()
	case screenHelp:
		content = m.helpView()
	case screenVersion:
		content = m.versionView()
	default:
		content = m.homeView()
	}

	return m.styles.panel.Render(content)
}

func (m Model) headerView() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.title.Render(config.AppName),
		m.styles.subtitle.Render("Modern terminal downloader"),
	)
}

func (m Model) homeView() string {
	var builder strings.Builder
	builder.WriteString(m.headerView())
	builder.WriteString("\n\n")

	for index, item := range homeMenu {
		label := fmt.Sprintf("%d. %s", index+1, item.title)
		if index == m.selected {
			builder.WriteString(m.styles.selected.Render("> " + label))
		} else {
			builder.WriteString(m.styles.item.Render("  " + label))
		}
		builder.WriteByte('\n')
	}

	builder.WriteString("\n")
	builder.WriteString(m.styles.muted.Render("↑/k ↓/j move • enter select • q quit"))
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func (m Model) plannedView() string {
	return m.subscreenView(
		m.screen.title(),
		"Planned for the download TUI milestone.\n\nUse the existing CLI download commands for now.",
	)
}

func (m Model) helpView() string {
	return m.subscreenView(
		"Help",
		"Use ↑/k and ↓/j to move through the home menu.\nPress enter to select an item.\nPress esc or backspace to return home.\nPress q or ctrl+c to quit.",
	)
}

func (m Model) versionView() string {
	info := m.version
	return m.subscreenView(
		"Version",
		fmt.Sprintf("%s version: %s\ncommit: %s\nbuilt: %s\ngo version: %s", config.AppName, info.Version, info.Commit, info.Date, info.GoVersion),
	)
}

func (m Model) subscreenView(title, body string) string {
	var builder strings.Builder
	builder.WriteString(m.headerView())
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.title.Render(title))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.body.Render(body))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.muted.Render("esc/backspace home • q quit"))
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func (m Model) footerView() string {
	return m.styles.muted.Render(config.FooterText)
}
