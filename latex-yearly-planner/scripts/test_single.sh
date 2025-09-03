#!/usr/bin/env bash

# Quick test script using single task example (uses central repo `input/test_single.csv`)
# Adjusted to use the existing input file at the repository root
PLANNER_CSV_FILE="../input/test_single.csv" \
PLANNER_YEAR=2025 \
PASSES=1 \
CFG="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml" \
NAME="test-single-task" \
./scripts/single.sh