package data

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
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
	
	// Check if start date is before end date
	if task.StartDate.After(task.EndDate) {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "StartDate",
			Value:     task.StartDate.Format("2006-01-02"),
			Message:   fmt.Sprintf("Start date (%s) cannot be after end date (%s)", task.StartDate.Format("2006-01-02"), task.EndDate.Format("2006-01-02")),
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Ensure start date is before end date"},
		})
	}
	
	// Check if task is in the past (warning only)
	if task.EndDate.Before(now) {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "EndDate",
			Value:     task.EndDate.Format("2006-01-02"),
			Message:   fmt.Sprintf("Task ends in the past (%s)", task.EndDate.Format("2006-01-02")),
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{"Consider updating the end date if this is a future task"},
		})
	}
	
	// Check if task starts too far in the future (warning)
	futureLimit := now.AddDate(2, 0, 0) // 2 years from now
	if task.StartDate.After(futureLimit) {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "StartDate",
			Value:     task.StartDate.Format("2006-01-02"),
			Message:   fmt.Sprintf("Task starts very far in the future (%s)", task.StartDate.Format("2006-01-02")),
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{"Verify the start date is correct"},
		})
	}
	
	// Check if task duration is reasonable (warning for very long tasks)
	taskDuration := task.EndDate.Sub(task.StartDate)
	maxDuration := 365 * 24 * time.Hour // 1 year
	if taskDuration > maxDuration {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "Duration",
			Value:     fmt.Sprintf("%.0f days", taskDuration.Hours()/24),
			Message:   fmt.Sprintf("Task duration is very long (%.0f days)", taskDuration.Hours()/24),
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{"Consider breaking this into smaller tasks"},
		})
	}
	
	// Check if task duration is too short (warning for very short tasks)
	minDuration := 24 * time.Hour // 1 day
	if taskDuration < minDuration {
		errors = append(errors, DataValidationError{
			Type:      "DATE_RANGE",
			TaskID:    task.ID,
			Field:     "Duration",
			Value:     fmt.Sprintf("%.0f hours", taskDuration.Hours()),
			Message:   fmt.Sprintf("Task duration is very short (%.0f hours)", taskDuration.Hours()),
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{"Consider if this should be a milestone or combined with another task"},
		})
	}
	
	// Check work day constraints if in strict mode
	if dv.strictMode {
		// Check if start date is a work day
		if !dv.IsWorkDay(task.StartDate) {
			errors = append(errors, DataValidationError{
				Type:      "WORK_DAY",
				TaskID:    task.ID,
				Field:     "StartDate",
				Value:     task.StartDate.Format("2006-01-02"),
				Message:   fmt.Sprintf("Task starts on a non-work day (%s)", task.StartDate.Weekday().String()),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{"Consider moving the start date to a work day"},
			})
		}
		
		// Check if end date is a work day
		if !dv.IsWorkDay(task.EndDate) {
			errors = append(errors, DataValidationError{
				Type:      "WORK_DAY",
				TaskID:    task.ID,
				Field:     "EndDate",
				Value:     task.EndDate.Format("2006-01-02"),
				Message:   fmt.Sprintf("Task ends on a non-work day (%s)", task.EndDate.Weekday().String()),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{"Consider moving the end date to a work day"},
			})
		}
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
					Value:     task1.Assignee,
					Message:   fmt.Sprintf("Task %s overlaps with task %s (same assignee: %s)", task1.ID, task2.ID, task1.Assignee),
					Severity:  "WARNING",
					Timestamp: now,
					Suggestions: []string{
						"Check if assignee can handle both tasks simultaneously",
						"Consider reassigning one of the tasks",
						"Verify if this is intentional parallel work",
					},
				})
			}
		}
	}
	
	return errors
}

// ValidateDateRanges validates date ranges for all tasks
func (dv *DateValidator) ValidateDateRanges(tasks []*Task) []DataValidationError {
	var errors []DataValidationError
	
	// Validate individual task dates
	for _, task := range tasks {
		taskErrors := dv.ValidateTaskDates(task)
		errors = append(errors, taskErrors...)
	}
	
	// Detect conflicts between tasks
	conflictErrors := dv.DetectDateConflicts(tasks)
	errors = append(errors, conflictErrors...)
	
	return errors
}

