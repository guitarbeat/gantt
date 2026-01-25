## 2024-05-23 - Date Normalization Optimization
**Learning:** `time.Date` construction is surprisingly expensive when called repeatedly inside tight loops (like `assignTaskTracks` and `findActiveTasks`).
**Action:** Normalize dates once at creation/assignment time (e.g., in `ApplySpanningTasksToMonth`) and store them in the struct, replacing dynamic `time.Date` calls with direct field access. This yielded a ~3.4x speedup in month rendering.
