# Configuration

Runtime configuration is planned but not implemented yet.

## Built-In Metadata Defaults

The current code defines these defaults:

- App name: `Daryaft`
- Binary name: `daryaft`
- Website: `https://xhesam.com`
- Project page: `https://xhesam.com/daryaft`
- Repository: `https://github.com/he8um/daryaft`
- Footer: `Developed with <3 by AmirHesam Piri`

## Planned Config Locations

Daryaft should follow platform conventions:

- macOS: `~/Library/Application Support/daryaft/config.toml`
- Linux: `${XDG_CONFIG_HOME}/daryaft/config.toml` or `~/.config/daryaft/config.toml`

## Planned Defaults

- Color enabled when the terminal supports it
- TUI enabled for interactive terminals
- Downloads saved to the current working directory unless configured
- Conservative retries
- Checksum verification when expected values are provided

Related docs:

- [Command Reference](command-reference.md)
- [Architecture Overview](architecture/overview.md)
