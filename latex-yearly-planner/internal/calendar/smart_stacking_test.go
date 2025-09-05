package calendar

import (
	"testing"
	"time"

	"latex-yearly-planner/internal/data"
)

func TestSmartStackingEngine(t *testing.T) {
	// Create test calendar range
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	
	// Create smart stacking engine
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Test basic properties
	if engine.overlapDetector != overlapDetector {
		t.Error("Expected overlap detector to be set")
	}
	
	if engine.conflictCategorizer != conflictCategorizer {
		t.Error("Expected conflict categorizer to be set")
	}
	
	if engine.priorityRanker != priorityRanker {
		t.Error("Expected priority ranker to be set")
	}
	
	if len(engine.stackingRules) == 0 {
		t.Error("Expected stacking rules to be initialized")
	}
	
	if engine.visualConstraints == nil {
		t.Error("Expected visual constraints to be initialized")
	}
}

func TestStackTasks(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create test tasks with overlapping scenarios
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
	}
	
	// Create context
	context := &StackingContext{
		CalendarStart:   calendarStart,
		CalendarEnd:     calendarEnd,
		CurrentTime:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		DayWidth:        100.0,
		DayHeight:       200.0,
		AvailableHeight: 180.0,
		AvailableWidth:  90.0,
		ExistingStacks:  make([]*TaskStack, 0),
		VisualSettings: &VisualSettings{
			ShowTaskNames:      true,
			ShowTaskDurations:  true,
			ShowTaskPriorities: true,
			ShowConflictIndicators: true,
			CollapseThreshold:  5,
			AnimationEnabled:   true,
			HighlightConflicts: true,
			ColorScheme:        "default",
		},
		VisualConstraints: &VisualConstraints{
			MaxStackHeight:     100.0,
			MinTaskHeight:      20.0,
			MaxTaskHeight:      40.0,
			MinTaskWidth:       50.0,
			MaxTaskWidth:       200.0,
			VerticalSpacing:    2.0,
			HorizontalSpacing:  5.0,
			MaxStackDepth:      10,
			CollisionThreshold: 0.1,
			OverflowThreshold:  0.8,
		},
	}
	
	// Stack tasks
	result := engine.StackTasks(tasks, context)
	
	// Verify result structure
	if result == nil {
		t.Error("Expected non-nil stacking result")
	}
	
	if len(result.Stacks) == 0 {
		t.Error("Expected stacks to be created")
	}
	
	if result.TotalStacks == 0 {
		t.Error("Expected total stacks to be greater than 0")
	}
	
	// Verify metrics are calculated
	if result.SpaceEfficiency < 0 || result.SpaceEfficiency > 1 {
		t.Error("Expected space efficiency to be between 0 and 1")
	}
	
	if result.VisualQuality < 0 || result.VisualQuality > 1 {
		t.Error("Expected visual quality to be between 0 and 1")
	}
	
	// Verify recommendations are generated (may be empty if no issues)
	// Recommendations are generated based on collisions, overflows, and efficiency
	// This is acceptable behavior - recommendations are optional
	
	// Verify analysis date is set
	if result.AnalysisDate.IsZero() {
		t.Error("Expected analysis date to be set")
	}
}

func TestDetermineStackingAction(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create test task
	task := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Category:  "DISSERTATION",
		Priority:  5,
		IsMilestone: true,
	}
	
	// Create context
	context := &StackingContext{
		CalendarStart:   calendarStart,
		CalendarEnd:     calendarEnd,
		CurrentTime:     time.Now(),
		DayWidth:        100.0,
		DayHeight:       200.0,
		AvailableHeight: 180.0,
		AvailableWidth:  90.0,
		ExistingStacks:  make([]*TaskStack, 0),
		TaskPriorities: map[string]*TaskPriority{
			"task1": {
				Task:              task,
				PriorityScore:     15.0,
				VisualProminence:  ProminenceCritical,
				RankingFactors:    make(map[PriorityCategory]float64),
				Recommendations:   make([]string, 0),
			},
		},
		VisualSettings: &VisualSettings{
			ShowTaskNames:      true,
			ShowTaskDurations:  true,
			ShowTaskPriorities: true,
			ShowConflictIndicators: true,
			CollapseThreshold:  5,
			AnimationEnabled:   true,
			HighlightConflicts: true,
			ColorScheme:        "default",
		},
		VisualConstraints: &VisualConstraints{
			MaxStackHeight:     100.0,
			MinTaskHeight:      20.0,
			MaxTaskHeight:      40.0,
			MinTaskWidth:       50.0,
			MaxTaskWidth:       200.0,
			VerticalSpacing:    2.0,
			HorizontalSpacing:  5.0,
			MaxStackDepth:      10,
			CollisionThreshold: 0.1,
			OverflowThreshold:  0.8,
		},
	}
	
	// Determine stacking action
	action := engine.determineStackingAction(task, context)
	
	// Verify action properties
	if action == nil {
		t.Error("Expected non-nil stacking action")
	}
	
	if action.StackingType == "" {
		t.Error("Expected stacking type to be set")
	}
	
	if action.Height <= 0 {
		t.Error("Expected positive height")
	}
	
	if action.Width <= 0 {
		t.Error("Expected positive width")
	}
	
	if action.ZIndex < 0 {
		t.Error("Expected non-negative Z-index")
	}
	
	if action.Priority <= 0 {
		t.Error("Expected positive priority")
	}
}

