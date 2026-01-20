## 2024-05-24 - Optimization of Task Track Assignment
**Learning:** Checking for task overlaps using a nested loop that re-scans the full task list (O(N^3)) is a major bottleneck. Grouping tasks by track in a temporary map (`map[int][]*Task`) reduces complexity to ~O(N^2) and eliminates redundant lookups.
**Action:** When optimizing "bin packing" or "track assignment" logic, always maintain a lookup index of the dimension you are checking against (e.g., track number) to avoid linear scans.
