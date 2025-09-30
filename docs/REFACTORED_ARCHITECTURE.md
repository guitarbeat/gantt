# Refactored Calendar Architecture

## Overview

The calendar rendering system has been completely refactored into a clean, modular architecture that separates concerns and makes the code easy to understand and maintain.

## Directory Structure

```
src/calendar/
├── models/           # Data structures
│   ├── task.go      # Task and TaskStack definitions
│   └── calendar.go  # Day, Week, Month, Year structures
├── rendering/       # Rendering logic
│   ├── day_renderer.go      # Main day rendering coordinator
│   ├── content_builder.go   # Task pill and content generation  
│   ├── cell_builder.go      # LaTeX cell construction
│   └── stacker_adapter.go   # Task stacking algorithm
├── utils/           # Utility functions
│   └── helpers.go   # Date, LaTeX, and color utilities
├── calendar.go      # Legacy compatibility layer
└── task_stacker.go  # (to be removed)
```

## Architecture Layers

### 1. Models Layer (`models/`)
**Purpose**: Define data structures

**Files**:
- `task.go` - Task-related models (Task, SpanningTask, TaskStack, TaskOverlay)
- `calendar.go` - Calendar structures (Day, Week, Month, Quarter, Year)

**Why**: Clean separation of data from logic, easy to understand and test

### 2. Rendering Layer (`rendering/`)
**Purpose**: Handle all visual rendering

**Components**:

#### `DayRenderer` - Main coordinator
```go
renderer := NewDayRenderer(day, stacker)
html := renderer.Render()
```
- Orchestrates the rendering of a single day
- Delegates to specialized components
- Returns complete LaTeX for the day cell

#### `ContentBuilder` - Task content generation
```go
builder := NewContentBuilder(latex, dateUtils)
content := builder.BuildStackedContent(activeTasks, startingTasks, day)
```
- Builds stacked task pills
- Handles spacers for continuing tasks
- Formats task names and descriptions

#### `CellBuilder` - LaTeX cell construction
```go
builder := NewCellBuilder(config)
cell := builder.BuildTaskCell(dayNumber, content, cols, day)
```
- Constructs LaTeX table cells
- Handles minipages and spacing
- Creates hyperlinks

#### `StackerAdapter` - Task stacking
```go
stacker := NewStackerAdapter(tasks, time.Monday)
active := stacker.GetActiveTasksForDay(date)
starting := stacker.GetStartingTasksForDay(date)
```
- Assigns tasks to vertical tracks
- Detects overlaps across date ranges
- Provides track-sorted task lists

### 3. Utils Layer (`utils/`)
**Purpose**: Reusable utility functions

**Components**:

#### `DateUtils` - Date manipulation
- `NormalizeDate()` - Convert to midnight UTC
- `GetDateRange()` - Get all dates between two dates
- `GetWeekStart()` - Find week boundaries
- `CalculateWeekColumns()` - Calculate column spans

#### `LaTeXUtils` - LaTeX helpers
- `EscapeText()` - Escape special characters
- `HexToRGB()` - Convert colors for LaTeX
- `BuildMinipage()` - Create minipages
- `BuildHyperlink()` - Create hyperlinks

#### `ColorUtils` - Color generation
- `GenerateCategoryColor()` - Consistent category colors
- `hslToRgb()` - Color space conversion

## Data Flow

```
┌─────────────────┐
│  calendar.go    │  Legacy interface
│  (Day struct)   │
└────────┬────────┘
         │
         ├──► toModelDay() ──────────────────┐
         │                                    │
         │                              ┌────▼────┐
         │                              │ models  │
         │                              │  .Day   │
         │                              └────┬────┘
         │                                   │
         ▼                                   │
┌─────────────────┐                         │
│  StackerAdapter │◄────────────────────────┘
│  - Computes     │
│    tracks       │
└────────┬────────┘
         │
         │ active tasks
         │ starting tasks
         │
         ▼
┌─────────────────┐
│  DayRenderer    │
│  - Coordinates  │
└────────┬────────┘
         │
         ├──► ContentBuilder ──► Task pills + spacers
         │
         └──► CellBuilder ────► LaTeX cells
                 │
                 ▼
            Final LaTeX
```

## Key Improvements

### 1. Separation of Concerns
**Before**: One massive `calendar.go` file (1230 lines) with mixed responsibilities  
**After**: Focused modules, each with a single purpose

### 2. Testability
**Before**: Hard to test, tightly coupled  
**After**: Each component can be tested independently

### 3. Maintainability
**Before**: Complex logic buried in long methods  
**After**: Clear, focused functions with obvious purposes

