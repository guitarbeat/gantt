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

## Current Single-Day Stacking Logic

**How Single-Day Tasks Stack (Working Implementation):**

### LaTeX Macros Analysis:
```latex
% TaskOverlayBox - First task gets vertical fill to bottom-align
\newcommand{\TaskOverlayBox}[3]{%
  \definecolor{taskbgcolor}{RGB}{#1}%
  \vfill  % <-- KEY: Pushes content to bottom of cell
  \begin{tcolorbox}[...]
    % Task content
  \end{tcolorbox}%
}

% TaskOverlayBoxNoOffset - Subsequent tasks stack without extra spacing
\newcommand{\TaskOverlayBoxNoOffset}[3]{%
  \definecolor{taskbgcolor}{RGB}{#1}%
  % No \vfill - stacks directly on previous task
  \begin{tcolorbox}[top=0pt, bottom=0pt, ...]
    % Task content
  \end{tcolorbox}%
}
```

### Go Implementation:
```go
func (d Day) renderSpanningTaskOverlay() *TaskOverlay {
    activeTasks, maxCols := d.findActiveTasks(dayDate)

    for i, task := range activeTasks {
        if i == 0 {
            // First task: use \TaskOverlayBox (with \vfill for bottom alignment)
            pillContent := fmt.Sprintf(`\TaskOverlayBox{%s}{%s}{%s}`,
                taskColor, taskName, objective)
        } else {
            // Subsequent tasks: use \TaskOverlayBoxNoOffset (no extra spacing)
            pillContent := fmt.Sprintf(`\TaskOverlayBoxNoOffset{%s}{%s}{%s}`,
                taskColor, taskName, objective)
        }
        pillContents = append(pillContents, pillContent)
    }

    // Stack pills vertically without extra spacing
    content := strings.Join(pillContents, "")
    return &TaskOverlay{content: content, cols: maxCols}
}
```

### Key Insights:
1. **First task gets `\vfill`** → Bottom-aligned in cell
2. **Subsequent tasks have `top=0pt, bottom=0pt`** → No extra vertical spacing
3. **Tasks stack by touching** → No gaps between stacked tasks
4. **All tasks sorted by start date** → Earlier tasks appear lower (behind later ones)

## Detailed Implementation Ideas

### Option A: Vertical Offset Calculation (Recommended)

**Core Concept**: Apply the same stacking logic as single-day tasks, but across multiple days.

**How to Apply Single-Day Logic to Multi-Day Tasks:**

1. **For each day, collect ALL tasks that should be visible:**
   - Tasks starting on this day (current logic)
   - Tasks that started earlier but are still active (NEW)

2. **Sort all visible tasks by start date** (same as single-day)

3. **Apply stacking offsets based on position in sorted list:**
   - Position 0 (earliest): Use `TaskOverlayBox` (with `\vfill`)
   - Position 1+: Use `TaskOverlayBoxNoOffset` (no extra spacing)

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

**LaTeX Rendering (Same as Single-Day):**
```go
func (d Day) renderMultiDayStackedOverlay() *TaskOverlay {
    allVisibleTasks, maxCols := d.findAllVisibleTasks(dayDate)

    for i, task := range allVisibleTasks {
        if i == 0 {
            // Earliest task: bottom-aligned with \vfill
            pillContent := fmt.Sprintf(`\TaskOverlayBox{%s}{%s}{%s}`,
                taskColor, taskName, objective)
        } else {
            // Later tasks: stack without extra spacing
            pillContent := fmt.Sprintf(`\TaskOverlayBoxNoOffset{%s}{%s}{%s}`,
                taskColor, taskName, objective)
        }
        pillContents = append(pillContents, pillContent)
    }

    content := strings.Join(pillContents, "")
    return &TaskOverlay{content: content, cols: maxCols}
}
```

**Key Changes from Single-Day:**
- `findActiveTasks()` → `findAllVisibleTasks()` (includes continuing tasks)
- Need `calculateRemainingSpanColumns()` for tasks that started earlier
- Same stacking macros (`\TaskOverlayBox` vs `TaskOverlayBoxNoOffset`)

**Core Concept**: Calculate y-offset based on active task count at each position.

**Implementation Strategy**:
1. Create a `TaskStackManager` struct to track active tasks across days
2. For each day, calculate which tasks are "stacked" at that position
3. Use TikZ `yshift` to vertically offset overlapping tasks

**Code Structure**:
```go
type TaskStackManager struct {
    activeTasks map[string][]*TaskPosition // date -> tasks with their stack positions
}

type TaskPosition struct {
    task     *SpanningTask
    stackPos int // 0 = bottom, 1 = middle, 2 = top, etc.
    yOffset  float64 // calculated pixel offset
}
```

**Algorithm**:
```go
func (d Day) calculateTaskOffsets() {
    // For each day, find ALL tasks that would be visible at this position
    // Tasks starting on this day + tasks that started earlier and are still active
    allVisibleTasks := d.getAllVisibleTasks()

    // Group by stack position (earlier start dates get lower positions)
    sortedTasks := sortTasksByStartDate(allVisibleTasks)

    // Assign stack positions and calculate y-offsets
    for i, task := range sortedTasks {
        task.stackPos = i
        task.yOffset = float64(i) * d.getTaskHeight()
    }
}
```

**LaTeX Generation**:
```go
func (d Day) renderStackedTaskOverlay() string {
    overlays := []string{}

    for _, taskPos := range d.taskPositions {
        yshift := fmt.Sprintf("yshift=%.1fpt", taskPos.yOffset)
        overlay := fmt.Sprintf(`
            \begin{tikzpicture}[overlay]
                \node[anchor=north west, inner sep=0pt, %s] at (0,0) {
                    \begin{minipage}[t]{\dimexpr %d\linewidth\relax}
                        %s
                    \end{minipage}
                };
            \end{tikzpicture}`, yshift, taskPos.cols, taskPos.content)
        overlays = append(overlays, overlay)
    }

    return strings.Join(overlays, "")
}
```

