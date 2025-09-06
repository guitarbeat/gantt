package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"latex-yearly-planner/internal/config"
)

// PDFPipeline handles the complete PDF generation process
type PDFPipeline struct {
	workDir         string
	outputDir       string
	templateEngine  *Tpl
	layoutInteg     *LayoutIntegration
	visualOptimizer *VisualOptimizer
	logger          PDFLogger
}

// PDFLogger interface for logging PDF generation events
type PDFLogger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// DefaultLogger provides basic logging functionality
type DefaultLogger struct{}

func (l DefaultLogger) Info(msg string, args ...interface{})  { fmt.Printf("[INFO] "+msg+"\n", args...) }
func (l DefaultLogger) Error(msg string, args ...interface{}) { fmt.Printf("[ERROR] "+msg+"\n", args...) }
func (l DefaultLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[DEBUG] "+msg+"\n", args...) }

// PDFGenerationOptions configures PDF generation
type PDFGenerationOptions struct {
	OutputFileName    string
	CleanupTempFiles  bool
	MaxRetries        int
	CompilationEngine string // "pdflatex", "xelatex", "lualatex"
	ExtraPackages     []string
	CustomPreamble    string
}

// PDFGenerationResult contains the results of PDF generation
type PDFGenerationResult struct {
	Success           bool
	OutputPath        string
	CompilationTime   time.Duration
	LaTeXLog          string
	Errors            []string
	Warnings          []string
	PageCount         int
	FileSize          int64
	TempFilesCreated  []string
	LayoutStatistics  interface{}
}

// NewPDFPipeline creates a new PDF generation pipeline
func NewPDFPipeline(workDir, outputDir string) *PDFPipeline {
	tpl := NewTpl()
	return &PDFPipeline{
		workDir:         workDir,
		outputDir:       outputDir,
		templateEngine:  &tpl,
		layoutInteg:     NewLayoutIntegration(),
		visualOptimizer: NewVisualOptimizer(),
		logger:          DefaultLogger{},
	}
}

// SetLogger sets a custom logger for the pipeline
func (p *PDFPipeline) SetLogger(logger PDFLogger) {
	p.logger = logger
}

// GetVisualOptimizer returns the visual optimizer
func (p *PDFPipeline) GetVisualOptimizer() *VisualOptimizer {
	return p.visualOptimizer
}

// GeneratePDF generates a PDF from configuration with enhanced error handling
func (p *PDFPipeline) GeneratePDF(cfg config.Config, options PDFGenerationOptions) (*PDFGenerationResult, error) {
	startTime := time.Now()
	result := &PDFGenerationResult{
		TempFilesCreated: make([]string, 0),
		Errors:          make([]string, 0),
		Warnings:        make([]string, 0),
	}

	p.logger.Info("Starting PDF generation with enhanced layout integration")

	// Step 1: Validate configuration
	if err := p.validateConfig(cfg); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Configuration validation failed: %v", err))
		return result, err
	}

	// Step 2: Create temporary working directory
	tempDir, err := p.createTempWorkspace()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create temp workspace: %v", err))
		return result, err
	}
	result.TempFilesCreated = append(result.TempFilesCreated, tempDir)

	// Step 3: Generate enhanced modules with layout integration
	modules, err := p.layoutInteg.EnhancedMonthly(cfg, []string{"monthly_page.tpl"})
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Enhanced monthly generation failed: %v", err))
		// Fallback to legacy generation
		p.logger.Info("Falling back to legacy monthly generation")
		modules, err = MonthlyLegacy(cfg, []string{"monthly_page.tpl"})
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Legacy monthly generation also failed: %v", err))
			return result, err
		}
		result.Warnings = append(result.Warnings, "Used legacy generation due to layout integration failure")
	}

	// Step 4: Generate LaTeX document
	latexFile := filepath.Join(tempDir, "document.tex")
	if err := p.generateLaTeXDocument(cfg, modules, latexFile, options); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("LaTeX document generation failed: %v", err))
		return result, err
	}
	result.TempFilesCreated = append(result.TempFilesCreated, latexFile)

	// Step 5: Compile PDF with retries
	pdfFile, compilationLog, err := p.compilePDFWithRetries(latexFile, options)
	result.LaTeXLog = compilationLog
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("PDF compilation failed: %v", err))
		return result, err
	}

	// Step 6: Validate and move output
	finalPath := filepath.Join(p.outputDir, options.OutputFileName)
	if err := p.validateAndMoveOutput(pdfFile, finalPath, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Output validation failed: %v", err))
		return result, err
	}

	// Step 7: Cleanup if requested
	if options.CleanupTempFiles {
		p.cleanupTempFiles(result.TempFilesCreated)
	}

	result.Success = true
	result.OutputPath = finalPath
	result.CompilationTime = time.Since(startTime)

	p.logger.Info("PDF generation completed successfully in %v", result.CompilationTime)
	return result, nil
}

