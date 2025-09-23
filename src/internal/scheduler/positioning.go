package scheduler

import (
	"fmt"
	"math"
	"sort"
	"time"

	"phd-dissertation-planner/internal/common"
)

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
		HighlightConflicts:     true,
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
