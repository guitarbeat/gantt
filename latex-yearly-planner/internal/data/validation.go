package data

import (
	"fmt"
	"sort"
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

// ValidationResult contains the results of validation
type ValidationResult struct {
	IsValid     bool                  // Overall validation status
	Errors      []DataValidationError // List of validation errors
	Warnings    []DataValidationError // List of validation warnings
	Info        []DataValidationError // List of validation info messages
	Summary     string                // Summary of validation results
	Timestamp   time.Time             // When validation was performed
	TaskCount   int                   // Number of tasks validated
	ErrorCount  int                   // Number of errors found
	WarningCount int                  // Number of warnings found
}

// DateValidator handles date range validation and conflict detection
type DateValidator struct {
	workDays    map[time.Weekday]bool
	holidays    []time.Time
	timezone    *time.Location
	strictMode  bool
}

// NewDateValidator creates a new date validator
func NewDateValidator() *DateValidator {
	dv := &DateValidator{
		workDays:   make(map[time.Weekday]bool),
		holidays:   make([]time.Time, 0),
		timezone:   time.UTC,
		strictMode: true,
	}
	
	// Set default work days (Monday to Friday)
	for i := 1; i <= 5; i++ {
		dv.workDays[time.Weekday(i)] = true
	}
	
	// Add common holidays
	dv.addCommonHolidays()
	
	return dv
}

// addCommonHolidays adds common US holidays
func (dv *DateValidator) addCommonHolidays() {
	year := time.Now().Year()
	
	// New Year's Day
	dv.holidays = append(dv.holidays, time.Date(year, 1, 1, 0, 0, 0, 0, dv.timezone))
	
	// Independence Day
	dv.holidays = append(dv.holidays, time.Date(year, 7, 4, 0, 0, 0, 0, dv.timezone))
	
	// Christmas Day
	dv.holidays = append(dv.holidays, time.Date(year, 12, 25, 0, 0, 0, 0, dv.timezone))
}

// AddHoliday adds a holiday to the validator
func (dv *DateValidator) AddHoliday(date time.Time) {
	dv.holidays = append(dv.holidays, date)
}

// SetTimezone sets the timezone for validation
func (dv *DateValidator) SetTimezone(tz *time.Location) {
	dv.timezone = tz
}

// SetStrictMode sets whether to use strict validation
func (dv *DateValidator) SetStrictMode(strict bool) {
	dv.strictMode = strict
}

// IsWorkDay checks if a date is a work day
func (dv *DateValidator) IsWorkDay(date time.Time) bool {
	// Check if it's a weekend
	if !dv.workDays[date.Weekday()] {
		return false
	}
	
	// Check if it's a holiday
	for _, holiday := range dv.holidays {
		if date.Year() == holiday.Year() && date.Month() == holiday.Month() && date.Day() == holiday.Day() {
			return false
		}
	}
	
	return true
}

// ValidateTaskDates validates a single task's dates
func (dv *DateValidator) ValidateTaskDates(task *Task) []DataValidationError {
	var errors []DataValidationError
	
	if task == nil {
		return errors
	}
	
	now := time.Now()
	
	// Check if start date is valid
	if task.StartDate.IsZero() {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "StartDate",
			Value:     "zero",
			Message:   "Start date is required and cannot be zero",
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Set a valid start date for the task"},
		})
		return errors
	}
	
	// Check if end date is valid
	if task.EndDate.IsZero() {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "EndDate",
			Value:     "zero",
			Message:   "End date is required and cannot be zero",
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Set a valid end date for the task"},
		})
		return errors
	}
	
	// Check if end date is after start date
	if task.EndDate.Before(task.StartDate) {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "EndDate",
			Value:     task.EndDate.Format("2006-01-02"),
			Message:   fmt.Sprintf("End date %s is before start date %s", task.EndDate.Format("2006-01-02"), task.StartDate.Format("2006-01-02")),
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Set end date after start date", "Check if dates were entered in wrong order"},
		})
	}
	
	// Check if task duration is reasonable
	duration := task.GetDuration()
	if duration < 1 {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "Duration",
			Value:     fmt.Sprintf("%d days", duration),
			Message:   "Task duration must be at least 1 day",
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Ensure task has at least 1 day duration"},
		})
	}
	
	// Check for very long tasks (warning)
	if duration > 365 {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "Duration",
			Value:     fmt.Sprintf("%d days", duration),
			Message:   fmt.Sprintf("Task duration is very long (%d days) - consider breaking into smaller tasks", duration),
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{"Consider breaking this task into smaller, manageable pieces"},
		})
	}
	
	// Check if task is in the past (warning)
	if task.EndDate.Before(now) && task.Status != "Completed" {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "EndDate",
			Value:     task.EndDate.Format("2006-01-02"),
			Message:   "Task end date is in the past but status is not completed",
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{"Update task status to 'Completed' or extend the end date"},
		})
	}
	
	// Check if task starts too far in the future (warning)
	futureLimit := now.AddDate(0, 0, 365) // 1 year from now
	if task.StartDate.After(futureLimit) {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "StartDate",
			Value:     task.StartDate.Format("2006-01-02"),
			Message:   "Task start date is more than 1 year in the future",
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{"Consider if this task should be scheduled sooner"},
		})
	}
	
	// Check for weekend-only tasks (info)
	if !dv.IsWorkDay(task.StartDate) && !dv.IsWorkDay(task.EndDate) {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "WorkDays",
			Value:     "weekend",
			Message:   "Task spans only weekend days",
			Severity:  "INFO",
			Timestamp: now,
			Suggestions: []string{"Consider if this task should be scheduled on work days"},
		})
	}
	
	return errors
}

