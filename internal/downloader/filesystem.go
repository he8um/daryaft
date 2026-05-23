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

	if _, err := os.Stat(finalPath); err == nil {
		return targetPaths{}, fmt.Errorf("%w: %s", ErrTargetExists, finalPath)
	} else if !os.IsNotExist(err) {
		return targetPaths{}, fmt.Errorf("check target file %q: %w", finalPath, err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return targetPaths{}, fmt.Errorf("create output directory %q: %w", outputDir, err)
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
