package calendar

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/data"
)

func TestCalendarGridIntegration(t *testing.T) {
	// Create test configuration
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
	
	// Create integration instance
	integration := NewCalendarGridIntegration(config)
	
	// Verify initialization
	if integration.smartStackingEngine == nil {
		t.Error("Expected smart stacking engine to be initialized")
	}
	
	if integration.verticalStackingEngine == nil {
		t.Error("Expected vertical stacking engine to be initialized")
	}
	
	if integration.taskPrioritizationEngine == nil {
		t.Error("Expected task prioritization engine to be initialized")
	}
	
	if integration.conflictResolutionEngine == nil {
		t.Error("Expected conflict resolution engine to be initialized")
	}
	
	if integration.multiDayLayoutEngine == nil {
		t.Error("Expected multi-day layout engine to be initialized")
	}
	
	if integration.gridConfig == nil {
		t.Error("Expected grid config to be set")
	}
	
	if integration.visualSettings == nil {
		t.Error("Expected visual settings to be initialized")
	}
}

func TestProcessTasksWithSmartStacking(t *testing.T) {
	// Create test configuration
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
	}
	
	// Create integration instance
	integration := NewCalendarGridIntegration(config)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Task 1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
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
			Name:      "MILESTONE: Complete Task 3",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Priority:  3,
		},
	}
	
	// Process tasks
	result, err := integration.ProcessTasksWithSmartStacking(tasks)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Verify result
	if result == nil {
		t.Fatal("Expected result to be non-nil")
	}
	
	if len(result.TaskBars) != len(tasks) {
		t.Errorf("Expected %d task bars, got %d", len(tasks), len(result.TaskBars))
	}
	
	if result.Statistics == nil {
		t.Error("Expected statistics to be non-nil")
	}
	
	if result.Statistics.TotalTasks != len(tasks) {
		t.Errorf("Expected %d total tasks, got %d", len(tasks), result.Statistics.TotalTasks)
	}
	
	if result.AnalysisDate.IsZero() {
		t.Error("Expected analysis date to be set")
	}
}

func TestCreateIntegratedTaskBars(t *testing.T) {
	// Create test configuration
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
	}
	
	// Create integration instance
	integration := NewCalendarGridIntegration(config)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Task 1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  1,
		},
	}
	
	// Create mock results
	stackingResult := &StackingResult{
		Stacks: []*TaskStack{},
	}
	
	verticalStackingResult := &VerticalStackingResult{
		Stacks: []*VerticalStack{},
	}
	
	conflictResolutionResult := &ConflictResolutionResult{}
	
	prioritizationResult := &TaskPrioritizationResult{
		PrioritizedTasks: []*PrioritizedTask{
			{
				Task: &data.Task{ID: "task1"},
				Priority: &TaskPriority{
					PriorityScore:     0.7,
					TimelineUrgency:   0.8,
					ResourceContention: 0.6,
					MilestoneWeight:   0.4,
				},
			},
		},
	}
	
	// Create integrated task bars
	bars := integration.createIntegratedTaskBars(tasks, stackingResult, verticalStackingResult, conflictResolutionResult, prioritizationResult)
	
	// Verify bars
	if len(bars) != len(tasks) {
		t.Errorf("Expected %d bars, got %d", len(tasks), len(bars))
	}
	
	bar := bars[0]
	if bar.TaskID != "task1" {
		t.Errorf("Expected task ID 'task1', got '%s'", bar.TaskID)
	}
	
	if bar.StartX < 0 {
		t.Error("Expected start X to be non-negative")
	}
	
	if bar.Width <= 0 {
		t.Error("Expected width to be positive")
	}
	
	if bar.Height <= 0 {
		t.Error("Expected height to be positive")
	}
	
	if bar.Color == "" {
		t.Error("Expected color to be set")
	}
	
	if bar.VisualWeight <= 0 {
		t.Error("Expected visual weight to be positive")
	}
	
	if bar.ProminenceScore <= 0 {
		t.Error("Expected prominence score to be positive")
	}
}

