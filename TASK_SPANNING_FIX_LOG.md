# Task Spanning & Stacking Fix Log

## Problem
Multi-day tasks fail to stack properly. Issues:
1. Tasks repeat instead of spanning days
2. Overlapping tasks visually overlap instead of stacking

Example: Nov 11 task "Write Methods" and Nov 12 task "Send Proposal" overlap visually on Nov 12-18.

## Root Cause
Each day only sees tasks that **start** on that day. Tasks that started earlier and are still active are invisible to subsequent days, causing visual overlap instead of stacking.

## Solution Implemented

### Final Fix: Account for Active Tasks in Stack Position
**Commit**: Latest

**Approach**: Modified `findActiveTasks()` to include ALL active tasks (both starting and continuing), then only render tasks that start on the current day, but use the full active list to determine stack position.

**Key Changes**:

1. **Enhanced `findActiveTasks()`**: Now returns ALL tasks active on a given day
   ```go
   func (d Day) findActiveTasks(dayDate time.Time) ([]*SpanningTask, int) {
       // Include tasks that START on this day
       // AND tasks that STARTED EARLIER but are still active
       // This ensures proper vertical stacking
   }
   ```

2. **Added `calculateRemainingSpanColumns()`**: For continuing tasks
   ```go
   func (d Day) calculateRemainingSpanColumns(dayDate, end time.Time) int {
       // Calculate how many columns a continuing task spans
       // from current day to its end (or end of week)
   }
   ```

3. **Smart Rendering Logic**: Only render tasks that START today, but know their position
   ```go
   func (d Day) renderSpanningTaskOverlay() *TaskOverlay {
       activeTasks, _ := d.findActiveTasks(dayDate)
       
       // For each task that STARTS today:
       // - If stackPos == 0: use TaskOverlayBox (with \vfill - bottom)
       // - If stackPos > 0: use TaskOverlayBoxNoOffset (no \vfill - stacked)
   }
   ```

4. **LaTeX Macro Distinction**:
   - `\TaskOverlayBox`: Uses `\vfill` to position at bottom of cell
   - `\TaskOverlayBoxNoOffset`: No `\vfill` - stacks naturally on top

**Result**:
- ✅ Fixed overlapping: Tasks stack vertically based on start date
- ✅ No repetition: Each task appears only once at its start
- ✅ Proper spanning: Tasks span multiple days correctly
- ✅ Efficient: 34 pages (reasonable size)

## How It Works

### Example: Nov 11-12 Week
1. **Nov 11 (Monday)**:
   - Active tasks: ["Write Methods" (starts today)]
   - "Write Methods" gets stackPos=0 → uses `\TaskOverlayBox` (bottom)
   
2. **Nov 12 (Tuesday)**:
   - Active tasks: ["Write Methods" (continuing), "Send Proposal" (starts today)]
   - Sorted by start date: "Write Methods" (Nov 11) comes first
   - "Send Proposal" gets stackPos=1 → uses `\TaskOverlayBoxNoOffset` (stacked on top)
   - Only "Send Proposal" renders (starts today), but it knows it's stackPos=1

### Vertical Positioning
- `\vfill` in `\TaskOverlayBox` pushes task to bottom of cell
- `\TaskOverlayBoxNoOffset` has no `\vfill`, so it flows naturally
- Multiple `\TaskOverlayBoxNoOffset` boxes stack vertically without overlap

## Testing
- Generated PDF: 34 pages
- Task stacking: Properly handles multi-day overlaps
- Visual appearance: Tasks don't overlap, colors preserved
- Performance: LaTeX compiles successfully

## Previous Attempts (for reference)

### Attempt #1: Spacer Approach (Failed)
- Added `\TaskOverlayBoxSpacer{}` for continuing tasks
- Result: 71 pages, excessive spacing
- Issue: `\vfill` + fixed height created too much vertical space

### Attempt #2: Render All Active (Failed)
- Rendered ALL active tasks on each day
- Result: Task repetition, performance issues
- Issue: Same task appeared multiple times across its span

### Attempt #3: Intelligent Stack Position (Success!)
- Find all active tasks but only render those that start today
- Use active list to determine stack position
- Apply correct LaTeX macro based on position
- Result: Perfect stacking without repetition
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
