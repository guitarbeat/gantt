package calendar

import (
	"fmt"
	"testing"
	"time"

	"latex-yearly-planner/internal/data"
)

// TestValidationScenarios tests various validation scenarios
func TestValidationScenarios(t *testing.T) {
	config := &GridConfig{
		CalendarStart:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:       time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:          20.0,
		DayHeight:         15.0,
		RowHeight:         12.0,
		MaxRowsPerDay:     6,
		OverlapThreshold:  0.1,
		MonthBoundaryGap:  2.0,
		TaskSpacing:       1.0,
		VisualConstraints: &VisualConstraints{
			MaxStackHeight:     90.0,
			MinTaskHeight:      6.0,
			MaxTaskHeight:      24.0,
			MinTaskWidth:       2.0,
			MaxTaskWidth:       140.0,
			VerticalSpacing:    1.0,
			HorizontalSpacing:  1.0,
			MaxStackDepth:      6,
			CollisionThreshold: 0.1,
			OverflowThreshold:  0.8,
		},
	}
	
	integration := NewCalendarGridIntegration(config)
	
	// Test 1: Validate layout accuracy
	t.Run("LayoutAccuracy", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "layout1",
				Name:      "Layout Test Task",
				StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  3,
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Validate task bar positioning
		if len(result.TaskBars) != 1 {
			t.Fatalf("Expected 1 task bar, got %d", len(result.TaskBars))
		}
		
		bar := result.TaskBars[0]
		
		// Validate X positioning (5 days from start = 5 * 20 = 100)
		expectedStartX := 4.0 * 20.0 // 4 days from start (0-indexed)
		if bar.StartX != expectedStartX {
			t.Errorf("Expected StartX %.2f, got %.2f", expectedStartX, bar.StartX)
		}
		
		// Validate width (5 days = 5 * 20 = 100)
		expectedWidth := 5.0 * 20.0
		if bar.Width != expectedWidth {
			t.Errorf("Expected Width %.2f, got %.2f", expectedWidth, bar.Width)
		}
		
		// Validate height is within constraints
		if bar.Height < config.VisualConstraints.MinTaskHeight {
			t.Errorf("Expected height >= %.2f, got %.2f", config.VisualConstraints.MinTaskHeight, bar.Height)
		}
		
		if bar.Height > config.VisualConstraints.MaxTaskHeight {
			t.Errorf("Expected height <= %.2f, got %.2f", config.VisualConstraints.MaxTaskHeight, bar.Height)
		}
	})
	
	// Test 2: Validate visual consistency
	t.Run("VisualConsistency", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "visual1",
				Name:      "Visual Test Task 1",
				StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  3,
			},
			{
				ID:        "visual2",
				Name:      "Visual Test Task 2",
				StartDate: time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  3,
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Validate visual consistency
		if len(result.TaskBars) != 2 {
			t.Fatalf("Expected 2 task bars, got %d", len(result.TaskBars))
		}
		
		bar1 := result.TaskBars[0]
		bar2 := result.TaskBars[1]
		
		// Both tasks should have the same category color
		if bar1.Color != bar2.Color {
			t.Errorf("Expected same color for same category, got %s and %s", bar1.Color, bar2.Color)
		}
		
		// Both tasks should have similar height (same priority)
		heightDiff := bar1.Height - bar2.Height
		if heightDiff < 0 {
			heightDiff = -heightDiff
		}
		if heightDiff > 2.0 { // Allow small variations
			t.Errorf("Expected similar heights for same priority, got %.2f and %.2f", bar1.Height, bar2.Height)
		}
	})
	
	// Test 3: Validate conflict resolution
	t.Run("ConflictResolution", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "conflict1",
				Name:      "Conflict Task 1",
				StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  4,
			},
			{
				ID:        "conflict2",
				Name:      "Conflict Task 2",
				StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 13, 0, 0, 0, 0, time.UTC),
				Category:  "IMAGING",
				Priority:  3,
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Validate conflict resolution
		if result.Statistics.CollisionCount > 0 {
			t.Logf("Conflicts detected: %d", result.Statistics.CollisionCount)
		}
		
		// Validate that higher priority task is positioned better
		bar1 := result.TaskBars[0]
		bar2 := result.TaskBars[1]
		
		// Higher priority task should have better positioning
		if bar1.Priority > bar2.Priority {
			if bar1.Y > bar2.Y {
				t.Error("Expected higher priority task to be positioned above lower priority task")
			}
		}
	})
	
	// Test 4: Validate month boundary handling
	t.Run("MonthBoundaryHandling", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "boundary1",
				Name:      "Month Boundary Task",
				StartDate: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  4,
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Validate month boundary handling
		if result.Statistics.MonthBoundaryCount == 0 {
			t.Error("Expected month boundary tasks to be detected")
		}
		
		// Validate task bar has month boundary flag
		if len(result.TaskBars) != 1 {
			t.Fatalf("Expected 1 task bar, got %d", len(result.TaskBars))
		}
		
		bar := result.TaskBars[0]
		if !bar.MonthBoundary {
			t.Error("Expected task bar to have month boundary flag")
		}
	})
	
	// Test 5: Validate LaTeX generation
	t.Run("LaTeXGeneration", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "latex1",
				Name:      "LaTeX Test Task",
				StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  3,
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Generate LaTeX
		latex := integration.GenerateIntegratedLaTeX(result)
		if latex == "" {
			t.Error("Expected LaTeX output to be generated")
		}
		
		// Validate LaTeX structure
		expectedElements := []string{
			"\\begin{integrated-calendar}",
			"\\end{integrated-calendar}",
			"\\node",
		}
		
		for _, element := range expectedElements {
			if !containsValidation(latex, element) {
				t.Errorf("Expected LaTeX to contain '%s'", element)
			}
		}
		
		// Validate LaTeX contains task information
		if !containsValidation(latex, "LaTeX Test Task") {
			t.Error("Expected LaTeX to contain task name")
		}
	})
}

