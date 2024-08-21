package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	homedir "github.com/mitchellh/go-homedir"
)

func ShowConflicts(
	logger logr.Logger,
	args []string,
	defaultObsidianPath string,
	skipPaths []string,
) error {
	paths, err := getConflictPaths(args, defaultObsidianPath)
	if err != nil {
		return fmt.Errorf("failed to get conflict paths: %w", err)
	}

	resolver := NewSyncConflictResolver(logger)
	return resolver.ResolveSyncConflicts(paths, skipPaths)
}

func getConflictPaths(
	args []string,
	defaultObsidianPath string,
) ([]string, error) {
	var paths []string
	if len(args) == 0 {
		paths = append(paths, defaultObsidianPath)
	} else {
		for _, arg := range args {
			expandedPath, err := homedir.Expand(arg)
			if err != nil {
				return nil, fmt.Errorf("failed to expand path %s: %w", arg, err)
			}
			paths = append(paths, expandedPath)
		}
	}
	return paths, nil
}

func GetDefaultObsidianPath() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return filepath.Join(home, "Documents", "Obsidian Vault")
}
