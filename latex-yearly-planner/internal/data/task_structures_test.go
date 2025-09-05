package data

import (
	"testing"
	"time"
)

func TestTaskCategory(t *testing.T) {
	// Test predefined categories
	category := GetCategory("PROPOSAL")
	if category.Name != "PROPOSAL" {
		t.Errorf("Expected PROPOSAL, got %s", category.Name)
	}
	if category.Color != "#4A90E2" {
		t.Errorf("Expected #4A90E2, got %s", category.Color)
	}

	// Test unknown category
	unknown := GetCategory("UNKNOWN")
	if unknown.Name != "UNKNOWN" {
		t.Errorf("Expected UNKNOWN, got %s", unknown.Name)
	}
	if unknown.Color != "#CCCCCC" {
		t.Errorf("Expected default color #CCCCCC, got %s", unknown.Color)
	}

	// Test all categories
	categories := GetAllCategories()
	if len(categories) != 7 {
		t.Errorf("Expected 7 categories, got %d", len(categories))
	}
}

func TestTaskCollection(t *testing.T) {
	collection := NewTaskCollection()
	
	// Create test tasks
	task1 := &Task{
		ID:        "A",
		Name:      "Task A",
		Category:  "PROPOSAL",
		Status:    "Planned",
		Assignee:  "John",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &Task{
		ID:        "B",
		Name:      "Task B",
		Category:  "LASER",
		Status:    "In Progress",
		Assignee:  "Jane",
		StartDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
	}
	
	// Add tasks
	collection.AddTask(task1)
	collection.AddTask(task2)
	
	// Test retrieval
	retrieved, exists := collection.GetTask("A")
	if !exists {
		t.Error("Expected task A to exist")
	}
	if retrieved.Name != "Task A" {
		t.Errorf("Expected Task A, got %s", retrieved.Name)
	}
	
	// Test category filtering
	proposalTasks := collection.GetTasksByCategory("PROPOSAL")
	if len(proposalTasks) != 1 {
		t.Errorf("Expected 1 PROPOSAL task, got %d", len(proposalTasks))
	}
	
	// Test status filtering
	inProgressTasks := collection.GetTasksByStatus("In Progress")
	if len(inProgressTasks) != 1 {
		t.Errorf("Expected 1 In Progress task, got %d", len(inProgressTasks))
	}
	
	// Test assignee filtering
	janeTasks := collection.GetTasksByAssignee("Jane")
	if len(janeTasks) != 1 {
		t.Errorf("Expected 1 Jane task, got %d", len(janeTasks))
	}
	
	// Test date range filtering
	start := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)
	overlappingTasks := collection.GetTasksByDateRange(start, end)
	if len(overlappingTasks) != 2 {
		t.Errorf("Expected 2 overlapping tasks, got %d", len(overlappingTasks))
	}
	
	// Test specific date filtering
	date := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)
	tasksOnDate := collection.GetTasksByDate(date)
	if len(tasksOnDate) != 2 {
		t.Errorf("Expected 2 tasks on date, got %d", len(tasksOnDate))
	}
}

func TestTaskDependencyGraph(t *testing.T) {
	graph := NewTaskDependencyGraph()
	
	// Create test tasks with dependencies
	taskA := &Task{
		ID:          "A",
		Name:        "Task A",
		Dependencies: []string{},
	}
	
	taskB := &Task{
		ID:          "B",
		Name:        "Task B",
		Dependencies: []string{"A"},
	}
	
	taskC := &Task{
		ID:          "C",
		Name:        "Task C",
		Dependencies: []string{"A", "B"},
	}
	
	// Add tasks
	graph.AddTask(taskA)
	graph.AddTask(taskB)
	graph.AddTask(taskC)
	
	// Test dependencies
	depsB := graph.GetDependencies("B")
	if len(depsB) != 1 {
		t.Errorf("Expected 1 dependency for B, got %d", len(depsB))
	}
	if depsB[0].ID != "A" {
		t.Errorf("Expected dependency A for B, got %s", depsB[0].ID)
	}
	
	// Test dependents
	depsOfA := graph.GetDependents("A")
	if len(depsOfA) != 2 {
		t.Errorf("Expected 2 dependents for A, got %d", len(depsOfA))
	}
	
	// Test levels
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
	
	// Test tasks by level
	level0Tasks := graph.GetTasksByLevel(0)
	if len(level0Tasks) != 1 {
		t.Errorf("Expected 1 level 0 task, got %d", len(level0Tasks))
	}
}

func TestTaskHierarchy(t *testing.T) {
	hierarchy := NewTaskHierarchy()
	
	// Create test tasks with parent-child relationships
	parent := &Task{
		ID:       "PARENT",
		Name:     "Parent Task",
		ParentID: "",
	}
	
	child1 := &Task{
		ID:       "CHILD1",
		Name:     "Child 1",
		ParentID: "PARENT",
	}
	
	child2 := &Task{
		ID:       "CHILD2",
		Name:     "Child 2",
		ParentID: "PARENT",
	}
	
	// Add tasks
	hierarchy.AddTask(parent)
	hierarchy.AddTask(child1)
	hierarchy.AddTask(child2)
	
	// Test root tasks
	roots := hierarchy.GetRootTasks()
	if len(roots) != 1 {
		t.Errorf("Expected 1 root task, got %d", len(roots))
	}
	
	// Test children
	children := hierarchy.GetChildren("PARENT")
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
	
	// Test parent
	parentOfChild := hierarchy.GetParent("CHILD1")
	if parentOfChild == nil || parentOfChild.ID != "PARENT" {
		t.Errorf("Expected parent PARENT for CHILD1, got %v", parentOfChild)
	}
	
	// Test ancestors
	ancestors := hierarchy.GetAncestors("CHILD1")
	if len(ancestors) != 1 {
		t.Errorf("Expected 1 ancestor for CHILD1, got %d", len(ancestors))
	}
	
	// Test descendants
	descendants := hierarchy.GetDescendants("PARENT")
	if len(descendants) != 2 {
		t.Errorf("Expected 2 descendants for PARENT, got %d", len(descendants))
	}
}

