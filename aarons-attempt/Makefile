PYTHON ?= python

.PHONY: build build-all clean list

# Standard build using main application
build:
	@TIMESTAMP=$$(date +"%Y%m%d_%H%M%S"); \
	$(PYTHON) main.py --input input/data.cleaned.csv --output output/tex/Calendar_$$TIMESTAMP.tex; \
	echo "Generated output/tex/Calendar_$$TIMESTAMP.tex"; \
	pdflatex -interaction=nonstopmode -output-directory output/tex output/tex/Calendar_$$TIMESTAMP.tex | cat; \
	mv output/tex/Calendar_$$TIMESTAMP.pdf output/pdf/; \
	echo "Built output/pdf/Calendar_$$TIMESTAMP.pdf"; \
	rm -f output/tex/Calendar_$$TIMESTAMP.aux output/tex/Calendar_$$TIMESTAMP.log output/tex/Calendar_$$TIMESTAMP.toc; \
	echo "Cleaned up build artifacts"

# Build all templates using enhanced build system
build-all:
	$(PYTHON) main.py build all-templates input/data.cleaned.csv --title "Project Timeline"

# Build for specific device
build-device:
	$(PYTHON) main.py build single input/data.cleaned.csv -d supernote_a5x --title "E-ink Planner"

# List available configurations
list:
	$(PYTHON) main.py build list

# Clean build artifacts
clean:
	rm -f output/pdf/*.pdf output/tex/*.tex output/tex/*.aux output/tex/*.log output/tex/*.out output/tex/*.toc
	rm -rf build/ temp/
