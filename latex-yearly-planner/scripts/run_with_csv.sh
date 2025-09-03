#!/usr/bin/env bash
set -euo pipefail

# Run the planner using a CSV file of tasks.
# Usage: scripts/run_with_csv.sh [CSV_PATH] [YEAR] [NAME]

CSV_PATH="${1:-}"
YEAR="${2:-}"
NAME="${3:-}"

if [[ -z "$CSV_PATH" ]]; then
  if [[ -f "../input/data.cleaned.csv" ]]; then
    CSV_PATH="../input/data.cleaned.csv"
  elif [[ -f "../input/data.csv" ]]; then
    CSV_PATH="../input/data.csv"
  else
    echo "CSV_PATH not provided and no default found at ../input/data.cleaned.csv or ../input/data.csv" >&2
    exit 2
  fi
fi

exec "$(dirname "$0")/build.sh" --csv "$CSV_PATH" ${YEAR:+-y "$YEAR"} ${NAME:+-n "$NAME"}
#!/usr/bin/env bash

set -eo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/run_with_csv.sh [options]

Runs the planner against CSV data. Set PLANNER_CSV_FILE to point to your data.

Environment variables:
  PLANNER_CSV_FILE     CSV file path (default: examples/sample_project_data.csv)
  PLANNER_YEAR         Base year (default: 2025)
  PASSES               XeLaTeX passes (default: 1)
  CFG                  Config chain (default: base+page_template+planner_config)
  NAME                 Output PDF name (default: adaptive-planner)
  PLANNERGEN_BINARY    Path to compiled generator (optional)

USAGE
}

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  usage
  exit 0
fi

export PLANNER_CSV_FILE="${PLANNER_CSV_FILE:-examples/sample_project_data.csv}"
export PLANNER_YEAR="${PLANNER_YEAR:-2025}"
export PASSES="${PASSES:-1}"
export CFG="${CFG:-configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml}"
export NAME="${NAME:-adaptive-planner}"

scripts/build.sh
