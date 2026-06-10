package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/he8um/daryaft/internal/download"
)

type KeyInfo struct {
	Key  string
	Type string
}

const (
	ThemeDefault = "default"
	ThemeMono    = "mono"

	keyDownloadDir = "download_dir"
	keyRetries     = "retries"
	keyResume      = "resume"
	keyNoColor     = "no_color"
	keyNoTUI       = "no_tui"
	keyTheme       = "theme"
	keyAnimations  = "animations"
	keyHyperlinks  = "hyperlinks"
	keyUserAgent   = "user_agent"
	keyTimeout     = "timeout"
)

func SupportedKeys() []KeyInfo {
	return []KeyInfo{
		{Key: keyDownloadDir, Type: "string"},
		{Key: keyRetries, Type: "int"},
		{Key: keyResume, Type: "bool"},
		{Key: keyNoColor, Type: "bool"},
		{Key: keyNoTUI, Type: "bool"},
		{Key: keyTheme, Type: "string"},
		{Key: keyAnimations, Type: "bool"},
		{Key: keyHyperlinks, Type: "bool"},
		{Key: keyUserAgent, Type: "string"},
		{Key: keyTimeout, Type: "string"},
	}
}

func Get(cfg Config, key string) (string, error) {
	switch key {
	case keyDownloadDir:
		return cfg.DownloadDir, nil
	case keyRetries:
		return strconv.Itoa(cfg.Retries), nil
	case keyResume:
		return strconv.FormatBool(cfg.Resume), nil
	case keyNoColor:
		return strconv.FormatBool(cfg.NoColor), nil
	case keyNoTUI:
		return strconv.FormatBool(cfg.NoTUI), nil
	case keyTheme:
		return cfg.Theme, nil
	case keyAnimations:
		return strconv.FormatBool(cfg.Animations), nil
	case keyHyperlinks:
		return strconv.FormatBool(cfg.Hyperlinks), nil
	case keyUserAgent:
		return cfg.UserAgent, nil
	case keyTimeout:
		return cfg.Timeout, nil
	default:
		return "", unknownKeyError(key)
	}
}

func Set(cfg Config, key string, value string) (Config, error) {
	switch key {
	case keyDownloadDir:
		cfg.DownloadDir = strings.TrimSpace(value)
	case keyRetries:
		retries, err := parseNonNegativeInt(key, value)
		if err != nil {
			return Config{}, err
		}
		cfg.Retries = retries
	case keyResume:
		resume, err := parseBoolValue(key, value)
		if err != nil {
			return Config{}, err
		}
		cfg.Resume = resume
	case keyNoColor:
		noColor, err := parseBoolValue(key, value)
		if err != nil {
			return Config{}, err
		}
		cfg.NoColor = noColor
	case keyNoTUI:
		noTUI, err := parseBoolValue(key, value)
		if err != nil {
			return Config{}, err
		}
		cfg.NoTUI = noTUI
	case keyTheme:
		theme, err := NormalizeTheme(value)
		if err != nil {
			return Config{}, err
		}
		cfg.Theme = theme
	case keyAnimations:
		animations, err := parseBoolValue(key, value)
		if err != nil {
			return Config{}, err
		}
		cfg.Animations = animations
	case keyHyperlinks:
		hyperlinks, err := parseBoolValue(key, value)
		if err != nil {
			return Config{}, err
		}
		cfg.Hyperlinks = hyperlinks
	case keyUserAgent:
		trimmed := strings.TrimSpace(value)
		if err := validateUserAgent(trimmed); err != nil {
			return Config{}, fmt.Errorf("invalid %s %q: %w", key, value, err)
		}
		cfg.UserAgent = trimmed
	case keyTimeout:
		trimmed := strings.TrimSpace(value)
		if err := validateTimeout(trimmed); err != nil {
			return Config{}, fmt.Errorf("invalid %s %q: %w", key, value, err)
		}
		cfg.Timeout = trimmed
	default:
		return Config{}, unknownKeyError(key)
	}
	return cfg, nil
}

func NormalizeTheme(value string) (string, error) {
	theme := strings.ToLower(strings.TrimSpace(value))
	switch theme {
	case ThemeDefault, ThemeMono:
		return theme, nil
	default:
		return "", fmt.Errorf("invalid theme %q: must be one of default, mono", value)
	}
}

func IsMonoTheme(value string) bool {
	theme, err := NormalizeTheme(value)
	return err == nil && theme == ThemeMono
}

func Reset() (string, error) {
	path, err := Path()
	if err != nil {
		return "", fmt.Errorf("resolve config path: %w", err)
	}
	if err := Save(Default()); err != nil {
		return "", err
	}
	return path, nil
}

func parseNonNegativeInt(name, value string) (int, error) {
	trimmed := strings.TrimSpace(value)
	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: must be an integer between 0 and %d", name, value, download.MaxRetries)
	}
	if err := download.ValidateRetries(parsed); err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", name, value, err)
	}
	return parsed, nil
}

func parseBoolValue(name, value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "y", "on":
		return true, nil
	case "false", "0", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid %s %q: must be one of true, false, 1, 0, yes, no, y, n, on, off", name, value)
	}
}

func unknownKeyError(key string) error {
	return fmt.Errorf("unknown config key %q", key)
}
