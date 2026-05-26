package config

import (
	"os"
	"path/filepath"
)

const configFileName = "config.yaml"

var userConfigDir = os.UserConfigDir

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
