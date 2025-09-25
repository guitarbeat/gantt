# Directory Organization Plan
## Task 4.4 - Directory Cleanup & Organization

### Current State Analysis
- **Total Files:** 149 files
- **Go Files:** 124 files (83.2%)
- **Test Files:** 18 files scattered in root
- **Backup Files:** 4 files (to be removed)
- **Temp Files:** 2 files (to be removed)
- **Test Output:** 3 directories with generated files

### Proposed Clean Directory Structure

```
latex-yearly-planner/
├── cmd/                                    # Application entry points
│   └── plannergen/
│       └── main.go
├── internal/                               # Core application logic
│   ├── app/
│   │   └── app.go
│   ├── calendar/                           # Calendar functionality
│   │   ├── *.go (core files)
│   │   └── *_test.go (test files)
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── data/
│   │   ├── *.go (core files)
│   │   └── *_test.go (test files)
│   ├── generator/                          # PDF generation
│   │   ├── *.go (core files)
│   │   └── *_test.go (test files)
│   ├── header/
│   │   └── *.go
│   ├── latex/
│   │   └── *.go
│   └── layout/
│       └── lengths.go
├── configs/                                # Configuration files
│   ├── base.yaml
│   ├── csv_config.yaml
│   ├── page_template.yaml
│   ├── planner_config.yaml
│   └── README.md
├── templates/                              # Template files
│   ├── embed.go
│   ├── monthly/
│   │   └── *.tpl
│   └── README.md
├── scripts/                                # Build and utility scripts
│   ├── README.md
│   └── simple.sh
├── tests/                                  # NEW: Organized test files
│   ├── integration/
│   │   ├── test_feedback_integration.go
│   │   ├── test_performance_integration.go
│   │   └── test_visual_integration.go
│   ├── quality/
│   │   ├── test_quality_assurance.go
│   │   ├── test_quality_issue_resolver.go
│   │   ├── test_quality_system.go
│   │   └── test_quality_validator.go
│   ├── validation/
│   │   ├── test_user_coordination.go
│   │   ├── test_user_validation.go
│   │   ├── test_final_approval.go
│   │   └── test_final_assessment.go
│   ├── performance/
│   │   ├── test_performance_optimization.go
│   │   ├── test_performance_simple.go
│   │   └── test_visual_spacing.go
│   ├── feedback/
│   │   ├── test_feedback_system.go
│   │   └── test_improvement_logic.go
│   └── parser/
│       ├── test_enhanced_parser.go
│       ├── test_enhanced_visual.go
│       ├── test_multi_format.go
│       └── test_parser.go
├── docs/                                   # NEW: Consolidated documentation
│   ├── README.md                           # Main project documentation
│   ├── reports/
│   │   ├── TASK_3_3_COMPLETION_REPORT.md
│   │   ├── TASK_3_4_COMPLETION_REPORT.md
│   │   └── PERFORMANCE_OPTIMIZATION_REPORT.md
│   ├── lessons/
│   │   └── LESSONS_LEARNED_FROM_AARONS_ATTEMPT.md
│   └── api/                                # API documentation
│       └── README.md
├── examples/                               # NEW: Example configurations
│   ├── sample_batch_config.json
│   └── README.md
├── build/                                  # NEW: Build artifacts (gitignored)
│   └── .gitkeep
├── dist/                                   # NEW: Distribution files (gitignored)
│   └── .gitkeep
├── go.mod
├── go.sum
├── Makefile
└── .gitignore                              # Updated to exclude build artifacts
```

### File Organization Actions

#### 1. Files to Remove (Cleanup)
**Backup Files (4 files):**
- `internal/calendar/integration_test.go.bak`
- `internal/calendar/complex_overlap_test.go.bak`
- `internal/calendar/integration_test.go.temp.bak`
- `internal/calendar/complex_overlap_test.go.temp.bak`

**Temporary Files (2 files):**
- `internal/calendar/complex_overlap_test.go.temp`
- `internal/calendar/integration_test.go.temp`

**Temporary Test Files (5 files):**
- `simple_quality_test.go`
- `simple_quality_validation.go`
- `simple_test.go`
- `simple_validation.go`
- `simple_visual_test.go`
- `simple_visual_validation.go`
- `validate_integration.go`

**Test Output Directories (3 directories):**
- `test_output/` (remove contents, keep directory)
- `test_single_output/` (remove contents, keep directory)
- `test_triple_output/` (remove contents, keep directory)

#### 2. Files to Move (Reorganization)

