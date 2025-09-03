# 📦 LaTeX Gantt Chart Generator - Source Package

This directory contains the refactored source code for the LaTeX Gantt Chart Generator, organized as a proper Python package.

## 📁 Package Structure

```
src/
├── __init__.py          # Package initialization and exports
├── app.py               # Main application entry point
├── config.py            # Configuration management
├── models.py            # Data models and validation
├── data_processor.py    # CSV processing and data pipeline
└── latex_generator.py   # LaTeX generation components
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
5. **Application**: Modify `app.py` for new CLI options

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
- **`Task`**: Represents a single task with validation and computed properties
- **`ProjectTimeline`**: Complete project timeline with metadata
- **`DataProcessor`**: Main coordinator for the data processing pipeline
- **`LaTeXGenerator`**: Main coordinator for complete document generation

### Configuration
- **`config`**: Global configuration instance
- **`ColorScheme`**: Color definitions for all task categories
- **`CalendarConfig`**: Layout and styling configuration
- **`TaskConfig`**: Task processing and categorization settings
- **`LaTeXConfig`**: LaTeX document generation settings

## 🔄 Migration from Legacy Code

The package maintains full backward compatibility with the original monolithic script while providing a much more maintainable and extensible architecture.

### Key Benefits
- **Modular Design**: Clear separation of concerns
- **Type Safety**: Comprehensive type hints and validation
- **Error Handling**: Robust error handling with helpful messages
- **Configuration**: Centralized, customizable settings
- **Testing**: Easy to unit test individual components
- **Documentation**: Well-documented code with clear structure

## 📝 License

Open source for academic and professional use.
