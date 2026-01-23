## 2024-05-22 - Date Normalization Optimization
**Learning:** Repeatedly normalizing dates (creating new `time.Time` values via `time.Date`) in hot loops (like rendering calendars or stacking tasks) is a significant performance bottleneck. Normalizing once at object creation allows direct field access, reducing allocations and CPU time.
**Action:** Ensure `SpanningTask` dates are normalized to UTC midnight at creation time and use direct field access in rendering/calculation loops.
