package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadURLFile(path string) ([]string, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("URL file path cannot be empty")
	}

	file, err := os.Open(path) // #nosec G304 -- users intentionally choose the URL list file path.
	if err != nil {
		return nil, fmt.Errorf("read URL file %q: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	var urls []string
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read URL file %q: line %d: %w", path, lineNumber, err)
	}

	return urls, nil
}
