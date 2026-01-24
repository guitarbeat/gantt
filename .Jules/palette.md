## 2026-01-14 - CLI Log Noise vs UI Output
**Learning:** Default loggers writing to `stderr` often interrupt formatted CLI UI output on `stdout`, creating a messy user experience.
**Action:** In CLI tools with rich UI, default logging to `Discard` or a file, or use `Debug` level for operational logs to keep the standard output clean.

## 2026-01-14 - Go Template Receiver Types
**Learning:** `text/template` fails to call methods defined on pointer receivers if the data passed is a value type.
**Action:** Ensure data passed to templates is a pointer, or access struct fields directly if public, to avoid runtime template execution errors.

## 2026-01-20 - CLI Progress Bar Artifacts
**Learning:** Overwriting CLI lines with `\r` leaves trailing characters if the new line is shorter than the previous one.
**Action:** Always append `\033[K` (Clear Line) after `\r` to cleanly erase the remaining line content before printing new text.
