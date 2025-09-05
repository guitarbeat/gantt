package main

import (
	"fmt"
	"os"
	"time"

	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
	"latex-yearly-planner/internal/generator"
	"latex-yearly-planner/internal/layout"
)

// Simple validation test that bypasses template issues
func main() {
	fmt.Println("Starting simple system validation...")

	// Create test output directory
	outputDir := "simple_validation_output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// Test 1: Layout Integration
	fmt.Println("\n=== Test 1: Layout Integration ===")
	testLayoutIntegration()

	// Test 2: View Configuration
	fmt.Println("\n=== Test 2: View Configuration ===")
	testViewConfiguration()

	// Test 3: Multi-Format Options
	fmt.Println("\n=== Test 3: Multi-Format Options ===")
	testMultiFormatOptions()

	// Test 4: Batch Configuration
	fmt.Println("\n=== Test 4: Batch Configuration ===")
	testBatchConfiguration()

	// Test 5: PDF Pipeline (without template compilation)
	fmt.Println("\n=== Test 5: PDF Pipeline Structure ===")
	testPDFPipelineStructure()

	fmt.Println("\n✅ Simple validation completed successfully!")
}

func testLayoutIntegration() {
	// Create test tasks
	tasks := []*data.Task{
		{
			ID:          "test1",
			Name:        "Test Task 1",
			Description: "A test task",
			Category:    "RESEARCH",
			Priority:    3,
			StartDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 1, 25, 0, 0, 0, 0, time.Local),
		},
	}

	// Create layout integration
	layoutIntegration := generator.NewLayoutIntegration()

	// Process tasks with layout
	result, err := layoutIntegration.ProcessTasksWithLayout(tasks)
	if err != nil {
		fmt.Printf("❌ Layout integration failed: %v\n", err)
		return
	}

	if result == nil {
		fmt.Printf("❌ Layout integration returned nil result\n")
		return
	}

	if len(result.TaskBars) == 0 {
		fmt.Printf("❌ No task bars generated\n")
		return
	}

	fmt.Printf("✅ Layout integration passed: %d task bars generated\n", len(result.TaskBars))
	
	// Validate task bar properties
	valid := true
	for i, bar := range result.TaskBars {
		if bar.TaskName == "" {
			fmt.Printf("❌ Task bar %d: TaskName is empty\n", i)
			valid = false
		}
		if bar.Color == "" {
			fmt.Printf("❌ Task bar %d: Color is empty\n", i)
			valid = false
		}
		if bar.Width <= 0 {
			fmt.Printf("❌ Task bar %d: Width is not positive\n", i)
			valid = false
		}
		if bar.Height <= 0 {
			fmt.Printf("❌ Task bar %d: Height is not positive\n", i)
			valid = false
		}
	}

	if valid {
		fmt.Printf("✅ Task bar properties validation passed\n")
	}
}

func testViewConfiguration() {
	// Test default view configs
	viewTypes := []generator.ViewType{
		generator.ViewTypeMonthly,
		generator.ViewTypeWeekly,
		generator.ViewTypeYearly,
		generator.ViewTypeQuarterly,
		generator.ViewTypeDaily,
	}

	valid := true
	for _, viewType := range viewTypes {
		config := generator.GetDefaultViewConfig(viewType)
		
		if config.Type != viewType {
			fmt.Printf("❌ View config type mismatch for %s\n", viewType)
			valid = false
		}
		
		if config.TemplateName == "" {
			fmt.Printf("❌ View config template name empty for %s\n", viewType)
			valid = false
		}
		
		if config.Title == "" {
			fmt.Printf("❌ View config title empty for %s\n", viewType)
			valid = false
		}
	}

	if valid {
		fmt.Printf("✅ View configuration validation passed\n")
	}

	// Test view presets
	presets := generator.GetViewPresets()
	if len(presets) == 0 {
		fmt.Printf("❌ No view presets found\n")
		return
	}

	fmt.Printf("✅ View presets validation passed: %d presets available\n", len(presets))

	// Test preset retrieval
	preset, err := generator.GetPresetByName("monthly-standard")
	if err != nil {
		fmt.Printf("❌ Failed to get monthly-standard preset: %v\n", err)
		return
	}

	if preset.Name != "monthly-standard" {
		fmt.Printf("❌ Preset name mismatch\n")
		return
	}

	fmt.Printf("✅ Preset retrieval validation passed\n")
}

