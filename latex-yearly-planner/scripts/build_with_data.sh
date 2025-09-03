#!/usr/bin/env bash

# Build script that uses CSV data to determine date range
PLANNER_CSV_FILE="../aarons-attempt/input/data.cleaned.csv" \
PLANNER_YEAR=2025 \
PASSES=1 \
CFG="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml" \
NAME="adaptive-planner" \
./scripts/single.sh
