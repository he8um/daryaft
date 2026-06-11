package tui

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/downloader"
	"github.com/he8um/daryaft/internal/inspect"
	"github.com/he8um/daryaft/pkg/version"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMenuNavigationWraps(t *testing.T) {
	model := NewModel(Options{NoColor: true})

	model = updateWithKey(t, model, tea.KeyUp)
	if model.selected != len(homeMenu)-1 {
		t.Fatalf("selected after up = %d, want %d", model.selected, len(homeMenu)-1)
	}

	model = updateWithKey(t, model, tea.KeyDown)
	if model.selected != 0 {
		t.Fatalf("selected after down = %d, want 0", model.selected)
	}
}

func TestJKNavigation(t *testing.T) {
	model := NewModel(Options{NoColor: true})

	model = updateWithRune(t, model, 'j')
	if model.selected != 1 {
		t.Fatalf("selected after j = %d, want 1", model.selected)
	}

	model = updateWithRune(t, model, 'k')
	if model.selected != 0 {
		t.Fatalf("selected after k = %d, want 0", model.selected)
	}
}

func TestScreenSwitchingAndBack(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 3

	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenHelp {
		t.Fatalf("screen = %v, want help", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyEsc)
	if model.screen != screenHome {
		t.Fatalf("screen after esc = %v, want home", model.screen)
	}

	model.selected = 4
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenVersion {
		t.Fatalf("screen = %v, want version", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenHome {
		t.Fatalf("screen after backspace = %v, want home", model.screen)
	}

	model.selected = 5
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenSettings {
		t.Fatalf("screen = %v, want settings", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyEsc)
	if model.screen != screenHome {
		t.Fatalf("screen after esc from settings = %v, want home", model.screen)
	}
}

func TestSelectURLInputScreen(t *testing.T) {
	model := NewModel(Options{NoColor: true})

	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenURLInput {
		t.Fatalf("screen = %v, want URL input", model.screen)
	}
	if !strings.Contains(model.View(), "Enter a download URL") {
		t.Fatalf("URL input view missing prompt:\n%s", model.View())
	}
}

func TestWindowSizeUpdatesResponsiveDimensions(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model = updateWithMsg(t, model, tea.WindowSizeMsg{Width: 90, Height: 24})

	if model.width != 90 || model.height != 24 {
		t.Fatalf("size = %dx%d, want 90x24", model.width, model.height)
	}
	if got := model.panelWidth(); got != 80 {
		t.Fatalf("panelWidth = %d, want 80", got)
	}
	if model.input.Width != 68 {
		t.Fatalf("input width = %d, want 68", model.input.Width)
	}
}

func TestResponsiveWidthsUseMinimumAndMaximum(t *testing.T) {
	model := NewModel(Options{NoColor: true})

	narrow := updateWithMsg(t, model, tea.WindowSizeMsg{Width: 35, Height: 20})
	if got := narrow.panelWidth(); got != minPanelWidth {
		t.Fatalf("narrow panelWidth = %d, want %d", got, minPanelWidth)
	}
	if narrow.input.Width != minPanelWidth-inputFrameWidth {
		t.Fatalf("narrow input width = %d", narrow.input.Width)
	}

	wide := updateWithMsg(t, model, tea.WindowSizeMsg{Width: 200, Height: 50})
	if got := wide.panelWidth(); got != maxPanelWidth {
		t.Fatalf("wide panelWidth = %d, want %d", got, maxPanelWidth)
	}
	if wide.input.Width != maxPanelWidth-inputFrameWidth {
		t.Fatalf("wide input width = %d", wide.input.Width)
	}
}

func TestInvalidURLShowsError(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "ftp://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenURLInput {
		t.Fatalf("screen = %v, want URL input", model.screen)
	}
	if model.errorMessage == "" {
		t.Fatal("errorMessage is empty, want validation error")
	}
	if !strings.Contains(model.View(), "scheme must be http or https") {
		t.Fatalf("URL input view missing validation error:\n%s", model.View())
	}
}

func TestValidURLAdvancesToOutputInput(t *testing.T) {
	home := t.TempDir()
	t.Cleanup(config.SetUserHomeDirForTest(home))
	wantOutput := filepath.Join(home, "Downloads")

	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenOutputInput {
		t.Fatalf("screen = %v, want output input", model.screen)
	}
	if len(model.plan.URLs) != 1 {
		t.Fatalf("plan URL count = %d, want 1", len(model.plan.URLs))
	}
	if !strings.Contains(model.View(), "Enter output directory") {
		t.Fatalf("output input view missing prompt:\n%s", model.View())
	}
	if model.outputDirInput != wantOutput {
		t.Fatalf("outputDirInput = %q, want %q", model.outputDirInput, wantOutput)
	}
}

func TestEmptyOutputDirectoryCreatesPlanWithDownloadsDefault(t *testing.T) {
	home := t.TempDir()
	t.Cleanup(config.SetUserHomeDirForTest(home))
	wantOutput := filepath.Join(home, "Downloads")

	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenFilenameInput {
		t.Fatalf("screen = %v, want filename input", model.screen)
	}
	if !strings.Contains(model.View(), "Leave empty to auto-detect") {
		t.Fatalf("filename input view missing help:\n%s", model.View())
	}
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenChecksumInput {
		t.Fatalf("screen = %v, want checksum input", model.screen)
	}
	if !strings.Contains(model.View(), "Leave empty to skip") {
		t.Fatalf("checksum input view missing help:\n%s", model.View())
	}
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if model.plan.Output != wantOutput {
		t.Fatalf("plan.Output = %q, want %q", model.plan.Output, wantOutput)
	}
	if !strings.Contains(model.View(), "Number of URLs: 1") {
		t.Fatalf("plan view missing URL count:\n%s", model.View())
	}
	if !strings.Contains(model.View(), "Filename: auto-detect") {
		t.Fatalf("plan view missing auto-detect filename:\n%s", model.View())
	}
	if !strings.Contains(model.View(), "Checksum: none") {
		t.Fatalf("plan view missing empty checksum:\n%s", model.View())
	}
	if !strings.Contains(model.View(), "enter start download") {
		t.Fatalf("plan view missing start action:\n%s", model.View())
	}
}

func TestCustomOutputDirectoryAppearsInPlan(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model.input.SetValue("")
	model = updateWithString(t, model, "downloads")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if model.plan.Output != "downloads" {
		t.Fatalf("plan.Output = %q, want downloads", model.plan.Output)
	}
	if !strings.Contains(model.View(), "Output: downloads") {
		t.Fatalf("plan view missing custom output:\n%s", model.View())
	}
}

func TestConfigDefaultsApplyToTUIPlan(t *testing.T) {
	model := NewModel(Options{
		NoColor:           true,
		DownloadDir:       "/tmp/daryaft-configured",
		Retries:           0,
		Resume:            false,
		UseConfigDefaults: true,
	})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenOutputInput {
		t.Fatalf("screen = %v, want output input", model.screen)
	}
	if model.outputDirInput != "/tmp/daryaft-configured" {
		t.Fatalf("outputDirInput = %q", model.outputDirInput)
	}
	if model.input.Value() != "/tmp/daryaft-configured" {
		t.Fatalf("input value = %q", model.input.Value())
	}

	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if model.plan.Output != "/tmp/daryaft-configured" {
		t.Fatalf("plan.Output = %q", model.plan.Output)
	}
	if model.plan.Retries != 0 {
		t.Fatalf("plan.Retries = %d, want 0", model.plan.Retries)
	}
	if model.plan.Resume {
		t.Fatal("plan.Resume = true, want false")
	}
}

func TestURLFlowAdvancesThroughFilenameInput(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0

	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenURLInput {
		t.Fatalf("screen = %v, want URL input", model.screen)
	}

	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenOutputInput {
		t.Fatalf("screen = %v, want output input", model.screen)
	}

	model.input.SetValue("")
	model = updateWithString(t, model, "downloads")
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenFilenameInput {
		t.Fatalf("screen = %v, want filename input", model.screen)
	}
	if !strings.Contains(model.View(), "Enter custom filename") {
		t.Fatalf("filename input view missing prompt:\n%s", model.View())
	}

	model = updateWithString(t, model, "custom.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenChecksumInput {
		t.Fatalf("screen = %v, want checksum input", model.screen)
	}
	if model.plan.Output != "downloads" {
		t.Fatalf("plan.Output = %q, want downloads", model.plan.Output)
	}
	if model.plan.Name != "custom.zip" {
		t.Fatalf("plan.Name = %q, want custom.zip", model.plan.Name)
	}
	if !strings.Contains(model.View(), "Enter checksum") {
		t.Fatalf("checksum input view missing prompt:\n%s", model.View())
	}
}

func TestCustomFilenameAppearsInPlan(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "  custom.txt  ")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if model.plan.Name != "custom.txt" {
		t.Fatalf("plan.Name = %q, want custom.txt", model.plan.Name)
	}
	if !strings.Contains(model.View(), "Filename: custom.txt") {
		t.Fatalf("plan view missing custom filename:\n%s", model.View())
	}
}

func TestValidChecksumAppearsInPlan(t *testing.T) {
	checksum := "sha256:" + strings.Repeat("a", 64)
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, checksum)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if model.plan.Checksum == nil {
		t.Fatal("plan.Checksum = nil")
	}
	if model.plan.Checksum.String() != checksum {
		t.Fatalf("plan.Checksum = %q, want %q", model.plan.Checksum.String(), checksum)
	}
	if !strings.Contains(model.View(), "Checksum:") || !strings.Contains(model.View(), "sha256:") {
		t.Fatalf("plan view missing checksum:\n%s", model.View())
	}
}

func TestInvalidChecksumStaysOnChecksumInput(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "sha256:not-hex")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenChecksumInput {
		t.Fatalf("screen = %v, want checksum input", model.screen)
	}
	if model.errorMessage == "" {
		t.Fatal("errorMessage is empty, want checksum validation error")
	}
	if !strings.Contains(model.View(), "sha256 checksum must be 64 hex characters") {
		t.Fatalf("checksum input view missing validation error:\n%s", model.View())
	}
}

