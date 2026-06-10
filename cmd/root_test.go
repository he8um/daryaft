package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	appconfig "github.com/he8um/daryaft/internal/config"
)

// resetRootCmdState resets global state that rootCmd execution may leave behind.
// Must be called at the start of each root_test that uses rootCmd.Execute().
func resetRootCmdState(t *testing.T) {
	t.Helper()
	// Reset configPath package var and the override in internal/config.
	configPath = ""
	t.Cleanup(appconfig.SetConfigPathForTest(""))
	t.Cleanup(func() {
		configPath = ""
		rootCmd.SetArgs(nil)
	})
}

func TestConfigFlagExplicitPathLoads(t *testing.T) {
	resetRootCmdState(t)

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "custom.yaml")
	if err := os.WriteFile(cfgPath, []byte("retries: 7\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"--config", cfgPath, "config", "get", "retries"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if strings.TrimSpace(output.String()) != "7" {
		t.Fatalf("output = %q, want 7", output.String())
	}
}

func TestConfigFlagExplicitPathMissingReturnsError(t *testing.T) {
	resetRootCmdState(t)

	dir := t.TempDir()
	missing := filepath.Join(dir, "nonexistent.yaml")

	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"--config", missing, "config", "path"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error for missing config file")
	}
	if !strings.Contains(err.Error(), "config file not found") {
		t.Fatalf("error = %q, want 'config file not found'", err)
	}
}

func TestConfigFlagOmittedMissingDefaultIsSilent(t *testing.T) {
	resetRootCmdState(t)

	// Point config dir at an empty temp dir so no config.yaml exists.
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"config", "get", "retries"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute returned error when default config missing: %v", err)
	}
	if strings.TrimSpace(output.String()) != "3" {
		t.Fatalf("output = %q, want 3 (default)", output.String())
	}
}

func TestConfigFlagAffectsConfigShow(t *testing.T) {
	resetRootCmdState(t)

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "show.yaml")
	if err := os.WriteFile(cfgPath, []byte("retries: 11\nuser_agent: ShowBot/1\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"--config", cfgPath, "config", "show"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	for _, want := range []string{"retries: 11", "user_agent: ShowBot/1"} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("config show output missing %q:\n%s", want, output.String())
		}
	}
}

func TestConfigFlagAffectsConfigPath(t *testing.T) {
	resetRootCmdState(t)

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "myconfig.yaml")
	if err := os.WriteFile(cfgPath, []byte("retries: 3\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"--config", cfgPath, "config", "path"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if strings.TrimSpace(output.String()) != cfgPath {
		t.Fatalf("config path output = %q, want %q", strings.TrimSpace(output.String()), cfgPath)
	}
}
