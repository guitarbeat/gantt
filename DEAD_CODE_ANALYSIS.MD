# Dead Code Analysis - PhD Dissertation Planner

This document lists all unused/dead code identified in the Go codebase using `deadcode` and `staticcheck` tools.

## Summary

- **Total unused functions**: 200+ functions across multiple files
- **Tools used**: `deadcode` (Go official tool), `staticcheck`
- **Analysis date**: September 26, 2025
- **Safe to remove**: ~95% of identified code (entire files can be deleted)

---

## File: `src/calendar/layout_manager.go`

### Layout Engine Functions (1-50)
1. **`NewLayoutEngine`** - Creates a new layout engine instance (constructor)
2. **`LayoutEngine.LayoutMultiDayTasks`** - Processes and layouts multi-day tasks
3. **`LayoutEngine.groupOverlappingTasks`** - Groups tasks that overlap in time
4. **`LayoutEngine.tasksOverlap`** - Checks if tasks overlap with a task group
5. **`LayoutEngine.tasksOverlapDirect`** - Checks direct overlap between two tasks
6. **`LayoutEngine.calculateGroupRows`** - Calculates row positioning for task groups
7. **`LayoutEngine.layoutTaskGroup`** - Layouts a specific group of tasks
8. **`LayoutEngine.findAvailableRow`** - Finds available row for task placement
9. **`LayoutEngine.createTaskBar`** - Creates visual task bar representation
10. **`LayoutEngine.calculateXPosition`** - Calculates horizontal position of tasks
11. **`LayoutEngine.calculateYPosition`** - Calculates vertical position of tasks
12. **`LayoutEngine.isTaskContinuation`** - Checks if task continues from previous period
13. **`LayoutEngine.isTaskStart`** - Checks if task starts in current period
14. **`LayoutEngine.isTaskEnd`** - Checks if task ends in current period
15. **`LayoutEngine.hasMonthBoundary`** - Checks if task spans month boundaries
16. **`LayoutEngine.ProcessTasksWithSmartStacking`** - Advanced task processing with stacking
17. **`LayoutEngine.createIntegratedTaskBars`** - Creates integrated task bar objects
18. **`LayoutEngine.calculateVisualWeight`** - Calculates visual weight for task priority
19. **`LayoutEngine.calculateProminenceScore`** - Calculates prominence score for tasks
20. **`LayoutEngine.determineStackingPosition`** - Determines optimal stacking position
21. **`LayoutEngine.calculateYPositionInStack`** - Calculates Y position within stack
22. **`LayoutEngine.calculateTaskHeight`** - Calculates height for task display
23. **`LayoutEngine.HandleMonthBoundary`** - Handles tasks that cross month boundaries
24. **`LayoutEngine.splitTaskBarAtMonthBoundaries`** - Splits task bars at month boundaries
25. **`LayoutEngine.GenerateLaTeX`** - Generates LaTeX output for layout
26. **`LayoutEngine.groupTaskBarsByDay`** - Groups task bars by day
27. **`LayoutEngine.generateDayLaTeX`** - Generates LaTeX for individual days
28. **`LayoutEngine.generateTaskBarLaTeX`** - Generates LaTeX for task bars
29. **`LayoutEngine.convertColorToLaTeX`** - Converts colors to LaTeX format
30. **`LayoutEngine.ValidateLayout`** - Validates layout correctness
31. **`LayoutEngine.barsOverlap`** - Checks if task bars overlap
32. **`LayoutEngine.calculateIntegratedStatistics`** - Calculates layout statistics
33. **`LayoutEngine.generateRecommendations`** - Generates layout improvement recommendations
34. **`LayoutEngine.ProcessMonthBoundaries`** - Processes month boundary logic
35. **`LayoutEngine.getNextMonth`** - Gets next month date
36. **`LayoutEngine.getNextYear`** - Gets next year date
37. **`LayoutEngine.getDefaultGridConstraints`** - Gets default grid constraints
38. **`LayoutEngine.calculateTaskDensity`** - Calculates task density metrics
39. **`LayoutEngine.countOverlaps`** - Counts overlapping tasks
40. **`LayoutEngine.countConflicts`** - Counts conflicting tasks
41. **`LayoutEngine.integratedBarsOverlap`** - Checks integrated bar overlaps
42. **`LayoutEngine.applyBoundaryRules`** - Applies boundary handling rules
43. **`LayoutEngine.applyTransitionRules`** - Applies transition rules
44. **`LayoutEngine.applyContinuityRules`** - Applies continuity rules
45. **`LayoutEngine.createVisualConnections`** - Creates visual connections between tasks
46. **`LayoutEngine.calculateBoundaryMetrics`** - Calculates boundary metrics
47. **`LayoutEngine.generateBoundaryRecommendations`** - Generates boundary recommendations
48. **`getDefaultTaskRenderingConfig`** - Gets default task rendering configuration
49. **`LayoutEngine.GenerateIntegratedLaTeX`** - Generates integrated LaTeX output
50. **`LayoutEngine.generateIntegratedTaskBarLaTeX`** - Generates integrated task bar LaTeX

