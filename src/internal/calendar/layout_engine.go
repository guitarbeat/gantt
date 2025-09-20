package calendar

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"phd-dissertation-planner/internal/shared"
)

// LayoutEngine handles both multi-day task layout and calendar grid integration
type LayoutEngine struct {
	// Multi-day layout engine components
	calendarStart    time.Time
	calendarEnd      time.Time
	dayWidth         float64
	dayHeight        float64
	rowHeight        float64
	maxRowsPerDay    int
	overlapThreshold float64

	// Calendar grid integration components
	stackingEngine           *StackingEngine
	taskPrioritizationEngine *TaskPrioritizationEngine
	conflictResolutionEngine *ConflictResolutionEngine
	positioningEngine        *PositioningEngine
	monthBoundaryEngine      *MonthBoundaryEngine
	gridConfig              *GridConfig
	visualSettings          *IntegratedVisualSettings
	dateValidator           *shared.DateValidator
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

// IntegratedTaskBar represents a task bar with integrated smart stacking
type IntegratedTaskBar struct {
	TaskID           string
	StartDate        time.Time
	EndDate          time.Time
	StartX           float64
	EndX             float64
	Y                float64
	Width            float64
	Height           float64
	Row              int
	StackIndex       int
	Color            string
	BorderColor      string
	Opacity          float64
	ZIndex           int
	IsContinuation   bool
	IsStart          bool
	IsEnd            bool
	MonthBoundary    bool
	StackingType     StackingType
	VisualWeight     float64
	ProminenceScore  float64
	IsCollapsed      bool
	IsVisible        bool
	CollisionLevel   int
	OverflowLevel    int
	Priority         int
	Category         string
	TaskName         string
	Description      string
}

// TaskGroup represents a group of overlapping tasks
type TaskGroup struct {
	Tasks     []*shared.Task
	StartDate time.Time
	EndDate   time.Time
	Rows      int
}

// GridConfig defines the configuration for the calendar grid
type GridConfig struct {
	CalendarStart      time.Time
	CalendarEnd        time.Time
	DayWidth           float64
	DayHeight          float64
	RowHeight          float64
	MaxRowsPerDay      int
	OverlapThreshold   float64
	MonthBoundaryGap   float64
	TaskSpacing        float64
	VisualConstraints  *VisualConstraints
}

// IntegratedVisualSettings defines visual settings for the integrated system
type IntegratedVisualSettings struct {
	ShowTaskNames        bool
	ShowTaskDurations    bool
	ShowTaskPriorities   bool
	ShowConflictIndicators bool
	CollapseThreshold    int
	AnimationEnabled     bool
	HighlightConflicts   bool
	ColorScheme          string
	FontSize             string
	TaskBarOpacity       float64
	BorderWidth          float64
}

// IntegratedLayoutResult contains the result of integrated layout operations
type IntegratedLayoutResult struct {
	TaskBars           []*IntegratedTaskBar
	Stacks             []*TaskStack
	Conflicts          []*ResolvedConflict
	OverflowResolutions []*OverflowResolution
	VisualOptimizations []*VisualOptimization
	LayoutAdjustments  []*LayoutAdjustment
	Statistics         *IntegratedLayoutStatistics
	Recommendations    []string
	AnalysisDate       time.Time
}

// IntegratedLayoutStatistics contains statistics about the integrated layout
type IntegratedLayoutStatistics struct {
	TotalTasks           int
	ProcessedBars        int
	TotalStacks          int
	ConflictsResolved    int
	OverflowResolutions  int
	VisualOptimizations  int
	LayoutAdjustments    int
	CollisionCount       int
	OverflowCount        int
	MonthBoundaryCount   int
	SpaceEfficiency      float64
	VisualQuality        float64
	AverageStackHeight   float64
	MaxStackHeight       float64
	AverageTaskHeight    float64
	AverageTaskWidth     float64
	AlignmentScore       float64
	SpacingScore         float64
	VisualBalance        float64
	GridUtilization      float64
}

// MultiDayLayoutResult contains the results of multi-day layout processing
type MultiDayLayoutResult struct {
	TaskBars        []*TaskBar
	ValidationResult []shared.DataValidationError
	LayoutIssues    []string
	TaskCount       int
	ProcessedCount  int
}

// LayoutStatistics contains statistics about the layout
type LayoutStatistics struct {
	TotalTasks         int
	ProcessedBars      int
	ValidationErrors   int
	LayoutIssues       int
	OverlapCount       int
	MonthBoundaryCount int
}

// NewLayoutEngine creates a new layout engine instance
func NewLayoutEngine(config *GridConfig) *LayoutEngine {
	// Create stacking engine
	overlapDetector := NewOverlapDetector(config.CalendarStart, config.CalendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	stackingEngine := NewStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create visibility manager and stacking optimizer
	visibilityManager := NewVisibilityManager()
	stackingOptimizer := NewStackingOptimizer()
	
	// Create task prioritization engine
	taskPrioritizationEngine := NewTaskPrioritizationEngine(stackingEngine, priorityRanker, visibilityManager, stackingOptimizer)
	
	// Create conflict resolution engine
	conflictResolutionEngine := NewConflictResolutionEngine(taskPrioritizationEngine, stackingEngine)
	
	// Create positioning engine
	positioningEngine := NewPositioningEngine(config)
	
	// Create month boundary engine
	monthBoundaryEngine := NewMonthBoundaryEngine(config)
	
	// Create date validator
	dateValidator := shared.NewDateValidator()
	
	// Set visual constraints
	if config.VisualConstraints == nil {
		config.VisualConstraints = &VisualConstraints{
			MaxStackHeight:     config.DayHeight * float64(config.MaxRowsPerDay),
			MinTaskHeight:      config.RowHeight * 0.5,
			MaxTaskHeight:      config.RowHeight * 2.0,
			MinTaskWidth:       config.DayWidth * 0.1,
			MaxTaskWidth:       config.DayWidth * 7.0, // Max 7 days
			VerticalSpacing:    config.TaskSpacing,
			HorizontalSpacing:  config.TaskSpacing,
			MaxStackDepth:      config.MaxRowsPerDay,
			CollisionThreshold: config.OverlapThreshold,
			OverflowThreshold:  0.8,
		}
	}
	
	return &LayoutEngine{
		calendarStart:            config.CalendarStart,
		calendarEnd:              config.CalendarEnd,
		dayWidth:                 config.DayWidth,
		dayHeight:                config.DayHeight,
		rowHeight:                config.RowHeight,
		maxRowsPerDay:            config.MaxRowsPerDay,
		overlapThreshold:         config.OverlapThreshold,
		stackingEngine:           stackingEngine,
		taskPrioritizationEngine: taskPrioritizationEngine,
		conflictResolutionEngine: conflictResolutionEngine,
		positioningEngine:        positioningEngine,
		monthBoundaryEngine:      monthBoundaryEngine,
		gridConfig:              config,
		visualSettings: &IntegratedVisualSettings{
			ShowTaskNames:         true,
			ShowTaskDurations:     true,
			ShowTaskPriorities:    true,
			ShowConflictIndicators: true,
			CollapseThreshold:     5,
			AnimationEnabled:      false,
			HighlightConflicts:    true,
			ColorScheme:           "default",
			FontSize:              "small",
			TaskBarOpacity:        1.0,
			BorderWidth:           0.5,
		},
		dateValidator:            dateValidator,
	}
}

// LayoutMultiDayTasks performs the two-step algorithm for multi-day task layout
func (le *LayoutEngine) LayoutMultiDayTasks(tasks []*shared.Task) []*TaskBar {
	// Step 1: Group overlapping tasks
	groups := le.groupOverlappingTasks(tasks)
	
	// Step 2: Layout calculation within groups
	var taskBars []*TaskBar
	for _, group := range groups {
		groupBars := le.layoutTaskGroup(group)
		taskBars = append(taskBars, groupBars...)
	}
	
	return taskBars
}

// groupOverlappingTasks implements Step 1: Grouping Overlapping Events
func (le *LayoutEngine) groupOverlappingTasks(tasks []*shared.Task) []*TaskGroup {
	// Sort tasks by start date and duration
	sortedTasks := make([]*shared.Task, len(tasks))
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
			Tasks:     []*shared.Task{task},
			StartDate: task.StartDate,
			EndDate:   task.EndDate,
		}
		used[task.ID] = true
		
		// Find all tasks that overlap with this group
		for _, otherTask := range sortedTasks {
			if used[otherTask.ID] {
				continue
			}
			
			if le.tasksOverlap(group, otherTask) {
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
		group.Rows = le.calculateGroupRows(group)
		groups = append(groups, group)
	}
	
	return groups
}

// tasksOverlap checks if a task overlaps with a group
func (le *LayoutEngine) tasksOverlap(group *TaskGroup, task *shared.Task) bool {
	// Check if task overlaps with any task in the group
	for _, groupTask := range group.Tasks {
		if le.tasksOverlapDirect(groupTask, task) {
			return true
		}
	}
	return false
}

// tasksOverlapDirect checks if two tasks overlap directly
func (le *LayoutEngine) tasksOverlapDirect(task1, task2 *shared.Task) bool {
	// Tasks overlap if one starts before the other ends
	return !task1.StartDate.After(task2.EndDate) && !task2.StartDate.After(task1.EndDate)
}

// calculateGroupRows calculates the number of rows needed for a group
func (le *LayoutEngine) calculateGroupRows(group *TaskGroup) int {
	// If no tasks, return 1 row
	if len(group.Tasks) == 0 {
		return 1
	}
	
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
	if rows > le.maxRowsPerDay {
		rows = le.maxRowsPerDay
	}
	
	return rows
}

// layoutTaskGroup implements Step 2: Layout Calculation within Groups
func (le *LayoutEngine) layoutTaskGroup(group *TaskGroup) []*TaskBar {
	var taskBars []*TaskBar
	
	// If no tasks in group, return empty result
	if len(group.Tasks) == 0 {
		return taskBars
	}
	
	// Debug: ensure group.Rows is at least 1
	if group.Rows <= 0 {
		group.Rows = 1
	}
	
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
		row := le.findAvailableRow(task, rowEndTimes)
		
		// Ensure row is within bounds
		if row >= len(rowEndTimes) {
			row = 0
		}
		
		// Create task bar
		taskBar := le.createTaskBar(task, row, group.Rows)
		taskBars = append(taskBars, taskBar)
		
		// Update row end time
		rowEndTimes[row] = task.EndDate
	}
	
	return taskBars
}

// findAvailableRow finds the first available row for a task
func (le *LayoutEngine) findAvailableRow(task *shared.Task, rowEndTimes []time.Time) int {
	// If no rows available, return 0
	if len(rowEndTimes) == 0 {
		return 0
	}
	
	for i, endTime := range rowEndTimes {
		if task.StartDate.After(endTime) || task.StartDate.Equal(endTime) {
			return i
		}
	}
	
	// If no row is available, use the first row (overlap will be handled visually)
	return 0
}

// createTaskBar creates a task bar with positioning information
func (le *LayoutEngine) createTaskBar(task *shared.Task, row, totalRows int) *TaskBar {
	// Calculate X coordinates based on start and end dates
	startX := le.calculateXPosition(task.StartDate)
	endX := le.calculateXPosition(task.EndDate)
	
	// Calculate Y position based on row
	y := le.calculateYPosition(row, totalRows)
	
	// Calculate width
	width := endX - startX
	
	// Get task color from category
	category := shared.GetCategory(task.Category)
	
	// Determine if this is a continuation, start, or end
	isContinuation := le.isTaskContinuation(task)
	isStart := le.isTaskStart(task)
	isEnd := le.isTaskEnd(task)
	monthBoundary := le.hasMonthBoundary(task)
	
	return &TaskBar{
		TaskID:         task.ID,
		StartDate:      task.StartDate,
		EndDate:        task.EndDate,
		StartX:         startX,
		EndX:           endX,
		Y:              y,
		Width:          width,
		Height:         le.rowHeight,
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
func (le *LayoutEngine) calculateXPosition(date time.Time) float64 {
	// Calculate days from calendar start
	daysFromStart := int(date.Sub(le.calendarStart).Hours() / 24)
	
	// Calculate X position (day width * days from start)
	return float64(daysFromStart) * le.dayWidth
}

// calculateYPosition calculates the Y position for a given row
func (le *LayoutEngine) calculateYPosition(row, totalRows int) float64 {
	// Distribute rows evenly within the day height
	rowSpacing := le.dayHeight / float64(totalRows+1)
	return float64(row+1) * rowSpacing
}

// isTaskContinuation checks if this task is a continuation from previous month
func (le *LayoutEngine) isTaskContinuation(task *shared.Task) bool {
	// Check if task started before calendar start
	return task.StartDate.Before(le.calendarStart)
}

// isTaskStart checks if this is the start of a multi-day task
func (le *LayoutEngine) isTaskStart(task *shared.Task) bool {
	// Check if task starts on or after calendar start
	return !task.StartDate.Before(le.calendarStart)
}

// isTaskEnd checks if this is the end of a multi-day task
func (le *LayoutEngine) isTaskEnd(task *shared.Task) bool {
	// Check if task ends on or before calendar end
	return !task.EndDate.After(le.calendarEnd)
}

// hasMonthBoundary checks if task spans across month boundaries
func (le *LayoutEngine) hasMonthBoundary(task *shared.Task) bool {
	startMonth := task.StartDate.Month()
	endMonth := task.EndDate.Month()
	return startMonth != endMonth
}

// ProcessTasksWithSmartStacking processes tasks with integrated smart stacking
func (le *LayoutEngine) ProcessTasksWithSmartStacking(tasks []*shared.Task) (*IntegratedLayoutResult, error) {
	// Step 1: Detect overlaps and conflicts
	overlapAnalysis := le.stackingEngine.overlapDetector.DetectOverlaps(tasks)
	
	// Step 2: Prioritize tasks
	priorityContext := &PriorityContext{
		CurrentTime: time.Now(),
		UserID:      "system",
	}
	
	// Create priority management engine for task prioritization
	priorityManagementEngine := NewPriorityManagementEngine(
		le.stackingEngine.conflictCategorizer,
		le.stackingEngine,
	)
	
	prioritizationResult := priorityManagementEngine.PrioritizeTasks(tasks, priorityContext)
	
	// Step 3: Create stacking context
	stackingContext := &StackingContext{
		CalendarStart:     le.gridConfig.CalendarStart,
		CalendarEnd:       le.gridConfig.CalendarEnd,
		CurrentTime:       time.Now(),
		DayWidth:          le.gridConfig.DayWidth,
		DayHeight:         le.gridConfig.DayHeight,
		AvailableHeight:   le.gridConfig.DayHeight * float64(le.gridConfig.MaxRowsPerDay),
		AvailableWidth:    le.gridConfig.DayWidth * 7.0, // Max 7 days
		ExistingStacks:    []*TaskStack{},
		TaskPriorities:    make(map[string]*TaskPriority),
		ConflictAnalysis:  nil,
		OverlapAnalysis:   overlapAnalysis,
		VisualSettings:    &VisualSettings{
			ShowTaskNames:         le.visualSettings.ShowTaskNames,
			ShowTaskDurations:     le.visualSettings.ShowTaskDurations,
			ShowTaskPriorities:    le.visualSettings.ShowTaskPriorities,
			ShowConflictIndicators: le.visualSettings.ShowConflictIndicators,
			CollapseThreshold:     le.visualSettings.CollapseThreshold,
			AnimationEnabled:      le.visualSettings.AnimationEnabled,
			HighlightConflicts:    le.visualSettings.HighlightConflicts,
			ColorScheme:           le.visualSettings.ColorScheme,
		},
		VisualConstraints: le.gridConfig.VisualConstraints,
	}
	
	// Step 4: Apply smart stacking
	stackingResult := le.stackingEngine.StackTasks(tasks, stackingContext)
	
	// Step 5: Apply vertical stacking
	verticalStackingResult := le.stackingEngine.StackTasksVertically(tasks, stackingContext)
	
	// Step 6: Resolve conflicts
	conflictResolutionResult := le.conflictResolutionEngine.ResolveConflicts(tasks, priorityContext)
	
	// Step 7: Create integrated task bars
	integratedBars := le.createIntegratedTaskBars(tasks, stackingResult, verticalStackingResult, conflictResolutionResult, prioritizationResult)
	
	// Step 8: Apply precise positioning
	positioningResult, err := le.positioningEngine.PositionTasks(tasks, integratedBars)
	if err != nil {
		return nil, fmt.Errorf("failed to position tasks: %v", err)
	}
	
	// Step 9: Handle month boundaries with dedicated engine
	monthBoundaryResult, err := le.monthBoundaryEngine.ProcessMonthBoundaries(positioningResult.TaskBars, time.Now().Month(), time.Now().Year())
	if err != nil {
		return nil, fmt.Errorf("failed to process month boundaries: %v", err)
	}
	
	processedBars := monthBoundaryResult.ProcessedBars
	
	// Step 10: Calculate statistics
	statistics := le.calculateIntegratedStatistics(processedBars, stackingResult, conflictResolutionResult)
	
	// Merge positioning metrics
	statistics.AlignmentScore = positioningResult.Metrics.AlignmentScore
	statistics.SpacingScore = positioningResult.Metrics.SpacingScore
	statistics.VisualBalance = positioningResult.Metrics.VisualBalance
	statistics.GridUtilization = positioningResult.Metrics.GridUtilization
	
	// Merge month boundary metrics
	statistics.MonthBoundaryCount = monthBoundaryResult.BoundaryMetrics.ContinuationsCreated
	
	// Step 11: Generate recommendations
	recommendations := le.generateRecommendations(statistics, conflictResolutionResult)
	recommendations = append(recommendations, positioningResult.Recommendations...)
	recommendations = append(recommendations, monthBoundaryResult.Recommendations...)
	
	return &IntegratedLayoutResult{
		TaskBars:            processedBars,
		Stacks:              stackingResult.Stacks,
		Conflicts:           conflictResolutionResult.ResolvedConflicts,
		OverflowResolutions: conflictResolutionResult.OverflowResolutions,
		VisualOptimizations: conflictResolutionResult.VisualOptimizations,
		LayoutAdjustments:   conflictResolutionResult.LayoutAdjustments,
		Statistics:          statistics,
		Recommendations:     recommendations,
		AnalysisDate:        time.Now(),
	}, nil
}

// createIntegratedTaskBars creates integrated task bars with smart stacking
func (le *LayoutEngine) createIntegratedTaskBars(
	tasks []*shared.Task,
	stackingResult *StackingResult,
	verticalStackingResult *VerticalStackingResult,
	conflictResolutionResult *ConflictResolutionResult,
	prioritizationResult *TaskPrioritizationResult,
) []*IntegratedTaskBar {
	var integratedBars []*IntegratedTaskBar
	
	// Create a map of task priorities for quick lookup
	priorityMap := make(map[string]*TaskPriority)
	for _, prioritizedTask := range prioritizationResult.PrioritizedTasks {
		priorityMap[prioritizedTask.Task.ID] = &TaskPriority{
			Value:       prioritizedTask.Task.Priority,
			Category:    prioritizedTask.Task.Category,
			Description: prioritizedTask.Task.Description,
			Weight:      prioritizedTask.PriorityScore.OverallScore,
			Urgency:     string(prioritizedTask.PriorityScore.VisualProminence),
			Importance:  string(prioritizedTask.PriorityScore.VisualProminence),
		}
	}
	
	// Process each task
	for _, task := range tasks {
		// Calculate basic positioning
		startX := le.calculateXPosition(task.StartDate)
		endX := le.calculateXPosition(task.EndDate)
		width := endX - startX
		
		// Get task priority
		priority := priorityMap[task.ID]
		if priority == nil {
			priority = &TaskPriority{
				Value:       task.Priority,
				Category:    task.Category,
				Description: task.Description,
				Weight:      0.5,
				Urgency:     "MEDIUM",
				Importance:  "MEDIUM",
			}
		}
		
		// Calculate visual weight and prominence
		visualWeight := le.calculateVisualWeight(task, priority)
		prominenceScore := le.calculateProminenceScore(task, priority, visualWeight)
		
		// Determine stacking type and position
		stackingType, stackIndex, y, height := le.determineStackingPosition(
			task, stackingResult, verticalStackingResult, visualWeight,
		)
		
		// Get task category and color
		category := shared.GetCategory(task.Category)
		
		// Create integrated task bar
		integratedBar := &IntegratedTaskBar{
			TaskID:          task.ID,
			StartDate:       task.StartDate,
			EndDate:         task.EndDate,
			StartX:          startX,
			EndX:            endX,
			Y:               y,
			Width:           width,
			Height:          height,
			Row:             0, // Will be calculated based on stacking
			StackIndex:      stackIndex,
			Color:           category.Color,
			BorderColor:     "#000000",
			Opacity:         le.visualSettings.TaskBarOpacity,
			ZIndex:          int(priority.Weight * 5),
			IsContinuation:  le.isTaskContinuation(task),
			IsStart:         le.isTaskStart(task),
			IsEnd:           le.isTaskEnd(task),
			MonthBoundary:   le.hasMonthBoundary(task),
			StackingType:    stackingType,
			VisualWeight:    visualWeight,
			ProminenceScore: prominenceScore,
			IsCollapsed:     false,
			IsVisible:       true,
			CollisionLevel:  0,
			OverflowLevel:   0,
			Priority:        int(priority.Weight * 5),
			Category:        task.Category,
			TaskName:        task.Name,
			Description:     task.Description,
		}
		
		integratedBars = append(integratedBars, integratedBar)
	}
	
	return integratedBars
}

// calculateVisualWeight calculates the visual weight of a task
func (le *LayoutEngine) calculateVisualWeight(task *shared.Task, priority *TaskPriority) float64 {
	// Base weight from priority
	weight := priority.Weight
	
	// Adjust based on task duration
	duration := task.EndDate.Sub(task.StartDate).Hours() / 24
	if duration > 7 {
		weight *= 1.2 // Longer tasks get more visual weight
	} else if duration < 1 {
		weight *= 0.8 // Shorter tasks get less visual weight
	}
	
	// Adjust based on category
	category := shared.GetCategory(task.Category)
	weight *= float64(category.Priority) / 5.0
	
	// Adjust based on milestone status
	if strings.Contains(strings.ToUpper(task.Name), "MILESTONE") {
		weight *= 1.5
	}
	
	return math.Min(weight, 1.0)
}

// calculateProminenceScore calculates the prominence score of a task
func (le *LayoutEngine) calculateProminenceScore(task *shared.Task, priority *TaskPriority, visualWeight float64) float64 {
	// Base prominence from visual weight
	prominence := visualWeight
	
	// Adjust based on priority (using weight as proxy)
	prominence *= priority.Weight
	
	// Adjust based on urgency (convert string to float)
	urgencyMultiplier := 0.5 // Default
	switch priority.Urgency {
	case "CRITICAL":
		urgencyMultiplier = 1.0
	case "HIGH":
		urgencyMultiplier = 0.8
	case "MEDIUM":
		urgencyMultiplier = 0.6
	case "LOW":
		urgencyMultiplier = 0.4
	case "MINIMAL":
		urgencyMultiplier = 0.2
	}
	prominence *= urgencyMultiplier
	
	// Adjust based on milestone priority
	if priority.Category == "MILESTONE" {
		prominence *= 1.2
	}
	
	return math.Min(prominence, 1.0)
}

// determineStackingPosition determines the stacking position for a task
func (le *LayoutEngine) determineStackingPosition(
	task *shared.Task,
	stackingResult *StackingResult,
	verticalStackingResult *VerticalStackingResult,
	visualWeight float64,
) (StackingType, int, float64, float64) {
	// Find the appropriate stack for this task
	for _, stack := range stackingResult.Stacks {
		for _, stackedTask := range stack.Tasks {
			if stackedTask.Task.ID == task.ID {
				// Calculate Y position based on stack and position within stack
				y := le.calculateYPositionInStack(stackedTask.Position.Y, stack.TotalHeight)
				height := le.calculateTaskHeight(task, visualWeight)
				
				return stack.StackingType, stackedTask.Position.ZIndex, y, height
			}
		}
	}
	
	// Default positioning if not found in stacks
	y := le.gridConfig.DayHeight * 0.1 // 10% from top
	height := le.calculateTaskHeight(task, visualWeight)
	
	return StackingTypeVertical, 0, y, height
}

// calculateYPositionInStack calculates the Y position within a stack
func (le *LayoutEngine) calculateYPositionInStack(relativeY, stackHeight float64) float64 {
	// Normalize relative Y position to actual Y position
	return relativeY * (le.gridConfig.DayHeight / stackHeight)
}

// calculateTaskHeight calculates the height of a task based on its properties
func (le *LayoutEngine) calculateTaskHeight(task *shared.Task, visualWeight float64) float64 {
	// Base height
	height := le.gridConfig.RowHeight
	
	// Adjust based on visual weight
	height *= visualWeight
	
	// Ensure within constraints
	minHeight := le.gridConfig.VisualConstraints.MinTaskHeight
	maxHeight := le.gridConfig.VisualConstraints.MaxTaskHeight
	
	if height < minHeight {
		height = minHeight
	} else if height > maxHeight {
		height = maxHeight
	}
	
	return height
}

// HandleMonthBoundary handles task bars that span across month boundaries
func (le *LayoutEngine) HandleMonthBoundary(taskBars []*TaskBar) []*TaskBar {
	var processedBars []*TaskBar
	
	for _, bar := range taskBars {
		if !bar.MonthBoundary {
			processedBars = append(processedBars, bar)
			continue
		}
		
		// Split task bar at month boundaries
		splitBars := le.splitTaskBarAtMonthBoundaries(bar)
		processedBars = append(processedBars, splitBars...)
	}
	
	return processedBars
}

// splitTaskBarAtMonthBoundaries splits a task bar at month boundaries
func (le *LayoutEngine) splitTaskBarAtMonthBoundaries(bar *TaskBar) []*TaskBar {
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
			StartX:         le.calculateXPosition(current),
			EndX:           le.calculateXPosition(monthEnd),
			Y:              bar.Y,
			Width:          le.calculateXPosition(monthEnd) - le.calculateXPosition(current),
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
func (le *LayoutEngine) GenerateLaTeX(taskBars []*TaskBar) string {
	var latex strings.Builder
	
	// Group task bars by day for efficient rendering
	dayGroups := le.groupTaskBarsByDay(taskBars)
	
	for day, bars := range dayGroups {
		latex.WriteString(le.generateDayLaTeX(day, bars))
	}
	
	return latex.String()
}

// groupTaskBarsByDay groups task bars by day
func (le *LayoutEngine) groupTaskBarsByDay(taskBars []*TaskBar) map[time.Time][]*TaskBar {
	dayGroups := make(map[time.Time][]*TaskBar)
	
	for _, bar := range taskBars {
		// Group by start date
		day := time.Date(bar.StartDate.Year(), bar.StartDate.Month(), bar.StartDate.Day(), 0, 0, 0, 0, bar.StartDate.Location())
		dayGroups[day] = append(dayGroups[day], bar)
	}
	
	return dayGroups
}

// generateDayLaTeX generates LaTeX code for a specific day
func (le *LayoutEngine) generateDayLaTeX(day time.Time, bars []*TaskBar) string {
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
		latex.WriteString(le.generateTaskBarLaTeX(bar))
	}
	
	return latex.String()
}

// generateTaskBarLaTeX generates LaTeX code for a single task bar
func (le *LayoutEngine) generateTaskBarLaTeX(bar *TaskBar) string {
	// Convert colors to LaTeX format
	color := le.convertColorToLaTeX(bar.Color)

    // Generate task bar LaTeX via centralized macro
    return fmt.Sprintf(`\\DrawTaskBar{%.2f}{%.2f}{%.2f}{%.2f}{%s}{%s}`,
        bar.StartX, bar.Y, bar.Width, bar.Height, color, bar.TaskID)
}

// convertColorToLaTeX converts hex color to LaTeX color name
func (le *LayoutEngine) convertColorToLaTeX(hexColor string) string {
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
func (le *LayoutEngine) ValidateLayout(taskBars []*TaskBar) []string {
	var issues []string
	
	// Check for overlapping task bars in the same row
	rowBars := make(map[int][]*TaskBar)
	for _, bar := range taskBars {
		rowBars[bar.Row] = append(rowBars[bar.Row], bar)
	}
	
	for row, bars := range rowBars {
		for i := 0; i < len(bars); i++ {
			for j := i + 1; j < len(bars); j++ {
				if le.barsOverlap(bars[i], bars[j]) {
					issues = append(issues, fmt.Sprintf("Task bars overlap in row %d: %s and %s", 
						row, bars[i].TaskID, bars[j].TaskID))
				}
			}
		}
	}
	
	// Check for bars extending beyond calendar bounds
	for _, bar := range taskBars {
		if bar.StartX < 0 || bar.EndX > float64(le.calendarEnd.Sub(le.calendarStart).Hours()/24)*le.dayWidth {
			issues = append(issues, fmt.Sprintf("Task bar %s extends beyond calendar bounds", bar.TaskID))
		}
	}
	
	return issues
}

// barsOverlap checks if two task bars overlap
func (le *LayoutEngine) barsOverlap(bar1, bar2 *TaskBar) bool {
	// Bars overlap if one starts before the other ends
	return !(bar1.EndX <= bar2.StartX || bar2.EndX <= bar1.StartX)
}

// calculateIntegratedStatistics calculates statistics for the integrated layout
func (le *LayoutEngine) calculateIntegratedStatistics(
	bars []*IntegratedTaskBar,
	stackingResult *StackingResult,
	conflictResolutionResult *ConflictResolutionResult,
) *IntegratedLayoutStatistics {
	stats := &IntegratedLayoutStatistics{
		TotalTasks:          len(bars),
		ProcessedBars:       len(bars),
		TotalStacks:         len(stackingResult.Stacks),
		ConflictsResolved:   len(conflictResolutionResult.ResolvedConflicts),
		OverflowResolutions: len(conflictResolutionResult.OverflowResolutions),
		VisualOptimizations: len(conflictResolutionResult.VisualOptimizations),
		LayoutAdjustments:   len(conflictResolutionResult.LayoutAdjustments),
		SpaceEfficiency:     stackingResult.SpaceEfficiency,
		VisualQuality:       stackingResult.VisualQuality,
	}
	
	// Calculate additional statistics
	var totalHeight, maxHeight, totalWidth float64
	monthBoundaryCount := 0
	
	for _, bar := range bars {
		totalHeight += bar.Height
		totalWidth += bar.Width
		
		if bar.Height > maxHeight {
			maxHeight = bar.Height
		}
		
		if bar.MonthBoundary {
			monthBoundaryCount++
		}
	}
	
	stats.MonthBoundaryCount = monthBoundaryCount
	stats.MaxStackHeight = maxHeight
	
	if len(bars) > 0 {
		stats.AverageTaskHeight = totalHeight / float64(len(bars))
		stats.AverageTaskWidth = totalWidth / float64(len(bars))
	}
	
	// Calculate average stack height
	if len(stackingResult.Stacks) > 0 {
		var totalStackHeight float64
		for _, stack := range stackingResult.Stacks {
			totalStackHeight += stack.TotalHeight
		}
		stats.AverageStackHeight = totalStackHeight / float64(len(stackingResult.Stacks))
	}
	
	return stats
}

// generateRecommendations generates recommendations based on the layout analysis
func (le *LayoutEngine) generateRecommendations(
	statistics *IntegratedLayoutStatistics,
	conflictResolutionResult *ConflictResolutionResult,
) []string {
	var recommendations []string
	
	// Space efficiency recommendations
	if statistics.SpaceEfficiency < 0.7 {
		recommendations = append(recommendations, "Consider reducing task spacing to improve space efficiency")
	}
	
	// Visual quality recommendations
	if statistics.VisualQuality < 0.8 {
		recommendations = append(recommendations, "Consider adjusting task heights and colors to improve visual quality")
	}
	
	// Stack height recommendations
	if statistics.AverageStackHeight > le.gridConfig.DayHeight*2 {
		recommendations = append(recommendations, "Consider using horizontal stacking for high-density days")
	}
	
	// Conflict recommendations
	if statistics.ConflictsResolved > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Resolved %d visual conflicts - consider reviewing task scheduling", statistics.ConflictsResolved))
	}
	
	// Overflow recommendations
	if statistics.OverflowResolutions > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Applied %d overflow resolutions - consider reducing task density", statistics.OverflowResolutions))
	}
	
	return recommendations
}

// GenerateIntegratedLaTeX generates LaTeX code for the integrated calendar
func (le *LayoutEngine) GenerateIntegratedLaTeX(result *IntegratedLayoutResult) string {
	var latex strings.Builder
	
	// Generate header
	latex.WriteString("\\begin{integrated-calendar}\n")
	
	// Generate task bars LaTeX
	for _, bar := range result.TaskBars {
		barLaTeX := le.generateIntegratedTaskBarLaTeX(bar)
		latex.WriteString(barLaTeX)
	}
	
	// Generate footer
	latex.WriteString("\\end{integrated-calendar}\n")
	
	return latex.String()
}

// generateIntegratedTaskBarLaTeX generates LaTeX code for a single integrated task bar
func (le *LayoutEngine) generateIntegratedTaskBarLaTeX(bar *IntegratedTaskBar) string {
	// Create TikZ node for the task bar
    var nodeOptions string
    if bar.Opacity >= 0.999 {
        nodeOptions = fmt.Sprintf(
            "anchor=west, inner sep=2pt, minimum height=%.2fpt, minimum width=%.2fpt, fill=%s",
            bar.Height,
            bar.Width,
            bar.Color,
        )
    } else {
        nodeOptions = fmt.Sprintf(
            "anchor=west, inner sep=2pt, minimum height=%.2fpt, minimum width=%.2fpt, fill=%s, opacity=%.2f",
            bar.Height,
            bar.Width,
            bar.Color,
            bar.Opacity,
        )
    }
	
	// Add border if specified
	if le.visualSettings.BorderWidth > 0 {
		nodeOptions += fmt.Sprintf(", draw=%s, line width=%.2fpt", bar.BorderColor, le.visualSettings.BorderWidth)
	}
	
	// Create the TikZ node
	tikzNode := fmt.Sprintf(
		"\\node[%s] at (%.2fpt, %.2fpt) {%s};",
		nodeOptions,
		bar.StartX,
		bar.Y,
		bar.TaskName,
	)
	
	return tikzNode + "\n"
}

// ProcessTasksWithValidation processes tasks with validation and creates multi-day layout
func (le *LayoutEngine) ProcessTasksWithValidation(tasks []*shared.Task) (*MultiDayLayoutResult, error) {
	// Validate tasks first
	validationResult := le.dateValidator.ValidateDateRanges(tasks)
	
	// Check for critical errors
	if len(validationResult) > 0 {
		// Log validation errors but continue with layout
		fmt.Printf("Warning: %d validation issues found\n", len(validationResult))
		for _, err := range validationResult {
			if err.Severity == "ERROR" {
				fmt.Printf("Error: %s\n", err.Message)
			}
		}
	}
	
	// Filter out tasks with critical errors for layout
	validTasks := le.filterValidTasks(tasks, validationResult)
	
	// Create multi-day layout
	taskBars := le.LayoutMultiDayTasks(validTasks)
	
	// Handle month boundaries
	processedBars := le.HandleMonthBoundary(taskBars)
	
	// Validate layout
	layoutIssues := le.ValidateLayout(processedBars)
	
	return &MultiDayLayoutResult{
		TaskBars:        processedBars,
		ValidationResult: validationResult,
		LayoutIssues:    layoutIssues,
		TaskCount:       len(validTasks),
		ProcessedCount:  len(processedBars),
	}, nil
}

// filterValidTasks filters out tasks with critical validation errors
func (le *LayoutEngine) filterValidTasks(tasks []*shared.Task, validationErrors []shared.DataValidationError) []*shared.Task {
	// Create map of tasks with critical errors
	errorTasks := make(map[string]bool)
	for _, err := range validationErrors {
		if err.Severity == "ERROR" {
			errorTasks[err.TaskID] = true
		}
	}
	
	// Filter out tasks with critical errors
	var validTasks []*shared.Task
	for _, task := range tasks {
		if !errorTasks[task.ID] {
			validTasks = append(validTasks, task)
		}
	}
	
	return validTasks
}

// GenerateCalendarLaTeX generates LaTeX code for the calendar with multi-day task bars
func (le *LayoutEngine) GenerateCalendarLaTeX(result *MultiDayLayoutResult) string {
	var latex strings.Builder
	
	// Generate header
	latex.WriteString("\\begin{calendar}\n")
	
	// Generate task bars LaTeX
	taskBarsLaTeX := le.GenerateLaTeX(result.TaskBars)
	latex.WriteString(taskBarsLaTeX)
	
	// Generate footer
	latex.WriteString("\\end{calendar}\n")
	
	return latex.String()
}

// GetLayoutStatistics returns statistics about the layout
func (le *LayoutEngine) GetLayoutStatistics(result *MultiDayLayoutResult) *LayoutStatistics {
	stats := &LayoutStatistics{
		TotalTasks:      result.TaskCount,
		ProcessedBars:   result.ProcessedCount,
		ValidationErrors: len(result.ValidationResult),
		LayoutIssues:    len(result.LayoutIssues),
		OverlapCount:    0,
		MonthBoundaryCount: 0,
	}
	
	// Count overlaps and month boundaries
	for _, bar := range result.TaskBars {
		if bar.MonthBoundary {
			stats.MonthBoundaryCount++
		}
	}
	
	// Count overlaps by checking for overlapping bars
	rowBars := make(map[int][]*TaskBar)
	for _, bar := range result.TaskBars {
		rowBars[bar.Row] = append(rowBars[bar.Row], bar)
	}
	
	for _, bars := range rowBars {
		for i := 0; i < len(bars); i++ {
			for j := i + 1; j < len(bars); j++ {
				if le.barsOverlap(bars[i], bars[j]) {
					stats.OverlapCount++
				}
			}
		}
	}
	
	return stats
}

// GetIntegratedStatistics returns statistics about the integrated layout
func (le *LayoutEngine) GetIntegratedStatistics(result *IntegratedLayoutResult) *IntegratedLayoutStatistics {
	return result.Statistics
}

// String returns a string representation of the integrated layout statistics
func (ils *IntegratedLayoutStatistics) String() string {
	return fmt.Sprintf("Integrated Layout Statistics:\n"+
		"  Total Tasks: %d\n"+
		"  Processed Bars: %d\n"+
		"  Total Stacks: %d\n"+
		"  Conflicts Resolved: %d\n"+
		"  Overflow Resolutions: %d\n"+
		"  Visual Optimizations: %d\n"+
		"  Layout Adjustments: %d\n"+
		"  Collision Count: %d\n"+
		"  Overflow Count: %d\n"+
		"  Month Boundary Count: %d\n"+
		"  Space Efficiency: %.2f\n"+
		"  Visual Quality: %.2f\n"+
		"  Average Stack Height: %.2f\n"+
		"  Max Stack Height: %.2f\n"+
		"  Average Task Height: %.2f\n"+
		"  Average Task Width: %.2f\n"+
		"  Alignment Score: %.2f\n"+
		"  Spacing Score: %.2f\n"+
		"  Visual Balance: %.2f\n"+
		"  Grid Utilization: %.2f\n",
		ils.TotalTasks, ils.ProcessedBars, ils.TotalStacks,
		ils.ConflictsResolved, ils.OverflowResolutions, ils.VisualOptimizations,
		ils.LayoutAdjustments, ils.CollisionCount, ils.OverflowCount,
		ils.MonthBoundaryCount, ils.SpaceEfficiency, ils.VisualQuality,
		ils.AverageStackHeight, ils.MaxStackHeight, ils.AverageTaskHeight, ils.AverageTaskWidth,
		ils.AlignmentScore, ils.SpacingScore, ils.VisualBalance, ils.GridUtilization)
}

// String returns a string representation of the statistics
func (ls *LayoutStatistics) String() string {
	return fmt.Sprintf("Layout Statistics:\n"+
		"  Total Tasks: %d\n"+
		"  Processed Bars: %d\n"+
		"  Validation Errors: %d\n"+
		"  Layout Issues: %d\n"+
		"  Overlaps: %d\n"+
		"  Month Boundaries: %d\n",
		ls.TotalTasks, ls.ProcessedBars, ls.ValidationErrors,
		ls.LayoutIssues, ls.OverlapCount, ls.MonthBoundaryCount)
}
