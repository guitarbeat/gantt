package core_test

import (
	"testing"
	"time"

	"phd-dissertation-planner/src/core"
)

func TestDateValidator_ValidateDateRanges(t *testing.T) {
	validator := core.NewDateValidator()

	tests := []struct {
		name     string
		tasks    []*core.Task
		expected int // number of validation errors expected
	}{
		{
			name: "valid date ranges",
			tasks: []*core.Task{
				{
					Name:      "Task 1",
					StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
				},
				{
					Name:      "Task 2", 
					StartDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: 0,
		},
		{
			name: "invalid date range - end before start",
			tasks: []*core.Task{
				{
					Name:      "Invalid Task",
					StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: 1,
		},
		{
			name: "mixed valid and invalid",
			tasks: []*core.Task{
				{
					Name:      "Valid Task",
					StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
				},
				{
					Name:      "Invalid Task",
					StartDate: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateDateRanges(tt.tasks)
			
			if len(errors) != tt.expected {
				t.Errorf("Expected %d validation errors, got %d", tt.expected, len(errors))
			}
			
			// If we expect errors, verify they contain meaningful messages
			if tt.expected > 0 && len(errors) > 0 {
				if errors[0].Message == "" {
					t.Errorf("Expected validation error to have a message")
				}
			}
		})
	}
}

func TestDateValidator_EmptyTaskList(t *testing.T) {
	validator := core.NewDateValidator()
	
	errors := validator.ValidateDateRanges([]*core.Task{})
	
	if len(errors) != 0 {
		t.Errorf("Expected no errors for empty task list, got %d", len(errors))
	}
}

func TestDateValidator_NilTaskList(t *testing.T) {
	validator := core.NewDateValidator()
	
	errors := validator.ValidateDateRanges(nil)
	
	if len(errors) != 0 {
		t.Errorf("Expected no errors for nil task list, got %d", len(errors))
	}
}