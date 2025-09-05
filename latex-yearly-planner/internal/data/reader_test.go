package data

import (
	"fmt"
	"os"
	"strings"
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
		StrictMode:            true,
		SkipInvalid:           false,
		MaxMemoryMB:           50,
		Logger:                nil, // Will use default
		ValidateDependencies:  true,
		DetectCircularDeps:    true,
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
	if reader.validateDependencies != true {
		t.Error("Expected validateDependencies to be true")
	}
	if reader.detectCircularDeps != true {
		t.Error("Expected detectCircularDeps to be true")
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

func TestParseDependencies(t *testing.T) {
	reader := NewReader("test.csv")

	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"A", []string{"A"}},
		{"A,B", []string{"A", "B"}},
		{"A, B, C", []string{"A", "B", "C"}},
		{"A,,B", []string{"A", "B"}},
		{" A , B ", []string{"A", "B"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := reader.parseDependencies(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d dependencies, got %d", len(tt.expected), len(result))
			}
			for i, dep := range result {
				if dep != tt.expected[i] {
					t.Errorf("Expected dependency %s, got %s", tt.expected[i], dep)
				}
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
		{"Regular Task", "Do something", false},
		{"MILESTONE: Complete", "Finish the project", true},
		{"Submit Proposal", "Submit the final proposal", true},
		{"Deadline Task", "Meet the deadline", true},
		{"Deliver Results", "Deliver final results", true},
		{"Complete Analysis", "Complete the analysis", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reader.isMilestoneTask(tt.name, tt.description)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for name='%s', description='%s'", tt.expected, result, tt.name, tt.description)
			}
		})
	}
}

func TestReadTasksWithDependencies(t *testing.T) {
	// Create a CSV file with dependencies
	csvContent := `Task ID,Task Name,Parent Task ID,Category,Start Date,Due Date,Dependencies,Description,Priority,Status,Assignee
A,Task A,,PROPOSAL,2024-01-15,2024-01-20,,Description A,1,Planned,John
B,Task B,,PROPOSAL,2024-01-21,2024-01-25,A,Description B,2,In Progress,Jane
C,Task C,A,PROPOSAL,2024-01-26,2024-01-30,"A,B",Description C,3,Planned,Bob
`

	tmpFile, err := os.CreateTemp("", "test_dependencies_*.csv")
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

	// Test task A (no dependencies)
	taskA := tasks[0]
	if taskA.ID != "A" {
		t.Errorf("Expected task A, got %s", taskA.ID)
	}
	if len(taskA.Dependencies) != 0 {
		t.Errorf("Expected 0 dependencies for task A, got %d", len(taskA.Dependencies))
	}
	if taskA.ParentID != "" {
		t.Errorf("Expected empty ParentID for task A, got %s", taskA.ParentID)
	}

	// Test task B (depends on A)
	taskB := tasks[1]
	if taskB.ID != "B" {
		t.Errorf("Expected task B, got %s", taskB.ID)
	}
	if len(taskB.Dependencies) != 1 || taskB.Dependencies[0] != "A" {
		t.Errorf("Expected task B to depend on A, got %v", taskB.Dependencies)
	}

	// Test task C (depends on A and B, parent is A)
	taskC := tasks[2]
	if taskC.ID != "C" {
		t.Errorf("Expected task C, got %s", taskC.ID)
	}
	if len(taskC.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies for task C, got %d", len(taskC.Dependencies))
	}
	if taskC.ParentID != "A" {
		t.Errorf("Expected ParentID A for task C, got %s", taskC.ParentID)
	}
}

func TestValidateTaskDependencies(t *testing.T) {
	reader := NewReader("test.csv")

	// Test valid dependencies
	validTasks := []Task{
		{ID: "A", Dependencies: []string{}},
		{ID: "B", Dependencies: []string{"A"}},
		{ID: "C", Dependencies: []string{"A", "B"}},
	}

	err := reader.validateTaskDependencies(validTasks)
	if err != nil {
		t.Errorf("Expected no error for valid dependencies, got %v", err)
	}

	// Test invalid dependencies
	invalidTasks := []Task{
		{ID: "A", Dependencies: []string{}},
		{ID: "B", Dependencies: []string{"X"}}, // X doesn't exist
	}

	err = reader.validateTaskDependencies(invalidTasks)
	if err == nil {
		t.Error("Expected error for invalid dependencies, got nil")
	}
}

func TestDetectCircularDependencies(t *testing.T) {
	reader := NewReader("test.csv")

	// Test no circular dependencies
	validTasks := []Task{
		{ID: "A", Dependencies: []string{}},
		{ID: "B", Dependencies: []string{"A"}},
		{ID: "C", Dependencies: []string{"B"}},
	}

	err := reader.detectCircularDependencies(validTasks)
	if err != nil {
		t.Errorf("Expected no error for valid dependencies, got %v", err)
	}

	// Test circular dependencies
	circularTasks := []Task{
		{ID: "A", Dependencies: []string{"B"}},
		{ID: "B", Dependencies: []string{"A"}},
	}

	err = reader.detectCircularDependencies(circularTasks)
	if err == nil {
		t.Error("Expected error for circular dependencies, got nil")
	}
	
	// Test that the error is a CircularDependencyError
	if _, ok := err.(*CircularDependencyError); !ok {
		t.Errorf("Expected CircularDependencyError, got %T", err)
	}
}

func TestErrorHandling(t *testing.T) {
	reader := NewReader("test.csv")

	// Test error collection
	reader.addError(fmt.Errorf("test error 1"))
	reader.addError(fmt.Errorf("test error 2"))

	if !reader.hasErrors() {
		t.Error("Expected hasErrors to return true")
	}

	errors := reader.getErrors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	// Test error summary
	summary := reader.getErrorSummary()
	if !strings.Contains(summary, "Found 2 errors") {
		t.Errorf("Expected error summary to contain 'Found 2 errors', got: %s", summary)
	}

	// Test clear errors
	reader.clearErrors()
	if reader.hasErrors() {
		t.Error("Expected hasErrors to return false after clearErrors")
	}
}

func TestParseError(t *testing.T) {
	parseErr := &ParseError{
		Row:     5,
		Column:  "Start Date",
		Value:   "invalid-date",
		Message: "invalid format",
	}

	expected := "row 5, column 'Start Date', value 'invalid-date': invalid format"
	if parseErr.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, parseErr.Error())
	}
}