// DetectDateConflicts detects conflicts between tasks
func (dv *DateValidator) DetectDateConflicts(tasks []*Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Sort tasks by start date for efficient conflict detection
	sortedTasks := make([]*Task, len(tasks))
	copy(sortedTasks, tasks)
	sort.Slice(sortedTasks, func(i, j int) bool {
		return sortedTasks[i].StartDate.Before(sortedTasks[j].StartDate)
	})
	
	// Check for overlapping tasks
	for i := 0; i < len(sortedTasks); i++ {
		for j := i + 1; j < len(sortedTasks); j++ {
			task1 := sortedTasks[i]
			task2 := sortedTasks[j]
			
			// Skip if tasks don't overlap
			if !task1.OverlapsWithDateRange(task2.StartDate, task2.EndDate) {
				continue
			}
			
			// Check if tasks are from the same category (potential conflict)
			if task1.Category == task2.Category && task1.Category != "" {
				errors = append(errors, DataValidationError{
					Type:      "CONFLICT",
					TaskID:    task1.ID,
					Field:     "Schedule",
					Value:     fmt.Sprintf("overlaps with %s", task2.ID),
					Message:   fmt.Sprintf("Task %s overlaps with task %s (both %s category)", task1.ID, task2.ID, task1.Category),
					Severity:  "WARNING",
					Timestamp: now,
					Suggestions: []string{
						"Consider rescheduling one of the tasks",
						"Check if tasks can be done in parallel",
						"Verify if this is intentional overlap",
					},
				})
			}
			
			// Check if tasks have the same assignee (potential conflict)
			if task1.Assignee == task2.Assignee && task1.Assignee != "" {
				errors = append(errors, DataValidationError{
					Type:      "CONFLICT",
					TaskID:    task1.ID,
					Field:     "Assignee",
					Value:     fmt.Sprintf("overlaps with %s", task2.ID),
					Message:   fmt.Sprintf("Task %s overlaps with task %s (both assigned to %s)", task1.ID, task2.ID, task1.Assignee),
					Severity:  "WARNING",
					Timestamp: now,
					Suggestions: []string{
						"Check if assignee can handle both tasks simultaneously",
						"Consider reassigning one of the tasks",
						"Verify if this is intentional overlap",
					},
				})
			}
		}
	}
	
	return errors
}

// ValidateDateRanges validates all task date ranges
func (dv *DateValidator) ValidateDateRanges(tasks []*Task) []DataValidationError {
	var errors []DataValidationError
	
	for _, task := range tasks {
		taskErrors := dv.ValidateTaskDates(task)
		errors = append(errors, taskErrors...)
	}
	
	// Detect conflicts between tasks
	conflictErrors := dv.DetectDateConflicts(tasks)
	errors = append(errors, conflictErrors...)
	
	return errors
}

