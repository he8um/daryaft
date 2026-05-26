package tui

import (
	"strings"

	"github.com/he8um/daryaft/internal/download"
)

const (
	tuiDefaultRetries = 3
	tuiDefaultResume  = true
)

func planFromURL(rawURL, output string) (download.Plan, error) {
	return download.BuildPlan(download.Options{
		URLs:    []string{strings.TrimSpace(rawURL)},
		Output:  output,
		DryRun:  true,
		Retries: tuiDefaultRetries,
		Resume:  tuiDefaultResume,
	})
}

func planFromFile(path, output string) (download.Plan, error) {
	return download.BuildPlan(download.Options{
		File:    strings.TrimSpace(path),
		Output:  output,
		DryRun:  true,
		Retries: tuiDefaultRetries,
		Resume:  tuiDefaultResume,
	})
}

func outputDirValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "."
	}
	return trimmed
}

func displayValue(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
