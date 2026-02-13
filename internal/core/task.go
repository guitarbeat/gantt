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
	Phase        string // * Combined: Phase with description (e.g., "1: Project Metadata")
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
		Color:       GenerateCategoryColor("PROPOSAL"),
		Description: "PhD proposal related tasks",
	}

	CategoryLASER = TaskCategory{
		Name:        "LASER",
		DisplayName: "Laser System",
		Color:       GenerateCategoryColor("LASER"),
		Description: "Laser system setup and maintenance",
	}

	CategoryIMAGING = TaskCategory{
		Name:        "IMAGING",
		DisplayName: "Imaging",
		Color:       GenerateCategoryColor("IMAGING"),
		Description: "Imaging experiments and data collection",
	}

	CategoryADMIN = TaskCategory{
		Name:        "ADMIN",
		DisplayName: "Administrative",
		Color:       GenerateCategoryColor("ADMIN"),
		Description: "Administrative tasks and paperwork",
	}

	CategoryDISSERTATION = TaskCategory{
		Name:        "DISSERTATION",
		DisplayName: "Dissertation",
		Color:       GenerateCategoryColor("DISSERTATION"),
		Description: "Dissertation writing and defense",
	}

	CategoryRESEARCH = TaskCategory{
		Name:        "RESEARCH",
		DisplayName: "Research",
		Color:       GenerateCategoryColor("RESEARCH"),
		Description: "General research activities",
	}

	CategoryPUBLICATION = TaskCategory{
		Name:        "PUBLICATION",
		DisplayName: "Publication",
		Color:       GenerateCategoryColor("PUBLICATION"),
		Description: "Publication and manuscript writing",
	}
)

// GenerateCategoryColor creates a consistent, visually distinct color based on the category name
func GenerateCategoryColor(category string) string {
	// Dynamic color assignment using golden angle for maximum visual distinction
	// This ensures each unique category gets a unique, well-distributed color

	// Normalize category name for consistency
	normalizedCategory := strings.ToUpper(strings.TrimSpace(category))

	// Create a hash of the category name for consistent color assignment
	hash := 0
	for i, char := range normalizedCategory {
		hash = hash*31 + int(char) + i*7 // Improved hash distribution
	}
	if hash < 0 {
		hash = -hash
	}

	// Use golden angle approximation (137.5 degrees) for optimal color distribution
	// This spreads colors evenly around the color wheel
	hue := float64(hash%360) * 137.5
	hue = hue - float64(int(hue/360.0)*360) // Keep hue in 0-360 range

	// Optimized saturation and lightness for accessibility and visual appeal
	saturation := 0.75 // High saturation for vibrancy
	lightness := 0.65  // Balanced lightness for good contrast

	// Convert HSL to RGB
	r, g, b := hslToRgb(hue, saturation, lightness)

	// Convert to hex
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// hslToRgb converts HSL color values to RGB
func hslToRgb(hue, saturation, lightness float64) (int, int, int) {
	// Normalize hue to 0-1 range
	normalizedHue := hue / 360.0

	var red, green, blue float64

	if saturation == 0 {
		// Grayscale
		red, green, blue = lightness, lightness, lightness
	} else {
		var chroma1, chroma2 float64
		if lightness < 0.5 {
			chroma1 = lightness * (1 + saturation)
		} else {
			chroma1 = lightness + saturation - lightness*saturation
		}
		chroma2 = 2*lightness - chroma1

		red = hueToRgb(chroma2, chroma1, normalizedHue+1.0/3.0)
		green = hueToRgb(chroma2, chroma1, normalizedHue)
		blue = hueToRgb(chroma2, chroma1, normalizedHue-1.0/3.0)
	}

	return int(red * 255), int(green * 255), int(blue * 255)
}

// hueToRgb helper function for HSL to RGB conversion
func hueToRgb(chroma2, chroma1, tempHue float64) float64 {
	if tempHue < 0 {
		tempHue += 1
	}
	if tempHue > 1 {
		tempHue -= 1
	}
	if tempHue < 1.0/6.0 {
		return chroma2 + (chroma1-chroma2)*6*tempHue
	}
	if tempHue < 1.0/2.0 {
		return chroma1
	}
	if tempHue < 2.0/3.0 {
		return chroma2 + (chroma1-chroma2)*(2.0/3.0-tempHue)*6
	}
	return chroma2
}

// GetCategory returns the TaskCategory for a given category name
// GetAllCategories returns all predefined categories
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

// Enhanced Task methods for calendar layout and rendering

// CalculateDateRange calculates the earliest and latest dates from a list of tasks
func CalculateDateRange(tasks []Task) DateRange {
	if len(tasks) == 0 {
		return DateRange{
			Earliest: time.Now(),
			Latest:   time.Now(),
		}
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

	return DateRange{
		Earliest: earliest,
		Latest:   latest,
	}
}

// GetMonthsWithTasks returns a list of unique months that have tasks
func GetMonthsWithTasks(tasks []Task, dateRange DateRange) []MonthYear {
	monthMap := make(map[string]MonthYear)

	for _, task := range tasks {
		// Add all months between start and end date
		current := time.Date(task.StartDate.Year(), task.StartDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(task.EndDate.Year(), task.EndDate.Month(), 1, 0, 0, 0, 0, time.UTC)

		for !current.After(end) {
			key := fmt.Sprintf("%d-%02d", current.Year(), current.Month())
			monthMap[key] = MonthYear{
				Year:  current.Year(),
				Month: current.Month(),
			}
			current = current.AddDate(0, 1, 0)
		}
	}

	// Convert map to sorted slice
	var months []MonthYear
	for _, month := range monthMap {
		months = append(months, month)
	}

	// Sort by year and month
	for i := 0; i < len(months)-1; i++ {
		for j := i + 1; j < len(months); j++ {
			if months[i].Year > months[j].Year ||
				(months[i].Year == months[j].Year && months[i].Month > months[j].Month) {
				months[i], months[j] = months[j], months[i]
			}
		}
	}

	return months
}
