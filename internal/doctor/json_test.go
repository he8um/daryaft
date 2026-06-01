package doctor

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	appconfig "github.com/he8um/daryaft/internal/config"
)

func TestFormatJSONHealthyReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "daryaft", "config.yaml")

	data, err := FormatJSON(Run(testOptions(dir, path, appconfig.Default())))
	if err != nil {
		t.Fatalf("FormatJSON returned error: %v", err)
	}

	var got JSONReport
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v\n%s", err, string(data))
	}
	if !got.OK {
		t.Fatalf("OK = false, want true:\n%s", string(data))
	}
	if got.Summary.Failures != 0 {
		t.Fatalf("Failures = %d, want 0", got.Summary.Failures)
	}
	if got.Summary.Checks == 0 {
		t.Fatal("Checks = 0, want checks")
	}
	assertJSONStatusPresent(t, got, "ok")
	assertJSONStatusPresent(t, got, "info")
	assertJSONStatusPresent(t, got, "skipped")
}

func TestFormatJSONFailureStatus(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "daryaft", "config.yaml")
	options := testOptions(dir, path, appconfig.Default())
	options.LoadConfig = func() (appconfig.Config, error) {
		return appconfig.Config{}, errors.New("parse config: invalid yaml")
	}

	got := ToJSONReport(Run(options))

	if got.OK {
		t.Fatal("OK = true, want false")
	}
	if got.Summary.Failures != 1 {
		t.Fatalf("Failures = %d, want 1", got.Summary.Failures)
	}
	assertJSONStatusPresent(t, got, "failure")
}

func TestJSONSummaryCountsMatchSections(t *testing.T) {
	report := Report{}
	report.Add("System", StatusOK, "OS", "darwin")
	report.Add("Download", StatusWarn, "Output writable", "directory missing")
	report.Add("Config", StatusFail, "Config load", "parse config")
	report.Add("Skipped", StatusSkipped, "GitHub release check", "skipped")

	got := ToJSONReport(report)

	checks := 0
	for _, section := range got.Sections {
		checks += len(section.Checks)
	}
	if got.Summary.Checks != checks {
		t.Fatalf("summary checks = %d, section checks = %d", got.Summary.Checks, checks)
	}
	if got.Summary.Checks != 4 {
		t.Fatalf("summary checks = %d, want 4", got.Summary.Checks)
	}
	if got.Summary.Failures != 1 {
		t.Fatalf("failures = %d, want 1", got.Summary.Failures)
	}
	if got.Summary.Warnings != 1 {
		t.Fatalf("warnings = %d, want 1", got.Summary.Warnings)
	}
}

func TestJSONStrictWarningSetsOKFalseWithoutChangingWarningStatus(t *testing.T) {
	report := Report{}
	report.Add("Download", StatusWarn, "Output writable", "directory missing")

	got := ToJSONReportWithOptions(report, true)

	if got.OK {
		t.Fatal("OK = true, want false in strict mode with warning")
	}
	if !got.Strict {
		t.Fatal("Strict = false, want true")
	}
	if got.Summary.Failures != 0 {
		t.Fatalf("failures = %d, want 0", got.Summary.Failures)
	}
	if got.Summary.Warnings != 1 {
		t.Fatalf("warnings = %d, want 1", got.Summary.Warnings)
	}
	if got.Sections[0].Checks[0].Status != "warning" {
		t.Fatalf("status = %q, want warning", got.Sections[0].Checks[0].Status)
	}
}

func TestReportOKStrictPolicy(t *testing.T) {
	healthy := Report{}
	healthy.Add("System", StatusOK, "OS", "darwin")
	if !healthy.OK(true) {
		t.Fatal("healthy report OK(true) = false")
	}

	warning := Report{}
	warning.Add("Download", StatusWarn, "Output writable", "directory missing")
	if !warning.OK(false) {
		t.Fatal("warning report OK(false) = false")
	}
	if warning.OK(true) {
		t.Fatal("warning report OK(true) = true")
	}

	failure := Report{}
	failure.Add("Config", StatusFail, "Config load", "parse config")
	if failure.OK(false) {
		t.Fatal("failure report OK(false) = true")
	}
	if failure.OK(true) {
		t.Fatal("failure report OK(true) = true")
	}
}

func assertJSONStatusPresent(t *testing.T, report JSONReport, status string) {
	t.Helper()

	for _, section := range report.Sections {
		for _, check := range section.Checks {
			if check.Status == status {
				return
			}
		}
	}
	t.Fatalf("status %q not found in %#v", status, report)
}
