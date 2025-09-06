package generator

import (
	"fmt"
	"time"

	cal "latex-yearly-planner/internal/calendar"
	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
)

// LayoutIntegration bridges the integrated layout system with the template system
type LayoutIntegration struct {
	gridConfig *cal.GridConfig
	integration *cal.CalendarGridIntegration
}

// NewLayoutIntegration creates a new layout integration instance
func NewLayoutIntegration() *LayoutIntegration {
	// Default grid configuration - can be made configurable later
	gridConfig := &cal.GridConfig{
		CalendarStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CalendarEnd:      time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		DayWidth:         20.0,
		DayHeight:        15.0,
		RowHeight:        12.0,
		MaxRowsPerDay:    5,
		OverlapThreshold: 0.1,
		MonthBoundaryGap: 2.0,
		TaskSpacing:      1.0,
	}

	integration := cal.NewCalendarGridIntegration(gridConfig)

	return &LayoutIntegration{
		gridConfig:  gridConfig,
		integration: integration,
	}
}

// ProcessTasksWithLayout processes tasks using the integrated layout system
func (li *LayoutIntegration) ProcessTasksWithLayout(tasks []*data.Task) (*cal.IntegratedLayoutResult, error) {
	// Convert data.Task to the format expected by the layout system
	layoutTasks := make([]*data.Task, len(tasks))
	copy(layoutTasks, tasks)

	// Process tasks with smart stacking and positioning
	result, err := li.integration.ProcessTasksWithSmartStacking(layoutTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to process tasks with layout: %w", err)
	}

	return result, nil
}

// GenerateTaskVisualization generates LaTeX code for task visualization
func (li *LayoutIntegration) GenerateTaskVisualization(result *cal.IntegratedLayoutResult) string {
	// Use the existing GenerateIntegratedLaTeX method
	return li.integration.GenerateIntegratedLaTeX(result)
}

// GetLayoutStatistics returns statistics about the layout processing
func (li *LayoutIntegration) GetLayoutStatistics(result *cal.IntegratedLayoutResult) *cal.IntegratedLayoutStatistics {
	return li.integration.GetIntegratedStatistics(result)
}

