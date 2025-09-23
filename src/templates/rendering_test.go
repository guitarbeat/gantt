package templates

import (
	"testing"
	"time"
	"phd-dissertation-planner/internal/common"
)

func TestNewRenderer(t *testing.T) {
	renderer := NewRenderer()
	if renderer == nil {
		t.Fatal("Expected renderer to be created, got nil")
	}
	
	if renderer.templates == nil {
		t.Error("Expected templates to be initialized")
	}
}

func TestRenderMonthly(t *testing.T) {
	renderer := NewRenderer()
	
	// Create test data
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	task := &common.Task{
		ID:        "1",
		Name:      "Test Task",
		StartDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	
	calendar := &common.CalendarLayout{
		startDate: start,
		endDate:   end,
		tasks:     []*common.Task{task},
	}
	
	// Test rendering
	output, err := renderer.RenderMonthly(calendar)
	if err != nil {
		t.Fatalf("Failed to render monthly calendar: %v", err)
	}
	
	if output == "" {
		t.Error("Expected non-empty output")
	}
	
	// Check if task name appears in output
	if !contains(output, "Test Task") {
		t.Error("Expected task name to appear in output")
	}
}

func TestRenderDocument(t *testing.T) {
	renderer := NewRenderer()
	
	// Create test data
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	
	calendar := &common.CalendarLayout{
		startDate: start,
		endDate:   end,
		tasks:     []*common.Task{},
	}
	
	// Test rendering
	output, err := renderer.RenderDocument(calendar)
	if err != nil {
		t.Fatalf("Failed to render document: %v", err)
	}
	
	if output == "" {
		t.Error("Expected non-empty output")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && contains(s[1:], substr)
}
