package calendar

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/core"
)

func TestApplySpanningTasksToMonth(t *testing.T) {
	// Setup
	cfg := core.DefaultConfig()
	year := 2023
	month := time.June
	wd := time.Sunday

	y := NewYear(wd, year, &cfg)
	// Find the month object in the year
	var targetMonth *Month
	for _, q := range y.Quarters {
		for _, m := range q.Months {
			if m.Month == month {
				targetMonth = m
				break
			}
		}
	}

	if targetMonth == nil {
		t.Fatal("Could not find target month")
	}

	// Create tasks
	// Task 1: Spans entire month
	task1 := SpanningTask{
		ID: "task-1",
		Name: "Task 1",
		StartDate: time.Date(year, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate: time.Date(year, 6, 30, 0, 0, 0, 0, time.UTC),
	}

	// Task 2: Starts before month, ends in middle
	task2 := SpanningTask{
		ID: "task-2",
		Name: "Task 2",
		StartDate: time.Date(year, 5, 20, 0, 0, 0, 0, time.UTC),
		EndDate: time.Date(year, 6, 15, 0, 0, 0, 0, time.UTC),
	}

	// Task 3: Starts in middle, ends after month
	task3 := SpanningTask{
		ID: "task-3",
		Name: "Task 3",
		StartDate: time.Date(year, 6, 15, 0, 0, 0, 0, time.UTC),
		EndDate: time.Date(year, 7, 10, 0, 0, 0, 0, time.UTC),
	}

	// Task 4: Completely outside (before)
	task4 := SpanningTask{
		ID: "task-4",
		Name: "Task 4",
		StartDate: time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate: time.Date(year, 2, 1, 0, 0, 0, 0, time.UTC),
	}

	// Task 5: Starts and ends within month
	task5 := SpanningTask{
		ID: "task-5",
		Name: "Task 5",
		StartDate: time.Date(year, 6, 5, 0, 0, 0, 0, time.UTC),
		EndDate: time.Date(year, 6, 6, 0, 0, 0, 0, time.UTC),
	}

	tasks := []SpanningTask{task1, task2, task3, task4, task5}

	ApplySpanningTasksToMonth(targetMonth, tasks)

	// Verification Helper
	checkDay := func(day int, expectedTaskIDs []string) {
		var d *Day
		for _, w := range targetMonth.Weeks {
			for i := range w.Days {
				if w.Days[i].Time.Month() == month && w.Days[i].Time.Day() == day {
					d = &w.Days[i]
					break
				}
			}
		}

		if d == nil {
			t.Fatalf("Day %d not found in month", day)
		}

		// Check tasks on this day
		if len(d.Tasks) != len(expectedTaskIDs) {
			t.Errorf("Day %d: Expected %d tasks, got %d", day, len(expectedTaskIDs), len(d.Tasks))
		}

		foundMap := make(map[string]bool)
		for _, task := range d.Tasks {
			foundMap[task.ID] = true
		}

		for _, id := range expectedTaskIDs {
			if !foundMap[id] {
				t.Errorf("Day %d: Expected task %s not found", day, id)
			}
		}
	}

	// June 1: Task 1, Task 2
	checkDay(1, []string{"task-1", "task-2"})

	// June 10: Task 1, Task 2
	checkDay(10, []string{"task-1", "task-2"})

	// June 15: Task 1, Task 2, Task 3
	checkDay(15, []string{"task-1", "task-2", "task-3"})

	// June 16: Task 1, Task 3
	checkDay(16, []string{"task-1", "task-3"})

	// June 5: Task 1, Task 2, Task 5
	checkDay(5, []string{"task-1", "task-2", "task-5"})

	// June 30: Task 1, Task 3
	checkDay(30, []string{"task-1", "task-3"})
}
