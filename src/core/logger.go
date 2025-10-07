// Package core - Logger provides centralized logging with level-based control.
//
// The logging system supports three levels:
//   - silent: No output
//   - info: Informational messages (default)
//   - debug: Detailed debugging information
//
// Control logging via environment variables:
//   - PLANNER_SILENT=1: Suppress all output (backward compatible)
//   - PLANNER_LOG_LEVEL=silent|info|debug: Explicit level control
//
// Example usage:
//
//	logger := core.NewDefaultLogger()
//	logger.Info("Processing file: %s", filename)
//	logger.Debug("Detailed state: %+v", state)
//	logger.Warn("Potential issue: %s", warning)
//	logger.Error("Operation failed: %v", err)
//
//	// Check if silent mode is active
//	if core.IsSilent() {
//	    // Skip expensive logging operations
//	}
package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	// LogLevelSilent suppresses all log output
	LogLevelSilent = "silent"
	// LogLevelInfo shows informational messages
	LogLevelInfo = "info"
	// LogLevelDebug shows detailed debugging information
	LogLevelDebug = "debug"

	// Environment variables for logging control
	envPlannerSilent   = "PLANNER_SILENT"
	envPlannerLogLevel = "PLANNER_LOG_LEVEL"
)

// Logger provides centralized logging functionality with silent mode support
type Logger struct {
	logger *log.Logger
	level  string
}

// NewLogger creates a new logger with the specified prefix
func NewLogger(prefix string) *Logger {
	level := getLogLevel()

	var out io.Writer = os.Stderr
	if level == LogLevelSilent {
		out = io.Discard
	}

	return &Logger{
		logger: log.New(out, prefix, log.LstdFlags|log.Lshortfile),
		level:  level,
	}
}

// NewDefaultLogger creates a logger with standard settings
func NewDefaultLogger() *Logger {
	return NewLogger("[planner] ")
}

// getLogLevel determines the logging level from environment variables
func getLogLevel() string {
	// Check PLANNER_SILENT first for backward compatibility
	if os.Getenv(envPlannerSilent) == "1" {
		return LogLevelSilent
	}

	// Check PLANNER_LOG_LEVEL for explicit level
	level := strings.ToLower(os.Getenv(envPlannerLogLevel))
	switch level {
	case LogLevelSilent, LogLevelInfo, LogLevelDebug:
		return level
	default:
		return LogLevelInfo // Default to info level
	}
}

// IsSilent returns true if logging is suppressed
func IsSilent() bool {
	return getLogLevel() == LogLevelSilent
}

// Info logs an informational message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level != LogLevelSilent {
		_ = l.logger.Output(2, fmt.Sprintf("[INFO] "+format, v...))
	}
}

// Debug logs a debug message (only when debug level is enabled)
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level == LogLevelDebug {
		_ = l.logger.Output(2, fmt.Sprintf("[DEBUG] "+format, v...))
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level != LogLevelSilent {
		_ = l.logger.Output(2, fmt.Sprintf("[ERROR] "+format, v...))
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level != LogLevelSilent {
		_ = l.logger.Output(2, fmt.Sprintf("[WARN] "+format, v...))
	}
}

// Printf provides compatibility with existing log.Logger interface
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Info(format, v...)
}
