package data

import (
	"fmt"
	"strings"
	"time"
)

// DependencyValidator handles dependency validation and conflict detection
type DependencyValidator struct {
	tasks       map[string]*Task
	graph       map[string][]string // task ID -> list of dependent task IDs
	reverse     map[string][]string // task ID -> list of tasks that depend on it
	visited     map[string]bool     // for cycle detection
	recursion   map[string]bool     // for cycle detection
}

// NewDependencyValidator creates a new dependency validator
func NewDependencyValidator() *DependencyValidator {
	return &DependencyValidator{
		tasks:     make(map[string]*Task),
		graph:     make(map[string][]string),
		reverse:   make(map[string][]string),
		visited:   make(map[string]bool),
		recursion: make(map[string]bool),
	}
}

// AddTask adds a task to the dependency validator
func (dv *DependencyValidator) AddTask(task *Task) {
	if task == nil {
		return
	}
	
	dv.tasks[task.ID] = task
	dv.graph[task.ID] = task.Dependencies
	
	// Update reverse graph
	for _, depID := range task.Dependencies {
		dv.reverse[depID] = append(dv.reverse[depID], task.ID)
	}
}

// AddTasks adds multiple tasks to the dependency validator
func (dv *DependencyValidator) AddTasks(tasks []*Task) {
	for _, task := range tasks {
		dv.AddTask(task)
	}
}

// ValidateDependencies validates all task dependencies
func (dv *DependencyValidator) ValidateDependencies() []DataValidationError {
	var errors []DataValidationError
	
	// Clear previous state
	dv.visited = make(map[string]bool)
	dv.recursion = make(map[string]bool)
	
	// Validate each task's dependencies
	for _, task := range dv.tasks {
		taskErrors := dv.validateTaskDependencies(task)
		errors = append(errors, taskErrors...)
	}
	
	// Detect circular dependencies
	circularErrors := dv.detectCircularDependencies()
	errors = append(errors, circularErrors...)
	
	// Detect dependency conflicts
	conflictErrors := dv.detectDependencyConflicts()
	errors = append(errors, conflictErrors...)
	
	// Detect orphaned dependencies
	orphanErrors := dv.detectOrphanedDependencies()
	errors = append(errors, orphanErrors...)
	
	// Detect dependency chains that are too long
	chainErrors := dv.detectLongDependencyChains()
	errors = append(errors, chainErrors...)
	
	return errors
}

// validateTaskDependencies validates a single task's dependencies
func (dv *DependencyValidator) validateTaskDependencies(task *Task) []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	if task == nil {
		return errors
	}
	
	// Check for self-dependency
	for _, depID := range task.Dependencies {
		if depID == task.ID {
			errors = append(errors, DataValidationError{
				Type:      "DEPENDENCY",
				TaskID:    task.ID,
				Field:     "Dependencies",
				Value:     depID,
				Message:   "Task cannot depend on itself",
				Severity:  "ERROR",
				Timestamp: now,
				Suggestions: []string{
					"Remove self-dependency from task dependencies",
					"Check if this is a data entry error",
				},
			})
		}
	}
	
	// Check for duplicate dependencies
	depCount := make(map[string]int)
	for _, depID := range task.Dependencies {
		depCount[depID]++
		if depCount[depID] > 1 {
			errors = append(errors, DataValidationError{
				Type:      "DEPENDENCY",
				TaskID:    task.ID,
				Field:     "Dependencies",
				Value:     depID,
				Message:   fmt.Sprintf("Duplicate dependency on task %s", depID),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Remove duplicate dependency entries",
					"Check if this is intentional",
				},
			})
		}
	}
	
	// Check for non-existent dependencies
	for _, depID := range task.Dependencies {
		if _, exists := dv.tasks[depID]; !exists {
			errors = append(errors, DataValidationError{
				Type:      "DEPENDENCY",
				TaskID:    task.ID,
				Field:     "Dependencies",
				Value:     depID,
				Message:   fmt.Sprintf("References non-existent task %s", depID),
				Severity:  "ERROR",
				Timestamp: now,
				Suggestions: []string{
					"Check if task ID is correct",
					"Add the missing task",
					"Remove the invalid dependency",
				},
			})
		}
	}
	
	// Check for too many dependencies
	if len(task.Dependencies) > 10 {
		errors = append(errors, DataValidationError{
			Type:      "DEPENDENCY",
			TaskID:    task.ID,
			Field:     "Dependencies",
			Value:     fmt.Sprintf("%d dependencies", len(task.Dependencies)),
			Message:   fmt.Sprintf("Task has too many dependencies (%d), consider breaking into smaller tasks", len(task.Dependencies)),
			Severity:  "WARNING",
			Timestamp: now,
			Suggestions: []string{
				"Consider breaking this task into smaller, more manageable pieces",
				"Review if all dependencies are necessary",
				"Check if some dependencies can be removed",
			},
		})
	}
	
	return errors
}

