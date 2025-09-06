package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"latex-yearly-planner/internal/config"
)

// MultiFormatGenerator handles generation of multiple output formats
type MultiFormatGenerator struct {
	pdfPipeline    *PDFPipeline
	workDir        string
	outputDir      string
	logger         PDFLogger
}

// MultiFormatOptions configures multi-format generation
type MultiFormatOptions struct {
	Formats        []OutputFormat
	ViewConfigs    []ViewConfig
	OutputPrefix   string
	IncludeStats   bool
	IncludeLegend  bool
	BatchMode      bool
	ParallelJobs   int
}

// MultiFormatResult contains results for multi-format generation
type MultiFormatResult struct {
	Success        bool
	TotalFiles     int
	SuccessfulFiles int
	FailedFiles    int
	Results        map[OutputFormat]map[string]*FormatResult
	Errors         []string
	Warnings       []string
	GenerationTime time.Duration
}

// FormatResult contains results for a specific format
type FormatResult struct {
	Format      OutputFormat
	ViewType    ViewType
	FilePath    string
	Success     bool
	FileSize    int64
	PageCount   int
	Error       string
	Warning     string
}

// NewMultiFormatGenerator creates a new multi-format generator
func NewMultiFormatGenerator(workDir, outputDir string) *MultiFormatGenerator {
	return &MultiFormatGenerator{
		pdfPipeline: NewPDFPipeline(workDir, outputDir),
		workDir:     workDir,
		outputDir:   outputDir,
		logger:      DefaultLogger{},
	}
}

// SetLogger sets a custom logger for the generator
func (mfg *MultiFormatGenerator) SetLogger(logger PDFLogger) {
	mfg.logger = logger
	mfg.pdfPipeline.SetLogger(logger)
}

// GenerateMultiFormat generates multiple output formats
func (mfg *MultiFormatGenerator) GenerateMultiFormat(cfg config.Config, options MultiFormatOptions) (*MultiFormatResult, error) {
	startTime := time.Now()
	result := &MultiFormatResult{
		Results: make(map[OutputFormat]map[string]*FormatResult),
		Errors:  make([]string, 0),
		Warnings: make([]string, 0),
	}

	mfg.logger.Info("Starting multi-format generation with %d formats", len(options.Formats))

	// Validate options
	if err := mfg.validateMultiFormatOptions(options); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Options validation failed: %v", err))
		return result, err
	}

	// Generate for each format
	for _, format := range options.Formats {
		mfg.logger.Info("Generating format: %s", format)
		
		formatResults := make(map[string]*FormatResult)
		
		for _, viewConfig := range options.ViewConfigs {
			formatResult, err := mfg.generateFormat(cfg, format, viewConfig, options)
			if err != nil {
				mfg.logger.Error("Failed to generate %s for view %s: %v", format, viewConfig.Type, err)
				formatResults[string(viewConfig.Type)] = &FormatResult{
					Format:  format,
					ViewType: viewConfig.Type,
					Success: false,
					Error:   err.Error(),
				}
				result.FailedFiles++
			} else {
				formatResults[string(viewConfig.Type)] = formatResult
				if formatResult.Success {
					result.SuccessfulFiles++
				} else {
					result.FailedFiles++
				}
			}
		}
		
		result.Results[format] = formatResults
	}

	result.Success = result.FailedFiles == 0
	result.TotalFiles = result.SuccessfulFiles + result.FailedFiles
	result.GenerationTime = time.Since(startTime)

	mfg.logger.Info("Multi-format generation completed: %d successful, %d failed", 
		result.SuccessfulFiles, result.FailedFiles)

	return result, nil
}

// generateFormat generates a specific format for a view
func (mfg *MultiFormatGenerator) generateFormat(cfg config.Config, format OutputFormat, viewConfig ViewConfig, options MultiFormatOptions) (*FormatResult, error) {
	result := &FormatResult{
		Format:   format,
		ViewType: viewConfig.Type,
	}

	// Generate filename
	filename := mfg.generateFilename(options.OutputPrefix, format, viewConfig)
	filePath := filepath.Join(mfg.outputDir, filename)

	switch format {
	case OutputFormatPDF:
		return mfg.generatePDF(cfg, viewConfig, filePath, result)
	case OutputFormatLaTeX:
		return mfg.generateLaTeX(cfg, viewConfig, filePath, result)
	case OutputFormatHTML:
		return mfg.generateHTML(cfg, viewConfig, filePath, result)
	case OutputFormatSVG:
		return mfg.generateSVG(cfg, viewConfig, filePath, result)
	case OutputFormatPNG:
		return mfg.generatePNG(cfg, viewConfig, filePath, result)
	default:
		return result, fmt.Errorf("unsupported output format: %s", format)
	}
}

