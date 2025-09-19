package data

import (
	"strings"
	"testing"
	"time"
)

func TestValidationReporter(t *testing.T) {
	vr := NewValidationReporter()
	
	// Test initial state
	if vr.GetErrorCount() != 0 {
		t.Error("Expected 0 errors initially")
	}
	if vr.GetWarningCount() != 0 {
		t.Error("Expected 0 warnings initially")
	}
	if vr.GetInfoCount() != 0 {
		t.Error("Expected 0 info messages initially")
	}
	if vr.HasErrors() {
		t.Error("Expected no errors initially")
	}
	if vr.HasWarnings() {
		t.Error("Expected no warnings initially")
	}
	if vr.HasIssues() {
		t.Error("Expected no issues initially")
	}
}

func TestAddValidationResult(t *testing.T) {
	vr := NewValidationReporter()
	
	// Create a validation result
	result := &ValidationResult{
		IsValid:      false,
		Errors:       []DataValidationError{
			{
				Type:      "REQUIRED_FIELD",
				TaskID:    "A",
				Field:     "ID",
				Value:     "",
				Message:   "Task ID is required",
				Severity:  "ERROR",
				Timestamp: time.Now(),
			},
		},
		Warnings:     []DataValidationError{
			{
				Type:      "FIELD_FORMAT",
				TaskID:    "B",
				Field:     "Category",
				Value:     "INVALID",
				Message:   "Invalid category",
				Severity:  "WARNING",
				Timestamp: time.Now(),
			},
		},
		Info:         []DataValidationError{
			{
				Type:      "BUSINESS_RULE",
				TaskID:    "C",
				Field:     "Description",
				Value:     "",
				Message:   "Consider adding description",
				Severity:  "INFO",
				Timestamp: time.Now(),
			},
		},
		TaskCount:    3,
		ErrorCount:   1,
		WarningCount: 1,
		Timestamp:    time.Now(),
	}
	
	vr.AddValidationResult(result)
	
	if vr.GetErrorCount() != 1 {
		t.Errorf("Expected 1 error, got %d", vr.GetErrorCount())
	}
	if vr.GetWarningCount() != 1 {
		t.Errorf("Expected 1 warning, got %d", vr.GetWarningCount())
	}
	if vr.GetInfoCount() != 1 {
		t.Errorf("Expected 1 info message, got %d", vr.GetInfoCount())
	}
	if vr.taskCount != 3 {
		t.Errorf("Expected 3 tasks, got %d", vr.taskCount)
	}
}

func TestAddIndividualErrors(t *testing.T) {
	vr := NewValidationReporter()
	
	// Add individual errors
	errors := []DataValidationError{
		{
			Type:      "REQUIRED_FIELD",
			TaskID:    "A",
			Field:     "ID",
			Value:     "",
			Message:   "Task ID is required",
			Severity:  "ERROR",
			Timestamp: time.Now(),
		},
		{
			Type:      "FIELD_FORMAT",
			TaskID:    "B",
			Field:     "Name",
			Value:     "",
			Message:   "Task name is required",
			Severity:  "ERROR",
			Timestamp: time.Now(),
		},
	}
	
	vr.AddErrors(errors)
	
	if vr.GetErrorCount() != 2 {
		t.Errorf("Expected 2 errors, got %d", vr.GetErrorCount())
	}
	
	// Add warnings
	warnings := []DataValidationError{
		{
			Type:      "FIELD_FORMAT",
			TaskID:    "C",
			Field:     "Category",
			Value:     "INVALID",
			Message:   "Invalid category",
			Severity:  "WARNING",
			Timestamp: time.Now(),
		},
	}
	
	vr.AddWarnings(warnings)
	
	if vr.GetWarningCount() != 1 {
		t.Errorf("Expected 1 warning, got %d", vr.GetWarningCount())
	}
	
	// Add info messages
	info := []DataValidationError{
		{
			Type:      "BUSINESS_RULE",
			TaskID:    "D",
			Field:     "Description",
			Value:     "",
			Message:   "Consider adding description",
			Severity:  "INFO",
			Timestamp: time.Now(),
		},
	}
	
	vr.AddInfo(info)
	
	if vr.GetInfoCount() != 1 {
		t.Errorf("Expected 1 info message, got %d", vr.GetInfoCount())
	}
}

