package download

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/checksum"
	"github.com/he8um/daryaft/internal/httpopts"
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
		"Mode: dry-run only, no network request performed",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("DryRunString() missing %q in:\n%s", want, got)
		}
	}
}

func TestDryRunStringIncludesChecksum(t *testing.T) {
	sum := sha256.Sum256([]byte("hello"))
	spec := checksum.Spec{
		Algorithm: checksum.AlgorithmSHA256,
		Expected:  fmt.Sprintf("%x", sum),
	}
	plan := Plan{
		URLs:     []string{"https://example.com/file.zip"},
		Checksum: &spec,
		Retries:  3,
		Resume:   true,
	}

	got := plan.DryRunString()
	if !strings.Contains(got, "Checksum: sha256:"+spec.Expected) {
		t.Fatalf("DryRunString() missing checksum in:\n%s", got)
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

func TestBuildPlanPassesHTTPOptions(t *testing.T) {
	headers, _ := httpopts.ParseHeaders([]string{"X-Custom: value"})
	opts := httpopts.Options{
		ProxyURL:  "http://proxy:8080",
		Headers:   headers,
		UserAgent: "TestAgent/1.0",
		Username:  "alice",
		Password:  "pass",
	}
	plan, err := BuildPlan(Options{
		URLs:        []string{"https://example.com/file.zip"},
		Retries:     3,
		Resume:      true,
		HTTPOptions: opts,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	if plan.HTTPOptions.ProxyURL != "http://proxy:8080" {
		t.Errorf("ProxyURL not passed through: %q", plan.HTTPOptions.ProxyURL)
	}
	if plan.HTTPOptions.UserAgent != "TestAgent/1.0" {
		t.Errorf("UserAgent not passed through: %q", plan.HTTPOptions.UserAgent)
	}
}

func TestBuildPlanRejectsInvalidHTTPOptions(t *testing.T) {
	tests := []struct {
		name string
		opts httpopts.Options
		want string
	}{
		{
			name: "invalid header",
			opts: httpopts.Options{Headers: []httpopts.Header{{Name: "Bad Header", Value: "v"}}},
			want: "invalid character",
		},
		{
			name: "invalid proxy",
			opts: httpopts.Options{ProxyURL: "socks5://proxy:1080"},
			want: "unsupported scheme",
		},
		{
			name: "password without username",
			opts: httpopts.Options{Password: "secret"},
			want: "--password requires --username",
		},
		{
			name: "auth header plus basic auth",
			opts: httpopts.Options{
				Username: "alice",
				Password: "pass",
				Headers:  []httpopts.Header{{Name: "Authorization", Value: "Bearer tok"}},
			},
			want: "Authorization header",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildPlan(Options{
				URLs:        []string{"https://example.com/file.zip"},
				Retries:     3,
				Resume:      true,
				HTTPOptions: tt.opts,
			})
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

func TestDryRunStringShowsHTTPOptions(t *testing.T) {
	headers, _ := httpopts.ParseHeaders([]string{"X-Custom: shown", "Authorization: Bearer secret"})
	plan := Plan{
		URLs:    []string{"https://example.com/file.zip"},
		Retries: 3,
		Resume:  true,
		HTTPOptions: httpopts.Options{
			ProxyURL:  "http://proxy:8080",
			Headers:   headers,
			UserAgent: "TestAgent/1.0",
			Username:  "alice",
			Password:  "topsecret",
		},
	}

	got := plan.DryRunString()

	for _, want := range []string{"http://proxy:8080", "X-Custom: shown", "TestAgent/1.0", "[REDACTED]"} {
		if !strings.Contains(got, want) {
			t.Errorf("DryRunString() missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "topsecret") {
		t.Errorf("DryRunString() must not contain raw password:\n%s", got)
	}
	if strings.Contains(got, "Bearer secret") {
		t.Errorf("DryRunString() must not contain raw Authorization value:\n%s", got)
	}
}

func TestDryRunStringAuthFormat(t *testing.T) {
	plan := Plan{
		URLs:    []string{"https://example.com/file.zip"},
		Retries: 3,
		Resume:  true,
		HTTPOptions: httpopts.Options{
			Username: "alice",
			Password: "secret",
		},
	}
	got := plan.DryRunString()
	if !strings.Contains(got, "Auth: alice:[REDACTED]") {
		t.Errorf("DryRunString() auth line should show username:[REDACTED], got:\n%s", got)
	}
	if strings.Contains(got, "secret") {
		t.Errorf("DryRunString() must not contain raw password:\n%s", got)
	}
}

func TestDryRunStringAuthUsernameOnlyFormat(t *testing.T) {
	plan := Plan{
		URLs:    []string{"https://example.com/file.zip"},
		Retries: 3,
		Resume:  true,
		HTTPOptions: httpopts.Options{
			Username: "alice",
		},
	}
	got := plan.DryRunString()
	if !strings.Contains(got, "Auth: alice") {
		t.Errorf("DryRunString() auth line should show username when no password, got:\n%s", got)
	}
}

func TestDryRunStringNoHTTPSectionWhenEmpty(t *testing.T) {
	plan := Plan{
		URLs:    []string{"https://example.com/file.zip"},
		Retries: 3,
		Resume:  true,
	}
	got := plan.DryRunString()
	if strings.Contains(got, "HTTP") {
		t.Errorf("DryRunString() should not include HTTP section when no options set:\n%s", got)
	}
}
