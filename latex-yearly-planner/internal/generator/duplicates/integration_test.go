package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
	"latex-yearly-planner/internal/layout"
)

// IntegrationTestSuite provides comprehensive testing for the integrated system
type IntegrationTestSuite struct {
	workDir    string
	outputDir  string
	logger     PDFLogger
}

// NewIntegrationTestSuite creates a new integration test suite
func NewIntegrationTestSuite() *IntegrationTestSuite {
	workDir, _ := os.Getwd()
	outputDir := filepath.Join(workDir, "test_integration_output")
	
	return &IntegrationTestSuite{
		workDir:   workDir,
		outputDir: outputDir,
		logger:    &TestLogger{},
	}
}

// TestLogger provides logging for integration tests
type TestLogger struct{}

func (l *TestLogger) Info(msg string, args ...interface{})  { fmt.Printf("[INFO] "+msg+"\n", args...) }
func (l *TestLogger) Error(msg string, args ...interface{}) { fmt.Printf("[ERROR] "+msg+"\n", args...) }
func (l *TestLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[DEBUG] "+msg+"\n", args...) }

// TestFullIntegration tests the complete integrated system
func TestFullIntegration(t *testing.T) {
	suite := NewIntegrationTestSuite()
	defer suite.cleanup()

	// Create test output directory
	if err := os.MkdirAll(suite.outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	t.Run("LayoutIntegration", func(t *testing.T) {
		suite.testLayoutIntegration(t)
	})

	t.Run("PDFPipeline", func(t *testing.T) {
		suite.testPDFPipeline(t)
	})

	t.Run("MultiFormatGeneration", func(t *testing.T) {
		suite.testMultiFormatGeneration(t)
	})

	t.Run("BatchProcessing", func(t *testing.T) {
		suite.testBatchProcessing(t)
	})

	t.Run("QualityValidation", func(t *testing.T) {
		suite.testQualityValidation(t)
	})
}

// testLayoutIntegration tests the layout integration system
func (suite *IntegrationTestSuite) testLayoutIntegration(t *testing.T) {
	// Create test tasks
	tasks := suite.createTestTasks()

	// Create layout integration
	layoutIntegration := NewLayoutIntegration()

	// Process tasks with layout
	result, err := layoutIntegration.ProcessTasksWithLayout(tasks)
	if err != nil {
		t.Fatalf("Layout integration failed: %v", err)
	}

	// Validate result
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.TaskBars) == 0 {
		t.Fatal("Expected task bars to be generated")
	}

	// Validate task bar properties
	for i, bar := range result.TaskBars {
		if bar.TaskName == "" {
			t.Errorf("Task bar %d: TaskName should not be empty", i)
		}
		if bar.Color == "" {
			t.Errorf("Task bar %d: Color should not be empty", i)
		}
		if bar.StartX < 0 {
			t.Errorf("Task bar %d: StartX should be non-negative", i)
		}
		if bar.Y < 0 {
			t.Errorf("Task bar %d: Y should be non-negative", i)
		}
		if bar.Width <= 0 {
			t.Errorf("Task bar %d: Width should be positive", i)
		}
		if bar.Height <= 0 {
			t.Errorf("Task bar %d: Height should be positive", i)
		}
	}

	// Test statistics
	stats := result.Statistics
	if stats.TotalTasks == 0 {
		t.Error("Expected non-zero total tasks")
	}
	if stats.ProcessedBars == 0 {
		t.Error("Expected non-zero processed bars")
	}

	suite.logger.Info("Layout integration test passed: %d task bars generated", len(result.TaskBars))
}

// testPDFPipeline tests the PDF generation pipeline
func (suite *IntegrationTestSuite) testPDFPipeline(t *testing.T) {
	// Create test configuration
	cfg := suite.createTestConfig()

	// Create PDF pipeline
	pipeline := NewPDFPipeline(suite.workDir, suite.outputDir)
	pipeline.SetLogger(suite.logger)

	// Set up generation options
	options := PDFGenerationOptions{
		OutputFileName:    "integration_test.pdf",
		CleanupTempFiles:  false, // Keep for debugging
		MaxRetries:        2,
		CompilationEngine: "pdflatex",
		ExtraPackages:     []string{"tikz", "tcolorbox"},
		CustomPreamble:    "% Integration Test PDF",
	}

	// Generate PDF
	result, err := pipeline.GeneratePDF(cfg, options)
	if err != nil {
		t.Fatalf("PDF generation failed: %v", err)
	}

	// Validate result
	if !result.Success {
		t.Fatal("PDF generation should succeed")
	}

	if result.OutputPath == "" {
		t.Fatal("Output path should not be empty")
	}

	// Check if PDF file exists and has reasonable size
	if info, err := os.Stat(result.OutputPath); err != nil {
		t.Fatalf("PDF file not found: %v", err)
	} else {
		if info.Size() < 1024 { // Less than 1KB is suspicious
			t.Errorf("PDF file too small: %d bytes", info.Size())
		}
	}

	suite.logger.Info("PDF pipeline test passed: %s (%d bytes)", result.OutputPath, result.FileSize)
}

// testMultiFormatGeneration tests multi-format generation
func (suite *IntegrationTestSuite) testMultiFormatGeneration(t *testing.T) {
	// Create test configuration
	cfg := suite.createTestConfig()

	// Create multi-format generator
	generator := NewMultiFormatGenerator(suite.workDir, suite.outputDir)
	generator.SetLogger(suite.logger)

	// Set up generation options
	options := MultiFormatOptions{
		Formats: []OutputFormat{
			OutputFormatPDF,
			OutputFormatLaTeX,
		},
		ViewConfigs: []ViewConfig{
			GetDefaultViewConfig(ViewTypeMonthly),
			GetDefaultViewConfig(ViewTypeWeekly),
		},
		OutputPrefix:   "integration_test",
		IncludeStats:   true,
		IncludeLegend:  true,
		BatchMode:      false,
		ParallelJobs:   2,
	}

	// Generate multi-format output
	result, err := generator.GenerateMultiFormat(cfg, options)
	if err != nil {
		t.Fatalf("Multi-format generation failed: %v", err)
	}

	// Validate result
	if !result.Success {
		t.Fatal("Multi-format generation should succeed")
	}

	if result.TotalFiles == 0 {
		t.Fatal("Expected at least one file to be generated")
	}

	// Validate each format
	for format, viewResults := range result.Results {
		for viewType, formatResult := range viewResults {
			if !formatResult.Success {
				t.Errorf("Format %s, view %s should succeed", format, viewType)
			}
			if formatResult.FilePath == "" {
				t.Errorf("Format %s, view %s should have file path", format, viewType)
			}
			if formatResult.FileSize <= 0 {
				t.Errorf("Format %s, view %s should have positive file size", format, viewType)
			}
		}
	}

	suite.logger.Info("Multi-format generation test passed: %d files generated", result.TotalFiles)
}

// testBatchProcessing tests batch processing
func (suite *IntegrationTestSuite) testBatchProcessing(t *testing.T) {
	// Create batch processor
	processor := NewBatchProcessor(suite.workDir, suite.outputDir)
	processor.SetLogger(suite.logger)

	// Create batch configuration
	batchConfig := BatchConfig{
		Name:        "Integration Test Batch",
		Description: "Test batch for integration testing",
		OutputDir:   suite.outputDir,
		Parallel:    false,
		MaxWorkers:  2,
		Configs: []BatchItem{
			{
				Name:       "test-monthly",
				ConfigFile: "configs/planner_config.yaml",
				ViewConfigs: []ViewConfig{
					GetDefaultViewConfig(ViewTypeMonthly),
				},
				Formats: []OutputFormat{OutputFormatPDF},
				Options: MultiFormatOptions{
					IncludeStats:  true,
					IncludeLegend: true,
					BatchMode:     true,
				},
			},
		},
	}

	// Process batch
	result, err := processor.ProcessBatch(batchConfig)
	if err != nil {
		t.Fatalf("Batch processing failed: %v", err)
	}

	// Validate result
	if !result.Success {
		t.Fatal("Batch processing should succeed")
	}

	if result.TotalBatches == 0 {
		t.Fatal("Expected at least one batch to be processed")
	}

	suite.logger.Info("Batch processing test passed: %d batches processed", result.TotalBatches)
}

// testQualityValidation tests output quality
func (suite *IntegrationTestSuite) testQualityValidation(t *testing.T) {
	// Create test configuration
	cfg := suite.createTestConfig()

	// Create PDF pipeline
	pipeline := NewPDFPipeline(suite.workDir, suite.outputDir)
	pipeline.SetLogger(suite.logger)

	// Generate PDF with quality validation
	options := PDFGenerationOptions{
		OutputFileName:    "quality_test.pdf",
		CleanupTempFiles:  false,
		MaxRetries:        2,
		CompilationEngine: "pdflatex",
		ExtraPackages:     []string{"tikz", "tcolorbox", "xcolor"},
	}

	result, err := pipeline.GeneratePDF(cfg, options)
	if err != nil {
		t.Fatalf("Quality test PDF generation failed: %v", err)
	}

	// Validate quality metrics
	if result.FileSize < 1024 {
		t.Errorf("PDF file too small for quality test: %d bytes", result.FileSize)
	}

	if result.PageCount <= 0 {
		t.Errorf("Expected positive page count, got: %d", result.PageCount)
	}

	// Check for compilation warnings
	if len(result.Warnings) > 0 {
		suite.logger.Info("Quality test warnings: %v", result.Warnings)
	}

	// Check compilation time (should be reasonable)
	if result.CompilationTime > time.Minute*5 {
		t.Errorf("Compilation time too long: %v", result.CompilationTime)
	}

	suite.logger.Info("Quality validation test passed: %s (%d bytes, %d pages)", 
		result.OutputPath, result.FileSize, result.PageCount)
}

// createTestTasks creates test tasks for integration testing
func (suite *IntegrationTestSuite) createTestTasks() []*data.Task {
	tasks := []*data.Task{
		{
			ID:          "task1",
			Name:        "Research Proposal",
			Description: "Write research proposal for funding",
			Category:    "PROPOSAL",
			Priority:    4,
			StartDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 1, 25, 0, 0, 0, 0, time.Local),
		},
		{
			ID:          "task2",
			Name:        "Laser Setup",
			Description: "Configure laser system for experiments",
			Category:    "LASER",
			Priority:    3,
			StartDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 1, 30, 0, 0, 0, 0, time.Local),
		},
		{
			ID:          "task3",
			Name:        "Data Analysis",
			Description: "Analyze experimental data",
			Category:    "IMAGING",
			Priority:    2,
			StartDate:   time.Date(2024, 1, 25, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 2, 5, 0, 0, 0, 0, time.Local),
		},
		{
			ID:          "task4",
			Name:        "Admin Meeting",
			Description: "Monthly admin meeting",
			Category:    "ADMIN",
			Priority:    1,
			StartDate:   time.Date(2024, 1, 30, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 1, 30, 0, 0, 0, 0, time.Local),
		},
		{
			ID:          "task5",
			Name:        "Dissertation Chapter",
			Description: "Write chapter 3 of dissertation",
			Category:    "DISSERTATION",
			Priority:    5,
			StartDate:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local),
			EndDate:     time.Date(2024, 2, 15, 0, 0, 0, 0, time.Local),
		},
	}

	return tasks
}

