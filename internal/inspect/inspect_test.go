package inspect

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/he8um/daryaft/internal/httpopts"
)

func TestInspectValidatesURL(t *testing.T) {
	for _, rawURL := range []string{"", "ftp://example.com/file.zip", "https://"} {
		t.Run(rawURL, func(t *testing.T) {
			_, err := URL(context.Background(), Options{URL: rawURL})
			if err == nil {
				t.Fatal("URL returned nil error")
			}
		})
	}
}

func TestInspectFollowsRedirectAndReportsFinalURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/source":
			http.Redirect(w, r, "/final/file.zip", http.StatusFound)
		case "/final/file.zip":
			w.Header().Set("Content-Length", "12")
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("ETag", `"abc123"`)
			w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	result, err := URL(context.Background(), Options{URL: server.URL + "/source"})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}

	if result.URL != server.URL+"/source" {
		t.Fatalf("URL = %q, want source URL", result.URL)
	}
	if result.FinalURL != server.URL+"/final/file.zip" {
		t.Fatalf("FinalURL = %q, want final URL", result.FinalURL)
	}
	if result.Filename != "file.zip" {
		t.Fatalf("Filename = %q, want file.zip", result.Filename)
	}
}

func TestInspectReadsHeadMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Fatalf("method = %s, want HEAD", r.Method)
		}
		w.Header().Set("Content-Length", "1048576")
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("ETag", `"abc123"`)
		w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
	}))
	defer server.Close()

	result, err := URL(context.Background(), Options{URL: server.URL + "/file.zip"})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}

	if result.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want 200", result.StatusCode)
	}
	if !result.ContentLengthKnown || result.ContentLength != 1048576 {
		t.Fatalf("ContentLength = %d known %t, want 1048576 true", result.ContentLength, result.ContentLengthKnown)
	}
	if result.ContentType != "application/zip" {
		t.Fatalf("ContentType = %q, want application/zip", result.ContentType)
	}
	if result.AcceptRanges != "bytes" || !result.ResumeSupported || !result.ResumeSupportKnown {
		t.Fatalf("resume fields = %q %t %t, want bytes true true", result.AcceptRanges, result.ResumeSupported, result.ResumeSupportKnown)
	}
	if result.ETag != `"abc123"` {
		t.Fatalf("ETag = %q", result.ETag)
	}
	if result.LastModified != "Tue, 01 Jun 2026 12:00:00 GMT" {
		t.Fatalf("LastModified = %q", result.LastModified)
	}
}

func TestInspectReportsUnknownMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result, err := URL(context.Background(), Options{URL: server.URL + "/download"})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}

	if result.ContentLengthKnown {
		t.Fatalf("ContentLengthKnown = true, want false")
	}
	if result.ContentType != "" {
		t.Fatalf("ContentType = %q, want empty", result.ContentType)
	}
	if result.ResumeSupportKnown {
		t.Fatalf("ResumeSupportKnown = true, want false")
	}
	if result.Format() == "" || !strings.Contains(result.Format(), "Content length: unknown") {
		t.Fatalf("human output missing unknown content length:\n%s", result.Format())
	}
}

func TestInspectFallsBackWhenHeadReturnsMethodNotAllowed(t *testing.T) {
	var sawRange bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodGet:
			if r.Header.Get("Range") == "bytes=0-0" {
				sawRange = true
			}
			w.Header().Set("Content-Range", "bytes 0-0/1048576")
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Accept-Ranges", "bytes")
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write([]byte("x"))
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	result, err := URL(context.Background(), Options{URL: server.URL + "/file.zip"})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}
	if !sawRange {
		t.Fatal("GET fallback did not send Range header")
	}
	if result.StatusCode != http.StatusPartialContent {
		t.Fatalf("StatusCode = %d, want 206", result.StatusCode)
	}
	if !result.ContentLengthKnown || result.ContentLength != 1048576 {
		t.Fatalf("ContentLength = %d known %t, want 1048576 true", result.ContentLength, result.ContentLengthKnown)
	}
	if !result.ResumeSupported || !result.ResumeSupportKnown {
		t.Fatalf("resume = %t known %t, want true true", result.ResumeSupported, result.ResumeSupportKnown)
	}
}

