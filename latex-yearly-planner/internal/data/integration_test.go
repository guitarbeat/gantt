package data

import (
	"testing"
	"time"
)

// TestIntegrationWithRealData tests all data structures with real CSV data
func TestIntegrationWithRealData(t *testing.T) {
	// Test with the fixed CSV data
	reader := NewReader("/Users/aaron/Downloads/gantt/input/data.cleaned.fixed.csv")
	tasks, err := reader.ReadTasks()
	if err != nil {
		t.Fatalf("Failed to read CSV data: %v", err)
	}
	
	if len(tasks) == 0 {
		t.Fatal("Expected tasks to be loaded")
	}
	
	t.Logf("Loaded %d tasks from CSV", len(tasks))
	
	// Test TaskCollection
	collection := NewTaskCollection()
	for _, task := range tasks {
		collection.AddTask(&task)
	}
	
	// Test categorization
	categoryManager := NewTaskCategoryManager()
	categorizedTasks := make(map[string]int)
	
	for _, task := range tasks {
		category := categoryManager.CategorizeTask(&task)
		categorizedTasks[category]++
		
		// Verify category info
		categoryInfo, exists := categoryManager.GetCategory(category)
		if !exists {
			t.Errorf("Category %s should exist", category)
		}
		if categoryInfo.Name != category {
			t.Errorf("Expected category name %s, got %s", category, categoryInfo.Name)
		}
	}
	
	t.Logf("Task categorization results:")
	for category, count := range categorizedTasks {
		t.Logf("  %s: %d tasks", category, count)
	}
	
	// Test dependency graph
	dependencyGraph := NewTaskDependencyGraph()
	for _, task := range tasks {
		dependencyGraph.AddTask(&task)
	}
	
	// Test hierarchy
	hierarchy := NewTaskHierarchy()
	for _, task := range tasks {
		hierarchy.AddTask(&task)
	}
	
	// Test calendar layout
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	taskPointers := make([]*Task, len(tasks))
	for i := range tasks {
		taskPointers[i] = &tasks[i]
	}
	calendarLayout := NewCalendarLayout(startDate, endDate, taskPointers)
	
	// Test timeline analysis
	timelineAnalyzer := NewTaskTimelineAnalyzer()
	
	// Analyze a few sample tasks
	sampleTasks := tasks[:min(5, len(tasks))]
	for _, task := range sampleTasks {
		analysis := timelineAnalyzer.AnalyzeTaskTimeline(&task)
		if analysis == nil {
			t.Errorf("Expected timeline analysis for task %s", task.ID)
			continue
		}
		
		t.Logf("Task %s (%s): Risk=%s, Progress=%.1f%%, WorkDays=%d/%d", 
			task.ID, task.Name, analysis.RiskLevel, analysis.ProgressPercent, 
			analysis.WorkDaysRemaining, analysis.WorkDays)
		
		// Verify analysis data
		if analysis.TaskID != task.ID {
			t.Errorf("Expected TaskID %s, got %s", task.ID, analysis.TaskID)
		}
		
		if analysis.TotalDays != task.GetDuration() {
			t.Errorf("Expected TotalDays %d, got %d", task.GetDuration(), analysis.TotalDays)
		}
	}
	
	// Test category filtering
	imagingTasks := collection.GetTasksByCategory("IMAGING")
	t.Logf("Found %d IMAGING tasks", len(imagingTasks))
	
	proposalTasks := collection.GetTasksByCategory("PROPOSAL")
	t.Logf("Found %d PROPOSAL tasks", len(proposalTasks))
	
	// Test date range filtering
	jan2025 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	dec2025 := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	tasks2025 := collection.GetTasksByDateRange(jan2025, dec2025)
	t.Logf("Found %d tasks in 2025", len(tasks2025))
	
	// Test dependency analysis
	rootTasks := hierarchy.GetRootTasks()
	t.Logf("Found %d root tasks", len(rootTasks))
	
	tasksWithDeps := 0
	for _, task := range tasks {
		if len(task.Dependencies) > 0 {
			tasksWithDeps++
		}
	}
	t.Logf("Found %d tasks with dependencies", tasksWithDeps)
	
	// Test milestone detection
	milestones := 0
	for _, task := range tasks {
		if task.IsMilestone {
			milestones++
		}
	}
	t.Logf("Found %d milestone tasks", milestones)
	
	// Test calendar layout
	months := calendarLayout.GetMonths()
	t.Logf("Calendar layout covers %d months", len(months))
	
	weeks := calendarLayout.GetWeeks()
	t.Logf("Calendar layout covers %d weeks", len(weeks))
	
	// Test task rendering
	for _, task := range sampleTasks {
		renderer := NewTaskRenderer(&task)
		if renderer.TaskID != task.ID {
			t.Errorf("Expected TaskID %s, got %s", task.ID, renderer.TaskID)
		}
		if !renderer.Visible {
			t.Error("Expected task to be visible")
		}
	}
}

