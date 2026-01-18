## 2026-01-18 - Quadratic Task Stacking
**Learning:** `calendar.go` used an O(N^3) algorithm for assigning tracks to tasks, involving repeated linear searches over `d.Tasks` inside nested loops.
**Action:** Use intermediate maps (e.g., `map[int][]*Task` keyed by track) to optimize collision checks, avoiding repeated linear scans and reducing complexity to O(N^2).
