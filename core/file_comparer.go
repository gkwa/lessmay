package core

import (
	"bytes"
	"fmt"
	"os"
)

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
