# Templates Overview

This project embeds its LaTeX templates and loads them from an embedded filesystem at runtime.
During development you can override this behavior to load from disk for fast iteration.

## Structure

Templates are flattened under `templates/monthly/`:

- `main_document.tpl` — Preamble and document wrapper; includes macros and pages.
- `macros.tpl` — Macros and length definitions shared by monthly templates.
- `calendar_table.tpl` — The monthly table (tabularx) structure.
- `monthly_body.tpl` — Monthly body (month grid + single full-width Notes area).
- `page_header.tpl` — Header rendering for the monthly page.
- `monthly_page.tpl` — Assembles monthly page sections.

Embedded FS is declared in `templates/embed.go`.

## Loading behavior

By default, the generator loads templates from the embedded FS for reproducible builds.
Set `DEV_TEMPLATES=1` to force loading from the local `templates/monthly/` directory instead.

Example:

```zsh
DEV_TEMPLATES=1 make preview
```

This will use templates from disk and run a preview build (unique pages only when supported).

## Debugging tips

- Use `make preview` to iterate quickly; pair with `DEV_TEMPLATES=1` while editing `.tpl` files.
- Generated LaTeX sources are written to `build/*.tex` (e.g., `build/planner_config.tex`).
- XeLaTeX warnings like Overfull/Underfull boxes can often be ignored; focus on fatal errors.

### Layout specifics

- Monthly grid uses `tabularx` for full width and `X` columns.
- Day cell content is composed of:
  - A compact corner day number (mini-tabular) overlaid on the left.
  - A right-side `minipage` with `\raggedright` content to avoid leaking table tokens.
- Notes area relies on vertical leaders of vboxed `\hrule`s with explicit `\vskip` glue, avoiding problematic dotted modes.
