package config

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
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
	envUserAgent   = "DARYAFT_USER_AGENT"
	envTimeout     = "DARYAFT_TIMEOUT"
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
		theme, err := NormalizeTheme(value)
		if err != nil {
			return Config{}, fmt.Errorf("invalid %s %q: %w", envTheme, value, err)
		}
		cfg.Theme = theme
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

	if value, ok := lookup(envUserAgent); ok {
		trimmed := strings.TrimSpace(value)
		if err := validateUserAgent(trimmed); err != nil {
			return Config{}, fmt.Errorf("invalid %s %q: %w", envUserAgent, value, err)
		}
		cfg.UserAgent = trimmed
	}
	if value, ok := lookup(envTimeout); ok {
		if err := validateTimeout(value); err != nil {
			return Config{}, fmt.Errorf("invalid %s %q: %w", envTimeout, value, err)
		}
		cfg.Timeout = strings.TrimSpace(value)
	}

	return cfg, nil
}

func validateUserAgent(ua string) error {
	for _, r := range ua {
		if unicode.IsControl(r) {
			return fmt.Errorf("user-agent contains control character")
		}
	}
	return nil
}

func validateTimeout(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	d, err := time.ParseDuration(trimmed)
	if err != nil {
		return fmt.Errorf("timeout must be a positive duration (for example 30s, 2m)")
	}
	if d <= 0 {
		return fmt.Errorf("timeout must be a positive duration (for example 30s, 2m)")
	}
	return nil
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
