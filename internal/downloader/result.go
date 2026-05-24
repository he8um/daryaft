package downloader

import (
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
	Item   BatchItem
	Result Result
	Err    error
}

type BatchResult struct {
	Items []BatchItemResult
}

func (r BatchResult) Total() int {
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
	return r.Total() - r.Completed()
}

func (r BatchResult) FailedItems() []BatchItemResult {
	failures := make([]BatchItemResult, 0)
	for _, item := range r.Items {
		if item.Err != nil {
			failures = append(failures, item)
		}
	}
	return failures
}

func (r BatchResult) Err() error {
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
