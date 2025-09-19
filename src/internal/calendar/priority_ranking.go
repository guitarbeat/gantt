package calendar

import (
	"fmt"
	"sort"
	"time"

	"phd-dissertation-planner/internal/data"
)

// PriorityRanker handles priority ranking and visual prominence decisions
type PriorityRanker struct {
	conflictCategorizer *ConflictCategorizer
	rankingRules        []PriorityRule
	visualWeights       map[VisualFactor]float64
}

// PriorityRule defines a rule for calculating task priority scores
type PriorityRule struct {
	Name        string
	Description string
	Weight      float64
	Calculator  func(*data.Task, *PriorityContext) float64
	Category    PriorityCategory
}

// PriorityCategory represents a category of priority factors
type PriorityCategory string

const (
	CategoryConflictPriority    PriorityCategory = "CONFLICT_PRIORITY"
	CategoryTaskPriority        PriorityCategory = "TASK_PRIORITY"
	CategoryTimelinePriority    PriorityCategory = "TIMELINE_PRIORITY"
	CategoryResourcePriority    PriorityCategory = "RESOURCE_PRIORITY"
	CategoryDependencyPriority  PriorityCategory = "DEPENDENCY_PRIORITY"
	CategoryMilestonePriority   PriorityCategory = "MILESTONE_PRIORITY"
	CategoryAssigneePriority    PriorityCategory = "ASSIGNEE_PRIORITY"
	CategoryCategoryPriority    PriorityCategory = "CATEGORY_PRIORITY"
	CategoryDeadlinePriority    PriorityCategory = "DEADLINE_PRIORITY"
	CategoryWorkloadPriority    PriorityCategory = "WORKLOAD_PRIORITY"
)

// VisualFactor represents factors that affect visual prominence
type VisualFactor string

const (
	FactorConflictSeverity    VisualFactor = "CONFLICT_SEVERITY"
	FactorTaskImportance      VisualFactor = "TASK_IMPORTANCE"
	FactorTimelineUrgency     VisualFactor = "TIMELINE_URGENCY"
	FactorResourceContention  VisualFactor = "RESOURCE_CONTENTION"
	FactorDependencyCriticality VisualFactor = "DEPENDENCY_CRITICALITY"
	FactorMilestoneStatus     VisualFactor = "MILESTONE_STATUS"
	FactorAssigneeWorkload    VisualFactor = "ASSIGNEE_WORKLOAD"
	FactorCategoryImportance  VisualFactor = "CATEGORY_IMPORTANCE"
	FactorDeadlineProximity   VisualFactor = "DEADLINE_PROXIMITY"
	FactorWorkloadBalance     VisualFactor = "WORKLOAD_BALANCE"
)

// PriorityContext provides context for priority calculations
type PriorityContext struct {
	Task                *data.Task
	Conflicts           []*CategorizedConflict
	OverlapAnalysis     *OverlapAnalysis
	ConflictAnalysis    *ConflictAnalysis
	CalendarStart       time.Time
	CalendarEnd         time.Time
	CurrentTime         time.Time
	AssigneeWorkloads   map[string]int
	CategoryImportance  map[string]float64
	ProjectMilestones   []*data.Task
}

// TaskPriority represents a task's calculated priority and visual prominence
type TaskPriority struct {
	Task                *data.Task
	PriorityScore       float64
	VisualProminence    VisualProminence
	RankingFactors      map[PriorityCategory]float64
	ConflictImpact      float64
	TimelineUrgency     float64
	ResourceContention  float64
	DependencyWeight    float64
	MilestoneWeight     float64
	AssigneeWeight      float64
	CategoryWeight      float64
	DeadlineWeight      float64
	WorkloadWeight      float64
	Recommendations     []string
	DisplayOrder        int
	VisualStyle         VisualStyle
}

// VisualProminence defines the visual prominence level
type VisualProminence string

const (
	ProminenceCritical  VisualProminence = "CRITICAL"
	ProminenceHigh      VisualProminence = "HIGH"
	ProminenceMedium    VisualProminence = "MEDIUM"
	ProminenceLow       VisualProminence = "LOW"
	ProminenceMinimal   VisualProminence = "MINIMAL"
)

