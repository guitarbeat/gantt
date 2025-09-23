package common

import (
	"fmt"
	"time"
)

// DataValidationError represents a validation error with detailed context
type DataValidationError struct {
	Type        string    // Error type (e.g., "DATE_RANGE", "CONFLICT", "DEPENDENCY")
	TaskID      string    // Task ID that has the error
	Field       string    // Field that has the error
	Value       string    // Value that caused the error
	Message     string    // Human-readable error message
	Severity    string    // Error severity (ERROR, WARNING, INFO)
	Timestamp   time.Time // When the error was detected
	Suggestions []string  // Suggested fixes
}

// Error returns the error message
func (ve *DataValidationError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", ve.Severity, ve.TaskID, ve.Message)
}

// DateValidator handles date range validation and conflict detection
type DateValidator struct {
	workDays   map[time.Weekday]bool
	holidays   []time.Time
	timezone   *time.Location
	strictMode bool
}

// NewDateValidator creates a new date validator
func NewDateValidator() *DateValidator {
	return &DateValidator{
		workDays: map[time.Weekday]bool{
			time.Monday:    true,
			time.Tuesday:   true,
			time.Wednesday: true,
			time.Thursday:  true,
			time.Friday:    true,
			time.Saturday:  false,
			time.Sunday:    false,
		},
		holidays:   []time.Time{},
		timezone:   time.UTC,
		strictMode: false,
	}
}

// ValidateDateRanges validates date ranges for a slice of tasks
func (dv *DateValidator) ValidateDateRanges(tasks []*Task) []DataValidationError {
	var errors []DataValidationError

	for _, task := range tasks {
		// Basic date validation
		if task.StartDate.After(task.EndDate) {
			errors = append(errors, DataValidationError{
				Type:        "DATE_RANGE",
				TaskID:      task.ID,
				Field:       "dates",
				Value:       fmt.Sprintf("%s - %s", task.StartDate.Format("2006-01-02"), task.EndDate.Format("2006-01-02")),
				Message:     "Start date is after end date",
				Severity:    "ERROR",
				Timestamp:   time.Now(),
				Suggestions: []string{"Correct the start or end date"},
			})
		}
	}

	return errors
}
