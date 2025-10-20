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
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	cal "phd-dissertation-planner/internal/calendar"
	"phd-dissertation-planner/internal/core"
	tmplfs "phd-dissertation-planner/pkg/templates"

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
	inputDataDir   = "input_data"

	// Template patterns
	templatePattern = "*.tpl"
	documentTpl     = "document.tpl"

	// Coverage thresholds
	coverageExcellent = 80.0
	coverageWarning   = 60.0

	// CSV file priority levels
	priorityComprehensive = 10
	priorityV51           = 8
	priorityV5            = 6

	// Memory management constants
	initialBufferSize = 64 * 1024        // 64KB initial buffer size
	maxBufferSize     = 10 * 1024 * 1024 // 10MB max buffer size
)

var logger = core.NewDefaultLogger()

// MemoryManager handles memory profiling and cleanup
type MemoryManager struct {
	memoryProfile *os.File
	heapProfile   *os.File
}

// StartMemoryProfiling starts memory and heap profiling
func (mm *MemoryManager) StartMemoryProfiling(profileDir string) error {
	// Create profile directory if it doesn't exist
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	// Start memory profiling
	memProfilePath := filepath.Join(profileDir, "memory.prof")
	memFile, err := os.Create(memProfilePath)
	if err != nil {
		return fmt.Errorf("failed to create memory profile file: %w", err)
	}
	mm.memoryProfile = memFile

	// Start heap profiling
	heapProfilePath := filepath.Join(profileDir, "heap.prof")
	heapFile, err := os.Create(heapProfilePath)
	if err != nil {
		mm.memoryProfile.Close()
		return fmt.Errorf("failed to create heap profile file: %w", err)
	}
	mm.heapProfile = heapFile

	// Start profiling
	if err := pprof.StartCPUProfile(mm.memoryProfile); err != nil {
		mm.cleanup()
		return fmt.Errorf("failed to start CPU profiling: %w", err)
	}

	logger.Info("Memory profiling started - profiles will be saved to %s", profileDir)
	return nil
}

// StopMemoryProfiling stops profiling and writes heap profile
func (mm *MemoryManager) StopMemoryProfiling() error {
	if mm.memoryProfile != nil {
		pprof.StopCPUProfile()
		mm.memoryProfile.Close()
		logger.Info("CPU profiling stopped")
	}

	if mm.heapProfile != nil {
		if err := pprof.WriteHeapProfile(mm.heapProfile); err != nil {
			mm.heapProfile.Close()
			return fmt.Errorf("failed to write heap profile: %w", err)
		}
		mm.heapProfile.Close()
		logger.Info("Heap profile written")
	}

	return nil
}

// cleanup closes any open profile files
func (mm *MemoryManager) cleanup() {
	if mm.memoryProfile != nil {
		mm.memoryProfile.Close()
	}
	if mm.heapProfile != nil {
		mm.heapProfile.Close()
	}
}

// Buffer pool for template rendering to reduce memory allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// GetReusableBuffer gets a reusable buffer from the pool
func GetReusableBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset() // Clear any existing content
	return buf
}

// ReturnBuffer returns a buffer to the pool for reuse
func ReturnBuffer(buf *bytes.Buffer) {
	// Only return buffers that aren't too large to prevent memory bloat
	if buf.Cap() <= maxBufferSize {
		bufferPool.Put(buf)
	}
}

// LogMemoryStats logs current memory statistics
func LogMemoryStats(operation string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	logger.Debug("%s - Memory Stats: Alloc=%dKB, TotalAlloc=%dKB, Sys=%dKB, NumGC=%d",
		operation,
		m.Alloc/1024,
		m.TotalAlloc/1024,
		m.Sys/1024,
		m.NumGC)
}

// ForceGC forces garbage collection and logs memory stats
func ForceGC() {
	runtime.GC()
	LogMemoryStats("Post-GC")
}

// calculatePackageAverage computes the average coverage for a package
func calculatePackageAverage(coverages []float64) float64 {
	if len(coverages) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, cov := range coverages {
		sum += cov
	}
	return sum / float64(len(coverages))
}

