package calendar

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"phd-dissertation-planner/internal/data"
)

// CalendarGridIntegration integrates smart stacking layout engine with calendar grid system
type CalendarGridIntegration struct {
	smartStackingEngine    *SmartStackingEngine
	verticalStackingEngine *VerticalStackingEngine
	taskPrioritizationEngine *TaskPrioritizationEngine
	conflictResolutionEngine *ConflictResolutionEngine
	multiDayLayoutEngine   *MultiDayLayoutEngine
	positioningEngine      *PositioningEngine
	monthBoundaryEngine    *MonthBoundaryEngine
	gridConfig            *GridConfig
	visualSettings        *IntegratedVisualSettings
	dateValidator         *data.DateValidator
}

// GridConfig defines the configuration for the calendar grid
type GridConfig struct {
	CalendarStart      time.Time
	CalendarEnd        time.Time
	DayWidth           float64
	DayHeight          float64
	RowHeight          float64
	MaxRowsPerDay      int
	OverlapThreshold   float64
	MonthBoundaryGap   float64
	TaskSpacing        float64
	VisualConstraints  *VisualConstraints
}

// IntegratedVisualSettings defines visual settings for the integrated system
type IntegratedVisualSettings struct {
	ShowTaskNames        bool
	ShowTaskDurations    bool
	ShowTaskPriorities   bool
	ShowConflictIndicators bool
	CollapseThreshold    int
	AnimationEnabled     bool
	HighlightConflicts   bool
	ColorScheme          string
	FontSize             string
	TaskBarOpacity       float64
	BorderWidth          float64
}

// IntegratedTaskBar represents a task bar with integrated smart stacking
type IntegratedTaskBar struct {
	TaskID           string
	StartDate        time.Time
	EndDate          time.Time
	StartX           float64
	EndX             float64
	Y                float64
	Width            float64
	Height           float64
	Row              int
	StackIndex       int
	Color            string
	BorderColor      string
	Opacity          float64
	ZIndex           int
	IsContinuation   bool
	IsStart          bool
	IsEnd            bool
	MonthBoundary    bool
	StackingType     StackingType
	VisualWeight     float64
	ProminenceScore  float64
	IsCollapsed      bool
	IsVisible        bool
	CollisionLevel   int
	OverflowLevel    int
	Priority         int
	Category         string
	TaskName         string
	Description      string
}

// IntegratedLayoutResult contains the result of integrated layout operations
type IntegratedLayoutResult struct {
	TaskBars           []*IntegratedTaskBar
	Stacks             []*TaskStack
	Conflicts          []*ResolvedConflict
	OverflowResolutions []*OverflowResolution
	VisualOptimizations []*VisualOptimization
	LayoutAdjustments  []*LayoutAdjustment
	Statistics         *IntegratedLayoutStatistics
	Recommendations    []string
	AnalysisDate       time.Time
}

// IntegratedLayoutStatistics contains statistics about the integrated layout
type IntegratedLayoutStatistics struct {
	TotalTasks           int
	ProcessedBars        int
	TotalStacks          int
	ConflictsResolved    int
	OverflowResolutions  int
	VisualOptimizations  int
	LayoutAdjustments    int
	CollisionCount       int
	OverflowCount        int
	MonthBoundaryCount   int
	SpaceEfficiency      float64
	VisualQuality        float64
	AverageStackHeight   float64
	MaxStackHeight       float64
	AverageTaskHeight    float64
	AverageTaskWidth     float64
	AlignmentScore       float64
	SpacingScore         float64
	VisualBalance        float64
	GridUtilization      float64
}

