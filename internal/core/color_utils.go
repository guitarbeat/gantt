// Package core provides core utilities for color manipulation and conversion
package core

import (
	"fmt"
	"strconv"
	"strings"
)

// HexToRGB converts a hex color string to RGB format for LaTeX compatibility.
// Accepts hex colors with or without the # prefix.
// Returns "128,128,128" (gray) for invalid hex strings.
func HexToRGB(hex string) string {
	// Remove # prefix if present
	hex = strings.TrimPrefix(hex, "#")

	// Validate hex length
	if len(hex) != 6 {
		return "128,128,128" // Default gray for invalid hex
	}

	// Parse hex values
	r, err1 := strconv.ParseInt(hex[0:2], 16, 64)
	g, err2 := strconv.ParseInt(hex[2:4], 16, 64)
	b, err3 := strconv.ParseInt(hex[4:6], 16, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return "128,128,128" // Default gray on parse error
	}

	return fmt.Sprintf("%d,%d,%d", r, g, b)
}
