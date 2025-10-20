package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

// * CSV Data Validation
// ============================================================================

// CSVValidator handles comprehensive validation of CSV task data
type CSVValidator struct {
	requiredFields     []string
	validStatuses      map[string]bool
	validPhases        map[string]bool
	validMilestoneValues map[string]bool
	logger             *Logger
}

// NewCSVValidator creates a new CSV validator with default validation rules
func NewCSVValidator() *CSVValidator {
	return &CSVValidator{
		requiredFields: []string{"Task ID", "Task", "Start Date", "End Date"},
		validStatuses: map[string]bool{
			"not started": true,
			"planned":     true,
			"in progress": true,
			"completed":   true,
			"on hold":     true,
			"cancelled":  true,
			"canceled":   true,
			"blocked":    true,
		},
		validPhases: map[string]bool{
			"1": true, "2": true, "3": true, "4": true, "5": true,
			"6": true, "7": true, "8": true, "9": true, "10": true,
		},
		validMilestoneValues: map[string]bool{
			"true":    true,
			"false":  true,
			"critical": true,
			"high":    true,
			"medium":  true,
			"low":     true,
		},
		logger: NewDefaultLogger(),
	}
}

// ValidateCSVFile validates a CSV file and returns detailed validation results
func (v *CSVValidator) ValidateCSVFile(filePath string) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     make([]ValidationIssue, 0),
		Warnings:   make([]ValidationIssue, 0),
		RowCount:   0,
		FieldCount: 0,
	}

	// Check file exists and is readable
	if err := v.validateFileAccess(filePath); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Type:    "file_access",
			Message: fmt.Sprintf("Cannot access CSV file: %v", err),
		})
		return result, err
	}

	// Read and validate CSV content
	reader := NewReader(filePath)
	tasks, err := reader.ReadTasks()
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Type:    "csv_parsing",
			Message: fmt.Sprintf("Failed to parse CSV: %v", err),
		})
		return result, err
	}

	result.RowCount = len(tasks)

	// Validate each task
	for i, task := range tasks {
		if errs := v.validateTask(task, i+2); len(errs) > 0 { // +2 for header row + 0-indexing
			result.Errors = append(result.Errors, errs...)
			result.IsValid = false
		}
		if warns := v.validateTaskWarnings(task, i+2); len(warns) > 0 {
			result.Warnings = append(result.Warnings, warns...)
		}
	}

	// Validate overall data consistency
	if overallErrs := v.validateDataConsistency(tasks); len(overallErrs) > 0 {
		result.Errors = append(result.Errors, overallErrs...)
		result.IsValid = false
	}

	return result, nil
}

// validateFileAccess checks if the CSV file can be accessed
func (v *CSVValidator) validateFileAccess(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	// Check if file is too large (arbitrary limit of 50MB)
	if info.Size() > 50*1024*1024 {
		return fmt.Errorf("CSV file is too large (%d MB), maximum allowed is 50MB", info.Size()/(1024*1024))
	}

	return nil
}

