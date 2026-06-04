package checksum

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseValidSHA256(t *testing.T) {
	expected := strings.Repeat("a", 64)
	spec, err := Parse(" sha256:" + expected + " ")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if spec.Algorithm != AlgorithmSHA256 {
		t.Fatalf("Algorithm = %q", spec.Algorithm)
	}
	if spec.Expected != expected {
		t.Fatalf("Expected = %q", spec.Expected)
	}
}

func TestParseValidSHA512(t *testing.T) {
	expected := strings.Repeat("B", 128)
	spec, err := Parse("sha512:" + expected)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if spec.Algorithm != AlgorithmSHA512 {
		t.Fatalf("Algorithm = %q", spec.Algorithm)
	}
	if spec.Expected != strings.ToLower(expected) {
		t.Fatalf("Expected = %q", spec.Expected)
	}
}

func TestParseRejectsMissingColon(t *testing.T) {
	_, err := Parse(strings.Repeat("a", 64))
	if err == nil {
		t.Fatal("Parse returned nil error")
	}
	if !strings.Contains(err.Error(), "format") {
		t.Fatalf("error = %q, want format context", err)
	}
}

func TestParseRejectsUnsupportedAlgorithm(t *testing.T) {
	_, err := Parse("md5:" + strings.Repeat("a", 32))
	if err == nil {
		t.Fatal("Parse returned nil error")
	}
	if !strings.Contains(err.Error(), "unsupported checksum algorithm") {
		t.Fatalf("error = %q", err)
	}
}

func TestParseRejectsInvalidHex(t *testing.T) {
	_, err := Parse("sha256:" + strings.Repeat("g", 64))
	if err == nil {
		t.Fatal("Parse returned nil error")
	}
	if !strings.Contains(err.Error(), "valid hexadecimal") {
		t.Fatalf("error = %q", err)
	}
}

func TestParseRejectsWrongLength(t *testing.T) {
	_, err := Parse("sha256:" + strings.Repeat("a", 63))
	if err == nil {
		t.Fatal("Parse returned nil error")
	}
	if !strings.Contains(err.Error(), "64 hex characters") {
		t.Fatalf("error = %q", err)
	}
}

func TestVerifyFileComparesCaseInsensitiveExpectedChecksum(t *testing.T) {
	path := writeChecksumTestFile(t, "hello")
	sum := sha256.Sum256([]byte("hello"))
	expected := strings.ToUpper(fmt.Sprintf("%x", sum))

	spec := Spec{
		Algorithm: AlgorithmSHA256,
		Expected:  expected,
	}
	actual, err := VerifyFile(path, spec)
	if err != nil {
		t.Fatalf("VerifyFile returned error: %v", err)
	}
	want := strings.ToLower(expected)
	if actual != want {
		t.Fatalf("actual = %q, want %q", actual, want)
	}
}

func TestComputeFileSHA256(t *testing.T) {
	path := writeChecksumTestFile(t, "hello")
	sum := sha256.Sum256([]byte("hello"))
	want := fmt.Sprintf("%x", sum)

	got, err := ComputeFile(path, AlgorithmSHA256)
	if err != nil {
		t.Fatalf("ComputeFile returned error: %v", err)
	}
	if got != want {
		t.Fatalf("ComputeFile = %q, want %q", got, want)
	}
}

func TestComputeFileSHA512(t *testing.T) {
	path := writeChecksumTestFile(t, "hello")
	sum := sha512.Sum512([]byte("hello"))
	want := fmt.Sprintf("%x", sum)

	got, err := ComputeFile(path, AlgorithmSHA512)
	if err != nil {
		t.Fatalf("ComputeFile returned error: %v", err)
	}
	if got != want {
		t.Fatalf("ComputeFile = %q, want %q", got, want)
	}
}

func writeChecksumTestFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}
	return path
}
