package data

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ValidationReporter handles comprehensive error reporting and validation feedback
type ValidationReporter struct {
	errors     []DataValidationError
	warnings   []DataValidationError
	info       []DataValidationError
	timestamp  time.Time
	taskCount  int
	config     *ValidationReportConfig
}

// ValidationReportConfig configures the validation report
type ValidationReportConfig struct {
	IncludeSuggestions bool
	IncludeTimestamps  bool
	IncludeSeverity    bool
	GroupByType        bool
	GroupBySeverity    bool
	GroupByTask        bool
	SortBySeverity     bool
	MaxErrors          int // 0 = no limit
	MaxWarnings        int // 0 = no limit
	MaxInfo            int // 0 = no limit
}

// DefaultValidationReportConfig returns the default configuration
func DefaultValidationReportConfig() *ValidationReportConfig {
	return &ValidationReportConfig{
		IncludeSuggestions: true,
		IncludeTimestamps:  false,
		IncludeSeverity:    true,
		GroupByType:        true,
		GroupBySeverity:    true,
		GroupByTask:        false,
		SortBySeverity:     true,
		MaxErrors:          0,
		MaxWarnings:        0,
		MaxInfo:            0,
	}
}

// NewValidationReporter creates a new validation reporter
func NewValidationReporter() *ValidationReporter {
	return &ValidationReporter{
		errors:    make([]DataValidationError, 0),
		warnings:  make([]DataValidationError, 0),
		info:      make([]DataValidationError, 0),
		timestamp: time.Now(),
		config:    DefaultValidationReportConfig(),
	}
}

// AddValidationResult adds a validation result to the reporter
func (vr *ValidationReporter) AddValidationResult(result *ValidationResult) {
	vr.errors = append(vr.errors, result.Errors...)
	vr.warnings = append(vr.warnings, result.Warnings...)
	vr.info = append(vr.info, result.Info...)
	vr.taskCount = result.TaskCount
}

// AddErrors adds individual errors to the reporter
func (vr *ValidationReporter) AddErrors(errors []DataValidationError) {
	vr.errors = append(vr.errors, errors...)
}

// AddWarnings adds individual warnings to the reporter
func (vr *ValidationReporter) AddWarnings(warnings []DataValidationError) {
	vr.warnings = append(vr.warnings, warnings...)
}

// AddInfo adds individual info messages to the reporter
func (vr *ValidationReporter) AddInfo(info []DataValidationError) {
	vr.info = append(vr.info, info...)
}

// SetConfig sets the validation report configuration
func (vr *ValidationReporter) SetConfig(config *ValidationReportConfig) {
	vr.config = config
}

// GenerateReport generates a comprehensive validation report
func (vr *ValidationReporter) GenerateReport() *ValidationReport {
	report := &ValidationReport{
		Summary:     vr.generateSummary(),
		Errors:      vr.filterAndSortErrors(vr.errors),
		Warnings:    vr.filterAndSortErrors(vr.warnings),
		Info:        vr.filterAndSortErrors(vr.info),
		TaskCount:   vr.taskCount,
		ErrorCount:  len(vr.errors),
		WarningCount: len(vr.warnings),
		InfoCount:   len(vr.info),
		Timestamp:   vr.timestamp,
		IsValid:     len(vr.errors) == 0,
	}
	
	// Generate detailed sections
	report.ErrorSummary = vr.generateErrorSummary(vr.errors)
	report.WarningSummary = vr.generateErrorSummary(vr.warnings)
	report.InfoSummary = vr.generateErrorSummary(vr.info)
	
	// Generate recommendations
	report.Recommendations = vr.generateRecommendations()
	
	// Generate statistics
	report.Statistics = vr.generateStatistics()
	
	return report
}

