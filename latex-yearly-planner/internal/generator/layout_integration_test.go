package generator

import (
	"testing"
	"time"

	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
)

func TestLayoutIntegration(t *testing.T) {

	// Create test tasks
	tasks := []data.Task{
		{
			ID:          "task1",
			Name:        "Test Task 1",
			Description: "A test task for integration",
			Category:    "RESEARCH",
			Priority:    3,
			StartDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 1, 20, 0, 0, 0, 0, time.Local),
		},
		{
			ID:          "task2",
			Name:        "Test Task 2",
			Description: "Another test task",
			Category:    "DISSERTATION",
			Priority:    4,
			StartDate:   time.Date(2024, 1, 18, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 1, 25, 0, 0, 0, 0, time.Local),
		},
	}

	// Create layout integration
	layoutIntegration := NewLayoutIntegration()

	// Test processing tasks with layout
	taskPointers := make([]*data.Task, len(tasks))
	for i := range tasks {
		taskPointers[i] = &tasks[i]
	}

	result, err := layoutIntegration.ProcessTasksWithLayout(taskPointers)
	if err != nil {
		t.Fatalf("Failed to process tasks with layout: %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.TaskBars) == 0 {
		t.Fatal("Expected task bars to be generated")
	}

	// Test task visualization generation
	visualization := layoutIntegration.GenerateTaskVisualization(result)
	if visualization == "" {
		t.Fatal("Expected non-empty visualization")
	}

	// Test statistics
	stats := layoutIntegration.GetLayoutStatistics(result)
	if stats == nil {
		t.Fatal("Expected non-nil statistics")
	}

	if stats.TotalTasks == 0 {
		t.Fatal("Expected non-zero total tasks")
	}
}

func TestEnhancedMonthly(t *testing.T) {
	// Create a test configuration with CSV file
	cfg := config.Config{
		CSVFilePath: "test_data.csv",
		WeekStart:   time.Monday,
		MonthsWithTasks: []data.MonthYear{
			{Month: time.January, Year: 2024},
		},
	}

	// Create layout integration
	layoutIntegration := NewLayoutIntegration()

	// Test enhanced monthly generation (this will fail without actual CSV file, but tests the structure)
	modules, err := layoutIntegration.EnhancedMonthly(cfg, []string{"monthly_page.tpl"})
	
	// We expect an error due to missing CSV file, but the structure should be correct
	if err == nil {
		t.Fatal("Expected error due to missing CSV file")
	}

	// Test with empty configuration
	cfgEmpty := config.Config{
		CSVFilePath: "",
		WeekStart:   time.Monday,
	}

	modules, err = layoutIntegration.EnhancedMonthly(cfgEmpty, []string{"monthly_page.tpl"})
	if err != nil {
		t.Fatalf("Unexpected error with empty config: %v", err)
	}

	if len(modules) == 0 {
		t.Fatal("Expected some modules to be generated")
	}

	// Verify module structure
	module := modules[0]
	if module.Cfg.WeekStart != time.Monday {
		t.Fatal("Expected correct week start")
	}

	// Check if layout data is present
	body, ok := module.Body.(map[string]interface{})
	if !ok {
		t.Fatal("Expected body to be a map")
	}

	// Layout data should be present even if empty
	if _, hasLayout := body["LayoutResult"]; !hasLayout {
		t.Fatal("Expected LayoutResult in body")
	}

	if _, hasTaskBars := body["TaskBars"]; !hasTaskBars {
		t.Fatal("Expected TaskBars in body")
	}

	if _, hasLayoutStats := body["LayoutStats"]; !hasLayoutStats {
		t.Fatal("Expected LayoutStats in body")
	}
}
