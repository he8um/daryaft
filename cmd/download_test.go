package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	appconfig "github.com/he8um/daryaft/internal/config"
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
