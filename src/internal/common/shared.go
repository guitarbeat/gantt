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

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug Debug

	Year                int `env:"PLANNER_YEAR"`
	WeekStart           time.Weekday
	Dotted              bool
	CalAfterSchedule    bool
	ClearTopRightCorner bool
	AMPMTime            bool
	AddLastHalfHour     bool

	// Data source configuration
	CSVFilePath string `env:"PLANNER_CSV_FILE"`
	StartYear   int    `env:"PLANNER_START_YEAR"`
	EndYear     int    `env:"PLANNER_END_YEAR"`

	// Months with tasks (populated from CSV)
	MonthsWithTasks []MonthYear

	Pages Pages

	Layout Layout

	// OutputDir is the directory where generated .tex and .pdf files will be written
	// Defaults to "build" when not provided via environment or config
	OutputDir string `env:"PLANNER_OUTPUT_DIR"`
}

type Debug struct {
	ShowFrame bool
	ShowLinks bool
}

type Pages []Page
type Page struct {
	Name         string
	RenderBlocks RenderBlocks
}

type RenderBlocks []RenderBlock

type Modules []Module
type Module struct {
	Cfg  Config
	Tpl  string
	Body interface{}
}

type RenderBlock struct {
	FuncName string
	Tpls     []string
}

type Colors struct {
	Gray      string
	LightGray string
}

type Layout struct {
	Paper Paper

	Numbers Numbers
	Lengths Lengths
	Colors  Colors
}

type Numbers struct {
	ArrayStretch float64
}

type Lengths struct {
	TabColSep             string
	LineThicknessDefault  string
	LineThicknessThick    string
	LineHeightButLine     string
	TwoColSep             string
	TriColSep             string
	FiveColSep            string
	MonthlyCellHeight     string
	HeaderResizeBox       string
	HeaderSideMonthsWidth string
	MonthlySpring         string
}

type Paper struct {
	Width  string `env:"PLANNER_LAYOUT_PAPER_WIDTH"`
	Height string `env:"PLANNER_LAYOUT_PAPER_HEIGHT"`

	Margin Margin

	ReverseMargins bool
	MarginParWidth string
	MarginParSep   string
}

type Margin struct {
	Top    string `env:"PLANNER_LAYOUT_PAPER_MARGIN_TOP"`
	Bottom string `env:"PLANNER_LAYOUT_PAPER_MARGIN_BOTTOM"`
	Left   string `env:"PLANNER_LAYOUT_PAPER_MARGIN_LEFT"`
	Right  string `env:"PLANNER_LAYOUT_PAPER_MARGIN_RIGHT"`
}

func NewConfig(pathConfigs ...string) (Config, error) {
	var (
		bts []byte
		err error
		cfg Config
	)

	for _, filepath := range pathConfigs {
		// * Skip missing files instead of failing
		if bts, err = os.ReadFile(strings.ToLower(filepath)); err != nil {
			if os.IsNotExist(err) {
				// * File doesn't exist, skip it
				continue
			}
			return cfg, fmt.Errorf("read file: %w", err)
		}

		// * Skip empty files
		if len(strings.TrimSpace(string(bts))) == 0 {
			continue
		}

		if err = yaml.Unmarshal(bts, &cfg); err != nil {
			return cfg, fmt.Errorf("yaml unmarshal: %w", err)
		}
	}

	if err = env.Parse(&cfg); err != nil {
		return cfg, fmt.Errorf("env parse: %w", err)
	}

	if cfg.Year == 0 {
		cfg.Year = time.Now().Year()
	}

	// Default output dir
	if strings.TrimSpace(cfg.OutputDir) == "" {
		cfg.OutputDir = "build"
	}

	// If CSV file is provided, determine date range dynamically
	if cfg.CSVFilePath != "" {
		if err := cfg.setDateRangeFromCSV(); err != nil {
			return cfg, fmt.Errorf("failed to set date range from CSV: %w", err)
		}
	}

	return cfg, nil
}

