package application_test

import (
	"phd-dissertation-planner/internal/application"
	"testing"
)

func TestNew(t *testing.T) {
	app := application.New()
	if app == nil {
		t.Fatal("Expected app to be created, got nil")
	}
	
	if app.Name != "plannergen" {
		t.Errorf("Expected app name 'plannergen', got %s", app.Name)
	}
}

func TestRootFilename(t *testing.T) {
	// Test with .yaml extension
	filename := application.RootFilename("config.yaml")
	if filename != "config.tex" {
		t.Errorf("Expected 'config.tex', got %s", filename)
	}
	
	// Test with .yml extension
	filename = application.RootFilename("config.yml")
	if filename != "config.tex" {
		t.Errorf("Expected 'config.tex', got %s", filename)
	}
	
	// Test with no extension
	filename = application.RootFilename("config")
	if filename != "config.tex" {
		t.Errorf("Expected 'config.tex', got %s", filename)
	}
}