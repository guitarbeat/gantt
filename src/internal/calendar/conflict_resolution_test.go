package calendar

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/data"
)

func TestConflictResolutionEngine(t *testing.T) {
	// Create test calendar range
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	taskPrioritizationEngine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	
	// Create conflict resolution engine
	engine := NewConflictResolutionEngine(taskPrioritizationEngine, verticalStackingEngine, smartStackingEngine)
	
	// Test basic properties
	if engine.taskPrioritizationEngine != taskPrioritizationEngine {
		t.Error("Expected task prioritization engine to be set")
	}
	
	if engine.verticalStackingEngine != verticalStackingEngine {
		t.Error("Expected vertical stacking engine to be set")
	}
	
	if engine.smartStackingEngine != smartStackingEngine {
		t.Error("Expected smart stacking engine to be set")
	}
	
	if engine.overflowManager == nil {
		t.Error("Expected overflow manager to be initialized")
	}
	
	if engine.visualConflictResolver == nil {
		t.Error("Expected visual conflict resolver to be initialized")
	}
}

func TestResolveConflicts(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	taskPrioritizationEngine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	engine := NewConflictResolutionEngine(taskPrioritizationEngine, verticalStackingEngine, smartStackingEngine)
	
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
	context := &PriorityContext{
		CalendarStart:      calendarStart,
		CalendarEnd:        calendarEnd,
		CurrentTime:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		AssigneeWorkloads:  map[string]int{"John Doe": 2, "Jane Doe": 1},
		CategoryImportance: map[string]float64{"DISSERTATION": 10.0, "PROPOSAL": 8.0, "LASER": 5.0},
	}
	
	// Resolve conflicts
	result := engine.ResolveConflicts(tasks, context)
	
	// Verify result structure
	if result == nil {
		t.Error("Expected non-nil conflict resolution result")
	}
	
	if result.ResolvedConflicts == nil {
		t.Error("Expected resolved conflicts to be initialized")
	}
	
	if result.OverflowResolutions == nil {
		t.Error("Expected overflow resolutions to be initialized")
	}
	
	if result.VisualOptimizations == nil {
		t.Error("Expected visual optimizations to be initialized")
	}
	
	if result.LayoutAdjustments == nil {
		t.Error("Expected layout adjustments to be initialized")
	}
	
	// Verify analysis date is set
	if result.AnalysisDate.IsZero() {
		t.Error("Expected analysis date to be set")
	}
}

func TestDetectConflicts(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	taskPrioritizationEngine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	engine := NewConflictResolutionEngine(taskPrioritizationEngine, verticalStackingEngine, smartStackingEngine)
	
	// Create test tasks with overlapping scenarios
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Overlapping Task 1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:  "DISSERTATION",
			Priority:  5,
		},
		{
			ID:        "task2",
			Name:      "Overlapping Task 2",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  4,
		},
	}
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      calendarStart,
		CalendarEnd:        calendarEnd,
		CurrentTime:        time.Now(),
		AssigneeWorkloads:  make(map[string]int),
		CategoryImportance: make(map[string]float64),
	}
	
	// Detect conflicts
	conflicts := engine.detectConflicts(tasks, context)
	
	// Verify conflicts are detected
	if len(conflicts) == 0 {
		t.Error("Expected conflicts to be detected for overlapping tasks")
	}
	
	// Verify conflict structure
	for _, conflict := range conflicts {
		if conflict.Task1ID == "" || conflict.Task2ID == "" {
			t.Error("Expected both task IDs to be set in conflict")
		}
		
		if conflict.OverlapType == "" {
			t.Error("Expected overlap type to be set")
		}
		
		if conflict.Severity == "" {
			t.Error("Expected severity to be set")
		}
	}
}

func TestOverflowManager(t *testing.T) {
	om := NewOverflowManager()
	
	// Test basic properties
	if om.overflowThresholds == nil {
		t.Error("Expected overflow thresholds to be initialized")
	}
	
	if om.resolutionStrategies == nil {
		t.Error("Expected resolution strategies to be initialized")
	}
	
	if !om.adaptiveThresholds {
		t.Error("Expected adaptive thresholds to be enabled")
	}
	
	if !om.smartCompression {
		t.Error("Expected smart compression to be enabled")
	}
	
	// Test threshold values
	expectedThresholds := map[OverflowType]float64{
		OverflowVertical:   0.8,
		OverflowHorizontal: 0.8,
		OverflowArea:       0.9,
		OverflowDensity:    0.7,
		OverflowCollision:  0.1,
	}
	
	for overflowType, expectedThreshold := range expectedThresholds {
		if threshold, exists := om.overflowThresholds[overflowType]; !exists {
			t.Errorf("Expected threshold for %s to be defined", overflowType)
		} else if threshold != expectedThreshold {
			t.Errorf("Expected threshold for %s to be %f, got %f", overflowType, expectedThreshold, threshold)
		}
	}
}

