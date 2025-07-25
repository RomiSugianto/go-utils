package housekeeper

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/romisugianto/go-utils/utils/logger"
)

func TestHousekeepFilesByAge(t *testing.T) {
	tests := []struct {
		name         string
		setupFiles   map[string]time.Time // filename -> modTime
		maxAgeDays   int
		wantRemoved  []string
		wantKept     []string
		expectError  bool
	}{
		{
			name: "basic age cleanup",
			setupFiles: map[string]time.Time{
				"old.txt": time.Now().Add(-48 * time.Hour),
				"new.txt": time.Now(),
			},
			maxAgeDays:  1,
			wantRemoved: []string{"old.txt"},
			wantKept:    []string{"new.txt"},
		},
		{
			name: "no files to clean",
			setupFiles: map[string]time.Time{
				"new1.txt": time.Now(),
				"new2.txt": time.Now().Add(-12 * time.Hour),
			},
			maxAgeDays: 2,
			wantKept:   []string{"new1.txt", "new2.txt"},
		},
		{
			name: "invalid directory",
			setupFiles: map[string]time.Time{},
			maxAgeDays: 1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			testDir := t.TempDir()

			// Create test files with specific mod times
			for name, modTime := range tt.setupFiles {
				path := filepath.Join(testDir, name)
				if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				if err := os.Chtimes(path, modTime, modTime); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			// Create housekeeper with test logger
			testLogger, _ := logger.NewLogger("housekeeper_test")
			defer testLogger.Close()
			hk := NewHousekeeper(testLogger)

			// For invalid directory test
			testPath := testDir
			if tt.expectError {
				testPath = filepath.Join(testDir, "nonexistent")
			}

			// Execute
			err := hk.HousekeepFilesByAge(testPath, tt.maxAgeDays)

			// Verify results
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check removed files
			for _, name := range tt.wantRemoved {
				path := filepath.Join(testDir, name)
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Errorf("file %q should have been removed", name)
				}
			}

			// Check kept files
			for _, name := range tt.wantKept {
				path := filepath.Join(testDir, name)
				if _, err := os.Stat(path); err != nil {
					t.Errorf("file %q should not have been removed", name)
				}
			}
		})
	}
}

func TestHousekeepFilesByCount(t *testing.T) {
	tests := []struct {
		name         string
		fileCount    int
		maxFiles     int
		wantRemoved  int
		expectError  bool
	}{
		{"keep 3 of 5", 5, 3, 2, false},
		{"keep all when under limit", 2, 5, 0, false},
		{"keep none when zero requested", 5, 0, 5, false},
		{"invalid directory", 0, 3, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()

			// Create test files
			for i := 0; i < tt.fileCount; i++ {
				path := filepath.Join(testDir, fmt.Sprintf("file%d.txt", i))
				if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				// Space out mod times by 1 hour
				modTime := time.Now().Add(-time.Duration(i) * time.Hour)
				if err := os.Chtimes(path, modTime, modTime); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			// Create housekeeper
			testLogger, _ := logger.NewLogger("housekeeper_test")
			defer testLogger.Close()
			hk := NewHousekeeper(testLogger)

			// For invalid directory test
			testPath := testDir
			if tt.expectError {
				testPath = filepath.Join(testDir, "nonexistent")
			}

			// Execute
			err := hk.HousekeepFilesByCount(testPath, tt.maxFiles)

			// Verify
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check remaining files
			files, err := os.ReadDir(testDir)
			if err != nil {
				t.Fatalf("failed to read dir: %v", err)
			}

			remaining := len(files)
			expected := tt.fileCount - tt.wantRemoved
			if remaining != expected {
				t.Errorf("expected %d files remaining, got %d", expected, remaining)
			}
		})
	}
}