func TestUnsafeFilenameShowsValidationError(t *testing.T) {
	for _, value := range []string{"../file.zip", `nested\file.zip`, ".", ".."} {
		t.Run(value, func(t *testing.T) {
			model := NewModel(Options{NoColor: true})
			model.selected = 0
			model = updateWithKey(t, model, tea.KeyEnter)
			model = updateWithString(t, model, "https://example.com/file.zip")
			model = updateWithKey(t, model, tea.KeyEnter)
			model = updateWithKey(t, model, tea.KeyEnter)
			model = updateWithString(t, model, value)
			model = updateWithKey(t, model, tea.KeyEnter)

			if model.screen != screenFilenameInput {
				t.Fatalf("screen = %v, want filename input", model.screen)
			}
			if model.errorMessage == "" {
				t.Fatal("errorMessage is empty, want validation error")
			}
			if !strings.Contains(model.View(), "Error:") {
				t.Fatalf("filename view missing validation error:\n%s", model.View())
			}
		})
	}
}

func TestSelectFileInputScreen(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenFileInput {
		t.Fatalf("screen = %v, want file input", model.screen)
	}
	if !strings.Contains(model.View(), ".txt file") {
		t.Fatalf("file input view missing prompt:\n%s", model.View())
	}
}

func TestHomeMenuIncludesInspectURL(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	view := model.View()

	for _, want := range []string{
		"Download from URL",
		"Download from .txt file",
		"Inspect URL",
		"View help",
		"Version",
		"Settings",
		"Quit",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("home view missing %q:\n%s", want, view)
		}
	}
	if len(homeMenu) != 7 {
		t.Fatalf("homeMenu len = %d, want 7", len(homeMenu))
	}
	if homeMenu[2].title != "Inspect URL" {
		t.Fatalf("homeMenu[2] = %q, want Inspect URL", homeMenu[2].title)
	}
	if homeMenu[5].title != "Settings" {
		t.Fatalf("homeMenu[5] = %q, want Settings", homeMenu[5].title)
	}
}

func TestSelectInspectURLInputScreen(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 2
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenInspectInput {
		t.Fatalf("screen = %v, want inspect input", model.screen)
	}
	if !strings.Contains(model.View(), "Enter URL to inspect") {
		t.Fatalf("inspect input view missing prompt:\n%s", model.View())
	}
}

func TestInvalidInspectURLShowsValidationError(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 2
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "ftp://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenInspectInput {
		t.Fatalf("screen = %v, want inspect input", model.screen)
	}
	if model.errorMessage == "" {
		t.Fatal("errorMessage is empty, want validation error")
	}
	if !strings.Contains(model.View(), "scheme must be http or https") {
		t.Fatalf("inspect input view missing validation error:\n%s", model.View())
	}
}

func TestValidInspectURLStartsInjectedRunner(t *testing.T) {
	var received string
	runner := func(ctx context.Context, rawURL string) (inspect.Result, error) {
		received = rawURL
		return sampleInspectResult(rawURL), nil
	}
	model := NewModelWithRunners(Options{NoColor: true}, nil, runner)
	model.selected = 2
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")

	model, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if model.screen != screenInspectExecution {
		t.Fatalf("screen = %v, want inspect execution", model.screen)
	}
	if !model.inspect.Running {
		t.Fatal("inspect.Running = false, want true")
	}
	if cmd == nil {
		t.Fatal("inspect command is nil")
	}

	model, _ = updateWithInspectCmd(t, model, cmd)
	if received != "https://example.com/file.zip" {
		t.Fatalf("runner URL = %q", received)
	}
	if model.screen != screenInspectResult {
		t.Fatalf("screen = %v, want inspect result", model.screen)
	}
	if model.inspect.Result.Filename != "file.zip" {
		t.Fatalf("filename = %q", model.inspect.Result.Filename)
	}
}

