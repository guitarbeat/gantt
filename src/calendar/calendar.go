// Package calendar handles calendar layout, task positioning, and LaTeX rendering
// for the PhD dissertation planner system.
//
// Key responsibilities:
// - Calendar grid generation with proper day/week/month structure
// - Task bar positioning and stacking for multi-day spanning tasks
// - Color management for task categories with LaTeX-safe escaping
// - PDF-optimized LaTeX template rendering with proper spacing
package calendar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"phd-dissertation-planner/src/core"
	"phd-dissertation-planner/src/shared/templates"
)

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// Days is a collection of Day pointers
type Days []*Day

// Day represents a single calendar day with its tasks
type Day struct {
	Time  time.Time
	Tasks []*SpanningTask // All tasks (even 1-day tasks are "spanning")
	Cfg   *core.Config
}

// TaskOverlay represents a spanning task overlay with LaTeX content
type TaskOverlay struct {
	content string // LaTeX content
	cols    int    // Number of columns to span
}

// ============================================================================
// DAY RENDERING
// ============================================================================

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

// renderLargeDay renders the day cell for large (monthly) view with tasks
func (d Day) renderLargeDay(day string) string {
	leftCell := d.buildDayNumberCell(day)

	// Check for tasks using intelligent stacking
	overlay := d.renderSpanningTaskOverlay()
	if overlay != nil {
		// Use spanning mode if any task spans more than 1 column
		isSpanning := overlay.cols > 1
		return d.buildTaskCell(leftCell, overlay.content, isSpanning, overlay.cols)
	}

	// No tasks: just the day number
	return d.buildSimpleDayCell(leftCell)
}

// ref generates a reference string for the day
func (d Day) ref(prefix ...string) string {
	p := ""

	if len(prefix) > 0 {
		p = prefix[0]
	}

	return p + d.Time.Format(time.RFC3339)
}

// ============================================================================
// LATEX CELL CONSTRUCTION
// ============================================================================

// cellConfig holds configuration values for cell rendering
type cellConfig struct {
	dayNumberWidth   string
	dayContentMargin string
	hyphenPenalty    int
	tolerance        int
	emergencyStretch string
}

// getCellConfig extracts cell configuration from Day config with fallbacks
func (d Day) getCellConfig() cellConfig {
	return cellConfig{
		dayNumberWidth:   d.Cfg.GetDayNumberWidth(),
		dayContentMargin: d.Cfg.GetDayContentMargin(),
		hyphenPenalty:    d.Cfg.GetHyphenPenalty(),
		tolerance:        d.Cfg.GetTolerance(),
		emergencyStretch: d.Cfg.GetEmergencyStretch(),
	}
}

// cellLayout defines the LaTeX layout parameters for a cell
type cellLayout struct {
	width          string
	spacing        string
	contentWrapper string
}

// buildDayNumberCell creates the basic day number cell with minimal padding and hypertarget
// Uses minipage instead of tabular to eliminate auto padding
func (d Day) buildDayNumberCell(day string) string {
	cfg := d.getCellConfig()
	// Create hypertarget for this day to enable hyperlink navigation
	hypertarget := fmt.Sprintf(`\hypertarget{%s}{}`, d.ref())
	return hypertarget + `\begin{minipage}[t]{` + cfg.dayNumberWidth + `}\centering{}` + day + `\end{minipage}`
}

// buildTaskCell creates a cell with either spanning tasks or regular tasks
func (d Day) buildTaskCell(leftCell, content string, isSpanning bool, cols int) string {
	cfg := d.getCellConfig()
	layout := d.determineCellLayout(content, isSpanning, cols, cfg)

	inner := d.buildCellInner(leftCell, layout)
	return d.wrapWithHyperlink(inner)
}

// determineCellLayout determines the appropriate layout based on task type
func (d Day) determineCellLayout(content string, isSpanning bool, cols int, cfg cellConfig) cellLayout {
	if isSpanning {
		return d.buildSpanningLayout(content, cols)
	} else if cols > 0 {
		return d.buildVerticalStackLayout(content)
	}
	return d.buildRegularLayout(content, cfg)
}

