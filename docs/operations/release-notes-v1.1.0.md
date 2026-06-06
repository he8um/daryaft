# Daryaft v1.1.0

> **Status: RELEASED** — Stable release. Tag: `v1.1.0`.

## Summary

Daryaft `v1.1.0` adds the read-only `daryaft update --check` command, which
queries the GitHub Releases API to report whether a newer stable release is
available. It is the only change relative to `v1.0.0`.

The release ships binary archives for linux/amd64, linux/arm64, darwin/amd64,
and darwin/arm64. The Homebrew formula at `he8um/homebrew-tap` has been updated
to `v1.1.0`.

## What's New

### `daryaft update --check` (read-only update check)

```bash
daryaft update --check
daryaft update --check --json
daryaft update --check --include-prerelease
```

- Queries the GitHub Releases API (`api.github.com/repos/he8um/daryaft`).
- Compares the current binary version to the latest stable release.
- Reports: current version, latest stable version, update availability, and
  install-channel-aware upgrade instructions.
- `--json` outputs machine-readable JSON with all fields.
- `--include-prerelease` widens the search to include pre-release versions.
- **Read-only**: does not download, install, or replace the current binary.
- **Auto-update is not implemented.** `daryaft update` without `--check` exits
  non-zero with a clear message directing users to `daryaft update --check`.

#### Human output example

```text
Daryaft update check

Current version:  1.1.0
Latest stable:    1.1.0
Status:           up to date

Install channel:  homebrew
Update command:   brew update && brew upgrade daryaft
```

#### JSON output example

```json
{
  "current_version": "1.1.0",
  "latest_version": "1.1.0",
  "update_available": false,
  "development_build": false,
  "include_prerelease": false,
  "install_channel": "homebrew",
  "release_url": "",
  "update_command": "brew update && brew upgrade daryaft",
  "message": ""
}
```

#### Install-channel awareness

| Channel | Update command suggested |
|---------|--------------------------|
| `homebrew` | `brew update && brew upgrade daryaft` |
| `goreleaser` | GitHub Releases URL |
| `source` | GitHub Releases URL |
| `unknown` | GitHub Releases URL |

Homebrew users should always upgrade through Homebrew to preserve Cellar
management, symlinks, and rollback capability.

## Upgrade Instructions

### Homebrew

```bash
brew update
brew upgrade daryaft
daryaft version
daryaft update --check
```

### GitHub Binary Archives

Download `v1.1.0` assets from the
[v1.1.0 GitHub release](https://github.com/he8um/daryaft/releases/tag/v1.1.0):

```bash
# Example: macOS Apple Silicon
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.1.0/daryaft_darwin_arm64.tar.gz
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.1.0/checksums.txt
shasum -a 256 --check checksums.txt
tar -xzf daryaft_darwin_arm64.tar.gz
./daryaft version
```

Available archives:

| File | Platform |
|------|----------|
| `daryaft_linux_amd64.tar.gz` | Linux x86-64 |
| `daryaft_linux_arm64.tar.gz` | Linux ARM64 |
| `daryaft_darwin_amd64.tar.gz` | macOS Intel |
| `daryaft_darwin_arm64.tar.gz` | macOS Apple Silicon |
| `checksums.txt` | SHA-256 checksums for all archives |

## Known Limitations

- **Auto-update not implemented.** `daryaft update` without `--check` exits
  non-zero. A future release will implement in-place binary upgrade.
- **Windows not officially supported.** No Windows binaries are published.
  Windows support is planned post-1.0.
- **GoReleaser Homebrew publishing remains disabled.** The `brews:` block in
  `.goreleaser.yml` is still commented out. The Homebrew formula is manually
  maintained.
- **No signed checksums.** `checksums.txt` contains plain SHA-256 hashes.
  Signed release assets are planned for a future release.
- **Concurrent batch downloads not implemented.** Sequential batch downloading
  continues as in v1.0.0.

## Validation Summary

All required checks passed before release:

| Check | Result |
|-------|--------|
| `go test ./...` | PASS |
| `go build ./...` | PASS |
| `go test -race ./internal/downloader` | PASS |
| `go test -race ./internal/tui` | PASS |
| `make lint` (`golangci-lint`) | 0 issues |
| `make security` (`govulncheck`) | no vulnerabilities |
| `make security` (`gosec`) | Issues: 0 |
| `goreleaser check` | 1 config validated |
| `make rc-check` | PASS |
| `git diff --check` | PASS |
| GoReleaser artifact build | PASS |
| Extracted binary version | `1.1.0` |
| `built_by` | `goreleaser` |
| Checksum verification | PASS |
| `daryaft update --check` smoke | PASS |
| `daryaft doctor` smoke | PASS |
| Homebrew formula update | PASS |
| `brew upgrade daryaft` | PASS |

## References

- [Changelog](../../CHANGELOG.md)
- [Self-Update Roadmap](../roadmap/self-update.md)
- [Homebrew Tap](homebrew-tap.md)
- [v1.0.0 Release Notes](release-notes-v1.0.0.md)
- [Release Assets Strategy](release-assets.md)
- [Post-1.0 Feature Packs](../roadmap/post-1-feature-packs.md)