// NewCalendarGridIntegration creates a new calendar grid integration instance
func NewCalendarGridIntegration(config *GridConfig) *CalendarGridIntegration {
	// Create smart stacking engine
	overlapDetector := NewOverlapDetector(config.CalendarStart, config.CalendarEnd)
	conflictCategorizer := NewConflictCategorizer(overlapDetector)
	priorityRanker := NewPriorityRanker(conflictCategorizer)
	smartStackingEngine := NewSmartStackingEngine(overlapDetector, conflictCategorizer, priorityRanker)
	
	// Create vertical stacking engine
	verticalStackingEngine := NewVerticalStackingEngine(smartStackingEngine)
	
	// Create task prioritization engine
	taskPrioritizationEngine := NewTaskPrioritizationEngine(verticalStackingEngine, priorityRanker)
	
	// Create conflict resolution engine
	conflictResolutionEngine := NewConflictResolutionEngine(taskPrioritizationEngine, verticalStackingEngine, smartStackingEngine)
	
	// Create multi-day layout engine
	multiDayLayoutEngine := NewMultiDayLayoutEngine(
		config.CalendarStart,
		config.CalendarEnd,
		config.DayWidth,
		config.DayHeight,
	)
	
	// Create positioning engine
	positioningEngine := NewPositioningEngine(config)
	
	// Create month boundary engine
	monthBoundaryEngine := NewMonthBoundaryEngine(config)
	
	// Create date validator
	dateValidator := data.NewDateValidator()
	
	// Set visual constraints
	if config.VisualConstraints == nil {
		config.VisualConstraints = &VisualConstraints{
			MaxStackHeight:     config.DayHeight * float64(config.MaxRowsPerDay),
			MinTaskHeight:      config.RowHeight * 0.5,
			MaxTaskHeight:      config.RowHeight * 2.0,
			MinTaskWidth:       config.DayWidth * 0.1,
			MaxTaskWidth:       config.DayWidth * 7.0, // Max 7 days
			VerticalSpacing:    config.TaskSpacing,
			HorizontalSpacing:  config.TaskSpacing,
			MaxStackDepth:      config.MaxRowsPerDay,
			CollisionThreshold: config.OverlapThreshold,
			OverflowThreshold:  0.8,
		}
	}
	
	return &CalendarGridIntegration{
		smartStackingEngine:      smartStackingEngine,
		verticalStackingEngine:   verticalStackingEngine,
		taskPrioritizationEngine: taskPrioritizationEngine,
		conflictResolutionEngine: conflictResolutionEngine,
		multiDayLayoutEngine:     multiDayLayoutEngine,
		positioningEngine:        positioningEngine,
		monthBoundaryEngine:      monthBoundaryEngine,
		gridConfig:              config,
		visualSettings: &IntegratedVisualSettings{
			ShowTaskNames:         true,
			ShowTaskDurations:     true,
			ShowTaskPriorities:    true,
			ShowConflictIndicators: true,
			CollapseThreshold:     5,
			AnimationEnabled:      false,
			HighlightConflicts:    true,
			ColorScheme:           "default",
			FontSize:              "small",
			TaskBarOpacity:        1.0,
			BorderWidth:           0.5,
		},
		dateValidator:            dateValidator,
	}
}

