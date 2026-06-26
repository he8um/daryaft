# Daryaft

Daryaft is a modern terminal downloader written in Go. It is similar in spirit
to `wget`, with a Bubble Tea terminal UI foundation, clean architecture,
packaging, and future self-update support.

> Developed with <3 by AmirHesam Piri

## Status

**v1.12.0 is the current stable release.** Download binary archives from the
[GitHub releases page](https://github.com/he8um/daryaft/releases/tag/v1.12.0).

Daryaft v1.12.0 polishes the TUI download input experience with cleaner URL and
file path validation, actionable guidance messages, and a compact safe defaults
preview on input screens. The full feature
set includes CLI and TUI HTTP/HTTPS downloading with dry-run planning, single URL
and sequential batch downloads, resume from `.part` files, retry with exponential
backoff, CLI checksum verification for single URL downloads (`--checksum`) and
batch downloads via `--checksum-file`, an interactive Bubble Tea TUI, YAML
configuration with environment overrides, `inspect` metadata preflight, `doctor`
diagnostics, shell completions, HTTP request customization (`--proxy`, `--header`,
`--user-agent`, Basic Auth), and a polished update check command. Running
`daryaft` with no arguments opens the interactive TUI home screen.

`daryaft update --check` is read-only: it queries the GitHub Releases API and
reports the current version, the latest stable release, and install-channel-aware
upgrade guidance. It never downloads, installs, or replaces the current binary.
Upgrade guidance depends on how Daryaft was installed:

| Install channel | Suggested command |
|----------------|-------------------|
| Homebrew | `brew update && brew upgrade daryaft` |
| Binary archive | Download the latest release archive from the GitHub release URL |
| Source build | `git pull && go build .` |
| Unknown | Download the latest release from the GitHub Releases page |

Auto-update (`daryaft update` without `--check`) is not yet implemented.

Known limitations (Windows, concurrent downloads, TUI HTTP options, auto-update,
package managers, checksum auto-discovery, PGP/attestation verification,
`Retry-After` not honored) are documented in the
[v1.12.0 release notes](docs/operations/release-notes-v1.11.0.md). Post-1.0
features are tracked in [post-1-feature-packs.md](docs/roadmap/post-1-feature-packs.md).

- Repository: https://github.com/he8um/daryaft
- Website: https://xhesam.com
- Project page: https://xhesam.com/daryaft
- Author: AmirHesam Piri <info@xhesam.com>
- License: MIT

## Install

**Homebrew** (macOS, first live package-manager channel):

```bash
brew tap he8um/tap
brew install daryaft
daryaft version
```

**GitHub binary archives** (all supported platforms):

```bash
# Example: macOS Apple Silicon
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.12.0/daryaft_darwin_arm64.tar.gz
curl -L -O https://github.com/he8um/daryaft/releases/download/v1.12.0/checksums.txt
shasum -a 256 --check checksums.txt
tar -xzf daryaft_darwin_arm64.tar.gz
./daryaft version
```

Available archives: `daryaft_linux_amd64.tar.gz`, `daryaft_linux_arm64.tar.gz`,
`daryaft_darwin_amd64.tar.gz`, `daryaft_darwin_arm64.tar.gz`.

For future releases, `scripts/update-homebrew-formula.sh` updates a local tap
clone from the published GitHub release checksums — the maintainer reviews the
diff and pushes manually. GoReleaser Homebrew publishing remains disabled.
Other package manager channels (deb, rpm, Arch) are later post-1.0 work.
Source builds on `main` report `1.13.0-dev` (current development default);
release builds inject the exact tag version via GoReleaser ldflags.

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
Snapshot versions are named like `1.13.0-dev-SNAPSHOT-<short-commit>`, and any
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
daryaft inspect https://example.com/file.zip --header "X-Token: abc" --user-agent "MyApp/1.0"
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
daryaft download https://example.com/file.zip --proxy http://proxy:8080
daryaft download https://example.com/file.zip --header "X-Custom: value" --header "Accept: application/json"
daryaft download https://example.com/file.zip --user-agent "MyApp/1.0"
daryaft download https://example.com/file.zip --username alice --password secret
daryaft update --check
daryaft update --check --json
daryaft update --check --include-prerelease
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
DARYAFT_USER_AGENT="MyBot/1.0" daryaft https://example.com/file.zip
DARYAFT_TIMEOUT=30s daryaft https://example.com/file.zip
```

`user_agent` sets the default User-Agent for downloads; empty means the
built-in `Daryaft/<version>` default. `timeout` sets the overall HTTP request
timeout as a Go duration string (for example `30s`, `2m`); empty means no
overall timeout. `theme` supports `default` and `mono`; `mono` uses monochrome
TUI styling like `no_color`. `animations` and `hyperlinks` are reserved config
fields stored for future behavior and do not currently change runtime output.

Use `--config <path>` as a global flag to select an explicit configuration file:

```bash
daryaft --config ~/my-daryaft.yaml download https://example.com/file.zip
daryaft --config ~/my-daryaft.yaml config show
```

Use `--timeout <duration>` to set a per-invocation request timeout:

```bash
daryaft download https://example.com/file.zip --timeout 2m
```

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
fields may be unknown when a server omits headers. Inspect supports the same
HTTP customization flags as download: `--proxy`, `--header`, `--user-agent`,
`--username`, and `--password`.

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
prints the checksum plan but does not compute a digest. `--checksum` is single
URL only. The TUI also accepts an optional checksum in the single URL flow.
Daryaft does not discover or download checksum files, and it does not delete the
completed file on mismatch.

For batch downloads, `--checksum-file <path>` verifies each target against a
manifest file of `<algorithm>:<hex> <url>` entries (one per line; blank lines
and `#` comments are ignored):

```bash
daryaft download URL1 URL2 --checksum-file checksums.txt
daryaft download --file urls.txt --checksum-file checksums.txt
```

Every target must have exactly one matching manifest entry, and URLs must match
exactly (no normalization). A mismatch fails that item, leaves the file in
place, and the command exits non-zero. The batch summary reports
`Checksum verified: N`, and the TUI shows `Checksum OK` / `Checksum Failed` per
item. `--checksum` and `--checksum-file` cannot be combined.

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
- Auto-update (`daryaft update` without `--check`) after the release model is ready
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
- [v1.11.0 Release Notes](docs/operations/release-notes-v1.11.0.md)
- [v1.10.0 Release Notes](docs/operations/release-notes-v1.10.0.md)
- [v1.9.0 Release Notes](docs/operations/release-notes-v1.9.0.md)
- [v1.8.0 Release Notes](docs/operations/release-notes-v1.8.0.md)
- [v1.7.0 Release Notes](docs/operations/release-notes-v1.7.0.md)
- [v1.4.0 Release Notes](docs/operations/release-notes-v1.4.0.md)
- [v1.3.0 Release Notes](docs/operations/release-notes-v1.3.0.md)
- [v1.2.0 Release Notes](docs/operations/release-notes-v1.2.0.md)
- [v1.1.0 Release Notes](docs/operations/release-notes-v1.1.0.md)
- [v1.0.0 Release Notes](docs/operations/release-notes-v1.0.0.md)
- [v1.0.0 Release Assets](docs/operations/release-assets.md)
- [Homebrew Tap](docs/operations/homebrew-tap.md)
- [v1.0.0 Go/No-Go Checklist](docs/operations/v1.0.0-go-no-go.md)
- [v1.0.0 Release Plan](docs/operations/v1.0.0-release-plan.md)
- [RC Release Notes: v0.6.0-rc.1](docs/operations/release-notes-v0.6.0-rc.1.md) (superseded)
- [Release Process](docs/operations/release-process.md)
- [Branch Protection](docs/operations/branch-protection.md)
- [Architecture Overview](docs/architecture/overview.md)
- [Roadmap](docs/roadmap/index.md)
- [Versioning Policy](docs/roadmap/versioning-policy.md)
