package calendar

import (
	"testing"
	"time"
)

func TestMonthBoundaryEngine(t *testing.T) {
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
	
	// Create month boundary engine
	engine := NewMonthBoundaryEngine(config)
	
	// Verify initialization
	if engine.gridConfig == nil {
		t.Error("Expected grid config to be set")
	}
	
	if engine.visualConstraints == nil {
		t.Error("Expected visual constraints to be set")
	}
	
	if len(engine.boundaryRules) == 0 {
		t.Error("Expected default boundary rules to be added")
	}
	
	if len(engine.transitionRules) == 0 {
		t.Error("Expected default transition rules to be added")
	}
	
	if len(engine.continuityRules) == 0 {
		t.Error("Expected default continuity rules to be added")
	}
}

func TestProcessMonthBoundaries(t *testing.T) {
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
	
	// Create month boundary engine
	engine := NewMonthBoundaryEngine(config)
	
	// Create test task bars
	taskBars := []*IntegratedTaskBar{
		{
			TaskID:        "task1",
			StartDate:     time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
			StartX:        480.0,
			EndX:          600.0,
			Y:             10.0,
			Width:         120.0,
			Height:        12.0,
			Priority:      5,
			MonthBoundary: true,
			Category:      "PROPOSAL",
			TaskName:      "High Priority Task",
		},
		{
			TaskID:        "task2",
			StartDate:     time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
			StartX:        580.0,
			EndX:          600.0,
			Y:             30.0,
			Width:         20.0,
			Height:        12.0,
			Priority:      3,
			MonthBoundary: true,
			Category:      "MILESTONE",
			TaskName:      "MILESTONE: Complete milestone",
		},
		{
			TaskID:        "task3",
			StartDate:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			StartX:        280.0,
			EndX:          400.0,
			Y:             50.0,
			Width:         120.0,
			Height:        12.0,
			Priority:      2,
			MonthBoundary: false,
			Category:      "LASER",
			TaskName:      "Regular Task",
		},
	}
	
	// Process month boundaries
	result, err := engine.ProcessMonthBoundaries(taskBars, time.January, 2024)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Verify result
	if result == nil {
		t.Fatal("Expected result to be non-nil")
	}
	
	if len(result.ProcessedBars) != len(taskBars) {
		t.Errorf("Expected %d processed bars, got %d", len(taskBars), len(result.ProcessedBars))
	}
	
	if result.BoundaryMetrics == nil {
		t.Error("Expected boundary metrics to be non-nil")
	}
	
	if result.BoundaryMetrics.TotalTasks != len(taskBars) {
		t.Errorf("Expected %d total tasks, got %d", len(taskBars), result.BoundaryMetrics.TotalTasks)
	}
	
	if result.AnalysisDate.IsZero() {
		t.Error("Expected analysis date to be set")
	}
}

func TestGetNextMonth(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	// Test normal month
	nextMonth := engine.getNextMonth(time.January)
	if nextMonth != time.February {
		t.Errorf("Expected February, got %v", nextMonth)
	}
	
	// Test December to January
	nextMonth = engine.getNextMonth(time.December)
	if nextMonth != time.January {
		t.Errorf("Expected January, got %v", nextMonth)
	}
}

func TestGetNextYear(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	// Test normal year
	nextYear := engine.getNextYear(time.January, 2024)
	if nextYear != 2024 {
		t.Errorf("Expected 2024, got %d", nextYear)
	}
	
	// Test December to next year
	nextYear = engine.getNextYear(time.December, 2024)
	if nextYear != 2025 {
		t.Errorf("Expected 2025, got %d", nextYear)
	}
}

func TestGetMonthEndDate(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	// Test January 2024 (31 days)
	endDate := engine.getMonthEndDate(time.January, 2024)
	expectedDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	if !endDate.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, endDate)
	}
	
	// Test February 2024 (29 days - leap year)
	endDate = engine.getMonthEndDate(time.February, 2024)
	expectedDate = time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	if !endDate.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, endDate)
	}
	
	// Test April 2024 (30 days)
	endDate = engine.getMonthEndDate(time.April, 2024)
	expectedDate = time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC)
	if !endDate.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, endDate)
	}
}

func TestCalculateTaskDensity(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	// Test with empty bars
	density := engine.calculateTaskDensity([]*IntegratedTaskBar{})
	if density != 0.0 {
		t.Errorf("Expected density 0.0 for empty bars, got %.2f", density)
	}
	
	// Test with task bars
	bars := []*IntegratedTaskBar{
		{
			StartX: 10.0,
			EndX:   30.0,
			Y:      5.0,
			Width:  20.0,
			Height: 10.0,
		},
		{
			StartX: 40.0,
			EndX:   60.0,
			Y:      5.0,
			Width:  20.0,
			Height: 10.0,
		},
	}
	
	density = engine.calculateTaskDensity(bars)
	if density <= 0 {
		t.Error("Expected density to be positive")
	}
	
	if density > 1.0 {
		t.Error("Expected density to be <= 1.0")
	}
}

