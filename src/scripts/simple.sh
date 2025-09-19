#!/usr/bin/env bash
set -euo pipefail

# Simple PDF generation script - one command does everything
# Usage: ./scripts/simple.sh [csv_file] [output_name]

CSV_FILE="${1:-../input/test_single.csv}"
OUTPUT_NAME="${2:-planner}"

echo "🎯 Generating PDF from: $CSV_FILE"
echo "📄 Output: ${OUTPUT_NAME}.pdf"

# Build if needed
if [ ! -f "build/plannergen" ]; then
    echo "🔨 Building plannergen..."
    go build -o build/plannergen ./cmd/plannergen
fi

# Generate LaTeX
echo "📝 Generating LaTeX..."
PLANNER_CSV_FILE="$CSV_FILE" \
./build/plannergen --config "configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml" --outdir build

# Compile PDF
echo "📚 Compiling PDF..."
cd build
xelatex -file-line-error -interaction=nonstopmode planner_config.tex > /dev/null 2>&1 || true
cd ..

# Copy to main directory (primary output)
cp "build/planner_config.pdf" "${OUTPUT_NAME}.pdf"
echo "✅ Created: ${OUTPUT_NAME}.pdf"

# Also copy to organized structure in parent directory
mkdir -p ../output/pdfs ../output/latex ../output/logs
cp "build/planner_config.pdf" "../output/pdfs/${OUTPUT_NAME}.pdf"
cp "build/planner_config.tex" "../output/latex/${OUTPUT_NAME}.tex" 2>/dev/null || true
cp "build/planner_config.log" "../output/logs/${OUTPUT_NAME}.log" 2>/dev/null || true
echo "📁 Also saved to: ../output/pdfs/${OUTPUT_NAME}.pdf"
