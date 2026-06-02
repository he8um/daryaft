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

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	dir := filepath.Dir(path)
	temp, err := os.CreateTemp(dir, ".config-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary config file: %w", err)
	}
	tempPath := temp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tempPath)
		}
	}()

	if err := temp.Chmod(0o600); err != nil {
		_ = temp.Close()
		return fmt.Errorf("set config permissions: %w", err)
	}
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write temporary config %q: %w", tempPath, err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary config %q: %w", tempPath, err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("write config %q: %w", path, err)
	}
	cleanup = false
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
