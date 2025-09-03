# Simple Makefile for latex-yearly-planner

GO ?= go
BINARY ?= build/plannergen
CFG ?= configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml

.PHONY: build run clean preview fmt vet test test-single run-csv

build:
	$(GO) build -o $(BINARY) ./cmd/plannergen

run:
	CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/build.sh

preview:
	CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" PREVIEW=1 scripts/build.sh

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test-single:
	CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/run_single.sh

run-csv:
	CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/run_with_csv.sh

clean:
	rm -rf build/*.pdf build/*.aux build/*.log build/*.out build/*.tex
	rm -f $(BINARY)
	rm -rf build/plannergen

