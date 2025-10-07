// Package app provides the CLI application interface and document generation.
//
// This package contains:
//   - CLI application setup and command handling
//   - Template loading and rendering
//   - Template helper functions
//   - Document generation orchestration
//
// The package is the main entry point for the planner generation process:
//
// CLI Application:
//
//	New() creates the CLI application with flags for configuration,
//	output directory, and preview mode.
//
// Template System:
//
//	Templates are loaded from embedded files or filesystem (for development).
//	TemplateFuncs() provides custom template functions (dict, incr, dec, is).
//
// Generation Pipeline:
//  1. Load configuration
//  2. Setup output directory
//  3. Generate root document
//  4. Generate individual pages
//
// Example usage:
//
//	// Create and run the application
//	app := app.New()
//	err := app.Run([]string{
//	    "plannergen",
//	    "--config", "base.yaml",
//	    "--outdir", "generated",
//	})
//
// Template functions available in templates:
//   - dict: Create maps from key-value pairs
//   - incr: Increment integers
//   - dec: Decrement integers
//   - is: Check truthiness
//
// Environment variables:
//   - DEV_TEMPLATES: Use filesystem templates instead of embedded
//   - PLANNER_SILENT: Suppress log output
//   - PLANNER_LOG_LEVEL: Set logging level (silent/info/debug)
package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
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