// ProcessTasksWithSmartStacking processes tasks with integrated smart stacking
func (cgi *CalendarGridIntegration) ProcessTasksWithSmartStacking(tasks []*data.Task) (*IntegratedLayoutResult, error) {
	// Step 1: Detect overlaps and conflicts
	overlapAnalysis := cgi.smartStackingEngine.overlapDetector.DetectOverlaps(tasks)
	conflictAnalysis := cgi.smartStackingEngine.conflictCategorizer.CategorizeConflicts(overlapAnalysis)
	
	// Step 2: Prioritize tasks
	priorityContext := &PriorityContext{
		CurrentTime:      time.Now(),
		CalendarStart:    cgi.gridConfig.CalendarStart,
		CalendarEnd:      cgi.gridConfig.CalendarEnd,
		OverlapAnalysis:  overlapAnalysis,
		ConflictAnalysis: conflictAnalysis,
	}
	
	prioritizationResult := cgi.taskPrioritizationEngine.PrioritizeTasks(tasks, priorityContext)
	
	// Step 3: Create stacking context
	stackingContext := &StackingContext{
		CalendarStart:     cgi.gridConfig.CalendarStart,
		CalendarEnd:       cgi.gridConfig.CalendarEnd,
		CurrentTime:       time.Now(),
		DayWidth:          cgi.gridConfig.DayWidth,
		DayHeight:         cgi.gridConfig.DayHeight,
		AvailableHeight:   cgi.gridConfig.DayHeight * float64(cgi.gridConfig.MaxRowsPerDay),
		AvailableWidth:    cgi.gridConfig.DayWidth * 7.0, // Max 7 days
		ExistingStacks:    []*TaskStack{},
		TaskPriorities:    make(map[string]*TaskPriority),
		ConflictAnalysis:  conflictAnalysis,
		OverlapAnalysis:   overlapAnalysis,
		VisualSettings:    &VisualSettings{
			ShowTaskNames:         cgi.visualSettings.ShowTaskNames,
			ShowTaskDurations:     cgi.visualSettings.ShowTaskDurations,
			ShowTaskPriorities:    cgi.visualSettings.ShowTaskPriorities,
			ShowConflictIndicators: cgi.visualSettings.ShowConflictIndicators,
			CollapseThreshold:     cgi.visualSettings.CollapseThreshold,
			AnimationEnabled:      cgi.visualSettings.AnimationEnabled,
			HighlightConflicts:    cgi.visualSettings.HighlightConflicts,
			ColorScheme:           cgi.visualSettings.ColorScheme,
		},
		VisualConstraints: cgi.gridConfig.VisualConstraints,
	}
	
	// Step 4: Apply smart stacking
	stackingResult := cgi.smartStackingEngine.StackTasks(tasks, stackingContext)
	
	// Step 5: Apply vertical stacking
	verticalStackingResult := cgi.verticalStackingEngine.StackTasksVertically(tasks, stackingContext)
	
	// Step 6: Resolve conflicts
	conflictResolutionResult := cgi.conflictResolutionEngine.ResolveConflicts(tasks, priorityContext)
	
	// Step 7: Create integrated task bars
	integratedBars := cgi.createIntegratedTaskBars(tasks, stackingResult, verticalStackingResult, conflictResolutionResult, prioritizationResult)
	
	// Step 8: Apply precise positioning
	positioningResult, err := cgi.positioningEngine.PositionTasks(tasks, integratedBars)
	if err != nil {
		return nil, fmt.Errorf("failed to position tasks: %v", err)
	}
	
	// Step 9: Handle month boundaries with dedicated engine
	monthBoundaryResult, err := cgi.monthBoundaryEngine.ProcessMonthBoundaries(positioningResult.TaskBars, time.Now().Month(), time.Now().Year())
	if err != nil {
		return nil, fmt.Errorf("failed to process month boundaries: %v", err)
	}
	
	processedBars := monthBoundaryResult.ProcessedBars
	
	// Step 10: Calculate statistics
	statistics := cgi.calculateIntegratedStatistics(processedBars, stackingResult, conflictResolutionResult)
	
	// Merge positioning metrics
	statistics.AlignmentScore = positioningResult.Metrics.AlignmentScore
	statistics.SpacingScore = positioningResult.Metrics.SpacingScore
	statistics.VisualBalance = positioningResult.Metrics.VisualBalance
	statistics.GridUtilization = positioningResult.Metrics.GridUtilization
	
	// Merge month boundary metrics
	statistics.MonthBoundaryCount = monthBoundaryResult.BoundaryMetrics.ContinuationsCreated
	
	// Step 11: Generate recommendations
	recommendations := cgi.generateRecommendations(statistics, conflictResolutionResult)
	recommendations = append(recommendations, positioningResult.Recommendations...)
	recommendations = append(recommendations, monthBoundaryResult.Recommendations...)
	
	return &IntegratedLayoutResult{
		TaskBars:            processedBars,
		Stacks:              stackingResult.Stacks,
		Conflicts:           conflictResolutionResult.ResolvedConflicts,
		OverflowResolutions: conflictResolutionResult.OverflowResolutions,
		VisualOptimizations: conflictResolutionResult.VisualOptimizations,
		LayoutAdjustments:   conflictResolutionResult.LayoutAdjustments,
		Statistics:          statistics,
		Recommendations:     recommendations,
		AnalysisDate:        time.Now(),
	}, nil
}

