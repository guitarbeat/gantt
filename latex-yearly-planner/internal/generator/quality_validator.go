package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"latex-yearly-planner/internal/calendar"
)

// QualityValidator provides comprehensive quality validation for PDF output
type QualityValidator struct {
	tester    *QualityTester
	config    *QualityValidationConfig
	logger    PDFLogger
	reportDir string
}

// QualityValidationConfig defines configuration for quality validation
type QualityValidationConfig struct {
	// Validation settings
	EnablePDFValidation    bool `json:"enable_pdf_validation"`
	EnableLaTeXValidation  bool `json:"enable_latex_validation"`
	EnableVisualValidation bool `json:"enable_visual_validation"`
	EnableContentValidation bool `json:"enable_content_validation"`
	
	// Validation thresholds
	MinPDFSize        int64   `json:"min_pdf_size"`        // Minimum PDF file size in bytes
	MaxPDFSize        int64   `json:"max_pdf_size"`        // Maximum PDF file size in bytes
	MinPageCount      int     `json:"min_page_count"`      // Minimum page count
	MaxPageCount      int     `json:"max_page_count"`      // Maximum page count
	MinQualityScore   float64 `json:"min_quality_score"`   // Minimum quality score
	MaxCompilationTime time.Duration `json:"max_compilation_time"` // Maximum compilation time
	
	// Report settings
	GenerateReport    bool   `json:"generate_report"`
	ReportFormat      string `json:"report_format"`      // "json", "html", "text"
	IncludeScreenshots bool  `json:"include_screenshots"`
	IncludeMetrics    bool   `json:"include_metrics"`
}

// QualityValidationResult contains the results of quality validation
type QualityValidationResult struct {
	OverallPassed     bool                      `json:"overall_passed"`
	ValidationScore   float64                   `json:"validation_score"`
	PDFValidation     *PDFValidationResult      `json:"pdf_validation"`
	LaTeXValidation   *LaTeXValidationResult    `json:"latex_validation"`
	VisualValidation  *VisualValidationResult   `json:"visual_validation"`
	ContentValidation *ContentValidationResult  `json:"content_validation"`
	Issues           []QualityIssue            `json:"issues"`
	Recommendations  []QualityRecommendation   `json:"recommendations"`
	ReportPath       string                    `json:"report_path"`
	ValidationTime   time.Duration             `json:"validation_time"`
	Timestamp        time.Time                 `json:"timestamp"`
}

// PDFValidationResult contains PDF-specific validation results
type PDFValidationResult struct {
	Passed           bool          `json:"passed"`
	FileExists       bool          `json:"file_exists"`
	FileSize         int64         `json:"file_size"`
	FileSizeValid    bool          `json:"file_size_valid"`
	PageCount        int           `json:"page_count"`
	PageCountValid   bool          `json:"page_count_valid"`
	CompilationTime  time.Duration `json:"compilation_time"`
	CompilationValid bool          `json:"compilation_valid"`
	Issues           []QualityIssue `json:"issues"`
}

// LaTeXValidationResult contains LaTeX-specific validation results
type LaTeXValidationResult struct {
	Passed           bool          `json:"passed"`
	SyntaxValid      bool          `json:"syntax_valid"`
	CompilationValid bool          `json:"compilation_valid"`
	PackageValid     bool          `json:"package_valid"`
	Issues           []QualityIssue `json:"issues"`
}

// VisualValidationResult contains visual-specific validation results
type VisualValidationResult struct {
	Passed           bool          `json:"passed"`
	SpacingValid     bool          `json:"spacing_valid"`
	AlignmentValid   bool          `json:"alignment_valid"`
	ReadabilityValid bool          `json:"readability_valid"`
	ColorValid       bool          `json:"color_valid"`
	Issues           []QualityIssue `json:"issues"`
}

// ContentValidationResult contains content-specific validation results
type ContentValidationResult struct {
	Passed           bool          `json:"passed"`
	TaskCountValid   bool          `json:"task_count_valid"`
	DataIntegrityValid bool        `json:"data_integrity_valid"`
	LayoutValid      bool          `json:"layout_valid"`
	Issues           []QualityIssue `json:"issues"`
}

