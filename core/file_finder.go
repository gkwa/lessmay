package core

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

type DefaultFileFinder struct{}

func (f *DefaultFileFinder) FindSyncConflictFiles(paths, skipPaths []string) ([]string, error) {
	var conflictFiles []string
	pattern := regexp.MustCompile(`\.sync-conflict-\d{8}-\d{6}-\w+\.md$`)

	for _, path := range paths {
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && pattern.MatchString(d.Name()) {
				if !shouldSkip(path, skipPaths) {
					conflictFiles = append(conflictFiles, path)
				}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", path, err)
		}
	}

	return conflictFiles, nil
}

func shouldSkip(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.Contains(path, skipPath) {
			return true
		}
	}
	return false
}