// validateTask validates a single task and returns validation errors
func (v *CSVValidator) validateTask(task Task, rowNum int) []ValidationIssue {
	var errors []ValidationIssue

	// Validate required fields
	if strings.TrimSpace(task.ID) == "" {
		errors = append(errors, ValidationIssue{
			Type:    "required_field",
			Field:   "Task ID",
			Row:     rowNum,
			Message: "Task ID is required",
		})
	}

	if strings.TrimSpace(task.Name) == "" {
		errors = append(errors, ValidationIssue{
			Type:    "required_field",
			Field:   "Task",
			Row:     rowNum,
			Message: "Task name is required",
		})
	}

	// Validate dates
	if task.StartDate.IsZero() {
		errors = append(errors, ValidationIssue{
			Type:    "required_field",
			Field:   "Start Date",
			Row:     rowNum,
			Message: "Start Date is required and must be a valid date",
		})
	}

	if task.EndDate.IsZero() {
		errors = append(errors, ValidationIssue{
			Type:    "required_field",
			Field:   "End Date",
			Row:     rowNum,
			Message: "End Date is required and must be a valid date",
		})
	}

	// Validate date logic
	if !task.StartDate.IsZero() && !task.EndDate.IsZero() {
		if task.EndDate.Before(task.StartDate) {
			errors = append(errors, ValidationIssue{
				Type:  "date_logic",
				Field: "End Date",
				Row:   rowNum,
				Message: fmt.Sprintf("End Date (%s) cannot be before Start Date (%s)",
					task.EndDate.Format("2006-01-02"), task.StartDate.Format("2006-01-02")),
			})
		}

		// Check for unreasonably long tasks (more than 2 years)
		if task.EndDate.Sub(task.StartDate).Hours() > 24*365*2 {
			errors = append(errors, ValidationIssue{
				Type:    "date_range",
				Field:   "End Date",
				Row:     rowNum,
				Message: "Task duration exceeds 2 years, please verify dates",
			})
		}
	}

	// Validate status
	if task.Status != "" {
		status := strings.ToLower(strings.TrimSpace(task.Status))
		if !v.validStatuses[status] {
			errors = append(errors, ValidationIssue{
				Type:  "invalid_value",
				Field: "Status",
				Row:   rowNum,
				Value: task.Status,
				Message: fmt.Sprintf("Invalid status '%s', must be one of: %s",
					task.Status, v.getValidStatusesString()),
			})
		}
	}

	// Validate phase
	if task.Phase != "" {
		if !v.validPhases[task.Phase] {
			errors = append(errors, ValidationIssue{
				Type:    "invalid_value",
				Field:   "Phase",
				Row:     rowNum,
				Value:   task.Phase,
				Message: "Phase must be a number between 1-10",
			})
		}
	}

	// Validate Task ID format
	if task.ID != "" {
		if matched, _ := regexp.MatchString(`^T\d+\.[A-Za-z0-9]+$`, task.ID); !matched {
			errors = append(errors, ValidationIssue{
				Type:    "invalid_format",
				Field:   "Task ID",
				Row:     rowNum,
				Value:   task.ID,
				Message: "Task ID must follow format T{phase}.{identifier} (e.g., T1.1, T2.M1, T3.4a)",
			})
		} else {
			// Validate that task ID phase matches the Phase column
			if task.Phase != "" {
				taskIDPhase := strings.TrimPrefix(strings.Split(task.ID, ".")[0], "T")
				if taskIDPhase != task.Phase {
					errors = append(errors, ValidationIssue{
						Type:    "phase_mismatch",
						Field:   "Task ID",
						Row:     rowNum,
						Value:   task.ID,
						Message: fmt.Sprintf("Task ID phase (%s) does not match Phase column (%s)", taskIDPhase, task.Phase),
					})
				}
			}
		}
	}

	// Note: Milestone field validation happens during CSV parsing, not here

	// Validate dependencies format (comma-separated task IDs)
	if task.Dependencies != nil {
		for _, dep := range task.Dependencies {
			dep = strings.TrimSpace(dep)
			if dep == "" {
				errors = append(errors, ValidationIssue{
					Type:    "invalid_format",
					Field:   "Dependencies",
					Row:     rowNum,
					Message: "Dependency entries cannot be empty",
				})
			} else {
				// Validate dependency task ID format
				if matched, _ := regexp.MatchString(`^T\d+\.[A-Za-z0-9]+$`, dep); !matched {
					errors = append(errors, ValidationIssue{
						Type:    "invalid_format",
						Field:   "Dependencies",
						Row:     rowNum,
						Value:   dep,
						Message: "Dependency must reference valid task ID format T{phase}.{identifier}",
					})
				}
			}
		}
	}

	// Validate reasonable date ranges for PhD timeline (2025-2027)
	if !task.StartDate.IsZero() {
		if task.StartDate.Year() < 2025 || task.StartDate.Year() > 2027 {
			errors = append(errors, ValidationIssue{
				Type:    "date_range",
				Field:   "Start Date",
				Row:     rowNum,
				Value:   task.StartDate.Format("2006-01-02"),
				Message: "Start date should be within PhD timeline range (2025-2027)",
			})
		}
	}

	if !task.EndDate.IsZero() {
		if task.EndDate.Year() < 2025 || task.EndDate.Year() > 2027 {
			errors = append(errors, ValidationIssue{
				Type:    "date_range",
				Field:   "End Date",
				Row:     rowNum,
				Value:   task.EndDate.Format("2006-01-02"),
				Message: "End date should be within PhD timeline range (2025-2027)",
			})
		}
	}

	return errors
}