// buildSpanningLayout creates layout for spanning tasks using tikzpicture overlay
func (d Day) buildSpanningLayout(content string, cols int) cellLayout {
	width := `\dimexpr ` + strconv.Itoa(cols) + `\linewidth\relax`
	spacing := `\makebox[0pt][l]{` + `\begin{tikzpicture}[overlay]` +
		`\node[anchor=north west, inner sep=0pt] at (0,0) {` + `\begin{minipage}[t]{` + width + `}` + content + `\end{minipage}` + `};` +
		`\end{tikzpicture}` + `}`

	return cellLayout{
		width:          width,
		spacing:        spacing,
		contentWrapper: "", // Don't add content twice for spanning tasks
	}
}

// buildVerticalStackLayout creates layout for vertically stacked tasks
func (d Day) buildVerticalStackLayout(content string) cellLayout {
	return cellLayout{
		width:          `\linewidth`, // Just use the cell width
		spacing:        "",           // No offset
		contentWrapper: content,      // Use content directly
	}
}

// buildRegularLayout creates layout for regular tasks with text flow
func (d Day) buildRegularLayout(content string, cfg cellConfig) cellLayout {
	width := `\dimexpr\linewidth - ` + cfg.dayContentMargin + `\relax`
	spacing := `\hspace*{` + cfg.dayNumberWidth + `}`
	contentWrapper := fmt.Sprintf(`{\sloppy\hyphenpenalty=%d\tolerance=%d\emergencystretch=%s\footnotesize\raggedright `,
		cfg.hyphenPenalty, cfg.tolerance, cfg.emergencyStretch) + content + `}`

	return cellLayout{
		width:          width,
		spacing:        spacing,
		contentWrapper: contentWrapper,
	}
}

// buildCellInner constructs the inner content of a cell
func (d Day) buildCellInner(leftCell string, layout cellLayout) string {
	return `{\begingroup` +
		`\makebox[0pt][l]{` + leftCell + `}` +
		layout.spacing +
		`\begin{minipage}[t]{` + layout.width + `}` +
		layout.contentWrapper +
		`\end{minipage}` +
		`\endgroup}`
}

// wrapWithHyperlink wraps content with a hyperlink to the day's reference
func (d Day) wrapWithHyperlink(inner string) string {
	return `\hyperlink{` + d.ref() + `}{` + inner + `}`
}

// buildSimpleDayCell creates a simple day cell without tasks
func (d Day) buildSimpleDayCell(leftCell string) string {
	inner := `{\begingroup\makebox[0pt][l]{` + leftCell + `}\endgroup}`
	return d.wrapWithHyperlink(inner)
}

// ============================================================================
// TASK RENDERING - SPANNING TASKS
// ============================================================================

// renderSpanningTaskOverlay creates a task overlay with proper vertical stacking
// Uses track-based positioning to prevent visual overlap of multi-day tasks
func (d Day) renderSpanningTaskOverlay() *TaskOverlay {
	dayDate := d.getDayDate()
	activeTasks, maxCols := d.findActiveTasks(dayDate)

	if len(activeTasks) == 0 {
		return nil
	}

	// Assign tracks to ALL active tasks (including continuing ones)
	// This ensures consistent track assignments across days
	trackAssignments := d.assignTaskTracks(activeTasks)

	// Build rendering lists: tasks that START today vs. tasks that CONTINUE today
	var startingTasks []*SpanningTask
	var continuingTasks []*SpanningTask

	// Categorize active tasks
	for _, task := range activeTasks {
		start := d.getTaskStartDate(task)
		if dayDate.Equal(start) {
			// This task starts today
			startingTasks = append(startingTasks, task)
		} else {
			// This task is continuing from a previous day
			continuingTasks = append(continuingTasks, task)
		}
	}

	// Combine all tasks that need rendering (starting tasks get full rendering, continuing tasks get continuation indicators)
	var allTasksToRender []*SpanningTask
	renderingTypes := make(map[*SpanningTask]string) // "start" or "continue"

	// Add starting tasks
	for _, task := range startingTasks {
		allTasksToRender = append(allTasksToRender, task)
		renderingTypes[task] = "start"
	}

	// Add continuing tasks that should show visual indicators
	for _, task := range continuingTasks {
		allTasksToRender = append(allTasksToRender, task)
		renderingTypes[task] = "continue"
	}

	// Sort tasks by their assigned track (lowest track first, renders at bottom)
	sortedTasks := make([]*SpanningTask, len(allTasksToRender))
	copy(sortedTasks, allTasksToRender)
	for i := 0; i < len(sortedTasks)-1; i++ {
		for j := 0; j < len(sortedTasks)-i-1; j++ {
			track1 := trackAssignments[sortedTasks[j].ID]
			track2 := trackAssignments[sortedTasks[j+1].ID]
			if track1 > track2 {
				sortedTasks[j], sortedTasks[j+1] = sortedTasks[j+1], sortedTasks[j]
			}
		}
	}

	// Render task pills with vertical offsets based on track
	var pillContents []string

	for i, task := range sortedTasks {
		renderType := renderingTypes[task]

		// Skip rendering text for continuing tasks - just show the colored bar
		if renderType == "continue" {
			// Don't render anything for continuing tasks
			// The visual bar will span automatically via the cols parameter
			continue
		}

		// Render starting task (original logic)
		taskName := d.escapeLatexSpecialChars(task.Name)
		if d.isMilestoneSpanningTask(task) {
			taskName = "★ " + taskName
		}

		objective := ""
		if task.Description != "" {
			objective = d.escapeLatexSpecialChars(task.Description)
		}

		taskColor := core.HexToRGB(task.Color)
		if taskColor == "" {
			taskColor = core.Defaults.DefaultTaskColor
		}

		// Add spacing between stacked tasks (except for the first task)
		var spacing string
		if i > 0 {
			spacing = `\vspace{1mm}` // Add 1mm spacing between stacked tasks
		}

		// Choose appropriate macro based on whether task is a milestone
		var macroName string
		if task.IsMilestone {
			macroName = `\MilestoneTaskOverlayBox`
		} else {
			macroName = `\TaskOverlayBox`
		}

		// Use appropriate macro - LaTeX will stack naturally with spacing
		pillContent := spacing + fmt.Sprintf(`%s{%s}{%s}{%s}`,
			macroName,
			taskColor,
			taskName,
			objective)
		pillContents = append(pillContents, pillContent)
	}

	// Join the pills - they will stack naturally bottom-to-top
	content := strings.Join(pillContents, "")

	return &TaskOverlay{
		content: content,
		cols:    maxCols,
	}
}

