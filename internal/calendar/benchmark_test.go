package calendar

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkAssignTaskTracks(b *testing.B) {
	// Create a day with many overlapping tasks
	count := 100 // Simulate 100 overlapping tasks
	tasks := make([]*SpanningTask, count)

	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < count; i++ {
		tasks[i] = &SpanningTask{
			ID:        fmt.Sprintf("T%d", i),
			StartDate: baseTime,
			EndDate:   baseTime.AddDate(0, 0, 10), // All span 10 days
		}
	}

	day := Day{
		Time:  baseTime,
		Tasks: tasks, // used by findLowestAvailableTrackForTask in original code
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		day.assignTaskTracks(tasks)
	}
}
