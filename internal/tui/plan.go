package tui

import (
	"fmt"
	"strings"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/download"
)

const (
	tuiDefaultRetries = 3
	tuiDefaultResume  = true
)

func planFromURL(rawURL, output, name, checksum string, retries int, resume bool) (download.Plan, error) {
	return download.BuildPlan(download.Options{
		URLs:     []string{strings.TrimSpace(rawURL)},
		Output:   output,
		Name:     name,
		DryRun:   true,
		Checksum: checksum,
		Retries:  retries,
		Resume:   resume,
	})
}

func planFromFile(path, output string, retries int, resume bool) (download.Plan, error) {
	return download.BuildPlan(download.Options{
		File:    strings.TrimSpace(path),
		Output:  output,
		DryRun:  true,
		Retries: retries,
		Resume:  resume,
	})
}

func defaultOutputDir(value string) string {
	return config.EffectiveDownloadDir(value)
}

func outputDirValue(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return defaultOutputDir(fallback)
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
