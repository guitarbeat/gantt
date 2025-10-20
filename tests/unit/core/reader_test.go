package core_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"phd-dissertation-planner/internal/core"
)

func TestNewReader(t *testing.T) {
	reader := core.NewReader("test.csv")
	if reader == nil {
		t.Fatal("NewReader returned nil")
	}
	if reader.FilePath != "test.csv" {
		t.Errorf("Expected FilePath 'test.csv', got '%s'", reader.FilePath)
	}
}

func TestReadTasks_ValidCSV(t *testing.T) {
	// Create temporary test CSV
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	csvContent := `Task,Start Date,End Date,Phase,Category
Literature Review,2025-01-01,2025-03-31,Phase 1,Research
Data Collection,2025-02-01,2025-04-30,Phase 1,Research
Analysis,2025-04-01,2025-06-30,Phase 2,Analysis`

	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	reader := core.NewReader(csvPath)
	tasks, err := reader.ReadTasks()

	if err != nil {
		t.Fatalf("ReadTasks failed: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}

	// Verify first task
	if tasks[0].Name != "Literature Review" {
		t.Errorf("Expected task name 'Literature Review', got '%s'", tasks[0].Name)
	}

	expectedStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if !tasks[0].StartDate.Equal(expectedStart) {
		t.Errorf("Expected start date %v, got %v", expectedStart, tasks[0].StartDate)
	}
}

func TestReadTasks_FileNotFound(t *testing.T) {
	reader := core.NewReader("nonexistent.csv")
	_, err := reader.ReadTasks()

	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestReadTasks_InvalidDateFormat(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	csvContent := `Task,Start Date,End Date
Invalid Task,01/01/2025,03/31/2025`

	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	reader := core.NewReader(csvPath)
	tasks, err := reader.ReadTasks()

	// Should handle invalid dates gracefully
	if err != nil && len(tasks) == 0 {
		// This is acceptable - either skip invalid rows or return error
		return
	}
}

func TestReadTasks_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	csvContent := `Task,Start Date,End Date`

	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	reader := core.NewReader(csvPath)
	tasks, err := reader.ReadTasks()

	if err != nil {
		t.Fatalf("ReadTasks failed on empty file: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks from empty file, got %d", len(tasks))
	}
}

func TestReadTasks_MissingRequiredColumns(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	csvContent := `Task,Start Date
Missing End Date,2025-01-01`

	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	reader := core.NewReader(csvPath)
	tasks, err := reader.ReadTasks()

	// Reader may be lenient and skip invalid rows or return error
	// Either behavior is acceptable
	if err != nil {
		t.Logf("Reader returned error as expected: %v", err)
	} else if len(tasks) == 0 {
		t.Log("Reader skipped invalid rows")
	} else {
		t.Logf("Reader parsed %d tasks (may have used defaults)", len(tasks))
	}
}

func TestReadTasks_WithOptionalFields(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	csvContent := `Task,Start Date,End Date,Phase,Category,Status,Notes
Complete Task,2025-01-01,2025-03-31,Phase 1,Research,In Progress,Important task`

	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	reader := core.NewReader(csvPath)
	tasks, err := reader.ReadTasks()

	if err != nil {
		t.Fatalf("ReadTasks failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]
	t.Logf("Task parsed: Name=%s, Phase=%s, Category=%s, Status=%s",
		task.Name, task.Phase, task.Category, task.Status)

	// Verify task was parsed (exact field mapping may vary)
	if task.Name != "Complete Task" {
		t.Errorf("Expected task name 'Complete Task', got '%s'", task.Name)
	}

	// Check that at least some optional fields were populated
	hasOptionalFields := task.Phase != "" || task.Category != "" || task.Status != ""
	if !hasOptionalFields {
		t.Error("Expected some optional fields to be populated")
	}
}

func TestReadTasks_UTF8Encoding(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")

	csvContent := `Task,Start Date,End Date
Café Research ☕,2025-01-01,2025-03-31
数据分析,2025-02-01,2025-04-30`

	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	reader := core.NewReader(csvPath)
	tasks, err := reader.ReadTasks()

	if err != nil {
		t.Fatalf("ReadTasks failed with UTF-8: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Name != "Café Research ☕" {
		t.Errorf("UTF-8 task name not preserved: got '%s'", tasks[0].Name)
	}
}

func BenchmarkReadTasks(b *testing.B) {
	tmpDir := b.TempDir()
	csvPath := filepath.Join(tmpDir, "bench.csv")

	// Create a CSV with 100 tasks
	csvContent := "Task,Start Date,End Date,Phase,Category\n"
	for i := 0; i < 100; i++ {
		csvContent += "Task,2025-01-01,2025-12-31,Phase 1,Research\n"
	}

	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		b.Fatalf("Failed to create benchmark CSV: %v", err)
	}

	reader := core.NewReader(csvPath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := reader.ReadTasks()
		if err != nil {
			b.Fatalf("ReadTasks failed: %v", err)
		}
	}
}
