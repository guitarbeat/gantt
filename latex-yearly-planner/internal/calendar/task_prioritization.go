package calendar

import (
	"fmt"
	"math"
	"sort"
	"time"

	"latex-yearly-planner/internal/data"
)

// TaskPrioritizationEngine handles intelligent task prioritization for stacking order
type TaskPrioritizationEngine struct {
	verticalStackingEngine *VerticalStackingEngine
	priorityRanker         *PriorityRanker
	visibilityManager      *VisibilityManager
	stackingOptimizer      *StackingOptimizer
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
	Condition   func(*data.Task, *PriorityContext) bool
	Action      func(*data.Task, *PriorityContext) *VisibilityAction
}

// VisibilityAction defines how a task should be made visible
type VisibilityAction struct {
	IsVisible       bool
	ProminenceLevel VisualProminence
	DisplayOrder    int
	VisualWeight    float64
	CollapseLevel   int
	HighlightLevel  int
}

// OptimizationRule defines a rule for stacking optimization
type OptimizationRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func([]*data.Task, *PriorityContext) bool
	Action      func([]*data.Task, *PriorityContext) []*data.Task
}

// StackingStrategy defines how tasks should be stacked based on priority category
type StackingStrategy struct {
	Category        PriorityCategory
	StackingOrder   []StackingOrder
	VisualProminence VisualProminence
	CollapseAllowed bool
	GroupingAllowed bool
}

// StackingOrder defines the order of stacking
type StackingOrder struct {
	Priority     int
	Condition    func(*data.Task) bool
	Position     int
	Grouping     bool
}

// TaskPrioritizationResult contains the result of task prioritization
type TaskPrioritizationResult struct {
	PrioritizedTasks    []*PrioritizedTask
	VisibilityAnalysis   *VisibilityAnalysis
	StackingOptimization *StackingOptimization
	Recommendations     []string
	AnalysisDate        time.Time
}

// PrioritizedTask represents a task with prioritization information
type PrioritizedTask struct {
	Task            *data.Task
	Priority        *TaskPriority
	Visibility      *VisibilityAction
	StackingOrder   int
	DisplayOrder    int
	VisualWeight    float64
	IsHidden        bool
	IsCollapsed     bool
	IsHighlighted   bool
	GroupingID      string
	ProminenceScore float64
}

// VisibilityAnalysis provides analysis of task visibility
type VisibilityAnalysis struct {
	TotalTasks        int
	VisibleTasks      int
	HiddenTasks       int
	CollapsedTasks    int
	HighlightedTasks  int
	ProminenceDistribution map[VisualProminence]int
	VisibilityScore   float64
	Recommendations   []string
}

// StackingOptimization provides analysis of stacking optimization
type StackingOptimization struct {
	OriginalOrder    []*data.Task
	OptimizedOrder   []*data.Task
	OptimizationGains map[string]float64
	VisualImprovement float64
	SpaceEfficiency   float64
	Recommendations  []string
}

// NewTaskPrioritizationEngine creates a new task prioritization engine
func NewTaskPrioritizationEngine(verticalStackingEngine *VerticalStackingEngine, priorityRanker *PriorityRanker) *TaskPrioritizationEngine {
	engine := &TaskPrioritizationEngine{
		verticalStackingEngine: verticalStackingEngine,
		priorityRanker:         priorityRanker,
		visibilityManager:      NewVisibilityManager(),
		stackingOptimizer:      NewStackingOptimizer(),
	}
	return engine
}

// NewVisibilityManager creates a new visibility manager
func NewVisibilityManager() *VisibilityManager {
	return &VisibilityManager{
		visibilityRules: make([]VisibilityRule, 0),
		prominenceWeights: map[VisualProminence]float64{
			ProminenceCritical: 1.0,
			ProminenceHigh:     0.8,
			ProminenceMedium:   0.6,
			ProminenceLow:      0.4,
			ProminenceMinimal:  0.2,
		},
		visibilityThreshold: 0.5,
		adaptiveVisibility:  true,
	}
}

