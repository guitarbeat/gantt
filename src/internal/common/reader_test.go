package common

import (
	"os"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	reader := NewReader("test.csv")
	if reader == nil {
		t.Fatal("Expected reader to be created, got nil")
	}
	
	if reader.filePath != "test.csv" {
		t.Errorf("Expected filePath to be 'test.csv', got %s", reader.filePath)
	}
	
	if reader.strictMode {
		t.Error("Expected strictMode to be false by default")
	}
	
	if !reader.skipInvalid {
		t.Error("Expected skipInvalid to be true by default")
	}
	
	if reader.maxMemoryMB != 100 {
		t.Errorf("Expected maxMemoryMB to be 100, got %d", reader.maxMemoryMB)
	}
}

func TestNewReaderWithOptions(t *testing.T) {
	opts := &ReaderOptions{
		StrictMode:  true,
		SkipInvalid: false,
		MaxMemoryMB: 50,
	}
	
	reader := NewReaderWithOptions("test.csv", opts)
	if reader == nil {
		t.Fatal("Expected reader to be created, got nil")
	}
	
	if reader.strictMode != opts.StrictMode {
		t.Errorf("Expected strictMode to be %v, got %v", opts.StrictMode, reader.strictMode)
	}
	
	if reader.skipInvalid != opts.SkipInvalid {
		t.Errorf("Expected skipInvalid to be %v, got %v", opts.SkipInvalid, reader.skipInvalid)
	}
	
	if reader.maxMemoryMB != opts.MaxMemoryMB {
		t.Errorf("Expected maxMemoryMB to be %d, got %d", opts.MaxMemoryMB, reader.maxMemoryMB)
	}
}

func TestNewReaderWithNilOptions(t *testing.T) {
	reader := NewReaderWithOptions("test.csv", nil)
	if reader == nil {
		t.Fatal("Expected reader to be created, got nil")
	}
	
	// Should use default options
	if reader.strictMode {
		t.Error("Expected strictMode to be false by default")
	}
	
	if !reader.skipInvalid {
		t.Error("Expected skipInvalid to be true by default")
	}
}

