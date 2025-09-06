package generator

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"

	cal "latex-yearly-planner/internal/calendar"
	"latex-yearly-planner/internal/config"
	"latex-yearly-planner/internal/data"
	tmplfs "latex-yearly-planner/templates"
)

// Generator handles LaTeX document generation with integrated layout processing
type Generator struct {
	templateEngine *template.Template
	layoutIntegration *LayoutIntegration
}

// NewGenerator creates a new generator instance
func NewGenerator() *Generator {
	return &Generator{
		templateEngine: createTemplateEngine(),
		layoutIntegration: NewLayoutIntegration(),
	}
}

// createTemplateEngine initializes the template engine with all necessary functions
func createTemplateEngine() *template.Template {
	t := template.New("").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"incr": func(i int) int { return i + 1 },
		"dec": func(i int) int { return i - 1 },
		"is": func(i interface{}) bool {
			if value, ok := i.(bool); ok {
				return value
			}
			return i != nil
		},
		// Layout integration functions
		"hasLayoutData": func(data interface{}) bool {
			if data == nil {
				return false
			}
			if m, ok := data.(map[string]interface{}); ok {
				_, hasLayout := m["LayoutResult"]
				_, hasTaskBars := m["TaskBars"]
				return hasLayout || hasTaskBars
			}
			return false
		},
		"getTaskBars": func(data interface{}) []*cal.IntegratedTaskBar {
			if m, ok := data.(map[string]interface{}); ok {
				if bars, ok := m["TaskBars"].([]*cal.IntegratedTaskBar); ok {
					return bars
				}
			}
			return nil
		},
		"getLayoutStats": func(data interface{}) *cal.IntegratedLayoutStatistics {
			if m, ok := data.(map[string]interface{}); ok {
				if stats, ok := m["LayoutStats"].(*cal.IntegratedLayoutStatistics); ok {
					return stats
				}
			}
			return nil
		},
		"formatTaskBar": func(bar *cal.IntegratedTaskBar) string {
			if bar == nil {
				return ""
			}
			var prominence string
			switch {
			case bar.Priority >= 4:
				prominence = "CRITICAL"
			case bar.Priority >= 3:
				prominence = "HIGH"
			case bar.Priority >= 2:
				prominence = "MEDIUM"
			case bar.Priority >= 1:
				prominence = "LOW"
			default:
				prominence = "MINIMAL"
			}
			
			return fmt.Sprintf("\\TaskOverlayBoxP{%s}{%s}{%s}{%s}",
				prominence, bar.Color, bar.TaskName, bar.Description)
		},
	})

	// Choose template source: embedded by default, filesystem when DEV_TEMPLATES is set
	var useFS fs.FS
	var err error

	if os.Getenv("DEV_TEMPLATES") != "" {
		useFS = os.DirFS(filepath.Join("templates", "monthly"))
	} else {
		sub, err := fs.Sub(tmplfs.FS, "monthly")
		if err != nil {
			panic(fmt.Sprintf("failed to sub FS for monthly templates: %v", err))
		}
		useFS = sub
	}

	// Parse all *.tpl templates
	t, err = t.ParseFS(useFS, "*.tpl")
	if err != nil {
		panic(fmt.Sprintf("failed to parse monthly templates: %v", err))
	}

	return t
}

// Document generates the main document template
func (g *Generator) Document(wr io.Writer, cfg config.Config) error {
	type pack struct {
		Cfg   config.Config
		Pages []config.Page
	}

	data := pack{Cfg: cfg, Pages: cfg.Pages}
	return g.templateEngine.ExecuteTemplate(wr, "main_document.tpl", data)
}

// Execute runs a specific template
func (g *Generator) Execute(wr io.Writer, name string, data interface{}) error {
	return g.templateEngine.ExecuteTemplate(wr, name, data)
}

// Monthly generates monthly modules with optional layout integration
func (g *Generator) Monthly(cfg config.Config, tpls []string) (config.Modules, error) {
	// Use enhanced monthly generation with layout processing
	return g.layoutIntegration.EnhancedMonthly(cfg, tpls)
}

// MonthlyLegacy provides the original monthly generation without layout integration
func (g *Generator) MonthlyLegacy(cfg config.Config, tpls []string) (config.Modules, error) {
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

			if targetMonth == nil {
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
					"TableType":    "tabularx",
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
						"TableType":    "tabularx",
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
	var spanningTasks []cal.SpanningTask

	for _, task := range tasks {
		// Check if task overlaps with this month
		monthStart := time.Date(month.Year.Number, month.Month, 1, 0, 0, 0, 0, time.Local)
		monthEnd := monthStart.AddDate(0, 1, -1)

		if task.StartDate.Before(monthEnd.AddDate(0, 0, 1)) && task.EndDate.After(monthStart.AddDate(0, 0, -1)) {
			spanningTask := cal.CreateSpanningTask(task, task.StartDate, task.EndDate)
			spanningTasks = append(spanningTasks, spanningTask)
		}
	}

	// Apply spanning tasks to the month for background coloring
	cal.ApplySpanningTasksToMonth(month, spanningTasks)
}

// LayoutIntegration bridges the integrated layout system with the template system
type LayoutIntegration struct {
	gridConfig   *cal.GridConfig
	integration  *cal.CalendarGridIntegration
}

// NewLayoutIntegration creates a new layout integration instance
func NewLayoutIntegration() *LayoutIntegration {
	// Default grid configuration
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
	layoutTasks := make([]*data.Task, len(tasks))
	copy(layoutTasks, tasks)

	result, err := li.integration.ProcessTasksWithSmartStacking(layoutTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to process tasks with layout: %w", err)
	}

	return result, nil
}

// GenerateTaskVisualization generates LaTeX code for task visualization
func (li *LayoutIntegration) GenerateTaskVisualization(result *cal.IntegratedLayoutResult) string {
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

			// Assign tasks to days in this month
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

	return result.Statistics
}

// UpdateGridConfig updates the grid configuration for the layout integration
func (li *LayoutIntegration) UpdateGridConfig(config *cal.GridConfig) {
	li.gridConfig = config
	li.integration = cal.NewCalendarGridIntegration(config)
}

// Legacy functions for backward compatibility
func Monthly(cfg config.Config, tpls []string) (config.Modules, error) {
	generator := NewGenerator()
	return generator.Monthly(cfg, tpls)
}

func MonthlyLegacy(cfg config.Config, tpls []string) (config.Modules, error) {
	generator := NewGenerator()
	return generator.MonthlyLegacy(cfg, tpls)
}

func NewTpl() Generator {
	return *NewGenerator()
}

