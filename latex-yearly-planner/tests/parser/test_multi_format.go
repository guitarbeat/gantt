package main

import (
	"fmt"
	"os"
	"time"

	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
	"latex-yearly-planner/internal/generator"
)

// Test program for multi-format generation
func main() {
	fmt.Println("Testing Multi-Format Generation...")

	// Create temporary output directory
	outputDir := "test_multi_output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// Create a test configuration
	cfg := config.Config{
		WeekStart: time.Monday,
		MonthsWithTasks: []data.MonthYear{
			{Month: time.January, Year: 2024},
		},
		Layout: config.Layout{
			Numbers: config.Numbers{
				ArrayStretch: "1.0",
			},
			Lengths: config.Lengths{
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

	// Create multi-format generator
	workDir, _ := os.Getwd()
	multiGen := generator.NewMultiFormatGenerator(workDir, outputDir)

	// Set up generation options
	options := generator.MultiFormatOptions{
		Formats: []generator.OutputFormat{
			generator.OutputFormatPDF,
			generator.OutputFormatLaTeX,
		},
		ViewConfigs: []generator.ViewConfig{
			generator.GetDefaultViewConfig(generator.ViewTypeMonthly),
			generator.GetDefaultViewConfig(generator.ViewTypeWeekly),
		},
		OutputPrefix:   "test_multi",
		IncludeStats:   true,
		IncludeLegend:  true,
		BatchMode:      false,
		ParallelJobs:   2,
	}

	// Generate multi-format output
	fmt.Printf("Generating %d formats for %d views\n", len(options.Formats), len(options.ViewConfigs))
	fmt.Printf("Output directory: %s\n", outputDir)

	result, err := multiGen.GenerateMultiFormat(cfg, options)
	if err != nil {
		fmt.Printf("Multi-format generation failed: %v\n", err)
		if result != nil {
			printMultiFormatResult(result)
		}
		os.Exit(1)
	}

	// Print success results
	fmt.Printf("Multi-format generation completed successfully!\n")
	printMultiFormatResult(result)
}

func printMultiFormatResult(result *generator.MultiFormatResult) {
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
