#!/usr/bin/env bash
# Clean output directory script

echo "🧹 Cleaning output directory..."

# Remove all generated files but keep directory structure
rm -f output/pdfs/*.pdf
rm -f output/latex/*.tex
rm -f output/logs/*.log

# Keep README.md
echo "✅ Output directory cleaned"
echo "📁 Directory structure preserved"
echo "📋 README.md kept"
