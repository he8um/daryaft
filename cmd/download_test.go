package cmd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/downloader"
	"github.com/he8um/daryaft/internal/httpopts"
	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "daryaft-cmd-test-*")
	if err != nil {
		panic(err)
	}
	clearDaryaftEnv()
	restore := appconfig.SetUserConfigDirForTest(dir)
	code := m.Run()
	restore()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

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

func TestRunDownloadDryRunIncludesChecksum(t *testing.T) {
	sum := sha256.Sum256([]byte("hello"))
	expected := fmt.Sprintf("%x", sum)

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)

	err := runDownload(cmd, []string{"https://example.com/file.txt"}, downloadFlagValues{
		dryRun:   true,
		checksum: "sha256:" + expected,
		retries:  3,
		resume:   true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Checksum: sha256:"+expected) {
		t.Fatalf("output missing checksum:\n%s", output.String())
	}
}

func TestRunDownloadInvalidChecksumFailsBeforeNetworkCall(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte("unexpected"))
	}))
	defer server.Close()

	cmd := &cobra.Command{Use: "download"}
	err := runDownload(cmd, []string{server.URL + "/file.txt"}, downloadFlagValues{
		checksum: "sha256:not-hex",
		retries:  3,
		resume:   true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
	if !strings.Contains(err.Error(), "64 hex characters") {
		t.Fatalf("error = %q", err)
	}
}

func TestRunDownloadUsesDownloadsDefaultWhenOutputNotProvided(t *testing.T) {
	restoreConfig := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restoreConfig)
	home := t.TempDir()
	t.Cleanup(appconfig.SetUserHomeDirForTest(home))
	wantOutput := filepath.Join(home, "Downloads")

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Output: "+wantOutput) {
		t.Fatalf("output missing Downloads default %q:\n%s", wantOutput, output.String())
	}
}

func TestRunDownloadOutputDotOverridesDownloadsDefault(t *testing.T) {
	restoreConfig := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restoreConfig)
	t.Cleanup(appconfig.SetUserHomeDirForTest(t.TempDir()))

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}
	if err := cmd.Flags().Set("output", "."); err != nil {
		t.Fatalf("set output: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Output: .") {
		t.Fatalf("output missing explicit current directory:\n%s", output.String())
	}
}

func TestRunDownloadEmptyConfigDownloadDirUsesDownloadsDefault(t *testing.T) {
	restoreConfig := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restoreConfig)
	home := t.TempDir()
	t.Cleanup(appconfig.SetUserHomeDirForTest(home))
	wantOutput := filepath.Join(home, "Downloads")

	cfg := appconfig.Default()
	cfg.DownloadDir = ""
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Output: "+wantOutput) {
		t.Fatalf("output missing Downloads default %q:\n%s", wantOutput, output.String())
	}
}

func TestRunDownloadConfigDotOverridesDownloadsDefault(t *testing.T) {
	restoreConfig := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restoreConfig)
	t.Cleanup(appconfig.SetUserHomeDirForTest(t.TempDir()))

	cfg := appconfig.Default()
	cfg.DownloadDir = "."
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Output: .") {
		t.Fatalf("output missing config current directory:\n%s", output.String())
	}
}

func TestRunDownloadUsesConfigDownloadDirWhenOutputNotProvided(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	cfg := appconfig.Default()
	cfg.DownloadDir = filepath.Join(t.TempDir(), "configured-downloads")
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Output: "+cfg.DownloadDir) {
		t.Fatalf("output missing configured download dir:\n%s", output.String())
	}
}

func TestRunDownloadUsesEnvDownloadDirBeforeConfig(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)
	t.Setenv("DARYAFT_DOWNLOAD_DIR", "/tmp/env-daryaft")

	cfg := appconfig.Default()
	cfg.DownloadDir = "/tmp/config-daryaft"
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Output: /tmp/env-daryaft") {
		t.Fatalf("output missing env download dir:\n%s", output.String())
	}
}

func TestRunDownloadUsesEnvRetriesBeforeConfig(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)
	t.Setenv("DARYAFT_RETRIES", "6")

	cfg := appconfig.Default()
	cfg.Retries = 2
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Retries: 6") {
		t.Fatalf("output missing env retries:\n%s", output.String())
	}
}