// validateConfig validates the configuration for PDF generation
func (p *PDFPipeline) validateConfig(cfg config.Config) error {
	if len(cfg.MonthsWithTasks) == 0 && len(cfg.GetYears()) == 0 {
		return fmt.Errorf("no months or years configured for generation")
	}

	if cfg.CSVFilePath != "" {
		if _, err := os.Stat(cfg.CSVFilePath); os.IsNotExist(err) {
			return fmt.Errorf("CSV file does not exist: %s", cfg.CSVFilePath)
		}
	}

	return nil
}

// createTempWorkspace creates a temporary directory for PDF generation
func (p *PDFPipeline) createTempWorkspace() (string, error) {
	tempDir, err := os.MkdirTemp(p.workDir, "pdf-gen-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	p.logger.Debug("Created temporary workspace: %s", tempDir)
	return tempDir, nil
}

// generateLaTeXDocument generates the complete LaTeX document
func (p *PDFPipeline) generateLaTeXDocument(cfg config.Config, modules config.Modules, outputFile string, options PDFGenerationOptions) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create LaTeX file: %w", err)
	}
	defer file.Close()

	// Write document preamble
	if err := p.writePreamble(file, cfg, options); err != nil {
		return fmt.Errorf("failed to write preamble: %w", err)
	}

	// Write document body
	if err := p.writeDocumentBody(file, cfg, modules); err != nil {
		return fmt.Errorf("failed to write document body: %w", err)
	}

	// Write document end
	if _, err := file.WriteString("\\end{document}\n"); err != nil {
		return fmt.Errorf("failed to write document end: %w", err)
	}

	p.logger.Debug("Generated LaTeX document: %s", outputFile)
	return nil
}

// writePreamble writes the LaTeX document preamble
func (p *PDFPipeline) writePreamble(w io.Writer, cfg config.Config, options PDFGenerationOptions) error {
	preamble := `\documentclass[11pt,letterpaper]{article}
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{geometry}
\usepackage{xcolor}
\usepackage{tikz}
\usepackage{tcolorbox}
\usepackage{tabularx}
\usepackage{array}
\usepackage{booktabs}
\usepackage{xparse}
\usepackage{expl3}

% Enhanced visual packages
\tcbuselibrary{skins}
\usetikzlibrary{shadows}

% Page geometry
\geometry{margin=0.5in}

% Color definitions for enhanced visual system
\definecolor{proposal}{RGB}{70,130,180}
\definecolor{laser}{RGB}{220,20,60}
\definecolor{imaging}{RGB}{34,139,34}
\definecolor{admin}{RGB}{128,128,128}
\definecolor{dissertation}{RGB}{138,43,226}
\definecolor{research}{RGB}{255,140,0}
\definecolor{publication}{RGB}{0,128,128}

`

	// Add custom preamble if provided
	if options.CustomPreamble != "" {
		preamble += options.CustomPreamble + "\n"
	}

	// Add extra packages
	for _, pkg := range options.ExtraPackages {
		preamble += fmt.Sprintf("\\usepackage{%s}\n", pkg)
	}

	preamble += "\n\\begin{document}\n"

	_, err := w.Write([]byte(preamble))
	return err
}

// writeDocumentBody writes the main document content
func (p *PDFPipeline) writeDocumentBody(w io.Writer, cfg config.Config, modules config.Modules) error {
	for i, module := range modules {
		p.logger.Debug("Processing module %d of %d", i+1, len(modules))

		// Use template engine to generate module content
		var buf bytes.Buffer
		if err := p.templateEngine.Execute(&buf, module.Tpl, module.Body); err != nil {
			// If template execution fails, write a simple fallback
			p.logger.Error("Template execution failed for module %d: %v", i, err)
			fallback := fmt.Sprintf("\\section*{Module %d - Template Error}\nTemplate execution failed: %v\n\\pagebreak\n", i+1, err)
			if _, err := w.Write([]byte(fallback)); err != nil {
				return fmt.Errorf("failed to write fallback content: %w", err)
			}
			continue
		}

		// Write the generated content
		if _, err := w.Write(buf.Bytes()); err != nil {
			return fmt.Errorf("failed to write module content: %w", err)
		}
	}

	return nil
}

