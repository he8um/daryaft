package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/downloader"
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
	model.selected = 2

	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenHelp {
		t.Fatalf("screen = %v, want help", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyEsc)
	if model.screen != screenHome {
		t.Fatalf("screen after esc = %v, want home", model.screen)
	}

	model.selected = 3
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenVersion {
		t.Fatalf("screen = %v, want version", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenHome {
		t.Fatalf("screen after backspace = %v, want home", model.screen)
	}
}

func TestSelectURLInputScreen(t *testing.T) {
	model := NewModel(Options{NoColor: true})

	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenURLInput {
		t.Fatalf("screen = %v, want URL input", model.screen)
	}
	if !strings.Contains(model.View(), "Enter download URL") {
		t.Fatalf("URL input view missing prompt:\n%s", model.View())
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
	if !strings.Contains(model.View(), "Default/current value: .") {
		t.Fatalf("output input view missing default:\n%s", model.View())
	}
}

func TestEmptyOutputDirectoryCreatesPlanWithCurrentDirectory(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if model.plan.Output != "." {
		t.Fatalf("plan.Output = %q, want .", model.plan.Output)
	}
	if !strings.Contains(model.View(), "Number of URLs: 1") {
		t.Fatalf("plan view missing URL count:\n%s", model.View())
	}
	if !strings.Contains(model.View(), "Output: .") {
		t.Fatalf("plan view missing current directory output:\n%s", model.View())
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
	model = updateWithString(t, model, "downloads")
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

func TestSelectFileInputScreen(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenFileInput {
		t.Fatalf("screen = %v, want file input", model.screen)
	}
	if !strings.Contains(model.View(), "Enter path to .txt file") {
		t.Fatalf("file input view missing prompt:\n%s", model.View())
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
	model = updateWithString(t, model, "/tmp/daryaft-out")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
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
	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenURLInput {
		t.Fatalf("screen after output backspace = %v, want URL input", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenHome {
		t.Fatalf("screen after input backspace = %v, want home", model.screen)
	}

	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "downloads")
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenOutputInput {
		t.Fatalf("screen after plan backspace = %v, want output input", model.screen)
	}
}

func TestPlanHomeKey(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
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
	updated, _ := model.Update(msg)
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want tui.Model", updated)
	}
	return next
}
