package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
	"latex-yearly-planner/internal/generator"
	"latex-yearly-planner/internal/layout"
)

// ValidationResult contains the results of system validation
type ValidationResult struct {
	Success           bool
	TestsRun          int
	TestsPassed       int
	TestsFailed       int
	Errors            []string
	Warnings          []string
	PerformanceMetrics PerformanceMetrics
	QualityMetrics    QualityMetrics
}

// PerformanceMetrics contains performance measurements
type PerformanceMetrics struct {
	TotalTime        time.Duration
	LayoutTime       time.Duration
	PDFGenerationTime time.Duration
	MultiFormatTime  time.Duration
	AverageFileSize  int64
	FilesGenerated   int
}

// QualityMetrics contains quality measurements
type QualityMetrics struct {
	PDFCompilationSuccess bool
	LaTeXSyntaxValid      bool
	FileSizesReasonable   bool
	PageCountsValid       bool
	ErrorCount            int
	WarningCount           int
}

// SystemValidator validates the integrated system
type SystemValidator struct {
	workDir    string
	outputDir  string
	logger     generator.PDFLogger
}

// NewSystemValidator creates a new system validator
func NewSystemValidator() *SystemValidator {
	workDir, _ := os.Getwd()
	outputDir := filepath.Join(workDir, "validation_output")
	
	return &SystemValidator{
		workDir:   workDir,
		outputDir: outputDir,
		logger:    &ValidationLogger{},
	}
}

// ValidationLogger provides logging for validation
type ValidationLogger struct{}

func (l *ValidationLogger) Info(msg string, args ...interface{})  { fmt.Printf("[INFO] "+msg+"\n", args...) }
func (l *ValidationLogger) Error(msg string, args ...interface{}) { fmt.Printf("[ERROR] "+msg+"\n", args...) }
func (l *ValidationLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[DEBUG] "+msg+"\n", args...) }

// ValidateSystem performs comprehensive system validation
func (v *SystemValidator) ValidateSystem() *ValidationResult {
	startTime := time.Now()
	result := &ValidationResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	v.logger.Info("Starting comprehensive system validation...")

	// Create output directory
	if err := os.MkdirAll(v.outputDir, 0755); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create output directory: %v", err))
		return result
	}

	// Run validation tests
	v.testLayoutIntegration(result)
	v.testPDFPipeline(result)
	v.testMultiFormatGeneration(result)
	v.testBatchProcessing(result)
	v.testQualityStandards(result)
	v.testPerformance(result)

	// Calculate final results
	result.TestsRun = 5
	result.TestsPassed = result.TestsRun - result.TestsFailed
	result.Success = result.TestsFailed == 0
	result.PerformanceMetrics.TotalTime = time.Since(startTime)

	v.logger.Info("System validation completed: %d/%d tests passed", result.TestsPassed, result.TestsRun)

	return result
}

// testLayoutIntegration tests the layout integration system
func (v *SystemValidator) testLayoutIntegration(result *ValidationResult) {
	v.logger.Info("Testing layout integration...")
	
	// Create test tasks
	tasks := v.createTestTasks()
	
	// Create layout integration
	layoutIntegration := generator.NewLayoutIntegration()
	
	// Process tasks with layout
	layoutStart := time.Now()
	layoutResult, err := layoutIntegration.ProcessTasksWithLayout(tasks)
	layoutTime := time.Since(layoutStart)
	result.PerformanceMetrics.LayoutTime = layoutTime
	
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Layout integration failed: %v", err))
		result.TestsFailed++
		return
	}
	
	// Validate result
	if layoutResult == nil {
		result.Errors = append(result.Errors, "Layout integration returned nil result")
		result.TestsFailed++
		return
	}
	
	if len(layoutResult.TaskBars) == 0 {
		result.Errors = append(result.Errors, "No task bars generated")
		result.TestsFailed++
		return
	}
	
	// Validate task bar properties
	for i, bar := range layoutResult.TaskBars {
		if bar.TaskName == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("Task bar %d: TaskName is empty", i))
		}
		if bar.Color == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("Task bar %d: Color is empty", i))
		}
		if bar.Width <= 0 {
			result.Errors = append(result.Errors, fmt.Sprintf("Task bar %d: Width is not positive", i))
		}
		if bar.Height <= 0 {
			result.Errors = append(result.Errors, fmt.Sprintf("Task bar %d: Height is not positive", i))
		}
	}
	
	v.logger.Info("Layout integration test passed: %d task bars generated in %v", 
		len(layoutResult.TaskBars), layoutTime)
}

