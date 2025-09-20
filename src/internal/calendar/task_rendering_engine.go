package calendar

import (
	"strconv"
	"strings"
)

// * Task rendering and overlay functions
// This file contains all the logic for rendering task overlays in the calendar.
// It handles both single tasks and multiple overlapping tasks with smart stacking.

// overlayInfo holds information about a spanning task overlay
type overlayInfo struct {
	content string
	cols    int
}

// TaskRenderingConfig holds configuration for task rendering
type TaskRenderingConfig struct {
	// Spacing configuration
	DefaultSpacing    string
	FirstTaskSpacing  string
	
	// Height configuration  
	DefaultHeight     string
	FirstTaskHeight   string
	
	// Text configuration
	MaxChars          int
	MaxCharsCompact   int
	MaxCharsVeryCompact int
	MaxTasksDisplay   int
}

// getDefaultTaskRenderingConfig returns the default configuration for task rendering
// * NO-OVERLAP CONFIGURATION: This configuration is optimized to prevent task overlap
//   by using increased spacing, larger heights, and limiting the number of displayed tasks
func getDefaultTaskRenderingConfig() TaskRenderingConfig {
	return TaskRenderingConfig{
		// Spacing configuration - increased to prevent overlap
		DefaultSpacing:   "0.8ex",
		FirstTaskSpacing: "0.5ex",
		
		// Height configuration - increased to prevent overlap
		DefaultHeight:    "3.0ex",
		FirstTaskHeight:  "3.5ex",
		
		// Text configuration - from constants in day.go
		MaxChars:          maxTaskChars,
		MaxCharsCompact:   maxTaskCharsCompact,
		MaxCharsVeryCompact: maxTaskCharsVeryCompact,
		MaxTasksDisplay:   2, // Reduced from 3 to prevent overlap
	}
}

// renderSpanningTaskOverlay renders spanning task overlays for multiple tasks starting on this day
// Returns nil if no spanning tasks exist or none start on this day
func (d Day) renderSpanningTaskOverlay() *overlayInfo {
	if len(d.SpanningTasks) == 0 {
		return nil
	}

	dayDate := d.getDayDate()
	startingTasks, maxCols := d.findStartingTasks(dayDate)

	if len(startingTasks) == 0 {
		return nil
	}

	// Build content for all starting tasks using smart stacking
	content := d.buildMultiTaskOverlayContent(startingTasks)

	return &overlayInfo{
		content: content,
		cols:    maxCols,
	}
}

// buildTaskOverlayContent creates the LaTeX content for a single task overlay
// Used when only one task starts on a given day
func (d Day) buildTaskOverlayContent(task *SpanningTask) string {
	nameText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Name))
	descText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Description))

	// Add star indicator for milestone tasks
	if d.isMilestoneSpanningTask(task) {
		nameText = "★ " + nameText
	}

	// Use calendar macros for overlay with proper spacing
	return `\vspace*{0.1ex}` + `\TaskOverlayBox{` + task.Color + `}{` + nameText + `}{` + descText + `}`
}

// buildMultiTaskOverlayContent creates compact stacked content for multiple tasks
// Uses smart stacking to prevent overlap and improve readability
func (d Day) buildMultiTaskOverlayContent(tasks []*SpanningTask) string {
	if len(tasks) == 0 {
		return ""
	}

	// Single task - use full overlay format
	if len(tasks) == 1 {
		return d.buildTaskOverlayContent(tasks[0])
	}

	config := getDefaultTaskRenderingConfig()
	
	// Sort tasks by category priority for better visual organization
	sortedTasks := d.sortTasksByPriority(tasks)

	var contentParts []string

	// Show up to maxTasksDisplay tasks in compact format
	for i := 0; i < config.MaxTasksDisplay && i < len(sortedTasks); i++ {
		task := sortedTasks[i]
		compactContent := d.buildCompactTaskOverlay(task, i, len(sortedTasks))
		contentParts = append(contentParts, compactContent)
	}

	// Add indicator if there are more tasks than we can display
	if len(sortedTasks) > config.MaxTasksDisplay {
		moreCount := len(sortedTasks) - config.MaxTasksDisplay
		indicator := d.buildMoreTasksIndicator(moreCount)
		contentParts = append(contentParts, indicator)
	}

	return strings.Join(contentParts, "")
}