// createTestConfig creates a test configuration
func (suite *IntegrationTestSuite) createTestConfig() config.Config {
	return config.Config{
		WeekStart: time.Monday,
		MonthsWithTasks: []data.MonthYear{
			{Month: time.January, Year: 2024},
			{Month: time.February, Year: 2024},
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

// cleanup cleans up test files
func (suite *IntegrationTestSuite) cleanup() {
	if err := os.RemoveAll(suite.outputDir); err != nil {
		suite.logger.Error("Failed to cleanup test directory: %v", err)
	}
}

// TestIntegrationWithRealData tests integration with real CSV data
func TestIntegrationWithRealData(t *testing.T) {
	// Check if test CSV files exist
	csvFiles := []string{
		"input/data.cleaned.csv",
		"input/test_single.csv",
		"input/test_triple.csv",
	}

	var testCSV string
	for _, csvFile := range csvFiles {
		if _, err := os.Stat(csvFile); err == nil {
			testCSV = csvFile
			break
		}
	}

	if testCSV == "" {
		t.Skip("No test CSV files found, skipping real data test")
	}

	suite := NewIntegrationTestSuite()
	defer suite.cleanup()

	// Create test output directory
	if err := os.MkdirAll(suite.outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Create configuration with CSV file
	cfg := config.Config{
		CSVFilePath: testCSV,
		WeekStart:   time.Monday,
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

	// Test PDF generation with real data
	pipeline := NewPDFPipeline(suite.workDir, suite.outputDir)
	pipeline.SetLogger(suite.logger)

	options := PDFGenerationOptions{
		OutputFileName:    "real_data_test.pdf",
		CleanupTempFiles:  false,
		MaxRetries:        2,
		CompilationEngine: "pdflatex",
		ExtraPackages:     []string{"tikz", "tcolorbox"},
	}

	result, err := pipeline.GeneratePDF(cfg, options)
	if err != nil {
		t.Fatalf("Real data PDF generation failed: %v", err)
	}

	if !result.Success {
		t.Fatal("Real data PDF generation should succeed")
	}

	t.Logf("Real data test passed: %s (%d bytes)", result.OutputPath, result.FileSize)
}
