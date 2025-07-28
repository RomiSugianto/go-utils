// Created by Romi Sugianto - https://romisugi.dev
package housekeeper

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/romisugianto/go-utils/utils/logger"
)

// Processor handles file splitting operations
type Housekeeper struct {
	logger *logger.Logger
}

// NewHousekeep creates a new Housekeep instance
func NewHousekeeper(log *logger.Logger) *Housekeeper {
	return &Housekeeper{logger: log}
}

// HousekeepFilesByAge manages the housekeeping of files in a directory based on their age
func (h *Housekeeper) HousekeepFilesByAge(dir string, maxAgeDays int, recursive ...bool) error {
	if maxAgeDays < 0 {
		return fmt.Errorf("maxAgeDays must be >= 0, got %d", maxAgeDays)
	}

	// Verify directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	recursiveFlag := false
	if len(recursive) > 0 {
		recursiveFlag = recursive[0]
	}

	var removed []string
	now := time.Now()
	cutoff := now.Add(-time.Duration(maxAgeDays*24) * time.Hour)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			h.logger.Error("Error accessing path %s: %v", path, err)
			return nil
		}

		// Skip directories unless we're processing them recursively
		if info.IsDir() {
			if !recursiveFlag && path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				h.logger.Error("Failed to remove file %s: %v", path, err)
				return nil
			}
			removed = append(removed, path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory %s: %w", dir, err)
	}

	h.logRemovals(removed, "age-based cleanup")
	return nil
}

func (h *Housekeeper) HousekeepFilesByCount(dir string, maxFiles int) error {
	if maxFiles < 0 {
		return fmt.Errorf("maxFiles must be >= 0, got %d", maxFiles)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	if len(files) <= maxFiles {
		h.logger.Info("No files to remove (current: %d, max: %d)", len(files), maxFiles)
		return nil
	}

	// Sort by mod time (oldest first)
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	var removed []string
	for i := 0; i < len(files)-maxFiles; i++ {
		path := filepath.Join(dir, files[i].Name())
		if err := os.Remove(path); err != nil {
			h.logger.Error("Failed to remove file %s: %v", path, err)
			continue
		}
		removed = append(removed, path)
	}

	h.logRemovals(removed, "count-based cleanup")
	return nil
}

func (h *Housekeeper) logRemovals(files []string, operation string) {
	if len(files) == 0 {
		h.logger.Summary("No files removed during %s", operation)
		return
	}

	sort.Strings(files)
	h.logger.Summary("Removed %d files during %s:", len(files), operation)
	for _, f := range files {
		h.logger.Summary("  - %s", f)
	}
}