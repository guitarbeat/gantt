package generator

import (
	"fmt"
)

// ViewType represents different calendar view types
type ViewType string

const (
	ViewTypeMonthly   ViewType = "monthly"
	ViewTypeWeekly    ViewType = "weekly"
	ViewTypeYearly    ViewType = "yearly"
	ViewTypeQuarterly ViewType = "quarterly"
	ViewTypeDaily     ViewType = "daily"
)

// OutputFormat represents different output formats
type OutputFormat string

const (
	OutputFormatPDF     OutputFormat = "pdf"
	OutputFormatLaTeX   OutputFormat = "latex"
	OutputFormatHTML    OutputFormat = "html"
	OutputFormatSVG     OutputFormat = "svg"
	OutputFormatPNG     OutputFormat = "png"
)

// ViewConfig defines configuration for different calendar views
type ViewConfig struct {
	Type                ViewType
	Title               string
	Description         string
	TemplateName        string
	PageSize            PageSize
	Orientation         Orientation
	ShowWeekNumbers     bool
	ShowTaskBars        bool
	ShowTaskDetails     bool
	ShowLayoutStats     bool
	ShowCategoryLegend  bool
	MaxTasksPerDay      int
	TaskBarHeight       float64
	TaskBarSpacing      float64
	CompactMode         bool
	ColorScheme          ColorScheme
	FontSize            FontSize
	LayoutDensity       LayoutDensity
	CustomCSS           string
	CustomLaTeX         string
}

// PageSize represents different page sizes
type PageSize string

const (
	PageSizeLetter    PageSize = "letter"
	PageSizeA4        PageSize = "a4"
	PageSizeA3        PageSize = "a3"
	PageSizeLegal     PageSize = "legal"
	PageSizeTabloid   PageSize = "tabloid"
	PageSizeCustom    PageSize = "custom"
)

// Orientation represents page orientation
type Orientation string

const (
	OrientationPortrait  Orientation = "portrait"
	OrientationLandscape Orientation = "landscape"
)

// ColorScheme represents different color schemes
type ColorScheme string

const (
	ColorSchemeDefault    ColorScheme = "default"
	ColorSchemeMinimal    ColorScheme = "minimal"
	ColorSchemeHighContrast ColorScheme = "high-contrast"
	ColorSchemeColorBlind  ColorScheme = "colorblind"
	ColorSchemeDark       ColorScheme = "dark"
	ColorSchemeCustom     ColorScheme = "custom"
)

// FontSize represents different font sizes
type FontSize string

const (
	FontSizeTiny    FontSize = "tiny"
	FontSizeSmall   FontSize = "small"
	FontSizeNormal  FontSize = "normal"
	FontSizeLarge   FontSize = "large"
	FontSizeHuge    FontSize = "huge"
)

// LayoutDensity represents different layout densities
type LayoutDensity string

const (
	LayoutDensityCompact    LayoutDensity = "compact"
	LayoutDensityNormal     LayoutDensity = "normal"
	LayoutDensitySpacious   LayoutDensity = "spacious"
	LayoutDensityMinimal    LayoutDensity = "minimal"
)

// ViewPreset represents predefined view configurations
type ViewPreset struct {
	Name        string
	Description string
	Config      ViewConfig
}