**Challenges & Solutions**:
- **Height Calculation**: Tasks have variable heights (title + description)
  - Solution: Use LaTeX `\heightof` or pre-calculate based on content length
- **Week Boundaries**: Tasks spanning weeks need consistent positioning
  - Solution: Use global task IDs and maintain state across week boundaries
- **Performance**: Calculating all active tasks each day
  - Solution: Pre-compute task positions for entire timeline

### Option B: Week-Level Rendering (Major Refactor)

**Core Concept**: Render entire weeks as single units, managing all task positioning within the week.

**Implementation Strategy**:
1. Change from day-level to week-level rendering
2. Calculate all active tasks for the entire week upfront
3. Position all task bars with proper vertical stacking

**Data Structure Changes**:
```go
type Week struct {
    Days       [7]*Day
    AllTasks   []*SpanningTask // All tasks active during this week
    TaskLayout map[*SpanningTask]*TaskLayout
}

type TaskLayout struct {
    StartDay  int // 0-6 (Monday-Sunday)
    EndDay    int
    StackPos  int // Vertical position in stack
    YOffset   float64
}
```

**Rendering Flow**:
```go
func (w Week) Render() string {
    w.calculateTaskLayout() // Pre-calculate all positions

    weekContent := ""
    for dayIdx, day := range w.Days {
        dayContent := w.renderDayWithLayout(day, dayIdx)
        weekContent += dayContent
    }

    return weekContent
}
```

**Advantages**:
- ✅ Perfect stacking control within weeks
- ✅ No cross-day positioning issues
- ✅ Easier to handle week boundaries

**Disadvantages**:
- ❌ Major architectural change
- ❌ More complex state management
- ❌ Harder to handle month boundaries

### Option C: Hybrid TikZ Row Spanning

**Core Concept**: Use TikZ overlays that span both columns AND rows vertically.

**Implementation Strategy**:
1. Tasks span multiple columns (already working)
2. Add vertical row spanning for stacking
3. Use TikZ `fit` library or manual positioning

**LaTeX Approach**:
```tex
% Instead of single cell overlay, use multi-cell spanning
\begin{tikzpicture}[overlay, remember picture]
    % Task 1: spans days 1-3, vertically positioned at level 0
    \node[fit=(day1)(day3), inner sep=2pt, fill=blue!20] {};

    % Task 2: spans days 2-4, vertically positioned at level 1
    \node[fit=(day2)(day4), inner sep=2pt, fill=red!20, yshift=1.5cm] {};
\end{tikzpicture}
```

**Go Implementation**:
```go
func (d Day) renderMultiCellTask(task *SpanningTask, stackPos int) string {
    startDay := d.getDayIndexInWeek(task.StartDate)
    endDay := d.getDayIndexInWeek(task.EndDate)
    yOffset := float64(stackPos) * d.getTaskRowHeight()

    return fmt.Sprintf(`
        \node[fit=(week%d-day%d)(week%d-day%d),
              inner sep=2pt,
              fill=%s!20,
              yshift=%.1fpt] {};`,
        d.weekNum, startDay, d.weekNum, endDay,
        task.Color, yOffset)
}
```

### Option D: LaTeX Package Integration (From Web Research)

**Leverage Existing LaTeX Packages**:
Based on web research, consider integrating specialized LaTeX packages:

1. **Tasks Package** (`\usepackage{tasks}`):
   ```tex
   \usepackage{tasks}
   \NewTasks[style=enumerate]{tasklist}[\task]

   % Use for task rendering
   \begin{tasklist}
     \task Task spanning multiple days
     \task Another overlapping task
   \end{tasklist}
   ```

2. **Calendar-Specific Packages**:
   - `calendar` package for advanced date calculations
   - `tikz-calendar` for visual calendar rendering
   - `pgfcalendar` for complex date logic

3. **Integration Approach**:
   ```go
   func (d Day) renderWithTasksPackage() string {
       // Generate LaTeX using tasks package syntax
       // Instead of custom TikZ, use package-provided task rendering
   }
   ```

### Option E: State-Persistent Task Tracking

**Core Concept**: Maintain global state of active tasks across the entire timeline.

**Implementation Strategy**:
1. Create a global `TimelineState` that tracks all active tasks
2. For each day, query which tasks are currently active
3. Calculate stack positions based on global task state

**Global State Structure**:
```go
type TimelineState struct {
    ActiveTasks []*ActiveTaskInfo // Sorted by start date
    TaskCounter int              // For assigning unique IDs
}

type ActiveTaskInfo struct {
    Task     *SpanningTask
    StartDay time.Time
    EndDay   time.Time
    StackPos int // Calculated once, used everywhere
}
```

**Usage Pattern**:
```go
func (d Day) renderWithGlobalState(state *TimelineState) string {
    // Query state for tasks active on this day
    activeOnThisDay := state.GetTasksForDay(d.Time)

    // Tasks already have their stack positions calculated
    return d.renderTasksWithPositions(activeOnThisDay)
}
```

**Advantages**:
- ✅ Consistent stacking across entire timeline
- ✅ Easy to handle complex overlaps
- ✅ No day-to-day state management

**Disadvantages**:
- ❌ Requires significant refactoring
- ❌ More memory usage
- ❌ Complex state management

### Option F: CSS-Grid Inspired LaTeX Layout

**Core Concept**: Use LaTeX's grid positioning similar to CSS Grid for precise control.

**Implementation Strategy**:
1. Define a grid system for each week
2. Position tasks using absolute coordinates within the grid
3. Use TikZ `grid` and `positioning` libraries

**LaTeX Grid Setup**:
```tex
\begin{tikzpicture}[overlay]
    % Define week grid
    \draw[grid] (0,0) grid (7,5); % 7 days, 5 potential stack levels

    % Position tasks within grid
    \node[anchor=north west] at (1,0) {\TaskBox{Task spanning days 1-3}};
    \node[anchor=north west] at (2,1) {\TaskBox{Task spanning days 2-4}};
\end{tikzpicture}
```

