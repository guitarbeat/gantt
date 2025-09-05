package generator

import (
	"fmt"
	"strings"
)

// VisualDesignSystem provides comprehensive visual design management
type VisualDesignSystem struct {
	colorScheme    *ColorScheme
	typography     *TypographySystem
	visualTokens   *VisualTokens
	accessibility  *AccessibilityConfig
	logger         PDFLogger
}

// ColorScheme defines the color palette and schemes
type ColorScheme struct {
	// Base colors
	Primary       ColorDefinition `json:"primary"`
	Secondary     ColorDefinition `json:"secondary"`
	Accent        ColorDefinition `json:"accent"`
	Neutral       ColorDefinition `json:"neutral"`
	Background    ColorDefinition `json:"background"`
	Surface       ColorDefinition `json:"surface"`
	Text          ColorDefinition `json:"text"`
	
	// Category colors
	Categories    map[string]CategoryColorSet `json:"categories"`
	
	// Status colors
	Success       ColorDefinition `json:"success"`
	Warning       ColorDefinition `json:"warning"`
	Error         ColorDefinition `json:"error"`
	Info          ColorDefinition `json:"info"`
	
	// Semantic colors
	Border        ColorDefinition `json:"border"`
	Shadow        ColorDefinition `json:"shadow"`
	Highlight     ColorDefinition `json:"highlight"`
}

// ColorDefinition defines a color with multiple variants
type ColorDefinition struct {
	Base      string `json:"base"`      // Primary color
	Light     string `json:"light"`     // Light variant
	Dark      string `json:"dark"`      // Dark variant
	Lighter   string `json:"lighter"`   // Lighter variant
	Darker    string `json:"darker"`    // Darker variant
	Alpha     string `json:"alpha"`     // Alpha variant
	Contrast  string `json:"contrast"`  // High contrast variant
}

// CategoryColorSet defines colors for a specific category
type CategoryColorSet struct {
	Primary   ColorDefinition `json:"primary"`
	Secondary ColorDefinition `json:"secondary"`
	Accent    ColorDefinition `json:"accent"`
	Text      ColorDefinition `json:"text"`
	Border    ColorDefinition `json:"border"`
}

// TypographySystem defines the typography system
type TypographySystem struct {
	// Font families
	PrimaryFont   FontFamily `json:"primary_font"`
	SecondaryFont FontFamily `json:"secondary_font"`
	MonospaceFont FontFamily `json:"monospace_font"`
	
	// Font scales
	Scale         FontScale  `json:"scale"`
	
	// Text styles
	Styles        map[string]TextStyle `json:"styles"`
	
	// Line heights
	LineHeights   map[string]float64   `json:"line_heights"`
	
	// Letter spacing
	LetterSpacing map[string]float64   `json:"letter_spacing"`
}

// FontFamily defines a font family with variants
type FontFamily struct {
	Name     string   `json:"name"`
	Fallback []string `json:"fallback"`
	Weights  []int    `json:"weights"`
	Styles   []string `json:"styles"`
}

// FontScale defines font size scales
type FontScale struct {
	Base    float64            `json:"base"`
	Ratio   float64            `json:"ratio"`
	Sizes   map[string]float64 `json:"sizes"`
}

// TextStyle defines a text style
type TextStyle struct {
	FontFamily    string  `json:"font_family"`
	FontSize      float64 `json:"font_size"`
	FontWeight    int     `json:"font_weight"`
	FontStyle     string  `json:"font_style"`
	LineHeight    float64 `json:"line_height"`
	LetterSpacing float64 `json:"letter_spacing"`
	Color         string  `json:"color"`
	TextAlign     string  `json:"text_align"`
}

// VisualTokens defines visual design tokens
type VisualTokens struct {
	// Spacing tokens
	Spacing map[string]float64 `json:"spacing"`
	
	// Border radius tokens
	BorderRadius map[string]float64 `json:"border_radius"`
	
	// Shadow tokens
	Shadows map[string]ShadowDefinition `json:"shadows"`
	
	// Border tokens
	Borders map[string]BorderDefinition `json:"borders"`
	
	// Animation tokens
	Animations map[string]AnimationDefinition `json:"animations"`
}

