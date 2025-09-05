package generator

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"latex-yearly-planner/internal/config"
)

// MultiFormatCLIOptions represents command line options for multi-format generation
type MultiFormatCLIOptions struct {
	ConfigFile        string
	OutputDir         string
	Formats           string
	Views             string
	Preset            string
	BatchConfig       string
	OutputPrefix      string
	IncludeStats      bool
	IncludeLegend     bool
	Verbose           bool
	Parallel          bool
	MaxWorkers        int
}

// ParseMultiFormatCLIArgs parses command line arguments for multi-format generation
func ParseMultiFormatCLIArgs() (*MultiFormatCLIOptions, error) {
	opts := &MultiFormatCLIOptions{}

	flag.StringVar(&opts.ConfigFile, "config", "", "Path to configuration file")
	flag.StringVar(&opts.OutputDir, "output-dir", "multi_output", "Output directory")
	flag.StringVar(&opts.Formats, "formats", "pdf,latex", "Comma-separated list of output formats (pdf,latex,html,svg,png)")
	flag.StringVar(&opts.Views, "views", "monthly", "Comma-separated list of view types (monthly,weekly,yearly,quarterly,daily)")
	flag.StringVar(&opts.Preset, "preset", "", "Use a predefined view preset")
	flag.StringVar(&opts.BatchConfig, "batch", "", "Path to batch configuration file")
	flag.StringVar(&opts.OutputPrefix, "prefix", "planner", "Output file prefix")
	flag.BoolVar(&opts.IncludeStats, "stats", false, "Include layout statistics")
	flag.BoolVar(&opts.IncludeLegend, "legend", true, "Include category legend")
	flag.BoolVar(&opts.Verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&opts.Parallel, "parallel", false, "Enable parallel processing")
	flag.IntVar(&opts.MaxWorkers, "workers", 4, "Maximum number of parallel workers")

	flag.Parse()

	// Validate options
	if opts.ConfigFile == "" && opts.BatchConfig == "" {
		return nil, fmt.Errorf("either --config or --batch must be specified")
	}

	if opts.ConfigFile != "" && opts.BatchConfig != "" {
		return nil, fmt.Errorf("cannot specify both --config and --batch")
	}

	// Validate formats
	validFormats := map[string]bool{
		"pdf":    true,
		"latex":  true,
		"html":   true,
		"svg":    true,
		"png":    true,
	}
	
	formats := strings.Split(opts.Formats, ",")
	for _, format := range formats {
		format = strings.TrimSpace(format)
		if !validFormats[format] {
			return nil, fmt.Errorf("invalid output format: %s", format)
		}
	}

	// Validate views
	validViews := map[string]bool{
		"monthly":   true,
		"weekly":    true,
		"yearly":    true,
		"quarterly": true,
		"daily":     true,
	}
	
	views := strings.Split(opts.Views, ",")
	for _, view := range views {
		view = strings.TrimSpace(view)
		if !validViews[view] {
			return nil, fmt.Errorf("invalid view type: %s", view)
		}
	}

	return opts, nil
}

// RunMultiFormatGeneration runs multi-format generation from command line
func RunMultiFormatGeneration(opts *MultiFormatCLIOptions) error {
	// Create output directory
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create multi-format generator
	workDir, _ := os.Getwd()
	generator := NewMultiFormatGenerator(workDir, opts.OutputDir)

	// Set up logging
	if opts.Verbose {
		generator.SetLogger(&VerboseLogger{})
	}

	if opts.BatchConfig != "" {
		return runBatchGeneration(generator, opts)
	} else {
		return runSingleConfigGeneration(generator, opts)
	}
}

// runBatchGeneration runs batch processing
func runBatchGeneration(generator *MultiFormatGenerator, opts *MultiFormatCLIOptions) error {
	fmt.Printf("Loading batch configuration from: %s\n", opts.BatchConfig)

	// Load batch configuration
	batchConfig, err := LoadBatchConfig(opts.BatchConfig)
	if err != nil {
		return fmt.Errorf("failed to load batch configuration: %w", err)
	}

	// Override output directory if specified
	if opts.OutputDir != "multi_output" {
		batchConfig.OutputDir = opts.OutputDir
	}

	// Override parallel processing settings
	batchConfig.Parallel = opts.Parallel
	batchConfig.MaxWorkers = opts.MaxWorkers

	// Create batch processor
	workDir, _ := os.Getwd()
	processor := NewBatchProcessor(workDir, batchConfig.OutputDir)
	if opts.Verbose {
		processor.SetLogger(&VerboseLogger{})
	}

	// Process batch
	fmt.Printf("Processing batch: %s\n", batchConfig.Name)
	fmt.Printf("Description: %s\n", batchConfig.Description)
	fmt.Printf("Output directory: %s\n", batchConfig.OutputDir)
	fmt.Printf("Items to process: %d\n", len(batchConfig.Configs))

	result, err := processor.ProcessBatch(*batchConfig)
	if err != nil {
		fmt.Printf("Batch processing failed: %v\n", err)
		if result != nil {
			PrintBatchResult(result)
		}
		return err
	}

	// Print results
	PrintBatchResult(result)

	return nil
}

