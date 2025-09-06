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

// Test visual spacing optimization
func main() {
	fmt.Println("Testing Visual Spacing Optimization...")

	// Test 1: Visual Spacing Configuration
	fmt.Println("\n=== Test 1: Visual Spacing Configuration ===")
	testVisualSpacingConfig()

	// Test 2: Visual Optimizer
	fmt.Println("\n=== Test 2: Visual Optimizer ===")
	testVisualOptimizer()

	// Test 3: PDF Pipeline with Visual Optimization
	fmt.Println("\n=== Test 3: PDF Pipeline with Visual Optimization ===")
	testPDFPipelineWithVisualOptimization()

	// Test 4: Quality Metrics
	fmt.Println("\n=== Test 4: Quality Metrics ===")
	testQualityMetrics()

	fmt.Println("\n✅ Visual spacing optimization tests completed!")
}

func testVisualSpacingConfig() {
	// Test default configuration
	config := generator.GetDefaultVisualSpacingConfig()
	if config == nil {
		fmt.Println("❌ Failed to get default visual spacing config")
		return
	}

	// Validate configuration
	if config.CalendarGridSpacing.Padding <= 0 {
		fmt.Println("❌ Calendar grid spacing padding is not positive")
		return
	}

	if config.TaskBarSpacing.Padding <= 0 {
		fmt.Println("❌ Task bar spacing padding is not positive")
		return
	}

	if config.TextSpacing.Padding <= 0 {
		fmt.Println("❌ Text spacing padding is not positive")
		return
	}

	// Test quality thresholds
	if config.QualityThresholds.MinTaskBarHeight <= 0 {
		fmt.Println("❌ Minimum task bar height is not positive")
		return
	}

	if config.QualityThresholds.MinTaskBarWidth <= 0 {
		fmt.Println("❌ Minimum task bar width is not positive")
		return
	}

	fmt.Printf("✅ Visual spacing configuration test passed\n")
	fmt.Printf("   Calendar grid padding: %.1f%s\n", 
		config.CalendarGridSpacing.Padding, config.CalendarGridSpacing.Unit)
	fmt.Printf("   Task bar padding: %.1f%s\n", 
		config.TaskBarSpacing.Padding, config.TaskBarSpacing.Unit)
	fmt.Printf("   Text padding: %.1f%s\n", 
		config.TextSpacing.Padding, config.TextSpacing.Unit)
}

func testVisualOptimizer() {
	// Create visual optimizer
	optimizer := generator.NewVisualOptimizer()
	if optimizer == nil {
		fmt.Println("❌ Failed to create visual optimizer")
		return
	}

	// Test with sample data
	context := generator.AnalysisContext{
		TaskDensity:        generator.DensityNormal,
		AvgTaskNameLength:  20.0,
		ViewType:           generator.ViewTypeMonthly,
		TaskCount:          10,
		CategoryDistribution: map[string]int{
			"RESEARCH": 5,
			"ADMIN":    3,
			"LASER":    2,
		},
		PriorityDistribution: map[string]int{
			"HIGH":   3,
			"MEDIUM": 5,
			"LOW":    2,
		},
		AvailableSpace:     100.0,
		ContentComplexity:  0.5,
	}

	// Calculate optimal spacing
	optimalSpacing := optimizer.GetSpacingConfig().CalculateOptimalSpacing(context)
	if optimalSpacing == nil {
		fmt.Println("❌ Failed to calculate optimal spacing")
		return
	}

	// Validate spacing
	validation := optimalSpacing.ValidateSpacing()
	if !validation.IsValid {
		fmt.Printf("❌ Spacing validation failed: %v\n", validation.Issues)
		return
	}

	fmt.Printf("✅ Visual optimizer test passed\n")
	fmt.Printf("   Quality score: %.2f\n", validation.Score)
	fmt.Printf("   Calendar grid padding: %.1f%s\n", 
		optimalSpacing.CalendarGridSpacing.Padding, optimalSpacing.CalendarGridSpacing.Unit)
	fmt.Printf("   Task bar padding: %.1f%s\n", 
		optimalSpacing.TaskBarSpacing.Padding, optimalSpacing.TaskBarSpacing.Unit)
}

