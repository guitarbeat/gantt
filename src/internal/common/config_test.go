package common

import (
	"os"
	"testing"
	"time"
)

func TestDefaultOutputDir(t *testing.T) {
	cfg, err := NewConfig() // no files, should fall back to defaults
	if err != nil {
		t.Fatalf("NewConfig error: %v", err)
	}
	if cfg.OutputDir != "build" {
		t.Fatalf("OutputDir should default to 'build', got %q", cfg.OutputDir)
	}
}

func TestNewConfigWithDefaults(t *testing.T) {
	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig error: %v", err)
	}

	// Test default values
	if cfg.Year != time.Now().Year() {
		t.Errorf("Expected Year to be current year %d, got %d", time.Now().Year(), cfg.Year)
	}
	if cfg.OutputDir != "build" {
		t.Errorf("Expected OutputDir to be 'build', got %q", cfg.OutputDir)
	}
	if cfg.WeekStart != 0 { // Sunday
		t.Errorf("Expected WeekStart to be 0 (Sunday), got %d", cfg.WeekStart)
	}
}

func TestNewConfigWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("PLANNER_YEAR", "2025")
	os.Setenv("PLANNER_OUTPUT_DIR", "/tmp/test")
	defer func() {
		os.Unsetenv("PLANNER_YEAR")
		os.Unsetenv("PLANNER_OUTPUT_DIR")
	}()

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig error: %v", err)
	}

	if cfg.Year != 2025 {
		t.Errorf("Expected Year to be 2025, got %d", cfg.Year)
	}
	if cfg.OutputDir != "/tmp/test" {
		t.Errorf("Expected OutputDir to be '/tmp/test', got %q", cfg.OutputDir)
	}
}

func TestNewConfigWithYAMLFile(t *testing.T) {
	// Create a temporary YAML file
	yamlContent := `
year: 2024
outputdir: "test_output"
dotted: true
calafterschedule: true
cleartoprightcorner: true
ampmtime: true
addlasthalfhour: true
debug:
  showframe: true
  showlinks: true
layout:
  paper:
    width: "8.5in"
    height: "11in"
    margin:
      top: "1in"
      bottom: "1in"
      left: "1in"
      right: "1in"
  numbers:
    arraystretch: 1.5
  colors:
    gray: "#808080"
    lightgray: "#C0C0C0"
pages:
  - name: "monthly"
    renderblocks:
      - funcname: "monthly"
        tpls: ["monthly.tpl"]
`

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(yamlContent); err != nil {
		t.Fatalf("Failed to write YAML content: %v", err)
	}
	tmpFile.Close()

	cfg, err := NewConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewConfig error: %v", err)
	}

	// Test YAML values
	if cfg.Year != 2024 {
		t.Errorf("Expected Year to be 2024, got %d", cfg.Year)
	}
	if cfg.OutputDir != "test_output" {
		t.Errorf("Expected OutputDir to be 'test_output', got %q", cfg.OutputDir)
	}
	if !cfg.Dotted {
		t.Error("Expected Dotted to be true")
	}
	if !cfg.CalAfterSchedule {
		t.Error("Expected CalAfterSchedule to be true")
	}
	if !cfg.ClearTopRightCorner {
		t.Error("Expected ClearTopRightCorner to be true")
	}
	if !cfg.AMPMTime {
		t.Error("Expected AMPMTime to be true")
	}
	if !cfg.AddLastHalfHour {
		t.Error("Expected AddLastHalfHour to be true")
	}
	if !cfg.Debug.ShowFrame {
		t.Error("Expected Debug.ShowFrame to be true")
	}
	if !cfg.Debug.ShowLinks {
		t.Error("Expected Debug.ShowLinks to be true")
	}
	if cfg.Layout.Paper.Width != "8.5in" {
		t.Errorf("Expected Paper.Width to be '8.5in', got %q", cfg.Layout.Paper.Width)
	}
	if cfg.Layout.Paper.Height != "11in" {
		t.Errorf("Expected Paper.Height to be '11in', got %q", cfg.Layout.Paper.Height)
	}
	if cfg.Layout.Paper.Margin.Top != "1in" {
		t.Errorf("Expected Paper.Margin.Top to be '1in', got %q", cfg.Layout.Paper.Margin.Top)
	}
	if cfg.Layout.Numbers.ArrayStretch != 1.5 {
		t.Errorf("Expected Numbers.ArrayStretch to be 1.5, got %f", cfg.Layout.Numbers.ArrayStretch)
	}
	if cfg.Layout.Colors.Gray != "#808080" {
		t.Errorf("Expected Colors.Gray to be '#808080', got %q", cfg.Layout.Colors.Gray)
	}
	if len(cfg.Pages) != 1 {
		t.Errorf("Expected 1 page, got %d", len(cfg.Pages))
	}
	if cfg.Pages[0].Name != "monthly" {
		t.Errorf("Expected page name to be 'monthly', got %q", cfg.Pages[0].Name)
	}
}