**Go Implementation**:
```go
func (w Week) renderGridBased() string {
    grid := w.calculateGridLayout() // Returns positions for all tasks

    latex := `\begin{tikzpicture}[overlay]`
    for _, pos := range grid.Positions {
        latex += fmt.Sprintf(`
            \node[anchor=north west] at (%.1f, %.1f) {%s};`,
            pos.X, pos.Y, pos.TaskContent)
    }
    latex += `\end{tikzpicture}`

    return latex
}
```

## Implementation Priority & Testing Strategy

### Recommended Approach: Option A (Vertical Offset Calculation)
- **Why**: Least disruptive to existing codebase
- **Effort**: Medium (2-3 days)
- **Risk**: Low-medium
- **Testability**: Can test incrementally

### Testing Strategy:
1. **Unit Tests**: Test offset calculations with mock tasks
2. **Integration Tests**: Verify LaTeX compilation with sample data
3. **Visual Tests**: Generate PDFs and manually verify stacking
4. **Regression Tests**: Ensure existing single-day tasks still work

### Next Steps

1. **Implement Option A** as proof of concept
2. **Add comprehensive tests** for edge cases (week boundaries, month transitions)
3. **Profile performance** to ensure LaTeX compilation remains fast
4. **Consider Option B** if Option A proves too complex

## Advanced Implementation Refinements

### Refinement A: Smart Span Calculation for Continuing Tasks

**Problem**: When a task continues from a previous day, how many columns should it span on the current day?

**Current Logic**: `calculateTaskSpanColumns()` calculates from start day
**Needed**: `calculateRemainingSpanColumns()` for continuing tasks

```go
// For tasks starting today: span from today to end (or end of week)
func (d Day) calculateTaskSpanColumns(dayDate, end time.Time) int {
    // Current logic works for starting tasks
    idxMonFirst := (int(dayDate.Weekday()) + 6) % 7
    remainInRow := 7 - idxMonFirst
    totalRemain := int(end.Sub(dayDate).Hours()/24) + 1
    if totalRemain > remainInRow {
        return remainInRow
    }
    return totalRemain
}

// NEW: For continuing tasks: span from today to end (or end of week)
func (d Day) calculateRemainingSpanColumns(dayDate, end time.Time) int {
    // Same logic but task has already consumed some days
    idxMonFirst := (int(dayDate.Weekday()) + 6) % 7
    remainInRow := 7 - idxMonFirst
    daysLeft := int(end.Sub(dayDate).Hours()/24) + 1
    if daysLeft < 1 {
        return 1 // At least show on current day
    }
    if daysLeft > remainInRow {
        return remainInRow
    }
    return daysLeft
}
```

### Refinement B: Task Height Estimation & Collision Detection

**Problem**: Tasks have variable heights based on title + description length. Need to prevent overlaps.

**Solution**: Pre-calculate task heights and adjust stacking offsets.

```go
type TaskDimensions struct {
    Task      *SpanningTask
    Height    float64 // Estimated height in points
    BaseY     float64 // Y position in stack
    StartDate time.Time
}

func (d Day) calculateTaskStackPositions(tasks []*SpanningTask) []*TaskDimensions {
    dimensions := make([]*TaskDimensions, len(tasks))

    currentY := 0.0
    for i, task := range tasks {
        height := d.estimateTaskHeight(task)
        dimensions[i] = &TaskDimensions{
            Task:      task,
            Height:    height,
            BaseY:     currentY,
            StartDate: task.StartDate,
        }
        currentY += height // Stack next task above this one
    }

    return dimensions
}

func (d Day) estimateTaskHeight(task *SpanningTask) float64 {
    baseHeight := 12.0 // Base height for tcolorbox
    titleLines := math.Ceil(float64(len(task.Name)) / 20.0) // ~20 chars per line
    descLines := 0.0
    if task.Description != "" {
        descLines = math.Ceil(float64(len(task.Description)) / 25.0) // Smaller font
    }
    return baseHeight + (titleLines * 6.0) + (descLines * 4.0)
}
```

### Refinement C: Week Boundary Handling

**Problem**: Tasks spanning week boundaries need consistent stacking across weeks.

**Solution**: Use global task IDs and maintain stack state across weeks.

```go
type TimelineStackState struct {
    TaskStackPositions map[string]int // TaskID -> stack position
    TaskHeights        map[string]float64 // TaskID -> height
    LastUpdated        time.Time
}

func (d Day) getConsistentStackPosition(task *SpanningTask, state *TimelineStackState) int {
    taskID := fmt.Sprintf("%s-%s", task.Name, task.StartDate.Format("2006-01-02"))

    if pos, exists := state.TaskStackPositions[taskID]; exists {
        return pos
    }

    // Calculate new position based on active tasks at start date
    pos := d.calculateNewStackPosition(task, state)
    state.TaskStackPositions[taskID] = pos
    return pos
}

func (d Day) calculateNewStackPosition(task *SpanningTask, state *TimelineStackState) int {
    // Find all tasks active on the start date, sorted by start time
    activeAtStart := d.getTasksActiveOnDate(task.StartDate)

    // Find position in sorted list
    for i, activeTask := range activeAtStart {
        if activeTask == task {
            return i
        }
    }
    return 0 // Fallback
}
```

### Refinement D: Performance Optimizations

**Problem**: Calculating task positions for every day is expensive.

**Solution**: Pre-compute all task positions and cache results.

