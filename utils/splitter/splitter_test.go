package splitter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/romisugianto/go-utils/utils/logger"
)

func createTestFile(t *testing.T, dir string) string {
	testFile := filepath.Join(dir, "testfile.csv")
	content := []byte("line1\nline2\nline3\nline4\nline5\n")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	return testFile
}

func TestSplitFileByLines(t *testing.T) {
	// Create logger instance
	testLogger, err := logger.NewLogger("splitter_test")
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer testLogger.Close()

	// Create Splitter with injected logger
	sp := NewSplitter(testLogger)

	tests := []struct {
		name          string
		linesPerFile  int
		expectedFiles int
	}{
		{"2 lines per file", 2, 3},
		{"1 line per file", 1, 5},
		{"more lines than file", 10, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup fresh directories and file for each test case
			testDir := t.TempDir()
			testFile := createTestFile(t, testDir)
			outputDir := filepath.Join(testDir, "output")
			processedDir := filepath.Join(testDir, "processed")

			// Execute
			if err := sp.SplitFileByLines(testFile, tt.linesPerFile, outputDir, processedDir); err != nil {
				t.Errorf("SplitFileByLines failed: %v", err)
			}

			// Verify output files
			outputFiles, err := os.ReadDir(outputDir)
			if err != nil {
				t.Fatalf("failed to read output directory: %v", err)
			}
			if len(outputFiles) != tt.expectedFiles {
				t.Errorf("expected %d output files, got %d", tt.expectedFiles, len(outputFiles))
			}

			// Verify processed file
			processedFiles, err := os.ReadDir(processedDir)
			if err != nil {
				t.Fatalf("failed to read processed directory: %v", err)
			}
			if len(processedFiles) != 1 {
				t.Errorf("expected 1 processed file, got %d", len(processedFiles))
			}
		})
	}
}

func TestSplitFileByLines_ErrorCases(t *testing.T) {
	testLogger, _ := logger.NewLogger("splitter_test")
	defer testLogger.Close()
	sp := NewSplitter(testLogger)

	tests := []struct {
		name         string
		filePath     string
		linesPerFile int
		expectError  bool
		errContains  string // Expected error message to contain
	}{
		{
			name:         "nonexistent file",
			filePath:     "/nonexistent/file",
			linesPerFile: 2,
			expectError:  true,
			errContains:  "failed to open file",
		},
		{
			name:         "zero lines",
			filePath:     "testfile.csv",
			linesPerFile: 0,
			expectError:  true,
			errContains:  "lines per file must be positive",
		},
		{
			name:         "negative lines",
			filePath:     "testfile.csv",
			linesPerFile: -1,
			expectError:  true,
			errContains:  "lines per file must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			outputDir := filepath.Join(testDir, "output")
			processedDir := filepath.Join(testDir, "processed")

			// Create test file if needed
			if tt.filePath == "testfile.csv" {
				tt.filePath = createTestFile(t, testDir)
			}

			err := sp.SplitFileByLines(tt.filePath, tt.linesPerFile, outputDir, processedDir)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}