// ValidationReport contains the comprehensive validation report
type ValidationReport struct {
	Summary         string                    `json:"summary"`
	Errors          []DataValidationError     `json:"errors"`
	Warnings        []DataValidationError     `json:"warnings"`
	Info            []DataValidationError     `json:"info"`
	ErrorSummary    map[string]int            `json:"error_summary"`
	WarningSummary  map[string]int            `json:"warning_summary"`
	InfoSummary     map[string]int            `json:"info_summary"`
	Recommendations []string                  `json:"recommendations"`
	Statistics      *ValidationStatistics     `json:"statistics"`
	TaskCount       int                       `json:"task_count"`
	ErrorCount      int                       `json:"error_count"`
	WarningCount    int                       `json:"warning_count"`
	InfoCount       int                       `json:"info_count"`
	Timestamp       time.Time                 `json:"timestamp"`
	IsValid         bool                      `json:"is_valid"`
}

// ValidationStatistics contains validation statistics
type ValidationStatistics struct {
	TotalIssues     int                       `json:"total_issues"`
	ErrorRate       float64                   `json:"error_rate"`
	WarningRate     float64                   `json:"warning_rate"`
	InfoRate        float64                   `json:"info_rate"`
	MostCommonError string                    `json:"most_common_error"`
	MostCommonField string                    `json:"most_common_field"`
	ErrorTypes      map[string]int            `json:"error_types"`
	FieldIssues     map[string]int            `json:"field_issues"`
	SeverityBreakdown map[string]int          `json:"severity_breakdown"`
}

// generateSummary generates a high-level summary of the validation
func (vr *ValidationReporter) generateSummary() string {
	if len(vr.errors) == 0 && len(vr.warnings) == 0 && len(vr.info) == 0 {
		return fmt.Sprintf("✅ Validation passed: %d tasks validated successfully", vr.taskCount)
	}
	
	var parts []string
	
	if len(vr.errors) > 0 {
		parts = append(parts, fmt.Sprintf("%d errors", len(vr.errors)))
	}
	if len(vr.warnings) > 0 {
		parts = append(parts, fmt.Sprintf("%d warnings", len(vr.warnings)))
	}
	if len(vr.info) > 0 {
		parts = append(parts, fmt.Sprintf("%d info messages", len(vr.info)))
	}
	
	summary := fmt.Sprintf("❌ Validation failed: %d tasks validated, %s", vr.taskCount, strings.Join(parts, ", "))
	
	if len(vr.errors) > 0 {
		summary += "\n🚨 Critical issues must be resolved before proceeding"
	}
	if len(vr.warnings) > 0 {
		summary += "\n⚠️  Warnings should be reviewed and addressed"
	}
	if len(vr.info) > 0 {
		summary += "\nℹ️  Info messages provide suggestions for improvement"
	}
	
	return summary
}

// generateErrorSummary generates a summary of errors by type
func (vr *ValidationReporter) generateErrorSummary(errors []DataValidationError) map[string]int {
	summary := make(map[string]int)
	
	for _, err := range errors {
		summary[err.Type]++
	}
	
	return summary
}