func TestCalendarGridIntegrationCalculateXPosition(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
	}
	
	integration := NewCalendarGridIntegration(config)
	
	// Test start date
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	x := integration.calculateXPosition(startDate)
	if x != 0 {
		t.Errorf("Expected X position 0 for start date, got %.2f", x)
	}
	
	// Test 5 days later
	date5Days := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)
	x = integration.calculateXPosition(date5Days)
	expectedX := 5.0 * 20.0
	if x != expectedX {
		t.Errorf("Expected X position %.2f for 5 days later, got %.2f", expectedX, x)
	}
}

func TestCalendarGridIntegrationCalculateVisualWeight(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	integration := NewCalendarGridIntegration(config)
	
	// Test task with normal priority
	task := &data.Task{
		ID:        "task1",
		Name:      "Task 1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
		Category:  "PROPOSAL",
		Priority:  1,
	}
	
	priority := &TaskPriority{
		PriorityScore:     0.7,
		TimelineUrgency:   0.8,
		ResourceContention: 0.6,
		MilestoneWeight:   0.4,
	}
	
	weight := integration.calculateVisualWeight(task, priority)
	if weight <= 0 {
		t.Error("Expected visual weight to be positive")
	}
	
	if weight > 1.0 {
		t.Error("Expected visual weight to be <= 1.0")
	}
	
	// Test milestone task
	milestoneTask := &data.Task{
		ID:        "milestone",
		Name:      "MILESTONE: Complete milestone",
		StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		Category:  "IMAGING",
		Priority:  3,
	}
	
	milestoneWeight := integration.calculateVisualWeight(milestoneTask, priority)
	if milestoneWeight <= weight {
		t.Error("Expected milestone task to have higher visual weight")
	}
}

func TestCalendarGridIntegrationCalculateProminenceScore(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	integration := NewCalendarGridIntegration(config)
	
	task := &data.Task{
		ID:        "task1",
		Name:      "Task 1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
		Category:  "PROPOSAL",
		Priority:  1,
	}
	
	priority := &TaskPriority{
		PriorityScore:     0.7,
		TimelineUrgency:   0.8,
		ResourceContention: 0.6,
		MilestoneWeight:   0.4,
	}
	
	visualWeight := 0.7
	prominence := integration.calculateProminenceScore(task, priority, visualWeight)
	
	if prominence <= 0 {
		t.Error("Expected prominence score to be positive")
	}
	
	if prominence > 1.0 {
		t.Error("Expected prominence score to be <= 1.0")
	}
}

func TestCalendarGridIntegrationCalculateTaskHeight(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
		RowHeight:     12.0,
		VisualConstraints: &VisualConstraints{
			MinTaskHeight: 6.0,
			MaxTaskHeight: 24.0,
		},
	}
	
	integration := NewCalendarGridIntegration(config)
	
	// Test normal visual weight
	height := integration.calculateTaskHeight(&data.Task{}, 0.5)
	if height < config.VisualConstraints.MinTaskHeight {
		t.Error("Expected height to be >= min task height")
	}
	
	if height > config.VisualConstraints.MaxTaskHeight {
		t.Error("Expected height to be <= max task height")
	}
	
	// Test high visual weight
	highWeightHeight := integration.calculateTaskHeight(&data.Task{}, 1.0)
	if highWeightHeight <= height {
		t.Error("Expected higher visual weight to result in greater height")
	}
}

