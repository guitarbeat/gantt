// Package calendar provides task rendering functionality for the PhD dissertation planner.
//
// This module handles:
// - Task pill rendering with proper colors and formatting
// - Spanning task overlay generation
// - Task content formatting and LaTeX escaping
package calendar

import (
	"fmt"
	"strings"

	"phd-dissertation-planner/src/core"
)

// TaskRenderer handles rendering of tasks and spanning tasks
type TaskRenderer struct {
	cfg *core.Config
}

// NewTaskRenderer creates a new task renderer with the given configuration
func NewTaskRenderer(cfg *core.Config) *TaskRenderer {
	return &TaskRenderer{cfg: cfg}
}

// TaskOverlayInfo represents a spanning task overlay
type TaskOverlayInfo struct {
	content string
	cols    int
}

// RenderSpanningTaskOverlay creates a spanning task overlay if this day starts a spanning task
func (tr *TaskRenderer) RenderSpanningTaskOverlay(day Day) *TaskOverlayInfo {
	dayDate := day.getDayDate()
	startingTasks, maxCols := day.findStartingTasks(dayDate)

	if len(startingTasks) == 0 {
		return nil
	}

	// Create separate pills for each spanning task
	var pillContents []string

	for _, spanningTask := range startingTasks {
		// Task name (will be bolded by the macro)
		taskName := tr.escapeLatexSpecialChars(spanningTask.Name)
		if tr.isMilestoneSpanningTask(spanningTask) {
			taskName = "★ " + taskName
		}

		// Objective (will be smaller by the macro)
		objective := ""
		if spanningTask.Description != "" {
			objective = tr.escapeLatexSpecialChars(spanningTask.Description)
		}

		// Get the color for this specific task
		taskColor := tr.hexToRGB(spanningTask.Color)
		if taskColor == "" {
			taskColor = "224,50,212" // Default fallback
		}

		// Create a separate pill for this task
		// All tasks use no offset to make them touch
		pillContent := fmt.Sprintf(`\TaskOverlayBoxNoOffset{%s}{%s}{%s}`,
			taskColor, // Use the task's specific color
			taskName,  // Task name (will be bolded by macro)
			objective) // Objective (will be smaller by macro)
		pillContents = append(pillContents, pillContent)
	}

	// Stack the pills vertically without extra spacing
	content := strings.Join(pillContents, "")

	return &TaskOverlayInfo{
		content: content,
		cols:    maxCols,
	}
}

// RenderTasksForDay returns a formatted string of tasks for this day
func (tr *TaskRenderer) RenderTasksForDay(day Day) string {
	var taskStrings []string

	// Add regular tasks (non-spanning tasks)
	for _, task := range day.Tasks {
		taskStr := tr.escapeLatexSpecialChars(task.Name)

		// Add star for milestone tasks
		if tr.isMilestoneTask(task) {
			taskStr = "★ " + taskStr
		}

		// Apply color styling based on category
		if task.Category != "" {
			color := tr.getColorForCategory(task.Category)
			if color != "" {
				rgbColor := tr.hexToRGB(color)
				taskStr = fmt.Sprintf(`\textcolor[RGB]{%s}{%s}`, rgbColor, taskStr)
			}
		}

		taskStrings = append(taskStrings, taskStr)
	}

	if len(taskStrings) == 0 {
		return ""
	}

	return strings.Join(taskStrings, "\\\\")
}

// Helper methods

func (tr *TaskRenderer) escapeLatexSpecialChars(text string) string {
	// Escape LaTeX special characters
	text = strings.ReplaceAll(text, "&", "\\&")
	text = strings.ReplaceAll(text, "%", "\\%")
	text = strings.ReplaceAll(text, "$", "\\$")
	text = strings.ReplaceAll(text, "#", "\\#")
	text = strings.ReplaceAll(text, "^", "\\textasciicircum{}")
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, "~", "\\textasciitilde{}")
	return text
}

func (tr *TaskRenderer) isMilestoneTask(task Task) bool {
	// Check if task name contains milestone indicators
	name := strings.ToLower(task.Name)
	return strings.Contains(name, "milestone") ||
		strings.Contains(name, "deliverable") ||
		strings.Contains(name, "deadline") ||
		strings.Contains(name, "★")
}

func (tr *TaskRenderer) isMilestoneSpanningTask(task *SpanningTask) bool {
	// Check if spanning task name contains milestone indicators
	name := strings.ToLower(task.Name)
	return strings.Contains(name, "milestone") ||
		strings.Contains(name, "deliverable") ||
		strings.Contains(name, "deadline") ||
		strings.Contains(name, "★")
}

func (tr *TaskRenderer) getColorForCategory(category string) string {
	// Color mapping for categories
	colorMap := map[string]string{
		"Research":     "#FF6B6B", // Red
		"Writing":      "#4ECDC4", // Teal
		"Analysis":     "#45B7D1", // Blue
		"Review":       "#96CEB4", // Green
		"Presentation": "#FFEAA7", // Yellow
		"Planning":     "#DDA0DD", // Plum
		"Data":         "#98D8C8", // Mint
		"Meeting":      "#F7DC6F", // Light Yellow
		"Admin":        "#BB8FCE", // Light Purple
		"Other":        "#85C1E9", // Light Blue
	}

	if color, exists := colorMap[category]; exists {
		return color
	}

	// Default color for unknown categories
	return "#E0E0E0" // Light Gray
}

func (tr *TaskRenderer) hexToRGB(hex string) string {
	if hex == "" {
		return ""
	}

	// Remove # if present
	if strings.HasPrefix(hex, "#") {
		hex = hex[1:]
	}

	// Convert hex to RGB
	if len(hex) == 6 {
		r := hex[0:2]
		g := hex[2:4]
		b := hex[4:6]
		return fmt.Sprintf("%s,%s,%s", r, g, b)
	}

	return ""
}