func TestInspectResultScreenRendersKeyFields(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenInspectResult
	model.inspect = inspectState{
		Done:   true,
		Status: "Completed",
		Result: inspect.Result{
			URL:                "https://example.com/file.zip",
			FinalURL:           "https://cdn.example.com/file.zip",
			Status:             "200 OK",
			StatusCode:         200,
			Filename:           "file.zip",
			ContentLength:      1048576,
			ContentLengthKnown: true,
			ContentType:        "application/zip",
			AcceptRanges:       "bytes",
			ResumeSupported:    true,
			ResumeSupportKnown: true,
			ETag:               `"abc123"`,
			LastModified:       "Tue, 01 Jun 2026 12:00:00 GMT",
		},
	}

	view := model.View()
	for _, want := range []string{
		"Inspect result",
		"URL: https://example.com/file.zip",
		"Final URL: https://cdn.example.com/file.zip",
		"Status: 200 OK",
		"Filename: file.zip",
		"Content length: 1048576 bytes",
		"Content type: application/zip",
		"Accept-Ranges: bytes",
		"Resume supported: yes",
		`ETag: "abc123"`,
		"Last-Modified: Tue, 01 Jun 2026 12:00:00 GMT",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("inspect result missing %q in:\n%s", want, view)
		}
	}
}

func TestInspectErrorScreenRendersError(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenInspectError
	model.inspect = inspectState{Error: "Inspect failed: timeout"}

	view := model.View()
	if !strings.Contains(view, "Inspect error") || !strings.Contains(view, "Inspect failed: timeout") {
		t.Fatalf("inspect error view missing error:\n%s", view)
	}
}

func TestInspectRunnerErrorShowsErrorScreenAndCanNavigateBack(t *testing.T) {
	runner := func(ctx context.Context, rawURL string) (inspect.Result, error) {
		return inspect.Result{}, context.DeadlineExceeded
	}
	model := NewModelWithRunners(Options{NoColor: true}, nil, runner)
	model.selected = 2
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)

	model, _ = updateWithInspectCmd(t, model, cmd)
	if model.screen != screenInspectError {
		t.Fatalf("screen = %v, want inspect error", model.screen)
	}
	if !strings.Contains(model.View(), "Inspect failed") {
		t.Fatalf("inspect error view missing failure:\n%s", model.View())
	}

	homeModel := updateWithRune(t, model, 'h')
	if homeModel.screen != screenHome {
		t.Fatalf("screen after h = %v, want home", homeModel.screen)
	}

	inputModel := updateWithKey(t, model, tea.KeyBackspace)
	if inputModel.screen != screenInspectInput {
		t.Fatalf("screen after backspace = %v, want inspect input", inputModel.screen)
	}
	if inputModel.input.Value() != "https://example.com/file.zip" {
		t.Fatalf("input value = %q", inputModel.input.Value())
	}
}

func TestBackspaceEditsInspectInput(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 2
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "abc")
	model = updateWithKey(t, model, tea.KeyBackspace)

	if model.screen != screenInspectInput {
		t.Fatalf("screen = %v, want inspect input", model.screen)
	}
	if model.input.Value() != "ab" {
		t.Fatalf("input value = %q, want ab", model.input.Value())
	}

	model.input.SetValue("")
	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenHome {
		t.Fatalf("screen after empty backspace = %v, want home", model.screen)
	}
}

func TestInspectDoesNotInvokeDownloadRunner(t *testing.T) {
	downloadInvoked := false
	downloadRunner := func(ctx context.Context, plan download.Plan, handlers downloader.BatchHandlers) downloader.BatchResult {
		downloadInvoked = true
		return downloader.BatchResult{Planned: len(plan.URLs)}
	}
	inspectRunner := func(ctx context.Context, rawURL string) (inspect.Result, error) {
		return sampleInspectResult(rawURL), nil
	}

	model := NewModelWithRunners(Options{NoColor: true}, downloadRunner, inspectRunner)
	model.selected = 2
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("inspect command is nil")
	}
	_, _ = updateWithInspectCmd(t, model, cmd)

	if downloadInvoked {
		t.Fatal("download runner was invoked for inspect")
	}
}

func TestQWhileInspectRunningDoesNotQuit(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenInspectExecution
	model.inspect = inspectState{Running: true}
	cancelled := false
	model.inspectCancel = func() {
		cancelled = true
	}

	model, cmd := updateWithRuneAndCmd(t, model, 'q')
	if cmd != nil {
		t.Fatal("q while inspect running returned a command, want nil")
	}
	if !cancelled {
		t.Fatal("inspect cancel func was not called")
	}
	if model.inspect.Status != "Cancelling" {
		t.Fatalf("status = %q, want Cancelling", model.inspect.Status)
	}
	if model.inspect.Message != "Cancelling..." {
		t.Fatalf("message = %q, want Cancelling...", model.inspect.Message)
	}
}

func TestInvalidFilePathShowsError(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, filepath.Join(t.TempDir(), "missing.txt"))
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenFileInput {
		t.Fatalf("screen = %v, want file input", model.screen)
	}
	if model.errorMessage == "" {
		t.Fatal("errorMessage is empty, want validation error")
	}
	if !strings.Contains(model.View(), "read URL file") {
		t.Fatalf("file input view missing validation error:\n%s", model.View())
	}
}

func TestValidFileAdvancesToOutputInput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "urls.txt")
	if err := os.WriteFile(path, []byte("https://example.com/a.txt\nhttps://example.com/b.txt\n"), 0o600); err != nil {
		t.Fatalf("write temp URL file: %v", err)
	}

	model := NewModel(Options{NoColor: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, path)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenOutputInput {
		t.Fatalf("screen = %v, want output input", model.screen)
	}
	if len(model.plan.URLs) != 2 {
		t.Fatalf("plan URL count = %d, want 2", len(model.plan.URLs))
	}
	if !strings.Contains(model.View(), "Enter output directory") {
		t.Fatalf("output input view missing prompt:\n%s", model.View())
	}
}

func TestFilePlanUsesCustomOutputDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "urls.txt")
	if err := os.WriteFile(path, []byte("https://example.com/a.txt\nhttps://example.com/b.txt\n"), 0o600); err != nil {
		t.Fatalf("write temp URL file: %v", err)
	}

	model := NewModel(Options{NoColor: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, path)
	model = updateWithKey(t, model, tea.KeyEnter)
	model.input.SetValue("")
	model = updateWithString(t, model, "/tmp/daryaft-out")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if model.plan.Name != "" {
		t.Fatalf("plan.Name = %q, want empty for batch", model.plan.Name)
	}
	if model.plan.Output != "/tmp/daryaft-out" {
		t.Fatalf("plan.Output = %q, want /tmp/daryaft-out", model.plan.Output)
	}
	if !strings.Contains(model.View(), "Number of URLs: 2") {
		t.Fatalf("plan view missing URL count:\n%s", model.View())
	}
	if !strings.Contains(model.View(), "Output: /tmp/daryaft-out") {
		t.Fatalf("plan view missing output:\n%s", model.View())
	}
	if !strings.Contains(model.View(), "Filename: auto-detect") {
		t.Fatalf("plan view missing auto-detect filename:\n%s", model.View())
	}
}

