package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LookupEnvFunc func(string) (string, bool)

const (
	envDownloadDir = "DARYAFT_DOWNLOAD_DIR"
	envRetries     = "DARYAFT_RETRIES"
	envResume      = "DARYAFT_RESUME"
	envNoColor     = "DARYAFT_NO_COLOR"
	envNoTUI       = "DARYAFT_NO_TUI"
	envTheme       = "DARYAFT_THEME"
	envAnimations  = "DARYAFT_ANIMATIONS"
	envHyperlinks  = "DARYAFT_HYPERLINKS"
)

func LoadEffective() (Config, error) {
	cfg, err := Load()
	if err != nil {
		return Config{}, err
	}
	return ApplyEnv(cfg, os.LookupEnv)
}

func ApplyEnv(cfg Config, lookup LookupEnvFunc) (Config, error) {
	if lookup == nil {
		lookup = os.LookupEnv
	}

	if value, ok := lookup(envDownloadDir); ok {
		cfg.DownloadDir = strings.TrimSpace(value)
	}
	if value, ok := lookup(envTheme); ok {
		cfg.Theme = strings.TrimSpace(value)
	}

	var err error
	if cfg.Retries, err = envInt(lookup, envRetries, cfg.Retries); err != nil {
		return Config{}, err
	}
	if cfg.Resume, err = envBool(lookup, envResume, cfg.Resume); err != nil {
		return Config{}, err
	}
	if cfg.NoColor, err = envBool(lookup, envNoColor, cfg.NoColor); err != nil {
		return Config{}, err
	}
	if cfg.NoTUI, err = envBool(lookup, envNoTUI, cfg.NoTUI); err != nil {
		return Config{}, err
	}
	if cfg.Animations, err = envBool(lookup, envAnimations, cfg.Animations); err != nil {
		return Config{}, err
	}
	if cfg.Hyperlinks, err = envBool(lookup, envHyperlinks, cfg.Hyperlinks); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func envInt(lookup LookupEnvFunc, name string, fallback int) (int, error) {
	value, ok := lookup(name)
	if !ok {
		return fallback, nil
	}

	trimmed := strings.TrimSpace(value)
	parsed, err := strconv.Atoi(trimmed)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("invalid %s %q: must be an integer greater than or equal to 0", name, value)
	}
	return parsed, nil
}

func envBool(lookup LookupEnvFunc, name string, fallback bool) (bool, error) {
	value, ok := lookup(name)
	if !ok {
		return fallback, nil
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "y", "on":
		return true, nil
	case "false", "0", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid %s %q: must be one of true, false, 1, 0, yes, no, y, n, on, off", name, value)
	}
}
