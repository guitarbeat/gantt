package calendar

import (
	"strconv"
	"strings"
)

// * Task rendering and overlay functions

// overlayInfo holds information about a spanning task overlay
type overlayInfo struct {
	content string
	cols    int
}

// renderSpanningTaskOverlay renders spanning task overlays for multiple tasks starting on this day
func (d Day) renderSpanningTaskOverlay() *overlayInfo {
	if len(d.SpanningTasks) == 0 {
		return nil
	}

	dayDate := d.getDayDate()
	startingTasks, maxCols := d.findStartingTasks(dayDate)

	if len(startingTasks) == 0 {
		return nil
	}

	// Build content for all starting tasks
	content := d.buildMultiTaskOverlayContent(startingTasks)

	return &overlayInfo{
		content: content,
		cols:    maxCols,
	}
}

// buildTaskOverlayContent creates the LaTeX content for a task overlay
func (d Day) buildTaskOverlayContent(task *SpanningTask) string {
	nameText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Name))
	descText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Description))

	// Add star for milestone tasks
	if d.isMilestoneSpanningTask(task) {
		nameText = "★ " + nameText
	}

	// Use calendar macros for overlay
	return `\vspace*{0.1ex}` + `\TaskOverlayBox{` + task.Color + `}{` + nameText + `}{` + descText + `}`
}

// buildMultiTaskOverlayContent creates compact stacked content for multiple tasks
func (d Day) buildMultiTaskOverlayContent(tasks []*SpanningTask) string {
	if len(tasks) == 0 {
		return ""
	}

	if len(tasks) == 1 {
		return d.buildTaskOverlayContent(tasks[0])
	}

	// Sort tasks by category priority for better visual organization
	sortedTasks := d.sortTasksByPriority(tasks)

	var contentParts []string

	// Show up to maxTasksDisplay tasks in compact format
	for i := 0; i < maxTasksDisplay && i < len(sortedTasks); i++ {
		task := sortedTasks[i]
		compactContent := d.buildCompactTaskOverlay(task, i, len(sortedTasks))
		contentParts = append(contentParts, compactContent)
	}

	// Add indicator if there are more tasks
	if len(sortedTasks) > maxTasksDisplay {
		moreCount := len(sortedTasks) - maxTasksDisplay
		indicator := d.buildMoreTasksIndicator(moreCount)
		contentParts = append(contentParts, indicator)
	}

	return strings.Join(contentParts, "")
}

// buildMoreTasksIndicator creates the "+X more" indicator for additional tasks
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
func (d Day) buildCompactTaskOverlay(task *SpanningTask, index, total int) string {
	nameText := d.prepareTaskName(task)
	nameText = d.truncateTaskName(nameText, total)

	spacing, boxHeight := d.getTaskSpacingAndHeight(index)
	textBody := d.buildTaskTextBody(nameText)

	return d.buildCompactTaskBox(spacing, boxHeight, task.Color, textBody)
}

// prepareTaskName prepares the task name with milestone indicator
func (d Day) prepareTaskName(task *SpanningTask) string {
	nameText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Name))
	if d.isMilestoneSpanningTask(task) {
		nameText = "★ " + nameText
	}
	return nameText
}

// truncateTaskName truncates task name based on total number of tasks
func (d Day) truncateTaskName(nameText string, total int) string {
	maxChars := maxTaskChars
	if total > 2 {
		maxChars = maxTaskCharsCompact
	}
	if total > 3 {
		maxChars = maxTaskCharsVeryCompact
	}

	if len(nameText) > maxChars {
		nameText = d.smartTruncateText(nameText, maxChars)
	}
	return nameText
}

// getTaskSpacingAndHeight returns spacing and height based on task index
func (d Day) getTaskSpacingAndHeight(index int) (string, string) {
	spacing := "0.05ex"
	boxHeight := "0.9ex"
	if index == 0 {
		spacing = "0.1ex"
		boxHeight = "1.0ex"
	}
	return spacing, boxHeight
}

// buildTaskTextBody creates the text body for a task
func (d Day) buildTaskTextBody(nameText string) string {
	return `{\hyphenpenalty=10000\exhyphenpenalty=10000\emergencystretch=2em\setstretch{0.7}` +
		`{\centering\color{black}\textbf{\tiny ` + nameText + `}}}`
}

// buildCompactTaskBox creates the tcolorbox for a compact task
func (d Day) buildCompactTaskBox(spacing, boxHeight, color, textBody string) string {
	// Use macro wrapper for compact bar
	return `\TaskCompactBox{` + spacing + `}{` + boxHeight + `}{` + color + `}{` + textBody + `}`
}
