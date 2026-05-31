package doctor

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/pkg/version"
)

func TestDoctorSucceedsWithValidDefaultEnvironment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "daryaft", "config.yaml")

	report := Run(testOptions(dir, path, appconfig.Default()))
	output := Format(report)

	if report.Failed() {
		t.Fatalf("report failed:\n%s", output)
	}
	for _, want := range []string{
		"Daryaft doctor",
		"✓ Config load: ok",
		"✓ Default output: current directory",
		"✓ Output writable: yes",
		"✓ Version: test-version",
		path,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestDoctorInvalidConfigCausesFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "daryaft", "config.yaml")

	options := testOptions(dir, path, appconfig.Default())
	options.LoadConfig = func() (appconfig.Config, error) {
		return appconfig.Config{}, errors.New("parse config: invalid yaml")
	}

	report := Run(options)
	output := Format(report)

	if !report.Failed() {
		t.Fatalf("report did not fail:\n%s", output)
	}
	if !strings.Contains(output, "✗ Config load: parse config: invalid yaml") {
		t.Fatalf("output missing config failure:\n%s", output)
	}
}

func TestDoctorMissingDownloadDirIsWarning(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "daryaft", "config.yaml")
	cfg := appconfig.Default()
	cfg.DownloadDir = filepath.Join(dir, "missing-downloads")

	report := Run(testOptions(dir, path, cfg))
	output := Format(report)

	if report.Failed() {
		t.Fatalf("report failed:\n%s", output)
	}
	if !strings.Contains(output, "! Output writable: directory missing") {
		t.Fatalf("output missing warning:\n%s", output)
	}
}

func TestDoctorExistingUnwritableDownloadDirCausesFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "daryaft", "config.yaml")
	downloads := filepath.Join(dir, "downloads")
	if err := os.Mkdir(downloads, 0o755); err != nil {
		t.Fatalf("mkdir downloads: %v", err)
	}
	cfg := appconfig.Default()
	cfg.DownloadDir = downloads

	options := testOptions(dir, path, cfg)
	options.CheckWritable = func(path string) error {
		if path == downloads {
			return errors.New("permission denied")
		}
		return nil
	}

	report := Run(options)
	output := Format(report)

	if !report.Failed() {
		t.Fatalf("report did not fail:\n%s", output)
	}
	if !strings.Contains(output, "✗ Output writable: no: permission denied") {
		t.Fatalf("output missing writable failure:\n%s", output)
	}
}

func TestDoctorClamscanDetectionReportsOptionalStatus(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "daryaft", "config.yaml")

	foundOptions := testOptions(dir, path, appconfig.Default())
	foundOptions.LookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}
	foundOutput := Format(Run(foundOptions))
	if !strings.Contains(foundOutput, "- clamscan: found: /usr/local/bin/clamscan") {
		t.Fatalf("output missing found clamscan:\n%s", foundOutput)
	}

	missingOptions := testOptions(dir, path, appconfig.Default())
	missingOptions.LookPath = func(name string) (string, error) {
		return "", exec.ErrNotFound
	}
	missingOutput := Format(Run(missingOptions))
	if !strings.Contains(missingOutput, "- clamscan: not found") {
		t.Fatalf("output missing missing clamscan:\n%s", missingOutput)
	}
}

func testOptions(dir string, configPath string, cfg appconfig.Config) Options {
	return Options{
		ConfigPath: func() (string, error) {
			return configPath, nil
		},
		LoadConfig: func() (appconfig.Config, error) {
			return cfg, nil
		},
		LookupEnv: func(name string) (string, bool) {
			switch name {
			case "TERM":
				return "xterm-256color", true
			default:
				return "", false
			}
		},
		LookPath: func(name string) (string, error) {
			return "", exec.ErrNotFound
		},
		Getwd: func() (string, error) {
			return dir, nil
		},
		StdoutStat: func() (os.FileInfo, error) {
			return fakeFileInfo{mode: 0}, nil
		},
		Version: version.Details{
			Version:   "test-version",
			Commit:    "test-commit",
			Date:      "test-date",
			GoVersion: "test-go",
		},
	}
}

type fakeFileInfo struct {
	mode os.FileMode
}

func (f fakeFileInfo) Name() string       { return "stdout" }
func (f fakeFileInfo) Size() int64        { return 0 }
func (f fakeFileInfo) Mode() os.FileMode  { return f.mode }
func (f fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fakeFileInfo) IsDir() bool        { return false }
func (f fakeFileInfo) Sys() any           { return nil }
