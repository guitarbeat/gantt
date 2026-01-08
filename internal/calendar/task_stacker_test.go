package calendar

import (
	"testing"
	"time"
)

func TestTaskStacker_ComputeStacks(t *testing.T) {
	// Create overlapping tasks
	// Task 1: 2023-01-01 to 2023-01-05
	// Task 2: 2023-01-03 to 2023-01-07 (Overlaps Task 1)
	// Task 3: 2023-01-02 to 2023-01-04 (Overlaps Task 1 and Task 2)

	t1 := &SpanningTask{ID: "1", Name: "T1", StartDate: date(2023, 1, 1), EndDate: date(2023, 1, 5)}
	t2 := &SpanningTask{ID: "2", Name: "T2", StartDate: date(2023, 1, 3), EndDate: date(2023, 1, 7)}
	t3 := &SpanningTask{ID: "3", Name: "T3", StartDate: date(2023, 1, 2), EndDate: date(2023, 1, 4)}

	tasks := []*SpanningTask{t1, t2, t3}
	stacker := NewTaskStacker(tasks, time.Monday)
	stacker.ComputeStacks()

	// Check 2023-01-03 (All 3 tasks active)
	day := date(2023, 1, 3)
	stacks := stacker.GetStacksForDay(day)

	if len(stacks) != 3 {
		t.Errorf("Expected 3 tasks on %v, got %d", day, len(stacks))
	}

	// Verify unique tracks
	tracks := make(map[int]bool)
	for _, s := range stacks {
		if tracks[s.Track] {
			t.Errorf("Duplicate track %d on %v", s.Track, day)
		}
		tracks[s.Track] = true
	}
}

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
