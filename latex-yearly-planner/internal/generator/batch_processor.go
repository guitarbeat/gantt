package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"latex-yearly-planner/internal/config"
)

// BatchConfig represents a batch processing configuration
type BatchConfig struct {
	Name        string
	Description string
	Configs     []BatchItem
	OutputDir   string
	Parallel    bool
	MaxWorkers  int
}

// BatchItem represents a single item in a batch
type BatchItem struct {
	Name        string
	ConfigFile  string
	ViewConfigs []ViewConfig
	Formats     []OutputFormat
	Options     MultiFormatOptions
}

// BatchResult represents the result of batch processing
type BatchResult struct {
	Success        bool
	TotalBatches   int
	SuccessfulBatches int
	FailedBatches  int
	TotalFiles     int
	SuccessfulFiles int
	FailedFiles    int
	Results        []BatchItemResult
	Errors         []string
	Warnings       []string
	ProcessingTime time.Duration
}

// BatchItemResult represents the result of processing a single batch item
type BatchItemResult struct {
	ItemName     string
	Success      bool
	FilesGenerated int
	FilesFailed  int
	Results      map[OutputFormat]map[string]*FormatResult
	Error        string
	Warning      string
	ProcessingTime time.Duration
}

// BatchProcessor handles batch processing of multiple configurations
type BatchProcessor struct {
	multiFormatGen *MultiFormatGenerator
	logger         PDFLogger
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(workDir, outputDir string) *BatchProcessor {
	return &BatchProcessor{
		multiFormatGen: NewMultiFormatGenerator(workDir, outputDir),
		logger:         DefaultLogger{},
	}
}

// SetLogger sets a custom logger for the batch processor
func (bp *BatchProcessor) SetLogger(logger PDFLogger) {
	bp.logger = logger
	bp.multiFormatGen.SetLogger(logger)
}

// ProcessBatch processes a batch configuration
func (bp *BatchProcessor) ProcessBatch(batchConfig BatchConfig) (*BatchResult, error) {
	startTime := time.Now()
	result := &BatchResult{
		Results: make([]BatchItemResult, 0),
		Errors:  make([]string, 0),
		Warnings: make([]string, 0),
	}

	bp.logger.Info("Starting batch processing: %s", batchConfig.Name)
	bp.logger.Info("Processing %d items", len(batchConfig.Configs))

	// Validate batch configuration
	if err := bp.validateBatchConfig(batchConfig); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Batch configuration validation failed: %v", err))
		return result, err
	}

	// Process each item in the batch
	for i, item := range batchConfig.Configs {
		bp.logger.Info("Processing batch item %d/%d: %s", i+1, len(batchConfig.Configs), item.Name)
		
		itemResult, err := bp.processBatchItem(item, batchConfig.OutputDir)
		if err != nil {
			bp.logger.Error("Failed to process batch item %s: %v", item.Name, err)
			itemResult = &BatchItemResult{
				ItemName: item.Name,
				Success:  false,
				Error:    err.Error(),
			}
			result.FailedBatches++
		} else {
			if itemResult.Success {
				result.SuccessfulBatches++
			} else {
				result.FailedBatches++
			}
		}
		
		result.Results = append(result.Results, *itemResult)
		result.TotalFiles += itemResult.FilesGenerated + itemResult.FilesFailed
		result.SuccessfulFiles += itemResult.FilesGenerated
		result.FailedFiles += itemResult.FilesFailed
	}

	result.Success = result.FailedBatches == 0
	result.TotalBatches = len(batchConfig.Configs)
	result.ProcessingTime = time.Since(startTime)

	bp.logger.Info("Batch processing completed: %d/%d batches successful", 
		result.SuccessfulBatches, result.TotalBatches)

	return result, nil
}

// processBatchItem processes a single batch item
func (bp *BatchProcessor) processBatchItem(item BatchItem, outputDir string) (*BatchItemResult, error) {
	itemStartTime := time.Now()
	result := &BatchItemResult{
		ItemName: item.Name,
	}

	// Load configuration
	cfg, err := config.NewConfig(item.ConfigFile)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to load configuration: %v", err)
		return result, err
	}

	// Set up multi-format options
	options := item.Options
	options.ViewConfigs = item.ViewConfigs
	options.Formats = item.Formats

	// Create item-specific output directory
	itemOutputDir := filepath.Join(outputDir, item.Name)
	if err := os.MkdirAll(itemOutputDir, 0755); err != nil {
		result.Error = fmt.Sprintf("Failed to create output directory: %v", err)
		return result, err
	}

	// Update multi-format generator output directory
	bp.multiFormatGen.outputDir = itemOutputDir

	// Generate multi-format output
	multiResult, err := bp.multiFormatGen.GenerateMultiFormat(cfg, options)
	if err != nil {
		result.Error = fmt.Sprintf("Multi-format generation failed: %v", err)
		return result, err
	}

	// Update result
	result.Success = multiResult.Success
	result.FilesGenerated = multiResult.SuccessfulFiles
	result.FilesFailed = multiResult.FailedFiles
	result.Results = multiResult.Results
	result.ProcessingTime = time.Since(itemStartTime)

	if len(multiResult.Warnings) > 0 {
		result.Warning = fmt.Sprintf("%d warnings: %v", len(multiResult.Warnings), multiResult.Warnings)
	}

	return result, nil
}

