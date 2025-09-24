package common

import (
	"testing"
	"time"
)

func TestTask_IsOnDate(t *testing.T) {
	// Create a test task
	task := &Task{
		ID:        "test-1",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		Category:  "TEST",
		Priority:  1,
		Status:    "Planned",
	}

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{
			name:     "task start date",
			date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "task end date",
			date:     time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "middle of task",
			date:     time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "before task",
			date:     time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "after task",
			date:     time.Date(2024, 1, 21, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := task.IsOnDate(tt.date)
			if result != tt.expected {
				t.Errorf("IsOnDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTask_GetDuration(t *testing.T) {
	task := &Task{
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}

	expected := 6 // 15th to 20th inclusive = 6 days
	result := task.GetDuration()

	if result != expected {
		t.Errorf("GetDuration() = %v, want %v", result, expected)
	}
}

func TestTask_IsOverdue(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		task     Task
		expected bool
	}{
		{
			name: "overdue task",
			task: Task{
				EndDate: now.AddDate(0, 0, -1), // Yesterday
				Status:  "Planned",
			},
			expected: true,
		},
		{
			name: "completed task",
			task: Task{
				EndDate: now.AddDate(0, 0, -1), // Yesterday
				Status:  "Completed",
			},
			expected: false,
		},
		{
			name: "future task",
			task: Task{
				EndDate: now.AddDate(0, 0, 1), // Tomorrow
				Status:  "Planned",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.task.IsOverdue()
			if result != tt.expected {
				t.Errorf("IsOverdue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected string
	}{
		{
			name:     "proposal category",
			category: "PROPOSAL",
			expected: "Proposal",
		},
		{
			name:     "research category",
			category: "RESEARCH",
			expected: "Research",
		},
		{
			name:     "unknown category",
			category: "UNKNOWN",
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCategory(tt.category)
			if result.DisplayName != tt.expected {
				t.Errorf("GetCategory() = %v, want %v", result.DisplayName, tt.expected)
			}
		})
	}
}
