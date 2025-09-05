package data

import (
	"testing"
	"time"
)

func TestDataIntegrityValidator(t *testing.T) {
	div := NewDataIntegrityValidator()
	
	// Test required fields
	requiredFields := div.GetRequiredFields()
	expectedFields := []string{"ID", "Name", "StartDate", "EndDate"}
	
	if len(requiredFields) != len(expectedFields) {
		t.Errorf("Expected %d required fields, got %d", len(expectedFields), len(requiredFields))
	}
	
	for _, field := range expectedFields {
		found := false
		for _, reqField := range requiredFields {
			if reqField == field {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected required field %s not found", field)
		}
	}
}

func TestValidateTaskIntegrity(t *testing.T) {
	div := NewDataIntegrityValidator()
	
	// Test valid task
	validTask := &Task{
		ID:          "A",
		Name:        "Valid Task",
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		Category:    "IMAGING",
		Status:      "Planned",
		Priority:    1, // Medium
		Assignee:    "John Doe",
		Description: "A valid task for testing",
	}
	
	errors := div.ValidateTaskIntegrity(validTask)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid task, got %d", len(errors))
	}
	
	// Test nil task
	errors = div.ValidateTaskIntegrity(nil)
	if len(errors) == 0 {
		t.Error("Expected errors for nil task")
	}
	
	// Test task with missing required fields
	invalidTask := &Task{
		ID:          "",
		Name:        "",
		StartDate:   time.Time{},
		EndDate:     time.Time{},
	}
	
	errors = div.ValidateTaskIntegrity(invalidTask)
	if len(errors) == 0 {
		t.Error("Expected errors for task with missing required fields")
	}
	
	// Check for specific required field errors
	hasIDError := false
	hasNameError := false
	hasStartDateError := false
	hasEndDateError := false
	
	for _, err := range errors {
		if err.Type == "REQUIRED_FIELD" {
			switch err.Field {
			case "ID":
				hasIDError = true
			case "Name":
				hasNameError = true
			case "StartDate":
				hasStartDateError = true
			case "EndDate":
				hasEndDateError = true
			}
		}
	}
	
	if !hasIDError {
		t.Error("Expected ID required field error")
	}
	if !hasNameError {
		t.Error("Expected Name required field error")
	}
	if !hasStartDateError {
		t.Error("Expected StartDate required field error")
	}
	if !hasEndDateError {
		t.Error("Expected EndDate required field error")
	}
}