func TestValidationError(t *testing.T) {
	validationErr := &ValidationError{
		TaskID:  "A",
		Field:   "Dependencies",
		Value:   "X",
		Message: "references non-existent task",
	}

	expected := "task A, field 'Dependencies', value 'X': references non-existent task"
	if validationErr.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, validationErr.Error())
	}
}

func TestCircularDependencyError(t *testing.T) {
	circularErr := &CircularDependencyError{
		Cycle: []string{"A", "B", "C", "A"},
	}

	expected := "circular dependency detected: A -> B -> C -> A"
	if circularErr.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, circularErr.Error())
	}
}

func TestReadTasksWithDetailedErrors(t *testing.T) {
	// Create a CSV file with various errors (but no invalid dependencies or circular deps to avoid validation failure)
	csvContent := `Task ID,Task Name,Parent Task ID,Category,Start Date,Due Date,Dependencies,Description,Priority,Status,Assignee
A,Valid Task,,PROPOSAL,2024-01-15,2024-01-20,,Description A,1,Planned,John
B,Invalid Date Task,,PROPOSAL,invalid-date,2024-01-25,,Description B,2,In Progress,Jane
C,Valid Deps Task,,PROPOSAL,2024-01-26,2024-01-30,A,Description C,3,Planned,Bob
D,Valid Task 2,,PROPOSAL,2024-02-01,2024-02-05,,Description D,4,Planned,Alice
E,Valid Task 3,,PROPOSAL,2024-02-06,2024-02-10,,Description E,5,Planned,Charlie
`

	tmpFile, err := os.CreateTemp("", "test_detailed_errors_*.csv")
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
	
	// Should have some tasks but also errors
	if err != nil {
		t.Fatalf("ReadTasks error: %v", err)
	}

	// Should have at least one valid task
	if len(tasks) < 1 {
		t.Errorf("Expected at least 1 valid task, got %d", len(tasks))
	}

	// Should have collected errors
	if !reader.hasErrors() {
		t.Error("Expected reader to have collected errors")
	}

	errors := reader.getErrors()
	if len(errors) == 0 {
		t.Error("Expected collected errors, got none")
	}

	// Test error summary
	summary := reader.getErrorSummary()
	if !strings.Contains(summary, "Found") {
		t.Errorf("Expected error summary to contain 'Found', got: %s", summary)
	}
}

func TestReadTasksWithInvalidDependencies(t *testing.T) {
	// Create a CSV file with invalid dependencies to test dependency validation
	csvContent := `Task ID,Task Name,Parent Task ID,Category,Start Date,Due Date,Dependencies,Description,Priority,Status,Assignee
A,Valid Task,,PROPOSAL,2024-01-15,2024-01-20,,Description A,1,Planned,John
B,Invalid Deps Task,,PROPOSAL,2024-01-26,2024-01-30,X,Description B,2,Planned,Bob
`

	tmpFile, err := os.CreateTemp("", "test_invalid_deps_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	reader := NewReader(tmpFile.Name())
	_, err = reader.ReadTasks()
	
	// Should fail due to invalid dependencies
	if err == nil {
		t.Error("Expected ReadTasks to fail due to invalid dependencies")
	}

	// Check that the error is about dependency validation
	if !strings.Contains(err.Error(), "dependency validation failed") {
		t.Errorf("Expected dependency validation error, got: %v", err)
	}
}
