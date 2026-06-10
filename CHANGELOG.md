# Changelog

All notable changes to Daryaft will be documented in this file.

Daryaft uses the project versioning policy described in `docs/roadmap/versioning-policy.md`.

## [Unreleased]

Post-v1.7.0 development begins. Source default version advanced to `1.8.0-dev`.

### Added

- Added `--checksum-file <path>` for per-target batch checksum verification.
  The manifest maps one `<algorithm>:<hex> <url>` entry to each download target.
- Added checksum manifest parsing (`checksum.ParseManifestFile`) with
  line-numbered validation errors, comment/blank-line handling, and
  duplicate-URL rejection.
- Added per-target checksum mapping to the download plan
  (`TargetChecksums`, `HasChecksumFile`).
- Added batch checksum verification inside the downloader for both single-target
  and per-target checksums, with checksum mismatches counted as failed items.
- Added TUI checksum status display for checksum-backed downloads: `Checksum OK`
  and `Checksum Failed` in the queue and a `Checksum verified: N` summary line.

### Changed

- Batch download summaries now include a `Checksum verified: N` count when
  applicable.
- `--checksum` and `--checksum-file` are mutually exclusive.

### Known limitations

- Checksum file URL matching is exact; URLs are not normalized.
- GNU `sha256sum` file compatibility and checksum auto-discovery are not
  implemented.
- TUI batch checksum input forms are not implemented; the TUI displays results
  only.
- Signature, PGP, and attestation verification remain out of scope.

## [1.7.0] - 2026-06-10

### Added

- Added `--checksum algorithm:hex` flag for single-target CLI download verification.
- Added SHA-256 and SHA-512 checksum verification after download completion.
- Added checksum format validation before the network request starts.
- Added checksum dry-run display in download plan output.
- Added `docs/features/checksum-verification.md` feature guide.
- Added `docs/operations/checksum-verification-qa.md` manual QA checklist.
- Added `SupportedAlgorithms()` to internal checksum package.

### Changed

- Root URL download mode now routes `--checksum` to CLI download mode instead of TUI.

### Known limitations

- Per-file batch checksums are not supported yet.
- TUI checksum entry is not implemented.
- Signature, PGP, and attestation verification are out of scope.

## [1.6.0] - 2026-06-10

### Changed

- Improved CLI download completion messages with final size and elapsed time (`Completed: file.zip (512 B in 1.2s)`).
- Added clearer CLI failure, resume, and restart messages for single-download event output.
- Renamed unfinished batch summary items from `Skipped` to `Not started`.
- Improved TUI queue visibility with per-item status history.
- Added clearer TUI queue status markers for color and no-color modes.
- Improved TUI post-run hint and summary labels.
- Expanded manual QA documentation for download and TUI queue UX.

## [1.5.0] - 2026-06-10

### Added

- Added optional `DARYAFT_USERNAME` / `DARYAFT_PASSWORD` fallback for Basic Auth
  in CLI download and inspect flows. CLI flag values always take priority. Same
  validation applies: `DARYAFT_PASSWORD` without a username fails with a clear
  error. Credentials are redacted in all output.
- Added `scripts/release-preflight.sh` and `make release-preflight` guardrail
  to validate a target version before tagging. Detects version skips, missing
  release notes, missing CHANGELOG entries, and pre-existing tags/releases.

### Changed

- Improved HTTP customization dry-run auth display to show `username:[REDACTED]`
  instead of only `[REDACTED]`, making it clear which account will be used.
- Improved invalid header error: `invalid header "...": expected "Name: Value" format`.
- Improved unsupported proxy scheme error:
  `invalid proxy "...": unsupported scheme "..."; supported schemes are http and https`.
- Updated `docs/features/http-request-customization.md` with env variable docs,
  corrected dry-run auth format, and security guidance.
- Updated `docs/operations/http-customization-qa.md` with env credential QA
  sections and corrected error message examples.

### Security

- Expanded sensitive-header redaction coverage: `Set-Cookie`, `X-Token`, header
  names containing `secret` or `password`, case-insensitive `Authorization`,
  username preservation, and proxy-authorization now all have explicit tests.
- Confirmed env-sourced Basic Auth passwords are redacted in dry-run and verbose
  output.

## [1.4.0] - 2026-06-07

### Added

- `daryaft update --check` error messages improved: 403 Forbidden now reports a
  clear rate-limit message; 404 reports repository not found; 5xx reports
  temporary unavailability; network timeout and connectivity failures each have
  distinct messages.