func TestRunDownloadCLIFlagsOverrideConfigDefaults(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)
	t.Setenv("DARYAFT_DOWNLOAD_DIR", "/tmp/env-daryaft")
	t.Setenv("DARYAFT_RETRIES", "7")
	t.Setenv("DARYAFT_RESUME", "false")

	cfg := appconfig.Default()
	cfg.DownloadDir = filepath.Join(t.TempDir(), "configured-downloads")
	cfg.Retries = 7
	cfg.Resume = false
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}
	if err := cmd.Flags().Set("output", "/tmp/explicit-daryaft"); err != nil {
		t.Fatalf("set output: %v", err)
	}
	if err := cmd.Flags().Set("retries", "2"); err != nil {
		t.Fatalf("set retries: %v", err)
	}
	if err := cmd.Flags().Set("resume", "true"); err != nil {
		t.Fatalf("set resume: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	for _, want := range []string{
		"Output: /tmp/explicit-daryaft",
		"Retries: 2",
		"Resume: true",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("output missing %q in:\n%s", want, output.String())
		}
	}
}

func TestRunDownloadNoResumeFlagOverridesEnvResume(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)
	t.Setenv("DARYAFT_RESUME", "true")

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}
	if err := cmd.Flags().Set("no-resume", "true"); err != nil {
		t.Fatalf("set no-resume: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Resume: false") {
		t.Fatalf("output missing no-resume override:\n%s", output.String())
	}
}

func TestRunDownloadUsesConfigRetryAndResumeDefaults(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	cfg := appconfig.Default()
	cfg.Retries = 0
	cfg.Resume = false
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatalf("set dry-run: %v", err)
	}

	var output bytes.Buffer
	cmd.SetOut(&output)
	err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	for _, want := range []string{
		"Retries: 0",
		"Resume: false",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("output missing %q in:\n%s", want, output.String())
		}
	}
}

func TestRunDownloadRejectsRetriesOutOfRange(t *testing.T) {
	for _, value := range []string{"-1", "21"} {
		t.Run(value, func(t *testing.T) {
			restore := appconfig.SetUserConfigDirForTest(t.TempDir())
			t.Cleanup(restore)

			var flags downloadFlagValues
			cmd := &cobra.Command{Use: "download"}
			addDownloadFlags(cmd, &flags)
			if err := cmd.Flags().Set("dry-run", "true"); err != nil {
				t.Fatalf("set dry-run: %v", err)
			}
			if err := cmd.Flags().Set("retries", value); err != nil {
				t.Fatalf("set retries: %v", err)
			}

			err := runDownload(cmd, []string{"https://example.com/file.zip"}, flags)
			if err == nil {
				t.Fatal("runDownload returned nil error")
			}
			if !strings.Contains(err.Error(), "retries") {
				t.Fatalf("error = %q, want retries context", err)
			}
		})
	}
}

func clearDaryaftEnv() {
	for _, name := range []string{
		"DARYAFT_DOWNLOAD_DIR",
		"DARYAFT_RETRIES",
		"DARYAFT_RESUME",
		"DARYAFT_NO_COLOR",
		"DARYAFT_NO_TUI",
		"DARYAFT_THEME",
		"DARYAFT_ANIMATIONS",
		"DARYAFT_HYPERLINKS",
	} {
		_ = os.Unsetenv(name)
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

func TestRunDownloadSingleURLVerifiesMatchingChecksum(t *testing.T) {
	content := []byte("verified content")
	sum := sha256.Sum256(content)
	expected := fmt.Sprintf("%x", sum)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(content)
	}))
	defer server.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)
	dir := t.TempDir()

	err := runDownload(cmd, []string{server.URL + "/file.txt"}, downloadFlagValues{
		output:   dir,
		checksum: "sha256:" + expected,
		retries:  3,
		resume:   true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Checksum verified: sha256") {
		t.Fatalf("output missing checksum success:\n%s", output.String())
	}
	if got, err := os.ReadFile(filepath.Join(dir, "file.txt")); err != nil || string(got) != string(content) {
		t.Fatalf("downloaded file = %q, err = %v", got, err)
	}
}

func TestRunDownloadSingleURLVerifiesMatchingSHA512Checksum(t *testing.T) {
	content := []byte("verified sha512 content")
	sum := sha512.Sum512(content)
	expected := fmt.Sprintf("%x", sum)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(content)
	}))
	defer server.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)

	err := runDownload(cmd, []string{server.URL + "/file.txt"}, downloadFlagValues{
		output:   t.TempDir(),
		checksum: "sha512:" + expected,
		retries:  3,
		resume:   true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Checksum verified: sha512") {
		t.Fatalf("output missing sha512 checksum success:\n%s", output.String())
	}
}

