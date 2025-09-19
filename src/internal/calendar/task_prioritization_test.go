package calendar

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/data"
)

func TestTaskPrioritizationEngine(t *testing.T) {
	// Create test calendar range
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Create task prioritization engine
	engine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	
	// Test basic properties
	if engine.verticalStackingEngine != verticalStackingEngine {
		t.Error("Expected vertical stacking engine to be set")
	}
	
	if engine.priorityRanker != priorityRanker {
		t.Error("Expected priority ranker to be set")
	}
	
	if engine.visibilityManager == nil {
		t.Error("Expected visibility manager to be initialized")
	}
	
	if engine.stackingOptimizer == nil {
		t.Error("Expected stacking optimizer to be initialized")
	}
}

func TestPrioritizeTasks(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	engine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	
	// Create test tasks with different priorities
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
			Assignee:  "John Doe",
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
			Name:      "Low Priority Task",
			StartDate: time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 24, 0, 0, 0, 0, time.UTC),
			Category:  "OTHER",
			Priority:  1,
			Assignee:  "Jane Doe",
		},
	}
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      calendarStart,
		CalendarEnd:        calendarEnd,
		CurrentTime:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		AssigneeWorkloads:  map[string]int{"John Doe": 2, "Jane Doe": 2},
		CategoryImportance: map[string]float64{"DISSERTATION": 10.0, "PROPOSAL": 8.0, "LASER": 5.0, "OTHER": 1.0},
	}
	
	// Prioritize tasks
	result := engine.PrioritizeTasks(tasks, context)
	
	// Verify result structure
	if result == nil {
		t.Error("Expected non-nil prioritization result")
	}
	
	if len(result.PrioritizedTasks) != len(tasks) {
		t.Errorf("Expected %d prioritized tasks, got %d", len(tasks), len(result.PrioritizedTasks))
	}
	
	if result.VisibilityAnalysis == nil {
		t.Error("Expected visibility analysis to be set")
	}
	
	if result.StackingOptimization == nil {
		t.Error("Expected stacking optimization to be set")
	}
	
	// Verify tasks are sorted by prominence score
	for i := 1; i < len(result.PrioritizedTasks); i++ {
		if result.PrioritizedTasks[i-1].ProminenceScore < result.PrioritizedTasks[i].ProminenceScore {
			t.Error("Expected tasks to be sorted by prominence score (highest first)")
		}
	}
	
	// Verify display order is set
	for i, task := range result.PrioritizedTasks {
		if task.DisplayOrder != i {
			t.Errorf("Expected display order %d, got %d", i, task.DisplayOrder)
		}
	}
	
	// Verify analysis date is set
	if result.AnalysisDate.IsZero() {
		t.Error("Expected analysis date to be set")
	}
}

func TestCalculateTaskVisualWeight(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	engine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	
	// Create test task
	task := &data.Task{
		ID:        "task1",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Category:  "DISSERTATION",
		Priority:  5,
		IsMilestone: true,
	}
	
	// Create priority and visibility
	priority := &TaskPriority{
		Task:              task,
		PriorityScore:     15.0,
		VisualProminence:  ProminenceCritical,
		RankingFactors:    make(map[PriorityCategory]float64),
		Recommendations:   make([]string, 0),
	}
	
	visibility := &VisibilityAction{
		IsVisible:       true,
		ProminenceLevel: ProminenceCritical,
		DisplayOrder:    0,
		VisualWeight:    2.0,
		CollapseLevel:   0,
		HighlightLevel:  2,
	}
	
	// Calculate visual weight
	weight := engine.calculateVisualWeight(task, priority, visibility)
	
	// Verify weight calculation
	if weight <= 0 {
		t.Error("Expected positive visual weight")
	}
	
	// Weight should be higher for milestone tasks
	if !task.IsMilestone || weight < 5.0 {
		t.Error("Expected higher weight for milestone tasks")
	}
}

func TestCalculateProminenceScore(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	engine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	
	// Create test task
	task := &data.Task{
		ID:        "task1",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Category:  "DISSERTATION",
		Priority:  5,
		IsMilestone: true,
	}
	
	// Create priority and visibility
	priority := &TaskPriority{
		Task:              task,
		PriorityScore:     15.0,
		VisualProminence:  ProminenceCritical,
		RankingFactors:    make(map[PriorityCategory]float64),
		Recommendations:   make([]string, 0),
		TimelineUrgency:   5.0,
	}
	
	visibility := &VisibilityAction{
		IsVisible:       true,
		ProminenceLevel: ProminenceCritical,
		DisplayOrder:    0,
		VisualWeight:    2.0,
		CollapseLevel:   0,
		HighlightLevel:  2,
	}
	
	// Calculate prominence score
	score := engine.calculateProminenceScore(task, priority, visibility)
	
	// Verify score calculation
	if score <= 0 {
		t.Error("Expected positive prominence score")
	}
	
	// Score should be high for critical milestone tasks
	if score < 100.0 {
		t.Error("Expected high prominence score for critical milestone tasks")
	}
}

