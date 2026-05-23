package input

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadURLFileIgnoresEmptyAndCommentLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "urls.txt")
	content := `
# primary downloads
  https://example.com/file.zip

http://example.com/archive.tar.gz
   # another comment
`

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	got, err := ReadURLFile(path)
	if err != nil {
		t.Fatalf("ReadURLFile returned error: %v", err)
	}

	want := []string{
		"https://example.com/file.zip",
		"http://example.com/archive.tar.gz",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ReadURLFile() = %#v, want %#v", got, want)
	}
}

func TestReadURLFileRejectsEmptyPath(t *testing.T) {
	if _, err := ReadURLFile(" "); err == nil {
		t.Fatal("ReadURLFile returned nil error for empty path")
	}
}