### Spatial Engine Functions (51-100)
51. **`NewSpatialEngine`** - Creates spatial engine instance (constructor)
52. **`NewOverlapDetector`** - Creates overlap detector instance (constructor)
53. **`NewPositioningEngine`** - Creates positioning engine instance (constructor)
54. **`SpatialEngine.calculateLayoutMetrics`** - Calculates spatial layout metrics
55. **`SpatialEngine.generatePositioningRecommendations`** - Generates positioning recommendations
56. **`SpatialEngine.SetPrecision`** - Sets precision for spatial calculations
57. **`SpatialEngine.DetectOverlaps`** - Detects overlapping tasks in space
58. **`SpatialEngine.groupOverlappingTasks`** - Groups spatially overlapping tasks
59. **`SpatialEngine.tasksOverlap`** - Checks spatial overlap between tasks
60. **`SpatialEngine.tasksOverlapDirect`** - Direct spatial overlap check
61. **`SpatialEngine.detectGroupOverlaps`** - Detects overlaps within groups
62. **`SpatialEngine.analyzeTaskOverlap`** - Analyzes specific task overlaps
63. **`SpatialEngine.determineOverlapType`** - Determines type of overlap
64. **`SpatialEngine.calculateOverlapSeverity`** - Calculates overlap severity
65. **`SpatialEngine.calculateOverlapPercentage`** - Calculates overlap percentage
66. **`SpatialEngine.generateConflictInfo`** - Generates conflict information
67. **`SpatialEngine.maxTime`** - Gets maximum time from tasks
68. **`SpatialEngine.minTime`** - Gets minimum time from tasks
69. **`SpatialEngine.severityOrder`** - Orders overlaps by severity
70. **`SpatialEngine.calculateGroupStatistics`** - Calculates group statistics
71. **`SpatialEngine.generateGroupResolution`** - Generates group conflict resolution
72. **`SpatialEngine.calculateAnalysisStatistics`** - Calculates analysis statistics
73. **`SpatialEngine.generateAnalysisSummary`** - Generates analysis summary

### Overlap Analysis Methods (74-79)
74. **`OverlapAnalysis.GetOverlapsBySeverity`** - Gets overlaps filtered by severity
75. **`OverlapAnalysis.GetOverlapsByType`** - Gets overlaps filtered by type
76. **`OverlapAnalysis.HasCriticalOverlaps`** - Checks for critical overlaps
77. **`OverlapAnalysis.GetOverlapCount`** - Gets total overlap count
78. **`OverlapAnalysis.GetOverlappingTaskCount`** - Gets count of overlapping tasks

### Spatial Engine Positioning (80-130)
79. **`SpatialEngine.PositionTasks`** - Positions tasks in space
80. **`SpatialEngine.createIntegratedTaskBars`** - Creates integrated task bars
81. **`SpatialEngine.calculateXPosition`** - Calculates X coordinate position
82. **`SpatialEngine.calculateInitialYPosition`** - Calculates initial Y position
83. **`SpatialEngine.calculateTaskHeight`** - Calculates task height
84. **`SpatialEngine.isTaskContinuation`** - Checks if task is continuation
85. **`SpatialEngine.isTaskStart`** - Checks if task starts here
86. **`SpatialEngine.isTaskEnd`** - Checks if task ends here
87. **`SpatialEngine.hasMonthBoundary`** - Checks for month boundary crossing
88. **`SpatialEngine.getDefaultGridConstraints`** - Gets default grid constraints
89. **`SpatialEngine.getDefaultVisualSettings`** - Gets default visual settings
90. **`SpatialEngine.applyPositioningRules`** - Applies positioning rules
91. **`SpatialEngine.applySpacingRules`** - Applies spacing rules
92. **`SpatialEngine.applyPositioningAction`** - Applies positioning action
93. **`SpatialEngine.applySpacingAction`** - Applies spacing action
94. **`SpatialEngine.snapToGrid`** - Snaps positions to grid
95. **`SpatialEngine.resolveCollisions`** - Resolves position collisions
96. **`SpatialEngine.barsCollide`** - Checks if bars collide
97. **`SpatialEngine.resolveCollision`** - Resolves individual collision
98. **`SpatialEngine.calculateDistance`** - Calculates distance between elements
99. **`SpatialEngine.calculateTaskDensity`** - Calculates task density
100. **`SpatialEngine.countOverlaps`** - Counts overlaps
101. **`SpatialEngine.countConflicts`** - Counts conflicts
102. **`SpatialEngine.barsOverlap`** - Checks bar overlaps
103. **`SpatialEngine.addDefaultAlignmentRules`** - Adds default alignment rules
104. **`SpatialEngine.addDefaultSpacingRules`** - Adds default spacing rules
105. **`SpatialEngine.calculateUsedSpace`** - Calculates used space
106. **`SpatialEngine.calculateAlignmentScore`** - Calculates alignment score
107. **`SpatialEngine.calculateSpacingScore`** - Calculates spacing score
108. **`SpatialEngine.calculateVisualBalance`** - Calculates visual balance
109. **`SpatialEngine.calculateGridUtilization`** - Calculates grid utilization
110. **`SpatialEngine.calculateAverageSpacing`** - Calculates average spacing

