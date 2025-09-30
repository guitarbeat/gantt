# Task Stacking System - User Guide

## Overview

The task stacking system intelligently handles multi-day tasks that overlap, even when they start on different days. This prevents visual overlap and ensures all tasks are clearly visible.

## How It Works

### Before (Original System)
- **Problem**: Tasks only stacked if they started on the **exact same day**
- **Result**: Tasks starting on different days would visually overlap if their date ranges overlapped
- **Example**: 
  - Task A: Sept 1-10
  - Task B: Sept 3-12
  - These would OVERLAP because they start on different days

### After (New Stacking System)
- **Solution**: Analyzes ALL tasks across ALL days to assign vertical "tracks"
- **Result**: Tasks that overlap in time are placed in different tracks, preventing visual overlap
- **Example**:
  - Task A: Sept 1-10 → Track 0 (top)
  - Task B: Sept 3-12 → Track 1 (below Task A)
  - They stack properly even though they start on different days!

## Key Features

### 1. Global Overlap Detection
The system examines all tasks to determine which ones have overlapping date ranges, regardless of start date.

### 2. Smart Track Assignment
Each task is assigned to the lowest available vertical track (layer) that's free for its entire duration.

### 3. Clean Visualization
Tasks are only shown on their **starting day** with their full label, preventing clutter from continuation indicators.

### 4. Consistent Stacking
Tasks maintain their track assignment throughout, so the vertical spacing remains consistent.

## Algorithm Details

### Track Assignment Process

1. **Sort Tasks**: Order by start date, then by duration (longer tasks first)
2. **Find Available Track**: For each task, find the lowest track number that's free for ALL days it spans
3. **Assign Track**: Place task in that track
4. **Repeat**: Continue for all tasks

### Example Scenario

```
Tasks:
- Task A: Jan 1-7   (7 days)
- Task B: Jan 3-10  (8 days)  
- Task C: Jan 5-12  (8 days)
- Task D: Jan 8-10  (3 days)

Track Assignment:
- Track 0: Task A (Jan 1-7), then Task D (Jan 8-10) [reuses track after A ends]
- Track 1: Task B (Jan 3-10)
- Track 2: Task C (Jan 5-12)

Visual Layout:
Day 1: [A starts]
Day 2: [A .......]
Day 3: [A .......] [B starts]
Day 4: [A .......] [B .......] 
Day 5: [A .......] [B .......] [C starts]
Day 6: [A .......] [B .......] [C .......]
Day 7: [A .......] [B .......] [C .......]
Day 8:             [B .......] [C .......] [D starts]
Day 9:             [B .......] [C .......] [D .......]
Day 10:            [B .......] [C .......] [D .......]
```

## Technical Implementation

### Core Components

1. **TaskStacker** (`task_stacker.go`)
   - `ComputeStacks()`: Main algorithm
   - `GetTasksStartingOnDay(date)`: Returns tasks with their track assignments
   - `GetStacksForDay(date)`: Returns all active tasks on a date

2. **Day Rendering** (`calendar.go`)
   - `renderLargeDayWithStacking()`: Uses the new stacker
   - `renderSpanningTaskOverlayWithStacking()`: Generates stacked task pills

3. **LaTeX Macros** (`macros.tpl`)
   - `\TaskOverlayBox`: First task in stack (with top margin)
   - `\TaskOverlayBoxNoOffset`: Subsequent tasks (touching previous)

### Data Structures

```go
type TaskStack struct {
    Track    int           // Vertical position (0=top, 1=second, etc.)
    Task     *SpanningTask // The actual task
    StartCol int           // Starting column in week grid
    EndCol   int           // Ending column in week grid
}
```

## Configuration

Currently uses default settings. Future enhancements will add:

```yaml
task_stacking:
  enabled: true
  max_tracks_per_day: 10
  min_task_height: 8pt
  stack_spacing: 1pt
```

## Testing

Run tests with:
```bash
go test -v ./src/calendar/task_stacker_test.go ./src/calendar/task_stacker.go
```

Test scenarios include:
- Two tasks with complete overlap
- Three tasks with partial overlap  
- Non-overlapping tasks (should reuse tracks)
- Single-day tasks
- Tasks crossing week boundaries

## Troubleshooting

### Tasks Still Overlapping?

1. **Rebuild**: Ensure you've rebuilt after the code changes
   ```bash
   make -f scripts/Makefile clean-build
   ```

2. **Check Stacking Method**: Verify `renderLargeDayWithStacking()` is being called
   ```bash
   grep "renderLargeDayWithStacking" src/calendar/calendar.go
   ```

3. **Verify Data**: Check that tasks have correct start/end dates in CSV

### Too Much Vertical Space?

The algorithm uses a greedy approach that minimizes tracks. If you see excessive vertical space, check for:
- Very long-running tasks forcing other tasks into higher tracks
- Many simultaneous tasks during busy periods

### Tasks Not Showing?

Verify tasks are:
- Loaded from CSV correctly
- Have valid date ranges
- Fall within the displayed month/year

## Future Enhancements

### Planned Features

1. **Continuation Indicators** (Optional)
   - Thin colored bars showing task duration
   - Configurable on/off

2. **Visual Task Spanning**
   - Tasks visually span across multiple days
   - Uses LaTeX multicolumn or TikZ overlays

3. **Compact Mode**
   - Auto-adjusts cell height based on number of tracks
   - Prevents excessive white space

4. **Priority Stacking**
   - Important tasks get lower track numbers (more prominent)
   - Based on urgency or milestone status

5. **Interactive PDF**
   - Click task to see full details
   - Hover to highlight entire task duration

## Performance Considerations

- **O(n²)** worst case for track assignment (n = number of tasks)
- **Optimized** for typical workloads (< 100 tasks per month)
- **Memory efficient**: Only stores task references, not copies

## References

- Source code: `src/calendar/task_stacker.go`
- Tests: `src/calendar/task_stacker_test.go`
- Rendering: `src/calendar/calendar.go`
- LaTeX macros: `src/shared/templates/monthly/macros.tpl`

## Contributing

To improve the stacking algorithm:

1. Add test cases to `task_stacker_test.go`
2. Implement improvement in `task_stacker.go`
3. Run tests to verify correctness
4. Rebuild and visually verify PDF output

## Questions?

See `docs/TASK_STACKING_NOTES.md` for implementation details and design decisions.
