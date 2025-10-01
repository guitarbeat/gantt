# Task Stacking Fix - Multi-Day Overlap Detection

## Summary

Fixed the issue where tasks that span multiple days and start on different days were not stacking properly. Tasks now correctly stack vertically when they overlap, even if they start on different days.

## Problem

Previously, when two or more tasks overlapped across multiple days but started on different days, the later task would overlay or cut off the earlier task instead of stacking above it.

**Example Case**:
- **Task 1**: "Write Methods & Timeline" (Nov 11-18, 2025)
- **Task 2**: "Send Proposal to Committee" (Nov 12 - Dec 1, 2025)

On November 12, both tasks are active. Task 1 started on Nov 11, so it's already displayed. Task 2 starts on Nov 12, so it should stack ABOVE Task 1, but it was instead overlaying it.

## Root Cause

The original implementation tried to use absolute positioning with `\raisebox` offsets calculated based on stack position. This approach had several issues:

1. **Coordinate system confusion**: LaTeX's `\raisebox` raises content UPWARD for positive values, but we needed tasks to stack DOWNWARD in cells
2. **Fixed offsets**: The `2.0em * stackPos` offset didn't account for variable content heights
3. **Complex positioning**: Trying to manually position each task made the code fragile and hard to maintain

## Solution

Replaced absolute positioning with **natural LaTeX flow positioning**:

1. **Removed** the complex `stackPos` offset calculation
2. **Simplified** the rendering to use `\TaskOverlayBox` for all tasks (no offset variant)
3. **Added** `\vspace{1mm}` spacing between tasks
4. **Let LaTeX** handle vertical stacking naturally (earliest tasks first, latest tasks last)

### Key Changes in `src/calendar/calendar.go`

**Before** (lines 185-275):
```go
// Complex struct to track stack positions
var tasksToRender []struct {
    task         *SpanningTask
    stackPos     int
    isBottomTask bool
}

// Calculate offsets based on stack position
verticalOffset := fmt.Sprintf("%.1fem", float64(tr.stackPos)*2.0)

// Use offset-based positioning macro
pillContent := fmt.Sprintf(`\TaskOverlayBoxWithOffset{%s}{%s}{%s}{%s}`,
    taskColor, taskName, objective, verticalOffset)
```

**After**:
```go
// Simple list of tasks to render
var tasksToRender []*SpanningTask

// Add spacing between tasks (except first)
var spacing string
if i > 0 {
    spacing = `\vspace{1mm}`
}

// Use simple box without offsets
pillContent := spacing + fmt.Sprintf(`\TaskOverlayBox{%s}{%s}{%s}`,
    taskColor, taskName, objective)
```

## Benefits

1. **Simpler code**: Removed ~30 lines of complex offset calculation logic
2. **More robust**: Works correctly with variable content heights
3. **Better maintainability**: Natural LaTeX flow is easier to understand and modify
4. **Correct stacking**: Tasks now stack properly regardless of when they start

## Testing

Verified the fix works with:

✅ Week 45 (November 2025): "Write Methods & Timeline" and "Send Proposal to Committee" now stack correctly

✅ Week 5 (February 2027): "Write Methods & Results Chapters" and "Write Discussion & Conclusions" stack correctly with proper spacing

✅ All generated PDFs compile successfully with no LaTeX errors

## Visual Result

Tasks now display with proper vertical stacking:
- Earlier tasks appear at the bottom
- Later tasks appear above earlier tasks
- Consistent 1mm spacing between stacked tasks
- Cell height expands naturally to fit all tasks

## Files Modified

- `src/calendar/calendar.go` - Simplified `renderSpanningTaskOverlay()` function
- Created analysis document: `OVERLAP_DETECTION_ANALYSIS.md`
- Created this fix summary: `TASK_STACKING_FIX.md`
