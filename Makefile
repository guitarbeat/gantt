.PHONY: build run validate clean help

# Binary name
BINARY=plannergen

# Build the application
build:
	@echo "🔨 Building $(BINARY)..."
	@go build -mod=mod -o $(BINARY) main.go
	@echo "✅ Build complete!"

# Build and run
run: build
	@echo "🚀 Running planner..."
	@./$(BINARY)

# Validate CSV files
validate: build
	@echo "🔍 Validating CSV files..."
	@./$(BINARY) --validate

# Clean build artifacts and output
clean:
	@echo "🧹 Cleaning..."
	@rm -f $(BINARY)
	@rm -rf output_data/*
	@echo "✅ Clean complete!"

# Clean everything including vendor
clean-all: clean
	@echo "🧹 Deep cleaning..."
	@rm -rf vendor/
	@go clean -cache
	@echo "✅ Deep clean complete!"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod vendor
	@echo "✅ Dependencies installed!"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -mod=mod ./...

# Show help
help:
	@echo "PhD Dissertation Planner - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build      - Build the plannergen binary"
	@echo "  run        - Build and run the planner"
	@echo "  validate   - Build and validate CSV files"
	@echo "  clean      - Remove binary and output files"
	@echo "  clean-all  - Remove binary, output, and vendor"
	@echo "  deps       - Download and vendor dependencies"
	@echo "  test       - Run tests"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run"
	@echo "  make validate"

# Default target
.DEFAULT_GOAL := help
