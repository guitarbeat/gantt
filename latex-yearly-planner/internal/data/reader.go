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
		Logger:      log.New(os.Stderr, "[data] ", log.LstdFlags|log.Lshortfile),
	}
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
		return time.Time{}, fmt.Errorf("empty date string")
	}

	// Clean the date string
	dateStr = strings.TrimSpace(dateStr)
	
	// Try each supported format
	for _, format := range supportedDateFormats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse date '%s' with any supported format", dateStr)
}

// ReadTasks reads all tasks from the CSV file with improved error handling and memory management
func (r *Reader) ReadTasks() ([]Task, error) {
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
		return task, fmt.Errorf("missing Task ID")
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

	// * Improved: Parse dates with flexible format support
	startDateStr := getField("Start Date")
	if startDateStr != "" {
		startDate, err := r.parseDate(startDateStr)
		if err != nil {
			return task, fmt.Errorf("invalid start date '%s': %w", startDateStr, err)
		}
		task.StartDate = startDate
	}

	endDateStr := getField("Due Date")
	if endDateStr != "" {
		endDate, err := r.parseDate(endDateStr)
		if err != nil {
			return task, fmt.Errorf("invalid end date '%s': %w", endDateStr, err)
		}
		task.EndDate = endDate
	}

	// * Added: Validate that end date is not before start date
	if !task.StartDate.IsZero() && !task.EndDate.IsZero() && task.EndDate.Before(task.StartDate) {
		return task, fmt.Errorf("end date %s is before start date %s", task.EndDate.Format("2006-01-02"), task.StartDate.Format("2006-01-02"))
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

	return nil
}
