package app

import "github.com/he8um/daryaft/internal/config"

// InteractivePlaceholder is used when the no-argument TUI is disabled.
func InteractivePlaceholder() string {
	return config.AppName + ` terminal UI is disabled.

Run ` + config.BinaryName + ` without --no-tui to open the interactive home screen,
or use ` + config.BinaryName + ` --help to see the current CLI commands and flags.

` + config.FooterText
}
