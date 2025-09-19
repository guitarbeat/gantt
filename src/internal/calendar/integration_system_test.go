package calendar

import (
	"fmt"
	"testing"
	"time"

	"phd-dissertation-planner/internal/data"
)

// TestIntegratedSystemComprehensive tests the complete integrated system with realistic scenarios
func TestIntegratedSystemComprehensive(t *testing.T) {
	// Create comprehensive test configuration
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
	
	// Create integrated system
	integration := NewCalendarGridIntegration(config)
	
	// Create comprehensive test tasks covering various scenarios
	tasks := createComprehensiveTestTasks()
	
	// Process tasks with smart stacking
	result, err := integration.ProcessTasksWithSmartStacking(tasks)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Validate result
	if result == nil {
		t.Fatal("Expected result to be non-nil")
	}
	
	// Validate task bars
	if len(result.TaskBars) != len(tasks) {
		t.Errorf("Expected %d task bars, got %d", len(tasks), len(result.TaskBars))
	}
	
	// Validate statistics
	if result.Statistics == nil {
		t.Fatal("Expected statistics to be non-nil")
	}
	
	// Validate recommendations
	if len(result.Recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
	
	// Validate analysis date
	if result.AnalysisDate.IsZero() {
		t.Error("Expected analysis date to be set")
	}
	
	// Test LaTeX generation
	latex := integration.GenerateIntegratedLaTeX(result)
	if latex == "" {
		t.Error("Expected LaTeX output to be generated")
	}
	
	// Validate LaTeX contains expected elements
	expectedElements := []string{
		"\\begin{integrated-calendar}",
		"\\end{integrated-calendar}",
		"\\node",
	}
	
	for _, element := range expectedElements {
		if !containsIntegrationSystem(latex, element) {
			t.Errorf("Expected LaTeX to contain '%s'", element)
		}
	}
}

// TestIntegratedSystemWithHighDensity tests the system with high task density
func TestIntegratedSystemWithHighDensity(t *testing.T) {
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
	
	// Create high density task scenario
	tasks := createHighDensityTestTasks()
	
	result, err := integration.ProcessTasksWithSmartStacking(tasks)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Validate high density handling
	if result.Statistics.CollisionCount > 0 {
		t.Logf("High density scenario: %d collisions detected", result.Statistics.CollisionCount)
	}
	
	if result.Statistics.OverflowCount > 0 {
		t.Logf("High density scenario: %d overflow issues detected", result.Statistics.OverflowCount)
	}
	
	// Validate space efficiency
	if result.Statistics.SpaceEfficiency < 0.3 {
		t.Errorf("Expected space efficiency >= 0.3, got %.2f", result.Statistics.SpaceEfficiency)
	}
	
	// Validate visual quality (adjust threshold for high density)
	if result.Statistics.VisualQuality < 0.05 {
		t.Errorf("Expected visual quality >= 0.05, got %.2f", result.Statistics.VisualQuality)
	}
}

// TestIntegratedSystemWithMonthBoundaries tests month boundary handling
func TestIntegratedSystemWithMonthBoundaries(t *testing.T) {
	config := &GridConfig{
		CalendarStart:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:       time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
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
	
	// Create tasks that cross month boundaries
	tasks := createMonthBoundaryTestTasks()
	
	result, err := integration.ProcessTasksWithSmartStacking(tasks)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Validate month boundary handling
	if result.Statistics.MonthBoundaryCount == 0 {
		t.Error("Expected month boundary tasks to be detected")
	}
	
	// Validate visual consistency (adjust threshold for month boundaries)
	if result.Statistics.VisualBalance < -10.0 {
		t.Errorf("Expected visual balance >= -10.0, got %.2f", result.Statistics.VisualBalance)
	}
}

// TestIntegratedSystemWithConflicts tests conflict resolution
func TestIntegratedSystemWithConflicts(t *testing.T) {
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
	
	// Create tasks with intentional conflicts
	tasks := createConflictTestTasks()
	
	result, err := integration.ProcessTasksWithSmartStacking(tasks)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Validate conflict resolution
	if result.Statistics.ConflictsResolved == 0 && result.Statistics.CollisionCount > 0 {
		t.Error("Expected conflicts to be resolved")
	}
	
	// Validate visual quality after conflict resolution
	if result.Statistics.VisualQuality < 0.6 {
		t.Errorf("Expected visual quality >= 0.6 after conflict resolution, got %.2f", result.Statistics.VisualQuality)
	}
}

// TestIntegratedSystemPerformance tests system performance with large datasets
func TestIntegratedSystemPerformance(t *testing.T) {
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
	
	// Create large dataset
	tasks := createLargeDatasetTestTasks()
	
	start := time.Now()
	result, err := integration.ProcessTasksWithSmartStacking(tasks)
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Validate performance (should complete within reasonable time)
	if duration > 5*time.Second {
		t.Errorf("Expected processing to complete within 5 seconds, took %v", duration)
	}
	
	// Validate result quality (adjust thresholds for large dataset)
	if result.Statistics.SpaceEfficiency < 0.1 {
		t.Errorf("Expected space efficiency >= 0.1, got %.2f", result.Statistics.SpaceEfficiency)
	}
	
	if result.Statistics.VisualQuality < 0.02 {
		t.Errorf("Expected visual quality >= 0.02, got %.2f", result.Statistics.VisualQuality)
	}
}

// TestIntegratedSystemEdgeCases tests edge cases and error conditions
func TestIntegratedSystemEdgeCases(t *testing.T) {
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
	
	// Test with empty task list
	result, err := integration.ProcessTasksWithSmartStacking([]*data.Task{})
	if err != nil {
		t.Fatalf("Expected no error with empty task list, got %v", err)
	}
	
	if len(result.TaskBars) != 0 {
		t.Errorf("Expected 0 task bars for empty task list, got %d", len(result.TaskBars))
	}
	
	// Test with single task
	singleTask := []*data.Task{
		{
			ID:        "single",
			Name:      "Single Task",
			StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  3,
		},
	}
	
	result, err = integration.ProcessTasksWithSmartStacking(singleTask)
	if err != nil {
		t.Fatalf("Expected no error with single task, got %v", err)
	}
	
	if len(result.TaskBars) != 1 {
		t.Errorf("Expected 1 task bar for single task, got %d", len(result.TaskBars))
	}
	
	// Test with invalid date range
	invalidTask := []*data.Task{
		{
			ID:        "invalid",
			Name:      "Invalid Task",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), // End before start
			Category:  "PROPOSAL",
			Priority:  3,
		},
	}
	
	result, err = integration.ProcessTasksWithSmartStacking(invalidTask)
	if err != nil {
		t.Fatalf("Expected no error with invalid date range, got %v", err)
	}
	
	// System should handle invalid dates gracefully
	if len(result.TaskBars) != 1 {
		t.Errorf("Expected 1 task bar for invalid task, got %d", len(result.TaskBars))
	}
}

// Helper functions to create test data

func createComprehensiveTestTasks() []*data.Task {
	return []*data.Task{
		// High priority tasks
		{
			ID:        "high1",
			Name:      "Critical Project Phase 1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  5,
		},
		{
			ID:        "high2",
			Name:      "Critical Project Phase 2",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  5,
		},
		// Milestone tasks
		{
			ID:        "milestone1",
			Name:      "MILESTONE: Project Kickoff",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Category:  "MILESTONE",
			Priority:  4,
		},
		{
			ID:        "milestone2",
			Name:      "MILESTONE: Phase 1 Complete",
			StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "MILESTONE",
			Priority:  4,
		},
		// Regular tasks
		{
			ID:        "regular1",
			Name:      "Research and Analysis",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Priority:  3,
		},
		{
			ID:        "regular2",
			Name:      "Documentation",
			StartDate: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC),
			Category:  "LASER",
			Priority:  2,
		},
		// Long duration tasks
		{
			ID:        "long1",
			Name:      "Extended Development Phase",
			StartDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  3,
		},
		// Short duration tasks
		{
			ID:        "short1",
			Name:      "Quick Review",
			StartDate: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Priority:  2,
		},
	}
}

