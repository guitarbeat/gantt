Project: Structured Tasks Repository

Overview
- Purpose: Normalize planning data into CSVs with unique task IDs, dependencies, and quarter groupings (Fall 2025–Summer 2027).
- Inputs: `repo/input/data.csv` (source). If starting with `data.csv` at workspace root, move it to `repo/input/data.csv` or pass `--input`.
- Outputs:
  - `repo/data/tasks.csv`: canonical tasks table
  - `repo/data/dependencies.csv`: task-to-task dependency edges
  - `repo/data/quarters.csv`: quarter definitions (Fall 2025 → Summer 2027)
  - `repo/data/task_quarters.csv`: task-to-quarter overlaps
  - `repo/reports/quarters/*.csv`: per-quarter task listings

Schema
- tasks.csv (columns):
  - task_id, name, parent_id, deliverable_type, group, owner, status, priority, start_date, due_date, duration_days, notes
- dependencies.csv (columns):
  - task_id, depends_on
- quarters.csv (columns):
  - quarter_id, label, start_date, end_date
- task_quarters.csv (columns):
  - task_id, quarter_id, overlap_days

Quarters
- 2025-FALL: 2025-08-01 → 2025-12-31
- 2026-WINTER: 2026-01-01 → 2026-03-31
- 2026-SPRING: 2026-04-01 → 2026-06-30
- 2026-SUMMER: 2026-07-01 → 2026-08-31
- 2026-FALL: 2026-09-01 → 2026-12-31
- 2027-WINTER: 2027-01-01 → 2027-03-31
- 2027-SPRING: 2027-04-01 → 2027-06-30
- 2027-SUMMER: 2027-07-01 → 2027-08-31

CLI Usage
- Import from CSV:
  - `python repo/scripts/build.py import --input data.csv --outdir repo/data`
- Validate referential integrity and dates:
  - `python repo/scripts/build.py validate --datadir repo/data`
- Generate per-quarter reports and task-to-quarter map:
  - `python repo/scripts/build.py quarters --datadir repo/data --reports repo/reports/quarters`

Workflow
1) Place your source CSV at `repo/input/data.csv` (or pass `--input`).
2) Run `import` to generate normalized CSVs.
3) Run `validate` to catch any issues.
4) Run `quarters` to create quarter-indexed reports.
5) When ready for a clean deliverables-only folder, remove everything except PDFs and the DOCX.

