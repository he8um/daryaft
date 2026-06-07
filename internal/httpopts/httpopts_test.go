package httpopts_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/he8um/daryaft/internal/httpopts"
)

// --- ParseHeader ---

func TestParseHeader_Valid(t *testing.T) {
	h, err := httpopts.ParseHeader("X-Custom: value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Name != "X-Custom" || h.Value != "value" {
		t.Errorf("got %+v", h)
	}
}

func TestParseHeader_ValueContainsColon(t *testing.T) {
	h, err := httpopts.ParseHeader("X-Custom: http://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Value != "http://example.com" {
		t.Errorf("got value %q", h.Value)
	}
}

func TestParseHeader_EmptyValue(t *testing.T) {
	h, err := httpopts.ParseHeader("X-Empty:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Name != "X-Empty" || h.Value != "" {
		t.Errorf("got %+v", h)
	}
}

func TestParseHeader_MissingColon(t *testing.T) {
	_, err := httpopts.ParseHeader("NoColon")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseHeader_EmptyName(t *testing.T) {
	_, err := httpopts.ParseHeader(": value")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestParseHeader_InvalidName(t *testing.T) {
	_, err := httpopts.ParseHeader("Bad Name: value")
	if err == nil {
		t.Fatal("expected error for invalid name with space")
	}
}

func TestParseHeader_ControlCharInName(t *testing.T) {
	_, err := httpopts.ParseHeader("X-\x01Bad: value")
	if err == nil {
		t.Fatal("expected error for control char in name")
	}
}

func TestParseHeader_ControlCharInValue(t *testing.T) {
	_, err := httpopts.ParseHeader("X-Good: val\x01ue")
	if err == nil {
		t.Fatal("expected error for control char in value")
	}
}

// --- ParseHeaders ---

func TestParseHeaders_Empty(t *testing.T) {
	headers, err := httpopts.ParseHeaders(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(headers) != 0 {
		t.Errorf("expected empty slice")
	}
}

func TestParseHeaders_MultipleValid(t *testing.T) {
	headers, err := httpopts.ParseHeaders([]string{"X-A: 1", "X-B: 2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(headers) != 2 {
		t.Errorf("expected 2, got %d", len(headers))
	}
}

func TestParseHeaders_OneInvalid(t *testing.T) {
	_, err := httpopts.ParseHeaders([]string{"X-A: 1", "NoColon"})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Validate proxy ---

func TestValidate_Proxy_ValidHTTP(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{ProxyURL: "http://proxy.example.com:8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_Proxy_ValidHTTPS(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{ProxyURL: "https://proxy.example.com:443"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_Proxy_MissingScheme(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{ProxyURL: "proxy.example.com:8080"})
	if err == nil {
		t.Fatal("expected error for missing scheme")
	}
}

func TestValidate_Proxy_MissingHost(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{ProxyURL: "http://"})
	if err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestValidate_Proxy_UnsupportedScheme(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{ProxyURL: "ftp://proxy.example.com"})
	if err == nil {
		t.Fatal("expected error for ftp scheme")
	}
}

func TestValidate_Proxy_SOCKS5Rejected(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{ProxyURL: "socks5://proxy.example.com:1080"})
	if err == nil {
		t.Fatal("expected error for socks5")
	}
}

func TestValidate_Proxy_InvalidURL(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{ProxyURL: "://bad"})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

// --- Validate user-agent ---

func TestValidate_UserAgent_Valid(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{UserAgent: "MyApp/1.0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_UserAgent_Empty(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{UserAgent: ""})
	if err != nil {
		t.Fatalf("empty user-agent should be ok: %v", err)
	}
}

func TestValidate_UserAgent_ControlChar(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{UserAgent: "Bad\x01Agent"})
	if err == nil {
		t.Fatal("expected error for control char in user-agent")
	}
}

// --- Validate auth ---

func TestValidate_Auth_UsernameAndPassword(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{Username: "alice", Password: "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_Auth_UsernameOnly(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{Username: "alice"})
	if err != nil {
		t.Fatalf("username with empty password should be ok: %v", err)
	}
}

func TestValidate_Auth_PasswordWithoutUsername(t *testing.T) {
	err := httpopts.Validate(httpopts.Options{Password: "secret"})
	if err == nil {
		t.Fatal("expected error: password without username")
	}
}

func TestValidate_Auth_BasicAuthPlusAuthorizationHeader(t *testing.T) {
	headers, _ := httpopts.ParseHeaders([]string{"Authorization: Bearer token"})
	err := httpopts.Validate(httpopts.Options{
		Username: "alice",
		Password: "secret",
		Headers:  headers,
	})
	if err == nil {
		t.Fatal("expected error: basic auth + Authorization header")
	}
}

// --- Redact ---

func TestRedact_PasswordRedacted(t *testing.T) {
	r := httpopts.Redact(httpopts.Options{Password: "secret"})
	if r.Password != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", r.Password)
	}
}

func TestRedact_EmptyPasswordNotRedacted(t *testing.T) {
	r := httpopts.Redact(httpopts.Options{})
	if r.Password != "" {
		t.Errorf("expected empty, got %q", r.Password)
	}
}

func TestRedact_AuthorizationHeaderRedacted(t *testing.T) {
	headers, _ := httpopts.ParseHeaders([]string{"Authorization: Bearer token123"})
	r := httpopts.Redact(httpopts.Options{Headers: headers})
	if r.Headers[0].Value != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", r.Headers[0].Value)
	}
}

func TestRedact_ProxyAuthorizationHeaderRedacted(t *testing.T) {
	headers, _ := httpopts.ParseHeaders([]string{"Proxy-Authorization: Basic xyz"})
	r := httpopts.Redact(httpopts.Options{Headers: headers})
	if r.Headers[0].Value != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", r.Headers[0].Value)
	}
}

func TestRedact_CookieHeaderRedacted(t *testing.T) {
	headers, _ := httpopts.ParseHeaders([]string{"Cookie: session=abc"})
	r := httpopts.Redact(httpopts.Options{Headers: headers})
	if r.Headers[0].Value != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", r.Headers[0].Value)
	}
}

func TestRedact_XAPIKeyRedacted(t *testing.T) {
	headers, _ := httpopts.ParseHeaders([]string{"X-Api-Key: mykey"})
	r := httpopts.Redact(httpopts.Options{Headers: headers})
	if r.Headers[0].Value != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", r.Headers[0].Value)
	}
}

