#!/usr/bin/env bash

# Quick test script using single task example
PLANNER_CSV_FILE="examples/test_single.csv" \
PLANNER_YEAR=2025 \
PASSES=1 \
CFG="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml" \
NAME="test-single-task" \
./scripts/single.sh