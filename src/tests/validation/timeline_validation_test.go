package validation

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// TimelineValidationResult holds the results of timeline validation
type TimelineValidationResult struct {
	Errors   []string
	Warnings []string
	Tasks    map[string]TaskInfo
}

// TaskInfo holds information about a task for validation
type TaskInfo struct {
	Row         int
	Task        string
	Dependencies string
	Phase       string
	SubPhase    string
	StartDate   string
	EndDate     string
	Objective   string
	Milestone   string
	Status      string
}

// ValidationError represents a validation error
type ValidationError struct {
	Type    string
	Message string
	Details string
}

// NewTimelineValidationTest creates a new timeline validation test
func NewTimelineValidationTest() *TimelineValidationTest {
	return &TimelineValidationTest{}
}

type TimelineValidationTest struct {
	tasks map[string]TaskInfo
}

// LoadTasksFromCSV loads tasks from a CSV file
func (t *TimelineValidationTest) LoadTasksFromCSV(csvFile string) error {
	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header and one data row")
	}

	// Parse header
	header := records[0]
	headerMap := make(map[string]int)
	for i, col := range header {
		headerMap[col] = i
	}

	// Validate required columns
	requiredColumns := []string{"Task ID", "Dependencies", "Task", "Phase", "Sub-Phase", "Start Date", "End Date", "Objective", "Milestone", "Status"}
	for _, col := range requiredColumns {
		if _, exists := headerMap[col]; !exists {
			return fmt.Errorf("missing required column: %s", col)
		}
	}

	// Parse data rows
	t.tasks = make(map[string]TaskInfo)
	for i, record := range records[1:] {
		if len(record) < len(header) {
			continue // Skip incomplete rows
		}

		taskID := strings.TrimSpace(record[headerMap["Task ID"]])
		if taskID == "" {
			continue // Skip rows without task ID
		}

		t.tasks[taskID] = TaskInfo{
			Row:         i + 2, // +2 because we skip header and 0-index
			Task:        strings.TrimSpace(record[headerMap["Task"]]),
			Dependencies: strings.TrimSpace(record[headerMap["Dependencies"]]),
			Phase:       strings.TrimSpace(record[headerMap["Phase"]]),
			SubPhase:    strings.TrimSpace(record[headerMap["Sub-Phase"]]),
			StartDate:   strings.TrimSpace(record[headerMap["Start Date"]]),
			EndDate:     strings.TrimSpace(record[headerMap["End Date"]]),
			Objective:   strings.TrimSpace(record[headerMap["Objective"]]),
			Milestone:   strings.TrimSpace(record[headerMap["Milestone"]]),
			Status:      strings.TrimSpace(record[headerMap["Status"]]),
		}
	}

	return nil
}

// ValidateTimeline performs comprehensive timeline validation
func (t *TimelineValidationTest) ValidateTimeline() *TimelineValidationResult {
	result := &TimelineValidationResult{
		Errors:   []string{},
		Warnings: []string{},
		Tasks:    t.tasks,
	}

	// 1. Check for duplicate task IDs
	t.checkDuplicateTaskIDs(result)

	// 2. Check for missing dependencies
	t.checkMissingDependencies(result)

	// 3. Check for circular dependencies
	t.checkCircularDependencies(result)

	// 4. Check for orphaned tasks
	t.checkOrphanedTasks(result)

	// 5. Check task ID format consistency
	t.checkTaskIDFormat(result)

	// 6. Check timeline logic based on dependencies
	t.checkTimelineLogic(result)

	// 7. Check for overlapping tasks in same phase/sub-phase
	t.checkOverlappingTasks(result)

	// 8. Check for gaps in sequential tasks
	t.checkGapsInSequentialTasks(result)

	return result
}

// checkDuplicateTaskIDs checks for duplicate task IDs
func (t *TimelineValidationTest) checkDuplicateTaskIDs(result *TimelineValidationResult) {
	seenIDs := make(map[string]bool)
	duplicates := []string{}

	for taskID := range t.tasks {
		if seenIDs[taskID] {
			duplicates = append(duplicates, taskID)
		}
		seenIDs[taskID] = true
	}

	if len(duplicates) > 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Duplicate Task IDs found: %v", duplicates))
	}
}