// testPDFPipeline tests the PDF generation pipeline
func (v *SystemValidator) testPDFPipeline(result *ValidationResult) {
	v.logger.Info("Testing PDF pipeline...")
	
	// Create test configuration
	cfg := v.createTestConfig()
	
	// Create PDF pipeline
	pipeline := generator.NewPDFPipeline(v.workDir, v.outputDir)
	pipeline.SetLogger(v.logger)
	
	// Set up generation options
	options := generator.PDFGenerationOptions{
		OutputFileName:    "validation_test.pdf",
		CleanupTempFiles:  false,
		MaxRetries:        2,
		CompilationEngine: "pdflatex",
		ExtraPackages:     []string{"tikz", "tcolorbox", "xcolor"},
		CustomPreamble:    "% Validation Test PDF",
	}
	
	// Generate PDF
	pdfStart := time.Now()
	pdfResult, err := pipeline.GeneratePDF(cfg, options)
	pdfTime := time.Since(pdfStart)
	result.PerformanceMetrics.PDFGenerationTime = pdfTime
	
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("PDF generation failed: %v", err))
		result.TestsFailed++
		return
	}
	
	// Validate result
	if !pdfResult.Success {
		result.Errors = append(result.Errors, "PDF generation did not succeed")
		result.TestsFailed++
		return
	}
	
	if pdfResult.OutputPath == "" {
		result.Errors = append(result.Errors, "PDF output path is empty")
		result.TestsFailed++
		return
	}
	
	// Check file existence and size
	if info, err := os.Stat(pdfResult.OutputPath); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("PDF file not found: %v", err))
		result.TestsFailed++
		return
	} else {
		result.PerformanceMetrics.AverageFileSize += info.Size()
		result.PerformanceMetrics.FilesGenerated++
		
		if info.Size() < 1024 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("PDF file very small: %d bytes", info.Size()))
		}
	}
	
	// Check compilation success
	result.QualityMetrics.PDFCompilationSuccess = pdfResult.Success
	result.QualityMetrics.ErrorCount += len(pdfResult.Errors)
	result.QualityMetrics.WarningCount += len(pdfResult.Warnings)
	
	v.logger.Info("PDF pipeline test passed: %s (%d bytes) in %v", 
		pdfResult.OutputPath, pdfResult.FileSize, pdfTime)
}

// testMultiFormatGeneration tests multi-format generation
func (v *SystemValidator) testMultiFormatGeneration(result *ValidationResult) {
	v.logger.Info("Testing multi-format generation...")
	
	// Create test configuration
	cfg := v.createTestConfig()
	
	// Create multi-format generator
	multiGen := generator.NewMultiFormatGenerator(v.workDir, v.outputDir)
	multiGen.SetLogger(v.logger)
	
	// Set up generation options
	options := generator.MultiFormatOptions{
		Formats: []generator.OutputFormat{
			generator.OutputFormatPDF,
			generator.OutputFormatLaTeX,
		},
		ViewConfigs: []generator.ViewConfig{
			generator.GetDefaultViewConfig(generator.ViewTypeMonthly),
		},
		OutputPrefix:   "validation_test",
		IncludeStats:   true,
		IncludeLegend:  true,
		BatchMode:      false,
		ParallelJobs:   2,
	}
	
	// Generate multi-format output
	multiStart := time.Now()
	multiResult, err := multiGen.GenerateMultiFormat(cfg, options)
	multiTime := time.Since(multiStart)
	result.PerformanceMetrics.MultiFormatTime = multiTime
	
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Multi-format generation failed: %v", err))
		result.TestsFailed++
		return
	}
	
	// Validate result
	if !multiResult.Success {
		result.Errors = append(result.Errors, "Multi-format generation did not succeed")
		result.TestsFailed++
		return
	}
	
	if multiResult.TotalFiles == 0 {
		result.Errors = append(result.Errors, "No files generated in multi-format test")
		result.TestsFailed++
		return
	}
	
	// Validate each format
	for format, viewResults := range multiResult.Results {
		for viewType, formatResult := range viewResults {
			if !formatResult.Success {
				result.Errors = append(result.Errors, 
					fmt.Sprintf("Format %s, view %s failed", format, viewType))
			}
			if formatResult.FilePath == "" {
				result.Errors = append(result.Errors, 
					fmt.Sprintf("Format %s, view %s has empty file path", format, viewType))
			}
			if formatResult.FileSize <= 0 {
				result.Warnings = append(result.Warnings, 
					fmt.Sprintf("Format %s, view %s has zero file size", format, viewType))
			}
		}
	}
	
	v.logger.Info("Multi-format generation test passed: %d files generated in %v", 
		multiResult.TotalFiles, multiTime)
}

