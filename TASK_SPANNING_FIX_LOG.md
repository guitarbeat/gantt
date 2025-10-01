# Task Spanning & Stacking Fix Log

## Problem
Multi-day tasks fail to stack properly. Issues:
1. Tasks repeat instead of spanning days
2. Overlapping tasks visually overlap instead of stacking

Example: Nov 11 task "Write Methods" and Nov 12 task "Send Proposal" overlap visually on Nov 12-18.

## Attempted Fixes

### Fix #1: Show Active Tasks Within Same Week
**Commit**: `4303f5d`

**Change**: Modified `findActiveTasks()` to include tasks that started within the same week and are still active.

**Result**:
- ✅ Fixed same-week stacking
- ❌ Caused task repetition (tasks appeared on every active day)
- ❌ Performance: 52 pages, LaTeX compilation hung

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
**Commit**: `948fc11`

**Change**: Enable TikZ overlay spanning by changing `isSpanning` from hardcoded `false` to `overlay.cols > 1`.

**Result**:
- ✅ Tasks span multiple days visually
- ✅ More efficient: 31 pages
- ❌ Stacking still broken (overlapping tasks from different days overlap)

## Current State

### Working
- ✅ Tasks appear once, span multiple days
- ✅ Single-day stacking works
- ✅ Colors and descriptions correct

### Broken
- ❌ Multi-day stacking: Tasks from different start days overlap visually

## Root Cause

Each day only sees tasks that **start** on that day. Tasks that started earlier and are still active are invisible to subsequent days, causing visual overlap instead of stacking.

## Solutions

### A: Vertical Offset Calculation
Expand `findActiveTasks()` to include continuing tasks, calculate stack positions, use TikZ `yshift`.

### B: Week-Level Rendering
Render entire weeks at once, managing all task positioning globally.

### C: Hybrid TikZ Row Spanning
Use TikZ `fit` to span both columns and rows for stacking.

## Single-Day Stacking Logic

**Working Implementation:**
- First task: `\TaskOverlayBox` (uses `\vfill` for bottom alignment)
- Subsequent tasks: `\TaskOverlayBoxNoOffset` (no extra spacing, `top=0pt, bottom=0pt`)
- Tasks sorted by start date (earlier = lower in stack)

## Implementation (Option A)

**Change:** Modify `findActiveTasks()` to include continuing tasks, use existing stacking logic.

**Updated Algorithm:**
```go
func (d Day) findAllVisibleTasks(dayDate time.Time) ([]*SpanningTask, int) {
    var allVisibleTasks []*SpanningTask
    var maxCols int

    // 1. Add tasks that START on this day
    for _, task := range d.Tasks {
        start := d.getTaskStartDate(task)
        if dayDate.Equal(start) {
            allVisibleTasks = append(allVisibleTasks, task)
            cols := d.calculateTaskSpanColumns(dayDate, end)
            if cols > maxCols { maxCols = cols }
        }
    }

    // 2. Add tasks that STARTED EARLIER but are still active
    for _, task := range d.Tasks {
        start := d.getTaskStartDate(task)
        end := d.getTaskEndDate(task)

        // Started before today AND still active AND not already added
        if start.Before(dayDate) && d.isTaskActiveOnDay(dayDate, start, end) {
            // Check if not already in list (avoid duplicates)
            alreadyAdded := false
            for _, existing := range allVisibleTasks {
                if existing == task { alreadyAdded = true; break }
            }
            if !alreadyAdded {
                allVisibleTasks = append(allVisibleTasks, task)
                // For continuing tasks, calculate remaining span
                cols := d.calculateRemainingSpanColumns(dayDate, end)
                if cols > maxCols { maxCols = cols }
            }
        }
    }

    // 3. Sort by start date (earliest first = bottom of stack)
    allVisibleTasks = d.sortTasksByStartDate(allVisibleTasks)

    return allVisibleTasks, maxCols
}
```

**LaTeX:** Use existing stacking macros unchanged.

## Implementation Priority

**Recommended:** Option A - expand `findActiveTasks()` to include continuing tasks.

**Why:** Minimal changes to existing working stacking logic.