// setDateRangeFromCSV reads the CSV file and sets the start and end years
func (cfg *Config) setDateRangeFromCSV() error {
	reader := NewReader(cfg.CSVFilePath)
	dateRange, err := reader.GetDateRange()
	if err != nil {
		return fmt.Errorf("failed to get date range: %w", err)
	}

	cfg.StartYear = dateRange.Earliest.Year()
	cfg.EndYear = dateRange.Latest.Year()

	// Get months with tasks
	monthsWithTasks, err := reader.GetMonthsWithTasks()
	if err != nil {
		return fmt.Errorf("failed to get months with tasks: %w", err)
	}
	cfg.MonthsWithTasks = monthsWithTasks

	// If we have months with tasks, limit the year range to only those years
	if len(monthsWithTasks) > 0 {
		// Find the unique years from the months with tasks
		yearSet := make(map[int]bool)
		for _, monthYear := range monthsWithTasks {
			yearSet[monthYear.Year] = true
		}

		// Set the year range to only include years with tasks
		years := make([]int, 0, len(yearSet))
		for year := range yearSet {
			years = append(years, year)
		}

		if len(years) > 0 {
			cfg.StartYear = years[0]
			cfg.EndYear = years[len(years)-1]
		}
	}

	// Update the main Year field to the start year if not explicitly set
	if cfg.Year == time.Now().Year() {
		cfg.Year = cfg.StartYear
	}

	return nil
}

// GetYears returns a slice of years to generate planners for
func (cfg *Config) GetYears() []int {
	if cfg.StartYear == 0 || cfg.EndYear == 0 {
		// Fallback to single year
		return []int{cfg.Year}
	}

	years := make([]int, 0, cfg.EndYear-cfg.StartYear+1)
	for year := cfg.StartYear; year <= cfg.EndYear; year++ {
		years = append(years, year)
	}

	return years
}

// Composer is a function type for generating modules
type Composer func(cfg Config, tpls []string) (Modules, error)

// ComposerMap maps function names to their implementations
// This will be populated by the app package to avoid circular imports
var ComposerMap = map[string]Composer{}

// FilterUniqueModules removes duplicate modules based on template name
func FilterUniqueModules(array []Module) []Module {
	filtered := make([]Module, 0)
	found := map[string]bool{}

	for _, val := range array {
		if _, present := found[val.Tpl]; !present {
			filtered = append(filtered, val)
			found[val.Tpl] = true
		}
	}

	return filtered
}

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

