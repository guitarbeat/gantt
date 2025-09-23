package scheduler

import (
	"math"
	"strconv"
	"strings"
	"time"

	"phd-dissertation-planner/internal/common"
	"phd-dissertation-planner/templates"
)

// * LaTeX rendering constants
const (
	dayCellWidth            = "5mm"
	maxTaskChars            = 18
	maxTaskCharsCompact     = 15
	maxTaskCharsVeryCompact = 12
	maxTasksDisplay         = 2 // Reduced to prevent overlap
)

// * Day types and methods (from day.go)

type Days []*Day
type Day struct {
	Time          time.Time
	Tasks         []Task
	SpanningTasks []*SpanningTask
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
		return d.buildTaskCell(leftCell, overlay.content, true, overlay.cols)
	}

	// Check for regular tasks
	if tasks := d.TasksForDay(); tasks != "" {
		return d.buildTaskCell(leftCell, tasks, false, 0)
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

// Add creates a new day by adding the specified number of days
func (d Day) Add(days int) Day {
	return Day{Time: d.Time.AddDate(0, 0, days), Tasks: nil, SpanningTasks: nil}
}

// WeekLink creates a link for the week view
func (d Day) WeekLink() string {
	return templates.Link(d.ref(), strconv.Itoa(d.Time.Day())+", "+d.Time.Weekday().String())
}

// Breadcrumb creates a breadcrumb navigation for the day
func (d Day) Breadcrumb(prefix string, leaf string, shorten bool) string {
	wpref := ""
	_, wn := d.Time.ISOWeek()
	if wn > 50 && d.Time.Month() == time.January {
		wpref = "fw"
	}

	dayLayout := "Monday, 2"
	if shorten {
		dayLayout = "Mon, 2"
	}

	dayItem := templates.NewTextItem(d.Time.Format(dayLayout)).RefText(d.Time.Format(time.RFC3339))
	items := templates.Items{
		templates.NewIntItem(d.Time.Year()),
		templates.NewTextItem("Q" + strconv.Itoa(int(math.Ceil(float64(d.Time.Month())/3.)))),
		templates.NewMonthItem(d.Time.Month()).Shorten(shorten),
		templates.NewTextItem("Week " + strconv.Itoa(wn)).RefPrefix(wpref),
	}

	if len(leaf) > 0 {
		items = append(items, dayItem, templates.NewTextItem(leaf).RefText(prefix+d.ref()).Ref(true))
	} else {
		items = append(items, dayItem.Ref(true))
	}

	return items.Table(true)
}

// LinkLeaf creates a link with a leaf text
func (d Day) LinkLeaf(prefix, leaf string) string {
	return templates.Link(prefix+d.ref(), leaf)
}

// PrevNext creates navigation items for previous and next days
func (d Day) PrevNext(prefix string) templates.Items {
	items := templates.Items{}

	if d.PrevExists() {
		prev := d.Prev()
		items = append(items, templates.NewTextItem(prev.Time.Format("Mon, 2")).RefText(prefix+prev.ref()))
	}

	if d.NextExists() {
		next := d.Next()
		items = append(items, templates.NewTextItem(next.Time.Format("Mon, 2")).RefText(prefix+next.ref()))
	}

	return items
}

// Next returns the next day
func (d Day) Next() Day { return d.Add(1) }

// Prev returns the previous day
func (d Day) Prev() Day { return d.Add(-1) }

// NextExists checks if the next day exists
func (d Day) NextExists() bool { return d.Time.Month() < time.December || d.Time.Day() < 31 }

// PrevExists checks if the previous day exists
func (d Day) PrevExists() bool { return d.Time.Month() > time.January || d.Time.Day() > 1 }

// Quarter returns the quarter number for this day
func (d Day) Quarter() int { return int(math.Ceil(float64(d.Time.Month()) / 3.)) }

// Month returns the month for this day
func (d Day) Month() time.Month { return d.Time.Month() }

// HeadingMOS creates a heading for the month-overview-single view
func (d Day) HeadingMOS(prefix, leaf string) string {
	day := strconv.Itoa(d.Time.Day())
	if len(leaf) > 0 {
		day = templates.Link(d.ref(), day)
	}

	anglesize := `\dimexpr\myLenHeaderResizeBox-0.86pt`
	var ll, rl string
	var r1, r2 []string
	if d.PrevExists() {
		ll = "l"
		leftNavBox := templates.ResizeBoxW(anglesize, `$\langle$`)
		r1 = append(r1, templates.Multirow(2, templates.Hyperlink(d.Prev().ref(prefix), leftNavBox)))
		r2 = append(r2, "")
	}
	r1 = append(r1, templates.Multirow(2, templates.ResizeBoxW(`\myLenHeaderResizeBox`, day)))
	r2 = append(r2, "")
	r1 = append(r1, templates.Bold(d.Time.Weekday().String()))
	r2 = append(r2, d.Time.Month().String())
	if d.NextExists() {
		rl = "l"
		rightNavBox := templates.ResizeBoxW(anglesize, `$\rangle$`)
		r1 = append(r1, templates.Multirow(2, templates.Hyperlink(d.Next().ref(prefix), rightNavBox)))
		r2 = append(r2, "")
	}
	contents := strings.Join(r1, ` & `) + `\\` + "\n" + strings.Join(r2, ` & `)
	return templates.Hypertarget(prefix+d.ref(), "") + templates.Tabular("@{}"+ll+"l|l"+rl, contents)
}

// * LaTeX cell construction functions

// buildDayNumberCell creates the basic day number cell
func (d Day) buildDayNumberCell(day string) string {
	return `\begin{tabular}{@{}p{` + dayCellWidth + `}@{}|}\hfil{}` + day + `\\ \hline\end{tabular}`
}

// buildTaskCell creates a cell with either spanning tasks or regular tasks
func (d Day) buildTaskCell(leftCell, content string, isSpanning bool, cols int) string {
	var width, spacing, contentWrapper string

	if isSpanning {
		// Spanning task: use tikzpicture overlay with calculated width
		width = `\dimexpr ` + strconv.Itoa(cols) + `\linewidth\relax`
		spacing = `\makebox[0pt][l]{` + `\begin{tikzpicture}[overlay]` +
			`\node[anchor=north west, inner sep=0pt] at (0,0) {` + `\begin{minipage}[t]{` + width + `}` + content + `\end{minipage}` + `};` +
			`\end{tikzpicture}` + `}`
		contentWrapper = content
	} else {
		// Regular task: use hspace and footnotesize
		width = `\dimexpr\linewidth\relax`
		spacing = `\hspace*{` + dayCellWidth + `}`
		contentWrapper = `\footnotesize{` + content + `}`
	}

	return `{\begingroup` +
		`\makebox[0pt][l]{` + leftCell + `}` +
		spacing +
		`\begin{minipage}[t]{` + width + `}` +
		contentWrapper +
		`\end{minipage}` +
		`\endgroup}`
}

// buildSimpleDayCell creates a simple day cell without tasks
func (d Day) buildSimpleDayCell(leftCell string) string {
	return leftCell
}

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

	// Enhanced sorting: priority first, then by duration (shorter tasks first), then by start time
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			priority1 := d.getTaskPriority(sorted[j].Category, priorityOrder)
			priority2 := d.getTaskPriority(sorted[j+1].Category, priorityOrder)

			// First sort by priority
			if priority1 != priority2 {
				if priority1 > priority2 {
					sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
				}
				continue
			}

			// If same priority, sort by duration (shorter tasks first for better stacking)
			duration1 := sorted[j].EndDate.Sub(sorted[j].StartDate)
			duration2 := sorted[j+1].EndDate.Sub(sorted[j+1].StartDate)
			if duration1 > duration2 {
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

// * Week types and methods (from week.go)

type Weeks []*Week
type Week struct {
	Days [7]Day

	Weekday  time.Weekday
	Year     *Year
	Months   Months
	Quarters Quarters
}

func NewWeeksForMonth(wd time.Weekday, year *Year, qrtr *Quarter, month *Month) Weeks {
	ptr := time.Date(year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
	weekday := ptr.Weekday()
	shift := (7 + weekday - wd) % 7

	week := &Week{Weekday: wd, Year: year, Months: Months{month}, Quarters: Quarters{qrtr}}

	for i := shift; i < 7; i++ {
		week.Days[i] = Day{Time: ptr, Tasks: nil, SpanningTasks: nil}
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

			week.Days[i] = Day{Time: ptr, Tasks: nil, SpanningTasks: nil}
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
	// Calculate week number based on the first day of the week
	firstDay := w.Days[0].Time
	if firstDay.IsZero() {
		return 0
	}

	// Get the week number of the year
	_, week := firstDay.ISOWeek()
	return week
}

func NewWeeksForYear(wd time.Weekday, year *Year) Weeks {
	var weeks Weeks
	ptr := time.Date(year.Number, 1, 1, 0, 0, 0, 0, time.Local)
	weekday := ptr.Weekday()
	_ = (7 + weekday - wd) % 7

	for i := 0; i < 53; i++ {
		week := &Week{Weekday: wd, Year: year}
		for j := 0; j < 7; j++ {
			week.Days[j] = Day{Time: ptr, Tasks: nil, SpanningTasks: nil}
			ptr = ptr.AddDate(0, 0, 1)
		}
		weeks = append(weeks, week)
	}

	return weeks
}

func (w Week) Breadcrumb() string {
	return templates.Items{
		templates.NewIntItem(w.Year.Number),
		templates.NewTextItem("Week " + strconv.Itoa(w.weekNumber())),
	}.Table(true)
}

func (w Week) WeekLink() string {
	return templates.Link(w.ref(), "Week "+strconv.Itoa(w.weekNumber()))
}

func (w Week) ref(prefix ...string) string {
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return p + "week-" + strconv.Itoa(w.Year.Number) + "-" + strconv.Itoa(w.weekNumber())
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
}

func NewMonth(wd time.Weekday, year *Year, qrtr *Quarter, month time.Month) *Month {
	m := &Month{
		Year:    year,
		Quarter: qrtr,
		Month:   month,
		Weekday: wd,
	}

	m.Weeks = NewWeeksForMonth(wd, year, qrtr, m)

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

	anglesize := `\dimexpr\myLenHeaderResizeBox-0.86pt`
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

func (m *Month) GetTaskColors() map[string]string {
	colorMap := make(map[string]string)

	// Iterate through all weeks and days to find unique task colors
	for _, week := range m.Weeks {
		for _, day := range week.Days {
			// Check spanning tasks
			for _, task := range day.SpanningTasks {
				if task.Color != "" {
					colorMap[task.Color] = task.Category
				}
			}
			// Check regular tasks
			for _, task := range day.Tasks {
				if task.Category != "" {
					// Get color for this category
					color := getColorForCategory(task.Category)
					if color != "" {
						colorMap[color] = task.Category
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

func NewYear(wd time.Weekday, year int) *Year {
	out := &Year{Number: year}
	out.Weeks = NewWeeksForYear(wd, out)
	for q := 1; q <= 4; q++ {
		out.Quarters = append(out.Quarters, NewQuarter(wd, out, q))
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

func NewQuarter(wd time.Weekday, year *Year, quarter int) *Quarter {
	out := &Quarter{Number: quarter, Year: year}
	for m := 1; m <= 3; m++ {
		month := time.Month((quarter-1)*3 + m)
		out.Months = append(out.Months, NewMonth(wd, year, out, month))
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
	Priority    int
	Progress    int    // Progress percentage (0-100)
	Status      string // Task status
	Assignee    string // Task assignee
}

// CreateSpanningTask creates a new spanning task from basic task data
func CreateSpanningTask(task common.Task, startDate, endDate time.Time) SpanningTask {
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