func TestDetermineGrouping(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	engine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      calendarStart,
		CalendarEnd:        calendarEnd,
		CurrentTime:        time.Now(),
		AssigneeWorkloads:  make(map[string]int),
		CategoryImportance: make(map[string]float64),
	}
	
	// Test different grouping scenarios
	testCases := []struct {
		task     *data.Task
		priority *TaskPriority
		expected string
	}{
		{
			task: &data.Task{
				ID:       "task1",
				Category: "DISSERTATION",
			},
			priority: nil,
			expected: "category_DISSERTATION",
		},
		{
			task: &data.Task{
				ID:       "task2",
				Category: "",
				Assignee: "John Doe",
			},
			priority: &TaskPriority{
				VisualProminence: ProminenceHigh,
			},
			expected: "priority_HIGH",
		},
		{
			task: &data.Task{
				ID:       "task3",
				Category: "",
				Assignee: "Jane Doe",
			},
			priority: nil,
			expected: "assignee_Jane Doe",
		},
		{
			task: &data.Task{
				ID:       "task4",
				Category: "",
				Assignee: "",
			},
			priority: nil,
			expected: "default",
		},
	}
	
	for _, tc := range testCases {
		groupingID := engine.determineGrouping(tc.task, tc.priority, context)
		if groupingID != tc.expected {
			t.Errorf("Expected grouping ID %s, got %s", tc.expected, groupingID)
		}
	}
}

func TestVisibilityManager(t *testing.T) {
	vm := NewVisibilityManager()
	
	// Test basic properties
	if vm.visibilityRules == nil {
		t.Error("Expected visibility rules to be initialized")
	}
	
	if vm.prominenceWeights == nil {
		t.Error("Expected prominence weights to be initialized")
	}
	
	if vm.visibilityThreshold != 0.5 {
		t.Error("Expected visibility threshold to be 0.5")
	}
	
	if !vm.adaptiveVisibility {
		t.Error("Expected adaptive visibility to be enabled")
	}
}

func TestAnalyzeVisibility(t *testing.T) {
	vm := NewVisibilityManager()
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Critical Task",
			Category:  "DISSERTATION",
			Priority:  5,
			IsMilestone: true,
		},
		{
			ID:       "task2",
			Name:     "High Priority Task",
			Category: "PROPOSAL",
			Priority: 4,
		},
		{
			ID:       "task3",
			Name:     "Low Priority Task",
			Category: "OTHER",
			Priority: 1,
		},
	}
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:        time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		CurrentTime:        time.Now(),
		AssigneeWorkloads:  make(map[string]int),
		CategoryImportance: make(map[string]float64),
	}
	
	// Analyze visibility
	analysis := vm.AnalyzeVisibility(tasks, context)
	
	// Verify analysis structure
	if analysis == nil {
		t.Error("Expected non-nil visibility analysis")
	}
	
	if analysis.TotalTasks != len(tasks) {
		t.Errorf("Expected %d total tasks, got %d", len(tasks), analysis.TotalTasks)
	}
	
	if analysis.VisibleTasks < 0 || analysis.VisibleTasks > analysis.TotalTasks {
		t.Error("Expected visible tasks to be between 0 and total tasks")
	}
	
	if analysis.HiddenTasks < 0 || analysis.HiddenTasks > analysis.TotalTasks {
		t.Error("Expected hidden tasks to be between 0 and total tasks")
	}
	
	if analysis.VisibilityScore < 0 || analysis.VisibilityScore > 1 {
		t.Error("Expected visibility score to be between 0 and 1")
	}
	
	// Verify prominence distribution is initialized
	for _, prominence := range []VisualProminence{
		ProminenceCritical, ProminenceHigh, ProminenceMedium, 
		ProminenceLow, ProminenceMinimal,
	} {
		if _, exists := analysis.ProminenceDistribution[prominence]; !exists {
			t.Errorf("Expected prominence distribution to include %s", prominence)
		}
	}
}

func TestDetermineVisibility(t *testing.T) {
	vm := NewVisibilityManager()
	
	// Create test task
	task := &data.Task{
		ID:        "task1",
		Name:      "Test Task",
		Category:  "DISSERTATION",
		Priority:  5,
		IsMilestone: true,
	}
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:        time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		CurrentTime:        time.Now(),
		AssigneeWorkloads:  make(map[string]int),
		CategoryImportance: make(map[string]float64),
	}
	
	// Determine visibility
	visibility := vm.DetermineVisibility(task, context)
	
	// Verify visibility structure
	if visibility == nil {
		t.Error("Expected non-nil visibility action")
	}
	
	if visibility.ProminenceLevel == "" {
		t.Error("Expected prominence level to be set")
	}
	
	if visibility.VisualWeight <= 0 {
		t.Error("Expected positive visual weight")
	}
}