// createIntegratedTaskBars creates integrated task bars with smart stacking
func (cgi *CalendarGridIntegration) createIntegratedTaskBars(
	tasks []*data.Task,
	stackingResult *StackingResult,
	verticalStackingResult *VerticalStackingResult,
	conflictResolutionResult *ConflictResolutionResult,
	prioritizationResult *TaskPrioritizationResult,
) []*IntegratedTaskBar {
	var integratedBars []*IntegratedTaskBar
	
	// Create a map of task priorities for quick lookup
	priorityMap := make(map[string]*TaskPriority)
	for _, prioritizedTask := range prioritizationResult.PrioritizedTasks {
		priorityMap[prioritizedTask.Task.ID] = prioritizedTask.Priority
	}
	
	// Process each task
	for _, task := range tasks {
		// Calculate basic positioning
		startX := cgi.calculateXPosition(task.StartDate)
		endX := cgi.calculateXPosition(task.EndDate)
		width := endX - startX
		
		// Get task priority
		priority := priorityMap[task.ID]
		if priority == nil {
			priority = &TaskPriority{
				TimelineUrgency:   0.5,
				ResourceContention: 0.5,
				MilestoneWeight:   0.5,
				PriorityScore:     0.5,
			}
		}
		
		// Calculate visual weight and prominence
		visualWeight := cgi.calculateVisualWeight(task, priority)
		prominenceScore := cgi.calculateProminenceScore(task, priority, visualWeight)
		
		// Determine stacking type and position
		stackingType, stackIndex, y, height := cgi.determineStackingPosition(
			task, stackingResult, verticalStackingResult, visualWeight,
		)
		
		// Get task category and color
		category := data.GetCategory(task.Category)
		
		// Create integrated task bar
		integratedBar := &IntegratedTaskBar{
			TaskID:          task.ID,
			StartDate:       task.StartDate,
			EndDate:         task.EndDate,
			StartX:          startX,
			EndX:            endX,
			Y:               y,
			Width:           width,
			Height:          height,
			Row:             0, // Will be calculated based on stacking
			StackIndex:      stackIndex,
			Color:           category.Color,
			BorderColor:     "#000000",
			Opacity:         cgi.visualSettings.TaskBarOpacity,
			ZIndex:          int(priority.PriorityScore * 5),
			IsContinuation:  cgi.isTaskContinuation(task),
			IsStart:         cgi.isTaskStart(task),
			IsEnd:           cgi.isTaskEnd(task),
			MonthBoundary:   cgi.hasMonthBoundary(task),
			StackingType:    stackingType,
			VisualWeight:    visualWeight,
			ProminenceScore: prominenceScore,
			IsCollapsed:     false,
			IsVisible:       true,
			CollisionLevel:  0,
			OverflowLevel:   0,
			Priority:        int(priority.PriorityScore * 5),
			Category:        task.Category,
			TaskName:        task.Name,
			Description:     task.Description,
		}
		
		integratedBars = append(integratedBars, integratedBar)
	}
	
	return integratedBars
}

// calculateXPosition calculates the X position for a given date
func (cgi *CalendarGridIntegration) calculateXPosition(date time.Time) float64 {
	daysFromStart := int(date.Sub(cgi.gridConfig.CalendarStart).Hours() / 24)
	return float64(daysFromStart) * cgi.gridConfig.DayWidth
}

// calculateVisualWeight calculates the visual weight of a task
func (cgi *CalendarGridIntegration) calculateVisualWeight(task *data.Task, priority *TaskPriority) float64 {
	// Base weight from priority
	weight := priority.PriorityScore
	
	// Adjust based on task duration
	duration := task.EndDate.Sub(task.StartDate).Hours() / 24
	if duration > 7 {
		weight *= 1.2 // Longer tasks get more visual weight
	} else if duration < 1 {
		weight *= 0.8 // Shorter tasks get less visual weight
	}
	
	// Adjust based on category
	category := data.GetCategory(task.Category)
	weight *= float64(category.Priority) / 5.0
	
	// Adjust based on milestone status
	if strings.Contains(strings.ToUpper(task.Name), "MILESTONE") {
		weight *= 1.5
	}
	
	return math.Min(weight, 1.0)
}