// checkMissingDependencies checks for missing dependencies
func (t *TimelineValidationTest) checkMissingDependencies(result *TimelineValidationResult) {
	missingDeps := []string{}

	for taskID, taskInfo := range t.tasks {
		if taskInfo.Dependencies == "" {
			continue
		}

		deps := t.parseDependencies(taskInfo.Dependencies)
		for _, dep := range deps {
			if dep != "" && !t.taskExists(dep) {
				missingDeps = append(missingDeps, fmt.Sprintf("Task '%s' references missing task '%s' (row %d)", taskID, dep, taskInfo.Row))
			}
		}
	}

	if len(missingDeps) > 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Missing dependencies found: %d", len(missingDeps)))
		for _, issue := range missingDeps {
			result.Errors = append(result.Errors, "  "+issue)
		}
	}
}

// checkCircularDependencies checks for circular dependencies
func (t *TimelineValidationTest) checkCircularDependencies(result *TimelineValidationResult) {
	circularDeps := []string{}

	for taskID := range t.tasks {
		cycle := t.findCircularDependency(taskID, make(map[string]bool), []string{})
		if len(cycle) > 0 {
			circularDeps = append(circularDeps, fmt.Sprintf("Cycle: %s", strings.Join(cycle, " → ")))
		}
	}

	if len(circularDeps) > 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Circular dependencies found: %d", len(circularDeps)))
		for _, cycle := range circularDeps {
			result.Errors = append(result.Errors, "  "+cycle)
		}
	}
}

// checkOrphanedTasks checks for orphaned tasks
func (t *TimelineValidationTest) checkOrphanedTasks(result *TimelineValidationResult) {
	dependents := make(map[string][]string)

	// Build dependency map
	for taskID, taskInfo := range t.tasks {
		if taskInfo.Dependencies == "" {
			continue
		}

		deps := t.parseDependencies(taskInfo.Dependencies)
		for _, dep := range deps {
			if dep != "" {
				dependents[dep] = append(dependents[dep], taskID)
			}
		}
	}

	// Find orphaned tasks
	orphaned := []string{}
	for taskID, taskInfo := range t.tasks {
		if taskInfo.Dependencies == "" && len(dependents[taskID]) == 0 {
			// Exclude milestones from orphaned check
			if !strings.HasSuffix(taskID, "M1") && !strings.HasSuffix(taskID, "M2") && !strings.HasSuffix(taskID, "M3") {
				orphaned = append(orphaned, fmt.Sprintf("%s: %s", taskID, taskInfo.Task))
			}
		}
	}

	if len(orphaned) > 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Orphaned tasks found: %d", len(orphaned)))
		for _, task := range orphaned {
			result.Warnings = append(result.Warnings, "  "+task)
		}
	}
}

// checkTaskIDFormat checks task ID format consistency
func (t *TimelineValidationTest) checkTaskIDFormat(result *TimelineValidationResult) {
	formatIssues := []string{}

	for taskID := range t.tasks {
		if taskID == "" {
			formatIssues = append(formatIssues, "Empty task ID found")
			continue
		}

		// Check for basic format patterns
		if !strings.HasPrefix(taskID, "T") || (!strings.Contains(taskID, ".") && !strings.HasSuffix(taskID, "M1") && !strings.HasSuffix(taskID, "M2") && !strings.HasSuffix(taskID, "M3")) {
			formatIssues = append(formatIssues, fmt.Sprintf("Unusual format: %s", taskID))
		}
	}

	if len(formatIssues) > 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Task ID format issues: %d", len(formatIssues)))
		for _, issue := range formatIssues {
			result.Warnings = append(result.Warnings, "  "+issue)
		}
	}
}

// checkTimelineLogic checks timeline logic based on dependencies
func (t *TimelineValidationTest) checkTimelineLogic(result *TimelineValidationResult) {
	timelineIssues := []string{}

	for taskID, taskInfo := range t.tasks {
		if taskInfo.Dependencies == "" {
			continue
		}

		taskStart, err1 := t.parseDate(taskInfo.StartDate)
		taskEnd, err2 := t.parseDate(taskInfo.EndDate)

		if err1 != nil || err2 != nil {
			timelineIssues = append(timelineIssues, fmt.Sprintf("Task '%s' has invalid date format: Start='%s', End='%s'", taskID, taskInfo.StartDate, taskInfo.EndDate))
			continue
		}

		// Check if task starts after it ends
		if taskStart.After(taskEnd) {
			timelineIssues = append(timelineIssues, fmt.Sprintf("Task '%s' starts after it ends: %s to %s", taskID, taskInfo.StartDate, taskInfo.EndDate))
		}

		deps := t.parseDependencies(taskInfo.Dependencies)
		for _, dep := range deps {
			if dep == "" || !t.taskExists(dep) {
				continue
			}

			depInfo := t.tasks[dep]
			_, err1 := t.parseDate(depInfo.StartDate)
			depEnd, err2 := t.parseDate(depInfo.EndDate)

			if err1 != nil || err2 != nil {
				timelineIssues = append(timelineIssues, fmt.Sprintf("Dependency '%s' of task '%s' has invalid date format", dep, taskID))
				continue
			}

			// Check if task starts before dependency ends
			if taskStart.Before(depEnd) {
				timelineIssues = append(timelineIssues, fmt.Sprintf("Task '%s' starts before dependency '%s' ends: %s vs %s", taskID, dep, taskInfo.StartDate, depInfo.EndDate))
			}

			// Check if task ends before dependency ends
			if taskEnd.Before(depEnd) {
				timelineIssues = append(timelineIssues, fmt.Sprintf("Task '%s' ends before dependency '%s' ends: %s vs %s", taskID, dep, taskInfo.EndDate, depInfo.EndDate))
			}
		}
	}

	if len(timelineIssues) > 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Timeline logic issues found: %d", len(timelineIssues)))
		for _, issue := range timelineIssues {
			result.Errors = append(result.Errors, "  "+issue)
		}
	}
}

