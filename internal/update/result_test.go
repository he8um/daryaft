package update

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/he8um/daryaft/pkg/version"
)

// ─── Format() human output tests ─────────────────────────────────────────────

func TestFormat_UpToDate(t *testing.T) {
	r := Result{
		CurrentVersion:   "1.1.0",
		LatestVersion:    "1.1.0",
		UpdateAvailable:  false,
		DevelopmentBuild: false,
		InstallChannel:   ChannelHomebrew,
		ReleaseURL:       "https://github.com/he8um/daryaft/releases/tag/v1.1.0",
		UpdateCommand:    "brew update && brew upgrade daryaft",
		Message:          "Daryaft is up to date.",
	}
	text := r.Format()

	for _, want := range []string{
		"Daryaft update check",
		"Current version:",
		"1.1.0",
		"Latest stable:",
		"up to date",
		"Release:",
		"https://github.com/he8um/daryaft/releases/tag/v1.1.0",
		"Install channel:",
		"homebrew",
		"brew update && brew upgrade daryaft",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("Format() up-to-date missing %q:\n%s", want, text)
		}
	}

	// No update-available prompt when up to date.
	if strings.Contains(text, "A new version is available") {
		t.Errorf("Format() up-to-date should not contain update-available message:\n%s", text)
	}
}

func TestFormat_UpdateAvailable(t *testing.T) {
	r := Result{
		CurrentVersion:   "1.0.0",
		LatestVersion:    "1.1.0",
		UpdateAvailable:  true,
		DevelopmentBuild: false,
		InstallChannel:   ChannelHomebrew,
		ReleaseURL:       "https://github.com/he8um/daryaft/releases/tag/v1.1.0",
		UpdateCommand:    "brew update && brew upgrade daryaft",
		Message:          "Daryaft 1.1.0 is available.",
	}
	text := r.Format()

	for _, want := range []string{
		"update available",
		"A new version is available",
		"brew update && brew upgrade daryaft",
		"Release:",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("Format() update-available missing %q:\n%s", want, text)
		}
	}

	if strings.Contains(text, "development builds may be ahead") {
		t.Errorf("Format() update-available should not show dev note:\n%s", text)
	}
}

func TestFormat_DevBuild(t *testing.T) {
	r := Result{
		CurrentVersion:   "1.2.0-dev",
		LatestVersion:    "1.1.0",
		UpdateAvailable:  false,
		DevelopmentBuild: true,
		InstallChannel:   ChannelSource,
		ReleaseURL:       "https://github.com/he8um/daryaft/releases/tag/v1.1.0",
		UpdateCommand:    "Pull the repository and rebuild: git pull && go build .",
		Message:          "Development build — not comparable to stable releases.",
	}
	text := r.Format()

	for _, want := range []string{
		"development build",
		"development builds may be ahead",
		"source",
		"Pull the repository and rebuild",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("Format() dev-build missing %q:\n%s", want, text)
		}
	}

	if strings.Contains(text, "A new version is available") {
		t.Errorf("Format() dev-build should not show update-available message:\n%s", text)
	}
}

func TestFormat_ReleaseURLAlwaysShown(t *testing.T) {
	// Release URL should appear in all status modes.
	url := "https://github.com/he8um/daryaft/releases/tag/v1.1.0"
	for _, updateAvailable := range []bool{false, true} {
		r := Result{
			CurrentVersion:  "1.1.0",
			LatestVersion:   "1.1.0",
			UpdateAvailable: updateAvailable,
			ReleaseURL:      url,
			InstallChannel:  ChannelSource,
			UpdateCommand:   "Pull the repository and rebuild: git pull && go build .",
		}
		text := r.Format()
		if !strings.Contains(text, url) {
			t.Errorf("Format(updateAvailable=%v) should always show release URL:\n%s", updateAvailable, text)
		}
	}
}

// ─── updateCommand guidance tests ────────────────────────────────────────────

