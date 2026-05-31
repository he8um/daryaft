package config

import (
	"os"
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

	return parseNonNegativeInt(name, value)
}

func envBool(lookup LookupEnvFunc, name string, fallback bool) (bool, error) {
	value, ok := lookup(name)
	if !ok {
		return fallback, nil
	}

	return parseBoolValue(name, value)
}
