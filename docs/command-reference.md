# Command Reference

## `daryaft`

Implemented. Opens the Bubble Tea interactive home screen when run with no
arguments.

The home screen shows:

- Daryaft
- Modern terminal downloader
- Download from URL
- Download from .txt file
- View help
- Version
- Quit

Use up/down arrows or `k`/`j` to move, enter to select, and `esc` or backspace
to return from sub-screens. `q` quits unless a download is running; ctrl+c exits
from anywhere. Download from URL and Download from .txt file open input forms,
validate with the existing download planner, then ask for an output directory
before showing dry-run plans. The single URL flow then asks for an optional
custom filename; leaving it empty means auto-detect. The `.txt` batch flow does
not offer one custom filename and keeps per-item auto-detect. Leaving the
output directory empty means `.`, the current directory. Press enter on the
plan screen to start a real download in the TUI. Existing CLI download commands
remain fully supported, and CLI `-o`/`--output` and `--name` behavior is
unchanged. Pressing `q` while a TUI download is running cancels it and keeps
partial state for resume.

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
  "version": "0.5.0-dev",
  "commit": "local",
  "date": "unknown",
  "built_by": "source",
  "go_version": "go1.xx.x"
}
```

Source builds default to version `0.5.0-dev`, commit `local`, date `unknown`,
and built by `source`. Release builds inject metadata through ldflags. Daryaft
is still pre-1.0; public stable install channels begin at `v1.0.0`.

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
  current directory. Existing unwritable output directories are critical
  failures. Missing configured output directories are warnings and are not
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

`retries` must be an integer greater than or equal to `0`. Boolean values
accept `true`, `1`, `yes`, `y`, `on`, `false`, `0`, `no`, `n`, and `off`,
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
```

If the server does not provide a known content length, progress uses:

```text
Progress: <downloaded> bytes | <speed>
```

Progress lines are generated from structured downloader events. The TUI
execution screen consumes the same event stream for status, target path,
downloaded bytes, percent, speed, retry/resume/restart messages, completion,
failure, and summaries.

Cancellation is supported by the downloader context path. TUI cancellation
emits a cancelled event, leaves the `.part` file and metadata sidecar in place,
does not rename to the final target, and does not retry.

`--retries` is implemented for transient failures. The value is the number of
retry attempts after the first try, so `--retries 0` means one total attempt and
`--retries 3` means up to four total attempts. Daryaft retries network errors,
timeouts, interrupted response bodies, HTTP `429`, `500`, `502`, `503`, and
`504`.

`--resume` is enabled by default. If a partial file exists, Daryaft sends
`Range: bytes=<partial_size>-` and appends only after `206 Partial Content`.
If the server returns a full response instead, Daryaft truncates the `.part`
file and restarts safely. If saved `ETag` or `Last-Modified` metadata shows the
remote file changed, Daryaft also restarts from byte `0`. `--no-resume` ignores
existing partial data and overwrites the partial file from byte `0`.

Config precedence is CLI flags, then environment variables, then config file
values, then built-in defaults. If `-o`/`--output` is omitted and
`DARYAFT_DOWNLOAD_DIR` or `download_dir` is non-empty, Daryaft uses that
directory. If `--retries` is omitted, Daryaft uses `DARYAFT_RETRIES` or config
`retries`. If neither `--resume` nor `--no-resume` is set, Daryaft uses
`DARYAFT_RESUME` or config `resume`.

## `daryaft download [url...] --dry-run`

Implemented. Explicit form of the same download validation and dry-run planner.

```bash
daryaft download https://example.com/file.zip --dry-run
daryaft download -f urls.txt --dry-run
```

## `daryaft download <url>`

Implemented. Explicit form of single URL real download.

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
- `--retries int`: retry attempts after the initial attempt, default `3`.
- `--resume`: resume interrupted `.part` files with HTTP Range, default `true`.
- `--no-resume`: disable resume and restart partial files from byte `0`.

Validation rules:

- At least one URL arg or `--file` is required for download mode.
- URL args and `--file` are combined.
- Only `http` and `https` URLs are accepted.
- Empty lines and `#` comments are ignored in URL files.
- `--name` is rejected when more than one URL is provided.
- `--retries` must be greater than or equal to `0`.

## Common Flags

Implemented:

- `--no-color`: avoid color styling in the TUI.
- `--no-tui`: skip the no-argument TUI and print the non-interactive placeholder.
- `-v`, `--verbose`: enable verbose output when verbose logging exists.

Config `no_color` and `no_tui`, or `DARYAFT_NO_COLOR` and `DARYAFT_NO_TUI`,
provide defaults for the no-argument TUI path. CLI flags still have priority.
Config `theme`, `animations`, and `hyperlinks`, plus their `DARYAFT_*`
environment variables, are stored for future TUI support.

Environment variables:

- `DARYAFT_DOWNLOAD_DIR`
- `DARYAFT_RETRIES`
- `DARYAFT_RESUME`
- `DARYAFT_NO_COLOR`
- `DARYAFT_NO_TUI`
- `DARYAFT_THEME`
- `DARYAFT_ANIMATIONS`
- `DARYAFT_HYPERLINKS`

Boolean environment values accept `true`, `1`, `yes`, `y`, `on`, `false`, `0`,
`no`, `n`, and `off`, case-insensitively. Empty boolean and integer values are
invalid.

Examples:

```bash
DARYAFT_DOWNLOAD_DIR=~/Downloads daryaft https://example.com/file.zip
DARYAFT_RETRIES=5 daryaft https://example.com/file.zip
DARYAFT_NO_TUI=true daryaft
```

## Planned Commands And Forms

These are planned and not implemented yet:

```bash
daryaft update
```

Batch concurrency, queue persistence, rich progress bars, segmented downloads,
and self-update are planned.

Related docs:

- [Usage](usage.md)
- [Architecture Overview](architecture/overview.md)
