package tui

import "github.com/charmbracelet/bubbles/textinput"

const (
	defaultPanelWidth = 56
	minPanelWidth     = 40
	maxPanelWidth     = 80
	minInputWidth     = 20
	panelFrameWidth   = 10
	inputFrameWidth   = 12
)

func newTextInput(s styles, width int) textinput.Model {
	input := textinput.New()
	input.Prompt = "> "
	input.Width = width
	input.CharLimit = 4096
	input.PromptStyle = s.muted
	input.TextStyle = s.body
	input.PlaceholderStyle = s.muted
	return input
}

func newOutputInput(s styles, value string, width int) textinput.Model {
	input := newTextInput(s, width)
	input.Placeholder = "."
	if value != "." {
		input.SetValue(value)
	}
	return input
}

func newFilenameInput(s styles, value string, width int) textinput.Model {
	input := newTextInput(s, width)
	input.Placeholder = "auto-detect"
	input.SetValue(value)
	return input
}

func newChecksumInput(s styles, value string, width int) textinput.Model {
	input := newTextInput(s, width)
	input.Placeholder = "sha256:<hex>"
	input.SetValue(value)
	return input
}

func (m Model) newTextInput() textinput.Model {
	return newTextInput(m.styles, m.inputWidth())
}

func (m Model) newOutputInput(value string) textinput.Model {
	return newOutputInput(m.styles, value, m.inputWidth())
}

func (m Model) newFilenameInput(value string) textinput.Model {
	return newFilenameInput(m.styles, value, m.inputWidth())
}

func (m Model) newChecksumInput(value string) textinput.Model {
	return newChecksumInput(m.styles, value, m.inputWidth())
}

func (m Model) panelWidth() int {
	if m.width <= 0 {
		return defaultPanelWidth
	}
	return clamp(m.width-panelFrameWidth, minPanelWidth, maxPanelWidth)
}

func (m Model) inputWidth() int {
	return max(minInputWidth, m.panelWidth()-inputFrameWidth)
}

func (m Model) withWindowSize(width, height int) Model {
	m.width = width
	m.height = height
	m.input.Width = m.inputWidth()
	return m
}

func clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func (m Model) inputPrompt() string {
	switch m.screen {
	case screenURLInput:
		return "Enter a download URL (https:// or http://)"
	case screenFileInput:
		return "Enter the absolute path to a .txt file with one URL per line"
	case screenInspectInput:
		return "Enter URL to inspect"
	case screenOutputInput:
		return "Enter output directory"
	case screenFilenameInput:
		return "Enter custom filename"
	case screenChecksumInput:
		return "Enter checksum"
	default:
		return ""
	}
}
