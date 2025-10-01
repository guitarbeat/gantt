# Task Spanning & Stacking Fix Log

## Problem Statement

Tasks that span multiple days were not rendering correctly. Two main issues:

1. **Repetition**: Tasks appeared multiple times (once per day) instead of showing once with a bar spanning multiple days
2. **Overlap**: Tasks starting on different days but overlapping in time would visually overlap instead of stacking vertically

### Example Issue
- **Week 45**: "Write Methods & Timeline" (Nov 11-18) and "Send Proposal to Committee" (Nov 12 - Dec 1)
  - "Write Methods" starts Nov 11
  - "Send Proposal" starts Nov 12 while "Write Methods" is still active
  - Expected: Both tasks visible and stacked on Nov 12
  - Actual: "Send Proposal" overlaps/covers "Write Methods"

## Attempted Fixes

### Fix #1: Show Active Tasks Within Same Week
**Commit**: `4303f5d` - "Refine: Show tasks from current week only, not all active tasks"

**Approach**: Modified `findActiveTasks()` to show tasks that:
- Start on this day, OR
- Started within the same week and are still active today

**Code Change**:
```go
// Calculate the start of this week (Monday)
weekStart := dayDate.AddDate(0, 0, -int((dayDate.Weekday()+6)%7))

for _, task := range d.Tasks {
    start := d.getTaskStartDate(task)
    end := d.getTaskEndDate(task)

    // Show tasks that:
    // 1. Start on this day, OR
    // 2. Started this week (after weekStart) and are still active today
    if dayDate.Equal(start) || (start.After(weekStart.AddDate(0, 0, -1)) && 
        start.Before(dayDate) && d.isTaskActiveOnDay(dayDate, start, end)) {
        activeTasks = append(activeTasks, task)
        // ...
    }
}
```

**Result**: 
- ✅ Fixed stacking within the same week
- ❌ Created repetition issue: Tasks appeared on every day they were active (e.g., "Draft Timeline v1" showed 3 times)
- Generated 285 task boxes, 7764 lines, 52 pages
- Performance issues: LaTeX compilation hung

### Fix #2: Show Tasks Only on Start Day
**Commit**: `facdead` - "Fix: Show each task only once at its start day"

**Approach**: Reverted logic to only show tasks on their START day

**Code Change**:
```go
func (d Day) findActiveTasks(dayDate time.Time) ([]*SpanningTask, int) {
    var activeTasks []*SpanningTask
    var maxCols int

    for _, task := range d.Tasks {
        start := d.getTaskStartDate(task)
        end := d.getTaskEndDate(task)

        // Only show tasks that START on this day
        if dayDate.Equal(start) {
            activeTasks = append(activeTasks, task)
            
            // Calculate how many columns this task spans from its start
            cols := d.calculateTaskSpanColumns(dayDate, end)
            if cols > maxCols {
                maxCols = cols
            }
        }
    }

    // Sort tasks by start date (earlier tasks appear first/on top)
    activeTasks = d.sortTasksByStartDate(activeTasks)

    return activeTasks, maxCols
}
```

**Result**:
- ✅ Fixed repetition: Each task appears only once
- ❌ Broke stacking: Tasks starting on different days don't see each other
- Generated 43 pages
- "Draft Timeline v1" appears only once on Aug 29 ✓

### Fix #3: Enable TikZ Spanning Mode
**Commit**: `948fc11` - "Fix: Enable actual spanning for multi-day tasks"

**Approach**: Fixed hardcoded `isSpanning=false` parameter to enable TikZ overlay spanning

**Code Change**:
```go
// renderLargeDay renders the day cell for large (monthly) view with tasks
func (d Day) renderLargeDay(day string) string {
    leftCell := d.buildDayNumberCell(day)

    // Check for tasks using intelligent stacking
    overlay := d.renderSpanningTaskOverlay()
    if overlay != nil {
        // Use spanning mode if any task spans more than 1 column
        isSpanning := overlay.cols > 1  // CHANGED FROM: false
        return d.buildTaskCell(leftCell, overlay.content, isSpanning, overlay.cols)
    }

    // No tasks: just the day number
    return d.buildSimpleDayCell(leftCell)
}
```