// generateRecommendations generates actionable recommendations
func (vr *ValidationReporter) generateRecommendations() []string {
	var recommendations []string
	
	// Analyze error patterns
	errorTypes := make(map[string]int)
	fieldIssues := make(map[string]int)
	
	allErrors := append(vr.errors, vr.warnings...)
	allErrors = append(allErrors, vr.info...)
	
	for _, err := range allErrors {
		errorTypes[err.Type]++
		fieldIssues[err.Field]++
	}
	
	// Generate recommendations based on error patterns
	if errorTypes["REQUIRED_FIELD"] > 0 {
		recommendations = append(recommendations, "🔧 Fix missing required fields (ID, Name, StartDate, EndDate)")
	}
	
	if errorTypes["FIELD_FORMAT"] > 0 {
		recommendations = append(recommendations, "🔧 Correct field format issues (invalid categories, statuses, priorities)")
	}
	
	if errorTypes["DATA_CONSISTENCY"] > 0 {
		recommendations = append(recommendations, "🔧 Resolve data consistency issues (self-references, duration problems)")
	}
	
	if errorTypes["CIRCULAR_DEPENDENCY"] > 0 {
		recommendations = append(recommendations, "🔧 Break circular dependencies in task relationships")
	}
	
	if errorTypes["DATE_RANGE"] > 0 {
		recommendations = append(recommendations, "🔧 Fix date range issues (invalid dates, conflicts)")
	}
	
	if errorTypes["DEPENDENCY"] > 0 {
		recommendations = append(recommendations, "🔧 Resolve dependency issues (missing tasks, invalid references)")
	}
	
	if fieldIssues["Category"] > 0 {
		recommendations = append(recommendations, "📝 Review task categorization and use valid categories")
	}
	
	if fieldIssues["Priority"] > 0 {
		recommendations = append(recommendations, "📝 Set appropriate task priorities (0-5 scale)")
	}
	
	if fieldIssues["Description"] > 0 {
		recommendations = append(recommendations, "📝 Add descriptions to tasks for better clarity")
	}
	
	// Add general recommendations
	if len(vr.errors) > 0 {
		recommendations = append(recommendations, "🚨 Address all errors before proceeding with calendar generation")
	}
	
	if len(vr.warnings) > 0 {
		recommendations = append(recommendations, "⚠️  Review warnings to improve data quality")
	}
	
	if len(vr.info) > 0 {
		recommendations = append(recommendations, "ℹ️  Consider implementing suggested improvements")
	}
	
	return recommendations
}

// generateStatistics generates validation statistics
func (vr *ValidationReporter) generateStatistics() *ValidationStatistics {
	allErrors := append(vr.errors, vr.warnings...)
	allErrors = append(allErrors, vr.info...)
	
	stats := &ValidationStatistics{
		TotalIssues:     len(allErrors),
		ErrorRate:       float64(len(vr.errors)) / float64(vr.taskCount) * 100,
		WarningRate:     float64(len(vr.warnings)) / float64(vr.taskCount) * 100,
		InfoRate:        float64(len(vr.info)) / float64(vr.taskCount) * 100,
		ErrorTypes:      make(map[string]int),
		FieldIssues:     make(map[string]int),
		SeverityBreakdown: make(map[string]int),
	}
	
	// Count error types and field issues
	for _, err := range allErrors {
		stats.ErrorTypes[err.Type]++
		stats.FieldIssues[err.Field]++
		stats.SeverityBreakdown[err.Severity]++
	}
	
	// Find most common error type
	maxCount := 0
	for errorType, count := range stats.ErrorTypes {
		if count > maxCount {
			maxCount = count
			stats.MostCommonError = errorType
		}
	}
	
	// Find most common field
	maxCount = 0
	for field, count := range stats.FieldIssues {
		if count > maxCount {
			maxCount = count
			stats.MostCommonField = field
		}
	}
	
	return stats
}

// filterAndSortErrors filters and sorts errors based on configuration
func (vr *ValidationReporter) filterAndSortErrors(errors []DataValidationError) []DataValidationError {
	// Apply limits
	filtered := errors
	if vr.config.MaxErrors > 0 && len(filtered) > vr.config.MaxErrors {
		filtered = filtered[:vr.config.MaxErrors]
	}
	if vr.config.MaxWarnings > 0 && len(filtered) > vr.config.MaxWarnings {
		filtered = filtered[:vr.config.MaxWarnings]
	}
	if vr.config.MaxInfo > 0 && len(filtered) > vr.config.MaxInfo {
		filtered = filtered[:vr.config.MaxInfo]
	}
	
	// Sort by severity if configured
	if vr.config.SortBySeverity {
		sort.Slice(filtered, func(i, j int) bool {
			severityOrder := map[string]int{"ERROR": 0, "WARNING": 1, "INFO": 2}
			return severityOrder[filtered[i].Severity] < severityOrder[filtered[j].Severity]
		})
	}
	
	return filtered
}

