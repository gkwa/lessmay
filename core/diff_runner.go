package core

import (
	"fmt"
	"path/filepath"
	"strings"
)

type DefaultDiffRunner struct{}

func (d *DefaultDiffRunner) RunDiff(
	conflictFile, originalFile string,
	count int,
) error {
	diffCmd := []string{
		"diff",
		"--unified",
		"--ignore-all-space",
		strings.ReplaceAll(conflictFile, "'", "'\"'\"'"),
		strings.ReplaceAll(originalFile, "'", "'\"'\"'"),
	}

	formattedCmd := strings.Join(diffCmd, " ")

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
