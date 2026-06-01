package version

import (
	"runtime"
	"testing"
)

func TestInfoDefaults(t *testing.T) {
	info := Info()

	if info.Version != "0.5.0-dev" {
		t.Fatalf("Version = %q, want 0.5.0-dev", info.Version)
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
	var _ string = Version
	var _ string = Commit
	var _ string = Date
	var _ string = BuiltBy
}