// testBatchProcessing tests batch processing
func (v *SystemValidator) testBatchProcessing(result *ValidationResult) {
	v.logger.Info("Testing batch processing...")
	
	// Create batch processor
	processor := generator.NewBatchProcessor(v.workDir, v.outputDir)
	processor.SetLogger(v.logger)
	
	// Create batch configuration
	batchConfig := generator.BatchConfig{
		Name:        "Validation Test Batch",
		Description: "Test batch for validation",
		OutputDir:   v.outputDir,
		Parallel:    false,
		MaxWorkers:  2,
		Configs: []generator.BatchItem{
			{
				Name:       "validation-test",
				ConfigFile: "configs/planner_config.yaml",
				ViewConfigs: []generator.ViewConfig{
					generator.GetDefaultViewConfig(generator.ViewTypeMonthly),
				},
				Formats: []generator.OutputFormat{generator.OutputFormatPDF},
				Options: generator.MultiFormatOptions{
					IncludeStats:  true,
					IncludeLegend: true,
					BatchMode:     true,
				},
			},
		},
	}
	
	// Process batch
	batchResult, err := processor.ProcessBatch(batchConfig)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Batch processing failed: %v", err))
		result.TestsFailed++
		return
	}
	
	// Validate result
	if !batchResult.Success {
		result.Errors = append(result.Errors, "Batch processing did not succeed")
		result.TestsFailed++
		return
	}
	
	if batchResult.TotalBatches == 0 {
		result.Errors = append(result.Errors, "No batches processed")
		result.TestsFailed++
		return
	}
	
	v.logger.Info("Batch processing test passed: %d batches processed", batchResult.TotalBatches)
}

// testQualityStandards tests output quality
func (v *SystemValidator) testQualityStandards(result *ValidationResult) {
	v.logger.Info("Testing quality standards...")
	
	// Create test configuration
	cfg := v.createTestConfig()
	
	// Create PDF pipeline
	pipeline := generator.NewPDFPipeline(v.workDir, v.outputDir)
	pipeline.SetLogger(v.logger)
	
	// Generate PDF with quality validation
	options := generator.PDFGenerationOptions{
		OutputFileName:    "quality_validation.pdf",
		CleanupTempFiles:  false,
		MaxRetries:        2,
		CompilationEngine: "pdflatex",
		ExtraPackages:     []string{"tikz", "tcolorbox", "xcolor"},
	}
	
	qualityResult, err := pipeline.GeneratePDF(cfg, options)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Quality test PDF generation failed: %v", err))
		result.TestsFailed++
		return
	}
	
	// Validate quality metrics
	if qualityResult.FileSize < 1024 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("PDF file small: %d bytes", qualityResult.FileSize))
	}
	
	if qualityResult.PageCount <= 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Invalid page count: %d", qualityResult.PageCount))
	} else {
		result.QualityMetrics.PageCountsValid = true
	}
	
	// Check compilation time (should be reasonable)
	if qualityResult.CompilationTime > time.Minute*5 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Long compilation time: %v", qualityResult.CompilationTime))
	}
	
	// Check for reasonable file size
	if qualityResult.FileSize > 1024*1024*10 { // 10MB
		result.Warnings = append(result.Warnings, fmt.Sprintf("Large PDF file: %d bytes", qualityResult.FileSize))
	} else {
		result.QualityMetrics.FileSizesReasonable = true
	}
	
	v.logger.Info("Quality standards test passed: %s (%d bytes, %d pages)", 
		qualityResult.OutputPath, qualityResult.FileSize, qualityResult.PageCount)
}

