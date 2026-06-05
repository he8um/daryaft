package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const configFileName = "config.yaml"

var userConfigDir = os.UserConfigDir
var userHomeDir = os.UserHomeDir

func Path() (string, error) {
	dir, err := userConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "daryaft", configFileName), nil
}

func Exists() (bool, error) {
	path, err := Path()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SetUserConfigDirForTest(dir string) func() {
	previous := userConfigDir
	userConfigDir = func() (string, error) {
		return dir, nil
	}
	return func() {
		userConfigDir = previous
	}
}

func BuiltinDownloadDir() string {
	home, err := userHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return "."
	}
	return filepath.Join(home, "Downloads")
}

func EffectiveDownloadDir(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return BuiltinDownloadDir()
	}
	return trimmed
}

func SetUserHomeDirForTest(dir string) func() {
	return SetUserHomeDirFuncForTest(func() (string, error) {
		return dir, nil
	})
}

func SetUserHomeDirErrorForTest() func() {
	return SetUserHomeDirFuncForTest(func() (string, error) {
		return "", errors.New("home unavailable")
	})
}

func SetUserHomeDirFuncForTest(fn func() (string, error)) func() {
	previous := userHomeDir
	userHomeDir = fn
	return func() {
		userHomeDir = previous
	}
}
