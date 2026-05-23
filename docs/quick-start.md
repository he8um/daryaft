# Quick Start

Daryaft is currently pre-1.0. Use the local Go toolchain for now.

## Run Locally

```bash
go mod download
go run . --help
go run . version
go run .
go run . https://example.com/file.zip --dry-run
go run . https://example.com/file.zip --output downloads
go run . download https://example.com/file.zip --dry-run
```

`go run .` prints the current placeholder message. Interactive mode is planned
for the TUI milestone.

The download command surface supports validation, dry-run planning, and real
single URL downloads. Batch real downloads, TUI, progress, resume, and retry
execution are planned.

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
