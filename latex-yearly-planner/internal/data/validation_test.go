package data

import (
	"testing"
	"time"
)

func TestDateValidator(t *testing.T) {
	dv := NewDateValidator()
	
	// Test work day detection
	monday := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday
	saturday := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC) // Saturday
	
	if !dv.IsWorkDay(monday) {
		t.Error("Expected Monday to be a work day")
	}
	
	if dv.IsWorkDay(saturday) {
		t.Error("Expected Saturday not to be a work day")
	}
	
	// Test holiday detection - add 2024 holiday
	dv.AddHoliday(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) // New Year's Day 2024
	holiday := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // New Year's Day
	if dv.IsWorkDay(holiday) {
		t.Error("Expected New Year's Day not to be a work day")
	}
}

func TestValidateTaskDates(t *testing.T) {
	dv := NewDateValidator()
	
	// Test valid task
	now := time.Now()
	validTask := &Task{
		ID:        "A",
		Name:      "Valid Task",
		StartDate: now.AddDate(0, 0, 1),  // Tomorrow
		EndDate:   now.AddDate(0, 0, 5),  // 5 days from now
		Status:    "Planned",
	}
	
	errors := dv.ValidateTaskDates(validTask)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid task, got %d: %v", len(errors), errors)
	}
	
	// Test task with zero start date
	zeroStartTask := &Task{
		ID:        "B",
		Name:      "Zero Start Task",
		StartDate: time.Time{},
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		Status:    "Planned",
	}
	
	errors = dv.ValidateTaskDates(zeroStartTask)
	if len(errors) == 0 {
		t.Error("Expected errors for task with zero start date")
	}
	
	// Test task with zero end date
	zeroEndTask := &Task{
		ID:        "C",
		Name:      "Zero End Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Time{},
		Status:    "Planned",
	}
	
	errors = dv.ValidateTaskDates(zeroEndTask)
	if len(errors) == 0 {
		t.Error("Expected errors for task with zero end date")
	}
	
	// Test task with end date before start date
	invalidDateTask := &Task{
		ID:        "D",
		Name:      "Invalid Date Task",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Status:    "Planned",
	}
	
	errors = dv.ValidateTaskDates(invalidDateTask)
	if len(errors) == 0 {
		t.Error("Expected errors for task with end date before start date")
	}
	
	// Test very long task (warning)
	longTask := &Task{
		ID:        "E",
		Name:      "Long Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		Status:    "Planned",
	}
	
	longErrors := dv.ValidateTaskDates(longTask)
	warnings := 0
	for _, err := range longErrors {
		if err.Severity == "WARNING" {
			warnings++
		}
	}
	if warnings == 0 {
		t.Error("Expected warnings for very long task")
	}
	
	// Test overdue task (warning)
	overdueTask := &Task{
		ID:        "F",
		Name:      "Overdue Task",
		StartDate: now.AddDate(0, 0, -10),
		EndDate:   now.AddDate(0, 0, -5),
		Status:    "In Progress",
	}
	
	errors = dv.ValidateTaskDates(overdueTask)
	warnings = 0
	for _, err := range errors {
		if err.Severity == "WARNING" {
			warnings++
		}
	}
	if warnings == 0 {
		t.Error("Expected warnings for overdue task")
	}
}

func TestDetectDateConflicts(t *testing.T) {
	dv := NewDateValidator()
	
	// Test overlapping tasks with same category
	tasks := []*Task{
		{
			ID:        "A",
			Name:      "Task A",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Assignee:  "John",
		},
		{
			ID:        "B",
			Name:      "Task B",
			StartDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Assignee:  "Jane",
		},
		{
			ID:        "C",
			Name:      "Task C",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Assignee:  "John",
		},
	}
	
	errors := dv.DetectDateConflicts(tasks)
	if len(errors) == 0 {
		t.Error("Expected conflicts to be detected")
	}
	
	// Check for category conflict
	categoryConflicts := 0
	for _, err := range errors {
		if err.Type == "CONFLICT" && err.Field == "Schedule" {
			categoryConflicts++
		}
	}
	if categoryConflicts == 0 {
		t.Error("Expected category conflicts to be detected")
	}
	
	// Check for assignee conflict
	assigneeConflicts := 0
	for _, err := range errors {
		if err.Type == "CONFLICT" && err.Field == "Assignee" {
			assigneeConflicts++
		}
	}
	if assigneeConflicts == 0 {
		t.Error("Expected assignee conflicts to be detected")
	}
}