func TestFileFlowSkipsFilenameInput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "urls.txt")
	if err := os.WriteFile(path, []byte("https://example.com/a.txt\nhttps://example.com/b.txt\n"), 0o600); err != nil {
		t.Fatalf("write temp URL file: %v", err)
	}

	model := NewModel(Options{NoColor: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, path)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if strings.Contains(model.View(), "Enter custom filename") {
		t.Fatalf("batch plan unexpectedly shows filename input prompt:\n%s", model.View())
	}
}

func TestEscAndBackspaceNavigation(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEsc)
	if model.screen != screenHome {
		t.Fatalf("screen after input esc = %v, want home", model.screen)
	}

	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model.input.SetValue("")
	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenURLInput {
		t.Fatalf("screen after output backspace = %v, want URL input", model.screen)
	}

	model.input.SetValue("")
	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenHome {
		t.Fatalf("screen after input backspace = %v, want home", model.screen)
	}

	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model.input.SetValue("")
	model = updateWithString(t, model, "downloads")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenOutputInput {
		t.Fatalf("screen after filename backspace = %v, want output input", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenChecksumInput {
		t.Fatalf("screen after plan backspace = %v, want checksum input", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyEsc)
	if model.screen != screenFilenameInput {
		t.Fatalf("screen after checksum esc = %v, want filename input", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyEsc)
	if model.screen != screenOutputInput {
		t.Fatalf("screen after filename esc = %v, want output input", model.screen)
	}
}

func TestBackspaceEditsNonEmptyInputs(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(Model) Model
		wantScreen screen
		wantValue  string
	}{
		{
			name: "url input",
			setup: func(model Model) Model {
				model.selected = 0
				model = updateWithKey(t, model, tea.KeyEnter)
				return updateWithString(t, model, "abc")
			},
			wantScreen: screenURLInput,
			wantValue:  "ab",
		},
		{
			name: "file input",
			setup: func(model Model) Model {
				model.selected = 1
				model = updateWithKey(t, model, tea.KeyEnter)
				return updateWithString(t, model, "abc")
			},
			wantScreen: screenFileInput,
			wantValue:  "ab",
		},
		{
			name: "output input",
			setup: func(model Model) Model {
				model.selected = 0
				model = updateWithKey(t, model, tea.KeyEnter)
				model = updateWithString(t, model, "https://example.com/file.zip")
				model = updateWithKey(t, model, tea.KeyEnter)
				model.input.SetValue("")
				return updateWithString(t, model, "abc")
			},
			wantScreen: screenOutputInput,
			wantValue:  "ab",
		},
		{
			name: "filename input",
			setup: func(model Model) Model {
				model.selected = 0
				model = updateWithKey(t, model, tea.KeyEnter)
				model = updateWithString(t, model, "https://example.com/file.zip")
				model = updateWithKey(t, model, tea.KeyEnter)
				model = updateWithKey(t, model, tea.KeyEnter)
				return updateWithString(t, model, "abc")
			},
			wantScreen: screenFilenameInput,
			wantValue:  "ab",
		},
		{
			name: "checksum input",
			setup: func(model Model) Model {
				model.selected = 0
				model = updateWithKey(t, model, tea.KeyEnter)
				model = updateWithString(t, model, "https://example.com/file.zip")
				model = updateWithKey(t, model, tea.KeyEnter)
				model = updateWithKey(t, model, tea.KeyEnter)
				model = updateWithKey(t, model, tea.KeyEnter)
				return updateWithString(t, model, "abc")
			},
			wantScreen: screenChecksumInput,
			wantValue:  "ab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setup(NewModel(Options{NoColor: true}))
			model = updateWithKey(t, model, tea.KeyBackspace)

			if model.screen != tt.wantScreen {
				t.Fatalf("screen = %v, want %v", model.screen, tt.wantScreen)
			}
			if model.input.Value() != tt.wantValue {
				t.Fatalf("input value = %q, want %q", model.input.Value(), tt.wantValue)
			}
		})
	}
}

func TestEscapeNavigatesBackWithNonEmptyInput(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "abc")
	model = updateWithKey(t, model, tea.KeyEsc)

	if model.screen != screenHome {
		t.Fatalf("screen after esc = %v, want home", model.screen)
	}
}

func TestThemeMonoUsesNoColorStyles(t *testing.T) {
	model := NewModel(Options{Theme: "mono"})
	if strings.Contains(model.View(), "\x1b[") {
		t.Fatalf("mono theme view contains ANSI color escapes:\n%q", model.View())
	}
}

func TestPlanHomeKey(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithRune(t, model, 'h')

	if model.screen != screenHome {
		t.Fatalf("screen after h = %v, want home", model.screen)
	}
}

func TestQQuits(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	_, cmd := updateWithRuneAndCmd(t, model, 'q')
	if cmd == nil {
		t.Fatal("q command is nil, want quit command")
	}
}

func TestEnterOnPlanStartsExecutionScreen(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenPlan
	model.plan = download.Plan{Output: "downloads"}

	model, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if model.screen != screenExecution {
		t.Fatalf("screen = %v, want execution", model.screen)
	}
	if !model.execution.Running {
		t.Fatal("execution.Running = false, want true")
	}
	if cmd == nil {
		t.Fatal("execution command is nil")
	}
	if model.plan.Output != "downloads" {
		t.Fatalf("plan.Output = %q, want downloads", model.plan.Output)
	}
}

func TestSingleURLExecutionPassesSelectedOutputDirectory(t *testing.T) {
	var received download.Plan
	model := NewModelWithRunner(Options{NoColor: true}, capturePlanRunner(&received))

	model = singleURLPlanModel(t, model, "https://example.com/file.zip", "downloads", "")
	model, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if model.screen != screenExecution {
		t.Fatalf("screen = %v, want execution", model.screen)
	}
	if cmd == nil {
		t.Fatal("execution command is nil")
	}

	_ = cmd()
	if received.Output != "downloads" {
		t.Fatalf("runner plan.Output = %q, want downloads", received.Output)
	}
	if !reflect.DeepEqual(received.URLs, []string{"https://example.com/file.zip"}) {
		t.Fatalf("runner plan.URLs = %#v", received.URLs)
	}
	if received.Retries != tuiDefaultRetries {
		t.Fatalf("runner plan.Retries = %d, want %d", received.Retries, tuiDefaultRetries)
	}
	if !received.Resume {
		t.Fatal("runner plan.Resume = false, want true")
	}
}

