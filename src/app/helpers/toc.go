// Package helpers provides utility functions extracted from the main generator.
//
// Table of Contents helpers provide functionality for generating LaTeX content
// for task indexes and navigation.
package helpers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"phd-dissertation-planner/src/core"
)

// TOCBuilder builds table of contents content for tasks
type TOCBuilder struct {
	phaseNames map[string]string
}

// NewTOCBuilder creates a new TOC builder with default phase names
func NewTOCBuilder() *TOCBuilder {
	return &TOCBuilder{
		phaseNames: map[string]string{
			"1": "Phase 1: Proposal \\& Setup",
			"2": "Phase 2: Research \\& Data Collection", 
			"3": "Phase 3: Publications",
			"4": "Phase 4: Dissertation",
		},
	}
}

// BuildTOCContent generates the complete LaTeX content for the table of contents
func (b *TOCBuilder) BuildTOCContent(tasks []core.Task) string {
	var content strings.Builder

	content.WriteString("% Table of Contents - Clickable Task Index\n")
	content.WriteString("\\hypertarget{task-index}{}\n")
	content.WriteString("{\\Large\\textbf{Task Index}}\n\n")
	content.WriteString("\\vspace{0.2cm}\n\n")

	// Group tasks by phase
	phaseTasks := b.groupTasksByPhase(tasks)

	// Create phase-based sections
	phases := []string{"1", "2", "3", "4"}
	for _, phase := range phases {
		if tasks, exists := phaseTasks[phase]; exists && len(tasks) > 0 {
			content.WriteString(b.buildPhaseSection(phase, tasks))
		}
	}

	content.WriteString(b.buildUsageLegend())

	return content.String()
}

// groupTasksByPhase groups tasks by their phase
func (b *TOCBuilder) groupTasksByPhase(tasks []core.Task) map[string][]core.Task {
	phaseTasks := make(map[string][]core.Task)
	for _, task := range tasks {
		phaseTasks[task.Phase] = append(phaseTasks[task.Phase], task)
	}
	return phaseTasks
}

// buildPhaseSection builds the LaTeX content for a single phase
func (b *TOCBuilder) buildPhaseSection(phase string, tasks []core.Task) string {
	var content strings.Builder

	// Phase header
	content.WriteString(fmt.Sprintf("{\\colorbox[RGB]{245,245,245}{\\makebox[\\linewidth][l]{\\textbf{%s}}}}\\\\\n", b.phaseNames[phase]))
	content.WriteString("\\vspace{0.05cm}\n\n")
	content.WriteString("\\vspace{0.1cm}\n\n")

	// Group tasks by sub-phase within this phase
	subPhaseTasks := b.groupTasksBySubPhase(tasks)

	// Sort sub-phases alphabetically for consistent ordering
	var subPhases []string
	for subPhase := range subPhaseTasks {
		subPhases = append(subPhases, subPhase)
	}
	sort.Strings(subPhases)

	// Render each sub-phase
	for _, subPhase := range subPhases {
		subPhaseTaskList := subPhaseTasks[subPhase]
		content.WriteString(b.buildSubPhaseSection(subPhase, subPhaseTaskList))
	}

	content.WriteString("\\vspace{0.2cm}\n\n")
	return content.String()
}

// groupTasksBySubPhase groups tasks by their sub-phase within a phase
func (b *TOCBuilder) groupTasksBySubPhase(tasks []core.Task) map[string][]core.Task {
	subPhaseTasks := make(map[string][]core.Task)
	for _, task := range tasks {
		subPhase := task.SubPhase
		if subPhase == "" {
			subPhase = "General" // Default for tasks without sub-phase
		}
		subPhaseTasks[subPhase] = append(subPhaseTasks[subPhase], task)
	}
	return subPhaseTasks
}

