package tui

import (
	"fmt"
	"strings"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/utils"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var content string
	switch m.screen {
	case screenHome:
		content = m.homeView()
	case screenURLInput, screenFileInput:
		content = m.inputView()
	case screenPlan:
		content = m.planView()
	case screenExecution:
		content = m.executionView()
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

func (m Model) inputView() string {
	var builder strings.Builder
	builder.WriteString(m.headerView())
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.title.Render(m.screen.title()))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.body.Render(m.inputPrompt()))
	builder.WriteString("\n")
	builder.WriteString(m.input.View())
	builder.WriteString("\n\n")
	if m.errorMessage != "" {
		builder.WriteString(m.styles.error.Render("Error: " + m.errorMessage))
	} else {
		builder.WriteString(m.styles.muted.Render(" "))
	}
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.muted.Render("enter plan • esc/backspace home • q quit"))
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func (m Model) planView() string {
	var builder strings.Builder
	builder.WriteString(m.headerView())
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.title.Render("Download plan"))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.body.Render(m.planBody()))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.muted.Render("Review the plan before starting the download."))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.muted.Render("enter start download • esc/backspace edit • h home • q quit"))
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func (m Model) planBody() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Number of URLs: %d\n", len(m.plan.URLs))
	for index, rawURL := range m.plan.URLs {
		if index >= 8 {
			fmt.Fprintf(&builder, "... and %d more\n", len(m.plan.URLs)-index)
			break
		}
		fmt.Fprintf(&builder, "%d. %s\n", index+1, rawURL)
	}
	fmt.Fprintf(&builder, "Output: %s\n", displayValue(m.plan.Output, "current directory"))
	fmt.Fprintf(&builder, "Filename: %s\n", displayValue(m.plan.Name, "auto-detect"))
	fmt.Fprintf(&builder, "Retries: %d\n", m.plan.Retries)
	fmt.Fprintf(&builder, "Resume: %t", m.plan.Resume)
	return builder.String()
}

func (m Model) executionView() string {
	var builder strings.Builder
	builder.WriteString(m.headerView())
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.title.Render("Downloading"))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.body.Render(m.executionBody()))
	builder.WriteString("\n\n")
	if m.execution.Running {
		builder.WriteString(m.styles.muted.Render("download running • q cancel • ctrl+c quit"))
	} else {
		builder.WriteString(m.styles.muted.Render("enter/h home • q quit"))
	}
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func (m Model) executionBody() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Item %d of %d\n", valueOrOne(m.execution.ItemIndex), valueOrOne(m.execution.ItemTotal))
	fmt.Fprintf(&builder, "Current URL: %s\n", displayValue(m.execution.CurrentURL, "pending"))
	fmt.Fprintf(&builder, "Target path: %s\n", displayValue(m.execution.TargetPath, "pending"))
	fmt.Fprintf(&builder, "Status: %s\n", displayValue(m.execution.Status, "Starting"))
	fmt.Fprintf(&builder, "Downloaded: %s\n", byteProgress(m.execution.DownloadedBytes, m.execution.TotalBytes))
	if m.execution.TotalBytes > 0 {
		fmt.Fprintf(&builder, "Percent: %.1f%%\n", m.execution.Percent)
	}
	if m.execution.Speed > 0 {
		fmt.Fprintf(&builder, "Speed: %s\n", utils.FormatSpeed(m.execution.Speed))
	}
	if m.execution.Message != "" {
		fmt.Fprintf(&builder, "Message: %s\n", m.execution.Message)
	}
	if m.execution.Done {
		builder.WriteString("\n")
		builder.WriteString(summaryView(m.execution.Summary))
	}
	return strings.TrimRight(builder.String(), "\n")
}

func byteProgress(downloaded, total int64) string {
	if total > 0 {
		return fmt.Sprintf("%s / %s", utils.FormatBytes(downloaded), utils.FormatBytes(total))
	}
	return utils.FormatBytes(downloaded)
}

func summaryView(summary executionSummary) string {
	var builder strings.Builder
	builder.WriteString("Summary\n")
	fmt.Fprintf(&builder, "Total: %d\n", summary.Total)
	fmt.Fprintf(&builder, "Completed: %d\n", summary.Completed)
	fmt.Fprintf(&builder, "Failed: %d\n", summary.Failed)
	fmt.Fprintf(&builder, "Cancelled: %d", summary.Cancelled)
	if summary.Skipped > 0 {
		fmt.Fprintf(&builder, "\nSkipped: %d", summary.Skipped)
	}
	if len(summary.Failures) == 0 {
		return builder.String()
	}

	builder.WriteString("\nFailures:\n")
	for index, failure := range summary.Failures {
		if index >= 3 {
			fmt.Fprintf(&builder, "... and %d more", len(summary.Failures)-index)
			break
		}
		fmt.Fprintf(&builder, "- %s: %s", failure.URL, failure.Error)
		if index < len(summary.Failures)-1 && index < 2 {
			builder.WriteString("\n")
		}
	}
	return strings.TrimRight(builder.String(), "\n")
}

func valueOrOne(value int) int {
	if value <= 0 {
		return 1
	}
	return value
}

func (m Model) footerView() string {
	return m.styles.muted.Render(config.FooterText)
}
