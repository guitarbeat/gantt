#!/usr/bin/env bash
set -euo pipefail

# Release builder for generating timestamped planner PDFs.
# Creates clean builds with timestamped filenames in the release/ directory.
#
# Usage:
#   scripts/build_release.sh [-c CFG_CHAIN] [-n NAME] [-y YEAR] [--csv CSV_PATH]
#
# Examples:
#   scripts/build_release.sh -c "configs/base.yaml,configs/page_template.yaml,configs/csv_config.yaml" -n "overlap_test"

CFG_DEFAULT="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml"

CFG="${CFG:-$CFG_DEFAULT}"
PLANNERGEN_BINARY="${PLANNERGEN_BINARY:-build/plannergen}"
NAME="${NAME:-planner}"
PLANNER_YEAR="${PLANNER_YEAR:-}"
CSV_PATH="${PLANNER_CSV_FILE:-}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    -c|--cfg)
      CFG="$2"; shift 2 ;;
    -n|--name)
      NAME="$2"; shift 2 ;;
    -y|--year)
      PLANNER_YEAR="$2"; shift 2 ;;
    --csv)
      CSV_PATH="$2"; shift 2 ;;
    -h|--help)
      echo "Usage: $0 [-c CFG_CHAIN] [-n NAME] [-y YEAR] [--csv CSV_PATH]"; exit 0 ;;
    *)
      echo "Unknown option: $1" >&2; exit 2 ;;
  esac
done

# Clean build directory first
echo "Cleaning build artifacts..."
rm -f build/*.pdf build/*.tex build/*.aux build/*.log build/*.out build/*.synctex.gz *.pdf

# Create timestamp
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
RELEASE_NAME="${NAME}_${TIMESTAMP}"

echo "Building release: ${RELEASE_NAME}"
echo "Using plannergen binary at \"$PLANNERGEN_BINARY\""

# Set up environment
export CFG
export PLANNERGEN_BINARY
export OUTDIR="build"
if [[ -n "${PLANNER_YEAR}" ]]; then export PLANNER_YEAR; fi
if [[ -n "${CSV_PATH}" ]]; then export PLANNER_CSV_FILE="$CSV_PATH"; fi

# Build the planner
./scripts/single.sh

# Move the generated PDF to release directory with timestamp
if [[ -f "build/csv_config.pdf" ]]; then
    cp "build/csv_config.pdf" "release/${RELEASE_NAME}.pdf"
    echo "✅ Release created: release/${RELEASE_NAME}.pdf"
elif [[ -f "build/planner_config.pdf" ]]; then
    cp "build/planner_config.pdf" "release/${RELEASE_NAME}.pdf"
    echo "✅ Release created: release/${RELEASE_NAME}.pdf"
else
    echo "❌ No PDF found in build directory"
    exit 1
fi

# Clean up build artifacts again
rm -f build/*.pdf build/*.tex build/*.aux build/*.log build/*.out build/*.synctex.gz

echo "🎉 Release complete! Open with: open release/${RELEASE_NAME}.pdf"