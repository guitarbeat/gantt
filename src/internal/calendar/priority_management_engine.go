package calendar

import (
	"fmt"
	"math"
	"sort"
	"time"

	"phd-dissertation-planner/internal/shared"
)

// PriorityManagementEngine handles both priority ranking and task prioritization
type PriorityManagementEngine struct {
	// Priority ranking components
	conflictCategorizer *ConflictCategorizer
	rankingRules        []PriorityRule
	visualWeights       map[VisualFactor]float64

	// Task prioritization components
	stackingEngine         *StackingEngine
	priorityRanker         *PriorityRanker
	visibilityManager      *VisibilityManager
	stackingOptimizer      *StackingOptimizer
}

// TaskPrioritizationEngine handles intelligent task prioritization for stacking order
type TaskPrioritizationEngine struct {
	stackingEngine         *StackingEngine
	priorityRanker         *PriorityRanker
	visibilityManager      *VisibilityManager
	stackingOptimizer      *StackingOptimizer
}

// PriorityRule defines a rule for calculating task priority scores
type PriorityRule struct {
	Name        string
	Description string
	Weight      float64
	Calculator  func(*shared.Task, *PriorityContext) float64
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
	FactorDependencyImpact    VisualFactor = "DEPENDENCY_IMPACT"
	FactorMilestoneSignificance VisualFactor = "MILESTONE_SIGNIFICANCE"
	FactorAssigneeWorkload    VisualFactor = "ASSIGNEE_WORKLOAD"
	FactorCategoryRelevance   VisualFactor = "CATEGORY_RELEVANCE"
	FactorDeadlinePressure    VisualFactor = "DEADLINE_PRESSURE"
	FactorWorkloadBalance     VisualFactor = "WORKLOAD_BALANCE"
)

// PriorityRanker handles priority ranking and visual prominence decisions
type PriorityRanker struct {
	conflictCategorizer *ConflictCategorizer
	rankingRules        []PriorityRule
	visualWeights       map[VisualFactor]float64
}

// VisibilityManager manages task visibility and prominence
type VisibilityManager struct {
	visibilityRules    []VisibilityRule
	prominenceWeights  map[VisualProminence]float64
	visibilityThreshold float64
	adaptiveVisibility bool
}

// StackingOptimizer optimizes stacking order based on priorities
type StackingOptimizer struct {
	optimizationRules []OptimizationRule
	stackingStrategies map[PriorityCategory]StackingStrategy
	adaptiveOrdering   bool
}

// VisibilityRule defines a rule for task visibility
type VisibilityRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*shared.Task, *PriorityContext) bool
	Action      func(*shared.Task, *PriorityContext) *VisibilityAction
}

// VisibilityAction defines how a task should be made visible
type VisibilityAction struct {
	IsVisible       bool
	ProminenceLevel VisualProminence
	DisplayOrder    int
	VisualWeight    float64
	CollapseLevel   int
}

// OptimizationRule defines a rule for stacking optimization
type OptimizationRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*shared.Task, *PriorityContext) bool
	Action      func(*shared.Task, *PriorityContext) *OptimizationAction
}

// OptimizationAction defines how a task should be optimized
type OptimizationAction struct {
	StackingOrder    int
	VisualProminence VisualProminence
	DisplayPriority  float64
	CollapseLevel    int
}

// StackingStrategy defines a strategy for stacking tasks
type StackingStrategy struct {
	StrategyType PriorityCategory
	Description  string
	Parameters   map[string]interface{}
	Rules        []OptimizationRule
}

