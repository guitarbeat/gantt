package calendar

import (
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
			// Create enhanced display for spanning tasks
			var taskInfo []string
			
			for _, task := range d.SpanningTasks {
				// Determine if this is the start of the task span
				dayDate := time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), 0, 0, 0, 0, time.UTC)
				taskStartDate := time.Date(task.StartDate.Year(), task.StartDate.Month(), task.StartDate.Day(), 0, 0, 0, 0, time.UTC)
				
				isStart := dayDate.Equal(taskStartDate)
				
				if isStart {
					// Show full task name and description on start day
					taskName := `\textbf{` + task.Name + `}`
					taskParts := []string{`\textcolor{` + task.Color + `}{\tiny ` + taskName + `}`}
					
					// Add description if present - use parbox for better wrapping
					if task.Description != "" {
						desc := `\textcolor{` + task.Color + `}{\scriptsize \parbox[t]{8mm}{` + task.Description + `}}`
						taskParts = append(taskParts, desc)
					}
					
					taskInfo = append(taskInfo, strings.Join(taskParts, `\\`))
				} else {
					// Show abbreviated indicator for continuation days
					indicator := `\textcolor{` + task.Color + `}{\tiny \textbf{` + string(task.Name[0]) + `}}`
					taskInfo = append(taskInfo, indicator)
				}
			}
			
			// Create a cell with colored background and better layout
			taskText := ""
			if len(taskInfo) > 0 {
				taskText = strings.Join(taskInfo, `\\`)
			}
			
			// Use wider cell for better text display
			cellContent := day
			if taskText != "" {
				cellContent = day + `\\` + taskText
			}
			
			// Add a colored background with better width
			color := d.SpanningTasks[0].Color // Use first task's color
			return `\hyperlink{` + d.ref() + `}{\colorbox{` + color + `!20}{\begin{tabular}{@{}p{8mm}@{}}\centering ` + cellContent + `\end{tabular}}}`
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
		// Create a well-formatted task entry with bold name and description
		var taskParts []string
		
		// Bold task name with category
		taskName := latex.Bold(task.Name)
		if task.Category != "" {
			taskName = latex.Bold("[" + task.Category + "] " + task.Name)
		}
		taskParts = append(taskParts, taskName)
		
		// Add description on new lines if present
		if task.Description != "" {
			// Use LaTeX text wrapping for long descriptions
			// Split description into multiple lines for better readability
			desc := task.Description
			
			// Use parbox for better text wrapping and smaller font
			descFormatted := "\\parbox[t]{10mm}{\\footnotesize " + desc + "}"
			taskParts = append(taskParts, descFormatted)
		}
		
		// Join task name and description with line breaks
		taskStr := strings.Join(taskParts, "\\\\")
		taskStrings = append(taskStrings, taskStr)
	}
	
	// Join multiple tasks with double line breaks for better separation
	return strings.Join(taskStrings, "\\\\[0.5ex]")
}
