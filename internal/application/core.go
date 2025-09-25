package application

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"phd-dissertation-planner/internal/common"
	cal "phd-dissertation-planner/internal/scheduler"
	tmplfs "phd-dissertation-planner/pkg/templates"

	"github.com/urfave/cli/v2"
)

const (
	fConfig = "config"
	pConfig = "preview"
	fOutDir = "outdir"
)

func New() *cli.App {
	// Initialize the composer map
	common.ComposerMap["monthly"] = Monthly

	return &cli.App{
		Name: "plannergen",

		Writer:    os.Stdout,
		ErrWriter: os.Stderr,

		Flags: []cli.Flag{
			&cli.PathFlag{Name: fConfig, Required: false, Value: "internal/common/base.yaml", Usage: "config file(s), comma-separated"},
			&cli.BoolFlag{Name: pConfig, Required: false, Usage: "render only one page per unique module"},
			&cli.PathFlag{Name: fOutDir, Required: false, Value: "", Usage: "output directory for generated files (overrides config)"},
		},

		Action: action,
	}
}

func action(c *cli.Context) error {
	var (
		fn  common.Composer
		ok  bool
		cfg common.Config
		err error
	)

	preview := c.Bool(pConfig)

	pathConfigs := strings.Split(c.Path(fConfig), ",")
	if cfg, err = common.NewConfig(pathConfigs...); err != nil {
		return fmt.Errorf("failed to load configuration from %v: %w", pathConfigs, err)
	}

	// If CLI flag for outdir provided, override config
	if od := strings.TrimSpace(c.Path(fOutDir)); od != "" {
		cfg.OutputDir = od
	}

	// Ensure output directory exists
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory %q: %w", cfg.OutputDir, err)
	}
	if os.Getenv("PLANNER_SILENT") != "1" {
		fmt.Fprintf(os.Stderr, "Output directory: %s\n", cfg.OutputDir)
	}

	wr := &bytes.Buffer{}

	t := NewTpl()

	if err = t.Document(wr, cfg); err != nil {
		return fmt.Errorf("failed to generate LaTeX document: %w", err)
	}

	outputFile := cfg.OutputDir + "/" + RootFilename(pathConfigs[len(pathConfigs)-1])
	if err = os.WriteFile(outputFile, wr.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write LaTeX file to %q: %w", outputFile, err)
	}
	if os.Getenv("PLANNER_SILENT") != "1" {
		fmt.Fprintf(os.Stderr, "Generated LaTeX file: %s\n", outputFile)
	}

	for _, file := range cfg.Pages {
		wr.Reset()

		var mom []common.Modules
		for _, block := range file.RenderBlocks {
			if fn, ok = common.ComposerMap[block.FuncName]; !ok {
				return fmt.Errorf("unknown composer function %q - check configuration", block.FuncName)
			}

			modules, err := fn(cfg, block.Tpls)
			if err != nil {
				return fmt.Errorf("failed to compose modules for %q: %w", block.FuncName, err)
			}

			// Only one page per unique module if preview flag is enabled
			if preview {
				modules = common.FilterUniqueModules(modules)
			}

			mom = append(mom, modules)
		}

		if len(mom) == 0 {
			return fmt.Errorf("no modules generated for page %q", file.Name)
		}

		allLen := len(mom[0])
		for _, mods := range mom {
			if len(mods) != allLen {
				return fmt.Errorf("module alignment error for page %q: expected %d modules, got %d", file.Name, allLen, len(mods))
			}
		}

		for i := 0; i < allLen; i++ {
			for j, mod := range mom {
				if err = t.Execute(wr, mod[i].Tpl, mod[i]); err != nil {
					return fmt.Errorf("failed to execute template %s for function %s: %w", mod[i].Tpl, file.RenderBlocks[j].FuncName, err)
				}
			}
		}

		pageFile := cfg.OutputDir + "/" + file.Name + ".tex"
		if err = os.WriteFile(pageFile, wr.Bytes(), 0o600); err != nil {
			return fmt.Errorf("failed to write page file %q: %w", pageFile, err)
		}
		if os.Getenv("PLANNER_SILENT") != "1" {
			fmt.Fprintf(os.Stderr, "Generated page: %s\n", pageFile)
		}
	}

	return nil
}

