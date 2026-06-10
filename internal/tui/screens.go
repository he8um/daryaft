package tui

type screen int

const (
	screenHome screen = iota
	screenURLInput
	screenFileInput
	screenInspectInput
	screenOutputInput
	screenFilenameInput
	screenChecksumInput
	screenPlan
	screenExecution
	screenInspectExecution
	screenInspectResult
	screenInspectError
	screenHelp
	screenVersion
	screenSettings
)

type menuItem struct {
	title  string
	screen screen
}

var homeMenu = []menuItem{
	{title: "Download from URL", screen: screenURLInput},
	{title: "Download from .txt file", screen: screenFileInput},
	{title: "Inspect URL", screen: screenInspectInput},
	{title: "View help", screen: screenHelp},
	{title: "Version", screen: screenVersion},
	{title: "Settings", screen: screenSettings},
	{title: "Quit", screen: screenHome},
}

func (s screen) title() string {
	switch s {
	case screenURLInput:
		return "Download from URL"
	case screenFileInput:
		return "Download from .txt file"
	case screenInspectInput:
		return "Inspect URL"
	case screenOutputInput:
		return "Output directory"
	case screenFilenameInput:
		return "Custom filename"
	case screenChecksumInput:
		return "Checksum"
	case screenPlan:
		return "Download plan"
	case screenExecution:
		return "Downloading"
	case screenInspectExecution:
		return "Inspecting"
	case screenInspectResult:
		return "Inspect result"
	case screenInspectError:
		return "Inspect error"
	case screenHelp:
		return "Help"
	case screenVersion:
		return "Version"
	case screenSettings:
		return "Settings"
	default:
		return "Home"
	}
}
