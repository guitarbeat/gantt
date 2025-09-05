package data

import (
	"fmt"
	"testing"
	"time"
)

func TestDependencyValidator(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Test adding tasks
	task1 := &Task{
		ID:          "A",
		Name:        "Task A",
		Dependencies: []string{},
	}
	
	task2 := &Task{
		ID:          "B",
		Name:        "Task B",
		Dependencies: []string{"A"},
	}
	
	dv.AddTask(task1)
	dv.AddTask(task2)
	
	// Test getting dependencies
	deps := dv.GetTaskDependencies("B")
	if len(deps) != 1 {
		t.Errorf("Expected 1 dependency for B, got %d", len(deps))
	}
	if deps[0].ID != "A" {
		t.Errorf("Expected dependency A for B, got %s", deps[0].ID)
	}
	
	// Test getting dependents
	dependents := dv.GetTaskDependents("A")
	if len(dependents) != 1 {
		t.Errorf("Expected 1 dependent for A, got %d", len(dependents))
	}
	if dependents[0].ID != "B" {
		t.Errorf("Expected dependent B for A, got %s", dependents[0].ID)
	}
}

func TestValidateTaskDependenciesDetailed(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Test self-dependency
	selfDepTask := &Task{
		ID:          "A",
		Name:        "Self Dep Task",
		Dependencies: []string{"A"},
	}
	
	dv.AddTask(selfDepTask)
	errors := dv.ValidateDependencies()
	
	hasSelfDepError := false
	for _, err := range errors {
		if err.Type == "DEPENDENCY" && err.Field == "Dependencies" && err.TaskID == "A" {
			hasSelfDepError = true
			break
		}
	}
	if !hasSelfDepError {
		t.Error("Expected self-dependency error")
	}
	
	// Test duplicate dependencies
	dupDepTask := &Task{
		ID:          "B",
		Name:        "Dup Dep Task",
		Dependencies: []string{"A", "A"},
	}
	
	dv.AddTask(dupDepTask)
	errors = dv.ValidateDependencies()
	
	hasDupDepError := false
	for _, err := range errors {
		if err.Type == "DEPENDENCY" && err.Field == "Dependencies" && err.TaskID == "B" {
			hasDupDepError = true
			break
		}
	}
	if !hasDupDepError {
		t.Error("Expected duplicate dependency error")
	}
	
	// Test non-existent dependency
	nonExistentDepTask := &Task{
		ID:          "C",
		Name:        "Non Existent Dep Task",
		Dependencies: []string{"X"},
	}
	
	dv.AddTask(nonExistentDepTask)
	errors = dv.ValidateDependencies()
	
	hasNonExistentError := false
	for _, err := range errors {
		if err.Type == "DEPENDENCY" && err.Field == "Dependencies" && err.TaskID == "C" {
			hasNonExistentError = true
			break
		}
	}
	if !hasNonExistentError {
		t.Error("Expected non-existent dependency error")
	}
	
	// Test too many dependencies
	manyDeps := make([]string, 15)
	for i := 0; i < 15; i++ {
		manyDeps[i] = fmt.Sprintf("Dep%d", i)
	}
	
	manyDepTask := &Task{
		ID:          "D",
		Name:        "Many Dep Task",
		Dependencies: manyDeps,
	}
	
	dv.AddTask(manyDepTask)
	errors = dv.ValidateDependencies()
	
	hasManyDepError := false
	for _, err := range errors {
		if err.Type == "DEPENDENCY" && err.Field == "Dependencies" && err.TaskID == "D" {
			hasManyDepError = true
			break
		}
	}
	if !hasManyDepError {
		t.Error("Expected too many dependencies error")
	}
}

func TestDetectCircularDependenciesDetailed(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Test simple circular dependency
	taskA := &Task{
		ID:          "A",
		Name:        "Task A",
		Dependencies: []string{"B"},
	}
	
	taskB := &Task{
		ID:          "B",
		Name:        "Task B",
		Dependencies: []string{"A"},
	}
	
	dv.AddTask(taskA)
	dv.AddTask(taskB)
	
	errors := dv.ValidateDependencies()
	
	hasCircularError := false
	for _, err := range errors {
		if err.Type == "CIRCULAR_DEPENDENCY" {
			hasCircularError = true
			break
		}
	}
	if !hasCircularError {
		t.Error("Expected circular dependency error")
	}
	
	// Test complex circular dependency
	dv2 := NewDependencyValidator()
	
	task1 := &Task{
		ID:          "1",
		Name:        "Task 1",
		Dependencies: []string{"2"},
	}
	
	task2 := &Task{
		ID:          "2",
		Name:        "Task 2",
		Dependencies: []string{"3"},
	}
	
	task3 := &Task{
		ID:          "3",
		Name:        "Task 3",
		Dependencies: []string{"1"},
	}
	
	dv2.AddTask(task1)
	dv2.AddTask(task2)
	dv2.AddTask(task3)
	
	errors2 := dv2.ValidateDependencies()
	
	hasComplexCircularError := false
	for _, err := range errors2 {
		if err.Type == "CIRCULAR_DEPENDENCY" {
			hasComplexCircularError = true
			break
		}
	}
	if !hasComplexCircularError {
		t.Error("Expected complex circular dependency error")
	}
}

