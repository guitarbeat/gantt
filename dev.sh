#!/bin/bash
# Development script for PhD Dissertation Planner

# Load development environment
if [ -f ".env.dev" ]; then
    source .env.dev
fi

# Set development environment variables
export PLANNER_CSV_FILE="${DEV_PLANNER_CSV_FILE:-input_data/research_timeline_v5.1_comprehensive.csv}"
export PLANNER_CONFIG_FILE="${DEV_PLANNER_CONFIG_FILE:-configs/base.yaml}"
export PLANNER_OUTPUT_DIR="${DEV_PLANNER_OUTPUT_DIR:-generated}"
export PLANNER_SILENT="${DEV_PLANNER_SILENT:-0}"
export PLANNER_DEBUG="${DEV_PLANNER_DEBUG:-1}"
export PLANNER_VERBOSE="${DEV_PLANNER_VERBOSE:-1}"

# Run the command
exec "$@"
