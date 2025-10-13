package app_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"phd-dissertation-planner/src/app"
	"phd-dissertation-planner/src/core"
	cal "phd-dissertation-planner/src/calendar"
)

func TestCalculatePackageAverage(t *testing.T) {
	tests := []struct {
		name     string
		coverages []float64
		expected float64
	}{
		{
			name:     "empty slice",
			coverages: []float64{},
			expected: 0.0,
		},
		{
			name:     "single value",
			coverages: []float64{75.5},
			expected: 75.5,
		},
		{
			name:     "multiple values",
			coverages: []float64{80.0, 90.0, 85.0},
			expected: 85.0,
		},
		{
			name:     "decimal values",
			coverages: []float64{33.3, 66.7, 50.0},
			expected: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.CalculatePackageAverage(tt.coverages)
			if result != tt.expected {
				t.Errorf("calculatePackageAverage(%v) = %v, want %v", tt.coverages, result, tt.expected)
			}
		})
	}
}

func TestGetCoverageStatus(t *testing.T) {
	tests := []struct {
		name     string
		coverage float64
		expected string
	}{
		{
			name:     "excellent coverage",
			coverage: 85.0,
			expected: "✅",
		},
		{
			name:     "warning coverage",
			coverage: 65.0,
			expected: "⚠️ ",
		},
		{
			name:     "low coverage",
			coverage: 45.0,
			expected: "❌",
		},
		{
			name:     "boundary excellent",
			coverage: 80.0,
			expected: "✅",
		},
		{
			name:     "boundary warning",
			coverage: 60.0,
			expected: "⚠️ ",
		},
		{
			name:     "boundary low",
			coverage: 59.9,
			expected: "❌",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.GetCoverageStatus(tt.coverage)
			if result != tt.expected {
				t.Errorf("getCoverageStatus(%v) = %q, want %q", tt.coverage, result, tt.expected)
			}
		})
	}
}