func TestRunDownloadDoesNotVerifyChecksumAfterDownloadFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failure", http.StatusInternalServerError)
	}))
	defer server.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)
	dir := t.TempDir()

	err := runDownload(cmd, []string{server.URL + "/file.txt"}, downloadFlagValues{
		output:   dir,
		checksum: "sha256:" + strings.Repeat("a", 64),
		retries:  0,
		resume:   true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if strings.Contains(err.Error(), "checksum") {
		t.Fatalf("error = %q, want download failure before checksum", err)
	}
	if strings.Contains(output.String(), "Checksum verified") {
		t.Fatalf("output unexpectedly contains checksum success:\n%s", output.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "file.txt")); !os.IsNotExist(err) {
		t.Fatalf("final file exists or stat failed: %v", err)
	}
}

func TestRunDownloadSingleURLReturnsErrorOnChecksumMismatch(t *testing.T) {
	content := []byte("actual content")
	wrong := sha256.Sum256([]byte("wrong content"))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(content)
	}))
	defer server.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)
	dir := t.TempDir()

	err := runDownload(cmd, []string{server.URL + "/file.txt"}, downloadFlagValues{
		output:   dir,
		checksum: "sha256:" + fmt.Sprintf("%x", wrong),
		retries:  3,
		resume:   true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if !strings.Contains(err.Error(), "checksum mismatch: expected") {
		t.Fatalf("error = %q", err)
	}
	if !strings.Contains(err.Error(), "got") {
		t.Fatalf("error missing actual checksum: %q", err)
	}
	if strings.Contains(output.String(), "Checksum verified") {
		t.Fatalf("output unexpectedly contains checksum success:\n%s", output.String())
	}
	final := filepath.Join(dir, "file.txt")
	if got, err := os.ReadFile(final); err != nil || string(got) != string(content) {
		t.Fatalf("final file after mismatch = %q, err = %v", got, err)
	}
	if _, err := os.Stat(final + ".part"); !os.IsNotExist(err) {
		t.Fatalf("partial exists after completed mismatch or stat failed: %v", err)
	}
}

func TestRunDownloadRejectsChecksumWithBatchInput(t *testing.T) {
	sum := sha256.Sum256([]byte("hello"))
	cmd := &cobra.Command{Use: "download"}

	err := runDownload(cmd, []string{"https://example.com/a.txt", "https://example.com/b.txt"}, downloadFlagValues{
		checksum: "sha256:" + fmt.Sprintf("%x", sum),
		retries:  3,
		resume:   true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if !strings.Contains(err.Error(), "--checksum is currently supported only for single URL downloads") {
		t.Fatalf("error = %q", err)
	}
}

func TestRunDownloadRejectsChecksumWithFileInput(t *testing.T) {
	dir := t.TempDir()
	urlFile := filepath.Join(dir, "urls.txt")
	if err := os.WriteFile(urlFile, []byte("https://example.com/file.txt\n"), 0o600); err != nil {
		t.Fatalf("write URL file: %v", err)
	}
	sum := sha256.Sum256([]byte("hello"))
	cmd := &cobra.Command{Use: "download"}

	err := runDownload(cmd, nil, downloadFlagValues{
		file:     urlFile,
		checksum: "sha256:" + fmt.Sprintf("%x", sum),
		retries:  3,
		resume:   true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if !strings.Contains(err.Error(), "--checksum is currently supported only for single URL downloads") {
		t.Fatalf("error = %q", err)
	}
}

func TestRunDownloadContextCancellationLeavesPartialAndMetadata(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "131072")
		for range 4 {
			_, _ = w.Write(bytes.Repeat([]byte("a"), 32*1024))
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
			select {
			case <-r.Context().Done():
				return
			case <-time.After(20 * time.Millisecond):
			}
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	output := &cancelOnWriteBuffer{needle: "Saving to:", cancel: cancel}
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(output)

	err := runDownloadWithContext(cmd, []string{server.URL + "/file.bin"}, downloadFlagValues{
		output:   dir,
		checksum: "sha256:" + strings.Repeat("a", 64),
		retries:  0,
		resume:   true,
	}, ctx)
	if !errors.Is(err, downloader.ErrCancelled) {
		t.Fatalf("runDownloadWithContext error = %v, want ErrCancelled", err)
	}
	if !strings.Contains(output.String(), "Download cancelled. Partial file kept for resume.") {
		t.Fatalf("output missing cancellation message:\n%s", output.String())
	}
	if strings.Contains(output.String(), "Checksum verified") {
		t.Fatalf("output unexpectedly contains checksum success:\n%s", output.String())
	}

	final := filepath.Join(dir, "file.bin")
	if _, err := os.Stat(final); !os.IsNotExist(err) {
		t.Fatalf("final exists or stat failed: %v", err)
	}
	if _, err := os.Stat(final + ".part"); err != nil {
		t.Fatalf("partial missing after cancellation: %v", err)
	}
	if _, err := os.Stat(final + ".part.daryaft.json"); err != nil {
		t.Fatalf("metadata missing after cancellation: %v", err)
	}
}

func TestRunDownloadBatchContextCancellationStopsRemainingItems(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	requestsByPath := make(map[string]int)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestsByPath[r.URL.Path]++
		switch r.URL.Path {
		case "/one.bin":
			w.Header().Set("Content-Length", "131072")
			for range 4 {
				_, _ = w.Write(bytes.Repeat([]byte("a"), 32*1024))
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
				select {
				case <-r.Context().Done():
					return
				case <-time.After(20 * time.Millisecond):
				}
			}
		case "/two.bin":
			_, _ = w.Write([]byte("should not start"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	output := &cancelOnWriteBuffer{needle: "Saving to:", cancel: cancel}
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(output)

	err := runDownloadWithContext(cmd, []string{server.URL + "/one.bin", server.URL + "/two.bin"}, downloadFlagValues{
		output:  dir,
		retries: 0,
		resume:  true,
	}, ctx)
	if !errors.Is(err, downloader.ErrCancelled) {
		t.Fatalf("runDownloadWithContext error = %v, want ErrCancelled", err)
	}
	if requestsByPath["/two.bin"] != 0 {
		t.Fatalf("second item requests = %d, want 0", requestsByPath["/two.bin"])
	}

	for _, want := range []string{
		"Download cancelled. Partial file kept for resume.",
		"Daryaft batch summary",
		"Total: 2",
		"Cancelled: 1",
		"Not started: 1",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("output missing %q in:\n%s", want, output.String())
		}
	}

	final := filepath.Join(dir, "one.bin")
	if _, err := os.Stat(final); !os.IsNotExist(err) {
		t.Fatalf("final exists or stat failed: %v", err)
	}
	if _, err := os.Stat(final + ".part"); err != nil {
		t.Fatalf("partial missing after cancellation: %v", err)
	}
	if _, err := os.Stat(final + ".part.daryaft.json"); err != nil {
		t.Fatalf("metadata missing after cancellation: %v", err)
	}
}

type cancelOnWriteBuffer struct {
	bytes.Buffer
	needle    string
	cancel    context.CancelFunc
	cancelled bool
}

func (b *cancelOnWriteBuffer) Write(p []byte) (int, error) {
	n, err := b.Buffer.Write(p)
	if !b.cancelled && strings.Contains(b.String(), b.needle) {
		b.cancelled = true
		b.cancel()
	}
	return n, err
}

func TestRunDownloadVerboseSingleURLAddsDiagnostics(t *testing.T) {
	verbose = true
	t.Cleanup(func() { verbose = false })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "2")
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)
	dir := t.TempDir()

	err := runDownload(cmd, []string{server.URL + "/file.txt"}, downloadFlagValues{
		output:  dir,
		retries: 3,
		resume:  true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}

	for _, want := range []string{
		"Verbose: output directory: " + dir,
		"Verbose: effective URL: " + server.URL + "/file.txt",
		"Verbose: HTTP status: 200 OK",
		"Verbose: target path: " + filepath.Join(dir, "file.txt"),
		"Verbose: selected filename: file.txt",
		"Verbose: completion duration:",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("verbose output missing %q in:\n%s", want, output.String())
		}
	}
	for _, notWant := range []string{"Authorization", "Cookie"} {
		if strings.Contains(output.String(), notWant) {
			t.Fatalf("verbose output contains sensitive header name %q:\n%s", notWant, output.String())
		}
	}
}

func TestRunDownloadNormalOutputDoesNotContainVerboseDiagnostics(t *testing.T) {
	verbose = false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)

	err := runDownload(cmd, []string{server.URL + "/file.txt"}, downloadFlagValues{
		output:  t.TempDir(),
		retries: 3,
		resume:  true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if strings.Contains(output.String(), "Verbose:") {
		t.Fatalf("normal output contains verbose diagnostics:\n%s", output.String())
	}
}

func TestRunDownloadVerboseBatchAddsItemDiagnostics(t *testing.T) {
	verbose = true
	t.Cleanup(func() { verbose = false })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a.txt", "/b.txt":
			_, _ = w.Write([]byte("ok"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)
	dir := t.TempDir()

	err := runDownload(cmd, []string{server.URL + "/a.txt", server.URL + "/b.txt"}, downloadFlagValues{
		output:  dir,
		retries: 3,
		resume:  true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}

	for _, want := range []string{
		"Verbose: output directory: " + dir,
		"Verbose: item 1/2 effective URL: " + server.URL + "/a.txt",
		"Verbose: item 2/2 effective URL: " + server.URL + "/b.txt",
		"Verbose: HTTP status: 200 OK",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("verbose batch output missing %q in:\n%s", want, output.String())
		}
	}
}

func TestVerboseURLDiagnosticsRedactUserInfoAndQuery(t *testing.T) {
	got := redactURL("https://user:pass@example.com/file.zip?token=secret#section")
	if got != "https://example.com/file.zip" {
		t.Fatalf("redactURL = %q", got)
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

// --- HTTP options CLI tests ---

func TestDownloadCLICustomHeader(t *testing.T) {
	var gotHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-CLI")
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", server.URL + "/file.txt", "--header", "X-CLI: clivalue", "--output", t.TempDir()})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	if gotHeader != "clivalue" {
		t.Errorf("X-CLI = %q, want clivalue", gotHeader)
	}
}

func TestDownloadCLIUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.UserAgent()
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", server.URL + "/file.txt", "--user-agent", "CLIAgent/4.0", "--output", t.TempDir()})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	if gotUA != "CLIAgent/4.0" {
		t.Errorf("User-Agent = %q, want CLIAgent/4.0", gotUA)
	}
}

func TestDownloadCLIBasicAuth(t *testing.T) {
	var gotUser, gotPass string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, _ = r.BasicAuth()
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", server.URL + "/file.txt", "--username", "alice", "--password", "secret", "--output", t.TempDir()})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	if gotUser != "alice" || gotPass != "secret" {
		t.Errorf("BasicAuth = %q/%q, want alice/secret", gotUser, gotPass)
	}
}

func TestDownloadCLIDryRunRedactsPassword(t *testing.T) {
	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{
		"download",
		"https://example.com/file.zip",
		"--username", "alice",
		"--password", "topsecret",
		"--dry-run",
	})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	output := out.String()
	if strings.Contains(output, "topsecret") {
		t.Errorf("dry-run output must not contain raw password:\n%s", output)
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Errorf("dry-run output must contain [REDACTED]:\n%s", output)
	}
}

func TestDownloadCLIDryRunRedactsSensitiveHeader(t *testing.T) {
	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{
		"download",
		"https://example.com/file.zip",
		"--header", "Authorization: Bearer supersecret",
		"--dry-run",
	})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	output := out.String()
	if strings.Contains(output, "supersecret") {
		t.Errorf("dry-run output must not contain raw Authorization value:\n%s", output)
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Errorf("dry-run output must contain [REDACTED]:\n%s", output)
	}
}

func TestDownloadCLIRejectsInvalidHeader(t *testing.T) {
	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", "https://example.com/file.zip", "--header", "NoColon", "--dry-run"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for invalid header")
	}
}

func TestDownloadCLIRejectsInvalidProxy(t *testing.T) {
	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", "https://example.com/file.zip", "--proxy", "socks5://proxy:1080", "--dry-run"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for socks5 proxy")
	}
}

func TestDownloadCLIRejectsPasswordWithoutUsername(t *testing.T) {
	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", "https://example.com/file.zip", "--password", "secret", "--dry-run"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error: password without username")
	}
}