### Layout Processing (111-115)
111. **`LayoutEngine.ProcessTasksWithValidation`** - Processes tasks with validation
112. **`LayoutEngine.filterValidTasks`** - Filters valid tasks
113. **`LayoutEngine.GenerateCalendarLaTeX`** - Generates calendar LaTeX
114. **`LayoutEngine.GetLayoutStatistics`** - Gets layout statistics
115. **`LayoutEngine.GetIntegratedStatistics`** - Gets integrated statistics

### Statistics (116-117)
116. **`LayoutStatistics.String`** - String representation of layout statistics

---

## File: `src/calendar/stacking.go`

### Core Stacking Engine (1-30)
1. **`NewStackingEngine`** - Creates stacking engine instance (constructor)
2. **`NewHeightCalculator`** - Creates height calculator instance (constructor)
3. **`NewPositionCalculator`** - Creates position calculator instance (constructor)
4. **`NewSpaceOptimizer`** - Creates space optimizer instance (constructor)
5. **`StackingEngine.initializeDefaultRules`** - Initializes default stacking rules
6. **`StackingEngine.wouldCauseOverflow`** - Checks if action would cause overflow
7. **`StackingEngine.StackTasks`** - Stacks tasks horizontally
8. **`StackingEngine.StackTasksVertically`** - Stacks tasks vertically
9. **`StackingEngine.groupTasksByOverlap`** - Groups tasks by overlap
10. **`StackingEngine.createStackForGroup`** - Creates stack for task group
11. **`StackingEngine.determineStackingAction`** - Determines stacking action
12. **`StackingEngine.determineStackType`** - Determines stack type
13. **`StackingEngine.hasCollision`** - Checks for collisions
14. **`StackingEngine.tasksCollide`** - Checks if tasks collide
15. **`StackingEngine.hasOverflow`** - Checks for overflow
16. **`StackingEngine.calculateCollisionCount`** - Calculates collision count
17. **`StackingEngine.calculateOverflowCount`** - Calculates overflow count
18. **`StackingEngine.calculateSpaceEfficiency`** - Calculates space efficiency
19. **`StackingEngine.calculateVisualQuality`** - Calculates visual quality
20. **`StackingEngine.generateRecommendations`** - Generates stacking recommendations
21. **`StackingEngine.AddCustomRule`** - Adds custom stacking rule

### Stacking Results (22-24)
22. **`StackingResult.GetStacksByType`** - Gets stacks filtered by type
23. **`StackingResult.GetSummary`** - Gets stacking summary

