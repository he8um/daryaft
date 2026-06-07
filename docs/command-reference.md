# Command Reference

## `daryaft`

Implemented. Opens the Bubble Tea interactive home screen when run with no
arguments.

The home screen shows:

- Daryaft
- Modern terminal downloader
- Download from URL
- Download from .txt file
- Inspect URL
- View help
- Version
- Quit

Use up/down arrows or `k`/`j` to move and enter to select. Escape navigates
back from sub-screens. Backspace edits text when the current input has content
and navigates back only when the input is empty. `q` quits unless a download is
running; ctrl+c exits from anywhere. Download from URL and Download from .txt
file open input forms, validate with the existing download planner, then ask
for an output directory before showing dry-run plans. The single URL flow then
asks for an optional custom filename; leaving it empty means auto-detect. The
`.txt` batch flow does not offer one custom filename and keeps per-item
auto-detect. Leaving the output directory empty uses the effective output
default, which falls back to `~/Downloads` when no environment or config value
is set. Enter `.` to use the current directory explicitly. Press enter on the
plan screen to start a real download in the TUI. Inspect URL prompts for one
HTTP/HTTPS URL and shows read-only metadata without starting a download or
writing files. The TUI panel and input width adapt to terminal resize messages
with bounded minimum and maximum widths. Existing CLI download commands remain
fully supported, and CLI `-o`/`--output` and `--name` behavior is unchanged.
Pressing `q` while a TUI download is running cancels it and keeps partial state
for resume.

When URL arguments or `--file` are provided, the root command enters the current
download validation mode.

## `daryaft --help`

Implemented. Shows help text with:

- short description
- usage
- available commands
- common flags
- roadmap examples
- footer line

## `daryaft version`

Implemented. Prints:

- Daryaft version
- commit
- build date
- built by
- Go version

Use `daryaft version --json` for stable machine-readable output:

```json
{
  "version": "1.1.0-dev",
  "commit": "local",
  "date": "unknown",
  "built_by": "source",
  "go_version": "go1.xx.x"
}
```

Source builds default to version `1.1.0-dev`, commit `local`, date `unknown`,
and built by `source`. Release builds inject metadata through ldflags. The
latest stable release is `v1.0.0`; the source development version is now
`1.1.0-dev`.

## `daryaft inspect <url>`

Implemented. Inspects one HTTP/HTTPS URL and prints metadata without saving a
file. Exactly one URL argument is required, and the URL scheme must be `http`
or `https`.

Human output includes:

- original URL
- final URL after redirects
- HTTP status
- inferred filename
- content length, or `unknown`
- content type, or `unknown`
- `Accept-Ranges`, or `unknown`
- resume support as `yes`, `no`, or `unknown`
- `ETag`, or `unknown`
- `Last-Modified`, or `unknown`

Use `--json` for stable machine-readable output:

```json
{
  "url": "https://example.com/file.zip",
  "final_url": "https://cdn.example.com/file.zip",
  "status": "200 OK",
  "status_code": 200,
  "filename": "file.zip",
  "content_length": 1048576,
  "content_length_known": true,
  "content_type": "application/zip",
  "accept_ranges": "bytes",
  "resume_supported": true,
  "resume_support_known": true,
  "etag": "\"abc123\"",
  "last_modified": "Tue, 01 Jun 2026 12:00:00 GMT"
}
```

Daryaft tries `HEAD` first and may fall back to `GET` with
`Range: bytes=0-0` when `HEAD` is not allowed or omits useful metadata. It
closes response bodies and does not write `.part` files, metadata sidecars, or
final downloads. Metadata may be unknown when the server omits headers.

`daryaft inspect` accepts the same HTTP customization flags as download:

| Flag | Description |
|------|-------------|
| `--proxy <url>` | HTTP or HTTPS proxy URL |
| `--header "Name: Value"` | Custom request header (repeatable) |
| `--user-agent <value>` | Override the default `User-Agent` |
| `--username <value>` | HTTP Basic Auth username |
| `--password <value>` | HTTP Basic Auth password |

