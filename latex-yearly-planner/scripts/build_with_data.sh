#!/usr/bin/env bash

# Build script that uses CSV data to determine date range
PLANNER_CSV_FILE="../aarons-attempt/input/data.cleaned.csv" \
PLANNER_YEAR=2025 \
PASSES=1 \
CFG="config/base.yaml,config/page_template.yaml,config/planner_config.yaml" \
NAME="adaptive-planner" \
./scripts/single.sh