func TestSingleURLExecutionPassesSelectedCustomFilename(t *testing.T) {
	var received download.Plan
	model := NewModelWithRunner(Options{NoColor: true}, capturePlanRunner(&received))

	model = singleURLPlanModel(t, model, "https://example.com/file.zip", "downloads", "custom.zip")
	_, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("execution command is nil")
	}

	_ = cmd()
	if received.Name != "custom.zip" {
		t.Fatalf("runner plan.Name = %q, want custom.zip", received.Name)
	}
}

func TestSingleURLExecutionPassesChecksum(t *testing.T) {
	var received download.Plan
	model := NewModelWithRunner(Options{NoColor: true}, capturePlanRunner(&received))
	checksum := "sha256:" + strings.Repeat("a", 64)

	model = singleURLChecksumPlanModel(t, model, "https://example.com/file.zip", "downloads", "custom.zip", checksum)
	_, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("execution command is nil")
	}

	_ = cmd()
	if received.Checksum == nil {
		t.Fatal("runner plan.Checksum = nil")
	}
	if received.Checksum.String() != checksum {
		t.Fatalf("runner plan.Checksum = %q, want %q", received.Checksum.String(), checksum)
	}
}

func TestBatchExecutionPassesSelectedOutputDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "urls.txt")
	if err := os.WriteFile(path, []byte("https://example.com/a.txt\nhttps://example.com/b.txt\n"), 0o600); err != nil {
		t.Fatalf("write temp URL file: %v", err)
	}

	var received download.Plan
	model := NewModelWithRunner(Options{NoColor: true}, capturePlanRunner(&received))
	model = batchPlanModel(t, model, path, "/tmp/daryaft-out")
	_, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("execution command is nil")
	}

	_ = cmd()
	if received.Output != "/tmp/daryaft-out" {
		t.Fatalf("runner plan.Output = %q, want /tmp/daryaft-out", received.Output)
	}
	if !reflect.DeepEqual(received.URLs, []string{"https://example.com/a.txt", "https://example.com/b.txt"}) {
		t.Fatalf("runner plan.URLs = %#v", received.URLs)
	}
	if received.Retries != tuiDefaultRetries {
		t.Fatalf("runner plan.Retries = %d, want %d", received.Retries, tuiDefaultRetries)
	}
	if !received.Resume {
		t.Fatal("runner plan.Resume = false, want true")
	}
}

func TestBatchExecutionDoesNotPassCustomFilename(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "urls.txt")
	if err := os.WriteFile(path, []byte("https://example.com/a.txt\nhttps://example.com/b.txt\n"), 0o600); err != nil {
		t.Fatalf("write temp URL file: %v", err)
	}

	var received download.Plan
	model := NewModelWithRunner(Options{NoColor: true}, capturePlanRunner(&received))
	model = batchPlanModel(t, model, path, "/tmp/daryaft-out")
	_, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("execution command is nil")
	}

	_ = cmd()
	if received.Name != "" {
		t.Fatalf("runner plan.Name = %q, want empty for batch", received.Name)
	}
}

func TestQWhileInjectedRunnerIsRunningCancelsContext(t *testing.T) {
	started := make(chan struct{})
	cancelled := make(chan struct{})
	runner := func(ctx context.Context, plan download.Plan, handlers downloader.BatchHandlers) downloader.BatchResult {
		close(started)
		<-ctx.Done()
		close(cancelled)
		return downloader.BatchResult{
			Planned: len(plan.URLs),
			Items: []downloader.BatchItemResult{{
				Item: downloader.BatchItem{Index: 1, Total: len(plan.URLs), URL: plan.URLs[0]},
				Err:  downloader.ErrCancelled,
			}},
		}
	}

	model := NewModelWithRunner(Options{NoColor: true}, runner)
	model = singleURLPlanModel(t, model, "https://example.com/file.zip", "downloads", "custom.zip")
	model, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("execution command is nil")
	}
	waitForSignal(t, started, "runner start")

	model, quitCmd := updateWithRuneAndCmd(t, model, 'q')
	if quitCmd != nil {
		t.Fatal("q while running returned a command, want nil")
	}
	if model.execution.Status != "Cancelling" {
		t.Fatalf("status = %q, want Cancelling", model.execution.Status)
	}
	waitForSignal(t, cancelled, "runner cancellation")
}

func TestInjectedRunnerForwardsExecutionEventsAndSummary(t *testing.T) {
	item := downloader.BatchItem{Index: 2, Total: 3, URL: "https://example.com/b.zip"}
	runner := func(ctx context.Context, plan download.Plan, handlers downloader.BatchHandlers) downloader.BatchResult {
		handlers.ItemStarted(item)
		handlers.Event(item, downloader.Event{
			Type:                downloader.EventProgress,
			URL:                 item.URL,
			TargetPath:          "downloads/b.zip",
			DownloadedBytes:     128,
			TotalBytes:          512,
			Percent:             25,
			SpeedBytesPerSecond: 2048,
		})
		handlers.Event(item, downloader.Event{
			Type:        downloader.EventRetrying,
			Error:       context.DeadlineExceeded,
			Attempt:     2,
			MaxAttempts: 4,
			NextDelay:   time.Second,
		})
		handlers.Event(item, downloader.Event{
			Type:    downloader.EventCancelled,
			Error:   downloader.ErrCancelled,
			Message: "Download cancelled. Partial file kept for resume.",
		})
		return downloader.BatchResult{
			Planned: 3,
			Items: []downloader.BatchItemResult{
				{Item: downloader.BatchItem{Index: 1, Total: 3, URL: "https://example.com/a.zip"}},
				{Item: item, Err: downloader.ErrCancelled},
			},
		}
	}

	model := NewModelWithRunner(Options{NoColor: true}, runner)
	model.plan = download.Plan{
		URLs:    []string{"https://example.com/a.zip", item.URL, "https://example.com/c.zip"},
		Output:  "downloads",
		Retries: tuiDefaultRetries,
		Resume:  tuiDefaultResume,
	}
	model.screen = screenPlan

	model, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if cmd == nil {
		t.Fatal("execution command is nil")
	}

	model, cmd = updateWithExecutionCmd(t, model, cmd)
	if model.execution.ItemIndex != 2 || model.execution.ItemTotal != 3 {
		t.Fatalf("item = %d/%d, want 2/3", model.execution.ItemIndex, model.execution.ItemTotal)
	}
	if model.execution.CurrentURL != item.URL {
		t.Fatalf("current URL = %q, want %q", model.execution.CurrentURL, item.URL)
	}
	if model.execution.Status != "Starting" {
		t.Fatalf("status after item start = %q, want Starting", model.execution.Status)
	}

	model, cmd = updateWithExecutionCmd(t, model, cmd)
	if model.execution.Status != "Downloading" {
		t.Fatalf("status after progress = %q, want Downloading", model.execution.Status)
	}
	if model.execution.DownloadedBytes != 128 || model.execution.TotalBytes != 512 {
		t.Fatalf("bytes = %d/%d, want 128/512", model.execution.DownloadedBytes, model.execution.TotalBytes)
	}
	if model.execution.Percent != 25 {
		t.Fatalf("percent = %.1f, want 25.0", model.execution.Percent)
	}
	if model.execution.Speed != 2048 {
		t.Fatalf("speed = %.1f, want 2048", model.execution.Speed)
	}
	if model.execution.TargetPath != "downloads/b.zip" {
		t.Fatalf("target path = %q, want downloads/b.zip", model.execution.TargetPath)
	}

	model, cmd = updateWithExecutionCmd(t, model, cmd)
	if model.execution.Status != "Retrying" {
		t.Fatalf("status after retry = %q, want Retrying", model.execution.Status)
	}
	if !strings.Contains(model.execution.Message, "Retrying 2/4 in 1s") {
		t.Fatalf("retry message = %q", model.execution.Message)
	}

	model, cmd = updateWithExecutionCmd(t, model, cmd)
	if model.execution.Status != "Cancelled" {
		t.Fatalf("status after cancelled event = %q, want Cancelled", model.execution.Status)
	}
	if model.execution.Message != "Download cancelled. Partial file kept for resume." {
		t.Fatalf("cancel message = %q", model.execution.Message)
	}

	model, _ = updateWithExecutionCmd(t, model, cmd)
	if model.execution.Running {
		t.Fatal("execution.Running = true, want false after final result")
	}
	if !model.execution.Done {
		t.Fatal("execution.Done = false, want true after final result")
	}
	if model.execution.Summary.Total != 3 {
		t.Fatalf("summary total = %d, want 3", model.execution.Summary.Total)
	}
	if model.execution.Summary.Completed != 1 {
		t.Fatalf("summary completed = %d, want 1", model.execution.Summary.Completed)
	}
	if model.execution.Summary.Cancelled != 1 {
		t.Fatalf("summary cancelled = %d, want 1", model.execution.Summary.Cancelled)
	}
	if model.execution.Summary.Skipped != 1 {
		t.Fatalf("summary skipped = %d, want 1", model.execution.Summary.Skipped)
	}
}

