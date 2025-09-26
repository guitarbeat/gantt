// Package scheduler handles calendar layout, task positioning, and LaTeX rendering
// for the PhD dissertation planner system.
//
// Key responsibilities:
// - Calendar grid generation with proper day/week/month structure
// - Task bar positioning and stacking for multi-day spanning tasks
// - Color management for task categories with LaTeX-safe escaping
// - PDF-optimized LaTeX template rendering with proper spacing
//
// This file specifically handles:
// - Day cell construction with responsive sizing
// - Task overlay rendering with smart truncation disabled per user request
// - Special character escaping for LaTeX compatibility (especially & in category names)
// - Dynamic color legend generation based on actual task categories present
package calendar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"phd-dissertation-planner/src/core"
	"phd-dissertation-planner/src/shared/templates"
)

// * Day types and methods (from day.go)

type Days []*Day
type Day struct {
	Time          time.Time
	Tasks         []Task
	SpanningTasks []*SpanningTask
	Cfg           *core.Config // * Reference to core configuration
}

// Task represents a task for a specific day
type Task struct {
	ID          string
	Name        string
	Description string
	Category    string
}

// Day renders the day cell for both small and large views
func (d Day) Day(today, large interface{}) string {
	if d.Time.IsZero() {
		return ""
	}

	day := strconv.Itoa(d.Time.Day())

	if larg, _ := large.(bool); larg {
		return d.renderLargeDay(day)
	}

	if td, ok := today.(Day); ok {
		if d.Time.Equal(td.Time) {
			return templates.EmphCell(day)
		}
	}

	return day
}

// renderLargeDay renders the day cell for large (monthly) view with tasks and spanning tasks
func (d Day) renderLargeDay(day string) string {
	leftCell := d.buildDayNumberCell(day)

	// Check for spanning tasks that start on this day
	overlay := d.renderSpanningTaskOverlay()
	if overlay != nil {
		// For spanning tasks, render them as regular content that stacks vertically
		// instead of using TikZ overlays that can stack in z-dimension
		return d.buildTaskCell(leftCell, overlay.content, false, overlay.cols)
	}

	// Check for regular tasks
	if tasks := d.TasksForDay(); tasks != "" {
		return d.buildTaskCell(leftCell, tasks, false, 0)
	}

	// No tasks: just the day number
	return d.buildSimpleDayCell(leftCell)
}

// renderLargeDayRefactored demonstrates how to use the refactored modules
// This is an example of how the code could be refactored to use the new modules
func (d Day) renderLargeDayRefactored(day string) string {
	// Create the refactored components
	taskRenderer := NewTaskRenderer(d.Cfg)
	cellBuilder := NewCellBuilder(d.Cfg)

	leftCell := cellBuilder.BuildDayNumberCell(day)

	// Check for spanning tasks that start on this day
	overlay := taskRenderer.RenderSpanningTaskOverlay(d)
	if overlay != nil {
		// For spanning tasks, render them as regular content that stacks vertically
		return cellBuilder.BuildTaskCell(leftCell, overlay.content, false, overlay.cols)
	}

	// Check for regular tasks
	if tasks := taskRenderer.RenderTasksForDay(d); tasks != "" {
		return cellBuilder.BuildTaskCell(leftCell, tasks, false, 0)
	}

	// No tasks: just the day number
	return cellBuilder.BuildSimpleDayCell(leftCell)
}

// ref generates a reference string for the day
func (d Day) ref(prefix ...string) string {
	p := ""

	if len(prefix) > 0 {
		p = prefix[0]
	}

	return p + d.Time.Format(time.RFC3339)
}

// * LaTeX cell construction functions

// buildDayNumberCell creates the basic day number cell with minimal padding
// Uses minipage instead of tabular to eliminate auto padding
func (d Day) buildDayNumberCell(day string) string {
	// * Use config-driven day number width
	dayWidth := "6mm" // Default fallback
	if d.Cfg.Layout.LayoutEngine.CalendarLayout.DayNumberWidth != "" {
		dayWidth = d.Cfg.Layout.LayoutEngine.CalendarLayout.DayNumberWidth
	}
	return `\begin{minipage}[t]{` + dayWidth + `}\centering{}` + day + `\end{minipage}`
}

