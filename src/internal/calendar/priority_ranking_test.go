package calendar

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/data"
)

func TestPriorityRanker(t *testing.T) {
	// Create test calendar range and overlap detector
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	
	// Create priority ranker
	ranker := NewPriorityRanker(conflictCategorizer)
	
	// Test basic properties
	if ranker.conflictCategorizer != conflictCategorizer {
		t.Error("Expected conflict categorizer to be set")
	}
	
	if len(ranker.rankingRules) == 0 {
		t.Error("Expected ranking rules to be initialized")
	}
	
	if len(ranker.visualWeights) == 0 {
		t.Error("Expected visual weights to be initialized")
	}
}

func TestRankTasks(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	// Create test tasks with various priority scenarios
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Critical Dissertation Task",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:  "DISSERTATION",
			Priority:  5,
			Assignee:  "John Doe",
			IsMilestone: true,
		},
		{
			ID:        "task2",
			Name:      "High Priority Proposal Task",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  4,
			Assignee:  "John Doe", // Same assignee - conflict
		},
		{
			ID:        "task3",
			Name:      "Medium Priority Laser Task",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			Category:  "LASER",
			Priority:  2,
			Assignee:  "Jane Doe",
		},
		{
			ID:        "task4",
			Name:      "Low Priority Admin Task",
			StartDate: time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 24, 0, 0, 0, 0, time.UTC),
			Category:  "ADMIN",
			Priority:  1,
			Assignee:  "Bob Smith",
		},
	}
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      calendarStart,
		CalendarEnd:        calendarEnd,
		CurrentTime:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		AssigneeWorkloads:  map[string]int{"John Doe": 3, "Jane Doe": 2, "Bob Smith": 1},
		CategoryImportance: map[string]float64{"DISSERTATION": 10.0, "PROPOSAL": 8.0, "LASER": 5.0, "ADMIN": 1.0},
	}
	
	// Rank tasks
	ranking := ranker.RankTasks(tasks, context)
	
	// Verify ranking structure
	if len(ranking.TaskPriorities) != len(tasks) {
		t.Errorf("Expected %d task priorities, got %d", len(tasks), len(ranking.TaskPriorities))
	}
	
	if len(ranking.RankingSummary) == 0 {
		t.Error("Expected ranking summary to be populated")
	}
	
	if len(ranking.Recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
	
	// Verify tasks are sorted by priority (highest first)
	for i := 1; i < len(ranking.TaskPriorities); i++ {
		if ranking.TaskPriorities[i-1].PriorityScore < ranking.TaskPriorities[i].PriorityScore {
			t.Error("Expected tasks to be sorted by priority score (highest first)")
		}
	}
	
	// Verify display order is assigned
	for i, taskPriority := range ranking.TaskPriorities {
		if taskPriority.DisplayOrder != i+1 {
			t.Errorf("Expected display order %d, got %d", i+1, taskPriority.DisplayOrder)
		}
	}
}

func TestCalculateTaskPriority(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	// Create test task
	task := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Category:  "DISSERTATION",
		Priority:  4,
		Assignee:  "John Doe",
		IsMilestone: true,
	}
	
	// Create context
	context := &PriorityContext{
		Task:              task,
		Conflicts:         []*CategorizedConflict{},
		CalendarStart:     calendarStart,
		CalendarEnd:       calendarEnd,
		CurrentTime:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		AssigneeWorkloads: map[string]int{"John Doe": 3},
		CategoryImportance: map[string]float64{"DISSERTATION": 10.0},
	}
	
	// Calculate priority
	taskPriority := ranker.calculateTaskPriority(task, context)
	
	// Verify basic properties
	if taskPriority.Task != task {
		t.Error("Expected task to be set")
	}
	
	if taskPriority.PriorityScore <= 0 {
		t.Error("Expected positive priority score")
	}
	
	if taskPriority.VisualProminence == "" {
		t.Error("Expected visual prominence to be determined")
	}
	
	if len(taskPriority.RankingFactors) == 0 {
		t.Error("Expected ranking factors to be calculated")
	}
	
	if taskPriority.DisplayOrder != 0 {
		t.Error("Expected display order to be 0 initially")
	}
	
	// Verify visual style is generated
	if taskPriority.VisualStyle.BorderColor == "" {
		t.Error("Expected visual style to be generated")
	}
	
	// Verify recommendations are generated
	if len(taskPriority.Recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
}

