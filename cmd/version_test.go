package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/he8um/daryaft/pkg/version"
)

func TestVersionCommandTextOutput(t *testing.T) {
	output, err := executeVersionCommand(t)
	if err != nil {
		t.Fatalf("version returned error: %v\n%s", err, output)
	}

	for _, want := range []string{
		"Daryaft version: " + version.Version,
		"commit: " + version.Commit,
		"build date: " + version.Date,
		"built by: " + version.BuiltBy,
		"go version: ",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestVersionCommandJSONOutput(t *testing.T) {
	output, err := executeVersionCommand(t, "--json")
	if err != nil {
		t.Fatalf("version --json returned error: %v\n%s", err, output)
	}

	var got version.Details
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v\n%s", err, output)
	}
	if got.Version != version.Version {
		t.Fatalf("Version = %q, want %q", got.Version, version.Version)
	}
	if got.Commit != version.Commit {
		t.Fatalf("Commit = %q, want %q", got.Commit, version.Commit)
	}
	if got.Date != version.Date {
		t.Fatalf("Date = %q, want %q", got.Date, version.Date)
	}
	if got.BuiltBy != version.BuiltBy {
		t.Fatalf("BuiltBy = %q, want %q", got.BuiltBy, version.BuiltBy)
	}
	if got.GoVersion == "" {
		t.Fatal("GoVersion is empty")
	}
	if strings.Contains(output, "Daryaft version:") {
		t.Fatalf("JSON output contains human text:\n%s", output)
	}
}

func executeVersionCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var output bytes.Buffer
	var stderr bytes.Buffer
	command := newVersionCommand()
	command.SetOut(&output)
	command.SetErr(&stderr)
	command.SetArgs(args)
	err := command.Execute()
	return output.String(), err
}
