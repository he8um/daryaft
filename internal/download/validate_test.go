package download

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestBuildPlanValidatesURLSchemes(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{name: "ftp", url: "ftp://example.com/file.zip"},
		{name: "missing scheme", url: "example.com/file.zip"},
		{name: "missing host", url: "https:///file.zip"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildPlan(Options{
				URLs:    []string{tt.url},
				Retries: 3,
				Resume:  true,
			})
			if err == nil {
				t.Fatal("BuildPlan returned nil error")
			}
			if !strings.Contains(err.Error(), "invalid URL") {
				t.Fatalf("error = %q, want invalid URL context", err)
			}
		})
	}
}

func TestBuildPlanCombinesArgsAndFileURLs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "urls.txt")
	if err := os.WriteFile(path, []byte("https://example.com/two.zip\n"), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	plan, err := BuildPlan(Options{
		URLs:    []string{"https://example.com/one.zip"},
		File:    path,
		Retries: 3,
		Resume:  true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}

	if len(plan.URLs) != 2 {
		t.Fatalf("len(plan.URLs) = %d, want 2", len(plan.URLs))
	}
}

func TestBuildPlanRejectsNameWithMultipleURLs(t *testing.T) {
	_, err := BuildPlan(Options{
		URLs: []string{
			"https://example.com/one.zip",
			"https://example.com/two.zip",
		},
		Name:    "file.zip",
		Retries: 3,
		Resume:  true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "--name") {
		t.Fatalf("error = %q, want --name context", err)
	}
}

func TestBuildPlanRejectsNegativeRetries(t *testing.T) {
	_, err := BuildPlan(Options{
		URLs:    []string{"https://example.com/file.zip"},
		Retries: -1,
		Resume:  true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "retries") {
		t.Fatalf("error = %q, want retries context", err)
	}
}

func TestBuildPlanAcceptsRetryBounds(t *testing.T) {
	for _, retries := range []int{0, MaxRetries} {
		t.Run(strconv.Itoa(retries), func(t *testing.T) {
			plan, err := BuildPlan(Options{
				URLs:    []string{"https://example.com/file.zip"},
				Retries: retries,
				Resume:  true,
			})
			if err != nil {
				t.Fatalf("BuildPlan returned error: %v", err)
			}
			if plan.Retries != retries {
				t.Fatalf("Retries = %d, want %d", plan.Retries, retries)
			}
		})
	}
}

func TestBuildPlanRejectsRetriesAboveMax(t *testing.T) {
	_, err := BuildPlan(Options{
		URLs:    []string{"https://example.com/file.zip"},
		Retries: MaxRetries + 1,
		Resume:  true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "retries") {
		t.Fatalf("error = %q, want retries context", err)
	}
}

func TestBuildPlanRequiresURLOrFile(t *testing.T) {
	_, err := BuildPlan(Options{
		Retries: 3,
		Resume:  true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "provide at least one URL") {
		t.Fatalf("error = %q, want missing input context", err)
	}
}

func TestBuildPlanNoResumeOverridesResume(t *testing.T) {
	plan, err := BuildPlan(Options{
		URLs:     []string{"https://example.com/file.zip"},
		Retries:  3,
		Resume:   true,
		NoResume: true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	if plan.Resume {
		t.Fatal("plan.Resume = true, want false")
	}
}
