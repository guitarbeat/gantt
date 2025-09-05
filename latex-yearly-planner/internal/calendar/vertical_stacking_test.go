package calendar

import (
	"testing"
	"time"

	"latex-yearly-planner/internal/data"
)

func TestVerticalStackingEngine(t *testing.T) {
	// Create test calendar range
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create vertical stacking engine
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Test basic properties
	if engine.smartStackingEngine != smartStackingEngine {
		t.Error("Expected smart stacking engine to be set")
	}
	
	if engine.heightCalculator == nil {
		t.Error("Expected height calculator to be initialized")
	}
	
	if engine.positionCalculator == nil {
		t.Error("Expected position calculator to be initialized")
	}
	
	if engine.spaceOptimizer == nil {
		t.Error("Expected space optimizer to be initialized")
	}
}

func TestStackTasksVertically(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
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
	
	// Stack tasks vertically
	result := engine.StackTasksVertically(tasks, context)
	
	// Verify result structure
	if result == nil {
		t.Error("Expected non-nil vertical stacking result")
	}
	
	if len(result.Stacks) == 0 {
		t.Error("Expected stacks to be created")
	}
	
	if result.TotalHeight <= 0 {
		t.Error("Expected positive total height")
	}
	
	// Verify metrics are calculated
	if result.SpaceEfficiency < 0 || result.SpaceEfficiency > 1 {
		t.Error("Expected space efficiency to be between 0 and 1")
	}
	
	if result.VisualBalance < 0 || result.VisualBalance > 1 {
		t.Error("Expected visual balance to be between 0 and 1")
	}
	
	// Verify analysis date is set
	if result.AnalysisDate.IsZero() {
		t.Error("Expected analysis date to be set")
	}
}

func TestCalculateTaskHeight(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
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
	
	// Calculate task height
	height := engine.calculateTaskHeight(task, context)
	
	// Verify height calculation
	if height <= 0 {
		t.Error("Expected positive height")
	}
	
	if height < context.VisualConstraints.MinTaskHeight {
		t.Error("Expected height to be at least minimum height")
	}
	
	if height > context.VisualConstraints.MaxTaskHeight {
		t.Error("Expected height to be at most maximum height")
	}
}

func TestCalculateVisualWeight(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
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
	
	// Calculate visual weight
	weight := engine.calculateVisualWeight(task, context)
	
	// Verify weight calculation
	if weight <= 0 {
		t.Error("Expected positive visual weight")
	}
	
	// Weight should be higher for milestone tasks
	if !task.IsMilestone || weight < 3.0 {
		t.Error("Expected higher weight for milestone tasks")
	}
}

func TestAssessContentComplexity(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Test different task complexities
	testCases := []struct {
		task     *data.Task
		expected string
	}{
		{
			task: &data.Task{
				Name: "Short",
				Category: "LASER",
			},
			expected: "minimal",
		},
		{
			task: &data.Task{
				Name: "Medium Length Task Name",
				Category: "LASER", // Changed from PROPOSAL to avoid complex category
			},
			expected: "normal",
		},
		{
			task: &data.Task{
				Name: "Very Long Task Name That Should Be Considered Complex",
				Category: "DISSERTATION",
				IsMilestone: true,
			},
			expected: "complex",
		},
	}
	
	for _, tc := range testCases {
		complexity := engine.assessContentComplexity(tc.task)
		if complexity != tc.expected {
			t.Errorf("Expected complexity %s, got %s", tc.expected, complexity)
		}
	}
}

