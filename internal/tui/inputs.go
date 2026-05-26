package tui

import "github.com/charmbracelet/bubbles/textinput"

func newTextInput(s styles) textinput.Model {
	input := textinput.New()
	input.Prompt = "> "
	input.Width = 44
	input.CharLimit = 4096
	input.PromptStyle = s.muted
	input.TextStyle = s.body
	input.PlaceholderStyle = s.muted
	return input
}

func newOutputInput(s styles, value string) textinput.Model {
	input := newTextInput(s)
	input.Placeholder = "."
	if value != "." {
		input.SetValue(value)
	}
	return input
}

func (m Model) inputPrompt() string {
	switch m.screen {
	case screenURLInput:
		return "Enter download URL"
	case screenFileInput:
		return "Enter path to .txt file"
	case screenOutputInput:
		return "Enter output directory"
	default:
		return ""
	}
}
