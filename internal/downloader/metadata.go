package downloader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type partialMetadata struct {
	URL             string    `json:"url"`
	TargetPath      string    `json:"target_path"`
	PartialPath     string    `json:"partial_path"`
	TotalBytes      int64     `json:"total_bytes"`
	DownloadedBytes int64     `json:"downloaded_bytes"`
	ETag            string    `json:"etag,omitempty"`
	LastModified    string    `json:"last_modified,omitempty"`
	AcceptRanges    string    `json:"accept_ranges,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func metadataPathForPartial(partialPath string) string {
	return partialPath + ".daryaft.json"
}

func loadPartialMetadata(path string) (partialMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return partialMetadata{}, err
	}

	var metadata partialMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return partialMetadata{}, fmt.Errorf("parse partial metadata %q: %w", path, err)
	}
	return metadata, nil
}

func savePartialMetadata(path string, metadata partialMetadata) error {
	now := time.Now().UTC()
	if metadata.CreatedAt.IsZero() {
		metadata.CreatedAt = now
	}
	metadata.UpdatedAt = now

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("encode partial metadata %q: %w", path, err)
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create metadata directory %q: %w", filepath.Dir(path), err)
	}

	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0o600); err != nil {
		return fmt.Errorf("write partial metadata %q: %w", path, err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("replace partial metadata %q: %w", path, err)
	}
	return nil
}

func removePartialMetadata(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove partial metadata %q: %w", path, err)
	}
	return nil
}
