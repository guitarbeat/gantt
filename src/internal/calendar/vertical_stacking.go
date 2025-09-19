package calendar

import (
	"fmt"
	"math"
	"time"

	"phd-dissertation-planner/internal/data"
)

// VerticalStackingEngine handles vertical stacking of tasks with intelligent positioning
type VerticalStackingEngine struct {
	smartStackingEngine *SmartStackingEngine
	heightCalculator    *HeightCalculator
	positionCalculator  *PositionCalculator
	spaceOptimizer      *SpaceOptimizer
}

// HeightCalculator calculates optimal heights for stacked tasks
type HeightCalculator struct {
	baseHeight        float64
	minHeight         float64
	maxHeight         float64
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
	AlignmentTop      AlignmentMode = "TOP"
	AlignmentCenter   AlignmentMode = "CENTER"
	AlignmentBottom   AlignmentMode = "BOTTOM"
	AlignmentJustify  AlignmentMode = "JUSTIFY"
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
	ID              string
	Tasks           []*VerticallyStackedTask
	StartTime       time.Time
	EndTime         time.Time
	TotalHeight     float64
	MaxWidth        float64
	AlignmentMode   AlignmentMode
	DistributionMode DistributionMode
	SpaceEfficiency float64
	VisualBalance   float64
	CollisionCount  int
	OverflowCount   int
}

// VerticallyStackedTask represents a task within a vertical stack
type VerticallyStackedTask struct {
	Task            *data.Task
	Position        *VerticalPosition
	CalculatedHeight float64
	IsCompressed    bool
	IsExpanded      bool
	CollisionLevel  int
	OverflowLevel   int
	VisualWeight    float64
}

// VerticalPosition represents the vertical position of a stacked task
type VerticalPosition struct {
	X           float64
	Y           float64
	Width       float64
	Height      float64
	ZIndex      int
	OffsetY     float64
	RelativeY   float64
	StackIndex  int
}

// NewVerticalStackingEngine creates a new vertical stacking engine
func NewVerticalStackingEngine(smartStackingEngine *SmartStackingEngine) *VerticalStackingEngine {
	engine := &VerticalStackingEngine{
		smartStackingEngine: smartStackingEngine,
		heightCalculator:    NewHeightCalculator(),
		positionCalculator:  NewPositionCalculator(),
		spaceOptimizer:      NewSpaceOptimizer(),
	}
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
			ProminenceLow:      0.8,
			ProminenceMinimal:  0.6,
		},
		durationMultiplier: map[string]float64{
			"short":  0.8,  // < 1 day
			"medium": 1.0,  // 1-7 days
			"long":   1.2,  // > 7 days
		},
		contentMultiplier: map[string]float64{
			"minimal": 0.7,  // Simple tasks
			"normal":  1.0,  // Standard tasks
			"complex": 1.3,  // Complex tasks with many details
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
		smartCollapsing:      true,
	}
}

// StackTasksVertically performs vertical stacking of overlapping tasks
func (vse *VerticalStackingEngine) StackTasksVertically(tasks []*data.Task, context *StackingContext) *VerticalStackingResult {
	// First, use the smart stacking engine to get initial stacks
	smartResult := vse.smartStackingEngine.StackTasks(tasks, context)
	
	// Convert smart stacks to vertical stacks
	var verticalStacks []*VerticalStack
	for _, smartStack := range smartResult.Stacks {
		verticalStack := vse.convertToVerticalStack(smartStack, context)
		if verticalStack != nil {
			verticalStacks = append(verticalStacks, verticalStack)
		}
	}
	
	// Optimize vertical stacking
	verticalStacks = vse.optimizeVerticalStacking(verticalStacks, context)
	
	// Calculate metrics
	result := &VerticalStackingResult{
		Stacks:          verticalStacks,
		TotalHeight:     vse.calculateTotalHeight(verticalStacks),
		SpaceEfficiency: vse.calculateSpaceEfficiency(verticalStacks, context),
		VisualBalance:   vse.calculateVisualBalance(verticalStacks),
		CollisionCount:  vse.calculateCollisionCount(verticalStacks),
		OverflowCount:   vse.calculateOverflowCount(verticalStacks, context),
		CompressionRatio: vse.calculateCompressionRatio(verticalStacks),
		Recommendations: vse.generateRecommendations(verticalStacks, context),
		AnalysisDate:    time.Now(),
	}
	
	return result
}

