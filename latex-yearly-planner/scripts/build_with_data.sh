#!/usr/bin/env bash

# Build script that uses CSV data to determine date range
PLANNER_CSV_FILE="../aarons-attempt/input/data.cleaned.csv" \
PLANNER_YEAR=2025 \
PASSES=1 \
CFG="cfg/base.yaml,cfg/template_breadcrumb.yaml,cfg/sn_a5x.breadcrumb.default.yaml" \
NAME="adaptive-planner" \
./scripts/single.sh
