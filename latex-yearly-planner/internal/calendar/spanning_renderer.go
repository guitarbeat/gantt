package calendar

import (
	"strconv"
	"strings"
	"time"
)

// SpanningTaskRenderer handles the rendering of spanning tasks with proper overlap management
type SpanningTaskRenderer struct {
	day       *Day
	dayDate   time.Time
	leftCell  string
	dayNumber string
}

// NewSpanningTaskRenderer creates a new renderer for spanning tasks
func NewSpanningTaskRenderer(day *Day, dayNumber string) *SpanningTaskRenderer {
	leftCell := `\begin{tabular}{@{}p{5mm}@{}|}\hfil{}` + dayNumber + `\\ \hline\end{tabular}`
	dayDate := time.Date(day.Time.Year(), day.Time.Month(), day.Time.Day(), 0, 0, 0, 0, time.UTC)

	return &SpanningTaskRenderer{
		day:       day,
		dayDate:   dayDate,
		leftCell:  leftCell,
		dayNumber: dayNumber,
	}
}

// RenderSpanningTasks renders all spanning tasks for this day with proper overlap management
func (r *SpanningTaskRenderer) RenderSpanningTasks() string {
	if len(r.day.SpanningTasks) == 0 {
		return `\hyperlink{` + r.day.ref() + `}{` + r.leftCell + `}`
	}

	// Group tasks by their start day to handle overlapping properly
	tasksStartingToday := r.getTasksStartingOnDay()
	tasksSpanningToday := r.getTasksSpanningOnDay()

	// If tasks start today, render them with overlays
	if len(tasksStartingToday) > 0 {
		return r.renderStartingTasks(tasksStartingToday)
	}

	// If only spanning tasks (started earlier), render simplified version
	if len(tasksSpanningToday) > 0 {
		return r.renderSpanningIndicator(tasksSpanningToday)
	}

	// No tasks on this day
	return `\hyperlink{` + r.day.ref() + `}{` + r.leftCell + `}`
}

// getTasksStartingOnDay returns tasks that start on this specific day
func (r *SpanningTaskRenderer) getTasksStartingOnDay() []*SpanningTask {
	var startingTasks []*SpanningTask
	for _, task := range r.day.SpanningTasks {
		start := time.Date(task.StartDate.Year(), task.StartDate.Month(), task.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		if r.dayDate.Equal(start) {
			startingTasks = append(startingTasks, task)
		}
	}
	return startingTasks
}

// getTasksSpanningOnDay returns tasks that span this day (but may have started earlier)
func (r *SpanningTaskRenderer) getTasksSpanningOnDay() []*SpanningTask {
	var spanningTasks []*SpanningTask
	for _, task := range r.day.SpanningTasks {
		start := time.Date(task.StartDate.Year(), task.StartDate.Month(), task.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(task.EndDate.Year(), task.EndDate.Month(), task.EndDate.Day(), 0, 0, 0, 0, time.UTC)
		if !r.dayDate.Before(start) && !r.dayDate.After(end) {
			spanningTasks = append(spanningTasks, task)
		}
	}
	return spanningTasks
}

// renderStartingTasks renders tasks that start on this day with proper overlays
func (r *SpanningTaskRenderer) renderStartingTasks(tasks []*SpanningTask) string {
	var overlayContents []string
	maxCols := 1

	for i, task := range tasks {
		// Calculate span width
		start := time.Date(task.StartDate.Year(), task.StartDate.Month(), task.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(task.EndDate.Year(), task.EndDate.Month(), task.EndDate.Day(), 0, 0, 0, 0, time.UTC)

		spanDays := int(end.Sub(start).Hours()/24) + 1
		dayOfWeek := (int(r.dayDate.Weekday()) + 6) % 7 // Monday=0
		remainingInWeek := 7 - dayOfWeek

		cols := spanDays
		if cols > remainingInWeek {
			cols = remainingInWeek
		}
		if cols > maxCols {
			maxCols = cols
		}

		// Create task overlay with text inside
		overlay := r.createTaskOverlay(task, i, len(tasks))
		overlayContents = append(overlayContents, overlay)
	}

	// Combine overlays with proper spacing
	spacer := `\vspace{0.1ex}`
	if len(overlayContents) > 1 {
		spacer = `\vspace{0.05ex}` // Tighter spacing for multiple tasks
	}
	combinedOverlays := strings.Join(overlayContents, spacer)

	// Calculate width for spanning
	width := `\dimexpr ` + strconv.Itoa(maxCols) + `\linewidth\relax`

	return `\hyperlink{` + r.day.ref() + `}{` +
		`{\begingroup` +
		`\makebox[0pt][l]{` + r.leftCell + `}` +
		`\makebox[0pt][l]{` +
		`\begin{tikzpicture}[overlay]` +
		`\node[anchor=north west, inner sep=0pt] at (0,0) {` +
		`\begin{minipage}[t]{` + width + `}` +
		combinedOverlays +
		`\end{minipage}};` +
		`\end{tikzpicture}}` +
		`\endgroup}}`
}

// createTaskOverlay creates a single task overlay with embedded text
func (r *SpanningTaskRenderer) createTaskOverlay(task *SpanningTask, index, total int) string {
	// Truncate task name based on available space
	nameText := strings.TrimSpace(task.Name)
	maxChars := 15 // Base character limit

	if total > 1 {
		maxChars = 12 - (total * 2) // Reduce for multiple tasks
		if maxChars < 6 {
			maxChars = 6
		}
	}

	if len(nameText) > maxChars {
		nameText = nameText[:maxChars-3] + "..."
	}

	// Adjust box size and font based on position
	fontSize := `\tiny`
	boxHeight := "1.0ex"
	verticalSpacing := "0.1ex"

	if total > 1 {
		fontSize = `\scriptsize` // Even smaller for multiple
		boxHeight = "0.8ex"
		if index > 0 {
			verticalSpacing = "0.05ex"
			boxHeight = "0.7ex"
		}
	}

	return `\vspace*{` + verticalSpacing + `}` +
		`{\begingroup\setlength{\fboxsep}{0pt}` +
		`\begin{tcolorbox}[enhanced, boxrule=0pt, arc=0pt,` +
		` left=0.3mm, right=0.3mm, top=0.1mm, bottom=0.1mm,` +
		` height=` + boxHeight + `,` +
		` colback=` + task.Color + `!40,` +
		` borderline west={1.0pt}{0pt}{` + task.Color + `!60!black}]` +
		`{\centering\color{black}\textbf{` + fontSize + ` ` + nameText + `}}` +
		`\end{tcolorbox}\endgroup}`
}

// renderSpanningIndicator renders a simple indicator for days where tasks span but don't start
func (r *SpanningTaskRenderer) renderSpanningIndicator(tasks []*SpanningTask) string {
	// For continuation days, just show the day number with subtle task indication
	return `\hyperlink{` + r.day.ref() + `}{` + r.leftCell + `}`
}
