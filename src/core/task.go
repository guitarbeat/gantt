package core

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
	Description string
}

// Predefined task categories with their properties
var (
	CategoryPROPOSAL = TaskCategory{
		Name:        "PROPOSAL",
		DisplayName: "Proposal",
		Color:       generateCategoryColor("PROPOSAL"),
		Description: "PhD proposal related tasks",
	}

	CategoryLASER = TaskCategory{
		Name:        "LASER",
		DisplayName: "Laser System",
		Color:       generateCategoryColor("LASER"),
		Description: "Laser system setup and maintenance",
	}

	CategoryIMAGING = TaskCategory{
		Name:        "IMAGING",
		DisplayName: "Imaging",
		Color:       generateCategoryColor("IMAGING"),
		Description: "Imaging experiments and data collection",
	}

	CategoryADMIN = TaskCategory{
		Name:        "ADMIN",
		DisplayName: "Administrative",
		Color:       generateCategoryColor("ADMIN"),
		Description: "Administrative tasks and paperwork",
	}

	CategoryDISSERTATION = TaskCategory{
		Name:        "DISSERTATION",
		DisplayName: "Dissertation",
		Color:       generateCategoryColor("DISSERTATION"),
		Description: "Dissertation writing and defense",
	}

	CategoryRESEARCH = TaskCategory{
		Name:        "RESEARCH",
		DisplayName: "Research",
		Color:       generateCategoryColor("RESEARCH"),
		Description: "General research activities",
	}

	CategoryPUBLICATION = TaskCategory{
		Name:        "PUBLICATION",
		DisplayName: "Publication",
		Color:       generateCategoryColor("PUBLICATION"),
		Description: "Publication and manuscript writing",
	}
)

// generateCategoryColor creates a consistent, visually distinct color based on the category name
func generateCategoryColor(category string) string {
	// Use a better hash function to generate consistent colors
	hash := 0
	for i, char := range category {
		hash = hash*31 + int(char) + i*7 // Add position to improve distribution
	}

	// Convert hash to a positive number
	if hash < 0 {
		hash = -hash
	}

	// Generate HSL color with good saturation and lightness for readability
	hue := float64(hash%360)                    // 0-360 degrees
	saturation := 0.7 + float64(hash%30)/100.0 // 0.7-1.0 for good saturation
	lightness := 0.5 + float64(hash%20)/100.0  // 0.5-0.7 for good contrast

	// Convert HSL to RGB
	r, g, b := hslToRgb(hue, saturation, lightness)
	
	// Convert to hex
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// hslToRgb converts HSL color values to RGB
func hslToRgb(h, s, l float64) (int, int, int) {
	// Normalize values
	h = h / 360.0
	
	var r, g, b float64
	
	if s == 0 {
		// Grayscale
		r, g, b = l, l, l
	} else {
		var q, p float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p = 2*l - q
		
		r = hueToRgb(p, q, h+1.0/3.0)
		g = hueToRgb(p, q, h)
		b = hueToRgb(p, q, h-1.0/3.0)
	}
	
	return int(r * 255), int(g * 255), int(b * 255)
}

// hueToRgb helper function for HSL to RGB conversion
func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

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
