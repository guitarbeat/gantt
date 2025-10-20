package calendar_test

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/calendar"
)

// TestTaskStackerBasic tests basic stacking functionality
func TestTaskStackerBasic(t *testing.T) {
	// Create test tasks
	task1 := &SpanningTask{
		ID:        "T1",
		Name:      "Task 1",
		StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		Color:     "#FF0000",
	}

	task2 := &SpanningTask{
		ID:        "T2",
		Name:      "Task 2",
		StartDate: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
		Color:     "#00FF00",
	}

	tasks := []*SpanningTask{task1, task2}
	stacker := calendar.NewTaskStacker(tasks, time.Monday)
	stacker.ComputeStacks()

	// Task 1 should be in track 0
	day1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	stacks1 := stacker.GetStacksForDay(day1)
	if len(stacks1) != 1 {
		t.Errorf("Expected 1 task on day 1, got %d", len(stacks1))
	}
	if stacks1[0].Track != 0 {
		t.Errorf("Expected Task 1 in track 0, got track %d", stacks1[0].Track)
	}

	// Day 3 should have both tasks but in different tracks
	day3 := time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)
	stacks3 := stacker.GetStacksForDay(day3)
	if len(stacks3) != 2 {
		t.Errorf("Expected 2 tasks on day 3, got %d", len(stacks3))
	}

	// Tasks should be in different tracks
	tracks := make(map[int]bool)
	for _, stack := range stacks3 {
		tracks[stack.Track] = true
	}
	if len(tracks) != 2 {
		t.Errorf("Expected 2 different tracks on day 3, got %d", len(tracks))
	}
}

// TestTaskStackerOverlap tests overlap detection
func TestTaskStackerOverlap(t *testing.T) {
	// Create three overlapping tasks that all start on different days
	task1 := &SpanningTask{
		ID:        "T1",
		Name:      "Task 1",
		StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		Color:     "#FF0000",
	}

	task2 := &SpanningTask{
		ID:        "T2",
		Name:      "Task 2",
		StartDate: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC),
		Color:     "#00FF00",
	}

	task3 := &SpanningTask{
		ID:        "T3",
		Name:      "Task 3",
		StartDate: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		Color:     "#0000FF",
	}

	tasks := []*SpanningTask{task1, task2, task3}
	stacker := calendar.NewTaskStacker(tasks, time.Monday)
	stacker.ComputeStacks()

	// Day 7 should have all three tasks
	day7 := time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC)
	stacks7 := stacker.GetStacksForDay(day7)

	if len(stacks7) != 3 {
		t.Errorf("Expected 3 tasks on day 7, got %d", len(stacks7))
	}

	// All tasks should be in different tracks
	tracks := make(map[int]bool)
	for _, stack := range stacks7 {
		tracks[stack.Track] = true
	}
	if len(tracks) != 3 {
		t.Errorf("Expected 3 different tracks on day 7, got %d unique tracks", len(tracks))
	}

	// Max tracks should be at least 3
	if stacker.GetMaxTracks() < 3 {
		t.Errorf("Expected at least 3 tracks, got %d", stacker.GetMaxTracks())
	}
}

// TestTaskStackerStartingTasks tests filtering for tasks starting on a specific day
func TestTaskStackerStartingTasks(t *testing.T) {
	task1 := &SpanningTask{
		ID:        "T1",
		Name:      "Task 1",
		StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		Color:     "#FF0000",
	}

	task2 := &SpanningTask{
		ID:        "T2",
		Name:      "Task 2",
		StartDate: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
		Color:     "#00FF00",
	}

	tasks := []*SpanningTask{task1, task2}
	stacker := calendar.NewTaskStacker(tasks, time.Monday)
	stacker.ComputeStacks()

	// Day 1: Only Task 1 starts
	day1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	starting1 := stacker.GetTasksStartingOnDay(day1)
	if len(starting1) != 1 {
		t.Errorf("Expected 1 task starting on day 1, got %d", len(starting1))
	}
	if starting1[0].Task.ID != "T1" {
		t.Errorf("Expected Task 1 starting on day 1, got %s", starting1[0].Task.ID)
	}

	// Day 3: Only Task 2 starts (Task 1 is continuing)
	day3 := time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)
	starting3 := stacker.GetTasksStartingOnDay(day3)
	if len(starting3) != 1 {
		t.Errorf("Expected 1 task starting on day 3, got %d", len(starting3))
	}
	if starting3[0].Task.ID != "T2" {
		t.Errorf("Expected Task 2 starting on day 3, got %s", starting3[0].Task.ID)
	}

	// But both tasks should be visible on day 3
	all3 := stacker.GetStacksForDay(day3)
	if len(all3) != 2 {
		t.Errorf("Expected 2 tasks visible on day 3, got %d", len(all3))
	}
}

// TestTaskStackerNoOverlap tests tasks that don't overlap
func TestTaskStackerNoOverlap(t *testing.T) {
	task1 := &SpanningTask{
		ID:        "T1",
		Name:      "Task 1",
		StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		Color:     "#FF0000",
	}

	task2 := &SpanningTask{
		ID:        "T2",
		Name:      "Task 2",
		StartDate: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
		Color:     "#00FF00",
	}

	tasks := []*SpanningTask{task1, task2}
	stacker := calendar.NewTaskStacker(tasks, time.Monday)
	stacker.ComputeStacks()

	// Both tasks should be in track 0 since they don't overlap
	day1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	stacks1 := stacker.GetStacksForDay(day1)
	if len(stacks1) != 1 || stacks1[0].Track != 0 {
		t.Errorf("Expected Task 1 in track 0 on day 1")
	}

	day5 := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
	stacks5 := stacker.GetStacksForDay(day5)
	if len(stacks5) != 1 || stacks5[0].Track != 0 {
		t.Errorf("Expected Task 2 in track 0 on day 5")
	}

	// Max tracks should be 1 (only track 0 is used)
	if stacker.GetMaxTracks() != 1 {
		t.Errorf("Expected 1 track for non-overlapping tasks, got %d", stacker.GetMaxTracks())
	}
}

// TestTaskStackerSingleDay tests single-day tasks
func TestTaskStackerSingleDay(t *testing.T) {
	task := &SpanningTask{
		ID:        "T1",
		Name:      "Task 1",
		StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Color:     "#FF0000",
	}

	tasks := []*SpanningTask{task}
	stacker := calendar.NewTaskStacker(tasks, time.Monday)
	stacker.ComputeStacks()

	day1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	stacks := stacker.GetStacksForDay(day1)

	if len(stacks) != 1 {
		t.Errorf("Expected 1 task on single day, got %d", len(stacks))
	}

	// Should be in track 0
	if stacks[0].Track != 0 {
		t.Errorf("Expected single task in track 0, got track %d", stacks[0].Track)
	}
}