func TestGenerateReport(t *testing.T) {
	vr := NewValidationReporter()
	
	// Add some validation data
	vr.AddErrors([]DataValidationError{
		{
			Type:      "REQUIRED_FIELD",
			TaskID:    "A",
			Field:     "ID",
			Value:     "",
			Message:   "Task ID is required",
			Severity:  "ERROR",
			Timestamp: time.Now(),
		},
	})
	
	vr.AddWarnings([]DataValidationError{
		{
			Type:      "FIELD_FORMAT",
			TaskID:    "B",
			Field:     "Category",
			Value:     "INVALID",
			Message:   "Invalid category",
			Severity:  "WARNING",
			Timestamp: time.Now(),
		},
	})
	
	vr.AddInfo([]DataValidationError{
		{
			Type:      "BUSINESS_RULE",
			TaskID:    "C",
			Field:     "Description",
			Value:     "",
			Message:   "Consider adding description",
			Severity:  "INFO",
			Timestamp: time.Now(),
		},
	})
	
	vr.taskCount = 3
	
	report := vr.GenerateReport()
	
	if report == nil {
		t.Error("Expected non-nil report")
	}
	
	if report.TaskCount != 3 {
		t.Errorf("Expected 3 tasks, got %d", report.TaskCount)
	}
	
	if report.ErrorCount != 1 {
		t.Errorf("Expected 1 error, got %d", report.ErrorCount)
	}
	
	if report.WarningCount != 1 {
		t.Errorf("Expected 1 warning, got %d", report.WarningCount)
	}
	
	if report.InfoCount != 1 {
		t.Errorf("Expected 1 info message, got %d", report.InfoCount)
	}
	
	if report.IsValid {
		t.Error("Expected validation to fail")
	}
	
	if report.Summary == "" {
		t.Error("Expected non-empty summary")
	}
	
	if len(report.Errors) != 1 {
		t.Errorf("Expected 1 error in report, got %d", len(report.Errors))
	}
	
	if len(report.Warnings) != 1 {
		t.Errorf("Expected 1 warning in report, got %d", len(report.Warnings))
	}
	
	if len(report.Info) != 1 {
		t.Errorf("Expected 1 info message in report, got %d", len(report.Info))
	}
}

func TestGenerateTextReport(t *testing.T) {
	vr := NewValidationReporter()
	
	// Add some validation data
	vr.AddErrors([]DataValidationError{
		{
			Type:      "REQUIRED_FIELD",
			TaskID:    "A",
			Field:     "ID",
			Value:     "",
			Message:   "Task ID is required",
			Severity:  "ERROR",
			Timestamp: time.Now(),
			Suggestions: []string{"Provide a unique task ID"},
		},
	})
	
	vr.AddWarnings([]DataValidationError{
		{
			Type:      "FIELD_FORMAT",
			TaskID:    "B",
			Field:     "Category",
			Value:     "INVALID",
			Message:   "Invalid category",
			Severity:  "WARNING",
			Timestamp: time.Now(),
			Suggestions: []string{"Use a valid category"},
		},
	})
	
	vr.taskCount = 2
	
	textReport := vr.GenerateTextReport()
	
	if textReport == "" {
		t.Error("Expected non-empty text report")
	}
	
	// Check for key sections
	if !strings.Contains(textReport, "VALIDATION REPORT") {
		t.Error("Expected 'VALIDATION REPORT' in text report")
	}
	
	if !strings.Contains(textReport, "SUMMARY") {
		t.Error("Expected 'SUMMARY' in text report")
	}
	
	if !strings.Contains(textReport, "ERRORS") {
		t.Error("Expected 'ERRORS' in text report")
	}
	
	if !strings.Contains(textReport, "WARNINGS") {
		t.Error("Expected 'WARNINGS' in text report")
	}
	
	if !strings.Contains(textReport, "RECOMMENDATIONS") {
		t.Error("Expected 'RECOMMENDATIONS' in text report")
	}
	
	// Check for error content
	if !strings.Contains(textReport, "Task ID is required") {
		t.Error("Expected error message in text report")
	}
	
	if !strings.Contains(textReport, "Invalid category") {
		t.Error("Expected warning message in text report")
	}
}

