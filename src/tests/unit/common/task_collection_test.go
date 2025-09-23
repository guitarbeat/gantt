package common_test

import (
	"phd-dissertation-planner/internal/common"
	"testing"
	"time"
)

func TestNewTaskCollection(t *testing.T) {
	collection := common.NewTaskCollection()
	if collection == nil {
		t.Fatal("Expected collection to be created, got nil")
	}
	
	// Test that collection was created successfully
	t.Log("Task collection created successfully")
}

func TestAddTask(t *testing.T) {
	collection := common.NewTaskCollection()
	
	task := &common.Task{
		ID:       "1",
		Name:     "Test Task",
		Category: "PROPOSAL",
		Status:   "Planned",
		Assignee: "John",
	}
	
	collection.AddTask(task)
	
	// Test that task was added by checking if we can retrieve it
	allTasks := collection.GetAllTasks()
	if len(allTasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(allTasks))
	}
	
	// Test category filtering
	proposalTasks := collection.GetTasksByCategory("PROPOSAL")
	if len(proposalTasks) != 1 {
		t.Errorf("Expected 1 task in PROPOSAL category, got %d", len(proposalTasks))
	}
	
	// Test status filtering
	plannedTasks := collection.GetTasksByStatus("Planned")
	if len(plannedTasks) != 1 {
		t.Errorf("Expected 1 task with Planned status, got %d", len(plannedTasks))
	}
	
	// Test assignee filtering
	johnTasks := collection.GetTasksByAssignee("John")
	if len(johnTasks) != 1 {
		t.Errorf("Expected 1 task assigned to John, got %d", len(johnTasks))
	}
}

func TestAddTaskNil(t *testing.T) {
	collection := common.NewTaskCollection()
	
	collection.AddTask(nil)
	
	// Test that no task was added
	allTasks := collection.GetAllTasks()
	if len(allTasks) != 0 {
		t.Errorf("Expected 0 tasks after adding nil, got %d", len(allTasks))
	}
}

func TestGetTask(t *testing.T) {
	collection := common.NewTaskCollection()
	
	task := &common.Task{
		ID:   "1",
		Name: "Test Task",
	}
	
	collection.AddTask(task)
	
	retrievedTask, found := collection.GetTask("Test Task")
	if !found {
		t.Error("Expected task to be found")
	}
	
	if retrievedTask != task {
		t.Error("Expected retrieved task to be the same as added task")
	}
	
	_, found = collection.GetTask("Non-existent Task")
	if found {
		t.Error("Expected non-existent task to not be found")
	}
}

func TestGetAllTasks(t *testing.T) {
	collection := common.NewTaskCollection()
	
	task1 := &common.Task{ID: "1", Name: "Task 1"}
	task2 := &common.Task{ID: "2", Name: "Task 2"}
	
	collection.AddTask(task1)
	collection.AddTask(task2)
	
	allTasks := collection.GetAllTasks()
	if len(allTasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(allTasks))
	}
}

func TestGetTasksByCategory(t *testing.T) {
	collection := common.NewTaskCollection()
	
	task1 := &common.Task{ID: "1", Name: "Task 1", Category: "PROPOSAL"}
	task2 := &common.Task{ID: "2", Name: "Task 2", Category: "RESEARCH"}
	task3 := &common.Task{ID: "3", Name: "Task 3", Category: "PROPOSAL"}
	
	collection.AddTask(task1)
	collection.AddTask(task2)
	collection.AddTask(task3)
	
	proposalTasks := collection.GetTasksByCategory("PROPOSAL")
	if len(proposalTasks) != 2 {
		t.Errorf("Expected 2 PROPOSAL tasks, got %d", len(proposalTasks))
	}
	
	researchTasks := collection.GetTasksByCategory("RESEARCH")
	if len(researchTasks) != 1 {
		t.Errorf("Expected 1 RESEARCH task, got %d", len(researchTasks))
	}
	
	adminTasks := collection.GetTasksByCategory("ADMIN")
	if len(adminTasks) != 0 {
		t.Errorf("Expected 0 ADMIN tasks, got %d", len(adminTasks))
	}
}