// ValidateWorkDayConstraints validates work day constraints for all tasks
func (dv *DateValidator) ValidateWorkDayConstraints(tasks []*Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	for _, task := range tasks {
		// Check if start date is a work day
		if !dv.IsWorkDay(task.StartDate) {
			errors = append(errors, DataValidationError{
				Type:      "WORK_DAY",
				TaskID:    task.ID,
				Field:     "StartDate",
				Value:     task.StartDate.Format("2006-01-02"),
				Message:   fmt.Sprintf("Task starts on a non-work day (%s)", task.StartDate.Weekday().String()),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{"Consider moving the start date to a work day"},
			})
		}
		
		// Check if end date is a work day
		if !dv.IsWorkDay(task.EndDate) {
			errors = append(errors, DataValidationError{
				Type:      "WORK_DAY",
				TaskID:    task.ID,
				Field:     "EndDate",
				Value:     task.EndDate.Format("2006-01-02"),
				Message:   fmt.Sprintf("Task ends on a non-work day (%s)", task.EndDate.Weekday().String()),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{"Consider moving the end date to a work day"},
			})
		}
		
		// Check if task spans only non-work days
		workDaysInRange := 0
		for d := task.StartDate; d.Before(task.EndDate) || d.Equal(task.EndDate); d = d.AddDate(0, 0, 1) {
			if dv.IsWorkDay(d) {
				workDaysInRange++
			}
		}
		
		if workDaysInRange == 0 {
			errors = append(errors, DataValidationError{
				Type:      "WORK_DAY",
				TaskID:    task.ID,
				Field:     "Duration",
				Value:     fmt.Sprintf("%d work days", workDaysInRange),
				Message:   "Task spans only non-work days",
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{"Consider adjusting the date range to include work days"},
			})
		}
	}
	
	return errors
}

// GetValidationSummary returns a summary of validation results
func (result *ValidationResult) GetValidationSummary() string {
	if result.IsValid {
		return fmt.Sprintf("Validation passed: %d tasks validated, %d warnings", result.TaskCount, result.WarningCount)
	}
	return fmt.Sprintf("Validation failed: %d tasks validated, %d errors, %d warnings", result.TaskCount, result.ErrorCount, result.WarningCount)
}