func TestGenerateCSVReport(t *testing.T) {
	vr := NewValidationReporter()
	
	// Add some validation data
	vr.AddErrors([]DataValidationError{
		{
			Type:      "REQUIRED_FIELD",
			TaskID:    "A",
			Field:     "ID",
			Value:     "",
			Message:   "Task ID is required",
			Severity:  "ERROR",
			Timestamp: time.Now(),
		},
	})
	
	vr.AddWarnings([]DataValidationError{
		{
			Type:      "FIELD_FORMAT",
			TaskID:    "B",
			Field:     "Category",
			Value:     "INVALID",
			Message:   "Invalid category",
			Severity:  "WARNING",
			Timestamp: time.Now(),
		},
	})
	
	csvReport := vr.GenerateCSVReport()
	
	if csvReport == "" {
		t.Error("Expected non-empty CSV report")
	}
	
	// Check for CSV header
	if !strings.Contains(csvReport, "Severity,Type,TaskID,Field,Value,Message,Timestamp") {
		t.Error("Expected CSV header in report")
	}
	
	// Check for error data
	if !strings.Contains(csvReport, "ERROR,REQUIRED_FIELD,A,ID") {
		t.Error("Expected error data in CSV report")
	}
	
	if !strings.Contains(csvReport, "WARNING,FIELD_FORMAT,B,Category") {
		t.Error("Expected warning data in CSV report")
	}
}

func TestGenerateJSONReport(t *testing.T) {
	vr := NewValidationReporter()
	
	// Add some validation data
	vr.AddErrors([]DataValidationError{
		{
			Type:      "REQUIRED_FIELD",
			TaskID:    "A",
			Field:     "ID",
			Value:     "",
			Message:   "Task ID is required",
			Severity:  "ERROR",
			Timestamp: time.Now(),
		},
	})
	
	vr.taskCount = 1
	
	jsonReport := vr.GenerateJSONReport()
	
	if jsonReport == "" {
		t.Error("Expected non-empty JSON report")
	}
	
	// Check for JSON structure
	if !strings.Contains(jsonReport, "\"summary\"") {
		t.Error("Expected 'summary' field in JSON report")
	}
	
	if !strings.Contains(jsonReport, "\"is_valid\"") {
		t.Error("Expected 'is_valid' field in JSON report")
	}
	
	if !strings.Contains(jsonReport, "\"task_count\"") {
		t.Error("Expected 'task_count' field in JSON report")
	}
	
	if !strings.Contains(jsonReport, "\"error_count\"") {
		t.Error("Expected 'error_count' field in JSON report")
	}
	
	if !strings.Contains(jsonReport, "\"statistics\"") {
		t.Error("Expected 'statistics' field in JSON report")
	}
}