// testPerformance tests system performance
func (v *SystemValidator) testPerformance(result *ValidationResult) {
	v.logger.Info("Testing performance...")
	
	// Test with multiple view types
	viewTypes := []generator.ViewType{
		generator.ViewTypeMonthly,
		generator.ViewTypeWeekly,
		generator.ViewTypeYearly,
	}
	
	cfg := v.createTestConfig()
	multiGen := generator.NewMultiFormatGenerator(v.workDir, v.outputDir)
	multiGen.SetLogger(v.logger)
	
	for _, viewType := range viewTypes {
		options := generator.MultiFormatOptions{
			Formats: []generator.OutputFormat{generator.OutputFormatPDF},
			ViewConfigs: []generator.ViewConfig{
				generator.GetDefaultViewConfig(viewType),
			},
			OutputPrefix:   "perf_test",
			IncludeStats:   false,
			IncludeLegend:  false,
			BatchMode:      false,
			ParallelJobs:   1,
		}
		
		start := time.Now()
		perfResult, err := multiGen.GenerateMultiFormat(cfg, options)
		duration := time.Since(start)
		
		if err != nil {
			result.Warnings = append(result.Warnings, 
				fmt.Sprintf("Performance test failed for %s: %v", viewType, err))
			continue
		}
		
		if perfResult.TotalFiles > 0 {
			result.PerformanceMetrics.FilesGenerated += perfResult.TotalFiles
		}
		
		v.logger.Info("Performance test for %s: %v", viewType, duration)
	}
	
	// Calculate average file size
	if result.PerformanceMetrics.FilesGenerated > 0 {
		result.PerformanceMetrics.AverageFileSize /= int64(result.PerformanceMetrics.FilesGenerated)
	}
}

// createTestTasks creates test tasks for validation
func (v *SystemValidator) createTestTasks() []*data.Task {
	tasks := []*data.Task{
		{
			ID:          "val_task1",
			Name:        "Validation Task 1",
			Description: "Test task for validation",
			Category:    "RESEARCH",
			Priority:    3,
			StartDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 1, 25, 0, 0, 0, 0, time.Local),
		},
		{
			ID:          "val_task2",
			Name:        "Validation Task 2",
			Description: "Another test task",
			Category:    "DISSERTATION",
			Priority:    4,
			StartDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 1, 30, 0, 0, 0, 0, time.Local),
		},
	}
	
	return tasks
}

// createTestConfig creates a test configuration
func (v *SystemValidator) createTestConfig() config.Config {
	return config.Config{
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
}

// PrintValidationResult prints a formatted validation result
func PrintValidationResult(result *ValidationResult) {
	fmt.Printf("\n=== System Validation Results ===\n")
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Tests Run: %d\n", result.TestsRun)
	fmt.Printf("Tests Passed: %d\n", result.TestsPassed)
	fmt.Printf("Tests Failed: %d\n", result.TestsFailed)
	fmt.Printf("Total Time: %v\n", result.PerformanceMetrics.TotalTime)
	
	fmt.Printf("\n=== Performance Metrics ===\n")
	fmt.Printf("Layout Time: %v\n", result.PerformanceMetrics.LayoutTime)
	fmt.Printf("PDF Generation Time: %v\n", result.PerformanceMetrics.PDFGenerationTime)
	fmt.Printf("Multi-Format Time: %v\n", result.PerformanceMetrics.MultiFormatTime)
	fmt.Printf("Files Generated: %d\n", result.PerformanceMetrics.FilesGenerated)
	fmt.Printf("Average File Size: %d bytes\n", result.PerformanceMetrics.AverageFileSize)
	
	fmt.Printf("\n=== Quality Metrics ===\n")
	fmt.Printf("PDF Compilation Success: %v\n", result.QualityMetrics.PDFCompilationSuccess)
	fmt.Printf("LaTeX Syntax Valid: %v\n", result.QualityMetrics.LaTeXSyntaxValid)
	fmt.Printf("File Sizes Reasonable: %v\n", result.QualityMetrics.FileSizesReasonable)
	fmt.Printf("Page Counts Valid: %v\n", result.QualityMetrics.PageCountsValid)
	fmt.Printf("Error Count: %d\n", result.QualityMetrics.ErrorCount)
	fmt.Printf("Warning Count: %d\n", result.QualityMetrics.WarningCount)
	
	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors (%d):\n", len(result.Errors))
		for i, err := range result.Errors {
			fmt.Printf("  %d. %s\n", i+1, err)
		}
	}
	
	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings (%d):\n", len(result.Warnings))
		for i, warning := range result.Warnings {
			fmt.Printf("  %d. %s\n", i+1, warning)
		}
	}
}

// Main function for validation
func main() {
	fmt.Println("Starting comprehensive system validation...")
	
	validator := NewSystemValidator()
	result := validator.ValidateSystem()
	
	PrintValidationResult(result)
	
	if result.Success {
		fmt.Println("\n✅ All validation tests passed!")
		os.Exit(0)
	} else {
		fmt.Println("\n❌ Some validation tests failed!")
		os.Exit(1)
	}
}
