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
- `DARYAFT_NO_COLOR`
- `DARYAFT_NO_TUI`
- `DARYAFT_THEME`
- `DARYAFT_ANIMATIONS`
- `DARYAFT_HYPERLINKS`

String values are trimmed. Empty string values are accepted for
`DARYAFT_DOWNLOAD_DIR` and `DARYAFT_THEME`.

`DARYAFT_RETRIES` must be an integer greater than or equal to `0`. Empty or
invalid integer values return an error.

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

- `download_dir`: default output directory. Empty means current directory unless
  the CLI or TUI explicitly sets output.
- `retries`: default retry attempts after the initial attempt.
- `resume`: default resume behavior for interrupted `.part` files.
- `no_color`: default no-color preference for the TUI.
- `no_tui`: when true, plain `daryaft` uses the non-TUI fallback unless command
  arguments or download flags are present.
- `theme`: stored for future theme support. Only `default` is meaningful now.
- `animations`: stored for future TUI polish.
- `hyperlinks`: stored for future TUI polish.

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
```

`retries` must be an integer greater than or equal to `0`. Boolean config set
values accept `true`, `1`, `yes`, `y`, `on`, `false`, `0`, `no`, `n`, and
`off`, case-insensitively.

Malformed YAML is reported as an error and is not silently ignored.

Related docs:

- [Command Reference](command-reference.md)
- [Usage](usage.md)
- [Architecture Overview](architecture/overview.md)