- `APIBaseURL` field on `update.CheckOptions` and hidden `--api-base-url` CLI
  flag for test injection, eliminating all real GitHub API calls from the default
  test suite.
- `DARYAFT_RUN_NETWORK_TESTS=1` opt-in guard for real-network update smoke tests
  (`TestUpdateCommand_RealNetwork_Check`, `TestUpdateCommand_RealNetwork_JSON`).
- `go test ./...` and `make rc-check` are now fully deterministic and will not
  fail due to GitHub API rate limits.
- Error-classification tests for 403, 404, 500, and invalid-JSON responses in
  `cmd/update_test.go`.
- `TestUpdateCommand_HiddenAPIBaseURLFlag` verifying the new hidden test flag.
- Post-1.3.0 development begins. Source default version advanced to `1.4.0-dev`.

## [1.3.0] - 2026-06-07

> **Note:** `v1.3.0` was tagged but no GitHub Release was published. The
> changes shipped as part of `v1.4.0`. Do not backfill a `v1.3.0` GitHub
> Release. See [versioning policy](docs/roadmap/versioning-policy.md).

### Added

- Post-1.2.0 development begins. Source default version advanced to `1.3.0-dev`.
- Added HTTP request customization for CLI download and inspect flows: `--proxy`,
  repeatable `--header "Name: Value"`, `--user-agent`, and `--username`/`--password`
  Basic Auth. Apply to root URL download mode, `download`, `inspect`, and batch
  downloads.
- Added redacted dry-run and verbose display for sensitive HTTP options (`[REDACTED]`
  for passwords and headers matching authorization, cookie, token, key, and related patterns).
- Added `internal/httpopts` package for shared HTTP options parsing, validation,
  redaction, request application, and proxy transport.
- Added `docs/features/http-request-customization.md`: feature guide with examples,
  validation rules, and security warnings.
- Added `docs/operations/http-customization-qa.md`: manual QA checklist for
  HTTP request customization.
- Added `docs/roadmap/v1.3.0-http-customization.md`: scope, quality gates, and
  future track.
- Updated `docs/command-reference.md`, `docs/usage.md`, `docs/index.md`,
  `docs/operations/manual-qa.md`, `README.md`, and `CHANGELOG.md`.

## [1.2.0] - 2026-06-07

### Added

- Post-1.1.0 development begins. Source default version advanced to `1.2.0-dev`.
- Added `scripts/update-homebrew-formula.sh`: a safe Homebrew formula update
  helper for future release maintenance. Fetches checksums from a GitHub
  release, updates a local tap checkout, and never pushes or commits
  automatically. Includes dry-run support and `make homebrew-formula-update`
  targets.
- `daryaft update --check` now shows `Release:` URL in all status modes, not
  only when an update is available.
- Per-channel update guidance: Homebrew → `brew update && brew upgrade daryaft`;
  source → `git pull && go build .`; goreleaser → archive download URL;
  unknown → GitHub release URL.
- `update --check` output adds a clear `A new version is available.` notice when
  an update is detected, and a development build note for pre-release versions.
- JSON output contract hardened: all fields always present (empty strings, not
  omitted), boolean fields are real JSON booleans, `update_command` reflects
  the correct channel guidance.
- Install-channel detection hardened: `source` BuiltBy takes fast path;
  goreleaser BuiltBy with Homebrew Cellar path → `homebrew`; goreleaser +
  non-Homebrew path → `goreleaser`; `brew --prefix daryaft` fallback for
  symlinked installs. All detection paths covered by injected-function unit
  tests (no real OS calls).
- Added `result_test.go` with comprehensive coverage of `Format()`, `Status()`,
  `updateCommand()`, `detectInstallChannel()`, and JSON contract.
- Added `docs/operations/update-check-qa.md`: manual QA checklist for all
  `update --check` commands, channels, and edge cases.
- Added `docs/roadmap/v1.2.0-update-ux.md`: v1.2.0 scope and quality gates.
- Updated `docs/command-reference.md`, `docs/usage.md`, `docs/roadmap/self-update.md`
  to reflect per-channel guidance and v1.2.0 scope.

## [1.1.0] - 2026-06-06

### Added

- Post-1.0 development begins. Source default version advanced to `1.1.0-dev`.
- Homebrew tap is live at `he8um/homebrew-tap`. Install with
  `brew tap he8um/tap && brew install daryaft`. Formula is manually maintained
  at `v1.0.0`; GoReleaser `brews:` publishing remains disabled. See
  `docs/operations/homebrew-tap.md` for details.
