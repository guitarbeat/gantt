package calendar

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/data"
)

func TestPositioningEngine(t *testing.T) {
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
	
	// Create positioning engine
	engine := NewPositioningEngine(config)
	
	// Verify initialization
	if engine.gridConfig == nil {
		t.Error("Expected grid config to be set")
	}
	
	if engine.visualConstraints == nil {
		t.Error("Expected visual constraints to be set")
	}
	
	if len(engine.alignmentRules) == 0 {
		t.Error("Expected default alignment rules to be added")
	}
	
	if len(engine.spacingRules) == 0 {
		t.Error("Expected default spacing rules to be added")
	}
	
	if engine.layoutMetrics == nil {
		t.Error("Expected layout metrics to be initialized")
	}
}

func TestPositionTasks(t *testing.T) {
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
	
	// Create positioning engine
	engine := NewPositioningEngine(config)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "High Priority Task",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  5,
		},
		{
			ID:        "task2",
			Name:      "MILESTONE: Complete milestone",
			StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Priority:  3,
		},
		{
			ID:        "task3",
			Name:      "Regular Task",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			Category:  "LASER",
			Priority:  2,
		},
	}
	
	// Position tasks
	result, err := engine.PositionTasks(tasks, []*IntegratedTaskBar{})
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
	
	if result.Metrics == nil {
		t.Error("Expected metrics to be non-nil")
	}
	
	if result.Metrics.TotalTasks != len(tasks) {
		t.Errorf("Expected %d total tasks, got %d", len(tasks), result.Metrics.TotalTasks)
	}
	
	if result.AnalysisDate.IsZero() {
		t.Error("Expected analysis date to be set")
	}
}

func TestPositioningEngineCalculateXPosition(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewPositioningEngine(config)
	
	context := &PositioningContext{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	// Test start date
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	x := engine.calculateXPosition(startDate, context)
	if x != 0 {
		t.Errorf("Expected X position 0 for start date, got %.2f", x)
	}
	
	// Test 5 days later
	date5Days := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)
	x = engine.calculateXPosition(date5Days, context)
	expectedX := 5.0 * 20.0
	if x != expectedX {
		t.Errorf("Expected X position %.2f for 5 days later, got %.2f", expectedX, x)
	}
}

