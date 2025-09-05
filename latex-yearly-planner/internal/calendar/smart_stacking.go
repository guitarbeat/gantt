package calendar

import (
	"fmt"
	"math"
	"sort"
	"time"

	"latex-yearly-planner/internal/data"
)

// SmartStackingEngine handles intelligent stacking of overlapping tasks
type SmartStackingEngine struct {
	overlapDetector     *OverlapDetector
	conflictCategorizer *ConflictCategorizer
	priorityRanker      *PriorityRanker
	stackingRules       []StackingRule
	visualConstraints   *VisualConstraints
}

// StackingRule defines a rule for task stacking behavior
type StackingRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*data.Task, *StackingContext) bool
	Action      func(*data.Task, *StackingContext) *StackingAction
}

// StackingAction defines how a task should be stacked
type StackingAction struct {
	StackingType    StackingType
	VerticalOffset  float64
	HorizontalOffset float64
	Height          float64
	Width           float64
	ZIndex          int
	VisualStyle     *VisualStyle
	CollisionAvoidance bool
	Priority        int
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
	MaxStackHeight      float64
	MinTaskHeight       float64
	MaxTaskHeight       float64
	MinTaskWidth        float64
	MaxTaskWidth        float64
	VerticalSpacing     float64
	HorizontalSpacing   float64
	MaxStackDepth       int
	CollisionThreshold  float64
	OverflowThreshold   float64
}

// StackingContext provides context for stacking decisions
type StackingContext struct {
	CalendarStart       time.Time
	CalendarEnd         time.Time
	CurrentTime         time.Time
	DayWidth            float64
	DayHeight           float64
	AvailableHeight     float64
	AvailableWidth      float64
	ExistingStacks      []*TaskStack
	TaskPriorities      map[string]*TaskPriority
	ConflictAnalysis    *ConflictAnalysis
	OverlapAnalysis     *OverlapAnalysis
	VisualSettings      *VisualSettings
	VisualConstraints   *VisualConstraints
}

// TaskStack represents a stack of overlapping tasks
type TaskStack struct {
	ID              string
	Tasks           []*StackedTask
	StartTime       time.Time
	EndTime         time.Time
	TotalHeight     float64
	MaxWidth        float64
	StackingType    StackingType
	Priority        int
	CollisionCount  int
	OverflowCount   int
	VisualStyle     *VisualStyle
}

// StackedTask represents a task within a stack
type StackedTask struct {
	Task            *data.Task
	StackingAction  *StackingAction
	Position        *Position
	IsVisible       bool
	IsCollapsed     bool
	CollisionLevel  int
	OverflowLevel   int
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
	ShowTaskNames      bool
	ShowTaskDurations  bool
	ShowTaskPriorities bool
	ShowConflictIndicators bool
	CollapseThreshold  int
	AnimationEnabled   bool
	HighlightConflicts bool
	ColorScheme        string
}

// StackingResult contains the result of stacking operations
type StackingResult struct {
	Stacks           []*TaskStack
	TotalStacks      int
	CollisionCount   int
	OverflowCount    int
	SpaceEfficiency  float64
	VisualQuality    float64
	Recommendations  []string
	AnalysisDate     time.Time
}

// NewSmartStackingEngine creates a new smart stacking engine
func NewSmartStackingEngine(overlapDetector *OverlapDetector, conflictCategorizer *ConflictCategorizer, priorityRanker *PriorityRanker) *SmartStackingEngine {
	engine := &SmartStackingEngine{
		overlapDetector:     overlapDetector,
		conflictCategorizer: conflictCategorizer,
		priorityRanker:      priorityRanker,
		stackingRules:       make([]StackingRule, 0),
		visualConstraints:   &VisualConstraints{
			MaxStackHeight:     100.0,
			MinTaskHeight:      20.0,
			MaxTaskHeight:      40.0,
			MinTaskWidth:       50.0,
			MaxTaskWidth:       200.0,
			VerticalSpacing:    2.0,
			HorizontalSpacing:  5.0,
			MaxStackDepth:      10,
			CollisionThreshold: 0.1,
			OverflowThreshold:  0.8,
		},
	}
	
	// Initialize default stacking rules
	engine.initializeDefaultRules()
	
	return engine
}

