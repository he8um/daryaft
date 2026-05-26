package tui

import (
	"fmt"
	"strings"

	"github.com/he8um/daryaft/internal/download"
)

const (
	tuiDefaultRetries = 3
	tuiDefaultResume  = true
)

func planFromURL(rawURL, output, name string) (download.Plan, error) {
	return download.BuildPlan(download.Options{
		URLs:    []string{strings.TrimSpace(rawURL)},
		Output:  output,
		Name:    name,
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

func filenameValue(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}
	if trimmed == "." || trimmed == ".." {
		return "", fmt.Errorf("filename cannot be %q", trimmed)
	}
	if strings.ContainsAny(trimmed, `/\`) {
		return "", fmt.Errorf("filename cannot contain path separators")
	}
	if strings.Contains(trimmed, "\x00") {
		return "", fmt.Errorf("filename cannot contain null bytes")
	}
	return trimmed, nil
}

func displayValue(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
