package downloader

import (
	"net/http"
	"testing"
)

func TestFilenameFromURL(t *testing.T) {
	got := FilenameFromResponse("https://example.com/releases/file%20name.zip", http.Header{}, "")
	if got != "file name.zip" {
		t.Fatalf("FilenameFromResponse() = %q, want %q", got, "file name.zip")
	}
}

func TestFilenameFromContentDisposition(t *testing.T) {
	header := http.Header{}
	header.Set("Content-Disposition", `attachment; filename="release.tar.gz"`)

	got := FilenameFromResponse("https://example.com/download", header, "")
	if got != "release.tar.gz" {
		t.Fatalf("FilenameFromResponse() = %q, want %q", got, "release.tar.gz")
	}
}

func TestFilenameSanitization(t *testing.T) {
	tests := []string{
		"",
		"   ",
		".",
		"..",
		"../bad",
		"..\\bad",
		"/tmp/bad",
		"dir/file",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			if got := sanitizeFilename(name); got != fallbackFilename {
				t.Fatalf("sanitizeFilename(%q) = %q, want %q", name, got, fallbackFilename)
			}
		})
	}
}

func TestCustomFilenameUsesFallbackWhenUnsafe(t *testing.T) {
	got := FilenameFromResponse("https://example.com/file.zip", http.Header{}, "../bad")
	if got != fallbackFilename {
		t.Fatalf("FilenameFromResponse() = %q, want %q", got, fallbackFilename)
	}
}
