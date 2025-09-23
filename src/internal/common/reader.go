package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
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

type ValidationError struct {
	TaskID  string
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("task %s, field '%s', value '%s': %s", e.TaskID, e.Field, e.Value, e.Message)
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
	filePath string
	logger   *log.Logger
	// * Added: Configuration options
	strictMode  bool // If true, fail on any parsing error
	skipInvalid bool // If true, skip invalid rows instead of failing
	maxMemoryMB int  // Maximum memory usage in MB for large files
	// * Added: Error collection
	errors []error // Collected errors during parsing
}

// ReaderOptions configures the CSV reader behavior
type ReaderOptions struct {
	StrictMode  bool
	SkipInvalid bool
	MaxMemoryMB int
	Logger      *log.Logger
}

// DefaultReaderOptions returns sensible defaults for the reader
func DefaultReaderOptions() *ReaderOptions {
	return &ReaderOptions{
		StrictMode:  false,
		SkipInvalid: true,
		MaxMemoryMB: 100, // 100MB default limit
		Logger:      newDefaultLogger(),
	}
}

// newDefaultLogger returns a logger that can be silenced via env flag.
// Set PLANNER_SILENT=1 to suppress data-layer logs.
func newDefaultLogger() *log.Logger {
	if os.Getenv("PLANNER_SILENT") == "1" || strings.EqualFold(os.Getenv("PLANNER_LOG_LEVEL"), "silent") {
		return log.New(io.Discard, "", 0)
	}
	return log.New(os.Stderr, "[data] ", log.LstdFlags|log.Lshortfile)
}

// NewReader creates a new CSV data reader with default options
func NewReader(filePath string) *Reader {
	opts := DefaultReaderOptions()
	return &Reader{
		filePath:    filePath,
		logger:      opts.Logger,
		strictMode:  opts.StrictMode,
		skipInvalid: opts.SkipInvalid,
		maxMemoryMB: opts.MaxMemoryMB,
	}
}

// NewReaderWithOptions creates a new CSV data reader with custom options
func NewReaderWithOptions(filePath string, opts *ReaderOptions) *Reader {
	if opts == nil {
		opts = DefaultReaderOptions()
	}
	return &Reader{
		filePath:    filePath,
		logger:      opts.Logger,
		strictMode:  opts.StrictMode,
		skipInvalid: opts.SkipInvalid,
		maxMemoryMB: opts.MaxMemoryMB,
	}
}

// parseDate attempts to parse a date string using multiple supported formats
func (r *Reader) parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, &ParseError{
			Column:  "Date",
			Value:   dateStr,
			Message: "empty date string",
		}
	}

	// Clean the date string
	dateStr = strings.TrimSpace(dateStr)

	// Try each supported format
	for _, format := range supportedDateFormats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, &ParseError{
		Column:  "Date",
		Value:   dateStr,
		Message: fmt.Sprintf("unable to parse with any supported format (tried: %v)", supportedDateFormats),
	}
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

// addError adds an error to the reader's error collection
func (r *Reader) addError(err error) {
	r.errors = append(r.errors, err)
}

// clearErrors clears all collected errors
func (r *Reader) clearErrors() {
	r.errors = nil
}

// hasErrors returns true if there are any collected errors
func (r *Reader) hasErrors() bool {
	return len(r.errors) > 0
}

// getErrorSummary returns a summary of all errors
func (r *Reader) getErrorSummary() string {
	if len(r.errors) == 0 {
		return "No errors"
	}

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Found %d errors:\n", len(r.errors)))

	for i, err := range r.errors {
		summary.WriteString(fmt.Sprintf("%d. %v\n", i+1, err))
	}

	return summary.String()
}

