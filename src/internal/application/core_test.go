package application

import (
	"testing"
	"time"
	"phd-dissertation-planner/internal/common"
)

func TestNewCore(t *testing.T) {
	core := NewCore()
	if core == nil {
		t.Fatal("Expected core to be created, got nil")
	}
	
	if core.calendar == nil {
		t.Error("Expected calendar to be initialized")
	}
	
	if core.taskCollection == nil {
		t.Error("Expected task collection to be initialized")
	}
	
	if core.taskHierarchy == nil {
		t.Error("Expected task hierarchy to be initialized")
	}
}

func TestAddTask(t *testing.T) {
	core := NewCore()
	
	task := &common.Task{
		ID:        "1",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	core.AddTask(task)
	
	// Check if task was added to collection
	allTasks := core.taskCollection.GetAllTasks()
	if len(allTasks) != 1 {
		t.Errorf("Expected 1 task in collection, got %d", len(allTasks))
	}
	
	// Check if task was added to calendar
	calendarTasks := core.calendar.GetAllTasks()
	if len(calendarTasks) != 1 {
		t.Errorf("Expected 1 task in calendar, got %d", len(calendarTasks))
	}
}

func TestGetAllTasks(t *testing.T) {
	core := NewCore()
	
	task1 := &common.Task{ID: "1", Name: "Task 1"}
	task2 := &common.Task{ID: "2", Name: "Task 2"}
	
	core.AddTask(task1)
	core.AddTask(task2)
	
	tasks := core.GetAllTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestGetTasksByCategory(t *testing.T) {
	core := NewCore()
	
	task1 := &common.Task{ID: "1", Name: "Task 1", Category: "PROPOSAL"}
	task2 := &common.Task{ID: "2", Name: "Task 2", Category: "RESEARCH"}
	
	core.AddTask(task1)
	core.AddTask(task2)
	
	proposalTasks := core.GetTasksByCategory("PROPOSAL")
	if len(proposalTasks) != 1 {
		t.Errorf("Expected 1 PROPOSAL task, got %d", len(proposalTasks))
	}
}

func TestGenerateCalendar(t *testing.T) {
	core := NewCore()
	
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	
	calendar := core.GenerateCalendar(start, end)
	if calendar == nil {
		t.Fatal("Expected calendar to be generated, got nil")
	}
}