func TestStackingOptimizer(t *testing.T) {
	so := NewStackingOptimizer()
	
	// Test basic properties
	if so.optimizationRules == nil {
		t.Error("Expected optimization rules to be initialized")
	}
	
	if so.stackingStrategies == nil {
		t.Error("Expected stacking strategies to be initialized")
	}
	
	if !so.adaptiveOrdering {
		t.Error("Expected adaptive ordering to be enabled")
	}
}

func TestOptimizeStackingOrder(t *testing.T) {
	so := NewStackingOptimizer()
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:       "task1",
			Name:     "Low Priority Task",
			Priority: 1,
		},
		{
			ID:       "task2",
			Name:     "High Priority Task",
			Priority: 5,
		},
		{
			ID:       "task3",
			Name:     "Medium Priority Task",
			Priority: 3,
		},
	}
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:        time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		CurrentTime:        time.Now(),
		AssigneeWorkloads:  make(map[string]int),
		CategoryImportance: make(map[string]float64),
	}
	
	// Optimize stacking order
	optimization := so.OptimizeStackingOrder(tasks, context)
	
	// Verify optimization structure
	if optimization == nil {
		t.Error("Expected non-nil stacking optimization")
	}
	
	if len(optimization.OriginalOrder) != len(tasks) {
		t.Errorf("Expected %d original tasks, got %d", len(tasks), len(optimization.OriginalOrder))
	}
	
	if len(optimization.OptimizedOrder) != len(tasks) {
		t.Errorf("Expected %d optimized tasks, got %d", len(tasks), len(optimization.OptimizedOrder))
	}
	
	if optimization.VisualImprovement < 0 || optimization.VisualImprovement > 1 {
		t.Error("Expected visual improvement to be between 0 and 1")
	}
	
	if optimization.SpaceEfficiency < 0 || optimization.SpaceEfficiency > 1 {
		t.Error("Expected space efficiency to be between 0 and 1")
	}
}

func TestTaskPrioritizationResultMethods(t *testing.T) {
	// Create test prioritized tasks
	prioritizedTasks := []*PrioritizedTask{
		{
			Task: &data.Task{
				ID: "task1",
				Category: "DISSERTATION",
			},
			Visibility: &VisibilityAction{
				ProminenceLevel: ProminenceCritical,
			},
			GroupingID: "category_DISSERTATION",
		},
		{
			Task: &data.Task{
				ID: "task2",
				Category: "PROPOSAL",
			},
			Visibility: &VisibilityAction{
				ProminenceLevel: ProminenceHigh,
			},
			GroupingID: "category_PROPOSAL",
		},
		{
			Task: &data.Task{
				ID: "task3",
				Category: "DISSERTATION",
			},
			Visibility: &VisibilityAction{
				ProminenceLevel: ProminenceMedium,
			},
			GroupingID: "category_DISSERTATION",
			IsHidden: true,
		},
	}
	
	// Create visibility analysis
	visibilityAnalysis := &VisibilityAnalysis{
		TotalTasks:        3,
		VisibleTasks:      2,
		HiddenTasks:       1,
		CollapsedTasks:    0,
		HighlightedTasks:  1,
		ProminenceDistribution: map[VisualProminence]int{
			ProminenceCritical: 1,
			ProminenceHigh:     1,
			ProminenceMedium:   1,
		},
		VisibilityScore: 0.67,
		Recommendations: []string{"Test recommendation"},
	}
	
	// Create stacking optimization
	stackingOptimization := &StackingOptimization{
		OriginalOrder:    []*data.Task{},
		OptimizedOrder:   []*data.Task{},
		OptimizationGains: make(map[string]float64),
		VisualImprovement: 0.1,
		SpaceEfficiency:   0.8,
		Recommendations:   []string{"Test optimization recommendation"},
	}
	
	// Create result
	result := &TaskPrioritizationResult{
		PrioritizedTasks:    prioritizedTasks,
		VisibilityAnalysis:   visibilityAnalysis,
		StackingOptimization: stackingOptimization,
		Recommendations:     []string{"Test recommendation"},
		AnalysisDate:        time.Now(),
	}
	
	// Test filtering methods
	criticalTasks := result.GetTasksByProminence(ProminenceCritical)
	if len(criticalTasks) != 1 {
		t.Errorf("Expected 1 critical task, got %d", len(criticalTasks))
	}
	
	visibleTasks := result.GetVisibleTasks()
	if len(visibleTasks) != 2 {
		t.Errorf("Expected 2 visible tasks, got %d", len(visibleTasks))
	}
	
	dissertationTasks := result.GetTasksByGroup("category_DISSERTATION")
	if len(dissertationTasks) != 2 {
		t.Errorf("Expected 2 dissertation tasks, got %d", len(dissertationTasks))
	}
	
	// Test summary method
	summary := result.GetSummary()
	if summary == "" {
		t.Error("Expected non-empty summary")
	}
	
	// Verify summary contains expected information
	if len(summary) < 50 {
		t.Error("Expected detailed summary")
	}
}
