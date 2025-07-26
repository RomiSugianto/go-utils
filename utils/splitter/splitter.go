package splitter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/romisugianto/go-utils/utils/logger"
)

// Splitter handles file splitting operations
type Splitter struct {
	logger *logger.Logger
}

// NewSplitter creates a new splitter instance
func NewSplitter(appName string) (*Splitter, error) {
	// Create a new logger instance
	log, err := logger.NewLogger(appName)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &Splitter{
		logger: log,
	}, nil
}

// SplitFileByLines splits a file into multiple files based on the number of lines specified.
// It creates output files in the specified output directory.
// If the processed directory is specified, it moves the original file there after splitting.
// It logs the splitting summary including the number of files created and the time taken.
// The original file is moved to the processed directory after splitting.
// It returns an error if any operation fails.
func (s *Splitter) SplitFileByLines(filePath string, linesPerFile int, outputDir string, processedDir string) error {
	startTime := time.Now()

	// Ensure the output and processed directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	if err := os.MkdirAll(processedDir, 0755); err != nil {
		return fmt.Errorf("failed to create processed directory: %w", err)
	}

	// Open the source file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Get file stats for size info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats: %w", err)
	}
	fileSize := fileInfo.Size()

	s.logger.Info("Starting to process file: %s (size: %.2f MB)", filePath, float64(fileSize)/1024/1024)

	// Get base filename without extension
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(filePath)
	baseName := strings.TrimSuffix(fileName, fileExt)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	linesCount := 0
	fileCount := 1
	var outputFile *os.File
	var writer *bufio.Writer

	// Process each line in the file
	for scanner.Scan() {
		// If we've reached the line limit or haven't created the first output file yet
		if linesCount%linesPerFile == 0 {
			// Close the previous file if it exists
			if outputFile != nil {
				writer.Flush()
				outputFile.Close()
				s.logger.Info("Created output file part %d", fileCount-1)
			}

			// Create a new output file
			outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_part%d%s", baseName, fileCount, fileExt))
			outputFile, err = os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
			}
			writer = bufio.NewWriter(outputFile)
			fileCount++
		}

		// Write the line to the output file
		_, err := writer.WriteString(scanner.Text() + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
		linesCount++
	}

	// Check if there was an error during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Make sure to flush and close the last file
	if writer != nil {
		writer.Flush()
	}
	if outputFile != nil {
		outputFile.Close()
		s.logger.Info("Created final output file part %d", fileCount-1)
	}

	// Calculate actual number of files created (could be one less if file ended exactly on a boundary)
	actualFileCount := fileCount - 1

		// Close any possible open handles to ensure we can move the file
	file.Close()

	// Move the original file to the processed directory using os.Rename
	processedPath := filepath.Join(processedDir, fileName)
	if err := os.Rename(filePath, processedPath); err != nil {
		return fmt.Errorf("failed to move file to processed directory: %w", err)
	}

	// Calculate processing duration
	duration := time.Since(startTime)

	// Log processing summary
	s.logger.Info("Processed file: %s", fileName)
	s.logger.Info("  - Original size: %.2f MB", float64(fileSize)/1024/1024)
	s.logger.Info("  - Files created: %d", actualFileCount)
	s.logger.Info("  - Processing time: %.2f seconds", duration.Seconds())
	s.logger.Info("  - Processing rate: %.2f MB/sec", (float64(fileSize)/1024/1024)/duration.Seconds())
	s.logger.Info("  - Processed file moved to: %s", processedPath)
	return nil
}