func TestCalculateVerticalPosition(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Create test task
	task := &data.Task{
		ID:        "task1",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Category:  "DISSERTATION",
		Priority:  5,
	}
	
	// Create vertically stacked task
	verticallyStackedTask := &VerticallyStackedTask{
		Task:            task,
		CalculatedHeight: 25.0,
		VisualWeight:    2.0,
	}
	
	// Create vertical stack
	stack := &VerticalStack{
		ID:        "stack1",
		StartTime: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		MaxWidth:  100.0,
		Tasks:     []*VerticallyStackedTask{},
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
	
	// Calculate vertical position
	position := engine.calculateVerticalPosition(verticallyStackedTask, stack, 0, context)
	
	// Verify position calculation
	if position == nil {
		t.Error("Expected non-nil position")
	}
	
	if position.X != 0.0 {
		t.Error("Expected X position to be 0")
	}
	
	if position.Y != 0.0 {
		t.Error("Expected Y position to be 0 for first task")
	}
	
	if position.Width != stack.MaxWidth {
		t.Error("Expected width to match stack max width")
	}
	
	if position.Height != verticallyStackedTask.CalculatedHeight {
		t.Error("Expected height to match calculated height")
	}
	
	if position.ZIndex != 1 {
		t.Error("Expected Z-index to be 1 for first task")
	}
	
	if position.StackIndex != 0 {
		t.Error("Expected stack index to be 0")
	}
}

func TestDetermineAlignmentMode(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Create test stack with critical task
	stack := &TaskStack{
		ID: "stack1",
		Tasks: []*StackedTask{
			{
				Task: &data.Task{
					ID: "task1",
					Name: "Critical Task",
				},
			},
		},
	}
	
	// Create context with critical task priority
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
				Task:              stack.Tasks[0].Task,
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
	
	// Determine alignment mode
	alignmentMode := engine.determineAlignmentMode(stack, context)
	
	// Verify alignment mode
	if alignmentMode != AlignmentTop {
		t.Error("Expected top alignment for critical tasks")
	}
}

func TestDetermineDistributionMode(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Create test stack with mixed priorities
	stack := &TaskStack{
		ID: "stack1",
		Tasks: []*StackedTask{
			{
				Task: &data.Task{
					ID: "task1",
					Name: "High Priority Task",
				},
			},
			{
				Task: &data.Task{
					ID: "task2",
					Name: "Low Priority Task",
				},
			},
		},
	}
	
	// Create context with mixed priorities
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
				Task:              stack.Tasks[0].Task,
				PriorityScore:     15.0,
				VisualProminence:  ProminenceHigh,
				RankingFactors:    make(map[PriorityCategory]float64),
				Recommendations:   make([]string, 0),
			},
			"task2": {
				Task:              stack.Tasks[1].Task,
				PriorityScore:     5.0,
				VisualProminence:  ProminenceLow,
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
	
	// Determine distribution mode
	distributionMode := engine.determineDistributionMode(stack, context)
	
	// Verify distribution mode
	if distributionMode != DistributionPriority {
		t.Error("Expected priority distribution for mixed priority tasks")
	}
}

func TestCalculateStackHeight(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Create test stack
	stack := &VerticalStack{
		ID: "stack1",
		Tasks: []*VerticallyStackedTask{
			{
				Position: &VerticalPosition{
					Y:      0.0,
					Height: 25.0,
				},
			},
			{
				Position: &VerticalPosition{
					Y:      27.0, // 25.0 + 2.0 spacing
					Height: 30.0,
				},
			},
		},
	}
	
	// Calculate stack height
	height := engine.calculateStackHeight(stack)
	
	// Verify height calculation
	expectedHeight := 27.0 + 30.0 // Last task Y + height
	if height != expectedHeight {
		t.Errorf("Expected height %f, got %f", expectedHeight, height)
	}
}

func TestCalculateVerticalSpaceEfficiency(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Create test stacks
	stacks := []*VerticalStack{
		{
			Tasks: []*VerticallyStackedTask{
				{
					Position: &VerticalPosition{
						Y:      0.0,
						Height: 50.0,
					},
				},
			},
		},
		{
			Tasks: []*VerticallyStackedTask{
				{
					Position: &VerticalPosition{
						Y:      0.0,
						Height: 30.0,
					},
				},
			},
		},
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

func TestCalculateVisualBalance(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	engine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Create test stacks with balanced visual weights
	stacks := []*VerticalStack{
		{
			Tasks: []*VerticallyStackedTask{
				{
					VisualWeight: 2.0,
				},
				{
					VisualWeight: 2.0,
				},
			},
		},
	}
	
	// Calculate visual balance
	balance := engine.calculateVisualBalance(stacks)
	
	// Verify balance calculation
	if balance < 0 || balance > 1 {
		t.Error("Expected visual balance to be between 0 and 1")
	}
	
	// Balanced weights should result in high balance
	if balance < 0.8 {
		t.Error("Expected high visual balance for balanced weights")
	}
}

func TestVerticalStackingResultMethods(t *testing.T) {
	// Create test stacks
	stacks := []*VerticalStack{
		{
			ID:            "stack1",
			AlignmentMode: AlignmentTop,
			DistributionMode: DistributionPriority,
		},
		{
			ID:            "stack2",
			AlignmentMode: AlignmentCenter,
			DistributionMode: DistributionEven,
		},
	}
	
	// Create result
	result := &VerticalStackingResult{
		Stacks:          stacks,
		TotalHeight:     100.0,
		SpaceEfficiency: 0.8,
		VisualBalance:   0.9,
		CollisionCount:  2,
		OverflowCount:   1,
		CompressionRatio: 0.3,
		AnalysisDate:    time.Now(),
	}
	
	// Test filtering methods
	topStacks := result.GetStacksByAlignment(AlignmentTop)
	if len(topStacks) != 1 {
		t.Errorf("Expected 1 top-aligned stack, got %d", len(topStacks))
	}
	
	priorityStacks := result.GetStacksByDistribution(DistributionPriority)
	if len(priorityStacks) != 1 {
		t.Errorf("Expected 1 priority-distributed stack, got %d", len(priorityStacks))
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