// PriorityContext provides context for priority calculations
type PriorityContext struct {
	CurrentTime    time.Time
	UserID         string
	ProjectID      string
	TeamMembers    []string
	ResourceLimits map[string]int
	DeadlineConstraints map[string]time.Time
	WorkloadLimits map[string]float64
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
	Category    PriorityCategory
	Factor      VisualFactor
	Value       float64
	Weight      float64
	Contribution float64
	Description string
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

// PriorityRankingResult contains the result of priority ranking
type PriorityRankingResult struct {
	TaskScores        []*PriorityScore
	RankingOrder      []string
	ConflictsDetected []string
	Recommendations   []string
	AnalysisDate      time.Time
}

// TaskPrioritizationResult contains the result of task prioritization
type TaskPrioritizationResult struct {
	PrioritizedTasks  []*PrioritizedTask
	StackingOrder     []string
	VisibilitySettings map[string]*VisibilityAction
	OptimizationResults []*OptimizationAction
	Recommendations   []string
	AnalysisDate      time.Time
}

// PrioritizedTask represents a task with prioritization information
type PrioritizedTask struct {
	Task              *shared.Task
	PriorityScore     *PriorityScore
	VisibilityAction  *VisibilityAction
	OptimizationAction *OptimizationAction
	StackingOrder     int
	DisplayPriority   float64
}

// NewPriorityManagementEngine creates a new priority management engine
func NewPriorityManagementEngine(
	conflictCategorizer *ConflictCategorizer,
	stackingEngine *StackingEngine,
) *PriorityManagementEngine {
	engine := &PriorityManagementEngine{
		conflictCategorizer: conflictCategorizer,
		stackingEngine:      stackingEngine,
		rankingRules:        make([]PriorityRule, 0),
		visualWeights:       make(map[VisualFactor]float64),
		priorityRanker:      NewPriorityRanker(conflictCategorizer),
		visibilityManager:   NewVisibilityManager(),
		stackingOptimizer:   NewStackingOptimizer(),
	}
	
	// Initialize default rules and weights
	engine.initializeDefaultRules()
	engine.initializeVisualWeights()
	
	return engine
}

// CalculatePriorityScores calculates priority scores for all tasks
func (pme *PriorityManagementEngine) CalculatePriorityScores(tasks []*shared.Task, context *PriorityContext) *PriorityRankingResult {
	return pme.priorityRanker.CalculatePriorityScores(tasks, context)
}

// PrioritizeTasks performs intelligent task prioritization for stacking order
func (pme *PriorityManagementEngine) PrioritizeTasks(tasks []*shared.Task, context *PriorityContext) *TaskPrioritizationResult {
	// Calculate priority scores
	priorityResult := pme.CalculatePriorityScores(tasks, context)
	
	// Apply visibility management
	visibilitySettings := pme.visibilityManager.ApplyVisibilityRules(tasks, context)
	
	// Optimize stacking order
	optimizationResults := pme.stackingOptimizer.OptimizeStackingOrder(tasks, context)
	
	// Create prioritized tasks
	prioritizedTasks := make([]*PrioritizedTask, 0, len(tasks))
	
	for i, task := range tasks {
		priorityScore := priorityResult.TaskScores[i]
		visibilityAction := visibilitySettings[task.ID]
		optimizationAction := optimizationResults[i]
		
		prioritizedTask := &PrioritizedTask{
			Task:              task,
			PriorityScore:     priorityScore,
			VisibilityAction:  visibilityAction,
			OptimizationAction: optimizationAction,
			StackingOrder:     optimizationAction.StackingOrder,
			DisplayPriority:   optimizationAction.DisplayPriority,
		}
		
		prioritizedTasks = append(prioritizedTasks, prioritizedTask)
	}
	
	// Sort by stacking order
	sort.Slice(prioritizedTasks, func(i, j int) bool {
		return prioritizedTasks[i].StackingOrder < prioritizedTasks[j].StackingOrder
	})
	
	// Generate recommendations
	recommendations := pme.generatePrioritizationRecommendations(prioritizedTasks, priorityResult)
	
	return &TaskPrioritizationResult{
		PrioritizedTasks:    prioritizedTasks,
		StackingOrder:       extractTaskIDs(prioritizedTasks),
		VisibilitySettings:  visibilitySettings,
		OptimizationResults: optimizationResults,
		Recommendations:     recommendations,
		AnalysisDate:        time.Now(),
	}
}

// initializeDefaultRules sets up default priority rules
func (pme *PriorityManagementEngine) initializeDefaultRules() {
	pme.rankingRules = []PriorityRule{
		{
			Name:        "Task Priority Score",
			Description: "Base priority from task definition",
			Weight:      0.3,
			Calculator:  pme.calculateTaskPriorityScore,
			Category:    CategoryTaskPriority,
		},
		{
			Name:        "Timeline Urgency",
			Description: "Urgency based on timeline proximity",
			Weight:      0.25,
			Calculator:  pme.calculateTimelineUrgencyScore,
			Category:    CategoryTimelinePriority,
		},
		{
			Name:        "Milestone Significance",
			Description: "Priority boost for milestone tasks",
			Weight:      0.2,
			Calculator:  pme.calculateMilestoneScore,
			Category:    CategoryMilestonePriority,
		},
		{
			Name:        "Deadline Pressure",
			Description: "Pressure based on deadline proximity",
			Weight:      0.15,
			Calculator:  pme.calculateDeadlinePressureScore,
			Category:    CategoryDeadlinePriority,
		},
		{
			Name:        "Dependency Impact",
			Description: "Impact based on task dependencies",
			Weight:      0.1,
			Calculator:  pme.calculateDependencyImpactScore,
			Category:    CategoryDependencyPriority,
		},
	}
}

// initializeVisualWeights sets up default visual weights
func (pme *PriorityManagementEngine) initializeVisualWeights() {
	pme.visualWeights = map[VisualFactor]float64{
		FactorConflictSeverity:      0.3,
		FactorTaskImportance:        0.25,
		FactorTimelineUrgency:       0.2,
		FactorResourceContention:    0.15,
		FactorDependencyImpact:      0.1,
		FactorMilestoneSignificance: 0.2,
		FactorAssigneeWorkload:      0.1,
		FactorCategoryRelevance:     0.05,
		FactorDeadlinePressure:      0.15,
		FactorWorkloadBalance:       0.1,
	}
}

// Priority calculation methods
func (pme *PriorityManagementEngine) calculateTaskPriorityScore(task *shared.Task, context *PriorityContext) float64 {
	return float64(task.Priority) / 5.0 // Normalize to 0-1 range
}

func (pme *PriorityManagementEngine) calculateTimelineUrgencyScore(task *shared.Task, context *PriorityContext) float64 {
	now := context.CurrentTime
	daysUntilStart := int(task.StartDate.Sub(now).Hours() / 24)
	
	if daysUntilStart <= 0 {
		return 1.0 // Already started or overdue
	} else if daysUntilStart <= 3 {
		return 0.9 // Very urgent
	} else if daysUntilStart <= 7 {
		return 0.7 // Urgent
	} else if daysUntilStart <= 14 {
		return 0.5 // Moderate urgency
	} else if daysUntilStart <= 30 {
		return 0.3 // Low urgency
	} else {
		return 0.1 // Very low urgency
	}
}

func (pme *PriorityManagementEngine) calculateMilestoneScore(task *shared.Task, context *PriorityContext) float64 {
	if task.IsMilestone {
		return 1.0
	}
	return 0.0
}

func (pme *PriorityManagementEngine) calculateDeadlinePressureScore(task *shared.Task, context *PriorityContext) float64 {
	now := context.CurrentTime
	daysUntilDeadline := int(task.EndDate.Sub(now).Hours() / 24)
	
	if daysUntilDeadline <= 0 {
		return 1.0 // Overdue
	} else if daysUntilDeadline <= 3 {
		return 0.9 // Very high pressure
	} else if daysUntilDeadline <= 7 {
		return 0.7 // High pressure
	} else if daysUntilDeadline <= 14 {
		return 0.5 // Moderate pressure
	} else if daysUntilDeadline <= 30 {
		return 0.3 // Low pressure
	} else {
		return 0.1 // Very low pressure
	}
}

func (pme *PriorityManagementEngine) calculateDependencyImpactScore(task *shared.Task, context *PriorityContext) float64 {
	// Higher score for tasks with more dependencies (blocking more other tasks)
	return math.Min(float64(len(task.Dependencies))/10.0, 1.0)
}

// generatePrioritizationRecommendations generates recommendations for task prioritization
func (pme *PriorityManagementEngine) generatePrioritizationRecommendations(prioritizedTasks []*PrioritizedTask, priorityResult *PriorityRankingResult) []string {
	var recommendations []string
	
	// High priority tasks
	highPriorityCount := 0
	for _, task := range prioritizedTasks {
		if task.PriorityScore.OverallScore >= 0.8 {
			highPriorityCount++
		}
	}
	
	if highPriorityCount > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🎯 %d high-priority tasks require immediate attention", highPriorityCount))
	}
	
	// Milestone tasks
	milestoneCount := 0
	for _, task := range prioritizedTasks {
		if task.Task.IsMilestone {
			milestoneCount++
		}
	}
	
	if milestoneCount > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🏁 %d milestone tasks should be prioritized", milestoneCount))
	}
	
	// Overdue tasks
	overdueCount := 0
	now := time.Now()
	for _, task := range prioritizedTasks {
		if task.Task.EndDate.Before(now) {
			overdueCount++
		}
	}
	
	if overdueCount > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("⚠️ %d overdue tasks need immediate resolution", overdueCount))
	}
	
	// General recommendations
	if len(prioritizedTasks) > 0 {
		recommendations = append(recommendations, 
			"📊 Review task priorities regularly to ensure optimal resource allocation")
		recommendations = append(recommendations, 
			"🔄 Consider adjusting task priorities based on changing project requirements")
	}
	
	return recommendations
}