// VisualStyle defines visual styling for task display
type VisualStyle struct {
	BorderColor      string
	FillColor        string
	BorderWidth      string
	Opacity          float64
	FontWeight       string
	FontSize         string
	ZIndex           int
	Animation        string
	Highlight        bool
	Blink            bool
	Pulse            bool
	Glow             bool
}

// PriorityRanking contains comprehensive priority ranking results
type PriorityRanking struct {
	TaskPriorities    []*TaskPriority
	RankingSummary    map[VisualProminence]int
	ConflictSummary   map[OverlapSeverity]int
	UrgencySummary    map[string]int
	Recommendations   []string
	VisualHierarchy   []*TaskPriority
	CriticalTasks     []*TaskPriority
	HighPriorityTasks []*TaskPriority
	MediumPriorityTasks []*TaskPriority
	LowPriorityTasks  []*TaskPriority
	MinimalTasks      []*TaskPriority
	AnalysisDate      time.Time
}

// NewPriorityRanker creates a new priority ranker
func NewPriorityRanker(conflictCategorizer *ConflictCategorizer) *PriorityRanker {
	pr := &PriorityRanker{
		conflictCategorizer: conflictCategorizer,
		rankingRules:        make([]PriorityRule, 0),
		visualWeights:       make(map[VisualFactor]float64),
	}
	
	// Initialize default ranking rules
	pr.initializeDefaultRules()
	
	// Initialize visual weights
	pr.initializeVisualWeights()
	
	return pr
}

// initializeDefaultRules sets up default priority ranking rules
func (pr *PriorityRanker) initializeDefaultRules() {
	pr.rankingRules = []PriorityRule{
		{
			Name:        "Conflict Severity Priority",
			Description: "Priority based on conflict severity and impact",
			Weight:      0.25,
			Calculator:  pr.calculateConflictPriority,
			Category:    CategoryConflictPriority,
		},
		{
			Name:        "Task Importance Priority",
			Description: "Priority based on task importance and category",
			Weight:      0.20,
			Calculator:  pr.calculateTaskImportancePriority,
			Category:    CategoryTaskPriority,
		},
		{
			Name:        "Timeline Urgency Priority",
			Description: "Priority based on timeline urgency and deadlines",
			Weight:      0.15,
			Calculator:  pr.calculateTimelinePriority,
			Category:    CategoryTimelinePriority,
		},
		{
			Name:        "Resource Contention Priority",
			Description: "Priority based on resource contention and assignee workload",
			Weight:      0.10,
			Calculator:  pr.calculateResourcePriority,
			Category:    CategoryResourcePriority,
		},
		{
			Name:        "Dependency Criticality Priority",
			Description: "Priority based on dependency criticality and blocking relationships",
			Weight:      0.10,
			Calculator:  pr.calculateDependencyPriority,
			Category:    CategoryDependencyPriority,
		},
		{
			Name:        "Milestone Status Priority",
			Description: "Priority based on milestone status and project importance",
			Weight:      0.10,
			Calculator:  pr.calculateMilestonePriority,
			Category:    CategoryMilestonePriority,
		},
		{
			Name:        "Assignee Workload Priority",
			Description: "Priority based on assignee workload and availability",
			Weight:      0.05,
			Calculator:  pr.calculateAssigneePriority,
			Category:    CategoryAssigneePriority,
		},
		{
			Name:        "Category Importance Priority",
			Description: "Priority based on category importance and project phase",
			Weight:      0.03,
			Calculator:  pr.calculateCategoryPriority,
			Category:    CategoryCategoryPriority,
		},
		{
			Name:        "Deadline Proximity Priority",
			Description: "Priority based on deadline proximity and urgency",
			Weight:      0.02,
			Calculator:  pr.calculateDeadlinePriority,
			Category:    CategoryDeadlinePriority,
		},
	}
}