func TestVisualConflictResolver(t *testing.T) {
	vcr := NewVisualConflictResolver()
	
	// Test basic properties
	if vcr.collisionDetector == nil {
		t.Error("Expected collision detector to be initialized")
	}
	
	if vcr.conflictResolvers == nil {
		t.Error("Expected conflict resolvers to be initialized")
	}
	
	if vcr.visualOptimizer == nil {
		t.Error("Expected visual optimizer to be initialized")
	}
	
	if !vcr.adaptiveResolution {
		t.Error("Expected adaptive resolution to be enabled")
	}
}

func TestCollisionDetector(t *testing.T) {
	cd := NewCollisionDetector()
	
	// Test basic properties
	if cd.collisionThreshold != 0.1 {
		t.Error("Expected collision threshold to be 0.1")
	}
	
	if cd.boundingBoxBuffer != 2.0 {
		t.Error("Expected bounding box buffer to be 2.0")
	}
	
	if cd.zIndexManager == nil {
		t.Error("Expected z-index manager to be initialized")
	}
}

func TestZIndexManager(t *testing.T) {
	zim := NewZIndexManager()
	
	// Test basic properties
	if zim.baseZIndex != 1000 {
		t.Error("Expected base z-index to be 1000")
	}
	
	if zim.layerSpacing != 10 {
		t.Error("Expected layer spacing to be 10")
	}
	
	if zim.priorityLayers == nil {
		t.Error("Expected priority layers to be initialized")
	}
	
	// Test priority layer values
	expectedLayers := map[VisualProminence]int{
		ProminenceCritical: 4,
		ProminenceHigh:     3,
		ProminenceMedium:   2,
		ProminenceLow:      1,
		ProminenceMinimal:  0,
	}
	
	for prominence, expectedLayer := range expectedLayers {
		if layer, exists := zim.priorityLayers[prominence]; !exists {
			t.Errorf("Expected layer for %s to be defined", prominence)
		} else if layer != expectedLayer {
			t.Errorf("Expected layer for %s to be %d, got %d", prominence, expectedLayer, layer)
		}
	}
}

func TestVisualOptimizer(t *testing.T) {
	vo := NewVisualOptimizer()
	
	// Test basic properties
	if vo.optimizationRules == nil {
		t.Error("Expected optimization rules to be initialized")
	}
	
	if vo.layoutStrategies == nil {
		t.Error("Expected layout strategies to be initialized")
	}
	
	if !vo.adaptiveLayout {
		t.Error("Expected adaptive layout to be enabled")
	}
}

func TestDetectOverflowIssues(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	taskPrioritizationEngine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	engine := NewConflictResolutionEngine(taskPrioritizationEngine, verticalStackingEngine, smartStackingEngine)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Long Task",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
			Category:  "DISSERTATION",
			Priority:  5,
		},
	}
	
	// Create context
	context := &PriorityContext{
		CalendarStart:      calendarStart,
		CalendarEnd:        calendarEnd,
		CurrentTime:        time.Now(),
		AssigneeWorkloads:  make(map[string]int),
		CategoryImportance: make(map[string]float64),
	}
	
	// Detect overflow issues
	issues := engine.detectOverflowIssues(tasks, context)
	
	// Verify issues are detected
	if len(issues) == 0 {
		t.Error("Expected overflow issues to be detected")
	}
	
	// Verify issue structure
	for _, issue := range issues {
		if issue.OverflowID == "" {
			t.Error("Expected overflow ID to be set")
		}
		
		if issue.OverflowType == "" {
			t.Error("Expected overflow type to be set")
		}
		
		if issue.Severity < 0 {
			t.Error("Expected severity to be non-negative")
		}
		
		if issue.SpaceConstraints == nil {
			t.Error("Expected space constraints to be set")
		}
	}
}

func TestDetermineConflictType(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create all components
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	taskPrioritizationEngine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	engine := NewConflictResolutionEngine(taskPrioritizationEngine, verticalStackingEngine, smartStackingEngine)
	
	// Test different overlap types
	testCases := []struct {
		overlapType    OverlapType
		expectedType   string
	}{
		{OverlapIdentical, "IDENTICAL"},
		{OverlapNested, "NESTED"},
		{OverlapComplete, "COMPLETE"},
		{OverlapPartial, "PARTIAL"},
		{OverlapAdjacent, "ADJACENT"},
		{"DEPENDENCY", "DEPENDENCY"},
	}
	
	for _, tc := range testCases {
		conflict := &TaskOverlap{
			OverlapType: tc.overlapType,
		}
		
		conflictType := engine.determineConflictType(conflict)
		if conflictType != tc.expectedType {
			t.Errorf("Expected conflict type %s for overlap type %s, got %s", 
				tc.expectedType, tc.overlapType, conflictType)
		}
	}
}

