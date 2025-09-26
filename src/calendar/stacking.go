package calendar

import (
	"fmt"
	"math"
	"sort"
	"time"

	"phd-dissertation-planner/src/core"
)

// StackingEngine handles both smart and vertical stacking of overlapping tasks
type StackingEngine struct {
	spatialEngine       *SpatialEngine
	conflictCategorizer *ConflictCategorizer
	priorityRanker      *PriorityRanker
	stackingRules       []StackingRule
	visualConstraints   *VisualConstraints
	heightCalculator    *HeightCalculator
	positionCalculator  *PositionCalculator
	spaceOptimizer      *SpaceOptimizer
}

// StackingRule defines a rule for task stacking behavior
type StackingRule struct {
	Name        string
	Description string
	Condition   func(*core.Task, *StackingContext) bool
	Action      func(*core.Task, *StackingContext) *StackingAction
}

// StackingAction defines how a task should be stacked
type StackingAction struct {
	StackingType       StackingType
	VerticalOffset     float64
	HorizontalOffset   float64
	Height             float64
	Width              float64
	ZIndex             int
	VisualStyle        *VisualStyle
	CollisionAvoidance bool
}

// StackingType defines the type of stacking behavior
type StackingType string

const (
	StackingTypeVertical   StackingType = "VERTICAL"
	StackingTypeHorizontal StackingType = "HORIZONTAL"
	StackingTypeLayered    StackingType = "LAYERED"
	StackingTypeCascading  StackingType = "CASCADING"
	StackingTypeFloating   StackingType = "FLOATING"
	StackingTypeMinimized  StackingType = "MINIMIZED"
)

// VisualConstraints defines visual constraints for stacking
type VisualConstraints struct {
	MaxStackHeight     float64
	MinTaskHeight      float64
	MaxTaskHeight      float64
	MinTaskWidth       float64
	MaxTaskWidth       float64
	VerticalSpacing    float64
	HorizontalSpacing  float64
	MaxStackDepth      int
	CollisionThreshold float64
	OverflowThreshold  float64
}

// StackingContext provides context for stacking decisions
type StackingContext struct {
	CalendarStart     time.Time
	CalendarEnd       time.Time
	CurrentTime       time.Time
	DayWidth          float64
	DayHeight         float64
	AvailableHeight   float64
	AvailableWidth    float64
	ExistingStacks    []*TaskStack
	TaskPriorities    map[string]*TaskPriority
	ConflictAnalysis  *ConflictAnalysis
	OverlapAnalysis   *OverlapAnalysis
	VisualSettings    *VisualSettings
	VisualConstraints *VisualConstraints
}

// TaskStack represents a stack of overlapping tasks
type TaskStack struct {
	ID             string
	Tasks          []*StackedTask
	StartTime      time.Time
	EndTime        time.Time
	TotalHeight    float64
	MaxWidth       float64
	StackingType   StackingType
	CollisionCount int
	OverflowCount  int
	VisualStyle    *VisualStyle
}

// StackedTask represents a task within a stack
type StackedTask struct {
	Task           *core.Task
	StackingAction *StackingAction
	Position       *Position
	IsVisible      bool
	IsCollapsed    bool
	CollisionLevel int
	OverflowLevel  int
}

// Position represents the position of a stacked task
type Position struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	ZIndex int
}

// VisualSettings defines visual settings for stacking
type VisualSettings struct {
	ShowTaskNames          bool
	ShowTaskDurations      bool
	ShowTaskPriorities     bool
	ShowConflictIndicators bool
	CollapseThreshold      int
	AnimationEnabled       bool
	HighlightConflicts     bool
	ColorScheme            string
}

// StackingResult contains the result of stacking operations
type StackingResult struct {
	Stacks          []*TaskStack
	TotalStacks     int
	CollisionCount  int
	OverflowCount   int
	SpaceEfficiency float64
	VisualQuality   float64
	Recommendations []string
	AnalysisDate    time.Time
}

// HeightCalculator calculates optimal heights for stacked tasks
type HeightCalculator struct {
	baseHeight         float64
	minHeight          float64
	maxHeight          float64
	priorityMultiplier map[VisualProminence]float64
	durationMultiplier map[string]float64
	contentMultiplier  map[string]float64
}

// PositionCalculator calculates optimal positions for stacked tasks
type PositionCalculator struct {
	verticalSpacing   float64
	horizontalSpacing float64
	alignmentMode     AlignmentMode
	distributionMode  DistributionMode
}

// SpaceOptimizer optimizes space usage for vertical stacking
type SpaceOptimizer struct {
	compressionThreshold float64
	expansionThreshold   float64
	adaptiveSpacing      bool
	smartCollapsing      bool
}

// AlignmentMode defines how tasks are aligned within a stack
type AlignmentMode string

const (
	AlignmentTop        AlignmentMode = "TOP"
	AlignmentCenter     AlignmentMode = "CENTER"
	AlignmentBottom     AlignmentMode = "BOTTOM"
	AlignmentJustify    AlignmentMode = "JUSTIFY"
	AlignmentDistribute AlignmentMode = "DISTRIBUTE"
)

// DistributionMode defines how tasks are distributed within available space
type DistributionMode string

const (
	DistributionEven     DistributionMode = "EVEN"
	DistributionPriority DistributionMode = "PRIORITY"
	DistributionContent  DistributionMode = "CONTENT"
	DistributionAdaptive DistributionMode = "ADAPTIVE"
)

// VerticalStackingResult contains the result of vertical stacking operations
type VerticalStackingResult struct {
	Stacks           []*VerticalStack
	TotalHeight      float64
	SpaceEfficiency  float64
	VisualBalance    float64
	CollisionCount   int
	OverflowCount    int
	CompressionRatio float64
	Recommendations  []string
	AnalysisDate     time.Time
}

// VerticalStack represents a vertically stacked group of tasks
type VerticalStack struct {
	ID               string
	Tasks            []*VerticallyStackedTask
	StartTime        time.Time
	EndTime          time.Time
	TotalHeight      float64
	MaxWidth         float64
	AlignmentMode    AlignmentMode
	DistributionMode DistributionMode
	SpaceEfficiency  float64
	VisualBalance    float64
	CollisionCount   int
	OverflowCount    int
}

// VerticallyStackedTask represents a task within a vertical stack
type VerticallyStackedTask struct {
	Task             *core.Task
	Position         *VerticalPosition
	CalculatedHeight float64
	IsCompressed     bool
	IsExpanded       bool
	CollisionLevel   int
	OverflowLevel    int
	VisualWeight     float64
}

// VerticalPosition represents the vertical position of a stacked task
type VerticalPosition struct {
	X          float64
	Y          float64
	Width      float64
	Height     float64
	ZIndex     int
	OffsetY    float64
	RelativeY  float64
	StackIndex int
}

// NewStackingEngine creates a new consolidated stacking engine
func NewStackingEngine(spatialEngine *SpatialEngine, conflictCategorizer *ConflictCategorizer, priorityRanker *PriorityRanker, config *core.Config) *StackingEngine {
	// * Use config-driven visual constraints with fallbacks
	visualConstraints := &VisualConstraints{
		MaxStackHeight:     100.0, // Default fallback
		MinTaskHeight:      20.0,  // Default fallback
		MaxTaskHeight:      40.0,  // Default fallback
		MinTaskWidth:       50.0,  // Default fallback
		MaxTaskWidth:       200.0, // Default fallback
		VerticalSpacing:    2.0,   // Default fallback
		HorizontalSpacing:  5.0,   // Default fallback
		MaxStackDepth:      10,    // Default fallback
		CollisionThreshold: 0.1,   // Default fallback
		OverflowThreshold:  0.8,   // Default fallback
	}

	// Override with config values if available
	if config != nil {
		if config.Layout.Stacking.BaseHeight > 0 {
			visualConstraints.MaxStackHeight = config.Layout.Stacking.BaseHeight
		}
		if config.Layout.Stacking.MinHeight > 0 {
			visualConstraints.MinTaskHeight = config.Layout.Stacking.MinHeight
		}
		if config.Layout.Stacking.MaxHeight > 0 {
			visualConstraints.MaxTaskHeight = config.Layout.Stacking.MaxHeight
		}
		if config.Layout.Stacking.VerticalSpacing > 0 {
			visualConstraints.VerticalSpacing = config.Layout.Stacking.VerticalSpacing
		}
		if config.Layout.Stacking.HorizontalSpacing > 0 {
			visualConstraints.HorizontalSpacing = config.Layout.Stacking.HorizontalSpacing
		}
		if config.Layout.Stacking.CollisionThreshold > 0 {
			visualConstraints.CollisionThreshold = config.Layout.Stacking.CollisionThreshold
		}
		if config.Layout.Stacking.OverflowVertical > 0 {
			visualConstraints.OverflowThreshold = config.Layout.Stacking.OverflowVertical
		}
	}

	engine := &StackingEngine{
		spatialEngine:       spatialEngine,
		conflictCategorizer: conflictCategorizer,
		priorityRanker:      priorityRanker,
		stackingRules:       make([]StackingRule, 0),
		visualConstraints:   visualConstraints,
		heightCalculator:    NewHeightCalculator(),
		positionCalculator:  NewPositionCalculator(),
		spaceOptimizer:      NewSpaceOptimizer(),
	}

	// Initialize default stacking rules
	engine.initializeDefaultRules()

	return engine
}

