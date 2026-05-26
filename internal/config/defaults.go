package config

const (
	AppName    = "Daryaft"
	BinaryName = "daryaft"
	WebsiteURL = "https://xhesam.com"
	ProjectURL = "https://xhesam.com/daryaft"
	GitHubRepo = "https://github.com/he8um/daryaft"
	FooterText = "Developed with <3 by AmirHesam Piri"
)

func Default() Config {
	return Config{
		DownloadDir: "",
		Retries:     3,
		Resume:      true,
		NoColor:     false,
		NoTUI:       false,
		Theme:       "default",
		Animations:  true,
		Hyperlinks:  true,
	}
}
