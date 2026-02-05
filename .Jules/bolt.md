# Bolt's Performance Journal ⚡

## 2024-05-23 - Redundant Filtering and Allocation in Hot Loop
**Learning:** In the `findActiveTasks` function, `d.Tasks` was being filtered into a new `activeTasks` slice, even though `d.Tasks` was already guaranteed to contain only valid tasks for the day (populated by `ApplySpanningTasksToMonth`). This resulted in unnecessary O(N) allocations and redundant checks in the hot rendering loop.
**Action:** Trust upstream data guarantees when possible. If a collection is pre-filtered, avoid re-filtering. Modifying the collection in-place (e.g., sorting) can save allocations if the order isn't required to be preserved for other callers.

## 2024-05-23 - Map Allocation vs Array in Tight Loops
**Learning:** `assignTaskTracks` was allocating a `map[int][]*SpanningTask` for every single day rendered. Since the key space (tracks) is small (<100) and integers, a fixed-size array `[100][]*SpanningTask` eliminates the map allocation overhead and hashing costs entirely.
**Action:** Use fixed-size arrays instead of maps for small, integer-keyed lookups in hot loops, especially when the maximum key value is known and small.

## 2024-05-23 - Safety in Fixed-Size Array Optimizations
**Learning:** When replacing dynamic maps with fixed-size arrays for performance, it's critical to define the array size as a named constant (e.g., `MaxTaskTracks`) and ensure that index access is bounded. Even if logically "it shouldn't happen", explicit bounds or constants prevent panic risks and improve maintainability.
**Action:** Always define named constants for array sizes and verify index bounds when optimizing with stack-allocated arrays.