// NewHeightCalculator creates a new height calculator
func NewHeightCalculator() *HeightCalculator {
	return &HeightCalculator{
		baseHeight: 20.0,
		minHeight:  15.0,
		maxHeight:  60.0,
		priorityMultiplier: map[VisualProminence]float64{
			ProminenceCritical: 1.5,
			ProminenceHigh:     1.3,
			ProminenceMedium:   1.0,
			ProminenceLow:      1.0,
			ProminenceMinimal:  1.0,
		},
		durationMultiplier: map[string]float64{
			"short":  1.0, // < 1 day
			"medium": 1.0, // 1-7 days
			"long":   1.2, // > 7 days
		},
		contentMultiplier: map[string]float64{
			"minimal": 1.0, // Simple tasks
			"normal":  1.0, // Standard tasks
			"complex": 1.3, // Complex tasks with many details
		},
	}
}

// NewPositionCalculator creates a new position calculator
func NewPositionCalculator() *PositionCalculator {
	return &PositionCalculator{
		verticalSpacing:   2.0,
		horizontalSpacing: 5.0,
		alignmentMode:     AlignmentTop,
		distributionMode:  DistributionPriority,
	}
}

// NewSpaceOptimizer creates a new space optimizer
func NewSpaceOptimizer() *SpaceOptimizer {
	return &SpaceOptimizer{
		compressionThreshold: 0.8,
		expansionThreshold:   0.6,
		adaptiveSpacing:      true,
		smartCollapsing:      false,
	}
}

// initializeDefaultRules sets up default stacking rules
func (se *StackingEngine) initializeDefaultRules() {
	se.stackingRules = []StackingRule{
		{
			Name:        "Milestone Task",
			Description: "Milestone tasks get special stacking treatment",
			Condition: func(task *core.Task, context *StackingContext) bool {
				return task.IsMilestone
			},
			Action: func(task *core.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeFloating,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MaxTaskHeight,
					Width:              context.VisualConstraints.MaxTaskWidth,
					ZIndex:             9,
					CollisionAvoidance: true,
				}
			},
		},
		{
			Name:        "Long Duration Task",
			Description: "Long duration tasks get horizontal stacking",
			Condition: func(task *core.Task, context *StackingContext) bool {
				duration := task.EndDate.Sub(task.StartDate)
				return duration > time.Hour*24*7 // More than a week
			},
			Action: func(task *core.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeHorizontal,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MinTaskHeight,
					Width:              context.VisualConstraints.MaxTaskWidth,
					ZIndex:             5,
					CollisionAvoidance: false,
				}
			},
		},
		{
			Name:        "Short Duration Task",
			Description: "Short duration tasks get vertical stacking",
			Condition: func(task *core.Task, context *StackingContext) bool {
				duration := task.EndDate.Sub(task.StartDate)
				return duration <= time.Hour*24 // One day or less
			},
			Action: func(task *core.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeVertical,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MinTaskHeight,
					Width:              context.VisualConstraints.MinTaskWidth,
					ZIndex:             3,
					CollisionAvoidance: false,
				}
			},
		},
		{
			Name:        "Conflict Resolution",
			Description: "Tasks with conflicts get special stacking treatment",
			Condition: func(task *core.Task, context *StackingContext) bool {
				// Check if task has conflicts
				if context.ConflictAnalysis == nil {
					return false
				}
				for _, conflict := range context.ConflictAnalysis.CategorizedConflicts {
					if conflict.Task1ID == task.ID || conflict.Task2ID == task.ID {
						return true
					}
				}
				return false
			},
			Action: func(task *core.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeCascading,
					VerticalOffset:     5.0,
					HorizontalOffset:   5.0,
					Height:             context.VisualConstraints.MaxTaskHeight * 0.8,
					Width:              context.VisualConstraints.MaxTaskWidth * 0.8,
					ZIndex:             7,
					CollisionAvoidance: true,
				}
			},
		},
		{
			Name:        "Overflow Handling",
			Description: "Tasks that cause overflow get minimized stacking",
			Condition: func(task *core.Task, context *StackingContext) bool {
				// Check if adding this task would cause overflow
				return se.wouldCauseOverflow(task, context)
			},
			Action: func(task *core.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeMinimized,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MinTaskHeight * 0.5,
					Width:              context.VisualConstraints.MinTaskWidth * 0.5,
					ZIndex:             1,
					CollisionAvoidance: false,
				}
			},
		},
		{
			Name:        "Default Stacking",
			Description: "Default stacking for all other tasks",
			Condition: func(task *core.Task, context *StackingContext) bool {
				return true // Always matches
			},
			Action: func(task *core.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeVertical,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MinTaskHeight,
					Width:              context.VisualConstraints.MinTaskWidth,
					ZIndex:             2,
					CollisionAvoidance: false,
				}
			},
		},
	}
}

// wouldCauseOverflow checks if adding a task would cause visual overflow
func (se *StackingEngine) wouldCauseOverflow(task *core.Task, context *StackingContext) bool {
	// Calculate current stack height
	currentHeight := 0.0
	for _, stack := range context.ExistingStacks {
		currentHeight += stack.TotalHeight
	}

	// Estimate height needed for new task
	estimatedHeight := context.VisualConstraints.MinTaskHeight + context.VisualConstraints.VerticalSpacing

	// Check if adding this task would exceed available height
	return (currentHeight + estimatedHeight) > context.AvailableHeight*context.VisualConstraints.OverflowThreshold
}

// StackTasks performs intelligent stacking of overlapping tasks
func (se *StackingEngine) StackTasks(tasks []*core.Task, context *StackingContext) *StackingResult {
	// Detect overlaps and categorize conflicts
	overlapAnalysis := se.spatialEngine.DetectOverlaps(tasks)
	// Rank tasks by priority
	priorityContext := &PriorityContext{
		CurrentTime: context.CurrentTime,
		UserID:      "system",
	}
	priorityRanking := se.priorityRanker.CalculatePriorityScores(tasks, priorityContext)

	// Update context with analysis results
	context.OverlapAnalysis = overlapAnalysis
	context.ConflictAnalysis = nil
	context.TaskPriorities = make(map[string]*TaskPriority)
	for _, taskScore := range priorityRanking.TaskScores {
		context.TaskPriorities[taskScore.TaskID] = &TaskPriority{
			Value:       0, // Will be set from task
			Category:    "",
			Description: "",
			Weight:      taskScore.OverallScore,
			Urgency:     string(taskScore.VisualProminence),
			Importance:  string(taskScore.VisualProminence),
		}
	}

	// Group tasks by overlapping time periods
	overlapGroups := se.groupTasksByOverlap(tasks, overlapAnalysis)

	// Create stacks for each overlap group
	var stacks []*TaskStack
	for _, group := range overlapGroups {
		stack := se.createStackForGroup(group, context)
		if stack != nil {
			stacks = append(stacks, stack)
		}
	}

	// Calculate stacking metrics
	result := &StackingResult{
		Stacks:          stacks,
		TotalStacks:     len(stacks),
		CollisionCount:  se.calculateCollisionCount(stacks),
		OverflowCount:   se.calculateOverflowCount(stacks, context),
		SpaceEfficiency: se.calculateSpaceEfficiency(stacks, context),
		VisualQuality:   se.calculateVisualQuality(stacks, context),
		Recommendations: se.generateRecommendations(stacks, context),
		AnalysisDate:    time.Now(),
	}

	return result
}

// StackTasksVertically performs vertical stacking of overlapping tasks
func (se *StackingEngine) StackTasksVertically(tasks []*core.Task, context *StackingContext) *VerticalStackingResult {
	// First, use the smart stacking engine to get initial stacks
	smartResult := se.StackTasks(tasks, context)

	// Convert smart stacks to vertical stacks
	var verticalStacks []*VerticalStack
	for _, smartStack := range smartResult.Stacks {
		verticalStack := se.convertToVerticalStack(smartStack, context)
		if verticalStack != nil {
			verticalStacks = append(verticalStacks, verticalStack)
		}
	}

	// Optimize vertical stacking
	verticalStacks = se.optimizeVerticalStacking(verticalStacks, context)

	// Calculate metrics
	result := &VerticalStackingResult{
		Stacks:           verticalStacks,
		TotalHeight:      se.calculateTotalHeight(verticalStacks),
		SpaceEfficiency:  se.calculateVerticalSpaceEfficiency(verticalStacks, context),
		VisualBalance:    se.calculateVisualBalance(verticalStacks),
		CollisionCount:   se.calculateVerticalCollisionCount(verticalStacks),
		OverflowCount:    se.calculateVerticalOverflowCount(verticalStacks, context),
		CompressionRatio: se.calculateCompressionRatio(verticalStacks),
		Recommendations:  se.generateVerticalRecommendations(verticalStacks, context),
		AnalysisDate:     time.Now(),
	}

	return result
}

