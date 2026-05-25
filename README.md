# Daryaft

Daryaft is a modern terminal downloader written in Go. It is similar in spirit
to `wget`, with a Bubble Tea terminal UI foundation, clean architecture,
packaging, and future self-update support.

> Developed with <3 by AmirHesam Piri

## Status

Daryaft is in pre-1.0 development and is not stable yet. The current codebase
has a CLI foundation, dry-run planning, and real single URL HTTP/HTTPS
downloads with simple text progress output. Multiple URLs can also be
downloaded sequentially in one command. Basic retry execution is implemented
for temporary network and server failures, and interrupted downloads can resume
from `.part` files when the server supports HTTP Range requests. Running
`daryaft` with no arguments opens the interactive TUI home screen. The TUI can
collect a single URL or `.txt` file path and show the same dry-run download
plan as the CLI, then start real single or sequential batch downloads from the
plan screen.

- Repository: https://github.com/he8um/daryaft
- Website: https://xhesam.com
- Project page: https://xhesam.com/daryaft
- Author: AmirHesam Piri <info@xhesam.com>
- License: MIT

## Install Policy

Public install channels begin at `v1.0.0`. Before that release, Homebrew, Debian,
RPM, Arch, and one-line install instructions are planned only and must not be
presented as stable user-facing install paths.

For local development:

```bash
go mod download
go run . --help
go run . version
go run .
go run . https://example.com/file.zip --dry-run
go run . https://example.com/file.zip
go run . https://example.com/a.txt https://example.com/b.txt
go run . download https://example.com/file.zip --dry-run
```

Build and test locally:

```bash
make test
make build
make run
```

## Current Commands

```bash
daryaft --help
daryaft version
daryaft
daryaft https://example.com/file.zip --dry-run
daryaft https://example.com/file.zip
daryaft https://example.com/a.txt https://example.com/b.txt
daryaft -f urls.txt --dry-run
daryaft -f urls.txt
daryaft download https://example.com/file.zip --dry-run
daryaft download https://example.com/file.zip
daryaft download -f urls.txt --dry-run
daryaft download -f urls.txt
```

With no arguments, Daryaft opens a Bubble Tea home screen with menu entries for
URL input, `.txt` file input, help, version information, and quit. Download
actions inside the TUI validate input, show a dry-run plan, and can start real
downloads from that plan. TUI downloads currently use the current directory as
the output path; output path input is planned for a later milestone. Existing
CLI download commands remain fully supported.

Download commands validate URLs and can print a dry-run plan. Real downloading
is implemented for one URL and for multiple URLs sequentially. Batch downloads
continue after item failures and print a final summary. `--retries` controls
retry attempts after the initial attempt; the default `3` means up to four total
attempts.

Downloads write to `<filename>.part` first and keep sidecar state in
`<filename>.part.daryaft.json` while incomplete. `--resume` is enabled by
default. When a partial file exists, Daryaft sends `Range: bytes=<size>-` and
appends only after the server replies with `206 Partial Content`. If the server
does not support resume, or saved validators show the remote file changed,
Daryaft truncates the `.part` file and restarts from byte `0`. `--no-resume`
always restarts the partial file from byte `0`.

The downloader now emits structured started, progress, completed, and failed
events plus retrying events. The CLI consumes those events for line-based
progress and retry output for single and sequential batch downloads. The TUI
execution screen consumes the same event stream for status, byte progress,
speed, retry/resume/restart messages, completion, failure, and batch summaries.
Cancellation from the TUI is planned; ctrl+c still exits the program.

## Planned Features

- Concurrent batch downloads
- Rich progress bars
- Checksum-aware behavior
- TUI cancellation controls
- Queue persistence and history management
- Structured automation output
- Self-update support after the release model is ready
- Public packages from `v1.0.0` onward

## Documentation

Repository documentation lives in `docs/`.

Start here:

- [Documentation Index](docs/index.md)
- [Quick Start](docs/quick-start.md)
- [Installation](docs/installation.md)
- [Usage](docs/usage.md)
- [Command Reference](docs/command-reference.md)
- [Architecture Overview](docs/architecture/overview.md)
- [Roadmap](docs/roadmap/index.md)
- [Versioning Policy](docs/roadmap/versioning-policy.md)
