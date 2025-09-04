PYTHON ?= python

.PHONY: build build-pdf clean help

# Standard build using main application
build:
	@TIMESTAMP=$$(date +"%Y%m%d_%H%M%S"); \
	$(PYTHON) main.py \
		--input ../input/data.cleaned.csv \
		--title "Project Timeline" \
		--output "output/tex/Timeline_$$TIMESTAMP.tex"

# Build and compile to PDF
build-pdf:
	@TIMESTAMP=$$(date +"%Y%m%d_%H%M%S"); \
	$(PYTHON) main.py \
		--input ../input/data.cleaned.csv \
		--title "Project Timeline" \
		--output "output/tex/Timeline_$$TIMESTAMP.tex" && \
	cd output/tex && pdflatex Timeline_$$TIMESTAMP.tex

# Clean build artifacts
clean:
	rm -f output/tex/*.tex output/tex/*.aux output/tex/*.log output/tex/*.out output/tex/*.pdf

# Help target
help:
	@echo "Available targets:"
	@echo "  build           - Build standard LaTeX timeline"
	@echo "  build-pdf       - Build LaTeX and compile to PDF"
	@echo "  clean           - Clean build artifacts"
	@echo "  help            - Show this help message"
