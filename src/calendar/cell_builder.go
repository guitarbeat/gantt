// Package calendar provides cell building functionality for the PhD dissertation planner.
//
// This module handles:
// - Day cell construction with proper spacing and alignment
// - Task cell building with different layouts
// - LaTeX cell formatting and minipage management
package calendar

import (
	"fmt"
	"strconv"

	"phd-dissertation-planner/src/core"
)

// CellBuilder handles construction of calendar cells
type CellBuilder struct {
	cfg *core.Config
}

// NewCellBuilder creates a new cell builder with the given configuration
func NewCellBuilder(cfg *core.Config) *CellBuilder {
	return &CellBuilder{cfg: cfg}
}

// BuildDayNumberCell creates the basic day number cell with minimal padding and hypertarget
func (cb *CellBuilder) BuildDayNumberCell(day string, ref ...string) string {
	dayNumberWidth := "6mm" // Default width
	dayNumberCell := fmt.Sprintf(`\begin{minipage}[t]{%s}\centering{}%s\end{minipage}`, dayNumberWidth, day)
	
	// Add hypertarget if reference is provided
	if len(ref) > 0 && ref[0] != "" {
		hypertarget := fmt.Sprintf(`\hypertarget{%s}{}`, ref[0])
		return hypertarget + dayNumberCell
	}
	
	return dayNumberCell
}

// BuildSimpleDayCell creates a simple day cell with just the day number
func (cb *CellBuilder) BuildSimpleDayCell(leftCell string) string {
	return fmt.Sprintf(`\hyperlink{%s}{%s}`, "", leftCell)
}

// BuildTaskCell constructs a task cell with proper spacing and alignment
func (cb *CellBuilder) BuildTaskCell(leftCell, content string, isSpanning bool, cols int) string {
	dayNumberWidth := "6mm"   // Default width
	dayContentMargin := "1mm" // Default margin

	var width, spacing, contentWrapper string

	// Get typography settings
	hyphenPenalty := 10000    // Default value
	tolerance := 1000         // Default value
	emergencyStretch := "2em" // Default value

	if cb.cfg.Layout.LaTeX.Typography.HyphenPenalty != 0 {
		hyphenPenalty = cb.cfg.Layout.LaTeX.Typography.HyphenPenalty
	}
	if cb.cfg.Layout.LaTeX.Typography.Tolerance != 0 {
		tolerance = cb.cfg.Layout.LaTeX.Typography.Tolerance
	}
	if cb.cfg.Layout.LaTeX.Typography.SloppyEmergencyStretch != "" {
		emergencyStretch = cb.cfg.Layout.LaTeX.Typography.SloppyEmergencyStretch
	}

	if isSpanning {
		// Spanning task: use tikzpicture overlay with calculated width (z-dimension stacking)
		width = `\dimexpr ` + strconv.Itoa(cols) + `\linewidth\relax`
		spacing = `\makebox[0pt][l]{` + `\begin{tikzpicture}[overlay]` +
			`\node[anchor=north west, inner sep=0pt] at (0,0) {` + `\begin{minipage}[t]{` + width + `}` + content + `\end{minipage}` + `};` +
			`\end{tikzpicture}` + `}`
		contentWrapper = "" // Don't add content twice for spanning tasks
	} else if cols > 0 {
		// Spanning task but rendered as regular content (vertical stacking)
		width = `\dimexpr ` + strconv.Itoa(cols) + `\linewidth\relax`
		spacing = ""             // No offset - start at the beginning of the cell
		contentWrapper = content // Use the content directly without additional wrapping
	} else {
		// Regular task: use full available width with fixed height container
		width = `\dimexpr\linewidth - ` + dayContentMargin + `\relax` // Leave space for day number + margins
		spacing = `\hspace*{` + dayNumberWidth + `}`                  // Spacing to align with day number cell width
		// Wrap in fixed-height minipage to prevent row expansion
		contentWrapper = `\begin{minipage}[t][\myLenMonthlyCellHeight][t]{` + width + `}` +
			fmt.Sprintf(`{\sloppy\hyphenpenalty=%d\tolerance=%d\emergencystretch=%s\footnotesize\raggedright `,
				hyphenPenalty, tolerance, emergencyStretch) + content + `}` +
			`\end{minipage}`
	}

	inner := `{\begingroup` +
		`\makebox[0pt][l]{` + leftCell + `}` +
		spacing +
		contentWrapper +
		`\endgroup}`

	// Wrap entire cell in hyperlink to the day's reference (restores link without visual borders via hypersetup)
	return fmt.Sprintf(`\hyperlink{%s}{%s}`, "", inner)
}

// BuildWeekHeaderCell creates a week header cell
func (cb *CellBuilder) BuildWeekHeaderCell(weekNum int) string {
	weekHeaderHeight := "\\myLenMonthlyCellHeight" // Default height

	return fmt.Sprintf(`\hyperlink{week-%d}{\rotatebox[origin=tr]{90}{\makebox[%s][c]{Week %d}}}`,
		weekNum, weekHeaderHeight, weekNum)
}

// BuildMonthHeaderCell creates a month header cell
func (cb *CellBuilder) BuildMonthHeaderCell(monthName string, monthNum int) string {
	monthHeaderHeight := "\\myLenMonthlyCellHeight" // Default height

	return fmt.Sprintf(`\hyperlink{month-%d}{\rotatebox[origin=tr]{90}{\makebox[%s][c]{%s}}}`,
		monthNum, monthHeaderHeight, monthName)
}

// BuildEmptyCell creates an empty cell
func (cb *CellBuilder) BuildEmptyCell() string {
	return `&`
}

// BuildCellSeparator creates a cell separator
func (cb *CellBuilder) BuildCellSeparator() string {
	return `\\`
}

// Helper method to get configuration values with defaults
func (cb *CellBuilder) getConfigValue(key string, defaultValue string) string {
	// This would be implemented based on your configuration structure
	// For now, returning the default value
	return defaultValue
}