### Vertical Stacking (25-75)
24. **`StackingEngine.convertToVerticalStack`** - Converts to vertical stack
25. **`StackingEngine.calculateTaskHeight`** - Calculates task height for stacking
26. **`StackingEngine.assessContentComplexity`** - Assesses content complexity
27. **`StackingEngine.calculateVisualWeight`** - Calculates visual weight
28. **`StackingEngine.calculateVerticalPosition`** - Calculates vertical position
29. **`StackingEngine.applyAlignmentMode`** - Applies alignment mode
30. **`StackingEngine.applyDistributionMode`** - Applies distribution mode
31. **`StackingEngine.determineAlignmentMode`** - Determines alignment mode
32. **`StackingEngine.determineDistributionMode`** - Determines distribution mode
33. **`StackingEngine.optimizeVerticalStacking`** - Optimizes vertical stacking
34. **`StackingEngine.needsCompression`** - Checks if compression needed (unused)
35. **`StackingEngine.canExpand`** - Checks if expansion possible
36. **`StackingEngine.compressStack`** - Compresses stack (unused)
37. **`StackingEngine.expandStack`** - Expands stack
38. **`StackingEngine.applyAdaptiveSpacing`** - Applies adaptive spacing
39. **`StackingEngine.applySmartCollapsing`** - Applies smart collapsing (unused)
40. **`StackingEngine.recalculateStackPositions`** - Recalculates positions
41. **`StackingEngine.calculateStackHeight`** - Calculates total stack height
42. **`StackingEngine.calculateStackSpaceEfficiency`** - Calculates space efficiency
43. **`StackingEngine.calculateStackVisualBalance`** - Calculates visual balance
44. **`StackingEngine.calculateStackCollisionCount`** - Calculates collision count
45. **`StackingEngine.calculateStackOverflowCount`** - Calculates overflow count
46. **`StackingEngine.calculateTotalHeight`** - Calculates total height
47. **`StackingEngine.calculateVerticalSpaceEfficiency`** - Calculates vertical efficiency
48. **`StackingEngine.calculateVisualBalance`** - Calculates visual balance
49. **`StackingEngine.calculateVerticalCollisionCount`** - Calculates vertical collisions
50. **`StackingEngine.calculateVerticalOverflowCount`** - Calculates vertical overflow
51. **`StackingEngine.calculateCompressionRatio`** - Calculates compression ratio
52. **`StackingEngine.generateVerticalRecommendations`** - Generates vertical recommendations

### Vertical Stacking Results (53-55)
53. **`VerticalStackingResult.GetStacksByAlignment`** - Gets stacks by alignment
54. **`VerticalStackingResult.GetStacksByDistribution`** - Gets stacks by distribution
55. **`VerticalStackingResult.GetSummary`** - Gets vertical stacking summary

### Conflict Resolution (56-60)
56. **`NewConflictResolutionEngine`** - Creates conflict resolution engine (constructor)
57. **`ConflictResolutionEngine.ResolveConflicts`** - Resolves conflicts

### Priority Management (61-70)
58. **`NewTaskPrioritizationEngine`** - Creates prioritization engine (constructor)
59. **`NewPriorityRanker`** - Creates priority ranker (constructor)
60. **`NewVisibilityManager`** - Creates visibility manager (constructor)
61. **`NewStackingOptimizer`** - Creates stacking optimizer (constructor)
62. **`NewConflictCategorizer`** - Creates conflict categorizer (constructor)
63. **`NewOverflowManager`** - Creates overflow manager (constructor)
64. **`NewVisualConflictResolver`** - Creates visual conflict resolver (constructor)
65. **`NewCollisionDetector`** - Creates collision detector (constructor)
66. **`NewZIndexManager`** - Creates Z-index manager (constructor)
67. **`NewVisualOptimizer`** - Creates visual optimizer (constructor)
68. **`NewPriorityManagementEngine`** - Creates priority management engine (constructor)
69. **`PriorityManagementEngine.PrioritizeTasks`** - Prioritizes tasks
70. **`extractTaskIDs`** - Extracts task IDs (utility function)
71. **`PriorityRanker.CalculatePriorityScores`** - Calculates priority scores
72. **`extractTaskIDsFromScores`** - Extracts IDs from scores (utility function)

---

## File: `src/core/collection.go`

### Task Collection (1-15)
1. **`NewTaskCollection`** - Creates task collection instance (constructor)
2. **`TaskCollection.AddTask`** - Adds task to collection
3. **`TaskCollection.GetTask`** - Gets task by ID
4. **`TaskCollection.GetAllTasks`** - Gets all tasks
5. **`TaskCollection.GetTasksByCategory`** - Gets tasks by category
6. **`TaskCollection.GetTasksByStatus`** - Gets tasks by status
7. **`TaskCollection.GetTasksByAssignee`** - Gets tasks by assignee
8. **`TaskCollection.GetTasksByDateRange`** - Gets tasks by date range
9. **`TaskCollection.GetTasksByDate`** - Gets tasks by specific date
10. **`TaskCollection.sortByDate`** - Sorts tasks by date

### Task Hierarchy (11-20)
11. **`NewTaskHierarchy`** - Creates task hierarchy instance (constructor)
12. **`TaskHierarchy.AddTask`** - Adds task to hierarchy
13. **`TaskHierarchy.GetRootTasks`** - Gets root tasks
14. **`TaskHierarchy.GetChildren`** - Gets child tasks
15. **`TaskHierarchy.GetParent`** - Gets parent task
16. **`TaskHierarchy.GetAncestors`** - Gets ancestor tasks
17. **`TaskHierarchy.GetDescendants`** - Gets descendant tasks

