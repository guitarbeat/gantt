package calendar

import (
	"fmt"
	"sort"
	"time"

	"phd-dissertation-planner/internal/data"
)

// ConflictCategorizer handles categorization and analysis of task conflicts
type ConflictCategorizer struct {
	overlapDetector *OverlapDetector
	rules           []ConflictRule
	severityWeights map[OverlapSeverity]int
}

// ConflictRule defines a rule for categorizing conflicts
type ConflictRule struct {
	Name        string
	Description string
	Condition   func(*TaskOverlap, *data.Task, *data.Task) bool
	Category    ConflictCategory
	Severity    OverlapSeverity
	Priority    int
}

// ConflictCategory represents a category of conflict
type ConflictCategory string

const (
	CategoryScheduleConflict    ConflictCategory = "SCHEDULE_CONFLICT"
	CategoryResourceConflict    ConflictCategory = "RESOURCE_CONFLICT"
	CategoryDependencyConflict  ConflictCategory = "DEPENDENCY_CONFLICT"
	CategoryPriorityConflict    ConflictCategory = "PRIORITY_CONFLICT"
	CategoryCategoryConflict    ConflictCategory = "CATEGORY_CONFLICT"
	CategoryAssigneeConflict    ConflictCategory = "ASSIGNEE_CONFLICT"
	CategoryTimelineConflict    ConflictCategory = "TIMELINE_CONFLICT"
	CategoryWorkloadConflict    ConflictCategory = "WORKLOAD_CONFLICT"
	CategoryDeadlineConflict    ConflictCategory = "DEADLINE_CONFLICT"
	CategoryMilestoneConflict   ConflictCategory = "MILESTONE_CONFLICT"
)

// ConflictResolution represents a resolution strategy for a conflict
type ConflictResolution struct {
	Strategy    string
	Description string
	Actions     []string
	Priority    int
	Effort      string // "LOW", "MEDIUM", "HIGH"
	Impact      string // "LOW", "MEDIUM", "HIGH"
}

// CategorizedConflict represents a conflict with detailed categorization
type CategorizedConflict struct {
	*TaskOverlap
	Category       ConflictCategory
	SubCategory    string
	RootCause      string
	Impact         string
	Resolution     ConflictResolution
	AlternativeResolutions []ConflictResolution
	RiskLevel      string
	Urgency        string
	Complexity     string
}

// ConflictAnalysis contains comprehensive conflict analysis results
type ConflictAnalysis struct {
	TotalConflicts      int
	ConflictsByCategory map[ConflictCategory]int
	ConflictsBySeverity map[OverlapSeverity]int
	ConflictsByUrgency  map[string]int
	ConflictsByRisk     map[string]int
	CategorizedConflicts []*CategorizedConflict
	ResolutionSummary   map[string]int
	RiskAssessment      string
	Recommendations     []string
	AnalysisDate        time.Time
}

// NewConflictCategorizer creates a new conflict categorizer
func NewConflictCategorizer(overlapDetector *OverlapDetector) *ConflictCategorizer {
	cc := &ConflictCategorizer{
		overlapDetector: overlapDetector,
		rules:           make([]ConflictRule, 0),
		severityWeights: map[OverlapSeverity]int{
			SeverityCritical: 5,
			SeverityHigh:     4,
			SeverityMedium:   3,
			SeverityLow:      2,
			SeverityNone:     1,
		},
	}
	
	// Initialize default conflict rules
	cc.initializeDefaultRules()
	
	return cc
}