// ValidateWorkDayConstraints validates work day constraints
func (dv *DateValidator) ValidateWorkDayConstraints(tasks []*Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	for _, task := range tasks {
		if task.StartDate.IsZero() || task.EndDate.IsZero() {
			continue
		}
		
		// Check if task starts on a work day
		if !dv.IsWorkDay(task.StartDate) {
			errors = append(errors, DataValidationError{
				Type:      "WORK_DAY",
				TaskID:    task.ID,
				Field:     "StartDate",
				Value:     task.StartDate.Format("2006-01-02"),
				Message:   fmt.Sprintf("Task starts on %s (not a work day)", task.StartDate.Weekday()),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Consider starting the task on a work day",
					"Check if this is intentional weekend work",
				},
			})
		}
		
		// Check if task ends on a work day
		if !dv.IsWorkDay(task.EndDate) {
			errors = append(errors, DataValidationError{
				Type:      "WORK_DAY",
				TaskID:    task.ID,
				Field:     "EndDate",
				Value:     task.EndDate.Format("2006-01-02"),
				Message:   fmt.Sprintf("Task ends on %s (not a work day)", task.EndDate.Weekday()),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Consider ending the task on a work day",
					"Check if this is intentional weekend work",
				},
			})
		}
		
		// Check for tasks that span only non-work days
		workDays := 0
		current := task.StartDate
		for !current.After(task.EndDate) {
			if dv.IsWorkDay(current) {
				workDays++
			}
			current = current.AddDate(0, 0, 1)
		}
		
		if workDays == 0 {
			errors = append(errors, DataValidationError{
				Type:      "WORK_DAY",
				TaskID:    task.ID,
				Field:     "Duration",
				Value:     fmt.Sprintf("%d work days", workDays),
				Message:   "Task spans only non-work days",
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Consider scheduling the task to include work days",
					"Check if this is intentional weekend work",
				},
			})
		}
	}
	
	return errors
}

// GetValidationSummary returns a summary of validation results
func (result *ValidationResult) GetValidationSummary() string {
	if result.IsValid {
		return fmt.Sprintf("Validation passed: %d tasks validated, %d warnings, %d info messages", 
			result.TaskCount, result.WarningCount, len(result.Info))
	}
	
	return fmt.Sprintf("Validation failed: %d tasks validated, %d errors, %d warnings, %d info messages", 
		result.TaskCount, result.ErrorCount, result.WarningCount, len(result.Info))
}

// GetErrorsBySeverity returns errors grouped by severity
func (result *ValidationResult) GetErrorsBySeverity(severity string) []DataValidationError {
	var filtered []DataValidationError
	
	switch severity {
	case "ERROR":
		filtered = result.Errors
	case "WARNING":
		filtered = result.Warnings
	case "INFO":
		filtered = result.Info
	default:
		// Return all errors
		filtered = append(filtered, result.Errors...)
		filtered = append(filtered, result.Warnings...)
		filtered = append(filtered, result.Info...)
	}
	
	return filtered
}

// GetErrorsByType returns errors grouped by type
func (result *ValidationResult) GetErrorsByType(errorType string) []DataValidationError {
	var filtered []DataValidationError
	allErrors := append(result.Errors, result.Warnings...)
	allErrors = append(allErrors, result.Info...)
	
	for _, err := range allErrors {
		if err.Type == errorType {
			filtered = append(filtered, err)
		}
	}
	
	return filtered
}

// HasErrors returns true if there are any errors
func (result *ValidationResult) HasErrors() bool {
	return len(result.Errors) > 0
}

// HasWarnings returns true if there are any warnings
func (result *ValidationResult) HasWarnings() bool {
	return len(result.Warnings) > 0
}

// GetErrorCount returns the total number of errors and warnings
func (result *ValidationResult) GetErrorCount() int {
	return len(result.Errors) + len(result.Warnings) + len(result.Info)
}
