package update

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/he8um/daryaft/pkg/version"
)

// InstallChannel is the detected install mechanism.
type InstallChannel string

const (
	ChannelHomebrew   InstallChannel = "homebrew"
	ChannelGoreleaser InstallChannel = "goreleaser"
	ChannelSource     InstallChannel = "source"
	ChannelUnknown    InstallChannel = "unknown"
)

// Status is the update check outcome.
type Status string

const (
	StatusUpToDate  Status = "up to date"
	StatusAvailable Status = "update available"
	StatusDevBuild  Status = "development build"
	StatusUnknown   Status = "unknown"
)

// Result holds the complete update check result.
type Result struct {
	CurrentVersion    string         `json:"current_version"`
	LatestVersion     string         `json:"latest_version"`
	UpdateAvailable   bool           `json:"update_available"`
	DevelopmentBuild  bool           `json:"development_build"`
	IncludePrerelease bool           `json:"include_prerelease"`
	InstallChannel    InstallChannel `json:"install_channel"`
	ReleaseURL        string         `json:"release_url"`
	UpdateCommand     string         `json:"update_command"`
	Message           string         `json:"message"`
}

// FormatJSON returns the result as indented JSON bytes.
func (r Result) FormatJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// Format returns a human-readable update check report.
func (r Result) Format() string {
	var b strings.Builder
	b.WriteString("Daryaft update check\n\n")
	fmt.Fprintf(&b, "Current version:  %s\n", r.CurrentVersion)
	fmt.Fprintf(&b, "Latest stable:    %s\n", r.LatestVersion)
	fmt.Fprintf(&b, "Status:           %s\n", string(r.Status()))

	if r.ReleaseURL != "" && (r.UpdateAvailable || r.DevelopmentBuild) {
		fmt.Fprintf(&b, "\nRelease: %s\n", r.ReleaseURL)
	}

	fmt.Fprintf(&b, "\nInstall channel:  %s\n", string(r.InstallChannel))
	if r.UpdateCommand != "" {
		fmt.Fprintf(&b, "Update command:   %s\n", r.UpdateCommand)
	}

	if r.DevelopmentBuild {
		b.WriteString("\nNote: development builds may be ahead of the latest stable release.\n")
	}

	return b.String()
}

// Status computes the display status.
func (r Result) Status() Status {
	if r.DevelopmentBuild {
		return StatusDevBuild
	}
	if r.UpdateAvailable {
		return StatusAvailable
	}
	return StatusUpToDate
}

// DetectInstallChannel returns the install channel based on the executable path
// and build metadata.
func DetectInstallChannel(info version.Details) InstallChannel {
	return detectInstallChannel(info, os.Executable, execBrewPrefix)
}

func detectInstallChannel(info version.Details, executableFn func() (string, error), brewPrefixFn func() string) InstallChannel {
	// BuiltBy is injected by GoReleaser or make; "source" means built from source.
	switch strings.ToLower(info.BuiltBy) {
	case "source":
		return ChannelSource
	case "goreleaser":
		// May still be a Homebrew install if GoReleaser-built binary is delivered via Homebrew.
		// Fall through to path check.
	}

	// Check executable path against Homebrew prefix/cellar.
	execPath, err := executableFn()
	if err == nil {
		execPath = filepath.Clean(execPath)
		// Resolve symlinks so /usr/local/bin/daryaft -> /opt/homebrew/Cellar/... resolves correctly.
		if resolved, err := filepath.EvalSymlinks(execPath); err == nil {
			execPath = resolved
		}
		if isHomebrewPath(execPath) {
			return ChannelHomebrew
		}
	}

	// Fallback: check `brew --prefix daryaft` if brew is on PATH.
	prefix := brewPrefixFn()
	if prefix != "" {
		execResolved, _ := filepath.EvalSymlinks(execPath)
		prefixResolved, _ := filepath.EvalSymlinks(prefix)
		if prefixResolved != "" && execResolved != "" && strings.HasPrefix(execResolved, prefixResolved) {
			return ChannelHomebrew
		}
	}

	if strings.ToLower(info.BuiltBy) == "goreleaser" {
		return ChannelGoreleaser
	}
	if strings.ToLower(info.BuiltBy) == "source" {
		return ChannelSource
	}
	return ChannelUnknown
}

// isHomebrewPath reports whether path looks like a Homebrew-managed path.
func isHomebrewPath(path string) bool {
	for _, prefix := range []string{
		"/opt/homebrew/",
		"/usr/local/Homebrew/",
		"/usr/local/Cellar/",
		"/home/linuxbrew/",
	} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func execBrewPrefix() string {
	out, err := exec.Command("brew", "--prefix", "daryaft").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// updateCommand returns the appropriate update command for the channel.
func updateCommand(channel InstallChannel, releaseURL string) string {
	switch channel {
	case ChannelHomebrew:
		return "brew update && brew upgrade daryaft"
	default:
		if releaseURL != "" {
			return "Download the latest release from: " + releaseURL
		}
		return "Download the latest release from: https://github.com/he8um/daryaft/releases/latest"
	}
}