// validateTaskWarnings validates a task and returns non-critical warnings
func (v *CSVValidator) validateTaskWarnings(task Task, rowNum int) []ValidationIssue {
	var warnings []ValidationIssue

	// Warn about tasks without descriptions
	if strings.TrimSpace(task.Description) == "" {
		warnings = append(warnings, ValidationIssue{
			Type:    "missing_description",
			Field:   "Objective",
			Row:     rowNum,
			Message: "Task has no description/objective",
		})
	}

	// Warn about very short tasks (less than 1 day)
	if !task.StartDate.IsZero() && !task.EndDate.IsZero() {
		duration := task.EndDate.Sub(task.StartDate)
		if duration.Hours() < 24 && !task.IsMilestone {
			warnings = append(warnings, ValidationIssue{
				Type:    "short_duration",
				Field:   "End Date",
				Row:     rowNum,
				Message: "Task duration is less than 1 day, consider if this should be a milestone",
			})
		}
	}

	// Note: Assignee validation removed - acceptable for single-person projects

	return warnings
}

// validateDataConsistency validates overall data consistency across all tasks
func (v *CSVValidator) validateDataConsistency(tasks []Task) []ValidationIssue {
	var errors []ValidationIssue

	// Check for duplicate task IDs
	idMap := make(map[string][]int)
	for i, task := range tasks {
		if task.ID != "" {
			idMap[task.ID] = append(idMap[task.ID], i+2) // +2 for header + 0-indexing
		}
	}

	for id, rows := range idMap {
		if len(rows) > 1 {
			errors = append(errors, ValidationIssue{
				Type:    "duplicate_id",
				Field:   "Task ID",
				Value:   id,
				Message: fmt.Sprintf("Task ID '%s' appears in multiple rows: %v", id, rows),
			})
		}
	}

	// Validate dependency references exist and check for cycles
	taskIDSet := make(map[string]bool)
	taskIndex := make(map[string]int)
	for i, task := range tasks {
		if task.ID != "" {
			taskIDSet[task.ID] = true
			taskIndex[task.ID] = i
		}
	}

	for i, task := range tasks {
		if task.Dependencies != nil {
			for _, dep := range task.Dependencies {
				if !taskIDSet[dep] {
					errors = append(errors, ValidationIssue{
						Type:    "invalid_dependency",
						Field:   "Dependencies",
						Row:     i + 2, // +2 for header + 0-indexing
						Value:   dep,
						Message: fmt.Sprintf("Dependency '%s' references non-existent task ID", dep),
					})
				}
			}
		}
	}

	// Check for dependency cycles
	if cycleErrs := v.detectDependencyCycles(tasks, taskIndex); len(cycleErrs) > 0 {
		errors = append(errors, cycleErrs...)
	}

	// Validate milestone task consistency
	if milestoneErrs := v.validateMilestoneConsistency(tasks); len(milestoneErrs) > 0 {
		errors = append(errors, milestoneErrs...)
	}

	return errors
}