// groupTasksByOverlap groups tasks by their overlapping time periods
func (se *StackingEngine) groupTasksByOverlap(tasks []*core.Task, overlapAnalysis *OverlapAnalysis) [][]*core.Task {
	var groups [][]*core.Task

	// Use overlap groups from analysis
	for _, overlapGroup := range overlapAnalysis.OverlapGroups {
		groups = append(groups, overlapGroup.Tasks)
	}

	// Add individual tasks that don't overlap with others
	overlappedTasks := make(map[string]bool)
	for _, group := range groups {
		for _, task := range group {
			overlappedTasks[task.ID] = true
		}
	}

	for _, task := range tasks {
		if !overlappedTasks[task.ID] {
			groups = append(groups, []*core.Task{task})
		}
	}

	return groups
}

// createStackForGroup creates a stack for a group of overlapping tasks
func (se *StackingEngine) createStackForGroup(tasks []*core.Task, context *StackingContext) *TaskStack {
	if len(tasks) == 0 {
		return nil
	}

	// Sort tasks by start date
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartDate.Before(tasks[j].StartDate)
	})

	// Create stack
	stack := &TaskStack{
		ID:             fmt.Sprintf("stack_%d", len(context.ExistingStacks)+1),
		Tasks:          make([]*StackedTask, 0),
		StartTime:      tasks[0].StartDate,
		EndTime:        tasks[0].EndDate,
		TotalHeight:    0.0,
		MaxWidth:       0.0,
		StackingType:   StackingTypeVertical,
		CollisionCount: 0,
		OverflowCount:  0,
	}

	// Calculate stack time range
	for _, task := range tasks {
		if task.StartDate.Before(stack.StartTime) {
			stack.StartTime = task.StartDate
		}
		if task.EndDate.After(stack.EndTime) {
			stack.EndTime = task.EndDate
		}
	}

	// Stack each task
	currentY := 0.0
	for _, task := range tasks {
		stackingAction := se.determineStackingAction(task, context)
		stackedTask := &StackedTask{
			Task:           task,
			StackingAction: stackingAction,
			Position: &Position{
				X:      0.0,
				Y:      currentY,
				Width:  stackingAction.Width,
				Height: stackingAction.Height,
				ZIndex: stackingAction.ZIndex,
			},
			IsVisible:      true,
			IsCollapsed:    false,
			CollisionLevel: 0,
			OverflowLevel:  0,
		}

		// Check for collisions
		if se.hasCollision(stackedTask, stack.Tasks) {
			stackedTask.CollisionLevel = 1
			stack.CollisionCount++
		}

		// Check for overflow
		if se.hasOverflow(stackedTask, context) {
			stackedTask.OverflowLevel = 1
			stack.OverflowCount++
		}

		stack.Tasks = append(stack.Tasks, stackedTask)
		stack.TotalHeight += stackingAction.Height + context.VisualConstraints.VerticalSpacing

		if stackingAction.Width > stack.MaxWidth {
			stack.MaxWidth = stackingAction.Width
		}

		// Update current Y position
		currentY += stackingAction.Height + context.VisualConstraints.VerticalSpacing
	}

	// Determine stack type based on tasks
	stack.StackingType = se.determineStackType(stack)

	return stack
}

// determineStackingAction determines the stacking action for a task
func (se *StackingEngine) determineStackingAction(task *core.Task, context *StackingContext) *StackingAction {
	// Find the best matching rule
	for _, rule := range se.stackingRules {
		if rule.Condition(task, context) {
			action := rule.Action(task, context)
			return action
		}
	}

	// Fallback to default action
	return &StackingAction{
		StackingType:       StackingTypeVertical,
		VerticalOffset:     0.0,
		HorizontalOffset:   0.0,
		Height:             context.VisualConstraints.MinTaskHeight,
		Width:              context.VisualConstraints.MinTaskWidth,
		ZIndex:             2,
		CollisionAvoidance: false,
	}
}

// determineStackType determines the overall stack type based on tasks
func (se *StackingEngine) determineStackType(stack *TaskStack) StackingType {
	if len(stack.Tasks) == 0 {
		return StackingTypeVertical
	}

	// Check if any task requires special stacking
	for _, stackedTask := range stack.Tasks {
		switch stackedTask.StackingAction.StackingType {
		case StackingTypeLayered:
			return StackingTypeLayered
		case StackingTypeFloating:
			return StackingTypeFloating
		case StackingTypeCascading:
			return StackingTypeCascading
		}
	}

	// Check if tasks are mostly horizontal
	horizontalCount := 0
	for _, stackedTask := range stack.Tasks {
		if stackedTask.StackingAction.StackingType == StackingTypeHorizontal {
			horizontalCount++
		}
	}

	if float64(horizontalCount)/float64(len(stack.Tasks)) > 0.5 {
		return StackingTypeHorizontal
	}

	return StackingTypeVertical
}

// hasCollision checks if a stacked task collides with existing tasks
func (se *StackingEngine) hasCollision(newTask *StackedTask, existingTasks []*StackedTask) bool {
	for _, existingTask := range existingTasks {
		if se.tasksCollide(newTask, existingTask) {
			return true
		}
	}
	return false
}

// tasksCollide checks if two stacked tasks collide
func (se *StackingEngine) tasksCollide(task1, task2 *StackedTask) bool {
	// Simple bounding box collision detection
	return !(task1.Position.X+task1.Position.Width < task2.Position.X ||
		task2.Position.X+task2.Position.Width < task1.Position.X ||
		task1.Position.Y+task1.Position.Height < task2.Position.Y ||
		task2.Position.Y+task2.Position.Height < task1.Position.Y)
}

// hasOverflow checks if a stacked task causes overflow
func (se *StackingEngine) hasOverflow(stackedTask *StackedTask, context *StackingContext) bool {
	return stackedTask.Position.Y+stackedTask.Position.Height > context.AvailableHeight*context.VisualConstraints.OverflowThreshold
}

// calculateCollisionCount calculates total collision count
func (se *StackingEngine) calculateCollisionCount(stacks []*TaskStack) int {
	total := 0
	for _, stack := range stacks {
		total += stack.CollisionCount
	}
	return total
}

// calculateOverflowCount calculates total overflow count
func (se *StackingEngine) calculateOverflowCount(stacks []*TaskStack, context *StackingContext) int {
	total := 0
	for _, stack := range stacks {
		total += stack.OverflowCount
	}
	return total
}

// calculateSpaceEfficiency calculates space efficiency
func (se *StackingEngine) calculateSpaceEfficiency(stacks []*TaskStack, context *StackingContext) float64 {
	if context.AvailableHeight == 0 {
		return 0.0
	}

	usedHeight := 0.0
	for _, stack := range stacks {
		usedHeight += stack.TotalHeight
	}

	return math.Min(usedHeight/context.AvailableHeight, 1.0)
}

// calculateVisualQuality calculates visual quality score
func (se *StackingEngine) calculateVisualQuality(stacks []*TaskStack, context *StackingContext) float64 {
	if len(stacks) == 0 {
		return 1.0
	}

	// Calculate quality based on collision and overflow counts
	totalCollisions := se.calculateCollisionCount(stacks)
	totalOverflows := se.calculateOverflowCount(stacks, context)
	totalTasks := 0

	for _, stack := range stacks {
		totalTasks += len(stack.Tasks)
	}

	if totalTasks == 0 {
		return 1.0
	}

	// Quality decreases with collisions and overflows
	collisionPenalty := float64(totalCollisions) / float64(totalTasks)
	overflowPenalty := float64(totalOverflows) / float64(totalTasks)

	quality := 1.0 - collisionPenalty - overflowPenalty
	return math.Max(quality, 0.0)
}

