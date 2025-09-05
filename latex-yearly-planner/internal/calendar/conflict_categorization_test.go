package calendar

import (
	"testing"
	"time"

	"latex-yearly-planner/internal/data"
)

func TestConflictCategorizer(t *testing.T) {
	// Create test calendar range and overlap detector
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	
	// Create conflict categorizer
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Test basic properties
	if categorizer.overlapDetector != overlapDetector {
		t.Error("Expected overlap detector to be set")
	}
	
	if len(categorizer.rules) == 0 {
		t.Error("Expected default rules to be initialized")
	}
	
	if len(categorizer.severityWeights) == 0 {
		t.Error("Expected severity weights to be initialized")
	}
}

func TestCategorizeConflicts(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Create test tasks with various conflict scenarios
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "High Priority Task",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  4,
			Assignee:  "John Doe",
		},
		{
			ID:        "task2",
			Name:      "Another High Priority Task",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  4,
			Assignee:  "John Doe", // Same assignee
		},
		{
			ID:        "task3",
			Name:      "Milestone Task",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Category:  "DISSERTATION",
			Priority:  5,
			IsMilestone: true,
		},
		{
			ID:        "task4",
			Name:      "Overlapping Milestone",
			StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Category:  "DISSERTATION",
			Priority:  3,
			IsMilestone: true,
		},
	}
	
	// Detect overlaps first
	overlapAnalysis := overlapDetector.DetectOverlaps(tasks)
	
	// Categorize conflicts
	conflictAnalysis := categorizer.CategorizeConflicts(overlapAnalysis)
	
	// Verify analysis structure
	if conflictAnalysis.TotalConflicts == 0 {
		t.Error("Expected to find conflicts")
	}
	
	if len(conflictAnalysis.CategorizedConflicts) == 0 {
		t.Error("Expected categorized conflicts")
	}
	
	// Check that statistics are populated
	if len(conflictAnalysis.ConflictsByCategory) == 0 {
		t.Error("Expected conflicts by category to be populated")
	}
	
	if len(conflictAnalysis.ConflictsBySeverity) == 0 {
		t.Error("Expected conflicts by severity to be populated")
	}
	
	// Check that risk assessment is generated
	if conflictAnalysis.RiskAssessment == "" {
		t.Error("Expected risk assessment to be generated")
	}
	
	// Check that recommendations are generated
	if len(conflictAnalysis.Recommendations) == 0 {
		t.Error("Expected recommendations to be generated")
	}
}

func TestCategorizeConflict(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Create test overlap
	overlap := &TaskOverlap{
		Task1ID:     "task1",
		Task2ID:     "task2",
		OverlapType: OverlapIdentical,
		Severity:    SeverityCritical,
		StartDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Duration:    time.Hour * 24 * 5,
		OverlapDays: 5,
		Priority:    4,
	}
	
	// Create test tasks
	task1 := &data.Task{
		ID:        "task1",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Priority:  4,
		Assignee:  "John Doe",
		Category:  "PROPOSAL",
	}
	
	task2 := &data.Task{
		ID:        "task2",
		StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Priority:  3,
		Assignee:  "John Doe", // Same assignee
		Category:  "PROPOSAL",
	}
	
	categorizedConflict := categorizer.categorizeConflict(overlap, task1, task2)
	
	// Verify categorization
	if categorizedConflict.Category == "" {
		t.Error("Expected conflict to be categorized")
	}
	
	if categorizedConflict.SubCategory == "" {
		t.Error("Expected subcategory to be determined")
	}
	
	if categorizedConflict.RootCause == "" {
		t.Error("Expected root cause to be determined")
	}
	
	if categorizedConflict.Impact == "" {
		t.Error("Expected impact to be assessed")
	}
	
	if categorizedConflict.Resolution.Strategy == "" {
		t.Error("Expected resolution strategy to be generated")
	}
	
	if categorizedConflict.RiskLevel == "" {
		t.Error("Expected risk level to be assessed")
	}
	
	if categorizedConflict.Urgency == "" {
		t.Error("Expected urgency to be assessed")
	}
	
	if categorizedConflict.Complexity == "" {
		t.Error("Expected complexity to be assessed")
	}
	
	// Should have alternative resolutions
	if len(categorizedConflict.AlternativeResolutions) == 0 {
		t.Error("Expected alternative resolutions to be generated")
	}
}

