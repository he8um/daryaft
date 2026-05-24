package tui

type screen int

const (
	screenHome screen = iota
	screenURLPlanned
	screenFilePlanned
	screenHelp
	screenVersion
)

type menuItem struct {
	title  string
	screen screen
}

var homeMenu = []menuItem{
	{title: "Download from URL", screen: screenURLPlanned},
	{title: "Download from .txt file", screen: screenFilePlanned},
	{title: "View help", screen: screenHelp},
	{title: "Version", screen: screenVersion},
	{title: "Quit", screen: screenHome},
}

func (s screen) title() string {
	switch s {
	case screenURLPlanned:
		return "Download from URL"
	case screenFilePlanned:
		return "Download from .txt file"
	case screenHelp:
		return "Help"
	case screenVersion:
		return "Version"
	default:
		return "Home"
	}
}
