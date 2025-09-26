# Tests Directory

This directory contains all test files for the PhD Dissertation Planner project.

## Directory Structure

```
tests/
├── unit/           # Unit tests for individual components
├── integration/    # Integration tests for component interactions
├── run_tests.sh    # Test runner script
└── README.md       # This file
```

## Running Tests

### Using the Test Runner Script
```bash
./tests/run_tests.sh
```

### Using Make
```bash
# Run all tests
make -f scripts/Makefile test

# Run only unit tests
make -f scripts/Makefile test-unit

# Run only integration tests
make -f scripts/Makefile test-integration
```

### Using Go Directly
```bash
# Run unit tests
go test -v ./tests/unit/...

# Run integration tests
go test -v ./tests/integration/...

# Run all tests
go test -v ./tests/... ./src/...
```

## Test Organization

### Unit Tests (`tests/unit/`)
- **`reader_test.go`** - Tests for CSV data reading and parsing
- **`validation_test.go`** - Tests for data validation logic

### Integration Tests (`tests/integration/`)
- Currently empty - add integration tests here as needed
- Examples: end-to-end PDF generation, configuration loading, etc.

## Adding New Tests

1. **Unit Tests**: Add new test files to `tests/unit/`
2. **Integration Tests**: Add new test files to `tests/integration/`
3. **Test Files**: Use `*_test.go` naming convention
4. **Package**: Use `package core_test` for external tests

## Test Guidelines

- Write tests for all public functions
- Include both positive and negative test cases
- Use descriptive test names
- Keep tests focused and independent
- Mock external dependencies when appropriate
