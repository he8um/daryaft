package inspect

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Result struct {
	URL                string `json:"url"`
	FinalURL           string `json:"final_url"`
	Status             string `json:"status"`
	StatusCode         int    `json:"status_code"`
	Filename           string `json:"filename"`
	ContentLength      int64  `json:"content_length"`
	ContentLengthKnown bool   `json:"content_length_known"`
	ContentType        string `json:"content_type"`
	AcceptRanges       string `json:"accept_ranges"`
	ResumeSupported    bool   `json:"resume_supported"`
	ResumeSupportKnown bool   `json:"resume_support_known"`
	ETag               string `json:"etag"`
	LastModified       string `json:"last_modified"`
}

func (r Result) Format() string {
	var builder strings.Builder

	fmt.Fprintln(&builder, "Daryaft inspect")
	fmt.Fprintln(&builder)
	fmt.Fprintf(&builder, "URL: %s\n", r.URL)
	fmt.Fprintf(&builder, "Final URL: %s\n", valueOrUnknown(r.FinalURL))
	fmt.Fprintf(&builder, "Status: %s\n", valueOrUnknown(r.Status))
	fmt.Fprintf(&builder, "Filename: %s\n", valueOrUnknown(r.Filename))
	if r.ContentLengthKnown {
		fmt.Fprintf(&builder, "Content length: %d bytes\n", r.ContentLength)
	} else {
		fmt.Fprintln(&builder, "Content length: unknown")
	}
	fmt.Fprintf(&builder, "Content type: %s\n", valueOrUnknown(r.ContentType))
	fmt.Fprintf(&builder, "Accept-Ranges: %s\n", valueOrUnknown(r.AcceptRanges))
	fmt.Fprintf(&builder, "Resume supported: %s\n", formatResumeSupport(r))
	fmt.Fprintf(&builder, "ETag: %s\n", valueOrUnknown(r.ETag))
	fmt.Fprintf(&builder, "Last-Modified: %s\n", valueOrUnknown(r.LastModified))

	return strings.TrimRight(builder.String(), "\n")
}

func (r Result) FormatJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func valueOrUnknown(value string) string {
	if strings.TrimSpace(value) == "" {
		return "unknown"
	}
	return value
}

func formatResumeSupport(result Result) string {
	if !result.ResumeSupportKnown {
		return "unknown"
	}
	if result.ResumeSupported {
		return "yes"
	}
	return "no"
}
