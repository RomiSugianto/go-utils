package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	t.Run("successful creation with app name", func(t *testing.T) {
		appName := "testapp"
		logger, err := NewLogger(appName)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()

		expectedFilename := appName + "_" + time.Now().Format("2006-01-02") + ".log"
		if !strings.HasSuffix(logger.GetLogFilePath(), expectedFilename) {
			t.Errorf("Expected log file to end with %q, got %q", expectedFilename, logger.GetLogFilePath())
		}

		// Verify logs directory was created
		if _, err := os.Stat("logs"); os.IsNotExist(err) {
			t.Error("Expected logs directory to be created")
		}
	})

	t.Run("successful creation with empty app name", func(t *testing.T) {
		logger, err := NewLogger("")
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()

		expectedFilename := "script_" + time.Now().Format("2006-01-02") + ".log"
		if !strings.HasSuffix(logger.GetLogFilePath(), expectedFilename) {
			t.Errorf("Expected log file to end with %q, got %q", expectedFilename, logger.GetLogFilePath())
		}
	})

	t.Run("invalid directory permissions", func(t *testing.T) {
        // Skip this test on Windows as permission handling is different
        if strings.ToLower(os.Getenv("GOOS")) == "windows" {
            t.Skip("Skipping permission test on Windows")
        }

        // Create a temporary directory
        tempDir, err := os.MkdirTemp("", "logger_test")
        if err != nil {
            t.Fatalf("Failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tempDir)

        // Create a subdirectory with no permissions
        noPermDir := filepath.Join(tempDir, "noperm")
        if err := os.Mkdir(noPermDir, 0000); err != nil {
            t.Fatalf("Failed to create no-permission dir: %v", err)
        }
        defer os.Chmod(noPermDir, 0755) // Cleanup

        // Try to create a file in the no-permission directory
        testPath := filepath.Join(noPermDir, "test.log")
        _, err = os.OpenFile(testPath, os.O_CREATE|os.O_WRONLY, 0644)
        if err == nil {
            t.Error("Expected error when creating file in no-permission directory")
            os.Remove(testPath)
        }
    })
}

func TestLoggerMethods(t *testing.T) {
	logger, err := NewLogger("testlogger")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Test log methods
	t.Run("info log", func(t *testing.T) {
		logger.Info("This is an info message: %d", 42)
	})

	t.Run("warning log", func(t *testing.T) {
		logger.Warning("This is a warning message: %s", "be careful")
	})

	t.Run("error log", func(t *testing.T) {
		logger.Error("This is an error message: %v", os.ErrNotExist)
	})

	t.Run("display credits", func(t *testing.T) {
		banner := `
=== %s ===
Version: %s
=============
`
		logger.DisplayCredits(banner, "testapp", "1.0.0")
	})
}

func TestClose(t *testing.T) {
	logger, err := NewLogger("testclose")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if err := logger.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// Verify file is actually closed by trying to write to it
	_, err = logger.logFile.WriteString("test")
	if err == nil {
		t.Error("Expected error when writing to closed file")
	}
}

func TestGetLogFilePath(t *testing.T) {
	logger, err := NewLogger("testpath")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	expectedPath := filepath.Join("logs", "testpath_"+time.Now().Format("2006-01-02")+".log")
	if logger.GetLogFilePath() != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, logger.GetLogFilePath())
	}
}

func TestLogFormatting(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "logger_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Create a proper logger instance with the temp file
	logger := &Logger{
		logFile: tmpFile,
		logPath: tmpFile.Name(),
	}
	defer logger.Close()

	// Test the log formatting
	testMsg := "test message %d"
	testArgs := []any{123}

	logger.log("TEST", testMsg, testArgs...)

	// Flush to disk
	if err := tmpFile.Sync(); err != nil {
		t.Fatalf("Failed to sync file: %v", err)
	}

	// Now open the file and check its contents
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("No content written to log file")
	}

	// Get the current time up to the second (same precision used in formatting)
	now := time.Now().Format("2006-01-02 15:04:05")
	expectedPrefix := fmt.Sprintf("[%s] [TEST] test message 123", now)
	if !strings.HasPrefix(string(content), expectedPrefix) {
		t.Errorf("Expected log line to start with %q, got %q", expectedPrefix, string(content))
	}
}