// initializeVisualWeights sets up visual factor weights
func (pr *PriorityRanker) initializeVisualWeights() {
	pr.visualWeights = map[VisualFactor]float64{
		FactorConflictSeverity:     0.30,
		FactorTaskImportance:       0.20,
		FactorTimelineUrgency:      0.15,
		FactorResourceContention:   0.10,
		FactorDependencyCriticality: 0.10,
		FactorMilestoneStatus:      0.05,
		FactorAssigneeWorkload:     0.03,
		FactorCategoryImportance:   0.03,
		FactorDeadlineProximity:    0.02,
		FactorWorkloadBalance:      0.02,
	}
}

// RankTasks performs comprehensive priority ranking of tasks
func (pr *PriorityRanker) RankTasks(tasks []*data.Task, context *PriorityContext) *PriorityRanking {
	// Detect overlaps and categorize conflicts
	overlapDetector := pr.conflictCategorizer.overlapDetector
	overlapAnalysis := overlapDetector.DetectOverlaps(tasks)
	conflictAnalysis := pr.conflictCategorizer.CategorizeConflicts(overlapAnalysis)
	
	// Update context with analysis results
	context.OverlapAnalysis = overlapAnalysis
	context.ConflictAnalysis = conflictAnalysis
	
	// Calculate priorities for each task
	var taskPriorities []*TaskPriority
	for _, task := range tasks {
		// Get conflicts for this task
		taskConflicts := pr.getConflictsForTask(task, conflictAnalysis)
		
		// Create task-specific context
		taskContext := *context
		taskContext.Task = task
		taskContext.Conflicts = taskConflicts
		
		// Calculate priority
		taskPriority := pr.calculateTaskPriority(task, &taskContext)
		taskPriorities = append(taskPriorities, taskPriority)
	}
	
	// Sort by priority score (highest first)
	sort.Slice(taskPriorities, func(i, j int) bool {
		return taskPriorities[i].PriorityScore > taskPriorities[j].PriorityScore
	})
	
	// Assign display order
	for i, taskPriority := range taskPriorities {
		taskPriority.DisplayOrder = i + 1
	}
	
	// Generate ranking summary
	ranking := &PriorityRanking{
		TaskPriorities:    taskPriorities,
		RankingSummary:    make(map[VisualProminence]int),
		ConflictSummary:   make(map[OverlapSeverity]int),
		UrgencySummary:    make(map[string]int),
		VisualHierarchy:   make([]*TaskPriority, 0),
		CriticalTasks:     make([]*TaskPriority, 0),
		HighPriorityTasks: make([]*TaskPriority, 0),
		MediumPriorityTasks: make([]*TaskPriority, 0),
		LowPriorityTasks:  make([]*TaskPriority, 0),
		MinimalTasks:      make([]*TaskPriority, 0),
		AnalysisDate:      time.Now(),
	}
	
	// Categorize tasks by prominence
	for _, taskPriority := range taskPriorities {
		ranking.RankingSummary[taskPriority.VisualProminence]++
		
		switch taskPriority.VisualProminence {
		case ProminenceCritical:
			ranking.CriticalTasks = append(ranking.CriticalTasks, taskPriority)
		case ProminenceHigh:
			ranking.HighPriorityTasks = append(ranking.HighPriorityTasks, taskPriority)
		case ProminenceMedium:
			ranking.MediumPriorityTasks = append(ranking.MediumPriorityTasks, taskPriority)
		case ProminenceLow:
			ranking.LowPriorityTasks = append(ranking.LowPriorityTasks, taskPriority)
		case ProminenceMinimal:
			ranking.MinimalTasks = append(ranking.MinimalTasks, taskPriority)
		}
	}
	
	// Generate visual hierarchy (most prominent first)
	ranking.VisualHierarchy = append(ranking.VisualHierarchy, ranking.CriticalTasks...)
	ranking.VisualHierarchy = append(ranking.VisualHierarchy, ranking.HighPriorityTasks...)
	ranking.VisualHierarchy = append(ranking.VisualHierarchy, ranking.MediumPriorityTasks...)
	ranking.VisualHierarchy = append(ranking.VisualHierarchy, ranking.LowPriorityTasks...)
	ranking.VisualHierarchy = append(ranking.VisualHierarchy, ranking.MinimalTasks...)
	
	// Generate recommendations
	ranking.Recommendations = pr.generateRankingRecommendations(ranking)
	
	return ranking
}

