package config

type Config struct {
	DownloadDir string `yaml:"download_dir"`
	Retries     int    `yaml:"retries"`
	Resume      bool   `yaml:"resume"`
	NoColor     bool   `yaml:"no_color"`
	NoTUI       bool   `yaml:"no_tui"`
	Theme       string `yaml:"theme"`
	Animations  bool   `yaml:"animations"`
	Hyperlinks  bool   `yaml:"hyperlinks"`
	UserAgent   string `yaml:"user_agent"`
	Timeout     string `yaml:"timeout"`
}
