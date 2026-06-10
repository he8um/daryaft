package checksum

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeManifestTestFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "checksums.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write manifest file: %v", err)
	}
	return path
}

func TestParseManifestFile_Valid(t *testing.T) {
	sha256Hex := strings.Repeat("a", 64)
	sha512Hex := strings.Repeat("b", 128)
	content := "sha256:" + sha256Hex + " https://example.com/file1.zip\n" +
		"sha512:" + sha512Hex + " https://example.com/file2.tar.gz\n"
	path := writeManifestTestFile(t, content)

	entries, err := ParseManifestFile(path)
	if err != nil {
		t.Fatalf("ParseManifestFile returned error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}
	if spec, ok := entries["https://example.com/file1.zip"]; !ok {
		t.Fatal("missing entry for file1.zip")
	} else if spec.Algorithm != AlgorithmSHA256 || spec.Expected != sha256Hex {
		t.Fatalf("file1 spec = %+v", spec)
	}
	if spec, ok := entries["https://example.com/file2.tar.gz"]; !ok {
		t.Fatal("missing entry for file2.tar.gz")
	} else if spec.Algorithm != AlgorithmSHA512 || spec.Expected != sha512Hex {
		t.Fatalf("file2 spec = %+v", spec)
	}
}

func TestParseManifestFile_IgnoresBlankLines(t *testing.T) {
	content := "\n\nsha256:" + strings.Repeat("a", 64) + " https://example.com/file.zip\n\n"
	path := writeManifestTestFile(t, content)

	entries, err := ParseManifestFile(path)
	if err != nil {
		t.Fatalf("ParseManifestFile returned error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(entries))
	}
}

func TestParseManifestFile_IgnoresCommentLines(t *testing.T) {
	content := "# primary downloads\n" +
		"  # indented comment\n" +
		"sha256:" + strings.Repeat("a", 64) + " https://example.com/file.zip\n"
	path := writeManifestTestFile(t, content)

	entries, err := ParseManifestFile(path)
	if err != nil {
		t.Fatalf("ParseManifestFile returned error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(entries))
	}
}

func TestParseManifestFile_MalformedLine_WrongFieldCount(t *testing.T) {
	content := "sha256:" + strings.Repeat("a", 64) + "\n"
	path := writeManifestTestFile(t, content)

	_, err := ParseManifestFile(path)
	if err == nil {
		t.Fatal("ParseManifestFile returned nil error")
	}
	if !strings.Contains(err.Error(), "manifest line 1") {
		t.Fatalf("error = %q, missing line number", err)
	}
	if !strings.Contains(err.Error(), "format") {
		t.Fatalf("error = %q, missing format context", err)
	}
}

func TestParseManifestFile_MalformedLine_ThreeFields(t *testing.T) {
	content := "sha256:" + strings.Repeat("a", 64) + " https://example.com/file.zip extra\n"
	path := writeManifestTestFile(t, content)

	_, err := ParseManifestFile(path)
	if err == nil {
		t.Fatal("ParseManifestFile returned nil error")
	}
	if !strings.Contains(err.Error(), "manifest line 1") {
		t.Fatalf("error = %q, missing line number", err)
	}
}

func TestParseManifestFile_MalformedLine_InvalidChecksum(t *testing.T) {
	content := "sha256:" + strings.Repeat("a", 64) + " https://example.com/file1.zip\n" +
		"md5:" + strings.Repeat("a", 32) + " https://example.com/file2.zip\n"
	path := writeManifestTestFile(t, content)

	_, err := ParseManifestFile(path)
	if err == nil {
		t.Fatal("ParseManifestFile returned nil error")
	}
	if !strings.Contains(err.Error(), "manifest line 2") {
		t.Fatalf("error = %q, missing line number", err)
	}
	if !strings.Contains(err.Error(), "unsupported checksum algorithm") {
		t.Fatalf("error = %q, missing checksum parse error", err)
	}
}

func TestParseManifestFile_DuplicateURL(t *testing.T) {
	content := "sha256:" + strings.Repeat("a", 64) + " https://example.com/file.zip\n" +
		"sha256:" + strings.Repeat("b", 64) + " https://example.com/file.zip\n"
	path := writeManifestTestFile(t, content)

	_, err := ParseManifestFile(path)
	if err == nil {
		t.Fatal("ParseManifestFile returned nil error")
	}
	if !strings.Contains(err.Error(), "manifest line 2") {
		t.Fatalf("error = %q, missing line number", err)
	}
	if !strings.Contains(err.Error(), "duplicate URL") {
		t.Fatalf("error = %q, missing duplicate context", err)
	}
	if !strings.Contains(err.Error(), "https://example.com/file.zip") {
		t.Fatalf("error = %q, missing URL", err)
	}
}

func TestParseManifestFile_EmptyFile(t *testing.T) {
	path := writeManifestTestFile(t, "")

	entries, err := ParseManifestFile(path)
	if err != nil {
		t.Fatalf("ParseManifestFile returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("len(entries) = %d, want 0", len(entries))
	}
}

func TestParseManifestFile_FileNotFound(t *testing.T) {
	_, err := ParseManifestFile(filepath.Join(t.TempDir(), "does-not-exist.txt"))
	if err == nil {
		t.Fatal("ParseManifestFile returned nil error for missing file")
	}
}

func TestParseManifestFile_OnlyComments(t *testing.T) {
	content := "# only comments\n# nothing else\n"
	path := writeManifestTestFile(t, content)

	entries, err := ParseManifestFile(path)
	if err != nil {
		t.Fatalf("ParseManifestFile returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("len(entries) = %d, want 0", len(entries))
	}
}
