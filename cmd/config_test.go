package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	appconfig "github.com/he8um/daryaft/internal/config"
)

func TestConfigPathCommandPrintsPath(t *testing.T) {
	dir := t.TempDir()
	restore := appconfig.SetUserConfigDirForTest(dir)
	t.Cleanup(restore)

	output, err := executeConfigCommand(t, "path")
	if err != nil {
		t.Fatalf("config path returned error: %v", err)
	}

	want := filepath.Join(dir, "daryaft", "config.yaml")
	if strings.TrimSpace(output) != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
}

func TestConfigShowCommandPrintsEffectiveConfig(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	output, err := executeConfigCommand(t, "show")
	if err != nil {
		t.Fatalf("config show returned error: %v", err)
	}

	for _, want := range []string{
		"download_dir: \"\"",
		"retries: 3",
		"resume: true",
		"no_color: false",
		"no_tui: false",
		"theme: default",
		"animations: true",
		"hyperlinks: true",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("config show missing %q in:\n%s", want, output)
		}
	}
}

func TestConfigShowCommandReflectsEnvOverrides(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)
	t.Setenv("DARYAFT_DOWNLOAD_DIR", "/tmp/env-daryaft")
	t.Setenv("DARYAFT_RETRIES", "9")
	t.Setenv("DARYAFT_RESUME", "false")
	t.Setenv("DARYAFT_NO_COLOR", "true")
	t.Setenv("DARYAFT_NO_TUI", "true")
	t.Setenv("DARYAFT_THEME", "env-theme")
	t.Setenv("DARYAFT_ANIMATIONS", "false")
	t.Setenv("DARYAFT_HYPERLINKS", "false")

	output, err := executeConfigCommand(t, "show")
	if err != nil {
		t.Fatalf("config show returned error: %v", err)
	}

	for _, want := range []string{
		"download_dir: /tmp/env-daryaft",
		"retries: 9",
		"resume: false",
		"no_color: true",
		"no_tui: true",
		"theme: env-theme",
		"animations: false",
		"hyperlinks: false",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("config show missing %q in:\n%s", want, output)
		}
	}
}

func TestConfigShowCommandReturnsEnvParseError(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)
	t.Setenv("DARYAFT_RETRIES", "abc")

	_, err := executeConfigCommand(t, "show")
	if err == nil {
		t.Fatal("config show returned nil error")
	}
	if !strings.Contains(err.Error(), "DARYAFT_RETRIES") {
		t.Fatalf("error = %q", err)
	}
}

func TestConfigInitCommandCreatesConfig(t *testing.T) {
	dir := t.TempDir()
	restore := appconfig.SetUserConfigDirForTest(dir)
	t.Cleanup(restore)

	output, err := executeConfigCommand(t, "init")
	if err != nil {
		t.Fatalf("config init returned error: %v", err)
	}

	path := filepath.Join(dir, "daryaft", "config.yaml")
	if !strings.Contains(output, "Created config: "+path) {
		t.Fatalf("output = %q, want created path", output)
	}
	exists, err := appconfig.Exists()
	if err != nil {
		t.Fatalf("Exists returned error: %v", err)
	}
	if !exists {
		t.Fatal("config file does not exist after init")
	}
}

func TestConfigInitCommandRefusesExistingConfig(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	if _, err := executeConfigCommand(t, "init"); err != nil {
		t.Fatalf("initial config init returned error: %v", err)
	}
	_, err := executeConfigCommand(t, "init")
	if err == nil {
		t.Fatal("second config init returned nil error")
	}
	if !strings.Contains(err.Error(), "config already exists") {
		t.Fatalf("error = %q", err)
	}
}

func TestConfigInitCommandForceOverwrites(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	if _, err := executeConfigCommand(t, "init"); err != nil {
		t.Fatalf("initial config init returned error: %v", err)
	}
	cfg := appconfig.Default()
	cfg.Theme = "custom"
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("save custom config: %v", err)
	}

	if _, err := executeConfigCommand(t, "init", "--force"); err != nil {
		t.Fatalf("forced config init returned error: %v", err)
	}

	loaded, err := appconfig.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if loaded.Theme != "default" {
		t.Fatalf("Theme = %q, want default after forced init", loaded.Theme)
	}
}

func executeConfigCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var output bytes.Buffer
	command := newConfigCommand()
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs(args)
	err := command.Execute()
	return output.String(), err
}
