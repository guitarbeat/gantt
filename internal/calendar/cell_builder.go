// Package calendar provides cell building functionality for the PhD dissertation planner.
//
// This module handles:
// - Day cell construction with proper spacing and alignment
// - Task cell building with different layouts
// - LaTeX cell formatting and minipage management
package calendar

import (
	"strconv"
	"strings"

	"phd-dissertation-planner/internal/core"
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
	var sb strings.Builder

	// Add hypertarget if reference is provided
	if len(ref) > 0 && ref[0] != "" {
		sb.WriteString(`\hypertarget{`)
		sb.WriteString(ref[0])
		sb.WriteString(`}{}`)
	}

	sb.WriteString(`\begin{minipage}[t]{`)
	sb.WriteString(dayNumberWidth)
	sb.WriteString(`}\centering{}`)
	sb.WriteString(day)
	sb.WriteString(`\end{minipage}`)

	return sb.String()
}

// BuildSimpleDayCell creates a simple day cell with just the day number
func (cb *CellBuilder) BuildSimpleDayCell(leftCell string) string {
	var sb strings.Builder
	sb.WriteString(`\hyperlink{}{`)
	sb.WriteString(leftCell)
	sb.WriteString(`}`)
	return sb.String()
}

// BuildTaskCell constructs a task cell with proper spacing and alignment
func (cb *CellBuilder) BuildTaskCell(leftCell, content string, isSpanning bool, cols int) string {
	dayNumberWidth := "6mm"   // Default width
	dayContentMargin := "1mm" // Default margin

	var sb strings.Builder

	// Pre-allocate buffer to avoid reallocations
	// Estimated size: headers + leftCell + content + wrapper overhead
	sb.Grow(512 + len(leftCell) + len(content))

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

	// Start with outer hyperlink wrapper
	sb.WriteString(`\hyperlink{}{`)

	// Inner group
	sb.WriteString(`{\begingroup\makebox[0pt][l]{`)
	sb.WriteString(leftCell)
	sb.WriteString(`}`)

	if isSpanning {
		// Spanning task: use tikzpicture overlay with calculated width (z-dimension stacking)
		sb.WriteString(`\makebox[0pt][l]{\begin{tikzpicture}[overlay]\node[anchor=north west, inner sep=0pt] at (0,0) {\begin{minipage}[t]{\dimexpr `)
		sb.WriteString(strconv.Itoa(cols))
		sb.WriteString(`\linewidth\relax}`)
		sb.WriteString(content)
		sb.WriteString(`\end{minipage}};\end{tikzpicture}}`)
	} else if cols > 0 {
		// Spanning task but rendered as regular content (vertical stacking)
		// No offset - start at the beginning of the cell
		sb.WriteString(content)
	} else {
		// Regular task: use full available width with fixed height container
		sb.WriteString(`\hspace*{`)
		sb.WriteString(dayNumberWidth)
		sb.WriteString(`}\begin{minipage}[t][\myLenMonthlyCellHeight][t]{\dimexpr\linewidth - `)
		sb.WriteString(dayContentMargin)
		sb.WriteString(`\relax}{\sloppy\hyphenpenalty=`)
		sb.WriteString(strconv.Itoa(hyphenPenalty))
		sb.WriteString(`\tolerance=`)
		sb.WriteString(strconv.Itoa(tolerance))
		sb.WriteString(`\emergencystretch=`)
		sb.WriteString(emergencyStretch)
		sb.WriteString(`\footnotesize\raggedright `)
		sb.WriteString(content)
		sb.WriteString(`}\end{minipage}`)
	}

	sb.WriteString(`\endgroup}}`)

	return sb.String()
}

// BuildWeekHeaderCell creates a week header cell
func (cb *CellBuilder) BuildWeekHeaderCell(weekNum int) string {
	weekHeaderHeight := "\\myLenMonthlyCellHeight" // Default height

	var sb strings.Builder
	sb.WriteString(`\hyperlink{week-`)
	weekStr := strconv.Itoa(weekNum)
	sb.WriteString(weekStr)
	sb.WriteString(`}{\rotatebox[origin=tr]{90}{\makebox[`)
	sb.WriteString(weekHeaderHeight)
	sb.WriteString(`][c]{Week `)
	sb.WriteString(weekStr)
	sb.WriteString(`}}}`)

	return sb.String()
}

// BuildMonthHeaderCell creates a month header cell
func (cb *CellBuilder) BuildMonthHeaderCell(monthName string, monthNum int) string {
	monthHeaderHeight := "\\myLenMonthlyCellHeight" // Default height

	var sb strings.Builder
	sb.WriteString(`\hyperlink{month-`)
	sb.WriteString(strconv.Itoa(monthNum))
	sb.WriteString(`}{\rotatebox[origin=tr]{90}{\makebox[`)
	sb.WriteString(monthHeaderHeight)
	sb.WriteString(`][c]{`)
	sb.WriteString(monthName)
	sb.WriteString(`}}}`)

	return sb.String()
}

// BuildEmptyCell creates an empty cell
func (cb *CellBuilder) BuildEmptyCell() string {
	return `&`
}

// BuildCellSeparator creates a cell separator
func (cb *CellBuilder) BuildCellSeparator() string {
	return `\\`
}
