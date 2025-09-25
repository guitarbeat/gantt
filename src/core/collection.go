package core

import (
	"sort"
	"time"
)

// TaskCollection represents a collection of tasks with efficient access patterns
type TaskCollection struct {
	tasks      []*Task
	byDate     []*Task
	byCategory map[string][]*Task
	byStatus   map[string][]*Task
	byAssignee map[string][]*Task
	sorted     bool
}

// NewTaskCollection creates a new empty task collection
func NewTaskCollection() *TaskCollection {
	return &TaskCollection{
		tasks:      make([]*Task, 0),
		byDate:     make([]*Task, 0),
		byCategory: make(map[string][]*Task),
		byStatus:   make(map[string][]*Task),
		byAssignee: make(map[string][]*Task),
		sorted:     false,
	}
}

// AddTask adds a task to the collection
func (tc *TaskCollection) AddTask(task *Task) {
	if task == nil {
		return
	}

	tc.tasks = append(tc.tasks, task)
	tc.byDate = append(tc.byDate, task)

	// Update category index
	if task.Category != "" {
		tc.byCategory[task.Category] = append(tc.byCategory[task.Category], task)
	}

	// Update status index
	if task.Status != "" {
		tc.byStatus[task.Status] = append(tc.byStatus[task.Status], task)
	}

	// Update assignee index
	if task.Assignee != "" {
		tc.byAssignee[task.Assignee] = append(tc.byAssignee[task.Assignee], task)
	}

	tc.sorted = false
}

// GetTask retrieves a task by name (since we removed ID)
func (tc *TaskCollection) GetTask(name string) (*Task, bool) {
	for _, task := range tc.tasks {
		if task.Name == name {
			return task, true
		}
	}
	return nil, false
}

// GetAllTasks returns all tasks in the collection
func (tc *TaskCollection) GetAllTasks() []*Task {
	return tc.tasks
}

// GetTasksByCategory returns all tasks in a specific category
func (tc *TaskCollection) GetTasksByCategory(category string) []*Task {
	return tc.byCategory[category]
}

// GetTasksByStatus returns all tasks with a specific status
func (tc *TaskCollection) GetTasksByStatus(status string) []*Task {
	return tc.byStatus[status]
}

// GetTasksByAssignee returns all tasks assigned to a specific person
func (tc *TaskCollection) GetTasksByAssignee(assignee string) []*Task {
	return tc.byAssignee[assignee]
}

// GetTasksByDateRange returns all tasks within a date range
func (tc *TaskCollection) GetTasksByDateRange(start, end time.Time) []*Task {
	if !tc.sorted {
		tc.sortByDate()
	}

	var result []*Task
	for _, task := range tc.byDate {
		if task.OverlapsWithDateRange(start, end) {
			result = append(result, task)
		}
	}
	return result
}

// GetTasksByDate returns all tasks on a specific date
func (tc *TaskCollection) GetTasksByDate(date time.Time) []*Task {
	if !tc.sorted {
		tc.sortByDate()
	}

	var result []*Task
	for _, task := range tc.byDate {
		if task.IsOnDate(date) {
			result = append(result, task)
		}
	}
	return result
}

// sortByDate sorts tasks by start date
func (tc *TaskCollection) sortByDate() {
	sort.Slice(tc.byDate, func(i, j int) bool {
		return tc.byDate[i].StartDate.Before(tc.byDate[j].StartDate)
	})
	tc.sorted = true
}

// TaskHierarchy represents the parent-child hierarchy of tasks
type TaskHierarchy struct {
	roots    []*Task
	parents  map[string]*Task
	children map[string][]*Task
	tasks    []*Task
}

// NewTaskHierarchy creates a new task hierarchy
func NewTaskHierarchy() *TaskHierarchy {
	return &TaskHierarchy{
		roots:    make([]*Task, 0),
		parents:  make(map[string]*Task),
		children: make(map[string][]*Task),
		tasks:    make([]*Task, 0),
	}
}

