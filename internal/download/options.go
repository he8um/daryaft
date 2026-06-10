package download

import "github.com/he8um/daryaft/internal/httpopts"

type Options struct {
	URLs         []string
	File         string
	Output       string
	Name         string
	DryRun       bool
	Checksum     string
	ChecksumFile string
	Retries      int
	Resume       bool
	NoResume     bool
	HTTPOptions  httpopts.Options
}