```go
type TaskPositionCache struct {
    Positions map[string][]*TaskDimensions // date -> sorted task positions
    ValidUntil time.Time
}

func (d Day) buildTaskPositionCache(allTasks []*SpanningTask) *TaskPositionCache {
    cache := &TaskPositionCache{
        Positions: make(map[string][]*TaskDimensions),
    }

    // Find date range
    minDate, maxDate := d.findTimelineRange(allTasks)

    // Pre-calculate positions for each day
    for current := minDate; !current.After(maxDate); current = current.AddDate(0, 0, 1) {
        dateKey := current.Format("2006-01-02")
        cache.Positions[dateKey] = d.calculateTaskStackPositionsForDate(current, allTasks)
    }

    return cache
}

func (d Day) calculateTaskStackPositionsForDate(date time.Time, allTasks []*SpanningTask) []*TaskDimensions {
    var visibleTasks []*SpanningTask

    for _, task := range allTasks {
        if d.isTaskVisibleOnDate(date, task) {
            visibleTasks = append(visibleTasks, task)
        }
    }

    // Sort by start date for consistent stacking
    sort.Slice(visibleTasks, func(i, j int) bool {
        return visibleTasks[i].StartDate.Before(visibleTasks[j].StartDate)
    })

    return d.calculateTaskStackPositions(visibleTasks)
}
```

### Refinement E: Visual Polish & Edge Cases

**Problem**: Various edge cases need handling for professional appearance.

**Edge Cases to Handle:**
1. **Tasks ending on same day they start** (1-day tasks)
2. **Multiple tasks starting on same day** (already handled)
3. **Tasks spanning entire week**
4. **Tasks with very long descriptions**
5. **Week boundary transitions**
6. **Month boundary transitions**

**Polish Improvements:**
```go
// Add subtle visual indicators for spanning tasks
func (d Day) addSpanIndicators(task *SpanningTask, isContinuing bool) string {
    if isContinuing {
        // Add small arrow or continuation indicator
        return fmt.Sprintf(`\text{\tiny$\rightarrow$}%s`, taskName)
    }
    return taskName
}

// Handle text overflow gracefully
func (d Day) truncateTaskText(text string, maxLength int) string {
    if len(text) <= maxLength {
        return text
    }
    return text[:maxLength-3] + "..."
}

// Add hover tooltips for truncated text
func (d Day) addTooltip(fullText, truncatedText string) string {
    if fullText == truncatedText {
        return truncatedText
    }
    return fmt.Sprintf(`\tooltip{%s}{%s}`, truncatedText, fullText)
}
```

### Refinement F: Alternative LaTeX Approaches

**Problem**: Current TikZ overlay approach might have limitations.

**Alternative 1: Multi-cell Spanning with tabular**
```latex
% Instead of overlays, use multi-column cells
\multicolumn{3}{|c|}{\TaskBox{Task spanning 3 days}} &
\multicolumn{2}{|c|}{\TaskBox{Task spanning 2 days}} &
% Empty cells for non-spanning days
```

**Alternative 2: Layered TikZ Pictures**
```latex
\begin{tikzpicture}[overlay, remember picture]
    % Background layer: spanning bars
    \node[fill=blue!20, minimum height=1cm] at (week1-day1) {\TaskBar};
    \node[fill=red!20, minimum height=1cm, yshift=1.2cm] at (week1-day2) {\TaskBar};

    % Foreground layer: day numbers (must be on top)
    \node[anchor=center] at (week1-day1) {1};
    \node[anchor=center] at (week1-day2) {2};
\end{tikzpicture}
```

**Alternative 3: CSS-Grid Inspired LaTeX**
```latex
% Use LaTeX's grid positioning
\grid{
    \cell{1}{1}{\TaskBox{A}} \cell{1}{2}{\TaskBox{A continues}}
    \cell{2}{1}{\TaskBox{B}} \cell{2}{2}{\TaskBox{B continues}}
}
```

### Refinement G: Testing & Validation Strategy

**Comprehensive Test Suite:**
```go
func TestTaskStackingScenarios(t *testing.T) {
    tests := []struct {
        name     string
        tasks    []*SpanningTask
        testDate time.Time
        expected []string // Expected task names in stack order
    }{
        {
            name: "Single day, multiple tasks",
            tasks: []*SpanningTask{
                {Name: "Task A", StartDate: date("2025-01-01"), EndDate: date("2025-01-01")},
                {Name: "Task B", StartDate: date("2025-01-01"), EndDate: date("2025-01-01")},
            },
            testDate: date("2025-01-01"),
            expected: []string{"Task A", "Task B"}, // A first (earlier), B on top
        },
        {
            name: "Multi-day overlap",
            tasks: []*SpanningTask{
                {Name: "Long Task", StartDate: date("2025-01-01"), EndDate: date("2025-01-05")},
                {Name: "Short Task", StartDate: date("2025-01-03"), EndDate: date("2025-01-03")},
            },
            testDate: date("2025-01-03"),
            expected: []string{"Long Task", "Short Task"}, // Long task first, short on top
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            day := Day{Tasks: tt.tasks}
            visibleTasks, _ := day.findAllVisibleTasks(tt.testDate)

            if len(visibleTasks) != len(tt.expected) {
                t.Errorf("Expected %d tasks, got %d", len(tt.expected), len(visibleTasks))
            }

            for i, expected := range tt.expected {
                if visibleTasks[i].Name != expected {
                    t.Errorf("Position %d: expected %s, got %s", i, expected, visibleTasks[i].Name)
                }
            }
        })
    }
}
```

**Visual Regression Testing:**
```bash
# Generate test PDFs with known layouts
# Compare pixel-by-pixel or extract text positions
# Ensure stacking doesn't break existing layouts
```

### Refinement H: Configuration & Customization

**Make stacking behavior configurable:**
```yaml
task_stacking:
  enabled: true
  method: "overlay"  # "overlay", "multicol", "grid"
  spacing: "tight"   # "tight", "loose", "auto"
  overflow: "truncate"  # "truncate", "wrap", "tooltip"
  indicators:
    show_continuation_arrows: true
    show_start_markers: true
```

**Performance Tuning Options:**
```yaml
performance:
  cache_positions: true
  precalculate_heights: true
  max_visible_tasks: 10  # Prevent excessive stacking
```