**Test Files to Move to `tests/` directory:**

**Integration Tests:**
- `test_feedback_integration.go` → `tests/integration/`
- `test_performance_integration.go` → `tests/integration/`
- `test_visual_integration.go` → `tests/integration/`

**Quality Tests:**
- `test_quality_assurance.go` → `tests/quality/`
- `test_quality_issue_resolver.go` → `tests/quality/`
- `test_quality_system.go` → `tests/quality/`
- `test_quality_validator.go` → `tests/quality/`

**Validation Tests:**
- `test_user_coordination.go` → `tests/validation/`
- `test_user_validation.go` → `tests/validation/`
- `test_final_approval.go` → `tests/validation/`
- `test_final_assessment.go` → `tests/validation/`

**Performance Tests:**
- `test_performance_optimization.go` → `tests/performance/`
- `test_performance_simple.go` → `tests/performance/`
- `test_visual_spacing.go` → `tests/performance/`

**Feedback Tests:**
- `test_feedback_system.go` → `tests/feedback/`
- `test_improvement_logic.go` → `tests/feedback/`

**Parser Tests:**
- `test_enhanced_parser.go` → `tests/parser/`
- `test_enhanced_visual.go` → `tests/parser/`
- `test_multi_format.go` → `tests/parser/`
- `test_parser.go` → `tests/parser/`

**Documentation Files to Move to `docs/` directory:**

**Reports:**
- `TASK_3_3_COMPLETION_REPORT.md` → `docs/reports/`
- `TASK_3_4_COMPLETION_REPORT.md` → `docs/reports/`
- `PERFORMANCE_OPTIMIZATION_REPORT.md` → `docs/reports/`

**Lessons:**
- `LESSONS_LEARNED_FROM_AARONS_ATTEMPT.md` → `docs/lessons/`

**Example Files to Move to `examples/` directory:**
- `sample_batch_config.json` → `examples/`

#### 3. Files to Keep in Place
**Core Application Structure:**
- `cmd/` - Application entry points
- `internal/` - Core application logic (with test files in place)
- `configs/` - Configuration files
- `templates/` - Template files
- `scripts/` - Build and utility scripts
- `go.mod`, `go.sum` - Go module files
- `Makefile` - Build configuration

#### 4. New Directories to Create
- `tests/` - Organized test files
- `docs/` - Consolidated documentation
- `examples/` - Example configurations
- `build/` - Build artifacts (gitignored)
- `dist/` - Distribution files (gitignored)

### Import Path Updates Required

#### Test Files Moving to `tests/` Directory
All test files moving to `tests/` subdirectories will need import path updates:

**Current imports in test files:**
```go
import "latex-yearly-planner/internal/..."
```

**Updated imports for moved test files:**
```go
import "../../internal/..."
```

**Specific files requiring import updates:**
- All files in `tests/integration/`
- All files in `tests/quality/`
- All files in `tests/validation/`
- All files in `tests/performance/`
- All files in `tests/feedback/`
- All files in `tests/parser/`

### .gitignore Updates

Add the following entries to `.gitignore`:

```gitignore
# Build artifacts
build/
dist/

# Test output directories
test_output/
test_single_output/
test_triple_output/

# Temporary files
*.bak
*.temp
*.tmp
*~

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db
```

### Benefits of This Organization

1. **Clean Root Directory:** Only essential files in root
2. **Organized Tests:** Tests grouped by functionality
3. **Consolidated Documentation:** All docs in one place
4. **Clear Separation:** Core code, tests, docs, examples separated
5. **Professional Structure:** Industry-standard Go project layout
6. **Maintainable:** Easy to find and manage files
7. **Scalable:** Structure supports future growth

### Implementation Steps

1. **Create new directories**
2. **Move test files to organized structure**
3. **Move documentation files**
4. **Move example files**
5. **Remove temporary and backup files**
6. **Update import paths in moved files**
7. **Update .gitignore**
8. **Clean test output directories**
9. **Validate all imports work correctly**

### Validation Checklist

- [ ] All test files moved to appropriate `tests/` subdirectories
- [ ] All documentation moved to `docs/` directory
- [ ] All example files moved to `examples/` directory
- [ ] All temporary and backup files removed
- [ ] All import paths updated and working
- [ ] Test output directories cleaned
- [ ] .gitignore updated
- [ ] All Go files compile successfully
- [ ] All tests run successfully
- [ ] Root directory contains only essential files

This organization plan will result in a clean, professional, and maintainable project structure that follows Go best practices and industry standards.