func TestGetTasksByStatus(t *testing.T) {
	collection := common.NewTaskCollection()
	
	task1 := &common.Task{ID: "1", Name: "Task 1", Status: "Planned"}
	task2 := &common.Task{ID: "2", Name: "Task 2", Status: "In Progress"}
	task3 := &common.Task{ID: "3", Name: "Task 3", Status: "Planned"}
	
	collection.AddTask(task1)
	collection.AddTask(task2)
	collection.AddTask(task3)
	
	plannedTasks := collection.GetTasksByStatus("Planned")
	if len(plannedTasks) != 2 {
		t.Errorf("Expected 2 Planned tasks, got %d", len(plannedTasks))
	}
	
	inProgressTasks := collection.GetTasksByStatus("In Progress")
	if len(inProgressTasks) != 1 {
		t.Errorf("Expected 1 In Progress task, got %d", len(inProgressTasks))
	}
}

func TestGetTasksByAssignee(t *testing.T) {
	collection := common.NewTaskCollection()
	
	task1 := &common.Task{ID: "1", Name: "Task 1", Assignee: "John"}
	task2 := &common.Task{ID: "2", Name: "Task 2", Assignee: "Jane"}
	task3 := &common.Task{ID: "3", Name: "Task 3", Assignee: "John"}
	
	collection.AddTask(task1)
	collection.AddTask(task2)
	collection.AddTask(task3)
	
	johnTasks := collection.GetTasksByAssignee("John")
	if len(johnTasks) != 2 {
		t.Errorf("Expected 2 tasks assigned to John, got %d", len(johnTasks))
	}
	
	janeTasks := collection.GetTasksByAssignee("Jane")
	if len(janeTasks) != 1 {
		t.Errorf("Expected 1 task assigned to Jane, got %d", len(janeTasks))
	}
}

