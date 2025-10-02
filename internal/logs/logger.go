package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	logger  *log.Logger
	logFile *os.File
)

// Init initializes the logging system
func Init(logDir string) error {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logDir, fmt.Sprintf("p-monitor_%s.log", timestamp))

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}

	logFile = file
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)

	// Also log to stdout for debugging
	logger.SetOutput(os.Stdout)

	Info("Logging initialized at %s", logPath)
	return nil
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf("[INFO] "+format, args...)
	}
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf("[ERROR] "+format, args...)
	}
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf("[DEBUG] "+format, args...)
	}
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