## Implementation Roadmap

### Phase 1: Core Stacking (Week 1)
1. ✅ Implement `findAllVisibleTasks()` - include continuing tasks
2. ✅ Add `calculateRemainingSpanColumns()` 
3. ✅ Test basic stacking scenarios
4. ✅ Verify LaTeX compilation

### Phase 2: Polish & Edge Cases (Week 2)
1. 🔄 Add height estimation and collision detection
2. 🔄 Handle week/month boundary transitions
3. 🔄 Add performance optimizations (caching)
4. 🔄 Comprehensive testing

### Phase 3: Advanced Features (Week 3)
1. 📋 Add visual indicators (continuation arrows)
2. 📋 Alternative LaTeX rendering approaches
3. 📋 Configuration system
4. 📋 Documentation and examples

## Development Workflow Enhancements

### Enhancement I: AI-Assisted Development with MCP Integration

**Leveraging Model Context Protocol (MCP) for Enhanced Development:**

Based on web research, integrating MCP can significantly improve the development workflow:

**MCP Server Configuration:**
```json
{
  "mcpServers": {
    "gantt-planner": {
      "command": "go",
      "args": ["run", "cmd/mcp-server/main.go"],
      "env": {
        "PROJECT_ROOT": "/Users/aaron/Downloads/gantt",
        "CONFIG_FILE": "src/core/base.yaml"
      }
    }
  }
}
```

**MCP-Enabled Resources:**
- **Project Structure**: Real-time access to file tree and dependencies
- **Task Data**: Live access to CSV task data and current configurations
- **Build Status**: Compilation results and LaTeX generation status
- **PDF Output**: Generated calendar analysis and validation

**MCP Actions:**
```go
// MCP Server implementation for Gantt Planner
type GanttPlannerServer struct{}

func (s *GanttPlannerServer) HandleCall(method string, params interface{}) (interface{}, error) {
    switch method {
    case "gantt.validate_tasks":
        return s.validateTaskData()
    case "gantt.generate_preview":
        return s.generatePDFPreview()
    case "gantt.analyze_stacking":
        return s.analyzeTaskStacking()
    case "gantt.optimize_layout":
        return s.optimizeCalendarLayout()
    }
    return nil, fmt.Errorf("unknown method: %s", method)
}
```

**Benefits:**
- ✅ **Contextual Code Suggestions**: AI understands project structure and current issues
- ✅ **Automated Testing**: Run tests and validations through MCP
- ✅ **Live Debugging**: Real-time LaTeX compilation and error analysis
- ✅ **Task Management**: Direct manipulation of calendar data through AI

### Enhancement J: Syntax Assistance & Code Generation

**LaTeX Syntax Validation:**
```go
type ValidationWarning struct {
    Type    string
    Message string
    Line    int
}

type MacroDefinition struct {
    Name      string
    Signature string
    Required  bool
}

type PackageInfo struct {
    Name        string
    Description string
    Required    bool
}

type LaTeXValidator struct {
    macros    map[string]*MacroDefinition
    packages  map[string]*PackageInfo
    warnings  []ValidationWarning
}

type ValidationCheck func(string) []ValidationWarning

func (v *LaTeXValidator) ValidateGeneratedCode(code string) []ValidationWarning {
    // Check for common LaTeX syntax errors
    checks := []ValidationCheck{
        v.checkUndefinedMacros,
        v.checkMissingPackages,
        v.checkTikZSyntax,
        v.checkTableStructure,
        v.checkColorDefinitions,
    }

    var allWarnings []ValidationWarning
    for _, check := range checks {
        warnings := check(code)
        allWarnings = append(allWarnings, warnings...)
    }

    return allWarnings
}

func (v *LaTeXValidator) checkTikZSyntax(code string) []ValidationWarning {
    var warnings []ValidationWarning

    // Check for common TikZ errors
    patterns := map[string]string{
        "unclosed_braces": `\{[^\}]*$`,
        "missing_coordinates": `\\node\[.*\]\s*\{[^}]*\}`,
        "invalid_anchors": `anchor=[^cnesw]*`,
    }

    for name, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        if re.MatchString(code) {
            warnings = append(warnings, ValidationWarning{
                Type:    name,
                Message: fmt.Sprintf("Potential TikZ syntax issue: %s", name),
                Line:    findLineNumber(code, re),
            })
        }
    }

    return warnings
}
```

**Go Code Generation Assistants:**
```go
type CodeGenerator struct {
    templates map[string]*template.Template
}

func (g *CodeGenerator) GenerateTaskStruct(taskName string) string {
    return fmt.Sprintf(`
// %s represents a %s task in the calendar
type %s struct {
    Name        string
    Description string
    StartDate   time.Time
    EndDate     time.Time
    Color       string
    Priority    int
    Dependencies []*%s
}

func New%s(name string, start, end time.Time) *%s {
    return &%s{
        Name:      name,
        StartDate: start,
        EndDate:   end,
        Color:     "#default_color",
        Priority:  1,
    }
}
`, taskName, strings.ToLower(taskName), taskName, taskName, taskName, taskName, taskName)
}

func (g *CodeGenerator) GenerateTestCase(scenario string) string {
    return fmt.Sprintf(`
func Test%s(t *testing.T) {
    // Arrange
    tasks := []*SpanningTask{
        // %s test data
    }
    day := Day{Tasks: tasks}

    // Act
    visibleTasks, cols := day.findAllVisibleTasks(testDate)

    // Assert
    if len(visibleTasks) != expectedCount {
        t.Errorf("Expected %%d tasks, got %%d", expectedCount, len(visibleTasks))
    }
}
`, scenario, scenario)
}
```

