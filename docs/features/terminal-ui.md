# Feature: Terminal UI

The Bubble Tea terminal UI is implemented for no-argument startup. It can
collect URL or `.txt` file input and show dry-run download plans. Full download
execution and progress screens are planned next.

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

Running `daryaft` with no URL arguments opens an interactive home screen. It is
implemented with Bubble Tea and styled with Lip Gloss.

The home screen shows:

- app title: `Daryaft`
- subtitle: `Modern terminal downloader`
- menu items:
  1. Download from URL
  2. Download from .txt file
  3. View help
  4. Version
  5. Quit
- footer: `Developed with <3 by AmirHesam Piri`

Navigation:

- up/down arrows or `k`/`j` move through menu items
- enter selects an item
- `esc` or backspace returns from sub-screens to home
- `q` or ctrl+c quits from anywhere

The View help and Version menu items render simple in-TUI screens. Download
from URL and Download from .txt file render text input forms. Pressing enter
validates the URL or URL file through the existing download planner and shows a
dry-run plan with URL count, the first URLs, output, filename, retries, and
resume settings.

The TUI does not execute downloads yet. CLI commands remain the execution path
for real downloads:

```bash
daryaft https://example.com/file.zip
daryaft -f urls.txt
daryaft download https://example.com/file.zip
```

`--no-color` keeps the panel layout but avoids color styling. `--no-tui` skips
the no-argument TUI and prints the non-interactive placeholder.

Existing CLI download commands are unchanged and remain the stable way to
download files.

## Planned

Future TUI work will add execution and a progress screen that consumes
downloader events for progress, retry, resume, failure, and completion states.
Queue persistence, concurrent downloads, history, and rich progress bars are not
implemented yet.

## Examples

```bash
daryaft -h
daryaft
daryaft https://example.com/file.zip
daryaft -f urls.txt
```