func TestRedact_CustomHeaderNotRedacted(t *testing.T) {
	headers, _ := httpopts.ParseHeaders([]string{"X-Custom: myvalue"})
	r := httpopts.Redact(httpopts.Options{Headers: headers})
	if r.Headers[0].Value != "myvalue" {
		t.Errorf("expected myvalue, got %q", r.Headers[0].Value)
	}
}

// --- ApplyToRequest ---

func TestApplyToRequest_CustomHeadersSet(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	headers, _ := httpopts.ParseHeaders([]string{"X-Custom: testval"})
	httpopts.ApplyToRequest(req, httpopts.Options{Headers: headers})
	if req.Header.Get("X-Custom") != "testval" {
		t.Errorf("expected testval, got %q", req.Header.Get("X-Custom"))
	}
}

func TestApplyToRequest_UserAgentSet(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	httpopts.ApplyToRequest(req, httpopts.Options{UserAgent: "TestAgent/1.0"})
	if req.Header.Get("User-Agent") != "TestAgent/1.0" {
		t.Errorf("expected TestAgent/1.0, got %q", req.Header.Get("User-Agent"))
	}
}

func TestApplyToRequest_UserAgentFlagWinsOverHeader(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	headers, _ := httpopts.ParseHeaders([]string{"User-Agent: HeaderAgent/1.0"})
	httpopts.ApplyToRequest(req, httpopts.Options{
		Headers:   headers,
		UserAgent: "FlagAgent/2.0",
	})
	if req.Header.Get("User-Agent") != "FlagAgent/2.0" {
		t.Errorf("expected FlagAgent/2.0, got %q", req.Header.Get("User-Agent"))
	}
}

func TestApplyToRequest_BasicAuthSet(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	httpopts.ApplyToRequest(req, httpopts.Options{Username: "alice", Password: "pass"})
	u, p, ok := req.BasicAuth()
	if !ok || u != "alice" || p != "pass" {
		t.Errorf("expected basic auth alice:pass, got ok=%v u=%q p=%q", ok, u, p)
	}
}

// --- NewClient ---

func TestNewClient_NoProxy_ReturnsSameClient(t *testing.T) {
	base := &http.Client{Timeout: 10 * time.Second}
	got, err := httpopts.NewClient(base, httpopts.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if got != base {
		t.Error("expected same client when no proxy is set")
	}
}

func TestNewClient_ProxySet(t *testing.T) {
	base := &http.Client{
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
	got, err := httpopts.NewClient(base, httpopts.Options{ProxyURL: "http://127.0.0.1:8888"})
	if err != nil {
		t.Fatal(err)
	}
	if got == base {
		t.Error("expected cloned client, not original")
	}
	transport, ok := got.Transport.(*http.Transport)
	if !ok {
		t.Fatal("transport not *http.Transport")
	}
	if transport.Proxy == nil {
		t.Error("expected Proxy to be set")
	}
}

func TestNewClient_ProxyPreservesTimeouts(t *testing.T) {
	base := &http.Client{Timeout: 42 * time.Second}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = 15 * time.Second
	base.Transport = transport

	got, err := httpopts.NewClient(base, httpopts.Options{ProxyURL: "http://127.0.0.1:8888"})
	if err != nil {
		t.Fatal(err)
	}
	gotTransport, ok := got.Transport.(*http.Transport)
	if !ok {
		t.Fatal("transport not *http.Transport")
	}
	if gotTransport.ResponseHeaderTimeout != 15*time.Second {
		t.Errorf("expected 15s, got %v", gotTransport.ResponseHeaderTimeout)
	}
}

func TestNewClient_ProxyActuallyUsed(t *testing.T) {
	var proxyCalled bool
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyCalled = true
		http.Error(w, "proxy ok", http.StatusOK)
	}))
	defer proxyServer.Close()

	base := &http.Client{
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
	client, err := httpopts.NewClient(base, httpopts.Options{ProxyURL: proxyServer.URL})
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/test", nil)
	resp, err := client.Do(req)
	if err == nil {
		_ = resp.Body.Close()
	}

	if !proxyCalled {
		t.Error("proxy server was not called")
	}
}