// checkOverlappingTasks checks for overlapping tasks in same phase/sub-phase
func (t *TimelineValidationTest) checkOverlappingTasks(result *TimelineValidationResult) {
	// Group tasks by phase and sub-phase
	phaseGroups := make(map[string][]string)
	for taskID, taskInfo := range t.tasks {
		phaseKey := fmt.Sprintf("%s|%s", taskInfo.Phase, taskInfo.SubPhase)
		phaseGroups[phaseKey] = append(phaseGroups[phaseKey], taskID)
	}

	overlapIssues := []string{}
	for phaseKey, taskIDs := range phaseGroups {
		if len(taskIDs) < 2 {
			continue
		}

		// Check for overlaps within the same phase/sub-phase
		for i, task1ID := range taskIDs {
			for _, task2ID := range taskIDs[i+1:] {
				task1Info := t.tasks[task1ID]
				task2Info := t.tasks[task2ID]

				start1, err1 := t.parseDate(task1Info.StartDate)
				end1, err2 := t.parseDate(task1Info.EndDate)
				start2, err3 := t.parseDate(task2Info.StartDate)
				end2, err4 := t.parseDate(task2Info.EndDate)

				if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
					continue
				}

				// Check for overlap
				if !(end1.Before(start2) || end1.Equal(start2) || end2.Before(start1) || end2.Equal(start1)) {
					overlapIssues = append(overlapIssues, fmt.Sprintf("Phase %s: %s (%s to %s) overlaps with %s (%s to %s)", phaseKey, task1ID, task1Info.StartDate, task1Info.EndDate, task2ID, task2Info.StartDate, task2Info.EndDate))
				}
			}
		}
	}

	if len(overlapIssues) > 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Overlapping tasks found: %d", len(overlapIssues)))
		for _, issue := range overlapIssues {
			result.Warnings = append(result.Warnings, "  "+issue)
		}
	}
}

// checkGapsInSequentialTasks checks for gaps in sequential tasks
func (t *TimelineValidationTest) checkGapsInSequentialTasks(result *TimelineValidationResult) {
	gapIssues := []string{}

	for taskID, taskInfo := range t.tasks {
		if taskInfo.Dependencies == "" {
			continue
		}

		taskStart, err1 := t.parseDate(taskInfo.StartDate)
		if err1 != nil {
			continue
		}

		deps := t.parseDependencies(taskInfo.Dependencies)
		for _, dep := range deps {
			if dep == "" || !t.taskExists(dep) {
				continue
			}

			depInfo := t.tasks[dep]
			depEnd, err2 := t.parseDate(depInfo.EndDate)
			if err2 != nil {
				continue
			}

			// Check for gaps (more than 7 days between tasks)
			gapDays := int(taskStart.Sub(depEnd).Hours() / 24)
			if gapDays > 7 {
				gapIssues = append(gapIssues, fmt.Sprintf("Task '%s' after '%s': %d days gap between %s and %s", taskID, dep, gapDays, depInfo.EndDate, taskInfo.StartDate))
			}
		}
	}

	if len(gapIssues) > 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Large gaps found: %d", len(gapIssues)))
		for _, issue := range gapIssues {
			result.Warnings = append(result.Warnings, "  "+issue)
		}
	}
}

// Helper methods

func (t *TimelineValidationTest) parseDependencies(depsStr string) []string {
	if depsStr == "" {
		return []string{}
	}

	// Handle comma-separated dependencies
	deps := strings.Split(depsStr, ",")
	result := make([]string, 0, len(deps))
	for _, dep := range deps {
		dep = strings.TrimSpace(dep)
		dep = strings.Trim(dep, "\"")
		if dep != "" {
			result = append(result, dep)
		}
	}
	return result
}