func TestDownloadCLIRejectsAuthHeaderPlusBasicAuth(t *testing.T) {
	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{
		"download",
		"https://example.com/file.zip",
		"--username", "alice",
		"--password", "pass",
		"--header", "Authorization: Bearer tok",
		"--dry-run",
	})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error: basic auth + Authorization header")
	}
}

func TestDownloadCLIBatchSendsHeaderToAllURLs(t *testing.T) {
	counts := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counts[r.Header.Get("X-Batch")]++
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{
		"download",
		server.URL + "/a.txt",
		server.URL + "/b.txt",
		"--header", "X-Batch: yes",
		"--output", t.TempDir(),
	})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	if counts["yes"] != 2 {
		t.Errorf("X-Batch header received %d times, want 2", counts["yes"])
	}
}

func TestDownloadCLIEnvCredentials(t *testing.T) {
	var gotUser, gotPass string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, _ = r.BasicAuth()
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	t.Setenv("DARYAFT_USERNAME", "envuser")
	t.Setenv("DARYAFT_PASSWORD", "envpass")

	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", server.URL + "/file.txt", "--output", t.TempDir()})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	if gotUser != "envuser" || gotPass != "envpass" {
		t.Errorf("BasicAuth = %q/%q, want envuser/envpass", gotUser, gotPass)
	}
}