func TestGetTasksByDateRange(t *testing.T) {
	collection := common.NewTaskCollection()
	
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	task1 := &common.Task{
		ID:        "1",
		Name:      "Task 1",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	task2 := &common.Task{
		ID:        "2",
		Name:      "Task 2",
		StartDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
	}
	task3 := &common.Task{
		ID:        "3",
		Name:      "Task 3",
		StartDate: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
	}
	
	collection.AddTask(task1)
	collection.AddTask(task2)
	collection.AddTask(task3)
	
	tasksInRange := collection.GetTasksByDateRange(start, end)
	if len(tasksInRange) != 2 {
		t.Errorf("Expected 2 tasks in date range, got %d", len(tasksInRange))
	}
}

func TestGetTasksByDate(t *testing.T) {
	collection := common.NewTaskCollection()
	
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	task1 := &common.Task{
		ID:        "1",
		Name:      "Task 1",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	task2 := &common.Task{
		ID:        "2",
		Name:      "Task 2",
		StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	task3 := &common.Task{
		ID:        "3",
		Name:      "Task 3",
		StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
	}
	
	collection.AddTask(task1)
	collection.AddTask(task2)
	collection.AddTask(task3)
	
	tasksOnDate := collection.GetTasksByDate(date)
	if len(tasksOnDate) != 2 {
		t.Errorf("Expected 2 tasks on date, got %d", len(tasksOnDate))
	}
}

func TestSortByDate(t *testing.T) {
	collection := common.NewTaskCollection()
	
	task1 := &common.Task{
		ID:        "1",
		Name:      "Task 1",
		StartDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
	}
	task2 := &common.Task{
		ID:        "2",
		Name:      "Task 2",
		StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	task3 := &common.Task{
		ID:        "3",
		Name:      "Task 3",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC),
	}
	
	collection.AddTask(task1)
	collection.AddTask(task2)
	collection.AddTask(task3)
	
	// Test that tasks were added successfully
	allTasks := collection.GetAllTasks()
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(allTasks))
	}
	
	// Test that we can get tasks by date range
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	tasksInRange := collection.GetTasksByDateRange(start, end)
	if len(tasksInRange) != 3 {
		t.Errorf("Expected 3 tasks in date range, got %d", len(tasksInRange))
	}
}

func TestNewTaskHierarchy(t *testing.T) {
	hierarchy := common.NewTaskHierarchy()
	if hierarchy == nil {
		t.Fatal("Expected hierarchy to be created, got nil")
	}
	
	// Test that hierarchy is properly initialized
	rootTasks := hierarchy.GetRootTasks()
	if len(rootTasks) != 0 {
		t.Errorf("Expected empty roots slice, got %d roots", len(rootTasks))
	}
}

func TestAddTaskToHierarchy(t *testing.T) {
	hierarchy := common.NewTaskHierarchy()
	
	// Add root task
	rootTask := &common.Task{ID: "1", Name: "Root Task"}
	hierarchy.AddTask(rootTask)
	
	rootTasks := hierarchy.GetRootTasks()
	if len(rootTasks) != 1 {
		t.Errorf("Expected 1 root task, got %d", len(rootTasks))
	}
	
	// Add child task
	childTask := &common.Task{ID: "2", Name: "Child Task", ParentID: "Root Task"}
	hierarchy.AddTask(childTask)
	
	// Test that child task was added
	children := hierarchy.GetChildren("Root Task")
	if len(children) != 1 {
		t.Errorf("Expected 1 child task, got %d", len(children))
	}
	
	// Check parent-child relationship
	parent := hierarchy.GetParent("Child Task")
	if parent == nil {
		t.Error("Expected parent to be found for child task")
	}
	
	if parent.Name != "Root Task" {
		t.Errorf("Expected parent to be 'Root Task', got %s", parent.Name)
	}
	
	children2 := hierarchy.GetChildren("Root Task")
	if len(children2) != 1 {
		t.Errorf("Expected 1 child task, got %d", len(children2))
	}
	
	if children2[0].Name != "Child Task" {
		t.Errorf("Expected child to be 'Child Task', got %s", children2[0].Name)
	}
}

func TestGetRootTasks(t *testing.T) {
	hierarchy := common.NewTaskHierarchy()
	
	rootTask1 := &common.Task{ID: "1", Name: "Root Task 1"}
	rootTask2 := &common.Task{ID: "2", Name: "Root Task 2"}
	childTask := &common.Task{ID: "3", Name: "Child Task", ParentID: "Root Task 1"}
	
	hierarchy.AddTask(rootTask1)
	hierarchy.AddTask(rootTask2)
	hierarchy.AddTask(childTask)
	
	rootTasks := hierarchy.GetRootTasks()
	if len(rootTasks) != 2 {
		t.Errorf("Expected 2 root tasks, got %d", len(rootTasks))
	}
}

func TestGetAncestors(t *testing.T) {
	hierarchy := common.NewTaskHierarchy()
	
	grandparent := &common.Task{ID: "1", Name: "Grandparent"}
	parent := &common.Task{ID: "2", Name: "Parent", ParentID: "Grandparent"}
	child := &common.Task{ID: "3", Name: "Child", ParentID: "Parent"}
	
	hierarchy.AddTask(grandparent)
	hierarchy.AddTask(parent)
	hierarchy.AddTask(child)
	
	ancestors := hierarchy.GetAncestors("Child")
	if len(ancestors) != 2 {
		t.Errorf("Expected 2 ancestors, got %d", len(ancestors))
	}
	
	// Check order (should be parent first, then grandparent)
	if ancestors[0].Name != "Parent" {
		t.Errorf("Expected first ancestor to be 'Parent', got %s", ancestors[0].Name)
	}
	
	if ancestors[1].Name != "Grandparent" {
		t.Errorf("Expected second ancestor to be 'Grandparent', got %s", ancestors[1].Name)
	}
}

func TestGetDescendants(t *testing.T) {
	hierarchy := common.NewTaskHierarchy()
	
	parent := &common.Task{ID: "1", Name: "Parent"}
	child1 := &common.Task{ID: "2", Name: "Child 1", ParentID: "Parent"}
	child2 := &common.Task{ID: "3", Name: "Child 2", ParentID: "Parent"}
	grandchild := &common.Task{ID: "4", Name: "Grandchild", ParentID: "Child 1"}
	
	hierarchy.AddTask(parent)
	hierarchy.AddTask(child1)
	hierarchy.AddTask(child2)
	hierarchy.AddTask(grandchild)
	
	descendants := hierarchy.GetDescendants("Parent")
	if len(descendants) != 3 {
		t.Errorf("Expected 3 descendants, got %d", len(descendants))
	}
}

func TestTaskOverlapsWithDateRange(t *testing.T) {
	task := &common.Task{
		StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	// Test overlapping range
	start := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)
	
	if !task.OverlapsWithDateRange(start, end) {
		t.Error("Expected task to overlap with date range")
	}
	
	// Test non-overlapping range
	start = time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)
	end = time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC)
	
	if task.OverlapsWithDateRange(start, end) {
		t.Error("Expected task to not overlap with date range")
	}
}