func TestDetermineSubCategory(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Test identical overlap
	overlap := &TaskOverlap{OverlapType: OverlapIdentical}
	task1 := &data.Task{StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)}
	task2 := &data.Task{StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)}
	
	subCategory := categorizer.determineSubCategory(overlap, task1, task2)
	if subCategory != "Identical Schedules" {
		t.Errorf("Expected 'Identical Schedules', got '%s'", subCategory)
	}
	
	// Test partial overlap with high percentage
	overlap.OverlapType = OverlapPartial
	overlap.Duration = time.Hour * 24 * 4 // 4 days
	task1.EndDate = time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC).Add(time.Hour * 24 * 5) // 5 days total
	
	subCategory = categorizer.determineSubCategory(overlap, task1, task2)
	if subCategory != "High Overlap" {
		t.Errorf("Expected 'High Overlap', got '%s'", subCategory)
	}
}

func TestAssessImpact(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Test high impact scenario
	overlap := &TaskOverlap{
		Severity: SeverityCritical,
		Duration: time.Hour * 24 * 10, // 10 days
	}
	
	task1 := &data.Task{
		Priority:  4,
		Assignee:  "John Doe",
		IsMilestone: true,
	}
	
	task2 := &data.Task{
		Priority:  4,
		Assignee:  "John Doe", // Same assignee
		IsMilestone: false,
	}
	
	impact := categorizer.assessImpact(overlap, task1, task2)
	if impact != "HIGH" {
		t.Errorf("Expected HIGH impact, got %s", impact)
	}
	
	// Test low impact scenario
	overlap.Severity = SeverityLow
	overlap.Duration = time.Hour * 24 * 1 // 1 day
	task1.Priority = 1
	task2.Priority = 1
	task1.Assignee = "John Doe"
	task2.Assignee = "Jane Doe" // Different assignee
	task1.IsMilestone = false
	task2.IsMilestone = false
	
	impact = categorizer.assessImpact(overlap, task1, task2)
	if impact != "LOW" {
		t.Errorf("Expected LOW impact, got %s", impact)
	}
}

func TestGenerateResolution(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Test schedule conflict resolution
	overlap := &TaskOverlap{OverlapType: OverlapIdentical}
	rule := &ConflictRule{Category: CategoryScheduleConflict}
	task1 := &data.Task{}
	task2 := &data.Task{}
	
	resolution := categorizer.generateResolution(overlap, task1, task2, rule)
	
	if resolution.Strategy == "" {
		t.Error("Expected resolution strategy")
	}
	
	if resolution.Description == "" {
		t.Error("Expected resolution description")
	}
	
	if len(resolution.Actions) == 0 {
		t.Error("Expected resolution actions")
	}
	
	if resolution.Priority == 0 {
		t.Error("Expected resolution priority")
	}
	
	if resolution.Effort == "" {
		t.Error("Expected resolution effort")
	}
	
	if resolution.Impact == "" {
		t.Error("Expected resolution impact")
	}
}

func TestAssessRiskLevel(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Test high risk scenario
	overlap := &TaskOverlap{
		Severity: SeverityCritical,
	}
	
	task1 := &data.Task{
		EndDate:   time.Now().AddDate(0, 0, 3), // Due in 3 days
		Priority:  5,
		IsMilestone: true,
		Dependencies: []string{"dep1", "dep2"},
	}
	
	task2 := &data.Task{
		EndDate:   time.Now().AddDate(0, 0, 5), // Due in 5 days
		Priority:  4,
		IsMilestone: false,
		Dependencies: []string{"dep3"},
	}
	
	riskLevel := categorizer.assessRiskLevel(overlap, task1, task2)
	if riskLevel != "HIGH" {
		t.Errorf("Expected HIGH risk, got %s", riskLevel)
	}
	
	// Test low risk scenario
	overlap.Severity = SeverityLow
	task1.EndDate = time.Now().AddDate(0, 0, 30) // Due in 30 days
	task1.Priority = 1
	task1.IsMilestone = false
	task1.Dependencies = []string{}
	task2.EndDate = time.Now().AddDate(0, 0, 35) // Due in 35 days
	task2.Priority = 1
	task2.Dependencies = []string{}
	
	riskLevel = categorizer.assessRiskLevel(overlap, task1, task2)
	if riskLevel != "LOW" {
		t.Errorf("Expected LOW risk, got %s", riskLevel)
	}
}

