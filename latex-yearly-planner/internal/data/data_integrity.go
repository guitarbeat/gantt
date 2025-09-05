package data

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

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