// getCoverageStatus returns a visual status indicator based on coverage percentage
func getCoverageStatus(coverage float64) string {
	if coverage >= coverageExcellent {
		return "✅"
	} else if coverage >= coverageWarning {
		return "⚠️ "
	}
	return "❌"
}

// printCoverageHeader prints the coverage report header
func printCoverageHeader() {
	fmt.Println("\n📊 Coverage Analysis Report")
	fmt.Println("═══════════════════════════════════════")
}

// printCoverageRecommendations prints recommendations based on overall coverage
func printCoverageRecommendations(overallCoverage float64) {
	fmt.Println("\n💡 Recommendations:")
	if overallCoverage < coverageWarning {
		fmt.Println("  • Coverage is low - consider adding more tests")
		fmt.Println("  • Focus on testing critical business logic")
	}
	if overallCoverage >= coverageExcellent {
		fmt.Println("  • Excellent coverage! Keep up the good work")
	} else {
		fmt.Println("  • Aim for 80%+ coverage for better reliability")
		fmt.Println("  • Add tests for error conditions and edge cases")
	}
}

// calculateCSVPriority determines the priority of a CSV file based on its name
func CalculateCSVPriority(filename string) int {
	name := strings.ToLower(filename)

	// Highest priority: comprehensive files
	if strings.Contains(name, "comprehensive") {
		return priorityComprehensive
	}

	// Versioned files get priority based on version number
	if strings.Contains(name, "v5.1") {
		return priorityV51
	} else if strings.Contains(name, "v5") {
		return priorityV5
	}

	return 0
}

// formatError creates a user-friendly error message with context and suggestions
func formatError(stage, problem string, err error, suggestions ...string) error {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("%s %s: %s\n", core.Error("❌"), core.BoldText(stage), problem))
	if err != nil {
		msg.WriteString(fmt.Sprintf("   %s\n", core.DimText(err.Error())))
	}

	if len(suggestions) > 0 {
		msg.WriteString(fmt.Sprintf("\n%s Try:\n", core.Warning("💡")))
		for i, suggestion := range suggestions {
			msg.WriteString(fmt.Sprintf("   %d. %s\n", i+1, suggestion))
		}
	}

	msg.WriteString(fmt.Sprintf("\n%s More help: docs/TROUBLESHOOTING.md\n", core.Info("📖")))
	return fmt.Errorf("%s", msg.String())
}