// convertToVerticalStack converts a smart stack to a vertical stack
func (vse *VerticalStackingEngine) convertToVerticalStack(smartStack *TaskStack, context *StackingContext) *VerticalStack {
	if len(smartStack.Tasks) == 0 {
		return nil
	}
	
	verticalStack := &VerticalStack{
		ID:               smartStack.ID,
		StartTime:        smartStack.StartTime,
		EndTime:          smartStack.EndTime,
		MaxWidth:         smartStack.MaxWidth,
		AlignmentMode:    vse.determineAlignmentMode(smartStack, context),
		DistributionMode: vse.determineDistributionMode(smartStack, context),
		Tasks:            make([]*VerticallyStackedTask, 0),
	}
	
	// Convert each stacked task
	for i, stackedTask := range smartStack.Tasks {
		verticallyStackedTask := &VerticallyStackedTask{
			Task:            stackedTask.Task,
			CalculatedHeight: vse.calculateTaskHeight(stackedTask.Task, context),
			IsCompressed:    false,
			IsExpanded:      false,
			CollisionLevel:  stackedTask.CollisionLevel,
			OverflowLevel:   stackedTask.OverflowLevel,
			VisualWeight:    vse.calculateVisualWeight(stackedTask.Task, context),
		}
		
		// Calculate position
		verticallyStackedTask.Position = vse.calculateVerticalPosition(
			verticallyStackedTask, 
			verticalStack, 
			i, 
			context,
		)
		
		verticalStack.Tasks = append(verticalStack.Tasks, verticallyStackedTask)
	}
	
	// Calculate stack metrics
	verticalStack.TotalHeight = vse.calculateStackHeight(verticalStack)
	verticalStack.SpaceEfficiency = vse.calculateStackSpaceEfficiency(verticalStack, context)
	verticalStack.VisualBalance = vse.calculateStackVisualBalance(verticalStack)
	verticalStack.CollisionCount = vse.calculateStackCollisionCount(verticalStack)
	verticalStack.OverflowCount = vse.calculateStackOverflowCount(verticalStack, context)
	
	return verticalStack
}

