package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/pkg/version"
)

const (
	defaultOwner   = "he8um"
	defaultRepo    = "daryaft"
	checkTimeout   = 10 * time.Second
	maxReleaseBody = 1 << 20 // 1 MiB
)

// CheckOptions controls the update check.
type CheckOptions struct {
	// Owner and Repo default to he8um/daryaft.
	Owner string
	Repo  string

	// IncludePrerelease includes pre-release versions when true.
	IncludePrerelease bool

	// Client is the HTTP client used for the GitHub API call.
	// If nil, a client with checkTimeout is used.
	Client *http.Client

	// CurrentVersion overrides version.Version for the comparison.
	// Empty string means use version.Info().Version.
	CurrentVersion string
}

// githubRelease is a minimal representation of a GitHub Releases API entry.
type githubRelease struct {
	TagName     string `json:"tag_name"`
	HTMLURL     string `json:"html_url"`
	Prerelease  bool   `json:"prerelease"`
	Draft       bool   `json:"draft"`
	PublishedAt string `json:"published_at"`
}

// Check performs a read-only update check against the GitHub Releases API.
func Check(ctx context.Context, opts CheckOptions) (Result, error) {
	if opts.Owner == "" {
		opts.Owner = defaultOwner
	}
	if opts.Repo == "" {
		opts.Repo = defaultRepo
	}
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: checkTimeout}
	}

	info := version.Info()
	currentRaw := opts.CurrentVersion
	if currentRaw == "" {
		currentRaw = info.Version
	}

	currentDisplay := normalizeVersion(currentRaw)
	isDev := isDevBuild(currentRaw)

	latest, err := fetchLatestRelease(ctx, client, opts.Owner, opts.Repo, opts.IncludePrerelease)
	if err != nil {
		return Result{}, err
	}

	latestTag := latest.TagName
	latestDisplay := normalizeVersion(latestTag)
	releaseURL := latest.HTMLURL

	channel := DetectInstallChannel(info)
	var updateAvailable bool

	if !isDev {
		currentSV, currentStable, currentOK := parseVersion(currentRaw)
		latestSV, latestStable, latestOK := parseVersion(latestTag)

		if currentOK && latestOK && currentStable && latestStable {
			updateAvailable = compareVersions(latestSV, currentSV) > 0
		}
	}

	updateCmd := updateCommand(channel, releaseURL)

	var msg string
	switch {
	case isDev:
		msg = "Development build — not comparable to stable releases."
	case updateAvailable:
		msg = fmt.Sprintf("%s %s is available.", config.AppName, latestDisplay)
	default:
		msg = fmt.Sprintf("%s is up to date.", config.AppName)
	}

	return Result{
		CurrentVersion:    currentDisplay,
		LatestVersion:     latestDisplay,
		UpdateAvailable:   updateAvailable,
		DevelopmentBuild:  isDev,
		IncludePrerelease: opts.IncludePrerelease,
		InstallChannel:    channel,
		ReleaseURL:        releaseURL,
		UpdateCommand:     updateCmd,
		Message:           msg,
	}, nil
}

func fetchLatestRelease(ctx context.Context, client *http.Client, owner, repo string, includePrerelease bool) (githubRelease, error) {
	if includePrerelease {
		return fetchNewestRelease(ctx, client, owner, repo)
	}
	return fetchStableRelease(ctx, client, owner, repo)
}

// fetchStableRelease uses /releases/latest which GitHub guarantees is the
// newest non-prerelease, non-draft release.
func fetchStableRelease(ctx context.Context, client *http.Client, owner, repo string) (githubRelease, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	var release githubRelease
	if err := githubGet(ctx, client, apiURL, &release); err != nil {
		return githubRelease{}, err
	}
	return release, nil
}

// fetchNewestRelease queries /releases and returns the newest non-draft entry,
// including prereleases.
func fetchNewestRelease(ctx context.Context, client *http.Client, owner, repo string) (githubRelease, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=20", owner, repo)
	var releases []githubRelease
	if err := githubGet(ctx, client, apiURL, &releases); err != nil {
		return githubRelease{}, err
	}

	for _, r := range releases {
		if r.Draft {
			continue
		}
		return r, nil
	}
	return githubRelease{}, fmt.Errorf("no releases found for %s/%s", owner, repo)
}

func githubGet(ctx context.Context, client *http.Client, url string, dest interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", config.AppName+"/"+version.Version)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("github releases query: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github releases API returned %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxReleaseBody))
	if err != nil {
		return fmt.Errorf("read github response: %w", err)
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("parse github response: %w", err)
	}
	return nil
}
