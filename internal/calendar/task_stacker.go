// Package calendar provides task stacking and overlap detection for multi-day tasks.
//
// This module handles:
// - Detection of overlapping tasks across multiple days
// - Intelligent stacking of tasks to prevent visual overlap
// - Assignment of vertical layers (tracks) to tasks
// - Optimization of track usage for compact visualization
package calendar

import (
	"sort"
	"time"
)

// TaskStack represents a vertical track assignment for tasks
type TaskStack struct {
	Track    int           // Vertical position (0 = top, 1 = second, etc.)
	Task     *SpanningTask // The task in this track
	StartCol int           // Starting column in the week (0-6, Monday-Sunday)
	EndCol   int           // Ending column in the week (0-6, Monday-Sunday)
}

// DayTaskStack represents all tasks that should be displayed on a specific day
type DayTaskStack struct {
	Date   time.Time
	Stacks []TaskStack // All tasks visible on this day, organized by track
}

// TaskStacker manages task overlap detection and track assignment
type TaskStacker struct {
	tasks        []*SpanningTask
	dayStacks    map[string]*DayTaskStack // Key: date string (YYYY-MM-DD)
	maxTracks    int                      // Maximum number of tracks needed for any day
	weekStartDay time.Weekday             // First day of week (Monday = 1)
}

// NewTaskStacker creates a new task stacker
func NewTaskStacker(tasks []*SpanningTask, weekStartDay time.Weekday) *TaskStacker {
	return &TaskStacker{
		tasks:        tasks,
		dayStacks:    make(map[string]*DayTaskStack),
		weekStartDay: weekStartDay,
	}
}

// ComputeStacks analyzes all tasks and computes optimal track assignments
// This is the main algorithm that prevents visual overlaps
func (ts *TaskStacker) ComputeStacks() {
	// Sort tasks by start date, then by duration (longer tasks first)
	sortedTasks := ts.sortTasksByPriority()

	// For each task, find the lowest available track across all its days
	for _, task := range sortedTasks {
		track := ts.findLowestAvailableTrack(task)
		ts.assignTaskToTrack(task, track)
	}
}

// sortTasksByPriority sorts tasks for optimal stacking
// Priority: earlier start date, then longer duration
func (ts *TaskStacker) sortTasksByPriority() []*SpanningTask {
	sorted := make([]*SpanningTask, len(ts.tasks))
	copy(sorted, ts.tasks)

	sort.Slice(sorted, func(i, j int) bool {
		// First, sort by start date
		if !sorted[i].StartDate.Equal(sorted[j].StartDate) {
			return sorted[i].StartDate.Before(sorted[j].StartDate)
		}

		// If start dates are equal, longer tasks first
		durationI := sorted[i].EndDate.Sub(sorted[i].StartDate)
		durationJ := sorted[j].EndDate.Sub(sorted[j].StartDate)
		return durationI > durationJ
	})

	return sorted
}

// findLowestAvailableTrack finds the lowest track number that's free for all days of the task
func (ts *TaskStacker) findLowestAvailableTrack(task *SpanningTask) int {
	// Get all dates this task spans
	dates := ts.getDateRange(task.StartDate, task.EndDate)

	// Check each track starting from 0
	for track := 0; track < 100; track++ { // Reasonable upper limit
		available := true

		// Check if this track is available for ALL days the task spans
		for _, date := range dates {
			if ts.isTrackOccupied(date, track) {
				available = false
				break
			}
		}

		if available {
			return track
		}
	}

	// Fallback: return a high track number
	return 99
}

// isTrackOccupied checks if a track is occupied on a specific date
func (ts *TaskStacker) isTrackOccupied(date time.Time, track int) bool {
	dateKey := ts.dateKey(date)
	dayStack, exists := ts.dayStacks[dateKey]

	if !exists {
		return false
	}

	// Check if any task in this track overlaps this date
	for _, stack := range dayStack.Stacks {
		if stack.Track == track {
			return true
		}
	}

	return false
}

// assignTaskToTrack assigns a task to a specific track for all its days
func (ts *TaskStacker) assignTaskToTrack(task *SpanningTask, track int) {
	dates := ts.getDateRange(task.StartDate, task.EndDate)

	for _, date := range dates {
		dateKey := ts.dateKey(date)

		// Get or create day stack
		dayStack, exists := ts.dayStacks[dateKey]
		if !exists {
			dayStack = &DayTaskStack{
				Date:   date,
				Stacks: []TaskStack{},
			}
			ts.dayStacks[dateKey] = dayStack
		}

		// Calculate column positions for this specific week
		startCol, endCol := ts.calculateWeekColumns(task, date)

		// Add task to this day's stack
		dayStack.Stacks = append(dayStack.Stacks, TaskStack{
			Track:    track,
			Task:     task,
			StartCol: startCol,
			EndCol:   endCol,
		})

		// Update max tracks
		if track+1 > ts.maxTracks {
			ts.maxTracks = track + 1
		}
	}
}

// calculateWeekColumns calculates which columns a task occupies in the week containing the given date
func (ts *TaskStacker) calculateWeekColumns(task *SpanningTask, date time.Time) (startCol, endCol int) {
	// Get the week boundaries for this date
	weekStart := ts.getWeekStart(date)
	weekEnd := weekStart.AddDate(0, 0, 6)

	// Clamp task dates to week boundaries
	taskStart := task.StartDate
	taskEnd := task.EndDate

	if taskStart.Before(weekStart) {
		taskStart = weekStart
	}
	if taskEnd.After(weekEnd) {
		taskEnd = weekEnd
	}

	// Calculate column indices (0-6 for Mon-Sun if weekStartDay is Monday)
	startCol = int((taskStart.Weekday() - ts.weekStartDay + 7) % 7)
	endCol = int((taskEnd.Weekday() - ts.weekStartDay + 7) % 7)

	return startCol, endCol
}

