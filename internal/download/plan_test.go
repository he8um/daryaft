package download

import (
	"strings"
	"testing"
)

func TestDryRunString(t *testing.T) {
	plan := Plan{
		URLs:    []string{"https://example.com/file.zip"},
		Retries: 3,
		Resume:  true,
	}

	got := plan.DryRunString()
	for _, want := range []string{
		"Daryaft download plan",
		"URLs: 1",
		"1. https://example.com/file.zip",
		"Output: current directory",
		"Filename: auto-detect",
		"Retries: 3",
		"Resume: true",
		"Mode: dry-run only, downloader engine not implemented yet",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("DryRunString() missing %q in:\n%s", want, got)
		}
	}
}

func TestBuildPlanCreatesDryRunPlan(t *testing.T) {
	plan, err := BuildPlan(Options{
		URLs:    []string{" https://example.com/file.zip "},
		Output:  "downloads",
		Name:    "file.zip",
		Retries: 3,
		Resume:  true,
		DryRun:  true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}

	if plan.URLs[0] != "https://example.com/file.zip" {
		t.Fatalf("plan.URLs[0] = %q", plan.URLs[0])
	}
	if plan.Output != "downloads" {
		t.Fatalf("plan.Output = %q", plan.Output)
	}
	if plan.Name != "file.zip" {
		t.Fatalf("plan.Name = %q", plan.Name)
	}
}
