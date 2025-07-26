package logger

import (
	"fmt"
	"os"
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
			t.Errorf("Expected log file to end with %s, got %s", expectedFilename, logger.GetLogFilePath())
		}

		// Verify log file was created
		if _, err := os.Stat(logger.GetLogFilePath()); os.IsNotExist(err) {
			t.Errorf("Log file was not created")
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
			t.Errorf("Expected log file to end with %s, got %s", expectedFilename, logger.GetLogFilePath())
		}
	})

	t.Run("fails when directory cannot be created", func(t *testing.T) {
		// Create a read-only directory to prevent creation of logs subdirectory
		dir := "readonly_dir"
		os.Mkdir(dir, 0444)
		defer os.RemoveAll(dir)

		oldDir, _ := os.Getwd()
		os.Chdir(dir)
		defer os.Chdir(oldDir)

		_, err := NewLogger("testapp")
		if err == nil {
			t.Error("Expected error when directory cannot be created")
		}
	})
}

func TestLoggingFunctions(t *testing.T) {
	// Setup
	appName := "testapp"
	logger, err := NewLogger(appName)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		logger.Close()
		// Clean up log file
		os.RemoveAll("logs")
	}()

	t.Run("Info logging", func(t *testing.T) {
		testMsg := "This is an info message"
		logger.Info("%s", testMsg)

		// Verify the log file contains the message
		content, err := os.ReadFile(logger.GetLogFilePath())
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		if !strings.Contains(string(content), "[INFO] "+testMsg) {
			t.Errorf("Log file does not contain expected info message")
		}
	})

	t.Run("Error logging", func(t *testing.T) {
		testMsg := "This is an error message"
		logger.Error("%s", testMsg)

		content, err := os.ReadFile(logger.GetLogFilePath())
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		if !strings.Contains(string(content), "[ERROR] "+testMsg) {
			t.Errorf("Log file does not contain expected error message")
		}
	})

	t.Run("Warning logging", func(t *testing.T) {
		testMsg := "This is a warning message"
		logger.Warning("%s", testMsg)

		content, err := os.ReadFile(logger.GetLogFilePath())
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		if !strings.Contains(string(content), "[WARNING] "+testMsg) {
			t.Errorf("Log file does not contain expected warning message")
		}
	})
}

func TestDisplayCredits(t *testing.T) {
	appName := "testapp"
	appVersion := "1.0.0"
	banner := `
=======================================
= %s - Version %s =
=======================================
`

	logger, err := NewLogger(appName)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		logger.Close()
		os.RemoveAll("logs")
	}()

	logger.DisplayCredits(banner, appName, appVersion)

	content, err := os.ReadFile(logger.GetLogFilePath())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	expectedBanner := fmt.Sprintf(banner, strings.ToUpper(appName), appVersion)
	if !strings.Contains(string(content), expectedBanner) {
		t.Errorf("Log file does not contain expected banner")
	}

	if !strings.Contains(string(content), "[INFO] TESTAPP v1.0.0 started") {
		t.Errorf("Log file does not contain expected startup message")
	}
}

func TestClose(t *testing.T) {
	logger, err := NewLogger("testapp")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if err := logger.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// Verify we can't write after closing
	testMsg := "This should fail"
	n, err := logger.logFile.WriteString(testMsg)
	if err == nil || n > 0 {
		t.Error("Should not be able to write after Close()")
	}
}