func TestDownloadCLIFlagCredentialOverridesEnv(t *testing.T) {
	var gotUser, gotPass string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, _ = r.BasicAuth()
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	t.Setenv("DARYAFT_USERNAME", "envuser")
	t.Setenv("DARYAFT_PASSWORD", "envpass")

	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", server.URL + "/file.txt", "--username", "flaguser", "--password", "flagpass", "--output", t.TempDir()})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	if gotUser != "flaguser" || gotPass != "flagpass" {
		t.Errorf("BasicAuth = %q/%q, want flaguser/flagpass (flag should override env)", gotUser, gotPass)
	}
}

func TestDownloadCLIEnvPasswordWithoutUsernameIsRejected(t *testing.T) {
	t.Setenv("DARYAFT_PASSWORD", "envpass")

	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", "https://example.com/file.zip", "--dry-run"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error: DARYAFT_PASSWORD without username")
	}
}

func TestDownloadCLIEnvCredentialsDryRunRedacted(t *testing.T) {
	t.Setenv("DARYAFT_USERNAME", "envuser")
	t.Setenv("DARYAFT_PASSWORD", "envpass")

	var out bytes.Buffer
	root := newDownloadCmdForTest(&out)
	root.SetArgs([]string{"download", "https://example.com/file.zip", "--dry-run"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\n%s", err, out.String())
	}
	output := out.String()
	if strings.Contains(output, "envpass") {
		t.Errorf("dry-run output must not contain raw DARYAFT_PASSWORD value:\n%s", output)
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Errorf("dry-run output must contain [REDACTED]:\n%s", output)
	}
}

func TestSingleDownloadCompleted_PrintsSizeAndElapsed(t *testing.T) {
	body := make([]byte, 512)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "512")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer ts.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)

	err := runDownload(cmd, []string{ts.URL + "/file.bin"}, downloadFlagValues{
		output:  t.TempDir(),
		retries: 0,
		resume:  true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}

	matched, _ := regexp.MatchString(`Completed: .+\(.+ in .+\)`, output.String())
	if !matched {
		t.Errorf("expected Completed with size and elapsed, got: %s", output.String())
	}
}