// buildSubPhaseSection builds the LaTeX content for a single sub-phase
func (b *TOCBuilder) buildSubPhaseSection(subPhase string, tasks []core.Task) string {
	var content strings.Builder

	// Sort tasks chronologically within this sub-phase
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartDate.Before(tasks[j].StartDate)
	})

	// Sub-phase header (smaller than phase header)
	content.WriteString("\\vspace{0.1cm}\n")
	content.WriteString(fmt.Sprintf("{\\colorbox[RGB]{250,250,250}{\\makebox[\\linewidth][l]{\\textbf{\\small %s}}}}\\\\\n", subPhase))
	content.WriteString("\\vspace{0.03cm}\n\n")

	// Tasks for this sub-phase - compact format
	content.WriteString("\\begin{itemize}[leftmargin=0.5cm,itemsep=0.1cm,parsep=0.05cm]\n")
	for _, task := range tasks {
		content.WriteString(b.formatTaskEntry(task))
	}
	content.WriteString("\\end{itemize}\n")
	content.WriteString("\\vspace{0.1cm}\n\n")

	return content.String()
}

// formatTaskEntry formats a single task entry with color and hyperlink
func (b *TOCBuilder) formatTaskEntry(task core.Task) string {
	// Create hyperlink reference to first occurrence of task
	// Use RFC3339 format to match calendar's d.ref() method
	dateRef := task.StartDate.Format(time.RFC3339)

	// Get color for the task category
	taskColor := core.GenerateCategoryColor(strings.ToUpper(task.Category))
	taskName := b.escapeLaTeXSpecialChars(task.Name)

	// Bold the task name if it's a milestone
	if task.IsMilestone {
		taskName = fmt.Sprintf("\\textbf{%s}", taskName)
	}

	// Format task with color and hyperlink
	if len(taskColor) >= 7 && taskColor[0] == '#' {
		rgbStr := b.hexToRGBString(taskColor)
		return fmt.Sprintf("\\item \\textcolor[RGB]{%s}{\\hyperlink{%s}{%s}}\n", rgbStr, dateRef, taskName)
	}
	return fmt.Sprintf("\\item \\hyperlink{%s}{%s}\n", dateRef, taskName)
}

// escapeLaTeXSpecialChars escapes special LaTeX characters in text
func (b *TOCBuilder) escapeLaTeXSpecialChars(text string) string {
	text = strings.ReplaceAll(text, "&", "\\&")
	text = strings.ReplaceAll(text, "%", "\\%")
	return text
}

// hexToRGBString converts a hex color string to RGB format for LaTeX
func (b *TOCBuilder) hexToRGBString(hex string) string {
	if len(hex) < 7 || hex[0] != '#' {
		return "0,0,0" // Default black
	}

	// Parse hex values
	r, err1 := parseHexByte(hex[1:3])
	g, err2 := parseHexByte(hex[3:5])
	blue, err3 := parseHexByte(hex[5:7])

	if err1 != nil || err2 != nil || err3 != nil {
		return "0,0,0" // Default black on error
	}

	return fmt.Sprintf("%d,%d,%d", r, g, blue)
}

// parseHexByte parses a two-character hex string to an integer
func parseHexByte(hex string) (int64, error) {
	if len(hex) != 2 {
		return 0, fmt.Errorf("hex string must be exactly 2 characters")
	}
	return strconv.ParseInt(hex, 16, 64)
}

// buildUsageLegend builds the usage legend for the table of contents
func (b *TOCBuilder) buildUsageLegend() string {
	var content strings.Builder

	content.WriteString("% Usage Legend\n")
	content.WriteString("{\\small\n")
	content.WriteString("\\textbf{How to use this index:}\\\\\n")
	content.WriteString("\\textbullet\\ \\textbf{Bold task names} indicate milestones with enhanced borders in timeline\\\\\n")
	content.WriteString("\\textbullet\\ Click on any task name to jump to its location in the timeline\n")
	content.WriteString("}\n\n")
	content.WriteString("\\pagebreak\n")

	return content.String()
}