func TestInspectFallsBackWhenHeadMetadataIsInsufficient(t *testing.T) {
	var methods []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Range", "bytes 0-0/2048")
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write([]byte("x"))
			return
		}
		w.Header().Set("Content-Length", "2048")
	}))
	defer server.Close()

	result, err := URL(context.Background(), Options{URL: server.URL + "/file.bin"})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}
	if strings.Join(methods, ",") != "HEAD,GET" {
		t.Fatalf("methods = %#v, want HEAD then GET", methods)
	}
	if result.ContentType != "application/octet-stream" {
		t.Fatalf("ContentType = %q, want fallback content type", result.ContentType)
	}
}

func TestInspectDoesNotWriteFiles(t *testing.T) {
	dir := t.TempDir()
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Accept-Ranges", "none")
		w.Header().Set("ETag", `"etag"`)
		w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
	}))
	defer server.Close()

	if _, err := URL(context.Background(), Options{URL: server.URL + "/file.txt"}); err != nil {
		t.Fatalf("URL returned error: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("inspect wrote files: %#v", entries)
	}
}

func TestInspectClosesResponseBodies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("ETag", `"etag"`)
		w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
	}))
	defer server.Close()

	transport := &closeTrackingTransport{base: server.Client().Transport}
	client := server.Client()
	client.Transport = transport

	if _, err := URL(context.Background(), Options{URL: server.URL + "/file.txt", Client: client}); err != nil {
		t.Fatalf("URL returned error: %v", err)
	}
	if transport.closed != 1 {
		t.Fatalf("closed bodies = %d, want 1", transport.closed)
	}
}

func TestInspectCustomHeader(t *testing.T) {
	var gotHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Custom")
		w.Header().Set("Content-Length", "4")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("ETag", `"e1"`)
		w.Header().Set("Last-Modified", "Tue, 01 Jun 2026 12:00:00 GMT")
	}))
	defer server.Close()

	headers, _ := httpopts.ParseHeaders([]string{"X-Custom: headerval"})
	_, err := URL(context.Background(), Options{
		URL:         server.URL + "/file.txt",
		HTTPOptions: httpopts.Options{Headers: headers},
	})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}
	if gotHeader != "headerval" {
		t.Errorf("X-Custom = %q, want headerval", gotHeader)
	}
}

func TestInspectCustomUserAgent(t *testing.T) {
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

	_, err := URL(context.Background(), Options{
		URL:         server.URL + "/file.txt",
		HTTPOptions: httpopts.Options{UserAgent: "InspectAgent/3.0"},
	})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}
	if gotUA != "InspectAgent/3.0" {
		t.Errorf("User-Agent = %q, want InspectAgent/3.0", gotUA)
	}
}

func TestInspectBasicAuth(t *testing.T) {
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

	_, err := URL(context.Background(), Options{
		URL:         server.URL + "/file.txt",
		HTTPOptions: httpopts.Options{Username: "bob", Password: "pass123"},
	})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}
	if gotUser != "bob" || gotPass != "pass123" {
		t.Errorf("BasicAuth = %q/%q, want bob/pass123", gotUser, gotPass)
	}
}

func TestInspectFallbackGETSendsCustomHeader(t *testing.T) {
	var headHeader, getHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			headHeader = r.Header.Get("X-Probe")
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodGet:
			getHeader = r.Header.Get("X-Probe")
			w.Header().Set("Content-Range", "bytes 0-0/100")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write([]byte("x"))
		}
	}))
	defer server.Close()

	headers, _ := httpopts.ParseHeaders([]string{"X-Probe: probe"})
	_, err := URL(context.Background(), Options{
		URL:         server.URL + "/file.txt",
		HTTPOptions: httpopts.Options{Headers: headers},
	})
	if err != nil {
		t.Fatalf("URL returned error: %v", err)
	}
	if headHeader != "probe" {
		t.Errorf("HEAD X-Probe = %q, want probe", headHeader)
	}
	if getHeader != "probe" {
		t.Errorf("GET fallback X-Probe = %q, want probe", getHeader)
	}
}

func TestInspectRejectsInvalidHTTPOptions(t *testing.T) {
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
			_, err := URL(context.Background(), Options{
				URL:         "https://example.com/file.txt",
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

type closeTrackingTransport struct {
	base   http.RoundTripper
	closed int
}

func (t *closeTrackingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.base.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	response.Body = &closeTrackingBody{ReadCloser: response.Body, onClose: func() { t.closed++ }}
	return response, nil
}

type closeTrackingBody struct {
	io.ReadCloser
	onClose func()
}

func (b *closeTrackingBody) Close() error {
	b.onClose()
	return b.ReadCloser.Close()
}
