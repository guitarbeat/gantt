package main

import (
	"github.com/kudrykv/latex-yearly-planner/cal"
)

func Monthly(cfg Config, tpls []string) (Modules, error) {
	// If we have months with tasks from CSV, use only those
	if len(cfg.MonthsWithTasks) > 0 {
		modules := make(Modules, 0, len(cfg.MonthsWithTasks))

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
				modules = append(modules, Module{
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
	modules := make(Modules, 0, totalMonths)

	for _, yearNum := range years {
		year := cal.NewYear(cfg.WeekStart, yearNum)

		for _, quarter := range year.Quarters {
			for _, month := range quarter.Months {
				modules = append(modules, Module{
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
