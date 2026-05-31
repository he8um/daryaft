package doctor

import "strings"

func Format(report Report) string {
	var builder strings.Builder
	builder.WriteString("Daryaft doctor\n")

	currentSection := ""
	for _, check := range report.Checks {
		if check.Section != currentSection {
			currentSection = check.Section
			builder.WriteString("\n")
			builder.WriteString(currentSection)
			builder.WriteString("\n")
		}
		builder.WriteString(statusSymbol(check.Status))
		builder.WriteString(" ")
		builder.WriteString(check.Label)
		builder.WriteString(": ")
		builder.WriteString(check.Value)
		builder.WriteString("\n")
	}

	return builder.String()
}

func statusSymbol(status Status) string {
	switch status {
	case StatusOK:
		return "✓"
	case StatusFail:
		return "✗"
	case StatusWarn:
		return "!"
	case StatusInfo:
		return "-"
	case StatusSkipped:
		return "-"
	default:
		return "?"
	}
}
