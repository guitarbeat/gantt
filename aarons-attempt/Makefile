PYTHON ?= python

.PHONY: build build-all clean list

# Standard build using main application
build:
	@TIMESTAMP=$$(date +"%Y%m%d_%H%M%S"); \
	$(PYTHON) main.py build single input/data.cleaned.csv \
		--title "Project Timeline" \
		--output "Timeline_$$TIMESTAMP"

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