// detectDependencyCycles detects circular dependencies in the task graph
func (v *CSVValidator) detectDependencyCycles(tasks []Task, taskIndex map[string]int) []ValidationIssue {
	var errors []ValidationIssue
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	var visit func(taskID string, path []string) bool
	visit = func(taskID string, path []string) bool {
		// Check if we're already processing this node (cycle detected)
		if recursionStack[taskID] {
			// Find where the cycle starts
			cycleStart := -1
			for i, id := range path {
				if id == taskID {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 {
				cycle := append(path[cycleStart:], taskID)
				row := taskIndex[taskID] + 2 // +2 for header + 0-indexing
				errors = append(errors, ValidationIssue{
					Type:    "dependency_cycle",
					Field:   "Dependencies",
					Row:     row,
					Value:   strings.Join(cycle, " → "),
					Message: fmt.Sprintf("Circular dependency detected: %s", strings.Join(cycle, " → ")),
				})
			}
			return true
		}

		// Check if already fully visited
		if visited[taskID] {
			return false
		}

		// Mark as visiting
		recursionStack[taskID] = true
		visited[taskID] = true

		// Visit dependencies
		taskIdx, exists := taskIndex[taskID]
		if exists && tasks[taskIdx].Dependencies != nil {
			newPath := append(path, taskID)
			for _, dep := range tasks[taskIdx].Dependencies {
				if visit(dep, newPath) {
					return true // Cycle found
				}
			}
		}

		// Remove from recursion stack
		recursionStack[taskID] = false
		return false
	}

	// Check all tasks for cycles
	for _, task := range tasks {
		if task.ID != "" && !visited[task.ID] {
			visit(task.ID, []string{})
		}
	}

	return errors
}

// isMilestoneTask determines if a task is a milestone based on its name or description
func (v *CSVValidator) isMilestoneTask(name, description string) bool {
	text := strings.ToLower(name + " " + description)
	milestoneKeywords := []string{"milestone", "deadline", "due", "complete", "finish", "submit", "deliver", "defense", "exam", "oral", "final", "critical"}

	for _, keyword := range milestoneKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}

// validateMilestoneConsistency validates that milestone tasks are properly identified
func (v *CSVValidator) validateMilestoneConsistency(tasks []Task) []ValidationIssue {
	var errors []ValidationIssue

	for i, task := range tasks {
		// Check if task is marked as milestone but not identified as such by keywords
		if task.IsMilestone && !v.isMilestoneTask(task.Name, task.Description) {
			errors = append(errors, ValidationIssue{
				Type:    "milestone_inconsistency",
				Field:   "Milestone",
				Row:     i + 2, // +2 for header + 0-indexing
				Message: "Task is marked as milestone but name/description doesn't indicate milestone characteristics",
			})
		}
	}

	return errors
}

// getValidStatusesString returns a comma-separated string of valid statuses
func (v *CSVValidator) getValidStatusesString() string {
	statuses := make([]string, 0, len(v.validStatuses))
	for status := range v.validStatuses {
		statuses = append(statuses, status)
	}
	return strings.Join(statuses, ", ")
}

// * Configuration Validation
// ============================================================================

// ConfigValidator handles YAML configuration file validation
type ConfigValidator struct {
	logger *Logger
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		logger: NewDefaultLogger(),
	}
}

// ValidateConfigFile validates a YAML configuration file
func (cv *ConfigValidator) ValidateConfigFile(filePath string) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   make([]ValidationIssue, 0),
		Warnings: make([]ValidationIssue, 0),
	}

	// Check file access
	content, err := os.ReadFile(filePath)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Type:    "file_access",
			Message: fmt.Sprintf("Cannot read config file: %v", err),
		})
		return result, err
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(content, &config); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationIssue{
			Type:    "yaml_parsing",
			Message: fmt.Sprintf("Invalid YAML syntax: %v", err),
		})
		return result, err
	}

	// Validate configuration structure and values
	if errs := cv.validateConfigStructure(config); len(errs) > 0 {
		result.Errors = append(result.Errors, errs...)
		result.IsValid = false
	}

	if warns := cv.validateConfigWarnings(config); len(warns) > 0 {
		result.Warnings = append(result.Warnings, warns...)
	}

	return result, nil
}

// ValidateConfigFileContent validates a configuration struct directly
func (cv *ConfigValidator) ValidateConfigFileContent(config *Config) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   make([]ValidationIssue, 0),
		Warnings: make([]ValidationIssue, 0),
	}

	// Validate configuration structure and values
	if errs := cv.validateConfigStructure(*config); len(errs) > 0 {
		result.Errors = append(result.Errors, errs...)
		result.IsValid = false
	}

	if warns := cv.validateConfigWarnings(*config); len(warns) > 0 {
		result.Warnings = append(result.Warnings, warns...)
	}

	return result, nil
}

