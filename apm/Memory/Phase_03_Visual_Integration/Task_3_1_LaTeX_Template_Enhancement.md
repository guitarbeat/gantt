# Task 3.1 - LaTeX Template Enhancement

- Date: 2025-09-05
- Agent: Agent_VisualRendering
- Context: Implements Google Calendar-style task bars and overflow handling within existing monthly layout without changing page geometry. Integrates with existing day cell rendering and spanning task overlays.

## What changed
- Added TikZ/tcolorbox macros for task bars, compact stacked bars, and continuation chevrons in `templates/monthly/macros.tpl`.
- Refactored day-level spanning task rendering to use macros in `internal/calendar/task_rendering.go`.
- Left page structure and table geometry unchanged; only visual components updated.

## Key details
- Rounded corners, subtle borders/shadows, and category colors are supported via `\TaskOverlayBox` and `\TaskCompactBox`.
- Simple category color registry and placeholder text-contrast macro introduced. Defaults to black text for readability.
- Multi-task overflow continues to render a “+N more” indicator; supports optional footnote expansion later without layout changes.

## Files edited
- `templates/monthly/macros.tpl`: Added macros and color registry.
- `internal/calendar/task_rendering.go`: Switched overlay rendering to macros and compact bar macro.
- `templates/monthly/calendar_table.tpl`: No structural change; minor safe update retained.

## Notes
- Month boundary chevrons macros added (`\TaskContLeft`, `\TaskContRight`), ready for use when month-boundary positions are surfaced to templates.
- Existing tests should pass; visuals are encapsulated in macros without affecting measurements.