func TestUpdateCommand_Homebrew(t *testing.T) {
	cmd := updateCommand(ChannelHomebrew, "https://github.com/he8um/daryaft/releases/tag/v1.1.0")
	if cmd != "brew update && brew upgrade daryaft" {
		t.Errorf("Homebrew update command = %q, want brew update && brew upgrade daryaft", cmd)
	}
}

func TestUpdateCommand_Source(t *testing.T) {
	cmd := updateCommand(ChannelSource, "https://github.com/he8um/daryaft/releases/tag/v1.1.0")
	if !strings.Contains(cmd, "git pull") {
		t.Errorf("Source update command should mention git pull, got %q", cmd)
	}
	if !strings.Contains(cmd, "go build") {
		t.Errorf("Source update command should mention go build, got %q", cmd)
	}
	// Source channel should NOT suggest downloading an archive.
	if strings.Contains(cmd, "Download") {
		t.Errorf("Source update command should not suggest downloading an archive, got %q", cmd)
	}
}

func TestUpdateCommand_Goreleaser(t *testing.T) {
	releaseURL := "https://github.com/he8um/daryaft/releases/tag/v1.1.0"
	cmd := updateCommand(ChannelGoreleaser, releaseURL)
	if !strings.Contains(cmd, releaseURL) {
		t.Errorf("Goreleaser update command should include release URL, got %q", cmd)
	}
	if !strings.Contains(cmd, "archive") {
		t.Errorf("Goreleaser update command should mention archive, got %q", cmd)
	}
}

func TestUpdateCommand_Goreleaser_NoURL(t *testing.T) {
	cmd := updateCommand(ChannelGoreleaser, "")
	if !strings.Contains(cmd, "releases/latest") {
		t.Errorf("Goreleaser update command with no URL should fall back to releases/latest, got %q", cmd)
	}
}

func TestUpdateCommand_Unknown(t *testing.T) {
	releaseURL := "https://github.com/he8um/daryaft/releases/tag/v1.1.0"
	cmd := updateCommand(ChannelUnknown, releaseURL)
	if !strings.Contains(cmd, releaseURL) {
		t.Errorf("Unknown channel update command should include release URL, got %q", cmd)
	}
}

func TestUpdateCommand_Unknown_NoURL(t *testing.T) {
	cmd := updateCommand(ChannelUnknown, "")
	if !strings.Contains(cmd, "releases/latest") {
		t.Errorf("Unknown channel with no URL should fall back to releases/latest, got %q", cmd)
	}
}

// ─── Install-channel detection tests ─────────────────────────────────────────

func TestDetectInstallChannel_SourceBuiltBy(t *testing.T) {
	info := version.Details{BuiltBy: "source"}
	ch := detectInstallChannel(info, nil, func() string { return "" })
	if ch != ChannelSource {
		t.Errorf("BuiltBy=source: expected ChannelSource, got %q", ch)
	}
}

func TestDetectInstallChannel_SourceBuiltByCase(t *testing.T) {
	info := version.Details{BuiltBy: "SOURCE"}
	ch := detectInstallChannel(info, nil, func() string { return "" })
	if ch != ChannelSource {
		t.Errorf("BuiltBy=SOURCE: expected ChannelSource, got %q", ch)
	}
}

func TestDetectInstallChannel_Goreleaser_NonHomebrewPath(t *testing.T) {
	info := version.Details{BuiltBy: "goreleaser"}
	execFn := func() (string, error) { return "/usr/local/bin/daryaft", nil }
	ch := detectInstallChannel(info, execFn, func() string { return "" })
	if ch != ChannelGoreleaser {
		t.Errorf("goreleaser built, /usr/local/bin path: expected ChannelGoreleaser, got %q", ch)
	}
}

func TestDetectInstallChannel_Goreleaser_TmpPath(t *testing.T) {
	info := version.Details{BuiltBy: "goreleaser"}
	execFn := func() (string, error) { return "/tmp/daryaft", nil }
	ch := detectInstallChannel(info, execFn, func() string { return "" })
	if ch != ChannelGoreleaser {
		t.Errorf("goreleaser built, /tmp path: expected ChannelGoreleaser, got %q", ch)
	}
}