func TestSingleDownloadEventFailed_PrintsFailedMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)

	err := runDownload(cmd, []string{ts.URL + "/file.bin"}, downloadFlagValues{
		output:  t.TempDir(),
		retries: 0,
		resume:  true,
	})
	// The download should fail (server returns 500).
	if err == nil {
		t.Fatal("runDownload returned nil error, want failure")
	}
	// The event-stream Failed: prefix or the returned error must surface the failure.
	hasFailedPrefix := strings.Contains(output.String(), "Failed:")
	hasError := err != nil
	if !hasFailedPrefix && !hasError {
		t.Errorf("expected Failed: in output or non-nil error; output: %s", output.String())
	}
}

// newDownloadCmdForTest builds a fresh root+download command pair for CLI flag tests.
// Returns the root command; callers should prepend "download" to args when needed,
// or use the download subcommand args directly via root.SetArgs([]string{"download", ...}).
// For convenience, returns the root so SetArgs("download", url, ...) routes correctly.
func newDownloadCmdForTest(out *bytes.Buffer) *cobra.Command {
	_ = httpopts.Options{} // ensure import is used

	flags := downloadFlagValues{}
	flags.resume = true

	sub := &cobra.Command{
		Use:           "download [url...]",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDownload(cmd, args, flags)
		},
	}
	addDownloadFlags(sub, &flags)

	root := &cobra.Command{Use: "daryaft", SilenceUsage: true, SilenceErrors: true}
	root.AddCommand(sub)
	root.SetOut(out)
	root.SetErr(out)
	return root
}

func writeCmdChecksumFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "checksums.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write checksum file: %v", err)
	}
	return path
}