// Helper function to extract task IDs from prioritized tasks
func extractTaskIDs(prioritizedTasks []*PrioritizedTask) []string {
	ids := make([]string, len(prioritizedTasks))
	for i, task := range prioritizedTasks {
		ids[i] = task.Task.ID
	}
	return ids
}

// NewPriorityRanker creates a new priority ranker
func NewPriorityRanker(conflictCategorizer *ConflictCategorizer) *PriorityRanker {
	return &PriorityRanker{
		conflictCategorizer: conflictCategorizer,
		rankingRules:        make([]PriorityRule, 0),
		visualWeights:       make(map[VisualFactor]float64),
	}
}

// CalculatePriorityScores calculates priority scores for all tasks
func (pr *PriorityRanker) CalculatePriorityScores(tasks []*shared.Task, context *PriorityContext) *PriorityRankingResult {
	taskScores := make([]*PriorityScore, 0, len(tasks))
	
	for _, task := range tasks {
		score := pr.calculateTaskPriorityScore(task, context)
		taskScores = append(taskScores, score)
	}
	
	// Sort by overall score (highest first)
	sort.Slice(taskScores, func(i, j int) bool {
		return taskScores[i].OverallScore > taskScores[j].OverallScore
	})
	
	// Assign rankings
	for i, score := range taskScores {
		score.Ranking = i + 1
	}
	
	// Generate recommendations
	recommendations := pr.generatePriorityRecommendations(taskScores)
	
	return &PriorityRankingResult{
		TaskScores:        taskScores,
		RankingOrder:      extractTaskIDsFromScores(taskScores),
		ConflictsDetected: []string{}, // Would be populated by conflict analysis
		Recommendations:   recommendations,
		AnalysisDate:      time.Now(),
	}
}

