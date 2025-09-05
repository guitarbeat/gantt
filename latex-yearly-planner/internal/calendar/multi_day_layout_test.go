package calendar

import (
	"testing"
	"time"

	"latex-yearly-planner/internal/data"
)

func TestMultiDayLayoutEngine(t *testing.T) {
	// Create test calendar range
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create layout engine
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Test basic properties
	if engine.calendarStart != calendarStart {
		t.Errorf("Expected calendar start %v, got %v", calendarStart, engine.calendarStart)
	}
	
	if engine.calendarEnd != calendarEnd {
		t.Errorf("Expected calendar end %v, got %v", calendarEnd, engine.calendarEnd)
	}
	
	if engine.dayWidth != 20.0 {
		t.Errorf("Expected day width 20.0, got %f", engine.dayWidth)
	}
	
	if engine.dayHeight != 15.0 {
		t.Errorf("Expected day height 15.0, got %f", engine.dayHeight)
	}
}

func TestGroupOverlappingTasks(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Task 1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  1,
		},
		{
			ID:        "task2",
			Name:      "Task 2",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "LASER",
			Priority:  2,
		},
		{
			ID:        "task3",
			Name:      "Task 3",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Priority:  3,
		},
	}
	
	groups := engine.groupOverlappingTasks(tasks)
	
	// Should have 2 groups: tasks 1&2 overlap, task 3 is separate
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
	
	// First group should contain tasks 1 and 2
	if len(groups[0].Tasks) != 2 {
		t.Errorf("Expected first group to have 2 tasks, got %d", len(groups[0].Tasks))
	}
	
	// Second group should contain task 3
	if len(groups[1].Tasks) != 1 {
		t.Errorf("Expected second group to have 1 task, got %d", len(groups[1].Tasks))
	}
}

func TestLayoutMultiDayTasks(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Task 1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  1,
		},
		{
			ID:        "task2",
			Name:      "Task 2",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "LASER",
			Priority:  2,
		},
	}
	
	taskBars := engine.LayoutMultiDayTasks(tasks)
	
	// Should have 2 task bars
	if len(taskBars) != 2 {
		t.Errorf("Expected 2 task bars, got %d", len(taskBars))
	}
	
	// Check task bar properties
	for _, bar := range taskBars {
		if bar.TaskID == "" {
			t.Error("Task bar should have TaskID")
		}
		
		if bar.Width <= 0 {
			t.Error("Task bar should have positive width")
		}
		
		if bar.Height <= 0 {
			t.Error("Task bar should have positive height")
		}
		
		if bar.Color == "" {
			t.Error("Task bar should have color")
		}
	}
}

func TestCalculateXPosition(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Test X position calculation
	testDate := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	xPos := engine.calculateXPosition(testDate)
	
	// Should be 4 days * 20.0 width = 80.0
	expectedX := 4.0 * 20.0
	if xPos != expectedX {
		t.Errorf("Expected X position %f, got %f", expectedX, xPos)
	}
}

func TestCalculateYPosition(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Test Y position calculation for different rows
	yPos0 := engine.calculateYPosition(0, 3) // Row 0 of 3 total rows
	yPos1 := engine.calculateYPosition(1, 3) // Row 1 of 3 total rows
	yPos2 := engine.calculateYPosition(2, 3) // Row 2 of 3 total rows
	
	// Should be evenly distributed
	expectedSpacing := 15.0 / 4.0 // dayHeight / (totalRows + 1)
	
	if yPos0 != expectedSpacing {
		t.Errorf("Expected Y position for row 0 %f, got %f", expectedSpacing, yPos0)
	}
	
	if yPos1 != expectedSpacing*2 {
		t.Errorf("Expected Y position for row 1 %f, got %f", expectedSpacing*2, yPos1)
	}
	
	if yPos2 != expectedSpacing*3 {
		t.Errorf("Expected Y position for row 2 %f, got %f", expectedSpacing*3, yPos2)
	}
}

