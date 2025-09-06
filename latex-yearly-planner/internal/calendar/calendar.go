package calendar

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"latex-yearly-planner/internal/data"
)

// Calendar represents the main calendar system
type Calendar struct {
	Year   *Year
	Tasks  []*data.Task
	Config *CalendarConfig
}

// CalendarConfig holds configuration for the calendar
type CalendarConfig struct {
	WeekStart           time.Weekday
	ShowTaskNames       bool
	ShowTaskDurations   bool
	ShowTaskPriorities  bool
	MaxTasksPerDay      int
	TaskSpacing         float64
	ColorScheme         string
}

// NewCalendar creates a new calendar instance
func NewCalendar(year int, weekStart time.Weekday) *Calendar {
	return &Calendar{
		Year: NewYear(weekStart, year),
		Config: &CalendarConfig{
			WeekStart:          weekStart,
			ShowTaskNames:      true,
			ShowTaskDurations:  true,
			ShowTaskPriorities: true,
			MaxTasksPerDay:     5,
			TaskSpacing:        1.0,
			ColorScheme:        "default",
		},
	}
}

// AddTasks adds tasks to the calendar
func (c *Calendar) AddTasks(tasks []*data.Task) {
	c.Tasks = append(c.Tasks, tasks...)
	c.assignTasksToCalendar()
}

// assignTasksToCalendar assigns tasks to appropriate days
func (c *Calendar) assignTasksToCalendar() {
	for _, task := range c.Tasks {
		c.assignTaskToDays(task)
	}
}

// assignTaskToDays assigns a single task to the appropriate days
func (c *Calendar) assignTaskToDays(task *data.Task) {
	// Find all months that this task spans
	for _, quarter := range c.Year.Quarters {
		for _, month := range quarter.Months {
			if c.taskOverlapsMonth(task, month) {
				c.assignTaskToMonth(task, month)
			}
		}
	}
}

