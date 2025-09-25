package core

import (
	"os"
	"testing"
	"time"
)

func TestReader_ParseDate(t *testing.T) {
	reader := NewReader("dummy.csv") // We won't actually read this file

	tests := []struct {
		name     string
		input    string
		expected time.Time
		hasError bool
	}{
		{
			name:     "ISO format",
			input:    "2024-01-15",
			expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name:     "US format",
			input:    "01/15/2024",
			expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name:     "Invalid format",
			input:    "invalid-date",
			expected: time.Time{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := reader.parseDate(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input '%s', but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tt.input, err)
				return
			}

			if !result.Equal(tt.expected) {
				t.Errorf("For input '%s', expected %v, but got %v", tt.input, tt.expected, result)
			}
		})
	}
}

func TestNewReader(t *testing.T) {
	// Test creating a reader (we can't easily test file reading without creating test files)
	reader := NewReader("dummy.csv")

	if reader == nil {
		t.Error("NewReader returned nil")
	}

	if reader.filePath != "dummy.csv" {
		t.Errorf("Expected filePath 'dummy.csv', got '%s'", reader.filePath)
	}
}

func TestReader_ReadTasks(t *testing.T) {
	// Create a temporary CSV file for testing using the actual format
	csvContent := `Phase,Sub-Phase,Task ID,Dependencies,Task,Start Date,End Date,Objective,Milestone,Status
1,PhD Proposal,TASK-001,,Research Proposal,2024-01-01,2024-01-31,Write research proposal,false,Not Started
1,Literature Review,TASK-002,TASK-001,Literature Review,2024-02-01,2024-02-28,Review existing literature,false,In Progress`

	// Create temporary file
	tempFile := "/tmp/test_tasks.csv"
	err := os.WriteFile(tempFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	reader := NewReader(tempFile)
	tasks, err := reader.ReadTasks()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Check first task
	if tasks[0].ID != "TASK-001" {
		t.Errorf("Expected task ID 'TASK-001', got '%s'", tasks[0].ID)
	}

	if tasks[0].Name != "Research Proposal" {
		t.Errorf("Expected task name 'Research Proposal', got '%s'", tasks[0].Name)
	}

	// Check dates are parsed correctly
	expectedStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !tasks[0].StartDate.Equal(expectedStart) {
		t.Errorf("Expected start date %v, got %v", expectedStart, tasks[0].StartDate)
	}

	// Check dependencies
	if len(tasks[1].Dependencies) != 1 || tasks[1].Dependencies[0] != "TASK-001" {
		t.Errorf("Expected TASK-002 to depend on TASK-001, got %v", tasks[1].Dependencies)
	}
}

func TestParseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      ParseError
		expected string
	}{
		{
			name: "With row number",
			err: ParseError{
				Row:     5,
				Column:  "Start Date",
				Value:   "invalid-date",
				Message: "invalid date format",
			},
			expected: "row 5, column 'Start Date', value 'invalid-date': invalid date format",
		},
		{
			name: "Without row number",
			err: ParseError{
				Column:  "Category",
				Value:   "",
				Message: "category is required",
			},
			expected: "column 'Category', value '': category is required",
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
