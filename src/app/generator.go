package app

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

	cal "phd-dissertation-planner/src/calendar"
	"phd-dissertation-planner/src/core"
	tmplfs "phd-dissertation-planner/src/shared/templates"

	"github.com/urfave/cli/v2"
)

// Constants for file operations and environment variables
const (
	// File extensions
	texExtension = ".tex"

	// Environment variables
	envDevTemplate = "DEV_TEMPLATES"

	// Directory paths
	templateSubDir = "monthly"
	templatePath   = "src/shared/templates/monthly"

	// Template patterns
	templatePattern = "*.tpl"
	documentTpl     = "document.tpl"
)

var logger = core.NewDefaultLogger()

// action is the main CLI action that orchestrates document generation
func action(c *cli.Context) error {
	// Load and prepare configuration
	cfg, pathConfigs, err := loadConfiguration(c)
	if err != nil {
		return err
	}

	// Setup output directory
	if err := setupOutputDirectory(cfg); err != nil {
		return err
	}

	// Generate root document
	if err := generateRootDocument(cfg, pathConfigs); err != nil {
		return err
	}

	// Generate pages
	preview := c.Bool(pConfig)
	if err := generatePages(cfg, preview); err != nil {
		return err
	}

	return nil
}

// loadConfiguration loads and validates the configuration from CLI context
func loadConfiguration(c *cli.Context) (core.Config, []string, error) {
	pathConfigs := strings.Split(c.Path(fConfig), ",")
	cfg, err := core.NewConfig(pathConfigs...)
	if err != nil {
		return core.Config{}, nil, core.NewConfigError(
			strings.Join(pathConfigs, ","),
			"",
			"failed to load configuration",
			err,
		)
	}

	// Override output directory from CLI flag if provided
	if od := strings.TrimSpace(c.Path(fOutDir)); od != "" {
		cfg.OutputDir = od
	}

	return cfg, pathConfigs, nil
}

// setupOutputDirectory ensures the output directory exists and logs its location
func setupOutputDirectory(cfg core.Config) error {
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return core.NewFileError(cfg.OutputDir, "create directory", err)
	}
	logger.Info("Output directory: %s", cfg.OutputDir)
	return nil
}

// generateRootDocument creates the main LaTeX document file
func generateRootDocument(cfg core.Config, pathConfigs []string) error {
	wr := &bytes.Buffer{}
	t := NewTpl()

	if err := t.Document(wr, cfg); err != nil {
		return core.NewTemplateError(documentTpl, 0, "failed to generate LaTeX document", err)
	}

	outputFile := filepath.Join(cfg.OutputDir, RootFilename(pathConfigs[len(pathConfigs)-1]))
	if err := os.WriteFile(outputFile, wr.Bytes(), 0o600); err != nil {
		return core.NewFileError(outputFile, "write", err)
	}
	logger.Info("Generated LaTeX file: %s", outputFile)
	return nil
}

// generatePages creates all page files from the configuration
func generatePages(cfg core.Config, preview bool) error {
	t := NewTpl()

	for _, file := range cfg.Pages {
		if err := generateSinglePage(cfg, file, t, preview); err != nil {
			return err
		}
	}

	return nil
}

// generateSinglePage generates a single page file
func generateSinglePage(cfg core.Config, file core.Page, t Tpl, preview bool) error {
	wr := &bytes.Buffer{}

	// Compose all modules for this page
	modules, err := composePageModules(cfg, file, preview)
	if err != nil {
		return err
	}

	// Validate module alignment
	if err := validateModuleAlignment(modules, file.Name); err != nil {
		return err
	}

	// Render modules to buffer
	if err := t.renderModules(wr, modules, file); err != nil {
		return err
	}

	// Write page file
	return writePageFile(cfg, file.Name, wr.Bytes())
}

// composePageModules composes all modules for a page by calling composer functions
func composePageModules(cfg core.Config, file core.Page, preview bool) ([]core.Modules, error) {
	var modules []core.Modules

	for _, block := range file.RenderBlocks {
		fn, ok := core.ComposerMap[block.FuncName]
		if !ok {
			return nil, fmt.Errorf("unknown composer function %q - check configuration", block.FuncName)
		}

		blockModules, err := fn(cfg, block.Tpls)
		if err != nil {
			return nil, fmt.Errorf("failed to compose modules for %q: %w", block.FuncName, err)
		}

		// Only one page per unique module if preview flag is enabled
		if preview {
			blockModules = core.FilterUniqueModules(blockModules)
		}

		modules = append(modules, blockModules)
	}

	if len(modules) == 0 {
		return nil, fmt.Errorf("no modules generated for page %q", file.Name)
	}

	return modules, nil
}

