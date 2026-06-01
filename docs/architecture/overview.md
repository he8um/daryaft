# Architecture Overview

Daryaft should stay small and testable. The current implementation has the CLI
foundation, version package, default metadata, validation and planning packages,
and a first single URL downloader engine.
The first Bubble Tea TUI foundation is also implemented for no-argument
startup.

## CLI Layer

`cmd/` owns Cobra commands, flags, help text, and command dispatch. It should not
contain downloader business logic. `pkg/version` owns build metadata defaults
and linker-injected release values; `daryaft version --json` exposes those
values for automation.

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
internals. The TUI now consumes this event boundary for its execution screen.

## TUI Renderer

Partly implemented. Running `daryaft` with no arguments opens a Bubble Tea home
screen with Lip Gloss styling, simple menu navigation, URL and `.txt` file input
forms, dry-run plan rendering, an execution/progress screen backed by downloader
events, TUI cancellation, help, version, and clean quit handling. Richer
progress, queue state, errors, and history are planned.

## Updater

Planned. The updater will check release metadata, download selected binaries,
verify checksums, and replace the executable safely.

## Config

Partly implemented. `internal/config` owns YAML defaults, platform config path
resolution through `os.UserConfigDir()`, load/save/init helpers, and config
commands. The current precedence is CLI flags, then config file values, then
built-in defaults. Environment-variable configuration and profiles are not
implemented yet.

## Storage

Partly implemented. The downloader writes temporary partial metadata sidecars
for interrupted downloads. Persistent history, queue state, and lockfiles are
planned. The current sequential batch runner does not persist queue state.

Related docs:

- [Command Reference](../command-reference.md)
- [Configuration](../configuration.md)
- [Roadmap](../roadmap/index.md)
