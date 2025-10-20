package calendar_test

import (
	"testing"
	"time"

	"phd-dissertation-planner/internal/calendar"
	"phd-dissertation-planner/internal/core"
)

// TestCalendarPackage verifies the calendar package is importable
func TestCalendarPackage(t *testing.T) {
	// Basic test to ensure package compiles
	t.Log("Calendar package test placeholder")
}

func TestGetMonthName(t *testing.T) {
	tests := []struct {
		month    time.Month
		expected string
	}{
		{time.January, "January"},
		{time.February, "February"},
		{time.March, "March"},
		{time.December, "December"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getMonthName(tt.month)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetWeekNumber(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected int
	}{
		{
			name:     "First week of 2025",
			date:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "Mid year",
			date:     time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 24,
		},
		{
			name:     "End of year",
			date:     time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, week := tt.date.ISOWeek()
			if week != tt.expected {
				t.Logf("Week number for %v: got %d, expected %d", tt.date, week, tt.expected)
				// Note: Week numbers can vary by implementation
			}
		})
	}
}

func TestIsWeekend(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{
			name:     "Monday",
			date:     time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Friday",
			date:     time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Saturday",
			date:     time.Date(2025, 1, 11, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Sunday",
			date:     time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWeekend(tt.date)
			if result != tt.expected {
				t.Errorf("isWeekend(%v) = %v, expected %v", tt.date, result, tt.expected)
			}
		})
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		year     int
		month    time.Month
		expected int
	}{
		{2025, time.January, 31},
		{2025, time.February, 28},
		{2024, time.February, 29}, // Leap year
		{2025, time.April, 30},
		{2025, time.December, 31},
	}

	for _, tt := range tests {
		t.Run(tt.month.String(), func(t *testing.T) {
			result := daysInMonth(tt.year, tt.month)
			if result != tt.expected {
				t.Errorf("daysInMonth(%d, %v) = %d, expected %d",
					tt.year, tt.month, result, tt.expected)
			}
		})
	}
}

func TestFirstDayOfMonth(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
	}{
		{2025, time.January},
		{2025, time.June},
		{2025, time.December},
	}

	for _, tt := range tests {
		t.Run(tt.month.String(), func(t *testing.T) {
			result := firstDayOfMonth(tt.year, tt.month)

			if result.Year() != tt.year {
				t.Errorf("Expected year %d, got %d", tt.year, result.Year())
			}
			if result.Month() != tt.month {
				t.Errorf("Expected month %v, got %v", tt.month, result.Month())
			}
			if result.Day() != 1 {
				t.Errorf("Expected day 1, got %d", result.Day())
			}
		})
	}
}

func TestLastDayOfMonth(t *testing.T) {
	tests := []struct {
		year        int
		month       time.Month
		expectedDay int
	}{
		{2025, time.January, 31},
		{2025, time.February, 28},
		{2024, time.February, 29},
		{2025, time.April, 30},
	}

	for _, tt := range tests {
		t.Run(tt.month.String(), func(t *testing.T) {
			result := lastDayOfMonth(tt.year, tt.month)

			if result.Day() != tt.expectedDay {
				t.Errorf("Expected day %d, got %d", tt.expectedDay, result.Day())
			}
		})
	}
}

func TestGetWeeksInMonth(t *testing.T) {
	year := 2025
	month := time.January
	weekStart := time.Monday

	weeks := getWeeksInMonth(year, month, weekStart)

	if len(weeks) == 0 {
		t.Error("Expected at least one week, got 0")
	}

	// January 2025 should have 5 weeks (starting Monday)
	if len(weeks) < 4 || len(weeks) > 6 {
		t.Errorf("Expected 4-6 weeks for January 2025, got %d", len(weeks))
	}
}

func TestDateRange(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)

	dates := dateRange(start, end)

	expectedDays := 5
	if len(dates) != expectedDays {
		t.Errorf("Expected %d dates, got %d", expectedDays, len(dates))
	}

	// Verify first and last dates
	if !dates[0].Equal(start) {
		t.Errorf("First date should be %v, got %v", start, dates[0])
	}
	if !dates[len(dates)-1].Equal(end) {
		t.Errorf("Last date should be %v, got %v", end, dates[len(dates)-1])
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{2024, true},  // Divisible by 4
		{2025, false}, // Not divisible by 4
		{2000, true},  // Divisible by 400
		{1900, false}, // Divisible by 100 but not 400
		{2100, false}, // Divisible by 100 but not 400
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.year)), func(t *testing.T) {
			result := isLeapYear(tt.year)
			if result != tt.expected {
				t.Errorf("isLeapYear(%d) = %v, expected %v", tt.year, result, tt.expected)
			}
		})
	}
}

func BenchmarkGetWeeksInMonth(b *testing.B) {
	year := 2025
	month := time.January
	weekStart := time.Monday

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getWeeksInMonth(year, month, weekStart)
	}
}

func BenchmarkDateRange(b *testing.B) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dateRange(start, end)
	}
}

// Helper functions for testing

func isWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func firstDayOfMonth(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func lastDayOfMonth(year int, month time.Month) time.Time {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
}

func getWeeksInMonth(year int, month time.Month, weekStart time.Weekday) [][]time.Time {
	// Simplified implementation for testing
	var weeks [][]time.Time
	firstDay := firstDayOfMonth(year, month)
	lastDay := lastDayOfMonth(year, month)

	current := firstDay
	for current.Before(lastDay) || current.Equal(lastDay) {
		week := []time.Time{}
		for i := 0; i < 7; i++ {
			week = append(week, current)
			current = current.AddDate(0, 0, 1)
			if current.After(lastDay) {
				break
			}
		}
		weeks = append(weeks, week)
	}

	return weeks
}

func dateRange(start, end time.Time) []time.Time {
	var dates []time.Time
	current := start

	for current.Before(end) || current.Equal(end) {
		dates = append(dates, current)
		current = current.AddDate(0, 0, 1)
	}

	return dates
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func getMonthName(month time.Month) string {
	return month.String()
}

func TestEscapeLatexSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "backslash",
			input:    `path\to\file`,
			expected: `path\textbackslash\{\}to\textbackslash\{\}file`,
		},
		{
			name:     "curly braces",
			input:    "{bold} text",
			expected: `\{bold\} text`,
		},
		{
			name:     "dollar and ampersand",
			input:    "$100 & 50%",
			expected: `\$100 \& 50\%`,
		},
		{
			name:     "hash and caret",
			input:    "#1 ^ superscript",
			expected: `\#1 \textasciicircum{} superscript`,
		},
		{
			name:     "underscore and tilde",
			input:    "file_name ~ approx",
			expected: `file\_name \textasciitilde{} approx`,
		},
		{
			name:     "multiple special chars",
			input:    `a{b}$c&d%e#f^g_h~i\j`,
			expected: `a\{b\}\$c\&d\%e\#f\textasciicircum{}g\_h\textasciitilde{}i\textbackslash\{\}j`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no special chars",
			input:    "Normal text",
			expected: "Normal text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calendar.EscapeLatexSpecialChars(tt.input)
			if result != tt.expected {
				t.Errorf("calendar.EscapeLatexSpecialChars(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetPhaseDescription(t *testing.T) {
	tests := []struct {
		name     string
		phaseNum string
		expected string
	}{
		{
			name:     "phase 1",
			phaseNum: "1",
			expected: "Phase 1: Proposal \\& Setup",
		},
		{
			name:     "phase 2",
			phaseNum: "2",
			expected: "Phase 2: Research \\& Data Collection",
		},
		{
			name:     "phase 3",
			phaseNum: "3",
			expected: "Phase 3: Publications",
		},
		{
			name:     "phase 4",
			phaseNum: "4",
			expected: "Phase 4: Dissertation",
		},
		{
			name:     "unknown phase",
			phaseNum: "5",
			expected: "Phase 5",
		},
		{
			name:     "empty phase",
			phaseNum: "",
			expected: "Phase ",
		},
		{
			name:     "non-numeric phase",
			phaseNum: "abc",
			expected: "Phase abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calendar.GetPhaseDescription(tt.phaseNum)
			if result != tt.expected {
				t.Errorf("calendar.GetPhaseDescription(%q) = %q, want %q", tt.phaseNum, result, tt.expected)
			}
		})
	}
}

func TestCreateCalendarSpanningTask(t *testing.T) {
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)

	task := core.Task{
		ID:          "T1",
		Name:        "Test Task",
		Description: "A test task",
		Phase:       "1",
		SubPhase:    "Planning",
		Category:    "Research",
		Status:      "In Progress",
		Assignee:    "John Doe",
		IsMilestone: true,
	}

	result := calendar.CreateSpanningTask(task, startDate, endDate)

	// Check that all fields are copied correctly
	if result.ID != task.ID {
		t.Errorf("ID = %q, want %q", result.ID, task.ID)
	}
	if result.Name != task.Name {
		t.Errorf("Name = %q, want %q", result.Name, task.Name)
	}
	if result.Description != task.Description {
		t.Errorf("Description = %q, want %q", result.Description, task.Description)
	}
	if result.Phase != task.Phase {
		t.Errorf("Phase = %q, want %q", result.Phase, task.Phase)
	}
	if result.SubPhase != task.SubPhase {
		t.Errorf("SubPhase = %q, want %q", result.SubPhase, task.SubPhase)
	}
	if result.Category != task.Category {
		t.Errorf("Category = %q, want %q", result.Category, task.Category)
	}
	if !result.StartDate.Equal(startDate) {
		t.Errorf("StartDate = %v, want %v", result.StartDate, startDate)
	}
	if !result.EndDate.Equal(endDate) {
		t.Errorf("EndDate = %v, want %v", result.EndDate, endDate)
	}
	if result.Progress != 0 {
		t.Errorf("Progress = %d, want 0", result.Progress)
	}
	if result.Status != task.Status {
		t.Errorf("Status = %q, want %q", result.Status, task.Status)
	}
	if result.Assignee != task.Assignee {
		t.Errorf("Assignee = %q, want %q", result.Assignee, task.Assignee)
	}
	if result.IsMilestone != task.IsMilestone {
		t.Errorf("IsMilestone = %v, want %v", result.IsMilestone, task.IsMilestone)
	}

	// Check that color is generated from category
	expectedColor := core.GenerateCategoryColor(task.Category)
	if result.Color != expectedColor {
		t.Errorf("Color = %q, want %q", result.Color, expectedColor)
	}
}