// initializeDefaultRules sets up default conflict categorization rules
func (cc *ConflictCategorizer) initializeDefaultRules() {
	cc.rules = []ConflictRule{
		{
			Name:        "Identical Schedule Conflict",
			Description: "Tasks have identical start and end dates",
			Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
				return overlap.OverlapType == OverlapIdentical
			},
			Category: CategoryScheduleConflict,
			Severity: SeverityCritical,
			Priority: 1,
		},
		{
			Name:        "High Priority Overlap",
			Description: "Both tasks have high priority (3+)",
			Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
				return task1.Priority >= 3 && task2.Priority >= 3
			},
			Category: CategoryPriorityConflict,
			Severity: SeverityHigh,
			Priority: 2,
		},
		{
			Name:        "Same Assignee Conflict",
			Description: "Both tasks are assigned to the same person",
			Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
				return task1.Assignee != "" && task1.Assignee == task2.Assignee
			},
			Category: CategoryAssigneeConflict,
			Severity: SeverityHigh,
			Priority: 3,
		},
		{
			Name:        "Same Category Conflict",
			Description: "Both tasks belong to the same category",
			Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
				return task1.Category != "" && task1.Category == task2.Category
			},
			Category: CategoryCategoryConflict,
			Severity: SeverityMedium,
			Priority: 4,
		},
		{
			Name:        "Milestone Overlap",
			Description: "One or both tasks are milestones",
			Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
				return task1.IsMilestone || task2.IsMilestone
			},
			Category: CategoryMilestoneConflict,
			Severity: SeverityHigh,
			Priority: 5,
		},
		{
			Name:        "Deadline Conflict",
			Description: "Tasks have overlapping deadlines with high priority",
			Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
				// Check if both tasks are due soon and have high priority
				now := time.Now()
				daysFromNow := 7
				task1DueSoon := task1.EndDate.Before(now.AddDate(0, 0, daysFromNow))
				task2DueSoon := task2.EndDate.Before(now.AddDate(0, 0, daysFromNow))
				return (task1DueSoon || task2DueSoon) && (task1.Priority >= 2 || task2.Priority >= 2)
			},
			Category: CategoryDeadlineConflict,
			Severity: SeverityHigh,
			Priority: 6,
		},
		{
			Name:        "Long Duration Overlap",
			Description: "Tasks have significant overlap duration",
			Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
				overlapPercentage := float64(overlap.Duration) / float64(task1.EndDate.Sub(task1.StartDate))
				return overlapPercentage >= 0.7
			},
			Category: CategoryTimelineConflict,
			Severity: SeverityMedium,
			Priority: 7,
		},
		{
			Name:        "Dependency Chain Conflict",
			Description: "Tasks are in a dependency chain",
			Condition: func(overlap *TaskOverlap, task1, task2 *data.Task) bool {
				// Check if one task depends on the other
				for _, dep := range task1.Dependencies {
					if dep == task2.ID {
						return true
					}
				}
				for _, dep := range task2.Dependencies {
					if dep == task1.ID {
						return true
					}
				}
				return false
			},
			Category: CategoryDependencyConflict,
			Severity: SeverityCritical,
			Priority: 8,
		},
	}
}

// CategorizeConflicts categorizes all conflicts in an overlap analysis
func (cc *ConflictCategorizer) CategorizeConflicts(analysis *OverlapAnalysis) *ConflictAnalysis {
	conflictAnalysis := &ConflictAnalysis{
		TotalConflicts:      0,
		ConflictsByCategory: make(map[ConflictCategory]int),
		ConflictsBySeverity: make(map[OverlapSeverity]int),
		ConflictsByUrgency:  make(map[string]int),
		ConflictsByRisk:     make(map[string]int),
		CategorizedConflicts: make([]*CategorizedConflict, 0),
		ResolutionSummary:   make(map[string]int),
		AnalysisDate:        time.Now(),
	}

	// Process each overlap group
	for _, group := range analysis.OverlapGroups {
		for _, overlap := range group.Overlaps {
			// Get the actual task objects
			task1, task2 := cc.getTasksFromOverlap(overlap, group.Tasks)
			if task1 == nil || task2 == nil {
				continue
			}

			// Categorize the conflict
			categorizedConflict := cc.categorizeConflict(overlap, task1, task2)
			conflictAnalysis.CategorizedConflicts = append(conflictAnalysis.CategorizedConflicts, categorizedConflict)
			conflictAnalysis.TotalConflicts++

			// Update statistics
			conflictAnalysis.ConflictsByCategory[categorizedConflict.Category]++
			conflictAnalysis.ConflictsBySeverity[categorizedConflict.Severity]++
			conflictAnalysis.ConflictsByUrgency[categorizedConflict.Urgency]++
			conflictAnalysis.ConflictsByRisk[categorizedConflict.RiskLevel]++
			conflictAnalysis.ResolutionSummary[categorizedConflict.Resolution.Strategy]++
		}
	}

	// Generate analysis summary
	conflictAnalysis.RiskAssessment = cc.generateRiskAssessment(conflictAnalysis)
	conflictAnalysis.Recommendations = cc.generateRecommendations(conflictAnalysis)

	return conflictAnalysis
}