func TestValidationReportConfig(t *testing.T) {
	vr := NewValidationReporter()
	
	// Test default config
	config := vr.config
	if config == nil {
		t.Error("Expected non-nil config")
	}
	
	if !config.IncludeSuggestions {
		t.Error("Expected IncludeSuggestions to be true by default")
	}
	
	if !config.IncludeSeverity {
		t.Error("Expected IncludeSeverity to be true by default")
	}
	
	if !config.GroupByType {
		t.Error("Expected GroupByType to be true by default")
	}
	
	if !config.GroupBySeverity {
		t.Error("Expected GroupBySeverity to be true by default")
	}
	
	if !config.SortBySeverity {
		t.Error("Expected SortBySeverity to be true by default")
	}
	
	// Test custom config
	customConfig := &ValidationReportConfig{
		IncludeSuggestions: false,
		IncludeTimestamps:  true,
		IncludeSeverity:    false,
		GroupByType:        false,
		GroupBySeverity:    false,
		GroupByTask:        true,
		SortBySeverity:     false,
		MaxErrors:          10,
		MaxWarnings:        5,
		MaxInfo:            3,
	}
	
	vr.SetConfig(customConfig)
	
	if vr.config.IncludeSuggestions {
		t.Error("Expected IncludeSuggestions to be false")
	}
	
	if !vr.config.IncludeTimestamps {
		t.Error("Expected IncludeTimestamps to be true")
	}
	
	if vr.config.IncludeSeverity {
		t.Error("Expected IncludeSeverity to be false")
	}
	
	if vr.config.GroupByType {
		t.Error("Expected GroupByType to be false")
	}
	
	if vr.config.GroupBySeverity {
		t.Error("Expected GroupBySeverity to be false")
	}
	
	if !vr.config.GroupByTask {
		t.Error("Expected GroupByTask to be true")
	}
	
	if vr.config.SortBySeverity {
		t.Error("Expected SortBySeverity to be false")
	}
	
	if vr.config.MaxErrors != 10 {
		t.Error("Expected MaxErrors to be 10")
	}
	
	if vr.config.MaxWarnings != 5 {
		t.Error("Expected MaxWarnings to be 5")
	}
	
	if vr.config.MaxInfo != 3 {
		t.Error("Expected MaxInfo to be 3")
	}
}

func TestClear(t *testing.T) {
	vr := NewValidationReporter()
	
	// Add some data
	vr.AddErrors([]DataValidationError{
		{
			Type:      "REQUIRED_FIELD",
			TaskID:    "A",
			Field:     "ID",
			Value:     "",
			Message:   "Task ID is required",
			Severity:  "ERROR",
			Timestamp: time.Now(),
		},
	})
	
	vr.AddWarnings([]DataValidationError{
		{
			Type:      "FIELD_FORMAT",
			TaskID:    "B",
			Field:     "Category",
			Value:     "INVALID",
			Message:   "Invalid category",
			Severity:  "WARNING",
			Timestamp: time.Now(),
		},
	})
	
	vr.taskCount = 2
	
	// Verify data exists
	if vr.GetErrorCount() != 1 {
		t.Error("Expected 1 error before clear")
	}
	
	if vr.GetWarningCount() != 1 {
		t.Error("Expected 1 warning before clear")
	}
	
	if vr.taskCount != 2 {
		t.Error("Expected 2 tasks before clear")
	}
	
	// Clear data
	vr.Clear()
	
	// Verify data is cleared
	if vr.GetErrorCount() != 0 {
		t.Error("Expected 0 errors after clear")
	}
	
	if vr.GetWarningCount() != 0 {
		t.Error("Expected 0 warnings after clear")
	}
	
	if vr.GetInfoCount() != 0 {
		t.Error("Expected 0 info messages after clear")
	}
	
	if vr.taskCount != 0 {
		t.Error("Expected 0 tasks after clear")
	}
	
	if vr.HasErrors() {
		t.Error("Expected no errors after clear")
	}
	
	if vr.HasWarnings() {
		t.Error("Expected no warnings after clear")
	}
	
	if vr.HasIssues() {
		t.Error("Expected no issues after clear")
	}
}