// calculateProminenceScore calculates the prominence score of a task
func (cgi *CalendarGridIntegration) calculateProminenceScore(task *data.Task, priority *TaskPriority, visualWeight float64) float64 {
	// Base prominence from visual weight
	prominence := visualWeight
	
	// Adjust based on priority (using priority score as proxy)
	prominence *= priority.PriorityScore
	
	// Adjust based on timeline urgency
	prominence *= priority.TimelineUrgency
	
	// Adjust based on milestone priority
	prominence *= priority.MilestoneWeight
	
	return math.Min(prominence, 1.0)
}

// determineStackingPosition determines the stacking position for a task
func (cgi *CalendarGridIntegration) determineStackingPosition(
	task *data.Task,
	stackingResult *StackingResult,
	verticalStackingResult *VerticalStackingResult,
	visualWeight float64,
) (StackingType, int, float64, float64) {
	// Find the appropriate stack for this task
	for _, stack := range stackingResult.Stacks {
		for _, stackedTask := range stack.Tasks {
			if stackedTask.Task.ID == task.ID {
				// Calculate Y position based on stack and position within stack
				y := cgi.calculateYPosition(stackedTask.Position.Y, stack.TotalHeight)
				height := cgi.calculateTaskHeight(task, visualWeight)
				
				return stack.StackingType, stackedTask.Position.ZIndex, y, height
			}
		}
	}
	
	// Default positioning if not found in stacks
	y := cgi.gridConfig.DayHeight * 0.1 // 10% from top
	height := cgi.calculateTaskHeight(task, visualWeight)
	
	return StackingTypeVertical, 0, y, height
}

// calculateYPosition calculates the Y position within a stack
func (cgi *CalendarGridIntegration) calculateYPosition(relativeY, stackHeight float64) float64 {
	// Normalize relative Y position to actual Y position
	return relativeY * (cgi.gridConfig.DayHeight / stackHeight)
}

// calculateTaskHeight calculates the height of a task based on its properties
func (cgi *CalendarGridIntegration) calculateTaskHeight(task *data.Task, visualWeight float64) float64 {
	// Base height
	height := cgi.gridConfig.RowHeight
	
	// Adjust based on visual weight
	height *= visualWeight
	
	// Ensure within constraints
	minHeight := cgi.gridConfig.VisualConstraints.MinTaskHeight
	maxHeight := cgi.gridConfig.VisualConstraints.MaxTaskHeight
	
	if height < minHeight {
		height = minHeight
	} else if height > maxHeight {
		height = maxHeight
	}
	
	return height
}

// isTaskContinuation checks if this task is a continuation from previous month
func (cgi *CalendarGridIntegration) isTaskContinuation(task *data.Task) bool {
	return task.StartDate.Before(cgi.gridConfig.CalendarStart)
}

// isTaskStart checks if this task starts on this day
func (cgi *CalendarGridIntegration) isTaskStart(task *data.Task) bool {
	return task.StartDate.Equal(cgi.gridConfig.CalendarStart) || 
		   task.StartDate.After(cgi.gridConfig.CalendarStart)
}

// isTaskEnd checks if this task ends on this day
func (cgi *CalendarGridIntegration) isTaskEnd(task *data.Task) bool {
	return task.EndDate.Equal(cgi.gridConfig.CalendarEnd) || 
		   task.EndDate.Before(cgi.gridConfig.CalendarEnd)
}

// hasMonthBoundary checks if this task crosses month boundaries
func (cgi *CalendarGridIntegration) hasMonthBoundary(task *data.Task) bool {
	startMonth := task.StartDate.Month()
	endMonth := task.EndDate.Month()
	return startMonth != endMonth
}

// handleMonthBoundaries handles month boundary transitions
func (cgi *CalendarGridIntegration) handleMonthBoundaries(bars []*IntegratedTaskBar) []*IntegratedTaskBar {
	// Sort bars by start date
	sort.Slice(bars, func(i, j int) bool {
		return bars[i].StartDate.Before(bars[j].StartDate)
	})
	
	// Process month boundaries
	for _, bar := range bars {
		if bar.MonthBoundary {
			// Add month boundary gap
			bar.StartX += cgi.gridConfig.MonthBoundaryGap
			bar.Width -= cgi.gridConfig.MonthBoundaryGap
		}
	}
	
	return bars
}

