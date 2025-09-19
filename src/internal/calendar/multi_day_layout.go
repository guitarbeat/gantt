package calendar

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"phd-dissertation-planner/internal/data"
)

// MultiDayLayoutEngine handles the layout and rendering of multi-day task bars
type MultiDayLayoutEngine struct {
	calendarStart    time.Time
	calendarEnd      time.Time
	dayWidth         float64
	dayHeight        float64
	rowHeight        float64
	maxRowsPerDay    int
	overlapThreshold float64
}

// TaskBar represents a rendered task bar with positioning information
type TaskBar struct {
	TaskID        string
	StartDate     time.Time
	EndDate       time.Time
	StartX        float64
	EndX          float64
	Y             float64
	Width         float64
	Height        float64
	Row           int
	Color         string
	BorderColor   string
	Opacity       float64
	ZIndex        int
	IsContinuation bool
	IsStart       bool
	IsEnd         bool
	MonthBoundary bool
}

// TaskGroup represents a group of overlapping tasks
type TaskGroup struct {
	Tasks     []*data.Task
	StartDate time.Time
	EndDate   time.Time
	Rows      int
}

// NewMultiDayLayoutEngine creates a new multi-day layout engine
func NewMultiDayLayoutEngine(calendarStart, calendarEnd time.Time, dayWidth, dayHeight float64) *MultiDayLayoutEngine {
	return &MultiDayLayoutEngine{
		calendarStart:    calendarStart,
		calendarEnd:      calendarEnd,
		dayWidth:         dayWidth,
		dayHeight:        dayHeight,
		rowHeight:        dayHeight * 0.8, // 80% of day height for task bars
		maxRowsPerDay:    4,               // Maximum 4 task rows per day
		overlapThreshold: 0.1,             // 10% overlap threshold
	}
}

// LayoutMultiDayTasks performs the two-step algorithm for multi-day task layout
func (mle *MultiDayLayoutEngine) LayoutMultiDayTasks(tasks []*data.Task) []*TaskBar {
	// Step 1: Group overlapping tasks
	groups := mle.groupOverlappingTasks(tasks)
	
	// Step 2: Layout calculation within groups
	var taskBars []*TaskBar
	for _, group := range groups {
		groupBars := mle.layoutTaskGroup(group)
		taskBars = append(taskBars, groupBars...)
	}
	
	return taskBars
}

// groupOverlappingTasks implements Step 1: Grouping Overlapping Events
func (mle *MultiDayLayoutEngine) groupOverlappingTasks(tasks []*data.Task) []*TaskGroup {
	// Sort tasks by start date and duration
	sortedTasks := make([]*data.Task, len(tasks))
	copy(sortedTasks, tasks)
	sort.Slice(sortedTasks, func(i, j int) bool {
		if sortedTasks[i].StartDate.Equal(sortedTasks[j].StartDate) {
			// If start dates are equal, sort by duration (longer first)
			return sortedTasks[i].GetDuration() > sortedTasks[j].GetDuration()
		}
		return sortedTasks[i].StartDate.Before(sortedTasks[j].StartDate)
	})
	
	var groups []*TaskGroup
	used := make(map[string]bool)
	
	for _, task := range sortedTasks {
		if used[task.ID] {
			continue
		}
		
		// Create a new group starting with this task
		group := &TaskGroup{
			Tasks:     []*data.Task{task},
			StartDate: task.StartDate,
			EndDate:   task.EndDate,
		}
		used[task.ID] = true
		
		// Find all tasks that overlap with this group
		for _, otherTask := range sortedTasks {
			if used[otherTask.ID] {
				continue
			}
			
			if mle.tasksOverlap(group, otherTask) {
				group.Tasks = append(group.Tasks, otherTask)
				used[otherTask.ID] = true
				
				// Update group date range
				if otherTask.StartDate.Before(group.StartDate) {
					group.StartDate = otherTask.StartDate
				}
				if otherTask.EndDate.After(group.EndDate) {
					group.EndDate = otherTask.EndDate
				}
			}
		}
		
		// Calculate number of rows needed for this group
		group.Rows = mle.calculateGroupRows(group)
		groups = append(groups, group)
	}
	
	return groups
}

// tasksOverlap checks if a task overlaps with a group
func (mle *MultiDayLayoutEngine) tasksOverlap(group *TaskGroup, task *data.Task) bool {
	// Check if task overlaps with any task in the group
	for _, groupTask := range group.Tasks {
		if mle.tasksOverlapDirect(groupTask, task) {
			return true
		}
	}
	return false
}