// calculateTaskPriority calculates comprehensive priority for a single task
func (pr *PriorityRanker) calculateTaskPriority(task *data.Task, context *PriorityContext) *TaskPriority {
	taskPriority := &TaskPriority{
		Task:           task,
		RankingFactors: make(map[PriorityCategory]float64),
		Recommendations: make([]string, 0),
	}
	
	// Calculate priority using all rules
	totalWeight := 0.0
	for _, rule := range pr.rankingRules {
		score := rule.Calculator(task, context)
		taskPriority.RankingFactors[rule.Category] = score
		taskPriority.PriorityScore += score * rule.Weight
		totalWeight += rule.Weight
	}
	
	// Normalize score
	if totalWeight > 0 {
		taskPriority.PriorityScore /= totalWeight
	}
	
	// Calculate specific weights
	taskPriority.ConflictImpact = taskPriority.RankingFactors[CategoryConflictPriority]
	taskPriority.TimelineUrgency = taskPriority.RankingFactors[CategoryTimelinePriority]
	taskPriority.ResourceContention = taskPriority.RankingFactors[CategoryResourcePriority]
	taskPriority.DependencyWeight = taskPriority.RankingFactors[CategoryDependencyPriority]
	taskPriority.MilestoneWeight = taskPriority.RankingFactors[CategoryMilestonePriority]
	taskPriority.AssigneeWeight = taskPriority.RankingFactors[CategoryAssigneePriority]
	taskPriority.CategoryWeight = taskPriority.RankingFactors[CategoryCategoryPriority]
	taskPriority.DeadlineWeight = taskPriority.RankingFactors[CategoryDeadlinePriority]
	taskPriority.WorkloadWeight = taskPriority.RankingFactors[CategoryWorkloadPriority]
	
	// Determine visual prominence
	taskPriority.VisualProminence = pr.determineVisualProminence(taskPriority)
	
	// Generate visual style
	taskPriority.VisualStyle = pr.generateVisualStyle(taskPriority)
	
	// Generate recommendations
	taskPriority.Recommendations = pr.generateTaskRecommendations(taskPriority)
	
	return taskPriority
}

// calculateConflictPriority calculates priority based on conflict severity
func (pr *PriorityRanker) calculateConflictPriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	
	// Base score from task priority
	score += float64(task.Priority) * 0.2
	
	// Conflict impact
	for _, conflict := range context.Conflicts {
		switch conflict.Severity {
		case SeverityCritical:
			score += 10.0
		case SeverityHigh:
			score += 7.0
		case SeverityMedium:
			score += 4.0
		case SeverityLow:
			score += 2.0
		}
		
		// Additional weight for high-impact conflicts
		if conflict.Impact == "HIGH" {
			score += 3.0
		} else if conflict.Impact == "MEDIUM" {
			score += 1.5
		}
	}
	
	// Risk level impact
	for _, conflict := range context.Conflicts {
		switch conflict.RiskLevel {
		case "HIGH":
			score += 5.0
		case "MEDIUM":
			score += 2.0
		case "LOW":
			score += 0.5
		}
	}
	
	return score
}

// calculateTaskImportancePriority calculates priority based on task importance
func (pr *PriorityRanker) calculateTaskImportancePriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	
	// Base task priority
	score += float64(task.Priority) * 2.0
	
	// Category importance
	if importance, exists := context.CategoryImportance[task.Category]; exists {
		score += importance * 3.0
	} else {
		// Default category weights
		switch task.Category {
		case "DISSERTATION":
			score += 8.0
		case "PROPOSAL":
			score += 6.0
		case "LASER":
			score += 4.0
		case "MEETING":
			score += 2.0
		default:
			score += 1.0
		}
	}
	
	// Milestone status
	if task.IsMilestone {
		score += 10.0
	}
	
	// Task duration (longer tasks get higher priority)
	duration := task.EndDate.Sub(task.StartDate)
	if duration > time.Hour*24*30 { // More than 30 days
		score += 3.0
	} else if duration > time.Hour*24*7 { // More than 7 days
		score += 2.0
	} else if duration > time.Hour*24 { // More than 1 day
		score += 1.0
	}
	
	return score
}