// ShadowDefinition defines shadow properties
type ShadowDefinition struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Blur   float64 `json:"blur"`
	Spread float64 `json:"spread"`
	Color  string  `json:"color"`
	Opacity float64 `json:"opacity"`
}

// BorderDefinition defines border properties
type BorderDefinition struct {
	Width float64 `json:"width"`
	Style string  `json:"style"`
	Color string  `json:"color"`
}

// AnimationDefinition defines animation properties
type AnimationDefinition struct {
	Duration float64 `json:"duration"`
	Easing   string  `json:"easing"`
	Delay    float64 `json:"delay"`
}

// AccessibilityConfig defines accessibility settings
type AccessibilityConfig struct {
	// Contrast ratios
	MinContrastRatio float64 `json:"min_contrast_ratio"`
	
	// Color accessibility
	ColorBlindSafe   bool    `json:"color_blind_safe"`
	HighContrast     bool    `json:"high_contrast"`
	
	// Typography accessibility
	MinFontSize      float64 `json:"min_font_size"`
	MaxLineLength    int     `json:"max_line_length"`
	
	// Focus indicators
	FocusVisible     bool    `json:"focus_visible"`
	FocusColor       string  `json:"focus_color"`
}

// NewVisualDesignSystem creates a new visual design system
func NewVisualDesignSystem() *VisualDesignSystem {
	return &VisualDesignSystem{
		colorScheme:   GetDefaultColorScheme(),
		typography:    GetDefaultTypographySystem(),
		visualTokens:  GetDefaultVisualTokens(),
		accessibility: GetDefaultAccessibilityConfig(),
		logger:        &VisualDesignLogger{},
	}
}

