# 📦 Source Package - LaTeX Gantt Chart Generator

This directory contains the refactored source code for the LaTeX Gantt Chart Generator, organized as a proper Python package.

## 📁 Package Structure

```
src/
├── __init__.py              # Package initialization and exports
├── app.py                   # Main application entry point
├── build.py                 # Enhanced build system
├── config.py                # Core configuration classes
├── config_manager.py        # YAML-based configuration management
├── data_processor.py        # CSV processing and data pipeline
├── export_system.py         # Multi-format export system
├── interactive_generator.py # Interactive features and UI
├── latex_generator.py       # LaTeX generation components
├── models.py                # Data models and validation
├── template_generators.py   # Template generation system
├── utils.py                 # Shared utilities
└── config/                  # YAML configuration files
    ├── templates.yaml       # Template definitions
    └── device_profiles.yaml # Device profiles
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
python -m src.app --help

# Run with custom options
python -m src.app --input data.csv --output timeline.tex --title "My Project"
```

## 🔧 Development

### Adding New Features
1. **Configuration**: Add new settings to `config.py`
2. **Data Models**: Extend classes in `models.py`
3. **Processing**: Add new processors to `data_processor.py`
4. **LaTeX Generation**: Extend generators in `latex_generator.py`
5. **Templates**: Add new templates in `template_generators.py`
6. **Application**: Modify `app.py` for new CLI options

### Testing
```bash
# Test the package
python -c "import src; print('Package imported successfully')"

# Test individual modules
python -c "from src.models import Task; print('Models imported successfully')"
```

## 📚 API Reference

### Main Classes
- **`Application`**: Main application class with logging and validation
- **`BuildSystem`**: Enhanced build system for multiple templates and devices
- **`Task`**: Represents a single task with validation and computed properties
- **`ProjectTimeline`**: Complete project timeline with metadata
- **`DataProcessor`**: Main coordinator for the data processing pipeline
- **`LaTeXGenerator`**: Main coordinator for complete document generation
- **`TemplateGeneratorFactory`**: Factory for creating template generators

### Configuration
- **`config`**: Global configuration instance
- **`ColorScheme`**: Color definitions for all task categories
- **`CalendarConfig`**: Layout and styling configuration
- **`TaskConfig`**: Task processing and categorization settings
- **`LaTeXConfig`**: LaTeX document generation settings
- **`ConfigManager`**: Enhanced configuration management with YAML support

### Template Generators
- **`GanttTimelineGenerator`**: Enhanced Gantt timeline with modern TikZ features
- **`MonthlyCalendarGenerator`**: Monthly calendar view with task overlays
- **`WeeklyPlannerGenerator`**: Weekly planner with detailed scheduling

## 🔄 Architecture

### Design Patterns
- **Factory Pattern**: `TemplateGeneratorFactory` for creating generators
- **Strategy Pattern**: Different template generators for different output types
- **Configuration Pattern**: Centralized configuration management
- **Pipeline Pattern**: Data processing pipeline from CSV to LaTeX

### Key Benefits
- **Modular Design**: Clear separation of concerns
- **Type Safety**: Comprehensive type hints and validation
- **Error Handling**: Robust error handling with helpful messages
- **Configuration**: Centralized, customizable settings
- **Testing**: Easy to unit test individual components
- **Documentation**: Well-documented code with clear structure

## 📝 License

Open source for academic and professional use.