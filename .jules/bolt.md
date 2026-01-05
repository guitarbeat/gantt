## 2024-05-23 - Task Stacker Optimization
**Learning:** In `TaskStacker`, using `time.Format("2006-01-02")` for map keys and allocating slices of `time.Time` for every date range check caused significant overhead.
**Action:** Replace date string keys with integer keys (YYYYMMDD) and replace range slice allocation with direct date iteration using `time.AddDate`. This yielded ~3.5x performance improvement (48ms -> 13.4ms).
