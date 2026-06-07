package version

import "runtime"

var (
	Version = "1.4.0-dev"
	Commit  = "local"
	Date    = "unknown"
	BuiltBy = "source"
)

type Details struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	BuiltBy   string `json:"built_by"`
	GoVersion string `json:"go_version"`
}

func Info() Details {
	return Details{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		BuiltBy:   BuiltBy,
		GoVersion: runtime.Version(),
	}
}
