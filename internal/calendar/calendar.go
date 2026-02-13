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
	"sort"
	"strconv"
	"strings"
	"time"

	"phd-dissertation-planner/internal/core"
	"phd-dissertation-planner/internal/templates"
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

	// Combine all tasks that need rendering (starting tasks get full rendering, continuing tasks get continuation indicators)
	// Use a struct to avoid map lookups
	type RenderedTask struct {
		Task  *SpanningTask
		Track int
		Type  string // "start" or "continue"
	}
	var allTasksToRender = make([]RenderedTask, 0, len(activeTasks))

	// Categorize active tasks
	for i, task := range activeTasks {
		track := trackAssignments[i]
		start := d.getTaskStartDate(task)
		if dayDate.Equal(start) {
			// This task starts today
			allTasksToRender = append(allTasksToRender, RenderedTask{task, track, "start"})
		} else {
			// This task is continuing from a previous day
			allTasksToRender = append(allTasksToRender, RenderedTask{task, track, "continue"})
		}
	}

	// Sort tasks by their assigned track (lowest track first, renders at bottom)
	sort.Slice(allTasksToRender, func(i, j int) bool {
		return allTasksToRender[i].Track < allTasksToRender[j].Track
	})

	// Render task pills with vertical offsets based on track
	// Use strings.Builder for efficient string concatenation
	var sb strings.Builder
	// Pre-allocate buffer if possible, but exact size is unknown.
	// Average pill is maybe 100-200 bytes.

	for i, rt := range allTasksToRender {
		task := rt.Task

		// Skip rendering text for continuing tasks - just show the colored bar
		if rt.Type == "continue" {
			// Don't render anything for continuing tasks
			// The visual bar will span automatically via the cols parameter
			continue
		}

		// Render starting task (original logic)
		// Optimization: Use pre-calculated escaped name
		taskName := task.EscapedName
		// UX/A11y: Use accessible star icon for milestones
		// Check both the boolean flag and legacy description prefix for backward compatibility
		if task.IsMilestone || d.isMilestoneSpanningTask(task) {
			taskName = `\BeginAccSupp{method=pdfstringdef,unicode,ActualText={Milestone: } }★\EndAccSupp{} ` + taskName
		}

		objective := ""
		if task.Description != "" {
			// Optimization: Use pre-calculated escaped description
			objective = task.EscapedDescription
		}

		taskColor := core.HexToRGB(task.Color)
		if taskColor == "" {
			taskColor = core.Defaults.DefaultTaskColor
		}

		// Add spacing between stacked tasks (except for the first task)
		if i > 0 {
			sb.WriteString(`\vspace{1mm}`) // Add 1mm spacing between stacked tasks
		}

		// Choose appropriate macro based on whether task is a milestone
		var macroName string
		if task.IsMilestone {
			macroName = `\MilestoneTaskOverlayBox`
		} else {
			macroName = `\TaskOverlayBox`
		}

		// Use appropriate macro - LaTeX will stack naturally with spacing
		// Optimization: Write directly to builder
		fmt.Fprintf(&sb, `%s{%s}{%s}{%s}`,
			macroName,
			taskColor,
			taskName,
			objective)
	}

	return &TaskOverlay{
		content: sb.String(),
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
	return task.StartDate
}