// calculateTimelinePriority calculates priority based on timeline urgency
func (pr *PriorityRanker) calculateTimelinePriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	now := context.CurrentTime
	
	// Days until start
	daysUntilStart := int(task.StartDate.Sub(now).Hours() / 24)
	if daysUntilStart <= 0 {
		score += 10.0 // Already started or overdue
	} else if daysUntilStart <= 1 {
		score += 8.0 // Starting tomorrow
	} else if daysUntilStart <= 3 {
		score += 6.0 // Starting in 3 days
	} else if daysUntilStart <= 7 {
		score += 4.0 // Starting in a week
	} else if daysUntilStart <= 14 {
		score += 2.0 // Starting in 2 weeks
	}
	
	// Days until end
	daysUntilEnd := int(task.EndDate.Sub(now).Hours() / 24)
	if daysUntilEnd <= 0 {
		score += 15.0 // Overdue
	} else if daysUntilEnd <= 1 {
		score += 12.0 // Due tomorrow
	} else if daysUntilEnd <= 3 {
		score += 8.0 // Due in 3 days
	} else if daysUntilEnd <= 7 {
		score += 5.0 // Due in a week
	} else if daysUntilEnd <= 14 {
		score += 3.0 // Due in 2 weeks
	}
	
	return score
}

// calculateResourcePriority calculates priority based on resource contention
func (pr *PriorityRanker) calculateResourcePriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	
	// Assignee workload
	if task.Assignee != "" {
		if workload, exists := context.AssigneeWorkloads[task.Assignee]; exists {
			// Higher workload = higher priority (more contention)
			score += float64(workload) * 0.5
		}
	}
	
	// Resource contention from conflicts
	for _, conflict := range context.Conflicts {
		if conflict.Category == CategoryAssigneeConflict {
			score += 5.0
		} else if conflict.Category == CategoryResourceConflict {
			score += 3.0
		}
	}
	
	return score
}

// calculateDependencyPriority calculates priority based on dependency criticality
func (pr *PriorityRanker) calculateDependencyPriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	
	// Number of dependencies (more dependencies = higher priority)
	score += float64(len(task.Dependencies)) * 1.0
	
	// Check if this task blocks other tasks
	if context.OverlapAnalysis != nil {
		for _, group := range context.OverlapAnalysis.OverlapGroups {
			for _, otherTask := range group.Tasks {
				for _, dep := range otherTask.Dependencies {
					if dep == task.ID {
						score += 3.0 // This task blocks another task
					}
				}
			}
		}
	}
	
	// Dependency conflicts
	for _, conflict := range context.Conflicts {
		if conflict.Category == CategoryDependencyConflict {
			score += 8.0
		}
	}
	
	return score
}

// calculateMilestonePriority calculates priority based on milestone status
func (pr *PriorityRanker) calculateMilestonePriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	
	if task.IsMilestone {
		score += 15.0
		
		// Additional weight for critical milestones
		if task.Priority >= 4 {
			score += 5.0
		}
		
		// Timeline urgency for milestones
		now := context.CurrentTime
		daysUntilEnd := int(task.EndDate.Sub(now).Hours() / 24)
		if daysUntilEnd <= 7 {
			score += 10.0
		} else if daysUntilEnd <= 30 {
			score += 5.0
		}
	}
	
	return score
}

// calculateAssigneePriority calculates priority based on assignee workload
func (pr *PriorityRanker) calculateAssigneePriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	
	if task.Assignee != "" {
		// Base score for having an assignee
		score += 2.0
		
		// Workload consideration
		if workload, exists := context.AssigneeWorkloads[task.Assignee]; exists {
			// Moderate workload gets higher priority
			if workload >= 3 && workload <= 6 {
				score += 3.0
			} else if workload > 6 {
				score += 1.0 // Overloaded assignee
			}
		}
	}
	
	return score
}

// calculateCategoryPriority calculates priority based on category importance
func (pr *PriorityRanker) calculateCategoryPriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	
	// Category-specific weights
	switch task.Category {
	case "DISSERTATION":
		score += 8.0
	case "PROPOSAL":
		score += 6.0
	case "LASER":
		score += 4.0
	case "MEETING":
		score += 2.0
	case "ADMIN":
		score += 1.0
	default:
		score += 0.5
	}
	
	return score
}

