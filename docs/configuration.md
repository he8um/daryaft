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
```

- `daryaft config` shows help for the config command group.
- `daryaft config path` prints the config file path.
- `daryaft config show` prints the effective config as YAML. If no file exists,
  it prints built-in defaults.
- `daryaft config init` creates a config file with defaults and refuses to
  overwrite an existing file.
- `daryaft config init --force` overwrites the existing config with defaults.

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
2. Config file values
3. Built-in defaults

Examples:

- `-o`/`--output` wins over `download_dir`.
- `--retries` wins over `retries`.
- `--resume` and `--no-resume` win over `resume`.
- `--no-color` wins over `no_color`.
- `--no-tui` wins over `no_tui`.

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

Malformed YAML is reported as an error and is not silently ignored.

Related docs:

- [Command Reference](command-reference.md)
- [Usage](usage.md)
- [Architecture Overview](architecture/overview.md)
