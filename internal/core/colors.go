// Package core - Colors provides terminal color utilities for better output formatting.
//
// This package provides ANSI color codes and utility functions for creating
// visually appealing terminal output with proper color support detection.
//
// Example usage:
//
//	fmt.Println(colors.Success("Operation completed successfully"))
//	fmt.Println(colors.Warning("This is a warning message"))
//	fmt.Println(colors.Error("An error occurred"))
//	fmt.Println(colors.Info("Informational message"))
package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Color codes for terminal output
const (
	// Reset all formatting
	Reset = "\033[0m"

	// Text colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright colors
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Background colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	// Text styles
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Strike    = "\033[9m"
)

// colorEnabled checks if colors should be used based on environment and terminal support
func colorEnabled() bool {
	// Check if NO_COLOR environment variable is set
	if noColor := os.Getenv("NO_COLOR"); noColor != "" {
		return false
	}

	// Check if FORCE_COLOR environment variable is set
	if forceColor := os.Getenv("FORCE_COLOR"); forceColor != "" {
		if val, err := strconv.Atoi(forceColor); err == nil && val > 0 {
			return true
		}
	}

	// Check if we're in a TTY
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}

	return false
}

// colorize applies color to text if colors are enabled
func colorize(color, text string) string {
	if !colorEnabled() {
		return text
	}
	return color + text + Reset
}

// Success returns green colored text for success messages
func Success(text string) string {
	return colorize(Green, text)
}

// Warning returns yellow colored text for warning messages
func Warning(text string) string {
	return colorize(Yellow, text)
}

// Error returns red colored text for error messages
func Error(text string) string {
	return colorize(Red, text)
}

// Info returns blue colored text for informational messages
func Info(text string) string {
	return colorize(Blue, text)
}

// DimText returns dimmed text for secondary information
func DimText(text string) string {
	return colorize(Dim, text)
}

// BoldText returns bold text for emphasis
func BoldText(text string) string {
	return colorize(Bold, text)
}

// Bright returns bright colored text for highlights
func Bright(text string) string {
	return colorize(BrightWhite, text)
}

// CyanText returns cyan colored text for special highlights
func CyanText(text string) string {
	return colorize(Cyan, text)
}

// MagentaText returns magenta colored text for special highlights
func MagentaText(text string) string {
	return colorize(Magenta, text)
}

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