func RootFilename(pathconfig string) string {
	filename := filepath.Base(pathconfig)
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + ".tex"
}

var tpl = func() *template.Template {
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

		"incr": func(i int) int {
			return i + 1
		},

		"dec": func(i int) int {
			return i - 1
		},

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
			// Check if data has layout-related fields
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

			// Generate LaTeX for individual task bar using the visual design system
			return fmt.Sprintf("\\TaskOverlayBox{%s}{%s}{%s}",
				bar.Color,       // category color
				bar.TaskName,    // task name
				bar.Description, // description
			)
		},
	})

	// Choose source of templates: embedded by default, filesystem when DEV_TEMPLATES is set
	var (
		err   error
		useFS fs.FS
	)

	if os.Getenv("DEV_TEMPLATES") != "" {
		// Use on-disk templates for development override
		useFS = os.DirFS(filepath.Join("templates", "monthly"))
	} else {
		// Use embedded templates from templates.FS
		// Narrow to monthly/ subdir
		var sub fs.FS
		sub, err = fs.Sub(tmplfs.FS, "monthly")
		if err != nil {
			panic(fmt.Sprintf("failed to sub FS for monthly templates: %v", err))
		}
		useFS = sub
	}

	// Parse all *.tpl templates from the selected FS
	t, err = t.ParseFS(useFS, "*.tpl")
	if err != nil {
		panic(fmt.Sprintf("failed to parse monthly templates: %v", err))
	}

	return t
}()

type Tpl struct {
	tpl *template.Template
}

func NewTpl() Tpl {
	return Tpl{
		tpl: tpl,
	}
}

func (t Tpl) Document(wr io.Writer, cfg common.Config) error {
	type pack struct {
		Cfg   common.Config
		Pages []common.Page
	}

	data := pack{Cfg: cfg, Pages: cfg.Pages}
	if err := t.tpl.ExecuteTemplate(wr, "document.tpl", data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func (t Tpl) Execute(wr io.Writer, name string, data interface{}) error {
	if err := t.tpl.ExecuteTemplate(wr, name, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
func Monthly(cfg common.Config, tpls []string) (common.Modules, error) {
	// Use legacy monthly generation without layout integration
	return MonthlyLegacy(cfg, tpls)
}

// MonthlyLegacy provides the original monthly generation without layout integration
func MonthlyLegacy(cfg common.Config, tpls []string) (common.Modules, error) {
	// Load tasks from CSV if available
	var tasks []common.Task
	if cfg.CSVFilePath != "" {
		reader := common.NewReader(cfg.CSVFilePath)
		var err error
		tasks, err = reader.ReadTasks()
		if err != nil {
			// Log error but continue without tasks
			return nil, fmt.Errorf("error reading tasks: %w", err)
		}
	}

	// If we have months with tasks from CSV, use only those
	if len(cfg.MonthsWithTasks) > 0 {
		modules := make(common.Modules, 0, len(cfg.MonthsWithTasks))

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

			modules = append(modules, common.Module{
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
					"Extra":        targetMonth.PrevNext().WithTopRightCorner(cfg.ClearTopRightCorner, cfg.Layout.Calendar.TaskKernSpacing),
					"Large":        true,
					"TableType":    "tabularx",
					"Today":        cal.Day{Time: time.Now()},
				},
			})
		}

		return modules, nil
	} else {
		// Fallback to original behavior if no CSV data
		years := cfg.GetYears()
		totalMonths := len(years) * 12
		modules := make(common.Modules, 0, totalMonths)

		for _, yearNum := range years {
			year := cal.NewYear(cfg.WeekStart, yearNum)

			for _, quarter := range year.Quarters {
				for _, month := range quarter.Months {
					modules = append(modules, common.Module{
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
							"Extra":        month.PrevNext().WithTopRightCorner(cfg.ClearTopRightCorner, cfg.Layout.Calendar.TaskKernSpacing),
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
}

// assignTasksToMonth assigns tasks to the appropriate days in a month
func assignTasksToMonth(month *cal.Month, tasks []common.Task) {
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