func TestProductionConstructorCreatesUsableModel(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	if model.screen != screenHome {
		t.Fatalf("screen = %v, want home", model.screen)
	}
	if model.executionRunner == nil {
		t.Fatal("executionRunner is nil")
	}
	if model.inspectRunner == nil {
		t.Fatal("inspectRunner is nil")
	}
	if !strings.Contains(model.View(), "Download from URL") {
		t.Fatalf("home view missing menu:\n%s", model.View())
	}
}

func TestExecutionKeepsSelectedCustomFilename(t *testing.T) {
	var received download.Plan
	model := NewModelWithRunner(Options{NoColor: true}, capturePlanRunner(&received))
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "custom.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)

	model, cmd := updateWithKeyAndCmd(t, model, tea.KeyEnter)
	if model.screen != screenExecution {
		t.Fatalf("screen = %v, want execution", model.screen)
	}
	if cmd == nil {
		t.Fatal("execution command is nil")
	}
	_ = cmd()
	if received.Name != "custom.zip" {
		t.Fatalf("runner plan.Name = %q, want custom.zip", received.Name)
	}
}

func TestExecutionScreenRendersProgressFields(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenExecution
	model.execution = executionState{
		Running:         true,
		ItemIndex:       1,
		ItemTotal:       2,
		CurrentURL:      "https://example.com/file.zip",
		TargetPath:      "file.zip",
		Status:          "Downloading",
		DownloadedBytes: 512,
		TotalBytes:      1024,
		Percent:         50,
		Speed:           2048,
		Message:         "Saving",
	}

	view := model.View()
	for _, want := range []string{
		"Downloading",
		"Item 1 of 2",
		"https://example.com/file.zip",
		"Target path: file.zip",
		"Status: Downloading",
		"Downloaded: 512 B / 1.0 KB",
		"Percent: 50.0%",
		"Speed: 2.0 KB/s",
		"Message: Saving",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("execution view missing %q in:\n%s", want, view)
		}
	}
}

func TestDownloaderEventMessageUpdatesExecutionState(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenExecution

	model = updateWithMsg(t, model, executionEventMsg{
		Item: downloader.BatchItem{Index: 2, Total: 3, URL: "https://example.com/file.zip"},
		Event: downloader.Event{
			Type:                downloader.EventProgress,
			URL:                 "https://example.com/file.zip",
			TargetPath:          "file.zip",
			DownloadedBytes:     100,
			TotalBytes:          200,
			Percent:             50,
			SpeedBytesPerSecond: 1000,
		},
	})

	if model.execution.ItemIndex != 2 || model.execution.ItemTotal != 3 {
		t.Fatalf("item = %d/%d, want 2/3", model.execution.ItemIndex, model.execution.ItemTotal)
	}
	if model.execution.Status != "Downloading" {
		t.Fatalf("status = %q, want Downloading", model.execution.Status)
	}
	if model.execution.TargetPath != "file.zip" {
		t.Fatalf("target = %q, want file.zip", model.execution.TargetPath)
	}
	if model.execution.DownloadedBytes != 100 || model.execution.TotalBytes != 200 {
		t.Fatalf("bytes = %d/%d, want 100/200", model.execution.DownloadedBytes, model.execution.TotalBytes)
	}
}

func TestCompletedEventChangesStatus(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenExecution

	model = updateWithMsg(t, model, executionEventMsg{
		Item:  downloader.BatchItem{Index: 1, Total: 1},
		Event: downloader.Event{Type: downloader.EventCompleted},
	})

	if model.execution.Status != "Completed" {
		t.Fatalf("status = %q, want Completed", model.execution.Status)
	}
}

func TestFailedEventChangesStatus(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenExecution

	model = updateWithMsg(t, model, executionEventMsg{
		Item:  downloader.BatchItem{Index: 1, Total: 1},
		Event: downloader.Event{Type: downloader.EventFailed, Error: os.ErrNotExist},
	})

	if model.execution.Status != "Failed" {
		t.Fatalf("status = %q, want Failed", model.execution.Status)
	}
	if model.execution.Message == "" {
		t.Fatal("message is empty, want failure reason")
	}
}

func TestBatchSummaryRendersCounts(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenExecution
	model.execution = executionState{
		Done:   true,
		Status: "Failed",
		Summary: executionSummary{
			Total:     4,
			Completed: 2,
			Failed:    1,
			Cancelled: 1,
			Failures: []executionFailure{{
				URL:   "https://example.com/missing.zip",
				Error: "server returned 404 Not Found",
			}},
		},
	}

	view := model.View()
	for _, want := range []string{
		"Summary",
		"Total: 4",
		"Completed: 2",
		"Failed: 1",
		"Cancelled: 1",
		"https://example.com/missing.zip",
		"server returned",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("summary view missing %q in:\n%s", want, view)
		}
	}
}