func TestNewConfigWithMultipleYAMLFiles(t *testing.T) {
	// Create base config
	baseYaml := `
year: 2024
outputdir: "base_output"
dotted: true
`

	// Create override config
	overrideYaml := `
outputdir: "override_output"
calafterschedule: true
debug:
  showframe: true
`

	// Create temporary files
	baseFile, err := os.CreateTemp("", "base_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create base temp file: %v", err)
	}
	defer os.Remove(baseFile.Name())

	overrideFile, err := os.CreateTemp("", "override_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create override temp file: %v", err)
	}
	defer os.Remove(overrideFile.Name())

	// Write content
	baseFile.WriteString(baseYaml)
	baseFile.Close()
	overrideFile.WriteString(overrideYaml)
	overrideFile.Close()

	cfg, err := NewConfig(baseFile.Name(), overrideFile.Name())
	if err != nil {
		t.Fatalf("NewConfig error: %v", err)
	}

	// Test that override values take precedence
	if cfg.Year != 2024 {
		t.Errorf("Expected Year to be 2024, got %d", cfg.Year)
	}
	if cfg.OutputDir != "override_output" {
		t.Errorf("Expected OutputDir to be 'override_output', got %q", cfg.OutputDir)
	}
	if !cfg.Dotted {
		t.Error("Expected Dotted to be true")
	}
	if !cfg.CalAfterSchedule {
		t.Error("Expected CalAfterSchedule to be true")
	}
	if !cfg.Debug.ShowFrame {
		t.Error("Expected Debug.ShowFrame to be true")
	}
}

