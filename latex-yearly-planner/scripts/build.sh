#!/usr/bin/env bash
set -euo pipefail

# Unified runner for generating planner PDFs.
# Wraps scripts/single.sh and passes through common options.
#
# Usage:
#   scripts/build.sh [--preview] [-c CFG_CHAIN] [-b BINARY] [-n NAME] [-y YEAR] [-p PASSES] [--csv CSV_PATH]
#
# Examples:
#   scripts/build.sh -c "configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml"
#   scripts/build.sh --preview -n demo

CFG_DEFAULT="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml"

CFG="${CFG:-$CFG_DEFAULT}"
PLANNERGEN_BINARY="${PLANNERGEN_BINARY:-build/plannergen}"
OUTDIR="${OUTDIR:-build}"
NAME="${NAME:-}"
PLANNER_YEAR="${PLANNER_YEAR:-}"
PASSES="${PASSES:-}"
CSV_PATH="${PLANNER_CSV_FILE:-}"
PREVIEW=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --preview)
      PREVIEW=1; shift ;;
    -c|--cfg)
      CFG="$2"; shift 2 ;;
    -b|--binary)
      PLANNERGEN_BINARY="$2"; shift 2 ;;
    -n|--name)
      NAME="$2"; shift 2 ;;
    -y|--year)
      PLANNER_YEAR="$2"; shift 2 ;;
    -p|--passes)
      PASSES="$2"; shift 2 ;;
    -o|--outdir)
      OUTDIR="$2"; shift 2 ;;
    --csv)
      CSV_PATH="$2"; shift 2 ;;
    -h|--help)
      echo "Usage: $0 [--preview] [-c CFG_CHAIN] [-b BINARY] [-n NAME] [-y YEAR] [-p PASSES] [--csv CSV_PATH]"; exit 0 ;;
    *)
      echo "Unknown option: $1" >&2; exit 2 ;;
  esac
done

echo "Building using plannergen binary at \"$PLANNERGEN_BINARY\""

# Export optional env vars only when set
export CFG
export PLANNERGEN_BINARY
if [[ -n "${PREVIEW}" ]]; then export PREVIEW; fi
if [[ -n "${NAME}" ]]; then export NAME; fi
if [[ -n "${PLANNER_YEAR}" ]]; then export PLANNER_YEAR; fi
if [[ -n "${PASSES}" ]]; then export PASSES; fi
if [[ -n "${CSV_PATH}" ]]; then export PLANNER_CSV_FILE="$CSV_PATH"; fi
if [[ -n "${OUTDIR}" ]]; then export OUTDIR; fi

exec ./scripts/single.sh