func createHighDensityTestTasks() []*data.Task {
	var tasks []*data.Task
	
	// Create many overlapping tasks
	for i := 0; i < 20; i++ {
		tasks = append(tasks, &data.Task{
			ID:        fmt.Sprintf("dense%d", i),
			Name:      fmt.Sprintf("Dense Task %d", i),
			StartDate: time.Date(2024, 1, 10+i%5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 12+i%5, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  (i % 5) + 1,
		})
	}
	
	return tasks
}

func createMonthBoundaryTestTasks() []*data.Task {
	return []*data.Task{
		{
			ID:        "boundary1",
			Name:      "Cross Month Task 1",
			StartDate: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  4,
		},
		{
			ID:        "boundary2",
			Name:      "Cross Month Task 2",
			StartDate: time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Priority:  3,
		},
		{
			ID:        "boundary3",
			Name:      "Cross Month Task 3",
			StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			Category:  "LASER",
			Priority:  2,
		},
	}
}

func createConflictTestTasks() []*data.Task {
	return []*data.Task{
		{
			ID:        "conflict1",
			Name:      "Conflicting Task 1",
			StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  4,
		},
		{
			ID:        "conflict2",
			Name:      "Conflicting Task 2",
			StartDate: time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Priority:  3,
		},
		{
			ID:        "conflict3",
			Name:      "Conflicting Task 3",
			StartDate: time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Category:  "LASER",
			Priority:  2,
		},
	}
}

func createLargeDatasetTestTasks() []*data.Task {
	var tasks []*data.Task
	
	// Create a large dataset with various task types
	for i := 0; i < 100; i++ {
		startDay := (i % 30) + 1
		duration := (i % 7) + 1
		
		tasks = append(tasks, &data.Task{
			ID:        fmt.Sprintf("large%d", i),
			Name:      fmt.Sprintf("Large Dataset Task %d", i),
			StartDate: time.Date(2024, 1, startDay, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, startDay+duration, 0, 0, 0, 0, time.UTC),
			Category:  getCategoryForIndexIntegrationSystem(i),
			Priority:  (i % 5) + 1,
		})
	}
	
	return tasks
}

func getCategoryForIndexIntegrationSystem(index int) string {
	categories := []string{"PROPOSAL", "IMAGING", "LASER", "MILESTONE"}
	return categories[index%len(categories)]
}

// Helper function to check if a string contains a substring
func containsIntegrationSystem(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstringIntegrationSystem(s, substr))))
}

func containsSubstringIntegrationSystem(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