func TestRunDownloadChecksumFileHelp(t *testing.T) {
	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "download"}
	addDownloadFlags(cmd, &flags)
	if cmd.Flags().Lookup("checksum-file") == nil {
		t.Fatal("--checksum-file flag not registered")
	}
	usage := cmd.Flags().FlagUsages()
	if !strings.Contains(usage, "checksum-file") {
		t.Fatalf("flag usage missing checksum-file:\n%s", usage)
	}
}

func TestRunDownloadChecksumFileMatch(t *testing.T) {
	contentA := []byte("alpha")
	contentB := []byte("beta")
	sumA := sha256.Sum256(contentA)
	sumB := sha256.Sum256(contentB)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a.txt":
			_, _ = w.Write(contentA)
		case "/b.txt":
			_, _ = w.Write(contentB)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	urlB := server.URL + "/b.txt"
	manifest := writeCmdChecksumFile(t,
		"sha256:"+fmt.Sprintf("%x", sumA)+" "+urlA+"\n"+
			"sha256:"+fmt.Sprintf("%x", sumB)+" "+urlB+"\n")

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)
	dir := t.TempDir()

	err := runDownload(cmd, []string{urlA, urlB}, downloadFlagValues{
		output:       dir,
		checksumFile: manifest,
		retries:      3,
		resume:       true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Checksum verified: 2") {
		t.Fatalf("output missing checksum verified count:\n%s", output.String())
	}
}

func TestRunDownloadChecksumFileMismatch(t *testing.T) {
	contentA := []byte("alpha")
	contentB := []byte("beta")
	sumB := sha256.Sum256(contentB)
	wrongA := sha256.Sum256([]byte("WRONG"))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a.txt":
			_, _ = w.Write(contentA)
		case "/b.txt":
			_, _ = w.Write(contentB)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	urlB := server.URL + "/b.txt"
	manifest := writeCmdChecksumFile(t,
		"sha256:"+fmt.Sprintf("%x", wrongA)+" "+urlA+"\n"+
			"sha256:"+fmt.Sprintf("%x", sumB)+" "+urlB+"\n")

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)
	dir := t.TempDir()

	err := runDownload(cmd, []string{urlA, urlB}, downloadFlagValues{
		output:       dir,
		checksumFile: manifest,
		retries:      3,
		resume:       true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error on mismatch")
	}
	out := output.String()
	if !strings.Contains(out, "checksum mismatch: expected") {
		t.Fatalf("summary missing checksum mismatch detail:\n%s", out)
	}
	if !strings.Contains(out, "got") {
		t.Fatalf("summary missing actual digest:\n%s", out)
	}
	// Downloaded file is left in place.
	if got, err := os.ReadFile(filepath.Join(dir, "a.txt")); err != nil || string(got) != string(contentA) {
		t.Fatalf("a.txt after mismatch = %q, err = %v", got, err)
	}
}

func TestRunDownloadChecksumFileDryRun(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte("unexpected"))
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	manifest := writeCmdChecksumFile(t, "sha256:"+strings.Repeat("a", 64)+" "+urlA+"\n")

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)

	err := runDownload(cmd, []string{urlA}, downloadFlagValues{
		output:       filepath.Join(t.TempDir(), "downloads"),
		dryRun:       true,
		checksumFile: manifest,
		retries:      3,
		resume:       true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0", requests)
	}
	if !strings.Contains(output.String(), "Checksums: from file (1 entries)") {
		t.Fatalf("dry-run output missing checksum file line:\n%s", output.String())
	}
}

func TestRunDownloadChecksumFileMissingURL(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte("unexpected"))
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	urlB := server.URL + "/b.txt"
	manifest := writeCmdChecksumFile(t, "sha256:"+strings.Repeat("a", 64)+" "+urlA+"\n")

	cmd := &cobra.Command{Use: "download"}
	err := runDownload(cmd, []string{urlA, urlB}, downloadFlagValues{
		output:       t.TempDir(),
		checksumFile: manifest,
		retries:      3,
		resume:       true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0 (validation before network)", requests)
	}
	if !strings.Contains(err.Error(), "no checksum provided for URL") {
		t.Fatalf("error = %q", err)
	}
}