// generatePDF generates PDF output
func (mfg *MultiFormatGenerator) generatePDF(cfg config.Config, viewConfig ViewConfig, filePath string, result *FormatResult) (*FormatResult, error) {
	pdfOptions := PDFGenerationOptions{
		OutputFileName:   filepath.Base(filePath),
		CleanupTempFiles: true,
		MaxRetries:       3,
		CompilationEngine: "pdflatex",
	}

	// Apply view configuration
	ApplyViewConfig(viewConfig, &pdfOptions)

	// Generate PDF
	pdfResult, err := mfg.pdfPipeline.GeneratePDF(cfg, pdfOptions)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Success = pdfResult.Success
	result.FilePath = pdfResult.OutputPath
	result.FileSize = pdfResult.FileSize
	result.PageCount = pdfResult.PageCount

	if len(pdfResult.Warnings) > 0 {
		result.Warning = strings.Join(pdfResult.Warnings, "; ")
	}

	return result, nil
}

// generateLaTeX generates LaTeX output
func (mfg *MultiFormatGenerator) generateLaTeX(cfg config.Config, viewConfig ViewConfig, filePath string, result *FormatResult) (*FormatResult, error) {
	// Generate LaTeX content
	latexContent, err := mfg.generateLaTeXContent(cfg, viewConfig)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Write to file
	if err := os.WriteFile(filePath, []byte(latexContent), 0644); err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Success = true
	result.FilePath = filePath
	result.FileSize = info.Size()

	return result, nil
}

// generateHTML generates HTML output
func (mfg *MultiFormatGenerator) generateHTML(cfg config.Config, viewConfig ViewConfig, filePath string, result *FormatResult) (*FormatResult, error) {
	// Generate HTML content
	htmlContent, err := mfg.generateHTMLContent(cfg, viewConfig)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Write to file
	if err := os.WriteFile(filePath, []byte(htmlContent), 0644); err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Success = true
	result.FilePath = filePath
	result.FileSize = info.Size()

	return result, nil
}

// generateSVG generates SVG output
func (mfg *MultiFormatGenerator) generateSVG(cfg config.Config, viewConfig ViewConfig, filePath string, result *FormatResult) (*FormatResult, error) {
	// First generate LaTeX
	latexContent, err := mfg.generateLaTeXContent(cfg, viewConfig)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Convert LaTeX to SVG using dvisvgm
	tempDir, err := os.MkdirTemp(mfg.workDir, "svg-gen-*")
	if err != nil {
		result.Error = err.Error()
		return result, err
	}
	defer os.RemoveAll(tempDir)

	// Write LaTeX file
	latexFile := filepath.Join(tempDir, "document.tex")
	if err := os.WriteFile(latexFile, []byte(latexContent), 0644); err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Compile to DVI
	dviFile := filepath.Join(tempDir, "document.dvi")
	cmd := exec.Command("pdflatex", "-output-format=dvi", "-interaction=nonstopmode", "-output-directory="+tempDir, latexFile)
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("LaTeX compilation failed: %v", err)
		return result, err
	}

	// Convert DVI to SVG
	cmd = exec.Command("dvisvgm", "--no-fonts", "--exact", dviFile)
	output, err := cmd.Output()
	if err != nil {
		result.Error = fmt.Sprintf("SVG conversion failed: %v", err)
		return result, err
	}

	// Write SVG file
	if err := os.WriteFile(filePath, output, 0644); err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Success = true
	result.FilePath = filePath
	result.FileSize = info.Size()

	return result, nil
}

// generatePNG generates PNG output
func (mfg *MultiFormatGenerator) generatePNG(cfg config.Config, viewConfig ViewConfig, filePath string, result *FormatResult) (*FormatResult, error) {
	// First generate PDF
	pdfOptions := PDFGenerationOptions{
		OutputFileName:   "temp.pdf",
		CleanupTempFiles: false,
		MaxRetries:       3,
		CompilationEngine: "pdflatex",
	}

	ApplyViewConfig(viewConfig, &pdfOptions)

	pdfResult, err := mfg.pdfPipeline.GeneratePDF(cfg, pdfOptions)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Convert PDF to PNG using ImageMagick
	cmd := exec.Command("convert", "-density", "300", "-quality", "100", pdfResult.OutputPath, filePath)
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("PNG conversion failed: %v", err)
		return result, err
	}

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Success = true
	result.FilePath = filePath
	result.FileSize = info.Size()

	return result, nil
}