// Task represents a single task from the CSV data
type Task struct {
	ID           string // * Added: Unique task identifier
	Name         string
	StartDate    time.Time
	EndDate      time.Time
	Category     string // * Fixed: Use Category instead of Priority for clarity
	Description  string
	Priority     int      // * Added: Separate priority field for task ordering
	Status       string   // * Added: Task status (Planned, In Progress, Completed, etc.)
	Assignee     string   // * Added: Task assignee
	ParentID     string   // * Added: Parent task ID for hierarchical relationships
	Dependencies []string // * Added: List of task IDs this task depends on
	IsMilestone  bool     // * Added: Whether this is a milestone task
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
		task.ID = getField("Task Name") // Fallback to name if no ID
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

// TaskCategory represents a task category with visual and organizational properties
type TaskCategory struct {
	Name        string
	DisplayName string
	Color       string
	Priority    int
	Description string
}

// Predefined task categories with their properties
var (
	CategoryPROPOSAL = TaskCategory{
		Name:        "PROPOSAL",
		DisplayName: "Proposal",
		Color:       "#4A90E2", // Blue
		Priority:    1,
		Description: "PhD proposal related tasks",
	}

	CategoryLASER = TaskCategory{
		Name:        "LASER",
		DisplayName: "Laser System",
		Color:       "#F5A623", // Orange
		Priority:    2,
		Description: "Laser system setup and maintenance",
	}

	CategoryIMAGING = TaskCategory{
		Name:        "IMAGING",
		DisplayName: "Imaging",
		Color:       "#7ED321", // Green
		Priority:    3,
		Description: "Imaging experiments and data collection",
	}

	CategoryADMIN = TaskCategory{
		Name:        "ADMIN",
		DisplayName: "Administrative",
		Color:       "#BD10E0", // Purple
		Priority:    4,
		Description: "Administrative tasks and paperwork",
	}

	CategoryDISSERTATION = TaskCategory{
		Name:        "DISSERTATION",
		DisplayName: "Dissertation",
		Color:       "#D0021B", // Red
		Priority:    5,
		Description: "Dissertation writing and defense",
	}

	CategoryRESEARCH = TaskCategory{
		Name:        "RESEARCH",
		DisplayName: "Research",
		Color:       "#50E3C2", // Teal
		Priority:    6,
		Description: "General research activities",
	}

	CategoryPUBLICATION = TaskCategory{
		Name:        "PUBLICATION",
		DisplayName: "Publication",
		Color:       "#B8E986", // Light Green
		Priority:    7,
		Description: "Publication and manuscript writing",
	}
)

// GetCategory returns the TaskCategory for a given category name
func GetCategory(categoryName string) TaskCategory {
	switch strings.ToUpper(categoryName) {
	case "PROPOSAL":
		return CategoryPROPOSAL
	case "LASER":
		return CategoryLASER
	case "IMAGING":
		return CategoryIMAGING
	case "ADMIN":
		return CategoryADMIN
	case "DISSERTATION":
		return CategoryDISSERTATION
	case "RESEARCH":
		return CategoryRESEARCH
	case "PUBLICATION":
		return CategoryPUBLICATION
	default:
		return TaskCategory{
			Name:        categoryName,
			DisplayName: categoryName,
			Color:       "#CCCCCC", // Default gray
			Priority:    99,
			Description: "Custom category",
		}
	}
}

// GetAllCategories returns all predefined categories
func GetAllCategories() []TaskCategory {
	return []TaskCategory{
		CategoryPROPOSAL,
		CategoryLASER,
		CategoryIMAGING,
		CategoryADMIN,
		CategoryDISSERTATION,
		CategoryRESEARCH,
		CategoryPUBLICATION,
	}
}

// TaskCollection represents a collection of tasks with efficient access patterns
type TaskCollection struct {
	tasks      []*Task
	byDate     []*Task
	byCategory map[string][]*Task
	byStatus   map[string][]*Task
	byAssignee map[string][]*Task
	sorted     bool
}

// NewTaskCollection creates a new empty task collection
func NewTaskCollection() *TaskCollection {
	return &TaskCollection{
		tasks:      make([]*Task, 0),
		byDate:     make([]*Task, 0),
		byCategory: make(map[string][]*Task),
		byStatus:   make(map[string][]*Task),
		byAssignee: make(map[string][]*Task),
		sorted:     false,
	}
}

// AddTask adds a task to the collection
func (tc *TaskCollection) AddTask(task *Task) {
	if task == nil {
		return
	}

	tc.tasks = append(tc.tasks, task)
	tc.byDate = append(tc.byDate, task)

	// Update category index
	if task.Category != "" {
		tc.byCategory[task.Category] = append(tc.byCategory[task.Category], task)
	}

	// Update status index
	if task.Status != "" {
		tc.byStatus[task.Status] = append(tc.byStatus[task.Status], task)
	}

	// Update assignee index
	if task.Assignee != "" {
		tc.byAssignee[task.Assignee] = append(tc.byAssignee[task.Assignee], task)
	}

	tc.sorted = false
}

// GetTask retrieves a task by name (since we removed ID)
func (tc *TaskCollection) GetTask(name string) (*Task, bool) {
	for _, task := range tc.tasks {
		if task.Name == name {
			return task, true
		}
	}
	return nil, false
}

// GetAllTasks returns all tasks in the collection
func (tc *TaskCollection) GetAllTasks() []*Task {
	return tc.tasks
}

// GetTasksByCategory returns all tasks in a specific category
func (tc *TaskCollection) GetTasksByCategory(category string) []*Task {
	return tc.byCategory[category]
}

// GetTasksByStatus returns all tasks with a specific status
func (tc *TaskCollection) GetTasksByStatus(status string) []*Task {
	return tc.byStatus[status]
}

// GetTasksByAssignee returns all tasks assigned to a specific person
func (tc *TaskCollection) GetTasksByAssignee(assignee string) []*Task {
	return tc.byAssignee[assignee]
}

// GetTasksByDateRange returns all tasks within a date range
func (tc *TaskCollection) GetTasksByDateRange(start, end time.Time) []*Task {
	if !tc.sorted {
		tc.sortByDate()
	}

	var result []*Task
	for _, task := range tc.byDate {
		if task.OverlapsWithDateRange(start, end) {
			result = append(result, task)
		}
	}
	return result
}

// GetTasksByDate returns all tasks on a specific date
func (tc *TaskCollection) GetTasksByDate(date time.Time) []*Task {
	if !tc.sorted {
		tc.sortByDate()
	}

	var result []*Task
	for _, task := range tc.byDate {
		if task.IsOnDate(date) {
			result = append(result, task)
		}
	}
	return result
}

// sortByDate sorts tasks by start date
func (tc *TaskCollection) sortByDate() {
	sort.Slice(tc.byDate, func(i, j int) bool {
		return tc.byDate[i].StartDate.Before(tc.byDate[j].StartDate)
	})
	tc.sorted = true
}

// TaskHierarchy represents the parent-child hierarchy of tasks
type TaskHierarchy struct {
	roots    []*Task
	parents  map[string]*Task
	children map[string][]*Task
	tasks    []*Task
}

// NewTaskHierarchy creates a new task hierarchy
func NewTaskHierarchy() *TaskHierarchy {
	return &TaskHierarchy{
		roots:    make([]*Task, 0),
		parents:  make(map[string]*Task),
		children: make(map[string][]*Task),
		tasks:    make([]*Task, 0),
	}
}

// AddTask adds a task to the hierarchy
func (th *TaskHierarchy) AddTask(task *Task) {
	if task == nil {
		return
	}

	th.tasks = append(th.tasks, task)

	if task.ParentID == "" {
		// This is a root task
		th.roots = append(th.roots, task)
	} else {
		// This is a child task - find parent by name
		for _, parent := range th.tasks {
			if parent.Name == task.ParentID {
				th.parents[task.Name] = parent
				th.children[task.ParentID] = append(th.children[task.ParentID], task)
				break
			}
		}
	}
}

// GetRootTasks returns all root tasks (tasks without parents)
func (th *TaskHierarchy) GetRootTasks() []*Task {
	return th.roots
}

// GetChildren returns all child tasks of a given task
func (th *TaskHierarchy) GetChildren(taskName string) []*Task {
	return th.children[taskName]
}

// GetParent returns the parent task of a given task
func (th *TaskHierarchy) GetParent(taskName string) *Task {
	return th.parents[taskName]
}

// GetAncestors returns all ancestor tasks of a given task
func (th *TaskHierarchy) GetAncestors(taskName string) []*Task {
	var ancestors []*Task
	current := th.GetParent(taskName)

	for current != nil {
		ancestors = append(ancestors, current)
		current = th.GetParent(current.Name)
	}

	return ancestors
}

// GetDescendants returns all descendant tasks of a given task
func (th *TaskHierarchy) GetDescendants(taskName string) []*Task {
	var descendants []*Task
	children := th.GetChildren(taskName)

	for _, child := range children {
		descendants = append(descendants, child)
		descendants = append(descendants, th.GetDescendants(child.Name)...)
	}

	return descendants
}

// CalendarLayout represents optimized date range calculations for calendar rendering
type CalendarLayout struct {
	startDate time.Time
	endDate   time.Time
	tasks     []*Task
	months    []MonthYear
	weeks     []time.Time
}

// NewCalendarLayout creates a new calendar layout for a date range
func NewCalendarLayout(startDate, endDate time.Time, tasks []*Task) *CalendarLayout {
	cl := &CalendarLayout{
		startDate: startDate,
		endDate:   endDate,
		tasks:     tasks,
	}

	cl.generateMonths()
	cl.generateWeeks()

	return cl
}

// generateMonths generates all months in the date range
func (cl *CalendarLayout) generateMonths() {
	cl.months = make([]MonthYear, 0)

	current := time.Date(cl.startDate.Year(), cl.startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(cl.endDate.Year(), cl.endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !current.After(end) {
		cl.months = append(cl.months, MonthYear{
			Year:  current.Year(),
			Month: current.Month(),
		})
		current = current.AddDate(0, 1, 0)
	}
}

// generateWeeks generates all weeks in the date range
func (cl *CalendarLayout) generateWeeks() {
	cl.weeks = make([]time.Time, 0)

	// Find the start of the first week
	start := cl.startDate
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, -1)
	}

	// Generate weeks until we cover the end date
	for !start.After(cl.endDate) {
		cl.weeks = append(cl.weeks, start)
		start = start.AddDate(0, 0, 7)
	}
}

// GetMonths returns all months in the layout
func (cl *CalendarLayout) GetMonths() []MonthYear {
	return cl.months
}

// GetWeeks returns all weeks in the layout
func (cl *CalendarLayout) GetWeeks() []time.Time {
	return cl.weeks
}

// GetTasksForMonth returns all tasks that occur in a specific month
func (cl *CalendarLayout) GetTasksForMonth(year int, month time.Month) []*Task {
	var result []*Task
	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, -1)

	for _, task := range cl.tasks {
		if task.OverlapsWithDateRange(monthStart, monthEnd) {
			result = append(result, task)
		}
	}

	return result
}

