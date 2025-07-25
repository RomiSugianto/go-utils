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
func NewHousekeeper(appName string) (*Housekeeper, error) {
	// Create a new logger instance
	log, err := logger.NewLogger(appName)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &Housekeeper{
		logger: log,
	}, nil
}

// HousekeepFilesByAge manages the housekeeping of files in a directory based on their age
func (h *Housekeeper) HousekeepFilesByAge(dir string, maxAgeDays int) error {
	// Ensure the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		h.logger.Error("directory does not exist: %s", dir)
		return nil
	}

	// Get current time
	now := time.Now()

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		h.logger.Error("failed to read directory: %v", err)
		return nil
	}

	var removedFiles []string

	// Iterate over files and check their modification time
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dir, file.Name())
			info, err := file.Info()
			if err != nil {
				h.logger.Error("Error getting info for file %s: %v", filePath, err)
				continue
			}

			// Check if the file is older than maxAgeDays
			if now.Sub(info.ModTime()).Hours() > float64(maxAgeDays*24) {
				if err := os.Remove(filePath); err != nil {
					h.logger.Error("Failed to remove file %s: %v", filePath, err)
				} else {
					removedFiles = append(removedFiles, filePath)
				}
			}
		}
	}

	sort.Strings(removedFiles)
	for _, removedFile := range removedFiles {
		h.logger.Info("Removed old file: %s", removedFile)
	}

	return nil
}

// housekeepFilesByCount manages the housekeeping of files in a directory based on a maximum count
func (h *Housekeeper) HousekeepFilesByCount(dir string, maxFiles int) error {
	// Ensure the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		h.logger.Error("directory does not exist: %s", dir)
		return nil
	}

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		h.logger.Error("failed to read directory: %v", err)
		return nil
	}

	// Sort files by modification time (oldest first)
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	if len(files) <= maxFiles {
		h.logger.Info("No files to remove, count is within limit.")
		return nil
	}

	var removedFiles []string

	// Remove oldest files until we reach the maxFiles limit
	for i := 0; i < len(files)-maxFiles; i++ {
		filePath := filepath.Join(dir, files[i].Name())
		if err := os.Remove(filePath); err != nil {
			h.logger.Error("Failed to remove file %s: %v", filePath, err)
			continue
		}
		removedFiles = append(removedFiles, filePath)
	}

	sort.Strings(removedFiles)
	for _, removedFile := range removedFiles {
		h.logger.Info("Removed old file: %s", removedFile)
	}

	return nil
}