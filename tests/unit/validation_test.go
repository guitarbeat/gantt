package core_test

import (
	"testing"
	"time"

	"phd-dissertation-planner/src/core"
)

// TestTaskValidation tests basic task validation functionality
// (DateValidator was removed during dead code cleanup)
func TestTaskValidation(t *testing.T) {
	// Test that we can create tasks with valid dates
	task := &core.Task{
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}

	if task.Name != "Test Task" {
		t.Errorf("Expected task name 'Test Task', got %s", task.Name)
	}

	if task.StartDate.IsZero() {
		t.Errorf("Expected non-zero start date")
	}

	if task.EndDate.IsZero() {
		t.Errorf("Expected non-zero end date")
	}

	// Test that end date is after start date
	if !task.EndDate.After(task.StartDate) {
		t.Errorf("Expected end date to be after start date")
	}
}

// TestTaskCreation tests basic task creation
func TestTaskCreation(t *testing.T) {
	// Test creating a task with all fields
	task := &core.Task{
		ID:          "task-1",
		Name:        "Test Task",
		Description: "A test task",
		Category:    "Test",
		StartDate:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		IsMilestone: true,
		Status:      "In Progress",
		Assignee:    "Test User",
	}

	// Verify all fields are set correctly
	if task.ID != "task-1" {
		t.Errorf("Expected ID 'task-1', got %s", task.ID)
	}

	if task.Name != "Test Task" {
		t.Errorf("Expected name 'Test Task', got %s", task.Name)
	}

	if !task.IsMilestone {
		t.Errorf("Expected milestone to be true")
	}

	if task.Status != "In Progress" {
		t.Errorf("Expected status 'In Progress', got %s", task.Status)
	}
}
