package downloader

import (
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

const fallbackFilename = "download.bin"

func FilenameFromResponse(rawURL string, header http.Header, customName string) string {
	if strings.TrimSpace(customName) != "" {
		return sanitizeFilename(customName)
	}

	if filename := filenameFromContentDisposition(header.Get("Content-Disposition")); filename != "" {
		return filename
	}

	if filename := filenameFromURL(rawURL); filename != "" {
		return filename
	}

	return fallbackFilename
}

func filenameFromContentDisposition(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}

	_, params, err := mime.ParseMediaType(value)
	if err != nil {
		return ""
	}

	filename, ok := params["filename"]
	if !ok || strings.TrimSpace(filename) == "" {
		return ""
	}

	return sanitizeFilename(filename)
}

func filenameFromURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	base := path.Base(parsed.EscapedPath())
	if base == "." || base == "/" {
		return ""
	}

	unescaped, err := url.PathUnescape(base)
	if err == nil {
		base = unescaped
	}

	return sanitizeFilename(base)
}

func sanitizeFilename(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || trimmed == "." || trimmed == ".." {
		return fallbackFilename
	}
	if filepath.IsAbs(trimmed) || strings.ContainsAny(trimmed, `/\`) || strings.Contains(trimmed, "\x00") {
		return fallbackFilename
	}
	return trimmed
}