// NewStackingOptimizer creates a new stacking optimizer
func NewStackingOptimizer() *StackingOptimizer {
	return &StackingOptimizer{
		optimizationRules: make([]OptimizationRule, 0),
		stackingStrategies: make(map[PriorityCategory]StackingStrategy),
		adaptiveOrdering:  true,
	}
}

// PrioritizeTasks performs intelligent task prioritization for stacking order
func (tpe *TaskPrioritizationEngine) PrioritizeTasks(tasks []*data.Task, context *PriorityContext) *TaskPrioritizationResult {
	// Step 1: Rank tasks by priority
	priorityRanking := tpe.priorityRanker.RankTasks(tasks, context)
	
	// Step 2: Determine visibility for each task
	visibilityAnalysis := tpe.visibilityManager.AnalyzeVisibility(tasks, context)
	
	// Step 3: Optimize stacking order
	stackingOptimization := tpe.stackingOptimizer.OptimizeStackingOrder(tasks, context)
	
	// Step 4: Create prioritized tasks
	prioritizedTasks := tpe.createPrioritizedTasks(tasks, priorityRanking, visibilityAnalysis, context)
	
	// Step 5: Apply stacking optimization
	prioritizedTasks = tpe.applyStackingOptimization(prioritizedTasks, stackingOptimization)
	
	// Step 6: Generate recommendations
	recommendations := tpe.generatePrioritizationRecommendations(prioritizedTasks, visibilityAnalysis, stackingOptimization)
	
	return &TaskPrioritizationResult{
		PrioritizedTasks:    prioritizedTasks,
		VisibilityAnalysis:   visibilityAnalysis,
		StackingOptimization: stackingOptimization,
		Recommendations:     recommendations,
		AnalysisDate:        time.Now(),
	}
}

// createPrioritizedTasks creates prioritized task objects
func (tpe *TaskPrioritizationEngine) createPrioritizedTasks(
	tasks []*data.Task,
	priorityRanking *PriorityRanking,
	visibilityAnalysis *VisibilityAnalysis,
	context *PriorityContext,
) []*PrioritizedTask {
	prioritizedTasks := make([]*PrioritizedTask, 0, len(tasks))
	
	for i, task := range tasks {
		// Get priority information
		var priority *TaskPriority
		for _, taskPriority := range priorityRanking.TaskPriorities {
			if taskPriority.Task.ID == task.ID {
				priority = taskPriority
				break
			}
		}
		
		// Get visibility information
		visibility := tpe.visibilityManager.DetermineVisibility(task, context)
		
		// Calculate visual weight
		visualWeight := tpe.calculateVisualWeight(task, priority, visibility)
		
		// Calculate prominence score
		prominenceScore := tpe.calculateProminenceScore(task, priority, visibility)
		
		// Determine grouping
		groupingID := tpe.determineGrouping(task, priority, context)
		
		prioritizedTask := &PrioritizedTask{
			Task:            task,
			Priority:        priority,
			Visibility:      visibility,
			StackingOrder:   i,
			DisplayOrder:    i,
			VisualWeight:    visualWeight,
			IsHidden:        !visibility.IsVisible,
			IsCollapsed:     visibility.CollapseLevel > 0,
			IsHighlighted:   visibility.HighlightLevel > 0,
			GroupingID:      groupingID,
			ProminenceScore: prominenceScore,
		}
		
		prioritizedTasks = append(prioritizedTasks, prioritizedTask)
	}
	
	// Sort by prominence score (highest first)
	sort.Slice(prioritizedTasks, func(i, j int) bool {
		return prioritizedTasks[i].ProminenceScore > prioritizedTasks[j].ProminenceScore
	})
	
	// Update display order
	for i, task := range prioritizedTasks {
		task.DisplayOrder = i
	}
	
	return prioritizedTasks
}

// calculateVisualWeight calculates the visual weight of a task
func (tpe *TaskPrioritizationEngine) calculateVisualWeight(task *data.Task, priority *TaskPriority, visibility *VisibilityAction) float64 {
	weight := 1.0
	
	// Base weight from priority
	if priority != nil {
		weight += priority.PriorityScore * 0.1
	}
	
	// Visibility weight
	if visibility != nil {
		weight += float64(visibility.VisualWeight)
	}
	
	// Task characteristics
	if task.IsMilestone {
		weight += 2.0
	}
	
	// Category weight
	switch task.Category {
	case "DISSERTATION":
		weight += 3.0
	case "PROPOSAL":
		weight += 2.0
	case "LASER":
		weight += 1.0
	}
	
	// Duration weight
	duration := task.EndDate.Sub(task.StartDate)
	weight += float64(duration.Hours()) * 0.01
	
	return weight
}