// TestEdgeCaseScenarios tests edge cases and error conditions
func TestEdgeCaseScenarios(t *testing.T) {
	config := &GridConfig{
		CalendarStart:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:       time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:          20.0,
		DayHeight:         15.0,
		RowHeight:         12.0,
		MaxRowsPerDay:     4,
		OverlapThreshold:  0.1,
		MonthBoundaryGap:  2.0,
		TaskSpacing:       1.0,
		VisualConstraints: &VisualConstraints{
			MaxStackHeight:     60.0,
			MinTaskHeight:      6.0,
			MaxTaskHeight:      24.0,
			MinTaskWidth:       2.0,
			MaxTaskWidth:       140.0,
			VerticalSpacing:    1.0,
			HorizontalSpacing:  1.0,
			MaxStackDepth:      4,
			CollisionThreshold: 0.1,
			OverflowThreshold:  0.8,
		},
	}
	
	integration := NewCalendarGridIntegration(config)
	
	// Test 1: Empty task list
	t.Run("EmptyTaskList", func(t *testing.T) {
		result, err := integration.ProcessTasksWithSmartStacking([]*data.Task{})
		if err != nil {
			t.Fatalf("Expected no error with empty task list, got %v", err)
		}
		
		if len(result.TaskBars) != 0 {
			t.Errorf("Expected 0 task bars for empty task list, got %d", len(result.TaskBars))
		}
		
		if result.Statistics.TotalTasks != 0 {
			t.Errorf("Expected 0 total tasks, got %d", result.Statistics.TotalTasks)
		}
	})
	
	// Test 2: Single task
	t.Run("SingleTask", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "single",
				Name:      "Single Task",
				StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  3,
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error with single task, got %v", err)
		}
		
		if len(result.TaskBars) != 1 {
			t.Errorf("Expected 1 task bar for single task, got %d", len(result.TaskBars))
		}
		
		if result.Statistics.TotalTasks != 1 {
			t.Errorf("Expected 1 total task, got %d", result.Statistics.TotalTasks)
		}
	})
	
	// Test 3: Invalid date range (end before start)
	t.Run("InvalidDateRange", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "invalid",
				Name:      "Invalid Task",
				StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), // End before start
				Category:  "PROPOSAL",
				Priority:  3,
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error with invalid date range, got %v", err)
		}
		
		// System should handle invalid dates gracefully
		if len(result.TaskBars) != 1 {
			t.Errorf("Expected 1 task bar for invalid task, got %d", len(result.TaskBars))
		}
	})
	
	// Test 4: Task outside calendar range
	t.Run("TaskOutsideCalendarRange", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "outside",
				Name:      "Outside Task",
				StartDate: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), // Before calendar start
				EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  3,
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error with task outside calendar range, got %v", err)
		}
		
		// System should handle tasks outside range gracefully
		if len(result.TaskBars) != 1 {
			t.Errorf("Expected 1 task bar for outside task, got %d", len(result.TaskBars))
		}
	})
	
	// Test 5: Very high priority task
	t.Run("VeryHighPriorityTask", func(t *testing.T) {
		tasks := []*data.Task{
			{
				ID:        "high",
				Name:      "Very High Priority Task",
				StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Category:  "PROPOSAL",
				Priority:  10, // Very high priority
			},
		}
		
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		if err != nil {
			t.Fatalf("Expected no error with very high priority task, got %v", err)
		}
		
		if len(result.TaskBars) != 1 {
			t.Errorf("Expected 1 task bar for very high priority task, got %d", len(result.TaskBars))
		}
		
		bar := result.TaskBars[0]
		if bar.Priority != 10 {
			t.Errorf("Expected priority 10, got %d", bar.Priority)
		}
	})
}

