## 2024-05-23 - Algorithmic Complexity in Calendar Rendering
**Learning:** The `assignTaskTracks` function had a hidden O(N³) complexity due to nested loops and repeated linear searches over slices to find tasks by ID. This became a bottleneck even with moderate task counts (N=100).
**Action:** Replace repeated linear lookups with pre-calculated maps or slices (e.g., `tracksUsage`) to reduce complexity to O(N²). Always check for hidden linear scans inside loops.