// AddTask adds a task to the hierarchy
func (th *TaskHierarchy) AddTask(task *Task) {
	if task == nil {
		return
	}

	th.tasks = append(th.tasks, task)

	if task.ParentID == "" {
		// This is a root task
		th.roots = append(th.roots, task)
	} else {
		// This is a child task - find parent by name
		for _, parent := range th.tasks {
			if parent.Name == task.ParentID {
				th.parents[task.Name] = parent
				th.children[task.ParentID] = append(th.children[task.ParentID], task)
				break
			}
		}
	}
}

// GetRootTasks returns all root tasks (tasks without parents)
func (th *TaskHierarchy) GetRootTasks() []*Task {
	return th.roots
}

// GetChildren returns all child tasks of a given task
func (th *TaskHierarchy) GetChildren(taskName string) []*Task {
	return th.children[taskName]
}

// GetParent returns the parent task of a given task
func (th *TaskHierarchy) GetParent(taskName string) *Task {
	return th.parents[taskName]
}

// GetAncestors returns all ancestor tasks of a given task
func (th *TaskHierarchy) GetAncestors(taskName string) []*Task {
	var ancestors []*Task
	current := th.GetParent(taskName)

	for current != nil {
		ancestors = append(ancestors, current)
		current = th.GetParent(current.Name)
	}

	return ancestors
}

// GetDescendants returns all descendant tasks of a given task
func (th *TaskHierarchy) GetDescendants(taskName string) []*Task {
	var descendants []*Task
	children := th.GetChildren(taskName)

	for _, child := range children {
		descendants = append(descendants, child)
		descendants = append(descendants, th.GetDescendants(child.Name)...)
	}

	return descendants
}

// CalendarLayout represents optimized date range calculations for calendar rendering
type CalendarLayout struct {
	startDate time.Time
	endDate   time.Time
	tasks     []*Task
	months    []MonthYear
	weeks     []time.Time
}

// NewCalendarLayout creates a new calendar layout for a date range
func NewCalendarLayout(startDate, endDate time.Time, tasks []*Task) *CalendarLayout {
	cl := &CalendarLayout{
		startDate: startDate,
		endDate:   endDate,
		tasks:     tasks,
	}

	cl.generateMonths()
	cl.generateWeeks()

	return cl
}

// generateMonths generates all months in the date range
func (cl *CalendarLayout) generateMonths() {
	cl.months = make([]MonthYear, 0)

	current := time.Date(cl.startDate.Year(), cl.startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(cl.endDate.Year(), cl.endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !current.After(end) {
		cl.months = append(cl.months, MonthYear{
			Year:  current.Year(),
			Month: current.Month(),
		})
		current = current.AddDate(0, 1, 0)
	}
}

// generateWeeks generates all weeks in the date range
func (cl *CalendarLayout) generateWeeks() {
	cl.weeks = make([]time.Time, 0)

	// Find the start of the first week
	start := cl.startDate
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, -1)
	}

	// Generate weeks until we cover the end date
	for !start.After(cl.endDate) {
		cl.weeks = append(cl.weeks, start)
		start = start.AddDate(0, 0, 7)
	}
}

// GetMonths returns all months in the layout
func (cl *CalendarLayout) GetMonths() []MonthYear {
	return cl.months
}

// GetWeeks returns all weeks in the layout
func (cl *CalendarLayout) GetWeeks() []time.Time {
	return cl.weeks
}

// GetTasksForMonth returns all tasks that occur in a specific month
func (cl *CalendarLayout) GetTasksForMonth(year int, month time.Month) []*Task {
	var result []*Task
	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, -1)

	for _, task := range cl.tasks {
		if task.OverlapsWithDateRange(monthStart, monthEnd) {
			result = append(result, task)
		}
	}

	return result
}

// GetTasksForWeek returns all tasks that occur in a specific week
func (cl *CalendarLayout) GetTasksForWeek(weekStart time.Time) []*Task {
	var result []*Task
	weekEnd := weekStart.AddDate(0, 0, 6)

	for _, task := range cl.tasks {
		if task.OverlapsWithDateRange(weekStart, weekEnd) {
			result = append(result, task)
		}
	}

	return result
}
