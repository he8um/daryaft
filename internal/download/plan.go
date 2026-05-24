package download

import (
	"fmt"
	"strings"
)

type Plan struct {
	URLs    []string
	Output  string
	Name    string
	Retries int
	Resume  bool
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
	fmt.Fprintf(&builder, "Retries: %d\n", p.Retries)
	fmt.Fprintf(&builder, "Resume: %t\n", p.Resume)
	fmt.Fprintln(&builder, "Mode: dry-run only, no network request performed")

	return strings.TrimRight(builder.String(), "\n")
}

func valueOrDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
