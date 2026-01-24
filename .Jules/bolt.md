## 2026-01-24 - Pre-calculating Date Normalization
**Learning:** `time.Date()` construction is surprisingly expensive (~100-200ns) when called inside O(N²) nested loops. Repeatedly normalizing the same date for comparison in `findLowestAvailableTrackForTask` was a major bottleneck.
**Action:** Normalize immutable dates (like StartDate/EndDate) once at creation time or batch processing time, and store them back in the struct, allowing hot loops to use direct field access. This yielded a ~3.4x speedup in day rendering.