func TestCalculateConflictPriority(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	// Create test task
	task := &data.Task{
		ID:        "task1",
		Priority:  3,
	}
	
	// Create context with conflicts
	conflicts := []*CategorizedConflict{
		{
			TaskOverlap: &TaskOverlap{
				Severity: SeverityCritical,
				Duration: time.Hour * 24 * 5,
			},
			Impact: "HIGH",
			RiskLevel: "HIGH",
		},
		{
			TaskOverlap: &TaskOverlap{
				Severity: SeverityMedium,
				Duration: time.Hour * 24 * 2,
			},
			Impact: "MEDIUM",
			RiskLevel: "MEDIUM",
		},
	}
	
	context := &PriorityContext{
		Task:      task,
		Conflicts: conflicts,
		CurrentTime: time.Now(),
	}
	
	// Calculate conflict priority
	score := ranker.calculateConflictPriority(task, context)
	
	// Should have positive score
	if score <= 0 {
		t.Error("Expected positive conflict priority score")
	}
	
	// Should be higher than base task priority due to conflicts
	expectedMin := float64(task.Priority) * 0.2 + 10.0 + 3.0 + 5.0 + 2.0 // Base + critical + high impact + high risk + medium
	if score < expectedMin {
		t.Errorf("Expected score >= %f, got %f", expectedMin, score)
	}
}

func TestCalculateTimelinePriority(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	
	// Test urgent task (starts tomorrow, ends in 2 days)
	task := &data.Task{
		StartDate: now.AddDate(0, 0, 1), // Tomorrow
		EndDate:   now.AddDate(0, 0, 3), // In 3 days
	}
	
	context := &PriorityContext{
		CurrentTime: now,
	}
	
	score := ranker.calculateTimelinePriority(task, context)
	
	// Should have high score for urgent task
	if score < 15.0 { // 8.0 (starts tomorrow) + 5.0 (due in 3 days)
		t.Errorf("Expected high timeline priority score, got %f", score)
	}
	
	// Test non-urgent task (starts in 30 days, ends in 35 days)
	task.StartDate = now.AddDate(0, 0, 30)
	task.EndDate = now.AddDate(0, 0, 35)
	
	score = ranker.calculateTimelinePriority(task, context)
	
	// Should have low score for non-urgent task
	if score > 5.0 {
		t.Errorf("Expected low timeline priority score, got %f", score)
	}
}

func TestCalculateResourcePriority(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	// Test task with assignee and workload
	task := &data.Task{
		Assignee: "John Doe",
	}
	
	context := &PriorityContext{
		AssigneeWorkloads: map[string]int{"John Doe": 5},
		Conflicts: []*CategorizedConflict{
			{
				Category: CategoryAssigneeConflict,
			},
		},
	}
	
	score := ranker.calculateResourcePriority(task, context)
	
	// Should have positive score
	if score <= 0 {
		t.Error("Expected positive resource priority score")
	}
	
	// Should include workload and conflict scores
	expectedMin := 2.5 + 5.0 // Workload (5 * 0.5) + conflict
	if score < expectedMin {
		t.Errorf("Expected score >= %f, got %f", expectedMin, score)
	}
}

func TestCalculateMilestonePriority(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	
	// Test milestone task
	task := &data.Task{
		IsMilestone: true,
		Priority:    4,
		EndDate:     now.AddDate(0, 0, 5), // Due in 5 days
	}
	
	context := &PriorityContext{
		CurrentTime: now,
	}
	
	score := ranker.calculateMilestonePriority(task, context)
	
	// Should have high score for milestone
	if score < 20.0 { // 15.0 (milestone) + 5.0 (high priority) + 10.0 (due soon)
		t.Errorf("Expected high milestone priority score, got %f", score)
	}
	
	// Test non-milestone task
	task.IsMilestone = false
	task.Priority = 1
	task.EndDate = now.AddDate(0, 0, 30) // Due in 30 days
	
	score = ranker.calculateMilestonePriority(task, context)
	
	// Should have zero score for non-milestone
	if score != 0.0 {
		t.Errorf("Expected zero milestone priority score, got %f", score)
	}
}

