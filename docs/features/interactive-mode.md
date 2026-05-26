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

The download menu entries open Bubble Tea text input forms:

- Download from URL prompts `Enter download URL`.
- Download from .txt file prompts `Enter path to .txt file`.
- Both flows then prompt `Enter output directory`.

Pressing enter validates the input with the existing download planning logic and
opens the output directory input. Empty output means `.`, the current
directory. Absolute and relative output paths are accepted and are not created
during planning. Invalid URLs or file paths keep the user on the source input
screen and show a validation error.

Pressing enter on the plan screen starts a real download. The execution screen
uses the same downloader event stream as the CLI and supports both one URL and
sequential `.txt` batch downloads. It shows status, target path, byte progress,
percent when known, speed, retry/resume/restart messages, failures, and final
summary counts. TUI downloads use the current directory as the output path for
empty output, or the selected output directory when one is entered. Pressing
`q` while a download is running cancels it, keeps the `.part` file and metadata
sidecar for resume, and stops without retrying. CLI `-o`/`--output` behavior is
unchanged.

CLI forms remain fully supported:

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
- enter from source input: continue to output directory input
- enter from output directory input: review plan
- enter on plan: start download
- enter after completion: return home
- `esc` or backspace: return to the previous input screen, or home from a
  source input
- `q`: quit, or cancel when a download is running
- ctrl+c: quit from anywhere

## Boundaries

Interactive mode does not implement custom filename input, queue persistence,
concurrency, self-update, or packaging flows yet. CLI ctrl+c behavior is
unchanged and may terminate the process directly.

## Examples

```bash
daryaft
daryaft --no-color
daryaft --no-tui
```
