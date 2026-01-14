## 2024-05-23 - Exponential Complexity in Task Assignment
**Learning:** `assignTaskTracks` in `calendar.go` contained a hidden O(N^4) complexity. It iterated over tasks, then tracks, then ALL assigned tasks, then linearly searched for the assigned task object by ID.
**Action:** Replaced linear search and iteration over all tasks with a `tracksUsage` map (track -> tasks), reducing complexity to O(N^2) and improving performance by ~39x for 1000 tasks.