// ReadTasks reads all tasks from the CSV file with improved error handling and memory management
func (r *Reader) ReadTasks() ([]Task, error) {
	// Clear any previous errors
	r.clearErrors()

	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	// * Added: Check file size for memory management
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	fileSizeMB := fileInfo.Size() / (1024 * 1024)
	if fileSizeMB > int64(r.maxMemoryMB) {
		r.logger.Printf("Warning: File size %dMB exceeds limit %dMB, consider using streaming mode", fileSizeMB, r.maxMemoryMB)
	}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1    // Allow variable number of fields
	reader.TrimLeadingSpace = true // * Added: Trim leading spaces

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// * Improved: Create field index map with case-insensitive matching
	fieldIndex := make(map[string]int)
	for i, field := range header {
		normalizedField := strings.ToLower(strings.TrimSpace(field))
		fieldIndex[normalizedField] = i
	}

	var tasks []Task
	var parseErrors []error
	rowNum := 1 // Start from 1 (header is row 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record at row %d: %w", rowNum, err)
		}

		rowNum++

		// Skip empty rows
		if len(record) == 0 || record[0] == "" {
			continue
		}

		task, err := r.parseTask(record, fieldIndex, rowNum)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("row %d: %w", rowNum, err))
			r.addError(fmt.Errorf("row %d: %w", rowNum, err))

			if r.strictMode {
				return nil, fmt.Errorf("strict mode: failed to parse task at row %d: %w", rowNum, err)
			}

			if !r.skipInvalid {
				return nil, fmt.Errorf("failed to parse task at row %d: %w", rowNum, err)
			}

			// Log error but continue processing other tasks
			r.logger.Printf("Warning: failed to parse task at row %d: %v", rowNum, err)
			continue
		}

		tasks = append(tasks, task)
	}

	// * Added: Log summary of parsing results
	if len(parseErrors) > 0 {
		r.logger.Printf("Parsed %d tasks successfully, %d errors encountered", len(tasks), len(parseErrors))
	} else {
		r.logger.Printf("Successfully parsed %d tasks", len(tasks))
	}

	// * Added: Log comprehensive error summary if there were any errors
	if r.hasErrors() {
		r.logger.Printf("Parsing completed with errors:\n%s", r.getErrorSummary())
	}

	return tasks, nil
}

// GetDateRange extracts the earliest and latest dates from the tasks
func (r *Reader) GetDateRange() (*DateRange, error) {
	tasks, err := r.ReadTasks()
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no tasks found in CSV file")
	}

	earliest := tasks[0].StartDate
	latest := tasks[0].EndDate

	for _, task := range tasks {
		if task.StartDate.Before(earliest) {
			earliest = task.StartDate
		}
		if task.EndDate.After(latest) {
			latest = task.EndDate
		}
	}

	return &DateRange{
		Earliest: earliest,
		Latest:   latest,
	}, nil
}

// GetMonthsWithTasks returns a slice of MonthYear structs for months that have tasks
func (r *Reader) GetMonthsWithTasks() ([]MonthYear, error) {
	tasks, err := r.ReadTasks()
	if err != nil {
		return nil, err
	}

	// Track which months have tasks using a map for deduplication
	monthsWithTasks := make(map[MonthYear]bool)

	for _, task := range tasks {
		// Add all months from start to end (inclusive)
		current := task.StartDate
		end := task.EndDate

		for !current.After(end) {
			month := MonthYear{
				Year:  current.Year(),
				Month: current.Month(),
			}
			monthsWithTasks[month] = true
			current = current.AddDate(0, 1, 0)
		}
	}

	// Convert to slice and sort
	months := make([]MonthYear, 0, len(monthsWithTasks))
	for month := range monthsWithTasks {
		months = append(months, month)
	}

	// Sort by year, then by month
	sort.Slice(months, func(i, j int) bool {
		if months[i].Year != months[j].Year {
			return months[i].Year < months[j].Year
		}
		return months[i].Month < months[j].Month
	})

	return months, nil
}

