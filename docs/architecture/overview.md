# Architecture Overview

Daryaft should stay small and testable. The current implementation has the CLI
foundation, version package, default metadata, validation and planning packages,
and a first single URL downloader engine.
The first Bubble Tea TUI foundation is also implemented for no-argument
startup.

## CLI Layer

`cmd/` owns Cobra commands, flags, help text, and command dispatch. It should not
contain downloader business logic.

## Downloader Engine

Partly implemented. The current engine performs one HTTP/HTTPS GET, accepts
HTTP 2xx responses, chooses a safe filename, writes to a `.part` file, stores
temporary `.part.daryaft.json` metadata, resumes with HTTP Range when safe, and
renames the partial file on success. The single URL path emits structured
lifecycle and progress events, including retrying, resuming, and restarting
events. Multiple URLs are supported through a simple sequential batch runner
that invokes the single URL downloader one item at a time and collects per-item
results. Batch concurrency, checksum validation, and richer file conflict
behavior are planned.

## Event System

Partly implemented. The downloader layer defines typed events for started,
progress, resuming, restarting, warning, retrying, completed, and failed states.
The current CLI consumes these events to print simple line-based progress,
resume, restart, and retry messages for single and sequential batch downloads.
The same event boundary is intended to support future TUI rendering and
structured automation output without coupling those interfaces to downloader
internals.

## TUI Renderer

Partly implemented. Running `daryaft` with no arguments opens a Bubble Tea home
screen with Lip Gloss styling, simple menu navigation, help, version, planned
download screens, and clean quit handling. It does not start downloads yet.
Future TUI download screens will subscribe to downloader events and display
richer progress, queue state, errors, and history.

## Updater

Planned. The updater will check release metadata, download selected binaries,
verify checksums, and replace the executable safely.

## Config

Planned. Configuration should merge defaults, config files, environment
variables, and flags in a predictable order.

## Storage

Partly implemented. The downloader writes temporary partial metadata sidecars
for interrupted downloads. Persistent history, queue state, and lockfiles are
planned. The current sequential batch runner does not persist queue state.

Related docs:

- [Command Reference](../command-reference.md)
- [Configuration](../configuration.md)
- [Roadmap](../roadmap/index.md)
