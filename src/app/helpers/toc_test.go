package helpers

import (
	"strings"
	"testing"
	"time"

	"phd-dissertation-planner/src/core"
)

func TestTOCBuilder_BuildTOCContent(t *testing.T) {
	tasks := []core.Task{
		{
			ID:          "1",
			Name:        "Task 1",
			Phase:       "1",
			SubPhase:    "Setup",
			Category:    "PROPOSAL",
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			IsMilestone: false,
		},
		{
			ID:          "2",
			Name:        "Milestone Task",
			Phase:       "1",
			SubPhase:    "Setup",
			Category:    "PROPOSAL",
			StartDate:   time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 1, 30, 0, 0, 0, 0, time.UTC),
			IsMilestone: true,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Phase:       "2",
			SubPhase:    "Research",
			Category:    "RESEARCH",
			StartDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
			IsMilestone: false,
		},
	}

	builder := NewTOCBuilder()
	content := builder.BuildTOCContent(tasks)

	// Check that content contains expected elements
	expectedElements := []string{
		"Table of Contents - Clickable Task Index",
		"\\hypertarget{task-index}{}",
		"Phase 1: Proposal \\& Setup",
		"Phase 2: Research \\& Data Collection",
		"Setup",
		"Research",
		"Task 1",
		"\\textbf{Milestone Task}", // Should be bold
		"Task 3",
		"\\hyperlink{", // Should contain hyperlinks
		"\\textcolor[RGB]{", // Should contain color formatting
		"How to use this index:",
		"\\pagebreak",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Expected content to contain '%s', but it didn't", element)
		}
	}

	// Check that milestone tasks are bold
	if !strings.Contains(content, "\\textbf{Milestone Task}") {
		t.Error("Expected milestone task to be bold")
	}

	// Check that regular tasks are not bold
	if strings.Contains(content, "\\textbf{Task 1}") {
		t.Error("Expected regular task not to be bold")
	}
}

func TestTOCBuilder_GroupTasksByPhase(t *testing.T) {
	tasks := []core.Task{
		{Phase: "1", Name: "Task 1"},
		{Phase: "1", Name: "Task 2"},
		{Phase: "2", Name: "Task 3"},
		{Phase: "3", Name: "Task 4"},
	}

	builder := NewTOCBuilder()
	phaseTasks := builder.groupTasksByPhase(tasks)

	// Check phase grouping
	if len(phaseTasks["1"]) != 2 {
		t.Errorf("Expected 2 tasks in phase 1, got %d", len(phaseTasks["1"]))
	}

	if len(phaseTasks["2"]) != 1 {
		t.Errorf("Expected 1 task in phase 2, got %d", len(phaseTasks["2"]))
	}

	if len(phaseTasks["3"]) != 1 {
		t.Errorf("Expected 1 task in phase 3, got %d", len(phaseTasks["3"]))
	}

	// Check that phase 4 is empty
	if len(phaseTasks["4"]) != 0 {
		t.Errorf("Expected 0 tasks in phase 4, got %d", len(phaseTasks["4"]))
	}
}

func TestTOCBuilder_GroupTasksBySubPhase(t *testing.T) {
	tasks := []core.Task{
		{Phase: "1", SubPhase: "Setup", Name: "Task 1"},
		{Phase: "1", SubPhase: "Setup", Name: "Task 2"},
		{Phase: "1", SubPhase: "Planning", Name: "Task 3"},
		{Phase: "1", SubPhase: "", Name: "Task 4"}, // Empty sub-phase should become "General"
	}

	builder := NewTOCBuilder()
	subPhaseTasks := builder.groupTasksBySubPhase(tasks)

	// Check sub-phase grouping
	if len(subPhaseTasks["Setup"]) != 2 {
		t.Errorf("Expected 2 tasks in Setup sub-phase, got %d", len(subPhaseTasks["Setup"]))
	}

	if len(subPhaseTasks["Planning"]) != 1 {
		t.Errorf("Expected 1 task in Planning sub-phase, got %d", len(subPhaseTasks["Planning"]))
	}

	if len(subPhaseTasks["General"]) != 1 {
		t.Errorf("Expected 1 task in General sub-phase, got %d", len(subPhaseTasks["General"]))
	}
}

func TestTOCBuilder_FormatTaskEntry(t *testing.T) {
	tests := []struct {
		name     string
		task     core.Task
		contains []string
	}{
		{
			name: "regular task",
			task: core.Task{
				Name:        "Regular Task",
				Category:    "PROPOSAL",
				StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				IsMilestone: false,
			},
			contains: []string{
				"\\item",
				"\\hyperlink{",
				"Regular Task",
				"\\textcolor[RGB]{",
			},
		},
		{
			name: "milestone task",
			task: core.Task{
				Name:        "Milestone Task",
				Category:    "PROPOSAL",
				StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				IsMilestone: true,
			},
			contains: []string{
				"\\item",
				"\\hyperlink{",
				"\\textbf{Milestone Task}",
				"\\textcolor[RGB]{",
			},
		},
		{
			name: "task with special characters",
			task: core.Task{
				Name:        "Task with & and % characters",
				Category:    "PROPOSAL",
				StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				IsMilestone: false,
			},
			contains: []string{
				"Task with \\& and \\% characters",
			},
		},
	}

	builder := NewTOCBuilder()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.formatTaskEntry(tt.task)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected task entry to contain '%s', but it didn't. Got: %s", expected, result)
				}
			}
		})
	}
}

func TestTOCBuilder_HexToRGBString(t *testing.T) {
	tests := []struct {
		hex      string
		expected string
	}{
		{"#FF0000", "255,0,0"},
		{"#00FF00", "0,255,0"},
		{"#0000FF", "0,0,255"},
		{"#FFFFFF", "255,255,255"},
		{"#000000", "0,0,0"},
		{"#123456", "18,52,86"},
		{"invalid", "0,0,0"},
		{"#GGGGGG", "0,0,0"},
		{"", "0,0,0"},
		{"#FF", "0,0,0"},
	}

	builder := NewTOCBuilder()

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			result := builder.hexToRGBString(tt.hex)
			if result != tt.expected {
				t.Errorf("For hex %s, expected %s, got %s", tt.hex, tt.expected, result)
			}
		})
	}
}

func TestTOCBuilder_EscapeLaTeXSpecialChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Normal text", "Normal text"},
		{"Text with & ampersand", "Text with \\& ampersand"},
		{"Text with % percent", "Text with \\% percent"},
		{"Text with & and %", "Text with \\& and \\%"},
		{"", ""},
	}

	builder := NewTOCBuilder()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := builder.escapeLaTeXSpecialChars(tt.input)
			if result != tt.expected {
				t.Errorf("For input '%s', expected '%s', got '%s'", tt.input, tt.expected, result)
			}
		})
	}
}

func TestParseHexByte(t *testing.T) {
	tests := []struct {
		hex      string
		expected int64
		hasError bool
	}{
		{"FF", 255, false},
		{"00", 0, false},
		{"1A", 26, false},
		{"a5", 165, false},
		{"GG", 0, true},
		{"", 0, true},
		{"F", 0, true},
		{"FFF", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			result, err := parseHexByte(tt.hex)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for hex '%s', but got none", tt.hex)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for hex '%s': %v", tt.hex, err)
				return
			}

			if result != tt.expected {
				t.Errorf("For hex '%s', expected %d, got %d", tt.hex, tt.expected, result)
			}
		})
	}
}