func TestPositioningEngineCalculateTaskHeight(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
		VisualConstraints: &VisualConstraints{
			MinTaskHeight: 6.0,
			MaxTaskHeight: 24.0,
		},
	}
	
	engine := NewPositioningEngine(config)
	
	context := &PositioningContext{
		DayHeight: 15.0,
		GridConstraints: &GridConstraints{
			MinRowHeight: 6.0,
			MaxRowHeight: 24.0,
		},
	}
	
	// Test normal duration task
	task := &data.Task{
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
	}
	
	height := engine.calculateTaskHeight(task, context)
	if height <= 0 {
		t.Error("Expected height to be positive")
	}
	
	if height < context.GridConstraints.MinRowHeight {
		t.Error("Expected height to be >= min row height")
	}
	
	if height > context.GridConstraints.MaxRowHeight {
		t.Error("Expected height to be <= max row height")
	}
	
	// Test short duration task
	shortTask := &data.Task{
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	shortHeight := engine.calculateTaskHeight(shortTask, context)
	if shortHeight >= height {
		t.Error("Expected short duration task to have smaller height")
	}
	
	// Test long duration task
	longTask := &data.Task{
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	longHeight := engine.calculateTaskHeight(longTask, context)
	if longHeight <= height {
		t.Error("Expected long duration task to have larger height")
	}
}

func TestBarsCollide(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewPositioningEngine(config)
	
	context := &PositioningContext{
		GridConstraints: &GridConstraints{
			CollisionBuffer: 2.0,
		},
	}
	
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
	
	// Test collision detection
	if !engine.barsCollide(bar1, bar2, context) {
		t.Error("Expected bars to collide")
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
	
	// Test non-collision
	if engine.barsCollide(bar3, bar4, context) {
		t.Error("Expected bars not to collide")
	}
}

func TestCalculateDistance(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewPositioningEngine(config)
	
	bar1 := &IntegratedTaskBar{
		StartX: 0.0,
		EndX:   10.0,
		Y:      0.0,
		Height: 10.0,
	}
	
	bar2 := &IntegratedTaskBar{
		StartX: 10.0,
		EndX:   20.0,
		Y:      0.0,
		Height: 10.0,
	}
	
	distance := engine.calculateDistance(bar1, bar2)
	if distance <= 0 {
		t.Error("Expected distance to be positive")
	}
	
	// Test with same position (should be 0)
	bar3 := &IntegratedTaskBar{
		StartX: 0.0,
		EndX:   10.0,
		Y:      0.0,
		Height: 10.0,
	}
	
	bar4 := &IntegratedTaskBar{
		StartX: 0.0,
		EndX:   10.0,
		Y:      0.0,
		Height: 10.0,
	}
	
	sameDistance := engine.calculateDistance(bar3, bar4)
	if sameDistance != 0 {
		t.Errorf("Expected distance 0 for same position, got %.2f", sameDistance)
	}
}

func TestSnapToGrid(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewPositioningEngine(config)
	
	context := &PositioningContext{
		GridConstraints: &GridConstraints{
			SnapToGrid:     true,
			GridResolution: 2.0,
		},
	}
	
	// Create bars with non-grid positions
	bars := []*IntegratedTaskBar{
		{
			StartX: 10.7,
			EndX:   20.3,
			Y:      5.8,
			Height: 8.1,
		},
		{
			StartX: 25.1,
			EndX:   35.9,
			Y:      7.2,
			Height: 9.7,
		},
	}
	
	// Snap to grid
	snappedBars := engine.snapToGrid(bars, context)
	
	// Verify snapping
	for _, bar := range snappedBars {
		// Check X position is snapped to grid
		if int(bar.StartX)%int(context.GridConstraints.GridResolution) != 0 {
			t.Error("Expected X position to be snapped to grid")
		}
		
		// Check Y position is snapped to grid
		if int(bar.Y)%int(context.GridConstraints.GridResolution) != 0 {
			t.Error("Expected Y position to be snapped to grid")
		}
		
		// Check height is snapped to grid
		if int(bar.Height)%int(context.GridConstraints.GridResolution) != 0 {
			t.Error("Expected height to be snapped to grid")
		}
	}
}

func TestResolveCollisions(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewPositioningEngine(config)
	
	context := &PositioningContext{
		GridConstraints: &GridConstraints{
			CollisionBuffer: 2.0,
		},
	}
	
	// Create colliding bars
	bars := []*IntegratedTaskBar{
		{
			StartX: 10.0,
			EndX:   30.0,
			Y:      5.0,
			Height: 10.0,
			Priority: 5, // Higher priority
		},
		{
			StartX: 25.0,
			EndX:   45.0,
			Y:      8.0,
			Height: 10.0,
			Priority: 3, // Lower priority
		},
	}
	
	// Resolve collisions
	resolvedBars := engine.resolveCollisions(bars, context)
	
	// Verify no collisions remain
	for i := 0; i < len(resolvedBars)-1; i++ {
		for j := i + 1; j < len(resolvedBars); j++ {
			if engine.barsCollide(resolvedBars[i], resolvedBars[j], context) {
				t.Error("Expected no collisions after resolution")
			}
		}
	}
}

func TestCalculateLayoutMetrics(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewPositioningEngine(config)
	
	context := &PositioningContext{
		AvailableWidth:  140.0,
		AvailableHeight: 60.0,
		GridConstraints: &GridConstraints{
			SnapToGrid:     true,
			GridResolution: 1.0,
			MinTaskSpacing: 1.0,
			MaxTaskSpacing: 10.0,
		},
	}
	
	// Create test bars
	bars := []*IntegratedTaskBar{
		{
			StartX: 10.0,
			EndX:   30.0,
			Y:      5.0,
			Height: 10.0,
			Width:  20.0,
		},
		{
			StartX: 40.0,
			EndX:   60.0,
			Y:      5.0,
			Height: 10.0,
			Width:  20.0,
		},
	}
	
	// Calculate metrics
	metrics := engine.calculateLayoutMetrics(bars, context)
	
	// Verify metrics
	if metrics.TotalTasks != len(bars) {
		t.Errorf("Expected %d total tasks, got %d", len(bars), metrics.TotalTasks)
	}
	
	if metrics.PositionedTasks != len(bars) {
		t.Errorf("Expected %d positioned tasks, got %d", len(bars), metrics.PositionedTasks)
	}
	
	if metrics.SpaceEfficiency < 0 || metrics.SpaceEfficiency > 1 {
		t.Error("Expected space efficiency to be between 0 and 1")
	}
	
	if metrics.AlignmentScore < 0 || metrics.AlignmentScore > 1 {
		t.Error("Expected alignment score to be between 0 and 1")
	}
	
	if metrics.SpacingScore < 0 || metrics.SpacingScore > 1 {
		t.Error("Expected spacing score to be between 0 and 1")
	}
	
	if metrics.VisualBalance < 0 || metrics.VisualBalance > 1 {
		t.Error("Expected visual balance to be between 0 and 1")
	}
	
	if metrics.GridUtilization < 0 || metrics.GridUtilization > 1 {
		t.Error("Expected grid utilization to be between 0 and 1")
	}
}

func TestGeneratePositioningRecommendations(t *testing.T) {
	config := &GridConfig{
		CalendarStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:      20.0,
		DayHeight:     15.0,
	}
	
	engine := NewPositioningEngine(config)
	
	context := &PositioningContext{
		AvailableWidth:  140.0,
		AvailableHeight: 60.0,
		GridConstraints: &GridConstraints{
			SnapToGrid:     true,
			GridResolution: 1.0,
			MinTaskSpacing: 1.0,
			MaxTaskSpacing: 10.0,
		},
	}
	
	// Test with poor metrics
	metrics := &PositioningLayoutMetrics{
		SpaceEfficiency: 0.5,
		AlignmentScore:  0.6,
		SpacingScore:    0.4,
		VisualBalance:   0.3,
		CollisionCount:  2,
	}
	
	recommendations := engine.generatePositioningRecommendations(metrics, context)
	
	if len(recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
	
	// Check for specific recommendations
	foundSpaceEfficiency := false
	foundAlignment := false
	foundSpacing := false
	foundBalance := false
	foundCollision := false
	
	for _, rec := range recommendations {
		if containsPositioning(rec, "space efficiency") {
			foundSpaceEfficiency = true
		}
		if containsPositioning(rec, "alignment") {
			foundAlignment = true
		}
		if containsPositioning(rec, "spacing") {
			foundSpacing = true
		}
		if containsPositioning(rec, "balance") {
			foundBalance = true
		}
		if containsPositioning(rec, "collision") {
			foundCollision = true
		}
	}
	
	if !foundSpaceEfficiency {
		t.Error("Expected space efficiency recommendation")
	}
	
	if !foundAlignment {
		t.Error("Expected alignment recommendation")
	}
	
	if !foundSpacing {
		t.Error("Expected spacing recommendation")
	}
	
	if !foundBalance {
		t.Error("Expected balance recommendation")
	}
	
	if !foundCollision {
		t.Error("Expected collision recommendation")
	}
}

// Helper function to check if a string contains a substring
func containsPositioning(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstringPositioning(s, substr))))
}

func containsSubstringPositioning(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