// runSingleConfigGeneration runs single configuration generation
func runSingleConfigGeneration(generator *MultiFormatGenerator, opts *MultiFormatCLIOptions) error {
	fmt.Printf("Loading configuration from: %s\n", opts.ConfigFile)

	// Load configuration
	cfg, err := config.NewConfig(opts.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Parse formats
	formatStrings := strings.Split(opts.Formats, ",")
	formats := make([]OutputFormat, len(formatStrings))
	for i, formatStr := range formatStrings {
		formats[i] = OutputFormat(strings.TrimSpace(formatStr))
	}

	// Parse views and create view configs
	viewStrings := strings.Split(opts.Views, ",")
	viewConfigs := make([]ViewConfig, len(viewStrings))
	
	for i, viewStr := range viewStrings {
		viewType := ViewType(strings.TrimSpace(viewStr))
		
		if opts.Preset != "" {
			// Use preset
			preset, err := GetPresetByName(opts.Preset)
			if err != nil {
				return fmt.Errorf("failed to get preset %s: %w", opts.Preset, err)
			}
			viewConfigs[i] = preset.Config
		} else {
			// Use default config for view type
			viewConfigs[i] = GetDefaultViewConfig(viewType)
		}
		
		// Apply CLI overrides
		viewConfigs[i].ShowLayoutStats = opts.IncludeStats
		viewConfigs[i].ShowCategoryLegend = opts.IncludeLegend
	}

	// Set up multi-format options
	options := MultiFormatOptions{
		Formats:        formats,
		ViewConfigs:    viewConfigs,
		OutputPrefix:   opts.OutputPrefix,
		IncludeStats:   opts.IncludeStats,
		IncludeLegend:  opts.IncludeLegend,
		BatchMode:      false,
		ParallelJobs:   opts.MaxWorkers,
	}

	// Generate multi-format output
	fmt.Printf("Generating %d formats for %d views\n", len(formats), len(viewConfigs))
	fmt.Printf("Output directory: %s\n", opts.OutputDir)

	result, err := generator.GenerateMultiFormat(cfg, options)
	if err != nil {
		fmt.Printf("Multi-format generation failed: %v\n", err)
		return err
	}

	// Print results
	printMultiFormatResult(result)

	return nil
}

// printMultiFormatResult prints a formatted multi-format result
func printMultiFormatResult(result *MultiFormatResult) {
	fmt.Printf("\n=== Multi-Format Generation Results ===\n")
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Total Files: %d\n", result.TotalFiles)
	fmt.Printf("Successful Files: %d\n", result.SuccessfulFiles)
	fmt.Printf("Failed Files: %d\n", result.FailedFiles)
	fmt.Printf("Generation Time: %v\n", result.GenerationTime)

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

	fmt.Printf("\n=== Format Results ===\n")
	for format, viewResults := range result.Results {
		fmt.Printf("\n%s:\n", format)
		for viewType, formatResult := range viewResults {
			status := "✓"
			if !formatResult.Success {
				status = "✗"
			}
			fmt.Printf("  %s %s: %s", status, viewType, formatResult.FilePath)
			if formatResult.FileSize > 0 {
				fmt.Printf(" (%d bytes)", formatResult.FileSize)
			}
			if formatResult.PageCount > 0 {
				fmt.Printf(" (%d pages)", formatResult.PageCount)
			}
			fmt.Printf("\n")
			
			if formatResult.Error != "" {
				fmt.Printf("    Error: %s\n", formatResult.Error)
			}
			if formatResult.Warning != "" {
				fmt.Printf("    Warning: %s\n", formatResult.Warning)
			}
		}
	}
}

// RunMultiFormatCLI is the main entry point for multi-format CLI
func RunMultiFormatCLI() {
	opts, err := ParseMultiFormatCLIArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nUsage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := RunMultiFormatGeneration(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Multi-format generation failed: %v\n", err)
		os.Exit(1)
	}
}

// CreateSampleBatchConfigFile creates a sample batch configuration file
func CreateSampleBatchConfigFile(filename string) error {
	sampleConfig := CreateSampleBatchConfig()
	return SaveBatchConfig(sampleConfig, filename)
}

// ListViewPresets lists available view presets
func ListViewPresets() {
	fmt.Printf("Available View Presets:\n\n")
	presets := GetViewPresets()
	for _, preset := range presets {
		fmt.Printf("  %s\n", preset.Name)
		fmt.Printf("    Description: %s\n", preset.Description)
		fmt.Printf("    View Type: %s\n", preset.Config.Type)
		fmt.Printf("    Page Size: %s\n", preset.Config.PageSize)
		fmt.Printf("    Orientation: %s\n", preset.Config.Orientation)
		fmt.Printf("    Color Scheme: %s\n", preset.Config.ColorScheme)
		fmt.Printf("\n")
	}
}

// ListOutputFormats lists available output formats
func ListOutputFormats() {
	fmt.Printf("Available Output Formats:\n\n")
	formats := []struct {
		Name        string
		Description string
		Extension   string
	}{
		{"pdf", "Portable Document Format", ".pdf"},
		{"latex", "LaTeX source code", ".tex"},
		{"html", "HyperText Markup Language", ".html"},
		{"svg", "Scalable Vector Graphics", ".svg"},
		{"png", "Portable Network Graphics", ".png"},
	}
	
	for _, format := range formats {
		fmt.Printf("  %s (%s)\n", format.Name, format.Extension)
		fmt.Printf("    Description: %s\n", format.Description)
		fmt.Printf("\n")
	}
}