// GetTasksForWeek returns all tasks that occur in a specific week
func (cl *CalendarLayout) GetTasksForWeek(weekStart time.Time) []*Task {
	var result []*Task
	weekEnd := weekStart.AddDate(0, 0, 6)

	for _, task := range cl.tasks {
		if task.OverlapsWithDateRange(weekStart, weekEnd) {
			result = append(result, task)
		}
	}

	return result
}

// TaskRenderer represents visual rendering properties for tasks
type TaskRenderer struct {
	TaskID      string
	X           float64 // X position in calendar
	Y           float64 // Y position in calendar
	Width       float64 // Width in calendar
	Height      float64 // Height in calendar
	Color       string  // Task color
	BorderColor string  // Border color
	Opacity     float64 // Opacity (0.0 to 1.0)
	Visible     bool    // Whether task is visible
	ZIndex      int     // Rendering order
}

// NewTaskRenderer creates a new task renderer
func NewTaskRenderer(task *Task) *TaskRenderer {
	category := GetCategory(task.Category)

	return &TaskRenderer{
		TaskID:      task.Name,
		Color:       category.Color,
		BorderColor: "#000000",
		Opacity:     1.0,
		Visible:     true,
		ZIndex:      category.Priority,
	}
}

// Enhanced Task methods for calendar layout and rendering

