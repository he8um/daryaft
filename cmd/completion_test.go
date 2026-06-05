package cmd

import (
	"bytes"
	"strings"
	"testing"

	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/spf13/cobra"
)

func TestCompletionCommandExists(t *testing.T) {
	command, _, err := rootCmd.Find([]string{"completion"})
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if command == nil || command.Name() != "completion" {
		t.Fatalf("completion command not found")
	}
}

func TestCompletionCommandUnsupportedShellReturnsError(t *testing.T) {
	_, err := executeCompletionCommand(t, "xonsh")
	if err == nil {
		t.Fatal("completion returned nil error")
	}
	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Fatalf("error = %q", err)
	}
}

func TestCompletionGenerationDoesNotError(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		t.Run(shell, func(t *testing.T) {
			output, err := executeCompletionCommand(t, shell)
			if err != nil {
				t.Fatalf("completion %s returned error: %v", shell, err)
			}
			if !strings.Contains(output, "daryaft") {
				t.Fatalf("completion %s output does not mention daryaft", shell)
			}
		})
	}
}

func TestConfigGetKeyCompletionIncludesAllSupportedKeys(t *testing.T) {
	command := newConfigGetCommand()

	got, directive := command.ValidArgsFunction(command, nil, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want NoFileComp", directive)
	}
	assertSameStrings(t, got, supportedConfigKeysForTest())
}

func TestConfigSetKeyCompletionIncludesAllSupportedKeys(t *testing.T) {
	command := newConfigSetCommand()

	got, directive := command.ValidArgsFunction(command, nil, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want NoFileComp", directive)
	}
	assertSameStrings(t, got, supportedConfigKeysForTest())
}

func TestConfigSetBoolValueCompletion(t *testing.T) {
	command := newConfigSetCommand()

	got, directive := command.ValidArgsFunction(command, []string{"resume"}, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want NoFileComp", directive)
	}
	assertSameStrings(t, got, []string{"true", "false"})
}

func TestChecksumFlagCompletion(t *testing.T) {
	got, directive := checksumFlagCompletions(&cobra.Command{}, nil, "")
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want NoFileComp", directive)
	}
	assertSameStrings(t, got, []string{"sha256:", "sha512:"})
}

func executeCompletionCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var output bytes.Buffer
	root := &cobra.Command{
		Use:           "daryaft",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newCompletionCommand())
	root.SetOut(&output)
	root.SetErr(&output)
	root.SetArgs(append([]string{"completion"}, args...))
	err := root.Execute()
	return output.String(), err
}

func supportedConfigKeysForTest() []string {
	keys := appconfig.SupportedKeys()
	names := make([]string, 0, len(keys))
	for _, key := range keys {
		names = append(names, key.Key)
	}
	return names
}

func assertSameStrings(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d; got %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("item %d = %q, want %q; got %#v", i, got[i], want[i], got)
		}
	}
}
