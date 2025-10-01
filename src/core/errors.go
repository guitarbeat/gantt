package core

import (
	"fmt"
	"strings"
)

// ConfigError represents an error that occurred during configuration loading or validation
type ConfigError struct {
	File    string // Configuration file that caused the error
	Field   string // Specific field that caused the error (optional)
	Message string // Human-readable error message
	Err     error  // Underlying error (optional)
}

func (e *ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("config error in %s, field '%s': %s", e.File, e.Field, e.Message)
	}
	return fmt.Sprintf("config error in %s: %s", e.File, e.Message)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new configuration error
func NewConfigError(file, field, message string, err error) *ConfigError {
	return &ConfigError{
		File:    file,
		Field:   field,
		Message: message,
		Err:     err,
	}
}

// FileError represents an error that occurred during file operations
type FileError struct {
	Path      string // File path
	Operation string // Operation that failed (read, write, open, etc.)
	Err       error  // Underlying error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("file error during %s of %s: %v", e.Operation, e.Path, e.Err)
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// NewFileError creates a new file error
func NewFileError(path, operation string, err error) *FileError {
	return &FileError{
		Path:      path,
		Operation: operation,
		Err:       err,
	}
}

// TemplateError represents an error that occurred during template processing
type TemplateError struct {
	Template string // Template name
	Line     int    // Line number (if known)
	Message  string // Error message
	Err      error  // Underlying error
}

func (e *TemplateError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("template error in %s at line %d: %s", e.Template, e.Line, e.Message)
	}
	return fmt.Sprintf("template error in %s: %s", e.Template, e.Message)
}

func (e *TemplateError) Unwrap() error {
	return e.Err
}

// NewTemplateError creates a new template error
func NewTemplateError(template string, line int, message string, err error) *TemplateError {
	return &TemplateError{
		Template: template,
		Line:     line,
		Message:  message,
		Err:      err,
	}
}

// DataError represents an error that occurred during data processing
type DataError struct {
	Source  string // Data source (CSV file, database, etc.)
	Row     int    // Row number (if applicable)
	Column  string // Column name (if applicable)
	Message string // Error message
	Err     error  // Underlying error
}

func (e *DataError) Error() string {
	if e.Row > 0 && e.Column != "" {
		return fmt.Sprintf("data error in %s at row %d, column '%s': %s", e.Source, e.Row, e.Column, e.Message)
	} else if e.Row > 0 {
		return fmt.Sprintf("data error in %s at row %d: %s", e.Source, e.Row, e.Message)
	}
	return fmt.Sprintf("data error in %s: %s", e.Source, e.Message)
}

func (e *DataError) Unwrap() error {
	return e.Err
}

// NewDataError creates a new data error
func NewDataError(source string, row int, column, message string, err error) *DataError {
	return &DataError{
		Source:  source,
		Row:     row,
		Column:  column,
		Message: message,
		Err:     err,
	}
}

// ErrorAggregator collects multiple errors and provides summary reporting
type ErrorAggregator struct {
	Errors   []error
	Warnings []error
}

// NewErrorAggregator creates a new error aggregator
func NewErrorAggregator() *ErrorAggregator {
	return &ErrorAggregator{
		Errors:   make([]error, 0),
		Warnings: make([]error, 0),
	}
}

// AddError adds an error to the aggregator
func (ea *ErrorAggregator) AddError(err error) {
	if err != nil {
		ea.Errors = append(ea.Errors, err)
	}
}

// AddWarning adds a warning to the aggregator
func (ea *ErrorAggregator) AddWarning(err error) {
	if err != nil {
		ea.Warnings = append(ea.Warnings, err)
	}
}

// HasErrors returns true if there are any errors
func (ea *ErrorAggregator) HasErrors() bool {
	return len(ea.Errors) > 0
}

// HasWarnings returns true if there are any warnings
func (ea *ErrorAggregator) HasWarnings() bool {
	return len(ea.Warnings) > 0
}

// ErrorCount returns the number of errors
func (ea *ErrorAggregator) ErrorCount() int {
	return len(ea.Errors)
}

// WarningCount returns the number of warnings
func (ea *ErrorAggregator) WarningCount() int {
	return len(ea.Warnings)
}

// Error implements the error interface, returns first error
func (ea *ErrorAggregator) Error() string {
	if len(ea.Errors) == 0 {
		return "no errors"
	}
	if len(ea.Errors) == 1 {
		return ea.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred (first: %v)", len(ea.Errors), ea.Errors[0])
}

// Summary returns a comprehensive summary of all errors and warnings
func (ea *ErrorAggregator) Summary() string {
	if !ea.HasErrors() && !ea.HasWarnings() {
		return "No errors or warnings"
	}

	var summary strings.Builder
	
	if ea.HasErrors() {
		summary.WriteString(fmt.Sprintf("Errors (%d):\n", len(ea.Errors)))
		for i, err := range ea.Errors {
			summary.WriteString(fmt.Sprintf("  %d. %v\n", i+1, err))
		}
	}

	if ea.HasWarnings() {
		if ea.HasErrors() {
			summary.WriteString("\n")
		}
		summary.WriteString(fmt.Sprintf("Warnings (%d):\n", len(ea.Warnings)))
		for i, err := range ea.Warnings {
			summary.WriteString(fmt.Sprintf("  %d. %v\n", i+1, err))
		}
	}

	return summary.String()
}

// Clear clears all errors and warnings
func (ea *ErrorAggregator) Clear() {
	ea.Errors = make([]error, 0)
	ea.Warnings = make([]error, 0)
}

