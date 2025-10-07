// Package helpers provides utility functions extracted from the main generator.
//
// This package contains helper functions that were extracted from large methods
// in generator.go to improve maintainability and testability. Functions are
// organized by domain (coverage, config, toc, monthly) for better organization.
//
// Coverage helpers provide functionality for test coverage analysis and reporting.
package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CoverageData represents parsed coverage information for a single file
type CoverageData struct {
	FilePath     string
	PackageName  string
	Coverage     float64
	Statements   int
	Covered      int
}

// CoverageReport represents the overall coverage analysis results
type CoverageReport struct {
	PackageCoverage map[string][]float64
	TotalStatements int
	TotalCovered    int
	OverallCoverage float64
}

// ParseCoverageFile parses a Go coverage file and returns structured data
func ParseCoverageFile(coverageFile string) (*CoverageReport, error) {
	file, err := os.Open(coverageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open coverage file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	report := &CoverageReport{
		PackageCoverage: make(map[string][]float64),
	}

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

		coverageData, err := parseCoverageLine(line)
		if err != nil {
			continue // Skip invalid lines
		}

		report.PackageCoverage[coverageData.PackageName] = append(
			report.PackageCoverage[coverageData.PackageName], 
			coverageData.Coverage,
		)
		report.TotalStatements++
		if coverageData.Coverage > 0 {
			report.TotalCovered++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading coverage file: %w", err)
	}

	// Calculate overall coverage
	if report.TotalStatements > 0 {
		report.OverallCoverage = float64(report.TotalCovered) / float64(report.TotalStatements) * 100
	}

	return report, nil
}

// parseCoverageLine parses a single line from the coverage file
func parseCoverageLine(line string) (*CoverageData, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid coverage line format")
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

	// Parse statement count and covered count
	statements, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse statement count: %w", err)
	}

	covered, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, fmt.Errorf("failed to parse covered count: %w", err)
	}

	// Calculate coverage percentage
	coverage := 0.0
	if statements > 0 {
		coverage = float64(covered) / float64(statements) * 100
	}

	return &CoverageData{
		FilePath:    filePath,
		PackageName: packageName,
		Coverage:    coverage,
		Statements:  statements,
		Covered:     covered,
	}, nil
}

// CalculatePackageAverage calculates the average coverage for a package
func CalculatePackageAverage(coverages []float64) float64 {
	if len(coverages) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, cov := range coverages {
		sum += cov
	}
	return sum / float64(len(coverages))
}

// GetCoverageStatus returns a status emoji based on coverage percentage
func GetCoverageStatus(coverage float64) string {
	if coverage >= 80 {
		return "✅"
	} else if coverage >= 60 {
		return "⚠️ "
	}
	return "❌"
}

// FormatCoverageRecommendations generates recommendations based on coverage
func FormatCoverageRecommendations(overallCoverage float64) []string {
	var recommendations []string

	if overallCoverage < 60 {
		recommendations = append(recommendations,
			"• Coverage is low - consider adding more tests",
			"• Focus on testing critical business logic",
		)
	}

	if overallCoverage >= 80 {
		recommendations = append(recommendations, "• Excellent coverage! Keep up the good work")
	} else {
		recommendations = append(recommendations,
			"• Aim for 80%+ coverage for better reliability",
			"• Add tests for error conditions and edge cases",
		)
	}

	return recommendations
}
