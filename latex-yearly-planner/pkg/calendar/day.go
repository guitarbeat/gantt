package calendar

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/kudrykv/latex-yearly-planner/pkg/header"
	"github.com/kudrykv/latex-yearly-planner/pkg/latex"
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

func (d Day) Day(today, large interface{}) string {
	if d.Time.IsZero() {
		return ""
	}

	day := strconv.Itoa(d.Time.Day())

	if larg, _ := large.(bool); larg {
		// Check for spanning task overlays
		if len(d.SpanningTasks) > 0 {
			// Create overlays for all spanning tasks
			var overlays []string

			for i, task := range d.SpanningTasks {
				// Determine if this is the start, middle, or end of the task span
				// Normalize dates to midnight for comparison
				dayDate := time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), 0, 0, 0, 0, time.UTC)
				taskStartDate := time.Date(task.StartDate.Year(), task.StartDate.Month(), task.StartDate.Day(), 0, 0, 0, 0, time.UTC)
				taskEndDate := time.Date(task.EndDate.Year(), task.EndDate.Month(), task.EndDate.Day(), 0, 0, 0, 0, time.UTC)

				isStart := dayDate.Equal(taskStartDate)
				isEnd := dayDate.Equal(taskEndDate)

				// Calculate vertical offset for multiple tasks
				yOffset := float64(i) * 0.8 // 0.8em spacing between tasks

				// Create a TikZ overlay that doesn't affect cell height
				var overlay string

				if isStart {
					// Start of task - show full task name, description, and left-rounded rectangle
					taskName := task.Name
					if len(taskName) > 15 {
						taskName = taskName[:12] + "..."
					}

					// Prepare description text (truncate if too long)
					description := task.Description
					if len(description) > 25 {
						description = description[:22] + "..."
					}

					// Add progress bar if progress > 0
					progressBar := ""
					if task.Progress > 0 {
						progressWidth := float64(task.Progress) / 100.0 * 5.0 // 5mm max width
						progressBar = `\draw[fill=` + task.Color + `!70, draw=` + task.Color + `!90, line width=0.2pt] 
							(cell-left) ++(0,-` + fmt.Sprintf("%.1f", 0.15+yOffset) + `em) rectangle ([xshift=` + fmt.Sprintf("%.2f", progressWidth) + `mm, yshift=-` + fmt.Sprintf("%.1f", 0.25+yOffset) + `em]cell-left);`
					}

					overlay = `\begin{tikzpicture}[overlay, remember picture]
						\coordinate (cell-top) at (0,0);
						\coordinate (cell-left) at (-2.5mm,0);
						\coordinate (cell-right) at (2.5mm,0);
						
						% Draw colored background bar
						\draw[fill=` + task.Color + `!30, draw=` + task.Color + `!60, line width=0.3pt] 
							(cell-left) ++(0,-` + fmt.Sprintf("%.1f", 0.2+yOffset) + `em) rectangle (cell-right) ++(0,-` + fmt.Sprintf("%.1f", 0.8+yOffset) + `em);
						
						% Add progress bar if applicable
						` + progressBar + `
						
						% Add task name and description on start day
						\node[anchor=west, font=\tiny, color=` + task.Color + `!80] 
							at ([xshift=-2.3mm, yshift=-` + fmt.Sprintf("%.1f", 0.4+yOffset) + `em]cell-top) {` + taskName + `};
						\node[anchor=west, font=\scriptsize, color=` + task.Color + `!60] 
							at ([xshift=-2.3mm, yshift=-` + fmt.Sprintf("%.1f", 0.6+yOffset) + `em]cell-top) {` + description + `};
					\end{tikzpicture}`
				} else if isEnd {
					// End of task - show right-rounded rectangle, no text
					overlay = `\begin{tikzpicture}[overlay, remember picture]
						\coordinate (cell-left) at (-2.5mm,0);
						\coordinate (cell-right) at (2.5mm,0);
						
						% Draw colored background bar
						\draw[fill=` + task.Color + `!30, draw=` + task.Color + `!60, line width=0.3pt] 
							(cell-left) ++(0,-` + fmt.Sprintf("%.1f", 0.2+yOffset) + `em) rectangle (cell-right) ++(0,-` + fmt.Sprintf("%.1f", 0.6+yOffset) + `em);
					\end{tikzpicture}`
				} else {
					// Middle of task - show plain rectangle, no text
					overlay = `\begin{tikzpicture}[overlay, remember picture]
						\coordinate (cell-left) at (-2.5mm,0);
						\coordinate (cell-right) at (2.5mm,0);
						
						% Draw colored background bar with no rounded corners
						\draw[fill=` + task.Color + `!30, draw=` + task.Color + `!60, line width=0.3pt] 
							(cell-left) ++(0,-` + fmt.Sprintf("%.1f", 0.2+yOffset) + `em) rectangle (cell-right) ++(0,-` + fmt.Sprintf("%.1f", 0.6+yOffset) + `em);
					\end{tikzpicture}`
				}

				overlays = append(overlays, overlay)
			}

			// Combine all overlays
			combinedOverlay := strings.Join(overlays, "")
			return `\hyperlink{` + d.ref() + `}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}` + day + `\\ \hline\end{tabular}}` + combinedOverlay
		}

		// For large view, include regular tasks if any
		tasks := d.TasksForDay()
		if tasks != "" {
			// Use wider column and better formatting for tasks with proper text wrapping
			return `\hyperlink{` + d.ref() + `}{\begin{tabular}{@{}p{15mm}@{}|}\centering\textbf{` + day + `}\\ \hline ` + tasks + `\end{tabular}}`
		}
		return `\hyperlink{` + d.ref() + `}{\begin{tabular}{@{}p{15mm}@{}|}\centering\textbf{` + day + `}\\ \hline\end{tabular}}`
	}

	if td, ok := today.(Day); ok {
		if d.Time.Equal(td.Time) {
			return latex.EmphCell(day)
		}
	}

	return latex.Link(d.ref(), day)
}

