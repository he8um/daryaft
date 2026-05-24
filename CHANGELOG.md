# Changelog

All notable changes to Daryaft will be documented in this file.

Daryaft uses the project versioning policy described in `docs/roadmap/versioning-policy.md`.

## [Unreleased]

### Added

- Initial project skeleton.
- Minimal Cobra CLI foundation.
- Root command placeholder for future interactive mode.
- Version command with build metadata defaults.
- Download command surface with validation and dry-run planning.
- Single URL HTTP/HTTPS downloader with safe filename selection and `.part` writes.
- Structured single URL downloader events for started, progress, completed, and failed states.
- CLI text progress output backed by downloader events.
- Sequential batch downloading for multiple URL args, URL files, and combined inputs.
- Batch summary output with continue-on-error failure reporting.
- Starter documentation, CI, Makefile, and future packaging configuration.

### Planned

- Concurrent batch downloader engine.
- Rich progress bars, resume, and retry execution.
- Terminal UI foundation.
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
