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
plan screen. The single URL TUI flow can also set an optional custom filename.
The current development version is `0.5.0-dev`.

- Repository: https://github.com/he8um/daryaft
- Website: https://xhesam.com
- Project page: https://xhesam.com/daryaft
- Author: AmirHesam Piri <info@xhesam.com>
- License: MIT

## Install Policy

Public install channels begin at `v1.0.0`. Before that release, Homebrew, Debian,
RPM, Arch, and one-line install instructions are planned only and must not be
presented as stable user-facing install paths. Source builds report
`0.5.0-dev`, commit `local`, build date `unknown`, and built by `source`.
Release builds inject version metadata with linker flags.

For local development:

```bash
go mod download
go run . --help
go run . version
go run . version --json
go run .
go run . config path
go run . config show
go run . config get retries
go run . config set retries 5
go run . https://example.com/file.zip --dry-run
go run . https://example.com/file.zip
go run . https://example.com/a.txt https://example.com/b.txt
go run . download https://example.com/file.zip --dry-run
```

Build and test locally:

```bash
make test
make ci
make build
make build-local
make version
make release-check
make run
```

`make release-check` requires GoReleaser v2 and runs a local snapshot release
check without publishing, creating tags, or enabling package-manager publishing.
Install GoReleaser with `brew install goreleaser` if the command is missing.
Snapshot versions are named like `0.5.0-dev-SNAPSHOT-<short-commit>`, and any
snapshot artifacts are written under ignored local build directories such as
`dist/`.

`make ci` runs the same local pre-release checks expected before opening a PR:
module tidy verification, tests, build, TUI race test, whitespace diff check,
and `goreleaser check` when GoReleaser is installed. GitHub Actions runs the Go
test/build/race matrix on Linux and macOS, plus a separate `goreleaser-check`
job that validates release configuration without publishing.

## Current Commands

```bash
daryaft --help
daryaft version
daryaft version --json
daryaft
daryaft doctor
daryaft completion zsh
daryaft config path
daryaft config show
daryaft config init
daryaft config get retries
daryaft config set retries 5
daryaft config reset
daryaft config keys
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
actions inside the TUI validate input, then let you set an output directory
before showing the dry-run plan. For single URL downloads, the TUI then asks
for an optional custom filename; leaving it empty means auto-detect. The TUI
does not offer one custom filename for `.txt` batch downloads, which continue
to auto-detect each item filename. Leaving the output directory empty means `.`,
the current directory. TUI downloads can start real single URL or sequential
batch downloads from that plan. The TUI resizes its panel and text inputs to
the terminal window. Backspace edits non-empty inputs and navigates back only
when the input is empty; Escape always navigates back. Existing CLI download
commands, including `-o`/`--output` and `--name`, remain fully supported and
unchanged.

Daryaft can read YAML configuration from the OS user config directory:
`<UserConfigDir>/daryaft/config.yaml`. On macOS this is usually
`~/Library/Application Support/daryaft/config.yaml`; on Linux it is usually
`~/.config/daryaft/config.yaml`. Use `daryaft config path` to print the exact
path, `daryaft config init` to create the default file, and
`daryaft config show` to print the effective config. `daryaft config get <key>`
prints one effective value, `daryaft config set <key> <value>` writes one file
value, `daryaft config reset` overwrites the file with defaults, and
`daryaft config keys` lists supported keys and types. Precedence is CLI flags,
then `DARYAFT_*` environment variables, then config file values, then built-in
defaults. For example:

```bash
DARYAFT_DOWNLOAD_DIR=~/Downloads daryaft https://example.com/file.zip
DARYAFT_RETRIES=5 daryaft https://example.com/file.zip
DARYAFT_NO_TUI=true daryaft
```

`theme` supports `default` and `mono`; `mono` uses monochrome TUI styling like
`no_color`. `animations` and `hyperlinks` are reserved config fields stored for
future behavior and do not currently change runtime output.

`daryaft doctor` prints a local diagnostic report for runtime details, version
metadata, config path and loading, default download directory writability,
terminal environment hints, and optional tools. `daryaft doctor --json` prints
the same checks as machine-readable JSON for automation and CI.
`daryaft doctor --strict` treats warnings as a non-zero exit status, while
keeping warning checks marked as warnings in the report. `clamscan` is reported
as an optional tool for future scan features. The GitHub release check is
currently listed as skipped and does not make a network request.

Daryaft can generate shell completion scripts with Cobra's standard completion
command. Installation paths depend on your OS and shell setup:

```bash
daryaft completion zsh > "${fpath[1]}/_daryaft"
daryaft completion bash > /etc/bash_completion.d/daryaft
daryaft completion fish > ~/.config/fish/completions/daryaft.fish
```

Download commands validate URLs and can print a dry-run plan. Real downloading
is implemented for one URL and for multiple URLs sequentially. Batch downloads
continue after item failures and print a final summary. `--retries` controls
retry attempts after the initial attempt; the default `3` means up to four total
attempts. Valid retry values are `0` through `20`.

Use `--verbose` or `-v` with CLI downloads to print additional diagnostic lines
prefixed with `Verbose:`. Verbose output includes the effective URL with user
info and query data redacted, output directory, selected filename, retry/resume
details, HTTP status when known, target path, and completion duration.

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
TUI downloads can be cancelled with `q`; Daryaft keeps the `.part` file and
metadata sidecar for a future resume and does not retry cancelled downloads.
CLI ctrl+c behavior is unchanged and may terminate the process directly.

## Planned Features

- Concurrent batch downloads
- Rich progress bars
- Checksum-aware behavior
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
- [Testing Strategy](docs/engineering/testing-strategy.md)
- [Release Process](docs/operations/release-process.md)
- [Branch Protection](docs/operations/branch-protection.md)
- [Architecture Overview](docs/architecture/overview.md)
- [Roadmap](docs/roadmap/index.md)
- [Versioning Policy](docs/roadmap/versioning-policy.md)
