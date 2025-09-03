package generator

import (
	"fmt"

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

			if targetMonth != nil {
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
					},
				})
			}
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
					},
				})
			}
		}
	}

	return modules, nil
}

// assignTasksToMonth assigns tasks to the appropriate days in a month
func assignTasksToMonth(month *cal.Month, tasks []data.Task) {
	for _, task := range tasks {
		// Check if task falls within this month
		if task.StartDate.Month() == month.Month && task.StartDate.Year() == month.Year.Number {
			// Find the day in the month and add the task
			for _, week := range month.Weeks {
				for _, day := range week.Days {
					if day.Time.Day() == task.StartDate.Day() {
						// Convert data.Task to calendar.Task
						calTask := cal.Task{
							ID:          task.ID,
							Name:        task.Name,
							Description: task.Description,
							Category:    task.Priority, // Use Priority field which now contains Category
						}
						day.Tasks = append(day.Tasks, calTask)
						break
					}
				}
			}
		}
	}
}
