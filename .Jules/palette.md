## 2025-05-23 - CLI Feedback Patterns
**Learning:** Users respond well to "magic" auto-detection when it's visually highlighted. Use `core.Info` with specific icons (🔍) to distinguish discovery steps from standard operational logs.
**Action:** When implementing auto-configuration or discovery logic, always wrap the success message in `core.Info` and include a relevant emoji/icon to make it feel like a helpful feature rather than a debug log.