// calculateIntegratedStatistics calculates statistics for the integrated layout
func (cgi *CalendarGridIntegration) calculateIntegratedStatistics(
	bars []*IntegratedTaskBar,
	stackingResult *StackingResult,
	conflictResolutionResult *ConflictResolutionResult,
) *IntegratedLayoutStatistics {
	stats := &IntegratedLayoutStatistics{
		TotalTasks:          len(bars),
		ProcessedBars:       len(bars),
		TotalStacks:         len(stackingResult.Stacks),
		ConflictsResolved:   len(conflictResolutionResult.ResolvedConflicts),
		OverflowResolutions: len(conflictResolutionResult.OverflowResolutions),
		VisualOptimizations: len(conflictResolutionResult.VisualOptimizations),
		LayoutAdjustments:   len(conflictResolutionResult.LayoutAdjustments),
		SpaceEfficiency:     stackingResult.SpaceEfficiency,
		VisualQuality:       stackingResult.VisualQuality,
	}
	
	// Calculate additional statistics
	var totalHeight, maxHeight, totalWidth float64
	monthBoundaryCount := 0
	
	for _, bar := range bars {
		totalHeight += bar.Height
		totalWidth += bar.Width
		
		if bar.Height > maxHeight {
			maxHeight = bar.Height
		}
		
		if bar.MonthBoundary {
			monthBoundaryCount++
		}
	}
	
	stats.MonthBoundaryCount = monthBoundaryCount
	stats.MaxStackHeight = maxHeight
	
	if len(bars) > 0 {
		stats.AverageTaskHeight = totalHeight / float64(len(bars))
		stats.AverageTaskWidth = totalWidth / float64(len(bars))
	}
	
	// Calculate average stack height
	if len(stackingResult.Stacks) > 0 {
		var totalStackHeight float64
		for _, stack := range stackingResult.Stacks {
			totalStackHeight += stack.TotalHeight
		}
		stats.AverageStackHeight = totalStackHeight / float64(len(stackingResult.Stacks))
	}
	
	return stats
}

// generateRecommendations generates recommendations based on the layout analysis
func (cgi *CalendarGridIntegration) generateRecommendations(
	statistics *IntegratedLayoutStatistics,
	conflictResolutionResult *ConflictResolutionResult,
) []string {
	var recommendations []string
	
	// Space efficiency recommendations
	if statistics.SpaceEfficiency < 0.7 {
		recommendations = append(recommendations, "Consider reducing task spacing to improve space efficiency")
	}
	
	// Visual quality recommendations
	if statistics.VisualQuality < 0.8 {
		recommendations = append(recommendations, "Consider adjusting task heights and colors to improve visual quality")
	}
	
	// Stack height recommendations
	if statistics.AverageStackHeight > cgi.gridConfig.DayHeight*2 {
		recommendations = append(recommendations, "Consider using horizontal stacking for high-density days")
	}
	
	// Conflict recommendations
	if statistics.ConflictsResolved > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Resolved %d visual conflicts - consider reviewing task scheduling", statistics.ConflictsResolved))
	}
	
	// Overflow recommendations
	if statistics.OverflowResolutions > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Applied %d overflow resolutions - consider reducing task density", statistics.OverflowResolutions))
	}
	
	return recommendations
}

// GenerateIntegratedLaTeX generates LaTeX code for the integrated calendar
func (cgi *CalendarGridIntegration) GenerateIntegratedLaTeX(result *IntegratedLayoutResult) string {
	var latex strings.Builder
	
	// Generate header
	latex.WriteString("\\begin{integrated-calendar}\n")
	
	// Generate task bars LaTeX
	for _, bar := range result.TaskBars {
		barLaTeX := cgi.generateTaskBarLaTeX(bar)
		latex.WriteString(barLaTeX)
	}
	
	// Generate footer
	latex.WriteString("\\end{integrated-calendar}\n")
	
	return latex.String()
}