// calculateProminenceScore calculates the prominence score of a task
func (tpe *TaskPrioritizationEngine) calculateProminenceScore(task *data.Task, priority *TaskPriority, visibility *VisibilityAction) float64 {
	score := 0.0
	
	// Priority score
	if priority != nil {
		score += priority.PriorityScore
	}
	
	// Visual prominence score
	if priority != nil {
		switch priority.VisualProminence {
		case ProminenceCritical:
			score += 100.0
		case ProminenceHigh:
			score += 80.0
		case ProminenceMedium:
			score += 60.0
		case ProminenceLow:
			score += 40.0
		case ProminenceMinimal:
			score += 20.0
		}
	}
	
	// Visibility score
	if visibility != nil {
		score += float64(visibility.VisualWeight) * 10.0
	}
	
	// Task characteristics
	if task.IsMilestone {
		score += 50.0
	}
	
	// Category score
	switch task.Category {
	case "DISSERTATION":
		score += 30.0
	case "PROPOSAL":
		score += 20.0
	case "LASER":
		score += 10.0
	}
	
	// Urgency score (based on timeline urgency)
	if priority != nil && priority.TimelineUrgency > 0 {
		score += priority.TimelineUrgency * 5.0
	}
	
	return score
}

// determineGrouping determines the grouping ID for a task
func (tpe *TaskPrioritizationEngine) determineGrouping(task *data.Task, priority *TaskPriority, context *PriorityContext) string {
	// Group by category
	if task.Category != "" {
		return fmt.Sprintf("category_%s", task.Category)
	}
	
	// Group by priority
	if priority != nil {
		return fmt.Sprintf("priority_%s", priority.VisualProminence)
	}
	
	// Group by assignee
	if task.Assignee != "" {
		return fmt.Sprintf("assignee_%s", task.Assignee)
	}
	
	return "default"
}

// applyStackingOptimization applies stacking optimization to prioritized tasks
func (tpe *TaskPrioritizationEngine) applyStackingOptimization(
	prioritizedTasks []*PrioritizedTask,
	stackingOptimization *StackingOptimization,
) []*PrioritizedTask {
	// Create a map of task ID to prioritized task
	taskMap := make(map[string]*PrioritizedTask)
	for _, task := range prioritizedTasks {
		taskMap[task.Task.ID] = task
	}
	
	// Apply optimized order
	optimizedTasks := make([]*PrioritizedTask, 0, len(prioritizedTasks))
	for _, task := range stackingOptimization.OptimizedOrder {
		if prioritizedTask, exists := taskMap[task.ID]; exists {
			optimizedTasks = append(optimizedTasks, prioritizedTask)
		}
	}
	
	// Update stacking order
	for i, task := range optimizedTasks {
		task.StackingOrder = i
	}
	
	return optimizedTasks
}