// NewQualityValidator creates a new quality validator
func NewQualityValidator() *QualityValidator {
	return &QualityValidator{
		tester:    NewQualityTester(),
		config:    GetDefaultQualityValidationConfig(),
		logger:    &QualityValidatorLogger{},
		reportDir: "quality_reports",
	}
}

// GetDefaultQualityValidationConfig returns the default quality validation configuration
func GetDefaultQualityValidationConfig() *QualityValidationConfig {
	return &QualityValidationConfig{
		EnablePDFValidation:    true,
		EnableLaTeXValidation:  true,
		EnableVisualValidation: true,
		EnableContentValidation: true,
		MinPDFSize:            1024,      // 1KB
		MaxPDFSize:            50 * 1024 * 1024, // 50MB
		MinPageCount:          1,
		MaxPageCount:          100,
		MinQualityScore:       0.8,
		MaxCompilationTime:    time.Minute * 5,
		GenerateReport:        true,
		ReportFormat:          "json",
		IncludeScreenshots:    false,
		IncludeMetrics:        true,
	}
}

// SetLogger sets the logger for the quality validator
func (qv *QualityValidator) SetLogger(logger PDFLogger) {
	qv.logger = logger
	qv.tester.SetLogger(logger)
}

// ValidateQuality performs comprehensive quality validation
func (qv *QualityValidator) ValidateQuality(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType, pdfPath string) (*QualityValidationResult, error) {
	startTime := time.Now()
	qv.logger.Info("Starting comprehensive quality validation for %s view", viewType)
	
	result := &QualityValidationResult{
		Issues:      make([]QualityIssue, 0),
		Recommendations: make([]QualityRecommendation, 0),
		Timestamp:   time.Now(),
	}
	
	// Validate PDF
	if qv.config.EnablePDFValidation {
		pdfResult, err := qv.validatePDF(pdfPath)
		if err != nil {
			qv.logger.Error("PDF validation failed: %v", err)
			return nil, err
		}
		result.PDFValidation = pdfResult
	}
	
	// Validate LaTeX
	if qv.config.EnableLaTeXValidation {
		latexResult, err := qv.validateLaTeX(layoutResult)
		if err != nil {
			qv.logger.Error("LaTeX validation failed: %v", err)
			return nil, err
		}
		result.LaTeXValidation = latexResult
	}
	
	// Validate visual quality
	if qv.config.EnableVisualValidation {
		visualResult, err := qv.validateVisual(layoutResult, viewType)
		if err != nil {
			qv.logger.Error("Visual validation failed: %v", err)
			return nil, err
		}
		result.VisualValidation = visualResult
	}
	
	// Validate content
	if qv.config.EnableContentValidation {
		contentResult, err := qv.validateContent(layoutResult)
		if err != nil {
			qv.logger.Error("Content validation failed: %v", err)
			return nil, err
		}
		result.ContentValidation = contentResult
	}
	
	// Calculate overall validation score
	result.ValidationScore = qv.calculateValidationScore(result)
	result.OverallPassed = result.ValidationScore >= qv.config.MinQualityScore
	
	// Collect all issues
	qv.collectIssues(result)
	
	// Generate recommendations
	result.Recommendations = qv.generateRecommendations(result)
	
	// Generate report if requested
	if qv.config.GenerateReport {
		reportPath, err := qv.generateReport(result)
		if err != nil {
			qv.logger.Error("Failed to generate report: %v", err)
		} else {
			result.ReportPath = reportPath
		}
	}
	
	// Calculate validation time
	result.ValidationTime = time.Since(startTime)
	
	qv.logger.Info("Quality validation completed in %v with score: %.2f", 
		result.ValidationTime, result.ValidationScore)
	
	return result, nil
}