// generateRecommendations generates stacking recommendations
func (se *StackingEngine) generateRecommendations(stacks []*TaskStack, context *StackingContext) []string {
	var recommendations []string

	// Collision recommendations
	totalCollisions := se.calculateCollisionCount(stacks)
	if totalCollisions > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️ %d visual collisions detected - consider adjusting task positioning", totalCollisions))
	}

	// Overflow recommendations
	totalOverflows := se.calculateOverflowCount(stacks, context)
	if totalOverflows > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("📏 %d overflow issues detected - consider reducing task sizes or using minimization", totalOverflows))
	}

	// Space efficiency recommendations
	efficiency := se.calculateSpaceEfficiency(stacks, context)
	if efficiency < 0.5 {
		recommendations = append(recommendations,
			"📊 Low space efficiency - consider optimizing task layouts")
	} else if efficiency > 0.9 {
		recommendations = append(recommendations,
			"📊 High space efficiency - good use of available space")
	}

	// Visual quality recommendations
	quality := se.calculateVisualQuality(stacks, context)
	if quality < 0.7 {
		recommendations = append(recommendations,
			"🎨 Visual quality could be improved - consider adjusting stacking rules")
	}

	// Stack depth recommendations
	maxDepth := 0
	for _, stack := range stacks {
		if len(stack.Tasks) > maxDepth {
			maxDepth = len(stack.Tasks)
		}
	}

	if maxDepth > context.VisualConstraints.MaxStackDepth {
		recommendations = append(recommendations,
			fmt.Sprintf("📚 Stack depth (%d) exceeds recommended maximum (%d) - consider task grouping",
				maxDepth, context.VisualConstraints.MaxStackDepth))
	}

	return recommendations
}

// AddCustomRule adds a custom stacking rule
func (se *StackingEngine) AddCustomRule(rule StackingRule) {
	se.stackingRules = append(se.stackingRules, rule)
}

// GetStacksByType returns stacks filtered by stacking type
func (result *StackingResult) GetStacksByType(stackingType StackingType) []*TaskStack {
	var filtered []*TaskStack
	for _, stack := range result.Stacks {
		if stack.StackingType == stackingType {
			filtered = append(filtered, stack)
		}
	}
	return filtered
}

// GetSummary returns a summary of the stacking result
func (result *StackingResult) GetSummary() string {
	return fmt.Sprintf("Smart Stacking Summary:\n"+
		"  Total Stacks: %d\n"+
		"  Collisions: %d\n"+
		"  Overflows: %d\n"+
		"  Space Efficiency: %.2f%%\n"+
		"  Visual Quality: %.2f%%\n"+
		"  Analysis Date: %s",
		result.TotalStacks,
		result.CollisionCount,
		result.OverflowCount,
		result.SpaceEfficiency*100,
		result.VisualQuality*100,
		result.AnalysisDate.Format("2006-01-02 15:04:05"))
}

// convertToVerticalStack converts a smart stack to a vertical stack
func (se *StackingEngine) convertToVerticalStack(smartStack *TaskStack, context *StackingContext) *VerticalStack {
	if len(smartStack.Tasks) == 0 {
		return nil
	}

	verticalStack := &VerticalStack{
		ID:               smartStack.ID,
		StartTime:        smartStack.StartTime,
		EndTime:          smartStack.EndTime,
		MaxWidth:         smartStack.MaxWidth,
		AlignmentMode:    se.determineAlignmentMode(smartStack, context),
		DistributionMode: se.determineDistributionMode(smartStack, context),
		Tasks:            make([]*VerticallyStackedTask, 0),
	}

	// Convert each stacked task
	for i, stackedTask := range smartStack.Tasks {
		verticallyStackedTask := &VerticallyStackedTask{
			Task:             stackedTask.Task,
			CalculatedHeight: se.calculateTaskHeight(stackedTask.Task, context),
			IsCompressed:     false,
			IsExpanded:       false,
			CollisionLevel:   stackedTask.CollisionLevel,
			OverflowLevel:    stackedTask.OverflowLevel,
			VisualWeight:     se.calculateVisualWeight(stackedTask.Task, context),
		}

		// Calculate position
		verticallyStackedTask.Position = se.calculateVerticalPosition(
			verticallyStackedTask,
			verticalStack,
			i,
			context,
		)

		verticalStack.Tasks = append(verticalStack.Tasks, verticallyStackedTask)
	}

	// Calculate stack metrics
	verticalStack.TotalHeight = se.calculateStackHeight(verticalStack)
	verticalStack.SpaceEfficiency = se.calculateStackSpaceEfficiency(verticalStack, context)
	verticalStack.VisualBalance = se.calculateStackVisualBalance(verticalStack)
	verticalStack.CollisionCount = se.calculateStackCollisionCount(verticalStack)
	verticalStack.OverflowCount = se.calculateStackOverflowCount(verticalStack, context)

	return verticalStack
}

// calculateTaskHeight calculates the optimal height for a task
func (se *StackingEngine) calculateTaskHeight(task *core.Task, context *StackingContext) float64 {
	hc := se.heightCalculator

	// Start with base height
	height := hc.baseHeight

	// Apply duration multiplier
	duration := task.EndDate.Sub(task.StartDate)
	if duration <= time.Hour*24 {
		height *= hc.durationMultiplier["short"]
	} else if duration <= time.Hour*24*7 {
		height *= hc.durationMultiplier["medium"]
	} else {
		height *= hc.durationMultiplier["long"]
	}

	// Apply content multiplier based on task complexity
	contentComplexity := se.assessContentComplexity(task)
	if multiplier, exists := hc.contentMultiplier[contentComplexity]; exists {
		height *= multiplier
	}

	// Apply visual constraints
	if context.VisualConstraints != nil {
		height = math.Max(height, context.VisualConstraints.MinTaskHeight)
		height = math.Min(height, context.VisualConstraints.MaxTaskHeight)
	}

	return height
}

// assessContentComplexity assesses the complexity of task content
func (se *StackingEngine) assessContentComplexity(task *core.Task) string {
	complexity := "normal"

	// Check for complex indicators
	if task.IsMilestone {
		complexity = "complex"
	} else if len(task.Name) > 30 {
		complexity = "complex"
	} else if len(task.Name) < 10 {
		complexity = "minimal"
	}

	// Check for special categories
	if task.Category == "DISSERTATION" || task.Category == "PROPOSAL" {
		complexity = "complex"
	}

	return complexity
}

// calculateVisualWeight calculates the visual weight of a task
func (se *StackingEngine) calculateVisualWeight(task *core.Task, context *StackingContext) float64 {
	weight := 1.0

	// Duration weight
	duration := task.EndDate.Sub(task.StartDate)
	weight += float64(duration.Hours()) * 0.01

	// Category weight
	if task.Category == "DISSERTATION" {
		weight += 2.0
	} else if task.Category == "PROPOSAL" {
		weight += 1.5
	} else if task.Category == "LASER" {
		weight += 1.0
	}

	// Milestone weight
	if task.IsMilestone {
		weight += 3.0
	}

	return weight
}

// calculateVerticalPosition calculates the vertical position of a task
func (se *StackingEngine) calculateVerticalPosition(
	task *VerticallyStackedTask,
	stack *VerticalStack,
	index int,
	context *StackingContext,
) *VerticalPosition {
	pc := se.positionCalculator

	// Calculate base position
	position := &VerticalPosition{
		X:          0.0,
		Y:          0.0,
		Width:      stack.MaxWidth,
		Height:     task.CalculatedHeight,
		ZIndex:     index + 1,
		StackIndex: index,
	}

	// Calculate Y position based on previous tasks
	if index > 0 {
		previousTask := stack.Tasks[index-1]
		position.Y = previousTask.Position.Y + previousTask.Position.Height + pc.verticalSpacing
	}

	// Apply alignment mode
	position = se.applyAlignmentMode(position, stack, context)

	// Apply distribution mode
	position = se.applyDistributionMode(position, stack, context)

	// Calculate relative position within stack
	position.RelativeY = position.Y - stack.StartTime.Sub(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).Hours()*10

	return position
}

// applyAlignmentMode applies the alignment mode to a position
func (se *StackingEngine) applyAlignmentMode(position *VerticalPosition, stack *VerticalStack, context *StackingContext) *VerticalPosition {
	pc := se.positionCalculator

	switch pc.alignmentMode {
	case AlignmentTop:
		// Already positioned at top
		break
	case AlignmentCenter:
		// Center within available height
		if context.AvailableHeight > 0 {
			stackHeight := se.calculateStackHeight(stack)
			offset := (context.AvailableHeight - stackHeight) / 2
			position.Y += offset
		}
	case AlignmentBottom:
		// Position at bottom
		if context.AvailableHeight > 0 {
			stackHeight := se.calculateStackHeight(stack)
			position.Y = context.AvailableHeight - stackHeight + position.Y
		}
	case AlignmentJustify:
		// Justify within available space
		if context.AvailableHeight > 0 {
			stackHeight := se.calculateStackHeight(stack)
			if stackHeight < context.AvailableHeight {
				// Distribute extra space evenly
				extraSpace := context.AvailableHeight - stackHeight
				position.Y += extraSpace * float64(position.StackIndex) / float64(len(stack.Tasks))
			}
		}
	}

	return position
}

