# Task Overlap Detection Analysis

## Problem Description

Tasks that span multiple days and overlap are not being stacked properly when they start on different days.

### Example Case from Input Data:
- **Task 1**: "Write Methods & Timeline" (Nov 11-18, 2025) - starts Tuesday, Week 45
- **Task 2**: "Send Proposal to Committee" (Nov 12 - Dec 1, 2025) - starts Wednesday, Week 45

**Expected Behavior**: Task 2 should stack ABOVE Task 1 since Task 1 is already occupying space when Task 2 starts.

**Actual Behavior**: Task 2 appears to cut off or overlay Task 1 instead of stacking properly.

## Root Cause Analysis

### Current Implementation (calendar.go)

The code has the correct logic structure:

1. **`findActiveTasks()`** (lines 321-355): Correctly identifies ALL tasks active on a given day
   - Includes tasks that started earlier (continuing tasks)
   - Includes tasks that start today
   - Returns them sorted by start date (earlier tasks first)

2. **`renderSpanningTaskOverlay()`** (lines 185-275): Renders task pills
   - Finds active tasks for the day
   - Only renders tasks that START on that day (correct!)
   - Assigns stackPos based on position in activeTasks array
   - Uses `\TaskOverlayBoxWithOffset{color}{name}{desc}{verticalOffset}` macro

### The Bug

The issue is in how the **vertical offset** is calculated (line 257):

```go
verticalOffset := fmt.Sprintf("%.1fem", float64(tr.stackPos)*2.0) // 2.0em per level
```

**Problem**: The `\TaskOverlayBoxWithOffset` LaTeX macro is defined to use `\raisebox{#4}` which accepts a dimension like `2.0em`. However, the macro implementation in `macros.tpl` (lines 86-100) uses `\raisebox{#4}{...}` which:

1. **Raises the box UPWARD** (positive values go up, negative go down in LaTeX)
2. But the visual stacking needs the boxes to go DOWN from the top of the cell

Additionally, the offset calculation doesn't account for:
- The actual height of task boxes (which varies based on content)
- The need for proper spacing between stacked tasks
- The coordinate system (LaTeX's \raisebox goes up for positive, but we want later tasks to appear below earlier tasks in the cell)

### Key Insight

The current approach uses **absolute positioning** with `\raisebox`, but calendar cells have **dynamic height** based on their content. This creates conflicts:

1. Earlier tasks (lower stackPos) should appear at the bottom of the cell
2. Later tasks (higher stackPos) should appear above earlier tasks
3. But `\raisebox` with positive offsets raises content UP, not down

## Solution Approach

We have several options:

### Option 1: Fix the offset calculation (Quick Fix)
- Use NEGATIVE offsets for later tasks: `-2.0em * stackPos`
- This pushes later tasks DOWN in the cell

### Option 2: Use TikZ absolute positioning (More Robust)
- Use TikZ overlays with precise coordinates
- Calculate exact Y positions for each task
- More complex but handles variable-height content better

### Option 3: Use LaTeX flow positioning (Simplest)
- Stack tasks vertically using `\vspace` or `\par`
- Let LaTeX handle the positioning naturally
- May need to ensure cells have enough height

## Recommended Fix

I recommend **Option 3** (LaTeX flow positioning) because:

1. **Simpler**: Let LaTeX handle vertical stacking naturally
2. **More Robust**: Works with variable content heights
3. **Easier to maintain**: No complex offset calculations

### Implementation Plan:

1. Remove the complex `stackPos` offset calculation
2. Stack tasks vertically in order (earliest first, latest last)
3. Use `\vspace` or similar to add spacing between tasks
4. Let the cell expand vertically to fit all tasks

This approach is already partially implemented - the `\TaskOverlayBox` macro uses `\vfill` which should work. The issue might be that we're trying to be too clever with offsets when we should just stack naturally.

## Testing

After implementing the fix, test with:
- Week 45 (November 2025) where the two tasks overlap
- Any other overlapping tasks in different weeks
- Tasks with varying content lengths (name + description)

## Files to Modify

1. `src/calendar/calendar.go` - `renderSpanningTaskOverlay()` function
2. Possibly `src/shared/templates/monthly/macros.tpl` if macro needs adjustment
