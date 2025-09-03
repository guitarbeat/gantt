package calendar

import (
	"fmt"
	"time"
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
}

// CreateSpanningTask creates a new spanning task from a calendar task and date range
func CreateSpanningTask(task Task, startDate, endDate time.Time) SpanningTask {
	return SpanningTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Category:    task.Category,
		StartDate:   startDate,
		EndDate:     endDate,
		Color:       getTaskColor(task.Category),
		Priority:    getTaskPriority(task.Category),
	}
}

// GetTaskOverlayForDay returns the appropriate overlay for a spanning task on a specific day
func GetTaskOverlayForDay(dayTime time.Time, task SpanningTask) string {
	// Check if this day falls within the task span
	if dayTime.Before(task.StartDate) || dayTime.After(task.EndDate) {
		return ""
	}

	// Create a simple colored background with task name
	// This creates a continuous visual element like Google Calendar
	color := task.Color
	taskName := task.Name
	
	// Escape LaTeX special characters
	taskName = escapeLaTeX(taskName)
	
	// Create a simple colored background that doesn't interfere with date numbers
	overlay := fmt.Sprintf(`\colorbox{%s!20}{\parbox{\linewidth}{\tiny %s}}`, color, taskName)

	return overlay
}

// getTaskColor returns a color based on task category
func getTaskColor(category string) string {
	switch category {
	case "Planning":
		return "blue"
	case "Research":
		return "green"
	case "Development":
		return "orange"
	case "Review":
		return "purple"
	case "Meeting":
		return "red"
	default:
		return "gray"
	}
}

// getTaskPriority returns a priority number based on task category
func getTaskPriority(category string) int {
	switch category {
	case "Planning":
		return 1
	case "Research":
		return 2
	case "Development":
		return 3
	case "Review":
		return 4
	case "Meeting":
		return 5
	default:
		return 6
	}
}

// escapeLaTeX escapes special LaTeX characters
func escapeLaTeX(text string) string {
	// Replace common LaTeX special characters
	replacements := map[string]string{
		"\\": "\\textbackslash{}",
		"{":  "\\{",
		"}":  "\\}",
		"$":  "\\$",
		"&":  "\\&",
		"%":  "\\%",
		"#":  "\\#",
		"^":  "\\textasciicircum{}",
		"_":  "\\_",
		"~":  "\\textasciitilde{}",
	}

	result := text
	for old, new := range replacements {
		result = replaceAll(result, old, new)
	}
	return result
}

// replaceAll replaces all occurrences of old with new in text
func replaceAll(text, old, new string) string {
	result := text
	for {
		index := findIndex(result, old)
		if index == -1 {
			break
		}
		result = result[:index] + new + result[index+len(old):]
	}
	return result
}

// findIndex finds the index of the first occurrence of substr in text
func findIndex(text, substr string) int {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetSpanningTasksForMonth returns all spanning tasks that overlap with a given month
func GetSpanningTasksForMonth(month *Month, tasks []SpanningTask) []SpanningTask {
	var monthTasks []SpanningTask
	
	monthStart := time.Date(month.Year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
	monthEnd := monthStart.AddDate(0, 1, -1)
	
	for _, task := range tasks {
		// Check if task overlaps with the month
		if task.StartDate.Before(monthEnd.AddDate(0, 0, 1)) && task.EndDate.After(monthStart.AddDate(0, 0, -1)) {
			monthTasks = append(monthTasks, task)
		}
	}
	
	return monthTasks
}

// ApplySpanningTasksToMonth applies spanning tasks to a month's calendar
func ApplySpanningTasksToMonth(month *Month, tasks []SpanningTask) {
	monthTasks := GetSpanningTasksForMonth(month, tasks)
	
	for _, week := range month.Weeks {
		for i := range week.Days {
			day := &week.Days[i]
			if day.Time.IsZero() {
				continue
			}
			
			// Find the highest priority task for this day
			var bestTask *SpanningTask
			bestPriority := 999
			
			for _, task := range monthTasks {
				if day.Time.After(task.StartDate.AddDate(0, 0, -1)) && day.Time.Before(task.EndDate.AddDate(0, 0, 1)) {
					if task.Priority < bestPriority {
						bestTask = &task
						bestPriority = task.Priority
					}
				}
			}
			
			// Apply the best task overlay to this day
			if bestTask != nil {
				overlay := GetTaskOverlayForDay(day.Time, *bestTask)
				if overlay != "" {
					// Store the overlay in a way that doesn't interfere with date display
					day.SpanningTask = bestTask
				}
			}
		}
	}
}