// buildTaskCell creates a cell with either spanning tasks or regular tasks
func (d Day) buildTaskCell(leftCell, content string, isSpanning bool, cols int) string {
	var width, spacing, contentWrapper string

	// * Get config values with fallbacks
	dayNumberWidth := "6mm"
	dayContentMargin := "8mm"
	hyphenPenalty := 50
	tolerance := 1000
	emergencyStretch := "3em"

	if d.Cfg.Layout.LayoutEngine.CalendarLayout.DayNumberWidth != "" {
		dayNumberWidth = d.Cfg.Layout.LayoutEngine.CalendarLayout.DayNumberWidth
	}
	if d.Cfg.Layout.LayoutEngine.CalendarLayout.DayContentMargin != "" {
		dayContentMargin = d.Cfg.Layout.LayoutEngine.CalendarLayout.DayContentMargin
	}
	if d.Cfg.Layout.LaTeX.Typography.HyphenPenalty > 0 {
		hyphenPenalty = d.Cfg.Layout.LaTeX.Typography.HyphenPenalty
	}
	if d.Cfg.Layout.LaTeX.Typography.Tolerance > 0 {
		tolerance = d.Cfg.Layout.LaTeX.Typography.Tolerance
	}
	if d.Cfg.Layout.LaTeX.Typography.SloppyEmergencyStretch != "" {
		emergencyStretch = d.Cfg.Layout.LaTeX.Typography.SloppyEmergencyStretch
	}

	if isSpanning {
		// Spanning task: use tikzpicture overlay with calculated width (z-dimension stacking)
		width = `\dimexpr ` + strconv.Itoa(cols) + `\linewidth\relax`
		spacing = `\makebox[0pt][l]{` + `\begin{tikzpicture}[overlay]` +
			`\node[anchor=north west, inner sep=0pt] at (0,0) {` + `\begin{minipage}[t]{` + width + `}` + content + `\end{minipage}` + `};` +
			`\end{tikzpicture}` + `}`
		contentWrapper = "" // Don't add content twice for spanning tasks
	} else if cols > 0 {
		// Spanning task but rendered as regular content (vertical stacking)
		width = `\dimexpr ` + strconv.Itoa(cols) + `\linewidth\relax`
		spacing = ""             // No offset - start at the beginning of the cell
		contentWrapper = content // Use the content directly without additional wrapping
	} else {
		// Regular task: use full available width and better text flow
		width = `\dimexpr\linewidth - ` + dayContentMargin + `\relax` // Leave space for day number + margins
		spacing = `\hspace*{` + dayNumberWidth + `}`                  // Spacing to align with day number cell width
		contentWrapper = fmt.Sprintf(`{\sloppy\hyphenpenalty=%d\tolerance=%d\emergencystretch=%s\footnotesize\raggedright `,
			hyphenPenalty, tolerance, emergencyStretch) + content + `}`
	}

	inner := `{\begingroup` +
		`\makebox[0pt][l]{` + leftCell + `}` +
		spacing +
		`\begin{minipage}[t]{` + width + `}` +
		contentWrapper +
		`\end{minipage}` +
		`\endgroup}`

	// Wrap entire cell in hyperlink to the day's reference (restores link without visual borders via hypersetup)
	return `\hyperlink{` + d.ref() + `}{` + inner + `}`
}

// buildSimpleDayCell creates a simple day cell without tasks
func (d Day) buildSimpleDayCell(leftCell string) string {
	return leftCell
}

// * Task processing and utility functions

// TaskOverlay represents a spanning task overlay
type TaskOverlay struct {
	content string
	cols    int
}

