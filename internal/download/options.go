package download

type Options struct {
	URLs     []string
	File     string
	Output   string
	Name     string
	DryRun   bool
	Retries  int
	Resume   bool
	NoResume bool
}
