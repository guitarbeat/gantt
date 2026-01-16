## 2026-01-16 - Consistent CLI Discovery Feedback
**Learning:** CLI "Auto-detected" messages used plain text while operational logs used colored icons (`core.Info`). This inconsistency made discovery steps feel less integrated and easier to miss.
**Action:** Use `core.Info` (Blue) and the `🔍` icon for all auto-detection/discovery messages to visual distinguish them from operational checks (Green/Red).