// compilePDFWithRetries compiles the LaTeX document with retry logic
func (p *PDFPipeline) compilePDFWithRetries(latexFile string, options PDFGenerationOptions) (string, string, error) {
	engine := options.CompilationEngine
	if engine == "" {
		engine = "pdflatex"
	}

	maxRetries := options.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var lastErr error
	var compilationLog string

	for attempt := 1; attempt <= maxRetries; attempt++ {
		p.logger.Info("PDF compilation attempt %d of %d using %s", attempt, maxRetries, engine)

		pdfFile, log, err := p.compilePDF(latexFile, engine)
		compilationLog += fmt.Sprintf("=== Attempt %d ===\n%s\n", attempt, log)

		if err == nil {
			p.logger.Info("PDF compilation successful on attempt %d", attempt)
			return pdfFile, compilationLog, nil
		}

		lastErr = err
		p.logger.Error("Compilation attempt %d failed: %v", attempt, err)

		// Wait before retry (except on last attempt)
		if attempt < maxRetries {
			time.Sleep(time.Second * time.Duration(attempt))
		}
	}

	return "", compilationLog, fmt.Errorf("PDF compilation failed after %d attempts: %w", maxRetries, lastErr)
}

// compilePDF compiles a single LaTeX file to PDF
func (p *PDFPipeline) compilePDF(latexFile, engine string) (string, string, error) {
	dir := filepath.Dir(latexFile)
	basename := strings.TrimSuffix(filepath.Base(latexFile), ".tex")

	cmd := exec.Command(engine, "-interaction=nonstopmode", "-output-directory="+dir, latexFile)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	log := stdout.String() + stderr.String()

	pdfFile := filepath.Join(dir, basename+".pdf")
	if err != nil {
		return "", log, fmt.Errorf("compilation failed: %w", err)
	}

	// Check if PDF was actually created
	if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
		return "", log, fmt.Errorf("PDF file was not created despite successful compilation")
	}

	return pdfFile, log, nil
}

// validateAndMoveOutput validates the generated PDF and moves it to the final location
func (p *PDFPipeline) validateAndMoveOutput(sourcePDF, targetPath string, result *PDFGenerationResult) error {
	// Check if source PDF exists and get its size
	info, err := os.Stat(sourcePDF)
	if err != nil {
		return fmt.Errorf("source PDF not found: %w", err)
	}

	result.FileSize = info.Size()

	// Basic validation: file size should be reasonable
	if result.FileSize < 1024 { // Less than 1KB is suspicious
		result.Warnings = append(result.Warnings, "Generated PDF is very small, may be corrupted")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy the PDF to the target location
	if err := p.copyFile(sourcePDF, targetPath); err != nil {
		return fmt.Errorf("failed to copy PDF to output location: %w", err)
	}

	// Try to get page count (best effort)
	if pageCount, err := p.getPDFPageCount(targetPath); err == nil {
		result.PageCount = pageCount
	} else {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Could not determine page count: %v", err))
	}

	p.logger.Info("PDF validated and moved to: %s (size: %d bytes)", targetPath, result.FileSize)
	return nil
}

// copyFile copies a file from source to destination
func (p *PDFPipeline) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// getPDFPageCount attempts to get the page count of a PDF file
func (p *PDFPipeline) getPDFPageCount(pdfPath string) (int, error) {
	// Try using pdfinfo if available
	cmd := exec.Command("pdfinfo", pdfPath)
	output, err := cmd.Output()
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Pages:") {
				var pages int
				if _, err := fmt.Sscanf(line, "Pages: %d", &pages); err == nil {
					return pages, nil
				}
			}
		}
	}

	// Fallback: estimate based on file size (very rough)
	info, err := os.Stat(pdfPath)
	if err != nil {
		return 0, err
	}

	// Very rough estimation: ~50KB per page for typical calendar PDFs
	estimatedPages := int(info.Size() / 51200)
	if estimatedPages < 1 {
		estimatedPages = 1
	}

	return estimatedPages, nil
}

// cleanupTempFiles removes temporary files created during generation
func (p *PDFPipeline) cleanupTempFiles(files []string) {
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			p.logger.Error("Failed to cleanup temp file %s: %v", file, err)
		} else {
			p.logger.Debug("Cleaned up temp file: %s", file)
		}
	}
}
