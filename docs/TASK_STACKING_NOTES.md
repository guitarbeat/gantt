# Task Stacking Improvements - Implementation Notes

## Problem Statement
The original implementation only stacked tasks that started on the **same day**. Tasks that started on different days but had overlapping date ranges would visually overlap because they weren't aware of each other.

## Solution Approach

### Key Insight
To properly stack multi-day tasks, we need to consider ALL tasks that are **active** on any given day, not just tasks that **start** on that day.

### Algorithm Overview

1. **Global Track Assignment**: Analyze all tasks across all days to assign each task to a vertical "track" (layer)
2. **Overlap Detection**: Check if tasks overlap by comparing their date ranges
3. **Smart Rendering**: 
   - Tasks that START on a day show their full label
   - Tasks continuing from previous days show thin continuation bars
   - All tasks maintain their track assignment for visual consistency

### Implementation Details

#### TaskStacker Class
- `ComputeStacks()`: Main algorithm that assigns tracks to all tasks
- `GetStacksForDay(date)`: Returns all tasks visible on a specific day
- `GetTasksStartingOnDay(date)`: Returns only tasks starting on that day
- `findLowestAvailableTrack(task)`: Finds the first track that's free for ALL days the task spans

#### Rendering Strategy
```
Day 1: [Task A starts] ─────────►
Day 2: [continuation] [Task B starts] ────►
Day 3: [continuation] [continuation] [Task C starts] ─►
Day 4: [continuation] [continuation] [continuation]
```

Each task stays in its assigned track throughout its duration.

## Current Issues to Address

### Issue 1: Too Many Continuation Bars
**Problem**: The current implementation shows a continuation bar for EVERY task on EVERY day it's active, which creates visual clutter.

**Better Approach**: Only show the task label on the starting day, and let the colored pill span across days visually (using LaTeX multicolumn or TikZ).

### Issue 2: Track Assignment Conflicts
**Problem**: The track assignment might not be optimal, leading to unnecessary vertical space.

**Solution**: Use a greedy algorithm that assigns tasks to the lowest available track.

### Issue 3: Week Boundaries
**Problem**: Tasks that span across week boundaries need special handling in the calendar grid.

**Solution**: Calculate column positions within each week separately.

## Recommended Simplifications

### Option A: Show Only Starting Day (Current Approach)
- ✅ Simple to implement
- ✅ Clear which day a task starts
- ❌ Doesn't show task continuation visually
- ❌ Hard to see multi-day tasks at a glance

### Option B: Show Continuation Bars (Implemented)
- ✅ Shows task duration visually
- ✅ Maintains track consistency
- ❌ Can be visually cluttered with many tasks
- ❌ Continuation bars might be too thin

### Option C: Colored Background Blocks (Recommended)
- ✅ Most intuitive visual representation
- ✅ Clear task duration
- ✅ Clean appearance
- ❌ More complex LaTeX implementation
- ❌ Requires TikZ overlays or table background coloring

## Next Steps

1. **Gather User Feedback**: Understand exactly what's "wrong" with the current output
2. **Adjust Visualization**: Based on feedback, adjust the rendering approach
3. **Optimize Track Assignment**: Improve the algorithm to minimize tracks used
4. **Add Configuration**: Allow users to choose visualization style

## Configuration Options to Add

```yaml
task_stacking:
  enabled: true
  show_continuation_bars: true  # or false for clean look
  continuation_bar_height: 3pt
  min_task_height: 8pt
  stack_spacing: 1pt
  max_tracks_per_day: 5
```

## Testing Recommendations

Test with these scenarios:
1. Two tasks with complete overlap (same start and end dates)
2. Three tasks with partial overlap (staggered starts)
3. Long task (10+ days) with multiple shorter tasks starting during it
4. Tasks that cross week boundaries
5. Month with many simultaneous tasks (stress test)