// categorizeConflict categorizes a single conflict
func (cc *ConflictCategorizer) categorizeConflict(overlap *TaskOverlap, task1, task2 *data.Task) *CategorizedConflict {
	// Find the best matching rule
	var bestRule *ConflictRule
	highestPriority := 0

	for _, rule := range cc.rules {
		if rule.Condition(overlap, task1, task2) && rule.Priority > highestPriority {
			bestRule = &rule
			highestPriority = rule.Priority
		}
	}

	// Default categorization if no rule matches
	if bestRule == nil {
		bestRule = &ConflictRule{
			Name:        "Generic Schedule Conflict",
			Description: "Generic schedule overlap",
			Category:    CategoryScheduleConflict,
			Severity:    overlap.Severity,
			Priority:    0,
		}
	}

	// Create categorized conflict
	categorizedConflict := &CategorizedConflict{
		TaskOverlap: overlap,
		Category:    bestRule.Category,
		SubCategory: cc.determineSubCategory(overlap, task1, task2),
		RootCause:   cc.determineRootCause(overlap, task1, task2, bestRule),
		Impact:      cc.assessImpact(overlap, task1, task2),
		Resolution:  cc.generateResolution(overlap, task1, task2, bestRule),
		RiskLevel:   cc.assessRiskLevel(overlap, task1, task2),
		Urgency:     cc.assessUrgency(overlap, task1, task2),
		Complexity:  cc.assessComplexity(overlap, task1, task2),
	}

	// Generate alternative resolutions
	categorizedConflict.AlternativeResolutions = cc.generateAlternativeResolutions(overlap, task1, task2, bestRule)

	return categorizedConflict
}

// determineSubCategory determines a more specific subcategory
func (cc *ConflictCategorizer) determineSubCategory(overlap *TaskOverlap, task1, task2 *data.Task) string {
	switch overlap.OverlapType {
	case OverlapIdentical:
		return "Identical Schedules"
	case OverlapNested:
		return "Nested Tasks"
	case OverlapComplete:
		return "Complete Overlap"
	case OverlapPartial:
		percentage := float64(overlap.Duration) / float64(task1.EndDate.Sub(task1.StartDate))
		if percentage >= 0.8 {
			return "High Overlap"
		} else if percentage >= 0.5 {
			return "Medium Overlap"
		} else {
			return "Low Overlap"
		}
	case OverlapAdjacent:
		return "Adjacent Tasks"
	default:
		return "Unknown Overlap"
	}
}

// determineRootCause determines the root cause of the conflict
func (cc *ConflictCategorizer) determineRootCause(overlap *TaskOverlap, task1, task2 *data.Task, rule *ConflictRule) string {
	switch rule.Category {
	case CategoryScheduleConflict:
		return "Tasks scheduled for identical time periods"
	case CategoryResourceConflict:
		return "Tasks competing for the same resources"
	case CategoryDependencyConflict:
		return "Tasks have conflicting dependency relationships"
	case CategoryPriorityConflict:
		return "Both tasks have high priority causing resource contention"
	case CategoryCategoryConflict:
		return "Tasks in the same category may have conflicting requirements"
	case CategoryAssigneeConflict:
		return "Same person assigned to overlapping tasks"
	case CategoryTimelineConflict:
		return "Tasks have significant timeline overlap"
	case CategoryWorkloadConflict:
		return "Tasks create excessive workload for assignee"
	case CategoryDeadlineConflict:
		return "Tasks have conflicting deadlines"
	case CategoryMilestoneConflict:
		return "Milestone tasks have overlapping schedules"
	default:
		return "Unknown conflict cause"
	}
}

// assessImpact assesses the impact of the conflict
func (cc *ConflictCategorizer) assessImpact(overlap *TaskOverlap, task1, task2 *data.Task) string {
	// Calculate impact based on multiple factors
	impactScore := 0

	// Severity impact
	impactScore += cc.severityWeights[overlap.Severity]

	// Priority impact
	impactScore += task1.Priority + task2.Priority

	// Duration impact
	if overlap.Duration > time.Hour*24*7 { // More than a week
		impactScore += 2
	} else if overlap.Duration > time.Hour*24*3 { // More than 3 days
		impactScore += 1
	}

	// Assignee impact
	if task1.Assignee == task2.Assignee && task1.Assignee != "" {
		impactScore += 3
	}

	// Category impact
	if task1.Category == task2.Category && task1.Category != "" {
		impactScore += 1
	}

	// Milestone impact
	if task1.IsMilestone || task2.IsMilestone {
		impactScore += 2
	}

	// Determine impact level
	if impactScore >= 10 {
		return "HIGH"
	} else if impactScore >= 6 {
		return "MEDIUM"
	} else {
		return "LOW"
	}
}