// buildMoreTasksIndicator creates the "+X more" indicator for additional tasks
// Shows when there are more tasks than can be displayed in the available space
func (d Day) buildMoreTasksIndicator(moreCount int) string {
	return `\vspace*{0.02ex}{\begingroup\setlength{\fboxsep}{0pt}` +
		`\begin{tcolorbox}[enhanced, boxrule=0pt, arc=0pt,` +
		` left=0.5mm, right=0.5mm, top=0.1mm, bottom=0.1mm,` +
		` colback=gray!15, height=0.5ex,` +
		` borderline west={0.5pt}{0pt}{gray!40}]` +
		`{\centering\color{gray}\textbf{\tiny +` + strconv.Itoa(moreCount) + ` more}}` +
		`\end{tcolorbox}\endgroup}`
}

// buildCompactTaskOverlay creates a compact task overlay for multiple tasks
// Used when multiple tasks start on the same day to create stacked display
func (d Day) buildCompactTaskOverlay(task *SpanningTask, index, total int) string {
	nameText := d.prepareTaskName(task)
	nameText = d.truncateTaskName(nameText, total)

	spacing, boxHeight := d.getTaskSpacingAndHeight(index)
    textBody := d.buildTaskTextBody(nameText)

	return d.buildCompactTaskBox(spacing, boxHeight, task.Color, textBody)
}

// prepareTaskName prepares the task name with milestone indicator
// Escapes LaTeX special characters and adds milestone star if applicable
func (d Day) prepareTaskName(task *SpanningTask) string {
	nameText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Name))
	if d.isMilestoneSpanningTask(task) {
		nameText = "★ " + nameText
	}
	return nameText
}

// truncateTaskName truncates task name based on total number of tasks
// Uses progressive truncation: more tasks = shorter text per task
func (d Day) truncateTaskName(nameText string, total int) string {
	config := getDefaultTaskRenderingConfig()
	
	// Progressive truncation based on number of tasks
	maxChars := config.MaxChars
	if total > 2 {
		maxChars = config.MaxCharsCompact
	}
	if total > 3 {
		maxChars = config.MaxCharsVeryCompact
	}

	// Apply truncation if needed
	if len(nameText) > maxChars {
		nameText = d.smartTruncateText(nameText, maxChars)
	}
	return nameText
}

// getTaskSpacingAndHeight returns spacing and height based on task index
// Uses configuration to ensure consistent spacing and readability
func (d Day) getTaskSpacingAndHeight(index int) (string, string) {
	config := getDefaultTaskRenderingConfig()
	
	// First task gets special treatment for better visual hierarchy
	if index == 0 {
		return config.FirstTaskSpacing, config.FirstTaskHeight
	}
	
	// Subsequent tasks use default spacing and height
	return config.DefaultSpacing, config.DefaultHeight
}

// buildTaskTextBody creates the text body for a task
func (d Day) buildTaskTextBody(nameText string) string {
    // * Use fixed task font size via LaTeX macro \TaskFontSize (defined in macros.tpl)
    return `{\hyphenpenalty=10000\exhyphenpenalty=10000\emergencystretch=2em\setstretch{0.7}` +
        `{\centering\color{black}\TaskFontSize\textbf{` + nameText + `}}}`
}

// buildCompactTaskBox creates the tcolorbox for a compact task
func (d Day) buildCompactTaskBox(spacing, boxHeight, color, textBody string) string {
	// Use macro wrapper for compact bar
	return `\TaskCompactBox{` + spacing + `}{` + boxHeight + `}{` + color + `}{` + textBody + `}`
}