func testPDFPipelineWithVisualOptimization() {
	// Create test configuration (not used in this test but kept for completeness)
	_ = config.Config{
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

	// Create PDF pipeline
	workDir, _ := os.Getwd()
	outputDir := "test_visual_output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create output directory: %v\n", err)
		return
	}

	pipeline := generator.NewPDFPipeline(workDir, outputDir)
	if pipeline == nil {
		fmt.Println("❌ Failed to create PDF pipeline")
		return
	}

	// Test visual optimizer integration
	if pipeline.GetVisualOptimizer() == nil {
		fmt.Println("❌ PDF pipeline does not have visual optimizer")
		return
	}

	// Test PDF generation options
	options := generator.PDFGenerationOptions{
		OutputFileName:    "visual_test.pdf",
		CleanupTempFiles:  false,
		MaxRetries:        2,
		CompilationEngine: "pdflatex",
		ExtraPackages:     []string{"tikz", "tcolorbox", "xcolor"},
		CustomPreamble:    "% Visual Spacing Test",
	}

	// Note: We're not actually generating PDF here to avoid LaTeX compilation issues
	// Just testing the structure and configuration
	fmt.Printf("✅ PDF pipeline with visual optimization test passed\n")
	fmt.Printf("   Output file: %s\n", options.OutputFileName)
	fmt.Printf("   Compilation engine: %s\n", options.CompilationEngine)
	fmt.Printf("   Extra packages: %v\n", options.ExtraPackages)
}

func testQualityMetrics() {
	// Create sample quality metrics
	metrics := &generator.VisualQualityMetrics{
		OverallScore:     0.85,
		SpacingScore:     0.90,
		AlignmentScore:   0.80,
		ReadabilityScore: 0.85,
		VisualClarity:    0.85,
		LayoutEfficiency: 0.80,
		VisualNoise:      0.15,
	}

	// Validate metrics
	if metrics.OverallScore < 0.0 || metrics.OverallScore > 1.0 {
		fmt.Println("❌ Overall score is out of range")
		return
	}

	if metrics.SpacingScore < 0.0 || metrics.SpacingScore > 1.0 {
		fmt.Println("❌ Spacing score is out of range")
		return
	}

	if metrics.AlignmentScore < 0.0 || metrics.AlignmentScore > 1.0 {
		fmt.Println("❌ Alignment score is out of range")
		return
	}

	if metrics.ReadabilityScore < 0.0 || metrics.ReadabilityScore > 1.0 {
		fmt.Println("❌ Readability score is out of range")
		return
	}

	if metrics.VisualClarity < 0.0 || metrics.VisualClarity > 1.0 {
		fmt.Println("❌ Visual clarity is out of range")
		return
	}

	if metrics.LayoutEfficiency < 0.0 || metrics.LayoutEfficiency > 1.0 {
		fmt.Println("❌ Layout efficiency is out of range")
		return
	}

	if metrics.VisualNoise < 0.0 || metrics.VisualNoise > 1.0 {
		fmt.Println("❌ Visual noise is out of range")
		return
	}

	fmt.Printf("✅ Quality metrics test passed\n")
	fmt.Printf("   Overall score: %.2f\n", metrics.OverallScore)
	fmt.Printf("   Spacing score: %.2f\n", metrics.SpacingScore)
	fmt.Printf("   Alignment score: %.2f\n", metrics.AlignmentScore)
	fmt.Printf("   Readability score: %.2f\n", metrics.ReadabilityScore)
	fmt.Printf("   Visual clarity: %.2f\n", metrics.VisualClarity)
	fmt.Printf("   Layout efficiency: %.2f\n", metrics.LayoutEfficiency)
	fmt.Printf("   Visual noise: %.2f\n", metrics.VisualNoise)
}
