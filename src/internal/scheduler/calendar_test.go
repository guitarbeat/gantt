package scheduler

import (
	"testing"
	"time"
	"phd-dissertation-planner/internal/common"
)

func TestNewCalendarLayout(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	
	calendar := NewCalendarLayout(start, end, []*common.Task{})
	if calendar == nil {
		t.Fatal("Expected calendar to be created, got nil")
	}
	
	if len(calendar.tasks) != 0 {
		t.Errorf("Expected empty tasks slice, got %d tasks", len(calendar.tasks))
	}
	
	if len(calendar.months) != 3 {
		t.Errorf("Expected 3 months, got %d", len(calendar.months))
	}
}

func TestGetTasksForMonth(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	
	task := &common.Task{
		ID:        "1",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	calendar := NewCalendarLayout(start, end, []*common.Task{task})
	
	tasks := calendar.GetTasksForMonth(2024, time.January)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task for January, got %d", len(tasks))
	}
}

func TestGetTasksForWeek(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	task := &common.Task{
		ID:        "1",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	calendar := NewCalendarLayout(start, end, []*common.Task{task})
	
	weekStart := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	tasks := calendar.GetTasksForWeek(weekStart)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task for week, got %d", len(tasks))
	}
}