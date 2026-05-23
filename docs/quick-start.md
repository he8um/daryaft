# Quick Start

Daryaft is currently pre-1.0. Use the local Go toolchain for now.

## Run Locally

```bash
go mod download
go run . --help
go run . version
go run .
go run . https://example.com/file.zip --dry-run
go run . download https://example.com/file.zip --dry-run
```

`go run .` prints the current placeholder message. Interactive mode is planned
for the TUI milestone.

The download command surface currently supports validation and dry-run planning.
Real downloading is planned for the next downloader engine milestone.

## Build And Test

```bash
make test
make build
./bin/daryaft version
./bin/daryaft https://example.com/file.zip --dry-run
```

Related docs:

- [Installation](installation.md)
- [Usage](usage.md)
- [Command Reference](command-reference.md)