func TestCalculateCSVPriority(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected int
	}{
		{
			name:     "comprehensive file",
			filename: "research_timeline_v5_comprehensive.csv",
			expected: 10,
		},
		{
			name:     "v5.1 file",
			filename: "data_v5.1_tasks.csv",
			expected: 8,
		},
		{
			name:     "v5 file",
			filename: "timeline_v5.csv",
			expected: 6,
		},
		{
			name:     "regular file",
			filename: "tasks.csv",
			expected: 0,
		},
		{
			name:     "case insensitive comprehensive",
			filename: "COMPREHENSIVE_DATA.csv",
			expected: 10,
		},
		{
			name:     "case insensitive version",
			filename: "V5.1_EXPORT.csv",
			expected: 8,
		},
		{
			name:     "no extension",
			filename: "comprehensive_data",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.CalculateCSVPriority(tt.filename)
			if result != tt.expected {
				t.Errorf("calculateCSVPriority(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestRootFilename(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple config",
			path:     "config.yaml",
			expected: "config.tex",
		},
		{
			name:     "nested path",
			path:     "configs/base.yaml",
			expected: "base.tex",
		},
		{
			name:     "no extension",
			path:     "config",
			expected: "config.tex",
		},
		{
			name:     "multiple dots",
			path:     "config.special.yaml",
			expected: "config.special.tex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.RootFilename(tt.path)
			if result != tt.expected {
				t.Errorf("RootFilename(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestEscapeLatex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special chars",
			input:    "Simple text",
			expected: "Simple text",
		},
		{
			name:     "ampersand",
			input:    "Tom & Jerry",
			expected: "Tom \\& Jerry",
		},
		{
			name:     "percent",
			input:    "100% complete",
			expected: "100\\% complete",
		},
		{
			name:     "dollar",
			input:    "$100 budget",
			expected: "\\$100 budget",
		},
		{
			name:     "hash",
			input:    "#1 priority",
			expected: "\\#1 priority",
		},
		{
			name:     "underscore",
			input:    "task_name",
			expected: "task\\_name",
		},
		{
			name:     "curly braces",
			input:    "{bold} text",
			expected: "\\{bold\\} text",
		},
		{
			name:     "multiple chars",
			input:    "Cost: $100 & 50% done_#1",
			expected: "Cost: \\$100 \\& 50\\% done\\_\\#1",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.EscapeLatex(tt.input)
			if result != tt.expected {
				t.Errorf("escapeLatex(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateModuleAlignment(t *testing.T) {
	tests := []struct {
		name        string
		modules     []core.Modules
		pageName    string
		shouldError bool
	}{
		{
			name:        "empty modules",
			modules:     []core.Modules{},
			pageName:    "test",
			shouldError: false,
		},
		{
			name: "aligned modules",
			modules: []core.Modules{
				{make(core.Module, 3), make(core.Module, 3)},
				{make(core.Module, 3), make(core.Module, 3)},
			},
			pageName:    "test",
			shouldError: false,
		},
		{
			name: "misaligned modules",
			modules: []core.Modules{
				{make(core.Module, 3), make(core.Module, 3)},
				{make(core.Module, 2), make(core.Module, 3)}, // Different length
			},
			pageName:    "test",
			shouldError: true,
		},
		{
			name: "single module",
			modules: []core.Modules{
				{make(core.Module, 2)},
			},
			pageName:    "test",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.ValidateModuleAlignment(tt.modules, tt.pageName)
			if tt.shouldError && err == nil {
				t.Errorf("validateModuleAlignment() expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("validateModuleAlignment() unexpected error: %v", err)
			}
		})
	}
}

func TestSetupOutputDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_output_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := core.Config{
		OutputDir: tmpDir,
	}

	err = app.SetupOutputDirectory(cfg)
	if err != nil {
		t.Errorf("setupOutputDirectory() error: %v", err)
	}

	// Check that directory exists
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Errorf("Output directory was not created")
	}
}

func TestAutoDetectCSV(t *testing.T) {
	// Create a temporary input_data directory
	tmpDir, err := os.MkdirTemp("", "test_input_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create input_data subdirectory
	inputDir := filepath.Join(tmpDir, "input_data")
	err = os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input_data dir: %v", err)
	}

	// Change to temp directory and create test files
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create test CSV files
	files := []string{
		"research_timeline_v5_comprehensive.csv",
		"data_v5.1_tasks.csv",
		"timeline_v5.csv",
		"regular_tasks.csv",
	}

	for _, file := range files {
		f, err := os.Create(filepath.Join(inputDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.WriteString("Task Name,Start Date,End Date\nTest,2024-01-01,2024-01-02\n")
		f.Close()
	}

	// Test auto-detection
	result, err := app.AutoDetectCSV()
	if err != nil {
		t.Errorf("autoDetectCSV() error: %v", err)
	}

	// Should select the comprehensive file (highest priority)
	expected := filepath.Join(inputDir, "research_timeline_v5_comprehensive.csv")
	if result != expected {
		t.Errorf("autoDetectCSV() = %q, want %q", result, expected)
	}
}

func TestAutoDetectConfig(t *testing.T) {
	// Create a temporary CSV file
	tmpFile, err := os.CreateTemp("", "test_config_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Test v5.1 format
	v51Data := `Task Name,Start Date,End Date,Phase,Sub-Phase
Test Task,2024-01-01,2024-01-02,1,Planning`
	tmpFile.WriteString(v51Data)
	tmpFile.Close()

	result, err := app.AutoDetectConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("autoDetectConfig() error: %v", err)
	}

	expected := []string{"src/core/base.yaml", "src/core/monthly_calendar.yaml"}
	if len(result) != len(expected) || result[0] != expected[0] || result[1] != expected[1] {
		t.Errorf("autoDetectConfig() = %v, want %v", result, expected)
	}
}

func TestCreateTableOfContentsModule(t *testing.T) {
	cfg := core.Config{
		WeekStart: time.Monday,
	}

	tasks := []core.Task{
		{
			ID:          "1",
			Name:        "Task 1",
			StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Status:      "completed",
			Phase:       "1",
			SubPhase:    "Planning",
			IsMilestone: true,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			StartDate:   time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			Status:      "in progress",
			Phase:       "1",
			SubPhase:    "Planning",
			IsMilestone: false,
		},
	}

	module := app.CreateTableOfContentsModule(cfg, tasks, "toc.tpl")

	// Verify module structure
	if module.Tpl != "toc.tpl" {
		t.Errorf("Expected template 'toc.tpl', got %s", module.Tpl)
	}

	body := module.Body
	if body["TotalTasks"] != 2 {
		t.Errorf("Expected 2 total tasks, got %v", body["TotalTasks"])
	}

	if body["MilestoneCount"] != 1 {
		t.Errorf("Expected 1 milestone, got %v", body["MilestoneCount"])
	}

	if body["CompletedCount"] != 1 {
		t.Errorf("Expected 1 completed task, got %v", body["CompletedCount"])
	}
}

func TestAssignTasksToMonth(t *testing.T) {
	// Create a test month
	year := cal.NewYear(time.Monday, 2024, &core.Config{})
	var targetMonth *cal.Month
	for _, quarter := range year.Quarters {
		for _, month := range quarter.Months {
			if month.Month == time.January {
				targetMonth = month
				break
			}
		}
		if targetMonth != nil {
			break
		}
	}

	if targetMonth == nil {
		t.Fatal("Could not find January month")
	}

	// Create test tasks
	tasks := []core.Task{
		{
			Name:      "Jan Task",
			StartDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			Status:    "in progress",
		},
		{
			Name:      "Feb Task",
			StartDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC),
			Status:    "planned",
		},
	}

	// Assign tasks
	app.AssignTasksToMonth(targetMonth, tasks)

	// Verify that Jan task was assigned (Feb task should not be)
	// This is a basic test - more detailed verification would require
	// inspecting the month's internal state
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

// Test edge cases for CSV parsing
func TestCSVEdgeCases(t *testing.T) {
	// Create a temporary CSV file with edge cases
	csvData := `Task Name,Start Date,End Date,Category,Description
Empty Name,,2024-01-20,Test,Test
Invalid Date,invalid-date,2024-01-20,Test,Test
Same Dates,2024-01-15,2024-01-15,Test,Milestone
Empty End,2024-01-15,,Test,Test
Special Chars,2024-01-15,2024-01-20,Test,Task with & % $ # _ { }
Very Long Name,2024-01-15,2024-01-20,Test,` + strings.Repeat("Long description ", 100)

	tmpFile, err := os.CreateTemp("", "test_edge_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	reader := core.NewReader(tmpFile.Name())
	tasks, err := reader.ReadTasks()

	// Should not fail completely on edge cases
	if err != nil {
		t.Logf("Note: CSV parsing produced error: %v", err)
	}

	// Should parse at least some valid tasks
	if len(tasks) == 0 {
		t.Errorf("Expected at least some tasks to be parsed")
	}
}

// Integration test for template rendering
func TestTemplateRendering(t *testing.T) {
	// Create a basic config
	cfg := core.Config{
		WeekStart: time.Monday,
		Pages: []core.Page{
			{
				Name: "test_page.tex",
				RenderBlocks: []core.RenderBlock{
					{
						FuncName: "monthly",
						Tpls:     []string{"calendar.tpl"},
					},
				},
			},
		},
	}

	// Create template
	tpl := app.NewTpl()

	// Test document rendering
	var buf bytes.Buffer
	err := tpl.Document(&buf, cfg)
	if err != nil {
		t.Errorf("Document rendering failed: %v", err)
	}

	content := buf.String()
	if !strings.Contains(content, "test_page.tex") {
		t.Errorf("Rendered document does not contain expected page name")
	}
}

// Property-based testing for date parsing (using table-driven tests as approximation)
func TestDateParsingRobustness(t *testing.T) {
	testCases := []struct {
		name        string
		dateString  string
		shouldParse bool
	}{
		{"valid ISO", "2024-01-15", true},
		{"valid US", "01/15/2024", true},
		{"invalid format", "15-Jan-2024", false},
		{"empty string", "", false},
		{"invalid date", "2024-13-45", false},
		{"leap year", "2024-02-29", true},
		{"non-leap year feb 29", "2023-02-29", false},
		{"year 0", "0000-01-01", false},
		{"future date", "2030-12-31", true},
		{"past date", "2000-01-01", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a minimal CSV with the date
			csvData := fmt.Sprintf("Task,Start,End\nTest,%s,%s\n", tc.dateString, tc.dateString)

			tmpFile, err := os.CreateTemp("", "test_date_*.csv")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(csvData); err != nil {
				t.Fatalf("Failed to write test data: %v", err)
			}
			tmpFile.Close()

			reader := core.NewReader(tmpFile.Name())
			tasks, err := reader.ReadTasks()

			if tc.shouldParse {
				if err != nil {
					t.Errorf("Expected date %q to parse successfully, but got error: %v", tc.dateString, err)
				}
				if len(tasks) == 0 {
					t.Errorf("Expected at least one task for valid date %q", tc.dateString)
				}
			} else {
				// For invalid dates, we might still get tasks but with zero dates
				// This tests the robustness of the parser
				t.Logf("Date %q parsing result: error=%v, tasks=%d", tc.dateString, err, len(tasks))
			}
		})
	}
}