// tasksOverlapDirect checks if two tasks overlap directly
func (mle *MultiDayLayoutEngine) tasksOverlapDirect(task1, task2 *data.Task) bool {
	// Tasks overlap if one starts before the other ends
	return !task1.StartDate.After(task2.EndDate) && !task2.StartDate.After(task1.EndDate)
}

// calculateGroupRows calculates the number of rows needed for a group
func (mle *MultiDayLayoutEngine) calculateGroupRows(group *TaskGroup) int {
	// Use greedy algorithm to determine minimum rows needed
	rows := 1
	rowEndTimes := make([]time.Time, 0)
	
	// Sort tasks within group by start date
	sort.Slice(group.Tasks, func(i, j int) bool {
		return group.Tasks[i].StartDate.Before(group.Tasks[j].StartDate)
	})
	
	for _, task := range group.Tasks {
		// Find first available row
		rowFound := false
		for i, endTime := range rowEndTimes {
			if task.StartDate.After(endTime) || task.StartDate.Equal(endTime) {
				// This row is available
				rowEndTimes[i] = task.EndDate
				rowFound = true
				break
			}
		}
		
		if !rowFound {
			// Need a new row
			rowEndTimes = append(rowEndTimes, task.EndDate)
			rows++
		}
	}
	
	// Cap at maximum rows per day
	if rows > mle.maxRowsPerDay {
		rows = mle.maxRowsPerDay
	}
	
	return rows
}

// layoutTaskGroup implements Step 2: Layout Calculation within Groups
func (mle *MultiDayLayoutEngine) layoutTaskGroup(group *TaskGroup) []*TaskBar {
	var taskBars []*TaskBar
	rowEndTimes := make([]time.Time, group.Rows)
	
	// Sort tasks by start date and priority
	sort.Slice(group.Tasks, func(i, j int) bool {
		if group.Tasks[i].StartDate.Equal(group.Tasks[j].StartDate) {
			// If start dates are equal, sort by priority (higher first)
			return group.Tasks[i].Priority > group.Tasks[j].Priority
		}
		return group.Tasks[i].StartDate.Before(group.Tasks[j].StartDate)
	})
	
	for _, task := range group.Tasks {
		// Find first available row
		row := mle.findAvailableRow(task, rowEndTimes)
		
		// Create task bar
		taskBar := mle.createTaskBar(task, row, group.Rows)
		taskBars = append(taskBars, taskBar)
		
		// Update row end time
		rowEndTimes[row] = task.EndDate
	}
	
	return taskBars
}

// findAvailableRow finds the first available row for a task
func (mle *MultiDayLayoutEngine) findAvailableRow(task *data.Task, rowEndTimes []time.Time) int {
	for i, endTime := range rowEndTimes {
		if task.StartDate.After(endTime) || task.StartDate.Equal(endTime) {
			return i
		}
	}
	
	// If no row is available, use the first row (overlap will be handled visually)
	return 0
}

// createTaskBar creates a task bar with positioning information
func (mle *MultiDayLayoutEngine) createTaskBar(task *data.Task, row, totalRows int) *TaskBar {
	// Calculate X coordinates based on start and end dates
	startX := mle.calculateXPosition(task.StartDate)
	endX := mle.calculateXPosition(task.EndDate)
	
	// Calculate Y position based on row
	y := mle.calculateYPosition(row, totalRows)
	
	// Calculate width
	width := endX - startX
	
	// Get task color from category
	category := data.GetCategory(task.Category)
	
	// Determine if this is a continuation, start, or end
	isContinuation := mle.isTaskContinuation(task)
	isStart := mle.isTaskStart(task)
	isEnd := mle.isTaskEnd(task)
	monthBoundary := mle.hasMonthBoundary(task)
	
	return &TaskBar{
		TaskID:         task.ID,
		StartDate:      task.StartDate,
		EndDate:        task.EndDate,
		StartX:         startX,
		EndX:           endX,
		Y:              y,
		Width:          width,
		Height:         mle.rowHeight,
		Row:            row,
		Color:          category.Color,
		BorderColor:    "#000000",
		Opacity:        1.0,
		ZIndex:         category.Priority,
		IsContinuation: isContinuation,
		IsStart:        isStart,
		IsEnd:          isEnd,
		MonthBoundary:  monthBoundary,
	}
}

