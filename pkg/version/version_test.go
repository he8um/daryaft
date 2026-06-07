package version

import (
	"runtime"
	"testing"
)

func TestInfoDefaults(t *testing.T) {
	info := Info()

	if info.Version != "1.3.0-dev" {
		t.Fatalf("Version = %q, want 1.3.0-dev", info.Version)
	}
	if info.Commit != "local" {
		t.Fatalf("Commit = %q, want local", info.Commit)
	}
	if info.Date != "unknown" {
		t.Fatalf("Date = %q, want unknown", info.Date)
	}
	if info.BuiltBy != "source" {
		t.Fatalf("BuiltBy = %q, want source", info.BuiltBy)
	}
	if info.GoVersion != runtime.Version() {
		t.Fatalf("GoVersion = %q, want %q", info.GoVersion, runtime.Version())
	}
}

func TestLdflagVariablesAreExportedStrings(t *testing.T) {
	values := []string{Version, Commit, Date, BuiltBy}
	if len(values) != 4 {
		t.Fatalf("ldflag variable count = %d, want 4", len(values))
	}
}