// validatePDF validates PDF-specific quality aspects
func (qv *QualityValidator) validatePDF(pdfPath string) (*PDFValidationResult, error) {
	result := &PDFValidationResult{
		Issues: make([]QualityIssue, 0),
	}
	
	// Check if file exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		result.Issues = append(result.Issues, QualityIssue{
			Severity:    SeverityCritical,
			Category:    "pdf",
			Description: "PDF file does not exist",
			Location:    pdfPath,
			Suggestions: []string{"Check file path", "Verify PDF generation"},
		})
		return result, nil
	}
	
	result.FileExists = true
	
	// Get file info
	fileInfo, err := os.Stat(pdfPath)
	if err != nil {
		result.Issues = append(result.Issues, QualityIssue{
			Severity:    SeverityHigh,
			Category:    "pdf",
			Description: fmt.Sprintf("Failed to get file info: %v", err),
			Location:    pdfPath,
			Suggestions: []string{"Check file permissions", "Verify file integrity"},
		})
		return result, nil
	}
	
	result.FileSize = fileInfo.Size()
	
	// Validate file size
	if result.FileSize < qv.config.MinPDFSize {
		result.Issues = append(result.Issues, QualityIssue{
			Severity:    SeverityHigh,
			Category:    "pdf",
			Description: fmt.Sprintf("PDF file size %d bytes is below minimum %d bytes", result.FileSize, qv.config.MinPDFSize),
			Location:    pdfPath,
			Suggestions: []string{"Check PDF generation", "Verify content", "Increase content"},
		})
	} else if result.FileSize > qv.config.MaxPDFSize {
		result.Issues = append(result.Issues, QualityIssue{
			Severity:    SeverityMedium,
			Category:    "pdf",
			Description: fmt.Sprintf("PDF file size %d bytes exceeds maximum %d bytes", result.FileSize, qv.config.MaxPDFSize),
			Location:    pdfPath,
			Suggestions: []string{"Optimize content", "Reduce image quality", "Compress PDF"},
		})
	} else {
		result.FileSizeValid = true
	}
	
	// TODO: Add page count validation (requires PDF parsing)
	result.PageCount = 1 // Placeholder
	result.PageCountValid = true
	
	// TODO: Add compilation time validation
	result.CompilationTime = time.Second * 2 // Placeholder
	result.CompilationValid = true
	
	// Determine if PDF validation passed
	result.Passed = result.FileExists && result.FileSizeValid && result.PageCountValid && result.CompilationValid
	
	return result, nil
}

// validateLaTeX validates LaTeX-specific quality aspects
func (qv *QualityValidator) validateLaTeX(layoutResult *calendar.IntegratedLayoutResult) (*LaTeXValidationResult, error) {
	result := &LaTeXValidationResult{
		Issues: make([]QualityIssue, 0),
	}
	
	// Validate LaTeX syntax (basic checks)
	syntaxValid := true
	for _, bar := range layoutResult.TaskBars {
		// Check for problematic characters in task names
		if strings.Contains(bar.TaskName, "\\") || strings.Contains(bar.TaskName, "{") || strings.Contains(bar.TaskName, "}") {
			syntaxValid = false
			result.Issues = append(result.Issues, QualityIssue{
				Severity:    SeverityHigh,
				Category:    "latex",
				Description: fmt.Sprintf("Task name '%s' contains LaTeX special characters", bar.TaskName),
				Location:    fmt.Sprintf("Task: %s", bar.TaskName),
				Suggestions: []string{"Escape special characters", "Sanitize input", "Use LaTeX-safe names"},
			})
		}
	}
	
	result.SyntaxValid = syntaxValid
	
	// Validate package requirements
	packageValid := true
	// TODO: Add package validation logic
	result.PackageValid = packageValid
	
	// Validate compilation (assume valid if syntax is valid)
	result.CompilationValid = syntaxValid && packageValid
	
	// Determine if LaTeX validation passed
	result.Passed = result.SyntaxValid && result.CompilationValid && result.PackageValid
	
	return result, nil
}