// generatePrioritizationRecommendations generates recommendations for task prioritization
func (tpe *TaskPrioritizationEngine) generatePrioritizationRecommendations(
	prioritizedTasks []*PrioritizedTask,
	visibilityAnalysis *VisibilityAnalysis,
	stackingOptimization *StackingOptimization,
) []string {
	var recommendations []string
	
	// Visibility recommendations
	if visibilityAnalysis.HiddenTasks > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("⚠️ %d tasks are hidden - consider increasing available space or adjusting visibility rules", 
				visibilityAnalysis.HiddenTasks))
	}
	
	if visibilityAnalysis.CollapsedTasks > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("📦 %d tasks are collapsed - consider expanding important tasks", 
				visibilityAnalysis.CollapsedTasks))
	}
	
	// Prominence recommendations
	criticalCount := visibilityAnalysis.ProminenceDistribution[ProminenceCritical]
	highCount := visibilityAnalysis.ProminenceDistribution[ProminenceHigh]
	
	if criticalCount > 5 {
		recommendations = append(recommendations, 
			"🔥 Too many critical tasks - consider prioritizing or breaking down large tasks")
	}
	
	if highCount > 10 {
		recommendations = append(recommendations, 
			"⭐ Many high-priority tasks - consider using task grouping or time management")
	}
	
	// Stacking optimization recommendations
	if stackingOptimization.VisualImprovement > 0.1 {
		recommendations = append(recommendations, 
			fmt.Sprintf("📈 Visual improvement of %.1f%% achieved through stacking optimization", 
				stackingOptimization.VisualImprovement*100))
	}
	
	if stackingOptimization.SpaceEfficiency < 0.7 {
		recommendations = append(recommendations, 
			"📏 Low space efficiency - consider adjusting task sizes or using compression")
	}
	
	// Task grouping recommendations
	groupingCounts := make(map[string]int)
	for _, task := range prioritizedTasks {
		groupingCounts[task.GroupingID]++
	}
	
	for groupID, count := range groupingCounts {
		if count > 5 {
			recommendations = append(recommendations, 
				fmt.Sprintf("📚 Group '%s' has %d tasks - consider sub-grouping or prioritization", 
					groupID, count))
		}
	}
	
	return recommendations
}

// AnalyzeVisibility analyzes task visibility
func (vm *VisibilityManager) AnalyzeVisibility(tasks []*data.Task, context *PriorityContext) *VisibilityAnalysis {
	analysis := &VisibilityAnalysis{
		TotalTasks:            len(tasks),
		VisibleTasks:          0,
		HiddenTasks:           0,
		CollapsedTasks:        0,
		HighlightedTasks:      0,
		ProminenceDistribution: make(map[VisualProminence]int),
		VisibilityScore:       0.0,
		Recommendations:       make([]string, 0),
	}
	
	// Initialize prominence distribution
	for _, prominence := range []VisualProminence{
		ProminenceCritical, ProminenceHigh, ProminenceMedium, 
		ProminenceLow, ProminenceMinimal,
	} {
		analysis.ProminenceDistribution[prominence] = 0
	}
	
	// Analyze each task
	for _, task := range tasks {
		visibility := vm.DetermineVisibility(task, context)
		
		if visibility.IsVisible {
			analysis.VisibleTasks++
		} else {
			analysis.HiddenTasks++
		}
		
		if visibility.CollapseLevel > 0 {
			analysis.CollapsedTasks++
		}
		
		if visibility.HighlightLevel > 0 {
			analysis.HighlightedTasks++
		}
		
		// Count prominence distribution
		analysis.ProminenceDistribution[visibility.ProminenceLevel]++
	}
	
	// Calculate visibility score
	if analysis.TotalTasks > 0 {
		analysis.VisibilityScore = float64(analysis.VisibleTasks) / float64(analysis.TotalTasks)
	}
	
	// Generate visibility recommendations
	analysis.Recommendations = vm.generateVisibilityRecommendations(analysis)
	
	return analysis
}

// DetermineVisibility determines the visibility of a task
func (vm *VisibilityManager) DetermineVisibility(task *data.Task, context *PriorityContext) *VisibilityAction {
	// Start with default visibility
	action := &VisibilityAction{
		IsVisible:       true,
		ProminenceLevel: ProminenceMedium,
		DisplayOrder:    0,
		VisualWeight:    1.0,
		CollapseLevel:   0,
		HighlightLevel:  0,
	}
	
	// Apply visibility rules
	for _, rule := range vm.visibilityRules {
		if rule.Condition(task, context) {
			ruleAction := rule.Action(task, context)
			if ruleAction != nil {
				action = ruleAction
				break
			}
		}
	}
	
	// Apply adaptive visibility
	if vm.adaptiveVisibility {
		action = vm.applyAdaptiveVisibility(action, task, context)
	}
	
	return action
}

// applyAdaptiveVisibility applies adaptive visibility based on context
func (vm *VisibilityManager) applyAdaptiveVisibility(action *VisibilityAction, task *data.Task, context *PriorityContext) *VisibilityAction {
	// Check if task should be hidden due to low priority
	if action.ProminenceLevel == ProminenceMinimal && vm.visibilityThreshold > 0.5 {
		action.IsVisible = false
		action.CollapseLevel = 1
	}
	
	// Check if task should be highlighted due to high priority
	if action.ProminenceLevel == ProminenceCritical {
		action.HighlightLevel = 2
		action.VisualWeight = 2.0
	}
	
	return action
}

