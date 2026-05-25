package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/config"
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

func TestValidURLCreatesPlanScreen(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
	model = updateWithKey(t, model, tea.KeyEnter)

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if len(model.plan.URLs) != 1 {
		t.Fatalf("plan URL count = %d, want 1", len(model.plan.URLs))
	}
	if !strings.Contains(model.View(), "Number of URLs: 1") {
		t.Fatalf("plan view missing URL count:\n%s", model.View())
	}
	if !strings.Contains(model.View(), "TUI download execution is planned") {
		t.Fatalf("plan view missing execution boundary:\n%s", model.View())
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

func TestValidFileCreatesPlanScreen(t *testing.T) {
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

	if model.screen != screenPlan {
		t.Fatalf("screen = %v, want plan", model.screen)
	}
	if len(model.plan.URLs) != 2 {
		t.Fatalf("plan URL count = %d, want 2", len(model.plan.URLs))
	}
	if !strings.Contains(model.View(), "Number of URLs: 2") {
		t.Fatalf("plan view missing URL count:\n%s", model.View())
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
		t.Fatalf("screen after plan backspace = %v, want URL input", model.screen)
	}

	model = updateWithKey(t, model, tea.KeyBackspace)
	if model.screen != screenHome {
		t.Fatalf("screen after input backspace = %v, want home", model.screen)
	}
}

func TestPlanHomeKey(t *testing.T) {
	model := NewModel(Options{NoColor: true})
	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	model = updateWithString(t, model, "https://example.com/file.zip")
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