func TestRunDownloadChecksumFileExtraURL(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte("unexpected"))
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	manifest := writeCmdChecksumFile(t,
		"sha256:"+strings.Repeat("a", 64)+" "+urlA+"\n"+
			"sha256:"+strings.Repeat("b", 64)+" "+server.URL+"/extra.txt\n")

	cmd := &cobra.Command{Use: "download"}
	err := runDownload(cmd, []string{urlA}, downloadFlagValues{
		output:       t.TempDir(),
		checksumFile: manifest,
		retries:      3,
		resume:       true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if requests != 0 {
		t.Fatalf("requests = %d, want 0 (validation before network)", requests)
	}
	if !strings.Contains(err.Error(), "manifest URL not in download targets") {
		t.Fatalf("error = %q", err)
	}
}

func TestRunDownloadChecksumAndChecksumFileTogether(t *testing.T) {
	manifest := writeCmdChecksumFile(t, "sha256:"+strings.Repeat("a", 64)+" https://example.com/a.txt\n")
	cmd := &cobra.Command{Use: "download"}

	err := runDownload(cmd, []string{"https://example.com/a.txt"}, downloadFlagValues{
		output:       t.TempDir(),
		checksum:     "sha256:" + strings.Repeat("a", 64),
		checksumFile: manifest,
		retries:      3,
		resume:       true,
	})
	if err == nil {
		t.Fatal("runDownload returned nil error")
	}
	if !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("error = %q", err)
	}
}

func TestRunDownloadChecksumFileSingleURL(t *testing.T) {
	content := []byte("alpha")
	sum := sha256.Sum256(content)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(content)
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	manifest := writeCmdChecksumFile(t, "sha256:"+fmt.Sprintf("%x", sum)+" "+urlA+"\n")

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)
	dir := t.TempDir()

	err := runDownload(cmd, []string{urlA}, downloadFlagValues{
		output:       dir,
		checksumFile: manifest,
		retries:      3,
		resume:       true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Checksum verified: 1") {
		t.Fatalf("output missing checksum verified count:\n%s", output.String())
	}
}

func TestRunDownloadChecksumFileWithURLFile(t *testing.T) {
	contentA := []byte("alpha")
	contentB := []byte("beta")
	sumA := sha256.Sum256(contentA)
	sumB := sha256.Sum256(contentB)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a.txt":
			_, _ = w.Write(contentA)
		case "/b.txt":
			_, _ = w.Write(contentB)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	urlB := server.URL + "/b.txt"
	urlFile := filepath.Join(t.TempDir(), "urls.txt")
	if err := os.WriteFile(urlFile, []byte(urlA+"\n"+urlB+"\n"), 0o600); err != nil {
		t.Fatalf("write URL file: %v", err)
	}
	manifest := writeCmdChecksumFile(t,
		"sha256:"+fmt.Sprintf("%x", sumA)+" "+urlA+"\n"+
			"sha256:"+fmt.Sprintf("%x", sumB)+" "+urlB+"\n")

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)

	err := runDownload(cmd, nil, downloadFlagValues{
		file:         urlFile,
		output:       t.TempDir(),
		checksumFile: manifest,
		retries:      3,
		resume:       true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if !strings.Contains(output.String(), "Checksum verified: 2") {
		t.Fatalf("output missing checksum verified count:\n%s", output.String())
	}
}

func TestRunDownloadRootURLChecksumFileRoutesToDownload(t *testing.T) {
	var flags downloadFlagValues
	cmd := &cobra.Command{Use: "daryaft"}
	addDownloadFlags(cmd, &flags)
	if err := cmd.Flags().Set("checksum-file", "checksums.txt"); err != nil {
		t.Fatalf("set checksum-file: %v", err)
	}
	if !hasDownloadFlagChanges(cmd) {
		t.Fatal("hasDownloadFlagChanges = false for --checksum-file, want true")
	}
}

func TestRunDownloadChecksumFileWithHTTPCustomization(t *testing.T) {
	content := []byte("alpha")
	sum := sha256.Sum256(content)
	gotHeader := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Custom")
		_, _ = w.Write(content)
	}))
	defer server.Close()

	urlA := server.URL + "/a.txt"
	manifest := writeCmdChecksumFile(t, "sha256:"+fmt.Sprintf("%x", sum)+" "+urlA+"\n")

	var output bytes.Buffer
	cmd := &cobra.Command{Use: "download"}
	cmd.SetOut(&output)

	err := runDownload(cmd, []string{urlA}, downloadFlagValues{
		output:       t.TempDir(),
		checksumFile: manifest,
		headers:      []string{"X-Custom: value"},
		retries:      3,
		resume:       true,
	})
	if err != nil {
		t.Fatalf("runDownload returned error: %v", err)
	}
	if gotHeader != "value" {
		t.Fatalf("X-Custom header = %q, want value", gotHeader)
	}
	if !strings.Contains(output.String(), "Checksum verified: 1") {
		t.Fatalf("output missing checksum verified count:\n%s", output.String())
	}
}