// ============================================================================
// HELPER FUNCTIONS - DATE AND TASK UTILITIES
// ============================================================================

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

// findActiveTasks finds ALL tasks that should reserve vertical space on this day
// This includes:
// 1. Tasks that START on this day (will show task bar)
// 2. Tasks that STARTED EARLIER but are still active (need space but don't show bar)
// This ensures proper vertical stacking to prevent visual overlap
func (d Day) findActiveTasks(dayDate time.Time) ([]*SpanningTask, int) {
	var activeTasks []*SpanningTask
	var maxCols int
	seen := make(map[*SpanningTask]bool)

	for _, task := range d.Tasks {
		start := d.getTaskStartDate(task)
		end := d.getTaskEndDate(task)

		// Include task if it's active on this day (either starting or continuing)
		if d.isTaskActiveOnDay(dayDate, start, end) && !seen[task] {
			activeTasks = append(activeTasks, task)
			seen[task] = true

			// Calculate columns differently based on whether task starts today
			var cols int
			if dayDate.Equal(start) {
				// Task starts today: span from today to end (or end of week)
				cols = d.calculateTaskSpanColumns(dayDate, end)
			} else {
				// Task started earlier: calculate remaining span
				cols = d.calculateRemainingSpanColumns(dayDate, end)
			}

			if cols > maxCols {
				maxCols = cols
			}
		}
	}

	// Sort tasks by start date (earlier tasks appear first/on bottom)
	activeTasks = d.sortTasksByStartDate(activeTasks)

	return activeTasks, maxCols
}

// assignTaskTracks assigns vertical tracks to tasks to prevent visual overlap
// Returns a map of task ID to track number (0-based, 0 is bottom)
func (d Day) assignTaskTracks(tasks []*SpanningTask) map[string]int {
	trackAssignments := make(map[string]int)

	// For each task, find the lowest available track
	for _, task := range tasks {
		track := d.findLowestAvailableTrackForTask(task, trackAssignments)
		trackAssignments[task.ID] = track
	}

	return trackAssignments
}