func (d Day) ref(prefix ...string) string {
	p := ""

	if len(prefix) > 0 {
		p = prefix[0]
	}

	return p + d.Time.Format(time.RFC3339)
}

func (d Day) Add(days int) Day {
	return Day{Time: d.Time.AddDate(0, 0, days), Tasks: nil, SpanningTasks: nil}
}

func (d Day) WeekLink() string {
	return latex.Link(d.ref(), strconv.Itoa(d.Time.Day())+", "+d.Time.Weekday().String())
}

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

func (d Day) LinkLeaf(prefix, leaf string) string {
	return latex.Link(prefix+d.ref(), leaf)
}

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

func (d Day) Next() Day {
	return d.Add(1)
}

func (d Day) Prev() Day {
	return d.Add(-1)
}

func (d Day) NextExists() bool {
	return d.Time.Month() < time.December || d.Time.Day() < 31
}

func (d Day) PrevExists() bool {
	return d.Time.Month() > time.January || d.Time.Day() > 1
}

func (d Day) Hours(bottom, top int) Days {
	moment := time.Date(1, 1, 1, bottom, 0, 0, 0, time.Local)
	list := make(Days, 0, top-bottom+1)

	for i := bottom; i <= top; i++ {
		list = append(list, &Day{Time: moment, Tasks: nil, SpanningTasks: nil})
		moment = moment.Add(time.Hour)
	}

	return list
}

func (d Day) FormatHour(ampm interface{}) string {
	if doAmpm, _ := ampm.(bool); doAmpm {
		return d.Time.Format("3 PM")
	}

	return d.Time.Format("15")
}

func (d Day) Quarter() int {
	return int(math.Ceil(float64(d.Time.Month()) / 3.))
}

func (d Day) Month() time.Month {
	return d.Time.Month()
}

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

// TasksForDay returns a formatted string of tasks for this day
func (d Day) TasksForDay() string {
	if len(d.Tasks) == 0 {
		return ""
	}

	var taskStrings []string
	for _, task := range d.Tasks {
		// Format: [Category] Task Name
		taskStr := "\\textbf{[" + task.Category + "]} " + task.Name
		taskStrings = append(taskStrings, taskStr)
	}

	return strings.Join(taskStrings, "\\\\")
}