func TestQWhileRunningDoesNotQuit(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenExecution
	model.execution = executionState{Running: true}
	cancelled := false
	model.executionCancel = func() {
		cancelled = true
	}

	model, cmd := updateWithRuneAndCmd(t, model, 'q')
	if cmd != nil {
		t.Fatal("q while running returned a command, want nil")
	}
	if !cancelled {
		t.Fatal("cancel func was not called")
	}
	if model.execution.Status != "Cancelling" {
		t.Fatalf("status = %q, want Cancelling", model.execution.Status)
	}
	if model.execution.Message != "Cancelling..." {
		t.Fatalf("message = %q, want Cancelling...", model.execution.Message)
	}
}

func TestCancelledEventChangesStatus(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenExecution

	model = updateWithMsg(t, model, executionEventMsg{
		Item: downloader.BatchItem{Index: 1, Total: 1},
		Event: downloader.Event{
			Type:    downloader.EventCancelled,
			Error:   downloader.ErrCancelled,
			Message: "Download cancelled. Partial file kept for resume.",
		},
	})

	if model.execution.Status != "Cancelled" {
		t.Fatalf("status = %q, want Cancelled", model.execution.Status)
	}
	if model.execution.Message != "Download cancelled. Partial file kept for resume." {
		t.Fatalf("message = %q", model.execution.Message)
	}
}

func TestEnterAfterCompletionReturnsHome(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenExecution
	model.execution = executionState{Done: true, Status: "Completed"}

	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenHome {
		t.Fatalf("screen after enter = %v, want home", model.screen)
	}
	if model.execution.Done {
		t.Fatal("execution state was not cleared")
	}
}

func TestFooterAppearsInView(t *testing.T) {
	model := NewModel(Options{NoColor: true})

	if !strings.Contains(model.View(), config.FooterText) {
		t.Fatalf("view missing footer:\n%s", model.View())
	}
}

func TestVersionScreenIncludesVersionValue(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenVersion

	if !strings.Contains(model.View(), version.Version) {
		t.Fatalf("version view missing %q:\n%s", version.Version, model.View())
	}
}

func TestOptionsConfigInfoPassedToModel(t *testing.T) {
	model := NewModel(Options{
		NoColor:    true,
		ConfigInfo: ConfigInfo{Path: "/tmp/config.yaml", Loaded: true},
	})
	if model.configInfo.Path != "/tmp/config.yaml" {
		t.Fatalf("configInfo.Path = %q, want /tmp/config.yaml", model.configInfo.Path)
	}
	if !model.configInfo.Loaded {
		t.Fatal("configInfo.Loaded = false, want true")
	}
}

func TestSettingsScreenRendersConfigPath(t *testing.T) {
	model := NewModel(Options{
		NoColor:    true,
		ConfigInfo: ConfigInfo{Path: "/tmp/daryaft/config.yaml", Loaded: true},
	})
	model.screen = screenSettings
	view := model.View()
	if !strings.Contains(view, "Config file:") {
		t.Fatalf("settings view missing Config file:\n%s", view)
	}
	if !strings.Contains(view, "/tmp/daryaft/config.yaml") {
		t.Fatalf("settings view missing path:\n%s", view)
	}
}

func TestSettingsScreenRendersConfigLoaded(t *testing.T) {
	model := NewModel(Options{
		NoColor:    true,
		ConfigInfo: ConfigInfo{Path: "/tmp/config.yaml", Loaded: true},
	})
	model.screen = screenSettings
	view := model.View()
	if !strings.Contains(view, "Config loaded: yes") {
		t.Fatalf("settings view missing loaded=yes:\n%s", view)
	}
}

func TestSettingsScreenRendersConfigNotLoaded(t *testing.T) {
	model := NewModel(Options{
		NoColor:    true,
		ConfigInfo: ConfigInfo{Path: "/tmp/config.yaml", Loaded: false},
	})
	model.screen = screenSettings
	view := model.View()
	if !strings.Contains(view, "Config loaded: no (using defaults)") {
		t.Fatalf("settings view missing loaded=no:\n%s", view)
	}
}

func TestSettingsScreenRendersAllSafeFields(t *testing.T) {
	model := NewModel(Options{
		NoColor:           true,
		Theme:             "default",
		DownloadDir:       "/tmp/downloads",
		Retries:           5,
		Resume:            true,
		UseConfigDefaults: true,
		NoTUI:             false,
		Animations:        true,
		Hyperlinks:        true,
		UserAgent:         "DaryaftTest/1.11",
		Timeout:           "30s",
		ConfigInfo:        ConfigInfo{Path: "/tmp/config.yaml", Loaded: true},
	})
	model.screen = screenSettings
	view := model.View()

	for _, want := range []string{
		"download_dir:",
		"retries:",
		"resume:",
		"no_color:",
		"no_tui:",
		"theme:",
		"animations:",
		"hyperlinks:",
		"user_agent:",
		"timeout:",
	} {
		if !strings.Contains(view, want) {
			t.Fatalf("settings view missing %q:\n%s", want, view)
		}
	}

	for _, forbidden := range []string{
		"password",
		"token",
		"authorization",
		"cookie",
		"proxy_authorization",
	} {
		if strings.Contains(strings.ToLower(view), forbidden) {
			t.Fatalf("settings view unexpectedly contains %q:\n%s", forbidden, view)
		}
	}
}

func TestSettingsScreenRendersUserAgentDefault(t *testing.T) {
	model := NewModel(Options{NoColor: true, UserAgent: ""})
	model.screen = screenSettings
	view := model.View()
	if !strings.Contains(view, "user_agent: (default)") {
		t.Fatalf("settings view missing user_agent default marker:\n%s", view)
	}
}

func TestSettingsScreenRendersTimeoutNone(t *testing.T) {
	model := NewModel(Options{NoColor: true, Timeout: ""})
	model.screen = screenSettings
	view := model.View()
	if !strings.Contains(view, "timeout: (none)") {
		t.Fatalf("settings view missing timeout none marker:\n%s", view)
	}
}

func TestSettingsScreenRendersAnimationsReserved(t *testing.T) {
	model := NewModel(Options{NoColor: true, Animations: true, Hyperlinks: true})
	model.screen = screenSettings
	view := model.View()
	if !strings.Contains(view, "animations: true (reserved)") {
		t.Fatalf("settings view missing animations reserved:\n%s", view)
	}
	if !strings.Contains(view, "hyperlinks: true (reserved)") {
		t.Fatalf("settings view missing hyperlinks reserved:\n%s", view)
	}
}

func TestSettingsScreenNavigationBack(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenSettings

	model = updateWithKey(t, model, tea.KeyEsc)
	if model.screen != screenHome {
		t.Fatalf("screen after esc from settings = %v, want home", model.screen)
	}

	model.screen = screenSettings
	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenHome {
		t.Fatalf("screen after backspace from settings = %v, want home", model.screen)
	}
}