func TestMonthBoundaryEngineBarsOverlap(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	// Create overlapping bars
	bar1 := &IntegratedTaskBar{
		StartX: 10.0,
		EndX:   30.0,
		Y:      5.0,
		Height: 10.0,
	}
	
	bar2 := &IntegratedTaskBar{
		StartX: 25.0,
		EndX:   45.0,
		Y:      8.0,
		Height: 10.0,
	}
	
	// Test overlap detection
	if !engine.barsOverlap(bar1, bar2) {
		t.Error("Expected bars to overlap")
	}
	
	// Create non-overlapping bars
	bar3 := &IntegratedTaskBar{
		StartX: 10.0,
		EndX:   30.0,
		Y:      5.0,
		Height: 10.0,
	}
	
	bar4 := &IntegratedTaskBar{
		StartX: 50.0,
		EndX:   70.0,
		Y:      5.0,
		Height: 10.0,
	}
	
	// Test non-overlap
	if engine.barsOverlap(bar3, bar4) {
		t.Error("Expected bars not to overlap")
	}
}

func TestCountOverlaps(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	// Test with no overlaps
	bars := []*IntegratedTaskBar{
		{
			StartX: 10.0,
			EndX:   30.0,
			Y:      5.0,
			Height: 10.0,
		},
		{
			StartX: 50.0,
			EndX:   70.0,
			Y:      5.0,
			Height: 10.0,
		},
	}
	
	overlapCount := engine.countOverlaps(bars)
	if overlapCount != 0 {
		t.Errorf("Expected 0 overlaps, got %d", overlapCount)
	}
	
	// Test with overlaps
	bars = []*IntegratedTaskBar{
		{
			StartX: 10.0,
			EndX:   30.0,
			Y:      5.0,
			Height: 10.0,
		},
		{
			StartX: 25.0,
			EndX:   45.0,
			Y:      8.0,
			Height: 10.0,
		},
		{
			StartX: 50.0,
			EndX:   70.0,
			Y:      5.0,
			Height: 10.0,
		},
	}
	
	overlapCount = engine.countOverlaps(bars)
	if overlapCount != 1 {
		t.Errorf("Expected 1 overlap, got %d", overlapCount)
	}
}

func TestCalculateContinuityScore(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	context := &MonthBoundaryContext{
		CurrentMonth: time.January,
		NextMonth:    time.February,
		CurrentYear:  2024,
		NextYear:     2024,
		DayWidth:     20.0,
		DayHeight:    15.0,
	}
	
	// Test with no month boundary tasks
	bars := []*IntegratedTaskBar{
		{
			TaskID:        "task1",
			MonthBoundary: false,
		},
		{
			TaskID:        "task2",
			MonthBoundary: false,
		},
	}
	
	continuations := []*TaskContinuation{}
	score := engine.calculateContinuityScore(bars, continuations, context)
	if score != 0.0 {
		t.Errorf("Expected continuity score 0.0, got %.2f", score)
	}
	
	// Test with month boundary tasks
	bars = []*IntegratedTaskBar{
		{
			TaskID:        "task1",
			MonthBoundary: true,
		},
		{
			TaskID:        "task2",
			MonthBoundary: false,
		},
	}
	
	score = engine.calculateContinuityScore(bars, continuations, context)
	if score != 0.5 {
		t.Errorf("Expected continuity score 0.5, got %.2f", score)
	}
}

func TestCalculateVisualConsistency(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	context := &MonthBoundaryContext{
		CurrentMonth: time.January,
		NextMonth:    time.February,
		CurrentYear:  2024,
		NextYear:     2024,
		DayWidth:     20.0,
		DayHeight:    15.0,
	}
	
	// Test with consistent visual properties
	bars := []*IntegratedTaskBar{
		{
			TaskID:       "task1",
			Color:        "#FF0000",
			BorderColor:  "#000000",
		},
		{
			TaskID:       "task2",
			Color:        "#00FF00",
			BorderColor:  "#000000",
		},
	}
	
	score := engine.calculateVisualConsistency(bars, context)
	if score != 1.0 {
		t.Errorf("Expected visual consistency score 1.0, got %.2f", score)
	}
	
	// Test with inconsistent visual properties
	bars = []*IntegratedTaskBar{
		{
			TaskID:       "task1",
			Color:        "",
			BorderColor:  "",
		},
		{
			TaskID:       "task2",
			Color:        "#00FF00",
			BorderColor:  "#000000",
		},
	}
	
	score = engine.calculateVisualConsistency(bars, context)
	if score != 0.5 {
		t.Errorf("Expected visual consistency score 0.5, got %.2f", score)
	}
}

