package common

import (
	"fmt"
	"strings"
	"time"
)

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