// generateResolution generates a resolution strategy for the conflict
func (cc *ConflictCategorizer) generateResolution(overlap *TaskOverlap, task1, task2 *data.Task, rule *ConflictRule) ConflictResolution {
	switch rule.Category {
	case CategoryScheduleConflict:
		return ConflictResolution{
			Strategy:    "Reschedule Tasks",
			Description: "Adjust start/end times to eliminate overlap",
			Actions: []string{
				"Move one task to a different time slot",
				"Extend timeline to accommodate both tasks",
				"Consider making one task a subtask",
			},
			Priority: 1,
			Effort:   "MEDIUM",
			Impact:   "HIGH",
		}
	case CategoryAssigneeConflict:
		return ConflictResolution{
			Strategy:    "Reassign Tasks",
			Description: "Assign tasks to different people",
			Actions: []string{
				"Find alternative assignee for one task",
				"Split tasks among multiple people",
				"Adjust workload distribution",
			},
			Priority: 1,
			Effort:   "HIGH",
			Impact:   "HIGH",
		}
	case CategoryPriorityConflict:
		return ConflictResolution{
			Strategy:    "Priority Adjustment",
			Description: "Adjust task priorities to resolve conflict",
			Actions: []string{
				"Lower priority of one task",
				"Escalate to management for priority decision",
				"Consider task dependencies",
			},
			Priority: 2,
			Effort:   "LOW",
			Impact:   "MEDIUM",
		}
	case CategoryDependencyConflict:
		return ConflictResolution{
			Strategy:    "Dependency Review",
			Description: "Review and adjust task dependencies",
			Actions: []string{
				"Remove circular dependencies",
				"Adjust task sequence",
				"Break down complex dependencies",
			},
			Priority: 1,
			Effort:   "HIGH",
			Impact:   "HIGH",
		}
	default:
		return ConflictResolution{
			Strategy:    "Generic Resolution",
			Description: "Apply standard conflict resolution approach",
			Actions: []string{
				"Review task requirements",
				"Adjust schedules",
				"Escalate if necessary",
			},
			Priority: 3,
			Effort:   "MEDIUM",
			Impact:   "MEDIUM",
		}
	}
}

// generateAlternativeResolutions generates alternative resolution strategies
func (cc *ConflictCategorizer) generateAlternativeResolutions(overlap *TaskOverlap, task1, task2 *data.Task, rule *ConflictRule) []ConflictResolution {
	var alternatives []ConflictResolution

	// Add generic alternatives
	alternatives = append(alternatives, ConflictResolution{
		Strategy:    "Task Merging",
		Description: "Combine overlapping tasks into a single task",
		Actions:     []string{"Merge task requirements", "Combine timelines", "Update dependencies"},
		Priority:    2,
		Effort:      "HIGH",
		Impact:      "HIGH",
	})

	alternatives = append(alternatives, ConflictResolution{
		Strategy:    "Timeline Extension",
		Description: "Extend project timeline to accommodate both tasks",
		Actions:     []string{"Adjust project deadlines", "Update dependencies", "Communicate changes"},
		Priority:    3,
		Effort:      "MEDIUM",
		Impact:      "MEDIUM",
	})

	alternatives = append(alternatives, ConflictResolution{
		Strategy:    "Resource Addition",
		Description: "Add additional resources to handle both tasks",
		Actions:     []string{"Hire additional staff", "Reallocate resources", "Outsource tasks"},
		Priority:    4,
		Effort:      "HIGH",
		Impact:      "LOW",
	})

	return alternatives
}