// renderSpanningTaskOverlay creates a spanning task overlay if this day starts a spanning task
func (d Day) renderSpanningTaskOverlay() *TaskOverlay {
	dayDate := d.getDayDate()
	startingTasks, maxCols := d.findStartingTasks(dayDate)

	if len(startingTasks) == 0 {
		return nil
	}

	// Build the task content for the overlay
	var taskStrings []string
	for _, spanningTask := range startingTasks {
		taskStr := d.escapeLatexSpecialChars(spanningTask.Name)

		// Add star for milestone spanning tasks
		if d.isMilestoneSpanningTask(spanningTask) {
			taskStr = "★ " + taskStr
		}

		// Apply color styling based on category
		if spanningTask.Category != "" && spanningTask.Color != "" {
			rgbColor := hexToRGB(spanningTask.Color)
			taskStr = fmt.Sprintf(`\textcolor[RGB]{%s}{%s}`, rgbColor, taskStr)
		}

		taskStrings = append(taskStrings, taskStr)
	}

	if len(taskStrings) == 0 {
		return nil
	}

	// Each task will get its own color, so we don't need a shared pillColor

	// Create separate pills for each spanning task
	var pillContents []string

	for i, spanningTask := range startingTasks {
		// Task name (will be bolded by the macro)
		taskName := d.escapeLatexSpecialChars(spanningTask.Name)
		if d.isMilestoneSpanningTask(spanningTask) {
			taskName = "★ " + taskName
		}

		// Objective (will be smaller by the macro)
		objective := ""
		if spanningTask.Description != "" {
			objective = d.escapeLatexSpecialChars(spanningTask.Description)
		}

		// Get the color for this specific task
		taskColor := hexToRGB(spanningTask.Color)
		if taskColor == "" {
			taskColor = "224,50,212" // Default fallback
		}

		// Create a separate pill for this task
		// Only the first task gets vertical offset, others touch
		if i == 0 {
			pillContent := fmt.Sprintf(`\TaskOverlayBox{%s}{%s}{%s}`,
				taskColor, // Use the task's specific color
				taskName,  // Task name (will be bolded by macro)
				objective) // Objective (will be smaller by macro)
			pillContents = append(pillContents, pillContent)
		} else {
			// For subsequent tasks, use a custom macro without vertical offset
			pillContent := fmt.Sprintf(`\TaskOverlayBoxNoOffset{%s}{%s}{%s}`,
				taskColor, // Use the task's specific color
				taskName,  // Task name (will be bolded by macro)
				objective) // Objective (will be smaller by macro)
			pillContents = append(pillContents, pillContent)
		}
	}

	// Stack the pills vertically without extra spacing
	content := strings.Join(pillContents, "")

	return &TaskOverlay{
		content: content,
		cols:    maxCols,
	}
}

// TasksForDay returns a formatted string of tasks for this day
func (d Day) TasksForDay() string {
	var taskStrings []string

	// Add regular tasks (non-spanning tasks)
	for _, task := range d.Tasks {
		taskStr := d.escapeLatexSpecialChars(task.Name)

		// Add star for milestone tasks
		if d.isMilestoneTask(task) {
			taskStr = "★ " + taskStr
		}

		// Apply color styling based on category
		if task.Category != "" {
			color := getColorForCategory(task.Category)
			if color != "" {
				rgbColor := hexToRGB(color)
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

// findStartingTasks finds tasks that should be displayed on this day and calculates max columns
func (d Day) findStartingTasks(dayDate time.Time) ([]*SpanningTask, int) {
	var tasksToShow []*SpanningTask
	var maxCols int

	// Find all tasks that are active on this day
	var activeTasks []*SpanningTask
	for _, task := range d.SpanningTasks {
		start := d.getTaskStartDate(task)
		end := d.getTaskEndDate(task)

		// Check if task is active on this day
		if (dayDate.Equal(start) || dayDate.After(start)) && (dayDate.Equal(end) || dayDate.Before(end)) {
			activeTasks = append(activeTasks, task)
		}
	}

	// If there are multiple active tasks, show them all stacked
	// If there's only one active task, only show it if it starts on this day
	if len(activeTasks) > 1 {
		// Multiple tasks active - show all of them stacked
		tasksToShow = activeTasks
		for _, task := range activeTasks {
			end := d.getTaskEndDate(task)
			cols := d.calculateTaskSpanColumns(dayDate, end)
			if cols > maxCols {
				maxCols = cols
			}
		}
	} else if len(activeTasks) == 1 {
		// Single task - only show if it starts on this day
		task := activeTasks[0]
		start := d.getTaskStartDate(task)
		if dayDate.Equal(start) {
			tasksToShow = activeTasks
			end := d.getTaskEndDate(task)
			maxCols = d.calculateTaskSpanColumns(dayDate, end)
		}
	}

	return tasksToShow, maxCols
}

// sortTasksByDuration sorts tasks by duration for better visual organization
func (d Day) sortTasksByDuration(tasks []*SpanningTask) []*SpanningTask {
	sorted := make([]*SpanningTask, len(tasks))
	copy(sorted, tasks)

	// Sort by duration (shorter tasks first for better stacking)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			duration1 := sorted[j].EndDate.Sub(sorted[j].StartDate)
			duration2 := sorted[j+1].EndDate.Sub(sorted[j+1].StartDate)
			if duration1 > duration2 {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// isMilestoneTask checks if a task is a milestone based on its description
func (d Day) isMilestoneTask(task Task) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(task.Description)), "MILESTONE:")
}

// isMilestoneSpanningTask checks if a spanning task is a milestone based on its description
func (d Day) isMilestoneSpanningTask(task *SpanningTask) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(task.Description)), "MILESTONE:")
}

