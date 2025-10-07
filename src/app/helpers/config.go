// Package helpers provides utility functions extracted from the main generator.
//
// Configuration helpers provide functionality for auto-detecting CSV files
// and configuration files based on content analysis.
package helpers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CSVFileInfo represents information about a CSV file for priority selection
type CSVFileInfo struct {
	File     os.DirEntry
	Priority int
	ModTime  time.Time
}

// ConfigDetectionResult represents the result of configuration auto-detection
type ConfigDetectionResult struct {
	BaseConfig    string
	AdditionalConfigs []string
	Reason       string
}

// AutoDetectCSV automatically finds the most appropriate CSV file in the input_data directory
func AutoDetectCSV() (string, error) {
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

	csvFiles := filterCSVFiles(files)
	if len(csvFiles) == 0 {
		return "", fmt.Errorf("no CSV files found in input_data directory")
	}

	// If only one CSV file, use it
	if len(csvFiles) == 1 {
		return filepath.Join(inputDir, csvFiles[0].Name()), nil
	}

	// Multiple CSV files - use priority selection
	bestFile, err := selectBestCSVFile(csvFiles)
	if err != nil {
		return "", err
	}

	return filepath.Join(inputDir, bestFile.Name()), nil
}

// filterCSVFiles filters directory entries to only include CSV files
func filterCSVFiles(files []os.DirEntry) []os.DirEntry {
	var csvFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			csvFiles = append(csvFiles, file)
		}
	}
	return csvFiles
}

// selectBestCSVFile selects the best CSV file based on priority rules
func selectBestCSVFile(csvFiles []os.DirEntry) (os.DirEntry, error) {
	var bestFile os.DirEntry
	bestPriority := 0

	for _, file := range csvFiles {
		priority, err := calculateCSVPriority(file)
		if err != nil {
			continue // Skip files with errors
		}

		if priority > bestPriority || (priority == bestPriority && bestFile == nil) {
			bestPriority = priority
			bestFile = file
		} else if priority == bestPriority && bestFile != nil {
			// Compare modification times as tiebreaker
			if isNewerFile(file, bestFile) {
				bestFile = file
			}
		}
	}

	if bestFile == nil {
		// Fallback to first file
		return csvFiles[0], nil
	}

	return bestFile, nil
}

// calculateCSVPriority calculates priority score for a CSV file
func calculateCSVPriority(file os.DirEntry) (int, error) {
	name := strings.ToLower(file.Name())
	priority := 0

	// Highest priority: comprehensive files
	if strings.Contains(name, "comprehensive") {
		priority = 10
	}

	// Versioned files get priority based on version number
	if strings.Contains(name, "v") && strings.Contains(name, ".") {
		if strings.Contains(name, "v5.1") {
			priority = 8
		} else if strings.Contains(name, "v5") {
			priority = 6
		}
	}

	return priority, nil
}

// isNewerFile checks if the first file is newer than the second
func isNewerFile(file1, file2 os.DirEntry) bool {
	info1, err1 := file1.Info()
	info2, err2 := file2.Info()
	
	if err1 != nil || err2 != nil {
		return false
	}
	
	return info1.ModTime().After(info2.ModTime())
}

// AutoDetectConfig automatically determines appropriate configuration files based on CSV content
func AutoDetectConfig(csvPath string) (*ConfigDetectionResult, error) {
	// Read first few lines to detect version/format
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV for config detection: %w", err)
	}
	defer file.Close()

	lines, err := readCSVHeaderLines(file, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV for config detection: %w", err)
	}

	result := &ConfigDetectionResult{
		BaseConfig: "src/core/base.yaml",
	}

	// Detect CSV version from filename or content
	csvName := strings.ToLower(filepath.Base(csvPath))
	
	if strings.Contains(csvName, "v5.1") {
		result.AdditionalConfigs = []string{"src/core/monthly_calendar.yaml"}
		result.Reason = "v5.1 format detected from filename"
	} else if strings.Contains(csvName, "v5") {
		result.AdditionalConfigs = []string{"src/core/calendar.yaml"}
		result.Reason = "v5 format detected from filename"
	} else if hasPhaseSubPhaseColumns(lines) {
		result.AdditionalConfigs = []string{"src/core/monthly_calendar.yaml"}
		result.Reason = "v5.1 format detected from content (phase/sub-phase columns)"
	} else {
		result.Reason = "using default configuration"
	}

	return result, nil
}

// readCSVHeaderLines reads the first n lines from a file
func readCSVHeaderLines(file *os.File, maxLines int) ([]string, error) {
	scanner := bufio.NewScanner(file)
	var lines []string
	
	for i := 0; i < maxLines && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return lines, nil
}

// hasPhaseSubPhaseColumns checks if the CSV header contains phase and sub-phase columns
func hasPhaseSubPhaseColumns(lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	
	header := strings.ToLower(lines[0])
	hasPhase := strings.Contains(header, "phase")
	hasSubPhase := strings.Contains(header, "sub-phase")
	return hasPhase && hasSubPhase
}

// GetConfigPaths returns the configuration file paths based on detection result
func (r *ConfigDetectionResult) GetConfigPaths() []string {
	paths := []string{r.BaseConfig}
	paths = append(paths, r.AdditionalConfigs...)
	return paths
}