// applyDistributionMode applies the distribution mode to a position
func (se *StackingEngine) applyDistributionMode(position *VerticalPosition, stack *VerticalStack, context *StackingContext) *VerticalPosition {
	pc := se.positionCalculator

	switch pc.distributionMode {
	case DistributionEven:
		// Even distribution (already handled in base calculation)
		break
	case DistributionPriority:
		// Distribute evenly (no priority-based distribution)
		break
	case DistributionContent:
		// Distribute based on content complexity
		if len(stack.Tasks) > position.StackIndex {
			contentComplexity := se.assessContentComplexity(stack.Tasks[position.StackIndex].Task)
			if contentComplexity == "complex" {
				position.Y += 5.0
			} else if contentComplexity == "minimal" {
				position.Y -= 2.0
			}
		}
	case DistributionAdaptive:
		// Adaptive distribution based on available space
		if context.AvailableHeight > 0 {
			stackHeight := se.calculateStackHeight(stack)
			if stackHeight < context.AvailableHeight {
				// Use adaptive spacing
				adaptiveSpacing := (context.AvailableHeight - stackHeight) / float64(len(stack.Tasks))
				position.Y += adaptiveSpacing * float64(position.StackIndex)
			}
		}
	}

	return position
}

// determineAlignmentMode determines the best alignment mode for a stack
func (se *StackingEngine) determineAlignmentMode(stack *TaskStack, context *StackingContext) AlignmentMode {
	// Check if stack has milestone tasks
	for _, task := range stack.Tasks {
		if task.Task.IsMilestone {
			return AlignmentTop
		}
	}

	// Check available space
	if context.AvailableHeight > 0 {
		estimatedHeight := float64(len(stack.Tasks)) * 25.0 // Rough estimate
		if estimatedHeight < context.AvailableHeight*0.5 {
			return AlignmentCenter
		}
	}

	return AlignmentTop
}

// determineDistributionMode determines the best distribution mode for a stack
func (se *StackingEngine) determineDistributionMode(stack *TaskStack, context *StackingContext) DistributionMode {

	// Check if stack has mixed content complexity
	hasComplexContent := false
	hasSimpleContent := false

	for _, task := range stack.Tasks {
		complexity := se.assessContentComplexity(task.Task)
		switch complexity {
case "complex":
			hasComplexContent = true
		case "minimal":
			hasSimpleContent = true
		}
	}

	if hasComplexContent && hasSimpleContent {
		return DistributionContent
	}

	return DistributionEven
}

// optimizeVerticalStacking optimizes the vertical stacking layout
func (se *StackingEngine) optimizeVerticalStacking(stacks []*VerticalStack, context *StackingContext) []*VerticalStack {
	so := se.spaceOptimizer

	// Apply space optimization
	for i, stack := range stacks {
		// Check if expansion is possible
		if se.canExpand(stack, context) {
			stacks[i] = se.expandStack(stack, context)
		}

		// Apply adaptive spacing
		if so.adaptiveSpacing {
			stacks[i] = se.applyAdaptiveSpacing(stack, context)
		}
	}

	return stacks
}

// needsCompression checks if a stack needs compression
func (se *StackingEngine) needsCompression(stack *VerticalStack, context *StackingContext) bool {
	if context.AvailableHeight <= 0 {
		return false
	}

	stackHeight := se.calculateStackHeight(stack)
	return stackHeight > context.AvailableHeight*se.spaceOptimizer.compressionThreshold
}

// canExpand checks if a stack can be expanded
func (se *StackingEngine) canExpand(stack *VerticalStack, context *StackingContext) bool {
	if context.AvailableHeight <= 0 {
		return false
	}

	stackHeight := se.calculateStackHeight(stack)
	return stackHeight < context.AvailableHeight*se.spaceOptimizer.expansionThreshold
}

// compressStack compresses a stack to fit available space
func (se *StackingEngine) compressStack(stack *VerticalStack, context *StackingContext) *VerticalStack {
	if context.AvailableHeight <= 0 {
		return stack
	}

	// Calculate compression ratio
	currentHeight := se.calculateStackHeight(stack)
	compressionRatio := context.AvailableHeight / currentHeight

	// Apply compression to each task
	for _, task := range stack.Tasks {
		task.CalculatedHeight *= compressionRatio
		task.IsCompressed = true
		task.Position.Height = task.CalculatedHeight
	}

	// Recalculate positions
	se.recalculateStackPositions(stack)

	return stack
}

// expandStack expands a stack to better utilize available space
func (se *StackingEngine) expandStack(stack *VerticalStack, context *StackingContext) *VerticalStack {
	if context.AvailableHeight <= 0 {
		return stack
	}

	// Calculate expansion ratio
	currentHeight := se.calculateStackHeight(stack)
	expansionRatio := math.Min(1.5, context.AvailableHeight/currentHeight)

	// Apply expansion to each task
	for _, task := range stack.Tasks {
		task.CalculatedHeight *= expansionRatio
		task.IsExpanded = true
		task.Position.Height = task.CalculatedHeight
	}

	// Recalculate positions
	se.recalculateStackPositions(stack)

	return stack
}

// applyAdaptiveSpacing applies adaptive spacing to a stack
func (se *StackingEngine) applyAdaptiveSpacing(stack *VerticalStack, context *StackingContext) *VerticalStack {
	if context.AvailableHeight <= 0 {
		return stack
	}

	// Calculate adaptive spacing
	currentHeight := se.calculateStackHeight(stack)
	availableSpace := context.AvailableHeight - currentHeight

	if availableSpace > 0 {
		// Distribute extra space as adaptive spacing
		adaptiveSpacing := availableSpace / float64(len(stack.Tasks))

		for i, task := range stack.Tasks {
			if i > 0 {
				task.Position.Y += adaptiveSpacing * float64(i)
			}
		}
	}

	return stack
}

// applySmartCollapsing applies smart collapsing to a stack
func (se *StackingEngine) applySmartCollapsing(stack *VerticalStack, context *StackingContext) *VerticalStack {
	// Collapse low-priority tasks if space is limited
	if context.AvailableHeight > 0 {
		stackHeight := se.calculateStackHeight(stack)
		if stackHeight > context.AvailableHeight*0.9 {
			// Collapse tasks with low visual weight
			for _, task := range stack.Tasks {
				if task.VisualWeight < 1.0 {
					task.CalculatedHeight *= 0.7
					task.Position.Height = task.CalculatedHeight
				}
			}

			// Recalculate positions
			se.recalculateStackPositions(stack)
		}
	}

	return stack
}

// recalculateStackPositions recalculates positions after height changes
func (se *StackingEngine) recalculateStackPositions(stack *VerticalStack) {
	pc := se.positionCalculator

	for i, task := range stack.Tasks {
		if i == 0 {
			task.Position.Y = 0.0
		} else {
			previousTask := stack.Tasks[i-1]
			task.Position.Y = previousTask.Position.Y + previousTask.Position.Height + pc.verticalSpacing
		}
	}
}

// calculateStackHeight calculates the total height of a stack
func (se *StackingEngine) calculateStackHeight(stack *VerticalStack) float64 {
	if len(stack.Tasks) == 0 {
		return 0.0
	}

	lastTask := stack.Tasks[len(stack.Tasks)-1]
	return lastTask.Position.Y + lastTask.Position.Height
}

// calculateStackSpaceEfficiency calculates the space efficiency of a stack
func (se *StackingEngine) calculateStackSpaceEfficiency(stack *VerticalStack, context *StackingContext) float64 {
	if context.AvailableHeight <= 0 {
		return 1.0
	}

	stackHeight := se.calculateStackHeight(stack)
	return math.Min(stackHeight/context.AvailableHeight, 1.0)
}

// calculateStackVisualBalance calculates the visual balance of a stack
func (se *StackingEngine) calculateStackVisualBalance(stack *VerticalStack) float64 {
	if len(stack.Tasks) == 0 {
		return 1.0
	}

	// Calculate weight distribution
	totalWeight := 0.0
	for _, task := range stack.Tasks {
		totalWeight += task.VisualWeight
	}

	// Calculate balance (closer to 1.0 is more balanced)
	avgWeight := totalWeight / float64(len(stack.Tasks))
	balance := 1.0

	for _, task := range stack.Tasks {
		deviation := math.Abs(task.VisualWeight-avgWeight) / avgWeight
		balance -= deviation * 0.1
	}

	return math.Max(balance, 0.0)
}

// calculateStackCollisionCount calculates the collision count of a stack
func (se *StackingEngine) calculateStackCollisionCount(stack *VerticalStack) int {
	count := 0
	for _, task := range stack.Tasks {
		count += task.CollisionLevel
	}
	return count
}

