package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/doctor"
)

func TestDoctorCommandExists(t *testing.T) {
	command, _, err := rootCmd.Find([]string{"doctor"})
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if command == nil || command.Name() != "doctor" {
		t.Fatal("doctor command not found")
	}
}

func TestDoctorJSONOutputsValidJSON(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	output, err := executeDoctorCommand(t, "--json")
	if err != nil {
		t.Fatalf("doctor --json returned error: %v\n%s", err, output)
	}

	var got doctor.JSONReport
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v\n%s", err, output)
	}
	if !got.OK {
		t.Fatalf("OK = false, want true:\n%s", output)
	}
	if got.Summary.Failures != 0 {
		t.Fatalf("Failures = %d, want 0", got.Summary.Failures)
	}
	if strings.Contains(output, "Daryaft doctor") {
		t.Fatalf("JSON output contains human heading:\n%s", output)
	}
}

func TestDoctorHumanOutputIsNotJSON(t *testing.T) {
	restore := appconfig.SetUserConfigDirForTest(t.TempDir())
	t.Cleanup(restore)

	output, err := executeDoctorCommand(t)
	if err != nil {
		t.Fatalf("doctor returned error: %v\n%s", err, output)
	}
	if !strings.HasPrefix(output, "Daryaft doctor\n") {
		t.Fatalf("human output missing heading:\n%s", output)
	}
	var got doctor.JSONReport
	if err := json.Unmarshal([]byte(output), &got); err == nil {
		t.Fatalf("human output unexpectedly parsed as JSON: %#v", got)
	}
}

func TestDoctorJSONCriticalFailureReturnsError(t *testing.T) {
	dir := t.TempDir()
	restore := appconfig.SetUserConfigDirForTest(dir)
	t.Cleanup(restore)

	path, err := appconfig.Path()
	if err != nil {
		t.Fatalf("Path returned error: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("download_dir: [\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	output, err := executeDoctorCommand(t, "--json")
	if err == nil {
		t.Fatalf("doctor --json returned nil error:\n%s", output)
	}

	var got doctor.JSONReport
	if unmarshalErr := json.Unmarshal([]byte(output), &got); unmarshalErr != nil {
		t.Fatalf("json.Unmarshal returned error: %v\n%s", unmarshalErr, output)
	}
	if got.OK {
		t.Fatalf("OK = true, want false:\n%s", output)
	}
	if got.Summary.Failures == 0 {
		t.Fatalf("Failures = 0, want at least one:\n%s", output)
	}
	if !doctorJSONStatusPresent(got, "failure") {
		t.Fatalf("failure status missing:\n%s", output)
	}
}

func executeDoctorCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var output bytes.Buffer
	var stderr bytes.Buffer
	command := newDoctorCommand()
	command.SetOut(&output)
	command.SetErr(&stderr)
	command.SetArgs(args)
	err := command.Execute()
	return output.String(), err
}

func doctorJSONStatusPresent(report doctor.JSONReport, status string) bool {
	for _, section := range report.Sections {
		for _, check := range section.Checks {
			if check.Status == status {
				return true
			}
		}
	}
	return false
}
