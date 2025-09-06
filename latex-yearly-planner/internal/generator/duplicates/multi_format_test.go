package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestViewConfig(t *testing.T) {
	t.Run("DefaultConfigs", func(t *testing.T) {
		viewTypes := []ViewType{ViewTypeMonthly, ViewTypeWeekly, ViewTypeYearly, ViewTypeQuarterly, ViewTypeDaily}
		
		for _, viewType := range viewTypes {
			config := GetDefaultViewConfig(viewType)
			
			if config.Type != viewType {
				t.Errorf("Expected view type %s, got %s", viewType, config.Type)
			}
			
			if config.TemplateName == "" {
				t.Errorf("Template name should not be empty for %s", viewType)
			}
			
			if config.Title == "" {
				t.Errorf("Title should not be empty for %s", viewType)
			}
		}
	})

	t.Run("ViewPresets", func(t *testing.T) {
		presets := GetViewPresets()
		
		if len(presets) == 0 {
			t.Error("Expected at least one view preset")
		}
		
		// Check for required presets
		requiredPresets := []string{"monthly-standard", "weekly-detailed", "yearly-overview"}
		for _, required := range requiredPresets {
			found := false
			for _, preset := range presets {
				if preset.Name == required {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Required preset not found: %s", required)
			}
		}
	})

	t.Run("PresetByName", func(t *testing.T) {
		preset, err := GetPresetByName("monthly-standard")
		if err != nil {
			t.Fatalf("Failed to get preset: %v", err)
		}
		
		if preset.Name != "monthly-standard" {
			t.Errorf("Expected preset name 'monthly-standard', got '%s'", preset.Name)
		}
		
		// Test non-existent preset
		_, err = GetPresetByName("non-existent")
		if err == nil {
			t.Error("Expected error for non-existent preset")
		}
	})

	t.Run("ConfigValidation", func(t *testing.T) {
		// Test valid config
		validConfig := GetDefaultViewConfig(ViewTypeMonthly)
		if err := ValidateViewConfig(validConfig); err != nil {
			t.Errorf("Valid config should not have errors: %v", err)
		}
		
		// Test invalid view type
		invalidConfig := validConfig
		invalidConfig.Type = "invalid"
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid view type")
		}
		
		// Test invalid page size
		invalidConfig = validConfig
		invalidConfig.PageSize = "invalid"
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid page size")
		}
		
		// Test invalid orientation
		invalidConfig = validConfig
		invalidConfig.Orientation = "invalid"
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid orientation")
		}
		
		// Test invalid color scheme
		invalidConfig = validConfig
		invalidConfig.ColorScheme = "invalid"
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid color scheme")
		}
		
		// Test invalid font size
		invalidConfig = validConfig
		invalidConfig.FontSize = "invalid"
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid font size")
		}
		
		// Test invalid layout density
		invalidConfig = validConfig
		invalidConfig.LayoutDensity = "invalid"
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid layout density")
		}
		
		// Test invalid numeric values
		invalidConfig = validConfig
		invalidConfig.MaxTasksPerDay = 0
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid max tasks per day")
		}
		
		invalidConfig = validConfig
		invalidConfig.TaskBarHeight = -1
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid task bar height")
		}
		
		invalidConfig = validConfig
		invalidConfig.TaskBarSpacing = -1
		if err := ValidateViewConfig(invalidConfig); err == nil {
			t.Error("Expected error for invalid task bar spacing")
		}
	})
}

func TestMultiFormatGenerator(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "multi-format-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Create multi-format generator
	generator := NewMultiFormatGenerator(tempDir, outputDir)

	t.Run("ConfigValidation", func(t *testing.T) {
		// Test empty options
		emptyOptions := MultiFormatOptions{}
		if err := generator.validateMultiFormatOptions(emptyOptions); err == nil {
			t.Error("Expected error for empty options")
		}
		
		// Test invalid format
		invalidOptions := MultiFormatOptions{
			Formats:     []OutputFormat{"invalid"},
			ViewConfigs: []ViewConfig{GetDefaultViewConfig(ViewTypeMonthly)},
		}
		if err := generator.validateMultiFormatOptions(invalidOptions); err == nil {
			t.Error("Expected error for invalid format")
		}
		
		// Test invalid view config
		invalidViewConfig := GetDefaultViewConfig(ViewTypeMonthly)
		invalidViewConfig.Type = "invalid"
		invalidOptions = MultiFormatOptions{
			Formats:     []OutputFormat{OutputFormatPDF},
			ViewConfigs: []ViewConfig{invalidViewConfig},
		}
		if err := generator.validateMultiFormatOptions(invalidOptions); err == nil {
			t.Error("Expected error for invalid view config")
		}
	})

	t.Run("FilenameGeneration", func(t *testing.T) {
		viewConfig := GetDefaultViewConfig(ViewTypeMonthly)
		
		filename := generator.generateFilename("test", OutputFormatPDF, viewConfig)
		expectedPrefix := "test-monthly-"
		if !contains(filename, expectedPrefix) {
			t.Errorf("Expected filename to contain '%s', got '%s'", expectedPrefix, filename)
		}
		
		if !contains(filename, ".pdf") {
			t.Errorf("Expected filename to contain '.pdf', got '%s'", filename)
		}
	})
}

