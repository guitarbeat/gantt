## 2023-10-27 - Spanning Task Redundancy
**Learning:** Tasks that span multiple days are represented by shared pointers in `Day` structs. Iterating over `Month.Weeks.Days.Tasks` processes the same task instance multiple times (once per day it spans).
**Action:** Use a `seen` map (keyed by Task ID or Category) when aggregating data across a month to avoid O(N*Days) complexity.
