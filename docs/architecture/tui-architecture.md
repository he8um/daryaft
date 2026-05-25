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
- `inputs.go`: Bubble Tea text input construction and prompt selection
- `plan.go`: TUI dry-run plan helpers backed by `internal/download`
- `messages.go`: Bubble Tea message types for downloader execution events
- `execution.go`: goroutine/channel bridge from `internal/downloader` to the
  TUI model

`cmd/root.go` calls `tui.Run` only for no-argument execution. URL arguments,
`--file`, and download flags continue through the existing CLI download path.

## Model

The model tracks:

- current screen
- selected home menu index
- Bubble Tea text input state
- validation error message
- generated dry-run download plan
- input screen to return to from the plan screen
- execution state for the current item, event data, messages, and summary
- style set
- version details

The first implemented screens are:

- home
- Download from URL input
- Download from .txt file input
- dry-run download plan
- execution/progress
- help
- version

The TUI calls `internal/download.BuildPlan` for URL and file inputs. It does
not import Cobra. The plan screen starts downloads through
`internal/downloader.DownloadBatch`, so single URL and `.txt` batch execution
share the same sequential runner as the CLI. TUI plans do not expose output path
input yet, so execution uses the current directory.

## Styling

Lip Gloss renders a bordered panel, highlighted selected menu item, muted
footer, and simple body text. `--no-color` builds styles without foreground or
background colors while preserving layout.

## Event Boundary

TUI execution starts a goroutine that runs the existing downloader. Downloader
item and progress events are copied into a channel and received by Bubble Tea
commands, keeping the update loop non-blocking and avoiding shared mutable
state. The TUI uses those messages to render status, target path, downloaded
bytes, percent, speed, retry/resume/restart messages, completion, failure, and
batch summaries.

Cancellation is not implemented yet. Pressing `q` while a download is running
updates the message area instead of silently quitting; ctrl+c still exits the
program.

## Testing

Current tests cover menu navigation, wrap behavior, screen switching, input
validation, dry-run plan creation, execution transitions, event message
handling, summary rendering, running-state quit behavior, back navigation,
footer rendering, and version rendering without brittle snapshots.

## Examples

```bash
daryaft
daryaft --no-color
daryaft https://example.com/file.zip --dry-run
```