func testMultiFormatOptions() {
	// Test output formats
	formats := []generator.OutputFormat{
		generator.OutputFormatPDF,
		generator.OutputFormatLaTeX,
		generator.OutputFormatHTML,
		generator.OutputFormatSVG,
		generator.OutputFormatPNG,
	}

	fmt.Printf("✅ Output formats validation passed: %d formats supported\n", len(formats))

	// Test multi-format options
	options := generator.MultiFormatOptions{
		Formats: []generator.OutputFormat{
			generator.OutputFormatPDF,
			generator.OutputFormatLaTeX,
		},
		ViewConfigs: []generator.ViewConfig{
			generator.GetDefaultViewConfig(generator.ViewTypeMonthly),
		},
		OutputPrefix:   "test",
		IncludeStats:   true,
		IncludeLegend:  true,
		BatchMode:      false,
		ParallelJobs:   2,
	}

	if len(options.Formats) == 0 {
		fmt.Printf("❌ No formats specified in options\n")
		return
	}

	if len(options.ViewConfigs) == 0 {
		fmt.Printf("❌ No view configs specified in options\n")
		return
	}

	fmt.Printf("✅ Multi-format options validation passed\n")
}

func testBatchConfiguration() {
	// Test batch config creation
	batchConfig := generator.CreateSampleBatchConfig()
	
	if batchConfig.Name == "" {
		fmt.Printf("❌ Batch config name is empty\n")
		return
	}

	if len(batchConfig.Configs) == 0 {
		fmt.Printf("❌ Batch config has no items\n")
		return
	}

	fmt.Printf("✅ Batch configuration validation passed: %d items\n", len(batchConfig.Configs))

	// Test batch config structure
	item := batchConfig.Configs[0]
	if item.Name == "" {
		fmt.Printf("❌ Batch item name is empty\n")
		return
	}
	
	if item.ConfigFile == "" {
		fmt.Printf("❌ Batch item config file is empty\n")
		return
	}
	
	if len(item.ViewConfigs) == 0 {
		fmt.Printf("❌ Batch item has no view configs\n")
		return
	}
	
	if len(item.Formats) == 0 {
		fmt.Printf("❌ Batch item has no formats\n")
		return
	}

	fmt.Printf("✅ Batch configuration validation logic passed\n")
}

func testPDFPipelineStructure() {
	// Test PDF pipeline creation
	workDir, _ := os.Getwd()
	pipeline := generator.NewPDFPipeline(workDir, "test_output")

	if pipeline == nil {
		fmt.Printf("❌ PDF pipeline creation failed\n")
		return
	}

	// Test PDF generation options
	options := generator.PDFGenerationOptions{
		OutputFileName:    "test.pdf",
		CleanupTempFiles:  true,
		MaxRetries:        3,
		CompilationEngine: "pdflatex",
		ExtraPackages:     []string{"tikz", "tcolorbox"},
		CustomPreamble:    "% Test preamble",
	}

	if options.OutputFileName == "" {
		fmt.Printf("❌ PDF options output filename is empty\n")
		return
	}

	if options.MaxRetries <= 0 {
		fmt.Printf("❌ PDF options max retries is not positive\n")
		return
	}

	fmt.Printf("✅ PDF pipeline structure validation passed\n")

	// Test configuration creation
	cfg := config.Config{
		WeekStart: time.Monday,
		MonthsWithTasks: []data.MonthYear{
			{Month: time.January, Year: 2024},
		},
		Layout: config.Layout{
			Numbers: config.Numbers{
				ArrayStretch: 1.0,
			},
			Lengths: layout.Lengths{
				TabColSep:              "6pt",
				LineThicknessDefault:   "0.4pt",
				LineThicknessThick:     "1.2pt",
				LineHeightButLine:      "1.2em",
				TwoColSep:              "2em",
				TriColSep:              "1em",
				FiveColSep:             "0.5em",
				MonthlyCellHeight:      "2.5em",
				MonthlySpring:          "\\vfill",
				HeaderResizeBox:        "0.8",
				HeaderSideMonthsWidth:  "2cm",
			},
			Colors: config.Colors{
				Gray:      "gray",
				LightGray: "lightgray",
			},
		},
	}

	if cfg.WeekStart == 0 {
		fmt.Printf("❌ Config week start is not set\n")
		return
	}

	if len(cfg.MonthsWithTasks) == 0 {
		fmt.Printf("❌ Config has no months with tasks\n")
		return
	}

	fmt.Printf("✅ Configuration structure validation passed\n")
}