// validateConfigStructure validates the structure and required fields of configuration
func (cv *ConfigValidator) validateConfigStructure(config Config) []ValidationIssue {
	var errors []ValidationIssue

	// Validate required layout fields
	if config.Layout.Paper.Width == "" {
		errors = append(errors, ValidationIssue{
			Type:    "missing_required",
			Field:   "layout.paper.width",
			Message: "Paper width is required",
		})
	}

	if config.Layout.Paper.Height == "" {
		errors = append(errors, ValidationIssue{
			Type:    "missing_required",
			Field:   "layout.paper.height",
			Message: "Paper height is required",
		})
	}

	// Validate numeric fields
	if err := cv.validateNumericField(config.Layout.LayoutEngine.InitialYPositionMultiplier, "layout.layout_engine.initial_y_position_multiplier", 0.0, 1.0); err != nil {
		errors = append(errors, *err)
	}

	if err := cv.validateNumericField(config.Layout.LayoutEngine.TaskHeightMultiplier, "layout.layout_engine.task_height_multiplier", 0.0, 1.0); err != nil {
		errors = append(errors, *err)
	}

	if err := cv.validateNumericField(config.Layout.LayoutEngine.MaxTaskWidthDays, "layout.layout_engine.max_task_width_days", 1.0, 365.0); err != nil {
		errors = append(errors, *err)
	}

	// Validate grid constraints
	if err := cv.validateNumericField(config.Layout.LayoutEngine.GridConstraints.MinTaskSpacing, "layout.layout_engine.grid_constraints.min_task_spacing", 0.0, 100.0); err != nil {
		errors = append(errors, *err)
	}

	if err := cv.validateNumericField(config.Layout.LayoutEngine.GridConstraints.MaxTaskSpacing, "layout.layout_engine.grid_constraints.max_task_spacing", 0.0, 100.0); err != nil {
		errors = append(errors, *err)
	}

	// Validate spacing relationships
	if config.Layout.LayoutEngine.GridConstraints.MinTaskSpacing > config.Layout.LayoutEngine.GridConstraints.MaxTaskSpacing {
		errors = append(errors, ValidationIssue{
			Type:    "invalid_relationship",
			Field:   "layout.layout_engine.grid_constraints",
			Message: "min_task_spacing cannot be greater than max_task_spacing",
		})
	}

	// Validate pages configuration
	if len(config.Pages) == 0 {
		errors = append(errors, ValidationIssue{
			Type:    "missing_required",
			Field:   "pages",
			Message: "At least one page must be defined",
		})
	}

	for i, page := range config.Pages {
		if strings.TrimSpace(page.Name) == "" {
			errors = append(errors, ValidationIssue{
				Type:    "missing_required",
				Field:   fmt.Sprintf("pages[%d].name", i),
				Message: "Page name is required",
			})
		}

		if len(page.RenderBlocks) == 0 {
			errors = append(errors, ValidationIssue{
				Type:    "missing_required",
				Field:   fmt.Sprintf("pages[%d].renderblocks", i),
				Message: "At least one render block must be defined per page",
			})
		}

		for j, block := range page.RenderBlocks {
			if strings.TrimSpace(block.FuncName) == "" {
				errors = append(errors, ValidationIssue{
					Type:    "missing_required",
					Field:   fmt.Sprintf("pages[%d].renderblocks[%d].funcname", i, j),
					Message: "Function name is required for render block",
				})
			}

			if len(block.Tpls) == 0 {
				errors = append(errors, ValidationIssue{
					Type:    "missing_required",
					Field:   fmt.Sprintf("pages[%d].renderblocks[%d].tpls", i, j),
					Message: "At least one template must be specified for render block",
				})
			}
		}
	}

	return errors
}

// validateConfigWarnings validates configuration and returns non-critical warnings
func (cv *ConfigValidator) validateConfigWarnings(config Config) []ValidationIssue {
	var warnings []ValidationIssue

	// Warn about very small multipliers
	if config.Layout.LayoutEngine.TaskHeightMultiplier < 0.3 {
		warnings = append(warnings, ValidationIssue{
			Type:    "performance_warning",
			Field:   "layout.layout_engine.task_height_multiplier",
			Message: "Very small task height multiplier may cause layout issues",
		})
	}

	// Warn about very large max task width
	if config.Layout.LayoutEngine.MaxTaskWidthDays > 30 {
		warnings = append(warnings, ValidationIssue{
			Type:    "performance_warning",
			Field:   "layout.layout_engine.max_task_width_days",
			Message: "Very large max task width may cause calendar overflow",
		})
	}

	return warnings
}

