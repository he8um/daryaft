# Configuration

Daryaft has a small YAML configuration system for user defaults.

## Config Path

Daryaft uses `os.UserConfigDir()` and stores its config at:

```text
<UserConfigDir>/daryaft/config.yaml
```

Common locations:

- macOS: `~/Library/Application Support/daryaft/config.yaml`
- Linux: `~/.config/daryaft/config.yaml`

Use the command below to print the exact path for the current machine:

```bash
daryaft config path
```

To use a different config file, pass `--config <path>` as a global flag before any subcommand:

```bash
daryaft --config ~/my-daryaft.yaml config show
daryaft --config ~/my-daryaft.yaml download https://example.com/file.zip
```

If `--config` is omitted, Daryaft uses the default path above. A missing default config is not an error — Daryaft continues with built-in defaults. A missing explicit `--config` path is an error. `daryaft config init` with `--config` creates the file at the specified path.

## Commands

```bash
daryaft config
daryaft config path
daryaft config show
daryaft config init
daryaft config init --force
daryaft config get <key>
daryaft config set <key> <value>
daryaft config reset
daryaft config keys
```

- `daryaft config` shows help for the config command group.
- `daryaft config path` prints the config file path.
- `daryaft config show` prints the effective config as YAML, including
  environment-variable overrides. If no file exists, it prints built-in
  defaults plus any environment overrides.
- `daryaft config init` creates a config file with defaults and refuses to
  overwrite an existing file.
- `daryaft config init --force` overwrites the existing config with defaults.
- `daryaft config get <key>` prints one effective value, including environment
  overrides.
- `daryaft config set <key> <value>` writes one value to the config file. It
  does not write environment-variable overrides.
- `daryaft config reset` overwrites the config file with built-in defaults.
- `daryaft config keys` lists supported keys and expected value types.

Examples:

```bash
daryaft config keys
daryaft config get retries
daryaft config set retries 5
daryaft config set download_dir ~/Downloads
daryaft config reset
```

## Default YAML

```yaml
download_dir: ""
retries: 3
resume: true
user_agent: ""
timeout: ""
no_color: false
no_tui: false
theme: default
animations: true
hyperlinks: true
```

## Precedence

Configuration is applied in this order:

1. CLI flags
2. Environment variables
3. Config file values
4. Built-in defaults

Examples:

- `-o`/`--output` wins over `download_dir`.
- `--retries` wins over `retries`.
- `--resume` and `--no-resume` win over `resume`.
- `--no-color` wins over `no_color`.
- `--no-tui` wins over `no_tui`.

## Environment Variables

Supported environment variables:

- `DARYAFT_DOWNLOAD_DIR`
- `DARYAFT_RETRIES`
- `DARYAFT_RESUME`
- `DARYAFT_USER_AGENT`
- `DARYAFT_TIMEOUT`
- `DARYAFT_NO_COLOR`
- `DARYAFT_NO_TUI`
- `DARYAFT_THEME`
- `DARYAFT_ANIMATIONS`
- `DARYAFT_HYPERLINKS`

String values are trimmed. Empty string values are accepted for
`DARYAFT_DOWNLOAD_DIR`, `DARYAFT_THEME`, `DARYAFT_USER_AGENT`, and
`DARYAFT_TIMEOUT`.

`DARYAFT_RETRIES` must be an integer from `0` through `20`. Empty, negative,
too-large, or otherwise invalid integer values return an error.

Boolean values are case-insensitive:

- true: `true`, `1`, `yes`, `y`, `on`
- false: `false`, `0`, `no`, `n`, `off`

Empty or invalid boolean values return an error.

Examples:

```bash
DARYAFT_DOWNLOAD_DIR=~/Downloads daryaft https://example.com/file.zip
DARYAFT_RETRIES=5 daryaft https://example.com/file.zip
DARYAFT_NO_TUI=true daryaft
```

## Fields

- `download_dir`: default output directory. Empty means Daryaft uses the
  built-in default `~/Downloads` unless the CLI, TUI, or environment explicitly
  sets output. Set `download_dir` to `.` when the current directory should be
  the saved default.