func TestMonthBoundaryHandling(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Create task that spans across month boundary
	task := &data.Task{
		ID:        "task1",
		Name:      "Cross-Month Task",
		StartDate: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
		Category:  "PROPOSAL",
		Priority:  1,
	}
	
	taskBars := engine.LayoutMultiDayTasks([]*data.Task{task})
	
	// Should have 1 task bar initially
	if len(taskBars) != 1 {
		t.Errorf("Expected 1 task bar, got %d", len(taskBars))
	}
	
	// Handle month boundary
	processedBars := engine.HandleMonthBoundary(taskBars)
	
	// Should be split into 2 segments
	if len(processedBars) != 2 {
		t.Errorf("Expected 2 task bar segments after month boundary handling, got %d", len(processedBars))
	}
	
	// Check that segments are properly marked
	firstSegment := processedBars[0]
	secondSegment := processedBars[1]
	
	if !firstSegment.IsStart {
		t.Error("First segment should be marked as start")
	}
	
	if !secondSegment.IsEnd {
		t.Error("Second segment should be marked as end")
	}
	
	if !secondSegment.IsContinuation {
		t.Error("Second segment should be marked as continuation")
	}
}

func TestTaskOverlapDetection(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Test overlapping tasks
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	
	task3 := &data.Task{
		ID:        "task3",
		StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
	}
	
	// Tasks 1 and 2 should overlap
	if !engine.tasksOverlapDirect(task1, task2) {
		t.Error("Tasks 1 and 2 should overlap")
	}
	
	// Tasks 1 and 3 should not overlap
	if engine.tasksOverlapDirect(task1, task3) {
		t.Error("Tasks 1 and 3 should not overlap")
	}
	
	// Tasks 2 and 3 should not overlap
	if engine.tasksOverlapDirect(task2, task3) {
		t.Error("Tasks 2 and 3 should not overlap")
	}
}

func TestLayoutValidation(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Create task bars with potential issues
	taskBars := []*TaskBar{
		{
			TaskID:    "task1",
			StartX:    0,
			EndX:      100,
			Y:         10,
			Row:       0,
		},
		{
			TaskID:    "task2",
			StartX:    50, // Overlaps with task1
			EndX:      150,
			Y:         10,
			Row:       0, // Same row as task1
		},
		{
			TaskID:    "task3",
			StartX:    -10, // Extends beyond left bound
			EndX:      50,
			Y:         20,
			Row:       1,
		},
	}
	
	issues := engine.ValidateLayout(taskBars)
	
	// Should detect overlapping bars and out-of-bounds bar
	if len(issues) < 2 {
		t.Errorf("Expected at least 2 validation issues, got %d", len(issues))
	}
	
	// Check for specific issues
	foundOverlap := false
	foundOutOfBounds := false
	
	for _, issue := range issues {
		if contains(issue, "overlap") {
			foundOverlap = true
		}
		if contains(issue, "bounds") {
			foundOutOfBounds = true
		}
	}
	
	if !foundOverlap {
		t.Error("Should detect overlapping task bars")
	}
	
	if !foundOutOfBounds {
		t.Error("Should detect out-of-bounds task bars")
	}
}

func TestColorConversion(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Test color conversion
	testCases := []struct {
		hexColor    string
		expectedLaTeX string
	}{
		{"#4A90E2", "blue"},
		{"#F5A623", "orange"},
		{"#7ED321", "green"},
		{"#BD10E0", "purple"},
		{"#D0021B", "red"},
		{"#50E3C2", "teal"},
		{"#B8E986", "lime"},
		{"#CCCCCC", "gray"},
		{"#UNKNOWN", "gray"}, // Unknown color should default to gray
	}
	
	for _, tc := range testCases {
		result := engine.convertColorToLaTeX(tc.hexColor)
		if result != tc.expectedLaTeX {
			t.Errorf("Expected color %s for hex %s, got %s", tc.expectedLaTeX, tc.hexColor, result)
		}
	}
}

func TestBarsOverlap(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	engine := NewMultiDayLayoutEngine(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Test overlapping bars
	bar1 := &TaskBar{StartX: 0, EndX: 100}
	bar2 := &TaskBar{StartX: 50, EndX: 150} // Overlaps with bar1
	
	if !engine.barsOverlap(bar1, bar2) {
		t.Error("Bars 1 and 2 should overlap")
	}
	
	// Test non-overlapping bars
	bar3 := &TaskBar{StartX: 0, EndX: 50}
	bar4 := &TaskBar{StartX: 60, EndX: 110} // No overlap with bar3
	
	if engine.barsOverlap(bar3, bar4) {
		t.Error("Bars 3 and 4 should not overlap")
	}
	
	// Test adjacent bars (touching but not overlapping)
	bar5 := &TaskBar{StartX: 0, EndX: 50}
	bar6 := &TaskBar{StartX: 50, EndX: 100} // Touches bar5 but doesn't overlap
	
	if engine.barsOverlap(bar5, bar6) {
		t.Error("Bars 5 and 6 should not overlap (they only touch)")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
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
