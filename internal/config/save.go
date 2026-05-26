package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return fmt.Errorf("resolve config path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config %q: %w", path, err)
	}
	return nil
}

func Init(overwrite bool) (string, error) {
	path, err := Path()
	if err != nil {
		return "", fmt.Errorf("resolve config path: %w", err)
	}

	exists, err := Exists()
	if err != nil {
		return "", fmt.Errorf("check config %q: %w", path, err)
	}
	if exists && !overwrite {
		return path, fmt.Errorf("config already exists: %s", path)
	}

	if err := Save(Default()); err != nil {
		return "", err
	}
	return path, nil
}