// generateTaskBarLaTeX generates LaTeX code for a single task bar
func (cgi *CalendarGridIntegration) generateTaskBarLaTeX(bar *IntegratedTaskBar) string {
	// Create TikZ node for the task bar
    var nodeOptions string
    if bar.Opacity >= 0.999 {
        nodeOptions = fmt.Sprintf(
            "anchor=west, inner sep=2pt, minimum height=%.2fpt, minimum width=%.2fpt, fill=%s",
            bar.Height,
            bar.Width,
            bar.Color,
        )
    } else {
        nodeOptions = fmt.Sprintf(
            "anchor=west, inner sep=2pt, minimum height=%.2fpt, minimum width=%.2fpt, fill=%s, opacity=%.2f",
            bar.Height,
            bar.Width,
            bar.Color,
            bar.Opacity,
        )
    }
	
	// Add border if specified
	if cgi.visualSettings.BorderWidth > 0 {
		nodeOptions += fmt.Sprintf(", draw=%s, line width=%.2fpt", bar.BorderColor, cgi.visualSettings.BorderWidth)
	}
	
	// Create the TikZ node
	tikzNode := fmt.Sprintf(
		"\\node[%s] at (%.2fpt, %.2fpt) {%s};",
		nodeOptions,
		bar.StartX,
		bar.Y,
		bar.TaskName,
	)
	
	return tikzNode + "\n"
}

// GetIntegratedStatistics returns statistics about the integrated layout
func (cgi *CalendarGridIntegration) GetIntegratedStatistics(result *IntegratedLayoutResult) *IntegratedLayoutStatistics {
	return result.Statistics
}

// String returns a string representation of the integrated layout statistics
func (ils *IntegratedLayoutStatistics) String() string {
	return fmt.Sprintf("Integrated Layout Statistics:\n"+
		"  Total Tasks: %d\n"+
		"  Processed Bars: %d\n"+
		"  Total Stacks: %d\n"+
		"  Conflicts Resolved: %d\n"+
		"  Overflow Resolutions: %d\n"+
		"  Visual Optimizations: %d\n"+
		"  Layout Adjustments: %d\n"+
		"  Collision Count: %d\n"+
		"  Overflow Count: %d\n"+
		"  Month Boundary Count: %d\n"+
		"  Space Efficiency: %.2f\n"+
		"  Visual Quality: %.2f\n"+
		"  Average Stack Height: %.2f\n"+
		"  Max Stack Height: %.2f\n"+
		"  Average Task Height: %.2f\n"+
		"  Average Task Width: %.2f\n"+
		"  Alignment Score: %.2f\n"+
		"  Spacing Score: %.2f\n"+
		"  Visual Balance: %.2f\n"+
		"  Grid Utilization: %.2f\n",
		ils.TotalTasks, ils.ProcessedBars, ils.TotalStacks,
		ils.ConflictsResolved, ils.OverflowResolutions, ils.VisualOptimizations,
		ils.LayoutAdjustments, ils.CollisionCount, ils.OverflowCount,
		ils.MonthBoundaryCount, ils.SpaceEfficiency, ils.VisualQuality,
		ils.AverageStackHeight, ils.MaxStackHeight, ils.AverageTaskHeight, ils.AverageTaskWidth,
		ils.AlignmentScore, ils.SpacingScore, ils.VisualBalance, ils.GridUtilization)
}

// ProcessTasksWithValidation processes tasks with validation and creates multi-day layout
func (cgi *CalendarGridIntegration) ProcessTasksWithValidation(tasks []*data.Task) (*MultiDayLayoutResult, error) {
	// Validate tasks first
	validationResult := cgi.dateValidator.ValidateDateRanges(tasks)
	
	// Check for critical errors
	if len(validationResult) > 0 {
		// Log validation errors but continue with layout
		fmt.Printf("Warning: %d validation issues found\n", len(validationResult))
		for _, err := range validationResult {
			if err.Severity == "ERROR" {
				fmt.Printf("Error: %s\n", err.Message)
			}
		}
	}
	
	// Filter out tasks with critical errors for layout
	validTasks := cgi.filterValidTasks(tasks, validationResult)
	
	// Create multi-day layout
	taskBars := cgi.multiDayLayoutEngine.LayoutMultiDayTasks(validTasks)
	
	// Handle month boundaries
	processedBars := cgi.multiDayLayoutEngine.HandleMonthBoundary(taskBars)
	
	// Validate layout
	layoutIssues := cgi.multiDayLayoutEngine.ValidateLayout(processedBars)
	
	return &MultiDayLayoutResult{
		TaskBars:        processedBars,
		ValidationResult: validationResult,
		LayoutIssues:    layoutIssues,
		TaskCount:       len(validTasks),
		ProcessedCount:  len(processedBars),
	}, nil
}