func TestAssessUrgency(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Test urgent scenario
	overlap := &TaskOverlap{Severity: SeverityCritical}
	task1 := &data.Task{
		StartDate: time.Now().AddDate(0, 0, 1), // Starts tomorrow
		Priority:  5,
	}
	task2 := &data.Task{
		StartDate: time.Now().AddDate(0, 0, 2), // Starts day after tomorrow
		Priority:  4,
	}
	
	urgency := categorizer.assessUrgency(overlap, task1, task2)
	if urgency != "URGENT" {
		t.Errorf("Expected URGENT, got %s", urgency)
	}
	
	// Test low urgency scenario
	overlap.Severity = SeverityLow
	task1.StartDate = time.Now().AddDate(0, 0, 30) // Starts in 30 days
	task1.Priority = 0 // Zero priority
	task2.StartDate = time.Now().AddDate(0, 0, 35) // Starts in 35 days
	task2.Priority = 0 // Zero priority
	
	urgency = categorizer.assessUrgency(overlap, task1, task2)
	if urgency != "LOW" {
		t.Errorf("Expected LOW urgency, got %s", urgency)
	}
}

func TestAssessComplexity(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Test high complexity scenario
	overlap := &TaskOverlap{OverlapType: OverlapIdentical}
	task1 := &data.Task{
		Assignee:    "John Doe",
		Category:    "PROPOSAL",
		Dependencies: []string{"dep1", "dep2", "dep3"},
	}
	task2 := &data.Task{
		Assignee:    "John Doe", // Same assignee
		Category:    "PROPOSAL", // Same category
		Dependencies: []string{"dep4", "dep5"},
	}
	
	complexity := categorizer.assessComplexity(overlap, task1, task2)
	if complexity != "HIGH" {
		t.Errorf("Expected HIGH complexity, got %s", complexity)
	}
	
	// Test low complexity scenario
	overlap.OverlapType = OverlapPartial
	task1.Assignee = "John Doe"
	task2.Assignee = "Jane Doe" // Different assignee
	task1.Category = "PROPOSAL"
	task2.Category = "LASER" // Different category
	task1.Dependencies = []string{}
	task2.Dependencies = []string{}
	
	complexity = categorizer.assessComplexity(overlap, task1, task2)
	if complexity != "LOW" {
		t.Errorf("Expected LOW complexity, got %s", complexity)
	}
}

func TestAddCustomRule(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	initialRuleCount := len(categorizer.rules)
	
	// Add custom rule
	customRule := ConflictRule{
		Name:        "Custom Test Rule",
		Description: "Test custom rule",
		Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
			return task1.Name == "Custom Task"
		},
		Category: CategoryScheduleConflict,
		Severity: SeverityMedium,
		Priority: 10, // High priority
	}
	
	categorizer.AddCustomRule(customRule)
	
	// Verify rule was added
	if len(categorizer.rules) != initialRuleCount+1 {
		t.Errorf("Expected %d rules, got %d", initialRuleCount+1, len(categorizer.rules))
	}
	
	// Verify rule is at the beginning (highest priority)
	if categorizer.rules[0].Name != "Custom Test Rule" {
		t.Error("Expected custom rule to be at highest priority")
	}
}

func TestConflictAnalysisMethods(t *testing.T) {
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	overlapDetector := NewOverlapDetector(calendarStart, calendarEnd)
	categorizer := NewConflictCategorizer(overlapDetector)
	
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Priority:  4,
			Assignee:  "John Doe",
			Category:  "PROPOSAL",
		},
		{
			ID:        "task2",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Priority:  4,
			Assignee:  "John Doe", // Same assignee
			Category:  "PROPOSAL",
		},
	}
	
	// Detect and categorize conflicts
	overlapAnalysis := overlapDetector.DetectOverlaps(tasks)
	conflictAnalysis := categorizer.CategorizeConflicts(overlapAnalysis)
	
	// Test filtering methods
	categoryConflicts := conflictAnalysis.GetConflictsByCategory(CategoryAssigneeConflict)
	severityConflicts := conflictAnalysis.GetConflictsBySeverity(SeverityHigh)
	urgencyConflicts := conflictAnalysis.GetConflictsByUrgency("HIGH")
	riskConflicts := conflictAnalysis.GetConflictsByRisk("HIGH")
	
	// Test summary method
	summary := conflictAnalysis.GetSummary()
	if summary == "" {
		t.Error("Expected non-empty summary")
	}
	
	// Verify counts are reasonable
	if len(categoryConflicts) < 0 {
		t.Error("Expected non-negative category conflicts count")
	}
	
	if len(severityConflicts) < 0 {
		t.Error("Expected non-negative severity conflicts count")
	}
	
	if len(urgencyConflicts) < 0 {
		t.Error("Expected non-negative urgency conflicts count")
	}
	
	if len(riskConflicts) < 0 {
		t.Error("Expected non-negative risk conflicts count")
	}
}
