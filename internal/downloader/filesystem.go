package downloader

import (
	"fmt"
	"os"
	"path/filepath"
)

type targetPaths struct {
	OutputDir string
	Final     string
	Partial   string
}

func prepareTarget(outputDir, filename string) (targetPaths, error) {
	target, err := targetPathsFor(outputDir, filename)
	if err != nil {
		return targetPaths{}, err
	}

	if _, err := os.Stat(target.Final); err == nil {
		return targetPaths{}, fmt.Errorf("%w: %s", ErrTargetExists, target.Final)
	} else if !os.IsNotExist(err) {
		return targetPaths{}, fmt.Errorf("check target file %q: %w", target.Final, err)
	}

	if err := os.MkdirAll(target.OutputDir, 0o755); err != nil {
		return targetPaths{}, fmt.Errorf("create output directory %q: %w", target.OutputDir, err)
	}

	return target, nil
}

func targetPathsFor(outputDir, filename string) (targetPaths, error) {
	if outputDir == "" {
		outputDir = "."
	}

	outputDir = filepath.Clean(outputDir)
	filename = sanitizeFilename(filename)

	finalPath := filepath.Clean(filepath.Join(outputDir, filename))
	partialPath := finalPath + ".part"

	if err := ensureInsideOutputDir(outputDir, finalPath); err != nil {
		return targetPaths{}, err
	}

	return targetPaths{
		OutputDir: outputDir,
		Final:     finalPath,
		Partial:   partialPath,
	}, nil
}

func ensureInsideOutputDir(outputDir, target string) error {
	rel, err := filepath.Rel(outputDir, target)
	if err != nil {
		return fmt.Errorf("validate target path %q: %w", target, err)
	}
	if rel == ".." || len(rel) >= 3 && rel[:3] == "../" {
		return fmt.Errorf("target path escapes output directory: %s", target)
	}
	return nil
}