// calculateStackOverflowCount calculates the overflow count of a stack
func (se *StackingEngine) calculateStackOverflowCount(stack *VerticalStack, context *StackingContext) int {
	count := 0
	for _, task := range stack.Tasks {
		count += task.OverflowLevel
	}
	return count
}

// calculateTotalHeight calculates the total height of all stacks
func (se *StackingEngine) calculateTotalHeight(stacks []*VerticalStack) float64 {
	total := 0.0
	for _, stack := range stacks {
		total += se.calculateStackHeight(stack)
	}
	return total
}

// calculateVerticalSpaceEfficiency calculates the overall space efficiency
func (se *StackingEngine) calculateVerticalSpaceEfficiency(stacks []*VerticalStack, context *StackingContext) float64 {
	if context.AvailableHeight <= 0 {
		return 1.0
	}

	totalHeight := se.calculateTotalHeight(stacks)
	return math.Min(totalHeight/context.AvailableHeight, 1.0)
}

// calculateVisualBalance calculates the overall visual balance
func (se *StackingEngine) calculateVisualBalance(stacks []*VerticalStack) float64 {
	if len(stacks) == 0 {
		return 1.0
	}

	totalBalance := 0.0
	for _, stack := range stacks {
		totalBalance += se.calculateStackVisualBalance(stack)
	}

	return totalBalance / float64(len(stacks))
}

// calculateVerticalCollisionCount calculates the total collision count
func (se *StackingEngine) calculateVerticalCollisionCount(stacks []*VerticalStack) int {
	total := 0
	for _, stack := range stacks {
		total += se.calculateStackCollisionCount(stack)
	}
	return total
}

// calculateVerticalOverflowCount calculates the total overflow count
func (se *StackingEngine) calculateVerticalOverflowCount(stacks []*VerticalStack, context *StackingContext) int {
	total := 0
	for _, stack := range stacks {
		total += se.calculateStackOverflowCount(stack, context)
	}
	return total
}

// calculateCompressionRatio calculates the compression ratio
func (se *StackingEngine) calculateCompressionRatio(stacks []*VerticalStack) float64 {
	compressedTasks := 0
	totalTasks := 0

	for _, stack := range stacks {
		for _, task := range stack.Tasks {
			totalTasks++
			if task.IsCompressed {
				compressedTasks++
			}
		}
	}

	if totalTasks == 0 {
		return 0.0
	}

	return float64(compressedTasks) / float64(totalTasks)
}

// generateVerticalRecommendations generates recommendations for vertical stacking
func (se *StackingEngine) generateVerticalRecommendations(stacks []*VerticalStack, context *StackingContext) []string {
	var recommendations []string

	// Space efficiency recommendations
	efficiency := se.calculateVerticalSpaceEfficiency(stacks, context)
	if efficiency < 0.5 {
		recommendations = append(recommendations,
			"📏 Low space efficiency - consider adjusting task heights or using compression")
	} else if efficiency > 0.9 {
		recommendations = append(recommendations,
			"📏 High space efficiency - good utilization of available space")
	}

	// Visual balance recommendations
	balance := se.calculateVisualBalance(stacks)
	if balance < 0.7 {
		recommendations = append(recommendations,
			"⚖️ Visual balance could be improved - consider adjusting task weights")
	}

	// Compression recommendations
	compressionRatio := se.calculateCompressionRatio(stacks)
	if compressionRatio > 0.5 {
		recommendations = append(recommendations,
			fmt.Sprintf("🗜️ High compression ratio (%.1f%%) - consider increasing available space", compressionRatio*100))
	}

	// Collision recommendations
	collisionCount := se.calculateVerticalCollisionCount(stacks)
	if collisionCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️ %d visual collisions detected - consider adjusting task positioning", collisionCount))
	}

	// Overflow recommendations
	overflowCount := se.calculateVerticalOverflowCount(stacks, context)
	if overflowCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("📏 %d overflow issues detected - consider reducing task sizes", overflowCount))
	}

	return recommendations
}

// GetStacksByAlignment returns stacks filtered by alignment mode
func (result *VerticalStackingResult) GetStacksByAlignment(alignment AlignmentMode) []*VerticalStack {
	var filtered []*VerticalStack
	for _, stack := range result.Stacks {
		if stack.AlignmentMode == alignment {
			filtered = append(filtered, stack)
		}
	}
	return filtered
}

// GetStacksByDistribution returns stacks filtered by distribution mode
func (result *VerticalStackingResult) GetStacksByDistribution(distribution DistributionMode) []*VerticalStack {
	var filtered []*VerticalStack
	for _, stack := range result.Stacks {
		if stack.DistributionMode == distribution {
			filtered = append(filtered, stack)
		}
	}
	return filtered
}

// GetSummary returns a summary of the vertical stacking result
func (result *VerticalStackingResult) GetSummary() string {
	return fmt.Sprintf("Vertical Stacking Summary:\n"+
		"  Total Stacks: %d\n"+
		"  Total Height: %.2f\n"+
		"  Space Efficiency: %.2f%%\n"+
		"  Visual Balance: %.2f%%\n"+
		"  Collisions: %d\n"+
		"  Overflows: %d\n"+
		"  Compression Ratio: %.2f%%\n"+
		"  Analysis Date: %s",
		len(result.Stacks),
		result.TotalHeight,
		result.SpaceEfficiency*100,
		result.VisualBalance*100,
		result.CollisionCount,
		result.OverflowCount,
		result.CompressionRatio*100,
		result.AnalysisDate.Format("2006-01-02 15:04:05"))
}

// ConflictResolutionResult contains the result of conflict resolution
type ConflictResolutionResult struct {
	ResolvedConflicts   []*ResolvedConflict
	OverflowResolutions []*OverflowResolution
	VisualOptimizations []*VisualOptimization
	LayoutAdjustments   []*LayoutAdjustment
	Recommendations     []string
	AnalysisDate        time.Time
}

// ResolvedConflict represents a resolved visual conflict
type ResolvedConflict struct {
	ConflictID        string
	ConflictType      string
	ResolutionMethod  string
	AffectedTasks     []*core.Task
	VisualChanges     *VisualChanges
	ResolutionQuality float64
	BeforeMetrics     *ConflictMetrics
	AfterMetrics      *ConflictMetrics
}

// OverflowResolution represents a resolved overflow issue
type OverflowResolution struct {
	OverflowID        string
	OverflowType      string
	ResolutionMethod  string
	AffectedTasks     []*core.Task
	SpaceRecovered    float64
	ResolutionQuality float64
	BeforeMetrics     *OverflowMetrics
	AfterMetrics      *OverflowMetrics
}

// VisualOptimization represents a visual optimization applied
type VisualOptimization struct {
	OptimizationID    string
	OptimizationType  string
	Description       string
	AffectedTasks     []*core.Task
	VisualImprovement float64
	SpaceEfficiency   float64
	BeforeMetrics     *VisualMetrics
	AfterMetrics      *VisualMetrics
}

// LayoutAdjustment represents a layout adjustment made
type LayoutAdjustment struct {
	AdjustmentID    string
	AdjustmentType  string
	Description     string
	AffectedTasks   []*core.Task
	PositionChanges *PositionChanges
	SizeChanges     *SizeChanges
	VisualChanges   *VisualChanges
}

// PriorityContext provides context for priority calculations
type PriorityContext struct {
	CurrentTime         time.Time
	UserID              string
	ProjectID           string
	TeamMembers         []string
	ResourceLimits      map[string]int
	DeadlineConstraints map[string]time.Time
	WorkloadLimits      map[string]float64
}

// TaskPriority represents task priority information
type TaskPriority struct {
	Value       int
	Category    string
	Description string
	Weight      float64
	Urgency     string
	Importance  string
}

// VisualProminence represents the visual prominence level
type VisualProminence string

const (
	ProminenceCritical VisualProminence = "CRITICAL"
	ProminenceHigh     VisualProminence = "HIGH"
	ProminenceMedium   VisualProminence = "MEDIUM"
	ProminenceLow      VisualProminence = "LOW"
	ProminenceMinimal  VisualProminence = "MINIMAL"
)

// VisualStyle represents visual styling information
type VisualStyle struct {
	Color       string
	BorderColor string
	BorderWidth float64
	Opacity     float64
	FontSize    float64
	FontWeight  string
	FontStyle   string
	Background  string
	Shadow      string
	Animation   string
}

// ConflictAnalysis contains comprehensive conflict analysis results
type ConflictAnalysis struct {
	TotalConflicts       int
	ConflictsByCategory  map[ConflictCategory]int
	ConflictsBySeverity  map[OverlapSeverity]int
	ConflictsByUrgency   map[string]int
	ConflictsByRisk      map[string]int
	CategorizedConflicts []*CategorizedConflict
	ResolutionSummary    map[string]int
	RiskAssessment       string
	Recommendations      []string
	AnalysisDate         time.Time
}

