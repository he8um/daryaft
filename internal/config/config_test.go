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

func TestSaveWritesValidYAMLWithPrivatePermissions(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	cfg := Default()
	cfg.Retries = 4
	if err := Save(cfg); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	path, err := Path()
	if err != nil {
		t.Fatalf("Path returned error: %v", err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if loaded.Retries != 4 {
		t.Fatalf("Retries = %d, want 4", loaded.Retries)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("config mode = %o, want 600", got)
	}
}

func TestSaveDoesNotLeaveTemporaryFiles(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	if err := Save(Default()); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	path, err := Path()
	if err != nil {
		t.Fatalf("Path returned error: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(filepath.Dir(path), ".config-*.tmp"))
	if err != nil {
		t.Fatalf("glob temp files: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary files left behind: %#v", matches)
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
	for _, value := range []string{"abc", "-1", "21", ""} {
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
	t.Setenv(envTheme, " mono ")
	t.Setenv(envAnimations, "false")
	t.Setenv(envHyperlinks, "false")

	cfg := Default()
	cfg.DownloadDir = "/config/downloads"
	cfg.Retries = 2
	cfg.Resume = true
	cfg.NoColor = false
	cfg.NoTUI = false
	cfg.Theme = "default"
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
	if got.Theme != "mono" {
		t.Fatalf("Theme = %q", got.Theme)
	}
	if got.Animations {
		t.Fatal("Animations = true, want false")
	}
	if got.Hyperlinks {
		t.Fatal("Hyperlinks = true, want false")
	}
}

func TestGetSupportedKeys(t *testing.T) {
	cfg := Config{
		DownloadDir: "/downloads",
		Retries:     7,
		Resume:      false,
		NoColor:     true,
		NoTUI:       true,
		Theme:       "plain",
		Animations:  false,
		Hyperlinks:  false,
	}

	tests := map[string]string{
		keyDownloadDir: "/downloads",
		keyRetries:     "7",
		keyResume:      "false",
		keyNoColor:     "true",
		keyNoTUI:       "true",
		keyTheme:       "plain",
		keyAnimations:  "false",
		keyHyperlinks:  "false",
	}
	for key, want := range tests {
		t.Run(key, func(t *testing.T) {
			got, err := Get(cfg, key)
			if err != nil {
				t.Fatalf("Get returned error: %v", err)
			}
			if got != want {
				t.Fatalf("Get(%q) = %q, want %q", key, got, want)
			}
		})
	}
}

func TestGetUnknownKeyReturnsError(t *testing.T) {
	_, err := Get(Default(), "missing")
	if err == nil {
		t.Fatal("Get returned nil error")
	}
	if !strings.Contains(err.Error(), "unknown config key") {
		t.Fatalf("error = %q", err)
	}
}

func TestSetStringValues(t *testing.T) {
	cfg, err := Set(Default(), keyDownloadDir, "  /downloads  ")
	if err != nil {
		t.Fatalf("Set download_dir returned error: %v", err)
	}
	if cfg.DownloadDir != "/downloads" {
		t.Fatalf("DownloadDir = %q", cfg.DownloadDir)
	}

	cfg, err = Set(cfg, keyTheme, "  mono  ")
	if err != nil {
		t.Fatalf("Set theme returned error: %v", err)
	}
	if cfg.Theme != "mono" {
		t.Fatalf("Theme = %q", cfg.Theme)
	}
}

func TestSetInvalidThemeReturnsError(t *testing.T) {
	_, err := Set(Default(), keyTheme, "high-contrast")
	if err == nil {
		t.Fatal("Set returned nil error")
	}
	if !strings.Contains(err.Error(), "invalid theme") {
		t.Fatalf("error = %q", err)
	}
}

func TestApplyEnvInvalidThemeReturnsError(t *testing.T) {
	_, err := ApplyEnv(Default(), mapLookup(map[string]string{
		envTheme: "high-contrast",
	}))
	if err == nil {
		t.Fatal("ApplyEnv returned nil error")
	}
	if !strings.Contains(err.Error(), envTheme) {
		t.Fatalf("error = %q", err)
	}
}

func TestSetIntValue(t *testing.T) {
	cfg, err := Set(Default(), keyRetries, " 5 ")
	if err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	if cfg.Retries != 5 {
		t.Fatalf("Retries = %d, want 5", cfg.Retries)
	}
}

func TestSetRetriesAllowsUpperBound(t *testing.T) {
	cfg, err := Set(Default(), keyRetries, "20")
	if err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	if cfg.Retries != 20 {
		t.Fatalf("Retries = %d, want 20", cfg.Retries)
	}
}

func TestSetInvalidIntReturnsError(t *testing.T) {
	for _, value := range []string{"abc", "-1", "21", ""} {
		t.Run(strconv.Quote(value), func(t *testing.T) {
			_, err := Set(Default(), keyRetries, value)
			if err == nil {
				t.Fatal("Set returned nil error")
			}
			if !strings.Contains(err.Error(), keyRetries) {
				t.Fatalf("error = %q", err)
			}
		})
	}
}

func TestSetBoolVariants(t *testing.T) {
	for _, value := range []string{"true", "TRUE", "1", "yes", "Y", "on"} {
		t.Run("true/"+value, func(t *testing.T) {
			cfg, err := Set(Default(), keyNoTUI, value)
			if err != nil {
				t.Fatalf("Set returned error: %v", err)
			}
			if !cfg.NoTUI {
				t.Fatalf("NoTUI = false for %q", value)
			}
		})
	}
	for _, value := range []string{"false", "FALSE", "0", "no", "N", "off"} {
		t.Run("false/"+value, func(t *testing.T) {
			cfg := Default()
			cfg.NoTUI = true
			got, err := Set(cfg, keyNoTUI, value)
			if err != nil {
				t.Fatalf("Set returned error: %v", err)
			}
			if got.NoTUI {
				t.Fatalf("NoTUI = true for %q", value)
			}
		})
	}
}

func TestSetInvalidBoolReturnsError(t *testing.T) {
	for _, value := range []string{"maybe", ""} {
		t.Run(strconv.Quote(value), func(t *testing.T) {
			_, err := Set(Default(), keyResume, value)
			if err == nil {
				t.Fatal("Set returned nil error")
			}
			if !strings.Contains(err.Error(), keyResume) {
				t.Fatalf("error = %q", err)
			}
		})
	}
}

func TestSetUnknownKeyReturnsError(t *testing.T) {
	_, err := Set(Default(), "missing", "value")
	if err == nil {
		t.Fatal("Set returned nil error")
	}
	if !strings.Contains(err.Error(), "unknown config key") {
		t.Fatalf("error = %q", err)
	}
}

func TestResetWritesDefaults(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(SetUserConfigDirForTest(dir))

	cfg := Default()
	cfg.Retries = 10
	if err := Save(cfg); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	path, err := Reset()
	if err != nil {
		t.Fatalf("Reset returned error: %v", err)
	}
	wantPath := filepath.Join(dir, "daryaft", "config.yaml")
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got != Default() {
		t.Fatalf("Load = %#v, want defaults %#v", got, Default())
	}
}

func TestSupportedKeysListsAllKeys(t *testing.T) {
	keys := SupportedKeys()
	want := []KeyInfo{
		{Key: keyDownloadDir, Type: "string"},
		{Key: keyRetries, Type: "int"},
		{Key: keyResume, Type: "bool"},
		{Key: keyNoColor, Type: "bool"},
		{Key: keyNoTUI, Type: "bool"},
		{Key: keyTheme, Type: "string"},
		{Key: keyAnimations, Type: "bool"},
		{Key: keyHyperlinks, Type: "bool"},
	}

	if len(keys) != len(want) {
		t.Fatalf("len(SupportedKeys) = %d, want %d", len(keys), len(want))
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Fatalf("SupportedKeys()[%d] = %#v, want %#v", i, keys[i], want[i])
		}
	}
}

func mapLookup(values map[string]string) LookupEnvFunc {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