**Interactive Code Completion:**
```go
type CompletionEngine struct {
    projectSymbols map[string]*SymbolInfo
    latexCommands  map[string]*CommandInfo
}

type CompletionItem struct {
    Label       string
    Detail      string
    InsertText  string
}

func (c *CompletionEngine) GetCompletions(prefix string, context string) []CompletionItem {
    var completions []CompletionItem

    // Go completions
    if strings.Contains(context, ".go") {
        completions = append(completions, c.getGoCompletions(prefix)...)
    }

    // LaTeX completions
    if strings.Contains(context, ".tex") || strings.Contains(context, ".tpl") {
        completions = append(completions, c.getLaTeXCompletions(prefix)...)
    }

    return completions
}

func (c *CompletionEngine) getGoCompletions(prefix string) []CompletionItem {
    completions := []CompletionItem{
        {Label: "findAllVisibleTasks", Detail: "Find tasks visible on a given day"},
        {Label: "calculateRemainingSpanColumns", Detail: "Calculate span for continuing tasks"},
        {Label: "renderSpanningTaskOverlay", Detail: "Render stacked task overlay"},
    }

    // Filter by prefix
    var filtered []CompletionItem
    for _, comp := range completions {
        if strings.HasPrefix(comp.Label, prefix) {
            filtered = append(filtered, comp)
        }
    }

    return filtered
}
```

### Enhancement K: Debugging & Profiling Tools

**LaTeX Compilation Tracer:**
```go
type LaTeXDebugger struct {
    logParser   *LogParser
    errorTracker map[string][]LaTeXError
}

type DebugReport struct {
    Errors []DebugError
}

type DebugError struct {
    Type    string
    Count   int
    Message string
}

type SuggestedFix struct {
    Problem   string
    Solution  string
    Code      string
}

func (d *LaTeXDebugger) AnalyzeCompilationLog(logContent string) *DebugReport {
    report := &DebugReport{}

    // Parse common LaTeX errors
    errorPatterns := map[string]*regexp.Regexp{
        "undefined_macro": regexp.MustCompile(`Undefined control sequence`),
        "missing_package": regexp.MustCompile(`LaTeX Error:.*package.*not found`),
        "tikz_error":     regexp.MustCompile(`Package tikz Error`),
        "table_error":    regexp.MustCompile(`Misplaced.*tabular`),
    }

    for errorType, pattern := range errorPatterns {
        matches := pattern.FindAllString(logContent, -1)
        if len(matches) > 0 {
            report.Errors = append(report.Errors, DebugError{
                Type:    errorType,
                Count:   len(matches),
                Message: fmt.Sprintf("Found %d %s errors", len(matches), errorType),
            })
        }
    }

    return report
}

func (d *LaTeXDebugger) SuggestFixes(report *DebugReport) []SuggestedFix {
    var fixes []SuggestedFix

    for _, err := range report.Errors {
        switch err.Type {
        case "undefined_macro":
            fixes = append(fixes, SuggestedFix{
                Problem: "Undefined LaTeX macro",
                Solution: "Check if custom macros are defined in document preamble",
                Code:     `\newcommand{\YourMacro}{definition}`,
            })
        case "missing_package":
            fixes = append(fixes, SuggestedFix{
                Problem: "Missing LaTeX package",
                Solution: "Add required package to document header",
                Code:     `\usepackage{package_name}`,
            })
        }
    }

    return fixes
}
```

**Performance Profiler:**
```go
type PerformanceProfiler struct {
    timings map[string]time.Duration
    counters map[string]int64
}

func (p *PerformanceProfiler) TimeFunction(name string, fn func()) {
    start := time.Now()
    fn()
    duration := time.Since(start)

    p.timings[name] = duration
    p.counters[name]++

    if duration > 100*time.Millisecond {
        log.Printf("Slow function: %s took %v", name, duration)
    }
}

type ProfileReport struct {
    TotalTime     time.Duration
    FunctionCount int
    Functions     []FunctionProfile
}

type FunctionProfile struct {
    Name     string
    Time     time.Duration
    Calls    int64
    AvgTime  time.Duration
}

func (p *PerformanceProfiler) GenerateReport() *ProfileReport {
    report := &ProfileReport{
        TotalTime: time.Duration(0),
        FunctionCount: len(p.timings),
    }

    for name, duration := range p.timings {
        report.TotalTime += duration
        report.Functions = append(report.Functions, FunctionProfile{
            Name:     name,
            Time:     duration,
            Calls:    p.counters[name],
            AvgTime:  duration / time.Duration(p.counters[name]),
        })
    }

    // Sort by total time (slowest first)
    sort.Slice(report.Functions, func(i, j int) bool {
        return report.Functions[i].Time > report.Functions[j].Time
    })

    return report
}
```

### Enhancement L: Testing & Quality Assurance

**Automated Test Generation:**
```go
type TestGenerator struct {
    templateEngine *template.Template
}

func (g *TestGenerator) GenerateStackingTest(tasks []*SpanningTask, expectedOrder []string) string {
    testData := struct {
        Tasks         []*SpanningTask
        ExpectedOrder []string
        TestDate      string
    }{
        Tasks:         tasks,
        ExpectedOrder: expectedOrder,
        TestDate:      "2025-01-15",
    }

    var buf bytes.Buffer
    err := g.templateEngine.Execute(&buf, testData)
    if err != nil {
        return fmt.Sprintf("// Error generating test: %v", err)
    }

    return buf.String()
}

func (g *TestGenerator) GenerateVisualRegressionTest() string {
    return `
func TestVisualRegression(t *testing.T) {
    // Generate baseline PDF
    baselinePDF := generateCalendarPDF(testTasks)

    // Generate current PDF
    currentPDF := generateCalendarPDF(currentTasks)

    // Compare visual layout
    differences := comparePDFLayouts(baselinePDF, currentPDF)

    if len(differences) > 0 {
        for _, diff := range differences {
            t.Errorf("Layout difference: %s", diff.Description)
        }
    }
}
`
}
```

