package core

type FileFinder interface {
	FindSyncConflictFiles(paths, skipPaths []string) ([]string, error)
}

type DiffRunner interface {
	RunDiff(conflictFile, originalFile string, count int) error
}

type FileComparer interface {
	CompareAndDelete(conflictFile, originalFile string) (bool, error)
}
