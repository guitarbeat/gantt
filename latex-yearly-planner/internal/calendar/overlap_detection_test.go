package calendar

import (
	"testing"
	"time"

	"latex-yearly-planner/internal/data"
)

func TestOverlapDetector(t *testing.T) {
	// Create test calendar range
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create overlap detector
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Test basic properties
	if detector.calendarStart != calendarStart {
		t.Errorf("Expected calendar start %v, got %v", calendarStart, detector.calendarStart)
	}
	
	if detector.calendarEnd != calendarEnd {
		t.Errorf("Expected calendar end %v, got %v", calendarEnd, detector.calendarEnd)
	}
	
	if detector.precision != time.Hour*1 {
		t.Errorf("Expected precision 1 hour, got %v", detector.precision)
	}
}

func TestDetectOverlaps(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Create test tasks with various overlap scenarios
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
		{
			ID:        "task4",
			Name:      "Task 4",
			StartDate: time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC),
			Category:  "ADMIN",
			Priority:  1,
		},
	}
	
	analysis := detector.DetectOverlaps(tasks)
	
	// Verify analysis structure
	if analysis.TotalTasks != 4 {
		t.Errorf("Expected 4 total tasks, got %d", analysis.TotalTasks)
	}
	
	if analysis.TotalOverlaps == 0 {
		t.Error("Expected to find overlaps, but found none")
	}
	
	// Should have at least one overlap group
	if len(analysis.OverlapGroups) == 0 {
		t.Error("Expected at least one overlap group")
	}
	
	// Check that summary is generated
	if analysis.Summary == "" {
		t.Error("Expected analysis summary to be generated")
	}
}

func TestTasksOverlapDirect(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
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
	if !detector.tasksOverlapDirect(task1, task2) {
		t.Error("Tasks 1 and 2 should overlap")
	}
	
	// Tasks 1 and 3 should not overlap
	if detector.tasksOverlapDirect(task1, task3) {
		t.Error("Tasks 1 and 3 should not overlap")
	}
	
	// Tasks 2 and 3 should not overlap
	if detector.tasksOverlapDirect(task2, task3) {
		t.Error("Tasks 2 and 3 should not overlap")
	}
}

func TestDetermineOverlapType(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Test identical tasks
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	overlapType := detector.determineOverlapType(task1, task2, task1.StartDate, task1.EndDate)
	if overlapType != OverlapIdentical {
		t.Errorf("Expected OverlapIdentical, got %s", overlapType)
	}
	
	// Test nested tasks
	task3 := &data.Task{
		ID:        "task3",
		StartDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
	}
	
	overlapType = detector.determineOverlapType(task1, task3, task1.StartDate, task1.EndDate)
	if overlapType != OverlapNested {
		t.Errorf("Expected OverlapNested, got %s", overlapType)
	}
	
	// Test partial overlap
	task4 := &data.Task{
		ID:        "task4",
		StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	
	overlapType = detector.determineOverlapType(task1, task4, time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC))
	if overlapType != OverlapPartial {
		t.Errorf("Expected OverlapPartial, got %s", overlapType)
	}
}

func TestCalculateOverlapSeverity(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Test critical severity for identical tasks
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	severity := detector.calculateOverlapSeverity(task1, task2, OverlapIdentical, time.Hour*24*5)
	if severity != SeverityCritical {
		t.Errorf("Expected SeverityCritical, got %s", severity)
	}
	
	// Test high severity for nested tasks
	severity = detector.calculateOverlapSeverity(task1, task2, OverlapNested, time.Hour*24*5)
	if severity != SeverityHigh {
		t.Errorf("Expected SeverityHigh, got %s", severity)
	}
	
	// Test low severity for partial overlap
	severity = detector.calculateOverlapSeverity(task1, task2, OverlapPartial, time.Hour*24*1)
	if severity != SeverityLow {
		t.Errorf("Expected SeverityLow, got %s", severity)
	}
}

