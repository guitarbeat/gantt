package common_test

import (
	"phd-dissertation-planner/internal/common"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	reader := common.NewReader("test.csv")
	if reader == nil {
		t.Fatal("Expected reader to be created, got nil")
	}
	
	// Test that reader was created successfully
	t.Log("Reader created successfully")
}

func TestParseDate(t *testing.T) {
	// Test date parsing directly
	date, err := time.Parse("2006-01-02", "2024-01-15")
	if err != nil {
		t.Fatalf("Failed to parse ISO date: %v", err)
	}
	
	expected := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !date.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, date)
	}
}
