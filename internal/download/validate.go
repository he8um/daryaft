package download

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/he8um/daryaft/internal/input"
)

func BuildPlan(options Options) (Plan, error) {
	if options.NoResume {
		options.Resume = false
	}

	if options.Retries < 0 {
		return Plan{}, fmt.Errorf("retries must be greater than or equal to 0")
	}

	urls, err := collectURLs(options)
	if err != nil {
		return Plan{}, err
	}

	if len(urls) == 0 {
		return Plan{}, fmt.Errorf("provide at least one URL or use --file")
	}

	for index, rawURL := range urls {
		if err := validateURL(rawURL); err != nil {
			return Plan{}, fmt.Errorf("invalid URL %d %q: %w", index+1, rawURL, err)
		}
	}

	if strings.TrimSpace(options.Name) != "" && len(urls) > 1 {
		return Plan{}, fmt.Errorf("--name can only be used with a single URL")
	}

	return Plan{
		URLs:    urls,
		Output:  strings.TrimSpace(options.Output),
		Name:    strings.TrimSpace(options.Name),
		Retries: options.Retries,
		Resume:  options.Resume,
	}, nil
}

func collectURLs(options Options) ([]string, error) {
	var urls []string
	for _, rawURL := range options.URLs {
		trimmed := strings.TrimSpace(rawURL)
		if trimmed == "" {
			return nil, fmt.Errorf("URL arguments cannot be empty")
		}
		urls = append(urls, trimmed)
	}

	if strings.TrimSpace(options.File) == "" {
		return urls, nil
	}

	fileURLs, err := input.ReadURLFile(options.File)
	if err != nil {
		return nil, err
	}

	urls = append(urls, fileURLs...)
	return urls, nil
}

func validateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("scheme must be http or https")
	}
	if parsed.Host == "" {
		return fmt.Errorf("host is required")
	}
	return nil
}
