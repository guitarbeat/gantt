package core_test

import (
	"testing"

	"phd-dissertation-planner/internal/core"
)

func TestDefaultcore.Config(t *testing.T) {
	cfg := core.Defaultcore.Config()

	// Test that defaults are set
	if cfg.OutputDir == "" {
		t.Error("core.Defaultcore.Config() OutputDir should not be empty")
	}

	if cfg.OutputDir != "generated" {
		t.Errorf("core.Defaultcore.Config() OutputDir = %s, want 'generated'", cfg.OutputDir)
	}

	// Test that Layout is initialized
	if cfg.Layout.LaTeX.Typography.HyphenPenalty == 0 {
		t.Error("core.Defaultcore.Config() should initialize Typography.HyphenPenalty")
	}
}

func TestConfigHelperMethods(t *testing.T) {
	cfg := core.Defaultcore.Config()

	t.Run("GetDayNumberWidth", func(t *testing.T) {
		width := cfg.GetDayNumberWidth()
		if width == "" {
			t.Error("GetDayNumberWidth() should return a non-empty value")
		}
		if width != Defaults.DayNumberWidth {
			t.Errorf("GetDayNumberWidth() = %s, want %s", width, Defaults.DayNumberWidth)
		}
	})

	t.Run("GetDayContentMargin", func(t *testing.T) {
		margin := cfg.GetDayContentMargin()
		if margin == "" {
			t.Error("GetDayContentMargin() should return a non-empty value")
		}
		if margin != Defaults.DayContentMargin {
			t.Errorf("GetDayContentMargin() = %s, want %s", margin, Defaults.DayContentMargin)
		}
	})

	t.Run("GetHyphenPenalty", func(t *testing.T) {
		penalty := cfg.GetHyphenPenalty()
		if penalty <= 0 {
			t.Error("GetHyphenPenalty() should return a positive value")
		}
		if penalty != Defaults.HyphenPenalty {
			t.Errorf("GetHyphenPenalty() = %d, want %d", penalty, Defaults.HyphenPenalty)
		}
	})

	t.Run("GetTolerance", func(t *testing.T) {
		tolerance := cfg.GetTolerance()
		if tolerance <= 0 {
			t.Error("GetTolerance() should return a positive value")
		}
		if tolerance != Defaults.Tolerance {
			t.Errorf("GetTolerance() = %d, want %d", tolerance, Defaults.Tolerance)
		}
	})

	t.Run("GetOutputDir", func(t *testing.T) {
		dir := cfg.GetOutputDir()
		if dir == "" {
			t.Error("GetOutputDir() should return a non-empty value")
		}
	})

	t.Run("IsDebugMode", func(t *testing.T) {
		// Default config should not be in debug mode
		if cfg.IsDebugMode() {
			t.Error("IsDebugMode() should return false for default config")
		}

		// Enable debug
		cfg.Debug.ShowFrame = true
		if !cfg.IsDebugMode() {
			t.Error("IsDebugMode() should return true when ShowFrame is enabled")
		}
	})

	t.Run("HasCSVData", func(t *testing.T) {
		// Default config has no CSV
		if cfg.HasCSVData() {
			t.Error("HasCSVData() should return false for default config")
		}

		// Set CSV path
		cfg.CSVFilePath = "test.csv"
		if !cfg.HasCSVData() {
			t.Error("HasCSVData() should return true when CSVFilePath is set")
		}
	})
}

func TestConfigGetYear(t *testing.T) {
	cfg := Config{}

	// Without year set, should return current year
	year := cfg.GetYear()
	if year == 0 {
		t.Error("GetYear() should never return 0")
	}

	// With year set
	cfg.Year = 2025
	if cfg.GetYear() != 2025 {
		t.Errorf("GetYear() = %d, want 2025", cfg.GetYear())
	}
}

func TestFilterUniquecore.Modules(t *testing.T) {
	modules := Modules{
		{Tpl: "template1", Body: "body1"},
		{Tpl: "template2", Body: "body2"},
		{Tpl: "template1", Body: "body3"}, // Duplicate template
	}

	filtered := FilterUniquecore.Modules(modules)

	if len(filtered) != 2 {
		t.Errorf("FilterUniquecore.Modules() returned %d modules, want 2", len(filtered))
	}

	// Check that we kept the first occurrence
	if filtered[0].Tpl != "template1" {
		t.Errorf("FilterUniquecore.Modules() first module = %s, want 'template1'", filtered[0].Tpl)
	}
	if filtered[1].Tpl != "template2" {
		t.Errorf("FilterUniquecore.Modules() second module = %s, want 'template2'", filtered[1].Tpl)
	}
}