// findLowestAvailableTrackForTask finds the lowest track that doesn't conflict with already-assigned tasks
func (d Day) findLowestAvailableTrackForTask(task *SpanningTask, existing map[string]int) int {
	taskStart := d.getTaskStartDate(task)
	taskEnd := d.getTaskEndDate(task)

	// Check each track starting from 0
	for track := 0; track < 100; track++ {
		occupied := false

		// Check if any existing task on this track overlaps with our task
		for otherTaskID, otherTrack := range existing {
			if otherTrack != track {
				continue // Different track, no conflict
			}

			// Find the other task
			for _, otherTask := range d.Tasks {
				if otherTask.ID == otherTaskID {
					otherStart := d.getTaskStartDate(otherTask)
					otherEnd := d.getTaskEndDate(otherTask)

					// Check if date ranges overlap
					if d.dateRangesOverlap(taskStart, taskEnd, otherStart, otherEnd) {
						occupied = true
						break
					}
				}
			}

			if occupied {
				break
			}
		}

		if !occupied {
			return track
		}
	}

	return 0 // Fallback
}

// dateRangesOverlap checks if two date ranges overlap
func (d Day) dateRangesOverlap(start1, end1, start2, end2 time.Time) bool {
	// Two ranges overlap if: start1 <= end2 AND start2 <= end1
	return !start1.After(end2) && !start2.After(end1)
}

// calculateRemainingSpanColumns calculates how many columns a continuing task spans
// from the current day to its end (or end of week, whichever is sooner)
func (d Day) calculateRemainingSpanColumns(dayDate, end time.Time) int {
	idxMonFirst := (int(dayDate.Weekday()) + 6) % 7 // Monday=0
	remainInRow := 7 - idxMonFirst
	daysLeft := int(end.Sub(dayDate).Hours()/24) + 1

	if daysLeft < 1 {
		return 1 // At least show on current day
	}
	if daysLeft > remainInRow {
		return remainInRow
	}
	return daysLeft
}

// sortTasksByStartDate sorts tasks by their start date (earliest first)
func (d Day) sortTasksByStartDate(tasks []*SpanningTask) []*SpanningTask {
	sorted := make([]*SpanningTask, len(tasks))
	copy(sorted, tasks)

	// Sort by start date (earliest first)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].StartDate.After(sorted[j+1].StartDate) {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// isMilestoneSpanningTask checks if a task is a milestone
func (d Day) isMilestoneSpanningTask(task *SpanningTask) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(task.Description)), "MILESTONE:")
}

// ============================================================================
// HELPER FUNCTIONS - LATEX UTILITIES
// ============================================================================

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

// ============================================================================
// WEEK STRUCTURES AND METHODS
// ============================================================================

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
		week.Days[i] = Day{Time: ptr, Tasks: nil, Cfg: cfg}
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

			week.Days[i] = Day{Time: ptr, Tasks: nil, Cfg: cfg}
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
			week.Days[j] = Day{Time: ptr, Tasks: nil, Cfg: cfg}
			ptr = ptr.AddDate(0, 0, 1)
		}
		weeks = append(weeks, week)
	}

	return weeks
}

// ============================================================================
// MONTH STRUCTURES AND METHODS
// ============================================================================

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

	// Use config helper for header angle size offset
	headerAngleOffset := m.Cfg.GetHeaderAngleSizeOffset()
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
			// Large mode: use zero-width paragraph column to force minimal width
			weekAlign = `|l!{\vrule width \myLenLineThicknessThick}`
			days = `@{}X@{}|`
		}

		return `\begin{tabularx}{\linewidth}{` + weekAlign + `*{7}{` + days + `}}`
	}

	return `\begin{tabular}[t]{c|*{7}{c}}`
}