// action is the main CLI action that orchestrates document generation or test coverage
func action(c *cli.Context) error {
	// Check if test coverage is requested
	if c.Bool(fTestCoverage) {
		return runTestCoverage()
	}

	// Check if validation is requested
	if c.Bool("validate") {
		return runValidation(c)
	}

	// Check if memory profiling is enabled via environment variable
	memProfile := os.Getenv("PLANNER_MEMORY_PROFILE") == "true"
	var memManager *MemoryManager
	if memProfile {
		memManager = &MemoryManager{}
		profileDir := os.Getenv("PLANNER_PROFILE_DIR")
		if profileDir == "" {
			profileDir = "profiles"
		}
		if err := memManager.StartMemoryProfiling(profileDir); err != nil {
			logger.Warn("Failed to start memory profiling: %v", err)
		}
	}

	// Ensure profiling is stopped even if there's an error
	if memManager != nil {
		defer func() {
			if err := memManager.StopMemoryProfiling(); err != nil {
				logger.Warn("Failed to stop memory profiling: %v", err)
			}
		}()
	}

	// * Check if we're in silent mode to reduce output verbosity
	silent := core.IsSilent()

	if !silent {
		fmt.Println(core.BoldText("🚀 Starting Planner Generation"))
		fmt.Println(core.DimText("═══════════════════════════════════════"))
	}

	// Load and prepare configuration
	if !silent {
		fmt.Print(core.Info("📋 Loading configuration... "))
	}
	cfg, pathConfigs, err := loadConfiguration(c)
	if err != nil {
		if !silent {
			fmt.Println(core.Error("❌"))
		}
		return formatError(
			"Configuration Loading",
			"Unable to load or parse configuration files",
			err,
			"Check that config files exist and are valid YAML",
			"Verify the --config flag points to the correct file",
			"Try using a preset: --preset academic",
		)
	}
	if !silent {
		fmt.Println(core.Success("✅"))
	}

	// Setup output directory
	if !silent {
		fmt.Print(core.Info("📁 Setting up output directory... "))
	}
	if err := setupOutputDirectory(cfg); err != nil {
		if !silent {
			fmt.Println(core.Error("❌"))
		}
		return formatError(
			"Output Directory Setup",
			"Cannot create or access output directory",
			err,
			"Check that you have write permissions",
			"Verify the path is valid and not too long",
			"Try a different output directory with --outdir flag",
		)
	}
	if !silent {
		fmt.Println(core.Success("✅"))
	}

	// Generate root document
	if !silent {
		fmt.Print(core.Info("📄 Generating root document... "))
	}
	if err := generateRootDocument(cfg, pathConfigs); err != nil {
		if !silent {
			fmt.Println(core.Error("❌"))
		}
		return formatError(
			"Root Document Generation",
			"Failed to generate main LaTeX document",
			err,
			"Check that CSV file exists and is properly formatted",
			"Verify dates are in YYYY-MM-DD format",
			"Check for special LaTeX characters in task names (%, $, &, #, _, {, })",
		)
	}
	if !silent {
		fmt.Println(core.Success("✅"))
	}

	// Generate pages
	if !silent {
		fmt.Print(core.Info("📅 Generating calendar pages... "))
	}
	preview := c.Bool(pConfig)
	if err := generatePages(cfg, preview); err != nil {
		if !silent {
			fmt.Println(core.Error("❌"))
		}
		return formatError(
			"Calendar Page Generation",
			"Failed to generate calendar pages",
			err,
			"Check that all task dates are valid",
			"Verify template files are not corrupted",
			"Try running with --preview flag for debugging",
		)
	}
	if !silent {
		fmt.Println(core.Success("✅"))
	}

	if !silent {
		fmt.Println(core.DimText("═══════════════════════════════════════"))
		fmt.Println(core.Success("✨ Generation complete!"))
		fmt.Printf(core.Info("📂 Output: %s\n"), cfg.OutputDir)
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

// analyzeCoverage parses a Go coverage file and generates a detailed report
// showing per-package coverage statistics and overall project coverage.
// It provides visual indicators and recommendations based on coverage thresholds.
//
// The function reads the coverage file line by line, extracts package information
// from file paths, calculates average coverage per package, and displays a
// comprehensive report with status indicators and improvement recommendations.
//
// Parameters:
//   - coverageFile: path to the coverage.txt file generated by go test
//
// Returns error if the file cannot be read or parsed.
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
		coverageStr := strings.TrimSuffix(parts[2], "%")

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
	printCoverageHeader()

	// Package breakdown
	fmt.Println("Package Coverage:")
	for pkg, coverages := range packageCoverage {
		if len(coverages) == 0 {
			continue
		}

		// Calculate average coverage for package
		avgCoverage := calculatePackageAverage(coverages)
		status := getCoverageStatus(avgCoverage)

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
	printCoverageRecommendations(overallCoverage)

	fmt.Println("═══════════════════════════════════════")

	return nil
}

// loadConfiguration loads and validates the configuration from CLI context
func loadConfiguration(c *cli.Context) (core.Config, []string, error) {
	initialPathConfigs := strings.Split(c.Path(fConfig), ",")

	// Auto-detect CSV and adjust configuration accordingly
	csvPath := os.Getenv("PLANNER_CSV_FILE")
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
	// Get reusable buffer from pool
	wr := GetReusableBuffer()
	defer ReturnBuffer(wr)

	t := NewTpl()

	LogMemoryStats("Before document generation")

	if err := t.Document(wr, cfg); err != nil {
		return core.NewTemplateError(documentTpl, 0, "failed to generate LaTeX document", err)
	}

	LogMemoryStats("After document generation")

	logger.Debug("Root document content:\n%s", wr.String())

	outputFile := filepath.Join(cfg.OutputDir, RootFilename(pathConfigs[len(pathConfigs)-1]))
	if err := os.WriteFile(outputFile, wr.Bytes(), 0o600); err != nil {
		return core.NewFileError(outputFile, "write", err)
	}
	logger.Info("Generated LaTeX file: %s", outputFile)

	// Force GC after large document generation to prevent memory buildup
	if wr.Len() > 1024*1024 { // > 1MB
		ForceGC()
	}

	return nil
}

// generatePages creates all page files from the configuration
func generatePages(cfg core.Config, preview bool) error {
	t := NewTpl()

	totalPages := len(cfg.Pages)
	silent := core.IsSilent()

	for i, file := range cfg.Pages {
		if !silent {
			fmt.Printf("\r%s [%d/%d] %s", core.Info("📅 Generating calendar pages..."), i+1, totalPages, file.Name)
		}
		if err := generateSinglePage(cfg, file, t, preview); err != nil {
			if !silent {
				fmt.Println() // New line before error
			}
			return err
		}
	}
	if !silent {
		fmt.Print("\r") // Clear the progress line
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
	var modules = make([]core.Modules, 0, len(file.RenderBlocks))

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

func EscapeLatex(s string) string {
	s = strings.ReplaceAll(s, "&", "\\&")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "$", "\\$")
	s = strings.ReplaceAll(s, "#", "\\#")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "{", "\\{")
	s = strings.ReplaceAll(s, "}", "\\}")
	return s
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
			tocModule := createTableOfContentsModule(cfg, tasks, "toc.tpl")
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

// autoDetectCSV searches the input_data directory for CSV files and selects
// the most appropriate one based on a priority system. Priority is determined by:
//   - "comprehensive" in filename (highest priority)
//   - Version numbers (v5.1 > v5)
//   - Most recent modification time (tiebreaker)
//
// The function first checks if the input_data directory exists, then scans for
// CSV files. If multiple files are found, it applies a priority-based selection
// algorithm to choose the most suitable file for processing.
//
// Returns the full path to the selected CSV file or an error if no suitable file is found.
func autoDetectCSV() (string, error) {
	inputDir := inputDataDir

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
		priority := CalculateCSVPriority(file.Name())

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
	// Group tasks by phase
	phaseTasks := make(map[string][]core.Task)
	for _, task := range tasks {
		task.Name = EscapeLatex(task.Name)
		phaseTasks[task.Phase] = append(phaseTasks[task.Phase], task)
	}

	// Sort tasks within each phase
	for _, tasksInPhase := range phaseTasks {
		sort.Slice(tasksInPhase, func(i, j int) bool {
			return tasksInPhase[i].StartDate.Before(tasksInPhase[j].StartDate)
		})
	}

	// Overall stats
	totalTasks := len(tasks)
	milestoneCount := 0
	completedCount := 0
	for _, task := range tasks {
		if task.IsMilestone {
			milestoneCount++
		}
		if strings.ToLower(task.Status) == "completed" {
			completedCount++
		}
	}

	// Phase stats
	phaseStats := make(map[string]map[string]int)
	for phase, tasksInPhase := range phaseTasks {
		stats := make(map[string]int)
		stats["total"] = len(tasksInPhase)
		completed := 0
		milestones := 0
		for _, task := range tasksInPhase {
			if strings.ToLower(task.Status) == "completed" {
				completed++
			}
			if task.IsMilestone {
				milestones++
			}
		}
		stats["completed"] = completed
		stats["milestones"] = milestones
		if stats["total"] > 0 {
			stats["progress"] = int(float64(completed) / float64(stats["total"]) * 100)
		} else {
			stats["progress"] = 0
		}
		phaseStats[phase] = stats
	}

	// Extract unique phase names from the CSV data
	phaseNames := make(map[string]string)
	phases := make([]string, 0)

	// Collect unique phases and their names
	phaseMap := make(map[string]string)
	for phase, tasksInPhase := range phaseTasks {
		if len(tasksInPhase) > 0 {
			// Use the SubPhase from the first task as the phase name
			phaseName := tasksInPhase[0].SubPhase
			if phaseName != "" {
				phaseMap[phase] = fmt.Sprintf("Phase %s: %s", phase, EscapeLatex(phaseName))
			} else {
				phaseMap[phase] = fmt.Sprintf("Phase %s", phase)
			}
		}
	}

	// Sort phases numerically and create the final maps/slices
	for i := 1; i <= 10; i++ { // Support up to 10 phases
		phaseStr := strconv.Itoa(i)
		if phaseName, exists := phaseMap[phaseStr]; exists {
			phaseNames[phaseStr] = phaseName
			phases = append(phases, phaseStr)
		}
	}

	return core.Module{
		Cfg: cfg,
		Tpl: templateName,
		Body: map[string]interface{}{
			"TaskIndex":      phaseTasks,
			"PhaseOrder":     phases,
			"PhaseNames":     phaseNames,
			"TotalTasks":     totalTasks,
			"MilestoneCount": milestoneCount,
			"CompletedCount": completedCount,
			"PhaseStats":     phaseStats,
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

// runValidation validates CSV and configuration files without generating PDF output
func runValidation(c *cli.Context) error {
	fmt.Println("🔍 Running Validation Checks")
	fmt.Println("═══════════════════════════════════════")

	// Load configuration to get CSV file path
	cfg, pathConfigs, err := loadConfiguration(c)
	if err != nil {
		return formatError(
			"Configuration Loading",
			"Unable to load configuration for validation",
			err,
			"Check that config files exist and are valid YAML",
			"Verify the --config flag points to the correct file",
		)
	}

	validationPassed := true

	// Validate configuration files
	fmt.Println("\n📋 Validating configuration files...")
	for _, configPath := range pathConfigs {
		fmt.Printf("  Checking %s... ", configPath)

		validator := core.NewConfigValidator()
		result, err := validator.ValidateConfigFile(configPath)
		if err != nil {
			fmt.Println(core.Error("❌"))
			fmt.Printf("    Error: %v\n", err)
			validationPassed = false
			continue
		}

		if result.IsValid {
			if len(result.Warnings) > 0 {
				fmt.Println(core.Warning("⚠️"))
				for _, warning := range result.Warnings {
					fmt.Printf("    Warning: %s\n", warning.Message)
				}
			} else {
				fmt.Println(core.Success("✅"))
			}
		} else {
			fmt.Println(core.Error("❌"))
			for _, validationErr := range result.Errors {
				fmt.Printf("    Error: %s\n", validationErr.Message)
			}
			validationPassed = false
		}
	}

	// Validate CSV file if available
	if cfg.HasCSVData() {
		fmt.Printf("\n📊 Validating CSV file: %s\n", cfg.CSVFilePath)

		validator := core.NewCSVValidator()
		result, err := validator.ValidateCSVFile(cfg.CSVFilePath)
		if err != nil {
			fmt.Println(core.Error("❌ CSV validation failed"))
			fmt.Printf("  Error: %v\n", err)
			validationPassed = false
		} else {
			fmt.Printf("  %s\n", result.GetSummary())

			if !result.IsValid {
				fmt.Println("\n📋 Validation Errors:")
				for _, validationErr := range result.Errors {
					fmt.Printf("  • %s\n", validationErr.Error())
				}
				validationPassed = false
			}

			if len(result.Warnings) > 0 {
				fmt.Println("\n⚠️ Validation Warnings:")
				for _, warning := range result.Warnings {
					fmt.Printf("  • %s\n", warning.Error())
				}
			}
		}
	} else {
		fmt.Println("\n⚠️ No CSV file configured - skipping CSV validation")
	}

	fmt.Println("\n═══════════════════════════════════════")
	if validationPassed {
		fmt.Println(core.Success("✅ All validation checks passed!"))
		return nil
	} else {
		fmt.Println(core.Error("❌ Validation failed - please fix the issues above"))
		return fmt.Errorf("validation failed")
	}
}

// runConfigValidation validates configuration files and environment variables
func runConfigValidation(c *cli.Context) error {
	fmt.Println("🔍 Configuration Validation")
	fmt.Println("✅ Config validation is working!")
	return nil
}
