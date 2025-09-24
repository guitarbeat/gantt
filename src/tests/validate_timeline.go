package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
	"path/filepath"
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

// TimelineValidationTest handles timeline validation
type TimelineValidationTest struct {
	tasks map[string]TaskInfo
}

// NewTimelineValidationTest creates a new timeline validation test
func NewTimelineValidationTest() *TimelineValidationTest {
	return &TimelineValidationTest{}
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

	fmt.Println("🔍 Starting Timeline Validation...")
	fmt.Println(strings.Repeat("=", 50))

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

	// Print summary
	fmt.Println()
	fmt.Println("📋 VALIDATION SUMMARY")
	fmt.Println(strings.Repeat("=", 50))

	if len(result.Errors) > 0 {
		fmt.Printf("❌ ERRORS: %d\n", len(result.Errors))
		for _, error := range result.Errors {
			fmt.Printf("   • %s\n", error)
		}
	} else {
		fmt.Println("✅ NO ERRORS FOUND")
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("⚠️  WARNINGS: %d\n", len(result.Warnings))
		for _, warning := range result.Warnings {
			fmt.Printf("   • %s\n", warning)
		}
	} else {
		fmt.Println("✅ NO WARNINGS")
	}

	fmt.Println()

	if len(result.Errors) > 0 {
		fmt.Println("🔧 RECOMMENDATIONS:")
		fmt.Println("   • Fix missing dependencies")
		fmt.Println("   • Resolve circular dependencies")
		fmt.Println("   • Check for typos in task IDs")
		return result
	} else {
		fmt.Println("🎉 Timeline validation passed!")
		return result
	}
}

// Helper methods (same as in the test file)

func (t *TimelineValidationTest) checkDuplicateTaskIDs(result *TimelineValidationResult) {
	fmt.Println("1️⃣ Checking for duplicate Task IDs...")
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
		fmt.Printf("❌ Found %d duplicate Task IDs: %v\n", len(duplicates), duplicates)
	} else {
		fmt.Println("✅ No duplicate Task IDs found")
	}
	fmt.Println()
}

func (t *TimelineValidationTest) checkMissingDependencies(result *TimelineValidationResult) {
	fmt.Println("2️⃣ Checking for missing dependencies...")
	missingDeps := []string{}

	for taskID, taskInfo := range t.tasks {
		if taskInfo.Dependencies == "" {
			continue
		}

		deps := t.parseDependencies(taskInfo.Dependencies)
		for _, dep := range deps {
			if dep != "" && !t.taskExists(dep) {
				missingDeps = append(missingDeps, fmt.Sprintf("Row %d: Task '%s' references missing task '%s'", taskInfo.Row, taskID, dep))
			}
		}
	}

	if len(missingDeps) > 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Missing dependencies found: %d", len(missingDeps)))
		fmt.Printf("❌ Found %d missing dependencies:\n", len(missingDeps))
		for _, issue := range missingDeps {
			fmt.Printf("   %s\n", issue)
		}
	} else {
		fmt.Println("✅ All dependencies reference existing tasks")
	}
	fmt.Println()
}

func (t *TimelineValidationTest) checkCircularDependencies(result *TimelineValidationResult) {
	fmt.Println("3️⃣ Checking for circular dependencies...")
	circularDeps := []string{}

	for taskID := range t.tasks {
		cycle := t.findCircularDependency(taskID, make(map[string]bool), []string{})
		if len(cycle) > 0 {
			circularDeps = append(circularDeps, fmt.Sprintf("Cycle: %s", strings.Join(cycle, " → ")))
		}
	}

	if len(circularDeps) > 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Circular dependencies found: %d", len(circularDeps)))
		fmt.Printf("❌ Found %d circular dependencies:\n", len(circularDeps))
		for i, cycle := range circularDeps {
			fmt.Printf("   Cycle %d: %s\n", i+1, cycle)
		}
	} else {
		fmt.Println("✅ No circular dependencies found")
	}
	fmt.Println()
}

func (t *TimelineValidationTest) checkOrphanedTasks(result *TimelineValidationResult) {
	fmt.Println("4️⃣ Checking for orphaned tasks...")
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
		fmt.Printf("⚠️  Found %d orphaned tasks (no dependencies, no dependents):\n", len(orphaned))
		for _, task := range orphaned {
			fmt.Printf("   %s\n", task)
		}
	} else {
		fmt.Println("✅ No orphaned tasks found")
	}
	fmt.Println()
}

func (t *TimelineValidationTest) checkTaskIDFormat(result *TimelineValidationResult) {
	fmt.Println("5️⃣ Checking task ID format consistency...")
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
		fmt.Printf("⚠️  Found %d task ID format issues:\n", len(formatIssues))
		for _, issue := range formatIssues {
			fmt.Printf("   %s\n", issue)
		}
	} else {
		fmt.Println("✅ Task ID formats look consistent")
	}
	fmt.Println()
}

func (t *TimelineValidationTest) checkTimelineLogic(result *TimelineValidationResult) {
	fmt.Println("6️⃣ Checking timeline logic based on dependencies...")
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
		fmt.Printf("❌ Found %d timeline logic issues:\n", len(timelineIssues))
		for _, issue := range timelineIssues {
			fmt.Printf("   %s\n", issue)
		}
	} else {
		fmt.Println("✅ Timeline logic looks consistent")
	}
	fmt.Println()
}

