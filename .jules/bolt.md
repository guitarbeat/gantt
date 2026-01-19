## 2024-05-22 - Nested Loop Complexity in Calendar Rendering
**Learning:** The `assignTaskTracks` function in `internal/calendar/calendar.go` used nested loops and repeated linear searches (`O(N^3)`), causing significant performance degradation (382ms for 1000 tasks).
**Action:** Replace nested searches with map-based lookups (`O(N^2)` worst case, `O(N)` average) when processing overlapping timeline items.
