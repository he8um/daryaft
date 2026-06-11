package download

import (
	"crypto/sha256"
	"fmt"
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

func TestBuildPlanParsesChecksum(t *testing.T) {
	sum := sha256.Sum256([]byte("hello"))
	plan, err := BuildPlan(Options{
		URLs:     []string{"https://example.com/file.zip"},
		Checksum: " sha256:" + fmt.Sprintf("%X", sum) + " ",
		Retries:  3,
		Resume:   true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	if plan.Checksum == nil {
		t.Fatal("plan.Checksum = nil")
	}
	if plan.Checksum.Algorithm != "sha256" {
		t.Fatalf("Algorithm = %q", plan.Checksum.Algorithm)
	}
	if plan.Checksum.Expected != fmt.Sprintf("%x", sum) {
		t.Fatalf("Expected = %q", plan.Checksum.Expected)
	}
}

func TestBuildPlanRejectsChecksumWithMultipleURLs(t *testing.T) {
	sum := sha256.Sum256([]byte("hello"))
	_, err := BuildPlan(Options{
		URLs: []string{
			"https://example.com/one.zip",
			"https://example.com/two.zip",
		},
		Checksum: "sha256:" + fmt.Sprintf("%x", sum),
		Retries:  3,
		Resume:   true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "--checksum is currently supported only for single URL downloads") {
		t.Fatalf("error = %q", err)
	}
}

func TestBuildPlanRejectsChecksumWithFileInput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "urls.txt")
	if err := os.WriteFile(path, []byte("https://example.com/file.zip\n"), 0o600); err != nil {
		t.Fatalf("write URL file: %v", err)
	}
	sum := sha256.Sum256([]byte("hello"))

	_, err := BuildPlan(Options{
		File:     path,
		Checksum: "sha256:" + fmt.Sprintf("%x", sum),
		Retries:  3,
		Resume:   true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "--checksum is currently supported only for single URL downloads") {
		t.Fatalf("error = %q", err)
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

func writeChecksumFileForTest(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "checksums.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write checksum file: %v", err)
	}
	return path
}

func TestBuildPlan_ChecksumAndChecksumFileTogether(t *testing.T) {
	sum := sha256.Sum256([]byte("hello"))
	manifest := writeChecksumFileForTest(t, "sha256:"+fmt.Sprintf("%x", sum)+" https://example.com/file.zip\n")

	_, err := BuildPlan(Options{
		URLs:         []string{"https://example.com/file.zip"},
		Checksum:     "sha256:" + fmt.Sprintf("%x", sum),
		ChecksumFile: manifest,
		Retries:      3,
		Resume:       true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("error = %q", err)
	}
}

func TestBuildPlan_ChecksumFile_Valid(t *testing.T) {
	sumA := sha256.Sum256([]byte("alpha"))
	sumB := sha256.Sum256([]byte("beta"))
	manifest := writeChecksumFileForTest(t,
		"sha256:"+fmt.Sprintf("%x", sumA)+" https://example.com/a.zip\n"+
			"sha256:"+fmt.Sprintf("%x", sumB)+" https://example.com/b.zip\n")

	plan, err := BuildPlan(Options{
		URLs: []string{
			"https://example.com/a.zip",
			"https://example.com/b.zip",
		},
		ChecksumFile: manifest,
		Retries:      3,
		Resume:       true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	if !plan.HasChecksumFile {
		t.Fatal("plan.HasChecksumFile = false")
	}
	if len(plan.TargetChecksums) != 2 {
		t.Fatalf("len(plan.TargetChecksums) = %d, want 2", len(plan.TargetChecksums))
	}
	if spec, ok := plan.TargetChecksums["https://example.com/a.zip"]; !ok {
		t.Fatal("missing checksum for a.zip")
	} else if spec.Expected != fmt.Sprintf("%x", sumA) {
		t.Fatalf("a.zip checksum = %q", spec.Expected)
	}
}

func TestBuildPlan_ChecksumFile_MissingTargetURL(t *testing.T) {
	sumA := sha256.Sum256([]byte("alpha"))
	manifest := writeChecksumFileForTest(t, "sha256:"+fmt.Sprintf("%x", sumA)+" https://example.com/a.zip\n")

	_, err := BuildPlan(Options{
		URLs: []string{
			"https://example.com/a.zip",
			"https://example.com/b.zip",
		},
		ChecksumFile: manifest,
		Retries:      3,
		Resume:       true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "no checksum provided for URL") {
		t.Fatalf("error = %q", err)
	}
	if !strings.Contains(err.Error(), "https://example.com/b.zip") {
		t.Fatalf("error = %q, missing URL", err)
	}
}

func TestBuildPlan_ChecksumFile_ExtraManifestURL(t *testing.T) {
	sumA := sha256.Sum256([]byte("alpha"))
	sumB := sha256.Sum256([]byte("beta"))
	manifest := writeChecksumFileForTest(t,
		"sha256:"+fmt.Sprintf("%x", sumA)+" https://example.com/a.zip\n"+
			"sha256:"+fmt.Sprintf("%x", sumB)+" https://example.com/extra.zip\n")

	_, err := BuildPlan(Options{
		URLs:         []string{"https://example.com/a.zip"},
		ChecksumFile: manifest,
		Retries:      3,
		Resume:       true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "manifest URL not in download targets") {
		t.Fatalf("error = %q", err)
	}
	if !strings.Contains(err.Error(), "https://example.com/extra.zip") {
		t.Fatalf("error = %q, missing URL", err)
	}
}

func TestBuildPlan_ChecksumFile_DuplicateManifestURL(t *testing.T) {
	sumA := sha256.Sum256([]byte("alpha"))
	sumB := sha256.Sum256([]byte("beta"))
	manifest := writeChecksumFileForTest(t,
		"sha256:"+fmt.Sprintf("%x", sumA)+" https://example.com/a.zip\n"+
			"sha256:"+fmt.Sprintf("%x", sumB)+" https://example.com/a.zip\n")

	_, err := BuildPlan(Options{
		URLs:         []string{"https://example.com/a.zip"},
		ChecksumFile: manifest,
		Retries:      3,
		Resume:       true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "duplicate URL") {
		t.Fatalf("error = %q", err)
	}
	if !strings.Contains(err.Error(), "checksum file") {
		t.Fatalf("error = %q, missing checksum file context", err)
	}
}

func TestBuildPlan_ChecksumFile_SingleURL(t *testing.T) {
	sumA := sha256.Sum256([]byte("alpha"))
	manifest := writeChecksumFileForTest(t, "sha256:"+fmt.Sprintf("%x", sumA)+" https://example.com/a.zip\n")

	plan, err := BuildPlan(Options{
		URLs:         []string{"https://example.com/a.zip"},
		ChecksumFile: manifest,
		Retries:      3,
		Resume:       true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	if !plan.HasChecksumFile {
		t.Fatal("plan.HasChecksumFile = false")
	}
	if len(plan.TargetChecksums) != 1 {
		t.Fatalf("len(plan.TargetChecksums) = %d, want 1", len(plan.TargetChecksums))
	}
}

func TestBuildPlan_ChecksumFile_WithURLFile(t *testing.T) {
	urlFile := filepath.Join(t.TempDir(), "urls.txt")
	if err := os.WriteFile(urlFile, []byte("https://example.com/a.zip\nhttps://example.com/b.zip\n"), 0o600); err != nil {
		t.Fatalf("write URL file: %v", err)
	}
	sumA := sha256.Sum256([]byte("alpha"))
	sumB := sha256.Sum256([]byte("beta"))
	manifest := writeChecksumFileForTest(t,
		"sha256:"+fmt.Sprintf("%x", sumA)+" https://example.com/a.zip\n"+
			"sha256:"+fmt.Sprintf("%x", sumB)+" https://example.com/b.zip\n")

	plan, err := BuildPlan(Options{
		File:         urlFile,
		ChecksumFile: manifest,
		Retries:      3,
		Resume:       true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	if len(plan.TargetChecksums) != 2 {
		t.Fatalf("len(plan.TargetChecksums) = %d, want 2", len(plan.TargetChecksums))
	}
}

func TestBuildPlan_ChecksumFile_EmptyManifest(t *testing.T) {
	manifest := writeChecksumFileForTest(t, "# only comments\n\n")

	_, err := BuildPlan(Options{
		URLs:         []string{"https://example.com/a.zip"},
		ChecksumFile: manifest,
		Retries:      3,
		Resume:       true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "no checksum entries found") {
		t.Fatalf("error = %q", err)
	}
}

func TestBuildPlan_ChecksumFile_NotFound(t *testing.T) {
	_, err := BuildPlan(Options{
		URLs:         []string{"https://example.com/a.zip"},
		ChecksumFile: filepath.Join(t.TempDir(), "missing.txt"),
		Retries:      3,
		Resume:       true,
	})
	if err == nil {
		t.Fatal("BuildPlan returned nil error")
	}
	if !strings.Contains(err.Error(), "checksum file") {
		t.Fatalf("error = %q, missing checksum file context", err)
	}
}

func TestBuildPlan_DryRun_ChecksumFileShownInOutput(t *testing.T) {
	sumA := sha256.Sum256([]byte("alpha"))
	manifest := writeChecksumFileForTest(t, "sha256:"+fmt.Sprintf("%x", sumA)+" https://example.com/a.zip\n")

	plan, err := BuildPlan(Options{
		URLs:         []string{"https://example.com/a.zip"},
		ChecksumFile: manifest,
		DryRun:       true,
		Retries:      3,
		Resume:       true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	out := plan.DryRunString()
	if !strings.Contains(out, "Checksums: from file (1 entries)") {
		t.Fatalf("dry-run output missing checksum file line:\n%s", out)
	}
}

func TestBuildPlan_ExistingSingleChecksumBehaviorUnchanged(t *testing.T) {
	sum := sha256.Sum256([]byte("hello"))
	plan, err := BuildPlan(Options{
		URLs:     []string{"https://example.com/file.zip"},
		Checksum: "sha256:" + fmt.Sprintf("%x", sum),
		Retries:  3,
		Resume:   true,
	})
	if err != nil {
		t.Fatalf("BuildPlan returned error: %v", err)
	}
	if plan.Checksum == nil {
		t.Fatal("plan.Checksum = nil")
	}
	if plan.HasChecksumFile {
		t.Fatal("plan.HasChecksumFile = true, want false")
	}
	if plan.TargetChecksums != nil {
		t.Fatal("plan.TargetChecksums != nil")
	}
}

func TestValidateURLExported(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid https", "https://example.com/file.zip", false},
		{"valid http", "http://example.com/file.zip", false},
		{"ftp scheme", "ftp://example.com/file.zip", true},
		{"no scheme", "example.com/file.zip", true},
		{"empty host", "https:///file.zip", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURL(tc.url)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateURL(%q) error = %v, wantErr %v", tc.url, err, tc.wantErr)
			}
		})
	}
}