// calculateTaskPriorityScore calculates priority score for a single task
func (pr *PriorityRanker) calculateTaskPriorityScore(task *shared.Task, context *PriorityContext) *PriorityScore {
	overallScore := 0.0
	categoryScores := make(map[PriorityCategory]float64)
	factors := make([]PriorityFactor, 0)
	
	// Apply each ranking rule
	for _, rule := range pr.rankingRules {
		score := rule.Calculator(task, context)
		weightedScore := score * rule.Weight
		overallScore += weightedScore
		categoryScores[rule.Category] = score
		
		factors = append(factors, PriorityFactor{
			Category:     rule.Category,
			Factor:       VisualFactor(rule.Name), // Simplified mapping
			Value:        score,
			Weight:       rule.Weight,
			Contribution: weightedScore,
			Description:  rule.Description,
		})
	}
	
	// Determine visual prominence
	var prominence VisualProminence
	if overallScore >= 0.9 {
		prominence = ProminenceCritical
	} else if overallScore >= 0.7 {
		prominence = ProminenceHigh
	} else if overallScore >= 0.5 {
		prominence = ProminenceMedium
	} else if overallScore >= 0.3 {
		prominence = ProminenceLow
	} else {
		prominence = ProminenceMinimal
	}
	
	return &PriorityScore{
		TaskID:           task.ID,
		OverallScore:     overallScore,
		CategoryScores:   categoryScores,
		VisualProminence: prominence,
		Ranking:          0, // Will be set later
		Confidence:       0.8, // Simplified confidence calculation
		Factors:          factors,
		Recommendations:  []string{}, // Would be populated by analysis
		CalculatedAt:     time.Now(),
	}
}

// generatePriorityRecommendations generates recommendations for priority ranking
func (pr *PriorityRanker) generatePriorityRecommendations(scores []*PriorityScore) []string {
	var recommendations []string
	
	// High priority tasks
	highPriorityCount := 0
	for _, score := range scores {
		if score.OverallScore >= 0.8 {
			highPriorityCount++
		}
	}
	
	if highPriorityCount > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🎯 %d high-priority tasks require immediate attention", highPriorityCount))
	}
	
	// Critical tasks
	criticalCount := 0
	for _, score := range scores {
		if score.VisualProminence == ProminenceCritical {
			criticalCount++
		}
	}
	
	if criticalCount > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🚨 %d critical tasks need urgent resolution", criticalCount))
	}
	
	// General recommendations
	if len(scores) > 0 {
		recommendations = append(recommendations, 
			"📊 Review priority scores regularly to ensure optimal task ordering")
		recommendations = append(recommendations, 
			"🔄 Consider adjusting task priorities based on changing project requirements")
	}
	
	return recommendations
}

