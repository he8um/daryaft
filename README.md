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
plan screen. The single URL TUI flow can also set an optional custom filename,
and the Inspect URL menu item can show read-only remote metadata without
downloading a file.
The current development version is `0.6.0-dev`.
This milestone represents the stabilized pre-release foundation after
downloader/config correctness work, CLI/TUI cancellation alignment, expanded
HTTP coverage, and local/CI quality gates. Daryaft remains pre-1.0, and
breaking changes may still happen before `v1.0.0`.

- Repository: https://github.com/he8um/daryaft
- Website: https://xhesam.com
- Project page: https://xhesam.com/daryaft
- Author: AmirHesam Piri <info@xhesam.com>
- License: MIT

## Install Policy

Public install channels begin at `v1.0.0`. Before that release, Homebrew, Debian,
RPM, Arch, and one-line install instructions are planned only and must not be
presented as stable user-facing install paths. Source builds report
`0.6.0-dev`, commit `local`, build date `unknown`, and built by `source`.
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
make lint
make security
make rc-check
make rc-info
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
Snapshot versions are named like `0.6.0-dev-SNAPSHOT-<short-commit>`, and any
snapshot artifacts are written under ignored local build directories such as
`dist/`.

`make lint` requires `golangci-lint` and runs the repository's practical lint
profile from `.golangci.yml`. `make security` requires `govulncheck` and
`gosec`, then runs both local security scans as blocking checks. GitHub Actions
also runs separate `lint` and `security` jobs; the security job uses
`govulncheck` and `gosec`, both blocking. `make rc-check` runs
release-candidate readiness checks including `govulncheck` and `gosec`.
`make rc-info` prints the current Git describe value, local RC tags, source
version metadata, and the next RC validation commands.

`make ci` runs the same local pre-release checks expected before opening a PR:
module tidy verification, tests, build, TUI race test, whitespace diff check,
and `goreleaser check` when GoReleaser is installed. GitHub Actions runs the Go
test/build/race matrix on Linux and macOS, plus separate `goreleaser-check`,
`lint`, and `security` jobs. None of these jobs publish releases.

## Current Commands

```bash
daryaft --help
daryaft version
daryaft version --json
daryaft
daryaft doctor
daryaft inspect https://example.com/file.zip
daryaft inspect https://example.com/file.zip --json
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
daryaft https://example.com/file.zip --checksum sha256:<hex>
daryaft https://example.com/a.txt https://example.com/b.txt
daryaft -f urls.txt --dry-run
daryaft -f urls.txt
daryaft download https://example.com/file.zip --dry-run
daryaft download https://example.com/file.zip
daryaft download -f urls.txt --dry-run
daryaft download -f urls.txt
```

With no arguments, Daryaft opens a Bubble Tea home screen with menu entries for
URL input, `.txt` file input, Inspect URL, help, version information, and quit.
Download actions inside the TUI validate input, then let you set an output
directory before showing the dry-run plan. For single URL downloads, the TUI
then asks for an optional custom filename and optional checksum; leaving them
empty means auto-detect filename and skip checksum verification. The TUI does
not offer one custom filename or checksum for `.txt` batch downloads, which
continue to auto-detect each item filename. Leaving the output directory empty
uses the effective default output directory, which falls back to `~/Downloads`
when no CLI flag, environment variable, or config value is set. Enter `.` to
download to the current directory explicitly. TUI downloads can start real
single URL or sequential batch downloads from that plan. Inspect URL accepts
one HTTP/HTTPS URL, shows remote metadata, and does not download or write
files. The TUI resizes its panel and text inputs to the terminal window.
Backspace edits non-empty inputs and navigates back only when the input is
empty; Escape always navigates back. Existing CLI download commands, including
`-o`/`--output` and `--name`, remain fully supported and unchanged.

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

`daryaft inspect <url>` prints HTTP metadata for one HTTP/HTTPS URL without
saving a file. It follows redirects, reports the final URL, status, inferred
filename, content length when known, content type, `Accept-Ranges`, resume
support, `ETag`, and `Last-Modified`. Use `daryaft inspect <url> --json` for
stable machine-readable CLI output. The no-argument TUI also has a read-only
Inspect URL flow for the same metadata, but JSON output remains CLI-only. Some
fields may be unknown when a server omits headers. Inspect does not implement
checksum, proxy, custom-header, or auth yet.

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

CLI single URL downloads can verify a manually supplied checksum after the
final file is completed and renamed into place:

```bash
daryaft https://example.com/file.zip --checksum sha256:<hex>
daryaft download https://example.com/file.zip --checksum sha512:<hex>
```

Supported checksum algorithms are `sha256` and `sha512`. Dry-run validates and
prints the checksum plan but does not compute a digest. `--checksum` is
currently single URL only. The TUI also accepts an optional checksum in the
single URL flow. Batch checksum remains unsupported, Daryaft does not discover
or download checksum files, and it does not delete the completed file on
mismatch.

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
CLI Ctrl+C/SIGTERM cancellation uses the same downloader context path, keeps
the `.part` file and metadata sidecar for resume, avoids renaming to the final
target, stops remaining batch items, and exits non-zero.

## Planned Features

- Concurrent batch downloads
- Rich progress bars
- Signed checksum handling
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
- [Pre-Release Readiness](docs/operations/pre-release-readiness.md)
- [Release-Candidate Validation](docs/operations/rc-validation.md)
- [RC Release Notes: v0.6.0-rc.2](docs/operations/release-notes-v0.6.0-rc.2.md)
- [v0.6.0-rc.2 Release Status](docs/operations/release-status-v0.6.0-rc.2.md)
- [v1.0.0 Release Notes (draft)](docs/operations/release-notes-v1.0.0.md)
- [v1.0.0 Release Assets](docs/operations/release-assets.md)
- [v1.0.0 Go/No-Go Checklist](docs/operations/v1.0.0-go-no-go.md)
- [v1.0.0 Release Plan](docs/operations/v1.0.0-release-plan.md)
- [RC Release Notes: v0.6.0-rc.1](docs/operations/release-notes-v0.6.0-rc.1.md) (superseded)
- [Release Process](docs/operations/release-process.md)
- [Branch Protection](docs/operations/branch-protection.md)
- [Architecture Overview](docs/architecture/overview.md)
- [Roadmap](docs/roadmap/index.md)
- [Versioning Policy](docs/roadmap/versioning-policy.md)
