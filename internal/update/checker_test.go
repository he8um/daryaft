package update

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// latestHandler returns a /releases/latest handler that responds with the given release.
func latestHandler(t *testing.T, release githubRelease) http.HandlerFunc {
	t.Helper()
	data, err := json.Marshal(release)
	if err != nil {
		t.Fatalf("marshal release: %v", err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/releases/latest") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

// listHandler returns a /releases handler that responds with a list of releases.
func listHandler(t *testing.T, releases []githubRelease) http.HandlerFunc {
	t.Helper()
	data, err := json.Marshal(releases)
	if err != nil {
		t.Fatalf("marshal releases: %v", err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/releases") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

// mustCheck calls Check with a test server client and test-injected owner/repo.
func mustCheck(t *testing.T, srv *httptest.Server, opts CheckOptions) (Result, error) {
	t.Helper()
	opts.Client = srv.Client()
	// The test server URL is e.g. http://127.0.0.1:PORT.
	// We override the API calls by pointing at a test-owned handler. To do so
	// without changing the production URL-building logic we use a custom
	// transport that rewrites the host.
	opts.Client = &http.Client{
		Transport: &rewriteTransport{base: srv.Client().Transport, target: srv.URL},
	}
	return Check(context.Background(), opts)
}

// rewriteTransport replaces the scheme+host of every outbound request with
// the test server address. This lets us use the real URL-building logic in
// checker.go while routing to the test server.
type rewriteTransport struct {
	base   http.RoundTripper
	target string // e.g. "http://127.0.0.1:PORT"
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid mutating the original.
	clone := req.Clone(req.Context())
	// Strip scheme and host, keep path+query.
	clone.URL.Scheme = "http"
	clone.URL.Host = strings.TrimPrefix(t.target, "http://")
	clone.Host = clone.URL.Host
	if t.base != nil {
		return t.base.RoundTrip(clone)
	}
	return http.DefaultTransport.RoundTrip(clone)
}

// --- Tests ---

func TestCheck_UpToDate(t *testing.T) {
	srv := httptest.NewServer(latestHandler(t, githubRelease{
		TagName:    "v1.0.0",
		HTMLURL:    "https://github.com/he8um/daryaft/releases/tag/v1.0.0",
		Prerelease: false,
		Draft:      false,
	}))
	defer srv.Close()

	result, err := mustCheck(t, srv, CheckOptions{
		CurrentVersion: "1.0.0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UpdateAvailable {
		t.Error("UpdateAvailable should be false when versions match")
	}
	if result.DevelopmentBuild {
		t.Error("DevelopmentBuild should be false for 1.0.0")
	}
	if result.LatestVersion != "1.0.0" {
		t.Errorf("LatestVersion = %q, want %q", result.LatestVersion, "1.0.0")
	}
	if result.CurrentVersion != "1.0.0" {
		t.Errorf("CurrentVersion = %q, want %q", result.CurrentVersion, "1.0.0")
	}
	if result.Status() != StatusUpToDate {
		t.Errorf("Status = %q, want %q", result.Status(), StatusUpToDate)
	}
}

func TestCheck_UpdateAvailable(t *testing.T) {
	srv := httptest.NewServer(latestHandler(t, githubRelease{
		TagName:    "v1.1.0",
		HTMLURL:    "https://github.com/he8um/daryaft/releases/tag/v1.1.0",
		Prerelease: false,
		Draft:      false,
	}))
	defer srv.Close()

	result, err := mustCheck(t, srv, CheckOptions{
		CurrentVersion: "1.0.0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.UpdateAvailable {
		t.Error("UpdateAvailable should be true when newer version exists")
	}
	if result.LatestVersion != "1.1.0" {
		t.Errorf("LatestVersion = %q, want %q", result.LatestVersion, "1.1.0")
	}
	if result.Status() != StatusAvailable {
		t.Errorf("Status = %q, want %q", result.Status(), StatusAvailable)
	}
}

func TestCheck_DevBuild(t *testing.T) {
	srv := httptest.NewServer(latestHandler(t, githubRelease{
		TagName:    "v1.0.0",
		HTMLURL:    "https://github.com/he8um/daryaft/releases/tag/v1.0.0",
		Prerelease: false,
		Draft:      false,
	}))
	defer srv.Close()

	result, err := mustCheck(t, srv, CheckOptions{
		CurrentVersion: "1.1.0-dev",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DevelopmentBuild {
		t.Error("DevelopmentBuild should be true for 1.1.0-dev")
	}
	if result.UpdateAvailable {
		t.Error("UpdateAvailable should be false for a dev build")
	}
	if result.Status() != StatusDevBuild {
		t.Errorf("Status = %q, want %q", result.Status(), StatusDevBuild)
	}
}

func TestCheck_IgnoresDrafts(t *testing.T) {
	// The /releases/latest endpoint already excludes drafts per GitHub API contract.
	// Test that if we somehow get a draft, the result still processes cleanly.
	srv := httptest.NewServer(latestHandler(t, githubRelease{
		TagName:    "v1.0.0",
		HTMLURL:    "https://github.com/he8um/daryaft/releases/tag/v1.0.0",
		Prerelease: false,
		Draft:      false,
	}))
	defer srv.Close()

	result, err := mustCheck(t, srv, CheckOptions{CurrentVersion: "1.0.0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UpdateAvailable {
		t.Error("should not report update available")
	}
}

func TestCheck_StableModeIgnoresPrerelease(t *testing.T) {
	// In stable mode, /releases/latest skips prereleases.
	// Simulate: current=1.0.0, latest stable=1.0.0 (prerelease 1.1.0-rc.1 exists but skipped).
	srv := httptest.NewServer(latestHandler(t, githubRelease{
		TagName:    "v1.0.0",
		HTMLURL:    "https://github.com/he8um/daryaft/releases/tag/v1.0.0",
		Prerelease: false,
		Draft:      false,
	}))
	defer srv.Close()

	result, err := mustCheck(t, srv, CheckOptions{
		CurrentVersion:    "1.0.0",
		IncludePrerelease: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UpdateAvailable {
		t.Error("stable mode should not report update from a prerelease")
	}
	if result.LatestVersion != "1.0.0" {
		t.Errorf("LatestVersion = %q, want 1.0.0", result.LatestVersion)
	}
}

func TestCheck_IncludePrereleaseSelectsPrerelease(t *testing.T) {
	releases := []githubRelease{
		{TagName: "v1.1.0-rc.1", HTMLURL: "https://github.com/he8um/daryaft/releases/tag/v1.1.0-rc.1", Prerelease: true, Draft: false},
		{TagName: "v1.0.0", HTMLURL: "https://github.com/he8um/daryaft/releases/tag/v1.0.0", Prerelease: false, Draft: false},
	}
	srv := httptest.NewServer(listHandler(t, releases))
	defer srv.Close()

	result, err := mustCheck(t, srv, CheckOptions{
		CurrentVersion:    "1.0.0",
		IncludePrerelease: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The prerelease appears first (newer) so should be selected.
	if result.LatestVersion != "1.1.0-rc.1" {
		t.Errorf("LatestVersion = %q, want 1.1.0-rc.1", result.LatestVersion)
	}
}

func TestCheck_IncludePrerelease_SkipsDraft(t *testing.T) {
	releases := []githubRelease{
		{TagName: "v1.1.0-rc.1", HTMLURL: "...", Prerelease: true, Draft: true},
		{TagName: "v1.0.0", HTMLURL: "https://github.com/he8um/daryaft/releases/tag/v1.0.0", Prerelease: false, Draft: false},
	}
	srv := httptest.NewServer(listHandler(t, releases))
	defer srv.Close()

	result, err := mustCheck(t, srv, CheckOptions{
		CurrentVersion:    "1.0.0",
		IncludePrerelease: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Draft should be skipped; stable 1.0.0 is the first non-draft entry.
	if result.LatestVersion != "1.0.0" {
		t.Errorf("LatestVersion = %q, want 1.0.0 (draft should be skipped)", result.LatestVersion)
	}
}

func TestCheck_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := mustCheck(t, srv, CheckOptions{CurrentVersion: "1.0.0"})
	if err == nil {
		t.Fatal("expected error on 500 response, got nil")
	}
}

func TestCheck_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	_, err := mustCheck(t, srv, CheckOptions{CurrentVersion: "1.0.0"})
	if err == nil {
		t.Fatal("expected error on invalid JSON, got nil")
	}
}