// GetDefaultViewConfig returns the default configuration for a view type
func GetDefaultViewConfig(viewType ViewType) ViewConfig {
	switch viewType {
	case ViewTypeMonthly:
		return ViewConfig{
			Type:               ViewTypeMonthly,
			Title:              "Monthly Calendar",
			Description:        "Traditional monthly calendar view with task bars",
			TemplateName:       "monthly_page.tpl",
			PageSize:           PageSizeLetter,
			Orientation:        OrientationPortrait,
			ShowWeekNumbers:    true,
			ShowTaskBars:       true,
			ShowTaskDetails:    true,
			ShowLayoutStats:    false,
			ShowCategoryLegend: true,
			MaxTasksPerDay:     5,
			TaskBarHeight:      12.0,
			TaskBarSpacing:     1.0,
			CompactMode:        false,
			ColorScheme:        ColorSchemeDefault,
			FontSize:           FontSizeNormal,
			LayoutDensity:      LayoutDensityNormal,
		}
	case ViewTypeWeekly:
		return ViewConfig{
			Type:               ViewTypeWeekly,
			Title:              "Weekly Calendar",
			Description:        "Detailed weekly view with enhanced task visualization",
			TemplateName:       "weekly_page.tpl",
			PageSize:           PageSizeLetter,
			Orientation:        OrientationLandscape,
			ShowWeekNumbers:    true,
			ShowTaskBars:       true,
			ShowTaskDetails:    true,
			ShowLayoutStats:    true,
			ShowCategoryLegend: true,
			MaxTasksPerDay:     10,
			TaskBarHeight:      15.0,
			TaskBarSpacing:     1.5,
			CompactMode:        false,
			ColorScheme:        ColorSchemeDefault,
			FontSize:           FontSizeNormal,
			LayoutDensity:      LayoutDensityNormal,
		}
	case ViewTypeYearly:
		return ViewConfig{
			Type:               ViewTypeYearly,
			Title:              "Yearly Overview",
			Description:        "High-level yearly calendar overview",
			TemplateName:       "yearly_page.tpl",
			PageSize:           PageSizeA3,
			Orientation:        OrientationLandscape,
			ShowWeekNumbers:    false,
			ShowTaskBars:       true,
			ShowTaskDetails:    false,
			ShowLayoutStats:    false,
			ShowCategoryLegend: true,
			MaxTasksPerDay:     3,
			TaskBarHeight:      8.0,
			TaskBarSpacing:     0.5,
			CompactMode:        true,
			ColorScheme:        ColorSchemeMinimal,
			FontSize:           FontSizeSmall,
			LayoutDensity:      LayoutDensityCompact,
		}
	case ViewTypeQuarterly:
		return ViewConfig{
			Type:               ViewTypeQuarterly,
			Title:              "Quarterly Calendar",
			Description:        "Three-month quarterly view",
			TemplateName:       "quarterly_page.tpl",
			PageSize:           PageSizeA3,
			Orientation:        OrientationLandscape,
			ShowWeekNumbers:    true,
			ShowTaskBars:       true,
			ShowTaskDetails:    true,
			ShowLayoutStats:    false,
			ShowCategoryLegend: true,
			MaxTasksPerDay:     4,
			TaskBarHeight:      10.0,
			TaskBarSpacing:     1.0,
			CompactMode:        false,
			ColorScheme:        ColorSchemeDefault,
			FontSize:           FontSizeSmall,
			LayoutDensity:      LayoutDensityNormal,
		}
	case ViewTypeDaily:
		return ViewConfig{
			Type:               ViewTypeDaily,
			Title:              "Daily Planner",
			Description:        "Detailed daily view with hourly breakdown",
			TemplateName:       "daily_page.tpl",
			PageSize:           PageSizeLetter,
			Orientation:        OrientationPortrait,
			ShowWeekNumbers:    false,
			ShowTaskBars:       true,
			ShowTaskDetails:    true,
			ShowLayoutStats:    true,
			ShowCategoryLegend: true,
			MaxTasksPerDay:     20,
			TaskBarHeight:      18.0,
			TaskBarSpacing:     2.0,
			CompactMode:        false,
			ColorScheme:        ColorSchemeDefault,
			FontSize:           FontSizeNormal,
			LayoutDensity:      LayoutDensitySpacious,
		}
	default:
		return GetDefaultViewConfig(ViewTypeMonthly)
	}
}

// GetViewPresets returns a list of predefined view presets
func GetViewPresets() []ViewPreset {
	return []ViewPreset{
		{
			Name:        "monthly-standard",
			Description: "Standard monthly calendar with task bars",
			Config:      GetDefaultViewConfig(ViewTypeMonthly),
		},
		{
			Name:        "monthly-compact",
			Description: "Compact monthly view for printing",
			Config: func() ViewConfig {
				config := GetDefaultViewConfig(ViewTypeMonthly)
				config.CompactMode = true
				config.LayoutDensity = LayoutDensityCompact
				config.FontSize = FontSizeSmall
				config.MaxTasksPerDay = 3
				return config
			}(),
		},
		{
			Name:        "weekly-detailed",
			Description: "Detailed weekly view with enhanced task visualization",
			Config:      GetDefaultViewConfig(ViewTypeWeekly),
		},
		{
			Name:        "yearly-overview",
			Description: "High-level yearly overview",
			Config:      GetDefaultViewConfig(ViewTypeYearly),
		},
		{
			Name:        "quarterly-planning",
			Description: "Quarterly planning view",
			Config:      GetDefaultViewConfig(ViewTypeQuarterly),
		},
		{
			Name:        "daily-planner",
			Description: "Detailed daily planner",
			Config:      GetDefaultViewConfig(ViewTypeDaily),
		},
		{
			Name:        "minimal-print",
			Description: "Minimal design optimized for printing",
			Config: func() ViewConfig {
				config := GetDefaultViewConfig(ViewTypeMonthly)
				config.ColorScheme = ColorSchemeMinimal
				config.ShowCategoryLegend = false
				config.ShowLayoutStats = false
				config.LayoutDensity = LayoutDensityMinimal
				return config
			}(),
		},
		{
			Name:        "high-contrast",
			Description: "High contrast design for accessibility",
			Config: func() ViewConfig {
				config := GetDefaultViewConfig(ViewTypeMonthly)
				config.ColorScheme = ColorSchemeHighContrast
				config.FontSize = FontSizeLarge
				return config
			}(),
		},
	}
}

