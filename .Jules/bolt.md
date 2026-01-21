## 2024-05-23 - Nested Slice Lookups
**Learning:** Looking up tasks by ID within nested loops (specifically in `assignTaskTracks`) creates O(N³) complexity, causing significant performance degradation (~4ms per day for 100 tasks).
**Action:** Replace linear slice searches with map-based lookups or maintain structured data (like `tracksUsage`) to reduce complexity to O(N²).
