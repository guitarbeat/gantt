package scheduler

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/common"
)

func TestTasksForDay_NoTasks(t *testing.T) {
	day := Day{Time: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)}
	if got := day.TasksForDay(); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestTasksForDay_WithTasksAndMilestone(t *testing.T) {
	day := Day{Time: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)}
	day.Tasks = []Task{
		{Name: "Write draft", Description: "normal"},
		{Name: "Submit", Description: "MILESTONE: submission deadline"},
	}

	got := day.TasksForDay()
	if got == "" {
		t.Fatal("expected non-empty tasks rendering")
	}
	// Ensure milestone star appears
	if !contains(got, "★") {
		t.Fatalf("expected milestone star prefix, got %q", got)
	}
}

func TestCalculateTaskSpanColumns(t *testing.T) {
	day := Day{Time: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)} // Monday
	start := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 9, 3, 0, 0, 0, 0, time.UTC)           // Wed
	cols := day.calculateTaskSpanColumns(start, end)
	if cols != 3 {
		t.Fatalf("expected 3 columns, got %d", cols)
	}

	// Spans to next row should be clamped to remaining columns in row
	end2 := time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC)
	cols2 := day.calculateTaskSpanColumns(start, end2)
	if cols2 != 7 { // Monday through Sunday (full row)
		t.Fatalf("expected 7 columns (row clamp), got %d", cols2)
	}
}

func TestCreateSpanningTask_UsesCategoryColor(t *testing.T) {
	base := common.Task{
		ID:          "1",
		Name:        "Experiment",
		Category:    "IMAGING",
		Description: "",
		Priority:    3,
		Status:      "Planned",
		Assignee:    "A",
	}
	st := CreateSpanningTask(base, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC))
	if st.Category != base.Category {
		t.Fatalf("expected Category to copy from base, got %q", st.Category)
	}
	if st.Color == "" {
		t.Fatal("expected color to be set based on category")
	}
}

// contains is a tiny helper to avoid importing strings just for one call
func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}


