package tui

type screen int

const (
	screenHome screen = iota
	screenURLInput
	screenFileInput
	screenPlan
	screenExecution
	screenHelp
	screenVersion
)

type menuItem struct {
	title  string
	screen screen
}

var homeMenu = []menuItem{
	{title: "Download from URL", screen: screenURLInput},
	{title: "Download from .txt file", screen: screenFileInput},
	{title: "View help", screen: screenHelp},
	{title: "Version", screen: screenVersion},
	{title: "Quit", screen: screenHome},
}

func (s screen) title() string {
	switch s {
	case screenURLInput:
		return "Download from URL"
	case screenFileInput:
		return "Download from .txt file"
	case screenPlan:
		return "Download plan"
	case screenExecution:
		return "Downloading"
	case screenHelp:
		return "Help"
	case screenVersion:
		return "Version"
	default:
		return "Home"
	}
}