// escapeLatexSpecialChars replaces special LaTeX characters with their escaped versions
func escapeLatexSpecialChars(text string) string {
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
// NOTE: Currently disabled - returning full text to avoid aggressive truncation
func (d Day) smartTruncateText(text string) string {
	// For now, return full text to avoid unwanted truncation
	// TODO: Implement better space utilization strategies
	return text
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

// * Week types and methods (from week.go)

type Weeks []*Week
type Week struct {
	Days [7]Day

	Weekday  time.Weekday
	Year     *Year
	Months   Months
	Quarters Quarters
}

func NewWeeksForMonth(wd time.Weekday, year *Year, qrtr *Quarter, month *Month, cfg *core.Config) Weeks {
	ptr := time.Date(year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
	weekday := ptr.Weekday()
	shift := (7 + weekday - wd) % 7

	week := &Week{Weekday: wd, Year: year, Months: Months{month}, Quarters: Quarters{qrtr}}

	for i := shift; i < 7; i++ {
		week.Days[i] = Day{Time: ptr, Tasks: nil, SpanningTasks: nil, Cfg: cfg}
		ptr = ptr.AddDate(0, 0, 1)
	}

	weeks := Weeks{}
	weeks = append(weeks, week)

	for ptr.Month() == month.Month {
		week = &Week{Weekday: weekday, Year: year, Months: Months{month}, Quarters: Quarters{qrtr}}

		for i := 0; i < 7; i++ {
			if ptr.Month() != month.Month {
				break
			}

			week.Days[i] = Day{Time: ptr, Tasks: nil, SpanningTasks: nil, Cfg: cfg}
			ptr = ptr.AddDate(0, 0, 1)
		}

		weeks = append(weeks, week)
	}

	return weeks
}

func (w *Week) HasDays() bool {
	for _, d := range w.Days {
		if !d.Time.IsZero() {
			return true
		}
	}
	return false
}

func (w *Week) WeekNumber(large interface{}) string {
	wn := w.weekNumber()
	larg, _ := large.(bool)

	itoa := strconv.Itoa(wn)
	ref := w.ref()
	if !larg {
		return templates.Link(ref, itoa)
	}

	text := `\rotatebox[origin=tr]{90}{\makebox[\myLenMonthlyCellHeight][c]{Week ` + itoa + `}}`

	return templates.Link(ref, text)
}

func (w *Week) weekNumber() int {
	// Calculate sequential week number for the entire year (1-based)
	// Find the first non-zero day in the week
	var firstDay time.Time
	for _, day := range w.Days {
		if !day.Time.IsZero() {
			firstDay = day.Time
			break
		}
	}

	if firstDay.IsZero() {
		return 0
	}

	// Find the first day of the year
	firstOfYear := time.Date(firstDay.Year(), 1, 1, 0, 0, 0, 0, firstDay.Location())

	// Calculate how many days from the start of the year to the first day of this week
	daysFromStart := int(firstDay.Sub(firstOfYear).Hours() / 24)

	// Calculate the week number (1-based, starting from the first week of the year)
	weekNum := (daysFromStart / 7) + 1

	// Ensure we don't go below 1
	if weekNum < 1 {
		weekNum = 1
	}

	return weekNum
}

func (w Week) ref(prefix ...string) string {
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return p + "week-" + strconv.Itoa(w.Year.Number) + "-" + strconv.Itoa(w.weekNumber())
}

func NewWeeksForYear(wd time.Weekday, year *Year, cfg *core.Config) Weeks {
	var weeks Weeks
	ptr := time.Date(year.Number, 1, 1, 0, 0, 0, 0, time.Local)
	weekday := ptr.Weekday()
	_ = (7 + weekday - wd) % 7

	for i := 0; i < 53; i++ {
		week := &Week{Weekday: wd, Year: year}
		for j := 0; j < 7; j++ {
			week.Days[j] = Day{Time: ptr, Tasks: nil, SpanningTasks: nil, Cfg: cfg}
			ptr = ptr.AddDate(0, 0, 1)
		}
		weeks = append(weeks, week)
	}

	return weeks
}

// * Month types and methods (from month.go)

type Months []*Month

func (m Months) Months() []time.Month {
	if len(m) == 0 {
		return nil
	}

	out := make([]time.Month, 0, len(m))

	for _, month := range m {
		out = append(out, month.Month)
	}

	return out
}

type Month struct {
	Year    *Year
	Quarter *Quarter
	Month   time.Month
	Weekday time.Weekday
	Weeks   Weeks
	Cfg     *core.Config // * Reference to core configuration
}

func NewMonth(wd time.Weekday, year *Year, qrtr *Quarter, month time.Month, cfg *core.Config) *Month {
	m := &Month{
		Year:    year,
		Quarter: qrtr,
		Month:   month,
		Weekday: wd,
		Cfg:     cfg,
	}

	m.Weeks = NewWeeksForMonth(wd, year, qrtr, m, cfg)

	return m
}

func (m Month) Breadcrumb() string {
	return templates.Items{
		templates.NewIntItem(m.Year.Number),
		templates.NewTextItem("Q" + strconv.Itoa(m.Quarter.Number)),
		templates.NewMonthItem(m.Month),
	}.Table(true)
}

func (m Month) MonthLink() string {
	return templates.Link(m.ref(), m.Month.String())
}

func (m Month) ref(prefix ...string) string {
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return p + "month-" + strconv.Itoa(m.Year.Number) + "-" + strconv.Itoa(int(m.Month))
}

// HeadingMOS creates a heading for the month-overview-single view
func (m Month) HeadingMOS(prefix ...string) string {
	leaf := ""
	if len(prefix) > 1 {
		leaf = prefix[1]
	}
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	monthStr := m.Month.String()
	if len(leaf) > 0 {
		monthStr = templates.Link(m.ref(p), monthStr)
	}

	// * Use config-driven header angle size offset
	headerAngleOffset := "0.86pt" // Default fallback
	if m.Cfg.Layout.LayoutEngine.CalendarLayout.HeaderAngleSizeOffset != "" {
		headerAngleOffset = m.Cfg.Layout.LayoutEngine.CalendarLayout.HeaderAngleSizeOffset
	}
	anglesize := `\dimexpr\myLenHeaderResizeBox-` + headerAngleOffset
	var ll, rl string
	var r1, r2 []string
	if m.PrevExists() {
		ll = "l"
		leftNavBox := templates.ResizeBoxW(anglesize, `$\langle$`)
		r1 = append(r1, templates.Multirow(2, templates.Hyperlink(m.Prev().ref(p), leftNavBox)))
		r2 = append(r2, "")
	}
	r1 = append(r1, templates.Multirow(2, templates.ResizeBoxW(`\myLenHeaderResizeBox`, monthStr)))
	r2 = append(r2, "")
	r1 = append(r1, templates.Bold(m.Month.String()))
	r2 = append(r2, strconv.Itoa(m.Year.Number))
	if m.NextExists() {
		rl = "l"
		rightNavBox := templates.ResizeBoxW(anglesize, `$\rangle$`)
		r1 = append(r1, templates.Multirow(2, templates.Hyperlink(m.Next().ref(p), rightNavBox)))
		r2 = append(r2, "")
	}
	contents := strings.Join(r1, ` & `) + `\\` + "\n" + strings.Join(r2, ` & `)
	return templates.Hypertarget(p+m.ref(), "") + templates.Tabular("@{}"+ll+"l|l"+rl, contents)
}

// PrevNext creates navigation items for previous and next months
func (m Month) PrevNext(prefix ...string) templates.Items {
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	items := templates.Items{}

	if m.PrevExists() {
		prev := m.Prev()
		items = append(items, templates.NewTextItem(prev.Month.String()).RefText(p+prev.ref()))
	}

	if m.NextExists() {
		next := m.Next()
		items = append(items, templates.NewTextItem(next.Month.String()).RefText(p+next.ref()))
	}

	return items
}

// Prev returns the previous month
func (m Month) Prev() Month {
	if m.Month == time.January {
		return Month{Year: m.Year, Quarter: m.Quarter, Month: time.December}
	}
	return Month{Year: m.Year, Quarter: m.Quarter, Month: m.Month - 1}
}

// Next returns the next month
func (m Month) Next() Month {
	if m.Month == time.December {
		return Month{Year: m.Year, Quarter: m.Quarter, Month: time.January}
	}
	return Month{Year: m.Year, Quarter: m.Quarter, Month: m.Month + 1}
}

// PrevExists checks if the previous month exists
func (m Month) PrevExists() bool {
	return m.Month > time.January
}

// NextExists checks if the next month exists
func (m Month) NextExists() bool {
	return m.Month < time.December
}

func (m *Month) DefineTable(typ interface{}, large interface{}) string {
	full, _ := large.(bool)

	typStr, ok := typ.(string)
	if !ok || typStr == "tabularx" {
		weekAlign := "Y|"
		days := "Y"
		if full {
			weekAlign = `|l!{\vrule width \myLenLineThicknessThick}`
			days = "@{}X@{}|"
		}

		return `\begin{tabularx}{\linewidth}{` + weekAlign + `*{7}{` + days + `}}`
	}

	return `\begin{tabular}[t]{c|*{7}{c}}`
}

func (m *Month) EndTable(typ interface{}) string {
	typStr, ok := typ.(string)
	if !ok || typStr == "tabularx" {
		return `\end{tabularx}`
	}

	return `\end{tabular}`
}

func (m *Month) MaybeName(large interface{}) string {
	larg, _ := large.(bool)

	if larg { // likely on a monthly page; no need to print it again
		return ""
	}

	return `\multicolumn{8}{c}{` + templates.Link(m.Month.String(), m.Month.String()) + `} \\ \hline`
}

func (m *Month) WeekHeader(large interface{}) string {
	full, _ := large.(bool)

	names := make([]string, 0, 8)

	if full {
		names = append(names, "")
	} else {
		names = append(names, "W")
	}

	for i := time.Sunday; i < 7; i++ {
		name := ((m.Weekday + i) % 7).String()
		if full {
			name = `\hfil{}` + name
		} else {
			name = name[:1]
		}

		names = append(names, name)
	}

	return strings.Join(names, " & ")
}

// stripHashPrefix removes the # prefix from hex colors for LaTeX compatibility (HTML colors work with both cases)
func stripHashPrefix(color string) string {
	if len(color) > 0 && color[0] == '#' {
		return color[1:]
	}
	return color
}

// hexToRGB converts hex color to RGB format for LaTeX
func hexToRGB(hex string) string {
	// Remove # prefix if present
	hex = stripHashPrefix(hex)

	// Convert hex to RGB
	if len(hex) == 6 {
		// Parse hex values
		r, _ := strconv.ParseInt(hex[0:2], 16, 64)
		g, _ := strconv.ParseInt(hex[2:4], 16, 64)
		b, _ := strconv.ParseInt(hex[4:6], 16, 64)
		return fmt.Sprintf("%d,%d,%d", r, g, b)
	}

	// Fallback for invalid hex
	return "128,128,128"
}

func (m *Month) GetTaskColors() map[string]string {
	colorMap := make(map[string]string)

	// Only add colors for task categories that are actually present in this month
	for _, week := range m.Weeks {
		for _, day := range week.Days {
			// Check spanning tasks
			for _, task := range day.SpanningTasks {
				if task.Category != "" {
					color := getColorForCategory(task.Category)
					if color != "" {
						// Convert to RGB for LaTeX compatibility
						// Escape LaTeX special characters in category name
						escapedCategory := escapeLatexSpecialChars(task.Category)
						colorMap[hexToRGB(color)] = escapedCategory
					}
				}
			}
			// Check regular tasks
			for _, task := range day.Tasks {
				if task.Category != "" {
					color := getColorForCategory(task.Category)
					if color != "" {
						// Convert to RGB for LaTeX compatibility
						// Escape LaTeX special characters in category name
						escapedCategory := escapeLatexSpecialChars(task.Category)
						colorMap[hexToRGB(color)] = escapedCategory
					}
				}
			}
		}
	}

	return colorMap
}

// * Year and Quarter types and methods (from time_units.go)

type Years []*Year

type Year struct {
	Number   int
	Quarters Quarters
	Weeks    Weeks
}

func NewYear(wd time.Weekday, year int, cfg *core.Config) *Year {
	out := &Year{Number: year}
	out.Weeks = NewWeeksForYear(wd, out, cfg)
	for q := 1; q <= 4; q++ {
		out.Quarters = append(out.Quarters, NewQuarter(wd, out, q, cfg))
	}
	return out
}

func (y Year) Breadcrumb() string {
	return templates.Items{
		templates.NewIntItem(y.Number),
	}.Table(true)
}

func (y Year) YearLink() string {
	return templates.Link(y.ref(), strconv.Itoa(y.Number))
}

func (y Year) ref(prefix ...string) string {
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return p + "year-" + strconv.Itoa(y.Number)
}

// SideQuarters returns the quarters for the year
func (y Year) SideQuarters(quarterNumber ...int) Quarters {
	return y.Quarters
}

// SideMonths returns all months for the year
func (y Year) SideMonths(month ...time.Month) Months {
	var months Months
	for _, quarter := range y.Quarters {
		months = append(months, quarter.Months...)
	}
	return months
}

type Quarters []*Quarter

type Quarter struct {
	Number int
	Year   *Year
	Months Months
}

func NewQuarter(wd time.Weekday, year *Year, quarter int, cfg *core.Config) *Quarter {
	out := &Quarter{Number: quarter, Year: year}
	for m := 1; m <= 3; m++ {
		month := time.Month((quarter-1)*3 + m)
		out.Months = append(out.Months, NewMonth(wd, year, out, month, cfg))
	}
	return out
}

func (q Quarter) Breadcrumb() string {
	return templates.Items{
		templates.NewIntItem(q.Year.Number),
		templates.NewTextItem("Q" + strconv.Itoa(q.Number)),
	}.Table(true)
}

func (q Quarter) QuarterLink() string {
	return templates.Link(q.ref(), "Q"+strconv.Itoa(q.Number))
}

func (q Quarter) ref(prefix ...string) string {
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return p + "quarter-" + strconv.Itoa(q.Year.Number) + "-" + strconv.Itoa(q.Number)
}

// SpanningTask represents a task that spans multiple days
type SpanningTask struct {
	ID          string
	Name        string
	Description string
	Category    string
	StartDate   time.Time
	EndDate     time.Time
	Color       string
	Progress    int    // Progress percentage (0-100)
	Status      string // Task status
	Assignee    string // Task assignee
}

// CreateSpanningTask creates a new spanning task from basic task data
func CreateSpanningTask(task core.Task, startDate, endDate time.Time) SpanningTask {
	// * Use Sub-Phase as category for better granularity
	color := getColorForCategory(task.Category)

	return SpanningTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Category:    task.Category, // * Fixed: Use Category field
		StartDate:   startDate,
		EndDate:     endDate,
		Color:       color,
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

// getColorForCategory returns a color for the given category using algorithmic generation
func getColorForCategory(category string) string {
	// Generate a consistent, visually distinct color algorithmically
	return generateDynamicColor(category)
}

// generateDynamicColor creates a consistent, visually distinct color based on the category name
func generateDynamicColor(category string) string {
	// Use a better hash function to generate consistent colors
	hash := 0
	for i, char := range category {
		hash = hash*31 + int(char) + i*7 // Add position to improve distribution
	}

	// Convert hash to a positive number
	if hash < 0 {
		hash = -hash
	}

	// Generate HSL color with good saturation and lightness for readability
	hue := float64(hash % 360)                 // 0-360 degrees
	saturation := 0.7 + float64(hash%30)/100.0 // 0.7-1.0 for good saturation
	lightness := 0.5 + float64(hash%20)/100.0  // 0.5-0.7 for good contrast

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
