# Changelog

All notable changes to Daryaft will be documented in this file.

Daryaft uses the project versioning policy described in `docs/roadmap/versioning-policy.md`.

## [Unreleased]

### Added

- Initial project skeleton.
- Minimal Cobra CLI foundation.
- Bubble Tea interactive home screen for no-argument `daryaft`.
- Version command with `0.5.0-dev` build metadata defaults, `built_by`
  reporting, JSON output, and release ldflags compatibility.
- Download command surface with validation and dry-run planning.
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
- TUI output directory input between source entry and plan review. Empty output
  means the current directory.
- Optional TUI custom filename input for single URL downloads. Empty filename
  input means auto-detect, and `.txt` batch downloads keep per-item auto-detect.
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
- Config management commands for reading, setting, resetting, and listing
  supported config keys.
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
  `0.5.0-dev-SNAPSHOT-<short-commit>` snapshot naming.
- CI GoReleaser config validation with `goreleaser check`, without publishing
  releases.
- Hardened pre-release CI with Linux/macOS Go test-build matrix, tidy check,
  TUI race test, and local `make ci`.
- Updated GitHub Actions workflow actions to newer Node 24-compatible majors.
- GitHub issue templates, pull request template, and branch protection
  documentation.
- Starter documentation, CI, Makefile, and future packaging configuration.

### Fixed

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
