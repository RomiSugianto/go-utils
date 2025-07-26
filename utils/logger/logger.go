package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Logger provides logging capabilities with file and console output
type Logger struct {
	logFile    *os.File
	logPath    string
}

// NewLogger creates a new logger instance
func NewLogger(appName string) (*Logger, error) {
	// Use provided app name or fallback to default
	if appName == "" {
		appName = "script"
	}

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logsDir, fmt.Sprintf("%s_%s.log", appName, timestamp))

	// open file in append mode or create if it doesn't exist
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		logFile:    logFile,
		logPath:    logPath,
	}, nil
}

// Close closes the logger's file handle
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// logRaw rewrites a raw message to the log file
func (l *Logger) logRaw(message string) {
	// Print to console
	fmt.Print(message)

	// Write to log file
	if l.logFile != nil {
		l.logFile.WriteString(message)
		l.logFile.Sync() // Ensure it's written to disk
	}
}

// GetLogFilePath returns the path to the current log file
func (l *Logger) GetLogFilePath() string {
	return l.logPath
}

// log writes a log message to both stdout and the log file
func (l *Logger) log(level, format string, args ...any) {
	// Format the message with timestamp and level
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	formattedMsg := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	// Write to stdout
	fmt.Print(formattedMsg)

	// Write to log file
	if l.logFile != nil {
		l.logFile.WriteString(formattedMsg)
		l.logFile.Sync() // Ensure it's written to disk
	}
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...any) {
	l.log("INFO", format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...any) {
	l.log("ERROR", format, args...)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...any) {
	l.log("WARNING", format, args...)
}

// Summary logs a summary message
func (l *Logger) Summary(format string, args ...any) {
	l.log("SUMMARY", format, args...)
}

// DisplayCredits displays the application credits/banner
func (l *Logger) DisplayCredits(banner string, appName string, appVersion string) {
	appNameUpper := strings.ToUpper(appName)
	formattedBanner := fmt.Sprintf(banner, appNameUpper, appVersion)

	// Log the banner to log file raw
	l.logRaw(formattedBanner)

	// Also log the banner to the file
	l.Info("%s v%s started", appNameUpper, appVersion)
}