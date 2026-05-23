package version

import "runtime"

var (
	Version = "0.1.0-dev"
	Commit  = "local"
	Date    = "unknown"
)

type Details struct {
	Version   string
	Commit    string
	Date      string
	GoVersion string
}

func Info() Details {
	return Details{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
	}
}
