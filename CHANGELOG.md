# Changelog

All notable changes to Daryaft will be documented in this file.

Daryaft uses the project versioning policy described in `docs/roadmap/versioning-policy.md`.

## [Unreleased]

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

### Planned

- Concurrent batch downloader engine.
- Rich progress bars.
- Queue persistence.
- Public installation channels from v1.0.0 onward.

## [v0.1.0-dev] - 2026-05-23

### Added

- Development version for the first local CLI foundation.
- `daryaft --help`.
- `daryaft version`.
- No-argument placeholder for planned interactive mode.

### Release policy

This version is a pre-1.0 development version and must not be promoted as a
public installable stable release.