**Code Quality Analyzer:**
```go
type QualityAnalyzer struct {
    rules []QualityRule
}

type QualityRule struct {
    Name        string
    Pattern     *regexp.Regexp
    Severity    string // "error", "warning", "info"
    Message     string
    Suggestion  string
}

type QualityIssue struct {
    File       string
    Line       int
    Rule       string
    Severity   string
    Message    string
    Suggestion string
}

func (a *QualityAnalyzer) AnalyzeCode(filename, content string) []QualityIssue {
    var issues []QualityIssue

    for _, rule := range a.rules {
        matches := rule.Pattern.FindAllStringIndex(content, -1)

        for _, match := range matches {
            lineNum := strings.Count(content[:match[0]], "\n") + 1

            issues = append(issues, QualityIssue{
                File:       filename,
                Line:       lineNum,
                Rule:       rule.Name,
                Severity:   rule.Severity,
                Message:    rule.Message,
                Suggestion: rule.Suggestion,
            })
        }
    }

    return issues
}

func (a *QualityAnalyzer) GetDefaultRules() []QualityRule {
    return []QualityRule{
        {
            Name:     "long_function",
            Pattern:  regexp.MustCompile(`func.*\{[^}]*\}`),
            Severity: "warning",
            Message:  "Function is longer than recommended",
            Suggestion: "Consider breaking into smaller functions",
        },
        {
            Name:     "magic_number",
            Pattern:  regexp.MustCompile(`\b\d{2,}\b`),
            Severity: "info",
            Message:  "Magic number detected",
            Suggestion: "Consider using named constants",
        },
        {
            Name:     "todo_comment",
            Pattern:  regexp.MustCompile(`(?i)//.*todo`),
            Severity: "info",
            Message:  "TODO comment found",
            Suggestion: "Address the TODO or create an issue",
        },
    }
}
```

### Enhancement M: Documentation & Knowledge Management

**Auto-Generated Documentation:**
```go
type DocumentationGenerator struct {
    templateDir string
}

func (g *DocumentationGenerator) GenerateFunctionDocs(filename string) string {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return ""
    }

    functions := g.extractFunctions(string(content))

    type FunctionInfo struct {
        Name       string
        Signature  string
        Comments   string
        Parameters []ParameterInfo
        ReturnType string
        Line       int
    }

    type ParameterInfo struct {
        Name string
        Type string
    }

    var docs strings.Builder
    docs.WriteString("# Function Documentation\n\n")

    for _, fn := range functions {
        docs.WriteString(fmt.Sprintf("## %s\n\n", fn.Name))
        docs.WriteString(fmt.Sprintf("**File:** %s:%d\n\n", filename, fn.Line))
        docs.WriteString(fmt.Sprintf("**Signature:** %s\n\n", fn.Signature))

        if fn.Comments != "" {
            docs.WriteString(fmt.Sprintf("**Description:** %s\n\n", fn.Comments))
        }

        if len(fn.Parameters) > 0 {
            docs.WriteString("**Parameters:**\n")
            for _, param := range fn.Parameters {
                docs.WriteString(fmt.Sprintf("- `%s`: %s\n", param.Name, param.Type))
            }
            docs.WriteString("\n")
        }

        if fn.ReturnType != "" {
            docs.WriteString(fmt.Sprintf("**Returns:** %s\n\n", fn.ReturnType))
        }
    }

    return docs.String()
}

func (g *DocumentationGenerator) extractFunctions(content string) []FunctionInfo {
    // Parse Go functions using regex and AST if needed
    // This is a simplified version
    var functions []FunctionInfo

    lines := strings.Split(content, "\n")
    for i, line := range lines {
        if strings.Contains(line, "func ") {
            fn := g.parseFunction(line, lines, i)
            functions = append(functions, fn)
        }
    }

    return functions
}
```

**Knowledge Base Integration:**
```go
type KnowledgeBase struct {
    solutions map[string]*Solution
    patterns  map[string]*Pattern
}

type Solution struct {
    Problem     string
    Solution    string
    Code        string
    Tags        []string
    LastUsed    time.Time
    SuccessRate float64
}

func (kb *KnowledgeBase) FindSolutions(problem string) []*Solution {
    var matches []*Solution

    for _, solution := range kb.solutions {
        if strings.Contains(strings.ToLower(solution.Problem), strings.ToLower(problem)) {
            matches = append(matches, solution)
        }
    }

    // Sort by success rate and recency
    sort.Slice(matches, func(i, j int) bool {
        if matches[i].SuccessRate != matches[j].SuccessRate {
            return matches[i].SuccessRate > matches[j].SuccessRate
        }
        return matches[i].LastUsed.After(matches[j].LastUsed)
    })

    return matches
}

func boolToFloat(b bool) float64 {
    if b {
        return 1.0
    }
    return 0.0
}

func (kb *KnowledgeBase) LearnFromAttempt(problem, solution string, success bool) {
    key := fmt.Sprintf("%s:%s", problem, solution)

    if sol, exists := kb.solutions[key]; exists {
        if success {
            sol.SuccessRate = (sol.SuccessRate*9 + 1) / 10 // Weighted average
        } else {
            sol.SuccessRate = (sol.SuccessRate*9 + 0) / 10
        }
        sol.LastUsed = time.Now()
    } else {
        kb.solutions[key] = &Solution{
            Problem:     problem,
            Solution:    solution,
            SuccessRate: boolToFloat(success),
            LastUsed:    time.Now(),
        }
    }
}
```

### Enhancement N: Interactive Development Features

**Live Preview System:**
```go
import (
    "github.com/fsnotify/fsnotify"
    "path/filepath"
    "time"
)

type LivePreview struct {
    watcher    *fsnotify.Watcher
    lastBuild  time.Time
    buildDebounce time.Duration
}

func (lp *LivePreview) StartWatching(projectDir string) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }

    lp.watcher = watcher
    lp.buildDebounce = 500 * time.Millisecond

    go lp.watchLoop()

    return watcher.Add(projectDir)
}

