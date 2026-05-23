# Daryaft

Daryaft is a modern terminal downloader written in Go. It is similar in spirit
to `wget`, with a planned terminal UI, clean architecture, packaging, and future
self-update support.

> Developed with <3 by AmirHesam Piri

## Status

Daryaft is in pre-1.0 development and is not stable yet. The current codebase
has a CLI foundation plus download input validation and dry-run planning.

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
daryaft -f urls.txt --dry-run
daryaft download https://example.com/file.zip --dry-run
daryaft download -f urls.txt --dry-run
```

With no arguments, Daryaft prints a friendly placeholder. Interactive mode is
planned for the TUI milestone.

Download commands currently validate URLs and print a dry-run plan. Real HTTP
downloading is planned for the next downloader engine milestone. Non-dry-run
download attempts return a clear not-implemented error.

## Planned Features

- Real single URL downloads
- Real batch downloads from files
- Resume, retry, and checksum-aware behavior
- Beautiful terminal UI
- Queue and history management
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