// TestComplexDependencyScenarios tests complex dependency scenarios
func TestComplexDependencyScenarios(t *testing.T) {
	// Create a complex dependency graph
	tasks := []Task{
		{
			ID:          "A",
			Name:        "Root Task A",
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			Category:    "RESEARCH",
			Dependencies: []string{},
		},
		{
			ID:          "B",
			Name:        "Task B depends on A",
			StartDate:   time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			Category:    "IMAGING",
			Dependencies: []string{"A"},
		},
		{
			ID:          "C",
			Name:        "Task C depends on A and B",
			StartDate:   time.Date(2025, 1, 11, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:    "PUBLICATION",
			Dependencies: []string{"A", "B"},
		},
		{
			ID:          "D",
			Name:        "Task D depends on C",
			StartDate:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			Category:    "DISSERTATION",
			Dependencies: []string{"C"},
		},
		{
			ID:          "E",
			Name:        "Parallel Task E",
			StartDate:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
			Category:    "ADMIN",
			Dependencies: []string{},
		},
	}
	
	// Test dependency graph
	graph := NewTaskDependencyGraph()
	for _, task := range tasks {
		// Create a copy to avoid pointer issues
		taskCopy := task
		graph.AddTask(&taskCopy)
		t.Logf("Added task %s with dependencies: %v", taskCopy.ID, taskCopy.Dependencies)
	}
	
	// Test dependency levels
	levelA := graph.GetTaskLevel("A")
	if levelA != 0 {
		t.Errorf("Expected level 0 for A, got %d", levelA)
	}
	
	levelB := graph.GetTaskLevel("B")
	if levelB != 1 {
		t.Errorf("Expected level 1 for B, got %d", levelB)
	}
	
	levelC := graph.GetTaskLevel("C")
	if levelC != 2 {
		t.Errorf("Expected level 2 for C, got %d", levelC)
	}
	
	levelD := graph.GetTaskLevel("D")
	if levelD != 3 {
		t.Errorf("Expected level 3 for D, got %d", levelD)
	}
	
	levelE := graph.GetTaskLevel("E")
	if levelE != 0 {
		t.Errorf("Expected level 0 for E, got %d", levelE)
	}
	
	// Test tasks by level
	level0Tasks := graph.GetTasksByLevel(0)
	if len(level0Tasks) != 2 {
		t.Errorf("Expected 2 level 0 tasks, got %d", len(level0Tasks))
	}
	
	level1Tasks := graph.GetTasksByLevel(1)
	if len(level1Tasks) != 1 {
		t.Errorf("Expected 1 level 1 task, got %d", len(level1Tasks))
	}
	
	// Test dependencies
	depsB := graph.GetDependencies("B")
	if len(depsB) != 1 || depsB[0].ID != "A" {
		t.Errorf("Expected B to depend on A, got %v", depsB)
		t.Logf("Dependencies for B: %v", depsB)
	}
	
	depsC := graph.GetDependencies("C")
	if len(depsC) != 2 {
		t.Errorf("Expected C to have 2 dependencies, got %d", len(depsC))
	}
	
	// Test dependents
	depsOfA := graph.GetDependents("A")
	if len(depsOfA) != 2 {
		t.Errorf("Expected A to have 2 dependents, got %d", len(depsOfA))
	}
}

// TestHierarchicalScenarios tests hierarchical task scenarios
func TestHierarchicalScenarios(t *testing.T) {
	// Create hierarchical tasks
	tasks := []Task{
		{
			ID:        "PARENT1",
			Name:      "Parent Task 1",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			Category:  "RESEARCH",
		},
		{
			ID:        "CHILD1",
			Name:      "Child Task 1",
			ParentID:  "PARENT1",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
		},
		{
			ID:        "CHILD2",
			Name:      "Child Task 2",
			ParentID:  "PARENT1",
			StartDate: time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			Category:  "PUBLICATION",
		},
		{
			ID:        "GRANDCHILD1",
			Name:      "Grandchild Task 1",
			ParentID:  "CHILD1",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
			Category:  "ADMIN",
		},
		{
			ID:        "ROOT2",
			Name:      "Root Task 2",
			StartDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
			Category:  "DISSERTATION",
		},
	}
	
	// Test hierarchy
	hierarchy := NewTaskHierarchy()
	for _, task := range tasks {
		// Create a copy to avoid pointer issues
		taskCopy := task
		hierarchy.AddTask(&taskCopy)
	}
	
	// Test root tasks
	roots := hierarchy.GetRootTasks()
	if len(roots) != 2 {
		t.Errorf("Expected 2 root tasks, got %d", len(roots))
	}
	
	// Test children
	children := hierarchy.GetChildren("PARENT1")
	if len(children) != 2 {
		t.Errorf("Expected 2 children for PARENT1, got %d", len(children))
	}
	
	// Test parent
	parent := hierarchy.GetParent("CHILD1")
	if parent == nil || parent.ID != "PARENT1" {
		t.Errorf("Expected parent PARENT1 for CHILD1, got %v", parent)
	}
	
	// Test ancestors
	ancestors := hierarchy.GetAncestors("GRANDCHILD1")
	if len(ancestors) != 2 {
		t.Errorf("Expected 2 ancestors for GRANDCHILD1, got %d", len(ancestors))
	}
	
	// Test descendants
	descendants := hierarchy.GetDescendants("PARENT1")
	if len(descendants) != 3 {
		t.Errorf("Expected 3 descendants for PARENT1, got %d", len(descendants))
	}
}

// TestCalendarLayoutScenarios tests calendar layout scenarios
func TestCalendarLayoutScenarios(t *testing.T) {
	// Create tasks spanning multiple months
	tasks := []Task{
		{
			ID:        "JAN_TASK",
			Name:      "January Task",
			StartDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			Category:  "RESEARCH",
		},
		{
			ID:        "FEB_TASK",
			Name:      "February Task",
			StartDate: time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
		},
		{
			ID:        "CROSS_MONTH",
			Name:      "Cross Month Task",
			StartDate: time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC),
			Category:  "PUBLICATION",
		},
	}
	
	// Create calendar layout
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)
	taskPointers := make([]*Task, len(tasks))
	for i := range tasks {
		taskPointers[i] = &tasks[i]
	}
	layout := NewCalendarLayout(startDate, endDate, taskPointers)
	
	// Test months
	months := layout.GetMonths()
	if len(months) != 3 {
		t.Errorf("Expected 3 months, got %d", len(months))
	}
	
	// Test tasks for January
	janTasks := layout.GetTasksForMonth(2025, time.January)
	if len(janTasks) != 2 {
		t.Errorf("Expected 2 tasks in January, got %d", len(janTasks))
	}
	
	// Test tasks for February
	febTasks := layout.GetTasksForMonth(2025, time.February)
	if len(febTasks) != 2 {
		t.Errorf("Expected 2 tasks in February, got %d", len(febTasks))
	}
	
	// Test tasks for March
	marTasks := layout.GetTasksForMonth(2025, time.March)
	if len(marTasks) != 0 {
		t.Errorf("Expected 0 tasks in March, got %d", len(marTasks))
	}
	
	// Test weeks
	weeks := layout.GetWeeks()
	if len(weeks) < 12 {
		t.Errorf("Expected at least 12 weeks, got %d", len(weeks))
	}
	
	// Test specific week
	weekStart := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC) // Monday
	weekTasks := layout.GetTasksForWeek(weekStart)
	if len(weekTasks) < 1 {
		t.Errorf("Expected at least 1 task in week starting %v, got %d", weekStart, len(weekTasks))
	}
}