// GetDefaultColorScheme returns the default color scheme
func GetDefaultColorScheme() *ColorScheme {
	return &ColorScheme{
		Primary: ColorDefinition{
			Base:     "#2563EB", // Indigo-600
			Light:    "#3B82F6", // Indigo-500
			Dark:     "#1D4ED8", // Indigo-700
			Lighter:  "#60A5FA", // Indigo-400
			Darker:   "#1E40AF", // Indigo-800
			Alpha:    "#2563EB20", // Indigo-600 with 20% opacity
			Contrast: "#FFFFFF", // White for contrast
		},
		Secondary: ColorDefinition{
			Base:     "#6B7280", // Gray-500
			Light:    "#9CA3AF", // Gray-400
			Dark:     "#4B5563", // Gray-600
			Lighter:  "#D1D5DB", // Gray-300
			Darker:   "#374151", // Gray-700
			Alpha:    "#6B728020", // Gray-500 with 20% opacity
			Contrast: "#FFFFFF", // White for contrast
		},
		Accent: ColorDefinition{
			Base:     "#F59E0B", // Amber-500
			Light:    "#FBBF24", // Amber-400
			Dark:     "#D97706", // Amber-600
			Lighter:  "#FCD34D", // Amber-300
			Darker:   "#B45309", // Amber-700
			Alpha:    "#F59E0B20", // Amber-500 with 20% opacity
			Contrast: "#000000", // Black for contrast
		},
		Neutral: ColorDefinition{
			Base:     "#9CA3AF", // Gray-400
			Light:    "#D1D5DB", // Gray-300
			Dark:     "#6B7280", // Gray-500
			Lighter:  "#F3F4F6", // Gray-100
			Darker:   "#4B5563", // Gray-600
			Alpha:    "#9CA3AF20", // Gray-400 with 20% opacity
			Contrast: "#000000", // Black for contrast
		},
		Background: ColorDefinition{
			Base:     "#FFFFFF", // White
			Light:    "#F9FAFB", // Gray-50
			Dark:     "#F3F4F6", // Gray-100
			Lighter:  "#FFFFFF", // White
			Darker:   "#E5E7EB", // Gray-200
			Alpha:    "#FFFFFF20", // White with 20% opacity
			Contrast: "#000000", // Black for contrast
		},
		Surface: ColorDefinition{
			Base:     "#F9FAFB", // Gray-50
			Light:    "#FFFFFF", // White
			Dark:     "#F3F4F6", // Gray-100
			Lighter:  "#FFFFFF", // White
			Darker:   "#E5E7EB", // Gray-200
			Alpha:    "#F9FAFB20", // Gray-50 with 20% opacity
			Contrast: "#000000", // Black for contrast
		},
		Text: ColorDefinition{
			Base:     "#111827", // Gray-900
			Light:    "#374151", // Gray-700
			Dark:     "#000000", // Black
			Lighter:  "#6B7280", // Gray-500
			Darker:   "#000000", // Black
			Alpha:    "#11182720", // Gray-900 with 20% opacity
			Contrast: "#FFFFFF", // White for contrast
		},
		Categories: map[string]CategoryColorSet{
			"PROPOSAL": {
				Primary:   ColorDefinition{Base: "#2563EB", Light: "#3B82F6", Dark: "#1D4ED8", Contrast: "#FFFFFF"},
				Secondary: ColorDefinition{Base: "#DBEAFE", Light: "#EFF6FF", Dark: "#BFDBFE", Contrast: "#1D4ED8"},
				Accent:    ColorDefinition{Base: "#1E40AF", Light: "#2563EB", Dark: "#1E3A8A", Contrast: "#FFFFFF"},
				Text:      ColorDefinition{Base: "#1E40AF", Light: "#2563EB", Dark: "#1E3A8A", Contrast: "#FFFFFF"},
				Border:    ColorDefinition{Base: "#3B82F6", Light: "#60A5FA", Dark: "#2563EB", Contrast: "#FFFFFF"},
			},
			"LASER": {
				Primary:   ColorDefinition{Base: "#EA580C", Light: "#F97316", Dark: "#C2410C", Contrast: "#FFFFFF"},
				Secondary: ColorDefinition{Base: "#FED7AA", Light: "#FFEDD5", Dark: "#FDBA74", Contrast: "#C2410C"},
				Accent:    ColorDefinition{Base: "#9A3412", Light: "#EA580C", Dark: "#7C2D12", Contrast: "#FFFFFF"},
				Text:      ColorDefinition{Base: "#9A3412", Light: "#EA580C", Dark: "#7C2D12", Contrast: "#FFFFFF"},
				Border:    ColorDefinition{Base: "#F97316", Light: "#FB923C", Dark: "#EA580C", Contrast: "#FFFFFF"},
			},
			"IMAGING": {
				Primary:   ColorDefinition{Base: "#16A34A", Light: "#22C55E", Dark: "#15803D", Contrast: "#FFFFFF"},
				Secondary: ColorDefinition{Base: "#BBF7D0", Light: "#DCFCE7", Dark: "#86EFAC", Contrast: "#15803D"},
				Accent:    ColorDefinition{Base: "#166534", Light: "#16A34A", Dark: "#14532D", Contrast: "#FFFFFF"},
				Text:      ColorDefinition{Base: "#166534", Light: "#16A34A", Dark: "#14532D", Contrast: "#FFFFFF"},
				Border:    ColorDefinition{Base: "#22C55E", Light: "#4ADE80", Dark: "#16A34A", Contrast: "#FFFFFF"},
			},
			"ADMIN": {
				Primary:   ColorDefinition{Base: "#6B7280", Light: "#9CA3AF", Dark: "#4B5563", Contrast: "#FFFFFF"},
				Secondary: ColorDefinition{Base: "#E5E7EB", Light: "#F3F4F6", Dark: "#D1D5DB", Contrast: "#4B5563"},
				Accent:    ColorDefinition{Base: "#374151", Light: "#6B7280", Dark: "#1F2937", Contrast: "#FFFFFF"},
				Text:      ColorDefinition{Base: "#374151", Light: "#6B7280", Dark: "#1F2937", Contrast: "#FFFFFF"},
				Border:    ColorDefinition{Base: "#9CA3AF", Light: "#D1D5DB", Dark: "#6B7280", Contrast: "#FFFFFF"},
			},
			"DISSERTATION": {
				Primary:   ColorDefinition{Base: "#7C3AED", Light: "#8B5CF6", Dark: "#6D28D9", Contrast: "#FFFFFF"},
				Secondary: ColorDefinition{Base: "#DDD6FE", Light: "#EDE9FE", Dark: "#C4B5FD", Contrast: "#6D28D9"},
				Accent:    ColorDefinition{Base: "#5B21B6", Light: "#7C3AED", Dark: "#4C1D95", Contrast: "#FFFFFF"},
				Text:      ColorDefinition{Base: "#5B21B6", Light: "#7C3AED", Dark: "#4C1D95", Contrast: "#FFFFFF"},
				Border:    ColorDefinition{Base: "#8B5CF6", Light: "#A78BFA", Dark: "#7C3AED", Contrast: "#FFFFFF"},
			},
			"RESEARCH": {
				Primary:   ColorDefinition{Base: "#0EA5E9", Light: "#38BDF8", Dark: "#0284C7", Contrast: "#FFFFFF"},
				Secondary: ColorDefinition{Base: "#BAE6FD", Light: "#E0F2FE", Dark: "#7DD3FC", Contrast: "#0284C7"},
				Accent:    ColorDefinition{Base: "#0C4A6E", Light: "#0EA5E9", Dark: "#075985", Contrast: "#FFFFFF"},
				Text:      ColorDefinition{Base: "#0C4A6E", Light: "#0EA5E9", Dark: "#075985", Contrast: "#FFFFFF"},
				Border:    ColorDefinition{Base: "#38BDF8", Light: "#7DD3FC", Dark: "#0EA5E9", Contrast: "#FFFFFF"},
			},
			"PUBLICATION": {
				Primary:   ColorDefinition{Base: "#DC2626", Light: "#EF4444", Dark: "#B91C1C", Contrast: "#FFFFFF"},
				Secondary: ColorDefinition{Base: "#FECACA", Light: "#FEE2E2", Dark: "#FCA5A5", Contrast: "#B91C1C"},
				Accent:    ColorDefinition{Base: "#991B1B", Light: "#DC2626", Dark: "#7F1D1D", Contrast: "#FFFFFF"},
				Text:      ColorDefinition{Base: "#991B1B", Light: "#DC2626", Dark: "#7F1D1D", Contrast: "#FFFFFF"},
				Border:    ColorDefinition{Base: "#EF4444", Light: "#F87171", Dark: "#DC2626", Contrast: "#FFFFFF"},
			},
		},
		Success: ColorDefinition{
			Base:     "#10B981", // Emerald-500
			Light:    "#34D399", // Emerald-400
			Dark:     "#059669", // Emerald-600
			Contrast: "#FFFFFF", // White for contrast
		},
		Warning: ColorDefinition{
			Base:     "#F59E0B", // Amber-500
			Light:    "#FBBF24", // Amber-400
			Dark:     "#D97706", // Amber-600
			Contrast: "#000000", // Black for contrast
		},
		Error: ColorDefinition{
			Base:     "#EF4444", // Red-500
			Light:    "#F87171", // Red-400
			Dark:     "#DC2626", // Red-600
			Contrast: "#FFFFFF", // White for contrast
		},
		Info: ColorDefinition{
			Base:     "#3B82F6", // Blue-500
			Light:    "#60A5FA", // Blue-400
			Dark:     "#2563EB", // Blue-600
			Contrast: "#FFFFFF", // White for contrast
		},
		Border: ColorDefinition{
			Base:     "#E5E7EB", // Gray-200
			Light:    "#F3F4F6", // Gray-100
			Dark:     "#D1D5DB", // Gray-300
			Contrast: "#000000", // Black for contrast
		},
		Shadow: ColorDefinition{
			Base:     "#000000", // Black
			Light:    "#00000010", // Black with 10% opacity
			Dark:     "#00000020", // Black with 20% opacity
			Contrast: "#FFFFFF", // White for contrast
		},
		Highlight: ColorDefinition{
			Base:     "#FEF3C7", // Amber-100
			Light:    "#FFFBEB", // Amber-50
			Dark:     "#FDE68A", // Amber-200
			Contrast: "#000000", // Black for contrast
		},
	}
}