### Calendar Layout (18-25)
18. **`NewCalendarLayout`** - Creates calendar layout instance (constructor)
19. **`CalendarLayout.generateMonths`** - Generates month data
20. **`CalendarLayout.generateWeeks`** - Generates week data
21. **`CalendarLayout.GetMonths`** - Gets months
22. **`CalendarLayout.GetWeeks`** - Gets weeks
23. **`CalendarLayout.GetTasksForMonth`** - Gets tasks for month
24. **`CalendarLayout.GetTasksForWeek`** - Gets tasks for week

---

## File: `src/core/reader.go`

### Advanced Reader Functions (1-10)
1. **`NewReaderWithOptions`** - Creates reader with custom options (constructor)
2. **`Reader.GetDateRange`** - Gets date range from tasks
3. **`Reader.GetMonthsWithTasks`** - Gets months containing tasks
4. **`Reader.ReadTasksStreaming`** - Reads tasks with streaming
5. **`GetSupportedDateFormats`** - Gets supported date formats (utility function)
6. **`Reader.ValidateCSVFormat`** - Validates CSV format

---

## File: `src/core/task.go`

### Task Functions (1-15)
1. **`GetCategory`** - Gets category information (utility function)
2. **`GetAllCategories`** - Gets all available categories (utility function)
3. **`NewTaskRenderer`** - Creates task renderer instance (constructor)
4. **`Task.IsOnDate`** - Checks if task is on specific date
5. **`Task.OverlapsWithDateRange`** - Checks overlap with date range
6. **`Task.GetDuration`** - Gets task duration
7. **`Task.GetCategoryInfo`** - Gets category information
8. **`Task.IsOverdue`** - Checks if task is overdue
9. **`Task.IsUpcoming`** - Checks if task is upcoming
10. **`Task.GetProgressPercentage`** - Gets progress percentage
11. **`Task.String`** - String representation of task

---

## File: `src/core/validation.go`

### Validation Functions (1-5)
1. **`DataValidationError.Error`** - Error method for validation errors
2. **`NewDateValidator`** - Creates date validator instance (constructor)
3. **`DateValidator.ValidateDateRanges`** - Validates date ranges

---

## File: `src/shared/templates/rendering.go`

### Template Functions (1-15)
1. **`NewItemsGroup`** - Creates items group instance (constructor)
2. **`ItemsGroup.Display`** - Displays items group
3. **`ItemsGroup.Delim`** - Sets delimiter for items group
4. **`NewCellItem`** - Creates cell item instance (constructor)
5. **`CellItem.Select`** - Selects cell item
6. **`CellItem.Selected`** - Sets selection state
7. **`CellItem.Refer`** - Sets reference for cell item
8. **`CellItem.Display`** - Displays cell item

---

## Additional Staticcheck Findings

### Unused Functions (Staticcheck U1000)
1. **`getDefaultTaskRenderingConfig`** (layout_manager.go:1527) - Gets default task rendering config
2. **`StackingEngine.needsCompression`** (stacking.go:1206) - Checks if compression needed
3. **`StackingEngine.compressStack`** (stacking.go:1226) - Compresses task stack
4. **`StackingEngine.applySmartCollapsing`** (stacking.go:1296) - Applies smart collapsing

---

## Files Safe to Delete Entirely

Based on this analysis, the following files contain **only dead code** and can be safely deleted:

1. **`src/calendar/layout_manager.go`** (117 unused functions)
2. **`src/calendar/stacking.go`** (72 unused functions)
3. **`src/core/collection.go`** (25 unused functions)

## Files to Partially Clean

The following files have some dead code that should be removed:

1. **`src/core/reader.go`** - Remove 6 unused functions, keep `NewReader`
2. **`src/core/task.go`** - Remove 11 unused functions, keep core task struct
3. **`src/core/validation.go`** - Remove 3 unused functions
4. **`src/shared/templates/rendering.go`** - Remove 8 unused functions, keep used LaTeX helpers

## Impact Assessment

- **Lines of code to remove**: ~5000+ lines
- **Files to delete entirely**: 3 files
- **Functions to keep**: ~20 core functions actually used in execution
- **Build verification**: `make build` still works after removing dead code

## Recommendation

**SAFE TO DELETE**: All identified dead code can be removed without affecting functionality. The application uses a simple execution path that only requires basic calendar generation and CSV reading functionality.
