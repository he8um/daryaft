# Daryaft v1.2.0

> **Status: RELEASED** — Stable release. Tag: `v1.2.0`.

## Summary

Daryaft `v1.2.0` is a focused polish sprint on the `daryaft update --check`
feature shipped in `v1.1.0`. It improves output clarity, hardens install-channel
detection, adds per-channel upgrade guidance, and strengthens test coverage. The
command remains read-only throughout.

The release ships binary archives for linux/amd64, linux/arm64, darwin/amd64,
and darwin/arm64. The Homebrew formula at `he8um/homebrew-tap` has been updated
to `v1.2.0`.

## What's New

### Update Check UX Polish

`daryaft update --check` output now includes:

- **Release URL shown in all status modes** — the GitHub release URL is always
  shown regardless of whether an update is available or the build is a
  development build.
- **Per-channel upgrade guidance** — the `Update command:` field now gives
  install-channel-specific instructions:

  | Channel | Guidance |
  |---------|----------|
  | `homebrew` | `brew update && brew upgrade daryaft` |
  | `goreleaser` | Download the latest release archive from the GitHub release URL |
  | `source` | Pull the repository and rebuild: `git pull && go build .` |
  | `unknown` | Download the latest release from the GitHub release URL |

- **Clear update notice** — when an update is available, output ends with:
  `A new version is available. Use the update command above to upgrade.`
- **Development build note** — development builds continue to show:
  `Note: development builds may be ahead of the latest stable release.`

#### Human output examples

Up-to-date Homebrew install:

```text
Daryaft update check

Current version:  1.2.0
Latest stable:    1.2.0
Status:           up to date

Release: https://github.com/he8um/daryaft/releases/tag/v1.2.0

Install channel:  homebrew
Update command:   brew update && brew upgrade daryaft
```

Update available (goreleaser binary):

```text
Daryaft update check

Current version:  1.1.0
Latest stable:    1.2.0
Status:           update available

Release: https://github.com/he8um/daryaft/releases/tag/v1.2.0

Install channel:  goreleaser
Update command:   Download the latest release archive from: https://github.com/he8um/daryaft/releases/tag/v1.2.0

A new version is available. Use the update command above to upgrade.
```

Source development build:

```text
Daryaft update check

Current version:  1.3.0-dev
Latest stable:    1.2.0
Status:           development build

Release: https://github.com/he8um/daryaft/releases/tag/v1.2.0

Install channel:  source
Update command:   Pull the repository and rebuild: git pull && go build .

Note: development builds may be ahead of the latest stable release.
```

### JSON Contract Hardening

The JSON output contract (`--json`) is now fully stable:

- All fields always present — empty strings (`""`) are included, never omitted.
- Boolean fields (`update_available`, `development_build`, `include_prerelease`)
  are real JSON booleans, not strings.
- `update_command` reflects the install-channel-specific guidance.
- `release_url` is always populated when a release is found.

#### JSON output example

```json
{
  "current_version": "1.2.0",
  "latest_version": "1.2.0",
  "update_available": false,
  "development_build": false,
  "include_prerelease": false,
  "install_channel": "homebrew",
  "release_url": "https://github.com/he8um/daryaft/releases/tag/v1.2.0",
  "update_command": "brew update && brew upgrade daryaft",
  "message": ""
}
```

### Install-Channel Detection Hardening

- `source` BuiltBy → `ChannelSource` fast path (no path check needed).
- `goreleaser` BuiltBy + Homebrew Cellar path → `ChannelHomebrew`.
- `goreleaser` BuiltBy + non-Homebrew path → `ChannelGoreleaser`.
- Unknown BuiltBy + non-Homebrew path → `ChannelUnknown`.
- `brew --prefix daryaft` fallback for symlinked Homebrew installs.

### Test Coverage

New `internal/update/result_test.go` covers:

- `Format()` for all status modes (up-to-date, update-available, dev-build).
- `updateCommand()` for all four channels, with and without release URL.
- `detectInstallChannel()` for Cellar paths, non-Homebrew paths, BuiltBy
  values, executable error, and brew-prefix fallback.
- JSON contract: all fields present, booleans are booleans, empty strings
  present.
- `Status()` precedence.

All detection tests use injected executable-path and brew-prefix functions —
no dependency on real OS state in unit tests.

### Documentation

- Added `docs/operations/update-check-qa.md`: manual QA checklist for all
  `update --check` commands, install channels, and edge cases.
- Added `docs/roadmap/v1.2.0-update-ux.md`: v1.2.0 scope and quality gates.
- Updated `docs/command-reference.md`, `docs/usage.md`,
  `docs/roadmap/self-update.md`, `CHANGELOG.md`, and `README.md` to reflect
  per-channel guidance.
- Added `scripts/update-homebrew-formula.sh`: helper script for updating the
  local Homebrew tap formula after a GitHub release is published.

## Upgrade Instructions

### Homebrew

```bash
brew update
brew upgrade daryaft
daryaft version
daryaft update --check
```

### GitHub Binary Archives

Download `v1.2.0` assets from the
[v1.2.0 GitHub release](https://github.com/he8um/daryaft/releases/tag/v1.2.0):

```bash
# Example: macOS Apple Silicon
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.2.0/daryaft_darwin_arm64.tar.gz
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.2.0/checksums.txt
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

### Source Build

```bash
git pull
go build .
./daryaft version
```

## Known Limitations

- **Auto-update not implemented.** `daryaft update` without `--check` exits
  non-zero. A future release will implement in-place binary upgrade.
- **TUI update screen not implemented.** The terminal UI has no update screen.
  Update check is CLI-only.
- **Proxy, custom headers, and authentication not implemented.** Planned as the
  next major feature track after this polish sprint.
- **Windows not officially supported.** No Windows binaries are published.
- **GoReleaser Homebrew publishing remains disabled.** The `brews:` block in
  `.goreleaser.yml` is still commented out. The Homebrew formula is manually
  maintained, assisted by `scripts/update-homebrew-formula.sh`.
- **No signed checksums.** `checksums.txt` contains plain SHA-256 hashes.

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
| Extracted binary version | `1.2.0` |
| `built_by` | `goreleaser` |
| Checksum verification | PASS |
| `daryaft update --check` smoke | PASS |
| `daryaft doctor` smoke | PASS |
| Homebrew formula update | PASS |
| `brew upgrade daryaft` | PASS |

## References

- [Changelog](../../CHANGELOG.md)
- [v1.2.0 Scope](../roadmap/v1.2.0-update-ux.md)
- [Update Check QA](update-check-qa.md)
- [Self-Update Roadmap](../roadmap/self-update.md)
- [Homebrew Tap](homebrew-tap.md)
- [v1.1.0 Release Notes](release-notes-v1.1.0.md)
- [Release Assets Strategy](release-assets.md)
- [Post-1.0 Feature Packs](../roadmap/post-1-feature-packs.md)