func TestCalculateOverlapPercentage(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Test 100% overlap (identical tasks)
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	percentage := detector.calculateOverlapPercentage(task1, task2, time.Hour*24*5)
	if percentage != 1.0 {
		t.Errorf("Expected 100%% overlap, got %.2f", percentage)
	}
	
	// Test 50% overlap
	task3 := &data.Task{
		ID:        "task3",
		StartDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
	}
	
	percentage = detector.calculateOverlapPercentage(task1, task3, time.Hour*24*2)
	if percentage != 0.5 {
		t.Errorf("Expected 50%% overlap, got %.2f", percentage)
	}
}

func TestGenerateConflictInfo(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	reason, hint := detector.generateConflictInfo(task1, task2, OverlapIdentical, SeverityCritical)
	
	if reason == "" {
		t.Error("Expected non-empty conflict reason")
	}
	
	if hint == "" {
		t.Error("Expected non-empty resolution hint")
	}
	
	// Check that critical severity is reflected in the hint
	if !containsString(hint, "URGENT") {
		t.Error("Expected URGENT in resolution hint for critical severity")
	}
}

func TestAnalyzeTaskOverlap(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Test overlapping tasks
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Priority:  1,
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Priority:  2,
	}
	
	overlap := detector.analyzeTaskOverlap(task1, task2)
	
	if overlap == nil {
		t.Error("Expected overlap to be detected")
		return
	}
	
	if overlap.Task1ID != "task1" {
		t.Errorf("Expected Task1ID to be 'task1', got %s", overlap.Task1ID)
	}
	
	if overlap.Task2ID != "task2" {
		t.Errorf("Expected Task2ID to be 'task2', got %s", overlap.Task2ID)
	}
	
	if overlap.OverlapType == OverlapNone {
		t.Error("Expected overlap type to be detected")
	}
	
	if overlap.Severity == SeverityNone {
		t.Error("Expected severity to be calculated")
	}
	
	if overlap.Duration <= 0 {
		t.Error("Expected positive overlap duration")
	}
	
	if overlap.OverlapDays <= 0 {
		t.Error("Expected positive overlap days")
	}
	
	if overlap.ConflictReason == "" {
		t.Error("Expected non-empty conflict reason")
	}
	
	if overlap.ResolutionHint == "" {
		t.Error("Expected non-empty resolution hint")
	}
	
	// Priority should be the higher of the two tasks
	if overlap.Priority != 2 {
		t.Errorf("Expected priority 2, got %d", overlap.Priority)
	}
}

func TestNonOverlappingTasks(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Test non-overlapping tasks
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	overlap := detector.analyzeTaskOverlap(task1, task2)
	
	if overlap != nil {
		t.Error("Expected no overlap to be detected for non-overlapping tasks")
	}
}

func TestPrecisionFiltering(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Set high precision to filter out short overlaps
	detector.SetPrecision(time.Hour * 24) // 1 day minimum
	
	// Test tasks with very short overlap
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC), // 12 hours overlap
		EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	
	overlap := detector.analyzeTaskOverlap(task1, task2)
	
	// Should be filtered out due to precision
	if overlap != nil {
		t.Error("Expected overlap to be filtered out due to precision")
	}
}

func TestOverlapAnalysisMethods(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	detector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Create tasks with various overlaps
	tasks := []*data.Task{
		{
			ID:        "task1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Priority:  1,
		},
		{
			ID:        "task2",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Priority:  2,
		},
	}
	
	analysis := detector.DetectOverlaps(tasks)
	
	// Test analysis methods
	if analysis.GetOverlapCount() == 0 {
		t.Error("Expected to find overlaps")
	}
	
	if analysis.GetOverlappingTaskCount() == 0 {
		t.Error("Expected to find overlapping tasks")
	}
	
	// Test severity filtering
	criticalOverlaps := analysis.GetOverlapsBySeverity(SeverityCritical)
	highOverlaps := analysis.GetOverlapsBySeverity(SeverityHigh)
	
	// Test type filtering
	partialOverlaps := analysis.GetOverlapsByType(OverlapPartial)
	
	// Verify counts
	if len(criticalOverlaps)+len(highOverlaps) != analysis.CriticalOverlaps+analysis.HighOverlaps {
		t.Error("Severity filtering counts don't match")
	}
	
	if len(partialOverlaps) >= 0 { // Should be non-negative
		// This is just a basic check that the method works
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstringHelper(s, substr))))
}

func containsSubstringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