// getWeekStart returns the first day of the week containing the given date
func (ts *TaskStacker) getWeekStart(date time.Time) time.Time {
	// Normalize to start of day
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Calculate days back to week start
	daysBack := int((date.Weekday() - ts.weekStartDay + 7) % 7)
	weekStart := date.AddDate(0, 0, -daysBack)

	return weekStart
}

// GetStacksForDay returns all task stacks for a specific day, sorted by track
func (ts *TaskStacker) GetStacksForDay(date time.Time) []TaskStack {
	dateKey := ts.dateKey(date)
	dayStack, exists := ts.dayStacks[dateKey]

	if !exists {
		return []TaskStack{}
	}

	// Sort by track number
	stacks := make([]TaskStack, len(dayStack.Stacks))
	copy(stacks, dayStack.Stacks)

	sort.Slice(stacks, func(i, j int) bool {
		return stacks[i].Track < stacks[j].Track
	})

	return stacks
}

// GetTasksStartingOnDay returns tasks that START on the given day, with their track assignments
func (ts *TaskStacker) GetTasksStartingOnDay(date time.Time) []TaskStack {
	dateKey := ts.dateKey(date)
	dayStack, exists := ts.dayStacks[dateKey]

	if !exists {
		return []TaskStack{}
	}

	// Filter to only tasks that start on this day
	var startingTasks []TaskStack
	normalizedDate := ts.normalizeDate(date)

	for _, stack := range dayStack.Stacks {
		taskStartDate := ts.normalizeDate(stack.Task.StartDate)
		if taskStartDate.Equal(normalizedDate) {
			startingTasks = append(startingTasks, stack)
		}
	}

	// Sort by track number
	sort.Slice(startingTasks, func(i, j int) bool {
		return startingTasks[i].Track < startingTasks[j].Track
	})

	return startingTasks
}

// GetMaxTracks returns the maximum number of tracks needed
func (ts *TaskStacker) GetMaxTracks() int {
	return ts.maxTracks
}

// Helper methods

// dateKey creates a unique key for a date (YYYY-MM-DD format)
func (ts *TaskStacker) dateKey(date time.Time) string {
	return date.Format("2006-01-02")
}

// normalizeDate normalizes a date to midnight UTC
func (ts *TaskStacker) normalizeDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

// getDateRange returns all dates between start and end (inclusive)
func (ts *TaskStacker) getDateRange(start, end time.Time) []time.Time {
	start = ts.normalizeDate(start)
	end = ts.normalizeDate(end)

	var dates []time.Time
	current := start

	for !current.After(end) {
		dates = append(dates, current)
		current = current.AddDate(0, 0, 1)
	}

	return dates
}

// TaskStackRenderer provides rendering information for stacked tasks
type TaskStackRenderer struct {
	stacker *TaskStacker
}

// NewTaskStackRenderer creates a new renderer for stacked tasks
func NewTaskStackRenderer(tasks []*SpanningTask, weekStartDay time.Weekday) *TaskStackRenderer {
	stacker := NewTaskStacker(tasks, weekStartDay)
	stacker.ComputeStacks()

	return &TaskStackRenderer{
		stacker: stacker,
	}
}

// GetRenderInfoForDay returns rendering information for tasks on a specific day
func (tsr *TaskStackRenderer) GetRenderInfoForDay(date time.Time) *DayRenderInfo {
	stacks := tsr.stacker.GetStacksForDay(date)
	startingTasks := tsr.stacker.GetTasksStartingOnDay(date)

	return &DayRenderInfo{
		Date:               date,
		AllStacks:          stacks,
		StartingTasks:      startingTasks,
		MaxTracks:          tsr.stacker.GetMaxTracks(),
		TasksVisibleOnDay:  len(stacks),
		TasksStartingOnDay: len(startingTasks),
	}
}

// DayRenderInfo contains all information needed to render tasks for a specific day
type DayRenderInfo struct {
	Date               time.Time
	AllStacks          []TaskStack // All tasks visible on this day (including those that started earlier)
	StartingTasks      []TaskStack // Only tasks that START on this day
	MaxTracks          int         // Maximum tracks for proper cell height
	TasksVisibleOnDay  int         // Total number of tasks visible
	TasksStartingOnDay int         // Number of tasks starting on this day
}

// ShouldShowContinuation returns true if this day shows continuation bars for tasks started earlier
func (dri *DayRenderInfo) ShouldShowContinuation() bool {
	return dri.TasksVisibleOnDay > dri.TasksStartingOnDay
}

// GetContinuationTasks returns tasks that are continuing from previous days
func (dri *DayRenderInfo) GetContinuationTasks() []TaskStack {
	if !dri.ShouldShowContinuation() {
		return []TaskStack{}
	}

	normalizedDate := time.Date(dri.Date.Year(), dri.Date.Month(), dri.Date.Day(), 0, 0, 0, 0, time.UTC)

	var continuations []TaskStack
	for _, stack := range dri.AllStacks {
		taskStartDate := time.Date(
			stack.Task.StartDate.Year(),
			stack.Task.StartDate.Month(),
			stack.Task.StartDate.Day(),
			0, 0, 0, 0, time.UTC,
		)

		if taskStartDate.Before(normalizedDate) {
			continuations = append(continuations, stack)
		}
	}

	return continuations
}
