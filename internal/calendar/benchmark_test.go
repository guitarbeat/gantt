package calendar

import (
	"phd-dissertation-planner/internal/core"
	"testing"
	"time"
)

func BenchmarkRenderMonthWithTasks(b *testing.B) {
	// Setup
	cfg := &core.Config{}
	year := &Year{Number: 2024}
	qrtr := &Quarter{Number: 1, Year: year}
	month := NewMonth(time.Monday, year, qrtr, time.January, cfg)

	// Create many spanning tasks
	tasks := make([]SpanningTask, 100)
	for i := 0; i < 100; i++ {
		tasks[i] = SpanningTask{
			ID:          "task" + string(rune(i)),
			Name:        "Task " + string(rune(i)),
			Category:    "Research", // Same category to test cache hit
			StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i%10),
			EndDate:     time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i%10+5),
		}
	}

	ApplySpanningTasksToMonth(month, tasks)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 1. Benchmark Color Generation (which was inefficient)
		month.GetTaskColors()

		// 2. Benchmark Rendering (which triggers findActiveTasks -> assignTaskTracks -> date normalization)
		for _, week := range month.Weeks {
			for _, day := range week.Days {
				day.Day(nil, true)
			}
		}
	}
}
