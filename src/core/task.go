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
	Phase        string   // * Added: Phase (1, 2, 3, 4)
	SubPhase     string   // * Added: Sub-Phase (detailed phase description)
	Category     string   // * Fixed: Use Category instead of Priority for clarity
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
