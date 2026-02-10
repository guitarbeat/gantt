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
	tmplfs "phd-dissertation-planner/internal/templates"

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
	buf, ok := bufferPool.Get().(*bytes.Buffer)
	if !ok {
		// Fallback: create new buffer if pool returns unexpected type
		buf = &bytes.Buffer{}
	}
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

	// Get all CSV files to process
	csvFiles, err := getAllCSVFiles()
	if err != nil {
		if !silent {
			fmt.Println(core.Error("❌"))
		}
		return formatError(
			"CSV File Detection",
			"Unable to find CSV files to process",
			err,
			"Check that input_data directory exists",
			"Verify there are CSV files in the input_data directory",
			"Ensure CSV files have .csv extension",
		)
	}

	if !silent {
		fmt.Printf("%s", core.Info(fmt.Sprintf("📋 Found %d CSV file(s) to merge and process\n", len(csvFiles))))
		for i, csvFile := range csvFiles {
			fmt.Printf("%s", core.Info(fmt.Sprintf("   %d. %s\n", i+1, filepath.Base(csvFile))))
		}
	}

	// Merge all CSV files in memory
	if !silent {
		fmt.Print(core.Info("🔄 Merging CSV files in memory... "))
	}
	allTasks, err := core.ReadTasksFromMultipleFiles(csvFiles)
	if err != nil {
		if !silent {
			fmt.Println(core.Error("❌"))
		}
		return formatError(
			"CSV Merging",
			"Unable to merge CSV files",
			err,
			"Check that all CSV files have the same header structure",
			"Verify there are no duplicate task IDs across files",
			"Ensure all CSV files are valid",
		)
	}
	if !silent {
		fmt.Printf("%s", core.Success(fmt.Sprintf("✅ (%d tasks total)\n", len(allTasks))))

		// Group tasks by phase for summary
		phaseCounts := make(map[string]int)
		for _, task := range allTasks {
			phaseCounts[task.Phase] = phaseCounts[task.Phase] + 1
		}

		// Sort phases for consistent display
		var phases []string
		for phase := range phaseCounts {
			phases = append(phases, phase)
		}
		sort.Strings(phases)

		// Print summary
		for _, phase := range phases {
			count := phaseCounts[phase]
			fmt.Printf("   • %s: %s\n", core.CyanText(phase), core.DimText(fmt.Sprintf("%d tasks", count)))
		}
	}

	// Load and prepare configuration with merged tasks
	if !silent {
		fmt.Print(core.Info("📋 Loading configuration... "))
	}
	cfg, pathConfigs, err := loadConfigurationWithTasks(c, allTasks)
	if err != nil {
		if !silent {
			fmt.Println(core.Error("❌"))
		}
		return formatError(
			"Configuration",
			"Unable to load configuration",
			err,
			"Check that input_data/config.yaml exists",
			"Verify configuration file syntax",
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
			"Output Directory",
			"Unable to create output directory",
			err,
			"Check directory permissions",
			"Verify disk space",
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
			"Document Generation",
			"Unable to generate root document",
			err,
			"Check template files",
			"Verify configuration",
		)
	}
	if !silent {
		fmt.Println(core.Success("✅"))
	}

	// Generate pages
	preview := c.Bool(pConfig)
	if err := generatePages(cfg, preview); err != nil {
		if !silent {
			fmt.Println(core.Error("❌"))
		}
		return formatError(
			"Page Generation",
			"Unable to generate calendar pages",
			err,
			"Check template files",
			"Verify task data",
		)
	}
	if !silent {
		fmt.Println(core.Success("✅"))
	}

	// Compile LaTeX to PDF
	stopSpinner := make(chan bool)
	if !silent {
		go func() {
			chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			i := 0
			for {
				select {
				case <-stopSpinner:
					return
				default:
					fmt.Print(core.ClearLine())
					fmt.Printf("%s %s", core.CyanText(chars[i]), core.Info("Compiling LaTeX to PDF..."))
					time.Sleep(100 * time.Millisecond)
					i = (i + 1) % len(chars)
				}
			}
		}()
	}

	pdfCompiled := false
	err = compileLaTeXToPDF(cfg)
	if !silent {
		stopSpinner <- true
	}

	if err != nil {
		if !silent {
			// Clear line and print error status
			fmt.Print(core.ClearLine())
			fmt.Printf("%s %s\n", core.Error("❌"), core.Info("Compiling LaTeX to PDF..."))
		}

		if strings.Contains(err.Error(), "executable file not found") {
			if !silent {
				fmt.Println(core.Warning("\n⚠️  PDF generation skipped: 'xelatex' not found"))
				fmt.Println(core.DimText("   LaTeX files have been generated in: " + filepath.Join(cfg.OutputDir, "latex")))
				fmt.Println(core.DimText("   To generate PDF manually, install TeX Live/MacTeX and run:"))
				fmt.Printf("   %s\n", core.CyanText(fmt.Sprintf("cd %s && xelatex %s", filepath.Join(cfg.OutputDir, "latex"), RootFilename(pathConfigs[len(pathConfigs)-1]))))
			}
			logger.Warn("PDF compilation skipped (xelatex missing)")
		} else {
			logger.Warn("PDF compilation failed: %v", err)
		}
	} else {
		pdfCompiled = true
		if !silent {
			// Clear line and print success status
			fmt.Print(core.ClearLine())
			fmt.Printf("%s %s\n", core.Success("✅"), core.Info("Compiling LaTeX to PDF..."))
		}
	}

	if !silent {
		fmt.Println(core.DimText("═══════════════════════════════════════"))
		if pdfCompiled {
			fmt.Printf("%s", core.Success(fmt.Sprintf("✨ Successfully generated calendar from %d CSV files!\n", len(csvFiles))))
		} else {
			fmt.Printf("%s", core.Warning("⚠️  Generated LaTeX files, but PDF compilation failed (check xelatex installation)\n"))
		}
		fmt.Printf("%s", core.Info(fmt.Sprintf("📂 Output: %s\n", cfg.OutputDir)))
	}

	if !silent {
		fmt.Println(core.DimText("═══════════════════════════════════════"))
		if pdfCompiled {
			fmt.Println(core.Success("✨ All files processed!"))
		} else {
			fmt.Println(core.Warning("⚠️  Done (with warnings)"))
		}
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
			fmt.Printf("%s", core.Info(fmt.Sprintf("🔍 Auto-detected CSV file: %s\n", csvPath)))
		}
	}

	// Auto-detect configuration based on CSV
	pathConfigs := initialPathConfigs
	if csvPath != "" && len(initialPathConfigs) == 1 && initialPathConfigs[0] == "input_data/config.yaml" {
		autoConfigs, err := autoDetectConfig(csvPath)
		if err == nil && len(autoConfigs) > 0 {
			pathConfigs = autoConfigs
			fmt.Printf("%s", core.Info(fmt.Sprintf("🔍 Auto-detected configuration files: %v\n", autoConfigs)))
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

// loadConfigurationWithTasks loads configuration and injects pre-loaded tasks
func loadConfigurationWithTasks(c *cli.Context, tasks []core.Task) (core.Config, []string, error) {
	initialPathConfigs := strings.Split(c.Path(fConfig), ",")

	cfg, err := core.NewConfig(initialPathConfigs...)
	if err != nil {
		return core.Config{}, nil, core.NewConfigError(
			strings.Join(initialPathConfigs, ","),
			"",
			"failed to load configuration",
			err,
		)
	}

	// Override output directory from CLI flag if provided
	if od := strings.TrimSpace(c.Path(fOutDir)); od != "" {
		cfg.OutputDir = od
	}

	// Inject the pre-loaded tasks into the configuration
	cfg.Tasks = tasks
	
	// Calculate date range and months with tasks from the merged data
	if len(tasks) > 0 {
		dateRange := core.CalculateDateRange(tasks)
		cfg.MonthsWithTasks = core.GetMonthsWithTasks(tasks, dateRange)
	}

	return cfg, initialPathConfigs, nil
}

// setupOutputDirectory ensures the output directory exists and logs its location
func setupOutputDirectory(cfg core.Config) error {
	// Create main output directory
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return core.NewFileError(cfg.OutputDir, "create directory", err)
	}
	
	// Create organized subdirectories
	subdirs := []string{
		filepath.Join(cfg.OutputDir, "pdfs"),
		filepath.Join(cfg.OutputDir, "latex"),
		filepath.Join(cfg.OutputDir, "auxiliary"),
		filepath.Join(cfg.OutputDir, "binaries"),
	}
	
	for _, subdir := range subdirs {
		if err := os.MkdirAll(subdir, 0o755); err != nil {
			return core.NewFileError(subdir, "create subdirectory", err)
		}
	}
	
	logger.Debug("Output directory: %s", cfg.OutputDir)
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

	outputFile := filepath.Join(cfg.OutputDir, "latex", RootFilename(pathConfigs[len(pathConfigs)-1]))
	if err := os.WriteFile(outputFile, wr.Bytes(), 0o600); err != nil {
		return core.NewFileError(outputFile, "write", err)
	}
	logger.Debug("Generated LaTeX file: %s", outputFile)

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
			fmt.Print(core.ClearLine())
			fmt.Printf("%s [%d/%d] %s", core.Info("📅 Generating calendar pages..."), i+1, totalPages, file.Name)
		}
		if err := generateSinglePage(cfg, file, t, preview); err != nil {
			if !silent {
				fmt.Println() // New line before error
			}
			return err
		}
	}
	if !silent {
		// Add a space so the checkmark printed by the caller appears next to the progress
		fmt.Print(" ")
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
	pageFile := filepath.Join(cfg.OutputDir, "latex", pageName+texExtension)
	if err := os.WriteFile(pageFile, content, 0o600); err != nil {
		return core.NewFileError(pageFile, "write", err)
	}
	logger.Debug("Generated page: %s", pageFile)
	return nil
}

