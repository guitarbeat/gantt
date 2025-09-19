package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"phd-dissertation-planner/internal/config"
	"phd-dissertation-planner/internal/data"
	"phd-dissertation-planner/internal/generator"
)

// Simple test program for the PDF pipeline
func main() {
	fmt.Println("Testing PDF Pipeline Integration...")

	// Create temporary output directory
	outputDir := "test_output_pipeline"
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

	// Create PDF pipeline
	workDir, _ := os.Getwd()
	pipeline := generator.NewPDFPipeline(workDir, outputDir)

	// Set up generation options
	options := generator.PDFGenerationOptions{
		OutputFileName:    "test_pipeline.pdf",
		CleanupTempFiles:  false, // Keep temp files for debugging
		MaxRetries:        2,
		CompilationEngine: "pdflatex",
		ExtraPackages:     []string{"tikz"},
		CustomPreamble:    "% Test PDF Pipeline Generation",
	}

	// Generate PDF
	fmt.Printf("Generating PDF with options:\n")
	fmt.Printf("  Output: %s\n", filepath.Join(outputDir, options.OutputFileName))
	fmt.Printf("  Engine: %s\n", options.CompilationEngine)

	result, err := pipeline.GeneratePDF(cfg, options)
	if err != nil {
		fmt.Printf("PDF generation failed: %v\n", err)
		if result != nil {
			printResult(result)
		}
		os.Exit(1)
	}

	// Print success results
	fmt.Printf("PDF generation completed successfully!\n")
	printResult(result)
}

func printResult(result *generator.PDFGenerationResult) {
	fmt.Printf("\n=== Generation Results ===\n")
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output Path: %s\n", result.OutputPath)
	fmt.Printf("Compilation Time: %v\n", result.CompilationTime)
	fmt.Printf("File Size: %d bytes\n", result.FileSize)
	
	if result.PageCount > 0 {
		fmt.Printf("Page Count: %d\n", result.PageCount)
	}

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

	if len(result.TempFilesCreated) > 0 {
		fmt.Printf("\nTemporary Files Created:\n")
		for _, file := range result.TempFilesCreated {
			fmt.Printf("  - %s\n", file)
		}
	}

	if result.LaTeXLog != "" {
		fmt.Printf("\n=== LaTeX Log (last 500 chars) ===\n")
		log := result.LaTeXLog
		if len(log) > 500 {
			log = "..." + log[len(log)-500:]
		}
		fmt.Printf("%s\n", log)
	}
}