func TestNewConfigWithInvalidYAML(t *testing.T) {
	// Create invalid YAML file
	invalidYaml := `
year: 2024
outputDir: "test"
invalid: [unclosed array
`

	tmpFile, err := os.CreateTemp("", "invalid_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(invalidYaml)
	tmpFile.Close()

	_, err = NewConfig(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestNewConfigWithNonExistentFile(t *testing.T) {
	// * Now we handle missing files gracefully, so this should succeed
	_, err := NewConfig("non_existent_file.yaml")
	if err != nil {
		t.Errorf("Expected no error for non-existent file (now handled gracefully), got %v", err)
	}
}

func TestGetYears(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		expected []int
	}{
		{
			name: "Single year",
			cfg: Config{
				Year:      2024,
				StartYear: 0,
				EndYear:   0,
			},
			expected: []int{2024},
		},
		{
			name: "Year range",
			cfg: Config{
				Year:      2024,
				StartYear: 2023,
				EndYear:   2025,
			},
			expected: []int{2023, 2024, 2025},
		},
		{
			name: "Same start and end year",
			cfg: Config{
				Year:      2024,
				StartYear: 2024,
				EndYear:   2024,
			},
			expected: []int{2024},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			years := tt.cfg.GetYears()
			if len(years) != len(tt.expected) {
				t.Errorf("Expected %d years, got %d", len(tt.expected), len(years))
				return
			}
			for i, year := range years {
				if year != tt.expected[i] {
					t.Errorf("Expected year %d at index %d, got %d", tt.expected[i], i, year)
				}
			}
		})
	}
}

func TestFilterUniqueModules(t *testing.T) {
	modules := []Module{
		{Cfg: Config{Year: 2024}, Tpl: "template1.tpl", Body: "body1"},
		{Cfg: Config{Year: 2024}, Tpl: "template2.tpl", Body: "body2"},
		{Cfg: Config{Year: 2025}, Tpl: "template1.tpl", Body: "body3"}, // Duplicate template
		{Cfg: Config{Year: 2024}, Tpl: "template3.tpl", Body: "body4"},
	}

	filtered := FilterUniqueModules(modules)

	if len(filtered) != 3 {
		t.Errorf("Expected 3 unique modules, got %d", len(filtered))
	}

	// Check that template1.tpl appears only once (first occurrence)
	template1Count := 0
	for _, module := range filtered {
		if module.Tpl == "template1.tpl" {
			template1Count++
		}
	}
	if template1Count != 1 {
		t.Errorf("Expected template1.tpl to appear once, got %d times", template1Count)
	}

	// Check that all expected templates are present
	expectedTemplates := map[string]bool{
		"template1.tpl": true,
		"template2.tpl": true,
		"template3.tpl": true,
	}
	for _, module := range filtered {
		if !expectedTemplates[module.Tpl] {
			t.Errorf("Unexpected template: %s", module.Tpl)
		}
	}
}

func TestComposerMap(t *testing.T) {
	// Test that ComposerMap is initialized
	if ComposerMap == nil {
		t.Error("Expected ComposerMap to be initialized")
	}

	// Test that we can add a composer
	testComposer := func(cfg Config, tpls []string) (Modules, error) {
		return []Module{{Cfg: cfg, Tpl: "test.tpl", Body: "test"}}, nil
	}

	ComposerMap["test"] = testComposer

	if ComposerMap["test"] == nil {
		t.Error("Expected test composer to be added to ComposerMap")
	}
}

func TestConfigStructFields(t *testing.T) {
	// Test that all struct fields have proper types and tags
	cfg := Config{}

	// Test that fields exist and have expected types
	_ = cfg.Debug
	_ = cfg.Year
	_ = cfg.WeekStart
	_ = cfg.Dotted
	_ = cfg.CalAfterSchedule
	_ = cfg.ClearTopRightCorner
	_ = cfg.AMPMTime
	_ = cfg.AddLastHalfHour
	_ = cfg.CSVFilePath
	_ = cfg.StartYear
	_ = cfg.EndYear
	_ = cfg.MonthsWithTasks
	_ = cfg.Pages
	_ = cfg.Layout
	_ = cfg.OutputDir

	// Test nested structs
	_ = cfg.Debug.ShowFrame
	_ = cfg.Debug.ShowLinks
	_ = cfg.Layout.Paper
	_ = cfg.Layout.Numbers
	_ = cfg.Layout.Colors
	_ = cfg.Layout.Paper.Width
	_ = cfg.Layout.Paper.Height
	_ = cfg.Layout.Paper.Margin
	_ = cfg.Layout.Paper.Margin.Top
	_ = cfg.Layout.Paper.Margin.Bottom
	_ = cfg.Layout.Paper.Margin.Left
	_ = cfg.Layout.Paper.Margin.Right
	_ = cfg.Layout.Numbers.ArrayStretch
	_ = cfg.Layout.Colors.Gray
	_ = cfg.Layout.Colors.LightGray
}

func TestConfigWithEmptyValues(t *testing.T) {
	// Test config with empty/zero values
	cfg := Config{}

	// Test default behavior
	if cfg.Year != 0 {
		t.Errorf("Expected Year to be 0, got %d", cfg.Year)
	}
	if cfg.OutputDir != "" {
		t.Errorf("Expected OutputDir to be empty, got %q", cfg.OutputDir)
	}
	if cfg.WeekStart != 0 {
		t.Errorf("Expected WeekStart to be 0, got %d", cfg.WeekStart)
	}
	if cfg.Dotted {
		t.Error("Expected Dotted to be false")
	}
	if cfg.CalAfterSchedule {
		t.Error("Expected CalAfterSchedule to be false")
	}
	if cfg.ClearTopRightCorner {
		t.Error("Expected ClearTopRightCorner to be false")
	}
	if cfg.AMPMTime {
		t.Error("Expected AMPMTime to be false")
	}
	if cfg.AddLastHalfHour {
		t.Error("Expected AddLastHalfHour to be false")
	}
	if cfg.Debug.ShowFrame {
		t.Error("Expected Debug.ShowFrame to be false")
	}
	if cfg.Debug.ShowLinks {
		t.Error("Expected Debug.ShowLinks to be false")
	}
}

func TestConfigWithCSVFile(t *testing.T) {
	// Create a test CSV file with correct format
	csvContent := `Task ID,Task Name,Description,Category,Start Date,Due Date
1,Test Task 1,Description 1,PROPOSAL,2024-01-15,2024-01-20
2,Test Task 2,Description 2,PROPOSAL,2024-02-10,2024-02-15
3,Test Task 3,Description 3,PROPOSAL,2025-03-05,2025-03-10
`

	tmpFile, err := os.CreateTemp("", "test_tasks_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	// Set environment variable for CSV file
	os.Setenv("PLANNER_CSV_FILE", tmpFile.Name())
	defer os.Unsetenv("PLANNER_CSV_FILE")

	cfg, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig error: %v", err)
	}

	// Test that CSV file path is set
	if cfg.CSVFilePath != tmpFile.Name() {
		t.Errorf("Expected CSVFilePath to be %s, got %s", tmpFile.Name(), cfg.CSVFilePath)
	}

	// Test that date range is set from CSV
	if cfg.StartYear != 2024 {
		t.Errorf("Expected StartYear to be 2024, got %d", cfg.StartYear)
	}
	if cfg.EndYear != 2025 {
		t.Errorf("Expected EndYear to be 2025, got %d", cfg.EndYear)
	}

	// Test that months with tasks are populated
	if len(cfg.MonthsWithTasks) == 0 {
		t.Error("Expected MonthsWithTasks to be populated")
	}
}