// NewVisibilityManager creates a new visibility manager
func NewVisibilityManager() *VisibilityManager {
	return &VisibilityManager{
		visibilityRules:    make([]VisibilityRule, 0),
		prominenceWeights:  make(map[VisualProminence]float64),
		visibilityThreshold: 0.5,
		adaptiveVisibility: true,
	}
}

// ApplyVisibilityRules applies visibility rules to tasks
func (vm *VisibilityManager) ApplyVisibilityRules(tasks []*shared.Task, context *PriorityContext) map[string]*VisibilityAction {
	visibilitySettings := make(map[string]*VisibilityAction)
	
	for _, task := range tasks {
		action := vm.determineVisibilityAction(task, context)
		visibilitySettings[task.ID] = action
	}
	
	return visibilitySettings
}

// determineVisibilityAction determines visibility action for a task
func (vm *VisibilityManager) determineVisibilityAction(task *shared.Task, context *PriorityContext) *VisibilityAction {
	// Simplified visibility logic
	isVisible := true
	prominence := ProminenceMedium
	displayOrder := 0
	visualWeight := 0.5
	collapseLevel := 0
	
	// Adjust based on task properties
	if task.IsMilestone {
		prominence = ProminenceHigh
		visualWeight = 0.9
	}
	
	if task.Priority >= 4 {
		prominence = ProminenceHigh
		visualWeight = 0.8
	}
	
	// Check if task is overdue
	if task.EndDate.Before(context.CurrentTime) {
		prominence = ProminenceCritical
		visualWeight = 1.0
	}
	
	return &VisibilityAction{
		IsVisible:       isVisible,
		ProminenceLevel: prominence,
		DisplayOrder:    displayOrder,
		VisualWeight:    visualWeight,
		CollapseLevel:   collapseLevel,
	}
}

// NewStackingOptimizer creates a new stacking optimizer
func NewStackingOptimizer() *StackingOptimizer {
	return &StackingOptimizer{
		optimizationRules:  make([]OptimizationRule, 0),
		stackingStrategies: make(map[PriorityCategory]StackingStrategy),
		adaptiveOrdering:   true,
	}
}

// OptimizeStackingOrder optimizes stacking order for tasks
func (so *StackingOptimizer) OptimizeStackingOrder(tasks []*shared.Task, context *PriorityContext) []*OptimizationAction {
	actions := make([]*OptimizationAction, len(tasks))
	
	for i, task := range tasks {
		action := so.determineOptimizationAction(task, context, i)
		actions[i] = action
	}
	
	return actions
}

// determineOptimizationAction determines optimization action for a task
func (so *StackingOptimizer) determineOptimizationAction(task *shared.Task, context *PriorityContext, index int) *OptimizationAction {
	// Simplified optimization logic
	stackingOrder := index
	prominence := ProminenceMedium
	displayPriority := 0.5
	collapseLevel := 0
	
	// Adjust based on task properties
	if task.IsMilestone {
		prominence = ProminenceHigh
		displayPriority = 0.9
		stackingOrder = 0 // Put milestones first
	}
	
	if task.Priority >= 4 {
		prominence = ProminenceHigh
		displayPriority = 0.8
	}
	
	// Check if task is overdue
	if task.EndDate.Before(context.CurrentTime) {
		prominence = ProminenceCritical
		displayPriority = 1.0
		stackingOrder = 0 // Put overdue tasks first
	}
	
	return &OptimizationAction{
		StackingOrder:    stackingOrder,
		VisualProminence: prominence,
		DisplayPriority:  displayPriority,
		CollapseLevel:    collapseLevel,
	}
}

// Helper function to extract task IDs from priority scores
func extractTaskIDsFromScores(scores []*PriorityScore) []string {
	ids := make([]string, len(scores))
	for i, score := range scores {
		ids[i] = score.TaskID
	}
	return ids
}

// Additional missing types

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

// TaskPriority represents task priority information
type TaskPriority struct {
	Value       int
	Category    string
	Description string
	Weight      float64
	Urgency     string
	Importance  string
}

// NewTaskPrioritizationEngine creates a new task prioritization engine
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
