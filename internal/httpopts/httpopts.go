package httpopts

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

// Header holds a parsed HTTP header name/value pair.
type Header struct {
	Name  string
	Value string
}

// Options holds CLI-provided HTTP customization options.
type Options struct {
	ProxyURL  string
	Headers   []Header
	UserAgent string
	Username  string
	Password  string
}

// RedactedOptions is a display-safe copy of Options.
type RedactedOptions struct {
	ProxyURL  string
	Headers   []Header
	UserAgent string
	Username  string
	Password  string
}

// ParseHeader parses a single "Name: Value" string.
func ParseHeader(input string) (Header, error) {
	idx := strings.Index(input, ":")
	if idx < 0 {
		return Header{}, fmt.Errorf("header %q must contain a colon", input)
	}
	name := strings.TrimSpace(input[:idx])
	value := strings.TrimSpace(input[idx+1:])

	if name == "" {
		return Header{}, fmt.Errorf("header name is empty in %q", input)
	}
	if err := validateHeaderName(name); err != nil {
		return Header{}, err
	}
	if err := validateHeaderValue(value); err != nil {
		return Header{}, err
	}
	return Header{Name: name, Value: value}, nil
}

// ParseHeaders parses a slice of "Name: Value" strings.
func ParseHeaders(inputs []string) ([]Header, error) {
	result := make([]Header, 0, len(inputs))
	for _, s := range inputs {
		h, err := ParseHeader(s)
		if err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

// Validate checks that all options are valid.
func Validate(opts Options) error {
	if err := validateProxy(opts.ProxyURL); err != nil {
		return err
	}
	for _, h := range opts.Headers {
		if err := validateHeaderName(h.Name); err != nil {
			return err
		}
		if err := validateHeaderValue(h.Value); err != nil {
			return err
		}
	}
	if err := validateUserAgent(opts.UserAgent); err != nil {
		return err
	}
	if err := validateAuth(opts); err != nil {
		return err
	}
	return nil
}

// HasAny reports whether any option is set.
func HasAny(opts Options) bool {
	return opts.ProxyURL != "" ||
		len(opts.Headers) > 0 ||
		opts.UserAgent != "" ||
		opts.Username != "" ||
		opts.Password != ""
}

// sensitiveHeaderNames lists header name patterns (lowercase) that must be redacted.
var sensitiveHeaderNames = []string{
	"authorization",
	"proxy-authorization",
	"cookie",
	"set-cookie",
	"x-api-key",
	"x-token",
	"x-auth-token",
	"password",
	"secret",
	"token",
	"key",
}

func isSensitiveHeader(name string) bool {
	lower := strings.ToLower(name)
	for _, s := range sensitiveHeaderNames {
		if lower == s || strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

// Redact returns a display-safe copy of Options.
func Redact(opts Options) RedactedOptions {
	redacted := RedactedOptions{
		ProxyURL:  opts.ProxyURL,
		UserAgent: opts.UserAgent,
		Username:  opts.Username,
	}
	if opts.Password != "" {
		redacted.Password = "[REDACTED]"
	}
	if len(opts.Headers) > 0 {
		redacted.Headers = make([]Header, len(opts.Headers))
		for i, h := range opts.Headers {
			if isSensitiveHeader(h.Name) {
				redacted.Headers[i] = Header{Name: h.Name, Value: "[REDACTED]"}
			} else {
				redacted.Headers[i] = h
			}
		}
	}
	return redacted
}

// ApplyToRequest sets headers, user-agent, and basic auth on req.
// Precedence: --user-agent overrides any User-Agent custom header.
// Basic auth is set if Username is non-empty.
func ApplyToRequest(req *http.Request, opts Options) {
	for _, h := range opts.Headers {
		req.Header.Set(h.Name, h.Value)
	}
	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}
	if opts.Username != "" {
		req.SetBasicAuth(opts.Username, opts.Password)
	}
}

// NewClient returns a client configured for proxy from opts, cloning base to
// avoid mutating shared state. Returns base unchanged when no proxy is set.
func NewClient(base *http.Client, opts Options) (*http.Client, error) {
	if opts.ProxyURL == "" {
		return base, nil
	}

	proxyURL, err := url.Parse(opts.ProxyURL)
	if err != nil {
		return nil, fmt.Errorf("parse proxy URL: %w", err)
	}

	var cloned http.Client
	if base != nil {
		cloned = *base
	}

	var transport *http.Transport
	if base != nil && base.Transport != nil {
		t, ok := base.Transport.(*http.Transport)
		if !ok {
			return nil, fmt.Errorf("base transport is not *http.Transport")
		}
		transport = t.Clone()
	} else {
		transport = http.DefaultTransport.(*http.Transport).Clone()
	}

	transport.Proxy = http.ProxyURL(proxyURL)
	cloned.Transport = transport
	return &cloned, nil
}

// --- internal validation helpers ---

func validateProxy(proxyURL string) error {
	if proxyURL == "" {
		return nil
	}
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("proxy scheme must be http or https (got %q)", parsed.Scheme)
	}
	if parsed.Host == "" {
		return fmt.Errorf("proxy URL must include a host")
	}
	return nil
}

func validateHeaderName(name string) error {
	if name == "" {
		return fmt.Errorf("header name cannot be empty")
	}
	for _, r := range name {
		if r <= 0x1F || r == 0x7F {
			return fmt.Errorf("header name %q contains control character", name)
		}
		if !isTokenChar(r) {
			return fmt.Errorf("header name %q contains invalid character %q", name, r)
		}
	}
	return nil
}

func validateHeaderValue(value string) error {
	for _, r := range value {
		if unicode.IsControl(r) && r != '\t' {
			return fmt.Errorf("header value contains control character")
		}
	}
	return nil
}

func validateUserAgent(ua string) error {
	if ua == "" {
		return nil
	}
	for _, r := range ua {
		if unicode.IsControl(r) {
			return fmt.Errorf("user-agent contains control character")
		}
	}
	return nil
}

func validateAuth(opts Options) error {
	if opts.Password != "" && opts.Username == "" {
		return fmt.Errorf("--password requires --username")
	}
	if opts.Username != "" {
		for _, h := range opts.Headers {
			if strings.ToLower(h.Name) == "authorization" {
				return fmt.Errorf("use either --username/--password or an Authorization header, not both")
			}
		}
	}
	return nil
}

// isTokenChar reports if r is a valid HTTP token character (RFC 7230).
func isTokenChar(r rune) bool {
	const tokenChars = "!#$%&'*+-.^_`|~" // #nosec G101 -- RFC 7230 token characters, not credentials
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}
	return strings.ContainsRune(tokenChars, r)
}
