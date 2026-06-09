package download

import (
	"fmt"
	"strings"

	"github.com/he8um/daryaft/internal/checksum"
	"github.com/he8um/daryaft/internal/httpopts"
)

type Plan struct {
	URLs        []string
	Output      string
	Name        string
	Checksum    *checksum.Spec
	Retries     int
	Resume      bool
	HTTPOptions httpopts.Options
}

func (p Plan) DryRunString() string {
	var builder strings.Builder

	fmt.Fprintln(&builder, "Daryaft download plan")
	fmt.Fprintf(&builder, "URLs: %d\n", len(p.URLs))
	for index, rawURL := range p.URLs {
		fmt.Fprintf(&builder, "%d. %s\n", index+1, rawURL)
	}
	fmt.Fprintf(&builder, "Output: %s\n", valueOrDefault(p.Output, "current directory"))
	fmt.Fprintf(&builder, "Filename: %s\n", valueOrDefault(p.Name, "auto-detect"))
	if p.Checksum != nil {
		fmt.Fprintf(&builder, "Checksum: %s\n", p.Checksum.String())
	}
	fmt.Fprintf(&builder, "Retries: %d\n", p.Retries)
	fmt.Fprintf(&builder, "Resume: %t\n", p.Resume)
	if httpopts.HasAny(p.HTTPOptions) {
		redacted := httpopts.Redact(p.HTTPOptions)
		fmt.Fprintln(&builder, "HTTP")
		if redacted.ProxyURL != "" {
			fmt.Fprintf(&builder, "  Proxy: %s\n", redacted.ProxyURL)
		}
		if len(redacted.Headers) > 0 {
			fmt.Fprintln(&builder, "  Headers:")
			for _, h := range redacted.Headers {
				fmt.Fprintf(&builder, "    %s: %s\n", h.Name, h.Value)
			}
		}
		if redacted.UserAgent != "" {
			fmt.Fprintf(&builder, "  User-Agent: %s\n", redacted.UserAgent)
		}
		if redacted.Username != "" || redacted.Password != "" {
			authDisplay := redacted.Username
			if redacted.Password != "" {
				authDisplay += ":[REDACTED]"
			}
			fmt.Fprintf(&builder, "  Auth: %s\n", authDisplay)
		}
	}
	fmt.Fprintln(&builder, "Mode: dry-run only, no network request performed")

	return strings.TrimRight(builder.String(), "\n")
}

func valueOrDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
