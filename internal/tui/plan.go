package tui

import (
	"strings"

	"github.com/he8um/daryaft/internal/download"
)

const (
	tuiDefaultRetries = 3
	tuiDefaultResume  = true
)

func planFromURL(rawURL string) (download.Plan, error) {
	return download.BuildPlan(download.Options{
		URLs:    []string{strings.TrimSpace(rawURL)},
		DryRun:  true,
		Retries: tuiDefaultRetries,
		Resume:  tuiDefaultResume,
	})
}

func planFromFile(path string) (download.Plan, error) {
	return download.BuildPlan(download.Options{
		File:    strings.TrimSpace(path),
		DryRun:  true,
		Retries: tuiDefaultRetries,
		Resume:  tuiDefaultResume,
	})
}

func displayValue(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
