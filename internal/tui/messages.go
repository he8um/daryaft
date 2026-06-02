package tui

import (
	"github.com/he8um/daryaft/internal/downloader"
	"github.com/he8um/daryaft/internal/inspect"
)

type executionItemStartedMsg struct {
	Item downloader.BatchItem
}

type executionEventMsg struct {
	Item  downloader.BatchItem
	Event downloader.Event
}

type executionFinishedMsg struct {
	Summary executionSummary
}

type executionClosedMsg struct{}

type inspectFinishedMsg struct {
	Result inspect.Result
	Err    error
}

type inspectClosedMsg struct{}