See [HTTP Request Customization](../features/http-request-customization.md) for full details.

## `daryaft doctor`

Implemented. Prints a local diagnostics report using simple text output.

Use `daryaft doctor --json` to print the same diagnostics as machine-readable
JSON for automation and CI. JSON mode does not print the human text report.
Use `daryaft doctor --strict` when warnings should fail the command, such as in
CI. `--json` and `--strict` can be combined.

Checks include:

- Go runtime OS, architecture, and version.
- Daryaft version, commit, and build date.
- Config path and whether the config directory is writable or appears
  creatable.
- Effective config loading. Invalid YAML is a critical failure.
- Default download directory. If `download_dir` is empty, this checks the
  built-in `~/Downloads` default. Existing unwritable output directories are
  critical failures. Missing output directories are warnings and are not
  created.
- Terminal environment hints: `TERM`, `NO_COLOR`, and stdout terminal status
  when available.
- Optional `clamscan` detection. Missing `clamscan` is informational only and
  reserved for future scan features.
- GitHub release check status. This foundation does not make a network request
  and reports the check as skipped.

Status markers:

```text
✓ ok
✗ critical failure
! warning
- informational
```

JSON output uses this stable shape:

```json
{
  "ok": true,
  "summary": {
    "failures": 0,
    "warnings": 0,
    "checks": 12
  },
  "sections": [
    {
      "name": "System",
      "checks": [
        {
          "status": "ok",
          "label": "OS",
          "message": "darwin"
        }
      ]
    }
  ]
}
```

JSON check statuses are `ok`, `warning`, `failure`, `info`, and `skipped`.
Both text and JSON modes exit non-zero when any critical failure is present.
With `--strict`, warnings also produce a non-zero exit status. JSON strict mode
keeps warning checks as `warning`, reports warnings separately from failures,
and sets top-level `ok` to `false` when warnings are present:

```json
{
  "ok": false,
  "strict": true,
  "summary": {
    "failures": 0,
    "warnings": 1,
    "checks": 16
  }
}
```

## `daryaft completion [bash|zsh|fish|powershell]`

Implemented. Generates shell completion scripts using Cobra's standard
generators.

```bash
daryaft completion bash
daryaft completion zsh
daryaft completion fish
daryaft completion powershell
```

Example setup commands:

```bash
daryaft completion zsh > "${fpath[1]}/_daryaft"
daryaft completion bash > /etc/bash_completion.d/daryaft
daryaft completion fish > ~/.config/fish/completions/daryaft.fish
```

Completion installation paths vary by OS, shell, and user permissions.
Unsupported shell names return a clear error.

## `daryaft config`

Implemented. Shows help for the configuration command group.

## `daryaft config path`

Implemented. Prints the YAML config path:

```text
<UserConfigDir>/daryaft/config.yaml
```

On macOS this is usually
`~/Library/Application Support/daryaft/config.yaml`. On Linux this is usually
`~/.config/daryaft/config.yaml`.

## `daryaft config show`

Implemented. Prints the effective config as YAML, including `DARYAFT_*`
environment overrides. If the file does not exist and no environment overrides
are set, the output is the built-in defaults:

```yaml
download_dir: ""
retries: 3
resume: true
no_color: false
no_tui: false
theme: default
animations: true
hyperlinks: true
```

## `daryaft config init`

Implemented. Creates the default config file and fails clearly if it already
exists.

## `daryaft config init --force`

Implemented. Rewrites the config file with default values.

## `daryaft config get <key>`

Implemented. Prints one effective config value. Environment variables are
reflected because `config get` reads the same effective config as
`config show`.

```bash
daryaft config get retries
daryaft config get download_dir
```

Unknown keys fail clearly.
Shell completion suggests all supported config keys.

## `daryaft config set <key> <value>`

Implemented. Sets one value in the YAML config file, creating the config
directory if needed. Environment variables are not written to the file.

```bash
daryaft config set retries 5
daryaft config set download_dir ~/Downloads
daryaft config set resume off
```

