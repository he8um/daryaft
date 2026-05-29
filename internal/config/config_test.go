package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := Default()

	if cfg.DownloadDir != "" {
		t.Fatalf("DownloadDir = %q, want empty", cfg.DownloadDir)
	}
	if cfg.Retries != 3 {
		t.Fatalf("Retries = %d, want 3", cfg.Retries)
	}
	if !cfg.Resume {
		t.Fatal("Resume = false, want true")
	}
	if cfg.NoColor {
		t.Fatal("NoColor = true, want false")
	}
	if cfg.NoTUI {
		t.Fatal("NoTUI = true, want false")
	}
	if cfg.Theme != "default" {
		t.Fatalf("Theme = %q, want default", cfg.Theme)
	}
	if !cfg.Animations {
		t.Fatal("Animations = false, want true")
	}
	if !cfg.Hyperlinks {
		t.Fatal("Hyperlinks = false, want true")
	}
}

func TestPathUsesUserConfigDir(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	path, err := Path()
	if err != nil {
		t.Fatalf("Path returned error: %v", err)
	}

	want := filepath.Join(dir, "daryaft", "config.yaml")
	if path != want {
		t.Fatalf("Path = %q, want %q", path, want)
	}
}

func TestInitCreatesConfigFile(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	path, err := Init(false)
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
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
		if !strings.Contains(string(data), want) {
			t.Fatalf("config missing %q in:\n%s", want, string(data))
		}
	}
}