func TestCKeyJumpsToSettings(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model = updateWithRune(t, model, 'c')
	if model.screen != screenSettings {
		t.Fatalf("screen after c = %v, want settings", model.screen)
	}
}

func TestSettingsScreenNoColorRendering(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenSettings
	view := model.View()
	if strings.Contains(view, "\x1b[") {
		t.Fatalf("no-color settings view contains ANSI escapes:\n%q", view)
	}
}

func capturePlanRunner(received *download.Plan) ExecutionRunner {
	return func(ctx context.Context, plan download.Plan, handlers downloader.BatchHandlers) downloader.BatchResult {
		*received = plan
		return downloader.BatchResult{Planned: len(plan.URLs)}
	}
}

func singleURLPlanModel(t *testing.T, model Model, rawURL, output, name string) Model {
	t.Helper()
	return singleURLChecksumPlanModel(t, model, rawURL, output, name, "")
}

func singleURLChecksumPlanModel(t *testing.T, model Model, rawURL, output, name, checksum string) Model {
	t.Helper()
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, rawURL)
	model = updateWithKey(t, model, tea.KeyEnter)
	model.input.SetValue("")
	model = updateWithString(t, model, output)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, name)
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, checksum)
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	return model
}

func batchPlanModel(t *testing.T, model Model, path, output string) Model {
	t.Helper()
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, path)
	model = updateWithKey(t, model, tea.KeyEnter)
	model.input.SetValue("")
	model = updateWithString(t, model, output)
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	return model
}

func waitForSignal(t *testing.T, signal <-chan struct{}, label string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for %s", label)
	}
}

func updateWithInspectCmd(t *testing.T, model Model, cmd tea.Cmd) (Model, tea.Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("inspect command is nil")
	}
	msg := cmd()
	return updateWithMsgAndCmd(t, model, msg)
}

func updateWithExecutionCmd(t *testing.T, model Model, cmd tea.Cmd) (Model, tea.Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("execution command is nil")
	}
	msg := cmd()
	return updateWithMsgAndCmd(t, model, msg)
}

func sampleInspectResult(rawURL string) inspect.Result {
	return inspect.Result{
		URL:                rawURL,
		FinalURL:           rawURL,
		Status:             "200 OK",
		StatusCode:         200,
		Filename:           "file.zip",
		ContentLength:      1048576,
		ContentLengthKnown: true,
		ContentType:        "application/zip",
		AcceptRanges:       "bytes",
		ResumeSupported:    true,
		ResumeSupportKnown: true,
		ETag:               `"abc123"`,
		LastModified:       "Tue, 01 Jun 2026 12:00:00 GMT",
	}
}

func TestEmptyURLShowsGuidanceError(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenURLInput {
		t.Fatalf("screen = %v, want URL input", model.screen)
	}
	if model.errorMessage == "" {
		t.Fatal("errorMessage is empty, want guidance error for empty URL")
	}
	if !strings.Contains(model.errorMessage, "https://") {
		t.Fatalf("errorMessage %q should mention https://", model.errorMessage)
	}
}

func TestEmptyFilePathShowsGuidanceError(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenFileInput {
		t.Fatalf("screen = %v, want file input", model.screen)
	}
	if model.errorMessage == "" {
		t.Fatal("errorMessage is empty, want guidance error for empty file path")
	}
	if !strings.Contains(model.errorMessage, ".txt") {
		t.Fatalf("errorMessage %q should mention .txt", model.errorMessage)
	}
}

func TestInvalidURLSchemeShowsError(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "ftp://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenURLInput {
		t.Fatalf("screen = %v, want URL input", model.screen)
	}
	if model.errorMessage == "" {
		t.Fatal("errorMessage is empty, want scheme validation error")
	}
}

func TestURLInputViewShowsDefaultsPreview(t *testing.T) {
	model := NewModel(Options{NoColor: true, DownloadDir: "/tmp/downloads", Retries: 5, Resume: false, UseConfigDefaults: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)

	view := model.View()
	if !strings.Contains(view, "Defaults:") {
		t.Fatalf("URL input view missing defaults preview:\n%s", view)
	}
	if !strings.Contains(view, "retries 5") {
		t.Fatalf("URL input view defaults preview missing retries:\n%s", view)
	}
	if !strings.Contains(view, "resume false") {
		t.Fatalf("URL input view defaults preview missing resume:\n%s", view)
	}
}

func TestFileInputViewShowsDefaultsPreview(t *testing.T) {
	model := NewModel(Options{NoColor: true, Retries: 2, Resume: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)

	view := model.View()
	if !strings.Contains(view, "Defaults:") {
		t.Fatalf("file input view missing defaults preview:\n%s", view)
	}
}

func TestHelpViewMentionsSettings(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.screen = screenHelp
	view := model.View()
	if !strings.Contains(view, "Settings") {
		t.Fatalf("help view does not mention Settings:\n%s", view)
	}
}

func TestURLInputPromptMentionsScheme(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	view := model.View()
	if !strings.Contains(view, "https://") {
		t.Fatalf("URL input prompt does not mention https://:\n%s", view)
	}
}

func TestFileInputPromptMentionsAbsolutePath(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	view := model.View()
	if !strings.Contains(view, "absolute") {
		t.Fatalf("file input prompt does not mention absolute path:\n%s", view)
	}
}

func TestErrorMessageClearsOnInputChange(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	// trigger an error
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.errorMessage == "" {
		t.Fatal("expected error message after empty submit")
	}
	// typing a character should clear the error
	model = updateWithString(t, model, "h")
	if model.errorMessage != "" {
		t.Fatalf("errorMessage = %q, want empty after typing", model.errorMessage)
	}
}

func updateWithKey(t *testing.T, model Model, key tea.KeyType) Model {
	t.Helper()
	next, _ := updateWithKeyAndCmd(t, model, key)
	return next
}

func updateWithKeyAndCmd(t *testing.T, model Model, key tea.KeyType) (Model, tea.Cmd) {
	t.Helper()
	updated, cmd := model.Update(tea.KeyMsg{Type: key})
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want tui.Model", updated)
	}
	return next, cmd
}

func updateWithRune(t *testing.T, model Model, key rune) Model {
	t.Helper()
	next, _ := updateWithRuneAndCmd(t, model, key)
	return next
}

func updateWithRuneAndCmd(t *testing.T, model Model, key rune) (Model, tea.Cmd) {
	t.Helper()
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}})
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want tui.Model", updated)
	}
	return next, cmd
}

func updateWithString(t *testing.T, model Model, value string) Model {
	t.Helper()
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)})
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want tui.Model", updated)
	}
	return next
}

func updateWithMsg(t *testing.T, model Model, msg tea.Msg) Model {
	t.Helper()
	next, _ := updateWithMsgAndCmd(t, model, msg)
	return next
}

func updateWithMsgAndCmd(t *testing.T, model Model, msg tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	updated, cmd := model.Update(msg)
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want tui.Model", updated)
	}
	return next, cmd
}