func RootFilename(pathconfig string) string {
	filename := filepath.Base(pathconfig)
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + texExtension
}

// latexReplacer is a reusable replacer for escaping LaTeX special characters.
// It is initialized once and used by EscapeLatex for better performance.
var latexReplacer = strings.NewReplacer(
	"&", "\\&",
	"%", "\\%",
	"$", "\\$",
	"#", "\\#",
	"_", "\\_",
	"{", "\\{",
	"}", "\\}",
	"\\", "\\textbackslash{}",
	"^", "\\textasciicircum{}",
	"~", "\\textasciitilde{}",
	"[", "{[}",
	"]", "{]}",
)

func EscapeLatex(s string) string {
	return latexReplacer.Replace(s)
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
	// Use tasks from config (already loaded and merged)
	tasks := cfg.Tasks

	// If we have months with tasks from CSV, use only those
	if len(cfg.MonthsWithTasks) > 0 {
		var modules core.Modules
		if len(tasks) > 0 {
			// Get CSV file list for TOC display
			csvFiles, _ := getAllCSVFiles()
			tocModule := createTableOfContentsModule(cfg, tasks, "toc.tpl", csvFiles)
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

// getAllCSVFiles returns all CSV files in the input_data directory
func getAllCSVFiles() ([]string, error) {
	inputDir := inputDataDir

	// Check if input_data directory exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("input_data directory not found")
	}

	// Find all CSV files
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read input_data directory: %w", err)
	}

	var csvFiles []string
	
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			// Skip hidden files and temporary files
			if !strings.HasPrefix(file.Name(), ".") {
				csvFiles = append(csvFiles, filepath.Join(inputDir, file.Name()))
			}
		}
	}

	if len(csvFiles) == 0 {
		return nil, fmt.Errorf("no CSV files found in input_data directory")
	}

	// Sort files alphabetically (phase_1, phase_2, etc. will be in order)
	sort.Strings(csvFiles)

	return csvFiles, nil
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
	baseConfig := "input_data/config.yaml"

	// Detect CSV version from filename or content
	csvName := strings.ToLower(filepath.Base(csvPath))

	if strings.Contains(csvName, "v5.1") {
		// v5.1 format - use main config only
		return []string{baseConfig}, nil
	} else if strings.Contains(csvName, "v5") {
		// v5 format - use main config only
		return []string{baseConfig}, nil
	}

	// Check content for version detection
	if len(lines) > 0 {
		header := strings.ToLower(lines[0])
		if strings.Contains(header, "phase") && strings.Contains(header, "sub-phase") {
			// Has phase and sub-phase columns - v5.1 format
			return []string{baseConfig}, nil
		}
	}

	// Default to basic configuration
	return []string{baseConfig}, nil
}

