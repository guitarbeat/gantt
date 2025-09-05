package data

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
	DateFormatISO     = "2006-01-02"     // ISO format: 2024-01-15
	DateFormatUS      = "01/02/2006"     // US format: 01/15/2024
	DateFormatEU      = "02/01/2006"     // EU format: 15/01/2024
	DateFormatSlash   = "2006/01/02"     // Slash format: 2024/01/15
	DateFormatDot     = "02.01.2006"     // Dot format: 15.01.2024
	DateFormatSpace   = "2006-01-02 15:04:05" // With time: 2024-01-15 10:30:00
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

type CircularDependencyError struct {
	Cycle []string
}

func (e *CircularDependencyError) Error() string {
	return fmt.Sprintf("circular dependency detected: %s", strings.Join(e.Cycle, " -> "))
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

// Task represents a single task from the CSV data
type Task struct {
	ID          string
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	Category    string    // * Fixed: Use Category instead of Priority for clarity
	Description string
	Priority    int       // * Added: Separate priority field for task ordering
	Status      string    // * Added: Task status (Planned, In Progress, Completed, etc.)
	Assignee    string    // * Added: Task assignee
	ParentID    string    // * Added: Parent task ID for hierarchical relationships
	Dependencies []string // * Added: List of task IDs this task depends on
	IsMilestone bool      // * Added: Whether this is a milestone task
}

// DateRange represents the earliest and latest dates from the task data
type DateRange struct {
	Earliest time.Time
	Latest   time.Time
}

// MonthYear represents a specific month and year
type MonthYear struct {
	Year  int
	Month time.Month
}

// Reader handles reading and parsing CSV task data
type Reader struct {
	filePath string
	logger   *log.Logger
	// * Added: Configuration options
	strictMode    bool // If true, fail on any parsing error
	skipInvalid   bool // If true, skip invalid rows instead of failing
	maxMemoryMB   int  // Maximum memory usage in MB for large files
	// * Added: Enhanced parsing options
	validateDependencies bool // If true, validate that all dependencies exist
	detectCircularDeps   bool // If true, detect circular dependencies
	// * Added: Error collection
	errors []error // Collected errors during parsing
}

// ReaderOptions configures the CSV reader behavior
type ReaderOptions struct {
	StrictMode            bool
	SkipInvalid           bool
	MaxMemoryMB           int
	Logger                *log.Logger
	ValidateDependencies  bool // * Added: Validate that all dependencies exist
	DetectCircularDeps    bool // * Added: Detect circular dependencies
}

// DefaultReaderOptions returns sensible defaults for the reader
func DefaultReaderOptions() *ReaderOptions {
	return &ReaderOptions{
		StrictMode:            false,
		SkipInvalid:           true,
		MaxMemoryMB:           100, // 100MB default limit
		Logger:                log.New(os.Stderr, "[data] ", log.LstdFlags|log.Lshortfile),
		ValidateDependencies:  true,  // * Added: Default to validating dependencies
		DetectCircularDeps:    true,  // * Added: Default to detecting circular dependencies
	}
}

// NewReader creates a new CSV data reader with default options
func NewReader(filePath string) *Reader {
	opts := DefaultReaderOptions()
	return &Reader{
		filePath:             filePath,
		logger:               opts.Logger,
		strictMode:           opts.StrictMode,
		skipInvalid:          opts.SkipInvalid,
		maxMemoryMB:          opts.MaxMemoryMB,
		validateDependencies: opts.ValidateDependencies,
		detectCircularDeps:   opts.DetectCircularDeps,
	}
}

// NewReaderWithOptions creates a new CSV data reader with custom options
func NewReaderWithOptions(filePath string, opts *ReaderOptions) *Reader {
	if opts == nil {
		opts = DefaultReaderOptions()
	}
	return &Reader{
		filePath:             filePath,
		logger:               opts.Logger,
		strictMode:           opts.StrictMode,
		skipInvalid:          opts.SkipInvalid,
		maxMemoryMB:          opts.MaxMemoryMB,
		validateDependencies: opts.ValidateDependencies,
		detectCircularDeps:   opts.DetectCircularDeps,
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

// parseDependencies parses comma-separated dependency task IDs
func (r *Reader) parseDependencies(depsStr string) []string {
	if depsStr == "" {
		return []string{}
	}
	
	// Split by comma and clean up each dependency
	deps := strings.Split(depsStr, ",")
	var cleanDeps []string
	for _, dep := range deps {
		cleanDep := strings.TrimSpace(dep)
		if cleanDep != "" {
			cleanDeps = append(cleanDeps, cleanDep)
		}
	}
	
	return cleanDeps
}

// validateTaskDependencies checks that all referenced task IDs exist
func (r *Reader) validateTaskDependencies(tasks []Task) error {
	if !r.validateDependencies {
		return nil
	}
	
	// Create a map of existing task IDs
	taskIDs := make(map[string]bool)
	for _, task := range tasks {
		taskIDs[task.ID] = true
	}
	
	// Check each task's dependencies
	var validationErrors []error
	for _, task := range tasks {
		for _, depID := range task.Dependencies {
			if !taskIDs[depID] {
				validationErrors = append(validationErrors, &ValidationError{
					TaskID:  task.ID,
					Field:   "Dependencies",
					Value:   depID,
					Message: "references non-existent task",
				})
			}
		}
	}
	
	if len(validationErrors) > 0 {
		// Return the first error for now, but log all errors
		for _, err := range validationErrors {
			r.logger.Printf("Validation error: %v", err)
		}
		return validationErrors[0]
	}
	
	return nil
}

// detectCircularDependencies checks for circular dependencies in the task graph
func (r *Reader) detectCircularDependencies(tasks []Task) error {
	if !r.detectCircularDeps {
		return nil
	}
	
	// Create adjacency list for dependency graph
	graph := make(map[string][]string)
	for _, task := range tasks {
		graph[task.ID] = task.Dependencies
	}
	
	// Use DFS to detect cycles and track the cycle path
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := make([]string, 0)
	
	var dfs func(string) (bool, []string)
	dfs = func(node string) (bool, []string) {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if hasCycle, cycle := dfs(neighbor); hasCycle {
					return true, cycle
				}
			} else if recStack[neighbor] {
				// Found a cycle, construct the cycle path
				cycleStart := -1
				for i, id := range path {
					if id == neighbor {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := append(path[cycleStart:], neighbor)
					return true, cycle
				}
			}
		}
		
		recStack[node] = false
		path = path[:len(path)-1]
		return false, nil
	}
	
	// Check each unvisited node
	for taskID := range graph {
		if !visited[taskID] {
			if hasCycle, cycle := dfs(taskID); hasCycle {
				return &CircularDependencyError{Cycle: cycle}
			}
		}
	}
	
	return nil
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

// getErrors returns all collected errors
func (r *Reader) getErrors() []error {
	return r.errors
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
	reader.FieldsPerRecord = -1 // Allow variable number of fields
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

	// * Added: Validate dependencies if enabled
	if r.validateDependencies {
		if err := r.validateTaskDependencies(tasks); err != nil {
			return nil, fmt.Errorf("dependency validation failed: %w", err)
		}
	}

	// * Added: Detect circular dependencies if enabled
	if r.detectCircularDeps {
		if err := r.detectCircularDependencies(tasks); err != nil {
			r.addError(err)
			return nil, fmt.Errorf("circular dependency detection failed: %w", err)
		}
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

	// Parse required fields
	task.ID = getField("Task ID")
	if task.ID == "" {
		return task, &ParseError{
			Row:     rowNum,
			Column:  "Task ID",
			Value:   "",
			Message: "missing required field",
		}
	}

	task.Name = getField("Task Name")
	task.Description = getField("Description")

	// * Fixed: Use Category field instead of Priority for clarity
	task.Category = getField("Category")
	
	// * Added: Parse Priority as integer if available
	if priorityStr := getField("Priority"); priorityStr != "" {
		// Try to parse as integer, default to 1 if parsing fails
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			task.Priority = priority
		} else {
			task.Priority = 1 // Default priority
			r.logger.Printf("Warning: Invalid priority '%s' for task %s, using default priority 1", priorityStr, task.ID)
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

	// * Added: Parse Dependencies field
	depsStr := getField("Dependencies")
	task.Dependencies = r.parseDependencies(depsStr)

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

	endDateStr := getField("Due Date")
	if endDateStr != "" {
		endDate, err := r.parseDate(endDateStr)
		if err != nil {
			parseErr := &ParseError{
				Row:     rowNum,
				Column:  "Due Date",
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
			TaskID:  task.ID,
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
	requiredFields := []string{"task id", "task name", "start date", "due date"}
	optionalFields := []string{"parent task id", "dependencies", "category", "description", "priority", "status", "assignee"}
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