// initializeDefaultRules sets up default stacking rules
func (sse *SmartStackingEngine) initializeDefaultRules() {
	sse.stackingRules = []StackingRule{
		{
			Name:        "Critical Task Priority",
			Description: "Critical tasks get top priority in stacking",
			Priority:    1,
			Condition: func(task *data.Task, context *StackingContext) bool {
				if priority, exists := context.TaskPriorities[task.ID]; exists {
					return priority.VisualProminence == ProminenceCritical
				}
				return false
			},
			Action: func(task *data.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeLayered,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MaxTaskHeight,
					Width:              context.VisualConstraints.MaxTaskWidth,
					ZIndex:             10,
					CollisionAvoidance: true,
					Priority:           1,
				}
			},
		},
		{
			Name:        "High Priority Task",
			Description: "High priority tasks get prominent stacking",
			Priority:    2,
			Condition: func(task *data.Task, context *StackingContext) bool {
				if priority, exists := context.TaskPriorities[task.ID]; exists {
					return priority.VisualProminence == ProminenceHigh
				}
				return false
			},
			Action: func(task *data.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeVertical,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MaxTaskHeight * 0.9,
					Width:              context.VisualConstraints.MaxTaskWidth * 0.9,
					ZIndex:             8,
					CollisionAvoidance: true,
					Priority:           2,
				}
			},
		},
		{
			Name:        "Milestone Task",
			Description: "Milestone tasks get special stacking treatment",
			Priority:    3,
			Condition: func(task *data.Task, context *StackingContext) bool {
				return task.IsMilestone
			},
			Action: func(task *data.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeFloating,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MaxTaskHeight,
					Width:              context.VisualConstraints.MaxTaskWidth,
					ZIndex:             9,
					CollisionAvoidance: true,
					Priority:           3,
				}
			},
		},
		{
			Name:        "Long Duration Task",
			Description: "Long duration tasks get horizontal stacking",
			Priority:    4,
			Condition: func(task *data.Task, context *StackingContext) bool {
				duration := task.EndDate.Sub(task.StartDate)
				return duration > time.Hour*24*7 // More than a week
			},
			Action: func(task *data.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeHorizontal,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MinTaskHeight,
					Width:              context.VisualConstraints.MaxTaskWidth,
					ZIndex:             5,
					CollisionAvoidance: false,
					Priority:           4,
				}
			},
		},
		{
			Name:        "Short Duration Task",
			Description: "Short duration tasks get vertical stacking",
			Priority:    5,
			Condition: func(task *data.Task, context *StackingContext) bool {
				duration := task.EndDate.Sub(task.StartDate)
				return duration <= time.Hour*24 // One day or less
			},
			Action: func(task *data.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeVertical,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MinTaskHeight,
					Width:              context.VisualConstraints.MinTaskWidth,
					ZIndex:             3,
					CollisionAvoidance: false,
					Priority:           5,
				}
			},
		},
		{
			Name:        "Conflict Resolution",
			Description: "Tasks with conflicts get special stacking treatment",
			Priority:    6,
			Condition: func(task *data.Task, context *StackingContext) bool {
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
			Action: func(task *data.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeCascading,
					VerticalOffset:     5.0,
					HorizontalOffset:   5.0,
					Height:             context.VisualConstraints.MaxTaskHeight * 0.8,
					Width:              context.VisualConstraints.MaxTaskWidth * 0.8,
					ZIndex:             7,
					CollisionAvoidance: true,
					Priority:           6,
				}
			},
		},
		{
			Name:        "Overflow Handling",
			Description: "Tasks that cause overflow get minimized stacking",
			Priority:    7,
			Condition: func(task *data.Task, context *StackingContext) bool {
				// Check if adding this task would cause overflow
				return sse.wouldCauseOverflow(task, context)
			},
			Action: func(task *data.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeMinimized,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MinTaskHeight * 0.5,
					Width:              context.VisualConstraints.MinTaskWidth * 0.5,
					ZIndex:             1,
					CollisionAvoidance: false,
					Priority:           7,
				}
			},
		},
		{
			Name:        "Default Stacking",
			Description: "Default stacking for all other tasks",
			Priority:    8,
			Condition: func(task *data.Task, context *StackingContext) bool {
				return true // Always matches
			},
			Action: func(task *data.Task, context *StackingContext) *StackingAction {
				return &StackingAction{
					StackingType:       StackingTypeVertical,
					VerticalOffset:     0.0,
					HorizontalOffset:   0.0,
					Height:             context.VisualConstraints.MinTaskHeight,
					Width:              context.VisualConstraints.MinTaskWidth,
					ZIndex:             2,
					CollisionAvoidance: false,
					Priority:           8,
				}
			},
		},
	}
}

