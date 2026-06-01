package doctor

import "encoding/json"

type JSONReport struct {
	OK       bool          `json:"ok"`
	Strict   bool          `json:"strict,omitempty"`
	Summary  JSONSummary   `json:"summary"`
	Sections []JSONSection `json:"sections"`
}

type JSONSummary struct {
	Failures int `json:"failures"`
	Warnings int `json:"warnings"`
	Checks   int `json:"checks"`
}

type JSONSection struct {
	Name   string      `json:"name"`
	Checks []JSONCheck `json:"checks"`
}

type JSONCheck struct {
	Status  string `json:"status"`
	Label   string `json:"label"`
	Message string `json:"message"`
}

func FormatJSON(report Report) ([]byte, error) {
	return FormatJSONWithOptions(report, false)
}

func FormatJSONWithOptions(report Report, strict bool) ([]byte, error) {
	dto := ToJSONReportWithOptions(report, strict)
	return json.MarshalIndent(dto, "", "  ")
}

func ToJSONReport(report Report) JSONReport {
	return ToJSONReportWithOptions(report, false)
}

func ToJSONReportWithOptions(report Report, strict bool) JSONReport {
	dto := JSONReport{
		OK:     report.OK(strict),
		Strict: strict,
		Summary: JSONSummary{
			Checks: len(report.Checks),
		},
	}

	sectionIndexes := make(map[string]int)
	for _, check := range report.Checks {
		switch check.Status {
		case StatusFail:
			dto.Summary.Failures++
		case StatusWarn:
			dto.Summary.Warnings++
		}

		index, ok := sectionIndexes[check.Section]
		if !ok {
			index = len(dto.Sections)
			sectionIndexes[check.Section] = index
			dto.Sections = append(dto.Sections, JSONSection{Name: check.Section})
		}
		dto.Sections[index].Checks = append(dto.Sections[index].Checks, JSONCheck{
			Status:  statusName(check.Status),
			Label:   check.Label,
			Message: check.Value,
		})
	}

	return dto
}

func statusName(status Status) string {
	switch status {
	case StatusOK:
		return "ok"
	case StatusFail:
		return "failure"
	case StatusWarn:
		return "warning"
	case StatusInfo:
		return "info"
	case StatusSkipped:
		return "skipped"
	default:
		return "info"
	}
}
