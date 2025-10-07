package core

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	// DateFormats supported date formats in CSV files (in order of preference)
	DateFormatISO   = "2006-01-02"          // ISO format: 2024-01-15
	DateFormatUS    = "01/02/2006"          // US format: 01/15/2024
	DateFormatEU    = "02/01/2006"          // EU format: 15/01/2024
	DateFormatSlash = "2006/01/02"          // Slash format: 2024/01/15
	DateFormatDot   = "02.01.2006"          // Dot format: 15.01.2024
	DateFormatSpace = "2006-01-02 15:04:05" // With time: 2024-01-15 10:30:00
)

// Error types for detailed error reporting
type ParseError struct {
	Row     int
	Column  string
	Value   string
	Message string
	Err     error
}

func (e *ParseError) Error() string {
	if e.Row > 0 {
		return fmt.Sprintf("row %d, column '%s', value '%s': %s", e.Row, e.Column, e.Value, e.Message)
	}
	return fmt.Sprintf("column '%s', value '%s': %s", e.Column, e.Value, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// NewParseError creates a new parse error with proper wrapping
func NewParseError(row int, column, value, message string, err error) *ParseError {
	return &ParseError{
		Row:     row,
		Column:  column,
		Value:   value,
		Message: message,
		Err:     err,
	}
}

type ValidationError struct {
	TaskID  string
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("task %s, field '%s', value '%s': %s", e.TaskID, e.Field, e.Value, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(taskID, field, value, message string) *ValidationError {
	return &ValidationError{
		TaskID:  taskID,
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// Supported date formats for parsing
var supportedDateFormats = []string{
	DateFormatISO,
	DateFormatUS,
	DateFormatEU,
	DateFormatSlash,
	DateFormatDot,
	DateFormatSpace,
}

// Reader handles reading and parsing CSV task data
type Reader struct {
	filePath    string
	logger      *Logger
	aggregator  *ErrorAggregator // Error aggregator for collecting multiple errors
	strictMode  bool             // If true, fail on any parsing error
	skipInvalid bool             // If true, skip invalid rows instead of failing
	maxMemoryMB int              // Maximum memory usage in MB for large files
}

// ReaderOptions configures the CSV reader behavior
type ReaderOptions struct {
	StrictMode  bool
	SkipInvalid bool
	MaxMemoryMB int
	Logger      *Logger
}

// DefaultReaderOptions returns sensible defaults for the reader
func DefaultReaderOptions() *ReaderOptions {
	return &ReaderOptions{
		StrictMode:  false,
		SkipInvalid: true,
		MaxMemoryMB: 100, // 100MB default limit
		Logger:      NewDefaultLogger(),
	}
}

// NewReader creates a new CSV data reader with default options
func NewReader(filePath string) *Reader {
	opts := DefaultReaderOptions()
	return &Reader{
		filePath:    filePath,
		logger:      opts.Logger,
		aggregator:  NewErrorAggregator(),
		strictMode:  opts.StrictMode,
		skipInvalid: opts.SkipInvalid,
		maxMemoryMB: opts.MaxMemoryMB,
	}
}

// parseDate attempts to parse a date string using multiple supported formats
func (r *Reader) parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, NewParseError(0, "Date", dateStr, "empty date string", nil)
	}

	// Clean the date string
	dateStr = strings.TrimSpace(dateStr)

	// Try each supported format
	for _, format := range supportedDateFormats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, NewParseError(0, "Date", dateStr, 
		fmt.Sprintf("unable to parse with any supported format (tried: %v)", supportedDateFormats), nil)
}

// isMilestoneTask determines if a task is a milestone based on its name or description
func (r *Reader) isMilestoneTask(name, description string) bool {
	text := strings.ToLower(name + " " + description)
	milestoneKeywords := []string{"milestone", "deadline", "due", "complete", "finish", "submit", "deliver"}

	for _, keyword := range milestoneKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}

// addError adds an error to the aggregator
func (r *Reader) addError(err error) {
	r.aggregator.AddError(err)
}

// addWarning adds a warning to the aggregator
func (r *Reader) addWarning(err error) {
	r.aggregator.AddWarning(err)
}

// clearErrors clears all collected errors
func (r *Reader) clearErrors() {
	r.aggregator.Clear()
}

// hasErrors returns true if there are any collected errors
func (r *Reader) hasErrors() bool {
	return r.aggregator.HasErrors()
}

// getErrorSummary returns a summary of all errors
func (r *Reader) getErrorSummary() string {
	return r.aggregator.Summary()
}

// ReadTasks reads all tasks from the CSV file with improved error handling and memory management
func (r *Reader) ReadTasks() ([]Task, error) {
	// Clear any previous errors
	r.clearErrors()

	// Open and validate file
	file, fileInfo, err := r.openAndValidateFile()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Check file size for memory management
	r.checkFileSize(fileInfo)

	// Create CSV reader with configuration
	csvReader := r.createCSVReader(file)

	// Read and parse header
	fieldIndex, err := r.readHeader(csvReader)
	if err != nil {
		return nil, err
	}

	// Parse all task records
	tasks, parseErrors := r.parseAllRecords(csvReader, fieldIndex)

	// Check for fatal errors (strict mode or non-skippable errors)
	if len(parseErrors) > 0 && (r.strictMode || !r.skipInvalid) {
		return tasks, parseErrors[0] // Return first error
	}

	// Log parsing summary
	r.logParsingSummary(tasks, parseErrors)

	return tasks, nil
}

// openAndValidateFile opens the CSV file and returns file info
func (r *Reader) openAndValidateFile() (*os.File, os.FileInfo, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open CSV file: %w", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return file, fileInfo, nil
}

// checkFileSize logs a warning if file size exceeds memory limit
func (r *Reader) checkFileSize(fileInfo os.FileInfo) {
	fileSizeMB := fileInfo.Size() / (1024 * 1024)
	if fileSizeMB > int64(r.maxMemoryMB) {
		r.logger.Warn("File size %dMB exceeds limit %dMB, consider using streaming mode", fileSizeMB, r.maxMemoryMB)
	}
}

// createCSVReader creates a configured CSV reader
func (r *Reader) createCSVReader(file *os.File) *csv.Reader {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1    // Allow variable number of fields
	reader.TrimLeadingSpace = true // Trim leading spaces
	return reader
}

// readHeader reads the CSV header and creates field index map
func (r *Reader) readHeader(reader *csv.Reader) (map[string]int, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	return r.createFieldIndexMap(header), nil
}

// createFieldIndexMap creates a case-insensitive field index map
func (r *Reader) createFieldIndexMap(header []string) map[string]int {
	fieldIndex := make(map[string]int)
	for i, field := range header {
		normalizedField := strings.ToLower(strings.TrimSpace(field))
		fieldIndex[normalizedField] = i
	}
	return fieldIndex
}

// parseAllRecords parses all CSV records into tasks
func (r *Reader) parseAllRecords(reader *csv.Reader, fieldIndex map[string]int) ([]Task, []error) {
	var tasks []Task
	var parseErrors []error
	rowNum := 1 // Start from 1 (header is row 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("row %d: %w", rowNum, err))
			r.addError(fmt.Errorf("row %d: %w", rowNum, err))
			break
		}

		rowNum++

		// Skip empty rows
		if len(record) == 0 || record[0] == "" {
			continue
		}

		task, err := r.parseTask(record, fieldIndex, rowNum)
		if err != nil {
			parseErr := fmt.Errorf("row %d: %w", rowNum, err)
			parseErrors = append(parseErrors, parseErr)
			r.addError(parseErr)

			if r.strictMode {
				// Return error immediately in strict mode
				return tasks, []error{fmt.Errorf("strict mode: failed to parse task at row %d: %w", rowNum, err)}
			}

			if !r.skipInvalid {
				// Return error if not skipping invalid rows
				return tasks, []error{fmt.Errorf("failed to parse task at row %d: %w", rowNum, err)}
			}

			// Log warning but continue processing other tasks
			r.logger.Warn("Skipping invalid task at row %d: %v", rowNum, err)
			r.addWarning(fmt.Errorf("skipped invalid task at row %d: %w", rowNum, err))
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, parseErrors
}

// logParsingSummary logs a summary of the parsing results
func (r *Reader) logParsingSummary(tasks []Task, _ []error) {
	errorCount := r.aggregator.ErrorCount()
	warningCount := r.aggregator.WarningCount()

	if errorCount == 0 && warningCount == 0 {
		r.logger.Info("Successfully parsed %d tasks with no issues", len(tasks))
		return
	}

	if errorCount > 0 && warningCount > 0 {
		r.logger.Info("Parsed %d tasks with %d errors and %d warnings", len(tasks), errorCount, warningCount)
	} else if errorCount > 0 {
		r.logger.Info("Parsed %d tasks with %d errors", len(tasks), errorCount)
	} else {
		r.logger.Info("Parsed %d tasks with %d warnings", len(tasks), warningCount)
	}

	if r.hasErrors() || r.aggregator.HasWarnings() {
		r.logger.Warn("Parsing details:\n%s", r.getErrorSummary())
	}
}

// fieldExtractor helps extract fields from CSV records with case-insensitive matching
type fieldExtractor struct {
	record     []string
	fieldIndex map[string]int
}

// newFieldExtractor creates a new field extractor
func newFieldExtractor(record []string, fieldIndex map[string]int) *fieldExtractor {
	return &fieldExtractor{
		record:     record,
		fieldIndex: fieldIndex,
	}
}

// get retrieves a field value by name (case-insensitive)
func (fe *fieldExtractor) get(fieldName string) string {
	normalizedField := strings.ToLower(strings.TrimSpace(fieldName))
	if index, exists := fe.fieldIndex[normalizedField]; exists && index < len(fe.record) {
		return strings.TrimSpace(fe.record[index])
	}
	return ""
}

// getWithDefault retrieves a field value with a default fallback
func (fe *fieldExtractor) getWithDefault(fieldName, defaultValue string) string {
	value := fe.get(fieldName)
	if value == "" {
		return defaultValue
	}
	return value
}

// getList retrieves a comma-separated list field as a slice
func (fe *fieldExtractor) getList(fieldName string) []string {
	value := fe.get(fieldName)
	if value == "" {
		return nil
	}

	var result []string
	parts := strings.Split(value, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// parseTask parses a single CSV record into a Task struct with improved field mapping
func (r *Reader) parseTask(record []string, fieldIndex map[string]int, rowNum int) (Task, error) {
	extractor := newFieldExtractor(record, fieldIndex)
	task := Task{}

	// Extract basic fields
	r.extractBasicFields(&task, extractor)

	// Extract phase and category
	r.extractPhaseFields(&task, extractor)

	// Extract status and assignment
	r.extractStatusFields(&task, extractor)

	// Extract dependencies
	task.Dependencies = extractor.getList("Dependencies")

	// Parse dates
	if err := r.extractDateFields(&task, extractor, rowNum); err != nil {
		return task, err
	}

	// Validate dates
	if err := r.validateDates(task); err != nil {
		return task, err
	}

	return task, nil
}

// extractBasicFields extracts ID, name, description, and milestone status
func (r *Reader) extractBasicFields(task *Task, extractor *fieldExtractor) {
	task.ID = extractor.get("Task ID")
	if task.ID == "" {
		task.ID = extractor.get("Task") // Fallback to name if no ID
	}
	task.Name = extractor.get("Task")
	task.Description = extractor.get("Objective")

	// Extract milestone status from CSV column or detect from content
	milestoneValue := extractor.get("Milestone")
	if milestoneValue != "" && strings.ToLower(milestoneValue) != "false" {
		task.IsMilestone = true
	} else {
		// Fallback to keyword detection
		task.IsMilestone = r.isMilestoneTask(task.Name, task.Description)
	}
}

// extractPhaseFields extracts phase and category information
func (r *Reader) extractPhaseFields(task *Task, extractor *fieldExtractor) {
	task.Phase = extractor.get("Phase")
	task.SubPhase = extractor.get("Sub-Phase")

	// Use Sub-Phase as primary category for better granularity
	task.Category = task.SubPhase
	if task.Category == "" {
		task.Category = task.Phase // Fallback to Phase if Sub-Phase is empty
	}
}

// extractStatusFields extracts status and assignee
func (r *Reader) extractStatusFields(task *Task, extractor *fieldExtractor) {
	task.Status = extractor.getWithDefault("Status", "Planned")
	task.Assignee = extractor.get("Assignee")
	task.ParentID = extractor.get("Parent Task ID")
}

// extractDateFields parses date fields from the extractor
func (r *Reader) extractDateFields(task *Task, extractor *fieldExtractor, rowNum int) error {
	startDateStr := extractor.get("Start Date")
	if startDateStr != "" {
		startDate, err := r.parseDate(startDateStr)
		if err != nil {
			return NewParseError(rowNum, "Start Date", startDateStr, "invalid date format", err)
		}
		task.StartDate = startDate
	}

	endDateStr := extractor.get("End Date")
	if endDateStr != "" {
		endDate, err := r.parseDate(endDateStr)
		if err != nil {
			return NewParseError(rowNum, "End Date", endDateStr, "invalid date format", err)
		}
		task.EndDate = endDate
	}

	return nil
}

// validateDates validates that end date is not before start date
func (r *Reader) validateDates(task Task) error {
	if !task.StartDate.IsZero() && !task.EndDate.IsZero() && task.EndDate.Before(task.StartDate) {
		return NewValidationError(
			task.ID,
			"Due Date",
			task.EndDate.Format("2006-01-02"),
			fmt.Sprintf("end date %s is before start date %s", 
				task.EndDate.Format("2006-01-02"), 
				task.StartDate.Format("2006-01-02")),
		)
	}
	return nil
}