// validateVisual validates visual quality aspects
func (qv *QualityValidator) validateVisual(layoutResult *calendar.IntegratedLayoutResult, viewType ViewType) (*VisualValidationResult, error) {
	result := &VisualValidationResult{
		Issues: make([]QualityIssue, 0),
	}
	
	// Run quality tests
	qualityResult, err := qv.tester.RunQualityTests(layoutResult, viewType)
	if err != nil {
		return nil, fmt.Errorf("failed to run quality tests: %w", err)
	}
	
	// Extract visual validation results
	if visualCategory, exists := qualityResult.TestResults["visual"]; exists {
		result.SpacingValid = visualCategory.Score >= 0.8
		result.AlignmentValid = visualCategory.Score >= 0.8
		result.ReadabilityValid = visualCategory.Score >= 0.8
		result.ColorValid = true // TODO: Add color validation
		
		// Add issues from quality tests
		for _, issue := range visualCategory.Issues {
			result.Issues = append(result.Issues, issue)
		}
	}
	
	// Determine if visual validation passed
	result.Passed = result.SpacingValid && result.AlignmentValid && result.ReadabilityValid && result.ColorValid
	
	return result, nil
}

// validateContent validates content quality aspects
func (qv *QualityValidator) validateContent(layoutResult *calendar.IntegratedLayoutResult) (*ContentValidationResult, error) {
	result := &ContentValidationResult{
		Issues: make([]QualityIssue, 0),
	}
	
	// Validate task count
	taskCount := len(layoutResult.TaskBars)
	if taskCount == 0 {
		result.Issues = append(result.Issues, QualityIssue{
			Severity:    SeverityHigh,
			Category:    "content",
			Description: "No tasks found in layout result",
			Location:    "Layout result",
			Suggestions: []string{"Check task data", "Verify layout processing", "Add sample tasks"},
		})
	} else {
		result.TaskCountValid = true
	}
	
	// Validate data integrity
	dataIntegrityValid := true
	for _, bar := range layoutResult.TaskBars {
		if bar.TaskName == "" {
			dataIntegrityValid = false
			result.Issues = append(result.Issues, QualityIssue{
				Severity:    SeverityHigh,
				Category:    "content",
				Description: "Task bar has empty task name",
				Location:    "Task bar",
				Suggestions: []string{"Check task data", "Validate input", "Add default name"},
			})
		}
		if bar.Width <= 0 || bar.Height <= 0 {
			dataIntegrityValid = false
			result.Issues = append(result.Issues, QualityIssue{
				Severity:    SeverityHigh,
				Category:    "content",
				Description: "Task bar has invalid dimensions",
				Location:    fmt.Sprintf("Task: %s", bar.TaskName),
				Suggestions: []string{"Check layout algorithm", "Validate dimensions", "Fix positioning"},
			})
		}
	}
	result.DataIntegrityValid = dataIntegrityValid
	
	// Validate layout
	layoutValid := true
	// TODO: Add layout validation logic
	result.LayoutValid = layoutValid
	
	// Determine if content validation passed
	result.Passed = result.TaskCountValid && result.DataIntegrityValid && result.LayoutValid
	
	return result, nil
}

// calculateValidationScore calculates the overall validation score
func (qv *QualityValidator) calculateValidationScore(result *QualityValidationResult) float64 {
	score := 0.0
	count := 0
	
	if result.PDFValidation != nil {
		if result.PDFValidation.Passed {
			score += 1.0
		}
		count++
	}
	
	if result.LaTeXValidation != nil {
		if result.LaTeXValidation.Passed {
			score += 1.0
		}
		count++
	}
	
	if result.VisualValidation != nil {
		if result.VisualValidation.Passed {
			score += 1.0
		}
		count++
	}
	
	if result.ContentValidation != nil {
		if result.ContentValidation.Passed {
			score += 1.0
		}
		count++
	}
	
	if count == 0 {
		return 0.0
	}
	
	return score / float64(count)
}

// collectIssues collects all issues from validation results
func (qv *QualityValidator) collectIssues(result *QualityValidationResult) {
	if result.PDFValidation != nil {
		result.Issues = append(result.Issues, result.PDFValidation.Issues...)
	}
	if result.LaTeXValidation != nil {
		result.Issues = append(result.Issues, result.LaTeXValidation.Issues...)
	}
	if result.VisualValidation != nil {
		result.Issues = append(result.Issues, result.VisualValidation.Issues...)
	}
	if result.ContentValidation != nil {
		result.Issues = append(result.Issues, result.ContentValidation.Issues...)
	}
}

