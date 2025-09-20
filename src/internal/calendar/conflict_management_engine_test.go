package calendar

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/shared"
)

func TestConflictManagementEngine(t *testing.T) {
	// Create test tasks
	task1 := &shared.Task{
		ID:          "task1",
		Name:        "Test Task 1",
		Description: "Test Description 1",
		Category:    "RESEARCH",
		Priority:    3,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		Assignee:    "user1",
		IsMilestone: false,
	}

	task2 := &shared.Task{
		ID:          "task2",
		Name:        "Test Task 2",
		Description: "Test Description 2",
		Category:    "RESEARCH",
		Priority:    4,
		StartDate:   time.Now().AddDate(0, 0, 3),
		EndDate:     time.Now().AddDate(0, 0, 10),
		Assignee:    "user1",
		IsMilestone: true,
	}

	// Create overlap detector
	overlapDetector := &OverlapDetector{
		calendarStart: time.Now().AddDate(0, 0, -30),
		calendarEnd:   time.Now().AddDate(0, 0, 30),
		precision:     time.Hour,
	}

	// Create conflict management engine
	engine := NewConflictManagementEngine(
		overlapDetector,
		nil, // taskPrioritizationEngine
		nil, // stackingEngine
	)

	// Test conflict categorization
	overlapAnalysis := &OverlapAnalysis{
		OverlapGroups: []*OverlapGroup{
			{
				Tasks: []*shared.Task{task1, task2},
				Overlaps: []*TaskOverlap{
					{
						Task1ID:     task1.ID,
						Task2ID:     task2.ID,
						OverlapType: OverlapPartial,
						Severity:    SeverityHigh,
						Duration:    time.Hour * 24 * 4,
						OverlapDays: 4,
					},
				},
			},
		},
	}

	conflictAnalysis := engine.CategorizeConflicts(overlapAnalysis)

	// Verify results
	if conflictAnalysis.TotalConflicts == 0 {
		t.Error("Expected conflicts to be detected")
	}

	if len(conflictAnalysis.CategorizedConflicts) == 0 {
		t.Error("Expected categorized conflicts")
	}

	// Test conflict resolution
	context := &PriorityContext{
		CurrentTime: time.Now(),
		UserID:      "test-user",
	}

	resolutionResult := engine.ResolveConflicts([]*shared.Task{task1, task2}, context)

	if resolutionResult == nil {
		t.Error("Expected conflict resolution result")
	}

	if len(resolutionResult.Recommendations) == 0 {
		t.Error("Expected recommendations")
	}
}

func TestConflictCategorizer(t *testing.T) {
	// Create test tasks
	task1 := &shared.Task{
		ID:          "task1",
		Name:        "Test Task 1",
		Description: "Test Description 1",
		Category:    "RESEARCH",
		Priority:    3,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		Assignee:    "user1",
		IsMilestone: false,
	}

	task2 := &shared.Task{
		ID:          "task2",
		Name:        "Test Task 2",
		Description: "Test Description 2",
		Category:    "RESEARCH",
		Priority:    4,
		StartDate:   time.Now().AddDate(0, 0, 3),
		EndDate:     time.Now().AddDate(0, 0, 10),
		Assignee:    "user1",
		IsMilestone: true,
	}

	// Create overlap detector
	overlapDetector := &OverlapDetector{
		calendarStart: time.Now().AddDate(0, 0, -30),
		calendarEnd:   time.Now().AddDate(0, 0, 30),
		precision:     time.Hour,
	}

	// Create conflict management engine
	engine := NewConflictManagementEngine(overlapDetector, nil, nil)

	// Test categorization
	overlap := &TaskOverlap{
		Task1ID:     task1.ID,
		Task2ID:     task2.ID,
		OverlapType: OverlapPartial,
		Severity:    SeverityHigh,
		Duration:    time.Hour * 24 * 4,
		OverlapDays: 4,
	}

	categorizedConflict := engine.categorizeConflict(overlap, task1, task2)

	if categorizedConflict == nil {
		t.Error("Expected categorized conflict")
	}

	if categorizedConflict.Category == "" {
		t.Error("Expected conflict category")
	}

	if categorizedConflict.RiskLevel == "" {
		t.Error("Expected risk level")
	}

	if categorizedConflict.Urgency == "" {
		t.Error("Expected urgency")
	}

	if categorizedConflict.Complexity == "" {
		t.Error("Expected complexity")
	}
}

func TestConflictResolutionEngine(t *testing.T) {
	// Create test tasks
	task1 := &shared.Task{
		ID:          "task1",
		Name:        "Test Task 1",
		Description: "Test Description 1",
		Category:    "RESEARCH",
		Priority:    3,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		Assignee:    "user1",
		IsMilestone: false,
	}

	task2 := &shared.Task{
		ID:          "task2",
		Name:        "Test Task 2",
		Description: "Test Description 2",
		Category:    "RESEARCH",
		Priority:    4,
		StartDate:   time.Now().AddDate(0, 0, 3),
		EndDate:     time.Now().AddDate(0, 0, 10),
		Assignee:    "user1",
		IsMilestone: true,
	}

	// Create conflict resolution engine
	engine := NewConflictResolutionEngine(nil, nil)

	// Test conflict resolution
	context := &PriorityContext{
		CurrentTime: time.Now(),
		UserID:      "test-user",
	}

	resolutionResult := engine.ResolveConflicts([]*shared.Task{task1, task2}, context)

	if resolutionResult == nil {
		t.Error("Expected conflict resolution result")
	}

	if resolutionResult.AnalysisDate.IsZero() {
		t.Error("Expected analysis date")
	}

	if len(resolutionResult.Recommendations) == 0 {
		t.Error("Expected recommendations")
	}
}
