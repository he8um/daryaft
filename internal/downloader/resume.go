package downloader

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	resumeNotSupportedMessage = "Resume not supported by server; restarting download"
	remoteChangedMessage      = "Remote file changed; restarting download"
)

type resumeCandidate struct {
	Target       targetPaths
	MetadataPath string
	Metadata     partialMetadata
	HasMetadata  bool
	PartialSize  int64
	CanResume    bool
}

func findResumeCandidate(outputDir, rawURL, customName string) (resumeCandidate, error) {
	candidate, found, err := findMetadataCandidate(outputDir, rawURL)
	if err != nil || found {
		return candidate, err
	}

	filename := customName
	if strings.TrimSpace(filename) == "" {
		filename = filenameFromURL(rawURL)
	}
	if strings.TrimSpace(filename) == "" {
		filename = fallbackFilename
	}

	target, err := targetPathsFor(outputDir, filename)
	if err != nil {
		return resumeCandidate{}, err
	}

	info, err := os.Stat(target.Partial)
	if err != nil {
		if os.IsNotExist(err) {
			return resumeCandidate{}, nil
		}
		return resumeCandidate{}, fmt.Errorf("check partial file %q: %w", target.Partial, err)
	}
	if info.IsDir() {
		return resumeCandidate{}, fmt.Errorf("partial path is a directory: %s", target.Partial)
	}

	if _, err := prepareTarget(outputDir, filename); err != nil {
		return resumeCandidate{}, err
	}

	return resumeCandidate{
		Target:       target,
		MetadataPath: metadataPathForPartial(target.Partial),
		PartialSize:  info.Size(),
		CanResume:    info.Size() > 0,
	}, nil
}

func findMetadataCandidate(outputDir, rawURL string) (resumeCandidate, bool, error) {
	if outputDir == "" {
		outputDir = "."
	}
	outputDir = filepath.Clean(outputDir)

	matches, err := filepath.Glob(filepath.Join(outputDir, "*.part.daryaft.json"))
	if err != nil {
		return resumeCandidate{}, false, fmt.Errorf("find partial metadata: %w", err)
	}

	for _, path := range matches {
		metadata, err := loadPartialMetadata(path)
		if err != nil {
			continue
		}
		if metadata.URL != rawURL {
			continue
		}

		target := targetPaths{
			OutputDir: outputDir,
			Final:     filepath.Clean(metadata.TargetPath),
			Partial:   filepath.Clean(metadata.PartialPath),
		}
		if err := ensureInsideOutputDir(outputDir, target.Final); err != nil {
			return resumeCandidate{}, false, err
		}
		if err := ensureInsideOutputDir(outputDir, target.Partial); err != nil {
			return resumeCandidate{}, false, err
		}

		if _, err := os.Stat(target.Final); err == nil {
			return resumeCandidate{}, false, fmt.Errorf("%w: %s", ErrTargetExists, target.Final)
		} else if !os.IsNotExist(err) {
			return resumeCandidate{}, false, fmt.Errorf("check target file %q: %w", target.Final, err)
		}

		info, err := os.Stat(target.Partial)
		if err != nil {
			if os.IsNotExist(err) {
				return resumeCandidate{
					Target:       target,
					MetadataPath: path,
					Metadata:     metadata,
					HasMetadata:  true,
				}, true, nil
			}
			return resumeCandidate{}, false, fmt.Errorf("check partial file %q: %w", target.Partial, err)
		}
		if info.IsDir() {
			return resumeCandidate{}, false, fmt.Errorf("partial path is a directory: %s", target.Partial)
		}

		return resumeCandidate{
			Target:       target,
			MetadataPath: path,
			Metadata:     metadata,
			HasMetadata:  true,
			PartialSize:  info.Size(),
			CanResume:    info.Size() > 0,
		}, true, nil
	}

	return resumeCandidate{}, false, nil
}

func remoteChanged(metadata partialMetadata, header http.Header) bool {
	if metadata.ETag != "" {
		if etag := header.Get("ETag"); etag != "" && etag != metadata.ETag {
			return true
		}
	}
	if metadata.LastModified != "" {
		if lastModified := header.Get("Last-Modified"); lastModified != "" && lastModified != metadata.LastModified {
			return true
		}
	}
	return false
}

func responseTotalBytes(response *http.Response, startOffset int64) int64 {
	if response == nil {
		return 0
	}

	if response.StatusCode == http.StatusPartialContent {
		if _, total, ok := parseContentRange(response.Header.Get("Content-Range")); ok {
			return total
		}
		if response.ContentLength > 0 {
			return startOffset + response.ContentLength
		}
		return 0
	}

	if response.ContentLength > 0 {
		return response.ContentLength
	}
	return 0
}

func contentRangeStartsAt(response *http.Response, offset int64) bool {
	start, _, ok := parseContentRange(response.Header.Get("Content-Range"))
	return !ok || start == offset
}

func parseContentRange(value string) (start int64, total int64, ok bool) {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "bytes ") {
		return 0, 0, false
	}

	parts := strings.SplitN(strings.TrimPrefix(value, "bytes "), "/", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}

	rangeParts := strings.SplitN(parts[0], "-", 2)
	if len(rangeParts) != 2 {
		return 0, 0, false
	}

	start, err := strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		return 0, 0, false
	}

	if parts[1] == "*" {
		return start, 0, true
	}

	total, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, false
	}

	return start, total, true
}