func TestDetermineVisualProminence(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	// Test critical prominence
	taskPriority := &TaskPriority{PriorityScore: 20.0}
	prominence := ranker.determineVisualProminence(taskPriority)
	if prominence != ProminenceCritical {
		t.Errorf("Expected CRITICAL prominence, got %s", prominence)
	}
	
	// Test high prominence
	taskPriority.PriorityScore = 12.0
	prominence = ranker.determineVisualProminence(taskPriority)
	if prominence != ProminenceHigh {
		t.Errorf("Expected HIGH prominence, got %s", prominence)
	}
	
	// Test medium prominence
	taskPriority.PriorityScore = 8.0
	prominence = ranker.determineVisualProminence(taskPriority)
	if prominence != ProminenceMedium {
		t.Errorf("Expected MEDIUM prominence, got %s", prominence)
	}
	
	// Test low prominence
	taskPriority.PriorityScore = 4.0
	prominence = ranker.determineVisualProminence(taskPriority)
	if prominence != ProminenceLow {
		t.Errorf("Expected LOW prominence, got %s", prominence)
	}
	
	// Test minimal prominence
	taskPriority.PriorityScore = 1.0
	prominence = ranker.determineVisualProminence(taskPriority)
	if prominence != ProminenceMinimal {
		t.Errorf("Expected MINIMAL prominence, got %s", prominence)
	}
}

func TestGenerateVisualStyle(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	// Test critical style
	taskPriority := &TaskPriority{VisualProminence: ProminenceCritical}
	style := ranker.generateVisualStyle(taskPriority)
	
	if style.BorderColor != "red" {
		t.Errorf("Expected red border color, got %s", style.BorderColor)
	}
	
	if !style.Highlight {
		t.Error("Expected highlight to be true")
	}
	
	if !style.Blink {
		t.Error("Expected blink to be true")
	}
	
	if !style.Glow {
		t.Error("Expected glow to be true")
	}
	
	// Test minimal style
	taskPriority.VisualProminence = ProminenceMinimal
	style = ranker.generateVisualStyle(taskPriority)
	
	if style.BorderColor != "gray" {
		t.Errorf("Expected gray border color, got %s", style.BorderColor)
	}
	
	if style.Opacity != 0.7 {
		t.Errorf("Expected opacity 0.7, got %f", style.Opacity)
	}
	
	if style.Highlight {
		t.Error("Expected highlight to be false")
	}
}

func TestAddCustomPriorityRule(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	initialRuleCount := len(ranker.rankingRules)
	
	// Add custom rule
	customRule := PriorityRule{
		Name:        "Custom Test Rule",
		Description: "Test custom rule",
		Weight:      0.5,
		Calculator: func(task *data.Task, context *PriorityContext) float64 {
			return 10.0
		},
		Category: CategoryTaskPriority,
	}
	
	ranker.AddCustomRule(customRule)
	
	// Verify rule was added
	if len(ranker.rankingRules) != initialRuleCount+1 {
		t.Errorf("Expected %d rules, got %d", initialRuleCount+1, len(ranker.rankingRules))
	}
	
	// Verify rule is at the beginning (highest weight)
	if ranker.rankingRules[0].Name != "Custom Test Rule" {
		t.Error("Expected custom rule to be at highest priority")
	}
}

func TestPriorityRankingMethods(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	ranker := NewPriorityRanker(conflictCategorizer)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:  "DISSERTATION",
			Priority:  5,
			IsMilestone: true,
		},
		{
			ID:        "task2",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  3,
		},
	}
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      calendarStart,
		CalendarEnd:        calendarEnd,
		CurrentTime:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		AssigneeWorkloads:  map[string]int{},
		CategoryImportance: map[string]float64{},
	}
	
	// Rank tasks
	ranking := ranker.RankTasks(tasks, context)
	
	// Test filtering methods
	criticalTasks := ranking.GetTasksByProminence(ProminenceCritical)
	highTasks := ranking.GetTasksByProminence(ProminenceHigh)
	topTasks := ranking.GetTopTasks(1)
	
	// Test summary method
	summary := ranking.GetSummary()
	if summary == "" {
		t.Error("Expected non-empty summary")
	}
	
	// Verify counts are reasonable
	if len(criticalTasks) < 0 {
		t.Error("Expected non-negative critical tasks count")
	}
	
	if len(highTasks) < 0 {
		t.Error("Expected non-negative high tasks count")
	}
	
	if len(topTasks) < 0 {
		t.Error("Expected non-negative top tasks count")
	}
	
	// Top tasks should not exceed total tasks
	if len(topTasks) > len(ranking.TaskPriorities) {
		t.Error("Expected top tasks count to not exceed total tasks")
	}
}
