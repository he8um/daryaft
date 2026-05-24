# TUI Architecture

Bubble Tea screens, models, styles, and footer behavior.

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

## Current Package

The TUI lives in `internal/tui` and is independent from Cobra except for the
small call from `cmd/root.go`.

Current files:

- `app.go`: starts the Bubble Tea program
- `model.go`: model state and menu movement helpers
- `update.go`: keyboard update logic
- `view.go`: home and sub-screen rendering
- `styles.go`: Lip Gloss style construction
- `screens.go`: screen and menu definitions
- `keys.go`: key classification helpers

`cmd/root.go` calls `tui.Run` only for no-argument execution. URL arguments,
`--file`, and download flags continue through the existing CLI download path.

## Model

The model tracks:

- current screen
- selected home menu index
- style set
- version details

The first implemented screens are:

- home
- planned Download from URL
- planned Download from .txt file
- help
- version

## Styling

Lip Gloss renders a bordered panel, highlighted selected menu item, muted
footer, and simple body text. `--no-color` builds styles without foreground or
background colors while preserving layout.

## Event Boundary

The current TUI does not start downloads. Future download screens should
subscribe to the existing downloader event stream instead of reaching into
downloader internals.

## Testing

Current tests cover menu navigation, wrap behavior, screen switching, back
navigation, footer rendering, and version rendering without brittle snapshots.

## Examples

```bash
daryaft
daryaft --no-color
daryaft https://example.com/file.zip --dry-run
```
