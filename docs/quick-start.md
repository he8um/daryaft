# Quick Start

Daryaft is currently pre-1.0. Use the local Go toolchain for now.

## Run Locally

```bash
go mod download
go run . --help
go run . version
go run .
```

`go run .` prints the current placeholder message. Interactive mode is planned
for the TUI milestone.

## Build And Test

```bash
make test
make build
./bin/daryaft version
```

Related docs:

- [Installation](installation.md)
- [Usage](usage.md)
- [Command Reference](command-reference.md)