func TestInitRefusesOverwriteWithoutForce(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	path, err := Init(false)
	if err != nil {
		t.Fatalf("initial Init returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte("theme: custom\n"), 0o644); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	_, err = Init(false)
	if err == nil {
		t.Fatal("Init returned nil error, want existing-file error")
	}
	if !strings.Contains(err.Error(), "config already exists") {
		t.Fatalf("error = %q", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if string(data) != "theme: custom\n" {
		t.Fatalf("config was overwritten:\n%s", string(data))
	}
}

func TestInitForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	path, err := Init(false)
	if err != nil {
		t.Fatalf("initial Init returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte("theme: custom\n"), 0o644); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	if _, err := Init(true); err != nil {
		t.Fatalf("forced Init returned error: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.Theme != "default" {
		t.Fatalf("Theme = %q, want default after force init", cfg.Theme)
	}
}

func TestLoadReturnsDefaultsWhenFileMissing(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg != Default() {
		t.Fatalf("Load = %#v, want defaults %#v", cfg, Default())
	}
}

func TestLoadReadsYAMLAndMergesDefaults(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	path, err := Path()
	if err != nil {
		t.Fatalf("Path returned error: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("download_dir: /tmp/daryaft\nretries: 5\nresume: false\nno_color: true\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.DownloadDir != "/tmp/daryaft" {
		t.Fatalf("DownloadDir = %q", cfg.DownloadDir)
	}
	if cfg.Retries != 5 {
		t.Fatalf("Retries = %d, want 5", cfg.Retries)
	}
	if cfg.Resume {
		t.Fatal("Resume = true, want false")
	}
	if !cfg.NoColor {
		t.Fatal("NoColor = false, want true")
	}
	if cfg.Theme != "default" {
		t.Fatalf("Theme = %q, want default from merge", cfg.Theme)
	}
	if !cfg.Animations || !cfg.Hyperlinks {
		t.Fatalf("future TUI fields not merged from defaults: %#v", cfg)
	}
}

func TestLoadInvalidYAMLReturnsClearError(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	path, err := Path()
	if err != nil {
		t.Fatalf("Path returned error: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("download_dir: [\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err = Load()
	if err == nil {
		t.Fatal("Load returned nil error, want parse error")
	}
	if !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("error = %q", err)
	}
}

func TestSaveCreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(filepath.Join(dir, "nested", "config-root")))

	cfg := Default()
	cfg.DownloadDir = "/tmp/daryaft"

	if err := Save(cfg); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	path, err := Path()
	if err != nil {
		t.Fatalf("Path returned error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stat config: %v", err)
	}
}

func TestApplyEnvOverridesDownloadDir(t *testing.T) {
	cfg := Default()
	cfg.DownloadDir = "/config/downloads"

	got, err := ApplyEnv(cfg, mapLookup(map[string]string{
		envDownloadDir: "  /env/downloads  ",
	}))
	if err != nil {
		t.Fatalf("ApplyEnv returned error: %v", err)
	}
	if got.DownloadDir != "/env/downloads" {
		t.Fatalf("DownloadDir = %q", got.DownloadDir)
	}
}

func TestApplyEnvOverridesRetries(t *testing.T) {
	cfg := Default()
	cfg.Retries = 2

	got, err := ApplyEnv(cfg, mapLookup(map[string]string{
		envRetries: "5",
	}))
	if err != nil {
		t.Fatalf("ApplyEnv returned error: %v", err)
	}
	if got.Retries != 5 {
		t.Fatalf("Retries = %d, want 5", got.Retries)
	}
}

func TestApplyEnvInvalidRetriesReturnsError(t *testing.T) {
	for _, value := range []string{"abc", "-1", ""} {
		t.Run(strconv.Quote(value), func(t *testing.T) {
			_, err := ApplyEnv(Default(), mapLookup(map[string]string{
				envRetries: value,
			}))
			if err == nil {
				t.Fatal("ApplyEnv returned nil error")
			}
			if !strings.Contains(err.Error(), envRetries) {
				t.Fatalf("error = %q", err)
			}
		})
	}
}

func TestApplyEnvBooleanTrueValues(t *testing.T) {
	for _, value := range []string{"true", "TRUE", "1", "yes", "Y", "on"} {
		t.Run(value, func(t *testing.T) {
			got, err := ApplyEnv(Default(), mapLookup(map[string]string{
				envNoTUI: value,
			}))
			if err != nil {
				t.Fatalf("ApplyEnv returned error: %v", err)
			}
			if !got.NoTUI {
				t.Fatalf("NoTUI = false for %q", value)
			}
		})
	}
}

func TestApplyEnvBooleanFalseValues(t *testing.T) {
	for _, value := range []string{"false", "FALSE", "0", "no", "N", "off"} {
		t.Run(value, func(t *testing.T) {
			cfg := Default()
			cfg.NoTUI = true
			got, err := ApplyEnv(cfg, mapLookup(map[string]string{
				envNoTUI: value,
			}))
			if err != nil {
				t.Fatalf("ApplyEnv returned error: %v", err)
			}
			if got.NoTUI {
				t.Fatalf("NoTUI = true for %q", value)
			}
		})
	}
}

func TestApplyEnvInvalidBooleanReturnsError(t *testing.T) {
	for _, value := range []string{"maybe", ""} {
		t.Run(strconv.Quote(value), func(t *testing.T) {
			_, err := ApplyEnv(Default(), mapLookup(map[string]string{
				envResume: value,
			}))
			if err == nil {
				t.Fatal("ApplyEnv returned nil error")
			}
			if !strings.Contains(err.Error(), envResume) {
				t.Fatalf("error = %q", err)
			}
		})
	}
}

func TestLoadEffectiveAppliesEnvOverFileConfig(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))
	t.Setenv(envDownloadDir, "/env/downloads")
	t.Setenv(envRetries, "8")
	t.Setenv(envResume, "false")
	t.Setenv(envNoColor, "true")
	t.Setenv(envNoTUI, "true")
	t.Setenv(envTheme, " env-theme ")
	t.Setenv(envAnimations, "false")
	t.Setenv(envHyperlinks, "false")

	cfg := Default()
	cfg.DownloadDir = "/config/downloads"
	cfg.Retries = 2
	cfg.Resume = true
	cfg.NoColor = false
	cfg.NoTUI = false
	cfg.Theme = "config-theme"
	cfg.Animations = true
	cfg.Hyperlinks = true
	if err := Save(cfg); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	got, err := LoadEffective()
	if err != nil {
		t.Fatalf("LoadEffective returned error: %v", err)
	}
	if got.DownloadDir != "/env/downloads" {
		t.Fatalf("DownloadDir = %q", got.DownloadDir)
	}
	if got.Retries != 8 {
		t.Fatalf("Retries = %d, want 8", got.Retries)
	}
	if got.Resume {
		t.Fatal("Resume = true, want false")
	}
	if !got.NoColor {
		t.Fatal("NoColor = false, want true")
	}
	if !got.NoTUI {
		t.Fatal("NoTUI = false, want true")
	}
	if got.Theme != "env-theme" {
		t.Fatalf("Theme = %q", got.Theme)
	}
	if got.Animations {
		t.Fatal("Animations = true, want false")
	}
	if got.Hyperlinks {
		t.Fatal("Hyperlinks = true, want false")
	}
}

func mapLookup(values map[string]string) LookupEnvFunc {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
