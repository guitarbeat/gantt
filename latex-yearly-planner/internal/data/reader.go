package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

const (
	// DateFormat is the expected date format in CSV files
	DateFormat = "2006-01-02"
)

// Task represents a single task from the CSV data
type Task struct {
	ID          string
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	Priority    string
	Description string
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
}

// NewReader creates a new CSV data reader
func NewReader(filePath string) *Reader {
	return &Reader{
		filePath: filePath,
	}
}

// ReadTasks reads all tasks from the CSV file
func (r *Reader) ReadTasks() ([]Task, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create field index map
	fieldIndex := make(map[string]int)
	for i, field := range header {
		fieldIndex[field] = i
	}

	var tasks []Task
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Skip empty rows
		if len(record) == 0 || record[0] == "" {
			continue
		}

		task, err := r.parseTask(record, fieldIndex)
		if err != nil {
			// Log error but continue processing other tasks
			log.Printf("Warning: failed to parse task: %v", err)
			continue
		}

		// Import all tasks (removed filter for task A only)

		tasks = append(tasks, task)
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

// parseTask parses a single CSV record into a Task struct
func (r *Reader) parseTask(record []string, fieldIndex map[string]int) (Task, error) {
	task := Task{}

	// Helper function to get field value safely
	getField := func(fieldName string) string {
		if index, exists := fieldIndex[fieldName]; exists && index < len(record) {
			return record[index]
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

	// Parse category from CSV (using Category field)
	if category := getField("Category"); category != "" {
		task.Priority = category // Using Priority field to store category
	}

	// Parse dates
	startDateStr := getField("Start Date")
	if startDateStr != "" {
		startDate, err := time.Parse(DateFormat, startDateStr)
		if err != nil {
			return task, fmt.Errorf("invalid start date format: %s", startDateStr)
		}
		task.StartDate = startDate
	}

	endDateStr := getField("Due Date")
	if endDateStr != "" {
		endDate, err := time.Parse(DateFormat, endDateStr)
		if err != nil {
			return task, fmt.Errorf("invalid end date format: %s", endDateStr)
		}
		task.EndDate = endDate
	}

	return task, nil
}
