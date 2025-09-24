package scheduler

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"phd-dissertation-planner/internal/common"
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
	spatialEngine            *SpatialEngine
	gridConfig               *GridConfig
	visualSettings           *IntegratedVisualSettings
	dateValidator            *common.DateValidator
}

// TaskBar represents a rendered task bar with positioning information
type TaskBar struct {
	TaskID         string
	StartDate      time.Time
	EndDate        time.Time
	StartX         float64
	EndX           float64
	Y              float64
	Width          float64
	Height         float64
	Row            int
	Color          string
	BorderColor    string
	Opacity        float64
	ZIndex         int
	IsContinuation bool
	IsStart        bool
	IsEnd          bool
	MonthBoundary  bool
}

// IntegratedTaskBar represents a task bar with integrated smart stacking
type IntegratedTaskBar struct {
	TaskID          string
	StartDate       time.Time
	EndDate         time.Time
	StartX          float64
	EndX            float64
	Y               float64
	Width           float64
	Height          float64
	Row             int
	StackIndex      int
	Color           string
	BorderColor     string
	Opacity         float64
	ZIndex          int
	IsContinuation  bool
	IsStart         bool
	IsEnd           bool
	MonthBoundary   bool
	StackingType    StackingType
	VisualWeight    float64
	ProminenceScore float64
	IsCollapsed     bool
	IsVisible       bool
	CollisionLevel  int
	OverflowLevel   int
	Priority        int
	Category        string
	TaskName        string
	Description     string
}

// TaskGroup represents a group of overlapping tasks
type TaskGroup struct {
	Tasks     []*common.Task
	StartDate time.Time
	EndDate   time.Time
	Rows      int
}

// GridConfig defines the configuration for the calendar grid
type GridConfig struct {
	CalendarStart     time.Time
	CalendarEnd       time.Time
	DayWidth          float64
	DayHeight         float64
	RowHeight         float64
	MaxRowsPerDay     int
	OverlapThreshold  float64
	MonthBoundaryGap  float64
	TaskSpacing       float64
	VisualConstraints *VisualConstraints
}

// IntegratedVisualSettings defines visual settings for the integrated system
type IntegratedVisualSettings struct {
	ShowTaskNames          bool
	ShowTaskDurations      bool
	ShowTaskPriorities     bool
	ShowConflictIndicators bool
	CollapseThreshold      int
	AnimationEnabled       bool
	HighlightConflicts     bool
	ColorScheme            string
	FontSize               string
	TaskBarOpacity         float64
	BorderWidth            float64
}

// IntegratedLayoutResult contains the result of integrated layout operations
type IntegratedLayoutResult struct {
	TaskBars            []*IntegratedTaskBar
	Stacks              []*TaskStack
	Conflicts           []*ResolvedConflict
	OverflowResolutions []*OverflowResolution
	VisualOptimizations []*VisualOptimization
	LayoutAdjustments   []*LayoutAdjustment
	Statistics          *IntegratedLayoutStatistics
	Recommendations     []string
	AnalysisDate        time.Time
}

// IntegratedLayoutStatistics contains statistics about the integrated layout
type IntegratedLayoutStatistics struct {
	TotalTasks          int
	ProcessedBars       int
	TotalStacks         int
	ConflictsResolved   int
	OverflowResolutions int
	VisualOptimizations int
	LayoutAdjustments   int
	CollisionCount      int
	OverflowCount       int
	MonthBoundaryCount  int
	SpaceEfficiency     float64
	VisualQuality       float64
	AverageStackHeight  float64
	MaxStackHeight      float64
	AverageTaskHeight   float64
	AverageTaskWidth    float64
	AlignmentScore      float64
	SpacingScore        float64
	VisualBalance       float64
	GridUtilization     float64
}

// MultiDayLayoutResult contains the results of multi-day layout processing
type MultiDayLayoutResult struct {
	TaskBars         []*TaskBar
	ValidationResult []common.DataValidationError
	LayoutIssues     []string
	TaskCount        int
	ProcessedCount   int
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

// Month boundary types and methods (consolidated from month_boundary_engine.go)

// BoundaryRule defines how tasks should behave at month boundaries
type BoundaryRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*IntegratedTaskBar, *MonthBoundaryContext) bool
	Action      func(*IntegratedTaskBar, *MonthBoundaryContext) *BoundaryAction
}

// TransitionRule defines how tasks should transition between months
type TransitionRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*IntegratedTaskBar, *MonthBoundaryContext) bool
	Action      func(*IntegratedTaskBar, *MonthBoundaryContext) *TransitionAction
}

// ContinuityRule defines how to maintain visual continuity across months
type ContinuityRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func([]*IntegratedTaskBar, *MonthBoundaryContext) bool
	Action      func([]*IntegratedTaskBar, *MonthBoundaryContext) *ContinuityAction
}

// MonthBoundaryContext provides context for month boundary decisions
type MonthBoundaryContext struct {
	CurrentMonth     time.Month
	NextMonth        time.Month
	CurrentYear      int
	NextYear         int
	CalendarStart    time.Time
	CalendarEnd      time.Time
	DayWidth         float64
	DayHeight        float64
	MonthBoundaryGap float64
	TransitionBuffer float64
	VisualSettings   *IntegratedVisualSettings
	GridConstraints  *GridConstraints
	TaskDensity      float64
	OverlapCount     int
	ConflictCount    int
}

// BoundaryAction defines how a task should be handled at month boundaries
type BoundaryAction struct {
	SplitTask          bool
	CreateContinuation bool
	AdjustPosition     bool
	AdjustWidth        bool
	AdjustHeight       bool
	NewX               float64
	NewY               float64
	NewWidth           float64
	NewHeight          float64
	ContinuationID     string
	VisualStyle        *BoundaryVisualStyle
	Priority           int
}

// TransitionAction defines how a task should transition between months
type TransitionAction struct {
	SmoothTransition bool
	FadeIn           bool
	FadeOut          bool
	SlideAnimation   bool
	ScaleAnimation   bool
	Duration         time.Duration
	EasingFunction   EasingFunction
	VisualEffects    []VisualEffect
	Priority         int
}

// ContinuityAction defines how to maintain visual continuity
type ContinuityAction struct {
	MaintainAlignment bool
	PreserveSpacing   bool
	ConsistentColors  bool
	ConsistentSizes   bool
	VisualConnections []VisualConnection
	Priority          int
}

// BoundaryVisualStyle defines visual styling for month boundaries
type BoundaryVisualStyle struct {
	BorderStyle      BorderStyle
	BorderColor      string
	BorderWidth      float64
	BackgroundColor  string
	Opacity          float64
	ShadowEnabled    bool
	ShadowColor      string
	ShadowBlur       float64
	ShadowOffsetX    float64
	ShadowOffsetY    float64
	HighlightEnabled bool
	HighlightColor   string
	HighlightWidth   float64
}

// BorderStyle defines the style of borders at month boundaries
type BorderStyle string

const (
	BorderSolid  BorderStyle = "SOLID"
	BorderDashed BorderStyle = "DASHED"
	BorderDotted BorderStyle = "DOTTED"
	BorderDouble BorderStyle = "DOUBLE"
	BorderGroove BorderStyle = "GROOVE"
	BorderRidge  BorderStyle = "RIDGE"
	BorderInset  BorderStyle = "INSET"
	BorderOutset BorderStyle = "OUTSET"
	BorderNone   BorderStyle = "NONE"
)

// EasingFunction defines easing functions for transitions
type EasingFunction string

const (
	EasingLinear    EasingFunction = "LINEAR"
	EasingEaseIn    EasingFunction = "EASE_IN"
	EasingEaseOut   EasingFunction = "EASE_OUT"
	EasingEaseInOut EasingFunction = "EASE_IN_OUT"
	EasingBounce    EasingFunction = "BOUNCE"
	EasingElastic   EasingFunction = "ELASTIC"
	EasingBack      EasingFunction = "BACK"
	EasingCubic     EasingFunction = "CUBIC"
	EasingQuart     EasingFunction = "QUART"
	EasingQuint     EasingFunction = "QUINT"
	EasingSine      EasingFunction = "SINE"
	EasingExpo      EasingFunction = "EXPO"
	EasingCirc      EasingFunction = "CIRC"
)

// VisualEffect defines visual effects for transitions
type VisualEffect string

const (
	EffectFadeIn     VisualEffect = "FADE_IN"
	EffectFadeOut    VisualEffect = "FADE_OUT"
	EffectSlideLeft  VisualEffect = "SLIDE_LEFT"
	EffectSlideRight VisualEffect = "SLIDE_RIGHT"
	EffectSlideUp    VisualEffect = "SLIDE_UP"
	EffectSlideDown  VisualEffect = "SLIDE_DOWN"
	EffectScaleIn    VisualEffect = "SCALE_IN"
	EffectScaleOut   VisualEffect = "SCALE_OUT"
	EffectRotateIn   VisualEffect = "ROTATE_IN"
	EffectRotateOut  VisualEffect = "ROTATE_OUT"
	EffectFlipIn     VisualEffect = "FLIP_IN"
	EffectFlipOut    VisualEffect = "FLIP_OUT"
	EffectZoomIn     VisualEffect = "ZOOM_IN"
	EffectZoomOut    VisualEffect = "ZOOM_OUT"
)

// VisualConnection defines visual connections between months
type VisualConnection struct {
	FromTaskID     string
	ToTaskID       string
	ConnectionType ConnectionType
	LineStyle      LineStyle
	LineColor      string
	LineWidth      float64
	ArrowStyle     ArrowStyle
	Label          string
	Priority       int
}

// ConnectionType defines the type of visual connection
type ConnectionType string

const (
	ConnectionArrow  ConnectionType = "ARROW"
	ConnectionLine   ConnectionType = "LINE"
	ConnectionCurve  ConnectionType = "CURVE"
	ConnectionDashed ConnectionType = "DASHED"
	ConnectionDotted ConnectionType = "DOTTED"
	ConnectionThick  ConnectionType = "THICK"
	ConnectionThin   ConnectionType = "THIN"
	ConnectionDouble ConnectionType = "DOUBLE"
)

// LineStyle defines the style of connection lines
type LineStyle string

const (
	LineSolid      LineStyle = "SOLID"
	LineDashed     LineStyle = "DASHED"
	LineDotted     LineStyle = "DOTTED"
	LineDashDot    LineStyle = "DASH_DOT"
	LineDashDotDot LineStyle = "DASH_DOT_DOT"
)

// ArrowStyle defines the style of arrows
type ArrowStyle string

