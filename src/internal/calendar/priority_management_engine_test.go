package calendar

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/shared"
)

func TestPriorityManagementEngine(t *testing.T) {
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

	// Create conflict categorizer
	conflictCategorizer := NewConflictCategorizer(overlapDetector)

	// Create priority management engine
	engine := NewPriorityManagementEngine(
		conflictCategorizer,
		nil, // verticalStackingEngine
	)

	// Test priority calculation
	context := &PriorityContext{
		CurrentTime: time.Now(),
		UserID:      "test-user",
	}

	priorityResult := engine.CalculatePriorityScores([]*shared.Task{task1, task2}, context)

	// Verify results
	if len(priorityResult.TaskScores) != 2 {
		t.Error("Expected 2 task scores")
	}

	if len(priorityResult.RankingOrder) != 2 {
		t.Error("Expected 2 tasks in ranking order")
	}

	if len(priorityResult.Recommendations) == 0 {
		t.Error("Expected recommendations")
	}

	// Test task prioritization
	prioritizationResult := engine.PrioritizeTasks([]*shared.Task{task1, task2}, context)

	// Verify results
	if len(prioritizationResult.PrioritizedTasks) != 2 {
		t.Error("Expected 2 prioritized tasks")
	}

	if len(prioritizationResult.StackingOrder) != 2 {
		t.Error("Expected 2 tasks in stacking order")
	}

	if len(prioritizationResult.Recommendations) == 0 {
		t.Error("Expected recommendations")
	}
}

func TestPriorityRanker(t *testing.T) {
	// Create test task
	task := &shared.Task{
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

	// Create overlap detector
	overlapDetector := &OverlapDetector{
		calendarStart: time.Now().AddDate(0, 0, -30),
		calendarEnd:   time.Now().AddDate(0, 0, 30),
		precision:     time.Hour,
	}

	// Create conflict categorizer
	conflictCategorizer := NewConflictCategorizer(overlapDetector)

	// Create priority ranker
	ranker := NewPriorityRanker(conflictCategorizer)

	// Test priority calculation
	context := &PriorityContext{
		CurrentTime: time.Now(),
		UserID:      "test-user",
	}

	result := ranker.CalculatePriorityScores([]*shared.Task{task}, context)

	// Verify results
	if len(result.TaskScores) != 1 {
		t.Error("Expected 1 task score")
	}

	score := result.TaskScores[0]
	if score.TaskID != task.ID {
		t.Error("Expected correct task ID")
	}

	if score.OverallScore < 0 || score.OverallScore > 1 {
		t.Error("Expected overall score between 0 and 1")
	}

	if score.VisualProminence == "" {
		t.Error("Expected visual prominence")
	}

	if score.Ranking == 0 {
		t.Error("Expected ranking to be set")
	}
}

func TestVisibilityManager(t *testing.T) {
	// Create test task
	task := &shared.Task{
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

	// Create visibility manager
	manager := NewVisibilityManager()

	// Test visibility rules
	context := &PriorityContext{
		CurrentTime: time.Now(),
		UserID:      "test-user",
	}

	visibilitySettings := manager.ApplyVisibilityRules([]*shared.Task{task}, context)

	// Verify results
	if len(visibilitySettings) != 1 {
		t.Error("Expected 1 visibility setting")
	}

	action := visibilitySettings[task.ID]
	if action == nil {
		t.Error("Expected visibility action")
	}

	if !action.IsVisible {
		t.Error("Expected task to be visible")
	}

	if action.ProminenceLevel == "" {
		t.Error("Expected prominence level")
	}

	if action.VisualWeight < 0 || action.VisualWeight > 1 {
		t.Error("Expected visual weight between 0 and 1")
	}
}

func TestStackingOptimizer(t *testing.T) {
	// Create test task
	task := &shared.Task{
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

	// Create stacking optimizer
	optimizer := NewStackingOptimizer()

	// Test stacking optimization
	context := &PriorityContext{
		CurrentTime: time.Now(),
		UserID:      "test-user",
	}

	actions := optimizer.OptimizeStackingOrder([]*shared.Task{task}, context)

	// Verify results
	if len(actions) != 1 {
		t.Error("Expected 1 optimization action")
	}

	action := actions[0]
	if action == nil {
		t.Error("Expected optimization action")
	}

	if action.StackingOrder < 0 {
		t.Error("Expected non-negative stacking order")
	}

	if action.VisualProminence == "" {
		t.Error("Expected visual prominence")
	}

	if action.DisplayPriority < 0 || action.DisplayPriority > 1 {
		t.Error("Expected display priority between 0 and 1")
	}
}
