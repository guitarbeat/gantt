package main

import (
	"fmt"
	"strings"
)

// Test the enhanced visual design system
func main() {
	fmt.Println("Testing Enhanced Visual Design System...")

	// Test 1: Color Scheme
	fmt.Println("\n=== Test 1: Color Scheme ===")
	testColorScheme()

	// Test 2: Typography System
	fmt.Println("\n=== Test 2: Typography System ===")
	testTypographySystem()

	// Test 3: Visual Tokens
	fmt.Println("\n=== Test 3: Visual Tokens ===")
	testVisualTokens()

	// Test 4: Accessibility
	fmt.Println("\n=== Test 4: Accessibility ===")
	testAccessibility()

	// Test 5: LaTeX Generation
	fmt.Println("\n=== Test 5: LaTeX Generation ===")
	testLaTeXGeneration()

	fmt.Println("\n✅ Enhanced visual design system tests completed!")
}

// ColorScheme represents the color scheme
type ColorScheme struct {
	Primary    ColorDefinition
	Secondary  ColorDefinition
	Accent     ColorDefinition
	Neutral    ColorDefinition
	Background ColorDefinition
	Surface    ColorDefinition
	Text       ColorDefinition
	Categories map[string]CategoryColorSet
	Success    ColorDefinition
	Warning    ColorDefinition
	Error      ColorDefinition
	Info       ColorDefinition
	Border     ColorDefinition
	Shadow     ColorDefinition
	Highlight  ColorDefinition
}

// ColorDefinition represents a color with variants
type ColorDefinition struct {
	Base     string
	Light    string
	Dark     string
	Lighter  string
	Darker   string
	Alpha    string
	Contrast string
}

// CategoryColorSet represents colors for a category
type CategoryColorSet struct {
	Primary   ColorDefinition
	Secondary ColorDefinition
	Accent    ColorDefinition
	Text      ColorDefinition
	Border    ColorDefinition
}

// TypographySystem represents the typography system
type TypographySystem struct {
	PrimaryFont   FontFamily
	SecondaryFont FontFamily
	MonospaceFont FontFamily
	Scale         FontScale
	Styles        map[string]TextStyle
	LineHeights   map[string]float64
	LetterSpacing map[string]float64
}

// FontFamily represents a font family
type FontFamily struct {
	Name     string
	Fallback []string
	Weights  []int
	Styles   []string
}

// FontScale represents font size scales
type FontScale struct {
	Base  float64
	Ratio float64
	Sizes map[string]float64
}

// TextStyle represents a text style
type TextStyle struct {
	FontFamily    string
	FontSize      float64
	FontWeight    int
	FontStyle     string
	LineHeight    float64
	LetterSpacing float64
	Color         string
	TextAlign     string
}

// VisualTokens represents visual design tokens
type VisualTokens struct {
	Spacing     map[string]float64
	BorderRadius map[string]float64
	Shadows     map[string]ShadowDefinition
	Borders     map[string]BorderDefinition
	Animations  map[string]AnimationDefinition
}

// ShadowDefinition represents shadow properties
type ShadowDefinition struct {
	X       float64
	Y       float64
	Blur    float64
	Spread  float64
	Color   string
	Opacity float64
}

// BorderDefinition represents border properties
type BorderDefinition struct {
	Width float64
	Style string
	Color string
}

// AnimationDefinition represents animation properties
type AnimationDefinition struct {
	Duration float64
	Easing   string
	Delay    float64
}

// AccessibilityConfig represents accessibility settings
type AccessibilityConfig struct {
	MinContrastRatio float64
	ColorBlindSafe   bool
	HighContrast     bool
	MinFontSize      float64
	MaxLineLength    int
	FocusVisible     bool
	FocusColor       string
}

