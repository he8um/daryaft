package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/inspect"
)

func TestInspectCommandRequiresExactlyOneURL(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"https://example.com/a", "https://example.com/b"},
	} {
		t.Run(strings.Join(args, ","), func(t *testing.T) {
			_, err := executeInspectCommand(t, args...)
			if err == nil {
				t.Fatal("inspect returned nil error")
			}
		})
	}
}

func TestInspectCommandRejectsInvalidURL(t *testing.T) {
	_, err := executeInspectCommand(t, "ftp://example.com/file.zip")
	if err == nil {
		t.Fatal("inspect returned nil error")
	}
	if !strings.Contains(err.Error(), "scheme must be http or https") {
		t.Fatalf("error = %q", err)
	}
}

func TestInspectCommandHumanOutput(t *testing.T) {
	server := inspectCommandServer()
	defer server.Close()

	output, err := executeInspectCommand(t, server.URL+"/file.zip")
	if err != nil {
		t.Fatalf("inspect returned error: %v\n%s", err, output)
	}

	for _, want := range []string{
		"Daryaft inspect",
		"URL: " + server.URL + "/file.zip",
		"Final URL: " + server.URL + "/file.zip",
		"Status: 200 OK",
		"Filename: file.zip",
		"Content length: 1048576 bytes",
		"Content type: application/zip",
		"Accept-Ranges: bytes",
		"Resume supported: yes",
		`ETag: "abc123"`,
		"Last-Modified: Tue, 01 Jun 2026 12:00:00 GMT",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestInspectCommandJSONOutput(t *testing.T) {
	server := inspectCommandServer()
	defer server.Close()

	output, err := executeInspectCommand(t, server.URL+"/file.zip", "--json")
	if err != nil {
		t.Fatalf("inspect --json returned error: %v\n%s", err, output)
	}

	var got inspect.Result
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v\n%s", err, output)
	}
	if got.URL != server.URL+"/file.zip" {
		t.Fatalf("URL = %q", got.URL)
	}
	if got.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want 200", got.StatusCode)
	}
	if got.ContentLength != 1048576 || !got.ContentLengthKnown {
		t.Fatalf("content length = %d known %t", got.ContentLength, got.ContentLengthKnown)
	}
	if !got.ResumeSupported || !got.ResumeSupportKnown {
		t.Fatalf("resume = %t known %t", got.ResumeSupported, got.ResumeSupportKnown)
	}
	if strings.Contains(output, "Daryaft inspect") {
		t.Fatalf("JSON output contains human text:\n%s", output)
	}
}

func executeInspectCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var output bytes.Buffer
	command := newInspectCommand()
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs(args)
	err := command.Execute()
	return output.String(), err
}

func inspectCommandServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1048576")
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("ETag", `"abc123"`)
		w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
	}))
}

func TestInspectCLICustomHeader(t *testing.T) {
	var gotHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Inspect")
		w.Header().Set("Content-Length", "4")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("ETag", `"e1"`)
		w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
	}))
	defer server.Close()

	_, err := executeInspectCommand(t, server.URL+"/file.txt", "--header", "X-Inspect: inspectval")
	if err != nil {
		t.Fatalf("inspect returned error: %v", err)
	}
	if gotHeader != "inspectval" {
		t.Errorf("X-Inspect = %q, want inspectval", gotHeader)
	}
}

func TestInspectCLIUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.UserAgent()
		w.Header().Set("Content-Length", "4")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("ETag", `"e1"`)
		w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
	}))
	defer server.Close()

	_, err := executeInspectCommand(t, server.URL+"/file.txt", "--user-agent", "InspectCLIAgent/5.0")
	if err != nil {
		t.Fatalf("inspect returned error: %v", err)
	}
	if gotUA != "InspectCLIAgent/5.0" {
		t.Errorf("User-Agent = %q, want InspectCLIAgent/5.0", gotUA)
	}
}

func TestInspectCLIBasicAuth(t *testing.T) {
	var gotUser, gotPass string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass, _ = r.BasicAuth()
		w.Header().Set("Content-Length", "4")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("ETag", `"e1"`)
		w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
	}))
	defer server.Close()

	_, err := executeInspectCommand(t, server.URL+"/file.txt", "--username", "carol", "--password", "pw123")
	if err != nil {
		t.Fatalf("inspect returned error: %v", err)
	}
	if gotUser != "carol" || gotPass != "pw123" {
		t.Errorf("BasicAuth = %q/%q, want carol/pw123", gotUser, gotPass)
	}
}

func TestInspectCLIRejectsInvalidHeader(t *testing.T) {
	_, err := executeInspectCommand(t, "https://example.com/file.txt", "--header", "NoColon")
	if err == nil {
		t.Fatal("expected error for invalid header")
	}
}

func TestInspectCLIRejectsInvalidProxy(t *testing.T) {
	_, err := executeInspectCommand(t, "https://example.com/file.txt", "--proxy", "socks5://proxy:1080")
	if err == nil {
		t.Fatal("expected error for socks5 proxy")
	}
}

func TestInspectCLIRejectsPasswordWithoutUsername(t *testing.T) {
	_, err := executeInspectCommand(t, "https://example.com/file.txt", "--password", "secret")
	if err == nil {
		t.Fatal("expected error: password without username")
	}
}

func TestInspectCLIRejectsAuthHeaderPlusBasicAuth(t *testing.T) {
	_, err := executeInspectCommand(t,
		"https://example.com/file.txt",
		"--username", "alice",
		"--password", "pass",
		"--header", "Authorization: Bearer tok",
	)
	if err == nil {
		t.Fatal("expected error: basic auth + Authorization header")
	}
}
