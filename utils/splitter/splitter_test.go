package splitter

import (
	"os"
	"testing"
)

func TestSplitFileByLines(t *testing.T) {
	// Setup: create a temp directory and a test file
	testDir := t.TempDir()
	testFile := testDir + "/testfile.csv"
	if err := os.WriteFile(testFile, []byte("line1\nline2\nline3\nline4\nline5\n"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	// Create output and processed directories
	outputDir := testDir + "/output"
	processedDir := testDir + "/processed"

	// Create a new Splitter instance
	sp, err := NewSplitter("testApp")
	if err != nil {
		t.Fatalf("failed to create splitter: %v", err)
	}
	defer sp.logger.Close()

	// Run the SplitFileByLines method with test parameters
	if err := sp.SplitFileByLines(testFile, 2, outputDir, processedDir); err != nil {
		t.Errorf("SplitFileByLines failed: %v", err)
	}

	// Check if output files are created
	outputFiles, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("failed to read output directory: %v", err)
	}
	if len(outputFiles) != 3 {
		t.Errorf("expected 3 output files, got %d", len(outputFiles))
	}

}