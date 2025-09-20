package calendar

import (
	"fmt"
	"math"
	"sort"
	"time"

	"phd-dissertation-planner/internal/shared"
)

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
func (se *SpatialEngine) PositionTasks(tasks []*shared.Task, existingBars []*IntegratedTaskBar) (*PositioningResult, error) {
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
func (se *SpatialEngine) createIntegratedTaskBars(tasks []*shared.Task, context *PositioningContext) []*IntegratedTaskBar {
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
		category := shared.GetCategory(task.Category)
		
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
func (se *SpatialEngine) calculateInitialYPosition(task *shared.Task, context *PositioningContext) float64 {
	// Start with a basic Y position based on task priority
	baseY := context.DayHeight * 0.1 // 10% from top
	
	// Adjust based on task priority
	priorityOffset := float64(task.Priority) * context.DayHeight * 0.1
	
	return baseY + priorityOffset
}

// calculateTaskHeight calculates the height of a task
func (se *SpatialEngine) calculateTaskHeight(task *shared.Task, context *PositioningContext) float64 {
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
func (se *SpatialEngine) isTaskContinuation(task *shared.Task, context *PositioningContext) bool {
	return task.StartDate.Before(context.CalendarStart)
}

func (se *SpatialEngine) isTaskStart(task *shared.Task, context *PositioningContext) bool {
	return task.StartDate.Equal(context.CalendarStart) || task.StartDate.After(context.CalendarStart)
}

func (se *SpatialEngine) isTaskEnd(task *shared.Task, context *PositioningContext) bool {
	return task.EndDate.Equal(context.CalendarEnd) || task.EndDate.Before(context.CalendarEnd)
}

func (se *SpatialEngine) hasMonthBoundary(task *shared.Task, context *PositioningContext) bool {
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
func (se *SpatialEngine) calculateTaskDensity(tasks []*shared.Task) float64 {
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
