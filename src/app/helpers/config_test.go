package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAutoDetectCSV(t *testing.T) {
	// Create a temporary input_data directory
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input_data")
	err := os.MkdirAll(inputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create input_data directory: %v", err)
	}

	tests := []struct {
		name        string
		files       []string
		expected    string
		expectError bool
	}{
		{
			name:        "no files",
			files:       []string{},
			expected:    "",
			expectError: true,
		},
		{
			name:        "single CSV file",
			files:       []string{"data.csv"},
			expected:    "data.csv",
			expectError: false,
		},
		{
			name:        "comprehensive file priority",
			files:       []string{"data.csv", "research_timeline_v5_comprehensive.csv", "other.csv"},
			expected:    "research_timeline_v5_comprehensive.csv",
			expectError: false,
		},
		{
			name:        "version priority",
			files:       []string{"data.csv", "research_timeline_v5.1.csv", "research_timeline_v5.csv"},
			expected:    "research_timeline_v5.1.csv",
			expectError: false,
		},
		{
			name:        "non-CSV files ignored",
			files:       []string{"data.txt", "config.yaml", "data.csv"},
			expected:    "data.csv",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up previous test files
			files, _ := os.ReadDir(inputDir)
			for _, file := range files {
				os.Remove(filepath.Join(inputDir, file.Name()))
			}

			// Create test files
			for _, filename := range tt.files {
				filePath := filepath.Join(inputDir, filename)
				err := os.WriteFile(filePath, []byte("test data"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", filename, err)
				}
			}

			// Change to temp directory for test
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			os.Chdir(tempDir)

			result, err := AutoDetectCSV()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// The function returns a relative path, so we need to compare just the filename
			if filepath.Base(result) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, filepath.Base(result))
			}
		})
	}
}

func TestAutoDetectConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		filename     string
		content      string
		expected     *ConfigDetectionResult
		expectError  bool
	}{
		{
			name:     "v5.1 filename detection",
			filename: "research_timeline_v5.1_comprehensive.csv",
			content:  "Name,StartDate,EndDate\n",
			expected: &ConfigDetectionResult{
				BaseConfig:        "src/core/base.yaml",
				AdditionalConfigs: []string{"src/core/monthly_calendar.yaml"},
				Reason:           "v5.1 format detected from filename",
			},
			expectError: false,
		},
		{
			name:     "v5 filename detection",
			filename: "research_timeline_v5.csv",
			content:  "Name,StartDate,EndDate\n",
			expected: &ConfigDetectionResult{
				BaseConfig:        "src/core/base.yaml",
				AdditionalConfigs: []string{"src/core/calendar.yaml"},
				Reason:           "v5 format detected from filename",
			},
			expectError: false,
		},
		{
			name:     "content-based detection",
			filename: "data.csv",
			content:  "Name,Phase,Sub-Phase,StartDate,EndDate\n",
			expected: &ConfigDetectionResult{
				BaseConfig:        "src/core/base.yaml",
				AdditionalConfigs: []string{"src/core/monthly_calendar.yaml"},
				Reason:           "v5.1 format detected from content (phase/sub-phase columns)",
			},
			expectError: false,
		},
		{
			name:     "default configuration",
			filename: "data.csv",
			content:  "Name,StartDate,EndDate\n",
			expected: &ConfigDetectionResult{
				BaseConfig:        "src/core/base.yaml",
				AdditionalConfigs: []string{},
				Reason:           "using default configuration",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test CSV file
			csvPath := filepath.Join(tempDir, tt.filename)
			err := os.WriteFile(csvPath, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test CSV file: %v", err)
			}

			result, err := AutoDetectConfig(csvPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.BaseConfig != tt.expected.BaseConfig {
				t.Errorf("Expected base config %s, got %s", tt.expected.BaseConfig, result.BaseConfig)
			}

			if len(result.AdditionalConfigs) != len(tt.expected.AdditionalConfigs) {
				t.Errorf("Expected %d additional configs, got %d", 
					len(tt.expected.AdditionalConfigs), len(result.AdditionalConfigs))
			}

			for i, expected := range tt.expected.AdditionalConfigs {
				if i < len(result.AdditionalConfigs) && result.AdditionalConfigs[i] != expected {
					t.Errorf("Expected additional config %s, got %s", expected, result.AdditionalConfigs[i])
				}
			}

			if result.Reason != tt.expected.Reason {
				t.Errorf("Expected reason %s, got %s", tt.expected.Reason, result.Reason)
			}
		})
	}
}

func TestConfigDetectionResult_GetConfigPaths(t *testing.T) {
	tests := []struct {
		name     string
		result   *ConfigDetectionResult
		expected []string
	}{
		{
			name: "with additional configs",
			result: &ConfigDetectionResult{
				BaseConfig:        "src/core/base.yaml",
				AdditionalConfigs: []string{"src/core/monthly_calendar.yaml"},
			},
			expected: []string{"src/core/base.yaml", "src/core/monthly_calendar.yaml"},
		},
		{
			name: "without additional configs",
			result: &ConfigDetectionResult{
				BaseConfig:        "src/core/base.yaml",
				AdditionalConfigs: []string{},
			},
			expected: []string{"src/core/base.yaml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := tt.result.GetConfigPaths()

			if len(paths) != len(tt.expected) {
				t.Errorf("Expected %d paths, got %d", len(tt.expected), len(paths))
			}

			for i, expected := range tt.expected {
				if i < len(paths) && paths[i] != expected {
					t.Errorf("Expected path %s, got %s", expected, paths[i])
				}
			}
		})
	}
}

func TestHasPhaseSubPhaseColumns(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected bool
	}{
		{
			name:     "has phase and sub-phase",
			lines:    []string{"Name,Phase,Sub-Phase,StartDate,EndDate"},
			expected: true,
		},
		{
			name:     "has phase but no sub-phase",
			lines:    []string{"Name,Phase,StartDate,EndDate"},
			expected: false,
		},
		{
			name:     "has sub-phase but no phase",
			lines:    []string{"Name,SubPhase,StartDate,EndDate"},
			expected: false,
		},
		{
			name:     "case insensitive",
			lines:    []string{"Name,PHASE,SUB-PHASE,StartDate,EndDate"},
			expected: true,
		},
		{
			name:     "empty lines",
			lines:    []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPhaseSubPhaseColumns(tt.lines)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