// GetPresetByName returns a view preset by name
func GetPresetByName(name string) (*ViewPreset, error) {
	presets := GetViewPresets()
	for _, preset := range presets {
		if preset.Name == name {
			return &preset, nil
		}
	}
	return nil, fmt.Errorf("preset not found: %s", name)
}

// ValidateViewConfig validates a view configuration
func ValidateViewConfig(config ViewConfig) error {
	// Validate view type
	validTypes := map[ViewType]bool{
		ViewTypeMonthly:   true,
		ViewTypeWeekly:    true,
		ViewTypeYearly:    true,
		ViewTypeQuarterly: true,
		ViewTypeDaily:     true,
	}
	if !validTypes[config.Type] {
		return fmt.Errorf("invalid view type: %s", config.Type)
	}

	// Validate page size
	validSizes := map[PageSize]bool{
		PageSizeLetter:  true,
		PageSizeA4:      true,
		PageSizeA3:      true,
		PageSizeLegal:   true,
		PageSizeTabloid: true,
		PageSizeCustom:  true,
	}
	if !validSizes[config.PageSize] {
		return fmt.Errorf("invalid page size: %s", config.PageSize)
	}

	// Validate orientation
	validOrientations := map[Orientation]bool{
		OrientationPortrait:  true,
		OrientationLandscape: true,
	}
	if !validOrientations[config.Orientation] {
		return fmt.Errorf("invalid orientation: %s", config.Orientation)
	}

	// Validate color scheme
	validColorSchemes := map[ColorScheme]bool{
		ColorSchemeDefault:      true,
		ColorSchemeMinimal:      true,
		ColorSchemeHighContrast: true,
		ColorSchemeColorBlind:   true,
		ColorSchemeDark:         true,
		ColorSchemeCustom:       true,
	}
	if !validColorSchemes[config.ColorScheme] {
		return fmt.Errorf("invalid color scheme: %s", config.ColorScheme)
	}

	// Validate font size
	validFontSizes := map[FontSize]bool{
		FontSizeTiny:   true,
		FontSizeSmall:  true,
		FontSizeNormal: true,
		FontSizeLarge:  true,
		FontSizeHuge:   true,
	}
	if !validFontSizes[config.FontSize] {
		return fmt.Errorf("invalid font size: %s", config.FontSize)
	}

	// Validate layout density
	validDensities := map[LayoutDensity]bool{
		LayoutDensityCompact:  true,
		LayoutDensityNormal:   true,
		LayoutDensitySpacious: true,
		LayoutDensityMinimal:  true,
	}
	if !validDensities[config.LayoutDensity] {
		return fmt.Errorf("invalid layout density: %s", config.LayoutDensity)
	}

	// Validate numeric values
	if config.MaxTasksPerDay < 1 {
		return fmt.Errorf("max tasks per day must be at least 1")
	}
	if config.TaskBarHeight <= 0 {
		return fmt.Errorf("task bar height must be positive")
	}
	if config.TaskBarSpacing < 0 {
		return fmt.Errorf("task bar spacing must be non-negative")
	}

	return nil
}

// ApplyViewConfig applies view configuration to PDF generation options
func ApplyViewConfig(config ViewConfig, options *PDFGenerationOptions) {
	// Note: PDFGenerationOptions doesn't have TemplateName field
	// This would need to be added to the struct if template selection is needed

	// Add view-specific packages
	extraPackages := []string{"tikz", "tcolorbox", "xcolor"}
	if config.ShowTaskBars {
		extraPackages = append(extraPackages, "pgfplots")
	}
	if config.ShowLayoutStats {
		extraPackages = append(extraPackages, "booktabs")
	}

	// Merge with existing packages
	existingPackages := make(map[string]bool)
	for _, pkg := range options.ExtraPackages {
		existingPackages[pkg] = true
	}
	for _, pkg := range extraPackages {
		if !existingPackages[pkg] {
			options.ExtraPackages = append(options.ExtraPackages, pkg)
		}
	}

	// Add view-specific LaTeX configuration
	viewLaTeX := generateViewLaTeX(config)
	if options.CustomPreamble == "" {
		options.CustomPreamble = viewLaTeX
	} else {
		options.CustomPreamble = viewLaTeX + "\n" + options.CustomPreamble
	}
}

