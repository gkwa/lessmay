package core

import (
	"fmt"
	"regexp"

	"github.com/go-logr/logr"
)

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

func (r *SyncConflictResolver) ResolveSyncConflicts(
	paths, skipPaths []string,
) error {
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
			r.logger.Error(
				err,
				"Failed to compare and delete files",
				"conflictFile",
				conflictFile,
				"originalFile",
				originalFile,
			)
			continue
		}

		if deleted {
			r.logger.Info(
				"Deleted identical sync conflict file",
				"conflictFile",
				conflictFile,
			)
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