// assessRiskLevel assesses the risk level of the conflict
func (cc *ConflictCategorizer) assessRiskLevel(overlap *TaskOverlap, task1, task2 *data.Task) string {
	riskScore := 0

	// Severity risk
	riskScore += cc.severityWeights[overlap.Severity]

	// Deadline risk
	now := time.Now()
	if task1.EndDate.Before(now.AddDate(0, 0, 7)) || task2.EndDate.Before(now.AddDate(0, 0, 7)) {
		riskScore += 3
	}

	// Priority risk
	if task1.Priority >= 4 || task2.Priority >= 4 {
		riskScore += 2
	}

	// Milestone risk
	if task1.IsMilestone || task2.IsMilestone {
		riskScore += 2
	}

	// Dependency risk
	if len(task1.Dependencies) > 0 || len(task2.Dependencies) > 0 {
		riskScore += 1
	}

	if riskScore >= 8 {
		return "HIGH"
	} else if riskScore >= 5 {
		return "MEDIUM"
	} else {
		return "LOW"
	}
}

// assessUrgency assesses the urgency of resolving the conflict
func (cc *ConflictCategorizer) assessUrgency(overlap *TaskOverlap, task1, task2 *data.Task) string {
	urgencyScore := 0

	// Severity urgency
	urgencyScore += cc.severityWeights[overlap.Severity]

	// Time urgency
	now := time.Now()
	daysUntilStart := int(task1.StartDate.Sub(now).Hours() / 24)
	if daysUntilStart <= 3 {
		urgencyScore += 3
	} else if daysUntilStart <= 7 {
		urgencyScore += 2
	} else if daysUntilStart <= 14 {
		urgencyScore += 1
	}

	// Priority urgency
	urgencyScore += task1.Priority + task2.Priority

	if urgencyScore >= 8 {
		return "URGENT"
	} else if urgencyScore >= 5 {
		return "HIGH"
	} else if urgencyScore >= 3 {
		return "MEDIUM"
	} else {
		return "LOW"
	}
}

// assessComplexity assesses the complexity of resolving the conflict
func (cc *ConflictCategorizer) assessComplexity(overlap *TaskOverlap, task1, task2 *data.Task) string {
	complexityScore := 0

	// Overlap type complexity
	switch overlap.OverlapType {
	case OverlapIdentical:
		complexityScore += 3
	case OverlapNested:
		complexityScore += 2
	case OverlapComplete:
		complexityScore += 2
	case OverlapPartial:
		complexityScore += 1
	default:
		complexityScore += 1
	}

	// Dependency complexity
	complexityScore += len(task1.Dependencies) + len(task2.Dependencies)

	// Assignee complexity
	if task1.Assignee == task2.Assignee && task1.Assignee != "" {
		complexityScore += 2
	}

	// Category complexity
	if task1.Category == task2.Category && task1.Category != "" {
		complexityScore += 1
	}

	if complexityScore >= 6 {
		return "HIGH"
	} else if complexityScore >= 3 {
		return "MEDIUM"
	} else {
		return "LOW"
	}
}

// getTasksFromOverlap retrieves task objects from overlap information
func (cc *ConflictCategorizer) getTasksFromOverlap(overlap *TaskOverlap, tasks []*data.Task) (*data.Task, *data.Task) {
	var task1, task2 *data.Task

	for _, task := range tasks {
		if task.ID == overlap.Task1ID {
			task1 = task
		}
		if task.ID == overlap.Task2ID {
			task2 = task
		}
	}

	return task1, task2
}

// generateRiskAssessment generates an overall risk assessment
func (cc *ConflictCategorizer) generateRiskAssessment(analysis *ConflictAnalysis) string {
	highRisk := analysis.ConflictsByRisk["HIGH"]
	mediumRisk := analysis.ConflictsByRisk["MEDIUM"]
	lowRisk := analysis.ConflictsByRisk["LOW"]

	if highRisk > 0 {
		return fmt.Sprintf("HIGH RISK: %d high-risk conflicts require immediate attention", highRisk)
	} else if mediumRisk > 0 {
		return fmt.Sprintf("MEDIUM RISK: %d medium-risk conflicts need resolution", mediumRisk)
	} else if lowRisk > 0 {
		return fmt.Sprintf("LOW RISK: %d low-risk conflicts can be addressed during planning", lowRisk)
	} else {
		return "NO RISK: No conflicts detected"
	}
}

