package data

import (
	"testing"
	"time"
)

func TestCompleteValidationSystem(t *testing.T) {
	// Create a comprehensive validation scenario with various issues
	tasks := []*Task{
		// Valid task
		{
			Name:        "Valid Task",
			StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Category:    "IMAGING",
			Status:      "Planned",
			Priority:    1,
			Assignee:    "John Doe",
			Description: "A valid task for testing",
		},
		// Task with missing required fields
		{
			Name:      "", // Missing Name
			StartDate: time.Time{}, // Missing StartDate
			EndDate:   time.Time{}, // Missing EndDate
		},
		// Task with invalid field formats
		{
			Name:      "Invalid Format Task",
			StartDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:  "INVALID_CATEGORY",
			Status:    "INVALID_STATUS",
			Priority:  10, // Invalid priority
		},
		// Task with data consistency issues
		{
			Name:      "Consistency Issues Task",
			StartDate: time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC), // Same day
			IsMilestone: true,
			ParentID:    "Consistency Issues Task", // Self-parent
		},
		// Task with business rule issues
		{
			Name:        "Dependency Issues Task",
			StartDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "Circular Dependency Task",
			StartDate:   time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 21, 0, 0, 0, 0, time.UTC),
		},
		// Task with business rule issues
		{
			Name:      "Short", // Very short name
			StartDate: time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
			// Missing category, status, priority, description
		},
	}
	
	// Create validators
	dateValidator := NewDateValidator()
	dataIntegrityValidator := NewDataIntegrityValidator()
	
	// Perform comprehensive validation
	dateErrors := dateValidator.ValidateDateRanges(tasks)
	workDayErrors := dateValidator.ValidateWorkDayConstraints(tasks)
	
	// Dependencies removed - no validation needed
	
	dataIntegrityErrors := dataIntegrityValidator.ValidateTasksIntegrity(tasks)
	
	// Create validation reporter
	reporter := NewValidationReporter()
	
	// Add all validation results
	reporter.AddErrors(dateErrors)
	reporter.AddErrors(workDayErrors)
	// Dependencies removed - no errors to add
	reporter.AddErrors(dataIntegrityErrors)
	
	reporter.taskCount = len(tasks)
	
	// Generate comprehensive report
	report := reporter.GenerateReport()
	
	// Verify report structure
	if report == nil {
		t.Error("Expected non-nil report")
	}
	
	if report.TaskCount != 7 {
		t.Errorf("Expected 7 tasks, got %d", report.TaskCount)
	}
	
	if report.IsValid {
		t.Error("Expected validation to fail due to various issues")
	}
	
	// Verify we have different types of errors
	errorTypes := make(map[string]int)
	allErrors := append(report.Errors, report.Warnings...)
	allErrors = append(allErrors, report.Info...)
	
	for _, err := range allErrors {
		errorTypes[err.Type]++
	}
	
	expectedErrorTypes := []string{
		"REQUIRED_FIELD",
		"FIELD_FORMAT", 
		"DATA_CONSISTENCY",
		"DEPENDENCY",
		"CIRCULAR_DEPENDENCY",
		"BUSINESS_RULE",
		"DATE_RANGE",
	}
	
	for _, expectedType := range expectedErrorTypes {
		if errorTypes[expectedType] == 0 {
			t.Errorf("Expected %s errors to be found", expectedType)
		}
	}
	
	// Verify statistics
	if report.Statistics == nil {
		t.Error("Expected non-nil statistics")
	}
	
	if report.Statistics.TotalIssues == 0 {
		t.Error("Expected total issues to be greater than 0")
	}
	
	if report.Statistics.ErrorRate < 0 {
		t.Error("Expected error rate to be non-negative")
	}
	
	if report.Statistics.WarningRate < 0 {
		t.Error("Expected warning rate to be non-negative")
	}
	
	if report.Statistics.InfoRate < 0 {
		t.Error("Expected info rate to be non-negative")
	}
	
	// Verify recommendations
	if len(report.Recommendations) == 0 {
		t.Error("Expected recommendations")
	}
	
	// Generate text report
	textReport := reporter.GenerateTextReport()
	
	if textReport == "" {
		t.Error("Expected non-empty text report")
	}
	
	// Verify text report contains key sections
	requiredSections := []string{
		"VALIDATION REPORT",
		"SUMMARY",
		"STATISTICS",
	}
	
	for _, section := range requiredSections {
		if !containsString(textReport, section) {
			t.Errorf("Expected '%s' section in text report", section)
		}
	}
	
	// Check for optional sections based on content
	if len(report.Errors) > 0 && !containsString(textReport, "ERRORS") {
		t.Error("Expected 'ERRORS' section in text report")
	}
	
	if len(report.Warnings) > 0 && !containsString(textReport, "WARNINGS") {
		t.Error("Expected 'WARNINGS' section in text report")
	}
	
	if len(report.Info) > 0 && !containsString(textReport, "INFO MESSAGES") {
		t.Error("Expected 'INFO MESSAGES' section in text report")
	}
	
	if len(report.Recommendations) > 0 && !containsString(textReport, "RECOMMENDATIONS") {
		t.Error("Expected 'RECOMMENDATIONS' section in text report")
	}
	
	// Generate CSV report
	csvReport := reporter.GenerateCSVReport()
	
	if csvReport == "" {
		t.Error("Expected non-empty CSV report")
	}
	
	// Verify CSV report contains header
			if !containsString(csvReport, "Severity,Type,TaskID,Field,Value,Message,Timestamp") {
		t.Error("Expected CSV header in report")
	}
	
	// Generate JSON report
	jsonReport := reporter.GenerateJSONReport()
	
	if jsonReport == "" {
		t.Error("Expected non-empty JSON report")
	}
	
	// Verify JSON report contains key fields
	expectedJSONFields := []string{
		"\"summary\"",
		"\"is_valid\"",
		"\"task_count\"",
		"\"error_count\"",
		"\"warning_count\"",
		"\"info_count\"",
		"\"statistics\"",
	}
	
	for _, field := range expectedJSONFields {
		if !containsString(jsonReport, field) {
			t.Errorf("Expected '%s' field in JSON report", field)
		}
	}
}

