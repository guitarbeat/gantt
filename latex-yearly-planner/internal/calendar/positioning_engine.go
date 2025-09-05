package calendar

import (
	"fmt"
	"math"
	"sort"
	"time"

	"latex-yearly-planner/internal/data"
)

// PositioningEngine handles precise positioning and alignment of tasks within the calendar grid
type PositioningEngine struct {
	gridConfig        *GridConfig
	visualConstraints *VisualConstraints
	alignmentRules    []AlignmentRule
	spacingRules      []SpacingRule
	layoutMetrics     *PositioningLayoutMetrics
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
	CalendarStart      time.Time
	CalendarEnd        time.Time
	DayWidth           float64
	DayHeight          float64
	AvailableHeight    float64
	AvailableWidth     float64
	ExistingBars       []*IntegratedTaskBar
	GridConstraints    *GridConstraints
	VisualSettings     *IntegratedVisualSettings
	CurrentTime        time.Time
	TaskDensity        float64
	OverlapCount       int
	ConflictCount      int
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
	PositioningAlignmentLeft      PositioningAlignmentMode = "LEFT"
	PositioningAlignmentCenter    PositioningAlignmentMode = "CENTER"
	PositioningAlignmentRight     PositioningAlignmentMode = "RIGHT"
	PositioningAlignmentJustify   PositioningAlignmentMode = "JUSTIFY"
	PositioningAlignmentStretch   PositioningAlignmentMode = "STRETCH"
	PositioningAlignmentTop       PositioningAlignmentMode = "TOP"
	PositioningAlignmentMiddle    PositioningAlignmentMode = "MIDDLE"
	PositioningAlignmentBottom    PositioningAlignmentMode = "BOTTOM"
)

// JustificationMode defines how tasks should be justified within available space
type JustificationMode string

const (
	JustifyStart      JustificationMode = "START"
	JustifyEnd        JustificationMode = "END"
	JustifyCenter     JustificationMode = "CENTER"
	JustifySpaceBetween JustificationMode = "SPACE_BETWEEN"
	JustifySpaceAround  JustificationMode = "SPACE_AROUND"
	JustifySpaceEvenly  JustificationMode = "SPACE_EVENLY"
)

// SpacingAction defines spacing adjustments between tasks
type SpacingAction struct {
	VerticalSpacing   float64
	HorizontalSpacing float64
	CollisionAvoidance bool
	OverlapResolution bool
	Priority          int
}

// PositioningLayoutMetrics contains metrics about the layout
type PositioningLayoutMetrics struct {
	TotalTasks        int
	PositionedTasks   int
	CollisionCount    int
	OverlapCount      int
	SpaceEfficiency   float64
	AlignmentScore    float64
	SpacingScore      float64
	VisualBalance     float64
	GridUtilization   float64
	AverageSpacing    float64
	MaxSpacing        float64
	MinSpacing        float64
	AlignmentErrors   int
	SpacingErrors     int
}