// MultiDayLayoutResult contains the results of multi-day layout processing
type MultiDayLayoutResult struct {
	TaskBars        []*TaskBar
	ValidationResult []data.DataValidationError
	LayoutIssues    []string
	TaskCount       int
	ProcessedCount  int
}

// filterValidTasks filters out tasks with critical validation errors
func (cgi *CalendarGridIntegration) filterValidTasks(tasks []*data.Task, validationErrors []data.DataValidationError) []*data.Task {
	// Create map of tasks with critical errors
	errorTasks := make(map[string]bool)
	for _, err := range validationErrors {
		if err.Severity == "ERROR" {
			errorTasks[err.TaskID] = true
		}
	}
	
	// Filter out tasks with critical errors
	var validTasks []*data.Task
	for _, task := range tasks {
		if !errorTasks[task.ID] {
			validTasks = append(validTasks, task)
		}
	}
	
	return validTasks
}

// GenerateCalendarLaTeX generates LaTeX code for the calendar with multi-day task bars
func (cgi *CalendarGridIntegration) GenerateCalendarLaTeX(result *MultiDayLayoutResult) string {
	var latex strings.Builder
	
	// Generate header
	latex.WriteString("\\begin{calendar}\n")
	
	// Generate task bars LaTeX
	taskBarsLaTeX := cgi.multiDayLayoutEngine.GenerateLaTeX(result.TaskBars)
	latex.WriteString(taskBarsLaTeX)
	
	// Generate footer
	latex.WriteString("\\end{calendar}\n")
	
	return latex.String()
}

// GetLayoutStatistics returns statistics about the layout
func (cgi *CalendarGridIntegration) GetLayoutStatistics(result *MultiDayLayoutResult) *LayoutStatistics {
	stats := &LayoutStatistics{
		TotalTasks:      result.TaskCount,
		ProcessedBars:   result.ProcessedCount,
		ValidationErrors: len(result.ValidationResult),
		LayoutIssues:    len(result.LayoutIssues),
		OverlapCount:    0,
		MonthBoundaryCount: 0,
	}
	
	// Count overlaps and month boundaries
	for _, bar := range result.TaskBars {
		if bar.MonthBoundary {
			stats.MonthBoundaryCount++
		}
	}
	
	// Count overlaps by checking for overlapping bars
	rowBars := make(map[int][]*TaskBar)
	for _, bar := range result.TaskBars {
		rowBars[bar.Row] = append(rowBars[bar.Row], bar)
	}
	
	for _, bars := range rowBars {
		for i := 0; i < len(bars); i++ {
			for j := i + 1; j < len(bars); j++ {
				if cgi.multiDayLayoutEngine.barsOverlap(bars[i], bars[j]) {
					stats.OverlapCount++
				}
			}
		}
	}
	
	return stats
}

// LayoutStatistics contains statistics about the layout
type LayoutStatistics struct {
	TotalTasks         int
	ProcessedBars      int
	ValidationErrors   int
	LayoutIssues       int
	OverlapCount       int
	MonthBoundaryCount int
}

// String returns a string representation of the statistics
func (ls *LayoutStatistics) String() string {
	return fmt.Sprintf("Layout Statistics:\n"+
		"  Total Tasks: %d\n"+
		"  Processed Bars: %d\n"+
		"  Validation Errors: %d\n"+
		"  Layout Issues: %d\n"+
		"  Overlaps: %d\n"+
		"  Month Boundaries: %d\n",
		ls.TotalTasks, ls.ProcessedBars, ls.ValidationErrors,
		ls.LayoutIssues, ls.OverlapCount, ls.MonthBoundaryCount)
}
