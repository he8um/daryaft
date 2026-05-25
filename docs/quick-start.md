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

`go run .` opens the interactive TUI home screen. Download actions inside the
TUI can collect a URL or `.txt` file path and show a dry-run plan, but they do
not start downloads yet.

The download command surface supports validation, dry-run planning, and real
single URL downloads. Sequential batch downloads, text progress, retry, and
resume are also implemented. Use the CLI download commands for real downloads
until TUI download execution is implemented.

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