// NewPositioningEngine creates a new positioning engine
func NewPositioningEngine(gridConfig *GridConfig) *PositioningEngine {
	// Create default grid constraints (will be set in getDefaultGridConstraints)
	
	engine := &PositioningEngine{
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

// PositionTasks positions all tasks within the calendar grid
func (pe *PositioningEngine) PositionTasks(tasks []*data.Task, existingBars []*IntegratedTaskBar) (*PositioningResult, error) {
	// Create positioning context
	context := &PositioningContext{
		CalendarStart:   pe.gridConfig.CalendarStart,
		CalendarEnd:     pe.gridConfig.CalendarEnd,
		DayWidth:        pe.gridConfig.DayWidth,
		DayHeight:       pe.gridConfig.DayHeight,
		AvailableHeight: pe.gridConfig.DayHeight * float64(pe.gridConfig.MaxRowsPerDay),
		AvailableWidth:  pe.gridConfig.DayWidth * 7.0,
		ExistingBars:    existingBars,
		GridConstraints: pe.getDefaultGridConstraints(),
		VisualSettings:  pe.getDefaultVisualSettings(),
		CurrentTime:     time.Now(),
		TaskDensity:     pe.calculateTaskDensity(tasks),
		OverlapCount:    pe.countOverlaps(existingBars),
		ConflictCount:   pe.countConflicts(existingBars),
	}
	
	// Create integrated task bars
	integratedBars := pe.createIntegratedTaskBars(tasks, context)
	
	// Apply positioning rules
	positionedBars := pe.applyPositioningRules(integratedBars, context)
	
	// Apply spacing rules
	spacedBars := pe.applySpacingRules(positionedBars, context)
	
	// Snap to grid if enabled
	if context.GridConstraints.SnapToGrid {
		spacedBars = pe.snapToGrid(spacedBars, context)
	}
	
	// Resolve collisions
	finalBars := pe.resolveCollisions(spacedBars, context)
	
	// Calculate metrics
	metrics := pe.calculateLayoutMetrics(finalBars, context)
	
	// Generate recommendations
	recommendations := pe.generatePositioningRecommendations(metrics, context)
	
	return &PositioningResult{
		TaskBars:        finalBars,
		Metrics:         metrics,
		Recommendations: recommendations,
		AnalysisDate:    time.Now(),
	}, nil
}

// createIntegratedTaskBars creates integrated task bars from tasks
func (pe *PositioningEngine) createIntegratedTaskBars(tasks []*data.Task, context *PositioningContext) []*IntegratedTaskBar {
	var bars []*IntegratedTaskBar
	
	for _, task := range tasks {
		// Calculate basic positioning
		startX := pe.calculateXPosition(task.StartDate, context)
		endX := pe.calculateXPosition(task.EndDate, context)
		width := endX - startX
		
		// Calculate Y position (will be refined by positioning rules)
		y := pe.calculateInitialYPosition(task, context)
		
		// Calculate height
		height := pe.calculateTaskHeight(task, context)
		
		// Get task category and color
		category := data.GetCategory(task.Category)
		
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
			IsContinuation:  pe.isTaskContinuation(task, context),
			IsStart:         pe.isTaskStart(task, context),
			IsEnd:           pe.isTaskEnd(task, context),
			MonthBoundary:   pe.hasMonthBoundary(task, context),
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
func (pe *PositioningEngine) calculateXPosition(date time.Time, context *PositioningContext) float64 {
	daysFromStart := int(date.Sub(context.CalendarStart).Hours() / 24)
	return float64(daysFromStart) * context.DayWidth
}

// calculateInitialYPosition calculates the initial Y position for a task
func (pe *PositioningEngine) calculateInitialYPosition(task *data.Task, context *PositioningContext) float64 {
	// Start with a basic Y position based on task priority
	baseY := context.DayHeight * 0.1 // 10% from top
	
	// Adjust based on task priority
	priorityOffset := float64(task.Priority) * context.DayHeight * 0.1
	
	return baseY + priorityOffset
}

// calculateTaskHeight calculates the height of a task
func (pe *PositioningEngine) calculateTaskHeight(task *data.Task, context *PositioningContext) float64 {
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

// applyPositioningRules applies alignment rules to task bars
func (pe *PositioningEngine) applyPositioningRules(bars []*IntegratedTaskBar, context *PositioningContext) []*IntegratedTaskBar {
	// Sort rules by priority
	sort.Slice(pe.alignmentRules, func(i, j int) bool {
		return pe.alignmentRules[i].Priority > pe.alignmentRules[j].Priority
	})
	
	// Apply rules to each bar
	for _, bar := range bars {
		for _, rule := range pe.alignmentRules {
			if rule.Condition(bar, context) {
				action := rule.Action(bar, context)
				pe.applyPositioningAction(bar, action)
				break // Apply only the first matching rule
			}
		}
	}
	
	return bars
}

// applySpacingRules applies spacing rules between task bars
func (pe *PositioningEngine) applySpacingRules(bars []*IntegratedTaskBar, context *PositioningContext) []*IntegratedTaskBar {
	// Sort rules by priority
	sort.Slice(pe.spacingRules, func(i, j int) bool {
		return pe.spacingRules[i].Priority > pe.spacingRules[j].Priority
	})
	
	// Apply spacing rules between adjacent bars
	for i := 0; i < len(bars)-1; i++ {
		for j := i + 1; j < len(bars); j++ {
			for _, rule := range pe.spacingRules {
				if rule.Condition(bars[i], bars[j], context) {
					action := rule.Action(bars[i], bars[j], context)
					pe.applySpacingAction(bars[i], bars[j], action)
					break // Apply only the first matching rule
				}
			}
		}
	}
	
	return bars
}

// applyPositioningAction applies a positioning action to a task bar
func (pe *PositioningEngine) applyPositioningAction(bar *IntegratedTaskBar, action *PositioningAction) {
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
func (pe *PositioningEngine) applySpacingAction(bar1, bar2 *IntegratedTaskBar, action *SpacingAction) {
	// Calculate distance between bars
	distance := pe.calculateDistance(bar1, bar2)
	
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
func (pe *PositioningEngine) snapToGrid(bars []*IntegratedTaskBar, context *PositioningContext) []*IntegratedTaskBar {
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
func (pe *PositioningEngine) resolveCollisions(bars []*IntegratedTaskBar, context *PositioningContext) []*IntegratedTaskBar {
	// Sort bars by priority (higher priority first)
	sort.Slice(bars, func(i, j int) bool {
		return bars[i].Priority > bars[j].Priority
	})
	
	// Resolve collisions
	for i := 0; i < len(bars); i++ {
		for j := i + 1; j < len(bars); j++ {
			if pe.barsCollide(bars[i], bars[j], context) {
				pe.resolveCollision(bars[i], bars[j], context)
			}
		}
	}
	
	return bars
}

// barsCollide checks if two task bars collide
func (pe *PositioningEngine) barsCollide(bar1, bar2 *IntegratedTaskBar, context *PositioningContext) bool {
	buffer := context.GridConstraints.CollisionBuffer
	
	// Check horizontal overlap
	horizontalOverlap := bar1.StartX < bar2.EndX+buffer && bar2.StartX < bar1.EndX+buffer
	
	// Check vertical overlap
	verticalOverlap := bar1.Y < bar2.Y+bar2.Height+buffer && bar2.Y < bar1.Y+bar1.Height+buffer
	
	return horizontalOverlap && verticalOverlap
}

// resolveCollision resolves a collision between two task bars
func (pe *PositioningEngine) resolveCollision(bar1, bar2 *IntegratedTaskBar, context *PositioningContext) {
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
func (pe *PositioningEngine) calculateDistance(bar1, bar2 *IntegratedTaskBar) float64 {
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
func (pe *PositioningEngine) calculateTaskDensity(tasks []*data.Task) float64 {
	if len(tasks) == 0 {
		return 0.0
	}
	
	// Calculate total calendar area
	totalArea := pe.gridConfig.DayWidth * pe.gridConfig.DayHeight * 7.0 * 4.0 // 7 days, 4 weeks
	
	// Calculate average task area
	var totalTaskArea float64
	for _, task := range tasks {
		duration := task.EndDate.Sub(task.StartDate).Hours() / 24
		taskArea := duration * pe.gridConfig.DayWidth * pe.gridConfig.DayHeight * 0.6
		totalTaskArea += taskArea
	}
	
	avgTaskArea := totalTaskArea / float64(len(tasks))
	
	// Calculate density
	return (avgTaskArea * float64(len(tasks))) / totalArea
}

// countOverlaps counts the number of overlapping task bars
func (pe *PositioningEngine) countOverlaps(bars []*IntegratedTaskBar) int {
	count := 0
	for i := 0; i < len(bars); i++ {
		for j := i + 1; j < len(bars); j++ {
			if pe.barsOverlap(bars[i], bars[j]) {
				count++
			}
		}
	}
	return count
}

// countConflicts counts the number of conflicts between task bars
func (pe *PositioningEngine) countConflicts(bars []*IntegratedTaskBar) int {
	// For now, use overlap count as conflict count
	return pe.countOverlaps(bars)
}

// barsOverlap checks if two task bars overlap
func (pe *PositioningEngine) barsOverlap(bar1, bar2 *IntegratedTaskBar) bool {
	// Check horizontal overlap
	horizontalOverlap := bar1.StartX < bar2.EndX && bar2.StartX < bar1.EndX
	
	// Check vertical overlap
	verticalOverlap := bar1.Y < bar2.Y+bar2.Height && bar2.Y < bar1.Y+bar1.Height
	
	return horizontalOverlap && verticalOverlap
}

// calculateLayoutMetrics calculates layout metrics
func (pe *PositioningEngine) calculateLayoutMetrics(bars []*IntegratedTaskBar, context *PositioningContext) *PositioningLayoutMetrics {
	metrics := &PositioningLayoutMetrics{
		TotalTasks:      len(bars),
		PositionedTasks: len(bars),
		CollisionCount:  pe.countOverlaps(bars),
		OverlapCount:    pe.countOverlaps(bars),
	}
	
	// Calculate space efficiency
	usedSpace := pe.calculateUsedSpace(bars)
	totalSpace := context.AvailableWidth * context.AvailableHeight
	metrics.SpaceEfficiency = usedSpace / totalSpace
	
	// Calculate alignment score
	metrics.AlignmentScore = pe.calculateAlignmentScore(bars, context)
	
	// Calculate spacing score
	metrics.SpacingScore = pe.calculateSpacingScore(bars, context)
	
	// Calculate visual balance
	metrics.VisualBalance = pe.calculateVisualBalance(bars, context)
	
	// Calculate grid utilization
	metrics.GridUtilization = pe.calculateGridUtilization(bars, context)
	
	// Calculate average spacing
	metrics.AverageSpacing = pe.calculateAverageSpacing(bars, context)
	
	return metrics
}

// calculateUsedSpace calculates the total space used by task bars
func (pe *PositioningEngine) calculateUsedSpace(bars []*IntegratedTaskBar) float64 {
	var totalSpace float64
	for _, bar := range bars {
		totalSpace += bar.Width * bar.Height
	}
	return totalSpace
}

// calculateAlignmentScore calculates the alignment score
func (pe *PositioningEngine) calculateAlignmentScore(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
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
func (pe *PositioningEngine) calculateSpacingScore(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
	if len(bars) < 2 {
		return 1.0
	}
	
	var spacingScore float64
	pairCount := 0
	
	for i := 0; i < len(bars)-1; i++ {
		for j := i + 1; j < len(bars); j++ {
			distance := pe.calculateDistance(bars[i], bars[j])
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
func (pe *PositioningEngine) calculateVisualBalance(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
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
func (pe *PositioningEngine) calculateGridUtilization(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
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
func (pe *PositioningEngine) calculateAverageSpacing(bars []*IntegratedTaskBar, context *PositioningContext) float64 {
	if len(bars) < 2 {
		return 0.0
	}
	
	var totalSpacing float64
	pairCount := 0
	
	for i := 0; i < len(bars)-1; i++ {
		for j := i + 1; j < len(bars); j++ {
			distance := pe.calculateDistance(bars[i], bars[j])
			totalSpacing += distance
			pairCount++
		}
	}
	
	return totalSpacing / float64(pairCount)
}

// generatePositioningRecommendations generates recommendations based on layout metrics
func (pe *PositioningEngine) generatePositioningRecommendations(metrics *PositioningLayoutMetrics, context *PositioningContext) []string {
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

// Helper methods for task positioning
func (pe *PositioningEngine) isTaskContinuation(task *data.Task, context *PositioningContext) bool {
	return task.StartDate.Before(context.CalendarStart)
}

func (pe *PositioningEngine) isTaskStart(task *data.Task, context *PositioningContext) bool {
	return task.StartDate.Equal(context.CalendarStart) || task.StartDate.After(context.CalendarStart)
}

func (pe *PositioningEngine) isTaskEnd(task *data.Task, context *PositioningContext) bool {
	return task.EndDate.Equal(context.CalendarEnd) || task.EndDate.Before(context.CalendarEnd)
}

func (pe *PositioningEngine) hasMonthBoundary(task *data.Task, context *PositioningContext) bool {
	startMonth := task.StartDate.Month()
	endMonth := task.EndDate.Month()
	return startMonth != endMonth
}

// Default configuration methods
func (pe *PositioningEngine) getDefaultGridConstraints() *GridConstraints {
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

func (pe *PositioningEngine) getDefaultVisualSettings() *IntegratedVisualSettings {
	return &IntegratedVisualSettings{
		ShowTaskNames:         true,
		ShowTaskDurations:     true,
		ShowTaskPriorities:    true,
		ShowConflictIndicators: true,
		CollapseThreshold:     5,
		AnimationEnabled:      false,
		HighlightConflicts:    true,
		ColorScheme:           "default",
		FontSize:              "small",
		TaskBarOpacity:        0.9,
		BorderWidth:           0.5,
	}
}

// PositioningResult contains the result of positioning operations
type PositioningResult struct {
	TaskBars        []*IntegratedTaskBar
	Metrics         *PositioningLayoutMetrics
	Recommendations []string
	AnalysisDate    time.Time
}

// Add default alignment and spacing rules
func (pe *PositioningEngine) addDefaultAlignmentRules() {
	// High priority tasks alignment rule
	pe.alignmentRules = append(pe.alignmentRules, AlignmentRule{
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
	pe.alignmentRules = append(pe.alignmentRules, AlignmentRule{
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
	pe.alignmentRules = append(pe.alignmentRules, AlignmentRule{
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

func (pe *PositioningEngine) addDefaultSpacingRules() {
	// High priority spacing rule
	pe.spacingRules = append(pe.spacingRules, SpacingRule{
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
	pe.spacingRules = append(pe.spacingRules, SpacingRule{
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