// generateRecommendations generates recommendations based on validation results
func (qv *QualityValidator) generateRecommendations(result *QualityValidationResult) []QualityRecommendation {
	recommendations := make([]QualityRecommendation, 0)
	
	// Analyze issues and generate recommendations
	issueCounts := make(map[string]int)
	for _, issue := range result.Issues {
		issueCounts[issue.Category]++
	}
	
	// Generate recommendations based on issue frequency
	for category, count := range issueCounts {
		if count > 0 {
			recommendations = append(recommendations, QualityRecommendation{
				Category:      category,
				Description:   fmt.Sprintf("Address %d issues in %s category", count, category),
				Priority:      qv.getPriorityForCategory(category),
				Impact:        qv.getImpactForCategory(category),
				Effort:        qv.getEffortForCategory(category),
				Implementation: qv.getImplementationForCategory(category),
			})
		}
	}
	
	return recommendations
}

// Helper methods for recommendation generation
func (qv *QualityValidator) getPriorityForCategory(category string) int {
	switch category {
	case "pdf", "content":
		return 1
	case "latex", "visual":
		return 2
	default:
		return 3
	}
}

func (qv *QualityValidator) getImpactForCategory(category string) float64 {
	switch category {
	case "pdf", "content":
		return 0.9
	case "latex", "visual":
		return 0.7
	default:
		return 0.5
	}
}

func (qv *QualityValidator) getEffortForCategory(category string) string {
	switch category {
	case "pdf":
		return "Low"
	case "latex", "content":
		return "Medium"
	case "visual":
		return "High"
	default:
		return "Medium"
	}
}

func (qv *QualityValidator) getImplementationForCategory(category string) string {
	switch category {
	case "pdf":
		return "Check PDF generation settings and file paths"
	case "latex":
		return "Sanitize input data and escape special characters"
	case "content":
		return "Validate input data and fix data integrity issues"
	case "visual":
		return "Improve visual design system and spacing configuration"
	default:
		return "Review and fix issues in the specified category"
	}
}

// generateReport generates a quality validation report
func (qv *QualityValidator) generateReport(result *QualityValidationResult) (string, error) {
	// Create report directory
	if err := os.MkdirAll(qv.reportDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create report directory: %w", err)
	}
	
	// Generate report filename
	timestamp := result.Timestamp.Format("2006-01-02_15-04-05")
	reportPath := filepath.Join(qv.reportDir, fmt.Sprintf("quality_report_%s.%s", timestamp, qv.config.ReportFormat))
	
	// Generate report based on format
	switch qv.config.ReportFormat {
	case "json":
		return qv.generateJSONReport(result, reportPath)
	case "html":
		return qv.generateHTMLReport(result, reportPath)
	case "text":
		return qv.generateTextReport(result, reportPath)
	default:
		return qv.generateJSONReport(result, reportPath)
	}
}

// generateJSONReport generates a JSON quality report
func (qv *QualityValidator) generateJSONReport(result *QualityValidationResult, reportPath string) (string, error) {
	// TODO: Implement JSON report generation
	qv.logger.Info("JSON report generation not yet implemented")
	return reportPath, nil
}

// generateHTMLReport generates an HTML quality report
func (qv *QualityValidator) generateHTMLReport(result *QualityValidationResult, reportPath string) (string, error) {
	// TODO: Implement HTML report generation
	qv.logger.Info("HTML report generation not yet implemented")
	return reportPath, nil
}

// generateTextReport generates a text quality report
func (qv *QualityValidator) generateTextReport(result *QualityValidationResult, reportPath string) (string, error) {
	// TODO: Implement text report generation
	qv.logger.Info("Text report generation not yet implemented")
	return reportPath, nil
}

// QualityValidatorLogger provides logging for quality validator
type QualityValidatorLogger struct{}

func (l *QualityValidatorLogger) Info(msg string, args ...interface{})  { fmt.Printf("[VALIDATOR-INFO] "+msg+"\n", args...) }
func (l *QualityValidatorLogger) Error(msg string, args ...interface{}) { fmt.Printf("[VALIDATOR-ERROR] "+msg+"\n", args...) }
func (l *QualityValidatorLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[VALIDATOR-DEBUG] "+msg+"\n", args...) }