// taskOverlapsMonth checks if a task overlaps with a month
func (c *Calendar) taskOverlapsMonth(task *data.Task, month *Month) bool {
	monthStart := time.Date(month.Year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
	monthEnd := monthStart.AddDate(0, 1, -1)
	
	return task.StartDate.Before(monthEnd.AddDate(0, 0, 1)) && 
		   task.EndDate.After(monthStart.AddDate(0, 0, -1))
}

// assignTaskToMonth assigns a task to a specific month
func (c *Calendar) assignTaskToMonth(task *data.Task, month *Month) {
	// Create spanning task
	spanningTask := CreateSpanningTask(*task, task.StartDate, task.EndDate)
	
	// Apply to month
	ApplySpanningTasksToMonth(month, []SpanningTask{spanningTask})
}

// GetTasksForMonth returns tasks for a specific month
func (c *Calendar) GetTasksForMonth(month time.Month) []*data.Task {
	var monthTasks []*data.Task
	monthStart := time.Date(c.Year.Number, month, 1, 0, 0, 0, 0, time.Local)
	monthEnd := monthStart.AddDate(0, 1, -1)
	
	for _, task := range c.Tasks {
		if task.StartDate.Before(monthEnd.AddDate(0, 0, 1)) && 
		   task.EndDate.After(monthStart.AddDate(0, 0, -1)) {
			monthTasks = append(monthTasks, task)
		}
	}
	
	return monthTasks
}

// GetTasksForDay returns tasks for a specific day
func (c *Calendar) GetTasksForDay(day time.Time) []*data.Task {
	var dayTasks []*data.Task
	
	for _, task := range c.Tasks {
		if task.StartDate.Before(day.AddDate(0, 0, 1)) && 
		   task.EndDate.After(day.AddDate(0, 0, -1)) {
			dayTasks = append(dayTasks, task)
		}
	}
	
	return dayTasks
}

// Year represents a calendar year
type Year struct {
	Number   int
	Quarters Quarters
	Weeks    Weeks
}

// NewYear creates a new year
func NewYear(wd time.Weekday, year int) *Year {
	y := &Year{Number: year}
	y.Quarters = make(Quarters, 4)
	y.Weeks = NewWeeksForYear(wd, y)
	
	for i := 0; i < 4; i++ {
		y.Quarters[i] = NewQuarter(wd, y, i+1)
	}
	
	return y
}

// Breadcrumb returns the year breadcrumb
func (y Year) Breadcrumb() string {
	return fmt.Sprintf("%d", y.Number)
}

// SideQuarters returns side quarter items
func (y Year) SideQuarters(sel ...int) []header.CellItem {
	var items []header.CellItem
	for _, q := range y.Quarters {
		if len(sel) == 0 || contains(sel, q.Number) {
			items = append(items, header.CellItem{
				Text: q.Name(),
				Ref:  q.Breadcrumb(),
			})
		}
	}
	return items
}

// SideMonths returns side month items
func (y Year) SideMonths(sel ...time.Month) []header.CellItem {
	var items []header.CellItem
	for _, q := range y.Quarters {
		for _, m := range q.Months {
			if len(sel) == 0 || containsMonth(sel, m.Month) {
				items = append(items, header.CellItem{
					Text: m.ShortName(),
					Ref:  m.Breadcrumb(),
				})
			}
		}
	}
	return items
}

// HeadingMOS returns the year heading
func (y Year) HeadingMOS() string {
	return fmt.Sprintf("Year %d", y.Number)
}

// Quarter represents a calendar quarter
type Quarter struct {
	Year   *Year
	Number int
	Months Months
}

// NewQuarter creates a new quarter
func NewQuarter(wd time.Weekday, year *Year, qrtr int) *Quarter {
	q := &Quarter{
		Year:   year,
		Number: qrtr,
		Months: make(Months, 3),
	}
	
	startMonth := time.Month((qrtr-1)*3 + 1)
	for i := 0; i < 3; i++ {
		month := startMonth + time.Month(i)
		q.Months[i] = NewMonth(wd, year, q, month)
	}
	
	return q
}

// Breadcrumb returns the quarter breadcrumb
func (q *Quarter) Breadcrumb() string {
	return fmt.Sprintf("%s %d", q.Name(), q.Year.Number)
}

// Name returns the quarter name
func (q *Quarter) Name() string { 
	return "Q" + strconv.Itoa(q.Number) 
}

// HeadingMOS returns the quarter heading
func (q *Quarter) HeadingMOS() string {
	return fmt.Sprintf("Quarter %d, %d", q.Number, q.Year.Number)
}

// Month represents a calendar month
type Month struct {
	Year    *Year
	Quarter *Quarter
	Month   time.Month
	Weekday time.Weekday
	Weeks   Weeks
}

// NewMonth creates a new month
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

// MaybeName returns the month name if large
func (m *Month) MaybeName(large interface{}) string {
	if isLarge, ok := large.(bool); ok && isLarge {
		return m.Month.String()
	}
	return ""
}

// WeekHeader returns the week header
func (m *Month) WeekHeader(large interface{}) string {
	if isLarge, ok := large.(bool); ok && isLarge {
		return "\\textbf{Week}"
	}
	return ""
}

// DefineTable defines the table structure
func (m *Month) DefineTable(typ interface{}, large interface{}) string {
	tableType := "tabular"
	if t, ok := typ.(string); ok && t != "" {
		tableType = t
	}
	
	if isLarge, ok := large.(bool); ok && isLarge {
		return fmt.Sprintf("\\begin{%s}{|c|*{7}{c|}}", tableType)
	}
	return fmt.Sprintf("\\begin{%s}{|c|*{7}{c|}}", tableType)
}

// EndTable ends the table
func (m *Month) EndTable(typ interface{}) string {
	tableType := "tabular"
	if t, ok := typ.(string); ok && t != "" {
		tableType = t
	}
	return fmt.Sprintf("\\end{%s}", tableType)
}

// Breadcrumb returns the month breadcrumb
func (m *Month) Breadcrumb() string {
	return fmt.Sprintf("%s %d", m.Month.String(), m.Year.Number)
}

// PrevNext returns previous/next navigation
func (m *Month) PrevNext() header.Items {
	items := header.Items{}
	
	// Previous month
	if m.Month > time.January || m.Year.Number > 1900 {
		prevMonth := m.Month - 1
		prevYear := m.Year.Number
		if prevMonth == 0 {
			prevMonth = time.December
			prevYear--
		}
		items = append(items, header.Item{
			Text: "← " + prevMonth.String()[:3],
			Ref:  fmt.Sprintf("%s %d", prevMonth.String(), prevYear),
		})
	}
	
	// Next month
	if m.Month < time.December || m.Year.Number < 2100 {
		nextMonth := m.Month + 1
		nextYear := m.Year.Number
		if nextMonth == 13 {
			nextMonth = time.January
			nextYear++
		}
		items = append(items, header.Item{
			Text: nextMonth.String()[:3] + " →",
			Ref:  fmt.Sprintf("%s %d", nextMonth.String(), nextYear),
		})
	}
	
	return items
}

// ShortName returns the short month name
func (m *Month) ShortName() string { 
	return m.Month.String()[:3] 
}

// HeadingMOS returns the month heading
func (m *Month) HeadingMOS() string {
	return fmt.Sprintf("%s %d", m.Month.String(), m.Year.Number)
}

// Week represents a calendar week
type Week struct {
	Days [7]Day
	Weekday  time.Weekday
	Year     *Year
	Months   Months
	Quarters Quarters
}

// NewWeeksForMonth creates weeks for a month
func NewWeeksForMonth(wd time.Weekday, year *Year, qrtr *Quarter, month *Month) Weeks {
	weeks := make(Weeks, 0, 6)
	
	// Find the first day of the month
	firstDay := time.Date(year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
	
	// Find the first day of the week containing the first day of the month
	daysSinceWeekStart := int(firstDay.Weekday() - wd)
	if daysSinceWeekStart < 0 {
		daysSinceWeekStart += 7
	}
	weekStart := firstDay.AddDate(0, 0, -daysSinceWeekStart)
	
	// Generate weeks
	currentWeekStart := weekStart
	for currentWeekStart.Month() == month.Month || 
		currentWeekStart.AddDate(0, 0, 6).Month() == month.Month {
		
		week := &Week{
			Weekday: wd,
			Year:    year,
		}
		
		// Fill days
		for i := 0; i < 7; i++ {
			dayTime := currentWeekStart.AddDate(0, 0, i)
			week.Days[i] = Day{Time: dayTime}
		}
		
		weeks = append(weeks, week)
		currentWeekStart = currentWeekStart.AddDate(0, 0, 7)
		
		// Prevent infinite loop
		if len(weeks) > 6 {
			break
		}
	}
	
	return weeks
}

// NewWeeksForYear creates weeks for a year
func NewWeeksForYear(wd time.Weekday, year *Year) Weeks {
	weeks := make(Weeks, 0, 53)
	
	// Start from January 1st
	start := time.Date(year.Number, time.January, 1, 0, 0, 0, 0, time.Local)
	
	// Find the first day of the week containing January 1st
	daysSinceWeekStart := int(start.Weekday() - wd)
	if daysSinceWeekStart < 0 {
		daysSinceWeekStart += 7
	}
	weekStart := start.AddDate(0, 0, -daysSinceWeekStart)
	
	// Generate weeks for the entire year
	currentWeekStart := weekStart
	for currentWeekStart.Year() == year.Number || 
		currentWeekStart.AddDate(0, 0, 6).Year() == year.Number {
		
		week := &Week{
			Weekday: wd,
			Year:    year,
		}
		
		// Fill days
		for i := 0; i < 7; i++ {
			dayTime := currentWeekStart.AddDate(0, 0, i)
			week.Days[i] = Day{Time: dayTime}
		}
		
		weeks = append(weeks, week)
		currentWeekStart = currentWeekStart.AddDate(0, 0, 7)
		
		// Prevent infinite loop
		if len(weeks) > 53 {
			break
		}
	}
	
	return weeks
}

// WeekNumber returns the week number
func (w *Week) WeekNumber(large interface{}) string {
	if isLarge, ok := large.(bool); ok && isLarge {
		return fmt.Sprintf("\\textbf{%d}", w.weekNumber())
	}
	return fmt.Sprintf("%d", w.weekNumber())
}

// weekNumber calculates the week number
func (w *Week) weekNumber() int {
	// Simple week number calculation
	_, week := w.Days[0].Time.ISOWeek()
	return week
}

// Breadcrumb returns the week breadcrumb
func (w *Week) Breadcrumb() string {
	return fmt.Sprintf("Week %d", w.weekNumber())
}

// monthOverlap checks if the week spans multiple months
func (w *Week) monthOverlap() bool { 
	return w.Days[0].Time.Month() != w.Days[6].Time.Month() 
}

// quarterOverlap checks if the week spans multiple quarters
func (w *Week) quarterOverlap() bool { 
	return w.leftQuarter() != w.rightQuarter() 
}

// leftQuarter returns the left quarter
func (w *Week) leftQuarter() int { 
	return int(math.Ceil(float64(w.Days[0].Time.Month()) / 3.)) 
}

// rightQuarter returns the right quarter
func (w *Week) rightQuarter() int { 
	return int(math.Ceil(float64(w.Days[6].Time.Month()) / 3.)) 
}

// rightMonth returns the right month
func (w *Week) rightMonth() time.Month {
	return w.Days[6].Time.Month()
}

// PrevNext returns previous/next navigation
func (w *Week) PrevNext() header.Items {
	items := header.Items{}
	
	// Previous week
	if w.Days[0].Time.Year() > 1900 {
		prevWeek := w.Days[0].Time.AddDate(0, 0, -7)
		items = append(items, header.Item{
			Text: "← Prev",
			Ref:  fmt.Sprintf("Week %d", prevWeek.ISOWeek()),
		})
	}
	
	// Next week
	if w.Days[6].Time.Year() < 2100 {
		nextWeek := w.Days[6].Time.AddDate(0, 0, 7)
		items = append(items, header.Item{
			Text: "Next →",
			Ref:  fmt.Sprintf("Week %d", nextWeek.ISOWeek()),
		})
	}
	
	return items
}

// NextExists checks if next week exists
func (w *Week) NextExists() bool {
	return w.Days[6].Time.Year() < 2100
}

// PrevExists checks if previous week exists
func (w *Week) PrevExists() bool {
	return w.Days[0].Time.Year() > 1900
}

// Next returns the next week
func (w *Week) Next() *Week { 
	return fillWeekly(w.Weekday, w.Year, w.Days[0].Add(7)) 
}

// Prev returns the previous week
func (w *Week) Prev() *Week { 
	return fillWeekly(w.Weekday, w.Year, w.Days[0].Add(-7)) 
}

// fillWeekly fills a week with days
func fillWeekly(wd time.Weekday, year *Year, ptr Day) *Week {
	week := &Week{
		Weekday: wd,
		Year:    year,
	}
	
	for i := 0; i < 7; i++ {
		dayTime := ptr.Time.AddDate(0, 0, i)
		week.Days[i] = Day{Time: dayTime}
	}
	
	return week
}

// QuartersBreadcrumb returns quarters breadcrumb
func (w *Week) QuartersBreadcrumb() header.ItemsGroup {
	items := make(header.ItemsGroup, 4)
	for i := 0; i < 4; i++ {
		items[i] = header.Items{
			header.Item{Text: fmt.Sprintf("Q%d", i+1), Ref: fmt.Sprintf("Q%d %d", i+1, w.Year.Number)},
		}
	}
	return items
}

// MonthsBreadcrumb returns months breadcrumb
func (w *Week) MonthsBreadcrumb() header.ItemsGroup {
	items := make(header.ItemsGroup, 12)
	months := []time.Month{
		time.January, time.February, time.March, time.April,
		time.May, time.June, time.July, time.August,
		time.September, time.October, time.November, time.December,
	}
	
	for i, month := range months {
		items[i] = header.Items{
			header.Item{Text: month.String()[:3], Ref: fmt.Sprintf("%s %d", month.String(), w.Year.Number)},
		}
	}
	return items
}

// ref returns the week reference
func (w *Week) ref() string {
	return fmt.Sprintf("week-%d-%d", w.Year.Number, w.weekNumber())
}

// leftMonth returns the left month
func (w *Week) leftMonth() time.Month {
	return w.Days[0].Time.Month()
}

// rightYear returns the right year
func (w *Week) rightYear() int {
	return w.Days[6].Time.Year()
}

// HeadingMOS returns the week heading
func (w *Week) HeadingMOS() string {
	return fmt.Sprintf("Week %d, %d", w.weekNumber(), w.Year.Number)
}

// Name returns the week name
func (w *Week) Name() string { 
	return "Week " + strconv.Itoa(w.weekNumber()) 
}

// Target returns the week target
func (w *Week) Target() string { 
	return latex.Hypertarget(w.ref(), w.Name()) 
}

// HasDays checks if the week has days
func (w *Week) HasDays() bool {
	return len(w.Days) > 0 && !w.Days[0].Time.IsZero()
}

// Day represents a calendar day
type Day struct {
	Time          time.Time
	Tasks         []Task
	SpanningTasks []*SpanningTask
}

// Day methods
func (d Day) Day(today, large interface{}) string {
	dayStr := strconv.Itoa(d.Time.Day())
	
	if isLarge, ok := large.(bool); ok && isLarge {
		return d.renderLargeDay(dayStr)
	}
	
	return dayStr
}

func (d Day) renderLargeDay(day string) string {
	// Check if today
	if d.Time.Year() == time.Now().Year() && 
	   d.Time.Month() == time.Now().Month() && 
	   d.Time.Day() == time.Now().Day() {
		return fmt.Sprintf("\\textbf{\\textcolor{red}{%s}}", day)
	}
	
	return fmt.Sprintf("\\textbf{%s}", day)
}

func (d Day) ref(prefix ...string) string {
	pre := ""
	if len(prefix) > 0 {
		pre = prefix[0] + "-"
	}
	return fmt.Sprintf("%s%d-%02d-%02d", pre, d.Time.Year(), d.Time.Month(), d.Time.Day())
}

func (d Day) Add(days int) Day {
	return Day{Time: d.Time.AddDate(0, 0, days)}
}

func (d Day) WeekLink() string {
	_, week := d.Time.ISOWeek()
	return fmt.Sprintf("Week %d", week)
}

func (d Day) Breadcrumb(prefix string, leaf string, shorten bool) string {
	year := d.Time.Year()
	month := d.Time.Month()
	day := d.Time.Day()
	
	if shorten {
		return fmt.Sprintf("%s %d, %d", month.String()[:3], day, year)
	}
	return fmt.Sprintf("%s %d, %d", month.String(), day, year)
}

func (d Day) LinkLeaf(prefix, leaf string) string {
	return fmt.Sprintf("%s-%s", prefix, leaf)
}

func (d Day) PrevNext(prefix string) header.Items {
	items := header.Items{}
	
	// Previous day
	if d.Time.Year() > 1900 {
		prevDay := d.Time.AddDate(0, 0, -1)
		items = append(items, header.Item{
			Text: "← " + strconv.Itoa(prevDay.Day()),
			Ref:  fmt.Sprintf("%s-%d-%02d-%02d", prefix, prevDay.Year(), prevDay.Month(), prevDay.Day()),
		})
	}
	
	// Next day
	if d.Time.Year() < 2100 {
		nextDay := d.Time.AddDate(0, 0, 1)
		items = append(items, header.Item{
			Text: strconv.Itoa(nextDay.Day()) + " →",
			Ref:  fmt.Sprintf("%s-%d-%02d-%02d", prefix, nextDay.Year(), nextDay.Month(), nextDay.Day()),
		})
	}
	
	return items
}

func (d Day) Next() Day { return d.Add(1) }
func (d Day) Prev() Day { return d.Add(-1) }

func (d Day) NextExists() bool { 
	return d.Time.Month() < time.December || d.Time.Day() < 31 
}

func (d Day) PrevExists() bool { 
	return d.Time.Month() > time.January || d.Time.Day() > 1 
}

func (d Day) Quarter() int { 
	return int(math.Ceil(float64(d.Time.Month()) / 3.)) 
}

func (d Day) Month() time.Month { 
	return d.Time.Month() 
}

func (d Day) HeadingMOS(prefix, leaf string) string {
	year := d.Time.Year()
	month := d.Time.Month()
	day := d.Time.Day()
	
	return fmt.Sprintf("%s %d, %d", month.String(), day, year)
}

// Task represents a simple task
type Task struct {
	ID          string
	Name        string
	Description string
	Category    string
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
	Progress    int
	Status      string
	Assignee    string
}

// CreateSpanningTask creates a spanning task from a data task
func CreateSpanningTask(task data.Task, startDate, endDate time.Time) SpanningTask {
	return SpanningTask{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Category:    task.Category,
		StartDate:   startDate,
		EndDate:     endDate,
		Color:       getColorForCategory(task.Category),
		Priority:    task.Priority,
		Progress:    0,
		Status:      task.Status,
		Assignee:    task.Assignee,
	}
}

// ApplySpanningTasksToMonth applies spanning tasks to a month
func ApplySpanningTasksToMonth(month *Month, tasks []SpanningTask) {
	// This is a simplified implementation
	// In a full implementation, you would assign tasks to specific days
}

// getColorForCategory returns a color for a task category
func getColorForCategory(category string) string {
	colors := map[string]string{
		"work":     "blue",
		"personal": "green",
		"urgent":   "red",
		"meeting":  "orange",
		"deadline": "purple",
	}
	
	if color, exists := colors[strings.ToLower(category)]; exists {
		return color
	}
	return "gray"
}

// Helper functions
func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsMonth(slice []time.Month, item time.Month) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Type definitions for collections
type Years []*Year
type Quarters []*Quarter
type Months []*Month
type Weeks []*Week
type Days []*Day

// Collection methods
func (q Quarters) Numbers() []int {
	numbers := make([]int, len(q))
	for i, quarter := range q {
		numbers[i] = quarter.Number
	}
	return numbers
}

func (m Months) Months() []time.Month {
	months := make([]time.Month, len(m))
	for i, month := range m {
		months[i] = month.Month
	}
	return months
}

// Additional helper functions for template compatibility
func selectStartWeek(year int, weekStart time.Weekday) Day {
	// Find the first day of the year
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	
	// Find the first day of the week containing January 1st
	daysSinceWeekStart := int(firstDay.Weekday() - weekStart)
	if daysSinceWeekStart < 0 {
		daysSinceWeekStart += 7
	}
	
	weekStartDay := firstDay.AddDate(0, 0, -daysSinceWeekStart)
	return Day{Time: weekStartDay}
}
