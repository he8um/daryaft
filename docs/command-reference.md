# Command Reference

## `daryaft`

Implemented. Prints the current placeholder message and exits successfully.
Interactive TUI mode is planned.

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
- Go version

Default local values are used unless release tooling injects build variables.

## Common Flags

Implemented harmless placeholders:

- `--no-color`: disable colored output when colorized output exists.
- `--no-tui`: disable terminal UI when the TUI exists.
- `-v`, `--verbose`: enable verbose output when verbose logging exists.

## Planned Commands And Forms

These are planned and not implemented yet:

```bash
daryaft https://example.com/file.zip
daryaft -f urls.txt
daryaft update
```

Related docs:

- [Usage](usage.md)
- [Architecture Overview](architecture/overview.md)
