package tui

import (
	"context"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/pkg/version"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Options struct {
	NoColor           bool
	Theme             string
	DownloadDir       string
	Retries           int
	Resume            bool
	UseConfigDefaults bool
}

type Model struct {
	screen            screen
	selected          int
	styles            styles
	noColor           bool
	width             int
	height            int
	version           version.Details
	input             textinput.Model
	sourceInput       string
	sourceScreen      screen
	outputDirInput    string
	defaultOutputDir  string
	retries           int
	resume            bool
	filenameInput     string
	checksumInput     string
	errorMessage      string
	plan              download.Plan
	executionRunner   ExecutionRunner
	inspectRunner     InspectRunner
	execution         executionState
	executionCancel   context.CancelFunc
	executionMessages <-chan tea.Msg
	inspectInput      string
	inspect           inspectState
	inspectCancel     context.CancelFunc
	inspectMessages   <-chan tea.Msg
}

func NewModel(options Options) Model {
	return NewModelWithRunners(options, defaultExecutionRunner, defaultInspectRunner)
}

func NewModelWithRunner(options Options, runner ExecutionRunner) Model {
	return NewModelWithRunners(options, runner, defaultInspectRunner)
}

func NewModelWithRunners(options Options, runner ExecutionRunner, inspectRunner InspectRunner) Model {
	noColorMode := options.NoColor || config.IsMonoTheme(options.Theme)
	styles := newStyles(noColorMode)
	if runner == nil {
		runner = defaultExecutionRunner
	}
	if inspectRunner == nil {
		inspectRunner = defaultInspectRunner
	}
	defaultOutput := defaultOutputDir(options.DownloadDir)
	retries := tuiDefaultRetries
	resume := tuiDefaultResume
	if options.UseConfigDefaults {
		retries = options.Retries
		resume = options.Resume
	}
	model := Model{
		screen:           screenHome,
		styles:           styles,
		noColor:          noColorMode,
		version:          version.Info(),
		sourceScreen:     screenURLInput,
		outputDirInput:   defaultOutput,
		defaultOutputDir: defaultOutput,
		retries:          retries,
		resume:           resume,
		executionRunner:  runner,
		inspectRunner:    inspectRunner,
	}
	model.input = model.newTextInput()
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) moveUp() Model {
	if m.selected <= 0 {
		m.selected = len(homeMenu) - 1
		return m
	}
	m.selected--
	return m
}

func (m Model) moveDown() Model {
	if m.selected >= len(homeMenu)-1 {
		m.selected = 0
		return m
	}
	m.selected++
	return m
}

func (m Model) enter() (Model, tea.Cmd) {
	item := homeMenu[m.selected]
	if item.title == "Quit" {
		return m, tea.Quit
	}
	if item.screen == screenURLInput || item.screen == screenFileInput || item.screen == screenInspectInput {
		return m.openInput(item.screen)
	}
	m.screen = item.screen
	return m, nil
}

func (m Model) openInput(next screen) (Model, tea.Cmd) {
	m.screen = next
	m.input = m.newTextInput()
	m.sourceInput = ""
	m.sourceScreen = next
	m.outputDirInput = m.defaultOutputDir
	m.filenameInput = ""
	m.checksumInput = ""
	m.errorMessage = ""
	m.plan = download.Plan{}
	m.execution = executionState{}
	m.executionCancel = nil
	m.executionMessages = nil
	m.inspectInput = ""
	m.inspect = inspectState{}
	m.inspectCancel = nil
	m.inspectMessages = nil
	return m, m.input.Focus()
}

func (m Model) back() (Model, tea.Cmd) {
	if m.screen == screenPlan {
		if m.sourceScreen == screenURLInput {
			m.screen = screenChecksumInput
			m.input = m.newChecksumInput(m.checksumInput)
		} else {
			m.screen = screenOutputInput
			m.input = m.newOutputInput(m.outputDirInput)
		}
		m.errorMessage = ""
		return m, m.input.Focus()
	}
	if m.screen == screenFilenameInput {
		m.screen = screenOutputInput
		m.input = m.newOutputInput(m.outputDirInput)
		m.errorMessage = ""
		return m, m.input.Focus()
	}
	if m.screen == screenChecksumInput {
		m.screen = screenFilenameInput
		m.input = m.newFilenameInput(m.filenameInput)
		m.errorMessage = ""
		return m, m.input.Focus()
	}
	if m.screen == screenOutputInput {
		m.screen = m.sourceScreen
		m.input = m.newTextInput()
		m.input.SetValue(m.sourceInput)
		m.errorMessage = ""
		return m, m.input.Focus()
	}
	if m.screen == screenInspectResult || m.screen == screenInspectError {
		m.screen = screenInspectInput
		m.input = m.newTextInput()
		m.input.SetValue(m.inspectInput)
		m.errorMessage = ""
		m.inspect = inspectState{}
		m.inspectCancel = nil
		m.inspectMessages = nil
		return m, m.input.Focus()
	}
	m.screen = screenHome
	m.errorMessage = ""
	return m, nil
}

func (m Model) home() Model {
	m.screen = screenHome
	m.errorMessage = ""
	m.execution = executionState{}
	m.executionCancel = nil
	m.executionMessages = nil
	m.inspect = inspectState{}
	m.inspectCancel = nil
	m.inspectMessages = nil
	return m
}

func (m Model) isInputScreen() bool {
	return m.screen == screenURLInput || m.screen == screenFileInput || m.screen == screenInspectInput || m.screen == screenOutputInput || m.screen == screenFilenameInput || m.screen == screenChecksumInput
}
