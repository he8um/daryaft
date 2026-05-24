package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDownloadDryRunDoesNotDownload(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte("unexpected"))
	}))
	defer server.Close()

	var output bytes.Buffer
	cmd := downloadCmd
	cmd.SetOut(&output)
	outputDir := filepath.Join(t.TempDir(), "downloads")

	err := runDownload(cmd, []string{server.URL + "/a.txt", server.URL + "/b.txt"}, downloadFlagValues{
		output:  outputDir,
		dryRun:  true,
		retries: 3,
		resume:  true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
	if !strings.Contains(output.String(), "Mode: dry-run only, no network request performed") {
		t.Fatalf("output missing dry-run marker:\n%s", output.String())
	}
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		t.Fatalf("dry-run touched output directory, stat err = %v", err)
	}
}

func TestRunDownloadBatchFromFileAndArgs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/arg.txt":
			_, _ = w.Write([]byte("arg"))
		case "/file.txt":
			_, _ = w.Write([]byte("file"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	urlFile := filepath.Join(dir, "urls.txt")
	if err := os.WriteFile(urlFile, []byte(server.URL+"/file.txt\n"), 0o600); err != nil {
		t.Fatalf("write URL file: %v", err)
	}

	var output bytes.Buffer
	cmd := downloadCmd
	cmd.SetOut(&output)

	err := runDownload(cmd, []string{server.URL + "/arg.txt"}, downloadFlagValues{
		file:    urlFile,
		output:  dir,
		retries: 3,
		resume:  true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}

	for _, want := range []string{
		"[1/2] Downloading: " + server.URL + "/arg.txt",
		"[2/2] Downloading: " + server.URL + "/file.txt",
		"Daryaft batch summary",
		"Total: 2",
		"Completed: 2",
		"Failed: 0",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("output missing %q in:\n%s", want, output.String())
		}
	}
}

func TestRunDownloadRejectsNameWithMultipleURLs(t *testing.T) {
	cmd := downloadCmd

	err := runDownload(cmd, []string{"https://example.com/a.txt", "https://example.com/b.txt"}, downloadFlagValues{
		name:    "file.txt",
		retries: 3,
		resume:  true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if !strings.Contains(err.Error(), "--name can only be used with a single URL") {
		t.Fatalf("error = %q", err)
	}
}
