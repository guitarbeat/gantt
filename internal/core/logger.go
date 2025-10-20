// Package core - Logger provides structured logging with level-based control.
//
// The logging system supports multiple output formats and levels:
//   - trace: Detailed execution tracing
//   - debug: Detailed debugging information
//   - info: Informational messages (default)
//   - warn: Warning messages
//   - error: Error messages
//   - fatal: Fatal errors that terminate the program
//   - silent: No output
//
// Control logging via environment variables:
//   - PLANNER_SILENT=1: Suppress all output (backward compatible)
//   - PLANNER_LOG_LEVEL=trace|debug|info|warn|error|fatal|silent: Explicit level control
//   - PLANNER_LOG_FORMAT=text|json: Output format (default: text)
//   - PLANNER_LOG_FILE=/path/to/logfile: Write logs to file instead of stderr
//
// Structured logging with key-value pairs:
//
//	logger := core.NewDefaultLogger()
//	logger.Info("Processing file",
//	    "file", filename,
//	    "size", fileSize,
//	    "operation", "read")
//	logger.WithField("request_id", "12345").Info("Starting operation")
//	logger.WithFields(map[string]interface{}{
//	    "user_id": 123,
//	    "action": "create_task",
//	}).Debug("User action")
//
// Context-aware logging:
//
//	ctx := logger.WithContext(context.Background(), "correlation_id", "abc-123")
//	logger.FromContext(ctx).Info("Operation started")
//
// Backward compatibility:
//
//	logger.Info("Simple message: %s", value) // Still supported
//	if core.IsSilent() { /* Skip expensive logging */ }
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Log levels in order of increasing severity
const (
	LogLevelTrace = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	LogLevelSilent = 999 // Special level for silent mode
)

// String representations of log levels
const (
	LogLevelTraceString  = "trace"
	LogLevelDebugString  = "debug"
	LogLevelInfoString   = "info"
	LogLevelWarnString   = "warn"
	LogLevelErrorString  = "error"
	LogLevelFatalString  = "fatal"
	LogLevelSilentString = "silent"
)

// Environment variables for logging control
const (
	envPlannerSilent    = "PLANNER_SILENT"
	envPlannerLogLevel  = "PLANNER_LOG_LEVEL"
	envPlannerLogFormat = "PLANNER_LOG_FORMAT"
	envPlannerLogFile   = "PLANNER_LOG_FILE"
)

// LogFormat represents the output format for logs
type LogFormat int

