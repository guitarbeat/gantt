package calendar

import (
	"strings"
	"time"
)

// * Task processing and utility functions

// TasksForDay returns a formatted string of tasks for this day
func (d Day) TasksForDay() string {
	if len(d.Tasks) == 0 {
		return ""
	}
	var taskStrings []string
	for _, task := range d.Tasks {
		// Only show task name, category is only used for color
		taskStr := d.escapeLatexSpecialChars(task.Name)
		
		// Add star for milestone tasks
		if d.isMilestoneTask(task) {
			taskStr = "★ " + taskStr
		}
		
		taskStrings = append(taskStrings, taskStr)
	}
	return strings.Join(taskStrings, "\\\\")
}

// getDayDate returns the day date normalized to UTC midnight
func (d Day) getDayDate() time.Time {
	return time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), 0, 0, 0, 0, time.UTC)
}

// getTaskStartDate returns the task start date normalized to UTC midnight
func (d Day) getTaskStartDate(task *SpanningTask) time.Time {
	return time.Date(task.StartDate.Year(), task.StartDate.Month(), task.StartDate.Day(), 0, 0, 0, 0, time.UTC)
}

// getTaskEndDate returns the task end date normalized to UTC midnight
func (d Day) getTaskEndDate(task *SpanningTask) time.Time {
	return time.Date(task.EndDate.Year(), task.EndDate.Month(), task.EndDate.Day(), 0, 0, 0, 0, time.UTC)
}

// isTaskActiveOnDay checks if a task is active on the given day
func (d Day) isTaskActiveOnDay(dayDate, start, end time.Time) bool {
	return !dayDate.Before(start) && !dayDate.After(end)
}

// calculateTaskSpanColumns calculates how many columns a task should span
func (d Day) calculateTaskSpanColumns(dayDate, end time.Time) int {
	idxMonFirst := (int(dayDate.Weekday()) + 6) % 7 // Monday=0
	remainInRow := 7 - idxMonFirst
	totalRemain := int(end.Sub(dayDate).Hours()/24) + 1
	if totalRemain < 1 {
		totalRemain = 1
	}
	overlayCols := totalRemain
	if overlayCols > remainInRow {
		overlayCols = remainInRow
	}
	return overlayCols
}

// findStartingTasks finds tasks that start on the given day and calculates max columns
func (d Day) findStartingTasks(dayDate time.Time) ([]*SpanningTask, int) {
	var startingTasks []*SpanningTask
	var maxCols int

	for _, task := range d.SpanningTasks {
		start := d.getTaskStartDate(task)
		end := d.getTaskEndDate(task)

		if !d.isTaskActiveOnDay(dayDate, start, end) {
			continue
		}

		if dayDate.Equal(start) {
			startingTasks = append(startingTasks, task)
			cols := d.calculateTaskSpanColumns(dayDate, end)
			if cols > maxCols {
				maxCols = cols
			}
		}
	}

	return startingTasks, maxCols
}

// sortTasksByPriority sorts tasks by category priority for better visual organization
func (d Day) sortTasksByPriority(tasks []*SpanningTask) []*SpanningTask {
	sorted := make([]*SpanningTask, len(tasks))
	copy(sorted, tasks)

	priorityOrder := d.getCategoryPriorityOrder()

	// Simple bubble sort by priority
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			priority1 := d.getTaskPriority(sorted[j].Category, priorityOrder)
			priority2 := d.getTaskPriority(sorted[j+1].Category, priorityOrder)
			if priority1 > priority2 {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// getCategoryPriorityOrder returns the priority order for task categories
func (d Day) getCategoryPriorityOrder() map[string]int {
	return map[string]int{
		"DISSERTATION": 1,
		"PROPOSAL":     2,
		"PUBLICATION":  3,
		"RESEARCH":     4,
		"IMAGING":      5,
		"LASER":        6,
		"ADMIN":        7,
	}
}

// getTaskPriority returns the priority for a task category
func (d Day) getTaskPriority(category string, priorityOrder map[string]int) int {
	if priority, exists := priorityOrder[category]; exists {
		return priority
	}
	return 99 // Unknown categories go last
}

// isMilestoneTask checks if a task is a milestone based on its description
func (d Day) isMilestoneTask(task Task) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(task.Description)), "MILESTONE:")
}

// isMilestoneSpanningTask checks if a spanning task is a milestone based on its description
func (d Day) isMilestoneSpanningTask(task *SpanningTask) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(task.Description)), "MILESTONE:")
}

// escapeLatexSpecialChars escapes special LaTeX characters in text
func (d Day) escapeLatexSpecialChars(text string) string {
	// Replace special LaTeX characters with their escaped versions
	text = strings.ReplaceAll(text, "\\", "\\textbackslash{}")
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, "$", "\\$")
	text = strings.ReplaceAll(text, "&", "\\&")
	text = strings.ReplaceAll(text, "%", "\\%")
	text = strings.ReplaceAll(text, "#", "\\#")
	text = strings.ReplaceAll(text, "^", "\\textasciicircum{}")
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "~", "\\textasciitilde{}")
	return text
}

// smartTruncateText intelligently truncates text at word boundaries when possible
func (d Day) smartTruncateText(text string, maxChars int) string {
	if len(text) <= maxChars {
		return text
	}
	
	// Try to break at word boundaries
	if maxChars > 8 {
		words := strings.Fields(text)
		result := ""
		for _, word := range words {
			if len(result)+len(word)+1 <= maxChars-3 {
				if result != "" {
					result += " "
				}
				result += word
			} else {
				break
			}
		}
		if result != "" {
			return result + "..."
		}
	}
	
	// Fallback to simple truncation
	return text[:maxChars-3] + "..."
}
