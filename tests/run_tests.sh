#!/bin/bash
# Test runner script for PhD Dissertation Planner

set -e

echo "🧪 Running PhD Dissertation Planner Tests"
echo "========================================"

# Change to project root
cd "$(dirname "$0")/.."

# Run unit tests
echo "📋 Running unit tests..."
go test -v ./tests/unit/...

# Run integration tests (if any)
if [ -d "tests/integration" ] && [ "$(ls -A tests/integration)" ]; then
    echo "🔗 Running integration tests..."
    go test -v ./tests/integration/...
fi

# Run all tests in src/ directory
echo "📦 Running source code tests..."
go test -v ./src/...

echo "✅ All tests completed successfully!"