// generateVisibilityRecommendations generates visibility recommendations
func (vm *VisibilityManager) generateVisibilityRecommendations(analysis *VisibilityAnalysis) []string {
	var recommendations []string
	
	// Hidden tasks recommendations
	if analysis.HiddenTasks > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("👁️ %d tasks are hidden - consider adjusting visibility threshold", 
				analysis.HiddenTasks))
	}
	
	// Collapsed tasks recommendations
	if analysis.CollapsedTasks > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("📦 %d tasks are collapsed - consider expanding important tasks", 
				analysis.CollapsedTasks))
	}
	
	// Prominence distribution recommendations
	criticalCount := analysis.ProminenceDistribution[ProminenceCritical]
	if criticalCount > 3 {
		recommendations = append(recommendations, 
			"🔥 Too many critical tasks - consider prioritizing or breaking down tasks")
	}
	
	// Visibility score recommendations
	if analysis.VisibilityScore < 0.8 {
		recommendations = append(recommendations, 
			"👁️ Low visibility score - consider adjusting visibility rules")
	}
	
	return recommendations
}

// OptimizeStackingOrder optimizes the stacking order of tasks
func (so *StackingOptimizer) OptimizeStackingOrder(tasks []*data.Task, context *PriorityContext) *StackingOptimization {
	optimization := &StackingOptimization{
		OriginalOrder:    make([]*data.Task, len(tasks)),
		OptimizedOrder:   make([]*data.Task, len(tasks)),
		OptimizationGains: make(map[string]float64),
		VisualImprovement: 0.0,
		SpaceEfficiency:   0.0,
		Recommendations:   make([]string, 0),
	}
	
	// Copy original order
	copy(optimization.OriginalOrder, tasks)
	
	// Start with original order
	optimizedTasks := make([]*data.Task, len(tasks))
	copy(optimizedTasks, tasks)
	
	// Apply optimization rules
	for _, rule := range so.optimizationRules {
		if rule.Condition(optimizedTasks, context) {
			optimizedTasks = rule.Action(optimizedTasks, context)
		}
	}
	
	// Apply adaptive ordering
	if so.adaptiveOrdering {
		optimizedTasks = so.applyAdaptiveOrdering(optimizedTasks, context)
	}
	
	// Set optimized order
	optimization.OptimizedOrder = optimizedTasks
	
	// Calculate optimization gains
	optimization.VisualImprovement = so.calculateVisualImprovement(tasks, optimizedTasks)
	optimization.SpaceEfficiency = so.calculateSpaceEfficiency(optimizedTasks)
	
	// Generate recommendations
	optimization.Recommendations = so.generateOptimizationRecommendations(optimization)
	
	return optimization
}

// applyAdaptiveOrdering applies adaptive ordering to tasks
func (so *StackingOptimizer) applyAdaptiveOrdering(tasks []*data.Task, context *PriorityContext) []*data.Task {
	// Sort by priority score (highest first)
	sort.Slice(tasks, func(i, j int) bool {
		// Get priority scores
		scoreI := so.getPriorityScore(tasks[i], context)
		scoreJ := so.getPriorityScore(tasks[j], context)
		return scoreI > scoreJ
	})
	
	return tasks
}

// getPriorityScore gets the priority score of a task
func (so *StackingOptimizer) getPriorityScore(task *data.Task, context *PriorityContext) float64 {
	// This would typically use the priority ranker
	// For now, use a simple scoring system
	score := 0.0
	
	// Base score from task priority
	score += float64(task.Priority)
	
	// Milestone bonus
	if task.IsMilestone {
		score += 10.0
	}
	
	// Category bonus
	switch task.Category {
	case "DISSERTATION":
		score += 5.0
	case "PROPOSAL":
		score += 3.0
	case "LASER":
		score += 1.0
	}
	
	return score
}

