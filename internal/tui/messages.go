package tui

import (
	"github.com/he8um/daryaft/internal/downloader"
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