// GetErrorsBySeverity returns errors filtered by severity
func (result *ValidationResult) GetErrorsBySeverity(severity string) []DataValidationError {
	var filtered []DataValidationError
	for _, err := range result.Errors {
		if err.Severity == severity {
			filtered = append(filtered, err)
		}
	}
	for _, err := range result.Warnings {
		if err.Severity == severity {
			filtered = append(filtered, err)
		}
	}
	for _, err := range result.Info {
		if err.Severity == severity {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// GetErrorsByType returns errors filtered by type
func (result *ValidationResult) GetErrorsByType(errorType string) []DataValidationError {
	var filtered []DataValidationError
	for _, err := range result.Errors {
		if err.Type == errorType {
			filtered = append(filtered, err)
		}
	}
	for _, err := range result.Warnings {
		if err.Type == errorType {
			filtered = append(filtered, err)
		}
	}
	for _, err := range result.Info {
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

// GetErrorCount returns the total number of errors
func (result *ValidationResult) GetErrorCount() int {
	return len(result.Errors)
}

// DataIntegrityValidator handles data integrity validation
type DataIntegrityValidator struct {
	requiredFields    []string
	fieldValidators   map[string]func(string) bool
	categoryManager   *TaskCategoryManager
	dateValidator     *DateValidator
	dependencyValidator *DependencyValidator
}

// NewDataIntegrityValidator creates a new data integrity validator
func NewDataIntegrityValidator() *DataIntegrityValidator {
	div := &DataIntegrityValidator{
		requiredFields: []string{"ID", "Name", "StartDate", "EndDate"},
		fieldValidators: make(map[string]func(string) bool),
		categoryManager: NewTaskCategoryManager(),
		dateValidator: NewDateValidator(),
		dependencyValidator: NewDependencyValidator(),
	}
	
	// Initialize field validators
	div.initializeFieldValidators()
	
	return div
}

// initializeFieldValidators sets up field validation functions
func (div *DataIntegrityValidator) initializeFieldValidators() {
	// ID validation - alphanumeric, underscore, hyphen allowed
	div.fieldValidators["ID"] = func(value string) bool {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, value)
		return matched && len(value) > 0 && len(value) <= 50
	}
	
	// Name validation - not empty, reasonable length
	div.fieldValidators["Name"] = func(value string) bool {
		trimmed := strings.TrimSpace(value)
		return len(trimmed) > 0 && len(trimmed) <= 200
	}
	
	// Category validation - must be valid category
	div.fieldValidators["Category"] = func(value string) bool {
		if value == "" {
			return true // Category is optional
		}
		_, exists := div.categoryManager.GetCategory(value)
		return exists
	}
	
	// Status validation - must be valid status
	div.fieldValidators["Status"] = func(value string) bool {
		if value == "" {
			return true // Status is optional
		}
		validStatuses := []string{"Planned", "In Progress", "Completed", "On Hold", "Cancelled", "Blocked"}
		for _, status := range validStatuses {
			if value == status {
				return true
			}
		}
		return false
	}
	
	// Priority validation - must be valid priority (int)
	div.fieldValidators["Priority"] = func(value string) bool {
		// This validator is not used for int fields, but kept for consistency
		return true
	}
	
	// Assignee validation - reasonable format
	div.fieldValidators["Assignee"] = func(value string) bool {
		if value == "" {
			return true // Assignee is optional
		}
		trimmed := strings.TrimSpace(value)
		return len(trimmed) > 0 && len(trimmed) <= 100
	}
	
	// Description validation - reasonable length
	div.fieldValidators["Description"] = func(value string) bool {
		if value == "" {
			return true // Description is optional
		}
		return len(value) <= 1000
	}
}

// ValidateTaskIntegrity validates a single task's data integrity
func (div *DataIntegrityValidator) ValidateTaskIntegrity(task *Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	if task == nil {
		errors = append(errors, DataValidationError{
			Type:      "DATA_INTEGRITY",
			TaskID:    "unknown",
			Field:     "Task",
			Value:     "nil",
			Message:   "Task is nil",
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Provide a valid task object"},
		})
		return errors
	}
	
	// Validate required fields
	requiredErrors := div.validateRequiredFields(task)
	errors = append(errors, requiredErrors...)
	
	// Validate field formats
	formatErrors := div.validateFieldFormats(task)
	errors = append(errors, formatErrors...)
	
	// Validate data consistency
	consistencyErrors := div.validateDataConsistency(task)
	errors = append(errors, consistencyErrors...)
	
	// Validate business rules
	businessErrors := div.validateBusinessRules(task)
	errors = append(errors, businessErrors...)
	
	return errors
}

// validateRequiredFields validates that all required fields are present
func (div *DataIntegrityValidator) validateRequiredFields(task *Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Check ID
	if task.ID == "" {
		errors = append(errors, DataValidationError{
			Type:      "REQUIRED_FIELD",
			TaskID:    "unknown",
			Field:     "ID",
			Value:     "",
			Message:   "Task ID is required",
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Provide a unique task ID"},
		})
	}
	
	// Check Name
	if strings.TrimSpace(task.Name) == "" {
		errors = append(errors, DataValidationError{
			Type:      "REQUIRED_FIELD",
			TaskID:    task.ID,
			Field:     "Name",
			Value:     "",
			Message:   "Task name is required",
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Provide a descriptive task name"},
		})
	}
	
	// Check StartDate
	if task.StartDate.IsZero() {
		errors = append(errors, DataValidationError{
			Type:      "REQUIRED_FIELD",
			TaskID:    task.ID,
			Field:     "StartDate",
			Value:     "",
			Message:   "Task start date is required",
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Provide a valid start date for the task"},
		})
	}
	
	// Check EndDate
	if task.EndDate.IsZero() {
		errors = append(errors, DataValidationError{
			Type:      "REQUIRED_FIELD",
			TaskID:    task.ID,
			Field:     "EndDate",
			Value:     "",
			Message:   "Task end date is required",
			Severity:  "ERROR",
			Timestamp: now,
			Suggestions: []string{"Provide a valid end date for the task"},
		})
	}
	
	return errors
}

// validateFieldFormats validates field formats and constraints
func (div *DataIntegrityValidator) validateFieldFormats(task *Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Validate ID format
	if task.ID != "" {
		if !div.fieldValidators["ID"](task.ID) {
			errors = append(errors, DataValidationError{
				Type:      "FIELD_FORMAT",
				TaskID:    task.ID,
				Field:     "ID",
				Value:     task.ID,
				Message:   "Task ID must be alphanumeric with underscores or hyphens, 1-50 characters",
				Severity:  "ERROR",
				Timestamp: now,
				Suggestions: []string{
					"Use only letters, numbers, underscores, and hyphens",
					"Keep ID length between 1 and 50 characters",
				},
			})
		}
	}
	
	// Validate Name format
	if task.Name != "" {
		if !div.fieldValidators["Name"](task.Name) {
			errors = append(errors, DataValidationError{
				Type:      "FIELD_FORMAT",
				TaskID:    task.ID,
				Field:     "Name",
				Value:     task.Name,
				Message:   "Task name must be 1-200 characters and not empty",
				Severity:  "ERROR",
				Timestamp: now,
				Suggestions: []string{
					"Provide a non-empty task name",
					"Keep name length under 200 characters",
				},
			})
		}
	}
	
	// Validate Category format
	if task.Category != "" {
		if !div.fieldValidators["Category"](task.Category) {
			errors = append(errors, DataValidationError{
				Type:      "FIELD_FORMAT",
				TaskID:    task.ID,
				Field:     "Category",
				Value:     task.Category,
				Message:   "Invalid task category",
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Use a valid category: PROPOSAL, LASER, IMAGING, ADMIN, DISSERTATION, RESEARCH, PUBLICATION",
					"Leave category empty for auto-categorization",
				},
			})
		}
	}
	
	// Validate Status format
	if task.Status != "" {
		if !div.fieldValidators["Status"](task.Status) {
			errors = append(errors, DataValidationError{
				Type:      "FIELD_FORMAT",
				TaskID:    task.ID,
				Field:     "Status",
				Value:     task.Status,
				Message:   "Invalid task status",
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Use a valid status: Planned, In Progress, Completed, On Hold, Cancelled, Blocked",
					"Leave status empty for default value",
				},
			})
		}
	}
	
	// Validate Priority format (int)
	if task.Priority < 0 || task.Priority > 5 {
		errors = append(errors, DataValidationError{
			Type:      "FIELD_FORMAT",
			TaskID:    task.ID,
			Field:     "Priority",
			Value:     fmt.Sprintf("%d", task.Priority),
			Message:   "Invalid task priority (must be 0-5)",
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{
				"Use priority values: 0=Low, 1=Medium, 2=High, 3=Critical, 4=Urgent, 5=Emergency",
				"Set priority to 0 for default value",
			},
		})
	}
	
	// Validate Assignee format
	if task.Assignee != "" {
		if !div.fieldValidators["Assignee"](task.Assignee) {
			errors = append(errors, DataValidationError{
				Type:      "FIELD_FORMAT",
				TaskID:    task.ID,
				Field:     "Assignee",
				Value:     task.Assignee,
				Message:   "Invalid assignee format",
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Use a valid assignee name (1-100 characters)",
					"Leave assignee empty if not assigned",
				},
			})
		}
	}
	
	// Validate Description format
	if task.Description != "" {
		if !div.fieldValidators["Description"](task.Description) {
			errors = append(errors, DataValidationError{
				Type:      "FIELD_FORMAT",
				TaskID:    task.ID,
				Field:     "Description",
				Value:     task.Description,
				Message:   "Description too long (max 1000 characters)",
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Keep description under 1000 characters",
					"Consider using a shorter description",
				},
			})
		}
	}
	
	return errors
}

