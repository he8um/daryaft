# Update Check QA Checklist

Manual validation checklist for `daryaft update --check`.

## Purpose

Verify that `daryaft update --check` behaves correctly across install channels,
output modes, and edge cases. This command is read-only: it must never download
files, replace binaries, or modify install state.

## Prerequisites

- Daryaft installed (Homebrew, source build, or binary archive).
- Network access to `api.github.com` for real smoke tests.
- `jq` optional but useful for JSON validation.

---

## Command Smoke Tests

Run each command and verify the expected behavior described below.

### `daryaft update --help`

```bash
daryaft update --help
```

Expected:
- Shows description, usage, flags, and examples.
- Flags shown: `--check`, `--json`, `--include-prerelease`.
- `--repo` flag is hidden (must not appear in help).
- Long description states the command is read-only.
- Auto-update is described as not yet implemented.
- Examples include `--check`, `--check --json`, `--check --include-prerelease`.

### `daryaft update` (without `--check`)

```bash
daryaft update
echo "Exit: $?"
```

Expected:
- Exits non-zero (exit code 1).
- Prints a clear error: `auto-update is not implemented`.
- Suggests running `daryaft update --check`.
- Does not download anything.
- Does not modify any files.

### `daryaft update --check`

```bash
daryaft update --check
echo "Exit: $?"
```

Expected:
- Exits 0.
- Prints `Daryaft update check` as the header.
- Shows `Current version:` with the installed version.
- Shows `Latest stable:` with the latest GitHub release version.
- Shows `Status:` as one of: `up to date`, `update available`, `development build`.
- Shows `Release:` with a GitHub release URL.
- Shows `Install channel:` with the detected channel.
- Shows `Update command:` with channel-appropriate guidance:
  - Homebrew: `brew update && brew upgrade daryaft`
  - GoReleaser binary: archive download URL
  - Source build: `git pull && go build .` guidance
  - Unknown: GitHub release URL
- If update is available: shows `A new version is available.`
- If development build: shows `Note: development builds may be ahead...`
- Does not download anything.
- Does not replace the binary.
- Does not write any files.

### `daryaft update --check --json`

```bash
daryaft update --check --json
echo "Exit: $?"
```

Expected:
- Exits 0.
- Prints valid JSON to stdout.
- JSON is parseable:

```bash
daryaft update --check --json | jq .
```

- JSON contains all required fields:

```json
{
  "current_version": "...",
  "latest_version": "...",
  "update_available": false,
  "development_build": false,
  "include_prerelease": false,
  "install_channel": "...",
  "release_url": "...",
  "update_command": "...",
  "message": "..."
}
```

- `update_available`, `development_build`, `include_prerelease` are JSON booleans
  (not strings).
- `release_url`, `update_command`, `message` are present even if empty (`""`),
  not omitted.
- `install_channel` is one of: `homebrew`, `goreleaser`, `source`, `unknown`.

### `daryaft update --check --include-prerelease`

```bash
daryaft update --check --include-prerelease
echo "Exit: $?"
```

Expected:
- Exits 0.
- Same format as `--check`.
- If a pre-release exists that is newer than the latest stable, it will appear
  as the latest version.
- Stable-only run and pre-release run may return different `latest_version`
  values when pre-releases exist.
- Does not download anything.

---

## Install-Channel Validation

### Homebrew-installed binary

```bash
# Install or upgrade via Homebrew
brew tap he8um/tap
brew install he8um/tap/daryaft   # or: brew upgrade daryaft

# Verify channel detection
daryaft version
daryaft update --check
daryaft update --check --json | jq .install_channel
```

Expected:
- `daryaft version` shows `built by: goreleaser`.
- `update --check` shows `Install channel: homebrew`.
- JSON `install_channel` is `"homebrew"`.
- `Update command:` shows `brew update && brew upgrade daryaft`.

### Source build

```bash
go run . update --check
go run . update --check --json | jq .install_channel
```

Expected:
- `version` shows `built by: source`.
- `install_channel` is `"source"`.
- `Update command:` contains `git pull` and `go build`.

### Binary archive (GoReleaser)

```bash
# Download and extract a release archive, then run directly:
./daryaft update --check
./daryaft update --check --json | jq .install_channel
```

Expected:
- `version` shows `built by: goreleaser`.
- `install_channel` is `"goreleaser"` (unless binary is under a Homebrew path).
- `Update command:` contains the GitHub release archive URL.

---

## Safety Verification

After any `update --check` invocation, confirm no side effects:

```bash
# No files were created or modified in the current directory
ls -la

# No files were modified in the download directory
ls ~/Downloads 2>/dev/null | tail -5

# Binary was not replaced
daryaft version
```

Expected:
- No new files.
- Version unchanged.
- `update --check` is strictly read-only.

---

## Network Failure Behavior

```bash
# Simulate no network by using an invalid repo override (hidden flag)
daryaft update --check --repo he8um/nonexistent-repo-xyz
echo "Exit: $?"
```

Expected:
- Exits non-zero.
- Prints a clear error message referencing the GitHub API.
- Does not crash or panic.
- No files written.

```bash
# JSON error output
daryaft update --check --repo he8um/nonexistent-repo-xyz --json
echo "Exit: $?"
```

Expected:
- Exits non-zero.
- Error printed to stderr as JSON `{"error": "..."}`.

---

## Known Limitations

- **Auto-update is not implemented.** `daryaft update` without `--check` exits
  non-zero. There is no binary replacement, download, or install path.
- **No TUI update flow.** The terminal UI does not include an update screen.
  Update check is CLI-only.
- **No package-manager automation from the command.** `daryaft update --check`
  never invokes `brew upgrade` or any other package manager. It only suggests
  the appropriate command.
- **No authentication required or supported.** The GitHub Releases API is
  queried without a token. Unauthenticated requests are subject to GitHub's
  rate limits (60 requests/hour per IP).
- **Timeout.** The API call has a 10-second timeout. On slow or unavailable
  networks, the command exits with a clear error after the timeout.

---

## References

- [Self-Update Roadmap](../roadmap/self-update.md)
- [v1.2.0 Scope](../roadmap/v1.2.0-update-ux.md)
- [Homebrew Tap](homebrew-tap.md)
- [Installation](../installation.md)
