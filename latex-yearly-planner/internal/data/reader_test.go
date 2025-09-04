package data

import (
	"os"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	reader := NewReader("test.csv")
	if reader == nil {
		t.Fatal("NewReader returned nil")
	}
	if reader.filePath != "test.csv" {
		t.Errorf("Expected filePath 'test.csv', got %s", reader.filePath)
	}
	if reader.logger == nil {
		t.Error("Expected logger to be initialized")
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
		Logger:      nil, // Will use default
	}

	reader := NewReaderWithOptions("test.csv", opts)
	if reader == nil {
		t.Fatal("NewReaderWithOptions returned nil")
	}
	if reader.strictMode != true {
		t.Error("Expected strictMode to be true")
	}
	if reader.skipInvalid != false {
		t.Error("Expected skipInvalid to be false")
	}
	if reader.maxMemoryMB != 50 {
		t.Errorf("Expected maxMemoryMB to be 50, got %d", reader.maxMemoryMB)
	}
}

func TestParseDate(t *testing.T) {
	reader := NewReader("test.csv")

	tests := []struct {
		input    string
		expected time.Time
		hasError bool
	}{
		{"2024-01-15", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"01/15/2024", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"15/01/2024", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"2024/01/15", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"15.01.2024", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"2024-01-15 10:30:00", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), false},
		{"", time.Time{}, true},
		{"invalid-date", time.Time{}, true},
		{" 2024-01-15 ", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), false}, // With spaces
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := reader.parseDate(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input '%s': %v", tt.input, err)
				}
				if !result.Equal(tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestReadTasks(t *testing.T) {
	// Create a test CSV file
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Test Task 1,Description 1,PROPOSAL,2024-01-15,2024-01-20,1,Planned,John
2,Test Task 2,Description 2,RESEARCH,2024-02-10,2024-02-15,2,In Progress,Jane
3,Test Task 3,Description 3,ADMIN,2025-03-05,2025-03-10,3,Completed,Bob
`

	tmpFile, err := os.CreateTemp("", "test_tasks_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	reader := NewReader(tmpFile.Name())
	tasks, err := reader.ReadTasks()
	if err != nil {
		t.Fatalf("ReadTasks error: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}

	// Test first task
	task1 := tasks[0]
	if task1.ID != "1" {
		t.Errorf("Expected ID '1', got %s", task1.ID)
	}
	if task1.Name != "Test Task 1" {
		t.Errorf("Expected Name 'Test Task 1', got %s", task1.Name)
	}
	if task1.Category != "PROPOSAL" {
		t.Errorf("Expected Category 'PROPOSAL', got %s", task1.Category)
	}
	if task1.Priority != 1 {
		t.Errorf("Expected Priority 1, got %d", task1.Priority)
	}
	if task1.Status != "Planned" {
		t.Errorf("Expected Status 'Planned', got %s", task1.Status)
	}
	if task1.Assignee != "John" {
		t.Errorf("Expected Assignee 'John', got %s", task1.Assignee)
	}
	if !task1.StartDate.Equal(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected StartDate 2024-01-15, got %v", task1.StartDate)
	}
	if !task1.EndDate.Equal(time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected EndDate 2024-01-20, got %v", task1.EndDate)
	}
}

func TestReadTasksWithInvalidData(t *testing.T) {
	// Create a CSV file with invalid data
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Valid Task,Description,PROPOSAL,2024-01-15,2024-01-20,1,Planned,John
2,Invalid Task,Description,RESEARCH,invalid-date,2024-02-15,2,In Progress,Jane
3,Another Valid Task,Description,ADMIN,2025-03-05,2025-03-10,3,Completed,Bob
`

	tmpFile, err := os.CreateTemp("", "test_invalid_tasks_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	reader := NewReader(tmpFile.Name())
	tasks, err := reader.ReadTasks()
	if err != nil {
		t.Fatalf("ReadTasks error: %v", err)
	}

	// Should have 2 valid tasks (invalid one skipped)
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks (1 invalid skipped), got %d", len(tasks))
	}
}

func TestReadTasksStrictMode(t *testing.T) {
	// Create a CSV file with invalid data
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Valid Task,Description,PROPOSAL,2024-01-15,2024-01-20,1,Planned,John
2,Invalid Task,Description,RESEARCH,invalid-date,2024-02-15,2,In Progress,Jane
`

	tmpFile, err := os.CreateTemp("", "test_strict_tasks_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	opts := &ReaderOptions{
		StrictMode:  true,
		SkipInvalid: false,
		MaxMemoryMB: 100,
		Logger:      nil,
	}

	reader := NewReaderWithOptions(tmpFile.Name(), opts)
	_, err = reader.ReadTasks()
	if err == nil {
		t.Error("Expected error in strict mode with invalid data, got nil")
	}
}

func TestGetDateRange(t *testing.T) {
	// Create a test CSV file
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Task 1,Description,PROPOSAL,2024-01-15,2024-01-20,1,Planned,John
2,Task 2,Description,RESEARCH,2024-02-10,2024-02-15,2,In Progress,Jane
3,Task 3,Description,ADMIN,2025-03-05,2025-03-10,3,Completed,Bob
`

	tmpFile, err := os.CreateTemp("", "test_date_range_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	reader := NewReader(tmpFile.Name())
	dateRange, err := reader.GetDateRange()
	if err != nil {
		t.Fatalf("GetDateRange error: %v", err)
	}

	expectedEarliest := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedLatest := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)

	if !dateRange.Earliest.Equal(expectedEarliest) {
		t.Errorf("Expected earliest date %v, got %v", expectedEarliest, dateRange.Earliest)
	}
	if !dateRange.Latest.Equal(expectedLatest) {
		t.Errorf("Expected latest date %v, got %v", expectedLatest, dateRange.Latest)
	}
}

func TestGetMonthsWithTasks(t *testing.T) {
	// Create a test CSV file
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Task 1,Description,PROPOSAL,2024-01-15,2024-01-20,1,Planned,John
2,Task 2,Description,RESEARCH,2024-02-10,2024-02-15,2,In Progress,Jane
3,Task 3,Description,ADMIN,2025-03-05,2025-03-10,3,Completed,Bob
`

	tmpFile, err := os.CreateTemp("", "test_months_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	reader := NewReader(tmpFile.Name())
	months, err := reader.GetMonthsWithTasks()
	if err != nil {
		t.Fatalf("GetMonthsWithTasks error: %v", err)
	}

	// Should have 3 months: Jan 2024, Feb 2024, Mar 2025
	if len(months) != 3 {
		t.Errorf("Expected 3 months, got %d", len(months))
	}

	// Check specific months
	expectedMonths := []MonthYear{
		{Year: 2024, Month: time.January},
		{Year: 2024, Month: time.February},
		{Year: 2025, Month: time.March},
	}

	for i, expected := range expectedMonths {
		if months[i].Year != expected.Year {
			t.Errorf("Expected year %d, got %d", expected.Year, months[i].Year)
		}
		if months[i].Month != expected.Month {
			t.Errorf("Expected month %v, got %v", expected.Month, months[i].Month)
		}
	}
}

func TestValidateCSVFormat(t *testing.T) {
	// Test valid CSV format
	validCSV := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Task 1,Description,PROPOSAL,2024-01-15,2024-01-20,1,Planned,John
`

	tmpFile, err := os.CreateTemp("", "test_valid_format_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(validCSV); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	reader := NewReader(tmpFile.Name())
	err = reader.ValidateCSVFormat()
	if err != nil {
		t.Errorf("Expected no error for valid CSV format, got %v", err)
	}

	// Test invalid CSV format (missing required fields)
	invalidCSV := `Task ID,Description,Category
1,Description,PROPOSAL
`

	tmpFile2, err := os.CreateTemp("", "test_invalid_format_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile2.Name())

	if _, err := tmpFile2.WriteString(invalidCSV); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile2.Close()

	reader2 := NewReader(tmpFile2.Name())
	err = reader2.ValidateCSVFormat()
	if err == nil {
		t.Error("Expected error for invalid CSV format, got nil")
	}
}

func TestGetSupportedDateFormats(t *testing.T) {
	formats := GetSupportedDateFormats()
	if len(formats) == 0 {
		t.Error("Expected supported date formats, got empty slice")
	}

	// Check that ISO format is included
	hasISO := false
	for _, format := range formats {
		if format == DateFormatISO {
			hasISO = true
			break
		}
	}
	if !hasISO {
		t.Error("Expected ISO date format to be supported")
	}
}

func TestReadTasksStreaming(t *testing.T) {
	// Create a test CSV file
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Task 1,Description,PROPOSAL,2024-01-15,2024-01-20,1,Planned,John
2,Task 2,Description,RESEARCH,2024-02-10,2024-02-15,2,In Progress,Jane
3,Task 3,Description,ADMIN,2025-03-05,2025-03-10,3,Completed,Bob
`

	tmpFile, err := os.CreateTemp("", "test_streaming_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
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
		t.Fatalf("ReadTasksStreaming error: %v", err)
	}

	if len(processedTasks) != 3 {
		t.Errorf("Expected 3 processed tasks, got %d", len(processedTasks))
	}
}

func TestTaskValidation(t *testing.T) {
	reader := NewReader("test.csv")

	// Test end date before start date validation
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date,Priority,Status,Assignee
1,Invalid Task,Description,PROPOSAL,2024-01-20,2024-01-15,1,Planned,John
`

	tmpFile, err := os.CreateTemp("", "test_invalid_dates_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	reader.filePath = tmpFile.Name()
	tasks, err := reader.ReadTasks()
	if err != nil {
		t.Fatalf("ReadTasks error: %v", err)
	}

	// Should skip the invalid task
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks (invalid date range), got %d", len(tasks))
	}
}