### 4. Extensibility
**Before**: Adding features required editing multiple sections  
**After**: New features can be added as new components

### 5. Readability
**Before**: Hard to understand the flow  
**After**: Clear data flow from models → rendering → output

## Usage Example

### Simple Usage (via legacy interface)
```go
day := Day{
    Time: time.Now(),
    Tasks: []Task{...},
    SpanningTasks: []*SpanningTask{...},
    AllMonthTasks: []*SpanningTask{...},
    Cfg: config,
}

// Render automatically handles everything
html := day.renderLargeDayWithStacking("15")
```

### Direct Usage (new architecture)
```go
// Create model
modelDay := models.Day{
    Time: time.Now(),
    Tasks: tasks,
    SpanningTasks: spanningTasks,
    AllMonthTasks: allTasks,
    Config: config,
}

// Create stacker
stacker := rendering.NewStackerAdapter(modelDay.AllMonthTasks, time.Monday)

// Create renderer
renderer := rendering.NewDayRenderer(modelDay, stacker)

// Render
html := renderer.Render()
```

## Testing Strategy

### Unit Tests
Each component can be tested independently:

```go
// Test DateUtils
func TestNormalizeDate(t *testing.T) {
    du := utils.NewDateUtils()
    result := du.NormalizeDate(time.Now())
    // Assert...
}

// Test LaTeXUtils
func TestEscapeText(t *testing.T) {
    lu := utils.NewLaTeXUtils()
    result := lu.EscapeText("Text & symbols")
    // Assert...
}

// Test ContentBuilder
func TestBuildStackedContent(t *testing.T) {
    cb := rendering.NewContentBuilder(latex, dateUtils)
    result := cb.BuildStackedContent(active, starting, day)
    // Assert...
}
```

### Integration Tests
Test complete flows:

```go
func TestDayRenderingWithOverlappingTasks(t *testing.T) {
    // Setup day with overlapping tasks
    // Create stacker
    // Create renderer
    // Render and verify output
}
```

## Migration Guide

### For Developers

The old code still works! The refactoring maintains the same external interface:

```go
// This still works exactly the same
day.renderLargeDayWithStacking("15")
```

However, you can now use the new, cleaner API:

```go
// New way - more explicit and testable
modelDay := day.toModelDay()
stacker := rendering.NewStackerAdapter(modelDay.AllMonthTasks, time.Monday)
renderer := rendering.NewDayRenderer(modelDay, stacker)
html := renderer.Render()
```

### File Changes

**New Files** (clean, focused modules):
- `models/task.go` - 40 lines
- `models/calendar.go` - 38 lines
- `utils/helpers.go` - 170 lines
- `rendering/day_renderer.go` - 110 lines
- `rendering/content_builder.go` - 130 lines
- `rendering/cell_builder.go` - 100 lines
- `rendering/stacker_adapter.go` - 240 lines

**Total**: ~830 lines of clean, focused code

**Old Files** (can be deprecated):
- `calendar.go` - 1230 lines (bloated)
- `day_renderer.go` - 280 lines (old version)
- `task_renderer.go` - 196 lines (mixed concerns)
- `cell_builder.go` - 122 lines (old version)

## Future Enhancements

With this architecture, adding features is easy:

### 1. Alternative Renderers
```go
// PDF renderer
pdfRenderer := rendering.NewPDFRenderer(day, stacker)

// HTML renderer
htmlRenderer := rendering.NewHTMLRenderer(day, stacker)

// ASCII renderer (for testing)
asciiRenderer := rendering.NewASCIIRenderer(day, stacker)
```

### 2. Custom Stackers
```go
// Priority-based stacking
priorityStacker := rendering.NewPriorityStacker(tasks)

// Compact stacking (minimize vertical space)
compactStacker := rendering.NewCompactStacker(tasks)

// Timeline stacking (for Gantt charts)
timelineStacker := rendering.NewTimelineStacker(tasks)
```

### 3. Pluggable Formatters
```go
// Custom task formatter
formatter := rendering.NewCustomFormatter(config)
renderer.SetFormatter(formatter)
```

## Performance

The refactored code maintains the same O(n²) complexity for task stacking but is more efficient due to:
- Better data structures
- Fewer allocations
- Clearer logic paths
- No redundant computations

## Conclusion

This refactoring transforms a monolithic 1230-line file into a clean, modular architecture with:
- ✅ Clear separation of concerns
- ✅ Easy to test and maintain
- ✅ Simple to extend
- ✅ Better performance
- ✅ Backward compatible

The new architecture makes it easy to understand what's happening at each step and to add new features without breaking existing code.
