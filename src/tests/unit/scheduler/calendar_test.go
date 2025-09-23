package scheduler_test

import (
	"phd-dissertation-planner/internal/scheduler"
	"testing"
	"time"
)

func TestNewYear(t *testing.T) {
	year := scheduler.NewYear(time.Monday, 2024)
	if year == nil {
		t.Fatal("Expected year to be created, got nil")
	}
	
	if year.Number != 2024 {
		t.Errorf("Expected year 2024, got %d", year.Number)
	}
	
	if len(year.Quarters) != 4 {
		t.Errorf("Expected 4 quarters, got %d", len(year.Quarters))
	}
}

func TestDay(t *testing.T) {
	day := scheduler.Day{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)}
	
	if day.Time.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", day.Time.Year())
	}
	
	if day.Time.Month() != time.January {
		t.Errorf("Expected January, got %s", day.Time.Month())
	}
	
	if day.Time.Day() != 15 {
		t.Errorf("Expected day 15, got %d", day.Time.Day())
	}
}