// detectCircularDependencies detects circular dependencies using DFS
func (dv *DependencyValidator) detectCircularDependencies() []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Clear visited and recursion states
	dv.visited = make(map[string]bool)
	dv.recursion = make(map[string]bool)
	
	// Check each task for cycles
	for taskID := range dv.tasks {
		if !dv.visited[taskID] {
			cycle := dv.detectCycleFromTask(taskID)
			if len(cycle) > 0 {
				errors = append(errors, DataValidationError{
					Type:      "CIRCULAR_DEPENDENCY",
					TaskID:    taskID,
					Field:     "Dependencies",
					Value:     strings.Join(cycle, " -> "),
					Message:   fmt.Sprintf("Circular dependency detected: %s", strings.Join(cycle, " -> ")),
					Severity:  "ERROR",
					Timestamp: now,
					Suggestions: []string{
						"Break the circular dependency by removing one dependency",
						"Restructure the task dependencies to avoid cycles",
						"Consider if the tasks can be combined or reordered",
					},
				})
			}
		}
	}
	
	return errors
}

// detectCycleFromTask detects a cycle starting from a specific task
func (dv *DependencyValidator) detectCycleFromTask(taskID string) []string {
	if dv.recursion[taskID] {
		// Found a cycle, reconstruct the path
		return dv.reconstructCycle(taskID)
	}
	
	if dv.visited[taskID] {
		return nil
	}
	
	dv.visited[taskID] = true
	dv.recursion[taskID] = true
	
	// Check all dependencies
	for _, depID := range dv.graph[taskID] {
		cycle := dv.detectCycleFromTask(depID)
		if len(cycle) > 0 {
			return cycle
		}
	}
	
	dv.recursion[taskID] = false
	return nil
}

// reconstructCycle reconstructs the cycle path
func (dv *DependencyValidator) reconstructCycle(startTaskID string) []string {
	var cycle []string
	current := startTaskID
	
	for {
		cycle = append(cycle, current)
		if len(cycle) > len(dv.tasks) {
			// Prevent infinite loop
			break
		}
		
		// Find the next task in the cycle
		found := false
		for _, depID := range dv.graph[current] {
			if dv.recursion[depID] {
				current = depID
				found = true
				break
			}
		}
		
		if !found || current == startTaskID {
			break
		}
	}
	
	return cycle
}