// calculateTaskHeight calculates the optimal height for a task
func (vse *VerticalStackingEngine) calculateTaskHeight(task *data.Task, context *StackingContext) float64 {
	hc := vse.heightCalculator
	
	// Start with base height
	height := hc.baseHeight
	
	// Apply priority multiplier
	if priority, exists := context.TaskPriorities[task.ID]; exists {
		if multiplier, exists := hc.priorityMultiplier[priority.VisualProminence]; exists {
			height *= multiplier
		}
	}
	
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
	contentComplexity := vse.assessContentComplexity(task)
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
func (vse *VerticalStackingEngine) assessContentComplexity(task *data.Task) string {
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
func (vse *VerticalStackingEngine) calculateVisualWeight(task *data.Task, context *StackingContext) float64 {
	weight := 1.0
	
	// Priority weight
	if priority, exists := context.TaskPriorities[task.ID]; exists {
		weight += priority.PriorityScore * 0.1
	}
	
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
func (vse *VerticalStackingEngine) calculateVerticalPosition(
	task *VerticallyStackedTask, 
	stack *VerticalStack, 
	index int, 
	context *StackingContext,
) *VerticalPosition {
	pc := vse.positionCalculator
	
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
	position = vse.applyAlignmentMode(position, stack, context)
	
	// Apply distribution mode
	position = vse.applyDistributionMode(position, stack, context)
	
	// Calculate relative position within stack
	position.RelativeY = position.Y - stack.StartTime.Sub(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).Hours() * 10
	
	return position
}

// applyAlignmentMode applies the alignment mode to a position
func (vse *VerticalStackingEngine) applyAlignmentMode(position *VerticalPosition, stack *VerticalStack, context *StackingContext) *VerticalPosition {
	pc := vse.positionCalculator
	
	switch pc.alignmentMode {
	case AlignmentTop:
		// Already positioned at top
		break
	case AlignmentCenter:
		// Center within available height
		if context.AvailableHeight > 0 {
			stackHeight := vse.calculateStackHeight(stack)
			offset := (context.AvailableHeight - stackHeight) / 2
			position.Y += offset
		}
	case AlignmentBottom:
		// Position at bottom
		if context.AvailableHeight > 0 {
			stackHeight := vse.calculateStackHeight(stack)
			position.Y = context.AvailableHeight - stackHeight + position.Y
		}
	case AlignmentJustify:
		// Justify within available space
		if context.AvailableHeight > 0 {
			stackHeight := vse.calculateStackHeight(stack)
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
func (vse *VerticalStackingEngine) applyDistributionMode(position *VerticalPosition, stack *VerticalStack, context *StackingContext) *VerticalPosition {
	pc := vse.positionCalculator
	
	switch pc.distributionMode {
	case DistributionEven:
		// Even distribution (already handled in base calculation)
		break
	case DistributionPriority:
		// Distribute based on priority
		if len(stack.Tasks) > position.StackIndex {
			if priority, exists := context.TaskPriorities[stack.Tasks[position.StackIndex].Task.ID]; exists {
				priorityOffset := priority.PriorityScore * 0.5
				position.Y += priorityOffset
			}
		}
	case DistributionContent:
		// Distribute based on content complexity
		if len(stack.Tasks) > position.StackIndex {
			contentComplexity := vse.assessContentComplexity(stack.Tasks[position.StackIndex].Task)
			if contentComplexity == "complex" {
				position.Y += 5.0
			} else if contentComplexity == "minimal" {
				position.Y -= 2.0
			}
		}
	case DistributionAdaptive:
		// Adaptive distribution based on available space
		if context.AvailableHeight > 0 {
			stackHeight := vse.calculateStackHeight(stack)
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
func (vse *VerticalStackingEngine) determineAlignmentMode(stack *TaskStack, context *StackingContext) AlignmentMode {
	// Check if stack has critical tasks
	for _, task := range stack.Tasks {
		if priority, exists := context.TaskPriorities[task.Task.ID]; exists {
			if priority.VisualProminence == ProminenceCritical {
				return AlignmentTop
			}
		}
	}
	
	// Check available space
	if context.AvailableHeight > 0 {
		estimatedHeight := float64(len(stack.Tasks)) * 25.0 // Rough estimate
		if estimatedHeight < context.AvailableHeight * 0.5 {
			return AlignmentCenter
		}
	}
	
	return AlignmentTop
}

// determineDistributionMode determines the best distribution mode for a stack
func (vse *VerticalStackingEngine) determineDistributionMode(stack *TaskStack, context *StackingContext) DistributionMode {
	// Check if stack has mixed priorities
	hasHighPriority := false
	hasLowPriority := false
	
	for _, task := range stack.Tasks {
		if priority, exists := context.TaskPriorities[task.Task.ID]; exists {
			if priority.VisualProminence == ProminenceCritical || priority.VisualProminence == ProminenceHigh {
				hasHighPriority = true
			} else if priority.VisualProminence == ProminenceLow || priority.VisualProminence == ProminenceMinimal {
				hasLowPriority = true
			}
		}
	}
	
	if hasHighPriority && hasLowPriority {
		return DistributionPriority
	}
	
	// Check if stack has mixed content complexity
	hasComplexContent := false
	hasSimpleContent := false
	
	for _, task := range stack.Tasks {
		complexity := vse.assessContentComplexity(task.Task)
		if complexity == "complex" {
			hasComplexContent = true
		} else if complexity == "minimal" {
			hasSimpleContent = true
		}
	}
	
	if hasComplexContent && hasSimpleContent {
		return DistributionContent
	}
	
	return DistributionEven
}

// optimizeVerticalStacking optimizes the vertical stacking layout
func (vse *VerticalStackingEngine) optimizeVerticalStacking(stacks []*VerticalStack, context *StackingContext) []*VerticalStack {
	so := vse.spaceOptimizer
	
	// Apply space optimization
	for _, stack := range stacks {
		// Check if compression is needed
		if vse.needsCompression(stack, context) {
			stack = vse.compressStack(stack, context)
		}
		
		// Check if expansion is possible
		if vse.canExpand(stack, context) {
			stack = vse.expandStack(stack, context)
		}
		
		// Apply adaptive spacing
		if so.adaptiveSpacing {
			stack = vse.applyAdaptiveSpacing(stack, context)
		}
		
		// Apply smart collapsing
		if so.smartCollapsing {
			stack = vse.applySmartCollapsing(stack, context)
		}
	}
	
	return stacks
}

// needsCompression checks if a stack needs compression
func (vse *VerticalStackingEngine) needsCompression(stack *VerticalStack, context *StackingContext) bool {
	if context.AvailableHeight <= 0 {
		return false
	}
	
	stackHeight := vse.calculateStackHeight(stack)
	return stackHeight > context.AvailableHeight * vse.spaceOptimizer.compressionThreshold
}

// canExpand checks if a stack can be expanded
func (vse *VerticalStackingEngine) canExpand(stack *VerticalStack, context *StackingContext) bool {
	if context.AvailableHeight <= 0 {
		return false
	}
	
	stackHeight := vse.calculateStackHeight(stack)
	return stackHeight < context.AvailableHeight * vse.spaceOptimizer.expansionThreshold
}

// compressStack compresses a stack to fit available space
func (vse *VerticalStackingEngine) compressStack(stack *VerticalStack, context *StackingContext) *VerticalStack {
	if context.AvailableHeight <= 0 {
		return stack
	}
	
	// Calculate compression ratio
	currentHeight := vse.calculateStackHeight(stack)
	compressionRatio := context.AvailableHeight / currentHeight
	
	// Apply compression to each task
	for _, task := range stack.Tasks {
		task.CalculatedHeight *= compressionRatio
		task.IsCompressed = true
		task.Position.Height = task.CalculatedHeight
	}
	
	// Recalculate positions
	vse.recalculateStackPositions(stack)
	
	return stack
}

// expandStack expands a stack to better utilize available space
func (vse *VerticalStackingEngine) expandStack(stack *VerticalStack, context *StackingContext) *VerticalStack {
	if context.AvailableHeight <= 0 {
		return stack
	}
	
	// Calculate expansion ratio
	currentHeight := vse.calculateStackHeight(stack)
	expansionRatio := math.Min(1.5, context.AvailableHeight / currentHeight)
	
	// Apply expansion to each task
	for _, task := range stack.Tasks {
		task.CalculatedHeight *= expansionRatio
		task.IsExpanded = true
		task.Position.Height = task.CalculatedHeight
	}
	
	// Recalculate positions
	vse.recalculateStackPositions(stack)
	
	return stack
}

// applyAdaptiveSpacing applies adaptive spacing to a stack
func (vse *VerticalStackingEngine) applyAdaptiveSpacing(stack *VerticalStack, context *StackingContext) *VerticalStack {
	if context.AvailableHeight <= 0 {
		return stack
	}
	
	// Calculate adaptive spacing
	currentHeight := vse.calculateStackHeight(stack)
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
func (vse *VerticalStackingEngine) applySmartCollapsing(stack *VerticalStack, context *StackingContext) *VerticalStack {
	// Collapse low-priority tasks if space is limited
	if context.AvailableHeight > 0 {
		stackHeight := vse.calculateStackHeight(stack)
		if stackHeight > context.AvailableHeight * 0.9 {
			// Collapse tasks with low visual weight
			for _, task := range stack.Tasks {
				if task.VisualWeight < 1.0 {
					task.CalculatedHeight *= 0.7
					task.Position.Height = task.CalculatedHeight
				}
			}
			
			// Recalculate positions
			vse.recalculateStackPositions(stack)
		}
	}
	
	return stack
}

// recalculateStackPositions recalculates positions after height changes
func (vse *VerticalStackingEngine) recalculateStackPositions(stack *VerticalStack) {
	pc := vse.positionCalculator
	
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
func (vse *VerticalStackingEngine) calculateStackHeight(stack *VerticalStack) float64 {
	if len(stack.Tasks) == 0 {
		return 0.0
	}
	
	lastTask := stack.Tasks[len(stack.Tasks)-1]
	return lastTask.Position.Y + lastTask.Position.Height
}

// calculateStackSpaceEfficiency calculates the space efficiency of a stack
func (vse *VerticalStackingEngine) calculateStackSpaceEfficiency(stack *VerticalStack, context *StackingContext) float64 {
	if context.AvailableHeight <= 0 {
		return 1.0
	}
	
	stackHeight := vse.calculateStackHeight(stack)
	return math.Min(stackHeight / context.AvailableHeight, 1.0)
}

// calculateStackVisualBalance calculates the visual balance of a stack
func (vse *VerticalStackingEngine) calculateStackVisualBalance(stack *VerticalStack) float64 {
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
		deviation := math.Abs(task.VisualWeight - avgWeight) / avgWeight
		balance -= deviation * 0.1
	}
	
	return math.Max(balance, 0.0)
}

// calculateStackCollisionCount calculates the collision count of a stack
func (vse *VerticalStackingEngine) calculateStackCollisionCount(stack *VerticalStack) int {
	count := 0
	for _, task := range stack.Tasks {
		count += task.CollisionLevel
	}
	return count
}

// calculateStackOverflowCount calculates the overflow count of a stack
func (vse *VerticalStackingEngine) calculateStackOverflowCount(stack *VerticalStack, context *StackingContext) int {
	count := 0
	for _, task := range stack.Tasks {
		count += task.OverflowLevel
	}
	return count
}

// calculateTotalHeight calculates the total height of all stacks
func (vse *VerticalStackingEngine) calculateTotalHeight(stacks []*VerticalStack) float64 {
	total := 0.0
	for _, stack := range stacks {
		total += vse.calculateStackHeight(stack)
	}
	return total
}

// calculateSpaceEfficiency calculates the overall space efficiency
func (vse *VerticalStackingEngine) calculateSpaceEfficiency(stacks []*VerticalStack, context *StackingContext) float64 {
	if context.AvailableHeight <= 0 {
		return 1.0
	}
	
	totalHeight := vse.calculateTotalHeight(stacks)
	return math.Min(totalHeight / context.AvailableHeight, 1.0)
}

// calculateVisualBalance calculates the overall visual balance
func (vse *VerticalStackingEngine) calculateVisualBalance(stacks []*VerticalStack) float64 {
	if len(stacks) == 0 {
		return 1.0
	}
	
	totalBalance := 0.0
	for _, stack := range stacks {
		totalBalance += vse.calculateStackVisualBalance(stack)
	}
	
	return totalBalance / float64(len(stacks))
}

// calculateCollisionCount calculates the total collision count
func (vse *VerticalStackingEngine) calculateCollisionCount(stacks []*VerticalStack) int {
	total := 0
	for _, stack := range stacks {
		total += vse.calculateStackCollisionCount(stack)
	}
	return total
}

// calculateOverflowCount calculates the total overflow count
func (vse *VerticalStackingEngine) calculateOverflowCount(stacks []*VerticalStack, context *StackingContext) int {
	total := 0
	for _, stack := range stacks {
		total += vse.calculateStackOverflowCount(stack, context)
	}
	return total
}

// calculateCompressionRatio calculates the compression ratio
func (vse *VerticalStackingEngine) calculateCompressionRatio(stacks []*VerticalStack) float64 {
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

// generateRecommendations generates recommendations for vertical stacking
func (vse *VerticalStackingEngine) generateRecommendations(stacks []*VerticalStack, context *StackingContext) []string {
	var recommendations []string
	
	// Space efficiency recommendations
	efficiency := vse.calculateSpaceEfficiency(stacks, context)
	if efficiency < 0.5 {
		recommendations = append(recommendations, 
			"📏 Low space efficiency - consider adjusting task heights or using compression")
	} else if efficiency > 0.9 {
		recommendations = append(recommendations, 
			"📏 High space efficiency - good utilization of available space")
	}
	
	// Visual balance recommendations
	balance := vse.calculateVisualBalance(stacks)
	if balance < 0.7 {
		recommendations = append(recommendations, 
			"⚖️ Visual balance could be improved - consider adjusting task weights")
	}
	
	// Compression recommendations
	compressionRatio := vse.calculateCompressionRatio(stacks)
	if compressionRatio > 0.5 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🗜️ High compression ratio (%.1f%%) - consider increasing available space", compressionRatio*100))
	}
	
	// Collision recommendations
	collisionCount := vse.calculateCollisionCount(stacks)
	if collisionCount > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("⚠️ %d visual collisions detected - consider adjusting task positioning", collisionCount))
	}
	
	// Overflow recommendations
	overflowCount := vse.calculateOverflowCount(stacks, context)
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