// GetDefaultTypographySystem returns the default typography system
func GetDefaultTypographySystem() *TypographySystem {
	return &TypographySystem{
		PrimaryFont: FontFamily{
			Name:     "Inter",
			Fallback: []string{"system-ui", "-apple-system", "BlinkMacSystemFont", "Segoe UI", "Roboto", "sans-serif"},
			Weights:  []int{300, 400, 500, 600, 700},
			Styles:   []string{"normal", "italic"},
		},
		SecondaryFont: FontFamily{
			Name:     "Inter",
			Fallback: []string{"system-ui", "-apple-system", "BlinkMacSystemFont", "Segoe UI", "Roboto", "sans-serif"},
			Weights:  []int{400, 500, 600},
			Styles:   []string{"normal"},
		},
		MonospaceFont: FontFamily{
			Name:     "JetBrains Mono",
			Fallback: []string{"SF Mono", "Monaco", "Inconsolata", "Roboto Mono", "monospace"},
			Weights:  []int{400, 500, 600},
			Styles:   []string{"normal", "italic"},
		},
		Scale: FontScale{
			Base:  16.0, // 16px base
			Ratio: 1.25, // Major third ratio
			Sizes: map[string]float64{
				"xs":   12.0, // 12px
				"sm":   14.0, // 14px
				"base": 16.0, // 16px
				"lg":   18.0, // 18px
				"xl":   20.0, // 20px
				"2xl":  24.0, // 24px
				"3xl":  30.0, // 30px
				"4xl":  36.0, // 36px
				"5xl":  48.0, // 48px
				"6xl":  60.0, // 60px
			},
		},
		Styles: map[string]TextStyle{
			"heading-1": {
				FontFamily:    "Inter",
				FontSize:      36.0,
				FontWeight:    700,
				FontStyle:     "normal",
				LineHeight:    1.2,
				LetterSpacing: -0.025,
				Color:         "#111827",
				TextAlign:     "left",
			},
			"heading-2": {
				FontFamily:    "Inter",
				FontSize:      30.0,
				FontWeight:    600,
				FontStyle:     "normal",
				LineHeight:    1.3,
				LetterSpacing: -0.025,
				Color:         "#111827",
				TextAlign:     "left",
			},
			"heading-3": {
				FontFamily:    "Inter",
				FontSize:      24.0,
				FontWeight:    600,
				FontStyle:     "normal",
				LineHeight:    1.4,
				LetterSpacing: 0,
				Color:         "#111827",
				TextAlign:     "left",
			},
			"body-large": {
				FontFamily:    "Inter",
				FontSize:      18.0,
				FontWeight:    400,
				FontStyle:     "normal",
				LineHeight:    1.6,
				LetterSpacing: 0,
				Color:         "#374151",
				TextAlign:     "left",
			},
			"body": {
				FontFamily:    "Inter",
				FontSize:      16.0,
				FontWeight:    400,
				FontStyle:     "normal",
				LineHeight:    1.5,
				LetterSpacing: 0,
				Color:         "#374151",
				TextAlign:     "left",
			},
			"body-small": {
				FontFamily:    "Inter",
				FontSize:      14.0,
				FontWeight:    400,
				FontStyle:     "normal",
				LineHeight:    1.5,
				LetterSpacing: 0,
				Color:         "#6B7280",
				TextAlign:     "left",
			},
			"caption": {
				FontFamily:    "Inter",
				FontSize:      12.0,
				FontWeight:    400,
				FontStyle:     "normal",
				LineHeight:    1.4,
				LetterSpacing: 0.025,
				Color:         "#6B7280",
				TextAlign:     "left",
			},
			"task-title": {
				FontFamily:    "Inter",
				FontSize:      14.0,
				FontWeight:    500,
				FontStyle:     "normal",
				LineHeight:    1.4,
				LetterSpacing: 0,
				Color:         "#FFFFFF",
				TextAlign:     "left",
			},
			"task-description": {
				FontFamily:    "Inter",
				FontSize:      12.0,
				FontWeight:    400,
				FontStyle:     "normal",
				LineHeight:    1.3,
				LetterSpacing: 0,
				Color:         "#FFFFFF",
				TextAlign:     "left",
			},
			"overflow": {
				FontFamily:    "Inter",
				FontSize:      10.0,
				FontWeight:    500,
				FontStyle:     "normal",
				LineHeight:    1.2,
				LetterSpacing: 0.025,
				Color:         "#6B7280",
				TextAlign:     "center",
			},
		},
		LineHeights: map[string]float64{
			"tight":   1.25,
			"snug":    1.375,
			"normal":  1.5,
			"relaxed": 1.625,
			"loose":   2.0,
		},
		LetterSpacing: map[string]float64{
			"tighter": -0.05,
			"tight":   -0.025,
			"normal":  0,
			"wide":    0.025,
			"wider":   0.05,
		},
	}
}

