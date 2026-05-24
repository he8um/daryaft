# Daryaft

Daryaft is a modern terminal downloader written in Go. It is similar in spirit
to `wget`, with a planned terminal UI, clean architecture, packaging, and future
self-update support.

> Developed with <3 by AmirHesam Piri

## Status

Daryaft is in pre-1.0 development and is not stable yet. The current codebase
has a CLI foundation, dry-run planning, and real single URL HTTP/HTTPS
downloads with simple text progress output. Multiple URLs can also be
downloaded sequentially in one command.

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

With no arguments, Daryaft prints a friendly placeholder. Interactive mode is
planned for the TUI milestone.

Download commands validate URLs and can print a dry-run plan. Real downloading
is implemented for one URL and for multiple URLs sequentially. Batch downloads
continue after item failures and print a final summary.

The downloader now emits structured started, progress, completed, and failed
events for the single URL path. The CLI consumes those events for line-based
progress output for single and sequential batch downloads; the full terminal UI
is still planned.

## Planned Features

- Concurrent batch downloads
- Rich progress bars
- Resume, retry, and checksum-aware behavior
- Beautiful terminal UI
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