// validateDataConsistency validates data consistency across fields
func (div *DataIntegrityValidator) validateDataConsistency(task *Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Check if task has valid dates for duration calculation
	if !task.StartDate.IsZero() && !task.EndDate.IsZero() {
		duration := task.GetDuration()
		
		// Check for single day duration (same start and end date)
		if duration == 1 && task.StartDate.Equal(task.EndDate) {
			errors = append(errors, DataValidationError{
				Type:      "DATA_CONSISTENCY",
				TaskID:    task.ID,
				Field:     "Duration",
				Value:     "1 day",
				Message:   "Task has single day duration (start and end dates are the same)",
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Consider if this should be a milestone task",
					"Check if end date should be one day after start date",
				},
			})
		}
		
		// Check for negative duration
		if duration < 0 {
			errors = append(errors, DataValidationError{
				Type:      "DATA_CONSISTENCY",
				TaskID:    task.ID,
				Field:     "Duration",
				Value:     fmt.Sprintf("%d days", duration),
				Message:   "Task has negative duration (end date before start date)",
				Severity:  "ERROR",
				Timestamp: now,
				Suggestions: []string{
					"Fix the start and end dates",
					"Check if dates were entered in wrong order",
				},
			})
		}
	}
	
	// Check milestone consistency
	if task.IsMilestone {
		if !task.StartDate.IsZero() && !task.EndDate.IsZero() && !task.StartDate.Equal(task.EndDate) {
			errors = append(errors, DataValidationError{
				Type:      "DATA_CONSISTENCY",
				TaskID:    task.ID,
				Field:     "IsMilestone",
				Value:     "true",
				Message:   "Milestone task should have same start and end date",
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Set end date to same as start date for milestone",
					"Remove milestone flag if this is a duration task",
				},
			})
		}
	}
	
	// Check parent-child consistency
	if task.ParentID != "" {
		if task.ParentID == task.ID {
			errors = append(errors, DataValidationError{
				Type:      "DATA_CONSISTENCY",
				TaskID:    task.ID,
				Field:     "ParentID",
				Value:     task.ParentID,
				Message:   "Task cannot be its own parent",
				Severity:  "ERROR",
				Timestamp: now,
				Suggestions: []string{
					"Remove self-reference from ParentID",
					"Set correct parent task ID",
				},
			})
		}
	}
	
	// Check dependency consistency
	for _, depID := range task.Dependencies {
		if depID == task.ID {
			errors = append(errors, DataValidationError{
				Type:      "DATA_CONSISTENCY",
				TaskID:    task.ID,
				Field:     "Dependencies",
				Value:     depID,
				Message:   "Task cannot depend on itself",
				Severity:  "ERROR",
				Timestamp: now,
				Suggestions: []string{
					"Remove self-dependency from Dependencies",
					"Set correct dependency task ID",
				},
			})
		}
	}
	
	return errors
}

