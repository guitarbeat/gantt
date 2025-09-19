package calendar

import (
	"math"
	"strconv"
	"strings"
	"time"

	"phd-dissertation-planner/internal/header"
	"phd-dissertation-planner/internal/latex"
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
			return latex.EmphCell(day)
		}
	}

	return latex.Link(d.ref(), day)
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
	return latex.Link(d.ref(), strconv.Itoa(d.Time.Day())+", "+d.Time.Weekday().String())
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

	dayItem := header.NewTextItem(d.Time.Format(dayLayout)).RefText(d.Time.Format(time.RFC3339))
	items := header.Items{
		header.NewIntItem(d.Time.Year()),
		header.NewTextItem("Q" + strconv.Itoa(int(math.Ceil(float64(d.Time.Month())/3.)))),
		header.NewMonthItem(d.Time.Month()).Shorten(shorten),
		header.NewTextItem("Week " + strconv.Itoa(wn)).RefPrefix(wpref),
	}

	if len(leaf) > 0 {
		items = append(items, dayItem, header.NewTextItem(leaf).RefText(prefix+d.ref()).Ref(true))
	} else {
		items = append(items, dayItem.Ref(true))
	}

	return items.Table(true)
}

// LinkLeaf creates a link with a leaf text
func (d Day) LinkLeaf(prefix, leaf string) string {
	return latex.Link(prefix+d.ref(), leaf)
}

// PrevNext creates navigation items for previous and next days
func (d Day) PrevNext(prefix string) header.Items {
	items := header.Items{}

	if d.PrevExists() {
		prev := d.Prev()
		items = append(items, header.NewTextItem(prev.Time.Format("Mon, 2")).RefText(prefix+prev.ref()))
	}

	if d.NextExists() {
		next := d.Next()
		items = append(items, header.NewTextItem(next.Time.Format("Mon, 2")).RefText(prefix+next.ref()))
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
		day = latex.Link(d.ref(), day)
	}

	anglesize := `\dimexpr\myLenHeaderResizeBox-0.86pt`
	var ll, rl string
	var r1, r2 []string
	if d.PrevExists() {
		ll = "l"
		leftNavBox := latex.ResizeBoxW(anglesize, `$\langle$`)
		r1 = append(r1, latex.Multirow(2, latex.Hyperlink(d.Prev().ref(prefix), leftNavBox)))
		r2 = append(r2, "")
	}
	r1 = append(r1, latex.Multirow(2, latex.ResizeBoxW(`\myLenHeaderResizeBox`, day)))
	r2 = append(r2, "")
	r1 = append(r1, latex.Bold(d.Time.Weekday().String()))
	r2 = append(r2, d.Time.Month().String())
	if d.NextExists() {
		rl = "l"
		rightNavBox := latex.ResizeBoxW(anglesize, `$\rangle$`)
		r1 = append(r1, latex.Multirow(2, latex.Hyperlink(d.Next().ref(prefix), rightNavBox)))
		r2 = append(r2, "")
	}
	contents := strings.Join(r1, ` & `) + `\\` + "\n" + strings.Join(r2, ` & `)
	return latex.Hypertarget(prefix+d.ref(), "") + latex.Tabular("@{}"+ll+"l|l"+rl, contents)
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