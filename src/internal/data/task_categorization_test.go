package data

import (
	"testing"
)

func TestTaskCategoryManager_GetCategory(t *testing.T) {
	manager := NewTaskCategoryManager()
	
	// Test valid categories
	validCategories := []string{
		"PROPOSAL", "LASER", "IMAGING", "ADMIN", 
		"DISSERTATION", "RESEARCH", "PUBLICATION",
	}
	
	for _, category := range validCategories {
		cat, exists := manager.GetCategory(category)
		if !exists {
			t.Errorf("Expected category %s to exist", category)
		}
		if cat.Name != category {
			t.Errorf("Expected category name %s, got %s", category, cat.Name)
		}
	}
	
	// Test invalid category
	_, exists := manager.GetCategory("INVALID")
	if exists {
		t.Error("Expected invalid category to not exist")
	}
}

func TestTaskCategoryManager_CategorizeTask(t *testing.T) {
	manager := NewTaskCategoryManager()
	
	// Test task with explicit category
	task := &Task{
		Name:     "Test Task",
		Category: "PROPOSAL",
	}
	
	category := manager.CategorizeTask(task)
	if category != "PROPOSAL" {
		t.Errorf("Expected PROPOSAL, got %s", category)
	}
	
	// Test task without category (should auto-categorize)
	task = &Task{
		Name:        "Write dissertation chapter",
		Description: "Complete chapter 3 of dissertation",
	}
	
	category = manager.CategorizeTask(task)
	if category == "" {
		t.Error("Expected auto-categorized category")
	}
}

func TestTaskCategoryManager_GetAllCategories(t *testing.T) {
	manager := NewTaskCategoryManager()
	
	categories := manager.GetAllCategories()
	if len(categories) == 0 {
		t.Error("Expected categories to be returned")
	}
	
	// Check for expected categories
	expectedCategories := []string{
		"PROPOSAL", "LASER", "IMAGING", "ADMIN", 
		"DISSERTATION", "RESEARCH", "PUBLICATION",
	}
	
	for _, expected := range expectedCategories {
		found := false
		for _, cat := range categories {
			if cat.Name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected category %s to be in all categories", expected)
		}
	}
}

func TestTaskCategoryManager_GetCategoryColor(t *testing.T) {
	manager := NewTaskCategoryManager()
	
	// Test valid category color
	color, exists := manager.GetCategoryColor("PROPOSAL")
	if !exists {
		t.Error("Expected PROPOSAL category color to exist")
	}
	if color == "" {
		t.Error("Expected non-empty color")
	}
	
	// Test invalid category color
	_, exists = manager.GetCategoryColor("INVALID")
	if exists {
		t.Error("Expected invalid category color to not exist")
	}
}

func TestTaskCategoryManager_GetCategoryDescription(t *testing.T) {
	manager := NewTaskCategoryManager()
	
	// Test valid category description
	desc, exists := manager.GetCategoryDescription("PROPOSAL")
	if !exists {
		t.Error("Expected PROPOSAL category description to exist")
	}
	if desc == "" {
		t.Error("Expected non-empty description")
	}
	
	// Test invalid category description
	_, exists = manager.GetCategoryDescription("INVALID")
	if exists {
		t.Error("Expected invalid category description to not exist")
	}
}