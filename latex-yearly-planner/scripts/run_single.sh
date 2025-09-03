#!/usr/bin/env bash
set -euo pipefail

# Convenience: run the single-task sample.
exec "$(dirname "$0")/build.sh" --csv "../input/test_single.csv" -y 2025 -p 1 -n test-single-task
#!/usr/bin/env bash

set -eo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/run_single.sh [options]

Runs the one-task demo using the sample CSV from ../input/test_single.csv.

Environment variables override defaults:
  PLANNER_YEAR         Year to render (default: 2025)
  PASSES               XeLaTeX passes (default: 1)
  CFG                  Config chain (default: base+page_template+planner_config)
  NAME                 Output PDF name (default: test-single-task)
  PLANNERGEN_BINARY    Path to compiled generator (optional)

USAGE
}

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  usage
  exit 0
fi

export PLANNER_CSV_FILE="../input/test_single.csv"
export PLANNER_YEAR="${PLANNER_YEAR:-2025}"
export PASSES="${PASSES:-1}"
export CFG="${CFG:-configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml}"
export NAME="${NAME:-test-single-task}"

scripts/build.sh
