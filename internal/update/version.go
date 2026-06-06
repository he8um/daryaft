package update

import (
	"fmt"
	"strconv"
	"strings"
)

// semver holds the numeric parts of a parsed version string.
type semver struct {
	major, minor, patch int
}

// parseVersion parses "v1.2.3", "1.2.3", "1.2.3-foo", etc.
// It returns the semver and true if the version is a clean stable release
// (no pre-release suffix). It returns ok=true for any parseable version;
// stable=false means a pre-release or dev suffix was present.
func parseVersion(v string) (sv semver, stable bool, ok bool) {
	v = strings.TrimPrefix(v, "v")

	// Check for pre-release suffix (anything after a hyphen).
	base := v
	if idx := strings.IndexByte(v, '-'); idx >= 0 {
		base = v[:idx]
		stable = false
	} else {
		stable = true
	}

	parts := strings.SplitN(base, ".", 3)
	if len(parts) != 3 {
		return semver{}, false, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return semver{}, false, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return semver{}, false, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return semver{}, false, false
	}
	return semver{major, minor, patch}, stable, true
}

// compareVersions returns:
//
//	-1 if a < b
//	 0 if a == b
//	+1 if a > b
func compareVersions(a, b semver) int {
	for _, pair := range [][2]int{
		{a.major, b.major},
		{a.minor, b.minor},
		{a.patch, b.patch},
	} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	return 0
}

// isDevBuild reports whether the version string looks like a development build.
func isDevBuild(v string) bool {
	v = strings.ToLower(strings.TrimPrefix(v, "v"))
	return strings.Contains(v, "dev") || strings.Contains(v, "snapshot")
}

// normalizeVersion strips a leading "v" for display consistency.
func normalizeVersion(v string) string {
	return strings.TrimPrefix(v, "v")
}

// versionTagURL returns the GitHub release tag URL.
func versionTagURL(owner, repo, tag string) string {
	return fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", owner, repo, tag)
}