// GetDefaultVisualTokens returns the default visual tokens
func GetDefaultVisualTokens() *VisualTokens {
	return &VisualTokens{
		Spacing: map[string]float64{
			"0":   0,
			"1":   4,
			"2":   8,
			"3":   12,
			"4":   16,
			"5":   20,
			"6":   24,
			"8":   32,
			"10":  40,
			"12":  48,
			"16":  64,
			"20":  80,
			"24":  96,
			"32":  128,
			"40":  160,
			"48":  192,
			"56":  224,
			"64":  256,
		},
		BorderRadius: map[string]float64{
			"none":   0,
			"sm":     2,
			"base":   4,
			"md":     6,
			"lg":     8,
			"xl":     12,
			"2xl":    16,
			"3xl":    24,
			"full":   9999,
		},
		Shadows: map[string]ShadowDefinition{
			"sm": {
				X:       0,
				Y:       1,
				Blur:    2,
				Spread:  0,
				Color:   "#000000",
				Opacity: 0.05,
			},
			"base": {
				X:       0,
				Y:       1,
				Blur:    3,
				Spread:  0,
				Color:   "#000000",
				Opacity: 0.1,
			},
			"md": {
				X:       0,
				Y:       4,
				Blur:    6,
				Spread:  -1,
				Color:   "#000000",
				Opacity: 0.1,
			},
			"lg": {
				X:       0,
				Y:       10,
				Blur:    15,
				Spread:  -3,
				Color:   "#000000",
				Opacity: 0.1,
			},
			"xl": {
				X:       0,
				Y:       20,
				Blur:    25,
				Spread:  -5,
				Color:   "#000000",
				Opacity: 0.1,
			},
		},
		Borders: map[string]BorderDefinition{
			"none":   {Width: 0, Style: "none", Color: "transparent"},
			"thin":   {Width: 1, Style: "solid", Color: "#E5E7EB"},
			"base":   {Width: 1, Style: "solid", Color: "#D1D5DB"},
			"thick":  {Width: 2, Style: "solid", Color: "#9CA3AF"},
			"thicker": {Width: 4, Style: "solid", Color: "#6B7280"},
		},
		Animations: map[string]AnimationDefinition{
			"fast":   {Duration: 0.15, Easing: "ease-out", Delay: 0},
			"base":   {Duration: 0.2, Easing: "ease-out", Delay: 0},
			"slow":   {Duration: 0.3, Easing: "ease-out", Delay: 0},
			"slower": {Duration: 0.5, Easing: "ease-out", Delay: 0},
		},
	}
}