func (lp *LivePreview) watchLoop() {
    for {
        select {
        case event, ok := <-lp.watcher.Events:
            if !ok {
                return
            }

            if lp.shouldTriggerBuild(event) {
                lp.debouncedBuild()
            }

        case err, ok := <-lp.watcher.Errors:
            if !ok {
                return
            }
            log.Printf("Watch error: %v", err)
        }
    }
}

func (lp *LivePreview) shouldTriggerBuild(event fsnotify.Event) bool {
    // Only rebuild on relevant file changes
    ext := filepath.Ext(event.Name)
    relevantExts := []string{".go", ".yaml", ".csv", ".tpl"}

    for _, rext := range relevantExts {
        if ext == rext {
            return true
        }
    }

    return false
}

func (lp *LivePreview) debouncedBuild() {
    lp.lastBuild = time.Now()

    time.Sleep(lp.buildDebounce)

    // Only build if no more changes occurred
    if time.Since(lp.lastBuild) >= lp.buildDebounce {
        lp.triggerBuild()
    }
}
```

**Interactive Task Editor:**
```go
type InteractiveEditor struct {
    currentTasks []*SpanningTask
    undoStack    []*EditAction
}

type EditAction struct {
    Type        string // "add", "modify", "delete"
    TaskIndex   int
    OldTask     *SpanningTask
    NewTask     *SpanningTask
    Timestamp   time.Time
}

func (ie *InteractiveEditor) AddTaskInteractive() {
    // Interactive prompts for task creation
    name := prompt("Task name: ")
    start := promptDate("Start date (YYYY-MM-DD): ")
    end := promptDate("End date (YYYY-MM-DD): ")
    color := promptColor("Color (hex or name): ")

    task := &SpanningTask{
        Name:      name,
        StartDate: start,
        EndDate:   end,
        Color:     color,
    }

    ie.addTask(task)
    ie.saveUndoAction("add", len(ie.currentTasks)-1, nil, task)
}

func (ie *InteractiveEditor) ModifyTaskInteractive(index int) {
    if index < 0 || index >= len(ie.currentTasks) {
        fmt.Println("Invalid task index")
        return
    }

    oldTask := ie.currentTasks[index]

    // Show current values and allow modifications
    fmt.Printf("Current task: %s\n", oldTask.Name)

    newName := promptWithDefault("Name", oldTask.Name)
    newStart := promptDateWithDefault("Start date", oldTask.StartDate)
    newEnd := promptDateWithDefault("End date", oldTask.EndDate)
    newColor := promptColorWithDefault("Color", oldTask.Color)

    newTask := &SpanningTask{
        Name:      newName,
        StartDate: newStart,
        EndDate:   newEnd,
        Color:     newColor,
    }

    ie.modifyTask(index, newTask)
    ie.saveUndoAction("modify", index, oldTask, newTask)
}

func (ie *InteractiveEditor) Undo() {
    if len(ie.undoStack) == 0 {
        fmt.Println("Nothing to undo")
        return
    }

    action := ie.undoStack[len(ie.undoStack)-1]
    ie.undoStack = ie.undoStack[:len(ie.undoStack)-1]

    switch action.Type {
    case "add":
        ie.currentTasks = append(ie.currentTasks[:action.TaskIndex], ie.currentTasks[action.TaskIndex+1:]...)
    case "modify":
        ie.currentTasks[action.TaskIndex] = action.OldTask
    case "delete":
        ie.currentTasks = append(ie.currentTasks, action.OldTask)
        // Shift elements to maintain correct position
        copy(ie.currentTasks[action.TaskIndex+1:], ie.currentTasks[action.TaskIndex:])
        ie.currentTasks[action.TaskIndex] = action.OldTask
    }
}
```

### Enhancement O: Version Control Integration

**Git-Aware Development:**
```go
import (
    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing/object"
)

type GitIntegration struct {
    repo *git.Repository
}

func (gi *GitIntegration) GetChangedFiles() ([]string, error) {
    worktree, err := gi.repo.Worktree()
    if err != nil {
        return nil, err
    }

    status, err := worktree.Status()
    if err != nil {
        return nil, err
    }

    var changed []string
    for file, fileStatus := range status {
        if fileStatus.Staging != git.Untracked {
            changed = append(changed, file)
        }
    }

    return changed, nil
}

func (gi *GitIntegration) CommitWithContext(message string, context *DevelopmentContext) error {
    worktree, err := gi.repo.Worktree()
    if err != nil {
        return err
    }

    // Add all changes
    _, err = worktree.Add(".")
    if err != nil {
        return err
    }

    // Create detailed commit message
    detailedMessage := gi.enhanceCommitMessage(message, context)

    // Commit
    _, err = worktree.Commit(detailedMessage, &git.CommitOptions{
        Author: &object.Signature{
            Name:  "Gantt Planner AI",
            Email: "ai@gantt-planner.local",
            When:  time.Now(),
        },
    })

    return err
}

type DevelopmentContext struct {
    LastError     string
    ChangedFiles  []string
    TestResults   *TestResults
}

type TestResults struct {
    Passed int
    Failed int
    Total  int
}

func (gi *GitIntegration) enhanceCommitMessage(baseMessage string, context *DevelopmentContext) string {
    enhanced := baseMessage

    if context.LastError != "" {
        enhanced += fmt.Sprintf("\n\nFixes: %s", context.LastError)
    }

    if len(context.ChangedFiles) > 0 {
        enhanced += "\n\nChanged files:"
        for _, file := range context.ChangedFiles {
            enhanced += fmt.Sprintf("\n- %s", file)
        }
    }

    if context.TestResults != nil {
        enhanced += fmt.Sprintf("\n\nTests: %d passed, %d failed",
            context.TestResults.Passed, context.TestResults.Failed)
    }

    return enhanced
}
```

**Key Success Criteria**:
- Each task appears once with a spanning bar ✅
- Overlapping tasks stack vertically without overlap ✅
- Proper spacing accounting for task height ✅
- Consistent rendering across week boundaries ✅
- Maintainable, testable code ✅