// validateBusinessRules validates business-specific rules
func (div *DataIntegrityValidator) validateBusinessRules(task *Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Check for tasks without category (suggest auto-categorization)
	if task.Category == "" {
		suggestedCategory := div.categoryManager.CategorizeTask(task)
		errors = append(errors, DataValidationError{
			Type:      "BUSINESS_RULE",
			TaskID:    task.ID,
			Field:     "Category",
			Value:     "",
			Message:   fmt.Sprintf("Task has no category, suggested: %s", suggestedCategory),
			Severity:  "INFO",
			Timestamp: now,
			Suggestions: []string{
				fmt.Sprintf("Consider setting category to: %s", suggestedCategory),
				"Leave empty for auto-categorization",
			},
		})
	}
	
	// Check for tasks without status (suggest default)
	if task.Status == "" {
		errors = append(errors, DataValidationError{
			Type:      "BUSINESS_RULE",
			TaskID:    task.ID,
			Field:     "Status",
			Value:     "",
			Message:   "Task has no status, suggested: Planned",
			Severity:  "INFO",
			Timestamp: now,
			Suggestions: []string{
				"Consider setting status to: Planned",
				"Leave empty for default value",
			},
		})
	}
	
	// Check for tasks without priority (suggest default)
	if task.Priority == 0 {
		errors = append(errors, DataValidationError{
			Type:      "BUSINESS_RULE",
			TaskID:    task.ID,
			Field:     "Priority",
			Value:     "0",
			Message:   "Task has default priority, suggested: 1 (Medium)",
			Severity:  "INFO",
			Timestamp: now,
			Suggestions: []string{
				"Consider setting priority to: 1 (Medium)",
				"Keep 0 for default value",
			},
		})
	}
	
	// Check for very short task names
	if len(strings.TrimSpace(task.Name)) < 3 {
		errors = append(errors, DataValidationError{
			Type:      "BUSINESS_RULE",
			TaskID:    task.ID,
			Field:     "Name",
			Value:     task.Name,
			Message:   "Task name is very short, consider a more descriptive name",
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{
				"Use a more descriptive task name",
				"Consider adding more detail to the name",
			},
		})
	}
	
	// Check for tasks with very long names
	if len(task.Name) > 100 {
		errors = append(errors, DataValidationError{
			Type:      "BUSINESS_RULE",
			TaskID:    task.ID,
			Field:     "Name",
			Value:     task.Name,
			Message:   "Task name is very long, consider shortening",
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{
				"Consider shortening the task name",
				"Move detailed information to description",
			},
		})
	}
	
	// Check for tasks without description (suggest adding one)
	if task.Description == "" {
		errors = append(errors, DataValidationError{
			Type:      "BUSINESS_RULE",
			TaskID:    task.ID,
			Field:     "Description",
			Value:     "",
			Message:   "Task has no description, consider adding one for clarity",
			Severity:  "INFO",
			Timestamp: now,
			Suggestions: []string{
				"Consider adding a description for better clarity",
				"Description helps with task understanding and tracking",
			},
		})
	}
	
	return errors
}

