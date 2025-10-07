// Package helpers provides utility functions extracted from the main generator.
//
// Monthly helpers provide functionality for processing monthly calendar data
// and task assignment.
package helpers

import (
	"fmt"
	"time"

	cal "phd-dissertation-planner/src/calendar"
	"phd-dissertation-planner/src/core"
)

// MonthlyProcessor handles monthly calendar generation and task processing
type MonthlyProcessor struct {
	cfg core.Config
}

// NewMonthlyProcessor creates a new monthly processor
func NewMonthlyProcessor(cfg core.Config) *MonthlyProcessor {
	return &MonthlyProcessor{cfg: cfg}
}

// ProcessMonthsWithTasks processes months that have tasks from CSV data
func (mp *MonthlyProcessor) ProcessMonthsWithTasks(tasks []core.Task, tpls []string) (core.Modules, error) {
	var modules core.Modules

	// Add table of contents if we have tasks
	if len(tasks) > 0 {
		tocModule := mp.createTableOfContentsModule(tasks, tpls[0])
		modules = append(modules, tocModule)
	}

	monthModules := make(core.Modules, 0, len(mp.cfg.MonthsWithTasks))

	for _, monthYear := range mp.cfg.MonthsWithTasks {
		monthModule, err := mp.processSingleMonth(monthYear, tasks, tpls[0])
		if err != nil {
			return nil, fmt.Errorf("failed to process month %s %d: %w", 
				monthYear.Month.String(), monthYear.Year, err)
		}
		monthModules = append(monthModules, monthModule)
	}

	// Combine TOC modules with month modules
	modules = append(modules, monthModules...)
	return modules, nil
}

// ProcessFallbackMonths processes months using the original fallback behavior
func (mp *MonthlyProcessor) ProcessFallbackMonths(tpls []string) (core.Modules, error) {
	years := mp.cfg.GetYears()
	totalMonths := len(years) * 12
	modules := make(core.Modules, 0, totalMonths)

	for _, yearNum := range years {
		year := cal.NewYear(mp.cfg.WeekStart, yearNum, &mp.cfg)

		for _, quarter := range year.Quarters {
			for _, month := range quarter.Months {
				module := mp.createMonthModule(*month, quarter, year, tpls[0])
				modules = append(modules, module)
			}
		}
	}

	return modules, nil
}

// processSingleMonth processes a single month with task assignment
func (mp *MonthlyProcessor) processSingleMonth(monthYear core.MonthYear, tasks []core.Task, templateName string) (core.Module, error) {
	year := cal.NewYear(mp.cfg.WeekStart, monthYear.Year, &mp.cfg)

	// Find the specific month in the year
	targetMonth := mp.findTargetMonth(year, monthYear)
	if targetMonth == nil {
		return core.Module{}, fmt.Errorf("month %s %d not found in calendar", 
			monthYear.Month.String(), monthYear.Year)
	}

	// Assign tasks to days in this month
	mp.assignTasksToMonth(targetMonth, tasks)

	return mp.createMonthModule(*targetMonth, targetMonth.Quarter, year, templateName), nil
}

// findTargetMonth finds the specific month in a year's calendar structure
func (mp *MonthlyProcessor) findTargetMonth(year *cal.Year, monthYear core.MonthYear) *cal.Month {
	if year == nil {
		return nil
	}
	
	for _, quarter := range year.Quarters {
		for _, month := range quarter.Months {
			if month.Month == monthYear.Month {
				return month
			}
		}
	}
	return nil
}

// createMonthModule creates a module for a single month
func (mp *MonthlyProcessor) createMonthModule(month cal.Month, quarter *cal.Quarter, year *cal.Year, templateName string) core.Module {
	// Handle nil inputs gracefully
	var sideQuarters interface{}
	var sideMonths interface{}
	var extra interface{}
	
	if year != nil {
		if quarter != nil {
			sideQuarters = year.SideQuarters(quarter.Number)
		}
		sideMonths = year.SideMonths(month.Month)
	}
	
	if month.Year != nil {
		extra = month.PrevNext().WithTopRightCorner(mp.cfg.ClearTopRightCorner, mp.cfg.Layout.Calendar.TaskKernSpacing)
	}

	return core.Module{
		Cfg: mp.cfg,
		Tpl: templateName,
		Body: map[string]interface{}{
			"Year":         year,
			"Quarter":      quarter,
			"Month":        month,
			"MonthRef":     fmt.Sprintf("month-%d-%d", month.Year.Number, int(month.Month)),
			"Breadcrumb":   month.Breadcrumb(),
			"HeadingMOS":   month.HeadingMOS(),
			"SideQuarters": sideQuarters,
			"SideMonths":   sideMonths,
			"Extra":        extra,
			"Large":        true,
			"TableType":    "tabularx",
			"Today":        cal.Day{Time: time.Now(), Cfg: &mp.cfg},
		},
	}
}

// createTableOfContentsModule creates a table of contents module
func (mp *MonthlyProcessor) createTableOfContentsModule(tasks []core.Task, templateName string) core.Module {
	builder := NewTOCBuilder()
	latexContent := builder.BuildTOCContent(tasks)

	return core.Module{
		Cfg: mp.cfg,
		Tpl: templateName,
		Body: map[string]interface{}{
			"TOCContent": latexContent,
		},
	}
}

// assignTasksToMonth assigns tasks to the appropriate days in a month
func (mp *MonthlyProcessor) assignTasksToMonth(month *cal.Month, tasks []core.Task) {
	if month == nil {
		return // Nothing to do with nil month
	}
	
	// Convert data.Task to SpanningTask and apply to month
	var spanningTasks []cal.SpanningTask

	for _, task := range tasks {
		// Check if task overlaps with this month
		monthStart := time.Date(month.Year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
		monthEnd := monthStart.AddDate(0, 1, -1)

		if task.StartDate.Before(monthEnd.AddDate(0, 0, 1)) && task.EndDate.After(monthStart.AddDate(0, 0, -1)) {
			// Create spanning task directly from common.Task
			// Rendering rules:
			// - Start day: show a thin colored bar + a single concise text label.
			// - Middle/end days: show only the bar (no repeated labels).
			// Therefore, we DO NOT add this task as a regular per-day entry to avoid duplication.
			spanningTask := cal.CreateSpanningTask(task, task.StartDate, task.EndDate)
			spanningTasks = append(spanningTasks, spanningTask)
		}
	}

	// Apply spanning tasks to the month for background coloring
	cal.ApplySpanningTasksToMonth(month, spanningTasks)
}
