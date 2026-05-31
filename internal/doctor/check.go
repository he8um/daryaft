package doctor

import (
	"os"
	"path/filepath"
)

type Status int

const (
	StatusOK Status = iota
	StatusFail
	StatusWarn
	StatusInfo
)

type Check struct {
	Section string
	Status  Status
	Label   string
	Value   string
}

type Report struct {
	Checks []Check
}

func (r *Report) Add(section string, status Status, label string, value string) {
	r.Checks = append(r.Checks, Check{
		Section: section,
		Status:  status,
		Label:   label,
		Value:   value,
	})
}

func (r Report) Failed() bool {
	for _, check := range r.Checks {
		if check.Status == StatusFail {
			return true
		}
	}
	return false
}

func writableDir(path string) error {
	file, err := os.CreateTemp(path, ".daryaft-doctor-*")
	if err != nil {
		return err
	}
	name := file.Name()
	closeErr := file.Close()
	removeErr := os.Remove(name)
	if closeErr != nil {
		return closeErr
	}
	return removeErr
}

func existingAncestor(path string, stat func(string) (os.FileInfo, error)) (string, error) {
	current := filepath.Clean(path)
	for {
		info, err := stat(current)
		if err == nil {
			if info.IsDir() {
				return current, nil
			}
			return "", errNotDirectory(current)
		}
		if !os.IsNotExist(err) {
			return "", err
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", err
		}
		current = parent
	}
}

type errNotDirectory string

func (e errNotDirectory) Error() string {
	return string(e) + " is not a directory"
}