**LaTeX Output Example**:
```tex
% Aug 29: Draft Timeline v1 (3 days: Fri-Sat-Sun)
\makebox[0pt][l]{\begin{tikzpicture}[overlay]
  \node[anchor=north west, inner sep=0pt] at (0,0) {
    \begin{minipage}[t]{\dimexpr 3\linewidth\relax}
      \TaskOverlayBox{224,50,212}{Draft Timeline v1}{...}
    \end{minipage}
  };
\end{tikzpicture}}

% Nov 11: Write Methods (6 days: Tue-Sun)
\begin{minipage}[t]{\dimexpr 6\linewidth\relax}
  \TaskOverlayBox{224,50,212}{Write Methods \& Timeline}{...}
\end{minipage}

% Nov 12: Send Proposal (5 days: Wed-Sun)
\begin{minipage}[t]{\dimexpr 5\linewidth\relax}
  \TaskOverlayBox{224,50,212}{Send Proposal to Committee}{...}
\end{minipage}
```

**Result**:
- ✅ Tasks now visually span multiple days using TikZ overlay
- ✅ More efficient: 31 pages (down from 43)
- ❌ **Stacking still broken**: Overlapping tasks from different days overlap visually

## Current State

### What Works
1. ✅ Each task appears only once (no repetition)
2. ✅ Tasks visually span across multiple days
3. ✅ Tasks have proper colors and descriptions
4. ✅ Within a single day with multiple starting tasks, they stack properly

### What's Still Broken
1. ❌ **Tasks starting on different days don't stack**
   - Nov 11: "Write Methods & Timeline" starts
   - Nov 12: "Send Proposal to Committee" starts
   - On Nov 12-18: Both tasks are active but "Send Proposal" doesn't know about "Write Methods"
   - Result: Visual overlap instead of vertical stacking

## Root Cause Analysis

The fundamental issue is in the **task rendering logic**:

```go
// Only show tasks that START on this day
if dayDate.Equal(start) {
    activeTasks = append(activeTasks, task)
}
```

Each day only knows about tasks that **start** on that day. It doesn't account for tasks that started earlier and are still active.

### Why This Breaks Stacking

1. **Nov 11**: Renders "Write Methods" with TikZ overlay spanning 6 columns
2. **Nov 12**: Renders "Send Proposal" with TikZ overlay spanning 5 columns
3. Both overlays start at `(0,0)` in their respective cells
4. TikZ overlays draw in z-space (on top of table), so they overlap visually
5. No vertical offset is applied because Nov 12 doesn't see Nov 11's task

### Why Fix #1 Failed

Fix #1 tried to show all active tasks within a week, which caused:
- Tasks to repeat on every day they're active
- Performance issues with too many task boxes
- Visual confusion with repeated labels

## Potential Solutions

### Option A: Calculate Proper Vertical Offsets
- On each day, check for ALL active tasks (not just starting)
- Calculate vertical offset for each based on when it started
- Use TikZ `yshift` to position overlays vertically
- Challenge: Need to calculate consistent heights across days

### Option B: Week-Level Rendering
- Render tasks at the week level instead of day level
- Calculate all active tasks for the entire week
- Position them with proper vertical offsets from the start
- Challenge: Major refactor of rendering logic

### Option C: Hybrid Approach
- Show tasks only on start day (no repetition)
- Use TikZ overlay to span columns AND rows
- Calculate y-offset based on "active task count at start position"
- Store task metadata to inform subsequent days
- Challenge: Need state tracking between days

## Next Steps

Need to implement one of the above solutions to achieve:
- Each task appears once with a spanning bar
- Overlapping tasks stack vertically without overlap
- Proper spacing accounting for task height (title + description)
- Consistent rendering across week boundaries
