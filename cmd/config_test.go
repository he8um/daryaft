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
	t.Setenv("DARYAFT_THEME", "mono")
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
		"theme: mono",
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

func TestConfigGetCommandPrintsEffectiveValue(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)
	t.Setenv("DARYAFT_RETRIES", "8")

	output, err := executeConfigCommand(t, "get", "retries")
	if err != nil {
		t.Fatalf("config get returned error: %v", err)
	}
	if strings.TrimSpace(output) != "8" {
		t.Fatalf("output = %q, want 8", output)
	}
}

func TestConfigGetCommandUnknownKey(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	_, err := executeConfigCommand(t, "get", "missing")
	if err == nil {
		t.Fatal("config get returned nil error")
	}
	if !strings.Contains(err.Error(), "unknown config key") {
		t.Fatalf("error = %q", err)
	}
}

func TestConfigSetCommandWritesFileConfig(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	output, err := executeConfigCommand(t, "set", "download_dir", " downloads ")
	if err != nil {
		t.Fatalf("config set returned error: %v", err)
	}
	if strings.TrimSpace(output) != "Updated config: download_dir=downloads" {
		t.Fatalf("output = %q", output)
	}

	cfg, err := appconfig.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.DownloadDir != "downloads" {
		t.Fatalf("DownloadDir = %q, want downloads", cfg.DownloadDir)
	}
}

func TestConfigSetCommandRejectsInvalidValue(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	_, err := executeConfigCommand(t, "set", "retries", "-1")
	if err == nil {
		t.Fatal("config set returned nil error")
	}
	if !strings.Contains(err.Error(), "retries") {
		t.Fatalf("error = %q", err)
	}
}

func TestConfigSetCommandDoesNotWriteEnv(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)
	t.Setenv("DARYAFT_RETRIES", "9")

	if _, err := executeConfigCommand(t, "set", "retries", "5"); err != nil {
		t.Fatalf("config set returned error: %v", err)
	}

	fileCfg, err := appconfig.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if fileCfg.Retries != 5 {
		t.Fatalf("file Retries = %d, want 5", fileCfg.Retries)
	}

	effective, err := appconfig.LoadEffective()
	if err != nil {
		t.Fatalf("LoadEffective returned error: %v", err)
	}
	if effective.Retries != 9 {
		t.Fatalf("effective Retries = %d, want env override 9", effective.Retries)
	}
}

func TestConfigResetCommandWritesDefaults(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	cfg := appconfig.Default()
	cfg.Retries = 12
	if err := appconfig.Save(cfg); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	output, err := executeConfigCommand(t, "reset")
	if err != nil {
		t.Fatalf("config reset returned error: %v", err)
	}
	if !strings.Contains(output, "Reset config: ") {
		t.Fatalf("output = %q", output)
	}

	loaded, err := appconfig.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if loaded != appconfig.Default() {
		t.Fatalf("Load = %#v, want defaults %#v", loaded, appconfig.Default())
	}
}

func TestConfigKeysCommandListsSupportedKeys(t *testing.T) {
	output, err := executeConfigCommand(t, "keys")
	if err != nil {
		t.Fatalf("config keys returned error: %v", err)
	}

	want := strings.Join([]string{
		"download_dir string",
		"retries int",
		"resume bool",
		"no_color bool",
		"no_tui bool",
		"theme string",
		"animations bool",
		"hyperlinks bool",
		"",
	}, "\n")
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
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