// generateLaTeXContent generates LaTeX content for a view
func (mfg *MultiFormatGenerator) generateLaTeXContent(cfg config.Config, viewConfig ViewConfig) (string, error) {
	// Create layout integration
	layoutIntegration := NewLayoutIntegration()

	// Generate enhanced modules
	modules, err := layoutIntegration.EnhancedMonthly(cfg, []string{viewConfig.TemplateName})
	if err != nil {
		// Fallback to legacy generation
		modules, err = MonthlyLegacy(cfg, []string{viewConfig.TemplateName})
		if err != nil {
			return "", fmt.Errorf("failed to generate modules: %w", err)
		}
	}

	// Generate LaTeX document
	var buf bytes.Buffer

	// Write preamble
	if err := mfg.writeLaTeXPreamble(&buf, viewConfig); err != nil {
		return "", fmt.Errorf("failed to write preamble: %w", err)
	}

	// Write document body
	for _, module := range modules {
		// Use template engine to generate module content
		var moduleBuf bytes.Buffer
		tpl := NewTpl()
		if err := tpl.Execute(&moduleBuf, module.Tpl, module.Body); err != nil {
			// Write fallback content
			fallback := fmt.Sprintf("\\section*{Module - Template Error}\nTemplate execution failed: %v\n\\pagebreak\n", err)
			buf.WriteString(fallback)
			continue
		}
		buf.Write(moduleBuf.Bytes())
	}

	// Write document end
	buf.WriteString("\\end{document}\n")

	return buf.String(), nil
}

// generateHTMLContent generates HTML content for a view
func (mfg *MultiFormatGenerator) generateHTMLContent(cfg config.Config, viewConfig ViewConfig) (string, error) {
	// Create layout integration
	layoutIntegration := NewLayoutIntegration()

	// Generate enhanced modules
	modules, err := layoutIntegration.EnhancedMonthly(cfg, []string{viewConfig.TemplateName})
	if err != nil {
		// Fallback to legacy generation
		modules, err = MonthlyLegacy(cfg, []string{viewConfig.TemplateName})
		if err != nil {
			return "", fmt.Errorf("failed to generate modules: %w", err)
		}
	}

	// Generate HTML template
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .calendar { border-collapse: collapse; width: 100%; }
        .calendar th, .calendar td { border: 1px solid #ddd; padding: 8px; text-align: center; }
        .calendar th { background-color: #f2f2f2; }
        .task-bar { background-color: #4CAF50; color: white; padding: 2px 4px; margin: 1px 0; border-radius: 2px; font-size: 12px; }
        .task-bar.proposal { background-color: #2196F3; }
        .task-bar.laser { background-color: #f44336; }
        .task-bar.imaging { background-color: #4CAF50; }
        .task-bar.admin { background-color: #9E9E9E; }
        .task-bar.dissertation { background-color: #9C27B0; }
        .task-bar.research { background-color: #FF9800; }
        .task-bar.publication { background-color: #00BCD4; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <p>{{.Description}}</p>
    
    {{range .Modules}}
    <div class="module">
        <h2>{{.Title}}</h2>
        <!-- Calendar content would go here -->
    </div>
    {{end}}
</body>
</html>`

	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var buf bytes.Buffer
	data := struct {
		Title       string
		Description string
		Modules     []config.Module
	}{
		Title:       viewConfig.Title,
		Description: viewConfig.Description,
		Modules:     modules,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return buf.String(), nil
}

// writeLaTeXPreamble writes the LaTeX document preamble
func (mfg *MultiFormatGenerator) writeLaTeXPreamble(w *bytes.Buffer, viewConfig ViewConfig) error {
	preamble := fmt.Sprintf(`\documentclass[11pt,%s]{article}
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage[%s,%s]{geometry}
\usepackage{xcolor}
\usepackage{tikz}
\usepackage{tcolorbox}
\usepackage{tabularx}
\usepackage{array}
\usepackage{booktabs}
\usepackage{xparse}
\usepackage{expl3}

%% Enhanced visual packages
\\tcbuselibrary{skins}
\\usetikzlibrary{shadows}

%% View-specific configuration
%s

\\begin{document}
`, viewConfig.PageSize, viewConfig.Orientation, viewConfig.PageSize, generateViewLaTeX(viewConfig))

	_, err := w.WriteString(preamble)
	return err
}

// generateFilename generates a filename for a format and view
func (mfg *MultiFormatGenerator) generateFilename(prefix string, format OutputFormat, viewConfig ViewConfig) string {
	if prefix == "" {
		prefix = "planner"
	}
	
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-%s-%s.%s", prefix, string(viewConfig.Type), timestamp, format)
}

// validateMultiFormatOptions validates multi-format generation options
func (mfg *MultiFormatGenerator) validateMultiFormatOptions(options MultiFormatOptions) error {
	if len(options.Formats) == 0 {
		return fmt.Errorf("at least one output format must be specified")
	}

	if len(options.ViewConfigs) == 0 {
		return fmt.Errorf("at least one view configuration must be specified")
	}

	// Validate formats
	validFormats := map[OutputFormat]bool{
		OutputFormatPDF:   true,
		OutputFormatLaTeX: true,
		OutputFormatHTML:  true,
		OutputFormatSVG:   true,
		OutputFormatPNG:   true,
	}

	for _, format := range options.Formats {
		if !validFormats[format] {
			return fmt.Errorf("invalid output format: %s", format)
		}
	}

	// Validate view configs
	for _, viewConfig := range options.ViewConfigs {
		if err := ValidateViewConfig(viewConfig); err != nil {
			return fmt.Errorf("invalid view configuration: %w", err)
		}
	}

	return nil
}