// wouldCauseOverflow checks if adding a task would cause visual overflow
func (sse *SmartStackingEngine) wouldCauseOverflow(task *data.Task, context *StackingContext) bool {
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
func (sse *SmartStackingEngine) StackTasks(tasks []*data.Task, context *StackingContext) *StackingResult {
	// Detect overlaps and categorize conflicts
	overlapAnalysis := sse.overlapDetector.DetectOverlaps(tasks)
	conflictAnalysis := sse.conflictCategorizer.CategorizeConflicts(overlapAnalysis)
	
	// Rank tasks by priority
	priorityContext := &PriorityContext{
		CalendarStart:      context.CalendarStart,
		CalendarEnd:        context.CalendarEnd,
		CurrentTime:        context.CurrentTime,
		AssigneeWorkloads:  make(map[string]int),
		CategoryImportance: make(map[string]float64),
	}
	priorityRanking := sse.priorityRanker.RankTasks(tasks, priorityContext)
	
	// Update context with analysis results
	context.OverlapAnalysis = overlapAnalysis
	context.ConflictAnalysis = conflictAnalysis
	context.TaskPriorities = make(map[string]*TaskPriority)
	for _, taskPriority := range priorityRanking.TaskPriorities {
		context.TaskPriorities[taskPriority.Task.ID] = taskPriority
	}
	
	// Group tasks by overlapping time periods
	overlapGroups := sse.groupTasksByOverlap(tasks, overlapAnalysis)
	
	// Create stacks for each overlap group
	var stacks []*TaskStack
	for _, group := range overlapGroups {
		stack := sse.createStackForGroup(group, context)
		if stack != nil {
			stacks = append(stacks, stack)
		}
	}
	
	// Calculate stacking metrics
	result := &StackingResult{
		Stacks:          stacks,
		TotalStacks:     len(stacks),
		CollisionCount:  sse.calculateCollisionCount(stacks),
		OverflowCount:   sse.calculateOverflowCount(stacks, context),
		SpaceEfficiency: sse.calculateSpaceEfficiency(stacks, context),
		VisualQuality:   sse.calculateVisualQuality(stacks, context),
		Recommendations: sse.generateRecommendations(stacks, context),
		AnalysisDate:    time.Now(),
	}
	
	return result
}

// groupTasksByOverlap groups tasks by their overlapping time periods
func (sse *SmartStackingEngine) groupTasksByOverlap(tasks []*data.Task, overlapAnalysis *OverlapAnalysis) [][]*data.Task {
	var groups [][]*data.Task
	
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
			groups = append(groups, []*data.Task{task})
		}
	}
	
	return groups
}

