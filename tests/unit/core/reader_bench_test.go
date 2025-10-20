package core_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"phd-dissertation-planner/internal/core"
)

// Benchmark data for consistent benchmarking
var benchmarkCSVContent = `Task Name,Start Date,End Date,Phase,Category,Status,Assignee,Description
Literature Review,2024-01-01,2024-03-31,1,Literature,In Progress,John Doe,Conduct comprehensive literature review
Data Collection,2024-04-01,2024-06-30,2,Methodology,Planned,Jane Smith,Collect experimental data
Analysis,2024-07-01,2024-09-30,3,Analysis,Planned,Bob Wilson,Analyze collected data
Writing Chapter 1,2024-10-01,2024-11-30,4,Writing,Planned,Alice Brown,Write introduction chapter
Writing Chapter 2,2024-12-01,2025-01-31,4,Writing,Planned,Charlie Davis,Write methodology chapter
Writing Chapter 3,2025-02-01,2025-03-31,4,Writing,Planned,Diana Evans,Write results chapter
Final Revisions,2025-04-01,2025-05-31,4,Revision,Planned,Frank Garcia,Final edits and revisions
Defense Preparation,2025-06-01,2025-07-31,4,Defense,Planned,Helen Harris,Prepare for dissertation defense
`

var benchmarkLargeCSVContent string

func init() {
	// Generate larger dataset for benchmarking
	benchmarkLargeCSVContent = benchmarkCSVContent
	for i := 0; i < 100; i++ {
		benchmarkLargeCSVContent += `Task Name,Start Date,End Date,Phase,Category,Status,Assignee,Description
Large Dataset Task ` + string(rune(i)) + `,2024-01-01,2024-03-31,1,Literature,In Progress,Test User,Test description for benchmarking
`
	}
}

// createBenchmarkCSV creates a temporary CSV file for benchmarking
func createBenchmarkCSV(b *testing.B, content string) string {
	dir := b.TempDir()
	file := filepath.Join(dir, "benchmark.csv")
	err := os.WriteFile(file, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}
	return file
}

// BenchmarkCSVReading benchmarks CSV reading performance
func BenchmarkCSVReading(b *testing.B) {
	file := createBenchmarkCSV(b, benchmarkCSVContent)
	defer os.Remove(file)

	reader := core.NewReader(file)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tasks, err := reader.ReadTasks()
		if err != nil {
			b.Fatal(err)
		}
		_ = tasks // Prevent optimization
	}
}

// BenchmarkCSVReadingLarge benchmarks CSV reading with larger dataset
func BenchmarkCSVReadingLarge(b *testing.B) {
	file := createBenchmarkCSV(b, benchmarkLargeCSVContent)
	defer os.Remove(file)

	reader := core.NewReader(file)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tasks, err := reader.ReadTasks()
		if err != nil {
			b.Fatal(err)
		}
		_ = tasks // Prevent optimization
	}
}

// BenchmarkDateParsing benchmarks date parsing performance
func BenchmarkDateParsing(b *testing.B) {
	dateStrings := []string{
		"2024-01-01",
		"2024-12-31",
		"2025-06-15",
		"2023-03-10",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, dateStr := range dateStrings {
			_, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// BenchmarkTaskValidation benchmarks task validation performance
func BenchmarkTaskValidation(b *testing.B) {
	tasks := []*core.Task{
		{
			Name:      "Test Task 1",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			Phase:     "1",
			Category:  "Test",
			Status:    "In Progress",
		},
		{
			Name:      "Test Task 2",
			StartDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
			Phase:     "2",
			Category:  "Test",
			Status:    "Completed",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, task := range tasks {
			if task.Name == "" {
				b.Fatal("Task name should not be empty")
			}
			if task.StartDate.After(task.EndDate) {
				b.Fatal("Start date should not be after end date")
			}
			if task.Phase == "" {
				b.Fatal("Phase should not be empty")
			}
		}
	}
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tasks := make([]*core.Task, 0, 100)
		for j := 0; j < 100; j++ {
			task := &core.Task{
				Name:      "Benchmark Task",
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				Phase:     "1",
				Category:  "Benchmark",
				Status:    "In Progress",
			}
			tasks = append(tasks, task)
		}
		_ = tasks // Prevent optimization
	}
}
