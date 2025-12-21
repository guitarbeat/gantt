## 2025-12-21 - [CLI Async Feedback]
**Learning:** For long-running blocking operations in CLI tools (like LaTeX compilation), users need immediate visual feedback. A simple ASCII spinner provides assurance the process hasn't hung, without cluttering the final output logs.
**Action:** Use a goroutine-based spinner that cleans up its own output (using backspaces) before the final success/failure message is printed. This maintains the clean "Step... ✅" log format while providing interactivity.
