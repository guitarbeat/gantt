package core

import (
	"strings"
	"testing"
	"time"
)

func TestSuggestCorrection(t *testing.T) {
	options := []string{
		"planned",
		"in progress",
		"completed",
		"cancelled",
		"on hold",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"planed", "planned"},
		{"plannd", "planned"},
		{"compelted", "completed"},
		{"cancel", "cancelled"},
		{"in prog", "in progress"},
		{"onhold", "on hold"}, // distance 1
		{"xyz", ""},           // too far
		{"", ""},
		{"completed", "completed"}, // exact match
	}

	for _, tt := range tests {
		result := SuggestCorrection(tt.input, options)
		if result != tt.expected {
			t.Errorf("SuggestCorrection(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		s1, s2   string
		expected int
	}{
		{"kitten", "sitting", 3},
		{"flaw", "lawn", 2},
		{"saturday", "sunday", 3},
		{"", "abc", 3},
		{"abc", "", 3},
		{"abc", "abc", 0},
	}

	for _, tt := range tests {
		result := levenshtein(tt.s1, tt.s2)
		if result != tt.expected {
			t.Errorf("levenshtein(%q, %q) = %d, want %d", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestCSVValidator_ValidationSuggestions(t *testing.T) {
	validator := NewCSVValidator()

	// Create a task with an invalid status that looks like a valid one
	task := Task{
		ID:        "T1.1",
		Name:      "Test Task",
		Status:    "planed", // Should be "planned"
		StartDate: parseDateOrPanic("2025-01-01"),
		EndDate:   parseDateOrPanic("2025-01-02"),
	}

	issues := validator.validateTask(task, 1)

	foundSuggestion := false
	for _, issue := range issues {
		if issue.Field == "Status" && strings.Contains(issue.Message, "Did you mean 'planned'?") {
			foundSuggestion = true
			break
		}
	}

	if !foundSuggestion {
		// Print all messages to see what happened
		var msgs []string
		for _, issue := range issues {
			msgs = append(msgs, issue.Message)
		}
		t.Errorf("Expected validation error to contain suggestion 'planned', but got messages: %v", msgs)
	}
}

// Helper to parse date for tests
func parseDateOrPanic(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