const (
	LogFormatText LogFormat = iota
	LogFormatJSON
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Time    time.Time              `json:"time"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
	Caller  string                 `json:"caller,omitempty"`
	Prefix  string                 `json:"prefix,omitempty"`
}

// Logger provides structured logging with context support
type Logger struct {
	mu     sync.RWMutex
	writer io.Writer
	level  int
	format LogFormat
	prefix string
	fields map[string]interface{}
}

// contextKey is used for storing logger in context
type contextKey struct{}

// globalLogger is the default logger instance
var globalLogger *Logger
var globalLoggerOnce sync.Once

// NewLogger creates a new logger with the specified prefix
func NewLogger(prefix string) *Logger {
	level := parseLogLevel(getLogLevelString())
	format := parseLogFormat(os.Getenv(envPlannerLogFormat))
	writer := getLogWriter()

	return &Logger{
		writer: writer,
		level:  level,
		format: format,
		prefix: strings.TrimSpace(prefix),
		fields: make(map[string]interface{}),
	}
}

// NewDefaultLogger creates a logger with standard settings
func NewDefaultLogger() *Logger {
	globalLoggerOnce.Do(func() {
		globalLogger = NewLogger("[planner] ")
	})
	return globalLogger
}

// getLogWriter returns the appropriate writer for logs
func getLogWriter() io.Writer {
	if logFile := os.Getenv(envPlannerLogFile); logFile != "" {
		if file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			return file
		}
	}
	return os.Stderr
}

// getLogLevelString gets the log level from environment variables
func getLogLevelString() string {
	// Check PLANNER_SILENT first for backward compatibility
	if os.Getenv(envPlannerSilent) == "1" {
		return LogLevelSilentString
	}

	// Check PLANNER_LOG_LEVEL for explicit level
	level := strings.ToLower(os.Getenv(envPlannerLogLevel))
	if level == "" {
		return LogLevelInfoString // Default to info level
	}
	return level
}

// parseLogLevel converts string level to int level
func parseLogLevel(level string) int {
	switch level {
	case LogLevelTraceString:
		return LogLevelTrace
	case LogLevelDebugString:
		return LogLevelDebug
	case LogLevelInfoString:
		return LogLevelInfo
	case LogLevelWarnString:
		return LogLevelWarn
	case LogLevelErrorString:
		return LogLevelError
	case LogLevelFatalString:
		return LogLevelFatal
	case LogLevelSilentString:
		return LogLevelSilent
	default:
		return LogLevelInfo
	}
}

// parseLogFormat converts string format to LogFormat
func parseLogFormat(format string) LogFormat {
	switch strings.ToLower(format) {
	case "json":
		return LogFormatJSON
	default:
		return LogFormatText
	}
}

// IsSilent returns true if logging is suppressed
func IsSilent() bool {
	return parseLogLevel(getLogLevelString()) == LogLevelSilent
}

// WithField creates a new logger with an additional field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.mu.RLock()
	newFields := make(map[string]interface{}, len(l.fields)+1)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	l.mu.RUnlock()

	return &Logger{
		writer: l.writer,
		level:  l.level,
		format: l.format,
		prefix: l.prefix,
		fields: newFields,
	}
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.mu.RLock()
	newFields := make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	l.mu.RUnlock()

	return &Logger{
		writer: l.writer,
		level:  l.level,
		format: l.format,
		prefix: l.prefix,
		fields: newFields,
	}
}

// WithContext adds logger to context
func (l *Logger) WithContext(ctx context.Context, key string, value interface{}) context.Context {
	logger := l.WithField(key, value)
	return context.WithValue(ctx, contextKey{}, logger)
}

// FromContext retrieves logger from context
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(contextKey{}).(*Logger); ok {
		return logger
	}
	return NewDefaultLogger()
}

// log sends a log entry to the output
func (l *Logger) log(level int, levelStr string, message string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Format message if needed
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	caller := ""
	if ok {
		parts := strings.Split(file, "/")
		if len(parts) > 0 {
			file = parts[len(parts)-1]
		}
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	entry := LogEntry{
		Time:    time.Now(),
		Level:   levelStr,
		Message: message,
		Fields:  l.fields,
		Caller:  caller,
		Prefix:  l.prefix,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var output string
	switch l.format {
	case LogFormatJSON:
		if jsonBytes, err := json.Marshal(entry); err == nil {
			output = string(jsonBytes)
		} else {
			output = fmt.Sprintf("{\"error\":\"failed to marshal log entry: %v\"}", err)
		}
	default:
		output = l.formatTextEntry(entry)
	}

	if level == LogLevelFatal {
		fmt.Fprintln(l.writer, output)
		os.Exit(1)
	} else {
		fmt.Fprintln(l.writer, output)
	}
}

// formatTextEntry formats a log entry as text
func (l *Logger) formatTextEntry(entry LogEntry) string {
	var parts []string

	// Add timestamp
	parts = append(parts, entry.Time.Format("2006/01/02 15:04:05"))

	// Add level
	levelStr := strings.ToUpper(entry.Level)
	parts = append(parts, fmt.Sprintf("[%s]", levelStr))

	// Add prefix if present
	if entry.Prefix != "" {
		parts = append(parts, entry.Prefix)
	}

	// Add message
	parts = append(parts, entry.Message)

	// Add fields
	if len(entry.Fields) > 0 {
		var fieldParts []string
		for k, v := range entry.Fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		parts = append(parts, fmt.Sprintf("{%s}", strings.Join(fieldParts, " ")))
	}

	// Add caller if in debug/trace mode
	if l.level <= LogLevelDebug && entry.Caller != "" {
		parts = append(parts, fmt.Sprintf("(%s)", entry.Caller))
	}

	return strings.Join(parts, " ")
}

// Structured logging methods

// Trace logs a trace message
func (l *Logger) Trace(message string, args ...interface{}) {
	l.log(LogLevelTrace, LogLevelTraceString, message, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(LogLevelDebug, LogLevelDebugString, message, args...)
}

// Info logs an informational message
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(LogLevelInfo, LogLevelInfoString, message, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(LogLevelWarn, LogLevelWarnString, message, args...)
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	l.log(LogLevelError, LogLevelErrorString, message, args...)
}

// Fatal logs a fatal error message and exits
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.log(LogLevelFatal, LogLevelFatalString, message, args...)
}

// Backward compatibility methods

// Printf provides compatibility with existing log.Logger interface
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Info(format, v...)
}