// parseTask parses a single CSV record into a Task struct with improved field mapping
func (r *Reader) parseTask(record []string, fieldIndex map[string]int, rowNum int) (Task, error) {
	task := Task{}

	// Helper function to get field value safely with case-insensitive matching
	getField := func(fieldName string) string {
		normalizedField := strings.ToLower(strings.TrimSpace(fieldName))
		if index, exists := fieldIndex[normalizedField]; exists && index < len(record) {
			return strings.TrimSpace(record[index])
		}
		return ""
	}

	// * Added: Parse Task ID field (use Name as ID if not provided)
	task.ID = getField("Task ID")
	if task.ID == "" {
		task.ID = getField("Task") // Fallback to name if no ID
	}

	task.Name = getField("Task")
	task.Description = getField("Objective")

	// * Use Sub-Phase as the primary category for better granularity
	task.Category = getField("Sub-Phase")
	if task.Category == "" {
		task.Category = getField("Phase") // Fallback to Phase if Sub-Phase is empty
	}

	// * Added: Parse Priority as integer if available
	if priorityStr := getField("Priority"); priorityStr != "" {
		// Try to parse as integer, default to 1 if parsing fails
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			task.Priority = priority
		} else {
			task.Priority = 1 // Default priority
			r.logger.Printf("Warning: Invalid priority '%s', using default priority 1", priorityStr)
		}
	} else {
		task.Priority = 1 // Default priority
	}

	// * Added: Parse Status field
	task.Status = getField("Status")
	if task.Status == "" {
		task.Status = "Planned" // Default status
	}

	// * Added: Parse Assignee field
	task.Assignee = getField("Assignee")

	// * Added: Parse Parent Task ID field
	task.ParentID = getField("Parent Task ID")

	// * Added: Parse Dependencies field (comma-separated task IDs)
	dependenciesStr := getField("Dependencies")
	if dependenciesStr != "" {
		// Split by comma and trim spaces
		deps := strings.Split(dependenciesStr, ",")
		for _, dep := range deps {
			trimmed := strings.TrimSpace(dep)
			if trimmed != "" {
				task.Dependencies = append(task.Dependencies, trimmed)
			}
		}
	}

	// * Added: Determine if this is a milestone task
	task.IsMilestone = r.isMilestoneTask(task.Name, task.Description)

	// * Improved: Parse dates with flexible format support
	startDateStr := getField("Start Date")
	if startDateStr != "" {
		startDate, err := r.parseDate(startDateStr)
		if err != nil {
			parseErr := &ParseError{
				Row:     rowNum,
				Column:  "Start Date",
				Value:   startDateStr,
				Message: "invalid date format",
				Err:     err,
			}
			return task, parseErr
		}
		task.StartDate = startDate
	}

	endDateStr := getField("End Date")
	if endDateStr != "" {
		endDate, err := r.parseDate(endDateStr)
		if err != nil {
			parseErr := &ParseError{
				Row:     rowNum,
				Column:  "End Date",
				Value:   endDateStr,
				Message: "invalid date format",
				Err:     err,
			}
			return task, parseErr
		}
		task.EndDate = endDate
	}

	// * Added: Validate that end date is not before start date
	if !task.StartDate.IsZero() && !task.EndDate.IsZero() && task.EndDate.Before(task.StartDate) {
		return task, &ValidationError{
			TaskID:  "unknown",
			Field:   "Due Date",
			Value:   endDateStr,
			Message: fmt.Sprintf("end date %s is before start date %s", task.EndDate.Format("2006-01-02"), task.StartDate.Format("2006-01-02")),
		}
	}

	return task, nil
}

// ReadTasksStreaming reads tasks from CSV file in streaming mode for large files
func (r *Reader) ReadTasksStreaming(processTask func(Task) error) error {
	file, err := os.Open(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create field index map with case-insensitive matching
	fieldIndex := make(map[string]int)
	for i, field := range header {
		normalizedField := strings.ToLower(strings.TrimSpace(field))
		fieldIndex[normalizedField] = i
	}

	rowNum := 1
	var parseErrors []error

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record at row %d: %w", rowNum, err)
		}

		rowNum++

		// Skip empty rows
		if len(record) == 0 || record[0] == "" {
			continue
		}

		task, err := r.parseTask(record, fieldIndex, rowNum)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Errorf("row %d: %w", rowNum, err))

			if r.strictMode {
				return fmt.Errorf("strict mode: failed to parse task at row %d: %w", rowNum, err)
			}

			if !r.skipInvalid {
				return fmt.Errorf("failed to parse task at row %d: %w", rowNum, err)
			}

			r.logger.Printf("Warning: failed to parse task at row %d: %v", rowNum, err)
			continue
		}

		// Process the task immediately
		if err := processTask(task); err != nil {
			return fmt.Errorf("failed to process task at row %d: %w", rowNum, err)
		}
	}

	// Log summary
	if len(parseErrors) > 0 {
		r.logger.Printf("Streaming completed with %d errors encountered", len(parseErrors))
	} else {
		r.logger.Printf("Streaming completed successfully")
	}

	return nil
}

// GetSupportedDateFormats returns the list of supported date formats
func GetSupportedDateFormats() []string {
	return append([]string{}, supportedDateFormats...)
}

// ValidateCSVFormat validates that a CSV file has the required columns
func (r *Reader) ValidateCSVFormat() error {
	file, err := os.Open(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Check for required fields
	requiredFields := []string{"task name", "start date", "due date"}
	optionalFields := []string{"parent task id", "category", "description", "priority", "status", "assignee"}
	fieldMap := make(map[string]bool)

	for _, field := range header {
		normalizedField := strings.ToLower(strings.TrimSpace(field))
		fieldMap[normalizedField] = true
	}

	var missingFields []string
	for _, required := range requiredFields {
		if !fieldMap[required] {
			missingFields = append(missingFields, required)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %v", missingFields)
	}

	// * Added: Log available optional fields
	var availableOptional []string
	for _, field := range optionalFields {
		if fieldMap[field] {
			availableOptional = append(availableOptional, field)
		}
	}
	if len(availableOptional) > 0 {
		r.logger.Printf("Detected optional fields: %v", availableOptional)
	}

	return nil
}