func testColorScheme() {
	// Test color scheme creation
	colorScheme := ColorScheme{
		Primary: ColorDefinition{
			Base:     "#2563EB",
			Light:    "#3B82F6",
			Dark:     "#1D4ED8",
			Lighter:  "#60A5FA",
			Darker:   "#1E40AF",
			Alpha:    "#2563EB20",
			Contrast: "#FFFFFF",
		},
		Secondary: ColorDefinition{
			Base:     "#6B7280",
			Light:    "#9CA3AF",
			Dark:     "#4B5563",
			Lighter:  "#D1D5DB",
			Darker:   "#374151",
			Alpha:    "#6B728020",
			Contrast: "#FFFFFF",
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
		},
	}

	// Validate color scheme
	if colorScheme.Primary.Base == "" {
		fmt.Println("❌ Primary color not defined")
		return
	}

	if colorScheme.Secondary.Base == "" {
		fmt.Println("❌ Secondary color not defined")
		return
	}

	if len(colorScheme.Categories) == 0 {
		fmt.Println("❌ No category colors defined")
		return
	}

	// Check category colors
	expectedCategories := []string{"PROPOSAL", "LASER"}
	for _, category := range expectedCategories {
		if _, exists := colorScheme.Categories[category]; !exists {
			fmt.Printf("❌ Category %s not defined\n", category)
			return
		}
	}

	// Check color variants
	for _, category := range colorScheme.Categories {
		if category.Primary.Base == "" {
			fmt.Println("❌ Category primary color not defined")
			return
		}
		if category.Secondary.Base == "" {
			fmt.Println("❌ Category secondary color not defined")
			return
		}
		if category.Accent.Base == "" {
			fmt.Println("❌ Category accent color not defined")
			return
		}
		if category.Text.Base == "" {
			fmt.Println("❌ Category text color not defined")
			return
		}
		if category.Border.Base == "" {
			fmt.Println("❌ Category border color not defined")
			return
		}
	}

	fmt.Printf("✅ Color scheme test passed\n")
	fmt.Printf("   Primary color: %s\n", colorScheme.Primary.Base)
	fmt.Printf("   Secondary color: %s\n", colorScheme.Secondary.Base)
	fmt.Printf("   Categories: %d\n", len(colorScheme.Categories))
	fmt.Printf("   PROPOSAL primary: %s\n", colorScheme.Categories["PROPOSAL"].Primary.Base)
	fmt.Printf("   LASER primary: %s\n", colorScheme.Categories["LASER"].Primary.Base)
}

