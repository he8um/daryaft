package doctor

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/pkg/version"
)

type Options struct {
	ConfigPath    func() (string, error)
	LoadConfig    func() (appconfig.Config, error)
	LookupEnv     func(string) (string, bool)
	LookPath      func(string) (string, error)
	Getwd         func() (string, error)
	Stat          func(string) (os.FileInfo, error)
	CheckWritable func(string) error
	StdoutStat    func() (os.FileInfo, error)
	Version       version.Details
}

func Run(options Options) Report {
	options = applyDefaults(options)

	var report Report
	addSystemChecks(&report)
	addVersionChecks(&report, options.Version)

	cfgPath, cfgPathOK := addConfigPathChecks(&report, options)
	cfg, cfgOK := addConfigLoadCheck(&report, options)
	if cfgPathOK && cfgPath != "" {
		addConfigDirCheck(&report, options, filepath.Dir(cfgPath))
	}
	if cfgOK {
		addDownloadChecks(&report, options, cfg)
	}

	addTerminalChecks(&report, options)
	addOptionalToolChecks(&report, options)
	report.Add("Skipped", StatusSkipped, "GitHub release check", "skipped")

	return report
}

func applyDefaults(options Options) Options {
	if options.ConfigPath == nil {
		options.ConfigPath = appconfig.Path
	}
	if options.LoadConfig == nil {
		options.LoadConfig = appconfig.LoadEffective
	}
	if options.LookupEnv == nil {
		options.LookupEnv = os.LookupEnv
	}
	if options.LookPath == nil {
		options.LookPath = exec.LookPath
	}
	if options.Getwd == nil {
		options.Getwd = os.Getwd
	}
	if options.Stat == nil {
		options.Stat = os.Stat
	}
	if options.CheckWritable == nil {
		options.CheckWritable = writableDir
	}
	if options.StdoutStat == nil {
		options.StdoutStat = os.Stdout.Stat
	}
	if options.Version.Version == "" {
		options.Version = version.Info()
	}
	return options
}

func addSystemChecks(report *Report) {
	report.Add("System", StatusOK, "OS", runtime.GOOS)
	report.Add("System", StatusOK, "Arch", runtime.GOARCH)
	report.Add("System", StatusOK, "Go", runtime.Version())
}

func addVersionChecks(report *Report, info version.Details) {
	report.Add("Version", StatusOK, "Version", info.Version)
	report.Add("Version", StatusOK, "Commit", info.Commit)
	report.Add("Version", StatusOK, "Build date", info.Date)
}

func addConfigPathChecks(report *Report, options Options) (string, bool) {
	path, err := options.ConfigPath()
	if err != nil {
		report.Add("Config", StatusFail, "Path", err.Error())
		return "", false
	}
	report.Add("Config", StatusOK, "Path", path)
	return path, true
}

func addConfigDirCheck(report *Report, options Options, dir string) {
	info, err := options.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			report.Add("Config", StatusWarn, "Config directory", "path exists but is not a directory")
			return
		}
		if err := options.CheckWritable(dir); err != nil {
			report.Add("Config", StatusWarn, "Config directory", "not writable: "+err.Error())
			return
		}
		report.Add("Config", StatusOK, "Config directory", "writable")
		return
	}

	if !os.IsNotExist(err) {
		report.Add("Config", StatusWarn, "Config directory", "cannot inspect: "+err.Error())
		return
	}

	ancestor, err := existingAncestor(dir, options.Stat)
	if err != nil {
		report.Add("Config", StatusWarn, "Config directory", "cannot find creatable parent: "+err.Error())
		return
	}
	if err := options.CheckWritable(ancestor); err != nil {
		report.Add("Config", StatusWarn, "Config directory", "missing; parent not writable: "+err.Error())
		return
	}
	report.Add("Config", StatusOK, "Config directory", "can create")
}

func addConfigLoadCheck(report *Report, options Options) (appconfig.Config, bool) {
	cfg, err := options.LoadConfig()
	if err != nil {
		report.Add("Config", StatusFail, "Config load", err.Error())
		return appconfig.Config{}, false
	}
	report.Add("Config", StatusOK, "Config load", "ok")
	return cfg, true
}

func addDownloadChecks(report *Report, options Options, cfg appconfig.Config) {
	target := cfg.DownloadDir
	if target == "" {
		wd, err := options.Getwd()
		if err != nil {
			report.Add("Download", StatusFail, "Default output", "current directory unavailable: "+err.Error())
			return
		}
		target = wd
		report.Add("Download", StatusOK, "Default output", "current directory")
	} else {
		report.Add("Download", StatusOK, "Default output", target)
	}

	info, err := options.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			report.Add("Download", StatusWarn, "Output writable", "directory missing")
			return
		}
		report.Add("Download", StatusFail, "Output writable", "cannot inspect: "+err.Error())
		return
	}
	if !info.IsDir() {
		report.Add("Download", StatusFail, "Output writable", "path exists but is not a directory")
		return
	}
	if err := options.CheckWritable(target); err != nil {
		report.Add("Download", StatusFail, "Output writable", "no: "+err.Error())
		return
	}
	report.Add("Download", StatusOK, "Output writable", "yes")
}

func addTerminalChecks(report *Report, options Options) {
	if term, ok := options.LookupEnv("TERM"); ok && term != "" {
		report.Add("Terminal", StatusInfo, "TERM", term)
	} else {
		report.Add("Terminal", StatusInfo, "TERM", "not set")
	}

	if _, ok := options.LookupEnv("NO_COLOR"); ok {
		report.Add("Terminal", StatusInfo, "NO_COLOR", "set")
	} else {
		report.Add("Terminal", StatusInfo, "NO_COLOR", "not set")
	}

	info, err := options.StdoutStat()
	if err != nil {
		report.Add("Terminal", StatusInfo, "Stdout terminal", "unknown: "+err.Error())
		return
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		report.Add("Terminal", StatusInfo, "Stdout terminal", "yes")
		return
	}
	report.Add("Terminal", StatusInfo, "Stdout terminal", "no")
}

func addOptionalToolChecks(report *Report, options Options) {
	path, err := options.LookPath("clamscan")
	if err == nil {
		report.Add("Optional tools", StatusInfo, "clamscan", "found: "+path)
		return
	}
	if errors.Is(err, exec.ErrNotFound) {
		report.Add("Optional tools", StatusInfo, "clamscan", "not found")
		return
	}
	report.Add("Optional tools", StatusInfo, "clamscan", "not found: "+err.Error())
}
