package tui

import (
	"fmt"
	"strings"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/inspect"
	"github.com/he8um/daryaft/internal/utils"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var content string
	switch m.screen {
	case screenHome:
		content = m.homeView()
	case screenURLInput, screenFileInput, screenInspectInput, screenOutputInput, screenFilenameInput, screenChecksumInput:
		content = m.inputView()
	case screenPlan:
		content = m.planView()
	case screenExecution:
		content = m.executionView()
	case screenInspectExecution:
		content = m.inspectExecutionView()
	case screenInspectResult:
		content = m.inspectResultView()
	case screenInspectError:
		content = m.inspectErrorView()
	case screenHelp:
		content = m.helpView()
	case screenVersion:
		content = m.versionView()
	default:
		content = m.homeView()
	}

	return m.styles.panel.Width(m.panelWidth()).Render(content)
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
	if m.screen == screenOutputInput {
		builder.WriteString(m.styles.muted.Render(fmt.Sprintf("Default/current value: %s", displayValue(m.outputDirInput, "."))))
		builder.WriteString("\n")
	}
	if m.screen == screenFilenameInput {
		builder.WriteString(m.styles.muted.Render("Leave empty to auto-detect"))
		builder.WriteString("\n")
	}
	if m.screen == screenChecksumInput {
		builder.WriteString(m.styles.muted.Render("Leave empty to skip. Format: sha256:<hex> or sha512:<hex>"))
		builder.WriteString("\n")
	}
	builder.WriteString(m.input.View())
	builder.WriteString("\n\n")
	if m.errorMessage != "" {
		builder.WriteString(m.styles.error.Render("Error: " + m.errorMessage))
	} else {
		builder.WriteString(m.styles.muted.Render(" "))
	}
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.muted.Render(m.inputHelp()))
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
	fmt.Fprintf(&builder, "Output: %s\n", displayValue(m.plan.Output, "."))
	fmt.Fprintf(&builder, "Filename: %s\n", displayValue(m.plan.Name, "auto-detect"))
	if m.plan.Checksum == nil {
		fmt.Fprintln(&builder, "Checksum: none")
	} else {
		fmt.Fprintf(&builder, "Checksum: %s\n", m.plan.Checksum.String())
	}
	fmt.Fprintf(&builder, "Retries: %d\n", m.plan.Retries)
	fmt.Fprintf(&builder, "Resume: %t", m.plan.Resume)
	return builder.String()
}

func (m Model) inputHelp() string {
	if m.screen == screenOutputInput {
		if m.sourceScreen == screenURLInput {
			return "enter next • esc previous • backspace empty previous • q quit"
		}
		return "enter plan • esc previous • backspace empty previous • q quit"
	}
	if m.screen == screenFilenameInput {
		return "enter next • esc previous • backspace empty previous • q quit"
	}
	if m.screen == screenChecksumInput {
		return "enter plan • esc previous • backspace empty previous • q quit"
	}
	if m.screen == screenInspectInput {
		return "enter inspect • esc home • backspace empty home • q quit"
	}
	return "enter next • esc home • backspace empty home • q quit"
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
		builder.WriteString(m.styles.muted.Render("enter/h new download • q quit"))
	}
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func (m Model) executionBody() string {
	var builder strings.Builder
	if len(m.execution.Items) > 0 {
		fmt.Fprintf(&builder, "Queue\n")

		displayItems := m.execution.Items
		const maxDisplay = 20
		if len(displayItems) > maxDisplay {
			displayItems = displayItems[len(displayItems)-maxDisplay:]
		}

		for _, item := range displayItems {
			marker := statusMarker(item.Status, m.noColor)
			fmt.Fprintf(&builder, "  %s %d. %s\n", marker, item.Index, truncateURL(item.URL, 55))
		}

		if len(m.execution.Items) > maxDisplay {
			fmt.Fprintf(&builder, "  ... and %d more\n", len(m.execution.Items)-maxDisplay)
		}

		builder.WriteString("\n")
	}
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

func (m Model) inspectExecutionView() string {
	var builder strings.Builder
	builder.WriteString(m.headerView())
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.title.Render("Inspecting"))
	builder.WriteString("\n\n")
	fmt.Fprintf(&builder, "URL: %s\n", displayValue(m.inspect.URL, "pending"))
	fmt.Fprintf(&builder, "Status: %s\n", displayValue(m.inspect.Status, "Inspecting"))
	if m.inspect.Message != "" {
		fmt.Fprintf(&builder, "Message: %s\n", m.inspect.Message)
	}
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.muted.Render("inspect running • q cancel • ctrl+c quit"))
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func (m Model) inspectResultView() string {
	var builder strings.Builder
	builder.WriteString(m.headerView())
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.title.Render("Inspect result"))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.body.Render(inspectResultBody(m.inspect.Result)))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.muted.Render("enter/h home • esc/backspace edit • q quit"))
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func (m Model) inspectErrorView() string {
	var builder strings.Builder
	builder.WriteString(m.headerView())
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.title.Render("Inspect error"))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.error.Render(displayValue(m.inspect.Error, "Inspect failed")))
	builder.WriteString("\n\n")
	builder.WriteString(m.styles.muted.Render("esc/backspace edit • h home • q quit"))
	builder.WriteString("\n\n")
	builder.WriteString(m.footerView())
	return strings.TrimRight(builder.String(), "\n")
}

func inspectResultBody(result inspect.Result) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "URL: %s\n", displayValue(result.URL, "unknown"))
	fmt.Fprintf(&builder, "Final URL: %s\n", displayValue(result.FinalURL, "unknown"))
	fmt.Fprintf(&builder, "Status: %s\n", displayValue(result.Status, "unknown"))
	fmt.Fprintf(&builder, "Filename: %s\n", displayValue(result.Filename, "unknown"))
	if result.ContentLengthKnown {
		fmt.Fprintf(&builder, "Content length: %d bytes\n", result.ContentLength)
	} else {
		builder.WriteString("Content length: unknown\n")
	}
	fmt.Fprintf(&builder, "Content type: %s\n", displayValue(result.ContentType, "unknown"))
	fmt.Fprintf(&builder, "Accept-Ranges: %s\n", displayValue(result.AcceptRanges, "unknown"))
	fmt.Fprintf(&builder, "Resume supported: %s\n", inspectResumeSupport(result))
	fmt.Fprintf(&builder, "ETag: %s\n", displayValue(result.ETag, "none"))
	fmt.Fprintf(&builder, "Last-Modified: %s", displayValue(result.LastModified, "none"))
	return builder.String()
}

func inspectResumeSupport(result inspect.Result) string {
	if !result.ResumeSupportKnown {
		return "unknown"
	}
	if result.ResumeSupported {
		return "yes"
	}
	return "no"
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
		fmt.Fprintf(&builder, "\nNot started: %d", summary.Skipped)
	}
	if summary.ChecksumVerified > 0 {
		fmt.Fprintf(&builder, "\nChecksum verified: %d", summary.ChecksumVerified)
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

func statusMarker(status string, noColor bool) string {
	switch status {
	case "Completed", "Checksum OK":
		if noColor {
			return "[ok]"
		}
		return "✓"
	case "Failed", "Cancelled", "Checksum Failed":
		if noColor {
			return "[!]"
		}
		return "✗"
	case "Downloading", "Starting", "Resuming", "Restarting", "Verifying":
		if noColor {
			return "[>]"
		}
		return "→"
	default:
		if noColor {
			return "[-]"
		}
		return "·"
	}
}

func truncateURL(u string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(u) <= max {
		return u
	}
	if max <= 3 {
		return u[:max]
	}
	return u[:max-3] + "..."
}