func TestConflictResolutionResultMethods(t *testing.T) {
	// Create test resolved conflicts
	resolvedConflicts := []*ResolvedConflict{
		{
			ConflictID:   "conflict1",
			ConflictType: "PARTIAL",
		},
		{
			ConflictID:   "conflict2",
			ConflictType: "COMPLETE",
		},
	}
	
	// Create test overflow resolutions
	overflowResolutions := []*OverflowResolution{
		{
			OverflowID:   "overflow1",
			OverflowType: OverflowVertical,
		},
		{
			OverflowID:   "overflow2",
			OverflowType: OverflowHorizontal,
		},
	}
	
	// Create test visual optimizations
	visualOptimizations := []*VisualOptimization{
		{
			OptimizationID:   "opt1",
			OptimizationType: "Visual Balance",
			Description:      "Improved visual balance",
			VisualImprovement: 0.1,
			SpaceEfficiency:   0.8,
		},
	}
	
	// Create test layout adjustments
	layoutAdjustments := []*LayoutAdjustment{
		{
			AdjustmentID:   "adj1",
			AdjustmentType: "STACK",
			Description:    "Applied stack layout",
		},
	}
	
	// Create result
	result := &ConflictResolutionResult{
		ResolvedConflicts:   resolvedConflicts,
		OverflowResolutions: overflowResolutions,
		VisualOptimizations: visualOptimizations,
		LayoutAdjustments:   layoutAdjustments,
		Recommendations:     []string{"Test recommendation"},
		AnalysisDate:        time.Now(),
	}
	
	// Test filtering methods
	partialConflicts := result.GetResolvedConflictsByType("PARTIAL")
	if len(partialConflicts) != 1 {
		t.Errorf("Expected 1 partial conflict, got %d", len(partialConflicts))
	}
	
	verticalOverflows := result.GetOverflowResolutionsByType(OverflowVertical)
	if len(verticalOverflows) != 1 {
		t.Errorf("Expected 1 vertical overflow, got %d", len(verticalOverflows))
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

func TestSpaceConstraints(t *testing.T) {
	sc := &SpaceConstraints{
		MaxWidth:     800.0,
		MaxHeight:    600.0,
		MinWidth:     100.0,
		MinHeight:    100.0,
		AvailableArea: 480000.0,
		UsedArea:     400000.0,
	}
	
	// Test basic properties
	if sc.MaxWidth != 800.0 {
		t.Error("Expected max width to be 800.0")
	}
	
	if sc.MaxHeight != 600.0 {
		t.Error("Expected max height to be 600.0")
	}
	
	if sc.AvailableArea != 480000.0 {
		t.Error("Expected available area to be 480000.0")
	}
	
	if sc.UsedArea != 400000.0 {
		t.Error("Expected used area to be 400000.0")
	}
}

func TestConflictMetrics(t *testing.T) {
	cm := &ConflictMetrics{
		CollisionCount:  2,
		OverlapCount:    3,
		SeverityScore:   0.5,
		VisualClarity:   0.8,
		SpaceEfficiency: 0.7,
	}
	
	// Test basic properties
	if cm.CollisionCount != 2 {
		t.Error("Expected collision count to be 2")
	}
	
	if cm.OverlapCount != 3 {
		t.Error("Expected overlap count to be 3")
	}
	
	if cm.SeverityScore != 0.5 {
		t.Error("Expected severity score to be 0.5")
	}
	
	if cm.VisualClarity != 0.8 {
		t.Error("Expected visual clarity to be 0.8")
	}
	
	if cm.SpaceEfficiency != 0.7 {
		t.Error("Expected space efficiency to be 0.7")
	}
}

func TestOverflowMetrics(t *testing.T) {
	om := &OverflowMetrics{
		OverflowAmount:     50.0,
		OverflowPercentage: 10.0,
		AffectedTasks:      2,
		SeverityScore:      0.3,
		SpaceWaste:         25.0,
	}
	
	// Test basic properties
	if om.OverflowAmount != 50.0 {
		t.Error("Expected overflow amount to be 50.0")
	}
	
	if om.OverflowPercentage != 10.0 {
		t.Error("Expected overflow percentage to be 10.0")
	}
	
	if om.AffectedTasks != 2 {
		t.Error("Expected affected tasks to be 2")
	}
	
	if om.SeverityScore != 0.3 {
		t.Error("Expected severity score to be 0.3")
	}
	
	if om.SpaceWaste != 25.0 {
		t.Error("Expected space waste to be 25.0")
	}
}

func TestVisualMetrics(t *testing.T) {
	vm := &VisualMetrics{
		VisualBalance:    0.8,
		SpaceEfficiency:  0.7,
		ClarityScore:     0.9,
		HarmonyScore:     0.8,
		ReadabilityScore: 0.85,
	}
	
	// Test basic properties
	if vm.VisualBalance != 0.8 {
		t.Error("Expected visual balance to be 0.8")
	}
	
	if vm.SpaceEfficiency != 0.7 {
		t.Error("Expected space efficiency to be 0.7")
	}
	
	if vm.ClarityScore != 0.9 {
		t.Error("Expected clarity score to be 0.9")
	}
	
	if vm.HarmonyScore != 0.8 {
		t.Error("Expected harmony score to be 0.8")
	}
	
	if vm.ReadabilityScore != 0.85 {
		t.Error("Expected readability score to be 0.85")
	}
}