func TestValidateFieldFormats(t *testing.T) {
	div := NewDataIntegrityValidator()
	
	// Test invalid ID format
	invalidIDTask := &Task{
		ID:        "Invalid ID with spaces!",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	errors := div.ValidateTaskIntegrity(invalidIDTask)
	
	hasIDFormatError := false
	for _, err := range errors {
		if err.Type == "FIELD_FORMAT" && err.Field == "ID" {
			hasIDFormatError = true
			break
		}
	}
	if !hasIDFormatError {
		t.Error("Expected ID format error")
	}
	
	// Test invalid category
	invalidCategoryTask := &Task{
		ID:        "B",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		Category:  "INVALID_CATEGORY",
	}
	
	errors = div.ValidateTaskIntegrity(invalidCategoryTask)
	
	hasCategoryFormatError := false
	for _, err := range errors {
		if err.Type == "FIELD_FORMAT" && err.Field == "Category" {
			hasCategoryFormatError = true
			break
		}
	}
	if !hasCategoryFormatError {
		t.Error("Expected category format error")
	}
	
	// Test invalid status
	invalidStatusTask := &Task{
		ID:        "C",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		Status:    "INVALID_STATUS",
	}
	
	errors = div.ValidateTaskIntegrity(invalidStatusTask)
	
	hasStatusFormatError := false
	for _, err := range errors {
		if err.Type == "FIELD_FORMAT" && err.Field == "Status" {
			hasStatusFormatError = true
			break
		}
	}
	if !hasStatusFormatError {
		t.Error("Expected status format error")
	}
	
	// Test invalid priority
	invalidPriorityTask := &Task{
		ID:        "D",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		Priority:  10, // Invalid priority
	}
	
	errors = div.ValidateTaskIntegrity(invalidPriorityTask)
	
	hasPriorityFormatError := false
	for _, err := range errors {
		if err.Type == "FIELD_FORMAT" && err.Field == "Priority" {
			hasPriorityFormatError = true
			break
		}
	}
	if !hasPriorityFormatError {
		t.Error("Expected priority format error")
	}
}

func TestValidateDataConsistency(t *testing.T) {
	div := NewDataIntegrityValidator()
	
	// Test zero duration task
	now := time.Now()
	zeroDurationTask := &Task{
		ID:        "A",
		Name:      "Zero Duration Task",
		StartDate: now.AddDate(0, 0, 1),  // Tomorrow
		EndDate:   now.AddDate(0, 0, 1),  // Same day
	}
	
	errors := div.ValidateTaskIntegrity(zeroDurationTask)
	
	// Debug: print all errors
	for _, err := range errors {
		t.Logf("Error: %s - %s - %s", err.Type, err.Field, err.Message)
	}
	
	hasZeroDurationError := false
	for _, err := range errors {
		if err.Type == "DATA_CONSISTENCY" && err.Field == "Duration" {
			hasZeroDurationError = true
			break
		}
	}
	if !hasZeroDurationError {
		t.Error("Expected zero duration error")
	}
	
	// Test negative duration task
	negativeDurationTask := &Task{
		ID:        "B",
		Name:      "Negative Duration Task",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	
	errors = div.ValidateTaskIntegrity(negativeDurationTask)
	
	hasNegativeDurationError := false
	for _, err := range errors {
		if err.Type == "DATA_CONSISTENCY" && err.Field == "Duration" {
			hasNegativeDurationError = true
			break
		}
	}
	if !hasNegativeDurationError {
		t.Error("Expected negative duration error")
	}
	
	// Test milestone with different start/end dates
	milestoneTask := &Task{
		ID:         "C",
		Name:       "Milestone Task",
		StartDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		IsMilestone: true,
	}
	
	errors = div.ValidateTaskIntegrity(milestoneTask)
	
	hasMilestoneError := false
	for _, err := range errors {
		if err.Type == "DATA_CONSISTENCY" && err.Field == "IsMilestone" {
			hasMilestoneError = true
			break
		}
	}
	if !hasMilestoneError {
		t.Error("Expected milestone consistency error")
	}
	
	// Test self-parent task
	selfParentTask := &Task{
		ID:       "D",
		Name:     "Self Parent Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:  time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		ParentID: "D",
	}
	
	errors = div.ValidateTaskIntegrity(selfParentTask)
	
	hasSelfParentError := false
	for _, err := range errors {
		if err.Type == "DATA_CONSISTENCY" && err.Field == "ParentID" {
			hasSelfParentError = true
			break
		}
	}
	if !hasSelfParentError {
		t.Error("Expected self-parent error")
	}
	
	// Test self-dependency task
	selfDepTask := &Task{
		ID:          "E",
		Name:        "Self Dep Task",
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		Dependencies: []string{"E"},
	}
	
	errors = div.ValidateTaskIntegrity(selfDepTask)
	
	hasSelfDepError := false
	for _, err := range errors {
		if err.Type == "DATA_CONSISTENCY" && err.Field == "Dependencies" {
			hasSelfDepError = true
			break
		}
	}
	if !hasSelfDepError {
		t.Error("Expected self-dependency error")
	}
}

func TestValidateBusinessRules(t *testing.T) {
	div := NewDataIntegrityValidator()
	
	// Test task without category
	noCategoryTask := &Task{
		ID:        "A",
		Name:      "No Category Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	errors := div.ValidateTaskIntegrity(noCategoryTask)
	
	hasCategorySuggestion := false
	for _, err := range errors {
		if err.Type == "BUSINESS_RULE" && err.Field == "Category" {
			hasCategorySuggestion = true
			break
		}
	}
	if !hasCategorySuggestion {
		t.Error("Expected category suggestion")
	}
	
	// Test task without status
	noStatusTask := &Task{
		ID:        "B",
		Name:      "No Status Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	errors = div.ValidateTaskIntegrity(noStatusTask)
	
	hasStatusSuggestion := false
	for _, err := range errors {
		if err.Type == "BUSINESS_RULE" && err.Field == "Status" {
			hasStatusSuggestion = true
			break
		}
	}
	if !hasStatusSuggestion {
		t.Error("Expected status suggestion")
	}
	
	// Test task with very short name
	shortNameTask := &Task{
		ID:        "C",
		Name:      "A",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	errors = div.ValidateTaskIntegrity(shortNameTask)
	
	hasShortNameWarning := false
	for _, err := range errors {
		if err.Type == "BUSINESS_RULE" && err.Field == "Name" && err.Severity == "WARNING" {
			hasShortNameWarning = true
			break
		}
	}
	if !hasShortNameWarning {
		t.Error("Expected short name warning")
	}
	
	// Test task with very long name
	longName := "This is a very long task name that exceeds the recommended length and should trigger a warning about being too long and suggesting to shorten it or move details to description"
	longNameTask := &Task{
		ID:        "D",
		Name:      longName,
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	errors = div.ValidateTaskIntegrity(longNameTask)
	
	hasLongNameWarning := false
	for _, err := range errors {
		if err.Type == "BUSINESS_RULE" && err.Field == "Name" && err.Severity == "WARNING" {
			hasLongNameWarning = true
			break
		}
	}
	if !hasLongNameWarning {
		t.Error("Expected long name warning")
	}
	
	// Test task without description
	noDescTask := &Task{
		ID:        "E",
		Name:      "No Description Task",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	errors = div.ValidateTaskIntegrity(noDescTask)
	
	hasDescSuggestion := false
	for _, err := range errors {
		if err.Type == "BUSINESS_RULE" && err.Field == "Description" {
			hasDescSuggestion = true
			break
		}
	}
	if !hasDescSuggestion {
		t.Error("Expected description suggestion")
	}
}

func TestValidateCrossTaskIntegrity(t *testing.T) {
	div := NewDataIntegrityValidator()
	
	// Test duplicate task IDs
	tasks := []*Task{
		{
			ID:        "A",
			Name:      "Task A",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "A", // Duplicate ID
			Name:      "Task B",
			StartDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	
	errors := div.ValidateTasksIntegrity(tasks)
	
	hasDuplicateIDError := false
	for _, err := range errors {
		if err.Type == "CROSS_TASK_INTEGRITY" && err.Field == "ID" {
			hasDuplicateIDError = true
			break
		}
	}
	if !hasDuplicateIDError {
		t.Error("Expected duplicate ID error")
	}
	
	// Test duplicate task names
	tasks2 := []*Task{
		{
			ID:        "A",
			Name:      "Same Name",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "B",
			Name:      "Same Name", // Duplicate name
			StartDate: time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	
	errors = div.ValidateTasksIntegrity(tasks2)
	
	hasDuplicateNameWarning := false
	for _, err := range errors {
		if err.Type == "CROSS_TASK_INTEGRITY" && err.Field == "Name" {
			hasDuplicateNameWarning = true
			break
		}
	}
	if !hasDuplicateNameWarning {
		t.Error("Expected duplicate name warning")
	}
}

func TestValidateDataIntegrity(t *testing.T) {
	div := NewDataIntegrityValidator()
	
	// Create tasks with various integrity issues
	now := time.Now()
	tasks := []*Task{
		{
			ID:        "A",
			Name:      "Valid Task",
			StartDate: now.AddDate(0, 0, 1),  // Tomorrow
			EndDate:   now.AddDate(0, 0, 5),  // 5 days from now
			Category:  "IMAGING",
			Status:    "Planned",
		},
		{
			ID:        "B",
			Name:      "Invalid Task",
			StartDate: time.Time{}, // Missing start date
			EndDate:   now.AddDate(0, 0, 5),  // 5 days from now
			Category:  "INVALID",
			Status:    "INVALID",
			Priority:  10, // Invalid priority
		},
		{
			ID:        "C",
			Name:      "Short",
			StartDate: now.AddDate(0, 0, 1),  // Tomorrow
			EndDate:   now.AddDate(0, 0, 1),  // Same day (zero duration)
			ParentID:  "C", // Self-parent
		},
		{
			ID:        "D",
			Name:      "Invalid Format Task",
			StartDate: now.AddDate(0, 0, 1),  // Tomorrow
			EndDate:   now.AddDate(0, 0, 5),  // 5 days from now
			Category:  "INVALID_CATEGORY",
			Status:    "INVALID_STATUS",
			Priority:  10, // Invalid priority
		},
	}
	
	result := div.ValidateDataIntegrity(tasks)
	
	if result.IsValid {
		t.Error("Expected validation to fail due to errors")
	}
	
	if result.ErrorCount == 0 {
		t.Error("Expected errors to be found")
	}
	
	if result.TaskCount != 4 {
		t.Errorf("Expected 4 tasks, got %d", result.TaskCount)
	}
	
	// Debug: print all errors
	for _, err := range result.Errors {
		t.Logf("Error: %s - %s - %s", err.Type, err.Field, err.Message)
	}
	
	// Check that we have the expected error types (including warnings and info)
	allErrors := append(result.Errors, result.Warnings...)
	allErrors = append(allErrors, result.Info...)
	
	errorTypes := make(map[string]int)
	for _, err := range allErrors {
		errorTypes[err.Type]++
	}
	
	if errorTypes["REQUIRED_FIELD"] == 0 {
		t.Error("Expected REQUIRED_FIELD errors")
	}
	
	if errorTypes["FIELD_FORMAT"] == 0 {
		t.Error("Expected FIELD_FORMAT errors")
	}
	
	if errorTypes["DATA_CONSISTENCY"] == 0 {
		t.Error("Expected DATA_CONSISTENCY errors")
	}
}

func TestDataIntegrityValidatorIntegration(t *testing.T) {
	div := NewDataIntegrityValidator()
	
	// Create a comprehensive test scenario
	tasks := []*Task{
		{
			ID:          "A",
			Name:        "Valid Task",
			StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Category:    "IMAGING",
			Status:      "Planned",
			Priority:    1, // Medium
			Assignee:    "John Doe",
			Description: "A valid task for testing",
		},
		{
			ID:          "B",
			Name:        "Invalid Format Task",
			StartDate:   time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:    "INVALID_CATEGORY",
			Status:      "INVALID_STATUS",
			Priority:    10, // Invalid priority
		},
		{
			ID:          "C",
			Name:        "Consistency Issues Task",
			StartDate:   time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC), // Zero duration
			IsMilestone: true,
			ParentID:    "C", // Self-parent
			Dependencies: []string{"C"}, // Self-dependency
		},
		{
			ID:        "D",
			Name:      "Business Rules Task",
			StartDate: time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			// Missing category, status, priority, description
		},
		{
			ID:        "", // Missing ID
			Name:      "", // Missing Name
			StartDate: time.Time{}, // Missing StartDate
			EndDate:   time.Time{}, // Missing EndDate
		},
	}
	
	result := div.ValidateDataIntegrity(tasks)
	
	if result.IsValid {
		t.Error("Expected validation to fail due to errors")
	}
	
	// Check error distribution
	if result.ErrorCount == 0 {
		t.Error("Expected errors to be found")
	}
	
	if result.WarningCount == 0 {
		t.Error("Expected warnings to be found")
	}
	
	// Check that we have various error types
	errorTypes := make(map[string]int)
	allErrors := append(result.Errors, result.Warnings...)
	allErrors = append(allErrors, result.Info...)
	
	for _, err := range allErrors {
		errorTypes[err.Type]++
	}
	
	// Debug: print all errors
	for _, err := range allErrors {
		t.Logf("Error: %s - %s - %s", err.Type, err.Field, err.Message)
	}
	
	expectedTypes := []string{"REQUIRED_FIELD", "FIELD_FORMAT", "DATA_CONSISTENCY", "BUSINESS_RULE"}
	for _, expectedType := range expectedTypes {
		if errorTypes[expectedType] == 0 {
			t.Errorf("Expected %s errors to be found", expectedType)
		}
	}
}
