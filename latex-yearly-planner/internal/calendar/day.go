package calendar

import (
	"math"
	"strconv"
	"strings"
	"time"

	"latex-yearly-planner/internal/header"
	"latex-yearly-planner/internal/latex"
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
	leftCell := `\begin{tabular}{@{}p{5mm}@{}|}\hfil{}` + day + `\\ \hline\end{tabular}`

	// Check for spanning tasks that start on this day
	overlay := d.renderSpanningTaskOverlay()
	if overlay != nil {
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

	// Check for regular tasks
	if tasks := d.TasksForDay(); tasks != "" {
		return `\hyperlink{` + d.ref() + `}{` +
			`{\begingroup` +
			`\makebox[0pt][l]{` + leftCell + `}` +
			`\hspace*{5mm}` +
			`\begin{minipage}[t]{\dimexpr\linewidth\relax}` +
			`\footnotesize{` + tasks + `}` +
			`\end{minipage}` +
			`\endgroup}` +
			`}`
	}

	// No tasks: just the day number
	return `\hyperlink{` + d.ref() + `}{` + leftCell + `}`
}

// overlayInfo holds information about a spanning task overlay
type overlayInfo struct {
	content string
	cols    int
}

// renderSpanningTaskOverlay renders spanning task overlays for multiple tasks starting on this day
func (d Day) renderSpanningTaskOverlay() *overlayInfo {
	if len(d.SpanningTasks) == 0 {
		return nil
	}

	dayDate := time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), 0, 0, 0, 0, time.UTC)
	
	var startingTasks []*SpanningTask
	var maxCols int
	
	// Find all tasks that start on this day
	for _, task := range d.SpanningTasks {
		start := time.Date(task.StartDate.Year(), task.StartDate.Month(), task.StartDate.Day(), 0, 0, 0, 0, time.UTC)
		end := time.Date(task.EndDate.Year(), task.EndDate.Month(), task.EndDate.Day(), 0, 0, 0, 0, time.UTC)
		
		if dayDate.Before(start) || dayDate.After(end) {
			continue
		}

		if dayDate.Equal(start) {
			startingTasks = append(startingTasks, task)
			
			// Calculate span width for this task
			idxMonFirst := (int(dayDate.Weekday()) + 6) % 7 // Monday=0
			remainInRow := 7 - idxMonFirst
			totalRemain := int(end.Sub(start).Hours()/24) + 1
			if totalRemain < 1 {
				totalRemain = 1
			}
			overlayCols := totalRemain
			if overlayCols > remainInRow {
				overlayCols = remainInRow
			}
			
			if overlayCols > maxCols {
				maxCols = overlayCols
			}
		}
	}
	
	if len(startingTasks) == 0 {
		return nil
	}
	
	// Build content for all starting tasks
	content := d.buildMultiTaskOverlayContent(startingTasks)
	
	return &overlayInfo{
		content: content,
		cols:    maxCols,
	}
}

// buildTaskOverlayContent creates the LaTeX content for a task overlay
func (d Day) buildTaskOverlayContent(task *SpanningTask) string {
	nameText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Name))
	descText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Description))

	// Add star for milestone tasks
	if d.isMilestoneSpanningTask(task) {
		nameText = "★ " + nameText
	}

	textBody := `{\hyphenpenalty=10000\exhyphenpenalty=10000\emergencystretch=2em\setstretch{0.75}` +
		`{\centering\color{black}\textbf{\scriptsize ` + nameText + `}}`
	if descText != "" {
		textBody += `\\[-0.3ex]{\color{black}\tiny ` + descText + `}`
	}
	textBody += `}`

	return `\vspace*{0.1ex}{\begingroup\setlength{\fboxsep}{0pt}` +
		`\begin{tcolorbox}[enhanced, boxrule=0pt, arc=0pt, drop shadow,` +
		` left=1.5mm, right=1.5mm, top=0.5mm, bottom=0.5mm,` +
		` colback=` + task.Color + `!26,` +
		` interior style={left color=` + task.Color + `!34, right color=` + task.Color + `!6},` +
		` borderline west={1.4pt}{0pt}{` + task.Color + `!60!black},` +
		` borderline east={1.0pt}{0pt}{` + task.Color + `!45}]` +
		textBody +
		`\end{tcolorbox}\endgroup}`
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

