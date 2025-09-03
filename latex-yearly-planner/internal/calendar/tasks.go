package calendar

import (
	"fmt"
	"time"

	"github.com/kudrykv/latex-yearly-planner/internal/data"
)

// SpanningTask represents a task that spans multiple days
type SpanningTask struct {
	ID          string
	Name        string
	Description string
	Category    string
	StartDate   time.Time
	EndDate     time.Time
	Color       string
	Priority    int
	Progress    int    // Progress percentage (0-100)
	Status      string // Task status
	Assignee    string // Task assignee
}

// CreateSpanningTask creates a new spanning task from basic task data
func CreateSpanningTask(task data.Task, startDate, endDate time.Time) SpanningTask {
	// Assign color based on category (stored in Priority field)
	color := getColorForCategory(task.Priority)

	return SpanningTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Category:    task.Priority, // Category is stored in Priority field
		StartDate:   startDate,
		EndDate:     endDate,
		Color:       color,
		Priority:    1,
		Progress:    task.Progress,
		Status:      task.Status,
		Assignee:    task.Assignee,
	}
}

// GetTaskOverlayForDay returns the overlay information for a task on a specific day
func GetTaskOverlayForDay(day time.Time, task SpanningTask) string {
	// Check if the day falls within the task span
	if day.Before(task.StartDate) || day.After(task.EndDate) {
		return ""
	}
	// For now, return a simple overlay
	return fmt.Sprintf("\\simpletaskoverlay{%s}{%s}{%s}", task.Color, task.Name, task.Category)
}

// ApplySpanningTasksToMonth applies spanning tasks to a month
func ApplySpanningTasksToMonth(month *Month, tasks []SpanningTask) {
	// Apply spanning tasks to the appropriate days in the month
	for _, task := range tasks {
		// Find all days in the month that this task spans
		current := task.StartDate
		for !current.After(task.EndDate) {
			// Check if this day is in the current month
			if current.Month() == month.Month && current.Year() == month.Year.Number {
				// Find the day in the month and set the spanning task
				for _, week := range month.Weeks {
					for i := range week.Days {
						if week.Days[i].Time.Day() == current.Day() &&
							week.Days[i].Time.Month() == current.Month() &&
							week.Days[i].Time.Year() == current.Year() {
							// Add the spanning task to this day
							week.Days[i].SpanningTasks = append(week.Days[i].SpanningTasks, &task)
							break
						}
					}
				}
			}
			current = current.AddDate(0, 0, 1)
		}
	}
}

// getColorForCategory returns a color for the given category
func getColorForCategory(category string) string {
	colorMap := map[string]string{
		"PROPOSAL":      "blue",
		"ADMIN":         "gray",
		"LASER":         "red",
		"IMAGING":       "green",
		"PUBLICATION":   "purple",
		"DISSERTATION":  "orange",
		"Planning":      "blue",
		"Research":      "green",
		"Development":   "orange",
		"Testing":       "red",
		"Documentation": "purple",
		"Meeting":       "teal",
		"Review":        "brown",
		"Default":       "gray",
	}
	if color, exists := colorMap[category]; exists {
		return color
	}
	return colorMap["Default"]
}