// IsOnDate checks if a task occurs on a specific date
func (t *Task) IsOnDate(date time.Time) bool {
	if t.StartDate.IsZero() || t.EndDate.IsZero() {
		return false
	}

	// Normalize dates to compare only the date part
	taskStart := time.Date(t.StartDate.Year(), t.StartDate.Month(), t.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	taskEnd := time.Date(t.EndDate.Year(), t.EndDate.Month(), t.EndDate.Day(), 0, 0, 0, 0, time.UTC)
	checkDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	return !checkDate.Before(taskStart) && !checkDate.After(taskEnd)
}

// OverlapsWithDateRange checks if a task overlaps with a date range
func (t *Task) OverlapsWithDateRange(start, end time.Time) bool {
	if t.StartDate.IsZero() || t.EndDate.IsZero() {
		return false
	}

	// Task overlaps if it starts before the range ends and ends after the range starts
	return !t.StartDate.After(end) && !t.EndDate.Before(start)
}

// GetDuration returns the duration of the task in days
func (t *Task) GetDuration() int {
	if t.StartDate.IsZero() || t.EndDate.IsZero() {
		return 0
	}

	duration := t.EndDate.Sub(t.StartDate)
	return int(duration.Hours()/24) + 1 // +1 to include both start and end days
}

// GetCategoryInfo returns the TaskCategory for this task
func (t *Task) GetCategoryInfo() TaskCategory {
	return GetCategory(t.Category)
}

// IsOverdue checks if the task is overdue
func (t *Task) IsOverdue() bool {
	if t.EndDate.IsZero() {
		return false
	}

	now := time.Now()
	return now.After(t.EndDate) && t.Status != "Completed"
}

// IsUpcoming checks if the task is starting soon (within 7 days)
func (t *Task) IsUpcoming() bool {
	if t.StartDate.IsZero() {
		return false
	}

	now := time.Now()
	sevenDaysFromNow := now.AddDate(0, 0, 7)
	return t.StartDate.After(now) && t.StartDate.Before(sevenDaysFromNow)
}

// GetProgressPercentage returns the progress percentage based on dates
func (t *Task) GetProgressPercentage() float64 {
	if t.StartDate.IsZero() || t.EndDate.IsZero() {
		return 0.0
	}

	now := time.Now()
	totalDuration := t.EndDate.Sub(t.StartDate)
	elapsed := now.Sub(t.StartDate)

	if elapsed < 0 {
		return 0.0
	}

	if elapsed >= totalDuration {
		return 100.0
	}

	return (elapsed.Hours() / totalDuration.Hours()) * 100.0
}

// String returns a string representation of the task
func (t *Task) String() string {
	return fmt.Sprintf("Task[%s (%s) %s - %s]",
		t.Name, t.Category,
		t.StartDate.Format("2006-01-02"),
		t.EndDate.Format("2006-01-02"))
}

// DataValidationError represents a validation error with detailed context
type DataValidationError struct {
	Type        string    // Error type (e.g., "DATE_RANGE", "CONFLICT", "DEPENDENCY")
	TaskID      string    // Task ID that has the error
	Field       string    // Field that has the error
	Value       string    // Value that caused the error
	Message     string    // Human-readable error message
	Severity    string    // Error severity (ERROR, WARNING, INFO)
	Timestamp   time.Time // When the error was detected
	Suggestions []string  // Suggested fixes
}

// Error returns the error message
func (ve *DataValidationError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", ve.Severity, ve.TaskID, ve.Message)
}

// DateValidator handles date range validation and conflict detection
type DateValidator struct {
	workDays   map[time.Weekday]bool
	holidays   []time.Time
	timezone   *time.Location
	strictMode bool
}

// NewDateValidator creates a new date validator
func NewDateValidator() *DateValidator {
	return &DateValidator{
		workDays: map[time.Weekday]bool{
			time.Monday:    true,
			time.Tuesday:   true,
			time.Wednesday: true,
			time.Thursday:  true,
			time.Friday:    true,
			time.Saturday:  false,
			time.Sunday:    false,
		},
		holidays:   []time.Time{},
		timezone:   time.UTC,
		strictMode: false,
	}
}

// ValidateDateRanges validates date ranges for a slice of tasks
func (dv *DateValidator) ValidateDateRanges(tasks []*Task) []DataValidationError {
	var errors []DataValidationError

	for _, task := range tasks {
		// Basic date validation
		if task.StartDate.After(task.EndDate) {
			errors = append(errors, DataValidationError{
				Type:        "DATE_RANGE",
				TaskID:      task.ID,
				Field:       "dates",
				Value:       fmt.Sprintf("%s - %s", task.StartDate.Format("2006-01-02"), task.EndDate.Format("2006-01-02")),
				Message:     "Start date is after end date",
				Severity:    "ERROR",
				Timestamp:   time.Now(),
				Suggestions: []string{"Correct the start or end date"},
			})
		}
	}

	return errors
}
