package utils

import "testing"

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{name: "bytes", bytes: 512, want: "512 B"},
		{name: "kilobytes", bytes: 1536, want: "1.5 KB"},
		{name: "megabytes", bytes: 2 * 1024 * 1024, want: "2.0 MB"},
		{name: "negative", bytes: -1, want: "0 B"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FormatBytes(test.bytes); got != test.want {
				t.Fatalf("FormatBytes(%d) = %q, want %q", test.bytes, got, test.want)
			}
		})
	}
}

func TestFormatSpeed(t *testing.T) {
	if got := FormatSpeed(1536); got != "1.5 KB/s" {
		t.Fatalf("FormatSpeed(1536) = %q", got)
	}
}