func TestDetectDependencyConflicts(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Test mutual dependency
	taskA := &Task{
		ID:          "A",
		Name:        "Task A",
		Dependencies: []string{"B"},
		StartDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	taskB := &Task{
		ID:          "B",
		Name:        "Task B",
		Dependencies: []string{"A"},
		StartDate:    time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	dv.AddTask(taskA)
	dv.AddTask(taskB)
	
	errors := dv.ValidateDependencies()
	
	hasMutualDepError := false
	for _, err := range errors {
		if err.Type == "DEPENDENCY_CONFLICT" {
			hasMutualDepError = true
			break
		}
	}
	if !hasMutualDepError {
		t.Error("Expected mutual dependency error")
	}
	
	// Test dependency with later start date
	dv2 := NewDependencyValidator()
	
	taskC := &Task{
		ID:          "C",
		Name:        "Task C",
		Dependencies: []string{"D"},
		StartDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	
	taskD := &Task{
		ID:          "D",
		Name:        "Task D",
		Dependencies: []string{},
		StartDate:    time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	
	dv2.AddTask(taskC)
	dv2.AddTask(taskD)
	
	errors2 := dv2.ValidateDependencies()
	
	hasLaterStartError := false
	for _, err := range errors2 {
		if err.Type == "DEPENDENCY_CONFLICT" {
			hasLaterStartError = true
			break
		}
	}
	if !hasLaterStartError {
		t.Error("Expected later start date error")
	}
}

func TestDetectOrphanedDependencies(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Test orphaned task
	taskA := &Task{
		ID:          "A",
		Name:        "Task A",
		Dependencies: []string{"B"},
	}
	
	taskB := &Task{
		ID:          "B",
		Name:        "Task B",
		Dependencies: []string{},
	}
	
	dv.AddTask(taskA)
	dv.AddTask(taskB)
	
	errors := dv.ValidateDependencies()
	
	hasOrphanedError := false
	for _, err := range errors {
		if err.Type == "ORPHANED_DEPENDENCY" {
			hasOrphanedError = true
			break
		}
	}
	if !hasOrphanedError {
		t.Error("Expected orphaned dependency error")
	}
}

func TestDetectLongDependencyChains(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Create a long dependency chain
	tasks := []*Task{
		{ID: "1", Name: "Task 1", Dependencies: []string{}},
		{ID: "2", Name: "Task 2", Dependencies: []string{"1"}},
		{ID: "3", Name: "Task 3", Dependencies: []string{"2"}},
		{ID: "4", Name: "Task 4", Dependencies: []string{"3"}},
		{ID: "5", Name: "Task 5", Dependencies: []string{"4"}},
		{ID: "6", Name: "Task 6", Dependencies: []string{"5"}},
		{ID: "7", Name: "Task 7", Dependencies: []string{"6"}},
	}
	
	for _, task := range tasks {
		dv.AddTask(task)
	}
	
	errors := dv.ValidateDependencies()
	
	hasLongChainError := false
	for _, err := range errors {
		if err.Type == "LONG_DEPENDENCY_CHAIN" {
			hasLongChainError = true
			break
		}
	}
	if !hasLongChainError {
		t.Error("Expected long dependency chain error")
	}
}

func TestGetDependencyLevels(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Create a dependency hierarchy
	tasks := []*Task{
		{ID: "A", Name: "Task A", Dependencies: []string{}},
		{ID: "B", Name: "Task B", Dependencies: []string{"A"}},
		{ID: "C", Name: "Task C", Dependencies: []string{"A"}},
		{ID: "D", Name: "Task D", Dependencies: []string{"B", "C"}},
		{ID: "E", Name: "Task E", Dependencies: []string{"D"}},
	}
	
	for _, task := range tasks {
		dv.AddTask(task)
	}
	
	levels := dv.GetDependencyLevels()
	
	// Check level 0 (root tasks)
	if len(levels[0]) != 1 || levels[0][0] != "A" {
		t.Errorf("Expected level 0 to have 1 task (A), got %v", levels[0])
	}
	
	// Check level 1
	if len(levels[1]) != 2 {
		t.Errorf("Expected level 1 to have 2 tasks, got %d", len(levels[1]))
	}
	
	// Check level 2
	if len(levels[2]) != 1 || levels[2][0] != "D" {
		t.Errorf("Expected level 2 to have 1 task (D), got %v", levels[2])
	}
	
	// Check level 3
	if len(levels[3]) != 1 || levels[3][0] != "E" {
		t.Errorf("Expected level 3 to have 1 task (E), got %v", levels[3])
	}
}

func TestGetTasksByDependencyLevel(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Create a dependency hierarchy
	tasks := []*Task{
		{ID: "A", Name: "Task A", Dependencies: []string{}},
		{ID: "B", Name: "Task B", Dependencies: []string{"A"}},
		{ID: "C", Name: "Task C", Dependencies: []string{"A"}},
		{ID: "D", Name: "Task D", Dependencies: []string{"B", "C"}},
	}
	
	for _, task := range tasks {
		dv.AddTask(task)
	}
	
	// Test level 0
	level0Tasks := dv.GetTasksByDependencyLevel(0)
	if len(level0Tasks) != 1 || level0Tasks[0].ID != "A" {
		t.Errorf("Expected level 0 to have 1 task (A), got %v", level0Tasks)
	}
	
	// Test level 1
	level1Tasks := dv.GetTasksByDependencyLevel(1)
	if len(level1Tasks) != 2 {
		t.Errorf("Expected level 1 to have 2 tasks, got %d", len(level1Tasks))
	}
	
	// Test level 2
	level2Tasks := dv.GetTasksByDependencyLevel(2)
	if len(level2Tasks) != 1 || level2Tasks[0].ID != "D" {
		t.Errorf("Expected level 2 to have 1 task (D), got %v", level2Tasks)
	}
}

func TestValidateDependencyIntegrity(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Create tasks with various dependency issues
	tasks := []*Task{
		{ID: "A", Name: "Task A", Dependencies: []string{}},
		{ID: "B", Name: "Task B", Dependencies: []string{"A"}},
		{ID: "C", Name: "Task C", Dependencies: []string{"X"}}, // Non-existent dependency
		{ID: "D", Name: "Task D", Dependencies: []string{"D"}}, // Self-dependency
	}
	
	for _, task := range tasks {
		dv.AddTask(task)
	}
	
	result := dv.ValidateDependencyIntegrity()
	
	if result.IsValid {
		t.Error("Expected validation to fail due to errors")
	}
	
	if result.ErrorCount == 0 {
		t.Error("Expected errors to be found")
	}
	
	if result.TaskCount != 4 {
		t.Errorf("Expected 4 tasks, got %d", result.TaskCount)
	}
	
	// Check that we have the expected error types
	errorTypes := make(map[string]int)
	for _, err := range result.Errors {
		errorTypes[err.Type]++
	}
	
	if errorTypes["DEPENDENCY"] == 0 {
		t.Error("Expected DEPENDENCY errors")
	}
}

func TestDependencyValidatorIntegration(t *testing.T) {
	dv := NewDependencyValidator()
	
	// Create a complex dependency scenario
	tasks := []*Task{
		{
			ID:          "A",
			Name:        "Root Task",
			Dependencies: []string{},
			StartDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          "B",
			Name:        "Dependent Task",
			Dependencies: []string{"A"},
			StartDate:    time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          "C",
			Name:        "Circular Task",
			Dependencies: []string{"B"},
			StartDate:    time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          "D",
			Name:        "Invalid Dep Task",
			Dependencies: []string{"X"},
			StartDate:    time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		},
	}
	
	// Add tasks
	for _, task := range tasks {
		dv.AddTask(task)
	}
	
	// Create circular dependency
	dv.tasks["B"].Dependencies = append(dv.tasks["B"].Dependencies, "C")
	dv.graph["B"] = dv.tasks["B"].Dependencies
	dv.reverse["C"] = append(dv.reverse["C"], "B")
	
	// Validate dependencies
	errors := dv.ValidateDependencies()
	
	if len(errors) == 0 {
		t.Error("Expected validation errors")
	}
	
	// Check for specific error types
	errorTypes := make(map[string]int)
	for _, err := range errors {
		errorTypes[err.Type]++
	}
	
	// Debug: print all errors
	for _, err := range errors {
		t.Logf("Error: %s - %s", err.Type, err.Message)
	}
	
	if errorTypes["CIRCULAR_DEPENDENCY"] == 0 {
		t.Error("Expected circular dependency error")
	}
	
	if errorTypes["DEPENDENCY"] == 0 {
		t.Error("Expected dependency error")
	}
	
	// Test dependency levels
	levels := dv.GetDependencyLevels()
	if len(levels) == 0 {
		t.Error("Expected dependency levels")
	}
	
	// Test getting tasks by level
	rootTasks := dv.GetTasksByDependencyLevel(0)
	if len(rootTasks) == 0 {
		t.Error("Expected root tasks")
	}
}
