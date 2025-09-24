package common

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	TaskID  string
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("task %s, field '%s', value '%s': %s", e.TaskID, e.Field, e.Value, e.Message)
}

// ParseError represents an error during CSV parsing
type ParseError struct {
	Row     int
	Column  string
	Value   string
	Message string
	Err     error
}

func (e *ParseError) Error() string {
	if e.Row > 0 {
		return fmt.Sprintf("parse error at row %d, column %q (value: %q): %s", e.Row, e.Column, e.Value, e.Message)
	}
	return fmt.Sprintf("parse error at column %q (value: %q): %s", e.Column, e.Value, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// ConfigError represents a configuration error
type ConfigError struct {
	File    string
	Section string
	Message string
}

func (e *ConfigError) Error() string {
	if e.Section != "" {
		return fmt.Sprintf("config error in %s, section %q: %s", e.File, e.Section, e.Message)
	}
	return fmt.Sprintf("config error in %s: %s", e.File, e.Message)
}

// MultiError represents multiple errors
type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	
	var messages []string
	for i, err := range e.Errors {
		messages = append(messages, fmt.Sprintf("error %d: %s", i+1, err.Error()))
	}
	return strings.Join(messages, "; ")
}

func (e *MultiError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

// NewMultiError creates a new MultiError
func NewMultiError() *MultiError {
	return &MultiError{
		Errors: make([]error, 0),
	}
}

// WrapError wraps an error with additional context
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsParseError checks if an error is a ParseError
func IsParseError(err error) bool {
	_, ok := err.(*ParseError)
	return ok
}

// IsConfigError checks if an error is a ConfigError
func IsConfigError(err error) bool {
	_, ok := err.(*ConfigError)
	return ok
}