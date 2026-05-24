package tui

import (
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

func TestPlannedScreens(t *testing.T) {
	model := NewModel(Options{NoColor: true})

	model.selected = 0
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenURLPlanned {
		t.Fatalf("screen = %v, want URL planned", model.screen)
	}
	if !strings.Contains(model.View(), "Planned for the download TUI milestone") {
		t.Fatalf("planned URL view missing planned message:\n%s", model.View())
	}

	model = model.back()
	model.selected = 1
	model = updateWithKey(t, model, tea.KeyEnter)
	if model.screen != screenFilePlanned {
		t.Fatalf("screen = %v, want file planned", model.screen)
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
	updated, _ := model.Update(tea.KeyMsg{Type: key})
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want tui.Model", updated)
	}
	return next
}

func updateWithRune(t *testing.T, model Model, key rune) Model {
	t.Helper()
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}})
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want tui.Model", updated)
	}
	return next
}
