package core

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	homedir "github.com/mitchellh/go-homedir"
)

type FileFinder interface {
	FindSyncConflictFiles(paths, skipPaths []string) ([]string, error)
}

type DiffRunner interface {
	RunDiff(conflictFile, originalFile string, count int) error
}

type FileComparer interface {
	CompareAndDelete(conflictFile, originalFile string) (bool, error)
}

type SyncConflictResolver struct {
	finder   FileFinder
	differ   DiffRunner
	comparer FileComparer
	logger   logr.Logger
}

func NewSyncConflictResolver(logger logr.Logger) *SyncConflictResolver {
	return &SyncConflictResolver{
		finder:   &DefaultFileFinder{},
		differ:   &DefaultDiffRunner{},
		comparer: &DefaultFileComparer{},
		logger:   logger,
	}
}

func (r *SyncConflictResolver) ResolveSyncConflicts(paths, skipPaths []string) error {
	r.logger.V(1).Info("Starting sync conflict resolution")

	conflictFiles, err := r.finder.FindSyncConflictFiles(paths, skipPaths)
	if err != nil {
		return fmt.Errorf("failed to find sync conflict files: %w", err)
	}

	for i, conflictFile := range conflictFiles {
		originalFile := regexp.MustCompile(`\.sync-conflict-\d{8}-\d{6}-\w+\.md$`).
			ReplaceAllString(conflictFile, ".md")

		deleted, err := r.comparer.CompareAndDelete(conflictFile, originalFile)
		if err != nil {
			r.logger.Error(err, "Failed to compare and delete files", "conflictFile", conflictFile, "originalFile", originalFile)
			continue
		}

		if deleted {
			r.logger.Info("Deleted identical sync conflict file", "conflictFile", conflictFile)
		} else {
			if err := r.differ.RunDiff(conflictFile, originalFile, i+1); err != nil {
				r.logger.Error(err, "Failed to run diff", "conflictFile", conflictFile, "originalFile", originalFile)
			}
		}
		fmt.Println()
	}

	r.logger.V(1).Info("Finished sync conflict resolution")
	return nil
}

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

type DefaultDiffRunner struct{}

func (d *DefaultDiffRunner) RunDiff(conflictFile, originalFile string, count int) error {
	formattedCmd := fmt.Sprintf("diff --unified --ignore-all-space '%s' '%s'",
		strings.ReplaceAll(conflictFile, "'", "'\"'\"'"),
		strings.ReplaceAll(originalFile, "'", "'\"'\"'"))

	absConflictFile, _ := filepath.Abs(conflictFile)
	absOriginalFile, _ := filepath.Abs(originalFile)

	fmt.Printf("# diff: %d\n", count)
	fmt.Printf("%s\n", formattedCmd)
	fmt.Printf("%s\n", absConflictFile)
	fmt.Printf("%s\n", absOriginalFile)
	fmt.Printf("open obsidian:'//open?path=%s'; ", absOriginalFile)
	fmt.Printf("open obsidian:'//open?path=%s'\n", absConflictFile)

	return nil
}

type DefaultFileComparer struct{}

func (c *DefaultFileComparer) CompareAndDelete(conflictFile, originalFile string) (bool, error) {
	conflictContent, err := os.ReadFile(conflictFile)
	if err != nil {
		return false, fmt.Errorf("error reading conflict file: %w", err)
	}

	originalContent, err := os.ReadFile(originalFile)
	if err != nil {
		return false, fmt.Errorf("error reading original file: %w", err)
	}

	if bytes.Equal(conflictContent, originalContent) {
		err := os.Remove(conflictFile)
		if err != nil {
			return false, fmt.Errorf("error deleting conflict file: %w", err)
		}
		return true, nil
	}

	return false, nil
}

func ShowConflicts(logger logr.Logger, args []string, defaultObsidianPath string, skipPaths []string) error {
	paths, err := getConflictPaths(args, defaultObsidianPath)
	if err != nil {
		return fmt.Errorf("failed to get conflict paths: %w", err)
	}

	resolver := NewSyncConflictResolver(logger)
	return resolver.ResolveSyncConflicts(paths, skipPaths)
}

func getConflictPaths(args []string, defaultObsidianPath string) ([]string, error) {
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