// LoadBatchConfig loads a batch configuration from a file
func LoadBatchConfig(filename string) (*BatchConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read batch config file: %w", err)
	}

	var batchConfig BatchConfig
	if err := json.Unmarshal(data, &batchConfig); err != nil {
		return nil, fmt.Errorf("failed to parse batch config JSON: %w", err)
	}

	return &batchConfig, nil
}

// SaveBatchConfig saves a batch configuration to a file
func SaveBatchConfig(batchConfig BatchConfig, filename string) error {
	data, err := json.MarshalIndent(batchConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal batch config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write batch config file: %w", err)
	}

	return nil
}

// CreateSampleBatchConfig creates a sample batch configuration
func CreateSampleBatchConfig() BatchConfig {
	return BatchConfig{
		Name:        "Sample Batch Configuration",
		Description: "Sample batch configuration with multiple views and formats",
		OutputDir:   "batch_output",
		Parallel:    false,
		MaxWorkers:  4,
		Configs: []BatchItem{
			{
				Name:       "monthly-calendar",
				ConfigFile: "configs/planner_config.yaml",
				ViewConfigs: []ViewConfig{
					GetDefaultViewConfig(ViewTypeMonthly),
					GetDefaultViewConfig(ViewTypeWeekly),
				},
				Formats: []OutputFormat{OutputFormatPDF, OutputFormatLaTeX},
				Options: MultiFormatOptions{
					IncludeStats:  true,
					IncludeLegend: true,
					BatchMode:     true,
				},
			},
			{
				Name:       "yearly-overview",
				ConfigFile: "configs/planner_config.yaml",
				ViewConfigs: []ViewConfig{
					GetDefaultViewConfig(ViewTypeYearly),
					GetDefaultViewConfig(ViewTypeQuarterly),
				},
				Formats: []OutputFormat{OutputFormatPDF, OutputFormatHTML},
				Options: MultiFormatOptions{
					IncludeStats:  false,
					IncludeLegend: true,
					BatchMode:     true,
				},
			},
		},
	}
}

// validateBatchConfig validates a batch configuration
func (bp *BatchProcessor) validateBatchConfig(batchConfig BatchConfig) error {
	if batchConfig.Name == "" {
		return fmt.Errorf("batch name is required")
	}

	if len(batchConfig.Configs) == 0 {
		return fmt.Errorf("at least one batch item must be specified")
	}

	for i, item := range batchConfig.Configs {
		if item.Name == "" {
			return fmt.Errorf("batch item %d: name is required", i)
		}

		if item.ConfigFile == "" {
			return fmt.Errorf("batch item %d: config file is required", i)
		}

		if _, err := os.Stat(item.ConfigFile); os.IsNotExist(err) {
			return fmt.Errorf("batch item %d: config file does not exist: %s", i, item.ConfigFile)
		}

		if len(item.ViewConfigs) == 0 {
			return fmt.Errorf("batch item %d: at least one view config must be specified", i)
		}

		if len(item.Formats) == 0 {
			return fmt.Errorf("batch item %d: at least one output format must be specified", i)
		}

		// Validate view configs
		for j, viewConfig := range item.ViewConfigs {
			if err := ValidateViewConfig(viewConfig); err != nil {
				return fmt.Errorf("batch item %d, view config %d: %w", i, j, err)
			}
		}

		// Validate formats
		validFormats := map[OutputFormat]bool{
			OutputFormatPDF:   true,
			OutputFormatLaTeX: true,
			OutputFormatHTML:  true,
			OutputFormatSVG:   true,
			OutputFormatPNG:   true,
		}

		for j, format := range item.Formats {
			if !validFormats[format] {
				return fmt.Errorf("batch item %d, format %d: invalid output format: %s", i, j, format)
			}
		}
	}

	return nil
}

// PrintBatchResult prints a formatted batch result
func PrintBatchResult(result *BatchResult) {
	fmt.Printf("\n=== Batch Processing Results ===\n")
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Total Batches: %d\n", result.TotalBatches)
	fmt.Printf("Successful Batches: %d\n", result.SuccessfulBatches)
	fmt.Printf("Failed Batches: %d\n", result.FailedBatches)
	fmt.Printf("Total Files: %d\n", result.TotalFiles)
	fmt.Printf("Successful Files: %d\n", result.SuccessfulFiles)
	fmt.Printf("Failed Files: %d\n", result.FailedFiles)
	fmt.Printf("Processing Time: %v\n", result.ProcessingTime)

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

	fmt.Printf("\n=== Individual Results ===\n")
	for i, itemResult := range result.Results {
		fmt.Printf("\n%d. %s\n", i+1, itemResult.ItemName)
		fmt.Printf("   Success: %v\n", itemResult.Success)
		fmt.Printf("   Files Generated: %d\n", itemResult.FilesGenerated)
		fmt.Printf("   Files Failed: %d\n", itemResult.FilesFailed)
		fmt.Printf("   Processing Time: %v\n", itemResult.ProcessingTime)
		
		if itemResult.Error != "" {
			fmt.Printf("   Error: %s\n", itemResult.Error)
		}
		
		if itemResult.Warning != "" {
			fmt.Printf("   Warning: %s\n", itemResult.Warning)
		}
	}
}
