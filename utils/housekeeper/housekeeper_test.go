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
		setupDirs    []string            // subdirectories to create
		maxAgeDays   int
		recursive    bool
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
			name:        "invalid directory",
			setupFiles:  map[string]time.Time{},
			maxAgeDays:  1,
			expectError: true,
		},
		{
			name: "recursive cleanup",
			setupFiles: map[string]time.Time{
				"root_old.txt":    time.Now().Add(-48 * time.Hour),
				"root_new.txt":     time.Now(),
				"sub1/old.txt":    time.Now().Add(-72 * time.Hour),
				"sub1/new.txt":     time.Now(),
				"sub2/sub_old.txt": time.Now().Add(-96 * time.Hour),
				"sub2/sub_new.txt": time.Now(),
			},
			setupDirs:    []string{"sub1", "sub2"},
			maxAgeDays:   1,
			recursive:    true,
			wantRemoved:  []string{"root_old.txt", "sub1/old.txt", "sub2/sub_old.txt"},
			wantKept:     []string{"root_new.txt", "sub1/new.txt", "sub2/sub_new.txt"},
		},
		{
			name: "non-recursive cleanup with subdirs",
			setupFiles: map[string]time.Time{
				"root_old.txt":    time.Now().Add(-48 * time.Hour),
				"root_new.txt":     time.Now(),
				"sub1/old.txt":    time.Now().Add(-72 * time.Hour),
				"sub1/new.txt":     time.Now(),
			},
			setupDirs:    []string{"sub1"},
			maxAgeDays:   1,
			recursive:    false,
			wantRemoved:  []string{"root_old.txt"},
			wantKept:     []string{"root_new.txt", "sub1/old.txt", "sub1/new.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			testDir := t.TempDir()

			// Create subdirectories
			for _, dir := range tt.setupDirs {
				fullPath := filepath.Join(testDir, dir)
				if err := os.MkdirAll(fullPath, 0755); err != nil {
					t.Fatalf("setup failed creating directory %s: %v", dir, err)
				}
			}

			// Create test files with specific mod times
			for name, modTime := range tt.setupFiles {
				path := filepath.Join(testDir, name)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatalf("setup failed creating parent directories for %s: %v", name, err)
				}
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
			hk, err := NewHousekeeper(testLogger)
			if err != nil {
				t.Fatalf("failed to create housekeeper: %v", err)
			}

			// For invalid directory test
			testPath := testDir
			if tt.expectError {
				testPath = filepath.Join(testDir, "nonexistent")
			}

			// Execute
			err = hk.HousekeepFilesByAge(testPath, tt.maxAgeDays, tt.recursive)

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
					t.Errorf("file %q should not have been removed: %v", name, err)
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
			hk, err := NewHousekeeper(testLogger)
			if err != nil {
				t.Fatalf("failed to create housekeeper: %v", err)
			}

			// For invalid directory test
			testPath := testDir
			if tt.expectError {
				testPath = filepath.Join(testDir, "nonexistent")
			}

			// Execute
			err = hk.HousekeepFilesByCount(testPath, tt.maxFiles)

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