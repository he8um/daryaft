# Feature: Interactive Mode

Running `daryaft` without arguments opens the first interactive home screen.

## Navigation

- [Documentation Index](../index.md)
- [Implementation Plan](../engineering/implementation-plan.md)
- [Agent Navigation](../engineering/agent-navigation.md)
- [Versioning Policy](../roadmap/versioning-policy.md)
- [Release Train](../roadmap/release-train.md)

## Project constants

```text
Project: Daryaft
Binary: daryaft
Module: github.com/he8um/daryaft
Repository: git@github.com:he8um/daryaft.git
Author: AmirHesam Piri <info@xhesam.com>
Website: https://xhesam.com
Project page: https://xhesam.com/daryaft
License: MIT
Footer: Developed with <3 by AmirHesam Piri
```

## Current Behavior

`daryaft` launches a Bubble Tea home screen. The current menu is intentionally
small:

1. Download from URL
2. Download from .txt file
3. View help
4. Version
5. Quit

The download menu entries are placeholders for now. They show a planned screen
inside the TUI and do not start downloads. Use these CLI forms for real
downloads:

```bash
daryaft https://example.com/file.zip
daryaft -f urls.txt
daryaft download https://example.com/file.zip
```

The help screen documents navigation keys. The version screen shows the same
build metadata as `daryaft version`.

## Keys

- up/down arrows: move selection
- `k`/`j`: move selection
- enter: select
- `esc` or backspace: return to home from a sub-screen
- `q` or ctrl+c: quit from anywhere

## Boundaries

Interactive mode does not implement download execution, queue persistence,
concurrency, self-update, or packaging flows yet. It is the foundation that will
later consume downloader events.

## Examples

```bash
daryaft
daryaft --no-color
daryaft --no-tui
```