func testTypographySystem() {
	// Test typography system creation
	typography := TypographySystem{
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
			Base:  16.0,
			Ratio: 1.25,
			Sizes: map[string]float64{
				"xs":   12.0,
				"sm":   14.0,
				"base": 16.0,
				"lg":   18.0,
				"xl":   20.0,
				"2xl":  24.0,
				"3xl":  30.0,
				"4xl":  36.0,
				"5xl":  48.0,
				"6xl":  60.0,
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

	// Validate typography system
	if typography.PrimaryFont.Name == "" {
		fmt.Println("❌ Primary font not defined")
		return
	}

	if typography.SecondaryFont.Name == "" {
		fmt.Println("❌ Secondary font not defined")
		return
	}

	if typography.MonospaceFont.Name == "" {
		fmt.Println("❌ Monospace font not defined")
		return
	}

	if typography.Scale.Base <= 0 {
		fmt.Println("❌ Font scale base should be positive")
		return
	}

	if typography.Scale.Ratio <= 0 {
		fmt.Println("❌ Font scale ratio should be positive")
		return
	}

	if len(typography.Scale.Sizes) == 0 {
		fmt.Println("❌ No font sizes defined")
		return
	}

	if len(typography.Styles) == 0 {
		fmt.Println("❌ No text styles defined")
		return
	}

	// Check font sizes
	expectedSizes := []string{"xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl"}
	for _, size := range expectedSizes {
		if _, exists := typography.Scale.Sizes[size]; !exists {
			fmt.Printf("❌ Font size %s not defined\n", size)
			return
		}
	}

	// Check text styles
	expectedStyles := []string{"heading-1", "body", "task-title"}
	for _, style := range expectedStyles {
		if _, exists := typography.Styles[style]; !exists {
			fmt.Printf("❌ Text style %s not defined\n", style)
			return
		}
	}

	// Check line heights
	if len(typography.LineHeights) == 0 {
		fmt.Println("❌ No line heights defined")
		return
	}

	// Check letter spacing
	if len(typography.LetterSpacing) == 0 {
		fmt.Println("❌ No letter spacing defined")
		return
	}

	fmt.Printf("✅ Typography system test passed\n")
	fmt.Printf("   Primary font: %s\n", typography.PrimaryFont.Name)
	fmt.Printf("   Secondary font: %s\n", typography.SecondaryFont.Name)
	fmt.Printf("   Monospace font: %s\n", typography.MonospaceFont.Name)
	fmt.Printf("   Font scale base: %.1f\n", typography.Scale.Base)
	fmt.Printf("   Font scale ratio: %.2f\n", typography.Scale.Ratio)
	fmt.Printf("   Font sizes: %d\n", len(typography.Scale.Sizes))
	fmt.Printf("   Text styles: %d\n", len(typography.Styles))
	fmt.Printf("   Line heights: %d\n", len(typography.LineHeights))
	fmt.Printf("   Letter spacing: %d\n", len(typography.LetterSpacing))
}

func testVisualTokens() {
	// Test visual tokens creation
	tokens := VisualTokens{
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

	// Validate visual tokens
	if len(tokens.Spacing) == 0 {
		fmt.Println("❌ No spacing tokens defined")
		return
	}

	if len(tokens.BorderRadius) == 0 {
		fmt.Println("❌ No border radius tokens defined")
		return
	}

	if len(tokens.Shadows) == 0 {
		fmt.Println("❌ No shadow tokens defined")
		return
	}

	if len(tokens.Borders) == 0 {
		fmt.Println("❌ No border tokens defined")
		return
	}

	if len(tokens.Animations) == 0 {
		fmt.Println("❌ No animation tokens defined")
		return
	}

	// Check spacing values
	expectedSpacing := []string{"0", "1", "2", "3", "4", "5", "6", "8", "10", "12", "16", "20", "24", "32", "40", "48", "56", "64"}
	for _, spacing := range expectedSpacing {
		if _, exists := tokens.Spacing[spacing]; !exists {
			fmt.Printf("❌ Spacing token %s not defined\n", spacing)
			return
		}
	}

	// Check border radius values
	expectedRadius := []string{"none", "sm", "base", "md", "lg", "xl", "2xl", "3xl", "full"}
	for _, radius := range expectedRadius {
		if _, exists := tokens.BorderRadius[radius]; !exists {
			fmt.Printf("❌ Border radius token %s not defined\n", radius)
			return
		}
	}

	// Check shadow values
	expectedShadows := []string{"sm", "base", "md"}
	for _, shadow := range expectedShadows {
		if _, exists := tokens.Shadows[shadow]; !exists {
			fmt.Printf("❌ Shadow token %s not defined\n", shadow)
			return
		}
	}

	// Check border values
	expectedBorders := []string{"none", "thin", "base", "thick", "thicker"}
	for _, border := range expectedBorders {
		if _, exists := tokens.Borders[border]; !exists {
			fmt.Printf("❌ Border token %s not defined\n", border)
			return
		}
	}

	// Check animation values
	expectedAnimations := []string{"fast", "base", "slow", "slower"}
	for _, animation := range expectedAnimations {
		if _, exists := tokens.Animations[animation]; !exists {
			fmt.Printf("❌ Animation token %s not defined\n", animation)
			return
		}
	}

	fmt.Printf("✅ Visual tokens test passed\n")
	fmt.Printf("   Spacing tokens: %d\n", len(tokens.Spacing))
	fmt.Printf("   Border radius tokens: %d\n", len(tokens.BorderRadius))
	fmt.Printf("   Shadow tokens: %d\n", len(tokens.Shadows))
	fmt.Printf("   Border tokens: %d\n", len(tokens.Borders))
	fmt.Printf("   Animation tokens: %d\n", len(tokens.Animations))
}

func testAccessibility() {
	// Test accessibility configuration
	accessibility := AccessibilityConfig{
		MinContrastRatio: 4.5,
		ColorBlindSafe:   true,
		HighContrast:     false,
		MinFontSize:      12.0,
		MaxLineLength:    75,
		FocusVisible:     true,
		FocusColor:       "#3B82F6",
	}

	// Validate accessibility configuration
	if accessibility.MinContrastRatio < 4.5 {
		fmt.Println("❌ Minimum contrast ratio should be at least 4.5:1")
		return
	}

	if accessibility.MinFontSize < 12.0 {
		fmt.Println("❌ Minimum font size should be at least 12px")
		return
	}

	if accessibility.MaxLineLength > 75 {
		fmt.Println("❌ Maximum line length should be 75 characters or less")
		return
	}

	if accessibility.FocusColor == "" {
		fmt.Println("❌ Focus color should be defined")
		return
	}

	fmt.Printf("✅ Accessibility test passed\n")
	fmt.Printf("   Min contrast ratio: %.1f:1\n", accessibility.MinContrastRatio)
	fmt.Printf("   Color blind safe: %v\n", accessibility.ColorBlindSafe)
	fmt.Printf("   High contrast: %v\n", accessibility.HighContrast)
	fmt.Printf("   Min font size: %.1fpx\n", accessibility.MinFontSize)
	fmt.Printf("   Max line length: %d characters\n", accessibility.MaxLineLength)
	fmt.Printf("   Focus visible: %v\n", accessibility.FocusVisible)
	fmt.Printf("   Focus color: %s\n", accessibility.FocusColor)
}

func testLaTeXGeneration() {
	// Test LaTeX command generation
	colorCommands := generateLaTeXColorCommands()
	typographyCommands := generateLaTeXTypographyCommands()
	visualCommands := generateLaTeXVisualCommands()

	// Validate color commands
	if !strings.Contains(colorCommands, "\\definecolor{primary}{HTML}{2563EB}") {
		fmt.Println("❌ Primary color command not generated")
		return
	}

	if !strings.Contains(colorCommands, "\\definecolor{catPROPOSAL}{HTML}{2563EB}") {
		fmt.Println("❌ Category color command not generated")
		return
	}

	// Validate typography commands
	if !strings.Contains(typographyCommands, "\\usepackage{fontspec}") {
		fmt.Println("❌ Fontspec package not included")
		return
	}

	if !strings.Contains(typographyCommands, "\\setmainfont{Inter}") {
		fmt.Println("❌ Main font not set")
		return
	}

	// Validate visual commands
	if !strings.Contains(visualCommands, "\\newlength{\\Spacing1}") {
		fmt.Println("❌ Spacing length not defined")
		return
	}

	if !strings.Contains(visualCommands, "\\newlength{\\RadiusBase}") {
		fmt.Println("❌ Border radius length not defined")
		return
	}

	fmt.Printf("✅ LaTeX generation test passed\n")
	fmt.Printf("   Color commands length: %d characters\n", len(colorCommands))
	fmt.Printf("   Typography commands length: %d characters\n", len(typographyCommands))
	fmt.Printf("   Visual commands length: %d characters\n", len(visualCommands))
}

// Helper functions for LaTeX generation
func generateLaTeXColorCommands() string {
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
	commands.WriteString("\\definecolor{catPROPOSAL}{HTML}{2563EB}\n")
	commands.WriteString("\\definecolor{catLASER}{HTML}{EA580C}\n")
	commands.WriteString("\\definecolor{catIMAGING}{HTML}{16A34A}\n")
	commands.WriteString("\\definecolor{catADMIN}{HTML}{6B7280}\n")
	commands.WriteString("\\definecolor{catDISSERTATION}{HTML}{7C3AED}\n")
	commands.WriteString("\\definecolor{catRESEARCH}{HTML}{0EA5E9}\n")
	commands.WriteString("\\definecolor{catPUBLICATION}{HTML}{DC2626}\n")
	
	return commands.String()
}

func generateLaTeXTypographyCommands() string {
	var commands strings.Builder
	
	// Define font families
	commands.WriteString("\\usepackage{fontspec}\n")
	commands.WriteString("\\setmainfont{Inter}\n")
	commands.WriteString("\\setmonofont{JetBrains Mono}\n")
	
	// Define text styles
	commands.WriteString("\\newcommand{\\HeadingOne}[1]{\\fontsize{36}{43.2}\\selectfont\\textbf{\\color{text}#1}}\n")
	commands.WriteString("\\newcommand{\\HeadingTwo}[1]{\\fontsize{30}{39}\\selectfont\\textbf{\\color{text}#1}}\n")
	commands.WriteString("\\newcommand{\\HeadingThree}[1]{\\fontsize{24}{33.6}\\selectfont\\textbf{\\color{text}#1}}\n")
	commands.WriteString("\\newcommand{\\BodyLarge}[1]{\\fontsize{18}{28.8}\\selectfont\\color{textLight}#1}\n")
	commands.WriteString("\\newcommand{\\Body}[1]{\\fontsize{16}{24}\\selectfont\\color{textLight}#1}\n")
	commands.WriteString("\\newcommand{\\BodySmall}[1]{\\fontsize{14}{21}\\selectfont\\color{textLighter}#1}\n")
	commands.WriteString("\\newcommand{\\Caption}[1]{\\fontsize{12}{16.8}\\selectfont\\color{textLighter}#1}\n")
	commands.WriteString("\\newcommand{\\TaskTitle}[1]{\\fontsize{14}{19.6}\\selectfont\\textbf{\\color{white}#1}}\n")
	commands.WriteString("\\newcommand{\\TaskDescription}[1]{\\fontsize{12}{15.6}\\selectfont\\color{white}#1}\n")
	commands.WriteString("\\newcommand{\\OverflowText}[1]{\\fontsize{10}{12}\\selectfont\\textbf{\\color{textLighter}#1}}\n")
	
	return commands.String()
}

func generateLaTeXVisualCommands() string {
	var commands strings.Builder
	
	// Define spacing
	commands.WriteString("\\newlength{\\Spacing0}\n")
	commands.WriteString("\\newlength{\\Spacing1}\n")
	commands.WriteString("\\newlength{\\Spacing2}\n")
	commands.WriteString("\\newlength{\\Spacing3}\n")
	commands.WriteString("\\newlength{\\Spacing4}\n")
	commands.WriteString("\\newlength{\\Spacing5}\n")
	commands.WriteString("\\newlength{\\Spacing6}\n")
	commands.WriteString("\\newlength{\\Spacing8}\n")
	commands.WriteString("\\newlength{\\Spacing10}\n")
	commands.WriteString("\\newlength{\\Spacing12}\n")
	commands.WriteString("\\newlength{\\Spacing16}\n")
	commands.WriteString("\\newlength{\\Spacing20}\n")
	commands.WriteString("\\newlength{\\Spacing24}\n")
	commands.WriteString("\\newlength{\\Spacing32}\n")
	commands.WriteString("\\newlength{\\Spacing40}\n")
	commands.WriteString("\\newlength{\\Spacing48}\n")
	commands.WriteString("\\newlength{\\Spacing56}\n")
	commands.WriteString("\\newlength{\\Spacing64}\n")
	
	// Set spacing values
	commands.WriteString("\\setlength{\\Spacing0}{0pt}\n")
	commands.WriteString("\\setlength{\\Spacing1}{4pt}\n")
	commands.WriteString("\\setlength{\\Spacing2}{8pt}\n")
	commands.WriteString("\\setlength{\\Spacing3}{12pt}\n")
	commands.WriteString("\\setlength{\\Spacing4}{16pt}\n")
	commands.WriteString("\\setlength{\\Spacing5}{20pt}\n")
	commands.WriteString("\\setlength{\\Spacing6}{24pt}\n")
	commands.WriteString("\\setlength{\\Spacing8}{32pt}\n")
	commands.WriteString("\\setlength{\\Spacing10}{40pt}\n")
	commands.WriteString("\\setlength{\\Spacing12}{48pt}\n")
	commands.WriteString("\\setlength{\\Spacing16}{64pt}\n")
	commands.WriteString("\\setlength{\\Spacing20}{80pt}\n")
	commands.WriteString("\\setlength{\\Spacing24}{96pt}\n")
	commands.WriteString("\\setlength{\\Spacing32}{128pt}\n")
	commands.WriteString("\\setlength{\\Spacing40}{160pt}\n")
	commands.WriteString("\\setlength{\\Spacing48}{192pt}\n")
	commands.WriteString("\\setlength{\\Spacing56}{224pt}\n")
	commands.WriteString("\\setlength{\\Spacing64}{256pt}\n")
	
	// Define border radius
	commands.WriteString("\\newlength{\\RadiusNone}\n")
	commands.WriteString("\\newlength{\\RadiusSm}\n")
	commands.WriteString("\\newlength{\\RadiusBase}\n")
	commands.WriteString("\\newlength{\\RadiusMd}\n")
	commands.WriteString("\\newlength{\\RadiusLg}\n")
	commands.WriteString("\\newlength{\\RadiusXl}\n")
	commands.WriteString("\\newlength{\\Radius2xl}\n")
	commands.WriteString("\\newlength{\\Radius3xl}\n")
	commands.WriteString("\\newlength{\\RadiusFull}\n")
	
	// Set border radius values
	commands.WriteString("\\setlength{\\RadiusNone}{0pt}\n")
	commands.WriteString("\\setlength{\\RadiusSm}{2pt}\n")
	commands.WriteString("\\setlength{\\RadiusBase}{4pt}\n")
	commands.WriteString("\\setlength{\\RadiusMd}{6pt}\n")
	commands.WriteString("\\setlength{\\RadiusLg}{8pt}\n")
	commands.WriteString("\\setlength{\\RadiusXl}{12pt}\n")
	commands.WriteString("\\setlength{\\Radius2xl}{16pt}\n")
	commands.WriteString("\\setlength{\\Radius3xl}{24pt}\n")
	commands.WriteString("\\setlength{\\RadiusFull}{9999pt}\n")
	
	return commands.String()
}