// validateNumericField validates a numeric field is within acceptable range
func (cv *ConfigValidator) validateNumericField(value float64, fieldName string, min, max float64) *ValidationIssue {
	if value < min || value > max {
		return &ValidationIssue{
			Type:    "invalid_range",
			Field:   fieldName,
			Value:   fmt.Sprintf("%.2f", value),
			Message: fmt.Sprintf("Value must be between %.2f and %.2f", min, max),
		}
	}
	return nil
}

// * Validation Middleware for API-like Operations
// ============================================================================

// ValidationMiddleware provides validation for operations that process data
type ValidationMiddleware struct {
	csvValidator    *CSVValidator
	configValidator *ConfigValidator
	logger          *Logger
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		csvValidator:    NewCSVValidator(),
		configValidator: NewConfigValidator(),
		logger:          NewDefaultLogger(),
	}
}

// ValidateTaskOperation validates a task creation/modification operation
func (vm *ValidationMiddleware) ValidateTaskOperation(task *Task, operation string) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   make([]ValidationIssue, 0),
		Warnings: make([]ValidationIssue, 0),
	}

	// Run task validation
	if errs := vm.csvValidator.validateTask(*task, 0); len(errs) > 0 {
		result.Errors = append(result.Errors, errs...)
		result.IsValid = false
	}

	if warns := vm.csvValidator.validateTaskWarnings(*task, 0); len(warns) > 0 {
		result.Warnings = append(result.Warnings, warns...)
	}

	// Additional operation-specific validation
	switch operation {
	case "create":
		if errs := vm.validateTaskCreation(task); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.IsValid = false
		}
	case "update":
		if errs := vm.validateTaskUpdate(task); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.IsValid = false
		}
	case "delete":
		if errs := vm.validateTaskDeletion(task); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.IsValid = false
		}
	}

	return result, nil
}

// ValidateConfigOperation validates a configuration operation
func (vm *ValidationMiddleware) ValidateConfigOperation(config *Config, operation string) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   make([]ValidationIssue, 0),
		Warnings: make([]ValidationIssue, 0),
	}

	// Run configuration validation
	if errs := vm.configValidator.validateConfigStructure(*config); len(errs) > 0 {
		result.Errors = append(result.Errors, errs...)
		result.IsValid = false
	}

	if warns := vm.configValidator.validateConfigWarnings(*config); len(warns) > 0 {
		result.Warnings = append(result.Warnings, warns...)
	}

	// Additional operation-specific validation
	switch operation {
	case "load":
		if errs := vm.validateConfigLoad(config); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.IsValid = false
		}
	case "save":
		if errs := vm.validateConfigSave(config); len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			result.IsValid = false
		}
	}

	return result, nil
}

// validateTaskCreation validates task creation specific rules
func (vm *ValidationMiddleware) validateTaskCreation(task *Task) []ValidationIssue {
	var errors []ValidationIssue

	// Check for required fields for new tasks
	if task.ID == "" {
		errors = append(errors, ValidationIssue{
			Type:    "creation_required",
			Field:   "Task ID",
			Message: "Task ID is required when creating a new task",
		})
	}

	// Validate ID format (alphanumeric, dashes, underscores only)
	if task.ID != "" {
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, task.ID); !matched {
			errors = append(errors, ValidationIssue{
				Type:    "invalid_format",
				Field:   "Task ID",
				Value:   task.ID,
				Message: "Task ID must contain only letters, numbers, dashes, and underscores",
			})
		}
	}

	return errors
}

// validateTaskUpdate validates task update specific rules
func (vm *ValidationMiddleware) validateTaskUpdate(task *Task) []ValidationIssue {
	var errors []ValidationIssue

	// For updates, we might allow partial data, but some fields are immutable
	// Task ID should not change during update
	if task.ID == "" {
		errors = append(errors, ValidationIssue{
			Type:    "update_required",
			Field:   "Task ID",
			Message: "Task ID is required when updating an existing task",
		})
	}

	return errors
}

// validateTaskDeletion validates task deletion specific rules
func (vm *ValidationMiddleware) validateTaskDeletion(task *Task) []ValidationIssue {
	var errors []ValidationIssue

	// Check if task can be safely deleted (not referenced by other tasks)
	if task.ID == "" {
		errors = append(errors, ValidationIssue{
			Type:    "deletion_invalid",
			Field:   "Task ID",
			Message: "Cannot delete task without valid Task ID",
		})
	}

	return errors
}

