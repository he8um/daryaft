package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/update"
)

// newUpdateTestCommand creates an isolated update command instance for testing.
func newUpdateTestCommand(t *testing.T) func(args ...string) (string, string, error) {
	t.Helper()
	return func(args ...string) (stdout, stderr string, err error) {
		var outBuf, errBuf bytes.Buffer
		cmd := newUpdateCommand()
		cmd.SetOut(&outBuf)
		cmd.SetErr(&errBuf)
		cmd.SetArgs(args)
		err = cmd.Execute()
		return outBuf.String(), errBuf.String(), err
	}
}

func TestUpdateCommand_WithoutCheck(t *testing.T) {
	run := newUpdateTestCommand(t)
	_, _, err := run()
	if err == nil {
		t.Fatal("expected error when running update without --check")
	}
	if !strings.Contains(err.Error(), "auto-update is not implemented") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestUpdateCommand_WithoutCheck_SuggestsCheck(t *testing.T) {
	run := newUpdateTestCommand(t)
	_, _, err := run()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "--check") {
		t.Errorf("error should mention --check flag, got: %v", err)
	}
}

func TestUpdateCommand_CheckFlag_NoRepo(t *testing.T) {
	// Running --check without a --repo override hits the real GitHub API.
	// Skip in short mode to keep unit tests offline.
	if testing.Short() {
		t.Skip("skipping real-network test in short mode")
	}

	run := newUpdateTestCommand(t)
	stdout, _, err := run("--check")
	if err != nil {
		t.Fatalf("update --check returned error: %v", err)
	}
	if !strings.Contains(stdout, "Daryaft update check") {
		t.Errorf("expected update check header, got:\n%s", stdout)
	}
}

func TestUpdateCommand_JSONFlag_NoRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real-network test in short mode")
	}

	run := newUpdateTestCommand(t)
	stdout, _, err := run("--check", "--json")
	if err != nil {
		t.Fatalf("update --check --json returned error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, stdout)
	}
	for _, field := range []string{"current_version", "latest_version", "update_available",
		"development_build", "install_channel", "release_url", "update_command", "message"} {
		if _, ok := out[field]; !ok {
			t.Errorf("JSON output missing field %q", field)
		}
	}
}

func TestUpdateCommand_HiddenRepoFlag(t *testing.T) {
	cmd := newUpdateCommand()
	flag := cmd.Flags().Lookup("repo")
	if flag == nil {
		t.Fatal("--repo flag not found")
	}
	if !flag.Hidden {
		t.Error("--repo flag should be hidden")
	}
}

func TestUpdateCommand_HelpNotContainsRepo(t *testing.T) {
	var helpBuf bytes.Buffer
	cmd := newUpdateCommand()
	cmd.SetOut(&helpBuf)
	cmd.SetArgs([]string{"--help"})
	_ = cmd.Execute()
	helpText := helpBuf.String()
	if strings.Contains(helpText, "--repo") {
		t.Error("hidden --repo flag should not appear in help output")
	}
}

func TestSplitRepo(t *testing.T) {
	tests := []struct {
		input     string
		wantOwner string
		wantRepo  string
	}{
		{"owner/repo", "owner", "repo"},
		{"he8um/daryaft", "he8um", "daryaft"},
		{"noslash", "", ""},
	}
	for _, tc := range tests {
		owner, repo := splitRepo(tc.input)
		if owner != tc.wantOwner || repo != tc.wantRepo {
			t.Errorf("splitRepo(%q) = (%q, %q), want (%q, %q)",
				tc.input, owner, repo, tc.wantOwner, tc.wantRepo)
		}
	}
}

// TestUpdateResult_Format exercises the human output formatting.
func TestUpdateResult_Format_UpToDate(t *testing.T) {
	r := update.Result{
		CurrentVersion:   "1.0.0",
		LatestVersion:    "1.0.0",
		UpdateAvailable:  false,
		DevelopmentBuild: false,
		InstallChannel:   update.ChannelHomebrew,
		ReleaseURL:       "https://github.com/he8um/daryaft/releases/tag/v1.0.0",
		UpdateCommand:    "brew update && brew upgrade daryaft",
		Message:          "Daryaft is up to date.",
	}
	text := r.Format()
	for _, want := range []string{
		"Daryaft update check",
		"Current version",
		"1.0.0",
		"up to date",
		"homebrew",
		"brew update && brew upgrade daryaft",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("Format() missing %q:\n%s", want, text)
		}
	}
}

func TestUpdateResult_Format_DevBuild(t *testing.T) {
	r := update.Result{
		CurrentVersion:   "1.1.0-dev",
		LatestVersion:    "1.0.0",
		UpdateAvailable:  false,
		DevelopmentBuild: true,
		InstallChannel:   update.ChannelSource,
		ReleaseURL:       "https://github.com/he8um/daryaft/releases/tag/v1.0.0",
		UpdateCommand:    "Download the latest release from: https://github.com/he8um/daryaft/releases/tag/v1.0.0",
		Message:          "Development build — not comparable to stable releases.",
	}
	text := r.Format()
	if !strings.Contains(text, "development build") {
		t.Errorf("Format() missing development build status:\n%s", text)
	}
	if !strings.Contains(text, "development builds may be ahead") {
		t.Errorf("Format() missing dev build note:\n%s", text)
	}
}