// EnhancedMonthly generates monthly modules with integrated layout processing
func (li *LayoutIntegration) EnhancedMonthly(cfg config.Config, tpls []string) (config.Modules, error) {
	// Load tasks from CSV if available
	var tasks []data.Task
	if cfg.CSVFilePath != "" {
		reader := data.NewReader(cfg.CSVFilePath)
		var err error
		tasks, err = reader.ReadTasks()
		if err != nil {
			return nil, fmt.Errorf("error reading tasks: %w", err)
		}
	}

	// Process tasks with integrated layout system
	var layoutResult *cal.IntegratedLayoutResult
	if len(tasks) > 0 {
		// Convert to pointer slice for layout processing
		var taskPointers []*data.Task
		for i := range tasks {
			taskPointers = append(taskPointers, &tasks[i])
		}

		var err error
		layoutResult, err = li.ProcessTasksWithLayout(taskPointers)
		if err != nil {
			return nil, fmt.Errorf("failed to process tasks with layout: %w", err)
		}
	}

	// Generate modules with enhanced layout data
	modules := make(config.Modules, 0)

	// If we have months with tasks from CSV, use only those
	if len(cfg.MonthsWithTasks) > 0 {
		for _, monthYear := range cfg.MonthsWithTasks {
			year := cal.NewYear(cfg.WeekStart, monthYear.Year)

			// Find the specific month in the year
			var targetMonth *cal.Month
			for _, quarter := range year.Quarters {
				for _, month := range quarter.Months {
					if month.Month == monthYear.Month {
						targetMonth = month
						break
					}
				}
				if targetMonth != nil {
					break
				}
			}

			if targetMonth == nil {
				fmt.Printf("Warning: Month %s %d not found in calendar, skipping\n",
					monthYear.Month.String(), monthYear.Year)
				continue
			}

			// Assign tasks to days in this month (existing logic)
			assignTasksToMonth(targetMonth, tasks)

			// Create enhanced module with layout data
			module := config.Module{
				Cfg: cfg,
				Tpl: tpls[0],
				Body: map[string]interface{}{
					"Year":         year,
					"Quarter":      targetMonth.Quarter,
					"Month":        targetMonth,
					"Breadcrumb":   targetMonth.Breadcrumb(),
					"HeadingMOS":   targetMonth.HeadingMOS(),
					"SideQuarters": year.SideQuarters(targetMonth.Quarter.Number),
					"SideMonths":   year.SideMonths(targetMonth.Month),
					"Extra":        targetMonth.PrevNext().WithTopRightCorner(cfg.ClearTopRightCorner),
					"Large":        true,
					"TableType":    "tabularx",
					"Today":        cal.Day{Time: time.Now()},
					// Enhanced layout data
					"LayoutResult": layoutResult,
					"TaskBars":     li.getTaskBarsForMonth(layoutResult, targetMonth),
					"LayoutStats":  li.getLayoutStatsForMonth(layoutResult, targetMonth),
				},
			}

			modules = append(modules, module)
		}

		return modules, nil
	}

	// Fallback to original behavior if no CSV data
	years := cfg.GetYears()
	for _, yearNum := range years {
		year := cal.NewYear(cfg.WeekStart, yearNum)

		for _, quarter := range year.Quarters {
			for _, month := range quarter.Months {
				module := config.Module{
					Cfg: cfg,
					Tpl: tpls[0],
					Body: map[string]interface{}{
						"Year":         year,
						"Quarter":      quarter,
						"Month":        month,
						"Breadcrumb":   month.Breadcrumb(),
						"HeadingMOS":   month.HeadingMOS(),
						"SideQuarters": year.SideQuarters(quarter.Number),
						"SideMonths":   year.SideMonths(month.Month),
						"Extra":        month.PrevNext().WithTopRightCorner(cfg.ClearTopRightCorner),
						"Large":        true,
						"TableType":    "tabularx",
						"Today":        cal.Day{Time: time.Now()},
						// Enhanced layout data
						"LayoutResult": layoutResult,
						"TaskBars":     li.getTaskBarsForMonth(layoutResult, month),
						"LayoutStats":  li.getLayoutStatsForMonth(layoutResult, month),
					},
				}

				modules = append(modules, module)
			}
		}
	}

	return modules, nil
}

// getTaskBarsForMonth filters task bars for a specific month
func (li *LayoutIntegration) getTaskBarsForMonth(result *cal.IntegratedLayoutResult, month *cal.Month) []*cal.IntegratedTaskBar {
	if result == nil {
		return nil
	}

	var monthBars []*cal.IntegratedTaskBar
	monthStart := time.Date(month.Year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
	monthEnd := monthStart.AddDate(0, 1, -1)

	for _, bar := range result.TaskBars {
		// Check if task bar overlaps with this month
		if bar.StartDate.Before(monthEnd.AddDate(0, 0, 1)) && bar.EndDate.After(monthStart.AddDate(0, 0, -1)) {
			monthBars = append(monthBars, bar)
		}
	}

	return monthBars
}

// getLayoutStatsForMonth returns layout statistics for a specific month
func (li *LayoutIntegration) getLayoutStatsForMonth(result *cal.IntegratedLayoutResult, month *cal.Month) *cal.IntegratedLayoutStatistics {
	if result == nil {
		return nil
	}

	// For now, return the overall statistics
	// In the future, we could calculate month-specific statistics
	return result.Statistics
}

// UpdateGridConfig updates the grid configuration for the layout integration
func (li *LayoutIntegration) UpdateGridConfig(config *cal.GridConfig) {
	li.gridConfig = config
	li.integration = cal.NewCalendarGridIntegration(config)
}