// GenerateTextReport generates a human-readable text report
func (vr *ValidationReporter) GenerateTextReport() string {
	report := vr.GenerateReport()
	
	var output strings.Builder
	
	// Header
	output.WriteString("=" + strings.Repeat("=", 60) + "\n")
	output.WriteString("VALIDATION REPORT\n")
	output.WriteString("=" + strings.Repeat("=", 60) + "\n")
	output.WriteString(fmt.Sprintf("Generated: %s\n", report.Timestamp.Format("2006-01-02 15:04:05")))
	output.WriteString(fmt.Sprintf("Tasks Validated: %d\n", report.TaskCount))
	output.WriteString(fmt.Sprintf("Status: %s\n", vr.getStatusEmoji(report.IsValid)))
	output.WriteString("\n")
	
	// Summary
	output.WriteString("SUMMARY\n")
	output.WriteString("-" + strings.Repeat("-", 20) + "\n")
	output.WriteString(report.Summary + "\n\n")
	
	// Statistics
	if report.Statistics != nil {
		output.WriteString("STATISTICS\n")
		output.WriteString("-" + strings.Repeat("-", 20) + "\n")
		output.WriteString(fmt.Sprintf("Total Issues: %d\n", report.Statistics.TotalIssues))
		output.WriteString(fmt.Sprintf("Error Rate: %.1f%%\n", report.Statistics.ErrorRate))
		output.WriteString(fmt.Sprintf("Warning Rate: %.1f%%\n", report.Statistics.WarningRate))
		output.WriteString(fmt.Sprintf("Info Rate: %.1f%%\n", report.Statistics.InfoRate))
		output.WriteString(fmt.Sprintf("Most Common Error: %s\n", report.Statistics.MostCommonError))
		output.WriteString(fmt.Sprintf("Most Common Field: %s\n", report.Statistics.MostCommonField))
		output.WriteString("\n")
	}
	
	// Errors
	if len(report.Errors) > 0 {
		output.WriteString("ERRORS\n")
		output.WriteString("-" + strings.Repeat("-", 20) + "\n")
		for i, err := range report.Errors {
			output.WriteString(fmt.Sprintf("%d. [%s] %s: %s\n", i+1, err.Severity, err.TaskID, err.Message))
			if vr.config.IncludeSuggestions && len(err.Suggestions) > 0 {
				output.WriteString("   Suggestions:\n")
				for _, suggestion := range err.Suggestions {
					output.WriteString(fmt.Sprintf("   - %s\n", suggestion))
				}
			}
			output.WriteString("\n")
		}
	}
	
	// Warnings
	if len(report.Warnings) > 0 {
		output.WriteString("WARNINGS\n")
		output.WriteString("-" + strings.Repeat("-", 20) + "\n")
		for i, err := range report.Warnings {
			output.WriteString(fmt.Sprintf("%d. [%s] %s: %s\n", i+1, err.Severity, err.TaskID, err.Message))
			if vr.config.IncludeSuggestions && len(err.Suggestions) > 0 {
				output.WriteString("   Suggestions:\n")
				for _, suggestion := range err.Suggestions {
					output.WriteString(fmt.Sprintf("   - %s\n", suggestion))
				}
			}
			output.WriteString("\n")
		}
	}
	
	// Info
	if len(report.Info) > 0 {
		output.WriteString("INFO MESSAGES\n")
		output.WriteString("-" + strings.Repeat("-", 20) + "\n")
		for i, err := range report.Info {
			output.WriteString(fmt.Sprintf("%d. [%s] %s: %s\n", i+1, err.Severity, err.TaskID, err.Message))
			if vr.config.IncludeSuggestions && len(err.Suggestions) > 0 {
				output.WriteString("   Suggestions:\n")
				for _, suggestion := range err.Suggestions {
					output.WriteString(fmt.Sprintf("   - %s\n", suggestion))
				}
			}
			output.WriteString("\n")
		}
	}
	
	// Recommendations
	if len(report.Recommendations) > 0 {
		output.WriteString("RECOMMENDATIONS\n")
		output.WriteString("-" + strings.Repeat("-", 20) + "\n")
		for i, rec := range report.Recommendations {
			output.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec))
		}
		output.WriteString("\n")
	}
	
	// Footer
	output.WriteString("=" + strings.Repeat("=", 60) + "\n")
	
	return output.String()
}

