package data

import (
	"testing"
	"time"
)

func TestDateValidator_ValidateTaskDates(t *testing.T) {
	validator := NewDateValidator()
	
	// Test valid task
	validTask := &Task{
		Name:      "Valid Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	errors := validator.ValidateTaskDates(validTask)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid task, got %d errors", len(errors))
	}
	
	// Test task with zero dates
	invalidTask := &Task{
		Name:      "Invalid Task",
		StartDate: time.Time{},
		EndDate:   time.Time{},
	}
	
	errors = validator.ValidateTaskDates(invalidTask)
	if len(errors) == 0 {
		t.Error("Expected errors for task with zero dates")
	}
	
	// Test task with start date after end date
	invalidTask2 := &Task{
		Name:      "Invalid Task 2",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	
	errors = validator.ValidateTaskDates(invalidTask2)
	if len(errors) == 0 {
		t.Error("Expected errors for task with start date after end date")
	}
}

func TestDateValidator_DetectDateConflicts(t *testing.T) {
	validator := NewDateValidator()
	
	tasks := []*Task{
		{
			Name:      "Task 1",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
		},
		{
			Name:      "Task 2",
			StartDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
		},
	}
	
	errors := validator.DetectDateConflicts(tasks)
	if len(errors) == 0 {
		t.Error("Expected conflicts between overlapping tasks")
	}
}

func TestDateValidator_IsWorkDay(t *testing.T) {
	validator := NewDateValidator()
	
	// Test work day (Monday)
	workDay := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday
	if !validator.IsWorkDay(workDay) {
		t.Error("Expected Monday to be a work day")
	}
	
	// Test weekend (Saturday)
	weekend := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC) // Saturday
	if validator.IsWorkDay(weekend) {
		t.Error("Expected Saturday to not be a work day")
	}
}

func TestDateValidator_ValidateDateRanges(t *testing.T) {
	validator := NewDateValidator()
	
	tasks := []*Task{
		{
			Name:      "Valid Task",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Invalid Task",
			StartDate: time.Time{},
			EndDate:   time.Time{},
		},
	}
	
	errors := validator.ValidateDateRanges(tasks)
	if len(errors) == 0 {
		t.Error("Expected errors for invalid tasks")
	}
}

func TestDateValidator_ValidateWorkDayConstraints(t *testing.T) {
	validator := NewDateValidator()
	
	tasks := []*Task{
		{
			Name:      "Weekend Task",
			StartDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), // Saturday
			EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), // Sunday
		},
	}
	
	errors := validator.ValidateWorkDayConstraints(tasks)
	if len(errors) == 0 {
		t.Error("Expected warnings for weekend tasks")
	}
}