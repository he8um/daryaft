package downloader

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareTargetCreatesCleanTargetPath(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "downloads")

	target, err := prepareTarget(dir, "file.zip")
	if err != nil {
		t.Fatalf("prepareTarget returned error: %v", err)
	}

	if target.Final != filepath.Join(dir, "file.zip") {
		t.Fatalf("target.Final = %q", target.Final)
	}
	if target.Partial != filepath.Join(dir, "file.zip.part") {
		t.Fatalf("target.Partial = %q", target.Partial)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("output directory was not created: %v", err)
	}
}

func TestPrepareTargetRejectsExistingTarget(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.zip")
	if err := os.WriteFile(path, []byte("existing"), 0o600); err != nil {
		t.Fatalf("write existing target: %v", err)
	}

	_, err := prepareTarget(dir, "file.zip")
	if !errors.Is(err, ErrTargetExists) {
		t.Fatalf("prepareTarget error = %v, want ErrTargetExists", err)
	}
}