// calculateDeadlinePriority calculates priority based on deadline proximity
func (pr *PriorityRanker) calculateDeadlinePriority(task *data.Task, context *PriorityContext) float64 {
	score := 0.0
	now := context.CurrentTime
	
	// Days until deadline
	daysUntilDeadline := int(task.EndDate.Sub(now).Hours() / 24)
	if daysUntilDeadline <= 0 {
		score += 15.0 // Overdue
	} else if daysUntilDeadline <= 1 {
		score += 12.0 // Due tomorrow
	} else if daysUntilDeadline <= 3 {
		score += 8.0 // Due in 3 days
	} else if daysUntilDeadline <= 7 {
		score += 5.0 // Due in a week
	} else if daysUntilDeadline <= 14 {
		score += 3.0 // Due in 2 weeks
	} else if daysUntilDeadline <= 30 {
		score += 1.0 // Due in a month
	}
	
	return score
}

// determineVisualProminence determines the visual prominence level
func (pr *PriorityRanker) determineVisualProminence(taskPriority *TaskPriority) VisualProminence {
	score := taskPriority.PriorityScore
	
	if score >= 15.0 {
		return ProminenceCritical
	} else if score >= 10.0 {
		return ProminenceHigh
	} else if score >= 6.0 {
		return ProminenceMedium
	} else if score >= 3.0 {
		return ProminenceLow
	} else {
		return ProminenceMinimal
	}
}

// generateVisualStyle generates visual styling for the task
func (pr *PriorityRanker) generateVisualStyle(taskPriority *TaskPriority) VisualStyle {
	style := VisualStyle{
		Opacity:    1.0,
		FontWeight: "normal",
		FontSize:   "normal",
		ZIndex:     1,
		Highlight:  false,
		Blink:      false,
		Pulse:      false,
		Glow:       false,
	}
	
	switch taskPriority.VisualProminence {
	case ProminenceCritical:
		style.BorderColor = "red"
		style.FillColor = "red!20"
		style.BorderWidth = "3pt"
		style.FontWeight = "bold"
		style.FontSize = "large"
		style.ZIndex = 10
		style.Highlight = true
		style.Blink = true
		style.Glow = true
	case ProminenceHigh:
		style.BorderColor = "orange"
		style.FillColor = "orange!15"
		style.BorderWidth = "2pt"
		style.FontWeight = "bold"
		style.FontSize = "large"
		style.ZIndex = 8
		style.Highlight = true
		style.Pulse = true
	case ProminenceMedium:
		style.BorderColor = "yellow"
		style.FillColor = "yellow!10"
		style.BorderWidth = "1.5pt"
		style.FontWeight = "semibold"
		style.ZIndex = 6
		style.Highlight = true
	case ProminenceLow:
		style.BorderColor = "blue"
		style.FillColor = "blue!5"
		style.BorderWidth = "1pt"
		style.ZIndex = 4
	case ProminenceMinimal:
		style.BorderColor = "gray"
		style.FillColor = "gray!5"
		style.BorderWidth = "0.5pt"
		style.Opacity = 0.7
		style.ZIndex = 2
	}
	
	return style
}

// generateTaskRecommendations generates recommendations for a specific task
func (pr *PriorityRanker) generateTaskRecommendations(taskPriority *TaskPriority) []string {
	var recommendations []string
	
	// Conflict-based recommendations
	if taskPriority.ConflictImpact > 5.0 {
		recommendations = append(recommendations, "🚨 High conflict impact - resolve conflicts immediately")
	}
	
	// Timeline recommendations
	if taskPriority.TimelineUrgency > 8.0 {
		recommendations = append(recommendations, "⏰ Urgent timeline - prioritize execution")
	}
	
	// Resource recommendations
	if taskPriority.ResourceContention > 3.0 {
		recommendations = append(recommendations, "👥 Resource contention - consider reassignment")
	}
	
	// Dependency recommendations
	if taskPriority.DependencyWeight > 5.0 {
		recommendations = append(recommendations, "🔗 Critical dependencies - ensure prerequisites are met")
	}
	
	// Milestone recommendations
	if taskPriority.MilestoneWeight > 10.0 {
		recommendations = append(recommendations, "🎯 Milestone task - ensure completion on time")
	}
	
	// Workload recommendations
	if taskPriority.WorkloadWeight > 3.0 {
		recommendations = append(recommendations, "⚖️ High workload - consider task splitting")
	}
	
	return recommendations
}