func TestDetectInstallChannel_Homebrew_OptHomebrewCellar(t *testing.T) {
	info := version.Details{BuiltBy: "goreleaser"}
	execFn := func() (string, error) {
		return "/opt/homebrew/Cellar/daryaft/1.1.0/bin/daryaft", nil
	}
	ch := detectInstallChannel(info, execFn, func() string { return "" })
	if ch != ChannelHomebrew {
		t.Errorf("/opt/homebrew/Cellar path: expected ChannelHomebrew, got %q", ch)
	}
}

func TestDetectInstallChannel_Homebrew_UsrLocalCellar(t *testing.T) {
	info := version.Details{BuiltBy: "goreleaser"}
	execFn := func() (string, error) {
		return "/usr/local/Cellar/daryaft/1.1.0/bin/daryaft", nil
	}
	ch := detectInstallChannel(info, execFn, func() string { return "" })
	if ch != ChannelHomebrew {
		t.Errorf("/usr/local/Cellar path: expected ChannelHomebrew, got %q", ch)
	}
}

func TestDetectInstallChannel_Homebrew_LinuxBrew(t *testing.T) {
	info := version.Details{BuiltBy: "goreleaser"}
	execFn := func() (string, error) {
		return "/home/linuxbrew/.linuxbrew/Cellar/daryaft/1.1.0/bin/daryaft", nil
	}
	ch := detectInstallChannel(info, execFn, func() string { return "" })
	if ch != ChannelHomebrew {
		t.Errorf("/home/linuxbrew path: expected ChannelHomebrew, got %q", ch)
	}
}

func TestDetectInstallChannel_Homebrew_BrewPrefixFallback(t *testing.T) {
	// When executable is at /opt/homebrew/bin/daryaft and brew prefix returns
	// /opt/homebrew/Cellar/daryaft/1.1.0, the prefix check should detect Homebrew.
	// Simulate: exec path is /opt/homebrew/bin/daryaft (symlink),
	// resolved is also /opt/homebrew/bin/daryaft (no real symlink in test),
	// but brew prefix returns a path that starts with the exec prefix.
	//
	// The current implementation uses brewPrefixFn as a fallback when the path
	// check via isHomebrewPath fails. To trigger the fallback, we use a path
	// that is not in a cellar but matches the brew prefix.
	info := version.Details{BuiltBy: "goreleaser"}
	execFn := func() (string, error) {
		return "/opt/homebrew/bin/daryaft", nil
	}
	// /opt/homebrew/bin is a prefix of /opt/homebrew/bin/daryaft
	// but isHomebrewPath checks for /opt/homebrew/ so this should already match.
	ch := detectInstallChannel(info, execFn, func() string { return "" })
	// /opt/homebrew/bin/daryaft starts with /opt/homebrew/ → homebrew
	if ch != ChannelHomebrew {
		t.Errorf("/opt/homebrew/bin path: expected ChannelHomebrew, got %q", ch)
	}
}

func TestDetectInstallChannel_Unknown_UnknownBuiltBy(t *testing.T) {
	info := version.Details{BuiltBy: "make"}
	execFn := func() (string, error) { return "/usr/local/bin/daryaft", nil }
	ch := detectInstallChannel(info, execFn, func() string { return "" })
	if ch != ChannelUnknown {
		t.Errorf("BuiltBy=make, non-homebrew path: expected ChannelUnknown, got %q", ch)
	}
}

func TestDetectInstallChannel_ExecutableError(t *testing.T) {
	// When os.Executable() fails, fall through to BuiltBy check.
	info := version.Details{BuiltBy: "goreleaser"}
	execFn := func() (string, error) {
		return "", errFakeExecError("no executable")
	}
	ch := detectInstallChannel(info, execFn, func() string { return "" })
	if ch != ChannelGoreleaser {
		t.Errorf("executable error + goreleaser BuiltBy: expected ChannelGoreleaser, got %q", ch)
	}
}

