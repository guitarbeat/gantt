package calendar

import (
	"fmt"
	"sort"
	"time"

	"phd-dissertation-planner/internal/shared"
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
	OverlapNone        OverlapType = "NONE"
	OverlapPartial     OverlapType = "PARTIAL"
	OverlapComplete    OverlapType = "COMPLETE"
	OverlapNested      OverlapType = "NESTED"
	OverlapAdjacent    OverlapType = "ADJACENT"
	OverlapIdentical   OverlapType = "IDENTICAL"
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
	Tasks        []*shared.Task
	Overlaps     []*TaskOverlap
	GroupID      string
	StartDate    time.Time
	EndDate      time.Time
	MaxSeverity  OverlapSeverity
	ConflictCount int
	Resolution   string
}

// OverlapAnalysis contains comprehensive overlap analysis results
type OverlapAnalysis struct {
	TotalTasks      int
	OverlappingTasks int
	OverlapGroups   []*OverlapGroup
	TotalOverlaps   int
	CriticalOverlaps int
	HighOverlaps    int
	MediumOverlaps  int
	LowOverlaps     int
	AnalysisDate    time.Time
	Summary         string
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

// SetPrecision sets the minimum overlap duration to consider
func (se *SpatialEngine) SetPrecision(precision time.Duration) {
	se.precision = precision
}

// DetectOverlaps detects all overlaps in a collection of tasks
func (se *SpatialEngine) DetectOverlaps(tasks []*shared.Task) *OverlapAnalysis {
	analysis := &OverlapAnalysis{
		TotalTasks:      len(tasks),
		OverlappingTasks: 0,
		OverlapGroups:   make([]*OverlapGroup, 0),
		AnalysisDate:    time.Now(),
	}

	// Create overlap groups using the existing grouping algorithm
	groups := se.groupOverlappingTasks(tasks)
	
	// Analyze each group for overlaps
	for i, group := range groups {
		overlapGroup := &OverlapGroup{
			Tasks:        group.Tasks,
			Overlaps:     make([]*TaskOverlap, 0),
			GroupID:      fmt.Sprintf("group_%d", i),
			StartDate:    group.StartDate,
			EndDate:      group.EndDate,
			MaxSeverity:  SeverityNone,
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
func (se *SpatialEngine) groupOverlappingTasks(tasks []*shared.Task) []*TaskGroup {
	// Sort tasks by start date
	sortedTasks := make([]*shared.Task, len(tasks))
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
func (se *SpatialEngine) tasksOverlap(group *TaskGroup, task *shared.Task) bool {
	for _, groupTask := range group.Tasks {
		if se.tasksOverlapDirect(groupTask, task) {
			return true
		}
	}
	return false
}

// tasksOverlapDirect checks if two tasks overlap directly
func (se *SpatialEngine) tasksOverlapDirect(task1, task2 *shared.Task) bool {
	// Tasks overlap if one starts before the other ends
	return !task1.StartDate.After(task2.EndDate) && !task2.StartDate.After(task1.EndDate)
}

// detectGroupOverlaps detects all overlaps within a group of tasks
func (se *SpatialEngine) detectGroupOverlaps(tasks []*shared.Task) []*TaskOverlap {
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
func (se *SpatialEngine) analyzeTaskOverlap(task1, task2 *shared.Task) *TaskOverlap {
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
func (se *SpatialEngine) determineOverlapType(task1, task2 *shared.Task, overlapStart, overlapEnd time.Time) OverlapType {
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
func (se *SpatialEngine) calculateOverlapSeverity(task1, task2 *shared.Task, overlapType OverlapType, duration time.Duration) OverlapSeverity {
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
func (se *SpatialEngine) calculateOverlapPercentage(task1, task2 *shared.Task, overlapDuration time.Duration) float64 {
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
func (se *SpatialEngine) calculateOverlapPriority(task1, task2 *shared.Task) int {
	// Higher priority task wins
	if task1.Priority > task2.Priority {
		return task1.Priority
	}
	return task2.Priority
}

// generateConflictInfo generates conflict reason and resolution hint
func (se *SpatialEngine) generateConflictInfo(task1, task2 *shared.Task, overlapType OverlapType, severity OverlapSeverity) (string, string) {
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