// TestPerformanceScenarios tests performance with various scenarios
func TestPerformanceScenarios(t *testing.T) {
	config := &GridConfig{
		CalendarStart:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:       time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:          20.0,
		DayHeight:         15.0,
		RowHeight:         12.0,
		MaxRowsPerDay:     6,
		OverlapThreshold:  0.1,
		MonthBoundaryGap:  2.0,
		TaskSpacing:       1.0,
		VisualConstraints: &VisualConstraints{
			MaxStackHeight:     90.0,
			MinTaskHeight:      6.0,
			MaxTaskHeight:      24.0,
			MinTaskWidth:       2.0,
			MaxTaskWidth:       140.0,
			VerticalSpacing:    1.0,
			HorizontalSpacing:  1.0,
			MaxStackDepth:      6,
			CollisionThreshold: 0.1,
			OverflowThreshold:  0.8,
		},
	}
	
	integration := NewCalendarGridIntegration(config)
	
	// Test 1: Performance with medium dataset
	t.Run("MediumDataset", func(t *testing.T) {
		tasks := createMediumDataset(50)
		
		start := time.Now()
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		duration := time.Since(start)
		
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Should complete within 2 seconds
		if duration > 2*time.Second {
			t.Errorf("Expected processing to complete within 2 seconds, took %v", duration)
		}
		
		// Validate result quality
		if result.Statistics.SpaceEfficiency < 0.3 {
			t.Errorf("Expected space efficiency >= 0.3, got %.2f", result.Statistics.SpaceEfficiency)
		}
	})
	
	// Test 2: Performance with large dataset
	t.Run("LargeDataset", func(t *testing.T) {
		tasks := createMediumDataset(200)
		
		start := time.Now()
		result, err := integration.ProcessTasksWithSmartStacking(tasks)
		duration := time.Since(start)
		
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Should complete within 5 seconds
		if duration > 5*time.Second {
			t.Errorf("Expected processing to complete within 5 seconds, took %v", duration)
		}
		
		// Validate result quality
		if result.Statistics.SpaceEfficiency < 0.2 {
			t.Errorf("Expected space efficiency >= 0.2, got %.2f", result.Statistics.SpaceEfficiency)
		}
	})
}

// Helper functions

func createMediumDataset(count int) []*data.Task {
	var tasks []*data.Task
	
	for i := 0; i < count; i++ {
		startDay := (i % 30) + 1
		duration := (i % 7) + 1
		
		tasks = append(tasks, &data.Task{
			ID:        fmt.Sprintf("medium%d", i),
			Name:      fmt.Sprintf("Medium Dataset Task %d", i),
			StartDate: time.Date(2024, 1, startDay, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, startDay+duration, 0, 0, 0, 0, time.UTC),
			Category:  getCategoryForIndexValidation(i),
			Priority:  (i % 5) + 1,
		})
	}
	
	return tasks
}

func getCategoryForIndexValidation(index int) string {
	categories := []string{"PROPOSAL", "IMAGING", "LASER", "MILESTONE"}
	return categories[index%len(categories)]
}

// Helper function to check if a string contains a substring
func containsValidation(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstringValidation(s, substr))))
}

func containsSubstringValidation(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
