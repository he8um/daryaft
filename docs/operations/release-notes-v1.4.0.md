# Daryaft v1.4.0

> **Status: RELEASED** — Stable release. Tag: `v1.4.0`.

## Summary

Daryaft `v1.4.0` is a reliability and test-determinism release. It makes
`go test ./...` and `make rc-check` fully deterministic by eliminating all
real GitHub API calls from the default test suite. It also improves the
user-facing error messages for the `daryaft update --check` command when the
GitHub API is unavailable, rate-limited, or returns unexpected responses.

No new end-user features are added. No new dependencies are introduced.

The release ships binary archives for linux/amd64, linux/arm64, darwin/amd64,
and darwin/arm64. The Homebrew formula at `he8um/homebrew-tap` has been updated
to `v1.4.0`.

## What's New

### Deterministic Test Suite

`go test ./...` no longer calls the real GitHub Releases API under any
circumstances. All update-check tests drive a local `httptest.Server` via the
new hidden `--api-base-url` flag and the `APIBaseURL` field on
`update.CheckOptions`. Real-network update tests are preserved but require the
opt-in environment variable:

```bash
DARYAFT_RUN_NETWORK_TESTS=1 go test ./cmd -run TestUpdateCommand_RealNetwork
```

### Deterministic `make rc-check`

`make rc-check` is now safe from GitHub API rate-limit failures. Previously,
running the full test suite without `-short` could trigger real GitHub API calls
and fail with `403 Forbidden` in CI or local environments with exhausted rate
limits. This is fully resolved.

### Improved Update Check Error Messages

`daryaft update --check` now reports clear, actionable error messages for each
failure category:

| Condition | Error Message |
|-----------|---------------|
| GitHub API rate limit / 403 | "GitHub API returned 403 Forbidden, likely due to rate limiting or access restrictions. Try again later" |
| Repository not found / 404 | "GitHub release repository was not found (404 Not Found)" |
| GitHub server error / 5xx | "GitHub API is temporarily unavailable (500 Internal Server Error). Try again later" |
| Request timeout | "request to GitHub timed out. Check your network and try again" |
| Network unreachable | "could not reach GitHub Releases API. Check your network and try again" |
| Invalid API response | "GitHub API returned an invalid response" |

### Error Classification Tests

New deterministic tests in `cmd/update_test.go` cover:

- `TestUpdateCommand_403Error` — rate-limit message
- `TestUpdateCommand_404Error` — not-found message
- `TestUpdateCommand_500Error` — unavailability message
- `TestUpdateCommand_InvalidJSONError` — invalid response message
- `TestUpdateCommand_403_JSONOutput` — `--json` error output is valid JSON
- `TestUpdateCommand_HiddenAPIBaseURLFlag` — hidden flag exists and is hidden

### Hidden Test Injection Flag

A hidden `--api-base-url` flag is available on `daryaft update` for test
injection only. It does not appear in `--help` output and is not intended for
end-user use.

## User-Facing Behavior

- `daryaft update --check` remains read-only. It never downloads, installs, or
  replaces the current binary.
- All existing CLI flags (`--check`, `--json`, `--include-prerelease`) are
  unchanged.
- Error messages are now clearer when GitHub is unavailable or rate-limited.
- No GitHub authentication or token support has been added.

## Known Limitations

- Auto-update is not implemented. Use `daryaft update --check` to check for
  updates and the appropriate command for your install channel to upgrade.
- GitHub auth/token support is not implemented. Unauthenticated API calls are
  subject to GitHub's public rate limits (60 requests/hour per IP).
- Real-network update smoke can still hit GitHub rate limits when run manually
  without `DARYAFT_RUN_NETWORK_TESTS=1`.
- GoReleaser Homebrew publishing remains disabled. The Homebrew formula is
  manually maintained at `he8um/homebrew-tap`.
- Windows is not officially supported.

## Upgrade Instructions

### Homebrew

```bash
brew update && brew upgrade daryaft
daryaft version
```

### GitHub Binary

Download the `v1.4.0` archive for your platform from
[GitHub Releases](https://github.com/he8um/daryaft/releases/tag/v1.4.0),
extract, and replace your existing binary.

### Source

```bash
git pull && go build .
./daryaft version
```

## Validation Summary

- `go test ./...` — PASS (deterministic, no real GitHub calls)
- `go test -short ./...` — PASS
- `go build ./...` — PASS
- `go test -race ./internal/downloader` — PASS
- `go test -race ./internal/tui` — PASS
- `make lint` — PASS (0 issues)
- `make security` — PASS (govulncheck: no vulnerabilities; gosec: 0 issues)
- `goreleaser check` — PASS
- `make rc-check` — PASS (deterministic)
- GoReleaser artifact build — PASS
- Extracted binary version: `1.4.0`, `built_by: goreleaser`
- Checksum verification — PASS
- Homebrew formula updated to `v1.4.0`
- Homebrew install/upgrade validated

## Security Behavior

HTTP request customization redaction (from v1.3.0) is unchanged:

- `--password` values are always displayed as `[REDACTED]`
- Headers matching authorization, cookie, token, key, and related patterns are
  displayed as `[REDACTED]`
- No credentials appear in error messages, logs, or dry-run output
