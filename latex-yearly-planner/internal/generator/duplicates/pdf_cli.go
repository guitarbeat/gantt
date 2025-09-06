package generator

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"latex-yearly-planner/internal/config"
)

// PDFCLIOptions represents command line options for PDF generation
type PDFCLIOptions struct {
	ConfigFile        string
	OutputFile        string
	OutputDir         string
	CleanupTemp       bool
	Verbose           bool
	Engine            string
	MaxRetries        int
	ExtraPackages     string
	CustomPreambleFile string
}

// ParseCLIArgs parses command line arguments for PDF generation
func ParseCLIArgs() (*PDFCLIOptions, error) {
	opts := &PDFCLIOptions{}

	flag.StringVar(&opts.ConfigFile, "config", "", "Path to configuration file")
	flag.StringVar(&opts.OutputFile, "output", "planner.pdf", "Output PDF filename")
	flag.StringVar(&opts.OutputDir, "output-dir", ".", "Output directory")
	flag.BoolVar(&opts.CleanupTemp, "cleanup", true, "Cleanup temporary files after generation")
	flag.BoolVar(&opts.Verbose, "verbose", false, "Enable verbose logging")
	flag.StringVar(&opts.Engine, "engine", "pdflatex", "LaTeX engine (pdflatex, xelatex, lualatex)")
	flag.IntVar(&opts.MaxRetries, "retries", 3, "Maximum compilation retries")
	flag.StringVar(&opts.ExtraPackages, "packages", "", "Comma-separated list of extra LaTeX packages")
	flag.StringVar(&opts.CustomPreambleFile, "preamble", "", "Path to custom LaTeX preamble file")

	flag.Parse()

	// Validate options
	if opts.ConfigFile == "" {
		return nil, fmt.Errorf("configuration file is required")
	}

	if _, err := os.Stat(opts.ConfigFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", opts.ConfigFile)
	}

	// Validate LaTeX engine
	validEngines := map[string]bool{
		"pdflatex": true,
		"xelatex":  true,
		"lualatex": true,
	}
	if !validEngines[opts.Engine] {
		return nil, fmt.Errorf("invalid LaTeX engine: %s (valid: pdflatex, xelatex, lualatex)", opts.Engine)
	}

	return opts, nil
}

// GeneratePDFFromCLI generates a PDF using command line options
func GeneratePDFFromCLI(opts *PDFCLIOptions) error {
	// Load configuration
	cfg, err := config.NewConfig(opts.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create PDF pipeline
	workDir, _ := os.Getwd()
	pipeline := NewPDFPipeline(workDir, opts.OutputDir)

	// Set up logging
	if opts.Verbose {
		pipeline.SetLogger(&VerboseLogger{})
	}

	// Parse extra packages
	var extraPackages []string
	if opts.ExtraPackages != "" {
		extraPackages = strings.Split(opts.ExtraPackages, ",")
		for i, pkg := range extraPackages {
			extraPackages[i] = strings.TrimSpace(pkg)
		}
	}

	// Load custom preamble if provided
	var customPreamble string
	if opts.CustomPreambleFile != "" {
		preambleBytes, err := os.ReadFile(opts.CustomPreambleFile)
		if err != nil {
			return fmt.Errorf("failed to read custom preamble: %w", err)
		}
		customPreamble = string(preambleBytes)
	}

	// Set up generation options
	genOptions := PDFGenerationOptions{
		OutputFileName:    opts.OutputFile,
		CleanupTempFiles:  opts.CleanupTemp,
		MaxRetries:        opts.MaxRetries,
		CompilationEngine: opts.Engine,
		ExtraPackages:     extraPackages,
		CustomPreamble:    customPreamble,
	}

	// Generate PDF
	fmt.Printf("Starting PDF generation...\n")
	fmt.Printf("  Config: %s\n", opts.ConfigFile)
	fmt.Printf("  Output: %s\n", filepath.Join(opts.OutputDir, opts.OutputFile))
	fmt.Printf("  Engine: %s\n", opts.Engine)

	result, err := pipeline.GeneratePDF(cfg, genOptions)
	if err != nil {
		fmt.Printf("PDF generation failed: %v\n", err)
		if result != nil {
			printGenerationResult(result, true)
		}
		return err
	}

	// Print results
	fmt.Printf("PDF generation completed successfully!\n")
	printGenerationResult(result, opts.Verbose)

	return nil
}

// VerboseLogger provides detailed logging output
type VerboseLogger struct{}

func (l *VerboseLogger) Info(msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] INFO: "+msg+"\n", append([]interface{}{timestamp}, args...)...)
}

func (l *VerboseLogger) Error(msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] ERROR: "+msg+"\n", append([]interface{}{timestamp}, args...)...)
}

func (l *VerboseLogger) Debug(msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] DEBUG: "+msg+"\n", append([]interface{}{timestamp}, args...)...)
}

// printGenerationResult prints the PDF generation results
func printGenerationResult(result *PDFGenerationResult, verbose bool) {
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

	if verbose && result.LaTeXLog != "" {
		fmt.Printf("\n=== LaTeX Compilation Log ===\n")
		fmt.Printf("%s\n", result.LaTeXLog)
	}

	if verbose && len(result.TempFilesCreated) > 0 {
		fmt.Printf("\nTemporary Files Created:\n")
		for _, file := range result.TempFilesCreated {
			fmt.Printf("  - %s\n", file)
		}
	}
}

// RunPDFGenerationCLI is the main entry point for the CLI
func RunPDFGenerationCLI() {
	opts, err := ParseCLIArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nUsage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := GeneratePDFFromCLI(opts); err != nil {
		fmt.Fprintf(os.Stderr, "PDF generation failed: %v\n", err)
		os.Exit(1)
	}
}
