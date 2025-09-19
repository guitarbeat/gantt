package calendar

import (
	"fmt"
	"sort"
	"time"

	"phd-dissertation-planner/internal/data"
)

// OverlapDetector handles detection and analysis of overlapping tasks
type OverlapDetector struct {
	calendarStart time.Time
	calendarEnd   time.Time
	precision     time.Duration // Minimum overlap duration to consider
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
	Tasks        []*data.Task
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

// NewOverlapDetector creates a new overlap detector
func NewOverlapDetector(calendarStart, calendarEnd time.Time) *OverlapDetector {
	return &OverlapDetector{
		calendarStart: calendarStart,
		calendarEnd:   calendarEnd,
		precision:     time.Hour * 1, // 1 hour minimum overlap
	}
}

// SetPrecision sets the minimum overlap duration to consider
func (od *OverlapDetector) SetPrecision(precision time.Duration) {
	od.precision = precision
}

// DetectOverlaps detects all overlaps in a collection of tasks
func (od *OverlapDetector) DetectOverlaps(tasks []*data.Task) *OverlapAnalysis {
	analysis := &OverlapAnalysis{
		TotalTasks:      len(tasks),
		OverlappingTasks: 0,
		OverlapGroups:   make([]*OverlapGroup, 0),
		AnalysisDate:    time.Now(),
	}

	// Create overlap groups using the existing grouping algorithm
	groups := od.groupOverlappingTasks(tasks)
	
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
		overlaps := od.detectGroupOverlaps(group.Tasks)
		overlapGroup.Overlaps = overlaps
		overlapGroup.ConflictCount = len(overlaps)

		// Calculate group statistics
		od.calculateGroupStatistics(overlapGroup)
		
		// Add to analysis
		analysis.OverlapGroups = append(analysis.OverlapGroups, overlapGroup)
		analysis.TotalOverlaps += len(overlaps)
	}

	// Calculate overall statistics
	od.calculateAnalysisStatistics(analysis)
	
	// Generate summary
	analysis.Summary = od.generateAnalysisSummary(analysis)

	return analysis
}

// groupOverlappingTasks groups tasks that have any temporal overlap
func (od *OverlapDetector) groupOverlappingTasks(tasks []*data.Task) []*TaskGroup {
	// Sort tasks by start date
	sortedTasks := make([]*data.Task, len(tasks))
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
			Tasks:     []*data.Task{task},
			StartDate: task.StartDate,
			EndDate:   task.EndDate,
		}
		used[task.ID] = true

		// Find all tasks that overlap with this group
		for _, otherTask := range sortedTasks {
			if used[otherTask.ID] {
				continue
			}

			if od.tasksOverlap(group, otherTask) {
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
func (od *OverlapDetector) tasksOverlap(group *TaskGroup, task *data.Task) bool {
	for _, groupTask := range group.Tasks {
		if od.tasksOverlapDirect(groupTask, task) {
			return true
		}
	}
	return false
}

// tasksOverlapDirect checks if two tasks overlap directly
func (od *OverlapDetector) tasksOverlapDirect(task1, task2 *data.Task) bool {
	// Tasks overlap if one starts before the other ends
	return !task1.StartDate.After(task2.EndDate) && !task2.StartDate.After(task1.EndDate)
}

// detectGroupOverlaps detects all overlaps within a group of tasks
func (od *OverlapDetector) detectGroupOverlaps(tasks []*data.Task) []*TaskOverlap {
	var overlaps []*TaskOverlap

	// Check all pairs of tasks in the group
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			overlap := od.analyzeTaskOverlap(tasks[i], tasks[j])
			if overlap != nil {
				overlaps = append(overlaps, overlap)
			}
		}
	}

	// Sort overlaps by severity and priority
	sort.Slice(overlaps, func(i, j int) bool {
		if overlaps[i].Severity != overlaps[j].Severity {
			return od.severityOrder(overlaps[i].Severity) < od.severityOrder(overlaps[j].Severity)
		}
		return overlaps[i].Priority > overlaps[j].Priority
	})

	return overlaps
}

// analyzeTaskOverlap analyzes the overlap between two specific tasks
func (od *OverlapDetector) analyzeTaskOverlap(task1, task2 *data.Task) *TaskOverlap {
	// Check if tasks actually overlap
	if !od.tasksOverlapDirect(task1, task2) {
		return nil
	}

	// Calculate overlap details
	overlapStart := od.maxTime(task1.StartDate, task2.StartDate)
	overlapEnd := od.minTime(task1.EndDate, task2.EndDate)
	
	// Check if overlap meets minimum precision requirement
	overlapDuration := overlapEnd.Sub(overlapStart)
	if overlapDuration < od.precision {
		return nil
	}

	// Determine overlap type
	overlapType := od.determineOverlapType(task1, task2, overlapStart, overlapEnd)
	
	// Calculate severity
	severity := od.calculateOverlapSeverity(task1, task2, overlapType, overlapDuration)
	
	// Calculate priority (higher priority task wins)
	priority := od.calculateOverlapPriority(task1, task2)
	
	// Generate conflict reason and resolution hint
	conflictReason, resolutionHint := od.generateConflictInfo(task1, task2, overlapType, severity)

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
func (od *OverlapDetector) determineOverlapType(task1, task2 *data.Task, overlapStart, overlapEnd time.Time) OverlapType {
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
func (od *OverlapDetector) calculateOverlapSeverity(task1, task2 *data.Task, overlapType OverlapType, duration time.Duration) OverlapSeverity {
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
		overlapPercentage := od.calculateOverlapPercentage(task1, task2, duration)
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
func (od *OverlapDetector) calculateOverlapPercentage(task1, task2 *data.Task, overlapDuration time.Duration) float64 {
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
func (od *OverlapDetector) calculateOverlapPriority(task1, task2 *data.Task) int {
	// Higher priority task wins
	if task1.Priority > task2.Priority {
		return task1.Priority
	}
	return task2.Priority
}

// generateConflictInfo generates conflict reason and resolution hint
func (od *OverlapDetector) generateConflictInfo(task1, task2 *data.Task, overlapType OverlapType, severity OverlapSeverity) (string, string) {
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

// calculateGroupStatistics calculates statistics for an overlap group
func (od *OverlapDetector) calculateGroupStatistics(group *OverlapGroup) {
	if len(group.Overlaps) == 0 {
		return
	}

	// Find maximum severity
	for _, overlap := range group.Overlaps {
		if od.severityOrder(overlap.Severity) < od.severityOrder(group.MaxSeverity) {
			group.MaxSeverity = overlap.Severity
		}
	}

	// Generate resolution suggestion
	group.Resolution = od.generateGroupResolution(group)
}

// generateGroupResolution generates resolution suggestions for a group
func (od *OverlapDetector) generateGroupResolution(group *OverlapGroup) string {
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
func (od *OverlapDetector) calculateAnalysisStatistics(analysis *OverlapAnalysis) {
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
func (od *OverlapDetector) generateAnalysisSummary(analysis *OverlapAnalysis) string {
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

// severityOrder returns the order value for severity sorting
func (od *OverlapDetector) severityOrder(severity OverlapSeverity) int {
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

// Helper functions for time comparisons
func (od *OverlapDetector) maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func (od *OverlapDetector) minTime(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}
	return t2
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