// buildMultiTaskOverlayContent creates compact stacked content for multiple tasks
func (d Day) buildMultiTaskOverlayContent(tasks []*SpanningTask) string {
	if len(tasks) == 0 {
		return ""
	}
	
	if len(tasks) == 1 {
		return d.buildTaskOverlayContent(tasks[0])
	}
	
	// Sort tasks by category priority for better visual organization
	sortedTasks := d.sortTasksByPriority(tasks)
	
	var contentParts []string
	
	// Show up to 3 tasks in compact format, with smart truncation
	maxTasks := 3
	for i := 0; i < maxTasks && i < len(sortedTasks); i++ {
		task := sortedTasks[i]
		compactContent := d.buildCompactTaskOverlay(task, i, len(sortedTasks))
		contentParts = append(contentParts, compactContent)
	}
	
	// Add indicator if there are more tasks
	if len(sortedTasks) > maxTasks {
		moreCount := len(sortedTasks) - maxTasks
		indicator := `\vspace*{0.02ex}{\begingroup\setlength{\fboxsep}{0pt}` +
			`\begin{tcolorbox}[enhanced, boxrule=0pt, arc=0pt,` +
			` left=0.5mm, right=0.5mm, top=0.1mm, bottom=0.1mm,` +
			` colback=gray!15, height=0.5ex,` +
			` borderline west={0.5pt}{0pt}{gray!40}]` +
			`{\centering\color{gray}\textbf{\tiny +` + strconv.Itoa(moreCount) + ` more}}` +
			`\end{tcolorbox}\endgroup}`
		contentParts = append(contentParts, indicator)
	}
	
	return strings.Join(contentParts, "")
}

// buildCompactTaskOverlay creates a compact task overlay for multiple tasks
func (d Day) buildCompactTaskOverlay(task *SpanningTask, index, total int) string {
	nameText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Name))
	
	// Add star for milestone tasks
	if d.isMilestoneSpanningTask(task) {
		nameText = "★ " + nameText
	}
	
	// Smart truncation for multiple tasks
	maxChars := 18
	if total > 2 {
		maxChars = 15
	}
	if total > 3 {
		maxChars = 12
	}
	
	if len(nameText) > maxChars {
		nameText = d.smartTruncateText(nameText, maxChars)
	}
	
	// Adjust spacing and size based on position
	spacing := "0.05ex"
	boxHeight := "0.9ex"
	if index == 0 {
		spacing = "0.1ex"
		boxHeight = "1.0ex"
	}
	
	textBody := `{\hyphenpenalty=10000\exhyphenpenalty=10000\emergencystretch=2em\setstretch{0.7}` +
		`{\centering\color{black}\textbf{\tiny ` + nameText + `}}}`
	
	return `\vspace*{` + spacing + `}{\begingroup\setlength{\fboxsep}{0pt}` +
		`\begin{tcolorbox}[enhanced, boxrule=0pt, arc=0pt,` +
		` left=1.0mm, right=1.0mm, top=0.2mm, bottom=0.2mm,` +
		` height=` + boxHeight + `,` +
		` colback=` + task.Color + `!20,` +
		` interior style={left color=` + task.Color + `!28, right color=` + task.Color + `!8},` +
		` borderline west={1.0pt}{0pt}{` + task.Color + `!50!black}]` +
		textBody +
		`\end{tcolorbox}\endgroup}`
}

// sortTasksByPriority sorts tasks by category priority for better visual organization
func (d Day) sortTasksByPriority(tasks []*SpanningTask) []*SpanningTask {
	sorted := make([]*SpanningTask, len(tasks))
	copy(sorted, tasks)
	
	// Define priority order for categories
	priorityOrder := map[string]int{
		"DISSERTATION": 1,
		"PROPOSAL":     2,
		"PUBLICATION":  3,
		"RESEARCH":     4,
		"IMAGING":      5,
		"LASER":        6,
		"ADMIN":        7,
	}
	
	// Simple bubble sort by priority
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			priority1 := priorityOrder[sorted[j].Category]
			priority2 := priorityOrder[sorted[j+1].Category]
			if priority1 == 0 {
				priority1 = 99 // Unknown categories go last
			}
			if priority2 == 0 {
				priority2 = 99
			}
			if priority1 > priority2 {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}
	
	return sorted
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
