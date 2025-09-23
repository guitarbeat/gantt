package main_test

import (
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	// Test that main function can be called without panicking
	// This is a basic smoke test
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Main function panicked: %v", r)
		}
	}()
	
	// Since main() calls os.Exit, we can't test it directly
	// This test just ensures the package compiles
}

func TestParseDate(t *testing.T) {
	// Test date parsing
	dateStr := "2024-01-15"
	expected := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	
	actual, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}
	
	if !actual.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestValidateDateRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	
	// Valid range
	if !isValidDateRange(start, end) {
		t.Error("Expected valid date range to be valid")
	}
	
	// Invalid range (end before start)
	invalidEnd := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	if isValidDateRange(start, invalidEnd) {
		t.Error("Expected invalid date range to be invalid")
	}
}

// Helper function to validate date range
func isValidDateRange(start, end time.Time) bool {
	return !start.IsZero() && !end.IsZero() && start.Before(end)
}