func TestCalendarLayout(t *testing.T) {
	// Create test tasks
	task1 := &Task{
		ID:        "A",
		Name:      "Task A",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	task2 := &Task{
		ID:        "B",
		Name:      "Task B",
		StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	
	tasks := []*Task{task1, task2}
	
	// Create calendar layout
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	layout := NewCalendarLayout(start, end, tasks)
	
	// Test months
	months := layout.GetMonths()
	if len(months) != 1 {
		t.Errorf("Expected 1 month, got %d", len(months))
	}
	if months[0].Year != 2024 || months[0].Month != time.January {
		t.Errorf("Expected January 2024, got %v", months[0])
	}
	
	// Test weeks
	weeks := layout.GetWeeks()
	if len(weeks) < 4 {
		t.Errorf("Expected at least 4 weeks, got %d", len(weeks))
	}
	
	// Test tasks for month
	monthTasks := layout.GetTasksForMonth(2024, time.January)
	if len(monthTasks) != 2 {
		t.Errorf("Expected 2 tasks for January, got %d", len(monthTasks))
	}
	
	// Test tasks for week
	weekStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	weekTasks := layout.GetTasksForWeek(weekStart)
	if len(weekTasks) != 1 {
		t.Errorf("Expected 1 task for first week, got %d", len(weekTasks))
	}
}

func TestTaskRenderer(t *testing.T) {
	task := &Task{
		ID:       "A",
		Name:     "Task A",
		Category: "PROPOSAL",
	}
	
	renderer := NewTaskRenderer(task)
	
	if renderer.TaskID != "A" {
		t.Errorf("Expected TaskID A, got %s", renderer.TaskID)
	}
	
	if renderer.Color != "#4A90E2" {
		t.Errorf("Expected PROPOSAL color #4A90E2, got %s", renderer.Color)
	}
	
	if !renderer.Visible {
		t.Error("Expected task to be visible")
	}
	
	if renderer.Opacity != 1.0 {
		t.Errorf("Expected opacity 1.0, got %f", renderer.Opacity)
	}
}

func TestTaskMethods(t *testing.T) {
	now := time.Now()
	task := &Task{
		ID:        "A",
		Name:      "Task A",
		Category:  "PROPOSAL",
		Status:    "Planned",
		StartDate: now.AddDate(0, 0, -2), // Started 2 days ago
		EndDate:   now.AddDate(0, 0, 3),  // Ends in 3 days
	}
	
	// Test IsOnDate
	date1 := now.AddDate(0, 0, 1) // Tomorrow
	if !task.IsOnDate(date1) {
		t.Error("Expected task to be on date tomorrow")
	}
	
	date2 := now.AddDate(0, 0, 10) // 10 days from now
	if task.IsOnDate(date2) {
		t.Error("Expected task not to be on date 10 days from now")
	}
	
	// Test OverlapsWithDateRange
	start := now.AddDate(0, 0, -1) // Yesterday
	end := now.AddDate(0, 0, 2)    // Day after tomorrow
	if !task.OverlapsWithDateRange(start, end) {
		t.Error("Expected task to overlap with date range")
	}
	
	// Test GetDuration
	duration := task.GetDuration()
	if duration != 6 {
		t.Errorf("Expected duration 6 days, got %d", duration)
	}
	
	// Test GetCategoryInfo
	category := task.GetCategoryInfo()
	if category.Name != "PROPOSAL" {
		t.Errorf("Expected category PROPOSAL, got %s", category.Name)
	}
	
	// Test IsOverdue (should not be overdue since end date is in the future)
	if task.IsOverdue() {
		t.Error("Expected task not to be overdue")
	}
	
	// Test String representation
	str := task.String()
	if !contains(str, "Task[A: Task A") {
		t.Errorf("Expected string to contain task info, got %s", str)
	}
}

func TestTaskProgress(t *testing.T) {
	// Test task in progress
	now := time.Now()
	task := &Task{
		ID:        "A",
		Name:      "Task A",
		StartDate: now.AddDate(0, 0, -5), // Started 5 days ago
		EndDate:   now.AddDate(0, 0, 5),  // Ends in 5 days
		Status:    "In Progress",
	}
	
	progress := task.GetProgressPercentage()
	if progress < 0 || progress > 100 {
		t.Errorf("Expected progress between 0 and 100, got %f", progress)
	}
	
	// Test upcoming task
	upcomingTask := &Task{
		ID:        "B",
		Name:      "Task B",
		StartDate: now.AddDate(0, 0, 3), // Starts in 3 days
		EndDate:   now.AddDate(0, 0, 8), // Ends in 8 days
		Status:    "Planned",
	}
	
	if !upcomingTask.IsUpcoming() {
		t.Error("Expected upcoming task to be upcoming")
	}
}

// Helper function for string contains
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && contains(s[1:], substr)
}