func TestCalculateTransitionSmoothness(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	context := &MonthBoundaryContext{
		CurrentMonth: time.January,
		NextMonth:    time.February,
		CurrentYear:  2024,
		NextYear:     2024,
		DayWidth:     20.0,
		DayHeight:    15.0,
	}
	
	// Test with smooth transitions
	transitions := []*TaskTransition{
		{
			TaskID:         "task1",
			TransitionType: TransitionSmooth,
		},
		{
			TaskID:         "task2",
			TransitionType: TransitionFade,
		},
	}
	
	score := engine.calculateTransitionSmoothness(transitions, context)
	if score != 0.9 {
		t.Errorf("Expected transition smoothness score 0.9, got %.2f", score)
	}
	
	// Test with no transitions
	score = engine.calculateTransitionSmoothness([]*TaskTransition{}, context)
	if score != 1.0 {
		t.Errorf("Expected transition smoothness score 1.0 for no transitions, got %.2f", score)
	}
}

func TestMonthBoundaryEngineCalculateSpaceEfficiency(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	context := &MonthBoundaryContext{
		CurrentMonth: time.January,
		NextMonth:    time.February,
		CurrentYear:  2024,
		NextYear:     2024,
		DayWidth:     20.0,
		DayHeight:    15.0,
	}
	
	// Test with empty bars
	score := engine.calculateSpaceEfficiency([]*IntegratedTaskBar{}, context)
	if score != 0.0 {
		t.Errorf("Expected space efficiency score 0.0 for empty bars, got %.2f", score)
	}
	
	// Test with task bars
	bars := []*IntegratedTaskBar{
		{
			StartX: 10.0,
			EndX:   30.0,
			Y:      5.0,
			Width:  20.0,
			Height: 10.0,
		},
		{
			StartX: 40.0,
			EndX:   60.0,
			Y:      5.0,
			Width:  20.0,
			Height: 10.0,
		},
	}
	
	score = engine.calculateSpaceEfficiency(bars, context)
	if score <= 0 {
		t.Error("Expected space efficiency score to be positive")
	}
	
	if score > 1.0 {
		t.Error("Expected space efficiency score to be <= 1.0")
	}
}

func TestGenerateBoundaryRecommendations(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewMonthBoundaryEngine(config)
	
	context := &MonthBoundaryContext{
		CurrentMonth: time.January,
		NextMonth:    time.February,
		CurrentYear:  2024,
		NextYear:     2024,
		DayWidth:     20.0,
		DayHeight:    15.0,
	}
	
	// Test with poor metrics
	metrics := &BoundaryMetrics{
		ContinuityScore:      0.5,
		VisualConsistency:    0.6,
		TransitionSmoothness: 0.4,
		GridContinuity:       0.7,
		SpaceEfficiency:      0.5,
		VisualBalance:        0.4,
	}
	
	recommendations := engine.generateBoundaryRecommendations(metrics, context)
	
	if len(recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
	
	// Check for specific recommendations
	foundContinuity := false
	foundVisual := false
	foundTransition := false
	foundGrid := false
	foundSpace := false
	foundBalance := false
	
	for _, rec := range recommendations {
		if containsMonthBoundary(rec, "continuity") {
			foundContinuity = true
		}
		if containsMonthBoundary(rec, "visual") {
			foundVisual = true
		}
		if containsMonthBoundary(rec, "transition") {
			foundTransition = true
		}
		if containsMonthBoundary(rec, "grid") {
			foundGrid = true
		}
		if containsMonthBoundary(rec, "space") {
			foundSpace = true
		}
		if containsMonthBoundary(rec, "balance") {
			foundBalance = true
		}
	}
	
	if !foundContinuity {
		t.Error("Expected continuity recommendation")
	}
	
	if !foundVisual {
		t.Error("Expected visual recommendation")
	}
	
	if !foundTransition {
		t.Error("Expected transition recommendation")
	}
	
	if !foundGrid {
		t.Error("Expected grid recommendation")
	}
	
	if !foundSpace {
		t.Error("Expected space recommendation")
	}
	
	if !foundBalance {
		t.Error("Expected balance recommendation")
	}
}

// Helper function to check if a string contains a substring
func containsMonthBoundary(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstringMonthBoundary(s, substr))))
}

func containsSubstringMonthBoundary(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
