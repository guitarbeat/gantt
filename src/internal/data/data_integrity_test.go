package data

import (
	"testing"
	"time"
)

func TestDataIntegrityValidator_ValidateTaskIntegrity(t *testing.T) {
	validator := NewDataIntegrityValidator()
	
	// Test with valid task
	validTask := &Task{
		Name:        "Test Task",
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		Category:    "PROPOSAL",
		Status:      "Planned",
		Priority:    1,
		Description: "Test task description",
	}
	
	errors := validator.ValidateTaskIntegrity(validTask)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid task, got %d errors", len(errors))
	}
	
	// Test with nil task
	errors = validator.ValidateTaskIntegrity(nil)
	if len(errors) == 0 {
		t.Error("Expected errors for nil task")
	}
}

func TestDataIntegrityValidator_ValidateTasksIntegrity(t *testing.T) {
	validator := NewDataIntegrityValidator()
	
	tasks := []*Task{
		{
			Name:      "Task 1",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Task 2",
			StartDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	
	errors := validator.ValidateTasksIntegrity(tasks)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid tasks, got %d errors", len(errors))
	}
}