// generateRecommendations generates actionable recommendations
func (cc *ConflictCategorizer) generateRecommendations(analysis *ConflictAnalysis) []string {
	var recommendations []string

	// Critical conflicts
	if analysis.ConflictsBySeverity[SeverityCritical] > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🚨 Address %d critical conflicts immediately", analysis.ConflictsBySeverity[SeverityCritical]))
	}

	// High priority conflicts
	if analysis.ConflictsBySeverity[SeverityHigh] > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🔴 Resolve %d high-priority conflicts", analysis.ConflictsBySeverity[SeverityHigh]))
	}

	// Assignee conflicts
	if analysis.ConflictsByCategory[CategoryAssigneeConflict] > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("👥 Review %d assignee conflicts for workload distribution", analysis.ConflictsByCategory[CategoryAssigneeConflict]))
	}

	// Dependency conflicts
	if analysis.ConflictsByCategory[CategoryDependencyConflict] > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🔗 Fix %d dependency conflicts in task relationships", analysis.ConflictsByCategory[CategoryDependencyConflict]))
	}

	// Schedule conflicts
	if analysis.ConflictsByCategory[CategoryScheduleConflict] > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("📅 Reschedule %d conflicting tasks", analysis.ConflictsByCategory[CategoryScheduleConflict]))
	}

	// General recommendations
	if analysis.TotalConflicts > 0 {
		recommendations = append(recommendations, 
			"📊 Review conflict patterns to identify systemic issues")
		recommendations = append(recommendations, 
			"🔄 Implement conflict prevention measures")
	}

	return recommendations
}

// AddCustomRule adds a custom conflict categorization rule
func (cc *ConflictCategorizer) AddCustomRule(rule ConflictRule) {
	cc.rules = append(cc.rules, rule)
	// Sort rules by priority (higher priority first)
	sort.Slice(cc.rules, func(i, j int) bool {
		return cc.rules[i].Priority > cc.rules[j].Priority
	})
}

// GetConflictsByCategory returns conflicts filtered by category
func (analysis *ConflictAnalysis) GetConflictsByCategory(category ConflictCategory) []*CategorizedConflict {
	var filtered []*CategorizedConflict
	for _, conflict := range analysis.CategorizedConflicts {
		if conflict.Category == category {
			filtered = append(filtered, conflict)
		}
	}
	return filtered
}

// GetConflictsBySeverity returns conflicts filtered by severity
func (analysis *ConflictAnalysis) GetConflictsBySeverity(severity OverlapSeverity) []*CategorizedConflict {
	var filtered []*CategorizedConflict
	for _, conflict := range analysis.CategorizedConflicts {
		if conflict.Severity == severity {
			filtered = append(filtered, conflict)
		}
	}
	return filtered
}

// GetConflictsByUrgency returns conflicts filtered by urgency
func (analysis *ConflictAnalysis) GetConflictsByUrgency(urgency string) []*CategorizedConflict {
	var filtered []*CategorizedConflict
	for _, conflict := range analysis.CategorizedConflicts {
		if conflict.Urgency == urgency {
			filtered = append(filtered, conflict)
		}
	}
	return filtered
}

// GetConflictsByRisk returns conflicts filtered by risk level
func (analysis *ConflictAnalysis) GetConflictsByRisk(riskLevel string) []*CategorizedConflict {
	var filtered []*CategorizedConflict
	for _, conflict := range analysis.CategorizedConflicts {
		if conflict.RiskLevel == riskLevel {
			filtered = append(filtered, conflict)
		}
	}
	return filtered
}

// GetSummary returns a summary of the conflict analysis
func (analysis *ConflictAnalysis) GetSummary() string {
	return fmt.Sprintf("Conflict Analysis Summary:\n"+
		"  Total Conflicts: %d\n"+
		"  By Severity: Critical=%d, High=%d, Medium=%d, Low=%d\n"+
		"  By Risk: High=%d, Medium=%d, Low=%d\n"+
		"  By Urgency: Urgent=%d, High=%d, Medium=%d, Low=%d\n"+
		"  Risk Assessment: %s",
		analysis.TotalConflicts,
		analysis.ConflictsBySeverity[SeverityCritical],
		analysis.ConflictsBySeverity[SeverityHigh],
		analysis.ConflictsBySeverity[SeverityMedium],
		analysis.ConflictsBySeverity[SeverityLow],
		analysis.ConflictsByRisk["HIGH"],
		analysis.ConflictsByRisk["MEDIUM"],
		analysis.ConflictsByRisk["LOW"],
		analysis.ConflictsByUrgency["URGENT"],
		analysis.ConflictsByUrgency["HIGH"],
		analysis.ConflictsByUrgency["MEDIUM"],
		analysis.ConflictsByUrgency["LOW"],
		analysis.RiskAssessment)
}