const (
	ArrowNone      ArrowStyle = "NONE"
	ArrowSimple    ArrowStyle = "SIMPLE"
	ArrowFilled    ArrowStyle = "FILLED"
	ArrowHollow    ArrowStyle = "HOLLOW"
	ArrowDouble    ArrowStyle = "DOUBLE"
	ArrowCurved    ArrowStyle = "CURVED"
	ArrowBarbed    ArrowStyle = "BARBED"
	ArrowFeathered ArrowStyle = "FEATHERED"
)

// MonthBoundaryResult contains the result of month boundary processing
type MonthBoundaryResult struct {
	ProcessedBars     []*IntegratedTaskBar
	Continuations     []*TaskContinuation
	Transitions       []*TaskTransition
	VisualConnections []*VisualConnection
	BoundaryMetrics   *BoundaryMetrics
	Recommendations   []string
	AnalysisDate      time.Time
}

// TaskContinuation represents a task continuation across month boundaries
type TaskContinuation struct {
	OriginalTaskID   string
	ContinuationID   string
	StartMonth       time.Month
	EndMonth         time.Month
	StartYear        int
	EndYear          int
	ContinuationType ContinuationType
	VisualStyle      *BoundaryVisualStyle
	ConnectionStyle  *VisualConnection
	Priority         int
}

// TaskTransition represents a task transition between months
type TaskTransition struct {
	TaskID         string
	FromMonth      time.Month
	ToMonth        time.Month
	TransitionType TransitionType
	Animation      *TransitionAnimation
	VisualEffects  []VisualEffect
	Duration       time.Duration
	EasingFunction EasingFunction
	Priority       int
}

// ContinuationType defines the type of task continuation
type ContinuationType string

const (
	ContinuationSplit    ContinuationType = "SPLIT"
	ContinuationExtend   ContinuationType = "EXTEND"
	ContinuationWrap     ContinuationType = "WRAP"
	ContinuationOverflow ContinuationType = "OVERFLOW"
	ContinuationTruncate ContinuationType = "TRUNCATE"
	ContinuationMinimize ContinuationType = "MINIMIZE"
	ContinuationCollapse ContinuationType = "COLLAPSE"
)

// TransitionType defines the type of task transition
type TransitionType string

const (
	TransitionSmooth  TransitionType = "SMOOTH"
	TransitionFade    TransitionType = "FADE"
	TransitionSlide   TransitionType = "SLIDE"
	TransitionScale   TransitionType = "SCALE"
	TransitionRotate  TransitionType = "ROTATE"
	TransitionFlip    TransitionType = "FLIP"
	TransitionZoom    TransitionType = "ZOOM"
	TransitionBounce  TransitionType = "BOUNCE"
	TransitionElastic TransitionType = "ELASTIC"
)

// TransitionAnimation defines animation properties for transitions
type TransitionAnimation struct {
	Type           TransitionType
	Duration       time.Duration
	EasingFunction EasingFunction
	Delay          time.Duration
	IterationCount int
	Direction      AnimationDirection
	FillMode       FillMode
	PlayState      PlayState
}

// AnimationDirection defines the direction of animation
type AnimationDirection string

const (
	DirectionNormal           AnimationDirection = "NORMAL"
	DirectionReverse          AnimationDirection = "REVERSE"
	DirectionAlternate        AnimationDirection = "ALTERNATE"
	DirectionAlternateReverse AnimationDirection = "ALTERNATE_REVERSE"
)

// FillMode defines how animations fill their target values
type FillMode string

const (
	FillNone      FillMode = "NONE"
	FillForwards  FillMode = "FORWARDS"
	FillBackwards FillMode = "BACKWARDS"
	FillBoth      FillMode = "BOTH"
)

// PlayState defines the play state of animations
type PlayState string

const (
	PlayRunning PlayState = "RUNNING"
	PlayPaused  PlayState = "PAUSED"
	PlayStopped PlayState = "STOPPED"
)

// BoundaryMetrics contains metrics about month boundary processing
type BoundaryMetrics struct {
	TotalTasks           int
	ProcessedTasks       int
	ContinuationsCreated int
	TransitionsApplied   int
	VisualConnections    int
	BoundaryConflicts    int
	TransitionErrors     int
	ContinuityScore      float64
	VisualConsistency    float64
	TransitionSmoothness float64
	GridContinuity       float64
	SpaceEfficiency      float64
	VisualBalance        float64
}

// Task rendering types and methods (consolidated from task_rendering_engine.go)

// overlayInfo holds information about a spanning task overlay
type overlayInfo struct {
	content string
	cols    int
}

// TaskRenderingConfig holds configuration for task rendering
type TaskRenderingConfig struct {
	// Spacing configuration
	DefaultSpacing   string
	FirstTaskSpacing string

	// Height configuration
	DefaultHeight   string
	FirstTaskHeight string

	// Text configuration
	MaxChars            int
	MaxCharsCompact     int
	MaxCharsVeryCompact int
}

// NewLayoutEngine creates a new layout engine instance
func NewLayoutEngine(config *GridConfig) *LayoutEngine {
	// Create spatial engine (handles both overlap detection and positioning)
	spatialEngine := NewSpatialEngine(config.CalendarStart, config.CalendarEnd, config)

	// Create stacking engine
	conflictCategorizer := NewConflictCategorizer(spatialEngine)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	stackingEngine := NewStackingEngine(spatialEngine, conflictCategorizer, priorityRanker)

	// Create visibility manager and stacking optimizer
	visibilityManager := NewVisibilityManager()
	stackingOptimizer := NewStackingOptimizer()

	// Create task prioritization engine
	taskPrioritizationEngine := NewTaskPrioritizationEngine(stackingEngine, priorityRanker, visibilityManager, stackingOptimizer)

	// Create conflict resolution engine
	conflictResolutionEngine := NewConflictResolutionEngine(taskPrioritizationEngine, stackingEngine)

	// Month boundary fields will be initialized in the struct

	// Create date validator
	dateValidator := common.NewDateValidator()

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
		spatialEngine:            spatialEngine,
		gridConfig:               config,
		visualSettings: &IntegratedVisualSettings{
			ShowTaskNames:          true,
			ShowTaskDurations:      true,
			ShowTaskPriorities:     true,
			ShowConflictIndicators: true,
			CollapseThreshold:      5,
			AnimationEnabled:       false,
			HighlightConflicts:     false,
			ColorScheme:            "default",
			FontSize:               "small",
			TaskBarOpacity:         1.0,
			BorderWidth:            0.5,
		},
		dateValidator: dateValidator,
	}
}

