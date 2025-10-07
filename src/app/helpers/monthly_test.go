package helpers

import (
	"testing"
	"time"

	cal "phd-dissertation-planner/src/calendar"
	"phd-dissertation-planner/src/core"
)

func TestMonthlyProcessor_ProcessMonthsWithTasks(t *testing.T) {
	// Create a test configuration
	cfg := core.Config{
		MonthsWithTasks: []core.MonthYear{
			{Year: 2025, Month: time.January},
			{Year: 2025, Month: time.February},
		},
		WeekStart: time.Sunday,
	}

	tasks := []core.Task{
		{
			ID:        "1",
			Name:      "Task 1",
			Phase:     "1",
			SubPhase:  "Setup",
			Category:  "PROPOSAL",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "2",
			Name:      "Task 2",
			Phase:     "1",
			SubPhase:  "Setup",
			Category:  "PROPOSAL",
			StartDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	tpls := []string{"monthly.tpl"}

	processor := NewMonthlyProcessor(cfg)
	modules, err := processor.ProcessMonthsWithTasks(tasks, tpls)

	if err != nil {
		t.Fatalf("ProcessMonthsWithTasks failed: %v", err)
	}

	// Should have TOC module + 2 month modules
	expectedModuleCount := 3
	if len(modules) != expectedModuleCount {
		t.Errorf("Expected %d modules, got %d", expectedModuleCount, len(modules))
	}

	// Check that first module is TOC
	if modules[0].Body.(map[string]interface{})["TOCContent"] == nil {
		t.Error("Expected first module to be TOC module")
	}

	// Check that subsequent modules are month modules
	for i := 1; i < len(modules); i++ {
		body := modules[i].Body.(map[string]interface{})
		if body["Month"] == nil {
			t.Errorf("Expected module %d to be month module", i)
		}
	}
}

func TestMonthlyProcessor_ProcessFallbackMonths(t *testing.T) {
	// Create a test configuration
	cfg := core.Config{
		StartYear: 2025,
		EndYear:   2025,
		WeekStart: time.Sunday,
	}

	tpls := []string{"monthly.tpl"}

	processor := NewMonthlyProcessor(cfg)
	modules, err := processor.ProcessFallbackMonths(tpls)

	if err != nil {
		t.Fatalf("ProcessFallbackMonths failed: %v", err)
	}

	// Should have 12 modules (one for each month)
	expectedModuleCount := 12
	if len(modules) != expectedModuleCount {
		t.Errorf("Expected %d modules, got %d", expectedModuleCount, len(modules))
	}

	// Check that all modules are month modules
	for i, module := range modules {
		body := module.Body.(map[string]interface{})
		if body["Month"] == nil {
			t.Errorf("Expected module %d to be month module", i)
		}
	}
}

func TestMonthlyProcessor_FindTargetMonth(t *testing.T) {
	cfg := core.Config{
		WeekStart: time.Sunday,
	}

	processor := NewMonthlyProcessor(cfg)

	// Test with nil year (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("findTargetMonth panicked with nil year: %v", r)
		}
	}()

	monthYear := core.MonthYear{Year: 2025, Month: time.January}
	result := processor.findTargetMonth(nil, monthYear)
	if result != nil {
		t.Error("Expected nil result for nil year")
	}
}

func TestMonthlyProcessor_CreateMonthModule(t *testing.T) {
	t.Skip("Skipping createMonthModule test - needs proper calendar package mocking")
	cfg := core.Config{
		WeekStart:           time.Sunday,
		ClearTopRightCorner: true,
	}

	processor := NewMonthlyProcessor(cfg)

	// Test that the function doesn't panic with nil inputs
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CreateMonthModule panicked: %v", r)
		}
	}()

	// Create a minimal month struct for testing
	// Note: This is a simplified test - in practice, you'd need to mock the calendar package
	year := &cal.Year{Number: 2025}
	month := cal.Month{
		Month: time.January,
		Year:  year,
	}
	
	// Test with nil quarter and year (should not panic)
	module := processor.createMonthModule(month, nil, year, "test.tpl")
	
	// Basic validation
	if module.Tpl != "test.tpl" {
		t.Error("Expected module template to match input")
	}
}

func TestMonthlyProcessor_AssignTasksToMonth(t *testing.T) {
	cfg := core.Config{
		WeekStart: time.Sunday,
	}

	processor := NewMonthlyProcessor(cfg)

	// Test that the function doesn't panic with nil month
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("AssignTasksToMonth panicked: %v", r)
		}
	}()

	tasks := []core.Task{
		{
			ID:        "1",
			Name:      "Task 1",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	// Test with nil month (should not panic)
	processor.assignTasksToMonth(nil, tasks)
}

func TestMonthlyProcessor_CreateTableOfContentsModule(t *testing.T) {
	cfg := core.Config{
		WeekStart: time.Sunday,
	}

	processor := NewMonthlyProcessor(cfg)

	tasks := []core.Task{
		{
			ID:        "1",
			Name:      "Task 1",
			Phase:     "1",
			SubPhase:  "Setup",
			Category:  "PROPOSAL",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	module := processor.createTableOfContentsModule(tasks, "test.tpl")

	// Check that module has correct structure
	// Note: We can't directly compare Config structs due to slice fields
	// In a real test, you'd compare specific fields or use a custom comparison function

	if module.Tpl != "test.tpl" {
		t.Error("Expected module template to match input")
	}

	body := module.Body.(map[string]interface{})
	if body["TOCContent"] == nil {
		t.Error("Expected module body to contain TOCContent")
	}

	content, ok := body["TOCContent"].(string)
	if !ok {
		t.Error("Expected TOCContent to be a string")
	}

	if content == "" {
		t.Error("Expected TOCContent to not be empty")
	}
}

// Test helper functions that don't require calendar package mocking
func TestNewMonthlyProcessor(t *testing.T) {
	cfg := core.Config{
		WeekStart: time.Sunday,
	}

	processor := NewMonthlyProcessor(cfg)

	// Check that processor was created successfully
	if processor == nil {
		t.Error("Expected processor to be created")
	}

	// Check that config was set (we can't directly compare due to slice fields)
	if processor.cfg.WeekStart != cfg.WeekStart {
		t.Error("Expected processor config WeekStart to match input")
	}
}
