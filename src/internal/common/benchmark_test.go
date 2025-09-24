package common

import (
	"testing"
	"time"
)

func BenchmarkTask_IsOnDate(b *testing.B) {
	task := &Task{
		ID:        "benchmark-task",
		Name:      "Benchmark Task",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		Category:  "BENCHMARK",
		Priority:  1,
		Status:    "Planned",
	}

	testDate := time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task.IsOnDate(testDate)
	}
}

func BenchmarkTask_GetDuration(b *testing.B) {
	task := &Task{
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task.GetDuration()
	}
}

func BenchmarkGetCategory(b *testing.B) {
	categories := []string{"PROPOSAL", "RESEARCH", "IMAGING", "ADMIN", "DISSERTATION", "UNKNOWN"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		category := categories[i%len(categories)]
		GetCategory(category)
	}
}

func BenchmarkTask_IsOverdue(b *testing.B) {
	now := time.Now()
	task := &Task{
		EndDate: now.AddDate(0, 0, -1), // Yesterday
		Status:  "Planned",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task.IsOverdue()
	}
}

func BenchmarkTask_GetProgressPercentage(b *testing.B) {
	now := time.Now()
	task := &Task{
		StartDate: now.AddDate(0, 0, -5), // 5 days ago
		EndDate:   now.AddDate(0, 0, 5),  // 5 days from now
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task.GetProgressPercentage()
	}
}

// Benchmark memory allocation for task creation
func BenchmarkTaskCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &Task{
			ID:           "task-" + string(rune(i)),
			Name:         "Test Task",
			StartDate:    time.Now(),
			EndDate:      time.Now().AddDate(0, 0, 7),
			Category:     "TEST",
			Description:  "Test description",
			Priority:     1,
			Status:       "Planned",
			Assignee:     "test-user",
			ParentID:     "",
			Dependencies: []string{},
			IsMilestone:  false,
		}
	}
}

// Benchmark string operations
func BenchmarkTask_String(b *testing.B) {
	task := &Task{
		ID:        "test-task",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		Category:  "TEST",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = task.String()
	}
}