// calculateVisualImprovement calculates the visual improvement from optimization
func (so *StackingOptimizer) calculateVisualImprovement(original, optimized []*data.Task) float64 {
	// Simple calculation based on task order changes
	// In a real implementation, this would be more sophisticated
	changes := 0
	for i, task := range original {
		if i < len(optimized) && task.ID != optimized[i].ID {
			changes++
		}
	}
	
	if len(original) == 0 {
		return 0.0
	}
	
	return float64(changes) / float64(len(original))
}

// calculateSpaceEfficiency calculates the space efficiency of the optimized order
func (so *StackingOptimizer) calculateSpaceEfficiency(tasks []*data.Task) float64 {
	// Simple calculation based on task characteristics
	// In a real implementation, this would consider actual space usage
	efficiency := 1.0
	
	// Penalize for too many high-priority tasks
	highPriorityCount := 0
	for _, task := range tasks {
		if task.Priority >= 4 {
			highPriorityCount++
		}
	}
	
	if highPriorityCount > 5 {
		efficiency -= 0.2
	}
	
	return math.Max(efficiency, 0.0)
}

// generateOptimizationRecommendations generates optimization recommendations
func (so *StackingOptimizer) generateOptimizationRecommendations(optimization *StackingOptimization) []string {
	var recommendations []string
	
	// Visual improvement recommendations
	if optimization.VisualImprovement > 0.1 {
		recommendations = append(recommendations, 
			fmt.Sprintf("📈 Visual improvement of %.1f%% achieved", 
				optimization.VisualImprovement*100))
	}
	
	// Space efficiency recommendations
	if optimization.SpaceEfficiency < 0.8 {
		recommendations = append(recommendations, 
			"📏 Low space efficiency - consider task prioritization")
	}
	
	// Order change recommendations
	changes := 0
	for i, task := range optimization.OriginalOrder {
		if i < len(optimization.OptimizedOrder) && task.ID != optimization.OptimizedOrder[i].ID {
			changes++
		}
	}
	
	if changes > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🔄 %d tasks reordered for better visibility", changes))
	}
	
	return recommendations
}

// GetTasksByProminence returns tasks filtered by prominence level
func (result *TaskPrioritizationResult) GetTasksByProminence(prominence VisualProminence) []*PrioritizedTask {
	var filtered []*PrioritizedTask
	for _, task := range result.PrioritizedTasks {
		if task.Visibility != nil && task.Visibility.ProminenceLevel == prominence {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// GetVisibleTasks returns only visible tasks
func (result *TaskPrioritizationResult) GetVisibleTasks() []*PrioritizedTask {
	var visible []*PrioritizedTask
	for _, task := range result.PrioritizedTasks {
		if !task.IsHidden {
			visible = append(visible, task)
		}
	}
	return visible
}

// GetTasksByGroup returns tasks filtered by grouping ID
func (result *TaskPrioritizationResult) GetTasksByGroup(groupID string) []*PrioritizedTask {
	var filtered []*PrioritizedTask
	for _, task := range result.PrioritizedTasks {
		if task.GroupingID == groupID {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// GetSummary returns a summary of the task prioritization result
func (result *TaskPrioritizationResult) GetSummary() string {
	return fmt.Sprintf("Task Prioritization Summary:\n"+
		"  Total Tasks: %d\n"+
		"  Visible Tasks: %d\n"+
		"  Hidden Tasks: %d\n"+
		"  Collapsed Tasks: %d\n"+
		"  Highlighted Tasks: %d\n"+
		"  Visibility Score: %.2f%%\n"+
		"  Visual Improvement: %.2f%%\n"+
		"  Space Efficiency: %.2f%%\n"+
		"  Analysis Date: %s",
		len(result.PrioritizedTasks),
		result.VisibilityAnalysis.VisibleTasks,
		result.VisibilityAnalysis.HiddenTasks,
		result.VisibilityAnalysis.CollapsedTasks,
		result.VisibilityAnalysis.HighlightedTasks,
		result.VisibilityAnalysis.VisibilityScore*100,
		result.StackingOptimization.VisualImprovement*100,
		result.StackingOptimization.SpaceEfficiency*100,
		result.AnalysisDate.Format("2006-01-02 15:04:05"))
}