func TestTaskIsOnDate(t *testing.T) {
	task := &common.Task{
		StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	// Test date within range
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !task.IsOnDate(date) {
		t.Error("Expected task to be on date")
	}
	
	// Test date before range
	date = time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	if task.IsOnDate(date) {
		t.Error("Expected task to not be on date")
	}
	
	// Test date after range
	date = time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)
	if task.IsOnDate(date) {
		t.Error("Expected task to not be on date")
	}
}

func TestTaskGetDuration(t *testing.T) {
	task := &common.Task{
		StartDate: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	duration := task.GetDuration()
	expected := 11 // 10 days + 1 (inclusive)
	
	if duration != expected {
		t.Errorf("Expected duration %d, got %d", expected, duration)
	}
}

func TestTaskGetCategoryInfo(t *testing.T) {
	task := &common.Task{Category: "PROPOSAL"}
	
	category := task.GetCategoryInfo()
	if category.Name != "PROPOSAL" {
		t.Errorf("Expected category name 'PROPOSAL', got %s", category.Name)
	}
	
	if category.Color != "#4A90E2" {
		t.Errorf("Expected category color '#4A90E2', got %s", category.Color)
	}
}

func TestTaskIsOverdue(t *testing.T) {
	// Test overdue task
	pastTask := &common.Task{
		EndDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Status:  "Planned",
	}
	
	if !pastTask.IsOverdue() {
		t.Error("Expected past task to be overdue")
	}
	
	// Test completed task (not overdue)
	completedTask := &common.Task{
		EndDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Status:  "Completed",
	}
	
	if completedTask.IsOverdue() {
		t.Error("Expected completed task to not be overdue")
	}
	
	// Test future task (not overdue)
	futureTask := &common.Task{
		EndDate: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
		Status:  "Planned",
	}
	
	if futureTask.IsOverdue() {
		t.Error("Expected future task to not be overdue")
	}
}

func TestTaskIsUpcoming(t *testing.T) {
	// Test upcoming task (within 7 days)
	upcomingTask := &common.Task{
		StartDate: time.Now().AddDate(0, 0, 3), // 3 days from now
	}
	
	if !upcomingTask.IsUpcoming() {
		t.Error("Expected task starting in 3 days to be upcoming")
	}
	
	// Test task starting in 10 days (not upcoming)
	futureTask := &common.Task{
		StartDate: time.Now().AddDate(0, 0, 10), // 10 days from now
	}
	
	if futureTask.IsUpcoming() {
		t.Error("Expected task starting in 10 days to not be upcoming")
	}
	
	// Test task that already started (not upcoming)
	pastTask := &common.Task{
		StartDate: time.Now().AddDate(0, 0, -5), // 5 days ago
	}
	
	if pastTask.IsUpcoming() {
		t.Error("Expected task that already started to not be upcoming")
	}
}

func TestTaskGetProgressPercentage(t *testing.T) {
	// Test task in progress
	now := time.Now()
	task := &common.Task{
		StartDate: now.AddDate(0, 0, -5), // Started 5 days ago
		EndDate:   now.AddDate(0, 0, 5),  // Ends in 5 days
	}
	
	progress := task.GetProgressPercentage()
	if progress < 40 || progress > 60 {
		t.Errorf("Expected progress around 50%%, got %.2f%%", progress)
	}
	
	// Test task that hasn't started
	notStartedTask := &common.Task{
		StartDate: now.AddDate(0, 0, 5), // Starts in 5 days
		EndDate:   now.AddDate(0, 0, 15), // Ends in 15 days
	}
	
	progress = notStartedTask.GetProgressPercentage()
	if progress != 0.0 {
		t.Errorf("Expected progress 0%%, got %.2f%%", progress)
	}
	
	// Test completed task
	completedTask := &common.Task{
		StartDate: now.AddDate(0, 0, -10), // Started 10 days ago
		EndDate:   now.AddDate(0, 0, -5),  // Ended 5 days ago
	}
	
	progress = completedTask.GetProgressPercentage()
	if progress != 100.0 {
		t.Errorf("Expected progress 100%%, got %.2f%%", progress)
	}
}

func TestTaskString(t *testing.T) {
	task := &common.Task{
		Name:      "Test Task",
		Category:  "PROPOSAL",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	expected := "Task[Test Task (PROPOSAL) 2024-01-15 - 2024-01-20]"
	result := task.String()
	
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
