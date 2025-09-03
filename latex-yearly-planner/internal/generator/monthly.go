package generator

import (
	"fmt"
	"time"

	"github.com/kudrykv/latex-yearly-planner/internal/config"
	"github.com/kudrykv/latex-yearly-planner/internal/data"
	cal "github.com/kudrykv/latex-yearly-planner/pkg/calendar"
)

func Monthly(cfg config.Config, tpls []string) (config.Modules, error) {
	// Load tasks from CSV if available
	var tasks []data.Task
	if cfg.CSVFilePath != "" {
		reader := data.NewReader(cfg.CSVFilePath)
		var err error
		tasks, err = reader.ReadTasks()
		if err != nil {
			// Log error but continue without tasks
			return nil, fmt.Errorf("error reading tasks: %w", err)
		}
	}

	// If we have months with tasks from CSV, use only those
	if len(cfg.MonthsWithTasks) > 0 {
		modules := make(config.Modules, 0, len(cfg.MonthsWithTasks))

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

			// * Check if targetMonth was found, log warning if not
			if targetMonth == nil {
				// Log warning but continue processing other months
				fmt.Printf("Warning: Month %s %d not found in calendar, skipping\n", 
					monthYear.Month.String(), monthYear.Year)
				continue
			}

			// Assign tasks to days in this month
			assignTasksToMonth(targetMonth, tasks)
			
			modules = append(modules, config.Module{
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
					"TableType":    "monthly",
					"Today":        cal.Day{Time: time.Now()},
				},
			})
		}

		return modules, nil
	}

	// Fallback to original behavior if no CSV data
	years := cfg.GetYears()
	totalMonths := len(years) * 12
	modules := make(config.Modules, 0, totalMonths)

	for _, yearNum := range years {
		year := cal.NewYear(cfg.WeekStart, yearNum)

		for _, quarter := range year.Quarters {
			for _, month := range quarter.Months {
				modules = append(modules, config.Module{
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
						"TableType":    "monthly",
						"Today":        cal.Day{Time: time.Now()},
					},
				})
			}
		}
	}

	return modules, nil
}

// assignTasksToMonth assigns tasks to the appropriate days in a month
func assignTasksToMonth(month *cal.Month, tasks []data.Task) {
	// Convert data.Task to SpanningTask and apply to month
	var spanningTasks []cal.SpanningTask
	
	for _, task := range tasks {
		// Check if task overlaps with this month
		monthStart := time.Date(month.Year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
		monthEnd := monthStart.AddDate(0, 1, -1)
		
		if task.StartDate.Before(monthEnd.AddDate(0, 0, 1)) && task.EndDate.After(monthStart.AddDate(0, 0, -1)) {
			// Create spanning task directly from data.Task
			spanningTask := cal.CreateSpanningTask(task, task.StartDate, task.EndDate)
			spanningTasks = append(spanningTasks, spanningTask)
			
			// Also add as regular tasks to individual days for better display
			current := task.StartDate
			for !current.After(task.EndDate) {
				// Check if this day is in the current month
				if current.Month() == month.Month && current.Year() == month.Year.Number {
					// Find the day in the month and add the task
					for _, week := range month.Weeks {
						for i := range week.Days {
							if week.Days[i].Time.Day() == current.Day() && 
							   week.Days[i].Time.Month() == current.Month() &&
							   week.Days[i].Time.Year() == current.Year() {
								// Add as regular task for detailed display
								dayTask := cal.Task{
									ID:          task.ID,
									Name:        task.Name,
									Description: task.Description,
									Category:    task.Priority, // Category stored in Priority field
								}
								week.Days[i].Tasks = append(week.Days[i].Tasks, dayTask)
								break
							}
						}
					}
				}
				current = current.AddDate(0, 0, 1)
			}
		}
	}
	
	// Apply spanning tasks to the month for background coloring
	cal.ApplySpanningTasksToMonth(month, spanningTasks)
}
