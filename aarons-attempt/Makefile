PYTHON ?= python

.PHONY: build build-all build-formats build-device clean list help

# Standard build using main application
build:
	@TIMESTAMP=$$(date +"%Y%m%d_%H%M%S"); \
	$(PYTHON) main.py build single input/data.cleaned.csv \
		--title "Project Timeline" \
		--output "Timeline_$$TIMESTAMP"

# Build all templates using enhanced build system
build-all:
	$(PYTHON) main.py build all-templates ../input/data.cleaned.csv --title "Project Timeline"

# Build in multiple export formats (PDF, SVG, HTML, PNG)
build-formats:
	$(PYTHON) main.py build multiple-formats ../input/data.cleaned.csv \
		--title "Project Timeline" \
		--formats pdf svg html png

# Build for specific device
build-device:
	$(PYTHON) main.py build single ../input/data.cleaned.csv -d supernote_a5x --title "E-ink Planner"

# Build with enhanced template
build-enhanced:
	$(PYTHON) main.py build single ../input/data.cleaned.csv \
		--title "Enhanced Project Timeline" \
		--template gantt_timeline \
		--color-scheme academic

# Build for all devices
build-all-devices:
	$(PYTHON) main.py build all-devices ../input/data.cleaned.csv --title "Multi-Device Planner"

# List available configurations
list:
	$(PYTHON) main.py build list

# Clean build artifacts
clean:
	rm -f output/pdf/*.pdf output/tex/*.tex output/tex/*.aux output/tex/*.log output/tex/*.out output/tex/*.toc
	rm -rf build/ temp/ output/svg/ output/html/ output/png/

# Help target
help:
	@echo "Available targets:"
	@echo "  build           - Build standard PDF timeline"
	@echo "  build-all       - Build all available templates"
	@echo "  build-formats   - Build in multiple formats (PDF, SVG, HTML, PNG)"
	@echo "  build-device    - Build for specific device (e-ink tablet)"
	@echo "  build-enhanced  - Build with enhanced template and academic colors"
	@echo "  build-all-devices - Build for all available device profiles"
	@echo "  list            - List available configurations"
	@echo "  clean           - Clean build artifacts"
	@echo "  help            - Show this help message"
