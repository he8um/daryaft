package checksum

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseManifestFile reads a checksum manifest file and returns a map of URL to
// the parsed checksum Spec for that URL.
//
// The manifest format is one entry per line:
//
//	<algorithm>:<hex> <url>
//
// Blank lines and lines whose trimmed content begins with "#" are ignored.
// Each remaining line must contain exactly two whitespace-separated fields: the
// checksum spec and the URL. Duplicate URLs are rejected. All errors include
// the offending line number.
//
// An empty manifest (or a comments-only manifest) returns an empty map and a
// nil error. Cross-validation against planned download targets is the caller's
// responsibility.
func ParseManifestFile(path string) (map[string]Spec, error) {
	// #nosec G304 -- users intentionally choose the checksum manifest file path.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	entries := make(map[string]Spec)
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("manifest line %d: expected \"<algorithm>:<hex> <url>\" format", lineNumber)
		}

		spec, err := Parse(fields[0])
		if err != nil {
			return nil, fmt.Errorf("manifest line %d: %w", lineNumber, err)
		}

		rawURL := fields[1]
		if _, exists := entries[rawURL]; exists {
			return nil, fmt.Errorf("manifest line %d: duplicate URL %s", lineNumber, rawURL)
		}
		entries[rawURL] = spec
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read checksum manifest %q: line %d: %w", path, lineNumber, err)
	}

	return entries, nil
}