// TestTimelineAnalysisScenarios tests timeline analysis scenarios
func TestTimelineAnalysisScenarios(t *testing.T) {
	now := time.Now()
	
	// Create various task scenarios
	tasks := []Task{
		{
			ID:        "COMPLETED",
			Name:      "Completed Task",
			StartDate: now.AddDate(0, 0, -10),
			EndDate:   now.AddDate(0, 0, -5),
			Status:    "Completed",
			Category:  "RESEARCH",
		},
		{
			ID:        "OVERDUE",
			Name:      "Overdue Task",
			StartDate: now.AddDate(0, 0, -10),
			EndDate:   now.AddDate(0, 0, -2),
			Status:    "In Progress",
			Category:  "IMAGING",
		},
		{
			ID:        "UPCOMING",
			Name:      "Upcoming Task",
			StartDate: now.AddDate(0, 0, 3),
			EndDate:   now.AddDate(0, 0, 8),
			Status:    "Planned",
			Category:  "PUBLICATION",
		},
		{
			ID:        "HIGH_RISK",
			Name:      "High Risk Task",
			StartDate: now.Add(-1 * time.Hour),
			EndDate:   now.AddDate(0, 0, 1),
			Status:    "In Progress",
			Category:  "DISSERTATION",
			Priority:  9,
		},
		{
			ID:           "WITH_DEPS",
			Name:         "Task with Dependencies",
			StartDate:    now.AddDate(0, 0, 1),
			EndDate:      now.AddDate(0, 0, 10),
			Status:       "Planned",
			Category:     "ADMIN",
			Dependencies: []string{"COMPLETED", "UPCOMING"},
		},
	}
	
	// Test timeline analysis
	analyzer := NewTaskTimelineAnalyzer()
	
	for _, task := range tasks {
		analysis := analyzer.AnalyzeTaskTimeline(&task)
		if analysis == nil {
			t.Errorf("Expected analysis for task %s", task.ID)
			continue
		}
		
		t.Logf("Task %s: Risk=%s, Progress=%.1f%%, Overdue=%v, Upcoming=%v", 
			task.ID, analysis.RiskLevel, analysis.ProgressPercent, 
			analysis.IsOverdue, analysis.IsUpcoming)
		
		// Verify specific scenarios
		switch task.ID {
		case "COMPLETED":
			if analysis.IsOverdue {
				t.Error("Completed task should not be overdue")
			}
			if analysis.IsUpcoming {
				t.Error("Completed task should not be upcoming")
			}
		case "OVERDUE":
			if !analysis.IsOverdue {
				t.Error("Overdue task should be overdue")
			}
		case "UPCOMING":
			if !analysis.IsUpcoming {
				t.Error("Upcoming task should be upcoming")
			}
		case "HIGH_RISK":
			if analysis.RiskLevel != "HIGH" {
				t.Errorf("Expected HIGH risk for high risk task, got %s", analysis.RiskLevel)
			}
		case "WITH_DEPS":
			hasDepRecommendation := false
			for _, rec := range analysis.Recommendations {
				if contains(rec, "dependencies") {
					hasDepRecommendation = true
					break
				}
			}
			if !hasDepRecommendation {
				t.Error("Expected dependency recommendation for task with dependencies")
			}
		}
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