// calculateXPosition calculates the X position for a given date
func (mle *MultiDayLayoutEngine) calculateXPosition(date time.Time) float64 {
	// Calculate days from calendar start
	daysFromStart := int(date.Sub(mle.calendarStart).Hours() / 24)
	
	// Calculate X position (day width * days from start)
	return float64(daysFromStart) * mle.dayWidth
}

// calculateYPosition calculates the Y position for a given row
func (mle *MultiDayLayoutEngine) calculateYPosition(row, totalRows int) float64 {
	// Distribute rows evenly within the day height
	rowSpacing := mle.dayHeight / float64(totalRows+1)
	return float64(row+1) * rowSpacing
}

// isTaskContinuation checks if this task is a continuation from previous month
func (mle *MultiDayLayoutEngine) isTaskContinuation(task *data.Task) bool {
	// Check if task started before calendar start
	return task.StartDate.Before(mle.calendarStart)
}

// isTaskStart checks if this is the start of a multi-day task
func (mle *MultiDayLayoutEngine) isTaskStart(task *data.Task) bool {
	// Check if task starts on or after calendar start
	return !task.StartDate.Before(mle.calendarStart)
}

// isTaskEnd checks if this is the end of a multi-day task
func (mle *MultiDayLayoutEngine) isTaskEnd(task *data.Task) bool {
	// Check if task ends on or before calendar end
	return !task.EndDate.After(mle.calendarEnd)
}

// hasMonthBoundary checks if task spans across month boundaries
func (mle *MultiDayLayoutEngine) hasMonthBoundary(task *data.Task) bool {
	startMonth := task.StartDate.Month()
	endMonth := task.EndDate.Month()
	return startMonth != endMonth
}

// HandleMonthBoundary handles task bars that span across month boundaries
func (mle *MultiDayLayoutEngine) HandleMonthBoundary(taskBars []*TaskBar) []*TaskBar {
	var processedBars []*TaskBar
	
	for _, bar := range taskBars {
		if !bar.MonthBoundary {
			processedBars = append(processedBars, bar)
			continue
		}
		
		// Split task bar at month boundaries
		splitBars := mle.splitTaskBarAtMonthBoundaries(bar)
		processedBars = append(processedBars, splitBars...)
	}
	
	return processedBars
}

// splitTaskBarAtMonthBoundaries splits a task bar at month boundaries
func (mle *MultiDayLayoutEngine) splitTaskBarAtMonthBoundaries(bar *TaskBar) []*TaskBar {
	var splitBars []*TaskBar
	
	// Find all month boundaries within the task duration
	current := bar.StartDate
	end := bar.EndDate
	
	for current.Before(end) {
		// Find the end of the current month
		monthEnd := time.Date(current.Year(), current.Month()+1, 0, 0, 0, 0, 0, current.Location())
		if monthEnd.After(end) {
			monthEnd = end
		}
		
		// Create a task bar segment
		segment := &TaskBar{
			TaskID:         bar.TaskID,
			StartDate:      current,
			EndDate:        monthEnd,
			StartX:         mle.calculateXPosition(current),
			EndX:           mle.calculateXPosition(monthEnd),
			Y:              bar.Y,
			Width:          mle.calculateXPosition(monthEnd) - mle.calculateXPosition(current),
			Height:         bar.Height,
			Row:            bar.Row,
			Color:          bar.Color,
			BorderColor:    bar.BorderColor,
			Opacity:        bar.Opacity,
			ZIndex:         bar.ZIndex,
			IsContinuation: current.Equal(bar.StartDate) && bar.IsContinuation || !current.Equal(bar.StartDate),
			IsStart:        current.Equal(bar.StartDate) && bar.IsStart,
			IsEnd:          monthEnd.Equal(bar.EndDate) && bar.IsEnd,
			MonthBoundary:  false, // Individual segments don't have month boundaries
		}
		
		splitBars = append(splitBars, segment)
		
		// Move to next month
		current = monthEnd.AddDate(0, 0, 1)
	}
	
	return splitBars
}

// GenerateLaTeX generates LaTeX code for multi-day task bars
func (mle *MultiDayLayoutEngine) GenerateLaTeX(taskBars []*TaskBar) string {
	var latex strings.Builder
	
	// Group task bars by day for efficient rendering
	dayGroups := mle.groupTaskBarsByDay(taskBars)
	
	for day, bars := range dayGroups {
		latex.WriteString(mle.generateDayLaTeX(day, bars))
	}
	
	return latex.String()
}

