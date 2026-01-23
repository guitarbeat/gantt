## 2026-01-14 - CLI Log Noise vs UI Output
**Learning:** Default loggers writing to `stderr` often interrupt formatted CLI UI output on `stdout`, creating a messy user experience.
**Action:** In CLI tools with rich UI, default logging to `Discard` or a file, or use `Debug` level for operational logs to keep the standard output clean.

## 2026-01-14 - Go Template Receiver Types
**Learning:** `text/template` fails to call methods defined on pointer receivers if the data passed is a value type.
**Action:** Ensure data passed to templates is a pointer, or access struct fields directly if public, to avoid runtime template execution errors.

## 2026-01-23 - Actionable Error Messages
**Learning:** When a tool dependency is missing (like `xelatex`), treating it as a hard failure (❌) discourages users. Presenting it as a warning (⚠️) with a high-contrast copy-pasteable command improves the recovery experience.
**Action:** Distinguish between fatal errors and "missing optional dependency" states, and format manual recovery commands to stand out visually.