func TestHandleMonthBoundaries(t *testing.T) {
	config := &GridConfig{
		CalendarStart:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:       time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:          20.0,
		DayHeight:         15.0,
		MonthBoundaryGap:  2.0,
	}
	
	integration := NewCalendarGridIntegration(config)
	
	// Create test bars
	bars := []*IntegratedTaskBar{
		{
			TaskID:        "task1",
			StartDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			StartX:        80.0,
			Width:         140.0,
			MonthBoundary: false,
		},
		{
			TaskID:        "task2",
			StartDate:     time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
			StartX:        480.0,
			Width:         140.0,
			MonthBoundary: true,
		},
	}
	
	// Handle month boundaries
	processedBars := integration.handleMonthBoundaries(bars)
	
	// Verify non-boundary task is unchanged
	if processedBars[0].StartX != 80.0 {
		t.Error("Expected non-boundary task start X to be unchanged")
	}
	
	if processedBars[0].Width != 140.0 {
		t.Error("Expected non-boundary task width to be unchanged")
	}
	
	// Verify boundary task is adjusted
	expectedStartX := 480.0 + config.MonthBoundaryGap
	if processedBars[1].StartX != expectedStartX {
		t.Errorf("Expected boundary task start X to be %.2f, got %.2f", expectedStartX, processedBars[1].StartX)
	}
	
	expectedWidth := 140.0 - config.MonthBoundaryGap
	if processedBars[1].Width != expectedWidth {
		t.Errorf("Expected boundary task width to be %.2f, got %.2f", expectedWidth, processedBars[1].Width)
	}
}

func TestGenerateIntegratedLaTeX(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	integration := NewCalendarGridIntegration(config)
	
	// Create test result
	result := &IntegratedLayoutResult{
		TaskBars: []*IntegratedTaskBar{
			{
				TaskID:      "task1",
				StartX:      80.0,
				Y:          10.0,
				Width:       140.0,
				Height:      12.0,
				Color:       "#FF0000",
				Opacity:     0.9,
				TaskName:    "Test Task",
			},
		},
		Statistics: &IntegratedLayoutStatistics{
			TotalTasks: 1,
		},
	}
	
	// Generate LaTeX
	latex := integration.GenerateIntegratedLaTeX(result)
	
	// Verify LaTeX contains expected elements
	if !containsIntegration(latex, "\\begin{integrated-calendar}") {
		t.Error("Expected LaTeX to contain integrated calendar begin tag")
	}
	
	if !containsIntegration(latex, "\\end{integrated-calendar}") {
		t.Error("Expected LaTeX to contain integrated calendar end tag")
	}
	
	if !containsIntegration(latex, "\\node[") {
		t.Error("Expected LaTeX to contain TikZ node")
	}
	
	if !containsIntegration(latex, "Test Task") {
		t.Error("Expected LaTeX to contain task name")
	}
}

func TestIntegratedLayoutStatistics(t *testing.T) {
	stats := &IntegratedLayoutStatistics{
		TotalTasks:           10,
		ProcessedBars:        10,
		TotalStacks:          3,
		ConflictsResolved:    2,
		OverflowResolutions:  1,
		VisualOptimizations:  1,
		LayoutAdjustments:    1,
		CollisionCount:       1,
		OverflowCount:        1,
		MonthBoundaryCount:   2,
		SpaceEfficiency:      0.85,
		VisualQuality:        0.90,
		AverageStackHeight:   15.0,
		MaxStackHeight:       30.0,
		AverageTaskHeight:    12.0,
		AverageTaskWidth:     50.0,
	}
	
	// Test string representation
	str := stats.String()
	if str == "" {
		t.Error("Expected string representation to be non-empty")
	}
	
	if !containsIntegration(str, "Integrated Layout Statistics") {
		t.Error("Expected string to contain header")
	}
	
	if !containsIntegration(str, "Total Tasks: 10") {
		t.Error("Expected string to contain total tasks")
	}
	
	if !containsIntegration(str, "Space Efficiency: 0.85") {
		t.Error("Expected string to contain space efficiency")
	}
}

// Helper function to check if a string contains a substring
func containsIntegration(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstringIntegration(s, substr))))
}

func containsSubstringIntegration(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