// generateViewLaTeX generates LaTeX configuration for a view
func generateViewLaTeX(config ViewConfig) string {
	latex := fmt.Sprintf("%% View Configuration: %s\n", config.Title)
	latex += fmt.Sprintf("%% %s\n\n", config.Description)

	// Page size and orientation
	latex += fmt.Sprintf("\\usepackage[%s,%s]{geometry}\n", 
		config.PageSize, config.Orientation)

	// Font size
	fontSizeMap := map[FontSize]string{
		FontSizeTiny:   "8pt",
		FontSizeSmall:  "9pt",
		FontSizeNormal: "10pt",
		FontSizeLarge:  "11pt",
		FontSizeHuge:   "12pt",
	}
	if fontSize, ok := fontSizeMap[config.FontSize]; ok {
		latex += fmt.Sprintf("\\documentclass[%s]{article}\n", fontSize)
	}

	// Color scheme
	latex += generateColorSchemeLaTeX(config.ColorScheme)

	// Layout density
	latex += generateLayoutDensityLaTeX(config.LayoutDensity)

	// Task bar configuration
	if config.ShowTaskBars {
		latex += fmt.Sprintf("\\setlength{\\TaskBarHeight}{%.1fpt}\n", config.TaskBarHeight)
		latex += fmt.Sprintf("\\setlength{\\TaskBarSpacing}{%.1fpt}\n", config.TaskBarSpacing)
		latex += fmt.Sprintf("\\setcounter{MaxTasksPerDay}{%d}\n", config.MaxTasksPerDay)
	}

	return latex
}

// generateColorSchemeLaTeX generates LaTeX for color schemes
func generateColorSchemeLaTeX(scheme ColorScheme) string {
	switch scheme {
	case ColorSchemeMinimal:
		return `
% Minimal color scheme
\\definecolor{taskbg}{RGB}{245,245,245}
\\definecolor{taskborder}{RGB}{200,200,200}
\\definecolor{tasktext}{RGB}{50,50,50}
`
	case ColorSchemeHighContrast:
		return `
% High contrast color scheme
\\definecolor{taskbg}{RGB}{255,255,255}
\\definecolor{taskborder}{RGB}{0,0,0}
\\definecolor{tasktext}{RGB}{0,0,0}
`
	case ColorSchemeColorBlind:
		return `
% Colorblind-friendly color scheme
\\definecolor{proposal}{RGB}{0,114,178}
\\definecolor{laser}{RGB}{213,94,0}
\\definecolor{imaging}{RGB}{0,158,115}
\\definecolor{admin}{RGB}{128,128,128}
\\definecolor{dissertation}{RGB}{204,121,167}
\\definecolor{research}{RGB}{230,159,0}
\\definecolor{publication}{RGB}{86,180,233}
`
	case ColorSchemeDark:
		return `
% Dark color scheme
\\definecolor{taskbg}{RGB}{40,40,40}
\\definecolor{taskborder}{RGB}{100,100,100}
\\definecolor{tasktext}{RGB}{220,220,220}
\\pagecolor{black}
\\color{white}
`
	default:
		return `
% Default color scheme (already defined in macros.tpl)
`
	}
}

// generateLayoutDensityLaTeX generates LaTeX for layout density
func generateLayoutDensityLaTeX(density LayoutDensity) string {
	switch density {
	case LayoutDensityCompact:
		return `
% Compact layout
\\setlength{\\parskip}{0pt}
\\setlength{\\parsep}{0pt}
\\setlength{\\itemsep}{0pt}
\\setlength{\\topsep}{0pt}
`
	case LayoutDensitySpacious:
		return `
% Spacious layout
\\setlength{\\parskip}{6pt}
\\setlength{\\parsep}{6pt}
\\setlength{\\itemsep}{6pt}
\\setlength{\\topsep}{6pt}
`
	case LayoutDensityMinimal:
		return `
% Minimal layout
\\setlength{\\parskip}{0pt}
\\setlength{\\parsep}{0pt}
\\setlength{\\itemsep}{0pt}
\\setlength{\\topsep}{0pt}
\\setlength{\\abovedisplayskip}{0pt}
\\setlength{\\belowdisplayskip}{0pt}
`
	default:
		return `
% Normal layout (default)
`
	}
}
