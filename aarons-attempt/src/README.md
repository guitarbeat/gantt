# 📦 LaTeX Gantt Chart Generator - Consolidated Package

A clean, consolidated tool for generating publication-quality LaTeX timelines from CSV data. Perfect for PhD research, formal reports, and advisor meetings.

## 📁 Package Structure

```text
src/
├── __init__.py              # Package initialization and exports
├── core.py                  # Core functionality (models, processing, generation)
├── config.py                # Configuration management
└── README.md                # This file
```

## 🚀 Usage

### As a Package

```python
from src import main, config, Task, ProjectTimeline

# Run the application
main()

# Access configuration
print(config.colors.researchcore)

# Create a task
task = Task(
    id="A",
    name="Sample Task",
    start_date=date(2025, 1, 1),
    due_date=date(2025, 1, 31),
    category="PROPOSAL",
    dependencies="",
    notes="Sample task description"
)
```

### As a Module

```bash
# Run the main application
python -m src.core --help

# Run with custom options
python -m src.core --input data.csv --output timeline.tex --title "My Project"
```

## 🔧 Core Components

### Data Models

- **`Task`**: Represents a single task with validation and computed properties
- **`ProjectTimeline`**: Complete project timeline with metadata

### Processing

- **`DataProcessor`**: Main coordinator for CSV processing and timeline building
- **`LaTeXGenerator`**: Generates complete LaTeX documents from timelines

### Application

- **`Application`**: Main application class with logging and validation
- **`main()`**: Entry point function

### Configuration

- **`config`**: Global configuration instance
- **`ColorScheme`**: Color definitions for all task categories
- **`TaskConfig`**: Task processing and categorization settings
- **`LaTeXConfig`**: LaTeX document generation settings

## 🔄 Architecture

### Design Principles

- **Simplicity**: Consolidated into essential components only
- **Type Safety**: Comprehensive type hints and validation
- **Error Handling**: Robust error handling with helpful messages
- **Configuration**: Centralized, customizable settings
- **Documentation**: Well-documented code with clear structure

### Key Benefits

- **Clean Structure**: No circular dependencies or complex imports
- **Easy to Use**: Simple API with sensible defaults
- **Maintainable**: Clear separation of concerns
- **Extensible**: Easy to add new features
- **Reliable**: Comprehensive validation and error handling

## 📝 License

Open source for academic and professional use.
