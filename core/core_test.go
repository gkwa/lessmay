package core

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-logr/logr/testr"
)

type mockFileFinder struct {
	files []string
	err   error
}

func (m *mockFileFinder) FindSyncConflictFiles(paths, skipPaths []string) ([]string, error) {
	return m.files, m.err
}

type mockDiffRunner struct {
	calls []struct {
		conflictFile string
		originalFile string
		count        int
	}
	err error
}

func (m *mockDiffRunner) RunDiff(conflictFile, originalFile string, count int) error {
	m.calls = append(m.calls, struct {
		conflictFile string
		originalFile string
		count        int
	}{conflictFile, originalFile, count})
	return m.err
}

func TestSyncConflictResolver_ResolveSyncConflicts(t *testing.T) {
	tests := []struct {
		name          string
		files         []string
		skipPaths     []string
		finderErr     error
		differErr     error
		expectedCalls int
		expectErr     bool
	}{
		{
			name:          "successful resolution",
			files:         []string{"/path/to/file.sync-conflict-20240818-215425-I2NUVZU.md"},
			skipPaths:     []string{".trash"},
			expectedCalls: 1,
		},
		{
			name:      "finder error",
			finderErr: errors.New("finder error"),
			skipPaths: []string{".trash"},
			expectErr: true,
		},
		{
			name:          "differ error",
			files:         []string{"/path/to/file.sync-conflict-20240818-215425-I2NUVZU.md"},
			skipPaths:     []string{".trash"},
			differErr:     errors.New("differ error"),
			expectedCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := &mockFileFinder{files: tt.files, err: tt.finderErr}
			differ := &mockDiffRunner{err: tt.differErr}
			logger := testr.New(t)

			resolver := &SyncConflictResolver{
				finder:   finder,
				differ:   differ,
				comparer: &DefaultFileComparer{},
				logger:   logger,
			}

			err := resolver.ResolveSyncConflicts([]string{"/test/path"}, tt.skipPaths)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if len(differ.calls) != tt.expectedCalls {
				t.Errorf("Expected %d diff calls, but got %d", tt.expectedCalls, len(differ.calls))
			}
		})
	}
}

func TestDefaultFileFinder_FindSyncConflictFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sync_conflict_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	files := []string{
		"file1.sync-conflict-20240818-215425-I2NUVZU.md",
		"subdir/file2.sync-conflict-20240818-215425-I2NUVZU.md",
		"file3.md",
		".trash/file4.sync-conflict-20240818-215425-I2NUVZU.md",
	}

	for _, file := range files {
		path := filepath.Join(tempDir, file)
		err := os.MkdirAll(filepath.Dir(path), 0o755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.WriteFile(path, []byte("test content"), 0o644)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	finder := &DefaultFileFinder{}
	foundFiles, err := finder.FindSyncConflictFiles([]string{tempDir}, []string{".trash"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(tempDir, "file1.sync-conflict-20240818-215425-I2NUVZU.md"),
		filepath.Join(tempDir, "subdir", "file2.sync-conflict-20240818-215425-I2NUVZU.md"),
	}

	if len(foundFiles) != len(expectedFiles) {
		t.Errorf("Expected %d files, but got %d", len(expectedFiles), len(foundFiles))
	}

	for i, file := range foundFiles {
		if file != expectedFiles[i] {
			t.Errorf("Expected file %s, but got %s", expectedFiles[i], file)
		}
	}
}

func TestDefaultFileComparer_CompareAndDelete(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "file_comparer_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name            string
		conflictContent string
		originalContent string
		expectDeleted   bool
	}{
		{
			name:            "identical files",
			conflictContent: "test content",
			originalContent: "test content",
			expectDeleted:   true,
		},
		{
			name:            "different files",
			conflictContent: "test content 1",
			originalContent: "test content 2",
			expectDeleted:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conflictFile := filepath.Join(tempDir, "conflict.md")
			originalFile := filepath.Join(tempDir, "original.md")

			err := os.WriteFile(conflictFile, []byte(tt.conflictContent), 0o644)
			if err != nil {
				t.Fatalf("Failed to write conflict file: %v", err)
			}

			err = os.WriteFile(originalFile, []byte(tt.originalContent), 0o644)
			if err != nil {
				t.Fatalf("Failed to write original file: %v", err)
			}

			comparer := &DefaultFileComparer{}
			deleted, err := comparer.CompareAndDelete(conflictFile, originalFile)
			if err != nil {
				t.Fatalf("CompareAndDelete failed: %v", err)
			}

			if deleted != tt.expectDeleted {
				t.Errorf("Expected deleted to be %v, got %v", tt.expectDeleted, deleted)
			}

			if tt.expectDeleted {
				if _, err := os.Stat(conflictFile); !os.IsNotExist(err) {
					t.Errorf("Expected conflict file to be deleted, but it still exists")
				}
			} else {
				if _, err := os.Stat(conflictFile); os.IsNotExist(err) {
					t.Errorf("Expected conflict file to exist, but it was deleted")
				}
			}
		})
	}
}