- `daryaft update --check`: read-only update check against the GitHub Releases
  API. Reports current version, latest stable release, update availability, and
  install-channel-aware upgrade instructions. `--json` for machine-readable
  output. `--include-prerelease` to include pre-release versions.
  `daryaft update` without `--check` exits non-zero (auto-update not yet
  implemented). See `docs/roadmap/self-update.md`.

## [1.0.0] - 2026-06-06

### Added

- Initial project skeleton.
- Minimal Cobra CLI foundation.
- Bubble Tea interactive home screen for no-argument `daryaft`.
- Version command with `0.6.0-dev` build metadata defaults, `built_by`
  reporting, JSON output, and release ldflags compatibility.
- Download command surface with validation and dry-run planning.
- CLI `--checksum` verification for completed single URL downloads, supporting
  manual `sha256:<hex>` and `sha512:<hex>` checksum specs.
- Single URL HTTP/HTTPS downloader with safe filename selection and `.part` writes.
- Structured single URL downloader events for started, progress, completed, and failed states.
- CLI text progress output backed by downloader events.
- Sequential batch downloading for multiple URL args, URL files, and combined inputs.
- Batch summary output with continue-on-error failure reporting.
- Basic retry execution with exponential backoff for transient network and server failures.
- Retry events and CLI retry messages.
- Real single URL resume support with `.part` files, `.part.daryaft.json`
  metadata sidecars, HTTP Range requests, safe server fallback restarts, and
  retry continuation from partial files.
- First TUI foundation with Lip Gloss styling, home menu navigation, help,
  version, and clean quit handling.
- TUI URL and `.txt` file input forms that validate existing download inputs
  and show dry-run plans without starting downloads.
- TUI download execution screen that starts real single URL or sequential batch
  downloads from the plan screen and consumes the downloader event stream.
- TUI output directory input between source entry and plan review, defaulting to
  the effective output directory and falling back to `~/Downloads`.
- Optional TUI custom filename input for single URL downloads. Empty filename
  input means auto-detect, and `.txt` batch downloads keep per-item auto-detect.
- Optional TUI checksum input for single URL downloads, using the existing
  `sha256:<hex>` and `sha512:<hex>` checksum validation.
- Context-aware downloader cancellation with cancelled events, preserved
  partial files and metadata, no retry after cancellation, and TUI `q`
  cancellation from the progress screen.
- Injectable TUI execution runner for deterministic plan and cancellation tests
  without changing user-facing TUI behavior.
- YAML user configuration with `config path`, `config show`, and `config init`
  commands plus download directory, retry, resume, no-color, and no-TUI
  defaults.
- `DARYAFT_*` environment variable overrides between CLI flags and config file
  values.
- Built-in `~/Downloads` output default when CLI flags, environment variables,
  and config do not set an output directory; explicit `.` still means the
  current directory.
- `make rc-check` for release-candidate validation including blocking
  `govulncheck` and `gosec` checks.
- `make rc-info` for printing local release-candidate tag and version metadata
  before RC validation.
- Config management commands for reading, setting, resetting, and listing
  supported config keys.
- `inspect` command for HTTP/HTTPS URL metadata preflight, with human and JSON
  output and no file writes.
- Read-only TUI Inspect URL flow backed by the shared inspect package and an
  injectable runner for deterministic tests.
- Shell completion generation for bash, zsh, fish, and PowerShell, with config
  key completion for `config get` and `config set`.
- `doctor` diagnostics command for local runtime, config, download directory,
  terminal, optional tool, and skipped release-check reporting.
- JSON output mode for `doctor` diagnostics with stable check status and summary
  fields for automation.
- Strict doctor mode for CI, where warnings cause a non-zero exit status without
  being converted into failures in the report.
- Local `make release-check` target for GoReleaser snapshot validation without
  publishing.
- GoReleaser v2 configuration with non-deprecated archive fields and
  `0.6.0-dev-SNAPSHOT-<short-commit>` snapshot naming.
- CI GoReleaser config validation with `goreleaser check`, without publishing
  releases.
- Hardened pre-release CI with Linux/macOS Go test-build matrix, tidy check,
  TUI race test, and local `make ci`.