// generateRankingRecommendations generates overall ranking recommendations
func (pr *PriorityRanker) generateRankingRecommendations(ranking *PriorityRanking) []string {
	var recommendations []string
	
	// Critical tasks
	if len(ranking.CriticalTasks) > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🚨 %d critical tasks require immediate attention", len(ranking.CriticalTasks)))
	}
	
	// High priority tasks
	if len(ranking.HighPriorityTasks) > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🔴 %d high priority tasks need focus", len(ranking.HighPriorityTasks)))
	}
	
	// Conflict summary
	if ranking.ConflictSummary[SeverityCritical] > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("⚠️ %d critical conflicts need resolution", ranking.ConflictSummary[SeverityCritical]))
	}
	
	// Visual hierarchy
	recommendations = append(recommendations, 
		"📊 Use visual hierarchy to guide attention to most important tasks")
	
	// Resource balancing
	totalTasks := len(ranking.TaskPriorities)
	if totalTasks > 0 {
		criticalRatio := float64(len(ranking.CriticalTasks)) / float64(totalTasks)
		if criticalRatio > 0.3 {
			recommendations = append(recommendations, 
				"⚖️ High ratio of critical tasks - consider resource reallocation")
		}
	}
	
	return recommendations
}

// getConflictsForTask retrieves conflicts related to a specific task
func (pr *PriorityRanker) getConflictsForTask(task *data.Task, conflictAnalysis *ConflictAnalysis) []*CategorizedConflict {
	var conflicts []*CategorizedConflict
	
	for _, conflict := range conflictAnalysis.CategorizedConflicts {
		if conflict.Task1ID == task.ID || conflict.Task2ID == task.ID {
			conflicts = append(conflicts, conflict)
		}
	}
	
	return conflicts
}

// AddCustomRule adds a custom priority ranking rule
func (pr *PriorityRanker) AddCustomRule(rule PriorityRule) {
	pr.rankingRules = append(pr.rankingRules, rule)
	// Sort rules by weight (highest first)
	sort.Slice(pr.rankingRules, func(i, j int) bool {
		return pr.rankingRules[i].Weight > pr.rankingRules[j].Weight
	})
}

// GetTasksByProminence returns tasks filtered by visual prominence
func (ranking *PriorityRanking) GetTasksByProminence(prominence VisualProminence) []*TaskPriority {
	var filtered []*TaskPriority
	for _, taskPriority := range ranking.TaskPriorities {
		if taskPriority.VisualProminence == prominence {
			filtered = append(filtered, taskPriority)
		}
	}
	return filtered
}

// GetTopTasks returns the top N tasks by priority
func (ranking *PriorityRanking) GetTopTasks(n int) []*TaskPriority {
	if n > len(ranking.TaskPriorities) {
		n = len(ranking.TaskPriorities)
	}
	return ranking.TaskPriorities[:n]
}

// GetSummary returns a summary of the priority ranking
func (ranking *PriorityRanking) GetSummary() string {
	return fmt.Sprintf("Priority Ranking Summary:\n"+
		"  Total Tasks: %d\n"+
		"  Critical: %d, High: %d, Medium: %d, Low: %d, Minimal: %d\n"+
		"  Critical Conflicts: %d\n"+
		"  High Conflicts: %d\n"+
		"  Analysis Date: %s",
		len(ranking.TaskPriorities),
		len(ranking.CriticalTasks),
		len(ranking.HighPriorityTasks),
		len(ranking.MediumPriorityTasks),
		len(ranking.LowPriorityTasks),
		len(ranking.MinimalTasks),
		ranking.ConflictSummary[SeverityCritical],
		ranking.ConflictSummary[SeverityHigh],
		ranking.AnalysisDate.Format("2006-01-02 15:04:05"))
}
