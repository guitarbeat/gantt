// Package services provides service classes for the document generation system.
//
// CoverageAnalyzer handles test coverage analysis and reporting functionality.
package services

import (
	"fmt"
	"os"
	"os/exec"

	"phd-dissertation-planner/src/app/helpers"
	"phd-dissertation-planner/src/core"
)

// CoverageAnalyzer handles test coverage analysis and reporting
type CoverageAnalyzer struct {
	logger core.Logger
}

// CoverageAnalysisResult represents the result of coverage analysis
type CoverageAnalysisResult struct {
	OverallCoverage float64
	PackageCount    int
	Recommendations []string
	Success         bool
	Errors          []error
}

// NewCoverageAnalyzer creates a new coverage analyzer
func NewCoverageAnalyzer() *CoverageAnalyzer {
	return &CoverageAnalyzer{
		logger: *core.NewDefaultLogger(),
	}
}

// RunTestCoverage executes tests with coverage analysis and provides formatted results
func (ca *CoverageAnalyzer) RunTestCoverage() (*CoverageAnalysisResult, error) {
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
		return &CoverageAnalysisResult{
			Success: false,
			Errors:  []error{err},
		}, err
	}

	// Check if coverage file was created
	if _, err := os.Stat(coverageFile); os.IsNotExist(err) {
		fmt.Println("⚠️  No coverage data generated")
		return &CoverageAnalysisResult{
			Success: true,
			Errors:  []error{fmt.Errorf("no coverage data generated")},
		}, nil
	}

	// Parse and display coverage report
	result, err := ca.analyzeCoverage(coverageFile)
	if err != nil {
		fmt.Printf("⚠️  Coverage analysis failed: %v\n", err)
		return &CoverageAnalysisResult{
			Success: false,
			Errors:  []error{err},
		}, err
	}

	return result, nil
}

// analyzeCoverage parses the coverage file and provides a formatted report
func (ca *CoverageAnalyzer) analyzeCoverage(coverageFile string) (*CoverageAnalysisResult, error) {
	report, err := helpers.ParseCoverageFile(coverageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coverage file: %w", err)
	}

	// Display results
	fmt.Println("\n📊 Coverage Analysis Report")
	fmt.Println("═══════════════════════════════════════")

	// Package breakdown
	fmt.Println("Package Coverage:")
	for pkg, coverages := range report.PackageCoverage {
		if len(coverages) == 0 {
			continue
		}

		// Calculate average coverage for package
		avgCoverage := helpers.CalculatePackageAverage(coverages)
		status := helpers.GetCoverageStatus(avgCoverage)

		fmt.Printf("  %s %-20s %.1f%% (%d files)\n", status, pkg, avgCoverage, len(coverages))
	}

	// Overall statistics
	fmt.Printf("\nOverall Coverage: %.1f%%\n", report.OverallCoverage)
	fmt.Printf("Files Analyzed: %d\n", len(report.PackageCoverage))

	// Generate recommendations
	recommendations := helpers.FormatCoverageRecommendations(report.OverallCoverage)
	fmt.Println("\n💡 Recommendations:")
	for _, rec := range recommendations {
		fmt.Printf("  %s\n", rec)
	}

	fmt.Println("═══════════════════════════════════════")

	return &CoverageAnalysisResult{
		OverallCoverage: report.OverallCoverage,
		PackageCount:    len(report.PackageCoverage),
		Recommendations: recommendations,
		Success:         true,
		Errors:          []error{},
	}, nil
}

// GetCoverageSummary returns a summary of the coverage analysis
func (ca *CoverageAnalyzer) GetCoverageSummary(result *CoverageAnalysisResult) string {
	if !result.Success {
		return "Coverage analysis failed"
	}

	summary := fmt.Sprintf("Coverage Summary:\n")
	summary += fmt.Sprintf("  Overall Coverage: %.1f%%\n", result.OverallCoverage)
	summary += fmt.Sprintf("  Packages Analyzed: %d\n", result.PackageCount)
	summary += fmt.Sprintf("  Recommendations: %d\n", len(result.Recommendations))
	
	if len(result.Errors) > 0 {
		summary += fmt.Sprintf("  Errors: %d\n", len(result.Errors))
	}

	return summary
}