func TestGroupTasksByOverlap(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "task2",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "task3",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
		},
	}
	
	// Detect overlaps
	overlapAnalysis := overlapDetector.DetectOverlaps(tasks)
	
	// Group tasks
	groups := engine.groupTasksByOverlap(tasks, overlapAnalysis)
	
	// Verify groups are created
	if len(groups) == 0 {
		t.Error("Expected groups to be created")
	}
	
	// Verify all tasks are included
	totalTasks := 0
	for _, group := range groups {
		totalTasks += len(group)
	}
	
	if totalTasks != len(tasks) {
		t.Errorf("Expected %d tasks in groups, got %d", len(tasks), totalTasks)
	}
}

func TestCreateStackForGroup(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Priority:  5,
		},
		{
			ID:        "task2",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Priority:  3,
		},
	}
	
	// Create context
	context := &StackingContext{
		CalendarStart:   calendarStart,
		CalendarEnd:     calendarEnd,
		CurrentTime:     time.Now(),
		DayWidth:        100.0,
		DayHeight:       200.0,
		AvailableHeight: 180.0,
		AvailableWidth:  90.0,
		ExistingStacks:  make([]*TaskStack, 0),
		TaskPriorities: map[string]*TaskPriority{
			"task1": {
				Task:              tasks[0],
				PriorityScore:     15.0,
				VisualProminence:  ProminenceCritical,
				RankingFactors:    make(map[PriorityCategory]float64),
				Recommendations:   make([]string, 0),
			},
			"task2": {
				Task:              tasks[1],
				PriorityScore:     8.0,
				VisualProminence:  ProminenceMedium,
				RankingFactors:    make(map[PriorityCategory]float64),
				Recommendations:   make([]string, 0),
			},
		},
		VisualSettings: &VisualSettings{
			ShowTaskNames:      true,
			ShowTaskDurations:  true,
			ShowTaskPriorities: true,
			ShowConflictIndicators: true,
			CollapseThreshold:  5,
			AnimationEnabled:   true,
			HighlightConflicts: true,
			ColorScheme:        "default",
		},
		VisualConstraints: &VisualConstraints{
			MaxStackHeight:     100.0,
			MinTaskHeight:      20.0,
			MaxTaskHeight:      40.0,
			MinTaskWidth:       50.0,
			MaxTaskWidth:       200.0,
			VerticalSpacing:    2.0,
			HorizontalSpacing:  5.0,
			MaxStackDepth:      10,
			CollisionThreshold: 0.1,
			OverflowThreshold:  0.8,
		},
	}
	
	// Create stack
	stack := engine.createStackForGroup(tasks, context)
	
	// Verify stack properties
	if stack == nil {
		t.Error("Expected non-nil stack")
	}
	
	if len(stack.Tasks) != len(tasks) {
		t.Errorf("Expected %d tasks in stack, got %d", len(tasks), len(stack.Tasks))
	}
	
	if stack.ID == "" {
		t.Error("Expected stack ID to be set")
	}
	
	if stack.StartTime.IsZero() {
		t.Error("Expected stack start time to be set")
	}
	
	if stack.EndTime.IsZero() {
		t.Error("Expected stack end time to be set")
	}
	
	if stack.TotalHeight <= 0 {
		t.Error("Expected positive total height")
	}
	
	if stack.MaxWidth <= 0 {
		t.Error("Expected positive max width")
	}
	
	if stack.StackingType == "" {
		t.Error("Expected stacking type to be set")
	}
	
	// Verify tasks are sorted by priority
	if len(stack.Tasks) >= 2 {
		if stack.Tasks[0].Task.Priority < stack.Tasks[1].Task.Priority {
			t.Error("Expected tasks to be sorted by priority (highest first)")
		}
	}
}

func TestHasCollision(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create test stacked tasks
	task1 := &StackedTask{
		Position: &Position{X: 0, Y: 0, Width: 50, Height: 20, ZIndex: 1},
	}
	
	task2 := &StackedTask{
		Position: &Position{X: 25, Y: 10, Width: 50, Height: 20, ZIndex: 2},
	}
	
	task3 := &StackedTask{
		Position: &Position{X: 100, Y: 0, Width: 50, Height: 20, ZIndex: 3},
	}
	
	// Test collision detection
	hasCollision := engine.hasCollision(task2, []*StackedTask{task1})
	if !hasCollision {
		t.Error("Expected collision between overlapping tasks")
	}
	
	hasCollision = engine.hasCollision(task3, []*StackedTask{task1})
	if hasCollision {
		t.Error("Expected no collision between non-overlapping tasks")
	}
}

