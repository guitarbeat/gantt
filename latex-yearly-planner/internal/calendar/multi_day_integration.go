package calendar

import (
	"fmt"
	"strings"
	"time"

	"latex-yearly-planner/internal/data"
)

// MultiDayIntegration provides integration between the multi-day layout engine and existing calendar system
type MultiDayIntegration struct {
	layoutEngine *MultiDayLayoutEngine
	dateValidator *data.DateValidator
}

// NewMultiDayIntegration creates a new multi-day integration instance
func NewMultiDayIntegration(calendarStart, calendarEnd time.Time, dayWidth, dayHeight float64) *MultiDayIntegration {
	return &MultiDayIntegration{
		layoutEngine:  NewMultiDayLayoutEngine(calendarStart, calendarEnd, dayWidth, dayHeight),
		dateValidator: data.NewDateValidator(),
	}
}

// ProcessTasksWithValidation processes tasks with validation and creates multi-day layout
func (mdi *MultiDayIntegration) ProcessTasksWithValidation(tasks []*data.Task) (*MultiDayLayoutResult, error) {
	// Validate tasks first
	validationResult := mdi.dateValidator.ValidateDateRanges(tasks)
	
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
	validTasks := mdi.filterValidTasks(tasks, validationResult)
	
	// Create multi-day layout
	taskBars := mdi.layoutEngine.LayoutMultiDayTasks(validTasks)
	
	// Handle month boundaries
	processedBars := mdi.layoutEngine.HandleMonthBoundary(taskBars)
	
	// Validate layout
	layoutIssues := mdi.layoutEngine.ValidateLayout(processedBars)
	
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
func (mdi *MultiDayIntegration) filterValidTasks(tasks []*data.Task, validationErrors []data.DataValidationError) []*data.Task {
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
func (mdi *MultiDayIntegration) GenerateCalendarLaTeX(result *MultiDayLayoutResult) string {
	var latex strings.Builder
	
	// Generate header
	latex.WriteString("\\begin{calendar}\n")
	
	// Generate task bars LaTeX
	taskBarsLaTeX := mdi.layoutEngine.GenerateLaTeX(result.TaskBars)
	latex.WriteString(taskBarsLaTeX)
	
	// Generate footer
	latex.WriteString("\\end{calendar}\n")
	
	return latex.String()
}

// GetLayoutStatistics returns statistics about the layout
func (mdi *MultiDayIntegration) GetLayoutStatistics(result *MultiDayLayoutResult) *LayoutStatistics {
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
				if mdi.layoutEngine.barsOverlap(bars[i], bars[j]) {
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

// ExampleUsage demonstrates how to use the multi-day layout integration
func ExampleUsage() {
	// Create calendar range
	calendarStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	calendarEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	// Create integration instance
	integration := NewMultiDayIntegration(calendarStart, calendarEnd, 20.0, 15.0)
	
	// Create sample tasks
	tasks := []*data.Task{
		{
			ID:        "task1",
			Name:      "Multi-day Task 1",
			StartDate: time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			Category:  "PROPOSAL",
			Priority:  1,
		},
		{
			ID:        "task2",
			Name:      "Overlapping Task",
			StartDate: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Category:  "LASER",
			Priority:  2,
		},
		{
			ID:        "task3",
			Name:      "Cross-Month Task",
			StartDate: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
			Category:  "IMAGING",
			Priority:  3,
		},
	}
	
	// Process tasks with validation
	result, err := integration.ProcessTasksWithValidation(tasks)
	if err != nil {
		fmt.Printf("Error processing tasks: %v\n", err)
		return
	}
	
	// Generate LaTeX
	latex := integration.GenerateCalendarLaTeX(result)
	fmt.Println("Generated LaTeX:")
	fmt.Println(latex)
	
	// Get statistics
	stats := integration.GetLayoutStatistics(result)
	fmt.Println(stats.String())
}