func (t *TimelineValidationTest) taskExists(taskID string) bool {
	_, exists := t.tasks[taskID]
	return exists
}

func (t *TimelineValidationTest) parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (t *TimelineValidationTest) findCircularDependency(taskID string, visited map[string]bool, path []string) []string {
	// Check if we've already visited this task in the current path
	for i, id := range path {
		if id == taskID {
			return path[i:] // Return the cycle
		}
	}

	// Check if we've already processed this task
	if visited[taskID] {
		return nil
	}

	visited[taskID] = true
	path = append(path, taskID)

	taskInfo, exists := t.tasks[taskID]
	if !exists || taskInfo.Dependencies == "" {
		return nil
	}

	deps := t.parseDependencies(taskInfo.Dependencies)
	for _, dep := range deps {
		if dep != "" {
			cycle := t.findCircularDependency(dep, visited, path)
			if len(cycle) > 0 {
				return cycle
			}
		}
	}

	return nil
}

// TestTimelineValidation tests the timeline validation functionality
func TestTimelineValidation(t *testing.T) {
	// Find the CSV file
	csvFile := "../../../input/Research Timeline v5 - Comprehensive.csv"
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		// Try alternative path
		csvFile = "../../input/Research Timeline v5 - Comprehensive.csv"
		if _, err := os.Stat(csvFile); os.IsNotExist(err) {
			// Try from project root
			csvFile = "input/Research Timeline v5 - Comprehensive.csv"
			if _, err := os.Stat(csvFile); os.IsNotExist(err) {
				t.Skip("CSV file not found, skipping timeline validation test")
			}
		}
	}

	validator := NewTimelineValidationTest()
	err := validator.LoadTasksFromCSV(csvFile)
	if err != nil {
		t.Fatalf("Failed to load CSV file: %v", err)
	}

	result := validator.ValidateTimeline()

	// Print validation results
	t.Logf("Timeline Validation Results:")
	t.Logf("Tasks loaded: %d", len(result.Tasks))

	if len(result.Errors) > 0 {
		t.Logf("❌ ERRORS: %d", len(result.Errors))
		for _, error := range result.Errors {
			t.Logf("   • %s", error)
		}
	} else {
		t.Logf("✅ NO ERRORS FOUND")
	}

	if len(result.Warnings) > 0 {
		t.Logf("⚠️  WARNINGS: %d", len(result.Warnings))
		for _, warning := range result.Warnings {
			t.Logf("   • %s", warning)
		}
	} else {
		t.Logf("✅ NO WARNINGS")
	}

	// Fail the test if there are errors
	if len(result.Errors) > 0 {
		t.Errorf("Timeline validation failed with %d errors", len(result.Errors))
	}
}

// TestTimelineValidationWithTestData tests validation with test data
func TestTimelineValidationWithTestData(t *testing.T) {
	validator := NewTimelineValidationTest()

	// Create test data with known issues
	validator.tasks = map[string]TaskInfo{
		"T1.0": {
			Row:         2,
			Task:        "Test Task 1",
			Dependencies: "",
			Phase:       "1",
			SubPhase:    "Test Phase",
			StartDate:   "2025-01-01",
			EndDate:     "2025-01-10",
			Objective: "Test objective 1",
			Milestone:    "false",
			Status:      "Not Started",
		},
		"T1.1": {
			Row:         3,
			Task:        "Test Task 2",
			Dependencies: "T1.0",
			Phase:       "1",
			SubPhase:    "Test Phase",
			StartDate:   "2025-01-05", // Overlaps with T1.0
			EndDate:     "2025-01-15",
			Objective:   "Test objective 2",
			Milestone:   "false",
			Status:      "Not Started",
		},
		"T1.2": {
			Row:         4,
			Task:        "Test Task 3",
			Dependencies: "T1.1",
			Phase:       "1",
			SubPhase:    "Test Phase",
			StartDate:   "2025-01-20", // Gap after T1.1
			EndDate:     "2025-01-25",
			Objective:   "Test objective 3",
			Milestone:   "false",
			Status:      "Not Started",
		},
	}

	result := validator.ValidateTimeline()

	// Should have warnings for overlapping tasks and gaps
	if len(result.Warnings) == 0 {
		t.Error("Expected warnings for overlapping tasks and gaps, but got none")
	}

	t.Logf("Test validation completed with %d errors and %d warnings", len(result.Errors), len(result.Warnings))
}