// ValidateTasksIntegrity validates multiple tasks for data integrity
func (div *DataIntegrityValidator) ValidateTasksIntegrity(tasks []*Task) []DataValidationError {
	var errors []DataValidationError
	
	// Validate each task individually
	for _, task := range tasks {
		taskErrors := div.ValidateTaskIntegrity(task)
		errors = append(errors, taskErrors...)
	}
	
	// Validate cross-task integrity
	crossTaskErrors := div.validateCrossTaskIntegrity(tasks)
	errors = append(errors, crossTaskErrors...)
	
	return errors
}

// validateCrossTaskIntegrity validates integrity across multiple tasks
func (div *DataIntegrityValidator) validateCrossTaskIntegrity(tasks []*Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Check for duplicate task IDs
	taskIDs := make(map[string]int)
	for _, task := range tasks {
		if task != nil && task.ID != "" {
			taskIDs[task.ID]++
		}
	}
	
	for taskID, count := range taskIDs {
		if count > 1 {
			errors = append(errors, DataValidationError{
				Type:      "CROSS_TASK_INTEGRITY",
				TaskID:    taskID,
				Field:     "ID",
				Value:     taskID,
				Message:   fmt.Sprintf("Duplicate task ID found (%d occurrences)", count),
				Severity:  "ERROR",
				Timestamp: now,
				Suggestions: []string{
					"Ensure all task IDs are unique",
					"Check for duplicate entries in the data source",
				},
			})
		}
	}
	
	// Check for duplicate task names
	taskNames := make(map[string]int)
	for _, task := range tasks {
		if task != nil && task.Name != "" {
			taskNames[task.Name]++
		}
	}
	
	for taskName, count := range taskNames {
		if count > 1 {
			errors = append(errors, DataValidationError{
				Type:      "CROSS_TASK_INTEGRITY",
				TaskID:    "multiple",
				Field:     "Name",
				Value:     taskName,
				Message:   fmt.Sprintf("Duplicate task name found (%d occurrences)", count),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Consider making task names more specific",
					"Check if these are actually different tasks",
				},
			})
		}
	}
	
	return errors
}

// ValidateDataIntegrity performs comprehensive data integrity validation
func (div *DataIntegrityValidator) ValidateDataIntegrity(tasks []*Task) *ValidationResult {
	// Validate individual task integrity
	errors := div.ValidateTasksIntegrity(tasks)
	
	// Add date validation
	dateErrors := div.dateValidator.ValidateDateRanges(tasks)
	errors = append(errors, dateErrors...)
	
	// Add dependency validation
	div.dependencyValidator.AddTasks(tasks)
	dependencyErrors := div.dependencyValidator.ValidateDependencies()
	errors = append(errors, dependencyErrors...)
	
	// Categorize errors by severity
	var errorList, warningList, infoList []DataValidationError
	
	for _, err := range errors {
		switch err.Severity {
		case "ERROR":
			errorList = append(errorList, err)
		case "WARNING":
			warningList = append(warningList, err)
		case "INFO":
			infoList = append(infoList, err)
		}
	}
	
	result := &ValidationResult{
		IsValid:      len(errorList) == 0,
		Errors:       errorList,
		Warnings:     warningList,
		Info:         infoList,
		TaskCount:    len(tasks),
		ErrorCount:   len(errorList),
		WarningCount: len(warningList),
		Timestamp:    time.Now(),
	}
	
	result.Summary = result.GetValidationSummary()
	
	return result
}

// GetFieldValidator returns the validator for a specific field
func (div *DataIntegrityValidator) GetFieldValidator(fieldName string) (func(string) bool, bool) {
	validator, exists := div.fieldValidators[fieldName]
	return validator, exists
}

// AddFieldValidator adds a custom field validator
func (div *DataIntegrityValidator) AddFieldValidator(fieldName string, validator func(string) bool) {
	div.fieldValidators[fieldName] = validator
}

// GetRequiredFields returns the list of required fields
func (div *DataIntegrityValidator) GetRequiredFields() []string {
	return div.requiredFields
}

// SetRequiredFields sets the list of required fields
func (div *DataIntegrityValidator) SetRequiredFields(fields []string) {
	div.requiredFields = fields
}
