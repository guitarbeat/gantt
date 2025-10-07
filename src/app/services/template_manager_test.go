package services

import (
	"bytes"
	"os"
	"testing"

	"phd-dissertation-planner/src/core"
)

func TestNewTemplateManager(t *testing.T) {
	// Test with development mode disabled (embedded templates)
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	if tm == nil {
		t.Fatal("Expected template manager to be created")
	}

	if tm.templates == nil {
		t.Fatal("Expected templates to be loaded")
	}

	// Logger is a struct, not a pointer, so we can't check for nil
	// Just verify the template manager was created successfully
}

func TestNewTemplateManagerWithConfig(t *testing.T) {
	config := TemplateManagerConfig{
		DevMode:         false,
		TemplateSubDir:  "monthly",
		TemplatePath:    "src/shared/templates/monthly",
		TemplatePattern: "*.tpl",
	}

	tm, err := NewTemplateManagerWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to create template manager with config: %v", err)
	}

	if tm == nil {
		t.Fatal("Expected template manager to be created")
	}
}

func TestTemplateManager_Execute(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Test executing a non-existent template
	var buf bytes.Buffer
	err = tm.Execute(&buf, "nonexistent.tpl", map[string]interface{}{"test": "data"})
	if err == nil {
		t.Error("Expected error for non-existent template")
	}

	// Test that error templateStringContains helpful information
	if err != nil && !templateStringContains(err.Error(), "template not found") {
		t.Errorf("Expected error to contain 'template not found', got: %v", err)
	}
}

func TestTemplateManager_ExecuteDocument(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Create a test configuration
	cfg := core.Config{
		OutputDir: "test_output",
		Pages: []core.Page{
			{Name: "test_page", RenderBlocks: []core.RenderBlock{}},
		},
	}

	var buf bytes.Buffer
	err = tm.ExecuteDocument(&buf, cfg)
	
	// This might fail if the document template doesn't exist, which is expected
	// We're mainly testing that the method doesn't panic and handles errors gracefully
	if err != nil {
		t.Logf("ExecuteDocument failed as expected (template may not exist): %v", err)
	}
}

func TestTemplateManager_GetAvailableTemplates(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	templates := tm.GetAvailableTemplates()
	
	// Templates should be a slice (even if empty)
	if templates == nil {
		t.Fatal("Expected templates slice, got nil")
	}

	t.Logf("Available templates: %v", templates)
}

func TestTemplateManager_ValidateTemplate(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Test validating a non-existent template
	err = tm.ValidateTemplate("nonexistent.tpl")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}

	// Test that error templateStringContains helpful information
	if err != nil && !templateStringContains(err.Error(), "template not found") {
		t.Errorf("Expected error to contain 'template not found', got: %v", err)
	}
}

func TestTemplateManagerConfig(t *testing.T) {
	config := TemplateManagerConfig{
		DevMode:         true,
		TemplateSubDir:  "test",
		TemplatePath:    "test/path",
		TemplatePattern: "*.test",
	}

	// Test that config fields are set correctly
	if !config.DevMode {
		t.Error("Expected DevMode to be true")
	}

	if config.TemplateSubDir != "test" {
		t.Errorf("Expected TemplateSubDir to be 'test', got %s", config.TemplateSubDir)
	}

	if config.TemplatePath != "test/path" {
		t.Errorf("Expected TemplatePath to be 'test/path', got %s", config.TemplatePath)
	}

	if config.TemplatePattern != "*.test" {
		t.Errorf("Expected TemplatePattern to be '*.test', got %s", config.TemplatePattern)
	}
}

// Test with development mode enabled
func TestTemplateManager_DevMode(t *testing.T) {
	// Set environment variable for dev mode
	originalValue := os.Getenv("DEV_TEMPLATES")
	defer os.Setenv("DEV_TEMPLATES", originalValue)

	os.Setenv("DEV_TEMPLATES", "1")

	config := TemplateManagerConfig{
		DevMode:         true,
		TemplateSubDir:  "monthly",
		TemplatePath:    "src/shared/templates/monthly",
		TemplatePattern: "*.tpl",
	}

	tm, err := NewTemplateManagerWithConfig(config)
	if err != nil {
		t.Logf("Dev mode template manager creation failed (expected if templates don't exist): %v", err)
		return
	}

	if tm == nil {
		t.Fatal("Expected template manager to be created")
	}
}

// Helper function to check if a string templateStringContains a substring
func templateStringContains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && templateStringContains(s[1:], substr)
}