`retries` must be an integer from `0` through `20`. Boolean values accept
`true`, `1`, `yes`, `y`, `on`, `false`, `0`, `no`, `n`, and `off`,
case-insensitively.
Shell completion suggests supported keys for the first argument and `true` or
`false` for boolean value arguments.

## `daryaft config reset`

Implemented. Overwrites the config file with built-in defaults and prints the
path that was reset.

## `daryaft config keys`

Implemented. Lists supported config keys and expected types:

```text
download_dir string
retries int
resume bool
no_color bool
no_tui bool
theme string
animations bool
hyperlinks bool
```

## `daryaft [url...] --dry-run`

Implemented. Validates one or more HTTP/HTTPS URLs and prints a dry-run plan.

```bash
daryaft https://example.com/file.zip --dry-run
daryaft https://example.com/a.txt https://example.com/b.txt --dry-run
daryaft -f urls.txt --dry-run
```

## `daryaft <url>`

Implemented for exactly one URL. Performs an HTTP GET, accepts HTTP 2xx
responses, writes to `<filename>.part`, then renames the file after completion.
Incomplete downloads keep sidecar metadata at
`<filename>.part.daryaft.json`.

```bash
daryaft https://example.com/file.zip
daryaft https://example.com/file.zip --output downloads
daryaft https://example.com/file.zip --name file.zip
daryaft https://example.com/file.zip --checksum sha256:<hex>
```

The command does not overwrite existing final files. It uses simple text output:

```text
Downloading: <url>
Saving to: <path>
Progress: <downloaded> / <total> bytes (<percent>%) | <speed>
Resuming from <bytes> bytes
Resume not supported by server; restarting download
Remote file changed; restarting download
Retrying <attempt>/<max> in <delay>: <reason>
Completed: <path>
Checksum verified: sha256
```

If the server does not provide a known content length, progress uses:

```text
Progress: <downloaded> bytes | <speed>
```

Progress lines are generated from structured downloader events. The TUI
execution screen consumes the same event stream for status, target path,
downloaded bytes, percent, speed, retry/resume/restart messages, completion,
failure, and summaries.

Cancellation is supported by the downloader context path. In CLI downloads,
Ctrl+C or SIGTERM cancels the active request/copy loop, prints
`Download cancelled. Partial file kept for resume.`, leaves the `.part` file
and metadata sidecar in place, does not rename to the final target, and returns
a non-zero exit code. In batch mode, cancellation stops the active item and
does not start remaining URLs; the summary reports cancelled and skipped items
when practical. TUI cancellation still uses `q` from the progress screen and
preserves the same partial state for resume.

`--retries` is implemented for transient failures. The value is the number of
retry attempts after the first try, so `--retries 0` means one total attempt and
`--retries 3` means up to four total attempts. Values must be from `0` through
`20`. Daryaft retries network errors, timeouts, interrupted response bodies,
HTTP `429`, `500`, `502`, `503`, and `504`.

`--resume` is enabled by default. If a partial file exists, Daryaft sends
`Range: bytes=<partial_size>-` and appends only after `206 Partial Content`.
If the server returns a full response instead, Daryaft truncates the `.part`
file and restarts safely. If saved `ETag` or `Last-Modified` metadata shows the
remote file changed, Daryaft also restarts from byte `0`. `--no-resume` ignores
existing partial data and overwrites the partial file from byte `0`.

Config precedence is CLI flags, then environment variables, then config file
values, then built-in defaults. If `-o`/`--output` is omitted, Daryaft uses
`DARYAFT_DOWNLOAD_DIR`, then config `download_dir`, then the built-in
`~/Downloads` default. Explicit `.` still means the current directory. If
`--retries` is omitted, Daryaft uses `DARYAFT_RETRIES` or config `retries`. If
neither `--resume` nor `--no-resume` is set, Daryaft uses `DARYAFT_RESUME` or
config `resume`.

With `--verbose` or `-v`, CLI downloads print additional diagnostic lines
prefixed with `Verbose:`. These include the effective URL with user info,
query, and fragment redacted, output directory, selected filename, retry/resume
details, HTTP status when known, target path, and completion duration. Normal
non-verbose output is unchanged.

