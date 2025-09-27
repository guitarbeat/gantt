// Package calendar provides color management functionality for the PhD dissertation planner.
//
// This module handles:
// - Color mapping for task categories
// - Hex to RGB conversion for LaTeX compatibility
// - Color legend generation
// - Default color assignment
package calendar

import (
	"fmt"
	"strings"
)

// ColorManager handles color-related operations
type ColorManager struct {
	categoryColors map[string]string
}

// NewColorManager creates a new color manager with default category colors
func NewColorManager() *ColorManager {
	return &ColorManager{
		categoryColors: getDefaultCategoryColors(),
	}
}

// GetColorForCategory returns the color for a given category
func (cm *ColorManager) GetColorForCategory(category string) string {
	if color, exists := cm.categoryColors[category]; exists {
		return color
	}
	// Default color for unknown categories
	return "#E0E0E0" // Light Gray
}

// GetRGBColorForCategory returns the RGB color for a given category
func (cm *ColorManager) GetRGBColorForCategory(category string) string {
	hexColor := cm.GetColorForCategory(category)
	return cm.HexToRGB(hexColor)
}

// HexToRGB converts a hex color to RGB format for LaTeX
func (cm *ColorManager) HexToRGB(hex string) string {
	if hex == "" {
		return ""
	}

	// Remove # if present
	hex = strings.TrimPrefix(hex, "#")

	// Convert hex to RGB
	if len(hex) == 6 {
		r := hex[0:2]
		g := hex[2:4]
		b := hex[4:6]
		return fmt.Sprintf("%s,%s,%s", r, g, b)
	}

	return ""
}

// GetCategoryColors returns all category colors
func (cm *ColorManager) GetCategoryColors() map[string]string {
	return cm.categoryColors
}

// SetCategoryColor sets a custom color for a category
func (cm *ColorManager) SetCategoryColor(category, color string) {
	cm.categoryColors[category] = color
}

// GenerateColorLegend generates a LaTeX color legend
func (cm *ColorManager) GenerateColorLegend() string {
	var legendItems []string

	for category, color := range cm.categoryColors {
		rgbColor := cm.HexToRGB(color)
		legendItem := fmt.Sprintf(`\textcolor[RGB]{%s}{%s}`, rgbColor, category)
		legendItems = append(legendItems, legendItem)
	}

	return strings.Join(legendItems, " \\quad ")
}

// GetDefaultCategoryColors returns the default color mapping
func getDefaultCategoryColors() map[string]string {
	return map[string]string{
		"Research":     "#FF6B6B", // Red
		"Writing":      "#4ECDC4", // Teal
		"Analysis":     "#45B7D1", // Blue
		"Review":       "#96CEB4", // Green
		"Presentation": "#FFEAA7", // Yellow
		"Planning":     "#DDA0DD", // Plum
		"Data":         "#98D8C8", // Mint
		"Meeting":      "#F7DC6F", // Light Yellow
		"Admin":        "#BB8FCE", // Light Purple
		"Other":        "#85C1E9", // Light Blue
	}
}

// Color constants for easy reference
const (
	ColorRed       = "#FF6B6B"
	ColorTeal      = "#4ECDC4"
	ColorBlue      = "#45B7D1"
	ColorGreen     = "#96CEB4"
	ColorYellow    = "#FFEAA7"
	ColorPlum      = "#DDA0DD"
	ColorMint      = "#98D8C8"
	ColorPurple    = "#BB8FCE"
	ColorLightBlue = "#85C1E9"
	ColorGray      = "#E0E0E0"
)