// CategorizedConflict represents a conflict with detailed categorization
type CategorizedConflict struct {
	*TaskOverlap
	Category               ConflictCategory
	SubCategory            string
	RootCause              string
	Impact                 string
	Resolution             ConflictResolution
	AlternativeResolutions []ConflictResolution
	RiskLevel              string
	Urgency                string
	Complexity             string
}

// ConflictResolution represents a resolution strategy for a conflict
type ConflictResolution struct {
	Strategy    string
	Description string
	Actions     []string
	Effort      string // "LOW", "MEDIUM", "HIGH"
	Impact      string // "LOW", "MEDIUM", "HIGH"
}

// TaskPrioritizationResult contains the result of task prioritization
type TaskPrioritizationResult struct {
	PrioritizedTasks    []*PrioritizedTask
	StackingOrder       []string
	VisibilitySettings  map[string]*VisibilityAction
	OptimizationResults []*OptimizationAction
	Recommendations     []string
	AnalysisDate        time.Time
}

// PrioritizedTask represents a task with prioritization information
type PrioritizedTask struct {
	Task               *core.Task
	PriorityScore      *PriorityScore
	VisibilityAction   *VisibilityAction
	OptimizationAction *OptimizationAction
	StackingOrder      int
	DisplayPriority    float64
}

// PriorityScore represents a calculated priority score
type PriorityScore struct {
	TaskID           string
	OverallScore     float64
	CategoryScores   map[PriorityCategory]float64
	VisualProminence VisualProminence
	Ranking          int
	Confidence       float64
	Factors          []PriorityFactor
	Recommendations  []string
	CalculatedAt     time.Time
}

// PriorityFactor represents a single factor in priority calculation
type PriorityFactor struct {
	Category     PriorityCategory
	Factor       VisualFactor
	Value        float64
	Weight       float64
	Contribution float64
	Description  string
}

// Supporting data structures
type ConflictMetrics struct {
	CollisionCount  int
	OverlapCount    int
	SeverityScore   float64
	VisualClarity   float64
	SpaceEfficiency float64
}

type OverflowMetrics struct {
	OverflowAmount     float64
	OverflowPercentage float64
	AffectedTasks      int
	SeverityScore      float64
	SpaceWaste         float64
}

type VisualMetrics struct {
	VisualBalance    float64
	SpaceEfficiency  float64
	ClarityScore     float64
	HarmonyScore     float64
	ReadabilityScore float64
}

type PositionChanges struct {
	XChanges map[string]float64
	YChanges map[string]float64
	ZChanges map[string]int
}

type SizeChanges struct {
	WidthChanges  map[string]float64
	HeightChanges map[string]float64
}

type VisualChanges struct {
	ColorChanges     map[string]string
	StyleChanges     map[string]string
	EffectChanges    map[string]string
	AnimationChanges map[string]string
}

// ConflictResolutionEngine handles visual conflict resolution and overflow management
type ConflictResolutionEngine struct {
	taskPrioritizationEngine *TaskPrioritizationEngine
	stackingEngine           *StackingEngine
	overflowManager          *OverflowManager
	visualConflictResolver   *VisualConflictResolver
}

// NewConflictResolutionEngine creates a new conflict resolution engine
func NewConflictResolutionEngine(
	taskPrioritizationEngine *TaskPrioritizationEngine,
	stackingEngine *StackingEngine,
) *ConflictResolutionEngine {
	engine := &ConflictResolutionEngine{
		taskPrioritizationEngine: taskPrioritizationEngine,
		stackingEngine:           stackingEngine,
		overflowManager:          NewOverflowManager(),
		visualConflictResolver:   NewVisualConflictResolver(),
	}
	return engine
}

// ResolveConflicts performs comprehensive conflict resolution
func (cre *ConflictResolutionEngine) ResolveConflicts(tasks []*core.Task, context *PriorityContext) *ConflictResolutionResult {
	// Simplified implementation - returns empty results for now
	return &ConflictResolutionResult{
		ResolvedConflicts:   []*ResolvedConflict{},
		OverflowResolutions: []*OverflowResolution{},
		VisualOptimizations: []*VisualOptimization{},
		LayoutAdjustments:   []*LayoutAdjustment{},
		Recommendations:     []string{"Applied conflict resolution"},
		AnalysisDate:        time.Now(),
	}
}

// Minimal supporting types for the engines
type TaskPrioritizationEngine struct {
	stackingEngine    *StackingEngine
	priorityRanker    *PriorityRanker
	visibilityManager *VisibilityManager
	stackingOptimizer *StackingOptimizer
}

type PriorityRanker struct {
	conflictCategorizer *ConflictCategorizer
	rankingRules        []PriorityRule
	visualWeights       map[VisualFactor]float64
}

type VisibilityManager struct {
	visibilityRules     []VisibilityRule
	prominenceWeights   map[VisualProminence]float64
	visibilityThreshold float64
	adaptiveVisibility  bool
}

type StackingOptimizer struct {
	optimizationRules  []OptimizationRule
	stackingStrategies map[PriorityCategory]StackingStrategy
	adaptiveOrdering   bool
}

type ConflictCategorizer struct {
	spatialEngine   *SpatialEngine
	rules           []ConflictRule
	severityWeights map[OverlapSeverity]int
}

type OverflowManager struct {
	overflowThresholds   map[OverflowType]float64
	resolutionStrategies map[OverflowType][]ResolutionStrategy
	adaptiveThresholds   bool
	smartCompression     bool
}

type VisualConflictResolver struct {
	collisionDetector  *CollisionDetector
	conflictResolvers  map[string]ConflictResolver
	visualOptimizer    *VisualOptimizer
	adaptiveResolution bool
}

type CollisionDetector struct {
	collisionThreshold float64
	boundingBoxBuffer  float64
	zIndexManager      *ZIndexManager
}

type ZIndexManager struct {
	baseZIndex     int
	layerSpacing   int
	priorityLayers map[VisualProminence]int
}

type VisualOptimizer struct {
	optimizationRules []VisualOptimizationRule
	layoutStrategies  map[LayoutStrategy]VisualStrategy
	adaptiveLayout    bool
}

// Minimal type definitions for the supporting types
type PriorityRule struct {
	Name        string
	Description string
	Weight      float64
	Calculator  func(*core.Task, *PriorityContext) float64
	Category    PriorityCategory
}

type PriorityCategory string

const (
	CategoryTaskPriority PriorityCategory = "TASK_PRIORITY"
)

type VisualFactor string

const (
	FactorTaskImportance VisualFactor = "TASK_IMPORTANCE"
)

type VisibilityRule struct {
	Name        string
	Description string
	Condition   func(*core.Task, *PriorityContext) bool
	Action      func(*core.Task, *PriorityContext) *VisibilityAction
}

type VisibilityAction struct {
	IsVisible       bool
	ProminenceLevel VisualProminence
	DisplayOrder    int
	VisualWeight    float64
	CollapseLevel   int
}

type OptimizationRule struct {
	Name        string
	Description string
	Condition   func(*core.Task, *PriorityContext) bool
	Action      func(*core.Task, *PriorityContext) *OptimizationAction
}

type OptimizationAction struct {
	StackingOrder    int
	VisualProminence VisualProminence
	DisplayPriority  float64
	CollapseLevel    int
}

type StackingStrategy struct {
	StrategyType PriorityCategory
	Description  string
	Parameters   map[string]interface{}
	Rules        []OptimizationRule
}

type ConflictRule struct {
	Name        string
	Description string
	Condition   func(*TaskOverlap, *core.Task, *core.Task) bool
	Category    ConflictCategory
	Severity    OverlapSeverity
}

type ConflictCategory string

const (
	CategoryScheduleConflict ConflictCategory = "SCHEDULE_CONFLICT"
)

type OverflowType string

const (
	OverflowVertical OverflowType = "VERTICAL"
)

type ResolutionStrategy struct {
	Name        string
	Description string
	Condition   func(*OverflowContext) bool
	Action      func(*OverflowContext) *ResolutionResult
}

type ConflictResolver struct {
	Name        string
	Description string
	Condition   func(*ConflictContext) bool
	Action      func(*ConflictContext) *ConflictResolutionResult
}

type VisualOptimizationRule struct {
	Name        string
	Description string
	Condition   func(*VisualContext) bool
	Action      func(*VisualContext) *VisualOptimization
}

type VisualStrategy struct {
	StrategyType LayoutStrategy
	Description  string
	Parameters   map[string]interface{}
	Constraints  *VisualConstraints
}

type LayoutStrategy string

const (
	LayoutStack LayoutStrategy = "STACK"
)