`--checksum` is implemented for manual single URL verification after a
successful completed download. Supported forms are `sha256:<hex>` and
`sha512:<hex>`. The checksum is validated before the network request starts.
After the final file is renamed into place, Daryaft computes the digest and
prints `Checksum verified: sha256` or `Checksum verified: sha512` on match. On
mismatch, it returns a non-zero error such as `checksum mismatch: expected
<expected>, got <actual>` and leaves the completed final file in place.
Dry-run validates and prints the checksum but does not compute it. `--checksum`
is rejected for multiple URLs and for `--file` input. The TUI also supports
optional checksum input for single URL downloads only.

## `daryaft download [url...] --dry-run`

Implemented. Explicit form of the same download validation and dry-run planner.

```bash
daryaft download https://example.com/file.zip --dry-run
daryaft download -f urls.txt --dry-run
```

## `daryaft download <url>`

Implemented. Explicit form of single URL real download.

The download commands accept HTTP customization flags:

| Flag | Description |
|------|-------------|
| `--proxy <url>` | HTTP or HTTPS proxy URL |
| `--header "Name: Value"` | Custom request header (repeatable) |
| `--user-agent <value>` | Override the default `User-Agent` |
| `--username <value>` | HTTP Basic Auth username |
| `--password <value>` | HTTP Basic Auth password |

See [HTTP Request Customization](../features/http-request-customization.md) for full details, validation rules, and security warnings.

## `daryaft [url...]`

Implemented for one or more URLs. When more than one URL is present, downloads
run sequentially in input order and continue after item failures.

```bash
daryaft https://example.com/a.txt https://example.com/b.txt
daryaft -f urls.txt
daryaft https://example.com/a.txt -f urls.txt
```

Batch output uses item headers, per-item progress, and a final summary:

```text
[1/2] Downloading: <url>
Saving to: <path>
Progress: <downloaded> / <total> bytes (<percent>%) | <speed>
Completed: <path>
Daryaft batch summary
Total: 2
Completed: 2
Failed: 0
```

If any item fails, Daryaft continues with the remaining URLs, lists failed
downloads in the summary, and returns a non-zero exit status at the end. Each
batch item has its own retry cycle.

## `daryaft download [url...]`

Implemented. Explicit form of sequential batch download.

## Download Flags

- `-f`, `--file string`: read URLs from a file.
- `-o`, `--output string`: output directory.
- `--name string`: filename for a single URL.
- `--dry-run`: validate inputs and print the download plan.
- `--checksum string`: verify a completed single URL download with
  `sha256:<hex>` or `sha512:<hex>`.
- `--retries int`: retry attempts after the initial attempt, default `3`,
  valid range `0` through `20`.
- `--resume`: resume interrupted `.part` files with HTTP Range, default `true`.
- `--no-resume`: disable resume and restart partial files from byte `0`.

Validation rules:

- At least one URL arg or `--file` is required for download mode.
- URL args and `--file` are combined.
- Only `http` and `https` URLs are accepted.
- Empty lines and `#` comments are ignored in URL files.
- `--name` is rejected when more than one URL is provided.
- `--checksum` is rejected when more than one URL is provided or when `--file`
  is used.
- `--checksum` must use `sha256:<64 hex chars>` or `sha512:<128 hex chars>`.
- `--retries` must be from `0` through `20`.

## Common Flags

Implemented:

- `--no-color`: avoid color styling in the TUI.
- `--no-tui`: skip the no-argument TUI and print the non-interactive placeholder.
- `-v`, `--verbose`: enable extra CLI download diagnostics.

Config `no_color` and `no_tui`, or `DARYAFT_NO_COLOR` and `DARYAFT_NO_TUI`,
provide defaults for the no-argument TUI path. CLI flags still have priority.
Config `theme` supports `default` and `mono`; `mono` uses monochrome TUI
styling like `no_color`. Config `animations` and `hyperlinks`, plus their
`DARYAFT_*` environment variables, are reserved fields stored for future
behavior and do not currently change runtime output.

