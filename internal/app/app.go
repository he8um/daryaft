package app

import "github.com/he8um/daryaft/internal/config"

// InteractivePlaceholder is the current no-argument startup message.
func InteractivePlaceholder() string {
	return config.AppName + ` is starting its terminal downloader foundation.

Interactive mode is planned for the TUI milestone and is not implemented yet.
Use ` + config.BinaryName + ` --help to see the current commands and flags.

` + config.FooterText
}