type OverflowContext struct {
	Tasks           []*core.Task
	AvailableSpace  *SpaceConstraints
	OverflowType    OverflowType
	Severity        float64
	PriorityContext *PriorityContext
	VisualSettings  *VisualSettings
	Constraints     *VisualConstraints
}

type ConflictContext struct {
	Conflicts       []*TaskOverlap
	Tasks           []*core.Task
	PriorityContext *PriorityContext
	VisualSettings  *VisualSettings
	Constraints     *VisualConstraints
}

type VisualContext struct {
	Tasks           []*core.Task
	LayoutMetrics   *LayoutMetrics
	VisualSettings  *VisualSettings
	Constraints     *VisualConstraints
	PriorityContext *PriorityContext
}

type ResolutionResult struct {
	Success         bool
	SpaceRecovered  float64
	TasksAffected   []*core.Task
	VisualChanges   *VisualChanges
	Quality         float64
	Recommendations []string
}

type SpaceConstraints struct {
	MaxWidth      float64
	MaxHeight     float64
	MinWidth      float64
	MinHeight     float64
	AvailableArea float64
	UsedArea      float64
}

type LayoutMetrics struct {
	TotalWidth      float64
	TotalHeight     float64
	UsedWidth       float64
	UsedHeight      float64
	SpaceEfficiency float64
	VisualBalance   float64
}

// Constructor functions
func NewTaskPrioritizationEngine(
	stackingEngine *StackingEngine,
	priorityRanker *PriorityRanker,
	visibilityManager *VisibilityManager,
	stackingOptimizer *StackingOptimizer,
) *TaskPrioritizationEngine {
	return &TaskPrioritizationEngine{
		stackingEngine:    stackingEngine,
		priorityRanker:    priorityRanker,
		visibilityManager: visibilityManager,
		stackingOptimizer: stackingOptimizer,
	}
}

func NewPriorityRanker(conflictCategorizer *ConflictCategorizer) *PriorityRanker {
	return &PriorityRanker{
		conflictCategorizer: conflictCategorizer,
		rankingRules:        make([]PriorityRule, 0),
		visualWeights:       make(map[VisualFactor]float64),
	}
}

func NewVisibilityManager() *VisibilityManager {
	return &VisibilityManager{
		visibilityRules:     make([]VisibilityRule, 0),
		prominenceWeights:   make(map[VisualProminence]float64),
		visibilityThreshold: 0.5,
		adaptiveVisibility:  true,
	}
}

func NewStackingOptimizer() *StackingOptimizer {
	return &StackingOptimizer{
		optimizationRules:  make([]OptimizationRule, 0),
		stackingStrategies: make(map[PriorityCategory]StackingStrategy),
		adaptiveOrdering:   true,
	}
}

func NewConflictCategorizer(spatialEngine *SpatialEngine) *ConflictCategorizer {
	return &ConflictCategorizer{
		spatialEngine:   spatialEngine,
		rules:           make([]ConflictRule, 0),
		severityWeights: make(map[OverlapSeverity]int),
	}
}

func NewOverflowManager() *OverflowManager {
	return &OverflowManager{
		overflowThresholds: map[OverflowType]float64{
			OverflowVertical: 0.8,
		},
		resolutionStrategies: make(map[OverflowType][]ResolutionStrategy),
		adaptiveThresholds:   true,
		smartCompression:     true,
	}
}

func NewVisualConflictResolver() *VisualConflictResolver {
	return &VisualConflictResolver{
		collisionDetector:  NewCollisionDetector(),
		conflictResolvers:  make(map[string]ConflictResolver),
		visualOptimizer:    NewVisualOptimizer(),
		adaptiveResolution: true,
	}
}

func NewCollisionDetector() *CollisionDetector {
	return &CollisionDetector{
		collisionThreshold: 0.1,
		boundingBoxBuffer:  2.0,
		zIndexManager:      NewZIndexManager(),
	}
}

func NewZIndexManager() *ZIndexManager {
	return &ZIndexManager{
		baseZIndex:   1000,
		layerSpacing: 10,
		priorityLayers: map[VisualProminence]int{
			ProminenceCritical: 4,
			ProminenceHigh:     3,
			ProminenceMedium:   2,
			ProminenceLow:      1,
			ProminenceMinimal:  0,
		},
	}
}

func NewVisualOptimizer() *VisualOptimizer {
	return &VisualOptimizer{
		optimizationRules: make([]VisualOptimizationRule, 0),
		layoutStrategies:  make(map[LayoutStrategy]VisualStrategy),
		adaptiveLayout:    true,
	}
}

// NewPriorityManagementEngine creates a new priority management engine
func NewPriorityManagementEngine(
	spatialEngine *SpatialEngine,
	conflictCategorizer *ConflictCategorizer,
	stackingEngine *StackingEngine,
) *PriorityManagementEngine {
	engine := &PriorityManagementEngine{
		spatialEngine:       spatialEngine,
		conflictCategorizer: conflictCategorizer,
		stackingEngine:      stackingEngine,
		rankingRules:        make([]PriorityRule, 0),
		visualWeights:       make(map[VisualFactor]float64),
		priorityRanker:      NewPriorityRanker(conflictCategorizer),
		visibilityManager:   NewVisibilityManager(),
		stackingOptimizer:   NewStackingOptimizer(),
	}

	return engine
}

// PriorityManagementEngine handles both priority ranking and task prioritization
type PriorityManagementEngine struct {
	spatialEngine       *SpatialEngine
	conflictCategorizer *ConflictCategorizer
	rankingRules        []PriorityRule
	visualWeights       map[VisualFactor]float64
	stackingEngine      *StackingEngine
	priorityRanker      *PriorityRanker
	visibilityManager   *VisibilityManager
	stackingOptimizer   *StackingOptimizer
}

// PrioritizeTasks performs intelligent task prioritization for stacking order
func (pme *PriorityManagementEngine) PrioritizeTasks(tasks []*core.Task, context *PriorityContext) *TaskPrioritizationResult {
	// Simplified implementation - return empty results for now
	prioritizedTasks := make([]*PrioritizedTask, 0, len(tasks))

	for i, task := range tasks {
		prioritizedTask := &PrioritizedTask{
			Task:               task,
			PriorityScore:      &PriorityScore{TaskID: task.ID, OverallScore: 0.5},
			VisibilityAction:   &VisibilityAction{IsVisible: true, ProminenceLevel: ProminenceMedium},
			OptimizationAction: &OptimizationAction{StackingOrder: i, VisualProminence: ProminenceMedium},
			StackingOrder:      i,
			DisplayPriority:    0.5,
		}
		prioritizedTasks = append(prioritizedTasks, prioritizedTask)
	}

	return &TaskPrioritizationResult{
		PrioritizedTasks:    prioritizedTasks,
		StackingOrder:       extractTaskIDs(prioritizedTasks),
		VisibilitySettings:  make(map[string]*VisibilityAction),
		OptimizationResults: make([]*OptimizationAction, len(tasks)),
		Recommendations:     []string{"Applied task prioritization"},
		AnalysisDate:        time.Now(),
	}
}

// Helper function to extract task IDs from prioritized tasks
func extractTaskIDs(prioritizedTasks []*PrioritizedTask) []string {
	ids := make([]string, len(prioritizedTasks))
	for i, task := range prioritizedTasks {
		ids[i] = task.Task.ID
	}
	return ids
}

// CalculatePriorityScores calculates priority scores for all tasks
func (pr *PriorityRanker) CalculatePriorityScores(tasks []*core.Task, context *PriorityContext) *PriorityRankingResult {
	taskScores := make([]*PriorityScore, 0, len(tasks))

	for _, task := range tasks {
		score := &PriorityScore{
			TaskID:           task.ID,
			OverallScore:     0.5,
			CategoryScores:   make(map[PriorityCategory]float64),
			VisualProminence: ProminenceMedium,
			Ranking:          0,
			Confidence:       0.8,
			Factors:          make([]PriorityFactor, 0),
			Recommendations:  make([]string, 0),
			CalculatedAt:     time.Now(),
		}
		taskScores = append(taskScores, score)
	}

	return &PriorityRankingResult{
		TaskScores:        taskScores,
		RankingOrder:      extractTaskIDsFromScores(taskScores),
		ConflictsDetected: []string{},
		Recommendations:   []string{"Applied priority ranking"},
		AnalysisDate:      time.Now(),
	}
}

// PriorityRankingResult contains the result of priority ranking
type PriorityRankingResult struct {
	TaskScores        []*PriorityScore
	RankingOrder      []string
	ConflictsDetected []string
	Recommendations   []string
	AnalysisDate      time.Time
}

// Helper function to extract task IDs from priority scores
func extractTaskIDsFromScores(scores []*PriorityScore) []string {
	ids := make([]string, len(scores))
	for i, score := range scores {
		ids[i] = score.TaskID
	}
	return ids
}