// GetDefaultAccessibilityConfig returns the default accessibility configuration
func GetDefaultAccessibilityConfig() *AccessibilityConfig {
	return &AccessibilityConfig{
		MinContrastRatio: 4.5,
		ColorBlindSafe:   true,
		HighContrast:     false,
		MinFontSize:      12.0,
		MaxLineLength:    75,
		FocusVisible:     true,
		FocusColor:       "#3B82F6",
	}
}

// SetLogger sets the logger for the visual design system
func (vds *VisualDesignSystem) SetLogger(logger PDFLogger) {
	vds.logger = logger
}

// GetCategoryColor returns the color for a specific category
func (vds *VisualDesignSystem) GetCategoryColor(category string, variant string) string {
	if categorySet, exists := vds.colorScheme.Categories[category]; exists {
		switch variant {
		case "primary":
			return categorySet.Primary.Base
		case "secondary":
			return categorySet.Secondary.Base
		case "accent":
			return categorySet.Accent.Base
		case "text":
			return categorySet.Text.Base
		case "border":
			return categorySet.Border.Base
		case "light":
			return categorySet.Primary.Light
		case "dark":
			return categorySet.Primary.Dark
		default:
			return categorySet.Primary.Base
		}
	}
	return vds.colorScheme.Primary.Base
}

