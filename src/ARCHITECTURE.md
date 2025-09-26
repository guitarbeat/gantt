# PhD Dissertation Planner - Clean Architecture

This document explains the clean architecture implementation of the PhD Dissertation Planner system.

## Overview

The system has been reorganized to follow clean architecture principles with clear separation of concerns. This makes the codebase more maintainable, testable, and easier to understand.

## Architecture Layers

### 1. Domain Layer (`src/internal/domain/`)

The core business logic and entities. This layer contains no dependencies on external frameworks or infrastructure.

**Files:**
- `task.go` - Task entity and business rules
- `calendar.go` - Calendar domain logic and entities
- `config.go` - Configuration domain models

**Key Entities:**
- `Task` - Represents a single task with business logic
- `Calendar` - Represents a calendar with months, weeks, and days
- `Config` - Application configuration
- `TaskCollection` - Collection of tasks with efficient access patterns

### 2. Use Cases Layer (`src/internal/usecases/`)

Contains application-specific business logic that orchestrates the domain entities.

**Directories:**
- `calendar/` - Calendar-related use cases
- `task/` - Task-related use cases  
- `rendering/` - Rendering and template use cases

**Key Use Cases:**
- `CalendarUseCase` - Handles calendar generation and task assignment
- `RenderingUseCase` - Handles LaTeX document generation

### 3. Interface Adapters Layer (`src/internal/adapters/`)

Implements the interfaces defined in the interfaces layer.

**Directories:**
- `repositories/` - Repository implementations (CSV, database, etc.)
- `services/` - Service implementations

**Key Adapters:**
- `CSVTaskRepository` - Implements task data access from CSV files
- `CalendarServiceImpl` - Implements calendar service operations

### 4. Interface Definitions (`src/internal/interfaces/`)

Defines the contracts that external layers must implement.

**Directories:**
- `repositories/` - Repository interfaces
- `services/` - Service interfaces

**Key Interfaces:**
- `TaskRepository` - Contract for task data access
- `CalendarService` - Contract for calendar operations

### 5. Infrastructure Layer (`src/internal/infrastructure/`)

Handles external concerns like file I/O, logging, and external services.

**Directories:**
- `csv/` - CSV file handling utilities
- `filesystem/` - File system operations
- `latex/` - LaTeX generation utilities
- `logging/` - Logging infrastructure

### 6. Application Layer (`src/internal/app/`)

Orchestrates the use cases and coordinates the overall application flow.

**Files:**
- `planner_app.go` - Main application orchestrator

### 7. Entry Points (`src/cmd/` and `src/app/`)

**Directories:**
- `cmd/planner/` - Command-line entry point
- `app/cli/` - CLI interface implementation

### 8. Public Packages (`src/pkg/`)

Reusable libraries that can be used by other projects.

**Directories:**
- `calendar/` - Calendar utilities (public)
- `templates/` - Template system (public)
- `utils/` - Utility functions (public)

### 9. Shared Resources (`src/shared/`)

**Directories:**
- `constants/` - Application constants
- `templates/` - Template files

## Key Benefits

### 1. **Separation of Concerns**
- Each layer has a single responsibility
- Business logic is isolated from infrastructure concerns
- Easy to understand what each component does

### 2. **Testability**
- Domain logic can be tested without external dependencies
- Use cases can be tested with mock repositories and services
- Clear interfaces make mocking straightforward

### 3. **Maintainability**
- Changes to one layer don't affect others
- Easy to swap implementations (e.g., CSV to database)
- Clear boundaries make refactoring safer

### 4. **Extensibility**
- New features can be added without changing existing code
- New data sources can be added by implementing interfaces
- New rendering formats can be added as new use cases

## Data Flow

```
CLI → App → Use Cases → Domain Entities
  ↓      ↓       ↓
Infrastructure ← Interfaces ← Repositories
```

1. **CLI** receives user input and parses configuration
2. **App** orchestrates the use cases
3. **Use Cases** contain business logic and coordinate domain entities
4. **Domain** contains pure business logic
5. **Repositories** handle data persistence
6. **Infrastructure** handles external concerns

## Example Usage

```go
// Create configuration
config, err := domain.NewConfig("config.yaml")

// Create application
app := app.NewPlannerApp(config)

// Generate calendar
err = app.GenerateCalendar()
```

## Migration from Old Architecture

The old architecture had mixed concerns:
- Business logic mixed with CLI logic
- Calendar logic mixed with LaTeX rendering
- Tight coupling between components

The new architecture:
- Separates concerns into distinct layers
- Uses dependency injection for loose coupling
- Makes testing and maintenance much easier

## Future Enhancements

With this clean architecture, it's easy to add:
- Database support (implement TaskRepository interface)
- Web interface (new CLI implementation)
- Different output formats (new rendering use cases)
- Caching layer (new service implementations)
- API endpoints (new application layer)

## Testing Strategy

Each layer can be tested independently:
- **Domain**: Unit tests with no dependencies
- **Use Cases**: Unit tests with mock repositories/services
- **Adapters**: Integration tests with real data sources
- **App**: End-to-end tests with full system

This architecture provides a solid foundation for a maintainable and extensible PhD dissertation planning system.
