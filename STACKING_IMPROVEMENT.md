# Multi-Day Task Stacking Improvement

## Issue Description

Previously, the task stacking system had a critical bug where tasks starting on different days but overlapping in time would visually overlap instead of stacking vertically. 

### Example Problem:
- **Task A**: "Write Methods" (Nov 11-18)
- **Task B**: "Send Proposal to Committee" (Nov 12 - Dec 1)

On November 12th:
- Task A is already active (started Nov 11)
- Task B starts on Nov 12

**Expected behavior**: Task B should stack **above** Task A to avoid visual overlap.

**Actual behavior (before fix)**: Task B was rendered at the same vertical position as Task A, causing them to overlap.

## Root Cause

The bug was in the `renderSpanningTaskOverlay()` function in `calendar.go`:

```go
// OLD CODE - BUGGY
for stackPos, task := range activeTasks {
    start := d.getTaskStartDate(task)
    if dayDate.Equal(start) {
        isBottom := (stackPos == 0) // BUG: Only true if first in list
        tasksToRender = append(tasksToRender, ...)
    }
}
```

### The Problem:
1. `activeTasks` correctly identified all tasks active on a given day (including those that started earlier)
2. Tasks were sorted by start date, so earlier tasks came first
3. BUT: Only tasks that **started today** were added to `tasksToRender`
4. The `isBottom` flag was based on `stackPos == 0`, which only worked if the first active task was also starting today

### Example of the Bug:
On Nov 12:
- `activeTasks[0]` = "Write Methods" (started Nov 11, **not rendered today**)
- `activeTasks[1]` = "Send Proposal" (starts Nov 12, **rendered today**)

When rendering "Send Proposal":
- Its `stackPos = 1`, so `isBottom = false`
- It used `TaskOverlayBoxNoOffset` (no vertical spacing)
- But since "Write Methods" wasn't rendered today, "Send Proposal" was the **only** task being rendered
- Result: Tasks overlapped visually because they were both at the same vertical position

## Solution

Changed the rendering approach to use **explicit vertical offsets** based on the `stackPos`:

```go
// NEW CODE - FIXED
for _, tr := range tasksToRender {
    // Use stackPos to determine vertical offset
    verticalOffset := fmt.Sprintf("%.1fem", float64(tr.stackPos)*2.0) // 2.0em per level
    
    // All tasks use offset-based positioning
    pillContent := fmt.Sprintf(`\TaskOverlayBoxWithOffset{%s}{%s}{%s}{%s}`,
        taskColor,
        taskName,
        objective,
        verticalOffset)
    pillContents = append(pillContents, pillContent)
}
```

### How It Works:
1. Each task gets a vertical offset based on its `stackPos` in the `activeTasks` list
2. `stackPos = 0` → offset = `0em` (bottom)
3. `stackPos = 1` → offset = `2.0em` (above the first task)
4. `stackPos = 2` → offset = `4.0em` (above the second task)
5. And so on...

This ensures that tasks are always positioned at the correct vertical level, regardless of whether earlier tasks are being rendered on the same day.

## Key Changes

### File: `src/calendar/calendar.go`

#### 1. Updated `renderSpanningTaskOverlay()` to track stack positions
```go
// Build rendering list: only tasks that START today
// However, we need to know the correct stack position based on ALL active tasks
var tasksToRender []struct {
    task         *SpanningTask
    stackPos     int
    isBottomTask bool
}

// Find which tasks actually start today and need to be rendered
var firstRenderIndex int = -1
for stackPos, task := range activeTasks {
    start := d.getTaskStartDate(task)
    if dayDate.Equal(start) {
        // This task starts today and should be rendered
        if firstRenderIndex == -1 {
            firstRenderIndex = stackPos
        }
        tasksToRender = append(tasksToRender, struct {
            task         *SpanningTask
            stackPos     int
            isBottomTask bool
        }{
            task:         task,
            stackPos:     stackPos,
            isBottomTask: false,
        })
    }
}
```

#### 2. Changed rendering to use `TaskOverlayBoxWithOffset`
```go
// Use stackPos to determine vertical offset
verticalOffset := fmt.Sprintf("%.1fem", float64(tr.stackPos)*2.0)

// All tasks use offset-based positioning
pillContent := fmt.Sprintf(`\TaskOverlayBoxWithOffset{%s}{%s}{%s}{%s}`,
    taskColor,
    taskName,
    objective,
    verticalOffset)
```

## Benefits

1. **Correct stacking**: Tasks that overlap in time are now always stacked vertically
2. **Consistent positioning**: Each task's vertical position is determined by its stack position, not by whether other tasks are being rendered
3. **Handles complex cases**: Works correctly when:
   - Multiple tasks overlap
   - Tasks start on different days
   - Tasks span multiple weeks
   - Tasks have different durations

## Testing

The fix was tested with the research timeline dataset which contains many overlapping tasks:
- PhD Proposal phase tasks that span multiple weeks
- Research execution tasks with complex dependencies
- Dissertation writing tasks that overlap with other phases

### Specific Test Case:
- **Week 45 (November 2025)**:
  - Nov 11: "Write Methods & Timeline" starts (8 days)
  - Nov 12: "Send Proposal to Committee" starts (20 days)
  - Result: Both tasks are now properly stacked without visual overlap

## Future Improvements

1. Make the vertical spacing (`2.0em`) configurable via the config file
2. Add automatic height adjustment for cells based on number of stacked tasks
3. Consider adding visual indicators (lines or shading) to show task continuation from previous days
4. Optimize spacing for tasks with very long descriptions

## Backward Compatibility

This change maintains backward compatibility:
- The `TaskOverlayBoxWithOffset` LaTeX command already existed
- No changes required to configuration files
- No changes required to CSV input format
- Existing PDFs will regenerate with the improved stacking