// detectDependencyConflicts detects conflicts in dependency relationships
func (dv *DependencyValidator) detectDependencyConflicts() []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Check for tasks that depend on each other (mutual dependency)
	for taskID, task := range dv.tasks {
		for _, depID := range task.Dependencies {
			if depTask, exists := dv.tasks[depID]; exists {
				// Check if the dependency also depends on this task
				for _, reverseDepID := range depTask.Dependencies {
					if reverseDepID == taskID {
						errors = append(errors, DataValidationError{
							Type:      "DEPENDENCY_CONFLICT",
							TaskID:    taskID,
							Field:     "Dependencies",
							Value:     depID,
							Message:   fmt.Sprintf("Mutual dependency detected between %s and %s", taskID, depID),
							Severity:  "WARNING",
							Timestamp: now,
							Suggestions: []string{
								"Remove one of the mutual dependencies",
								"Consider if these tasks can be combined",
								"Restructure the dependency relationship",
							},
						})
					}
				}
			}
		}
	}
	
	// Check for tasks that depend on tasks with later start dates
	for taskID, task := range dv.tasks {
		for _, depID := range task.Dependencies {
			if depTask, exists := dv.tasks[depID]; exists {
				if depTask.StartDate.After(task.StartDate) {
					errors = append(errors, DataValidationError{
						Type:      "DEPENDENCY_CONFLICT",
						TaskID:    taskID,
						Field:     "Dependencies",
						Value:     depID,
						Message:   fmt.Sprintf("Task %s depends on %s which starts later (%s vs %s)", taskID, depID, task.StartDate.Format("2006-01-02"), depTask.StartDate.Format("2006-01-02")),
						Severity:  "WARNING",
						Timestamp: now,
						Suggestions: []string{
							"Adjust the start date of the dependent task",
							"Adjust the start date of the dependency",
							"Check if the dependency relationship is correct",
						},
					})
				}
			}
		}
	}
	
	return errors
}

// detectOrphanedDependencies detects dependencies that are not referenced by any task
func (dv *DependencyValidator) detectOrphanedDependencies() []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Find all tasks that are referenced as dependencies
	referencedTasks := make(map[string]bool)
	for _, task := range dv.tasks {
		for _, depID := range task.Dependencies {
			referencedTasks[depID] = true
		}
	}
	
	// Check for tasks that are not referenced by any other task
	for taskID, task := range dv.tasks {
		if !referencedTasks[taskID] && len(task.Dependencies) > 0 {
			// This task has dependencies but is not referenced by any other task
			errors = append(errors, DataValidationError{
				Type:      "ORPHANED_DEPENDENCY",
				TaskID:    taskID,
				Field:     "Dependencies",
				Value:     "orphaned",
				Message:   fmt.Sprintf("Task %s has dependencies but is not referenced by any other task", taskID),
				Severity:  "INFO",
				Timestamp: now,
				Suggestions: []string{
					"Check if this task should be referenced by other tasks",
					"Consider if this is a root task in the dependency chain",
					"Verify if the task is part of the main workflow",
				},
			})
		}
	}
	
	return errors
}

// detectLongDependencyChains detects dependency chains that are too long
func (dv *DependencyValidator) detectLongDependencyChains() []DataValidationError {
	var errors []DataValidationError
	now := time.Now()
	
	// Find the longest dependency chain for each task
	for taskID := range dv.tasks {
		chainLength := dv.getDependencyChainLength(taskID)
		if chainLength > 5 {
			errors = append(errors, DataValidationError{
				Type:      "LONG_DEPENDENCY_CHAIN",
				TaskID:    taskID,
				Field:     "Dependencies",
				Value:     fmt.Sprintf("%d levels", chainLength),
				Message:   fmt.Sprintf("Task %s has a dependency chain of %d levels, which may be too complex", taskID, chainLength),
				Severity:  "WARNING",
				Timestamp: now,
				Suggestions: []string{
					"Consider breaking the dependency chain into smaller parts",
					"Review if all dependencies in the chain are necessary",
					"Check if some tasks can be parallelized",
				},
			})
		}
	}
	
	return errors
}

// getDependencyChainLength calculates the length of the dependency chain for a task
func (dv *DependencyValidator) getDependencyChainLength(taskID string) int {
	visited := make(map[string]bool)
	return dv.calculateChainLength(taskID, visited)
}

