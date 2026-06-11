package download

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/he8um/daryaft/internal/checksum"
	"github.com/he8um/daryaft/internal/httpopts"
	"github.com/he8um/daryaft/internal/input"
)

const MaxRetries = 20

func BuildPlan(options Options) (Plan, error) {
	if options.NoResume {
		options.Resume = false
	}

	if err := ValidateRetries(options.Retries); err != nil {
		return Plan{}, err
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

	parsedChecksum, err := parsePlanChecksum(options, urls)
	if err != nil {
		return Plan{}, err
	}

	targetChecksums, hasChecksumFile, err := parsePlanChecksumFile(options, urls)
	if err != nil {
		return Plan{}, err
	}

	if err := httpopts.Validate(options.HTTPOptions); err != nil {
		return Plan{}, err
	}

	return Plan{
		URLs:            urls,
		Output:          strings.TrimSpace(options.Output),
		Name:            strings.TrimSpace(options.Name),
		Checksum:        parsedChecksum,
		TargetChecksums: targetChecksums,
		HasChecksumFile: hasChecksumFile,
		Retries:         options.Retries,
		Resume:          options.Resume,
		HTTPOptions:     options.HTTPOptions,
	}, nil
}

func parsePlanChecksum(options Options, urls []string) (*checksum.Spec, error) {
	rawChecksum := strings.TrimSpace(options.Checksum)
	if rawChecksum == "" {
		return nil, nil
	}

	if strings.TrimSpace(options.ChecksumFile) != "" {
		return nil, fmt.Errorf("--checksum and --checksum-file cannot be used together")
	}

	spec, err := checksum.Parse(rawChecksum)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(options.File) != "" || len(urls) != 1 {
		return nil, fmt.Errorf("--checksum is currently supported only for single URL downloads")
	}
	return &spec, nil
}

// parsePlanChecksumFile parses and cross-validates a checksum manifest file
// against the planned download targets. It returns a map of URL to checksum
// Spec and a flag indicating that a checksum file was supplied.
//
// The manifest must cover every planned target exactly once, and every manifest
// URL must match a planned target exactly (no normalization).
func parsePlanChecksumFile(options Options, urls []string) (map[string]checksum.Spec, bool, error) {
	rawPath := strings.TrimSpace(options.ChecksumFile)
	if rawPath == "" {
		return nil, false, nil
	}

	entries, err := checksum.ParseManifestFile(rawPath)
	if err != nil {
		return nil, false, fmt.Errorf("checksum file: %w", err)
	}
	if len(entries) == 0 {
		return nil, false, fmt.Errorf("checksum file: no checksum entries found")
	}

	targetSet := make(map[string]struct{}, len(urls))
	for _, rawURL := range urls {
		targetSet[rawURL] = struct{}{}
	}

	for _, rawURL := range urls {
		if _, ok := entries[rawURL]; !ok {
			return nil, false, fmt.Errorf("checksum file: no checksum provided for URL: %s", rawURL)
		}
	}

	for manifestURL := range entries {
		if _, ok := targetSet[manifestURL]; !ok {
			return nil, false, fmt.Errorf("checksum file: manifest URL not in download targets: %s", manifestURL)
		}
	}

	return entries, true, nil
}

func ValidateRetries(retries int) error {
	if retries < 0 || retries > MaxRetries {
		return fmt.Errorf("retries must be between 0 and %d", MaxRetries)
	}
	return nil
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

// ValidateURL validates a single URL for download use.
func ValidateURL(rawURL string) error {
	return validateURL(rawURL)
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