// GetTextStyle returns the text style for a specific style name
func (vds *VisualDesignSystem) GetTextStyle(styleName string) *TextStyle {
	if style, exists := vds.typography.Styles[styleName]; exists {
		return &style
	}
	return &vds.typography.Styles["body"]
}

// GetSpacing returns the spacing value for a specific token
func (vds *VisualDesignSystem) GetSpacing(token string) float64 {
	if spacing, exists := vds.visualTokens.Spacing[token]; exists {
		return spacing
	}
	return 0
}

// GetBorderRadius returns the border radius value for a specific token
func (vds *VisualDesignSystem) GetBorderRadius(token string) float64 {
	if radius, exists := vds.visualTokens.BorderRadius[token]; exists {
		return radius
	}
	return 0
}

// GetShadow returns the shadow definition for a specific token
func (vds *VisualDesignSystem) GetShadow(token string) *ShadowDefinition {
	if shadow, exists := vds.visualTokens.Shadows[token]; exists {
		return &shadow
	}
	return &vds.visualTokens.Shadows["base"]
}

// GetBorder returns the border definition for a specific token
func (vds *VisualDesignSystem) GetBorder(token string) *BorderDefinition {
	if border, exists := vds.visualTokens.Borders[token]; exists {
		return &border
	}
	return &vds.visualTokens.Borders["base"]
}

// CalculateContrastRatio calculates the contrast ratio between two colors
func (vds *VisualDesignSystem) CalculateContrastRatio(color1, color2 string) float64 {
	// Simplified contrast ratio calculation
	// In a real implementation, you would convert hex colors to RGB and calculate luminance
	return 4.5 // Placeholder - should be calculated based on actual color values
}

// ValidateAccessibility validates accessibility requirements
func (vds *VisualDesignSystem) ValidateAccessibility() []string {
	issues := make([]string, 0)
	
	// Check minimum font size
	if vds.accessibility.MinFontSize < 12.0 {
		issues = append(issues, "Minimum font size should be at least 12px for accessibility")
	}
	
	// Check contrast ratio
	if vds.accessibility.MinContrastRatio < 4.5 {
		issues = append(issues, "Minimum contrast ratio should be at least 4.5:1 for accessibility")
	}
	
	// Check line length
	if vds.accessibility.MaxLineLength > 75 {
		issues = append(issues, "Maximum line length should be 75 characters or less for readability")
	}
	
	return issues
}