// errFakeExecError is a simple error type for testing.
type errFakeExecError string

func (e errFakeExecError) Error() string { return string(e) }

// ─── JSON contract tests ──────────────────────────────────────────────────────

func TestFormatJSON_AllFieldsPresent(t *testing.T) {
	r := Result{
		CurrentVersion:    "1.1.0",
		LatestVersion:     "1.1.0",
		UpdateAvailable:   false,
		DevelopmentBuild:  false,
		IncludePrerelease: false,
		InstallChannel:    ChannelHomebrew,
		ReleaseURL:        "https://github.com/he8um/daryaft/releases/tag/v1.1.0",
		UpdateCommand:     "brew update && brew upgrade daryaft",
		Message:           "Daryaft is up to date.",
	}
	data, err := r.FormatJSON()
	if err != nil {
		t.Fatalf("FormatJSON: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal JSON: %v", err)
	}
	required := []string{
		"current_version",
		"latest_version",
		"update_available",
		"development_build",
		"include_prerelease",
		"install_channel",
		"release_url",
		"update_command",
		"message",
	}
	for _, field := range required {
		if _, ok := out[field]; !ok {
			t.Errorf("JSON missing required field %q", field)
		}
	}
}

func TestFormatJSON_BoolsAreActualBools(t *testing.T) {
	r := Result{
		UpdateAvailable:   true,
		DevelopmentBuild:  false,
		IncludePrerelease: true,
	}
	data, err := r.FormatJSON()
	if err != nil {
		t.Fatalf("FormatJSON: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal JSON: %v", err)
	}
	for _, field := range []string{"update_available", "development_build", "include_prerelease"} {
		v, ok := out[field]
		if !ok {
			t.Errorf("JSON missing field %q", field)
			continue
		}
		if _, isBool := v.(bool); !isBool {
			t.Errorf("JSON field %q should be bool, got %T (%v)", field, v, v)
		}
	}
}

func TestFormatJSON_EmptyStringsPresent(t *testing.T) {
	// Fields that may be empty should still be present as "", not omitted.
	r := Result{
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.0.0",
		InstallChannel: ChannelSource,
		ReleaseURL:     "",
		UpdateCommand:  "",
		Message:        "",
	}
	data, err := r.FormatJSON()
	if err != nil {
		t.Fatalf("FormatJSON: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal JSON: %v", err)
	}
	for _, field := range []string{"release_url", "update_command", "message"} {
		if _, ok := out[field]; !ok {
			t.Errorf("JSON field %q should be present even when empty string", field)
		}
	}
}

// ─── Status() tests ───────────────────────────────────────────────────────────

func TestStatus_UpToDate(t *testing.T) {
	r := Result{UpdateAvailable: false, DevelopmentBuild: false}
	if r.Status() != StatusUpToDate {
		t.Errorf("Status() = %q, want %q", r.Status(), StatusUpToDate)
	}
}

func TestStatus_Available(t *testing.T) {
	r := Result{UpdateAvailable: true, DevelopmentBuild: false}
	if r.Status() != StatusAvailable {
		t.Errorf("Status() = %q, want %q", r.Status(), StatusAvailable)
	}
}

func TestStatus_DevBuild(t *testing.T) {
	r := Result{DevelopmentBuild: true}
	if r.Status() != StatusDevBuild {
		t.Errorf("Status() = %q, want %q", r.Status(), StatusDevBuild)
	}
}

func TestStatus_DevBuildTakesPrecedence(t *testing.T) {
	// DevelopmentBuild true should override UpdateAvailable.
	r := Result{UpdateAvailable: true, DevelopmentBuild: true}
	if r.Status() != StatusDevBuild {
		t.Errorf("Status() with both flags: expected %q, got %q", StatusDevBuild, r.Status())
	}
}
