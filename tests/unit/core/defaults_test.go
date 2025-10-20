package core_test

import (
	"testing"

	"phd-dissertation-planner/internal/core"
)

func Testcore.Defaults(t *testing.T) {
	t.Run("DayNumberWidth", func(t *testing.T) {
		if Defaults.DayNumberWidth == "" {
			t.Error("Defaults.DayNumberWidth should not be empty")
		}
	})

	t.Run("DayContentMargin", func(t *testing.T) {
		if Defaults.DayContentMargin == "" {
			t.Error("Defaults.DayContentMargin should not be empty")
		}
	})

	t.Run("HyphenPenalty", func(t *testing.T) {
		if Defaults.HyphenPenalty <= 0 {
			t.Error("Defaults.HyphenPenalty should be positive")
		}
	})

	t.Run("Tolerance", func(t *testing.T) {
		if Defaults.Tolerance <= 0 {
			t.Error("Defaults.Tolerance should be positive")
		}
	})

	t.Run("DefaultOutputDir", func(t *testing.T) {
		if Defaults.DefaultOutputDir == "" {
			t.Error("Defaults.DefaultOutputDir should not be empty")
		}
	})

	t.Run("DefaultTaskColor", func(t *testing.T) {
		if Defaults.DefaultTaskColor == "" {
			t.Error("Defaults.DefaultTaskColor should not be empty")
		}
	})
}

func TestDefaultTypography(t *testing.T) {
	typo := DefaultTypography()

	if typo.HyphenPenalty <= 0 {
		t.Error("DefaultTypography() HyphenPenalty should be positive")
	}

	if typo.Tolerance <= 0 {
		t.Error("DefaultTypography() Tolerance should be positive")
	}

	if typo.EmergencyStretch == "" {
		t.Error("DefaultTypography() EmergencyStretch should not be empty")
	}

	if typo.SloppyEmergencyStretch == "" {
		t.Error("DefaultTypography() SloppyEmergencyStretch should not be empty")
	}
}

func TestDefaultLayoutCalendarLayout(t *testing.T) {
	layout := DefaultLayoutCalendarLayout()

	fields := []struct {
		name  string
		value string
	}{
		{"DayNumberWidth", layout.DayNumberWidth},
		{"DayContentMargin", layout.DayContentMargin},
		{"TaskCellMargin", layout.TaskCellMargin},
		{"TaskCellSpacing", layout.TaskCellSpacing},
		{"DayCellMinipageWidth", layout.DayCellMinipageWidth},
		{"HeaderAngleSizeOffset", layout.HeaderAngleSizeOffset},
	}

	for _, f := range fields {
		t.Run(f.name, func(t *testing.T) {
			if f.value == "" {
				t.Errorf("DefaultLayoutCalendarLayout() %s should not be empty", f.name)
			}
		})
	}
}

func TestDefaultDocument(t *testing.T) {
	doc := DefaultDocument()

	if doc.FontSize == "" {
		t.Error("DefaultDocument() FontSize should not be empty")
	}

	if doc.ParIndent == "" {
		t.Error("DefaultDocument() ParIndent should not be empty")
	}
}

func TestDefaultLaTeX(t *testing.T) {
	latex := DefaultLaTeX()

	// Check that numeric values are set
	if latex.ArrayStretch == 0 {
		t.Error("DefaultLaTeX() ArrayStretch should not be zero")
	}

	// Check string values
	stringFields := []struct {
		name  string
		value string
	}{
		{"TabColSep", latex.TabColSep},
		{"HeaderSideMonthsWidth", latex.HeaderSideMonthsWidth},
		{"MonthlyCellHeight", latex.MonthlyCellHeight},
		{"HeaderResizeBox", latex.HeaderResizeBox},
		{"LineThicknessDefault", latex.LineThicknessDefault},
		{"LineThicknessThick", latex.LineThicknessThick},
		{"ColSep", latex.ColSep},
	}

	for _, f := range stringFields {
		t.Run(f.name, func(t *testing.T) {
			if f.value == "" {
				t.Errorf("DefaultLaTeX() %s should not be empty", f.name)
			}
		})
	}

	// Check nested structures are initialized
	if latex.Document.FontSize == "" {
		t.Error("DefaultLaTeX() Document should be initialized")
	}

	if latex.Typography.HyphenPenalty == 0 {
		t.Error("DefaultLaTeX() Typography should be initialized")
	}
}

func TestDefaultLayout(t *testing.T) {
	layout := DefaultLayout()

	// Check that LaTeX is initialized
	if layout.LaTeX.TabColSep == "" {
		t.Error("DefaultLayout() LaTeX should be initialized")
	}

	// Check that LayoutEngine is initialized
	if layout.LayoutEngine.CalendarLayout.DayNumberWidth == "" {
		t.Error("DefaultLayout() LayoutEngine should be initialized")
	}
}

func TestDefaultConfigCompleteness(t *testing.T) {
	cfg := core.Defaultcore.Config()

	// Test all major sections are initialized
	tests := []struct {
		name  string
		check func() bool
	}{
		{"Debug initialized", func() bool { return cfg.Debug.ShowFrame == false }},
		{"OutputDir set", func() bool { return cfg.OutputDir != "" }},
		{"Layout initialized", func() bool { return cfg.Layout.LaTeX.TabColSep != "" }},
		{"Typography initialized", func() bool { return cfg.Layout.LaTeX.Typography.HyphenPenalty > 0 }},
		{"CalendarLayout initialized", func() bool { return cfg.Layout.LayoutEngine.CalendarLayout.DayNumberWidth != "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Errorf("core.Defaultcore.Config() %s check failed", tt.name)
			}
		})
	}
}