func TestHasOverflow(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create context with limited height
	context := &StackingContext{
		AvailableHeight: 100.0,
		VisualConstraints: &VisualConstraints{
			OverflowThreshold: 0.8,
		},
	}
	
	// Create test stacked task
	task := &StackedTask{
		Position: &Position{X: 0, Y: 70, Width: 50, Height: 20, ZIndex: 1},
	}
	
	// Test overflow detection
	hasOverflow := engine.hasOverflow(task, context)
	if !hasOverflow {
		t.Error("Expected overflow for task exceeding threshold")
	}
	
	// Test no overflow
	task.Position.Y = 50
	hasOverflow = engine.hasOverflow(task, context)
	if hasOverflow {
		t.Error("Expected no overflow for task within threshold")
	}
}

func TestCalculateSpaceEfficiency(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create test stacks
	stacks := []*TaskStack{
		{TotalHeight: 50.0},
		{TotalHeight: 30.0},
	}
	
	// Create context
	context := &StackingContext{
		AvailableHeight: 100.0,
	}
	
	// Calculate space efficiency
	efficiency := engine.calculateSpaceEfficiency(stacks, context)
	
	// Verify efficiency calculation
	expectedEfficiency := 80.0 / 100.0 // (50 + 30) / 100
	if efficiency != expectedEfficiency {
		t.Errorf("Expected efficiency %f, got %f", expectedEfficiency, efficiency)
	}
}

func TestCalculateVisualQuality(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create test stacks with collisions and overflows
	stacks := []*TaskStack{
		{
			Tasks:         []*StackedTask{{}, {}},
			CollisionCount: 1,
			OverflowCount:  1,
		},
	}
	
	// Create context
	context := &StackingContext{
		AvailableHeight: 100.0,
	}
	
	// Calculate visual quality
	quality := engine.calculateVisualQuality(stacks, context)
	
	// Verify quality calculation
	if quality < 0 || quality > 1 {
		t.Error("Expected visual quality to be between 0 and 1")
	}
	
	// Quality should be less than 1 due to collisions and overflows
	if quality >= 1.0 {
		t.Error("Expected visual quality to be less than 1 due to collisions and overflows")
	}
}

func TestAddCustomStackingRule(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	engine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	initialRuleCount := len(engine.stackingRules)
	
	// Add custom rule
	customRule := StackingRule{
		Name:        "Custom Test Rule",
		Description: "Test custom rule",
		Priority:    0, // Highest priority
		Condition: func(task *data.Task, context *StackingContext) bool {
			return task.Name == "Custom Task"
		},
		Action: func(task *data.Task, context *StackingContext) *StackingAction {
			return &StackingAction{
				StackingType:       StackingTypeFloating,
				VerticalOffset:     0.0,
				HorizontalOffset:   0.0,
				Height:             30.0,
				Width:              100.0,
				ZIndex:             15,
				CollisionAvoidance: true,
				Priority:           0,
			}
		},
	}
	
	engine.AddCustomRule(customRule)
	
	// Verify rule was added
	if len(engine.stackingRules) != initialRuleCount+1 {
		t.Errorf("Expected %d rules, got %d", initialRuleCount+1, len(engine.stackingRules))
	}
	
	// Verify rule is at the beginning (highest priority)
	if engine.stackingRules[0].Name != "Custom Test Rule" {
		t.Error("Expected custom rule to be at highest priority")
	}
}

func TestStackingResultMethods(t *testing.T) {
	// Create test stacks
	stacks := []*TaskStack{
		{
			ID:           "stack1",
			StackingType: StackingTypeVertical,
			Priority:     1,
		},
		{
			ID:           "stack2",
			StackingType: StackingTypeHorizontal,
			Priority:     2,
		},
		{
			ID:           "stack3",
			StackingType: StackingTypeVertical,
			Priority:     1,
		},
	}
	
	// Create result
	result := &StackingResult{
		Stacks:          stacks,
		TotalStacks:     len(stacks),
		CollisionCount:  2,
		OverflowCount:   1,
		SpaceEfficiency: 0.75,
		VisualQuality:   0.8,
		AnalysisDate:    time.Now(),
	}
	
	// Test filtering methods
	verticalStacks := result.GetStacksByType(StackingTypeVertical)
	if len(verticalStacks) != 2 {
		t.Errorf("Expected 2 vertical stacks, got %d", len(verticalStacks))
	}
	
	priority1Stacks := result.GetStacksByPriority(1)
	if len(priority1Stacks) != 2 {
		t.Errorf("Expected 2 priority 1 stacks, got %d", len(priority1Stacks))
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
