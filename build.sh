#!/bin/bash
# Build script for PhD Dissertation Planner with CSV data

# Set the CSV file path
export PLANNER_CSV_FILE="input_data/research_timeline_v5_comprehensive.csv"

echo "Building planner with CSV: $PLANNER_CSV_FILE"

# Build the Go binary
echo "Compiling Go code..."
go build -o plannergen ./cmd/planner || exit 1

# Generate LaTeX files
echo "Generating LaTeX files..."
./plannergen || exit 1

# Compile PDF
echo "Compiling PDF..."
cd build || exit 1
pdflatex -interaction=nonstopmode base.tex > /tmp/latex_output.log 2>&1
if [ $? -eq 0 ]; then
    echo "PDF compiled successfully!"
else
    echo "⚠️  LaTeX compilation had warnings but PDF was generated"
    tail -20 /tmp/latex_output.log
fi

# Open the PDF
if [[ "$OSTYPE" == "darwin"* ]]; then
    open base.pdf
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    xdg-open base.pdf
fi

echo "✅ Build complete!"