// validateModuleAlignment ensures all module arrays have the same length
func validateModuleAlignment(modules []core.Modules, pageName string) error {
	if len(modules) == 0 {
		return nil
	}

	expectedLen := len(modules[0])
	for _, mods := range modules {
		if len(mods) != expectedLen {
			return fmt.Errorf("module alignment error for page %q: expected %d modules, got %d", pageName, expectedLen, len(mods))
		}
	}

	return nil
}

// renderModules renders all modules to the writer using the template
func (t Tpl) renderModules(wr io.Writer, modules []core.Modules, file core.Page) error {
	if len(modules) == 0 {
		return nil
	}

	moduleCount := len(modules[0])
	for i := 0; i < moduleCount; i++ {
		for j, mod := range modules {
			if err := t.Execute(wr, mod[i].Tpl, mod[i]); err != nil {
				return core.NewTemplateError(
					mod[i].Tpl,
					0,
					fmt.Sprintf("failed to execute template for function %s", file.RenderBlocks[j].FuncName),
					err,
				)
			}
		}
	}

	return nil
}

// writePageFile writes the page content to a file
func writePageFile(cfg core.Config, pageName string, content []byte) error {
	pageFile := filepath.Join(cfg.OutputDir, pageName+texExtension)
	if err := os.WriteFile(pageFile, content, 0o600); err != nil {
		return core.NewFileError(pageFile, "write", err)
	}
	logger.Info("Generated page: %s", pageFile)
	return nil
}

func RootFilename(pathconfig string) string {
	filename := filepath.Base(pathconfig)
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + texExtension
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
	})

	// Choose source of templates: embedded by default, filesystem when DEV_TEMPLATES is set
	var (
		err   error
		useFS fs.FS
	)

	if os.Getenv(envDevTemplate) != "" {
		// Use on-disk templates for development override
		useFS = os.DirFS(filepath.Join("src", "shared", "templates", templateSubDir))
	} else {
		// Use embedded templates from templates.FS
		// Narrow to monthly/ subdir
		var sub fs.FS
		sub, err = fs.Sub(tmplfs.FS, templateSubDir)
		if err != nil {
			panic(fmt.Sprintf("failed to sub FS for monthly templates: %v", err))
		}
		useFS = sub
	}

	// Parse all *.tpl templates from the selected FS
	t, err = t.ParseFS(useFS, templatePattern)
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

func (t Tpl) Document(wr io.Writer, cfg core.Config) error {
	type pack struct {
		Cfg   core.Config
		Pages []core.Page
	}

	data := pack{Cfg: cfg, Pages: cfg.Pages}
	if err := t.tpl.ExecuteTemplate(wr, documentTpl, data); err != nil {
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

func Monthly(cfg core.Config, tpls []string) (core.Modules, error) {
	// Use legacy monthly generation without layout integration
	return MonthlyLegacy(cfg, tpls)
}

// MonthlyLegacy provides the original monthly generation without layout integration
func MonthlyLegacy(cfg core.Config, tpls []string) (core.Modules, error) {
	// Load tasks from CSV if available
	var tasks []core.Task
	if cfg.CSVFilePath != "" {
		reader := core.NewReader(cfg.CSVFilePath)
		var err error
		tasks, err = reader.ReadTasks()
		if err != nil {
			// Log error but continue without tasks
			return nil, fmt.Errorf("error reading tasks: %w", err)
		}
	}

	// If we have months with tasks from CSV, use only those
	if len(cfg.MonthsWithTasks) > 0 {
		modules := make(core.Modules, 0, len(cfg.MonthsWithTasks))

		for _, monthYear := range cfg.MonthsWithTasks {
			year := cal.NewYear(cfg.WeekStart, monthYear.Year, &cfg)

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

			modules = append(modules, core.Module{
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
					"Today":        cal.Day{Time: time.Now(), Cfg: &cfg},
				},
			})
		}

		return modules, nil
	} else {
		// Fallback to original behavior if no CSV data
		years := cfg.GetYears()
		totalMonths := len(years) * 12
		modules := make(core.Modules, 0, totalMonths)

		for _, yearNum := range years {
			year := cal.NewYear(cfg.WeekStart, yearNum, &cfg)

			for _, quarter := range year.Quarters {
				for _, month := range quarter.Months {
					modules = append(modules, core.Module{
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
							"Today":        cal.Day{Time: time.Now(), Cfg: &cfg},
						},
					})
				}
			}
		}

		return modules, nil
	}
}

// assignTasksToMonth assigns tasks to the appropriate days in a month
func assignTasksToMonth(month *cal.Month, tasks []core.Task) {
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