- `retries`: default retry attempts after the initial attempt. Valid range:
  `0` through `20`.
- `resume`: default resume behavior for interrupted `.part` files.
- `user_agent`: default User-Agent header for downloads. Empty means Daryaft uses
  its built-in default (`Daryaft/<version>`). `--user-agent` and `DARYAFT_USER_AGENT`
  override this value. Must not contain control characters.
- `timeout`: overall HTTP request timeout as a Go duration string (for example
  `30s`, `2m`, `1m30s`). Empty means no overall timeout is set. `--timeout` and
  `DARYAFT_TIMEOUT` override this value. Must be a positive duration when set.
- `no_color`: default no-color preference for the TUI.
- `no_tui`: when true, plain `daryaft` uses the non-TUI fallback unless command
  arguments or download flags are present.
- `theme`: TUI theme. Supported values are `default` and `mono`; `mono` uses
  monochrome styling like `no_color`.
- `animations`: reserved for future TUI polish. Stored and shown in config, but
  it does not currently change runtime behavior.
- `hyperlinks`: reserved for future OSC 8 hyperlink support. Stored and shown
  in config, but Daryaft does not currently emit terminal hyperlinks.

## Supported Keys

```text
download_dir string
retries int
resume bool
no_color bool
no_tui bool
theme string
animations bool
hyperlinks bool
user_agent string
timeout string
```

`retries` must be an integer from `0` through `20`. Boolean config set values
accept `true`, `1`, `yes`, `y`, `on`, `false`, `0`, `no`, `n`, and `off`,
case-insensitively.

`theme` must be `default` or `mono`. Invalid `theme` values passed through
`config set` or `DARYAFT_THEME` return a clear error.

`user_agent` accepts any string without control characters. Empty clears the
override and restores the built-in default.

`timeout` must be a positive Go duration string such as `30s`, `2m`, or
`1m30s`. Zero, negative, or unparseable values return a clear error. Empty
clears the timeout override (no overall client timeout).

## Security

Credentials, tokens, cookies, authorization headers, proxy URLs, and arbitrary
headers are intentionally not supported in persistent configuration. The strict
YAML decoder (`KnownFields: true`) rejects unknown keys, so attempting to write
`username`, `password`, or `authorization` to the config file produces a parse
error. Use `DARYAFT_USERNAME` and `DARYAFT_PASSWORD` environment variables for
credentials — they are never persisted.

Config saves are written through a temporary file in the same directory and
renamed into place with `0600` permissions. This keeps normal config updates
atomic on local filesystems and avoids leaving partially written config files
behind after successful saves.

Shell completion suggests these keys for `daryaft config get` and
`daryaft config set`. For boolean fields, `daryaft config set` also suggests
`true` and `false` as value completions.

Malformed YAML is reported as an error and is not silently ignored.

Use `daryaft doctor` to verify that the config path can be resolved, the config
directory is writable or appears creatable, and the effective config can be
loaded. Invalid YAML is reported as a critical doctor failure. The effective
download directory, including the built-in `~/Downloads` default, is reported
and checked. A download directory that does not exist is reported as a warning;
`doctor` does not create download directories.

## TUI Settings screen

The TUI includes a read-only Settings screen that shows the active config path, whether
a config file was loaded, and safe effective settings.

Open it from the TUI home screen by selecting **Settings** (menu item 6) or pressing
`c` from the home screen. Press `esc` or `backspace` to return home.

The screen shows:

```text
Config file: <path>
Config loaded: yes / no (using defaults)
---
download_dir, retries, resume, no_color, no_tui, theme,
animations, hyperlinks, user_agent, timeout
```

The screen is read-only. It does not edit config files. It does not display credentials,
tokens, cookies, auth headers, proxy values, or credential environment variables
(`DARYAFT_USERNAME` / `DARYAFT_PASSWORD`).

Related docs:

- [Command Reference](command-reference.md)
- [Usage](usage.md)
- [Doctor Troubleshooting](troubleshooting/doctor.md)
- [Architecture Overview](architecture/overview.md)
