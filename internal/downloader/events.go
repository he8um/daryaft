package downloader

import "time"

type EventType string

const (
	EventStarted    EventType = "started"
	EventProgress   EventType = "progress"
	EventCompleted  EventType = "completed"
	EventFailed     EventType = "failed"
	EventRetrying   EventType = "retrying"
	EventResuming   EventType = "resuming"
	EventRestarting EventType = "restarting"
	EventWarning    EventType = "warning"
	EventCancelled  EventType = "cancelled"
)

type Event struct {
	Type                EventType
	URL                 string
	TargetPath          string
	PartialPath         string
	DownloadedBytes     int64
	TotalBytes          int64
	Percent             float64
	SpeedBytesPerSecond float64
	Message             string
	Error               error
	Attempt             int
	MaxAttempts         int
	NextDelay           time.Duration
	Timestamp           time.Time
}

type EventHandler func(Event)

func emitEvent(handler EventHandler, event Event) {
	if handler == nil {
		return
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	handler(event)
}

func newProgressEvent(rawURL, targetPath, partialPath string, downloadedBytes, transferredBytes, totalBytes int64, startedAt time.Time) Event {
	event := Event{
		Type:            EventProgress,
		URL:             rawURL,
		TargetPath:      targetPath,
		PartialPath:     partialPath,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      totalBytes,
		Timestamp:       time.Now(),
	}

	if totalBytes > 0 {
		event.Percent = float64(downloadedBytes) / float64(totalBytes) * 100
	}

	elapsed := event.Timestamp.Sub(startedAt).Seconds()
	if elapsed > 0 {
		event.SpeedBytesPerSecond = float64(transferredBytes) / elapsed
	}

	return event
}