// createStackForGroup creates a stack for a group of overlapping tasks
func (sse *SmartStackingEngine) createStackForGroup(tasks []*data.Task, context *StackingContext) *TaskStack {
	if len(tasks) == 0 {
		return nil
	}
	
	// Sort tasks by priority
	sort.Slice(tasks, func(i, j int) bool {
		priorityI := context.TaskPriorities[tasks[i].ID]
		priorityJ := context.TaskPriorities[tasks[j].ID]
		return priorityI.PriorityScore > priorityJ.PriorityScore
	})
	
	// Create stack
	stack := &TaskStack{
		ID:           fmt.Sprintf("stack_%d", len(context.ExistingStacks)+1),
		Tasks:        make([]*StackedTask, 0),
		StartTime:    tasks[0].StartDate,
		EndTime:      tasks[0].EndDate,
		TotalHeight:  0.0,
		MaxWidth:     0.0,
		StackingType: StackingTypeVertical,
		Priority:     0,
		CollisionCount: 0,
		OverflowCount: 0,
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
		stackingAction := sse.determineStackingAction(task, context)
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
		if sse.hasCollision(stackedTask, stack.Tasks) {
			stackedTask.CollisionLevel = 1
			stack.CollisionCount++
		}
		
		// Check for overflow
		if sse.hasOverflow(stackedTask, context) {
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
	stack.StackingType = sse.determineStackType(stack)
	
	// Set stack priority
	if len(stack.Tasks) > 0 {
		stack.Priority = stack.Tasks[0].StackingAction.Priority
	}
	
	return stack
}

// determineStackingAction determines the stacking action for a task
func (sse *SmartStackingEngine) determineStackingAction(task *data.Task, context *StackingContext) *StackingAction {
	// Find the best matching rule
	for _, rule := range sse.stackingRules {
		if rule.Condition(task, context) {
			action := rule.Action(task, context)
			action.Priority = rule.Priority
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
		Priority:           8,
	}
}

// determineStackType determines the overall stack type based on tasks
func (sse *SmartStackingEngine) determineStackType(stack *TaskStack) StackingType {
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
func (sse *SmartStackingEngine) hasCollision(newTask *StackedTask, existingTasks []*StackedTask) bool {
	for _, existingTask := range existingTasks {
		if sse.tasksCollide(newTask, existingTask) {
			return true
		}
	}
	return false
}

// tasksCollide checks if two stacked tasks collide
func (sse *SmartStackingEngine) tasksCollide(task1, task2 *StackedTask) bool {
	// Simple bounding box collision detection
	return !(task1.Position.X+task1.Position.Width < task2.Position.X ||
		task2.Position.X+task2.Position.Width < task1.Position.X ||
		task1.Position.Y+task1.Position.Height < task2.Position.Y ||
		task2.Position.Y+task2.Position.Height < task1.Position.Y)
}

// hasOverflow checks if a stacked task causes overflow
func (sse *SmartStackingEngine) hasOverflow(stackedTask *StackedTask, context *StackingContext) bool {
	return stackedTask.Position.Y+stackedTask.Position.Height > context.AvailableHeight*context.VisualConstraints.OverflowThreshold
}

// calculateCollisionCount calculates total collision count
func (sse *SmartStackingEngine) calculateCollisionCount(stacks []*TaskStack) int {
	total := 0
	for _, stack := range stacks {
		total += stack.CollisionCount
	}
	return total
}

// calculateOverflowCount calculates total overflow count
func (sse *SmartStackingEngine) calculateOverflowCount(stacks []*TaskStack, context *StackingContext) int {
	total := 0
	for _, stack := range stacks {
		total += stack.OverflowCount
	}
	return total
}

// calculateSpaceEfficiency calculates space efficiency
func (sse *SmartStackingEngine) calculateSpaceEfficiency(stacks []*TaskStack, context *StackingContext) float64 {
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
func (sse *SmartStackingEngine) calculateVisualQuality(stacks []*TaskStack, context *StackingContext) float64 {
	if len(stacks) == 0 {
		return 1.0
	}
	
	// Calculate quality based on collision and overflow counts
	totalCollisions := sse.calculateCollisionCount(stacks)
	totalOverflows := sse.calculateOverflowCount(stacks, context)
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
func (sse *SmartStackingEngine) generateRecommendations(stacks []*TaskStack, context *StackingContext) []string {
	var recommendations []string
	
	// Collision recommendations
	totalCollisions := sse.calculateCollisionCount(stacks)
	if totalCollisions > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("⚠️ %d visual collisions detected - consider adjusting task positioning", totalCollisions))
	}
	
	// Overflow recommendations
	totalOverflows := sse.calculateOverflowCount(stacks, context)
	if totalOverflows > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("📏 %d overflow issues detected - consider reducing task sizes or using minimization", totalOverflows))
	}
	
	// Space efficiency recommendations
	efficiency := sse.calculateSpaceEfficiency(stacks, context)
	if efficiency < 0.5 {
		recommendations = append(recommendations, 
			"📊 Low space efficiency - consider optimizing task layouts")
	} else if efficiency > 0.9 {
		recommendations = append(recommendations, 
			"📊 High space efficiency - good use of available space")
	}
	
	// Visual quality recommendations
	quality := sse.calculateVisualQuality(stacks, context)
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
func (sse *SmartStackingEngine) AddCustomRule(rule StackingRule) {
	sse.stackingRules = append(sse.stackingRules, rule)
	// Sort rules by priority (highest first)
	sort.Slice(sse.stackingRules, func(i, j int) bool {
		return sse.stackingRules[i].Priority < sse.stackingRules[j].Priority
	})
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

// GetStacksByPriority returns stacks filtered by priority
func (result *StackingResult) GetStacksByPriority(priority int) []*TaskStack {
	var filtered []*TaskStack
	for _, stack := range result.Stacks {
		if stack.Priority == priority {
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