// getTaskEndDate returns the task end date normalized to UTC midnight
func (d Day) getTaskEndDate(task *SpanningTask) time.Time {
	return task.EndDate
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

const MaxTaskTracks = 100

// findActiveTasks finds ALL tasks that should reserve vertical space on this day
// This includes:
// 1. Tasks that START on this day (will show task bar)
// 2. Tasks that STARTED EARLIER but are still active (need space but don't show bar)
// This ensures proper vertical stacking to prevent visual overlap
func (d Day) findActiveTasks(dayDate time.Time) ([]*SpanningTask, int) {
	// Bolt optimization: d.Tasks already contains only active tasks (populated by ApplySpanningTasksToMonth).
	// Since tasks are pre-sorted in ApplySpanningTasksToMonth, d.Tasks is already sorted by StartDate.
	// We iterate directly to avoid unnecessary allocations and checks.

	var maxCols int

	for _, task := range d.Tasks {
		start := d.getTaskStartDate(task)
		end := d.getTaskEndDate(task)

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

	return d.Tasks, maxCols
}

// assignTaskTracks assigns vertical tracks to tasks to prevent visual overlap
// Returns a slice of track numbers corresponding to the input tasks indices (0-based, 0 is bottom)
func (d Day) assignTaskTracks(tasks []*SpanningTask) []int {
	trackAssignments := make([]int, len(tasks))
	// Bolt Optimization: Use fixed-size array instead of map.
	// MaxTaskTracks limits the number of concurrent tasks we can stack visually.
	var tracksUsage [MaxTaskTracks][]*SpanningTask

	// For each task, find the lowest available track
	for i, task := range tasks {
		track := d.findLowestAvailableTrackForTask(task, &tracksUsage)
		trackAssignments[i] = track
		tracksUsage[track] = append(tracksUsage[track], task)
	}

	return trackAssignments
}

// findLowestAvailableTrackForTask finds the lowest track that doesn't conflict with already-assigned tasks
func (d Day) findLowestAvailableTrackForTask(task *SpanningTask, tracksUsage *[MaxTaskTracks][]*SpanningTask) int {
	taskStart := d.getTaskStartDate(task)
	taskEnd := d.getTaskEndDate(task)

	// Check each track starting from 0 up to MaxTaskTracks
	for track := 0; track < MaxTaskTracks; track++ {
		occupied := false

		// Optimization: Only check tasks already assigned to this specific track
		// Bolt Optimization: Direct array access (via pointer) avoids map lookup overhead
		tasksOnTrack := tracksUsage[track]
		for _, otherTask := range tasksOnTrack {
			otherStart := d.getTaskStartDate(otherTask)
			otherEnd := d.getTaskEndDate(otherTask)

			// Check if date ranges overlap
			if d.dateRangesOverlap(taskStart, taskEnd, otherStart, otherEnd) {
				occupied = true
				break
			}
		}

		if !occupied {
			return track
		}
	}

	return 0 // Fallback: if all tracks occupied, default to 0 (overlap will occur)
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

// sortTasksInPlace sorts tasks by their start date (earliest first)
// Optimization: Sorts in-place to avoid allocation
func (d Day) sortTasksInPlace(tasks []*SpanningTask) {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartDate.Before(tasks[j].StartDate)
	})
}

// isMilestoneSpanningTask checks if a task is a milestone
func (d Day) isMilestoneSpanningTask(task *SpanningTask) bool {
	return strings.HasPrefix(strings.ToUpper(strings.TrimSpace(task.Description)), "MILESTONE:")
}

// ============================================================================
// HELPER FUNCTIONS - LATEX UTILITIES
// ============================================================================

// latexReplacer is a pre-compiled replacer for LaTeX special characters
var latexReplacer = strings.NewReplacer(
	"\\", "\\textbackslash{}",
	"{", "\\{",
	"}", "\\}",
	"$", "\\$",
	"&", "\\&",
	"%", "\\%",
	"#", "\\#",
	"^", "\\textasciicircum{}",
	"_", "\\_",
	"~", "\\textasciitilde{}",
)

// EscapeLatexSpecialChars replaces special LaTeX characters with their escaped versions
func EscapeLatexSpecialChars(text string) string {
	// Optimize: Use pre-compiled replacer for significantly better performance
	// and fewer allocations compared to chained strings.ReplaceAll
	return latexReplacer.Replace(text)
}

// EscapeLatexSpecialChars escapes special LaTeX characters in text
func (d Day) EscapeLatexSpecialChars(text string) string {
	return EscapeLatexSpecialChars(text)
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

	text := `\rotatebox[origin=tr]{90}{\makebox[70pt][c]{Week ` + itoa + `}}`

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
			// Add vertical padding with \rule for equal top/bottom spacing
			name = `\hfil{}\rule{0pt}{2.5ex}\rule[-1ex]{0pt}{0pt}` + name
		} else {
			name = name[:1]
		}

		names = append(names, name)
	}

	return strings.Join(names, " & ")
}