// calculateChainLength recursively calculates the chain length
func (dv *DependencyValidator) calculateChainLength(taskID string, visited map[string]bool) int {
	if visited[taskID] {
		return 0 // Prevent infinite recursion
	}
	
	visited[taskID] = true
	maxLength := 0
	
	for _, depID := range dv.graph[taskID] {
		length := dv.calculateChainLength(depID, visited)
		if length >= maxLength {
			maxLength = length + 1
		}
	}
	
	return maxLength
}

// GetDependencyLevels returns all tasks grouped by their dependency level
func (dv *DependencyValidator) GetDependencyLevels() map[int][]string {
	levels := make(map[int][]string)
	
	for taskID := range dv.tasks {
		level := dv.getDependencyLevel(taskID)
		levels[level] = append(levels[level], taskID)
	}
	
	return levels
}

// getDependencyLevel calculates the dependency level of a task
func (dv *DependencyValidator) getDependencyLevel(taskID string) int {
	visited := make(map[string]bool)
	return dv.calculateDependencyLevel(taskID, visited)
}

// calculateDependencyLevel recursively calculates the dependency level
func (dv *DependencyValidator) calculateDependencyLevel(taskID string, visited map[string]bool) int {
	if visited[taskID] {
		return 0 // Prevent infinite recursion
	}
	
	visited[taskID] = true
	maxLevel := 0
	
	for _, depID := range dv.graph[taskID] {
		level := dv.calculateDependencyLevel(depID, visited)
		if level >= maxLevel {
			maxLevel = level + 1
		}
	}
	
	return maxLevel
}

// GetTasksByDependencyLevel returns tasks at a specific dependency level
func (dv *DependencyValidator) GetTasksByDependencyLevel(level int) []*Task {
	var tasks []*Task
	
	for taskID := range dv.tasks {
		if dv.getDependencyLevel(taskID) == level {
			tasks = append(tasks, dv.tasks[taskID])
		}
	}
	
	return tasks
}

// GetDependencyGraph returns the dependency graph
func (dv *DependencyValidator) GetDependencyGraph() map[string][]string {
	return dv.graph
}

// GetReverseDependencyGraph returns the reverse dependency graph
func (dv *DependencyValidator) GetReverseDependencyGraph() map[string][]string {
	return dv.reverse
}

// GetTaskDependencies returns the dependencies of a specific task
func (dv *DependencyValidator) GetTaskDependencies(taskID string) []*Task {
	var deps []*Task
	
	for _, depID := range dv.graph[taskID] {
		if task, exists := dv.tasks[depID]; exists {
			deps = append(deps, task)
		}
	}
	
	return deps
}

// GetTaskDependents returns the tasks that depend on a specific task
func (dv *DependencyValidator) GetTaskDependents(taskID string) []*Task {
	var deps []*Task
	
	for _, depID := range dv.reverse[taskID] {
		if task, exists := dv.tasks[depID]; exists {
			deps = append(deps, task)
		}
	}
	
	return deps
}

// ValidateDependencyIntegrity performs comprehensive dependency validation
func (dv *DependencyValidator) ValidateDependencyIntegrity() *ValidationResult {
	errors := dv.ValidateDependencies()
	
	// Categorize errors by severity
	var errorList, warningList, infoList []DataValidationError
	
	for _, err := range errors {
		switch err.Severity {
		case "ERROR":
			errorList = append(errorList, err)
		case "WARNING":
			warningList = append(warningList, err)
		case "INFO":
			infoList = append(infoList, err)
		}
	}
	
	result := &ValidationResult{
		IsValid:      len(errorList) == 0,
		Errors:       errorList,
		Warnings:     warningList,
		Info:         infoList,
		TaskCount:    len(dv.tasks),
		ErrorCount:   len(errorList),
		WarningCount: len(warningList),
		Timestamp:    time.Now(),
	}
	
	result.Summary = result.GetValidationSummary()
	
	return result
}
