package downloader

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/download"
)

func TestDownloadSingleURL(t *testing.T) {
	var userAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.UserAgent()
		w.Header().Set("Content-Disposition", `attachment; filename="server-name.txt"`)
		_, _ = w.Write([]byte("hello daryaft"))
	}))
	defer server.Close()

	dir := t.TempDir()
	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/ignored-path.bin"},
		Output: dir,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	wantPath := filepath.Join(dir, "server-name.txt")
	if result.Path != wantPath {
		t.Fatalf("result.Path = %q, want %q", result.Path, wantPath)
	}

	content, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}
	if string(content) != "hello daryaft" {
		t.Fatalf("downloaded content = %q", content)
	}
	if _, err := os.Stat(wantPath + ".part"); !os.IsNotExist(err) {
		t.Fatalf("partial file still exists or stat failed: %v", err)
	}
	if userAgent != "Daryaft/0.1.0-dev" {
		t.Fatalf("User-Agent = %q", userAgent)
	}
}

func TestDownloadRejectsNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusTeapot)
	}))
	defer server.Close()

	_, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: t.TempDir(),
	})
	if err == nil {
		t.Fatal("Download returned nil error")
	}
	if !strings.Contains(err.Error(), "server returned 418") {
		t.Fatalf("error = %q, want status context", err)
	}
}

func TestDownloadCreatesOutputDirectory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	dir := filepath.Join(t.TempDir(), "nested", "downloads")
	result, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
	})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}

	if result.Path != filepath.Join(dir, "file.txt") {
		t.Fatalf("result.Path = %q", result.Path)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("output directory was not created: %v", err)
	}
}

func TestDownloadRejectsExistingTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("new"))
	}))
	defer server.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(target, []byte("old"), 0o600); err != nil {
		t.Fatalf("write existing target: %v", err)
	}

	_, err := New().Download(download.Plan{
		URLs:   []string{server.URL + "/file.txt"},
		Output: dir,
	})
	if err == nil {
		t.Fatal("Download returned nil error")
	}
	if !strings.Contains(err.Error(), "target file already exists") {
		t.Fatalf("error = %q, want existing target context", err)
	}
}