// GenerateLaTeXColorCommands generates LaTeX color commands
func (vds *VisualDesignSystem) GenerateLaTeXColorCommands() string {
	var commands strings.Builder
	
	// Define base colors
	commands.WriteString("\\definecolor{primary}{HTML}{2563EB}\n")
	commands.WriteString("\\definecolor{secondary}{HTML}{6B7280}\n")
	commands.WriteString("\\definecolor{accent}{HTML}{F59E0B}\n")
	commands.WriteString("\\definecolor{neutral}{HTML}{9CA3AF}\n")
	commands.WriteString("\\definecolor{background}{HTML}{FFFFFF}\n")
	commands.WriteString("\\definecolor{surface}{HTML}{F9FAFB}\n")
	commands.WriteString("\\definecolor{text}{HTML}{111827}\n")
	commands.WriteString("\\definecolor{success}{HTML}{10B981}\n")
	commands.WriteString("\\definecolor{warning}{HTML}{F59E0B}\n")
	commands.WriteString("\\definecolor{error}{HTML}{EF4444}\n")
	commands.WriteString("\\definecolor{info}{HTML}{3B82F6}\n")
	commands.WriteString("\\definecolor{border}{HTML}{E5E7EB}\n")
	commands.WriteString("\\definecolor{shadow}{HTML}{000000}\n")
	commands.WriteString("\\definecolor{highlight}{HTML}{FEF3C7}\n")
	
	// Define category colors
	for category, colorSet := range vds.colorScheme.Categories {
		commands.WriteString(fmt.Sprintf("\\definecolor{cat%s}{HTML}{%s}\n", category, strings.TrimPrefix(colorSet.Primary.Base, "#")))
		commands.WriteString(fmt.Sprintf("\\definecolor{cat%sLight}{HTML}{%s}\n", category, strings.TrimPrefix(colorSet.Primary.Light, "#")))
		commands.WriteString(fmt.Sprintf("\\definecolor{cat%sDark}{HTML}{%s}\n", category, strings.TrimPrefix(colorSet.Primary.Dark, "#")))
	}
	
	return commands.String()
}

// GenerateLaTeXTypographyCommands generates LaTeX typography commands
func (vds *VisualDesignSystem) GenerateLaTeXTypographyCommands() string {
	var commands strings.Builder
	
	// Define font families
	commands.WriteString("\\usepackage{fontspec}\n")
	commands.WriteString("\\setmainfont{Inter}\n")
	commands.WriteString("\\setmonofont{JetBrains Mono}\n")
	
	// Define text styles
	for styleName, style := range vds.typography.Styles {
		commands.WriteString(fmt.Sprintf("\\newcommand{\\%s}[1]{\\fontsize{%.1f}{%.1f}\\selectfont\\textbf{%s}}\n", 
			styleName, style.FontSize, style.FontSize*style.LineHeight, "#1"))
	}
	
	return commands.String()
}

// GenerateLaTeXVisualCommands generates LaTeX visual commands
func (vds *VisualDesignSystem) GenerateLaTeXVisualCommands() string {
	var commands strings.Builder
	
	// Define spacing
	for token, value := range vds.visualTokens.Spacing {
		commands.WriteString(fmt.Sprintf("\\newlength{\\spacing%s}\n", strings.Title(token)))
		commands.WriteString(fmt.Sprintf("\\setlength{\\spacing%s}{%.1fpt}\n", strings.Title(token), value))
	}
	
	// Define border radius
	for token, value := range vds.visualTokens.BorderRadius {
		commands.WriteString(fmt.Sprintf("\\newlength{\\radius%s}\n", strings.Title(token)))
		commands.WriteString(fmt.Sprintf("\\setlength{\\radius%s}{%.1fpt}\n", strings.Title(token), value))
	}
	
	// Define shadows
	for token, shadow := range vds.visualTokens.Shadows {
		commands.WriteString(fmt.Sprintf("\\tikzset{shadow%s/.style={shadow={xshift=%.1fpt,yshift=%.1fpt,blur=%.1fpt,spread=%.1fpt,opacity=%.2f,color=%s}}}\n", 
			strings.Title(token), shadow.X, shadow.Y, shadow.Blur, shadow.Spread, shadow.Opacity, shadow.Color))
	}
	
	return commands.String()
}

// VisualDesignLogger provides logging for visual design system
type VisualDesignLogger struct{}

func (l *VisualDesignLogger) Info(msg string, args ...interface{})  { fmt.Printf("[DESIGN-INFO] "+msg+"\n", args...) }
func (l *VisualDesignLogger) Error(msg string, args ...interface{}) { fmt.Printf("[DESIGN-ERROR] "+msg+"\n", args...) }
func (l *VisualDesignLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[DESIGN-DEBUG] "+msg+"\n", args...) }