// LayoutMultiDayTasks performs the two-step algorithm for multi-day task layout
func (le *LayoutEngine) LayoutMultiDayTasks(tasks []*common.Task) []*TaskBar {
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
func (le *LayoutEngine) groupOverlappingTasks(tasks []*common.Task) []*TaskGroup {
	// Sort tasks by start date and duration
	sortedTasks := make([]*common.Task, len(tasks))
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
			Tasks:     []*common.Task{task},
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
func (le *LayoutEngine) tasksOverlap(group *TaskGroup, task *common.Task) bool {
	// Check if task overlaps with any task in the group
	for _, groupTask := range group.Tasks {
		if le.tasksOverlapDirect(groupTask, task) {
			return true
		}
	}
	return false
}

// tasksOverlapDirect checks if two tasks overlap directly
func (le *LayoutEngine) tasksOverlapDirect(task1, task2 *common.Task) bool {
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
func (le *LayoutEngine) findAvailableRow(task *common.Task, rowEndTimes []time.Time) int {
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
func (le *LayoutEngine) createTaskBar(task *common.Task, row, totalRows int) *TaskBar {
	// Calculate X coordinates based on start and end dates
	startX := le.calculateXPosition(task.StartDate)
	endX := le.calculateXPosition(task.EndDate)

	// Calculate Y position based on row
	y := le.calculateYPosition(row, totalRows)

	// Calculate width
	width := endX - startX

	// Get task color from category
	category := common.GetCategory(task.Category)

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
func (le *LayoutEngine) isTaskContinuation(task *common.Task) bool {
	// Check if task started before calendar start
	return task.StartDate.Before(le.calendarStart)
}

// isTaskStart checks if this is the start of a multi-day task
func (le *LayoutEngine) isTaskStart(task *common.Task) bool {
	// Check if task starts on or after calendar start
	return !task.StartDate.Before(le.calendarStart)
}

// isTaskEnd checks if this is the end of a multi-day task
func (le *LayoutEngine) isTaskEnd(task *common.Task) bool {
	// Check if task ends on or before calendar end
	return !task.EndDate.After(le.calendarEnd)
}

// hasMonthBoundary checks if task spans across month boundaries
func (le *LayoutEngine) hasMonthBoundary(task *common.Task) bool {
	startMonth := task.StartDate.Month()
	endMonth := task.EndDate.Month()
	return startMonth != endMonth
}

// ProcessTasksWithSmartStacking processes tasks with integrated smart stacking
func (le *LayoutEngine) ProcessTasksWithSmartStacking(tasks []*common.Task) (*IntegratedLayoutResult, error) {
	// Step 1: Detect overlaps and conflicts
	overlapAnalysis := le.spatialEngine.DetectOverlaps(tasks)

	// Step 2: Prioritize tasks
	priorityContext := &PriorityContext{
		CurrentTime: time.Now(),
		UserID:      "system",
	}

	// Create priority management engine for task prioritization
	priorityManagementEngine := NewPriorityManagementEngine(
		le.spatialEngine,
		le.stackingEngine.conflictCategorizer,
		le.stackingEngine,
	)

	prioritizationResult := priorityManagementEngine.PrioritizeTasks(tasks, priorityContext)

	// Step 3: Create stacking context
	stackingContext := &StackingContext{
		CalendarStart:    le.gridConfig.CalendarStart,
		CalendarEnd:      le.gridConfig.CalendarEnd,
		CurrentTime:      time.Now(),
		DayWidth:         le.gridConfig.DayWidth,
		DayHeight:        le.gridConfig.DayHeight,
		AvailableHeight:  le.gridConfig.DayHeight * float64(le.gridConfig.MaxRowsPerDay),
		AvailableWidth:   le.gridConfig.DayWidth * 7.0, // Max 7 days
		ExistingStacks:   []*TaskStack{},
		TaskPriorities:   make(map[string]*TaskPriority),
		ConflictAnalysis: nil,
		OverlapAnalysis:  overlapAnalysis,
		VisualSettings: &VisualSettings{
			ShowTaskNames:          le.visualSettings.ShowTaskNames,
			ShowTaskDurations:      le.visualSettings.ShowTaskDurations,
			ShowTaskPriorities:     le.visualSettings.ShowTaskPriorities,
			ShowConflictIndicators: le.visualSettings.ShowConflictIndicators,
			CollapseThreshold:      le.visualSettings.CollapseThreshold,
			AnimationEnabled:       le.visualSettings.AnimationEnabled,
			HighlightConflicts:     le.visualSettings.HighlightConflicts,
			ColorScheme:            le.visualSettings.ColorScheme,
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
	positioningResult, err := le.spatialEngine.PositionTasks(tasks, integratedBars)
	if err != nil {
		return nil, fmt.Errorf("failed to position tasks: %v", err)
	}

	// Step 9: Handle month boundaries with dedicated engine
	monthBoundaryResult, err := le.ProcessMonthBoundaries(positioningResult.TaskBars, time.Now().Month(), time.Now().Year())
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
	tasks []*common.Task,
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
		category := common.GetCategory(task.Category)

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
func (le *LayoutEngine) calculateVisualWeight(task *common.Task, priority *TaskPriority) float64 {
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
	category := common.GetCategory(task.Category)
	weight *= float64(category.Priority) / 5.0

	// Adjust based on milestone status
	if strings.Contains(strings.ToUpper(task.Name), "MILESTONE") {
		weight *= 1.5
	}

	return math.Min(weight, 1.0)
}

// calculateProminenceScore calculates the prominence score of a task
func (le *LayoutEngine) calculateProminenceScore(task *common.Task, priority *TaskPriority, visualWeight float64) float64 {
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
	task *common.Task,
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
func (le *LayoutEngine) calculateTaskHeight(task *common.Task, visualWeight float64) float64 {
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
		"#4A90E2": "blue",   // PROPOSAL
		"#F5A623": "orange", // LASER
		"#7ED321": "green",  // IMAGING
		"#BD10E0": "purple", // ADMIN
		"#D0021B": "red",    // DISSERTATION
		"#50E3C2": "teal",   // RESEARCH
		"#B8E986": "lime",   // PUBLICATION
		"#CCCCCC": "gray",   // Default
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

// Month boundary methods (consolidated from month_boundary_engine.go)

// ProcessMonthBoundaries processes month boundary transitions for task bars
func (le *LayoutEngine) ProcessMonthBoundaries(taskBars []*IntegratedTaskBar, currentMonth time.Month, currentYear int) (*MonthBoundaryResult, error) {
	// Create month boundary context
	context := &MonthBoundaryContext{
		CurrentMonth:     currentMonth,
		NextMonth:        le.getNextMonth(currentMonth),
		CurrentYear:      currentYear,
		NextYear:         le.getNextYear(currentMonth, currentYear),
		CalendarStart:    le.calendarStart,
		CalendarEnd:      le.calendarEnd,
		DayWidth:         le.dayWidth,
		DayHeight:        le.dayHeight,
		MonthBoundaryGap: le.gridConfig.MonthBoundaryGap,
		TransitionBuffer: 2.0,
		VisualSettings:   le.visualSettings,
		GridConstraints:  le.getDefaultGridConstraints(),
		TaskDensity:      le.calculateTaskDensity(taskBars),
		OverlapCount:     le.countOverlaps(taskBars),
		ConflictCount:    le.countConflicts(taskBars),
	}

	// Process boundary rules
	processedBars := le.applyBoundaryRules(taskBars, context)

	// Process transition rules
	transitions := le.applyTransitionRules(processedBars, context)

	// Process continuity rules
	continuations := le.applyContinuityRules(processedBars, context)

	// Create visual connections
	visualConnections := le.createVisualConnections(processedBars, continuations, context)

	// Calculate metrics
	metrics := le.calculateBoundaryMetrics(processedBars, continuations, transitions, context)

	// Generate recommendations
	recommendations := le.generateBoundaryRecommendations(metrics, context)

	return &MonthBoundaryResult{
		ProcessedBars:     processedBars,
		Continuations:     continuations,
		Transitions:       transitions,
		VisualConnections: visualConnections,
		BoundaryMetrics:   metrics,
		Recommendations:   recommendations,
		AnalysisDate:      time.Now(),
	}, nil
}

// Helper methods for month boundary processing
func (le *LayoutEngine) getNextMonth(month time.Month) time.Month {
	if month == time.December {
		return time.January
	}
	return month + 1
}

func (le *LayoutEngine) getNextYear(month time.Month, year int) int {
	if month == time.December {
		return year + 1
	}
	return year
}

func (le *LayoutEngine) getDefaultGridConstraints() *GridConstraints {
	return &GridConstraints{
		MinTaskSpacing:     1.0,
		MaxTaskSpacing:     10.0,
		MinRowHeight:       8.0,
		MaxRowHeight:       20.0,
		MinColumnWidth:     5.0,
		MaxColumnWidth:     50.0,
		SnapToGrid:         true,
		GridResolution:     1.0,
		AlignmentTolerance: 0.5,
		CollisionBuffer:    2.0,
	}
}

func (le *LayoutEngine) calculateTaskDensity(bars []*IntegratedTaskBar) float64 {
	if len(bars) == 0 {
		return 0.0
	}

	// Calculate total calendar area
	totalArea := le.dayWidth * le.dayHeight * 7.0 * 4.0

	// Calculate average task area
	var totalTaskArea float64
	for _, bar := range bars {
		totalTaskArea += bar.Width * bar.Height
	}

	avgTaskArea := totalTaskArea / float64(len(bars))

	// Calculate density
	return (avgTaskArea * float64(len(bars))) / totalArea
}

func (le *LayoutEngine) countOverlaps(bars []*IntegratedTaskBar) int {
	count := 0
	for i := 0; i < len(bars); i++ {
		for j := i + 1; j < len(bars); j++ {
			if le.integratedBarsOverlap(bars[i], bars[j]) {
				count++
			}
		}
	}
	return count
}

func (le *LayoutEngine) countConflicts(bars []*IntegratedTaskBar) int {
	// For now, use overlap count as conflict count
	return le.countOverlaps(bars)
}

func (le *LayoutEngine) integratedBarsOverlap(bar1, bar2 *IntegratedTaskBar) bool {
	// Check horizontal overlap
	horizontalOverlap := bar1.StartX < bar2.EndX && bar2.StartX < bar1.EndX

	// Check vertical overlap
	verticalOverlap := bar1.Y < bar2.Y+bar2.Height && bar2.Y < bar1.Y+bar1.Height

	return horizontalOverlap && verticalOverlap
}

// Placeholder methods for month boundary processing
func (le *LayoutEngine) applyBoundaryRules(bars []*IntegratedTaskBar, context *MonthBoundaryContext) []*IntegratedTaskBar {
	// Simplified implementation - just return the bars as-is
	return bars
}

func (le *LayoutEngine) applyTransitionRules(bars []*IntegratedTaskBar, context *MonthBoundaryContext) []*TaskTransition {
	// Simplified implementation - return empty transitions
	return []*TaskTransition{}
}

func (le *LayoutEngine) applyContinuityRules(bars []*IntegratedTaskBar, context *MonthBoundaryContext) []*TaskContinuation {
	// Simplified implementation - return empty continuations
	return []*TaskContinuation{}
}

func (le *LayoutEngine) createVisualConnections(bars []*IntegratedTaskBar, continuations []*TaskContinuation, context *MonthBoundaryContext) []*VisualConnection {
	// Simplified implementation - return empty connections
	return []*VisualConnection{}
}

func (le *LayoutEngine) calculateBoundaryMetrics(bars []*IntegratedTaskBar, continuations []*TaskContinuation, transitions []*TaskTransition, context *MonthBoundaryContext) *BoundaryMetrics {
	return &BoundaryMetrics{
		TotalTasks:           len(bars),
		ProcessedTasks:       len(bars),
		ContinuationsCreated: len(continuations),
		TransitionsApplied:   len(transitions),
		VisualConnections:    len(continuations),
		ContinuityScore:      1.0,
		VisualConsistency:    1.0,
		TransitionSmoothness: 1.0,
		GridContinuity:       1.0,
		SpaceEfficiency:      1.0,
		VisualBalance:        1.0,
	}
}

func (le *LayoutEngine) generateBoundaryRecommendations(metrics *BoundaryMetrics, context *MonthBoundaryContext) []string {
	// Simplified implementation - return empty recommendations
	return []string{}
}

// Task rendering methods (consolidated from task_rendering_engine.go)

// getDefaultTaskRenderingConfig returns the default configuration for task rendering
//   - NO-OVERLAP CONFIGURATION: This configuration is optimized to prevent task overlap
//     by using increased spacing, larger heights, and limiting the number of displayed tasks
func getDefaultTaskRenderingConfig() TaskRenderingConfig {
	return TaskRenderingConfig{
		// Spacing configuration - increased to prevent overlap
		DefaultSpacing:   "0.8ex",
		FirstTaskSpacing: "0.5ex",

		// Height configuration - increased to prevent overlap
		DefaultHeight:   "3.0ex",
		FirstTaskHeight: "3.5ex",

		// Text configuration - from constants in day.go
		MaxChars:            maxTaskChars,
		MaxCharsCompact:     maxTaskCharsCompact,
		MaxCharsVeryCompact: maxTaskCharsVeryCompact,
	}
}

// renderSpanningTaskOverlay renders spanning task overlays for multiple tasks starting on this day
// Returns nil if no spanning tasks exist or none start on this day
func (d Day) renderSpanningTaskOverlay() *overlayInfo {
	if len(d.SpanningTasks) == 0 {
		return nil
	}
	

	dayDate := d.getDayDate()
	startingTasks, maxCols := d.findStartingTasks(dayDate)

	if len(startingTasks) == 0 {
		return nil
	}

	// Build content for all starting tasks using smart stacking
	content := d.buildMultiTaskOverlayContent(startingTasks)

	return &overlayInfo{
		content: content,
		cols:    maxCols,
	}
}

// buildTaskOverlayContent creates the LaTeX content for a single task overlay
// Used when only one task starts on a given day
func (d Day) buildTaskOverlayContent(task *SpanningTask) string {
	nameText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Name))
	descText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Description))

	// Add star indicator for milestone tasks
	if d.isMilestoneSpanningTask(task) {
		nameText = "★ " + nameText
	}

	// Use calendar macros for overlay with proper spacing
		// Convert hex color to RGB for LaTeX compatibility
		color := hexToRGB(task.Color)
		return `\vspace*{0.1ex}` + `\TaskOverlayBox{` + color + `}{` + nameText + `}{` + descText + `}`
}

// buildMultiTaskOverlayContent creates compact stacked content for multiple tasks
// Uses smart stacking to prevent overlap and improve readability
func (d Day) buildMultiTaskOverlayContent(tasks []*SpanningTask) string {
	if len(tasks) == 0 {
		return ""
	}

	// Single task - use full overlay format
	if len(tasks) == 1 {
		return d.buildTaskOverlayContent(tasks[0])
	}

	// Sort tasks by category priority for better visual organization
	sortedTasks := d.sortTasksByPriority(tasks)

	var contentParts []string

	// Show all tasks in compact format (no limit)
	for i := 0; i < len(sortedTasks); i++ {
		task := sortedTasks[i]
		compactContent := d.buildCompactTaskOverlay(task, i, len(sortedTasks))
		contentParts = append(contentParts, compactContent)
	}

	return strings.Join(contentParts, "")
}


// buildCompactTaskOverlay creates a compact task overlay for multiple tasks
// Used when multiple tasks start on the same day to create stacked display
func (d Day) buildCompactTaskOverlay(task *SpanningTask, index, total int) string {
	nameText := d.prepareTaskName(task)
	nameText = d.truncateTaskName(nameText, total)

	spacing, boxHeight := d.getTaskSpacingAndHeight(index)
	textBody := d.buildTaskTextBody(nameText)

	// Convert hex color to RGB for LaTeX compatibility
	color := hexToRGB(task.Color)
	return d.buildCompactTaskBox(spacing, boxHeight, color, textBody)
}

// prepareTaskName prepares the task name with milestone indicator
// Escapes LaTeX special characters and adds milestone star if applicable
func (d Day) prepareTaskName(task *SpanningTask) string {
	nameText := d.escapeLatexSpecialChars(strings.TrimSpace(task.Name))
	if d.isMilestoneSpanningTask(task) {
		nameText = "★ " + nameText
	}
	return nameText
}

// truncateTaskName truncates task name based on total number of tasks
// Uses progressive truncation: more tasks = shorter text per task
func (d Day) truncateTaskName(nameText string, total int) string {
	config := getDefaultTaskRenderingConfig()

	// Progressive truncation based on number of tasks
	maxChars := config.MaxChars
	if total > 2 {
		maxChars = config.MaxCharsCompact
	}
	if total > 3 {
		maxChars = config.MaxCharsVeryCompact
	}

	// Apply truncation if needed
	if len(nameText) > maxChars {
		nameText = d.smartTruncateText(nameText, maxChars)
	}
	return nameText
}

// getTaskSpacingAndHeight returns spacing and height based on task index
// Uses configuration to ensure consistent spacing and readability
func (d Day) getTaskSpacingAndHeight(index int) (string, string) {
	config := getDefaultTaskRenderingConfig()

	// First task gets special treatment for better visual hierarchy
	if index == 0 {
		return config.FirstTaskSpacing, config.FirstTaskHeight
	}

	// Subsequent tasks use default spacing and height
	return config.DefaultSpacing, config.DefaultHeight
}

// buildTaskTextBody creates the text body for a task
func (d Day) buildTaskTextBody(nameText string) string {
	// * Use improved text formatting with better line breaking and left alignment
	return `{\sloppy\hyphenpenalty=50\tolerance=1000\emergencystretch=2em\color{black}\TaskFontSize\raggedright\textbf{` + nameText + `}}`
}

// buildCompactTaskBox creates the tcolorbox for a compact task
func (d Day) buildCompactTaskBox(spacing, boxHeight, color, textBody string) string {
	// Use macro wrapper for compact bar
	return `\TaskCompactBox{` + spacing + `}{` + boxHeight + `}{` + color + `}{` + textBody + `}`
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
func (le *LayoutEngine) ProcessTasksWithValidation(tasks []*common.Task) (*MultiDayLayoutResult, error) {
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
		TaskBars:         processedBars,
		ValidationResult: validationResult,
		LayoutIssues:     layoutIssues,
		TaskCount:        len(validTasks),
		ProcessedCount:   len(processedBars),
	}, nil
}

// filterValidTasks filters out tasks with critical validation errors
func (le *LayoutEngine) filterValidTasks(tasks []*common.Task, validationErrors []common.DataValidationError) []*common.Task {
	// Create map of tasks with critical errors
	errorTasks := make(map[string]bool)
	for _, err := range validationErrors {
		if err.Severity == "ERROR" {
			errorTasks[err.TaskID] = true
		}
	}

	// Filter out tasks with critical errors
	var validTasks []*common.Task
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
		TotalTasks:         result.TaskCount,
		ProcessedBars:      result.ProcessedCount,
		ValidationErrors:   len(result.ValidationResult),
		LayoutIssues:       len(result.LayoutIssues),
		OverlapCount:       0,
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
// SpatialEngine handles both overlap detection and positioning of tasks within the calendar grid
type SpatialEngine struct {
	// Overlap detection components
	calendarStart time.Time
	calendarEnd   time.Time
	precision     time.Duration

	// Positioning components
	gridConfig        *GridConfig
	visualConstraints *VisualConstraints
	alignmentRules    []AlignmentRule
	spacingRules      []SpacingRule
	layoutMetrics     *PositioningLayoutMetrics
}

// OverlapType represents the type of overlap between two tasks
type OverlapType string

const (
	OverlapNone      OverlapType = "NONE"
	OverlapPartial   OverlapType = "PARTIAL"
	OverlapComplete  OverlapType = "COMPLETE"
	OverlapNested    OverlapType = "NESTED"
	OverlapAdjacent  OverlapType = "ADJACENT"
	OverlapIdentical OverlapType = "IDENTICAL"
)

// OverlapSeverity represents the severity level of an overlap
type OverlapSeverity string

const (
	SeverityNone     OverlapSeverity = "NONE"
	SeverityLow      OverlapSeverity = "LOW"
	SeverityMedium   OverlapSeverity = "MEDIUM"
	SeverityHigh     OverlapSeverity = "HIGH"
	SeverityCritical OverlapSeverity = "CRITICAL"
)

// TaskOverlap represents an overlap between two tasks
type TaskOverlap struct {
	Task1ID        string
	Task2ID        string
	OverlapType    OverlapType
	Severity       OverlapSeverity
	StartDate      time.Time
	EndDate        time.Time
	Duration       time.Duration
	OverlapDays    int
	ConflictReason string
	ResolutionHint string
	Priority       int
}

// OverlapGroup represents a group of overlapping tasks
type OverlapGroup struct {
	Tasks         []*common.Task
	Overlaps      []*TaskOverlap
	GroupID       string
	StartDate     time.Time
	EndDate       time.Time
	MaxSeverity   OverlapSeverity
	ConflictCount int
	Resolution    string
}

// OverlapAnalysis contains comprehensive overlap analysis results
type OverlapAnalysis struct {
	TotalTasks       int
	OverlappingTasks int
	OverlapGroups    []*OverlapGroup
	TotalOverlaps    int
	CriticalOverlaps int
	HighOverlaps     int
	MediumOverlaps   int
	LowOverlaps      int
	AnalysisDate     time.Time
	Summary          string
}

// AlignmentRule defines how tasks should be aligned within the grid
type AlignmentRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*IntegratedTaskBar, *PositioningContext) bool
	Action      func(*IntegratedTaskBar, *PositioningContext) *PositioningAction
}

// SpacingRule defines spacing rules between tasks
type SpacingRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*IntegratedTaskBar, *IntegratedTaskBar, *PositioningContext) bool
	Action      func(*IntegratedTaskBar, *IntegratedTaskBar, *PositioningContext) *SpacingAction
}

// PositioningContext provides context for positioning decisions
type PositioningContext struct {
	CalendarStart   time.Time
	CalendarEnd     time.Time
	DayWidth        float64
	DayHeight       float64
	AvailableHeight float64
	AvailableWidth  float64
	ExistingBars    []*IntegratedTaskBar
	GridConstraints *GridConstraints
	VisualSettings  *IntegratedVisualSettings
	CurrentTime     time.Time
	TaskDensity     float64
	OverlapCount    int
	ConflictCount   int
}

// GridConstraints defines constraints for grid positioning
type GridConstraints struct {
	MinTaskSpacing     float64
	MaxTaskSpacing     float64
	MinRowHeight       float64
	MaxRowHeight       float64
	MinColumnWidth     float64
	MaxColumnWidth     float64
	SnapToGrid         bool
	GridResolution     float64
	AlignmentTolerance float64
	CollisionBuffer    float64
}

// PositioningAction defines how a task should be positioned
type PositioningAction struct {
	X                float64
	Y                float64
	Width            float64
	Height           float64
	Row              int
	Column           int
	AlignmentMode    PositioningAlignmentMode
	Justification    JustificationMode
	VerticalOffset   float64
	HorizontalOffset float64
	ZIndex           int
	SnapToGrid       bool
	Priority         int
}

// PositioningAlignmentMode defines the alignment mode for tasks
type PositioningAlignmentMode string

const (
	PositioningAlignmentLeft    PositioningAlignmentMode = "LEFT"
	PositioningAlignmentCenter  PositioningAlignmentMode = "CENTER"
	PositioningAlignmentRight   PositioningAlignmentMode = "RIGHT"
	PositioningAlignmentJustify PositioningAlignmentMode = "JUSTIFY"
	PositioningAlignmentStretch PositioningAlignmentMode = "STRETCH"
	PositioningAlignmentTop     PositioningAlignmentMode = "TOP"
	PositioningAlignmentMiddle  PositioningAlignmentMode = "MIDDLE"
	PositioningAlignmentBottom  PositioningAlignmentMode = "BOTTOM"
)

// JustificationMode defines how tasks should be justified within available space
type JustificationMode string

const (
	JustifyStart        JustificationMode = "START"
	JustifyEnd          JustificationMode = "END"
	JustifyCenter       JustificationMode = "CENTER"
	JustifySpaceBetween JustificationMode = "SPACE_BETWEEN"
	JustifySpaceAround  JustificationMode = "SPACE_AROUND"
	JustifySpaceEvenly  JustificationMode = "SPACE_EVENLY"
)

// SpacingAction defines spacing adjustments between tasks
type SpacingAction struct {
	VerticalSpacing    float64
	HorizontalSpacing  float64
	CollisionAvoidance bool
	OverlapResolution  bool
	Priority           int
}

// PositioningLayoutMetrics contains metrics about the layout
type PositioningLayoutMetrics struct {
	TotalTasks      int
	PositionedTasks int
	CollisionCount  int
	OverlapCount    int
	SpaceEfficiency float64
	AlignmentScore  float64
	SpacingScore    float64
	VisualBalance   float64
	GridUtilization float64
	AverageSpacing  float64
	MaxSpacing      float64
	MinSpacing      float64
	AlignmentErrors int
	SpacingErrors   int
}

// PositioningResult contains the result of positioning operations
type PositioningResult struct {
	TaskBars        []*IntegratedTaskBar
	Metrics         *PositioningLayoutMetrics
	Recommendations []string
	AnalysisDate    time.Time
}

// NewSpatialEngine creates a new spatial engine that handles both overlap detection and positioning
func NewSpatialEngine(calendarStart, calendarEnd time.Time, gridConfig *GridConfig) *SpatialEngine {
	engine := &SpatialEngine{
		calendarStart:     calendarStart,
		calendarEnd:       calendarEnd,
		precision:         time.Hour * 1, // 1 hour minimum overlap
		gridConfig:        gridConfig,
		visualConstraints: gridConfig.VisualConstraints,
		alignmentRules:    []AlignmentRule{},
		spacingRules:      []SpacingRule{},
		layoutMetrics:     &PositioningLayoutMetrics{},
	}

	// Add default alignment rules
	engine.addDefaultAlignmentRules()

	// Add default spacing rules
	engine.addDefaultSpacingRules()

	return engine
}

// NewOverlapDetector creates a new overlap detector (for backward compatibility)
func NewOverlapDetector(calendarStart, calendarEnd time.Time) *SpatialEngine {
	return NewSpatialEngine(calendarStart, calendarEnd, nil)
}

// NewPositioningEngine creates a new positioning engine (for backward compatibility)
func NewPositioningEngine(gridConfig *GridConfig) *SpatialEngine {
	return NewSpatialEngine(gridConfig.CalendarStart, gridConfig.CalendarEnd, gridConfig)
}

// calculateLayoutMetrics calculates layout metrics
func (se *SpatialEngine) calculateLayoutMetrics(bars []*IntegratedTaskBar, context *PositioningContext) *PositioningLayoutMetrics {
	metrics := &PositioningLayoutMetrics{
		TotalTasks:      len(bars),
		PositionedTasks: len(bars),
		CollisionCount:  se.countOverlaps(bars),
		OverlapCount:    se.countOverlaps(bars),
	}

	// Calculate space efficiency
	usedSpace := se.calculateUsedSpace(bars)
	totalSpace := context.AvailableWidth * context.AvailableHeight
	metrics.SpaceEfficiency = usedSpace / totalSpace

	// Calculate alignment score
	metrics.AlignmentScore = se.calculateAlignmentScore(bars, context)

	// Calculate spacing score
	metrics.SpacingScore = se.calculateSpacingScore(bars, context)

	// Calculate visual balance
	metrics.VisualBalance = se.calculateVisualBalance(bars, context)

	// Calculate grid utilization
	metrics.GridUtilization = se.calculateGridUtilization(bars, context)

	// Calculate average spacing
	metrics.AverageSpacing = se.calculateAverageSpacing(bars, context)

	return metrics
}

// generatePositioningRecommendations generates recommendations based on layout metrics
func (se *SpatialEngine) generatePositioningRecommendations(metrics *PositioningLayoutMetrics, context *PositioningContext) []string {
	var recommendations []string

	// Space efficiency recommendations
	if metrics.SpaceEfficiency < 0.7 {
		recommendations = append(recommendations, "Consider reducing task spacing to improve space efficiency")
	}

	// Alignment recommendations
	if metrics.AlignmentScore < 0.8 {
		recommendations = append(recommendations, "Enable grid snapping to improve alignment consistency")
	}

	// Spacing recommendations
	if metrics.SpacingScore < 0.7 {
		recommendations = append(recommendations, "Adjust spacing rules to improve visual consistency")
	}

	// Visual balance recommendations
	if metrics.VisualBalance < 0.6 {
		recommendations = append(recommendations, "Redistribute tasks to improve visual balance")
	}

	// Collision recommendations
	if metrics.CollisionCount > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Resolve %d collisions to improve layout clarity", metrics.CollisionCount))
	}

	return recommendations
}

// SetPrecision sets the minimum overlap duration to consider
func (se *SpatialEngine) SetPrecision(precision time.Duration) {
	se.precision = precision
}

// DetectOverlaps detects all overlaps in a collection of tasks
func (se *SpatialEngine) DetectOverlaps(tasks []*common.Task) *OverlapAnalysis {
	analysis := &OverlapAnalysis{
		TotalTasks:       len(tasks),
		OverlappingTasks: 0,
		OverlapGroups:    make([]*OverlapGroup, 0),
		AnalysisDate:     time.Now(),
	}

	// Create overlap groups using the existing grouping algorithm
	groups := se.groupOverlappingTasks(tasks)

	// Analyze each group for overlaps
	for i, group := range groups {
		overlapGroup := &OverlapGroup{
			Tasks:         group.Tasks,
			Overlaps:      make([]*TaskOverlap, 0),
			GroupID:       fmt.Sprintf("group_%d", i),
			StartDate:     group.StartDate,
			EndDate:       group.EndDate,
			MaxSeverity:   SeverityNone,
			ConflictCount: 0,
		}

		// Detect overlaps within the group
		overlaps := se.detectGroupOverlaps(group.Tasks)
		overlapGroup.Overlaps = overlaps
		overlapGroup.ConflictCount = len(overlaps)

		// Calculate group statistics
		se.calculateGroupStatistics(overlapGroup)

		// Add to analysis
		analysis.OverlapGroups = append(analysis.OverlapGroups, overlapGroup)
		analysis.TotalOverlaps += len(overlaps)
	}

	// Calculate overall statistics
	se.calculateAnalysisStatistics(analysis)

	// Generate summary
	analysis.Summary = se.generateAnalysisSummary(analysis)

	return analysis
}

// groupOverlappingTasks groups tasks that have any temporal overlap
func (se *SpatialEngine) groupOverlappingTasks(tasks []*common.Task) []*TaskGroup {
	// Sort tasks by start date
	sortedTasks := make([]*common.Task, len(tasks))
	copy(sortedTasks, tasks)
	sort.Slice(sortedTasks, func(i, j int) bool {
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
			Tasks:     []*common.Task{task},
			StartDate: task.StartDate,
			EndDate:   task.EndDate,
		}
		used[task.ID] = true

		// Find all tasks that overlap with this group
		for _, otherTask := range sortedTasks {
			if used[otherTask.ID] {
				continue
			}

			if se.tasksOverlap(group, otherTask) {
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

		groups = append(groups, group)
	}

	return groups
}

// tasksOverlap checks if a task overlaps with a group
func (se *SpatialEngine) tasksOverlap(group *TaskGroup, task *common.Task) bool {
	for _, groupTask := range group.Tasks {
		if se.tasksOverlapDirect(groupTask, task) {
			return true
		}
	}
	return false
}

// tasksOverlapDirect checks if two tasks overlap directly
func (se *SpatialEngine) tasksOverlapDirect(task1, task2 *common.Task) bool {
	// Tasks overlap if one starts before the other ends
	return !task1.StartDate.After(task2.EndDate) && !task2.StartDate.After(task1.EndDate)
}

// detectGroupOverlaps detects all overlaps within a group of tasks
func (se *SpatialEngine) detectGroupOverlaps(tasks []*common.Task) []*TaskOverlap {
	var overlaps []*TaskOverlap

	// Check all pairs of tasks in the group
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			overlap := se.analyzeTaskOverlap(tasks[i], tasks[j])
			if overlap != nil {
				overlaps = append(overlaps, overlap)
			}
		}
	}

	// Sort overlaps by severity and priority
	sort.Slice(overlaps, func(i, j int) bool {
		if overlaps[i].Severity != overlaps[j].Severity {
			return se.severityOrder(overlaps[i].Severity) < se.severityOrder(overlaps[j].Severity)
		}
		return overlaps[i].Priority > overlaps[j].Priority
	})

	return overlaps
}

// analyzeTaskOverlap analyzes the overlap between two specific tasks
func (se *SpatialEngine) analyzeTaskOverlap(task1, task2 *common.Task) *TaskOverlap {
	// Check if tasks actually overlap
	if !se.tasksOverlapDirect(task1, task2) {
		return nil
	}

	// Calculate overlap details
	overlapStart := se.maxTime(task1.StartDate, task2.StartDate)
	overlapEnd := se.minTime(task1.EndDate, task2.EndDate)

	// Check if overlap meets minimum precision requirement
	overlapDuration := overlapEnd.Sub(overlapStart)
	if overlapDuration < se.precision {
		return nil
	}

	// Determine overlap type
	overlapType := se.determineOverlapType(task1, task2, overlapStart, overlapEnd)

	// Calculate severity
	severity := se.calculateOverlapSeverity(task1, task2, overlapType, overlapDuration)

	// Calculate priority (higher priority task wins)
	priority := se.calculateOverlapPriority(task1, task2)

	// Generate conflict reason and resolution hint
	conflictReason, resolutionHint := se.generateConflictInfo(task1, task2, overlapType, severity)

	// Calculate overlap days
	overlapDays := int(overlapDuration.Hours()/24) + 1

	return &TaskOverlap{
		Task1ID:        task1.ID,
		Task2ID:        task2.ID,
		OverlapType:    overlapType,
		Severity:       severity,
		StartDate:      overlapStart,
		EndDate:        overlapEnd,
		Duration:       overlapDuration,
		OverlapDays:    overlapDays,
		ConflictReason: conflictReason,
		ResolutionHint: resolutionHint,
		Priority:       priority,
	}
}

// determineOverlapType determines the type of overlap between two tasks
func (se *SpatialEngine) determineOverlapType(task1, task2 *common.Task, overlapStart, overlapEnd time.Time) OverlapType {
	// Check for identical tasks
	if task1.StartDate.Equal(task2.StartDate) && task1.EndDate.Equal(task2.EndDate) {
		return OverlapIdentical
	}

	// Check for complete overlap (one task completely contains the other)
	if (task1.StartDate.Before(task2.StartDate) || task1.StartDate.Equal(task2.StartDate)) &&
		(task1.EndDate.After(task2.EndDate) || task1.EndDate.Equal(task2.EndDate)) {
		return OverlapNested
	}
	if (task2.StartDate.Before(task1.StartDate) || task2.StartDate.Equal(task1.StartDate)) &&
		(task2.EndDate.After(task1.EndDate) || task2.EndDate.Equal(task1.EndDate)) {
		return OverlapNested
	}

	// Check for adjacent tasks (touching but not overlapping)
	if task1.EndDate.Equal(task2.StartDate) || task2.EndDate.Equal(task1.StartDate) {
		return OverlapAdjacent
	}

	// Check for complete overlap (same start and end dates)
	if overlapStart.Equal(task1.StartDate) && overlapEnd.Equal(task1.EndDate) &&
		overlapStart.Equal(task2.StartDate) && overlapEnd.Equal(task2.EndDate) {
		return OverlapComplete
	}

	// Default to partial overlap
	return OverlapPartial
}

// calculateOverlapSeverity calculates the severity of an overlap
func (se *SpatialEngine) calculateOverlapSeverity(task1, task2 *common.Task, overlapType OverlapType, duration time.Duration) OverlapSeverity {
	// Base severity on overlap type
	switch overlapType {
	case OverlapIdentical:
		return SeverityCritical
	case OverlapNested:
		return SeverityHigh
	case OverlapComplete:
		return SeverityHigh
	case OverlapPartial:
		// Severity based on overlap percentage
		overlapPercentage := se.calculateOverlapPercentage(task1, task2, duration)
		if overlapPercentage >= 0.8 {
			return SeverityHigh
		} else if overlapPercentage >= 0.5 {
			return SeverityMedium
		} else {
			return SeverityLow
		}
	case OverlapAdjacent:
		return SeverityLow
	default:
		return SeverityNone
	}
}

// calculateOverlapPercentage calculates the percentage of overlap
func (se *SpatialEngine) calculateOverlapPercentage(task1, task2 *common.Task, overlapDuration time.Duration) float64 {
	task1Duration := task1.EndDate.Sub(task1.StartDate)
	task2Duration := task2.EndDate.Sub(task2.StartDate)

	// Use the shorter task duration as the base
	baseDuration := task1Duration
	if task2Duration < task1Duration {
		baseDuration = task2Duration
	}

	if baseDuration == 0 {
		return 0.0
	}

	return float64(overlapDuration) / float64(baseDuration)
}

// calculateOverlapPriority calculates priority for overlap resolution
func (se *SpatialEngine) calculateOverlapPriority(task1, task2 *common.Task) int {
	// Higher priority task wins
	if task1.Priority > task2.Priority {
		return task1.Priority
	}
	return task2.Priority
}

// generateConflictInfo generates conflict reason and resolution hint
func (se *SpatialEngine) generateConflictInfo(task1, task2 *common.Task, overlapType OverlapType, severity OverlapSeverity) (string, string) {
	var reason, hint string

	switch overlapType {
	case OverlapIdentical:
		reason = fmt.Sprintf("Tasks %s and %s have identical schedules", task1.ID, task2.ID)
		hint = "Consider merging tasks or adjusting one task's schedule"
	case OverlapNested:
		reason = fmt.Sprintf("Task %s is completely contained within task %s", task1.ID, task2.ID)
		hint = "Consider making the nested task a subtask or adjusting schedules"
	case OverlapComplete:
		reason = fmt.Sprintf("Tasks %s and %s have complete schedule overlap", task1.ID, task2.ID)
		hint = "Tasks cannot run simultaneously - reschedule one task"
	case OverlapPartial:
		reason = fmt.Sprintf("Tasks %s and %s have partial schedule overlap", task1.ID, task2.ID)
		hint = "Consider adjusting start/end times to reduce overlap"
	case OverlapAdjacent:
		reason = fmt.Sprintf("Tasks %s and %s are adjacent in schedule", task1.ID, task2.ID)
		hint = "Consider adding buffer time between tasks"
	default:
		reason = fmt.Sprintf("Tasks %s and %s have unknown overlap type", task1.ID, task2.ID)
		hint = "Review task schedules for potential conflicts"
	}

	// Add severity context
	if severity == SeverityCritical {
		reason += " (CRITICAL)"
		hint = "URGENT: " + hint
	} else if severity == SeverityHigh {
		reason += " (HIGH)"
		hint = "Important: " + hint
	}

	return reason, hint
}

// Helper functions for time comparisons
func (se *SpatialEngine) maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func (se *SpatialEngine) minTime(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}
	return t2
}

// severityOrder returns the order value for severity sorting
func (se *SpatialEngine) severityOrder(severity OverlapSeverity) int {
	switch severity {
	case SeverityCritical:
		return 0
	case SeverityHigh:
		return 1
	case SeverityMedium:
		return 2
	case SeverityLow:
		return 3
	case SeverityNone:
		return 4
	default:
		return 5
	}
}

// calculateGroupStatistics calculates statistics for an overlap group
func (se *SpatialEngine) calculateGroupStatistics(group *OverlapGroup) {
	if len(group.Overlaps) == 0 {
		return
	}

	// Find maximum severity
	for _, overlap := range group.Overlaps {
		if se.severityOrder(overlap.Severity) < se.severityOrder(group.MaxSeverity) {
			group.MaxSeverity = overlap.Severity
		}
	}

	// Generate resolution suggestion
	group.Resolution = se.generateGroupResolution(group)
}

// generateGroupResolution generates resolution suggestions for a group
func (se *SpatialEngine) generateGroupResolution(group *OverlapGroup) string {
	if len(group.Overlaps) == 0 {
		return "No conflicts detected"
	}

	criticalCount := 0
	highCount := 0
	for _, overlap := range group.Overlaps {
		if overlap.Severity == SeverityCritical {
			criticalCount++
		} else if overlap.Severity == SeverityHigh {
			highCount++
		}
	}

	if criticalCount > 0 {
		return fmt.Sprintf("URGENT: %d critical conflicts require immediate attention", criticalCount)
	} else if highCount > 0 {
		return fmt.Sprintf("Important: %d high-priority conflicts need resolution", highCount)
	} else {
		return fmt.Sprintf("Moderate: %d conflicts can be addressed during planning", len(group.Overlaps))
	}
}

// calculateAnalysisStatistics calculates overall analysis statistics
func (se *SpatialEngine) calculateAnalysisStatistics(analysis *OverlapAnalysis) {
	overlappingTaskIDs := make(map[string]bool)

	for _, group := range analysis.OverlapGroups {
		for _, overlap := range group.Overlaps {
			overlappingTaskIDs[overlap.Task1ID] = true
			overlappingTaskIDs[overlap.Task2ID] = true

			switch overlap.Severity {
			case SeverityCritical:
				analysis.CriticalOverlaps++
			case SeverityHigh:
				analysis.HighOverlaps++
			case SeverityMedium:
				analysis.MediumOverlaps++
			case SeverityLow:
				analysis.LowOverlaps++
			}
		}
	}

	analysis.OverlappingTasks = len(overlappingTaskIDs)
}

// generateAnalysisSummary generates a summary of the overlap analysis
func (se *SpatialEngine) generateAnalysisSummary(analysis *OverlapAnalysis) string {
	if analysis.TotalOverlaps == 0 {
		return fmt.Sprintf("✅ No task overlaps detected in %d tasks", analysis.TotalTasks)
	}

	summary := fmt.Sprintf("⚠️  Detected %d overlaps affecting %d tasks:\n",
		analysis.TotalOverlaps, analysis.OverlappingTasks)

	if analysis.CriticalOverlaps > 0 {
		summary += fmt.Sprintf("  🚨 %d critical overlaps (immediate action required)\n", analysis.CriticalOverlaps)
	}
	if analysis.HighOverlaps > 0 {
		summary += fmt.Sprintf("  🔴 %d high-priority overlaps\n", analysis.HighOverlaps)
	}
	if analysis.MediumOverlaps > 0 {
		summary += fmt.Sprintf("  🟡 %d medium-priority overlaps\n", analysis.MediumOverlaps)
	}
	if analysis.LowOverlaps > 0 {
		summary += fmt.Sprintf("  🟢 %d low-priority overlaps\n", analysis.LowOverlaps)
	}

	summary += fmt.Sprintf("  📊 %d overlap groups identified", len(analysis.OverlapGroups))

	return summary
}

// GetOverlapsBySeverity returns overlaps filtered by severity
func (analysis *OverlapAnalysis) GetOverlapsBySeverity(severity OverlapSeverity) []*TaskOverlap {
	var filtered []*TaskOverlap
	for _, group := range analysis.OverlapGroups {
		for _, overlap := range group.Overlaps {
			if overlap.Severity == severity {
				filtered = append(filtered, overlap)
			}
		}
	}
	return filtered
}

// GetOverlapsByType returns overlaps filtered by type
func (analysis *OverlapAnalysis) GetOverlapsByType(overlapType OverlapType) []*TaskOverlap {
	var filtered []*TaskOverlap
	for _, group := range analysis.OverlapGroups {
		for _, overlap := range group.Overlaps {
			if overlap.OverlapType == overlapType {
				filtered = append(filtered, overlap)
			}
		}
	}
	return filtered
}

// HasCriticalOverlaps returns true if there are any critical overlaps
func (analysis *OverlapAnalysis) HasCriticalOverlaps() bool {
	return analysis.CriticalOverlaps > 0
}

// GetOverlapCount returns the total number of overlaps
func (analysis *OverlapAnalysis) GetOverlapCount() int {
	return analysis.TotalOverlaps
}

// GetOverlappingTaskCount returns the number of tasks involved in overlaps
func (analysis *OverlapAnalysis) GetOverlappingTaskCount() int {
	return analysis.OverlappingTasks
}

// PositionTasks positions all tasks within the calendar grid
func (se *SpatialEngine) PositionTasks(tasks []*common.Task, existingBars []*IntegratedTaskBar) (*PositioningResult, error) {
	// Create positioning context
	context := &PositioningContext{
		CalendarStart:   se.gridConfig.CalendarStart,
		CalendarEnd:     se.gridConfig.CalendarEnd,
		DayWidth:        se.gridConfig.DayWidth,
		DayHeight:       se.gridConfig.DayHeight,
		AvailableHeight: se.gridConfig.DayHeight * float64(se.gridConfig.MaxRowsPerDay),
		AvailableWidth:  se.gridConfig.DayWidth * 7.0,
		ExistingBars:    existingBars,
		GridConstraints: se.getDefaultGridConstraints(),
		VisualSettings:  se.getDefaultVisualSettings(),
		CurrentTime:     time.Now(),
		TaskDensity:     se.calculateTaskDensity(tasks),
		OverlapCount:    se.countOverlaps(existingBars),
		ConflictCount:   se.countConflicts(existingBars),
	}

	// Create integrated task bars
	integratedBars := se.createIntegratedTaskBars(tasks, context)

	// Apply positioning rules
	positionedBars := se.applyPositioningRules(integratedBars, context)

	// Apply spacing rules
	spacedBars := se.applySpacingRules(positionedBars, context)

	// Snap to grid if enabled
	if context.GridConstraints.SnapToGrid {
		spacedBars = se.snapToGrid(spacedBars, context)
	}

	// Resolve collisions
	finalBars := se.resolveCollisions(spacedBars, context)

	// Calculate metrics
	metrics := se.calculateLayoutMetrics(finalBars, context)

	// Generate recommendations
	recommendations := se.generatePositioningRecommendations(metrics, context)

	return &PositioningResult{
		TaskBars:        finalBars,
		Metrics:         metrics,
		Recommendations: recommendations,
		AnalysisDate:    time.Now(),
	}, nil
}

// createIntegratedTaskBars creates integrated task bars from tasks
func (se *SpatialEngine) createIntegratedTaskBars(tasks []*common.Task, context *PositioningContext) []*IntegratedTaskBar {
	var bars []*IntegratedTaskBar

	for _, task := range tasks {
		// Calculate basic positioning
		startX := se.calculateXPosition(task.StartDate, context)
		endX := se.calculateXPosition(task.EndDate, context)
		width := endX - startX

		// Calculate Y position (will be refined by positioning rules)
		y := se.calculateInitialYPosition(task, context)

		// Calculate height
		height := se.calculateTaskHeight(task, context)

		// Get task category and color
		category := common.GetCategory(task.Category)

		// Create integrated task bar
		bar := &IntegratedTaskBar{
			TaskID:          task.ID,
			StartDate:       task.StartDate,
			EndDate:         task.EndDate,
			StartX:          startX,
			EndX:            endX,
			Y:               y,
			Width:           width,
			Height:          height,
			Row:             0,
			StackIndex:      0,
			Color:           category.Color,
			BorderColor:     "#000000",
			Opacity:         0.9,
			ZIndex:          category.Priority,
			IsContinuation:  se.isTaskContinuation(task, context),
			IsStart:         se.isTaskStart(task, context),
			IsEnd:           se.isTaskEnd(task, context),
			MonthBoundary:   se.hasMonthBoundary(task, context),
			StackingType:    StackingTypeVertical,
			VisualWeight:    0.5,
			ProminenceScore: 0.5,
			IsCollapsed:     false,
			IsVisible:       true,
			CollisionLevel:  0,
			OverflowLevel:   0,
			Priority:        task.Priority,
			Category:        task.Category,
			TaskName:        task.Name,
			Description:     task.Description,
		}

		bars = append(bars, bar)
	}

	return bars
}

// calculateXPosition calculates the X position for a given date
func (se *SpatialEngine) calculateXPosition(date time.Time, context *PositioningContext) float64 {
	daysFromStart := int(date.Sub(context.CalendarStart).Hours() / 24)
	return float64(daysFromStart) * context.DayWidth
}

// calculateInitialYPosition calculates the initial Y position for a task
func (se *SpatialEngine) calculateInitialYPosition(task *common.Task, context *PositioningContext) float64 {
	// Start with a basic Y position based on task priority
	baseY := context.DayHeight * 0.1 // 10% from top

	// Adjust based on task priority
	priorityOffset := float64(task.Priority) * context.DayHeight * 0.1

	return baseY + priorityOffset
}

// calculateTaskHeight calculates the height of a task
func (se *SpatialEngine) calculateTaskHeight(task *common.Task, context *PositioningContext) float64 {
	// Base height
	height := context.DayHeight * 0.6 // 60% of day height

	// Adjust based on task duration
	duration := task.EndDate.Sub(task.StartDate).Hours() / 24
	if duration > 7 {
		height *= 1.2 // Longer tasks get more height
	} else if duration < 1 {
		height *= 0.8 // Shorter tasks get less height
	}

	// Ensure within constraints
	if height < context.GridConstraints.MinRowHeight {
		height = context.GridConstraints.MinRowHeight
	} else if height > context.GridConstraints.MaxRowHeight {
		height = context.GridConstraints.MaxRowHeight
	}

	return height
}

// Helper methods for task positioning
func (se *SpatialEngine) isTaskContinuation(task *common.Task, context *PositioningContext) bool {
	return task.StartDate.Before(context.CalendarStart)
}

func (se *SpatialEngine) isTaskStart(task *common.Task, context *PositioningContext) bool {
	return task.StartDate.Equal(context.CalendarStart) || task.StartDate.After(context.CalendarStart)
}

func (se *SpatialEngine) isTaskEnd(task *common.Task, context *PositioningContext) bool {
	return task.EndDate.Equal(context.CalendarEnd) || task.EndDate.Before(context.CalendarEnd)
}

func (se *SpatialEngine) hasMonthBoundary(task *common.Task, context *PositioningContext) bool {
	startMonth := task.StartDate.Month()
	endMonth := task.EndDate.Month()
	return startMonth != endMonth
}

// Default configuration methods
func (se *SpatialEngine) getDefaultGridConstraints() *GridConstraints {
	return &GridConstraints{
		MinTaskSpacing:     1.0,
		MaxTaskSpacing:     10.0,
		MinRowHeight:       8.0,
		MaxRowHeight:       20.0,
		MinColumnWidth:     5.0,
		MaxColumnWidth:     50.0,
		SnapToGrid:         true,
		GridResolution:     1.0,
		AlignmentTolerance: 0.5,
		CollisionBuffer:    2.0,
	}
}

func (se *SpatialEngine) getDefaultVisualSettings() *IntegratedVisualSettings {
	return &IntegratedVisualSettings{
		ShowTaskNames:          true,
		ShowTaskDurations:      true,
		ShowTaskPriorities:     true,
		ShowConflictIndicators: true,
		CollapseThreshold:      5,
		AnimationEnabled:       false,
        HighlightConflicts:     false,
		ColorScheme:            "default",
		FontSize:               "small",
		TaskBarOpacity:         0.9,
		BorderWidth:            0.5,
	}
}

// applyPositioningRules applies alignment rules to task bars
func (se *SpatialEngine) applyPositioningRules(bars []*IntegratedTaskBar, context *PositioningContext) []*IntegratedTaskBar {
	// Sort rules by priority
	sort.Slice(se.alignmentRules, func(i, j int) bool {
		return se.alignmentRules[i].Priority > se.alignmentRules[j].Priority
	})

	// Apply rules to each bar
	for _, bar := range bars {
		for _, rule := range se.alignmentRules {
			if rule.Condition(bar, context) {
				action := rule.Action(bar, context)
				se.applyPositioningAction(bar, action)
				break // Apply only the first matching rule
			}
		}
	}

	return bars
}

// applySpacingRules applies spacing rules between task bars
func (se *SpatialEngine) applySpacingRules(bars []*IntegratedTaskBar, context *PositioningContext) []*IntegratedTaskBar {
	// Sort rules by priority
	sort.Slice(se.spacingRules, func(i, j int) bool {
		return se.spacingRules[i].Priority > se.spacingRules[j].Priority
	})

	// Apply spacing rules between adjacent bars
	for i := 0; i < len(bars)-1; i++ {
		for j := i + 1; j < len(bars); j++ {
			for _, rule := range se.spacingRules {
				if rule.Condition(bars[i], bars[j], context) {
					action := rule.Action(bars[i], bars[j], context)
					se.applySpacingAction(bars[i], bars[j], action)
					break // Apply only the first matching rule
				}
			}
		}
	}

	return bars
}

// applyPositioningAction applies a positioning action to a task bar
func (se *SpatialEngine) applyPositioningAction(bar *IntegratedTaskBar, action *PositioningAction) {
	if action.X > 0 {
		bar.StartX = action.X
		bar.EndX = action.X + bar.Width
	}

	if action.Y > 0 {
		bar.Y = action.Y
	}

	if action.Width > 0 {
		bar.Width = action.Width
		bar.EndX = bar.StartX + bar.Width
	}

	if action.Height > 0 {
		bar.Height = action.Height
	}

	if action.Row >= 0 {
		bar.Row = action.Row
	}

	if action.ZIndex > 0 {
		bar.ZIndex = action.ZIndex
	}

	// Apply offsets
	bar.StartX += action.HorizontalOffset
	bar.EndX += action.HorizontalOffset
	bar.Y += action.VerticalOffset
}

// applySpacingAction applies a spacing action between two task bars
func (se *SpatialEngine) applySpacingAction(bar1, bar2 *IntegratedTaskBar, action *SpacingAction) {
	// Calculate distance between bars
	distance := se.calculateDistance(bar1, bar2)

	// Apply vertical spacing
	if action.VerticalSpacing > 0 && distance < action.VerticalSpacing {
		adjustment := action.VerticalSpacing - distance
		bar2.Y += adjustment
	}

	// Apply horizontal spacing
	if action.HorizontalSpacing > 0 && distance < action.HorizontalSpacing {
		adjustment := action.HorizontalSpacing - distance
		bar2.StartX += adjustment
		bar2.EndX += adjustment
	}
}

// snapToGrid snaps task bars to grid positions
func (se *SpatialEngine) snapToGrid(bars []*IntegratedTaskBar, context *PositioningContext) []*IntegratedTaskBar {
	resolution := context.GridConstraints.GridResolution

	for _, bar := range bars {
		// Snap X position
		bar.StartX = math.Round(bar.StartX/resolution) * resolution
		bar.EndX = math.Round(bar.EndX/resolution) * resolution
		bar.Width = bar.EndX - bar.StartX

		// Snap Y position
		bar.Y = math.Round(bar.Y/resolution) * resolution

		// Snap height
		bar.Height = math.Round(bar.Height/resolution) * resolution
	}

	return bars
}

// resolveCollisions resolves collisions between task bars
func (se *SpatialEngine) resolveCollisions(bars []*IntegratedTaskBar, context *PositioningContext) []*IntegratedTaskBar {
	// Sort bars by priority (higher priority first)
	sort.Slice(bars, func(i, j int) bool {
		return bars[i].Priority > bars[j].Priority
	})

	// Resolve collisions
	for i := 0; i < len(bars); i++ {
		for j := i + 1; j < len(bars); j++ {
			if se.barsCollide(bars[i], bars[j], context) {
				se.resolveCollision(bars[i], bars[j], context)
			}
		}
	}

	return bars
}

// barsCollide checks if two task bars collide
func (se *SpatialEngine) barsCollide(bar1, bar2 *IntegratedTaskBar, context *PositioningContext) bool {
	buffer := context.GridConstraints.CollisionBuffer

	// Check horizontal overlap
	horizontalOverlap := bar1.StartX < bar2.EndX+buffer && bar2.StartX < bar1.EndX+buffer

	// Check vertical overlap
	verticalOverlap := bar1.Y < bar2.Y+bar2.Height+buffer && bar2.Y < bar1.Y+bar1.Height+buffer

	return horizontalOverlap && verticalOverlap
}

// resolveCollision resolves a collision between two task bars
func (se *SpatialEngine) resolveCollision(bar1, bar2 *IntegratedTaskBar, context *PositioningContext) {
	// Move the lower priority bar
	if bar1.Priority > bar2.Priority {
		// Move bar2
		bar2.Y = bar1.Y + bar1.Height + context.GridConstraints.CollisionBuffer
	} else {
		// Move bar1
		bar1.Y = bar2.Y + bar2.Height + context.GridConstraints.CollisionBuffer
	}
}

// calculateDistance calculates the distance between two task bars
func (se *SpatialEngine) calculateDistance(bar1, bar2 *IntegratedTaskBar) float64 {
	// Calculate center points
	center1X := bar1.StartX + bar1.Width/2
	center1Y := bar1.Y + bar1.Height/2
	center2X := bar2.StartX + bar2.Width/2
	center2Y := bar2.Y + bar2.Height/2

	// Calculate Euclidean distance
	dx := center2X - center1X
	dy := center2Y - center1Y
	return math.Sqrt(dx*dx + dy*dy)
}

// calculateTaskDensity calculates the task density in the calendar
func (se *SpatialEngine) calculateTaskDensity(tasks []*common.Task) float64 {
	if len(tasks) == 0 {
		return 0.0
	}

	// Calculate total calendar area
	totalArea := se.gridConfig.DayWidth * se.gridConfig.DayHeight * 7.0 * 4.0 // 7 days, 4 weeks

	// Calculate average task area
	var totalTaskArea float64
	for _, task := range tasks {
		duration := task.EndDate.Sub(task.StartDate).Hours() / 24
		taskArea := duration * se.gridConfig.DayWidth * se.gridConfig.DayHeight * 0.6
		totalTaskArea += taskArea
	}

	avgTaskArea := totalTaskArea / float64(len(tasks))

	// Calculate density
	return (avgTaskArea * float64(len(tasks))) / totalArea
}

// countOverlaps counts the number of overlapping task bars
func (se *SpatialEngine) countOverlaps(bars []*IntegratedTaskBar) int {
	count := 0
	for i := 0; i < len(bars); i++ {
		for j := i + 1; j < len(bars); j++ {
			if se.barsOverlap(bars[i], bars[j]) {
				count++
			}
		}
	}
	return count
}

// countConflicts counts the number of conflicts between task bars
func (se *SpatialEngine) countConflicts(bars []*IntegratedTaskBar) int {
	// For now, use overlap count as conflict count
	return se.countOverlaps(bars)
}

// barsOverlap checks if two task bars overlap
func (se *SpatialEngine) barsOverlap(bar1, bar2 *IntegratedTaskBar) bool {
	// Check horizontal overlap
	horizontalOverlap := bar1.StartX < bar2.EndX && bar2.StartX < bar1.EndX

	// Check vertical overlap
	verticalOverlap := bar1.Y < bar2.Y+bar2.Height && bar2.Y < bar1.Y+bar1.Height

	return horizontalOverlap && verticalOverlap
}

// Add default alignment and spacing rules
func (se *SpatialEngine) addDefaultAlignmentRules() {
	// High priority tasks alignment rule
	se.alignmentRules = append(se.alignmentRules, AlignmentRule{
		Name:        "High Priority Alignment",
		Description: "Align high priority tasks to the top",
		Priority:    10,
		Condition: func(bar *IntegratedTaskBar, context *PositioningContext) bool {
			return bar.Priority >= 4
		},
		Action: func(bar *IntegratedTaskBar, context *PositioningContext) *PositioningAction {
			return &PositioningAction{
				Y:             context.DayHeight * 0.1,
				AlignmentMode: PositioningAlignmentTop,
				Priority:      10,
			}
		},
	})

	// Milestone tasks alignment rule
	se.alignmentRules = append(se.alignmentRules, AlignmentRule{
		Name:        "Milestone Alignment",
		Description: "Center milestone tasks vertically",
		Priority:    8,
		Condition: func(bar *IntegratedTaskBar, context *PositioningContext) bool {
			return bar.Category == "MILESTONE" ||
				(bar.TaskName != "" && len(bar.TaskName) > 10 &&
					bar.TaskName[:10] == "MILESTONE:")
		},
		Action: func(bar *IntegratedTaskBar, context *PositioningContext) *PositioningAction {
			return &PositioningAction{
				Y:             context.DayHeight * 0.4,
				AlignmentMode: PositioningAlignmentCenter,
				Priority:      8,
			}
		},
	})

	// Default alignment rule
	se.alignmentRules = append(se.alignmentRules, AlignmentRule{
		Name:        "Default Alignment",
		Description: "Default alignment for all tasks",
		Priority:    1,
		Condition: func(bar *IntegratedTaskBar, context *PositioningContext) bool {
			return true
		},
		Action: func(bar *IntegratedTaskBar, context *PositioningContext) *PositioningAction {
			return &PositioningAction{
				Y:             context.DayHeight * 0.2,
				AlignmentMode: PositioningAlignmentLeft,
				Priority:      1,
			}
		},
	})
}

func (se *SpatialEngine) addDefaultSpacingRules() {
	// High priority spacing rule
	se.spacingRules = append(se.spacingRules, SpacingRule{
		Name:        "High Priority Spacing",
		Description: "Extra spacing around high priority tasks",
		Priority:    10,
		Condition: func(bar1, bar2 *IntegratedTaskBar, context *PositioningContext) bool {
			return bar1.Priority >= 4 || bar2.Priority >= 4
		},
		Action: func(bar1, bar2 *IntegratedTaskBar, context *PositioningContext) *SpacingAction {
			return &SpacingAction{
				VerticalSpacing:   3.0,
				HorizontalSpacing: 2.0,
				Priority:          10,
			}
		},
	})

	// Default spacing rule
	se.spacingRules = append(se.spacingRules, SpacingRule{
		Name:        "Default Spacing",
		Description: "Default spacing between tasks",
		Priority:    1,
		Condition: func(bar1, bar2 *IntegratedTaskBar, context *PositioningContext) bool {
			return true
		},
		Action: func(bar1, bar2 *IntegratedTaskBar, context *PositioningContext) *SpacingAction {
			return &SpacingAction{
				VerticalSpacing:   1.0,
				HorizontalSpacing: 0.5,
				Priority:          1,
			}
		},
	})
}

// calculateUsedSpace calculates the total space used by task bars
func (se *SpatialEngine) calculateUsedSpace(bars []*IntegratedTaskBar) float64 {
	var totalSpace float64
	for _, bar := range bars {
		totalSpace += bar.Width * bar.Height
	}
	return totalSpace
}

// calculateAlignmentScore calculates the alignment score
func (se *SpatialEngine) calculateAlignmentScore(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
	if len(bars) == 0 {
		return 1.0
	}

	// Calculate alignment consistency
	var alignmentScore float64
	for _, bar := range bars {
		// Check if bar is aligned to grid
		if context.GridConstraints.SnapToGrid {
			xAligned := math.Mod(bar.StartX, context.GridConstraints.GridResolution) < context.GridConstraints.AlignmentTolerance
			yAligned := math.Mod(bar.Y, context.GridConstraints.GridResolution) < context.GridConstraints.AlignmentTolerance
			if xAligned && yAligned {
				alignmentScore += 1.0
			}
		}
	}

	return alignmentScore / float64(len(bars))
}

// calculateSpacingScore calculates the spacing score
func (se *SpatialEngine) calculateSpacingScore(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
	if len(bars) < 2 {
		return 1.0
	}

	var spacingScore float64
	pairCount := 0

	for i := 0; i < len(bars)-1; i++ {
		for j := i + 1; j < len(bars); j++ {
			distance := se.calculateDistance(bars[i], bars[j])
			minSpacing := context.GridConstraints.MinTaskSpacing
			maxSpacing := context.GridConstraints.MaxTaskSpacing

			if distance >= minSpacing && distance <= maxSpacing {
				spacingScore += 1.0
			}
			pairCount++
		}
	}

	return spacingScore / float64(pairCount)
}

// calculateVisualBalance calculates the visual balance score
func (se *SpatialEngine) calculateVisualBalance(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
	if len(bars) == 0 {
		return 1.0
	}

	// Calculate center of mass
	var totalX, totalY, totalWeight float64
	for _, bar := range bars {
		weight := bar.VisualWeight
		centerX := bar.StartX + bar.Width/2
		centerY := bar.Y + bar.Height/2

		totalX += centerX * weight
		totalY += centerY * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 1.0
	}

	centerOfMassX := totalX / totalWeight
	centerOfMassY := totalY / totalWeight

	// Calculate distance from center of available space
	availableCenterX := context.AvailableWidth / 2
	availableCenterY := context.AvailableHeight / 2

	distanceFromCenter := math.Sqrt(
		math.Pow(centerOfMassX-availableCenterX, 2) +
			math.Pow(centerOfMassY-availableCenterY, 2),
	)

	// Normalize to 0-1 scale
	maxDistance := math.Sqrt(
		math.Pow(context.AvailableWidth/2, 2) +
			math.Pow(context.AvailableHeight/2, 2),
	)

	return 1.0 - (distanceFromCenter / maxDistance)
}

// calculateGridUtilization calculates the grid utilization percentage
func (se *SpatialEngine) calculateGridUtilization(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
	if len(bars) == 0 {
		return 0.0
	}

	// Calculate total grid cells
	gridCellsX := int(context.AvailableWidth / context.GridConstraints.GridResolution)
	gridCellsY := int(context.AvailableHeight / context.GridConstraints.GridResolution)
	totalCells := gridCellsX * gridCellsY

	// Calculate occupied cells
	occupiedCells := make(map[string]bool)
	for _, bar := range bars {
		startCellX := int(bar.StartX / context.GridConstraints.GridResolution)
		endCellX := int(bar.EndX / context.GridConstraints.GridResolution)
		startCellY := int(bar.Y / context.GridConstraints.GridResolution)
		endCellY := int((bar.Y + bar.Height) / context.GridConstraints.GridResolution)

		for x := startCellX; x < endCellX; x++ {
			for y := startCellY; y < endCellY; y++ {
				cellKey := fmt.Sprintf("%d,%d", x, y)
				occupiedCells[cellKey] = true
			}
		}
	}

	return float64(len(occupiedCells)) / float64(totalCells)
}

// calculateAverageSpacing calculates the average spacing between task bars
func (se *SpatialEngine) calculateAverageSpacing(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
	if len(bars) < 2 {
		return 0.0
	}

	var totalSpacing float64
	pairCount := 0

	for i := 0; i < len(bars)-1; i++ {
		for j := i + 1; j < len(bars); j++ {
			distance := se.calculateDistance(bars[i], bars[j])
			totalSpacing += distance
			pairCount++
		}
	}

	return totalSpacing / float64(pairCount)
}