func (m *Month) GetTaskColors() map[string]string {
	colorMap := make(map[string]string)
	seen := make(map[string]struct{})

	// Only add colors for task categories that are actually present in this month
	for _, week := range m.Weeks {
		for _, day := range week.Days {
			for _, task := range day.Tasks {
				if task.Category != "" {
					if _, ok := seen[task.Category]; ok {
						continue
					}
					seen[task.Category] = struct{}{}

					color := core.GenerateCategoryColor(task.Category)
					if color != "" {
						// Convert to RGB for LaTeX compatibility
						// Optimization: Use pre-calculated escaped category
						colorMap[core.HexToRGB(color)] = task.EscapedCategory
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
	// Map to track unique phases and their colors
	phaseMap := make(map[string]string) // phase name -> color
	phaseOrder := make([]string, 0)     // track order of phases

	// Collect all unique phases in this month
	for _, week := range m.Weeks {
		for _, day := range week.Days {
			for _, task := range day.Tasks {
				if task.Phase != "" {
					// Use the phase name directly (no number extraction needed)
					phaseName := task.Phase
					
					// Get color for this phase
					if _, exists := phaseMap[phaseName]; !exists {
						color := core.GenerateCategoryColor(phaseName)
						if color != "" {
							phaseMap[phaseName] = core.HexToRGB(color)
							phaseOrder = append(phaseOrder, phaseName)
						}
					}
				}
			}
		}
	}

	// Convert to sorted structure
	var phases []PhaseGroup

	for _, phaseName := range phaseOrder {
		if color, exists := phaseMap[phaseName]; exists {
			phase := PhaseGroup{
				PhaseNumber: "",  // No longer using phase numbers
				PhaseName:   EscapeLatexSpecialChars(phaseName),
			}

			// Add the phase as a "subphase" for consistency with template
			phase.SubPhases = append(phase.SubPhases, SubPhaseLegendItem{
				// Optimization: We could use task.EscapedPhase here, but we are iterating unique phases.
				// Since we don't have the task object here (we are iterating phaseOrder), we must escape again.
				// However, we can optimize by storing escaped name in the map or just doing it here.
				// Given the low cardinality of phases, this is not a hot path.
				Name:  EscapeLatexSpecialChars(phaseName),
				Color: color,
			})

			phases = append(phases, phase)
		}
	}

	return phases
}

// GetPhaseDescription returns a human-readable phase description
func GetPhaseDescription(phaseNum string) string {
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
	Phase       string // Combined: Phase with description
	Category    string
	StartDate   time.Time
	EndDate     time.Time
	Color       string
	Progress    int    // Progress percentage (0-100)
	Status      string // Task status
	Assignee    string // Task assignee
	IsMilestone bool   // Whether this is a milestone task

	// Memoized escaped strings for LaTeX rendering
	EscapedName        string
	EscapedDescription string
	EscapedCategory    string
	EscapedPhase       string
}

// CreateSpanningTask creates a new spanning task from basic task data
func CreateSpanningTask(task core.Task, startDate, endDate time.Time) SpanningTask {
	// * Use Sub-Phase as category for better granularity
	color := core.GenerateCategoryColor(task.Category)

	return SpanningTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Phase:       task.Phase,    // Combined: Phase with description
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
	// Optimization: Create a map of day numbers to Day pointers for O(1) lookup
	// This avoids nested loops searching for the correct day cell
	dayMap := make(map[int]*Day, 31)
	for _, week := range month.Weeks {
		for i := range week.Days {
			// Only map days that belong to the current month
			if week.Days[i].Time.Month() == month.Month &&
				week.Days[i].Time.Year() == month.Year.Number {
				dayMap[week.Days[i].Time.Day()] = &week.Days[i]
			}
		}
	}

	// Clone tasks to avoid mutating the input slice and to ensure memory ownership
	localTasks := make([]SpanningTask, len(tasks))
	copy(localTasks, tasks)

	// 1. Normalize dates for all tasks first
	// This ensures consistent comparison logic and avoids re-normalization
	for i := range localTasks {
		localTasks[i].StartDate = time.Date(localTasks[i].StartDate.Year(), localTasks[i].StartDate.Month(), localTasks[i].StartDate.Day(), 0, 0, 0, 0, time.UTC)
		localTasks[i].EndDate = time.Date(localTasks[i].EndDate.Year(), localTasks[i].EndDate.Month(), localTasks[i].EndDate.Day(), 0, 0, 0, 0, time.UTC)

		// Pre-calculate escaped strings to avoid repeated work during rendering
		localTasks[i].EscapedName = EscapeLatexSpecialChars(localTasks[i].Name)
		localTasks[i].EscapedDescription = EscapeLatexSpecialChars(localTasks[i].Description)
		localTasks[i].EscapedCategory = EscapeLatexSpecialChars(localTasks[i].Category)
		localTasks[i].EscapedPhase = EscapeLatexSpecialChars(localTasks[i].Phase)
	}

	// 2. Sort tasks by StartDate
	// This ensures that when we append tasks to days, they are already sorted by start date.
	// This eliminates the need to sort tasks in the hot loop (findActiveTasks) for every day.
	sort.Slice(localTasks, func(i, j int) bool {
		return localTasks[i].StartDate.Before(localTasks[j].StartDate)
	})

	monthStart := time.Date(month.Year.Number, month.Month, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, -1) // Last day of month

	// 3. Apply sorted tasks to the appropriate days in the month
	for i := range localTasks {
		// Optimization: Since tasks are sorted by StartDate, if current task starts after month end,
		// all subsequent tasks will also start after month end.
		if localTasks[i].StartDate.After(monthEnd) {
			break
		}

		// Quick check if task overlaps with this month
		// Use direct array access to avoid copying
		if localTasks[i].EndDate.Before(monthStart) {
			continue
		}

		// Calculate intersection range
		start := localTasks[i].StartDate
		if start.Before(monthStart) {
			start = monthStart
		}

		end := localTasks[i].EndDate
		if end.After(monthEnd) {
			end = monthEnd
		}

		// Use pointer to the task in the local slice.
		// The local slice persists because it is referenced by the Month pointers.
		taskPtr := &localTasks[i]

		// Iterate directly through the days in the range
		startDay := start.Day()
		endDay := end.Day()

		for d := startDay; d <= endDay; d++ {
			if dayCell, exists := dayMap[d]; exists {
				dayCell.Tasks = append(dayCell.Tasks, taskPtr)
			}
		}
	}
}