func TestParseDate(t *testing.T) {
	reader := NewReader("test.csv")
	
	tests := []struct {
		name     string
		dateStr  string
		expected time.Time
		hasError bool
	}{
		{
			name:     "ISO format",
			dateStr:  "2024-01-15",
			expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name:     "US format",
			dateStr:  "01/15/2024",
			expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name:     "EU format",
			dateStr:  "15/01/2024",
			expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name:     "Empty string",
			dateStr:  "",
			expected: time.Time{},
			hasError: true,
		},
		{
			name:     "Invalid format",
			dateStr:  "invalid-date",
			expected: time.Time{},
			hasError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := reader.parseDate(tt.dateStr)
			
			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsMilestoneTask(t *testing.T) {
	reader := NewReader("test.csv")
	
	tests := []struct {
		name        string
		description string
		expected    bool
	}{
		{
			name:        "Milestone in description",
			description: "MILESTONE: Complete project",
			expected:    true,
		},
		{
			name:        "Milestone in name",
			description: "Complete project",
			expected:    false,
		},
		{
			name:        "Deadline in description",
			description: "DEADLINE: Submit report",
			expected:    true,
		},
		{
			name:        "Due in description",
			description: "DUE: Final review",
			expected:    true,
		},
		{
			name:        "Complete in description",
			description: "COMPLETE: Implementation",
			expected:    true,
		},
		{
			name:        "Finish in description",
			description: "FINISH: Documentation",
			expected:    true,
		},
		{
			name:        "Submit in description",
			description: "SUBMIT: Proposal",
			expected:    true,
		},
		{
			name:        "Deliver in description",
			description: "DELIVER: Results",
			expected:    true,
		},
		{
			name:        "Regular task",
			description: "Write code",
			expected:    false,
		},
		{
			name:        "Empty description",
			description: "",
			expected:    false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reader.isMilestoneTask("Task Name", tt.description)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	reader := NewReader("test.csv")
	
	// Test addError
	reader.addError(nil)
	if reader.hasErrors() {
		t.Error("Expected no errors after adding nil error")
	}
	
	reader.addError(&ParseError{Message: "Test error"})
	if !reader.hasErrors() {
		t.Error("Expected errors after adding error")
	}
	
	// Test clearErrors
	reader.clearErrors()
	if reader.hasErrors() {
		t.Error("Expected no errors after clearing")
	}
	
	// Test getErrorSummary
	reader.addError(&ParseError{Message: "Error 1"})
	reader.addError(&ParseError{Message: "Error 2"})
	
	summary := reader.getErrorSummary()
	if summary == "No errors" {
		t.Error("Expected error summary, got 'No errors'")
	}
}

func TestReadTasksWithValidCSV(t *testing.T) {
	// Create a temporary CSV file
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Test Task 1,Description 1,PROPOSAL,2024-01-15,2024-01-20,1,Planned,John
2,Test Task 2,Description 2,RESEARCH,2024-02-10,2024-02-15,2,In Progress,Jane
`
	
	tmpFile, err := os.CreateTemp("", "test_tasks_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()
	
	reader := NewReader(tmpFile.Name())
	tasks, err := reader.ReadTasks()
	if err != nil {
		t.Fatalf("Failed to read tasks: %v", err)
	}
	
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	
	// Check first task
	if tasks[0].ID != "1" {
		t.Errorf("Expected task ID '1', got %s", tasks[0].ID)
	}
	if tasks[0].Name != "Test Task 1" {
		t.Errorf("Expected task name 'Test Task 1', got %s", tasks[0].Name)
	}
	if tasks[0].Category != "PROPOSAL" {
		t.Errorf("Expected category 'PROPOSAL', got %s", tasks[0].Category)
	}
	if tasks[0].Priority != 1 {
		t.Errorf("Expected priority 1, got %d", tasks[0].Priority)
	}
	if tasks[0].Status != "Planned" {
		t.Errorf("Expected status 'Planned', got %s", tasks[0].Status)
	}
	if tasks[0].Assignee != "John" {
		t.Errorf("Expected assignee 'John', got %s", tasks[0].Assignee)
	}
}

func TestReadTasksWithInvalidCSV(t *testing.T) {
	// Create a temporary CSV file with invalid data
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date
1,Test Task 1,Description 1,PROPOSAL,invalid-date,2024-01-20
2,Test Task 2,Description 2,RESEARCH,2024-02-10,invalid-date
`
	
	tmpFile, err := os.CreateTemp("", "test_tasks_invalid_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()
	
	reader := NewReader(tmpFile.Name())
	tasks, err := reader.ReadTasks()
	
	// Should have errors but continue processing
	if err != nil {
		t.Logf("Expected errors in processing: %v", err)
	}
	
	// Should still return some tasks (those that could be parsed)
	if len(tasks) == 0 {
		t.Error("Expected some tasks to be parsed despite errors")
	}
}

func TestReadTasksWithEmptyFile(t *testing.T) {
	// Create an empty CSV file
	tmpFile, err := os.CreateTemp("", "test_tasks_empty_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	
	reader := NewReader(tmpFile.Name())
	tasks, err := reader.ReadTasks()
	
	if err == nil {
		t.Error("Expected error for empty file, got nil")
	}
	
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

func TestReadTasksWithNonExistentFile(t *testing.T) {
	reader := NewReader("non_existent_file.csv")
	tasks, err := reader.ReadTasks()
	
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

func TestGetDateRange(t *testing.T) {
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date
1,Task 1,Description 1,PROPOSAL,2024-01-15,2024-01-20
2,Task 2,Description 2,RESEARCH,2024-02-10,2024-02-15
3,Task 3,Description 3,ADMIN,2024-03-05,2024-03-10
`
	
	tmpFile, err := os.CreateTemp("", "test_tasks_daterange_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()
	
	reader := NewReader(tmpFile.Name())
	dateRange, err := reader.GetDateRange()
	if err != nil {
		t.Fatalf("Failed to get date range: %v", err)
	}
	
	expectedEarliest := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedLatest := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
	
	if !dateRange.Earliest.Equal(expectedEarliest) {
		t.Errorf("Expected earliest date %v, got %v", expectedEarliest, dateRange.Earliest)
	}
	
	if !dateRange.Latest.Equal(expectedLatest) {
		t.Errorf("Expected latest date %v, got %v", expectedLatest, dateRange.Latest)
	}
}

func TestGetMonthsWithTasks(t *testing.T) {
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date
1,Task 1,Description 1,PROPOSAL,2024-01-15,2024-01-20
2,Task 2,Description 2,RESEARCH,2024-02-10,2024-02-15
3,Task 3,Description 3,ADMIN,2024-01-25,2024-02-05
`
	
	tmpFile, err := os.CreateTemp("", "test_tasks_months_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()
	
	reader := NewReader(tmpFile.Name())
	months, err := reader.GetMonthsWithTasks()
	if err != nil {
		t.Fatalf("Failed to get months with tasks: %v", err)
	}
	
	// Should have January and February
	if len(months) != 2 {
		t.Errorf("Expected 2 months, got %d", len(months))
	}
	
	// Check that months are sorted
	if months[0].Year != 2024 || months[0].Month != time.January {
		t.Errorf("Expected first month to be January 2024, got %v %d", months[0].Month, months[0].Year)
	}
	
	if months[1].Year != 2024 || months[1].Month != time.February {
		t.Errorf("Expected second month to be February 2024, got %v %d", months[1].Month, months[1].Year)
	}
}

func TestValidateCSVFormat(t *testing.T) {
	// Test with valid CSV
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date
1,Task 1,Description 1,PROPOSAL,2024-01-15,2024-01-20
`
	
	tmpFile, err := os.CreateTemp("", "test_tasks_valid_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()
	
	reader := NewReader(tmpFile.Name())
	err = reader.ValidateCSVFormat()
	if err != nil {
		t.Errorf("Expected no error for valid CSV, got %v", err)
	}
}

func TestValidateCSVFormatMissingRequiredFields(t *testing.T) {
	// Test with missing required fields
	csvContent := `Task ID,Description,Category
1,Description 1,PROPOSAL
`
	
	tmpFile, err := os.CreateTemp("", "test_tasks_invalid_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()
	
	reader := NewReader(tmpFile.Name())
	err = reader.ValidateCSVFormat()
	if err == nil {
		t.Error("Expected error for missing required fields, got nil")
	}
}

func TestReadTasksStreaming(t *testing.T) {
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date
1,Task 1,Description 1,PROPOSAL,2024-01-15,2024-01-20
2,Task 2,Description 2,RESEARCH,2024-02-10,2024-02-15
`
	
	tmpFile, err := os.CreateTemp("", "test_tasks_streaming_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()
	
	reader := NewReader(tmpFile.Name())
	
	var processedTasks []Task
	err = reader.ReadTasksStreaming(func(task Task) error {
		processedTasks = append(processedTasks, task)
		return nil
	})
	
	if err != nil {
		t.Fatalf("Failed to read tasks in streaming mode: %v", err)
	}
	
	if len(processedTasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(processedTasks))
	}
}

func TestGetSupportedDateFormats(t *testing.T) {
	formats := GetSupportedDateFormats()
	
	expectedFormats := []string{
		"2006-01-02",
		"01/02/2006",
		"02/01/2006",
		"2006/01/02",
		"02.01.2006",
		"2006-01-02 15:04:05",
	}
	
	if len(formats) != len(expectedFormats) {
		t.Errorf("Expected %d formats, got %d", len(expectedFormats), len(formats))
	}
	
	for i, format := range formats {
		if format != expectedFormats[i] {
			t.Errorf("Expected format %s, got %s", expectedFormats[i], format)
		}
	}
}
