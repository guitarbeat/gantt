package calendar

import (
	"time"

	"latex-yearly-planner/internal/data"
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
	// * Fixed: Use Category field instead of Priority
	color := getColorForCategory(task.Category)

	return SpanningTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Category:    task.Category, // * Fixed: Use Category field
		StartDate:   startDate,
		EndDate:     endDate,
		Color:       color,
		Priority:    task.Priority, // * Fixed: Use actual Priority field
		Progress:    0,             // Default progress
		Status:      task.Status,   // * Fixed: Use actual Status field
		Assignee:    task.Assignee, // * Fixed: Use actual Assignee field
	}
}

// ApplySpanningTasksToMonth applies spanning tasks to a month
func ApplySpanningTasksToMonth(month *Month, tasks []SpanningTask) {
	// Apply spanning tasks to the appropriate days in the month
	for taskIndex, task := range tasks {
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
							// Create a copy of the task to avoid pointer issues
							taskCopy := tasks[taskIndex]
							// Add the spanning task to this day
							week.Days[i].SpanningTasks = append(week.Days[i].SpanningTasks, &taskCopy)
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
		"RESEARCH":      "green",
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
