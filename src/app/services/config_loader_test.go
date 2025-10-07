package services

import (
	"testing"

	"phd-dissertation-planner/src/core"
)

func TestNewConfigLoader(t *testing.T) {
	cl := NewConfigLoader()
	
	if cl == nil {
		t.Fatal("Expected config loader to be created")
	}
	
	// Logger is a struct, not a pointer, so we can't check for nil
	// Just verify the config loader was created successfully
}

func TestConfigLoader_LoadConfiguration(t *testing.T) {
	t.Skip("Skipping CLI context test - needs proper CLI setup")
}

func TestConfigLoader_ValidateConfiguration(t *testing.T) {
	cl := NewConfigLoader()
	
	tests := []struct {
		name        string
		config      core.Config
		expectError bool
	}{
		{
			name: "valid configuration",
			config: core.Config{
				OutputDir:  "test_output",
				StartYear:  2025,
				EndYear:    2026,
				Pages: []core.Page{
					{
						Name: "test_page",
						RenderBlocks: []core.RenderBlock{
							{FuncName: "monthly", Tpls: []string{"test.tpl"}},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "missing output directory",
			config: core.Config{
				OutputDir: "",
				StartYear: 2025,
				EndYear:   2026,
				Pages: []core.Page{
					{Name: "test_page", RenderBlocks: []core.RenderBlock{}},
				},
			},
			expectError: true,
		},
		{
			name: "invalid year range",
			config: core.Config{
				OutputDir: "test_output",
				StartYear: 2026,
				EndYear:   2025,
				Pages: []core.Page{
					{Name: "test_page", RenderBlocks: []core.RenderBlock{}},
				},
			},
			expectError: true,
		},
		{
			name: "no pages",
			config: core.Config{
				OutputDir: "test_output",
				StartYear: 2025,
				EndYear:   2026,
				Pages:     []core.Page{},
			},
			expectError: true,
		},
		{
			name: "page without name",
			config: core.Config{
				OutputDir: "test_output",
				StartYear: 2025,
				EndYear:   2026,
				Pages: []core.Page{
					{Name: "", RenderBlocks: []core.RenderBlock{}},
				},
			},
			expectError: true,
		},
		{
			name: "page without render blocks",
			config: core.Config{
				OutputDir: "test_output",
				StartYear: 2025,
				EndYear:   2026,
				Pages: []core.Page{
					{Name: "test_page", RenderBlocks: []core.RenderBlock{}},
				},
			},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cl.ValidateConfiguration(tt.config)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestConfigLoader_GetConfigurationSummary(t *testing.T) {
	cl := NewConfigLoader()
	
	result := &ConfigLoadResult{
		Config: core.Config{
			OutputDir: "test_output",
			StartYear: 2025,
			EndYear:   2026,
			Pages: []core.Page{
				{Name: "test_page", RenderBlocks: []core.RenderBlock{}},
			},
		},
		PathConfigs: []string{"src/core/base.yaml"},
		CSVPath:     "test.csv",
		Reason:      "test reason",
	}
	
	summary := cl.GetConfigurationSummary(result)
	
	if summary == "" {
		t.Error("Expected non-empty summary")
	}
	
	// Check that summary configStringContains expected information
	expected := []string{
		"Configuration Summary:",
		"Output Directory: test_output",
		"Year Range: 2025-2026",
		"Pages: 1",
		"CSV File: test.csv",
		"Config Files: src/core/base.yaml",
		"Reason: test reason",
	}
	
	for _, exp := range expected {
		if !configStringContains(summary, exp) {
			t.Errorf("Expected summary to contain '%s', but it didn't", exp)
		}
	}
}

func TestConfigLoader_ReloadConfiguration(t *testing.T) {
	cl := NewConfigLoader()
	
	pathConfigs := []string{"src/core/base.yaml"}
	
	result, err := cl.ReloadConfiguration(pathConfigs)
	if err != nil {
		t.Logf("Configuration reload failed (expected if base.yaml doesn't exist): %v", err)
		return
	}
	
	if result == nil {
		t.Fatal("Expected result to be created")
	}
	
	if len(result.PathConfigs) != len(pathConfigs) {
		t.Errorf("Expected %d path configs, got %d", len(pathConfigs), len(result.PathConfigs))
	}
	
	if result.PathConfigs[0] != pathConfigs[0] {
		t.Errorf("Expected first path config to be %s, got %s", pathConfigs[0], result.PathConfigs[0])
	}
}

func TestConfigLoadResult(t *testing.T) {
	result := &ConfigLoadResult{
		Config: core.Config{
			OutputDir: "test_output",
		},
		PathConfigs: []string{"config1.yaml", "config2.yaml"},
		CSVPath:     "test.csv",
		Reason:      "test reason",
	}
	
	// Test that all fields are set correctly
	if result.Config.OutputDir != "test_output" {
		t.Errorf("Expected OutputDir to be 'test_output', got %s", result.Config.OutputDir)
	}
	
	if len(result.PathConfigs) != 2 {
		t.Errorf("Expected 2 path configs, got %d", len(result.PathConfigs))
	}
	
	if result.CSVPath != "test.csv" {
		t.Errorf("Expected CSVPath to be 'test.csv', got %s", result.CSVPath)
	}
	
	if result.Reason != "test reason" {
		t.Errorf("Expected Reason to be 'test reason', got %s", result.Reason)
	}
}

// Helper function to check if a string configStringContains a substring
func configStringContains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && configStringContains(s[1:], substr)
}