func TestValidateWorkDayConstraints(t *testing.T) {
	dv := NewDateValidator()
	
	// Test task starting on weekend
	weekendTask := &Task{
		ID:        "A",
		Name:      "Weekend Task",
		StartDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), // Saturday
		EndDate:   time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
		Status:    "Planned",
	}
	
	tasks := []*Task{weekendTask}
	errors := dv.ValidateWorkDayConstraints(tasks)
	
	warnings := 0
	for _, err := range errors {
		if err.Severity == "WARNING" && err.Type == "WORK_DAY" {
			warnings++
		}
	}
	if warnings == 0 {
		t.Error("Expected warnings for weekend task")
	}
	
	// Test task spanning only non-work days
	nonWorkTask := &Task{
		ID:        "B",
		Name:      "Non-Work Task",
		StartDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), // Saturday
		EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), // Sunday
		Status:    "Planned",
	}
	
	tasks = []*Task{nonWorkTask}
	errors = dv.ValidateWorkDayConstraints(tasks)
	
	warnings = 0
	for _, err := range errors {
		if err.Severity == "WARNING" && err.Type == "WORK_DAY" {
			warnings++
		}
	}
	if warnings == 0 {
		t.Error("Expected warnings for non-work day task")
	}
}

func TestValidationResult(t *testing.T) {
	// Test validation result
	result := &ValidationResult{
		IsValid:      false,
		Errors:       []DataValidationError{},
		Warnings:     []DataValidationError{},
		Info:         []DataValidationError{},
		TaskCount:    10,
		ErrorCount:   2,
		WarningCount: 3,
		Timestamp:    time.Now(),
	}
	
	// Test summary
	summary := result.GetValidationSummary()
	if summary == "" {
		t.Error("Expected non-empty summary")
	}
	
	// Test error count
	if result.GetErrorCount() != 0 {
		t.Errorf("Expected 0 errors, got %d", result.GetErrorCount())
	}
	
	// Test has errors
	if result.HasErrors() {
		t.Error("Expected no errors")
	}
	
	if result.HasWarnings() {
		t.Error("Expected no warnings")
	}
}

func TestDataValidationError(t *testing.T) {
	// Test validation error
	err := DataValidationError{
		Type:        "DATE_RANGE",
		TaskID:      "A",
		Field:       "StartDate",
		Value:       "invalid",
		Message:     "Invalid start date",
		Severity:    "ERROR",
		Timestamp:   time.Now(),
		Suggestions: []string{"Fix the date"},
	}
	
	// Test error message
	msg := err.Error()
	if msg == "" {
		t.Error("Expected non-empty error message")
	}
	
	// Test error type
	if err.Type != "DATE_RANGE" {
		t.Errorf("Expected DATE_RANGE, got %s", err.Type)
	}
	
	// Test severity
	if err.Severity != "ERROR" {
		t.Errorf("Expected ERROR, got %s", err.Severity)
	}
}

func TestDateValidatorIntegration(t *testing.T) {
	dv := NewDateValidator()
	
	// Create a comprehensive test scenario
	tasks := []*Task{
		{
			ID:        "A",
			Name:      "Valid Task",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Assignee:  "John",
			Status:    "Planned",
		},
		{
			ID:        "B",
			Name:      "Overlapping Task",
			StartDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Assignee:  "John",
			Status:    "Planned",
		},
		{
			ID:        "C",
			Name:      "Weekend Task",
			StartDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), // Saturday
			EndDate:   time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
			Category:  "ADMIN",
			Assignee:  "Jane",
			Status:    "Planned",
		},
		{
			ID:        "D",
			Name:      "Invalid Date Task",
			StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Before start
			Category:  "RESEARCH",
			Assignee:  "Bob",
			Status:    "Planned",
		},
	}
	
	// Test date range validation
	dateErrors := dv.ValidateDateRanges(tasks)
	if len(dateErrors) == 0 {
		t.Error("Expected date validation errors")
	}
	
	// Test work day constraints
	workDayErrors := dv.ValidateWorkDayConstraints(tasks)
	if len(workDayErrors) == 0 {
		t.Error("Expected work day constraint errors")
	}
	
	// Test conflict detection
	conflictErrors := dv.DetectDateConflicts(tasks)
	if len(conflictErrors) == 0 {
		t.Error("Expected conflict detection errors")
	}
	
	// Verify error types
	errorTypes := make(map[string]int)
	for _, err := range dateErrors {
		errorTypes[err.Type]++
	}
	
	if errorTypes["DATE_RANGE"] == 0 {
		t.Error("Expected DATE_RANGE errors")
	}
	
	if errorTypes["CONFLICT"] == 0 {
		t.Error("Expected CONFLICT errors")
	}
}
