# Changelog

All notable changes to Daryaft will be documented in this file.

Daryaft uses the project versioning policy described in `docs/roadmap/versioning-policy.md`.

## [Unreleased]

### Added

- Initial project skeleton.
- Minimal Cobra CLI foundation.
- Bubble Tea interactive home screen for no-argument `daryaft`.
- Version command with build metadata defaults.
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
- Starter documentation, CI, Makefile, and future packaging configuration.

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