// groupTaskBarsByDay groups task bars by day
func (mle *MultiDayLayoutEngine) groupTaskBarsByDay(taskBars []*TaskBar) map[time.Time][]*TaskBar {
	dayGroups := make(map[time.Time][]*TaskBar)
	
	for _, bar := range taskBars {
		// Group by start date
		day := time.Date(bar.StartDate.Year(), bar.StartDate.Month(), bar.StartDate.Day(), 0, 0, 0, 0, bar.StartDate.Location())
		dayGroups[day] = append(dayGroups[day], bar)
	}
	
	return dayGroups
}

// generateDayLaTeX generates LaTeX code for a specific day
func (mle *MultiDayLayoutEngine) generateDayLaTeX(day time.Time, bars []*TaskBar) string {
	var latex strings.Builder
	
	// Sort bars by row and start time
	sort.Slice(bars, func(i, j int) bool {
		if bars[i].Row == bars[j].Row {
			return bars[i].StartX < bars[j].StartX
		}
		return bars[i].Row < bars[j].Row
	})
	
	// Generate LaTeX for each bar
	for _, bar := range bars {
		latex.WriteString(mle.generateTaskBarLaTeX(bar))
	}
	
	return latex.String()
}

// generateTaskBarLaTeX generates LaTeX code for a single task bar
func (mle *MultiDayLayoutEngine) generateTaskBarLaTeX(bar *TaskBar) string {
	// Convert colors to LaTeX format
	color := mle.convertColorToLaTeX(bar.Color)
	
	// Generate task bar LaTeX
	return fmt.Sprintf(`
		\\begin{tikzpicture}[overlay]
			\\node[anchor=north west, inner sep=0pt] at (%.2f,%.2f) {
				\\begin{tcolorbox}[enhanced, boxrule=0pt, arc=2pt, drop shadow,
					left=1.5mm, right=1.5mm, top=0.5mm, bottom=0.5mm,
					width=%.2fmm, height=%.2fmm,
					colback=%s!26,
					interior style={left color=%s!34, right color=%s!6},
					borderline west={1.4pt}{0pt}{%s!60!black},
					borderline east={1.0pt}{0pt}{%s!45}]
					{\\footnotesize %s}
				\\end{tcolorbox}
			};
		\\end{tikzpicture}
	`, bar.StartX, bar.Y, bar.Width, bar.Height, color, color, color, color, color, bar.TaskID)
}

// convertColorToLaTeX converts hex color to LaTeX color name
func (mle *MultiDayLayoutEngine) convertColorToLaTeX(hexColor string) string {
	// Map hex colors to LaTeX color names
	colorMap := map[string]string{
		"#4A90E2": "blue",      // PROPOSAL
		"#F5A623": "orange",    // LASER
		"#7ED321": "green",     // IMAGING
		"#BD10E0": "purple",    // ADMIN
		"#D0021B": "red",       // DISSERTATION
		"#50E3C2": "teal",      // RESEARCH
		"#B8E986": "lime",      // PUBLICATION
		"#CCCCCC": "gray",      // Default
	}
	
	if color, exists := colorMap[hexColor]; exists {
		return color
	}
	return "gray" // Default fallback
}

// ValidateLayout validates the layout for potential issues
func (mle *MultiDayLayoutEngine) ValidateLayout(taskBars []*TaskBar) []string {
	var issues []string
	
	// Check for overlapping task bars in the same row
	rowBars := make(map[int][]*TaskBar)
	for _, bar := range taskBars {
		rowBars[bar.Row] = append(rowBars[bar.Row], bar)
	}
	
	for row, bars := range rowBars {
		for i := 0; i < len(bars); i++ {
			for j := i + 1; j < len(bars); j++ {
				if mle.barsOverlap(bars[i], bars[j]) {
					issues = append(issues, fmt.Sprintf("Task bars overlap in row %d: %s and %s", 
						row, bars[i].TaskID, bars[j].TaskID))
				}
			}
		}
	}
	
	// Check for bars extending beyond calendar bounds
	for _, bar := range taskBars {
		if bar.StartX < 0 || bar.EndX > float64(mle.calendarEnd.Sub(mle.calendarStart).Hours()/24)*mle.dayWidth {
			issues = append(issues, fmt.Sprintf("Task bar %s extends beyond calendar bounds", bar.TaskID))
		}
	}
	
	return issues
}

// barsOverlap checks if two task bars overlap
func (mle *MultiDayLayoutEngine) barsOverlap(bar1, bar2 *TaskBar) bool {
	// Bars overlap if one starts before the other ends
	return !(bar1.EndX <= bar2.StartX || bar2.EndX <= bar1.StartX)
}