// action is the main CLI action that orchestrates document generation or test coverage
func action(c *cli.Context) error {
	// Check if test coverage is requested
	if c.Bool(fTestCoverage) {
		return runTestCoverage()
	}

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

// runTestCoverage executes tests with coverage analysis and provides formatted results
func runTestCoverage() error {
	fmt.Println("🧪 Running Test Coverage Analysis")
	fmt.Println("═══════════════════════════════════════")

	// Create coverage output file
	coverageFile := "coverage.out"

	// Run tests with coverage
	cmd := exec.Command("go", "test", "-mod=vendor", "-coverprofile="+coverageFile, "-covermode=count", "./...")
	output, err := cmd.CombinedOutput()

	// Print test results
	if len(output) > 0 {
		fmt.Println("Test Results:")
		fmt.Println(string(output))
	}

	if err != nil {
		fmt.Printf("❌ Tests failed: %v\n", err)
		return err
	}

	// Check if coverage file was created
	if _, err := os.Stat(coverageFile); os.IsNotExist(err) {
		fmt.Println("⚠️  No coverage data generated")
		return nil
	}

	// Parse and display coverage report
	if err := analyzeCoverage(coverageFile); err != nil {
		fmt.Printf("⚠️  Coverage analysis failed: %v\n", err)
		return err
	}

	return nil
}

// analyzeCoverage parses the coverage file and provides a formatted report
func analyzeCoverage(coverageFile string) error {
	file, err := os.Open(coverageFile)
	if err != nil {
		return fmt.Errorf("failed to open coverage file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Maps to store coverage data
	packageCoverage := make(map[string][]float64)
	totalStatements := 0
	totalCovered := 0

	// Skip the first line (mode)
	if scanner.Scan() {
		// Skip mode line
	}

	// Parse coverage data
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		// Extract package name from file path
		filePath := parts[0]
		pathParts := strings.Split(filePath, "/")
		var packageName string
		for i, part := range pathParts {
			if strings.HasSuffix(part, ".go") {
				if i > 0 {
					packageName = pathParts[i-1]
				}
				break
			}
		}

		if packageName == "" {
			packageName = "main"
		}

		// Parse coverage percentage
		coverageStr := parts[2]
		if strings.HasSuffix(coverageStr, "%") {
			coverageStr = coverageStr[:len(coverageStr)-1]
		}

		coverage, err := strconv.ParseFloat(coverageStr, 64)
		if err != nil {
			continue
		}

		packageCoverage[packageName] = append(packageCoverage[packageName], coverage)
		totalStatements++
		if coverage > 0 {
			totalCovered++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading coverage file: %w", err)
	}

	// Calculate and display results
	fmt.Println("\n📊 Coverage Analysis Report")
	fmt.Println("═══════════════════════════════════════")

	// Package breakdown
	fmt.Println("Package Coverage:")
	for pkg, coverages := range packageCoverage {
		if len(coverages) == 0 {
			continue
		}

		// Calculate average coverage for package
		sum := 0.0
		for _, cov := range coverages {
			sum += cov
		}
		avgCoverage := sum / float64(len(coverages))

		status := "❌"
		if avgCoverage >= 80 {
			status = "✅"
		} else if avgCoverage >= 60 {
			status = "⚠️ "
		}

		fmt.Printf("  %s %-20s %.1f%% (%d files)\n", status, pkg, avgCoverage, len(coverages))
	}

	// Overall statistics
	overallCoverage := 0.0
	if totalStatements > 0 {
		overallCoverage = float64(totalCovered) / float64(totalStatements) * 100
	}

	fmt.Printf("\nOverall Coverage: %.1f%%\n", overallCoverage)
	fmt.Printf("Files Analyzed: %d\n", len(packageCoverage))

	// Provide recommendations
	fmt.Println("\n💡 Recommendations:")
	if overallCoverage < 60 {
		fmt.Println("  • Coverage is low - consider adding more tests")
		fmt.Println("  • Focus on testing critical business logic")
	}
	if overallCoverage >= 80 {
		fmt.Println("  • Excellent coverage! Keep up the good work")
	} else {
		fmt.Println("  • Aim for 80%+ coverage for better reliability")
		fmt.Println("  • Add tests for error conditions and edge cases")
	}

	fmt.Println("═══════════════════════════════════════")

	return nil
}

// loadConfiguration loads and validates the configuration from CLI context
func loadConfiguration(c *cli.Context) (core.Config, []string, error) {
	initialPathConfigs := strings.Split(c.Path(fConfig), ",")

	// Auto-detect CSV and adjust configuration accordingly
	csvPath := c.String("PLANNER_CSV_FILE")
	if csvPath == "" {
		autoPath, err := autoDetectCSV()
		if err == nil {
			csvPath = autoPath
			// Set the CSV path for later use
			os.Setenv("PLANNER_CSV_FILE", csvPath)
			fmt.Printf("Auto-detected CSV file: %s\n", csvPath)
		}
	}

	// Auto-detect configuration based on CSV
	pathConfigs := initialPathConfigs
	if csvPath != "" && len(initialPathConfigs) == 1 && initialPathConfigs[0] == "src/core/base.yaml" {
		autoConfigs, err := autoDetectConfig(csvPath)
		if err == nil && len(autoConfigs) > 0 {
			pathConfigs = autoConfigs
			fmt.Printf("Auto-detected configuration files: %v\n", autoConfigs)
		}
	}

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
	// Create template with custom functions
	t := template.New("").Funcs(TemplateFuncs())

	// Choose source of templates: embedded by default, filesystem when DEV_TEMPLATES is set
	var (
		err   error
		useFS fs.FS
	)

	if os.Getenv(envDevTemplate) != "" {
		// Use on-disk templates for development override
		logger.Debug("Loading templates from filesystem: %s", templatePath)
		useFS = os.DirFS(filepath.Join("src", "shared", "templates", templateSubDir))
	} else {
		// Use embedded templates from templates.FS
		logger.Debug("Loading embedded templates from: %s", templateSubDir)
		// Narrow to monthly/ subdir
		var sub fs.FS
		sub, err = fs.Sub(tmplfs.FS, templateSubDir)
		if err != nil {
			panic(fmt.Sprintf("failed to access embedded templates directory '%s': %v (check that templates are properly embedded)", templateSubDir, err))
		}
		useFS = sub
	}

	// Parse all *.tpl templates from the selected FS
	t, err = t.ParseFS(useFS, templatePattern)
	if err != nil {
		// Provide detailed error message with troubleshooting hints
		if os.Getenv(envDevTemplate) != "" {
			panic(fmt.Sprintf("failed to parse templates from filesystem '%s' with pattern '%s': %v\n"+
				"Check that template files exist and have valid syntax", templatePath, templatePattern, err))
		} else {
			panic(fmt.Sprintf("failed to parse embedded templates with pattern '%s': %v\n"+
				"This may indicate a build issue - ensure templates are embedded correctly", templatePattern, err))
		}
	}

	logger.Debug("Successfully loaded templates with pattern: %s", templatePattern)
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
		return core.NewTemplateError(documentTpl, 0, "failed to execute document template", err)
	}

	return nil
}

func (t Tpl) Execute(wr io.Writer, name string, data interface{}) error {
	// Check if template exists before trying to execute
	if t.tpl.Lookup(name) == nil {
		availableTemplates := make([]string, 0)
		for _, tmpl := range t.tpl.Templates() {
			availableTemplates = append(availableTemplates, tmpl.Name())
		}
		return core.NewTemplateError(
			name,
			0,
			fmt.Sprintf("template not found (available: %v)", availableTemplates),
			nil,
		)
	}

	if err := t.tpl.ExecuteTemplate(wr, name, data); err != nil {
		return core.NewTemplateError(name, 0, "failed to execute template", err)
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
	csvPath := cfg.CSVFilePath

	if csvPath != "" {
		reader := core.NewReader(csvPath)
		var err error
		tasks, err = reader.ReadTasks()
		if err != nil {
			// Log error but continue without tasks
			return nil, fmt.Errorf("error reading tasks: %w", err)
		}
	}

	// If we have months with tasks from CSV, use only those
	if len(cfg.MonthsWithTasks) > 0 {
		var modules core.Modules
		if len(tasks) > 0 {
			tocModule := createTableOfContentsModule(cfg, tasks, tpls[0])
			modules = append(modules, tocModule)
		}

		monthModules := make(core.Modules, 0, len(cfg.MonthsWithTasks))

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

			monthModules = append(monthModules, core.Module{
				Cfg: cfg,
				Tpl: tpls[0],
				Body: map[string]interface{}{
					"Year":         year,
					"Quarter":      targetMonth.Quarter,
					"Month":        targetMonth,
					"MonthRef":     fmt.Sprintf("month-%d-%d", targetMonth.Year.Number, int(targetMonth.Month)),
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

		// Combine TOC modules with month modules
		modules = append(modules, monthModules...)
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
					"MonthRef":     fmt.Sprintf("month-%d-%d", month.Year.Number, int(month.Month)),
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

// hexToRGBString converts a hex color string to RGB format for LaTeX
func hexToRGBString(hex string) string {
	if len(hex) < 7 || hex[0] != '#' {
		return "0,0,0" // Default black
	}

	// Parse hex values
	r, err1 := strconv.ParseInt(hex[1:3], 16, 64)
	g, err2 := strconv.ParseInt(hex[3:5], 16, 64)
	b, err3 := strconv.ParseInt(hex[5:7], 16, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return "0,0,0" // Default black on error
	}

	return fmt.Sprintf("%d,%d,%d", r, g, b)
}

// autoDetectCSV automatically finds the most appropriate CSV file in the input_data directory
func autoDetectCSV() (string, error) {
	inputDir := "input_data"

	// Check if input_data directory exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return "", fmt.Errorf("input_data directory not found")
	}

	// Find all CSV files
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return "", fmt.Errorf("failed to read input_data directory: %w", err)
	}

	var csvFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			csvFiles = append(csvFiles, file)
		}
	}

	if len(csvFiles) == 0 {
		return "", fmt.Errorf("no CSV files found in input_data directory")
	}

	// If only one CSV file, use it
	if len(csvFiles) == 1 {
		return filepath.Join(inputDir, csvFiles[0].Name()), nil
	}

	// Multiple CSV files - use priority selection
	// Priority: comprehensive > numbered versions > others
	var bestFile os.DirEntry
	bestPriority := 0

	for _, file := range csvFiles {
		name := strings.ToLower(file.Name())
		priority := 0

		// Highest priority: comprehensive files
		if strings.Contains(name, "comprehensive") {
			priority = 10
		}

		// Versioned files get priority based on version number
		if strings.Contains(name, "v") && strings.Contains(name, ".") {
			// Extract version numbers (simple heuristic)
			if strings.Contains(name, "v5.1") {
				priority = 8
			} else if strings.Contains(name, "v5") {
				priority = 6
			}
		}

		// Most recent modification time as tiebreaker
		if priority > bestPriority ||
		   (priority == bestPriority && bestFile == nil) {
			bestPriority = priority
			bestFile = file
		} else if priority == bestPriority && bestFile != nil {
			// Compare modification times
			currentInfo, err1 := file.Info()
			bestInfo, err2 := bestFile.Info()
			if err1 == nil && err2 == nil && currentInfo.ModTime().After(bestInfo.ModTime()) {
				bestFile = file
			}
		}
	}

	if bestFile != nil {
		return filepath.Join(inputDir, bestFile.Name()), nil
	}

	// Fallback to first file
	return filepath.Join(inputDir, csvFiles[0].Name()), nil
}

// autoDetectConfig automatically determines appropriate configuration files based on CSV content
func autoDetectConfig(csvPath string) ([]string, error) {
	// Read first few lines to detect version/format
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV for config detection: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for i := 0; i < 5 && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read CSV for config detection: %w", err)
	}

	// Default configuration
	baseConfig := "src/core/base.yaml"

	// Detect CSV version from filename or content
	csvName := strings.ToLower(filepath.Base(csvPath))

	if strings.Contains(csvName, "v5.1") {
		// v5.1 format - use monthly calendar config
		return []string{baseConfig, "src/core/monthly_calendar.yaml"}, nil
	} else if strings.Contains(csvName, "v5") {
		// v5 format - use basic calendar config
		return []string{baseConfig, "src/core/calendar.yaml"}, nil
	}

	// Check content for version detection
	if len(lines) > 0 {
		header := strings.ToLower(lines[0])
		if strings.Contains(header, "phase") && strings.Contains(header, "sub-phase") {
			// Has phase and sub-phase columns - v5.1 format
			return []string{baseConfig, "src/core/monthly_calendar.yaml"}, nil
		}
	}

	// Default to basic configuration
	return []string{baseConfig}, nil
}

// createTableOfContentsModule creates a table of contents module with links to all tasks
func createTableOfContentsModule(cfg core.Config, tasks []core.Task, templateName string) core.Module {

	// Generate LaTeX content directly for the TOC
	var latexContent strings.Builder

	latexContent.WriteString("% Table of Contents - Clickable Task Index\n")
	latexContent.WriteString("\\hypertarget{task-index}{}\n")
	latexContent.WriteString("{\\Large\\textbf{Task Index}}\n\n")
	latexContent.WriteString("\\vspace{0.2cm}\n\n")

	// Group tasks by phase
	phaseTasks := make(map[string][]core.Task)
	phaseNames := map[string]string{
		"1": "Phase 1: Proposal \\& Setup",
		"2": "Phase 2: Research \\& Data Collection",
		"3": "Phase 3: Publications",
		"4": "Phase 4: Dissertation",
	}

	for _, task := range tasks {
		phaseTasks[task.Phase] = append(phaseTasks[task.Phase], task)
	}

		// Create phase-based sections
		phases := []string{"1", "2", "3", "4"}
		for _, phase := range phases {
			if tasks, exists := phaseTasks[phase]; exists && len(tasks) > 0 {
				// Phase header
				latexContent.WriteString(fmt.Sprintf("{\\colorbox[RGB]{245,245,245}{\\makebox[\\linewidth][l]{\\textbf{%s}}}}\\\\\n", phaseNames[phase]))
				latexContent.WriteString("\\vspace{0.05cm}\n\n")


				latexContent.WriteString("\\vspace{0.1cm}\n\n")

				// Group tasks by sub-phase within this phase
				subPhaseTasks := make(map[string][]core.Task)
				for _, task := range tasks {
					subPhase := task.SubPhase
					if subPhase == "" {
						subPhase = "General" // Default for tasks without sub-phase
					}
					subPhaseTasks[subPhase] = append(subPhaseTasks[subPhase], task)
				}

				// Sort sub-phases alphabetically for consistent ordering
				var subPhases []string
				for subPhase := range subPhaseTasks {
					subPhases = append(subPhases, subPhase)
				}
				sort.Strings(subPhases)

				// Render each sub-phase
				for _, subPhase := range subPhases {
					subPhaseTaskList := subPhaseTasks[subPhase]
					
					// Sort tasks chronologically within this sub-phase
					sort.Slice(subPhaseTaskList, func(i, j int) bool {
						return subPhaseTaskList[i].StartDate.Before(subPhaseTaskList[j].StartDate)
					})
					
					// Sub-phase header (smaller than phase header)
					latexContent.WriteString(fmt.Sprintf("\\vspace{0.1cm}\n"))
					latexContent.WriteString(fmt.Sprintf("{\\colorbox[RGB]{250,250,250}{\\makebox[\\linewidth][l]{\\textbf{\\small %s}}}}\\\\\n", subPhase))
					latexContent.WriteString("\\vspace{0.03cm}\n\n")

					// Tasks for this sub-phase - compact format
					latexContent.WriteString("\\begin{itemize}[leftmargin=0.5cm,itemsep=0.1cm,parsep=0.05cm]\n")
					for _, task := range subPhaseTaskList {
						// Create hyperlink reference to first occurrence of task
						// Use RFC3339 format to match calendar's d.ref() method
						dateRef := task.StartDate.Format(time.RFC3339)

						// Get color for the task category
						taskColor := core.GenerateCategoryColor(strings.ToUpper(task.Category))
						taskName := strings.ReplaceAll(task.Name, "&", "\\&")
						taskName = strings.ReplaceAll(taskName, "%", "\\%")

						// Bold the task name if it's a milestone
						if task.IsMilestone {
							taskName = fmt.Sprintf("\\textbf{%s}", taskName)
						}

						// Format task with color and hyperlink
						if len(taskColor) >= 7 && taskColor[0] == '#' {
							rgbStr := hexToRGBString(taskColor)
							latexContent.WriteString(fmt.Sprintf("\\item \\textcolor[RGB]{%s}{\\hyperlink{%s}{%s}}\n", rgbStr, dateRef, taskName))
						} else {
							latexContent.WriteString(fmt.Sprintf("\\item \\hyperlink{%s}{%s}\n", dateRef, taskName))
						}
					}
					latexContent.WriteString("\\end{itemize}\n")
					latexContent.WriteString("\\vspace{0.1cm}\n\n")
				}

				latexContent.WriteString("\\vspace{0.2cm}\n\n")
			}
		}

	latexContent.WriteString("% Usage Legend\n")
	latexContent.WriteString("{\\small\n")
	latexContent.WriteString("\\textbf{How to use this index:}\\\\\n")
	latexContent.WriteString("\\textbullet\\ \\textbf{Bold task names} indicate milestones with enhanced borders in timeline\\\\\n")
	latexContent.WriteString("\\textbullet\\ Click on any task name to jump to its location in the timeline\n")
	latexContent.WriteString("}\n\n")
	latexContent.WriteString("\\pagebreak\n")

	return core.Module{
		Cfg: cfg,
		Tpl: templateName,
		Body: map[string]interface{}{
			"TOCContent": latexContent.String(),
		},
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
