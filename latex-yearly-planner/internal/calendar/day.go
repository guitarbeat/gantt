package calendar

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/kudrykv/latex-yearly-planner/internal/header"
	"github.com/kudrykv/latex-yearly-planner/internal/latex"
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
		// Compact day number block drawn as a small tabular, overlaid on the left
		leftCell := `\begin{tabular}{@{}p{5mm}@{}|}\hfil{}` + day + `\\ \hline\end{tabular}`

		// Right side: bars for spanning tasks and/or regular tasks list
		var rightLines []string
		// Optional overlay for multi-day spanning task rendered once on the start day
		overlayStart := false
		overlayCols := 1
		overlayContent := ""

		// Spanning task rendering
		if len(d.SpanningTasks) > 0 {
			dayDate := time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), 0, 0, 0, 0, time.UTC)
			for _, task := range d.SpanningTasks {
				start := time.Date(task.StartDate.Year(), task.StartDate.Month(), task.StartDate.Day(), 0, 0, 0, 0, time.UTC)
				end := time.Date(task.EndDate.Year(), task.EndDate.Month(), task.EndDate.Day(), 0, 0, 0, 0, time.UTC)
				if dayDate.Before(start) || dayDate.After(end) {
					continue
				}

				if dayDate.Equal(start) {
					// Start day: render a single overlay block that visually spans multiple columns in the week row
					// Determine how many columns we can span in this row (up to week boundary or task end)
					idxMonFirst := (int(dayDate.Weekday()) + 6) % 7 // Monday=0
					remainInRow := 7 - idxMonFirst
					totalRemain := int(end.Sub(start).Hours()/24) + 1
					if totalRemain < 1 {
						totalRemain = 1
					}
					overlayCols = totalRemain
					if overlayCols > remainInRow {
						overlayCols = remainInRow
					}

					// Build name/description separately
					nameText := strings.TrimSpace(task.Name)
					descText := strings.TrimSpace(task.Description)

					// Compose overlay content with a single highlight bar and a legible colorboxed block
					// - Bar uses task color
					// - Background uses a light tint of the task color
					// - Text is black; name bold on top, description below
					boxBody := `\parbox{\linewidth}{` + `{\color{black}\textbf{\small ` + nameText + `}}`
					if descText != "" {
						boxBody += `\\[-0.15ex]{\color{black}\footnotesize ` + descText + `}`
					}
					boxBody += `}`

					overlayContent = `\textcolor{` + task.Color + `}{\rule{\linewidth}{0.6pt}}` + `\\[0.25ex]` +
						`{\begingroup\setlength{\fboxsep}{2pt}\colorbox{` + task.Color + `!12}{` + boxBody + `}\endgroup}`
					overlayStart = true
				} else {
					// Mid/end days: do not add duplicate bars/text; the overlay from the start day will visually cover
				}
			}
		}

		// Regular (non-spanning) tasks, if any
		if tasks := d.TasksForDay(); tasks != "" {
			if len(rightLines) > 0 {
				// add a subtle separator, avoid custom macros to prevent color errors
				rightLines = append(rightLines, `\vspace{0.1ex}\textcolor{black!30}{\rule{\linewidth}{0.3pt}}`)
			}
			rightLines = append(rightLines, `\footnotesize{`+tasks+`}`)
		}

		// If we built an overlay for a spanning task start, render it to span multiple columns
		if overlayStart {
			// Compute width across overlayCols columns using \dimexpr N\linewidth
			width := `\dimexpr ` + strconv.Itoa(overlayCols) + `\linewidth-5mm\relax`
			return `\hyperlink{` + d.ref() + `}{` +
				`{\begingroup` +
				`\makebox[0pt][l]{` + leftCell + `}` +
				`\makebox[0pt][l]{` + `\hspace*{5mm}` + `\begin{minipage}[t]{` + width + `}\raggedright` + overlayContent + `\end{minipage}` + `}` +
				`\endgroup}` +
				`}`
		}

		if len(rightLines) > 0 {
			right := strings.Join(rightLines, `\\[0.25ex]`)
			// Use an overlayed left mini-tabular and a right minipage to avoid &/\\ at outer level
			return `\hyperlink{` + d.ref() + `}{` +
				`{\begingroup` +
				`\makebox[0pt][l]{` + leftCell + `}` +
				`\hspace*{5mm}` +
				`\begin{minipage}[t]{\dimexpr\linewidth-5mm\relax}\raggedright` +
				right +
				`\end{minipage}` +
				`\endgroup}` +
				`}`
		}

		// No tasks: just the compact day number
		return `\hyperlink{` + d.ref() + `}{` + leftCell + `}`
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

func (d Day) Next() Day { return d.Add(1) }
func (d Day) Prev() Day { return d.Add(-1) }

func (d Day) NextExists() bool { return d.Time.Month() < time.December || d.Time.Day() < 31 }
func (d Day) PrevExists() bool { return d.Time.Month() > time.January || d.Time.Day() > 1 }

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

func (d Day) Quarter() int      { return int(math.Ceil(float64(d.Time.Month()) / 3.)) }
func (d Day) Month() time.Month { return d.Time.Month() }

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
		safeCat := strings.ReplaceAll(task.Category, "-", "\\allowbreak-\\allowbreak")
		cat := "\\textbf{\\scriptsize[" + safeCat + "]}"
		taskStr := cat + " " + task.Name
		taskStrings = append(taskStrings, taskStr)
	}
	return strings.Join(taskStrings, "\\\\")
}

// segmentByWords splits a sentence into N roughly equal word segments and returns the idx-th segment.
// If idx is out of range or there are no words left for that segment, returns an empty string.
// (segmentation helper removed; overlay approach renders a single block across columns)