func (t *TimelineValidationTest) checkOverlappingTasks(result *TimelineValidationResult) {
	fmt.Println("7️⃣ Checking for overlapping tasks in same phase/sub-phase...")
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
		fmt.Printf("⚠️  Found %d overlapping tasks:\n", len(overlapIssues))
		for _, issue := range overlapIssues {
			fmt.Printf("   %s\n", issue)
		}
	} else {
		fmt.Println("✅ No overlapping tasks found")
	}
	fmt.Println()
}

func (t *TimelineValidationTest) checkGapsInSequentialTasks(result *TimelineValidationResult) {
	fmt.Println("8️⃣ Checking for gaps in sequential tasks...")
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
		fmt.Printf("⚠️  Found %d large gaps between sequential tasks:\n", len(gapIssues))
		for _, issue := range gapIssues {
			fmt.Printf("   %s\n", issue)
		}
	} else {
		fmt.Println("✅ No large gaps found between sequential tasks")
	}
	fmt.Println()
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

// saveValidationLog saves validation results to a log file
func saveValidationLog(result *TimelineValidationResult, logPath string) error {
	// Create output directory if it doesn't exist
	dir := filepath.Dir(logPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open log file for writing
	file, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "Timeline Validation Report\n")
	fmt.Fprintf(file, "Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "========================================\n\n")

	// Write summary
	fmt.Fprintf(file, "VALIDATION SUMMARY\n")
	fmt.Fprintf(file, "========================================\n")
	fmt.Fprintf(file, "❌ ERRORS: %d\n", len(result.Errors))
	fmt.Fprintf(file, "⚠️  WARNINGS: %d\n", len(result.Warnings))
	fmt.Fprintf(file, "\n")

	// Write errors
	if len(result.Errors) > 0 {
		fmt.Fprintf(file, "ERRORS\n")
		fmt.Fprintf(file, "----------------------------------------\n")
		for i, err := range result.Errors {
			fmt.Fprintf(file, "%d. %s\n", i+1, err)
		}
		fmt.Fprintf(file, "\n")
	}

	// Write warnings
	if len(result.Warnings) > 0 {
		fmt.Fprintf(file, "WARNINGS\n")
		fmt.Fprintf(file, "----------------------------------------\n")
		for i, warning := range result.Warnings {
			fmt.Fprintf(file, "%d. %s\n", i+1, warning)
		}
		fmt.Fprintf(file, "\n")
	}

	// Write task details
	fmt.Fprintf(file, "TASK DETAILS\n")
	fmt.Fprintf(file, "========================================\n")
	fmt.Fprintf(file, "Total Tasks: %d\n\n", len(result.Tasks))

	// Group tasks by phase for better organization
	phaseGroups := make(map[string][]TaskInfo)
	for _, task := range result.Tasks {
		phaseGroups[task.Phase] = append(phaseGroups[task.Phase], task)
	}

	for phase, tasks := range phaseGroups {
		fmt.Fprintf(file, "Phase: %s (%d tasks)\n", phase, len(tasks))
		fmt.Fprintf(file, "----------------------------------------\n")
		for _, task := range tasks {
			fmt.Fprintf(file, "  • %s (%s to %s)\n", task.Task, task.StartDate, task.EndDate)
			if task.Status != "" {
				fmt.Fprintf(file, "    Status: %s\n", task.Status)
			}
			if task.Milestone != "" {
				fmt.Fprintf(file, "    Milestone: %s\n", task.Milestone)
			}
		}
		fmt.Fprintf(file, "\n")
	}

	// Write recommendations
	fmt.Fprintf(file, "RECOMMENDATIONS\n")
	fmt.Fprintf(file, "========================================\n")
	if len(result.Errors) > 0 {
		fmt.Fprintf(file, "• Fix missing dependencies\n")
		fmt.Fprintf(file, "• Resolve circular dependencies\n")
		fmt.Fprintf(file, "• Check for typos in task IDs\n")
		fmt.Fprintf(file, "• Review overlapping tasks that cannot run simultaneously\n")
	}
	if len(result.Warnings) > 0 {
		fmt.Fprintf(file, "• Consider adjusting task schedules to reduce overlaps\n")
		fmt.Fprintf(file, "• Review orphaned tasks for proper integration\n")
		fmt.Fprintf(file, "• Evaluate large gaps between sequential tasks\n")
	}
	if len(result.Errors) == 0 && len(result.Warnings) == 0 {
		fmt.Fprintf(file, "✅ No issues detected - timeline looks good!\n")
	}

	return nil
}

func main() {
	csvFile := "../input/Research Timeline v5 - Comprehensive.csv"
	if len(os.Args) > 1 {
		csvFile = os.Args[1]
	}

	validator := NewTimelineValidationTest()
	err := validator.LoadTasksFromCSV(csvFile)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📊 Found %d tasks\n", len(validator.tasks))
	fmt.Println()

	result := validator.ValidateTimeline()

	// Save validation results to log file
	err = saveValidationLog(result, "../output/validation.log")
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to save validation log: %v\n", err)
	} else {
		fmt.Printf("📝 Validation results saved to: ../output/validation.log\n")
	}

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}