func (m *Month) EndTable(typ interface{}) string {
	typStr, ok := typ.(string)
	if !ok || typStr == "tabularx" {
		// Check if we're in a tcolorbox (large calendar)
		// This is a simple check - in practice, we'd need to track state
		return `\end{tabularx}\end{tcolorbox}`
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
			// Add vertical padding with \rule for equal top/bottom spacing
			name = `\hfil{}\rule{0pt}{2.5ex}\rule[-1ex]{0pt}{0pt}` + name
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

func (m *Month) GetTaskColors() map[string]string {
	colorMap := make(map[string]string)

	// Only add colors for task categories that are actually present in this month
	for _, week := range m.Weeks {
		for _, day := range week.Days {
			for _, task := range day.Tasks {
				if task.Category != "" {
					color := core.GenerateCategoryColor(task.Category)
					if color != "" {
						// Convert to RGB for LaTeX compatibility
						// Escape LaTeX special characters in category name
						escapedCategory := escapeLatexSpecialChars(task.Category)
						colorMap[core.HexToRGB(color)] = escapedCategory
					}
				}
			}
		}
	}

	return colorMap
}

// PhaseGroup represents a phase with its sub-phases and colors for the legend
type PhaseGroup struct {
	PhaseNumber string
	PhaseName   string
	SubPhases   []SubPhaseLegendItem
}

// SubPhaseLegendItem represents a sub-phase entry in the legend
type SubPhaseLegendItem struct {
	Name  string
	Color string // RGB format for LaTeX
}

// GetTaskColorsByPhase returns tasks grouped by phase for a structured legend
func (m *Month) GetTaskColorsByPhase() []PhaseGroup {
	// Map to track unique phase -> subphase -> color
	phaseMap := make(map[string]map[string]string)
	phaseNames := make(map[string]string) // phase number -> full phase description

	// Collect all unique phase/subphase combinations in this month
	for _, week := range m.Weeks {
		for _, day := range week.Days {
			for _, task := range day.Tasks {
				if task.Phase != "" && task.SubPhase != "" {
					// Initialize phase map if needed
					if phaseMap[task.Phase] == nil {
						phaseMap[task.Phase] = make(map[string]string)
					}

					// Get color for this subphase (category)
					color := core.GenerateCategoryColor(task.SubPhase)
					if color != "" {
						phaseMap[task.Phase][task.SubPhase] = core.HexToRGB(color)
					}

					// Store phase name (use first non-empty one)
					if phaseNames[task.Phase] == "" && task.SubPhase != "" {
						phaseNames[task.Phase] = task.SubPhase
					}
				}
			}
		}
	}

	// Convert to sorted structure
	var phases []PhaseGroup
	phaseOrder := []string{"1", "2", "3", "4"} // Typical phase numbering

	for _, phaseNum := range phaseOrder {
		if subPhases, exists := phaseMap[phaseNum]; exists {
			phase := PhaseGroup{
				PhaseNumber: phaseNum,
				PhaseName:   getPhaseDescription(phaseNum),
			}

			// Add each subphase
			for subPhaseName, color := range subPhases {
				phase.SubPhases = append(phase.SubPhases, SubPhaseLegendItem{
					Name:  escapeLatexSpecialChars(subPhaseName),
					Color: color,
				})
			}

			phases = append(phases, phase)
		}
	}

	return phases
}

// getPhaseDescription returns a human-readable phase description
func getPhaseDescription(phaseNum string) string {
	descriptions := map[string]string{
		"1": "Phase 1: Proposal \\& Setup",
		"2": "Phase 2: Research \\& Data Collection",
		"3": "Phase 3: Publications",
		"4": "Phase 4: Dissertation",
	}

	if desc, exists := descriptions[phaseNum]; exists {
		return desc
	}
	return "Phase " + phaseNum
}

// ============================================================================
// YEAR AND QUARTER STRUCTURES
// ============================================================================

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

// ============================================================================
// SPANNING TASK DATA STRUCTURES AND UTILITIES
// ============================================================================

// SpanningTask represents a task that spans multiple days
type SpanningTask struct {
	ID          string
	Name        string
	Description string
	Phase       string // Added: Phase number
	SubPhase    string // Added: Sub-Phase description
	Category    string
	StartDate   time.Time
	EndDate     time.Time
	Color       string
	Progress    int    // Progress percentage (0-100)
	Status      string // Task status
	Assignee    string // Task assignee
	IsMilestone bool   // Whether this is a milestone task
}

// CreateSpanningTask creates a new spanning task from basic task data
func CreateSpanningTask(task core.Task, startDate, endDate time.Time) SpanningTask {
	// * Use Sub-Phase as category for better granularity
	color := core.GenerateCategoryColor(task.Category)

	return SpanningTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Phase:       task.Phase,    // * Added: Include Phase
		SubPhase:    task.SubPhase, // * Added: Include SubPhase
		Category:    task.Category, // * Fixed: Use Category field
		StartDate:   startDate,
		EndDate:     endDate,
		Color:       color,
		Progress:    0,                // Default progress
		Status:      task.Status,      // * Fixed: Use actual Status field
		Assignee:    task.Assignee,    // * Fixed: Use actual Assignee field
		IsMilestone: task.IsMilestone, // * Added: Pass milestone status
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
							week.Days[i].Tasks = append(week.Days[i].Tasks, &taskCopy)
							break
						}
					}
				}
			}
			current = current.AddDate(0, 0, 1)
		}
	}
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