func TestValidationSystemWithRealData(t *testing.T) {
	// Test with real CSV data
	reader := NewReader("input/data.cleaned.fixed.csv")
	tasks, err := reader.ReadTasks()
	
	if err != nil {
		t.Fatalf("Failed to read CSV data: %v", err)
	}
	
	if len(tasks) == 0 {
		t.Fatal("Expected tasks to be loaded")
	}
	
	// Convert []Task to []*Task
	taskPointers := make([]*Task, len(tasks))
	for i := range tasks {
		taskPointers[i] = &tasks[i]
	}
	
	// Create validators
	dateValidator := NewDateValidator()
	dataIntegrityValidator := NewDataIntegrityValidator()
	
	// Perform comprehensive validation
	dateErrors := dateValidator.ValidateDateRanges(taskPointers)
	workDayErrors := dateValidator.ValidateWorkDayConstraints(taskPointers)
	
	// Dependencies removed - no validation needed
	
	dataIntegrityErrors := dataIntegrityValidator.ValidateTasksIntegrity(taskPointers)
	
	// Create validation reporter
	reporter := NewValidationReporter()
	
	// Add all validation results
	reporter.AddErrors(dateErrors)
	reporter.AddErrors(workDayErrors)
	// Dependencies removed - no errors to add
	reporter.AddErrors(dataIntegrityErrors)
	
	reporter.taskCount = len(taskPointers)
	
	// Generate comprehensive report
	report := reporter.GenerateReport()
	
	// Verify report structure
	if report == nil {
		t.Error("Expected non-nil report")
	}
	
	if report.TaskCount != len(taskPointers) {
		t.Errorf("Expected %d tasks, got %d", len(taskPointers), report.TaskCount)
	}
	
	// Verify we have some validation results
	if report.Statistics == nil {
		t.Error("Expected non-nil statistics")
	}
	
	// Generate text report
	textReport := reporter.GenerateTextReport()
	
	if textReport == "" {
		t.Error("Expected non-empty text report")
	}
	
	// Verify text report contains key sections
	expectedSections := []string{
		"VALIDATION REPORT",
		"SUMMARY",
		"STATISTICS",
	}
	
	for _, section := range expectedSections {
		if !containsString(textReport, section) {
			t.Errorf("Expected '%s' section in text report", section)
		}
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
