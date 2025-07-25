package housekeeper

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHousekeepFilesByAge(t *testing.T) {
    // Setup: create a temp directory and some files
    dir := t.TempDir()
    oldFile := filepath.Join(dir, "old.txt")
    newFile := filepath.Join(dir, "new.txt")

    // Create an old file (mod time 48 hours ago)
    if err := os.WriteFile(oldFile, []byte("old"), 0644); err != nil {
        t.Fatalf("failed to create old file: %v", err)
    }
    twoDaysAgo := time.Now().Add(-48 * time.Hour)
    if err := os.Chtimes(oldFile, twoDaysAgo, twoDaysAgo); err != nil {
        t.Fatalf("failed to set old file time: %v", err)
    }

    // Create a new file (mod time now)
    if err := os.WriteFile(newFile, []byte("new"), 0644); err != nil {
        t.Fatalf("failed to create new file: %v", err)
    }

    // Create a Housekeeper
    hk, err := NewHousekeeper("testApp")
    if err != nil {
        t.Fatalf("failed to create housekeeper: %v", err)
    }
    defer hk.logger.Close()

    // Run housekeeping: should remove oldFile, keep newFile
    if err := hk.HousekeepFilesByAge(dir, 1); err != nil {
        t.Errorf("HousekeepFilesByAge failed: %v", err)
    }

    // Check results
    if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
        t.Errorf("old file was not removed")
    }
    if _, err := os.Stat(newFile); err != nil {
        t.Errorf("new file was removed unexpectedly")
    }
}

func TestHousekeepFilesByCount(t *testing.T) {
		// Setup: create a temp directory and some files
		dir := t.TempDir()
		for i := 0; i < 5; i++ {
				filePath := filepath.Join(dir, fmt.Sprintf("file_%d.txt", i))
				if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
						t.Fatalf("failed to create file %d: %v", i, err)
				}
		}

		// Create a Housekeeper
		hk, err := NewHousekeeper("testApp")
		if err != nil {
				t.Fatalf("failed to create housekeeper: %v", err)
		}
		defer hk.logger.Close()

		// Run housekeeping: should keep only 3 files
		if err := hk.HousekeepFilesByCount(dir, 3); err != nil {
				t.Errorf("HousekeepFilesByCount failed: %v", err)
		}

		// Check results
		files, err := os.ReadDir(dir)
		if err != nil {
				t.Fatalf("failed to read directory: %v", err)
		}
		if len(files) > 3 {
				t.Errorf("expected at most 3 files, got %d", len(files))
		}
}