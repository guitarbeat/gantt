package calendar

import (
	"math"
	"strconv"
	"strings"
	"time"

	"phd-dissertation-planner/internal/rendering"
)

// * LaTeX rendering constants
const (
	dayCellWidth           = "5mm"
	maxTaskChars           = 18
	maxTaskCharsCompact    = 15
	maxTaskCharsVeryCompact = 12
	maxTasksDisplay        = 2  // Reduced to prevent overlap
)

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
			return rendering.EmphCell(day)
		}
	}

	return rendering.Link(d.ref(), day)
}

// renderLargeDay renders the day cell for large (monthly) view with tasks and spanning tasks
func (d Day) renderLargeDay(day string) string {
	leftCell := d.buildDayNumberCell(day)

	// Check for spanning tasks that start on this day
	overlay := d.renderSpanningTaskOverlay()
	if overlay != nil {
		return d.buildSpanningTaskCell(leftCell, overlay)
	}

	// Check for regular tasks
	if tasks := d.TasksForDay(); tasks != "" {
		return d.buildRegularTaskCell(leftCell, tasks)
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
	return rendering.Link(d.ref(), strconv.Itoa(d.Time.Day())+", "+d.Time.Weekday().String())
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

	dayItem := rendering.NewTextItem(d.Time.Format(dayLayout)).RefText(d.Time.Format(time.RFC3339))
	items := rendering.Items{
		rendering.NewIntItem(d.Time.Year()),
		rendering.NewTextItem("Q" + strconv.Itoa(int(math.Ceil(float64(d.Time.Month())/3.)))),
		rendering.NewMonthItem(d.Time.Month()).Shorten(shorten),
		rendering.NewTextItem("Week " + strconv.Itoa(wn)).RefPrefix(wpref),
	}

	if len(leaf) > 0 {
		items = append(items, dayItem, rendering.NewTextItem(leaf).RefText(prefix+d.ref()).Ref(true))
	} else {
		items = append(items, dayItem.Ref(true))
	}

	return items.Table(true)
}

// LinkLeaf creates a link with a leaf text
func (d Day) LinkLeaf(prefix, leaf string) string {
	return rendering.Link(prefix+d.ref(), leaf)
}

// PrevNext creates navigation items for previous and next days
func (d Day) PrevNext(prefix string) rendering.Items {
	items := rendering.Items{}

	if d.PrevExists() {
		prev := d.Prev()
		items = append(items, rendering.NewTextItem(prev.Time.Format("Mon, 2")).RefText(prefix+prev.ref()))
	}

	if d.NextExists() {
		next := d.Next()
		items = append(items, rendering.NewTextItem(next.Time.Format("Mon, 2")).RefText(prefix+next.ref()))
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
		day = rendering.Link(d.ref(), day)
	}

	anglesize := `\dimexpr\myLenHeaderResizeBox-0.86pt`
	var ll, rl string
	var r1, r2 []string
	if d.PrevExists() {
		ll = "l"
		leftNavBox := rendering.ResizeBoxW(anglesize, `$\langle$`)
		r1 = append(r1, rendering.Multirow(2, rendering.Hyperlink(d.Prev().ref(prefix), leftNavBox)))
		r2 = append(r2, "")
	}
	r1 = append(r1, rendering.Multirow(2, rendering.ResizeBoxW(`\myLenHeaderResizeBox`, day)))
	r2 = append(r2, "")
	r1 = append(r1, rendering.Bold(d.Time.Weekday().String()))
	r2 = append(r2, d.Time.Month().String())
	if d.NextExists() {
		rl = "l"
		rightNavBox := rendering.ResizeBoxW(anglesize, `$\rangle$`)
		r1 = append(r1, rendering.Multirow(2, rendering.Hyperlink(d.Next().ref(prefix), rightNavBox)))
		r2 = append(r2, "")
	}
	contents := strings.Join(r1, ` & `) + `\\` + "\n" + strings.Join(r2, ` & `)
	return rendering.Hypertarget(prefix+d.ref(), "") + rendering.Tabular("@{}"+ll+"l|l"+rl, contents)
}

// * LaTeX cell construction functions

// buildDayNumberCell creates the basic day number cell
func (d Day) buildDayNumberCell(day string) string {
	return `\begin{tabular}{@{}p{` + dayCellWidth + `}@{}|}\hfil{}` + day + `\\ \hline\end{tabular}`
}

// buildSpanningTaskCell creates a cell with spanning task overlay
func (d Day) buildSpanningTaskCell(leftCell string, overlay *overlayInfo) string {
	width := `\dimexpr ` + strconv.Itoa(overlay.cols) + `\linewidth\relax`
	return `\hyperlink{` + d.ref() + `}{` +
		`{\begingroup` +
		`\makebox[0pt][l]{` + leftCell + `}` +
		`\makebox[0pt][l]{` + `\begin{tikzpicture}[overlay]` +
		`\node[anchor=north west, inner sep=0pt] at (0,0) {` + `\begin{minipage}[t]{` + width + `}` + overlay.content + `\end{minipage}` + `};` +
		`\end{tikzpicture}` + `}` +
		`\endgroup}` +
		`}`
}

// buildRegularTaskCell creates a cell with regular tasks
func (d Day) buildRegularTaskCell(leftCell, tasks string) string {
	return `\hyperlink{` + d.ref() + `}{` +
		`{\begingroup` +
		`\makebox[0pt][l]{` + leftCell + `}` +
		`\hspace*{` + dayCellWidth + `}` +
		`\begin{minipage}[t]{\dimexpr\linewidth\relax}` +
		`\footnotesize{` + tasks + `}` +
		`\end{minipage}` +
		`\endgroup}` +
		`}`
}

// buildSimpleDayCell creates a simple day cell without tasks
func (d Day) buildSimpleDayCell(leftCell string) string {
	return `\hyperlink{` + d.ref() + `}{` + leftCell + `}`
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