Environment variables:

- `DARYAFT_DOWNLOAD_DIR`
- `DARYAFT_RETRIES`
- `DARYAFT_RESUME`
- `DARYAFT_NO_COLOR`
- `DARYAFT_NO_TUI`
- `DARYAFT_THEME`
- `DARYAFT_ANIMATIONS`
- `DARYAFT_HYPERLINKS`

`DARYAFT_RETRIES` must be from `0` through `20`. `DARYAFT_THEME` must be
`default` or `mono`. Boolean environment values accept `true`, `1`, `yes`, `y`,
`on`, `false`, `0`, `no`, `n`, and `off`, case-insensitively. Empty boolean and
integer values are invalid.

Examples:

```bash
DARYAFT_DOWNLOAD_DIR=~/Downloads daryaft https://example.com/file.zip
DARYAFT_RETRIES=5 daryaft https://example.com/file.zip
DARYAFT_NO_TUI=true daryaft
```

## `daryaft update --check`

Implemented. Queries the GitHub Releases API and compares the current version to
the latest stable release. Read-only: it does not download, install, or replace
the current binary. Auto-update is not implemented.

```bash
daryaft update --check
daryaft update --check --json
daryaft update --check --include-prerelease
```

Human output (up to date):

```text
Daryaft update check

Current version:  1.1.0
Latest stable:    1.1.0
Status:           up to date

Release: https://github.com/he8um/daryaft/releases/tag/v1.1.0

Install channel:  homebrew
Update command:   brew update && brew upgrade daryaft
```

Human output (update available):

```text
Daryaft update check

Current version:  1.0.0
Latest stable:    1.1.0
Status:           update available

Release: https://github.com/he8um/daryaft/releases/tag/v1.1.0

Install channel:  homebrew
Update command:   brew update && brew upgrade daryaft

A new version is available. Use the update command above to upgrade.
```

Human output (development build):

```text
Daryaft update check

Current version:  1.3.0-dev
Latest stable:    1.2.0
Status:           development build

Release: https://github.com/he8um/daryaft/releases/tag/v1.2.0

Install channel:  source
Update command:   Pull the repository and rebuild: git pull && go build .

Note: development builds may be ahead of the latest stable release.
```

The `Release:` URL is shown in all status modes.

Use `--json` for stable machine-readable output:

```json
{
  "current_version": "1.1.0",
  "latest_version": "1.1.0",
  "update_available": false,
  "development_build": false,
  "include_prerelease": false,
  "install_channel": "homebrew",
  "release_url": "https://github.com/he8um/daryaft/releases/tag/v1.1.0",
  "update_command": "brew update && brew upgrade daryaft",
  "message": "Daryaft is up to date."
}
```

All JSON fields are always present. Boolean fields are JSON booleans.
String fields that are empty are present as `""`, not omitted.

`--include-prerelease` queries the releases list and may return a newer
pre-release version as `latest_version`.

Exits `0` on success regardless of whether an update is available.
Exits non-zero only on network or API failure.

`daryaft update` without `--check` exits non-zero with a clear message
directing the user to `daryaft update --check`.

Install channel detection and update guidance:

| Channel | Detected when | Update command |
|---------|--------------|----------------|
| `homebrew` | Executable path is under a Homebrew prefix/Cellar | `brew update && brew upgrade daryaft` |
| `goreleaser` | Built by GoReleaser, path not under Homebrew | Download the latest release archive from the GitHub release URL |
| `source` | Built from source (`built_by: source`) | `git pull && go build .` |
| `unknown` | Cannot determine install channel | Download the latest release from the GitHub release URL |

Auto-update is not implemented. See [Self-Update Roadmap](../roadmap/self-update.md)
and [Update Check QA](../operations/update-check-qa.md).

## Planned Commands And Forms

These are not yet implemented:

Batch concurrency, queue persistence, rich progress bars, segmented downloads,
and self-update are planned.

Related docs:

- [Usage](usage.md)
- [Architecture Overview](architecture/overview.md)