// validateConfigLoad validates configuration loading
func (vm *ValidationMiddleware) validateConfigLoad(config *Config) []ValidationIssue {
	var errors []ValidationIssue

	// Validate year range
	if config.StartYear > 0 && config.EndYear > 0 && config.StartYear > config.EndYear {
		errors = append(errors, ValidationIssue{
			Type:    "invalid_range",
			Field:   "year_range",
			Message: "Start year cannot be greater than end year",
		})
	}

	// Validate output directory path
	if config.OutputDir != "" {
		if strings.Contains(config.OutputDir, "..") {
			errors = append(errors, ValidationIssue{
				Type:    "security_violation",
				Field:   "output_dir",
				Value:   config.OutputDir,
				Message: "Output directory path cannot contain '..' for security reasons",
			})
		}
	}

	return errors
}

// validateConfigSave validates configuration saving
func (vm *ValidationMiddleware) validateConfigSave(config *Config) []ValidationIssue {
	var errors []ValidationIssue

	// Similar to load validation but may have additional constraints
	loadErrors := vm.validateConfigLoad(config)
	errors = append(errors, loadErrors...)

	// Additional save-specific validations could go here

	return errors
}

// * Common Validation Types and Results
// ============================================================================

// ValidationIssue represents a single validation error or warning
type ValidationIssue struct {
	Type    string `json:"type"`            // Error type (required_field, invalid_value, etc.)
	Field   string `json:"field,omitempty"` // Field name that caused the error
	Row     int    `json:"row,omitempty"`   // Row number (for CSV validation)
	Value   string `json:"value,omitempty"` // Invalid value that caused the error
	Message string `json:"message"`         // Human-readable error message
}

func (ve ValidationIssue) Error() string {
	var parts []string

	if ve.Row > 0 {
		parts = append(parts, fmt.Sprintf("Row %d", ve.Row))
	}

	if ve.Field != "" {
		parts = append(parts, fmt.Sprintf("Field '%s'", ve.Field))
	}

	if ve.Value != "" {
		parts = append(parts, fmt.Sprintf("Value '%s'", ve.Value))
	}

	location := strings.Join(parts, ", ")
	if location != "" {
		return fmt.Sprintf("%s: %s", location, ve.Message)
	}

	return ve.Message
}

// ValidationResult contains the complete results of a validation operation
type ValidationResult struct {
	IsValid    bool              `json:"is_valid"`
	Errors     []ValidationIssue `json:"errors,omitempty"`
	Warnings   []ValidationIssue `json:"warnings,omitempty"`
	RowCount   int               `json:"row_count,omitempty"`   // For CSV validation
	FieldCount int               `json:"field_count,omitempty"` // For CSV validation
	Summary    string            `json:"summary,omitempty"`     // Human-readable summary
}

// Summary generates a human-readable summary of the validation results
func (vr *ValidationResult) GetSummary() string {
	if vr.Summary != "" {
		return vr.Summary
	}

	var summary strings.Builder

	if vr.IsValid && len(vr.Warnings) == 0 {
		summary.WriteString("✅ Validation successful")
		if vr.RowCount > 0 {
			summary.WriteString(fmt.Sprintf(" - %d rows validated", vr.RowCount))
		}
	} else {
		if !vr.IsValid {
			summary.WriteString(fmt.Sprintf("❌ Validation failed with %d errors", len(vr.Errors)))
		} else {
			summary.WriteString("⚠️  Validation passed with warnings")
		}

		if vr.RowCount > 0 {
			summary.WriteString(fmt.Sprintf(" - %d rows validated", vr.RowCount))
		}

		if len(vr.Warnings) > 0 {
			summary.WriteString(fmt.Sprintf(", %d warnings", len(vr.Warnings)))
		}
	}

	return summary.String()
}

// HasErrors returns true if there are any validation errors
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are any validation warnings
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// Error implements the error interface for ValidationResult
func (vr *ValidationResult) Error() string {
	if !vr.HasErrors() {
		return ""
	}

	if len(vr.Errors) == 1 {
		return vr.Errors[0].Error()
	}

	return fmt.Sprintf("%d validation errors occurred", len(vr.Errors))
}
