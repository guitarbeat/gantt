# Multi-Day Task Stacking Fix - Summary

## Problem Statement
Tasks spanning multiple days with different start dates were visually overlapping instead of stacking vertically. Specifically:
- Task A: "Write Methods & Timeline" (Nov 11-18)
- Task B: "Send Proposal to Committee" (Nov 12 - Dec 1)

When rendered, Task B would draw over Task A on Nov 12-18, making both tasks hard to read.

## Root Cause
The calendar rendering system only considered tasks that **started** on each specific day. When Task B started on Nov 12, it didn't know that Task A (which started Nov 11) was still active, so it couldn't position itself properly above Task A.

## Solution
Enhanced the task detection algorithm to find ALL active tasks on a given day (both starting and continuing), then use this information to calculate proper vertical stack positions.

### Key Components:

1. **Enhanced Task Detection** (`findActiveTasks`)
   - Now returns ALL tasks active on a day, not just those that start
   - Includes a deduplication mechanism to avoid counting tasks twice
   - Calculates proper column spans for both starting and continuing tasks

2. **Stack Position Calculation**
   - Tasks are sorted by start date (earlier tasks are lower in the stack)
   - Each task knows its position: 0 = bottom, 1 = next level, etc.
   - Only tasks that START on a given day are rendered, but all active tasks contribute to positioning

3. **LaTeX Macro Selection**
   - Bottom task (stackPos=0): Uses `\TaskOverlayBox` with `\vfill` (positions at cell bottom)
   - Stacked tasks (stackPos>0): Use `\TaskOverlayBoxNoOffset` without `\vfill` (stack naturally)

## Visual Result
```
Nov 11:  [Write Methods ----------------]
Nov 12:  [Write Methods --------]
         [Send Proposal ------------------]
Nov 18:                 [Send Proposal ---]
```

Tasks properly stack vertically without overlap, each appearing once at its start position and spanning the appropriate number of days.

## Technical Details
- PDF size: 34 pages (efficient)
- LaTeX compilation: Successful
- Color preservation: Yes
- Task descriptions: Visible
- Spanning logic: Maintained

## Files Modified
1. `src/calendar/calendar.go` - Core stacking logic
2. `src/shared/templates/monthly/macros.tpl` - LaTeX macros
3. `TASK_SPANNING_FIX_LOG.md` - Development log

## Testing
To verify the fix:
1. Build: `./build.sh`
2. Open: `build/base.pdf`
3. Navigate to Week 45 (November 2025)
4. Verify: "Write Methods" and "Send Proposal" stack properly

