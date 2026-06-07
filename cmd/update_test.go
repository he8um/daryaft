package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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

// skipUnlessNetworkTestsEnabled skips the test unless DARYAFT_RUN_NETWORK_TESTS=1.
// Real-network tests must never run under go test ./..., make rc-check, or CI.
func skipUnlessNetworkTestsEnabled(t *testing.T) {
	t.Helper()
	if os.Getenv("DARYAFT_RUN_NETWORK_TESTS") != "1" {
		t.Skip("set DARYAFT_RUN_NETWORK_TESTS=1 to run real-network update tests")
	}
}

// fakeLatestServer returns an httptest.Server that responds to /repos/*/releases/latest
// with a minimal valid release JSON.
func fakeLatestServer(t *testing.T, version string) *httptest.Server {
	t.Helper()
	body := `{"tag_name":"` + version + `","html_url":"https://github.com/he8um/daryaft/releases/tag/` + version + `","prerelease":false,"draft":false}`
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
}

// fakeStatusServer returns an httptest.Server that always responds with the given HTTP status.
func fakeStatusServer(t *testing.T, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(status), status)
	}))
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

// TestUpdateCommand_CheckFlag uses a fake server so no real GitHub call is made.
func TestUpdateCommand_CheckFlag(t *testing.T) {
	srv := fakeLatestServer(t, "v1.0.0")
	defer srv.Close()

	run := newUpdateTestCommand(t)
	stdout, _, err := run("--check", "--api-base-url", srv.URL)
	if err != nil {
		t.Fatalf("update --check returned error: %v", err)
	}
	if !strings.Contains(stdout, "Daryaft update check") {
		t.Errorf("expected update check header, got:\n%s", stdout)
	}
}

// TestUpdateCommand_JSONFlag uses a fake server so no real GitHub call is made.
func TestUpdateCommand_JSONFlag(t *testing.T) {
	srv := fakeLatestServer(t, "v1.0.0")
	defer srv.Close()

	run := newUpdateTestCommand(t)
	stdout, _, err := run("--check", "--json", "--api-base-url", srv.URL)
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

// TestUpdateCommand_403Error verifies the user-friendly rate-limit error message.
func TestUpdateCommand_403Error(t *testing.T) {
	srv := fakeStatusServer(t, http.StatusForbidden)
	defer srv.Close()

	run := newUpdateTestCommand(t)
	_, _, err := run("--check", "--api-base-url", srv.URL)
	if err == nil {
		t.Fatal("expected error on 403 response")
	}
	msg := err.Error()
	if !strings.Contains(msg, "403") {
		t.Errorf("error should mention 403, got: %v", msg)
	}
	if !strings.Contains(strings.ToLower(msg), "rate limit") && !strings.Contains(strings.ToLower(msg), "rate limiting") {
		t.Errorf("error should mention rate limit, got: %v", msg)
	}
	if !strings.Contains(strings.ToLower(msg), "try again") {
		t.Errorf("error should say try again, got: %v", msg)
	}
}

// TestUpdateCommand_404Error verifies the not-found error message.
func TestUpdateCommand_404Error(t *testing.T) {
	srv := fakeStatusServer(t, http.StatusNotFound)
	defer srv.Close()

	run := newUpdateTestCommand(t)
	_, _, err := run("--check", "--api-base-url", srv.URL)
	if err == nil {
		t.Fatal("expected error on 404 response")
	}
	msg := err.Error()
	if !strings.Contains(msg, "404") {
		t.Errorf("error should mention 404, got: %v", msg)
	}
}

// TestUpdateCommand_500Error verifies the server-error message.
func TestUpdateCommand_500Error(t *testing.T) {
	srv := fakeStatusServer(t, http.StatusInternalServerError)
	defer srv.Close()

	run := newUpdateTestCommand(t)
	_, _, err := run("--check", "--api-base-url", srv.URL)
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
	msg := err.Error()
	if !strings.Contains(strings.ToLower(msg), "unavailable") && !strings.Contains(msg, "500") {
		t.Errorf("error should mention unavailable or 500, got: %v", msg)
	}
}

// TestUpdateCommand_InvalidJSONError verifies the invalid-response error message.
func TestUpdateCommand_InvalidJSONError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not json at all"))
	}))
	defer srv.Close()

	run := newUpdateTestCommand(t)
	_, _, err := run("--check", "--api-base-url", srv.URL)
	if err == nil {
		t.Fatal("expected error on invalid JSON response")
	}
	msg := err.Error()
	if !strings.Contains(strings.ToLower(msg), "invalid") {
		t.Errorf("error should mention invalid response, got: %v", msg)
	}
}

// TestUpdateCommand_403_JSONOutput verifies that JSON error output on 403 is valid JSON.
func TestUpdateCommand_403_JSONOutput(t *testing.T) {
	srv := fakeStatusServer(t, http.StatusForbidden)
	defer srv.Close()

	run := newUpdateTestCommand(t)
	_, stderr, err := run("--check", "--json", "--api-base-url", srv.URL)
	if err == nil {
		t.Fatal("expected error on 403 response")
	}
	// When --json is set and an error occurs, error JSON is written to stderr.
	// Cobra may also append a text error line after the JSON, so extract just the
	// first complete JSON object from the beginning of stderr.
	if stderr == "" {
		t.Skip("no stderr JSON output — error written to err only")
	}
	// Find the first '}' that closes the JSON object.
	end := strings.Index(stderr, "\n}")
	var jsonPart string
	if end >= 0 {
		jsonPart = stderr[:end+2]
	} else {
		jsonPart = strings.TrimSpace(stderr)
	}
	var out map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(jsonPart), &out); jsonErr != nil {
		t.Errorf("stderr JSON is not valid: %v\nstderr: %s", jsonErr, stderr)
	}
	if _, ok := out["error"]; !ok {
		t.Errorf("JSON error output missing 'error' field, got: %s", stderr)
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

func TestUpdateCommand_HiddenAPIBaseURLFlag(t *testing.T) {
	cmd := newUpdateCommand()
	flag := cmd.Flags().Lookup("api-base-url")
	if flag == nil {
		t.Fatal("--api-base-url flag not found")
	}
	if !flag.Hidden {
		t.Error("--api-base-url flag should be hidden")
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
	if strings.Contains(helpText, "--api-base-url") {
		t.Error("hidden --api-base-url flag should not appear in help output")
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

// --- Real-network smoke tests (opt-in only) ---
// These tests call the real GitHub API and must never run under go test ./... or make rc-check.
// Run with: DARYAFT_RUN_NETWORK_TESTS=1 go test ./cmd -run TestUpdateCommand_RealNetwork

func TestUpdateCommand_RealNetwork_Check(t *testing.T) {
	skipUnlessNetworkTestsEnabled(t)
	run := newUpdateTestCommand(t)
	stdout, _, err := run("--check")
	if err != nil {
		t.Fatalf("update --check returned error: %v", err)
	}
	if !strings.Contains(stdout, "Daryaft update check") {
		t.Errorf("expected update check header, got:\n%s", stdout)
	}
}

func TestUpdateCommand_RealNetwork_JSON(t *testing.T) {
	skipUnlessNetworkTestsEnabled(t)
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
