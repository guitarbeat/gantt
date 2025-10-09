# PhD Dissertation Planner - Makefile
# Common development tasks

.PHONY: help build test clean install lint fmt run organize status

# Default target
help:
	@echo "PhD Dissertation Planner - Available Commands:"
	@echo ""
	@echo "  make build       - Build the planner binary"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make install     - Install dependencies"
	@echo "  make lint        - Run linters"
	@echo "  make fmt         - Format code"
	@echo "  make run         - Build and run with default config"
	@echo "  make hooks       - Install pre-commit hooks"
	@echo "  make organize    - Clean up and organize project files"
	@echo "  make status      - Show project organization status"
	@echo ""

# Build the binary
build:
	@echo "🔨 Building plannergen..."
	go build -o plannergen.exe ./cmd/planner
	@echo "✅ Build complete!"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...
	@echo "✅ Tests complete!"

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f plannergen.exe plannergen
	rm -f coverage.txt coverage.html
	rm -rf generated/
	@echo "✅ Clean complete!"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

# Run linters
lint:
	@echo "🔍 Running linters..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
	@echo "✅ Lint complete!"

# Run basic linters (without golangci-lint)
lint-basic:
	@echo "🔍 Running basic linters..."
	go vet ./...
	gofmt -l .
	@echo "✅ Basic lint complete!"

# Format code
fmt:
	@echo "✨ Formatting code..."
	gofmt -w .
	@echo "✅ Format complete!"

# Build and run with default config
run: build
	@echo "🚀 Running planner..."
	./plannergen.exe

# Install pre-commit hooks
hooks:
	@echo "🪝 Installing pre-commit hooks..."
	pre-commit install
	@echo "✅ Hooks installed!"

# Run pre-commit on all files
check:
	@echo "✅ Running pre-commit checks..."
	pre-commit run --all-files

# Organize project files
organize:
	@echo "🧹 Organizing project files..."
	./scripts/cleanup_and_organize.sh
	@echo "✅ Organization complete!"

# Show project status
status:
	@echo "📊 Project Status:"
	./scripts/cleanup_and_organize.sh --status

# Include build configuration
-include .config/build/makefile.conf