// getStatusEmoji returns an emoji for the validation status
func (vr *ValidationReporter) getStatusEmoji(isValid bool) string {
	if isValid {
		return "✅ PASSED"
	}
	return "❌ FAILED"
}

// GenerateCSVReport generates a CSV report of validation issues
func (vr *ValidationReporter) GenerateCSVReport() string {
	report := vr.GenerateReport()
	
	var output strings.Builder
	
	// CSV Header
	output.WriteString("Severity,Type,TaskID,Field,Value,Message,Timestamp\n")
	
	// Add all errors
	allErrors := append(report.Errors, report.Warnings...)
	allErrors = append(allErrors, report.Info...)
	
	for _, err := range allErrors {
		output.WriteString(fmt.Sprintf("%s,%s,%s,%s,\"%s\",\"%s\",%s\n",
			err.Severity,
			err.Type,
			err.TaskID,
			err.Field,
			err.Value,
			err.Message,
			err.Timestamp.Format("2006-01-02 15:04:05"),
		))
	}
	
	return output.String()
}

// GenerateJSONReport generates a JSON report of validation issues
func (vr *ValidationReporter) GenerateJSONReport() string {
	report := vr.GenerateReport()
	
	// Simple JSON marshaling (in a real implementation, use json.Marshal)
	json := fmt.Sprintf(`{
  "summary": "%s",
  "is_valid": %t,
  "task_count": %d,
  "error_count": %d,
  "warning_count": %d,
  "info_count": %d,
  "timestamp": "%s",
  "statistics": {
    "total_issues": %d,
    "error_rate": %.2f,
    "warning_rate": %.2f,
    "info_rate": %.2f,
    "most_common_error": "%s",
    "most_common_field": "%s"
  }
}`,
		report.Summary,
		report.IsValid,
		report.TaskCount,
		report.ErrorCount,
		report.WarningCount,
		report.InfoCount,
		report.Timestamp.Format("2006-01-02 15:04:05"),
		report.Statistics.TotalIssues,
		report.Statistics.ErrorRate,
		report.Statistics.WarningRate,
		report.Statistics.InfoRate,
		report.Statistics.MostCommonError,
		report.Statistics.MostCommonField,
	)
	
	return json
}

// Clear clears all validation data
func (vr *ValidationReporter) Clear() {
	vr.errors = make([]DataValidationError, 0)
	vr.warnings = make([]DataValidationError, 0)
	vr.info = make([]DataValidationError, 0)
	vr.taskCount = 0
	vr.timestamp = time.Now()
}

// GetErrorCount returns the total number of errors
func (vr *ValidationReporter) GetErrorCount() int {
	return len(vr.errors)
}

// GetWarningCount returns the total number of warnings
func (vr *ValidationReporter) GetWarningCount() int {
	return len(vr.warnings)
}

// GetInfoCount returns the total number of info messages
func (vr *ValidationReporter) GetInfoCount() int {
	return len(vr.info)
}

// GetTotalIssueCount returns the total number of issues
func (vr *ValidationReporter) GetTotalIssueCount() int {
	return len(vr.errors) + len(vr.warnings) + len(vr.info)
}

// HasErrors returns true if there are any errors
func (vr *ValidationReporter) HasErrors() bool {
	return len(vr.errors) > 0
}

// HasWarnings returns true if there are any warnings
func (vr *ValidationReporter) HasWarnings() bool {
	return len(vr.warnings) > 0
}

// HasIssues returns true if there are any issues (errors, warnings, or info)
func (vr *ValidationReporter) HasIssues() bool {
	return vr.GetTotalIssueCount() > 0
}