- CI lint and security gates for `golangci-lint`, `govulncheck`, and `gosec`.
- CI lint/security jobs use a newer Go toolchain for quality-tool installation
  while the Go test/build matrix remains tied to the module Go version.
- Updated GitHub Actions workflow actions to newer Node 24-compatible majors.
- GitHub issue templates, pull request template, and branch protection
  documentation.
- Pre-release readiness review for the `0.6.0-dev` internal/manual validation
  milestone.
- QA results record for the completed `0.6.0-dev` internal validation readiness
  pass.
- `v0.6.0-rc.2` GitHub pre-release published (source/tag-only, marked
  pre-release, not stable).
- Post-RC2 release status document, v1.0.0 readiness roadmap, and doc sync.
- Clean install validation record for `v0.6.0-rc.2` — PASS WITH NOTES.
- v1.0.0 release notes draft and binary asset strategy documented; binary
  archives (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64) and
  checksums.txt will be attached to the v1.0.0 GitHub release.
- v1.0.0 go/no-go checklist and release execution plan documented; version
  policy confirmed (source default unchanged, GoReleaser ldflags inject release
  version from tag); GoReleaser publish path analyzed and Option A (build
  locally with `--skip=publish`, publish via `gh`) recommended as safest path.
- Clarified v1.0.0 release strategy: stable baseline release of current
  feature set; post-1.0 features (Windows, self-update, proxy/auth,
  concurrency, package managers, etc.) are explicitly deferred and not
  blockers.
- Internal `v0.6.0-rc.2` validation docs and release notes; supersedes
  `v0.6.0-rc.1`.
- Internal `v0.6.0-rc.1` validation docs and release notes (historical).
- Safe CLI Ctrl+C/SIGTERM cancellation for single and batch downloads, with
  partial files and metadata preserved for resume.
- Local `make lint` and `make security` quality gates with a practical
  GolangCI-Lint profile plus govulncheck/gosec scanning.
- Starter documentation, CI, Makefile, and future packaging configuration.

### Fixed

- Reviewed checksum QA coverage, strengthened manual QA docs, and confirmed
  narrow `gosec` `#nosec` suppressions remain justified.
- Restored CI `govulncheck ./...` to blocking after upgrading CI to Go `1.26.4`
  or newer, resolving the previous standard-library advisory gap. Both
  `govulncheck` and `gosec` are blocking in CI and in `make rc-check`.
- Expanded downloader HTTP integration coverage for redirects, unknown-length
  responses, slow-stream cancellation, and exhausted retry failures, and cleaned
  docs so planned update/install/Windows support is not overstated.
- Added responsive TUI sizing from terminal window dimensions, including
  adaptive input widths for narrow, normal, and wide terminals.
- Fixed TUI Backspace behavior so it edits non-empty text inputs and navigates
  back only when the current input is empty.
- Implemented CLI `--verbose` download diagnostics for URL, output, filename,
  HTTP status, target path, resume/retry decisions, and completion duration.
- Made config `theme` truthful by supporting `default` and `mono`; `mono`
  uses the same no-color TUI styling path. `animations` and `hyperlinks`
  remain stored reserved fields.
- Reworked downloader resume/restart response ownership to avoid in-place
  `http.Response` mutation and clarify body close responsibility.
- Made retry backoff cancellation context-aware without leaving sleeper
  goroutines behind.
- Replaced the default total HTTP client timeout with transport phase timeouts
  so large downloads are not killed by a fixed overall duration.
- Saved YAML config files through same-directory temporary files plus rename
  with private `0600` permissions.
- Tightened Daryaft-created config, metadata, and output directories to
  `0750` and documented narrow local gosec suppressions for intentional file
  reads.
- Hardened output-directory traversal checks with normalized relative paths.
- Added a retry upper bound of `20` across CLI planning, config values, and
  environment overrides.

### Post-1.0 Roadmap

- Concurrent batch downloader engine.
- Rich progress bars.
- Queue persistence.
- Package manager publishing (Homebrew, deb, rpm, Arch) from v1.0.0 onward.
- Windows official support and CI.
- Self-update mechanism.
- Proxy, custom headers, and authentication.
- Checksum file auto-discovery and signed checksum verification.

## [v0.1.0-dev] - 2026-05-23

### Added

- Development version for the first local CLI foundation.
- `daryaft --help`.
- `daryaft version`.
- No-argument placeholder for planned interactive mode.

### Release policy

This version is a pre-1.0 development version and must not be promoted as a
public installable stable release.
