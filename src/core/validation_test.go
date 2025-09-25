package core

import (
	"testing"
	"time"
)

func TestDateValidator_ValidateDateRanges(t *testing.T) {
	validator := NewDateValidator()

	tests := []struct {
		name      string
		tasks     []*Task
		hasErrors bool
	}{
		{
			name: "Valid date ranges",
			tasks: []*Task{
				{
					ID:        "TASK-001",
					Name:      "Research",
					StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
					Category:  "Planning",
				},
				{
					ID:        "TASK-002",
					Name:      "Development",
					StartDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
					Category:  "Development",
				},
			},
			hasErrors: false,
		},
		{
			name: "Invalid date ranges",
			tasks: []*Task{
				{
					ID:        "TASK-001",
					Name:      "Invalid Task",
					StartDate: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // End before start
					Category:  "Planning",
				},
			},
			hasErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateDateRanges(tt.tasks)

			if tt.hasErrors {
				if len(errors) == 0 {
					t.Errorf("Expected validation errors for test case '%s', but got none", tt.name)
				}
				return
			}

			if len(errors) > 0 {
				t.Errorf("Unexpected validation errors for test case '%s': %v", tt.name, errors)
			}
		})
	}
}

func TestDataValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      DataValidationError
		expected string
	}{
		{
			name: "Error with task ID",
			err: DataValidationError{
				Type:     "DATE_RANGE",
				TaskID:   "TASK-001",
				Severity: "ERROR",
				Message:  "End date before start date",
			},
			expected: "[ERROR] TASK-001: End date before start date",
		},
		{
			name: "Warning without task ID",
			err: DataValidationError{
				Type:     "OVERLAP",
				Severity: "WARNING",
				Message:  "Tasks may overlap",
			},
			expected: "[WARNING] : Tasks may overlap",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