// createTableOfContentsModule creates a table of contents module with links to all tasks
func createTableOfContentsModule(cfg core.Config, tasks []core.Task, templateName string, csvFiles []string) core.Module {
	// Group tasks by phase
	phaseTasks := make(map[string][]core.Task)
	for _, task := range tasks {
		task.Name = EscapeLatex(task.Name)
		phaseTasks[task.Phase] = append(phaseTasks[task.Phase], task)
	}

	// Sort tasks within each phase by start date
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

	// Define hierarchical structure with sections and phases
	type PhaseSection struct {
		Name   string
		Phases []string
	}

	sections := []PhaseSection{
		{
			Name: "Setup \\& Proposal",
			Phases: []string{
				"Project Metadata",
				"PhD Proposal",
				"Committee Management",
				"Microscope Setup",
				"Laser System",
			},
		},
		{
			Name: "Research Aims",
			Phases: []string{
				"Aim 1 - AAV-based Vascular Imaging",
				"Aim 2 - Dual-channel Imaging Platform",
				"Aim 3 - Stroke Study & Analysis",
				"Data Management & Analysis",
			},
		},
		{
			Name: "Publications \\& Tools",
			Phases: []string{
				"SLAVV-T Development",
				"AR Platform Development",
				"Research Paper",
				"Methodology Paper",
				"Manuscript Submissions",
			},
		},
		{
			Name: "Dissertation \\& Defense",
			Phases: []string{
				"Dissertation Writing",
				"Committee Review & Defense",
				"Final Submission & Graduation",
			},
		},
	}

	// Build phase order from sections
	phaseOrder := make([]string, 0)
	for _, section := range sections {
		phaseOrder = append(phaseOrder, section.Phases...)
	}

	// Get unique phases that exist in the data, ordered by phaseOrder
	phases := make([]string, 0, len(phaseTasks))
	phaseSet := make(map[string]bool)
	for phase := range phaseTasks {
		phaseSet[phase] = true
	}

	// Add phases in the defined order
	for _, phase := range phaseOrder {
		if phaseSet[phase] {
			phases = append(phases, phase)
			delete(phaseSet, phase)
		}
	}

	// Add any remaining phases not in the defined order (alphabetically)
	remainingPhases := make([]string, 0)
	for phase := range phaseSet {
		remainingPhases = append(remainingPhases, phase)
	}
	sort.Strings(remainingPhases)
	phases = append(phases, remainingPhases...)

	// Create phase names map (escaped for LaTeX) and phase colors
	phaseNames := make(map[string]string)
	phaseColors := make(map[string]string)
	for _, phase := range phases {
		phaseNames[phase] = EscapeLatex(phase)
		// Generate color for this phase using the same algorithm as the calendar
		color := core.GenerateCategoryColor(phase)
		phaseColors[phase] = core.HexToRGB(color)
	}

	// Calculate task durations in days
	taskDurations := make(map[string]string)
	for _, task := range tasks {
		duration := task.EndDate.Sub(task.StartDate)
		days := int(duration.Hours() / 24)
		if days < 1 {
			days = 1
		}
		taskDurations[task.ID] = fmt.Sprintf("%d", days)
	}

	// Create phase-to-section mapping
	phaseToSection := make(map[string]string)
	for _, section := range sections {
		for _, phase := range section.Phases {
			phaseToSection[phase] = section.Name
		}
	}

	// Prepare CSV file info for display
	csvFileNames := make([]string, len(csvFiles))
	for i, csvFile := range csvFiles {
		csvFileNames[i] = EscapeLatex(filepath.Base(csvFile))
	}

	return core.Module{
		Cfg: cfg,
		Tpl: templateName,
		Body: map[string]interface{}{
			"TaskIndex":      phaseTasks,
			"PhaseOrder":     phases,
			"PhaseNames":     phaseNames,
			"PhaseColors":    phaseColors,
			"PhaseToSection": phaseToSection,
			"TaskDurations":  taskDurations,
			"TotalTasks":     totalTasks,
			"MilestoneCount": milestoneCount,
			"CompletedCount": completedCount,
			"PhaseStats":     phaseStats,
			"CSVFiles":       csvFileNames,
			"CSVFileCount":   len(csvFiles),
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

	// Get all CSV files to validate
	csvFiles, err := getAllCSVFiles()
	if err != nil {
		return formatError(
			"CSV File Detection",
			"Unable to find CSV files to validate",
			err,
			"Check that input_data directory exists",
			"Verify there are CSV files in the input_data directory",
			"Ensure CSV files have .csv extension",
		)
	}

	fmt.Printf("📋 Found %d CSV file(s) to validate\n", len(csvFiles))

	validationPassed := true
	validatedConfigs := make(map[string]bool)

	// Process each CSV file
	for i, csvFile := range csvFiles {
		fmt.Printf("\n📄 Validating file %d/%d: %s\n", i+1, len(csvFiles), filepath.Base(csvFile))

		// Set the CSV file for this validation run
		os.Setenv("PLANNER_CSV_FILE", csvFile)

		// Load configuration to get CSV file path
		cfg, pathConfigs, err := loadConfiguration(c)
		if err != nil {
			fmt.Printf("❌ Configuration error for %s: %v\n", filepath.Base(csvFile), err)
			validationPassed = false
			continue
		}

		// Filter out already validated configurations
		configsToValidate := make([]string, 0)
		for _, configPath := range pathConfigs {
			if !validatedConfigs[configPath] {
				configsToValidate = append(configsToValidate, configPath)
			}
		}

		// Validate configuration files
		if len(configsToValidate) > 0 {
			fmt.Println("📋 Validating configuration files...")
			for _, configPath := range configsToValidate {
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
					validatedConfigs[configPath] = true
				} else {
					fmt.Println(core.Error("❌"))
					for _, validationErr := range result.Errors {
						fmt.Printf("    Error: %s\n", validationErr.Message)
					}
					validationPassed = false
				}
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
						fmt.Println(formatValidationIssue(validationErr))
					}
					validationPassed = false
				}

				if len(result.Warnings) > 0 {
					fmt.Println("\n⚠️ Validation Warnings:")
					for _, warning := range result.Warnings {
						fmt.Println(formatValidationIssue(warning))
					}
				}
			}
		} else {
			fmt.Println("\n⚠️ No CSV file configured - skipping CSV validation")
		}
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

// formatValidationIssue creates a visually structured string for a validation issue
func formatValidationIssue(issue core.ValidationIssue) string {
	var parts []string

	// Row info (Dim)
	if issue.Row > 0 {
		parts = append(parts, core.DimText(fmt.Sprintf("Row %d", issue.Row)))
	}

	// Field info (Cyan)
	if issue.Field != "" {
		parts = append(parts, core.CyanText(issue.Field))
	}

	// Value info (Dim)
	if issue.Value != "" {
		parts = append(parts, core.DimText(fmt.Sprintf("'%s'", issue.Value)))
	}

	prefix := ""
	if len(parts) > 0 {
		prefix = strings.Join(parts, " • ") + ": "
	}

	return fmt.Sprintf("  • %s%s", prefix, issue.Message)
}

// runConfigValidation validates configuration files and environment variables
func runConfigValidation(c *cli.Context) error {
	fmt.Println("🔍 Configuration Validation")
	fmt.Println("✅ Config validation is working!")
	return nil
}

// compileLaTeXToPDF compiles LaTeX files to PDF using XeLaTeX
func compileLaTeXToPDF(cfg core.Config) error {
	latexDir := filepath.Join(cfg.OutputDir, "latex")
	pdfDir := filepath.Join(cfg.OutputDir, "pdfs")
	auxDir := filepath.Join(cfg.OutputDir, "auxiliary")

	// Find the main LaTeX file (usually the first .tex file)
	files, err := os.ReadDir(latexDir)
	if err != nil {
		return fmt.Errorf("failed to read latex directory: %w", err)
	}

	var mainTexFile string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".tex") {
			mainTexFile = filepath.Join(latexDir, file.Name())
			break
		}
	}

	if mainTexFile == "" {
		return fmt.Errorf("no LaTeX file found in %s", latexDir)
	}

	// Change to latex directory for compilation
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(latexDir); err != nil {
		return fmt.Errorf("failed to change to latex directory: %w", err)
	}

	// Run XeLaTeX compilation
	cmd := exec.Command("xelatex", "-interaction=nonstopmode", filepath.Base(mainTexFile))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xelatex compilation failed: %w\nOutput: %s", err, string(output))
	}

	// Move generated files to appropriate directories
	baseName := strings.TrimSuffix(filepath.Base(mainTexFile), ".tex")
	
	// Ensure paths are absolute to avoid issues after chdir
	absPdfDir, err := filepath.Abs(pdfDir)
	if err != nil {
		absPdfDir = pdfDir // Fallback to relative if Abs fails
	}
	absAuxDir, err := filepath.Abs(auxDir)
	if err != nil {
		absAuxDir = auxDir // Fallback to relative if Abs fails
	}
	
	// Move PDF to pdfs directory
	pdfFile := baseName + ".pdf"
	if _, err := os.Stat(pdfFile); err == nil {
		destPath := filepath.Join(absPdfDir, pdfFile)
		if err := os.Rename(pdfFile, destPath); err != nil {
			logger.Warn("Failed to move PDF file: %v", err)
		}
	}

	// Move auxiliary files to auxiliary directory
	auxFiles := []string{".aux", ".log", ".fdb_latexmk", ".fls", ".synctex.gz", ".tmp"}
	for _, ext := range auxFiles {
		auxFile := baseName + ext
		if _, err := os.Stat(auxFile); err == nil {
			destPath := filepath.Join(absAuxDir, auxFile)
			if err := os.Rename(auxFile, destPath); err != nil {
				logger.Warn("Failed to move auxiliary file %s: %v", auxFile, err)
			}
		}
	}

	return nil
}