func TestBatchProcessor(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "batch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create batch processor
	processor := NewBatchProcessor(tempDir, tempDir)

	t.Run("BatchConfigValidation", func(t *testing.T) {
		// Test empty batch config
		emptyConfig := BatchConfig{}
		if err := processor.validateBatchConfig(emptyConfig); err == nil {
			t.Error("Expected error for empty batch config")
		}
		
		// Test batch config with empty name
		configWithEmptyName := BatchConfig{
			Configs: []BatchItem{
				{
					Name:       "test",
					ConfigFile: "test.yaml",
					ViewConfigs: []ViewConfig{GetDefaultViewConfig(ViewTypeMonthly)},
					Formats:    []OutputFormat{OutputFormatPDF},
				},
			},
		}
		if err := processor.validateBatchConfig(configWithEmptyName); err == nil {
			t.Error("Expected error for empty batch name")
		}
		
		// Test batch config with empty items
		configWithEmptyItems := BatchConfig{
			Name: "test",
		}
		if err := processor.validateBatchConfig(configWithEmptyItems); err == nil {
			t.Error("Expected error for empty batch items")
		}
		
		// Test batch config with invalid item
		configWithInvalidItem := BatchConfig{
			Name: "test",
			Configs: []BatchItem{
				{
					Name:       "", // Empty name
					ConfigFile: "test.yaml",
					ViewConfigs: []ViewConfig{GetDefaultViewConfig(ViewTypeMonthly)},
					Formats:    []OutputFormat{OutputFormatPDF},
				},
			},
		}
		if err := processor.validateBatchConfig(configWithInvalidItem); err == nil {
			t.Error("Expected error for invalid batch item")
		}
	})

	t.Run("SampleBatchConfig", func(t *testing.T) {
		sampleConfig := CreateSampleBatchConfig()
		
		if sampleConfig.Name == "" {
			t.Error("Sample config should have a name")
		}
		
		if len(sampleConfig.Configs) == 0 {
			t.Error("Sample config should have at least one item")
		}
		
		// Validate the sample config
		if err := processor.validateBatchConfig(sampleConfig); err != nil {
			t.Errorf("Sample config should be valid: %v", err)
		}
	})
}

func TestApplyViewConfig(t *testing.T) {
	viewConfig := GetDefaultViewConfig(ViewTypeMonthly)
	options := &PDFGenerationOptions{}
	
	ApplyViewConfig(viewConfig, options)
	
	// Note: PDFGenerationOptions doesn't have TemplateName field
	// This test would need to be updated if TemplateName is added to the struct
	
	// Check that extra packages are added
	if len(options.ExtraPackages) == 0 {
		t.Error("Expected extra packages to be added")
	}
	
	// Check for required packages
	requiredPackages := []string{"tikz", "tcolorbox", "xcolor"}
	for _, pkg := range requiredPackages {
		found := false
		for _, addedPkg := range options.ExtraPackages {
			if addedPkg == pkg {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Required package %s not found in extra packages", pkg)
		}
	}
	
	// Check that custom preamble is set
	if options.CustomPreamble == "" {
		t.Error("Expected custom preamble to be set")
	}
}

func TestColorSchemeLaTeX(t *testing.T) {
	schemes := []ColorScheme{
		ColorSchemeDefault,
		ColorSchemeMinimal,
		ColorSchemeHighContrast,
		ColorSchemeColorBlind,
		ColorSchemeDark,
	}
	
	for _, scheme := range schemes {
		latex := generateColorSchemeLaTeX(scheme)
		if latex == "" {
			t.Errorf("Color scheme %s should generate LaTeX", scheme)
		}
	}
}

func TestLayoutDensityLaTeX(t *testing.T) {
	densities := []LayoutDensity{
		LayoutDensityCompact,
		LayoutDensityNormal,
		LayoutDensitySpacious,
		LayoutDensityMinimal,
	}
	
	for _, density := range densities {
		latex := generateLayoutDensityLaTeX(density)
		if latex == "" {
			t.Errorf("Layout density %s should generate LaTeX", density)
		}
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    contains(s[1:], substr) || 
		    (len(s) > 0 && s[:len(substr)] == substr))
}

// Note: simpleContains function is already defined above, removing duplicate
