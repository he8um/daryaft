package downloader

import (
	"errors"
	"fmt"
	"strings"
)

type Result struct {
	URL  string
	Path string
}

type BatchItem struct {
	Index int
	Total int
	URL   string
}

type BatchItemResult struct {
	Item           BatchItem
	Result         Result
	Err            error
	ChecksumStatus string // "", "verified", "failed"
}

type BatchResult struct {
	Planned          int
	ChecksumVerified int
	Items            []BatchItemResult
}

func (r BatchResult) Total() int {
	if r.Planned > 0 {
		return r.Planned
	}
	return len(r.Items)
}

func (r BatchResult) Completed() int {
	completed := 0
	for _, item := range r.Items {
		if item.Err == nil {
			completed++
		}
	}
	return completed
}

func (r BatchResult) Failed() int {
	failed := 0
	for _, item := range r.Items {
		if item.Err != nil && !errors.Is(item.Err, ErrCancelled) {
			failed++
		}
	}
	return failed
}

func (r BatchResult) Cancelled() int {
	cancelled := 0
	for _, item := range r.Items {
		if errors.Is(item.Err, ErrCancelled) {
			cancelled++
		}
	}
	return cancelled
}

func (r BatchResult) Skipped() int {
	skipped := r.Total() - r.Completed() - r.Failed() - r.Cancelled()
	if skipped < 0 {
		return 0
	}
	return skipped
}

func (r BatchResult) FailedItems() []BatchItemResult {
	failures := make([]BatchItemResult, 0)
	for _, item := range r.Items {
		if item.Err != nil && !errors.Is(item.Err, ErrCancelled) {
			failures = append(failures, item)
		}
	}
	return failures
}

func (r BatchResult) Err() error {
	if r.Cancelled() > 0 {
		return ErrCancelled
	}
	if r.Failed() == 0 {
		return nil
	}

	plural := ""
	if r.Failed() != 1 {
		plural = "s"
	}
	return fmt.Errorf("batch download completed with %d failure%s", r.Failed(), plural)
}

func (r BatchResult) SummaryString() string {
	var builder strings.Builder

	fmt.Fprintln(&builder, "Daryaft batch summary")
	fmt.Fprintf(&builder, "Total: %d\n", r.Total())
	fmt.Fprintf(&builder, "Completed: %d\n", r.Completed())
	fmt.Fprintf(&builder, "Failed: %d", r.Failed())
	if r.Cancelled() > 0 || r.Skipped() > 0 {
		fmt.Fprintf(&builder, "\nCancelled: %d", r.Cancelled())
	}
	if r.Skipped() > 0 {
		fmt.Fprintf(&builder, "\nNot started: %d", r.Skipped())
	}
	if r.ChecksumVerified > 0 {
		fmt.Fprintf(&builder, "\nChecksum verified: %d", r.ChecksumVerified)
	}

	failures := r.FailedItems()
	if len(failures) > 0 {
		fmt.Fprintln(&builder)
		fmt.Fprintln(&builder)
		fmt.Fprintln(&builder, "Failed downloads:")
		for _, failure := range failures {
			fmt.Fprintf(&builder, "- %s: %v\n", failure.Item.URL, failure.Err)
		}
		return strings.TrimRight(builder.String(), "\n")
	}

	return builder.String()
}