func TestValidationReporterIntegration(t *testing.T) {
	vr := NewValidationReporter()
	
	// Create a comprehensive validation scenario
	errors := []DataValidationError{
		{
			Type:      "REQUIRED_FIELD",
			TaskID:    "A",
			Field:     "ID",
			Value:     "",
			Message:   "Task ID is required",
			Severity:  "ERROR",
			Timestamp: time.Now(),
			Suggestions: []string{"Provide a unique task ID"},
		},
		{
			Type:      "FIELD_FORMAT",
			TaskID:    "B",
			Field:     "Category",
			Value:     "INVALID",
			Message:   "Invalid category",
			Severity:  "WARNING",
			Timestamp: time.Now(),
			Suggestions: []string{"Use a valid category"},
		},
		{
			Type:      "DATA_CONSISTENCY",
			TaskID:    "C",
			Field:     "ParentID",
			Value:     "C",
			Message:   "Task cannot be its own parent",
			Severity:  "ERROR",
			Timestamp: time.Now(),
			Suggestions: []string{"Remove self-reference from ParentID"},
		},
		{
			Type:      "CIRCULAR_DEPENDENCY",
			TaskID:    "D",
			Field:     "Dependencies",
			Value:     "E -> F -> D",
			Message:   "Circular dependency detected",
			Severity:  "ERROR",
			Timestamp: time.Now(),
			Suggestions: []string{"Break the circular dependency"},
		},
		{
			Type:      "BUSINESS_RULE",
			TaskID:    "E",
			Field:     "Description",
			Value:     "",
			Message:   "Consider adding description",
			Severity:  "INFO",
			Timestamp: time.Now(),
			Suggestions: []string{"Add a description for better clarity"},
		},
	}
	
	// Add errors to reporter
	vr.AddErrors(errors[:3]) // First 3 are errors
	vr.AddWarnings(errors[1:2]) // Second one is also a warning
	vr.AddInfo(errors[4:5]) // Last one is info
	
	vr.taskCount = 5
	
	// Generate report
	report := vr.GenerateReport()
	
	if report == nil {
		t.Error("Expected non-nil report")
	}
	
	if report.TaskCount != 5 {
		t.Errorf("Expected 5 tasks, got %d", report.TaskCount)
	}
	
	if report.ErrorCount != 3 {
		t.Errorf("Expected 3 errors, got %d", report.ErrorCount)
	}
	
	if report.WarningCount != 1 {
		t.Errorf("Expected 1 warning, got %d", report.WarningCount)
	}
	
	if report.InfoCount != 1 {
		t.Errorf("Expected 1 info message, got %d", report.InfoCount)
	}
	
	if report.IsValid {
		t.Error("Expected validation to fail")
	}
	
	// Check statistics
	if report.Statistics == nil {
		t.Error("Expected non-nil statistics")
	}
	
	if report.Statistics.TotalIssues != 5 {
		t.Errorf("Expected 5 total issues, got %d", report.Statistics.TotalIssues)
	}
	
	if report.Statistics.ErrorRate != 60.0 {
		t.Errorf("Expected 60.0%% error rate, got %.1f%%", report.Statistics.ErrorRate)
	}
	
	if report.Statistics.WarningRate != 20.0 {
		t.Errorf("Expected 20.0%% warning rate, got %.1f%%", report.Statistics.WarningRate)
	}
	
	if report.Statistics.InfoRate != 20.0 {
		t.Errorf("Expected 20.0%% info rate, got %.1f%%", report.Statistics.InfoRate)
	}
	
	// Check recommendations
	if len(report.Recommendations) == 0 {
		t.Error("Expected recommendations")
	}
	
	// Generate text report
	textReport := vr.GenerateTextReport()
	
	if textReport == "" {
		t.Error("Expected non-empty text report")
	}
	
	// Check for key content
	if !strings.Contains(textReport, "VALIDATION REPORT") {
		t.Error("Expected 'VALIDATION REPORT' in text report")
	}
	
	if !strings.Contains(textReport, "ERRORS") {
		t.Error("Expected 'ERRORS' section in text report")
	}
	
	if !strings.Contains(textReport, "WARNINGS") {
		t.Error("Expected 'WARNINGS' section in text report")
	}
	
	if !strings.Contains(textReport, "INFO MESSAGES") {
		t.Error("Expected 'INFO MESSAGES' section in text report")
	}
	
	if !strings.Contains(textReport, "RECOMMENDATIONS") {
		t.Error("Expected 'RECOMMENDATIONS' section in text report")
	}
	
	if !strings.Contains(textReport, "STATISTICS") {
		t.Error("Expected 'STATISTICS' section in text report")
	}
}
