// Package core - Defaults provides centralized default values for all configuration.
//
// This module serves as the single source of truth for all default configuration
// values used throughout the application. When no configuration is provided or
// when values are missing, these defaults are used.
//
// Key functions:
//   - DefaultConfig(): Complete configuration with all defaults
//   - DefaultLayout(): Layout configuration defaults
//   - DefaultLaTeX(): LaTeX-specific defaults
//   - DefaultTypography(): Typography settings defaults
//   - DefaultLayoutCalendarLayout(): Calendar layout defaults
//
// The Defaults struct provides easy access to commonly used default values
// as constants, useful for fallback logic throughout the codebase.
//
// Example usage:
//
//	// Get a complete default configuration
//	cfg := core.DefaultConfig()
//
//	// Override specific values
//	cfg.Year = 2025
//	cfg.OutputDir = "custom/output"
//
//	// Access default constants directly
//	width := core.Defaults.DayNumberWidth  // "6mm"
//	color := core.Defaults.DefaultTaskColor // "224,50,212"
//
// Design philosophy:
//   - Sensible defaults that work for most use cases
//   - Easy to understand and modify in one place
//   - No magic numbers scattered throughout code
//   - New users can see all options with defaults
package core

import "time"

// DefaultConfig returns a Config with sensible defaults
// This serves as the baseline configuration before loading from YAML or environment
func DefaultConfig() Config {
	return Config{
		Debug: Debug{
			ShowFrame: false,
			ShowLinks: false,
		},
		Year:                time.Now().Year(),
		WeekStart:           time.Sunday,
		Dotted:              false,
		CalAfterSchedule:    false,
		ClearTopRightCorner: false,
		AMPMTime:            false,
		AddLastHalfHour:     false,
		OutputDir:           "generated",
		Layout:              DefaultLayout(),
	}
}

// DefaultLayout returns a Layout with sensible defaults
func DefaultLayout() Layout {
	return Layout{
		Paper: Paper{
			Width:  "8.5in",
			Height: "11in",
			Margin: Margin{
				Top:    "1in",
				Bottom: "1in",
				Left:   "1in",
				Right:  "1in",
			},
			MarginParWidth: "1in",
			MarginParSep:   "0.2in",
			ReverseMargins: false,
		},
		LaTeX:        DefaultLaTeX(),
		LayoutEngine: DefaultLayoutEngine(),
	}
}

// DefaultLaTeX returns LaTeX configuration with sensible defaults
func DefaultLaTeX() LaTeX {
	return LaTeX{
		TabColSep:             "1pt",
		HeaderSideMonthsWidth: "3em",
		ArrayStretch:          1.0,
		MonthlyCellHeight:     "4em",
		HeaderResizeBox:       "0.9",
		LineThicknessDefault:  "0.4pt",
		LineThicknessThick:    "1.2pt",
		ColSep:                "3pt",
		Document:              DefaultDocument(),
		Typography:            DefaultTypography(),
	}
}

// DefaultDocument returns Document configuration defaults
func DefaultDocument() Document {
	return Document{
		FontSize:  "10pt",
		ParIndent: "0pt",
	}
}

// DefaultTypography returns Typography configuration defaults
func DefaultTypography() Typography {
	return Typography{
		HyphenPenalty:          50,
		Tolerance:              1000,
		EmergencyStretch:       "3em",
		SloppyEmergencyStretch: "3em",
	}
}

// DefaultLayoutEngine returns LayoutEngine configuration defaults
func DefaultLayoutEngine() LayoutEngine {
	return LayoutEngine{
		CalendarLayout: DefaultLayoutCalendarLayout(),
	}
}

// DefaultLayoutCalendarLayout returns LayoutCalendarLayout configuration defaults
func DefaultLayoutCalendarLayout() LayoutCalendarLayout {
	return LayoutCalendarLayout{
		DayNumberWidth:        "6mm",
		DayContentMargin:      "8mm",
		TaskCellMargin:        "1mm",
		TaskCellSpacing:       "0.5mm",
		DayCellMinipageWidth:  "8mm",
		HeaderAngleSizeOffset: "2pt",
	}
}

// DefaultTaskStyling returns TaskStyling configuration defaults
func DefaultTaskStyling() TaskStyling {
	return TaskStyling{
		FontSize:          "\\footnotesize",
		BarHeight:         "1.5mm",
		BorderWidth:       "0.5pt",
		ShowObjectives:    true,
		BackgroundOpacity: 20,
		BorderOpacity:     80,
		Spacing:           DefaultTaskStylingSpacing(),
		TColorBox:         DefaultTaskStylingTColorBox(),
	}
}

// DefaultTaskStylingSpacing returns TaskStylingSpacing defaults
func DefaultTaskStylingSpacing() TaskStylingSpacing {
	return TaskStylingSpacing{
		VerticalOffset:    "0.2ex",
		ContentVspace:     "0.2ex",
		PaddingHorizontal: "2pt",
		PaddingVertical:   "1pt",
	}
}

// DefaultTaskStylingTColorBox returns TaskStylingTColorBox defaults
func DefaultTaskStylingTColorBox() TaskStylingTColorBox {
	return TaskStylingTColorBox{
		Overlay: TColorBoxOverlay{
			Arc:     "2pt",
			Left:    "1pt",
			Right:   "1pt",
			Top:     "1pt",
			Bottom:  "1pt",
			BoxRule: "0.5pt",
		},
	}
}

// DefaultColors returns Colors configuration defaults
func DefaultColors() Colors {
	return Colors{
		Gray:      "0.5",
		LightGray: "0.8",
	}
}

// ConfigDefaults holds all default values as constants for easy reference
type ConfigDefaults struct {
	// Calendar layout defaults
	DayNumberWidth        string
	DayContentMargin      string
	TaskCellMargin        string
	TaskCellSpacing       string
	HeaderAngleSizeOffset string

	// Typography defaults
	HyphenPenalty    int
	Tolerance        int
	EmergencyStretch string

	// Output defaults
	DefaultOutputDir string

	// Task color defaults
	DefaultTaskColor string
}

// Defaults provides easy access to default values
var Defaults = ConfigDefaults{
	// Calendar layout
	DayNumberWidth:        "6mm",
	DayContentMargin:      "8mm",
	TaskCellMargin:        "1mm",
	TaskCellSpacing:       "0.5mm",
	HeaderAngleSizeOffset: "2pt",

	// Typography
	HyphenPenalty:    50,
	Tolerance:        1000,
	EmergencyStretch: "3em",

	// Output
	DefaultOutputDir: "generated",

	// Task colors
	DefaultTaskColor: "224,50,212", // Magenta fallback
}
