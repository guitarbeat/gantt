# Ad-Hoc Debugging System

A comprehensive implementation of the Agentic Project Management (APM) Ad-Hoc debugging workflow for the LaTeX Gantt Chart Generator.

## 🎯 Overview

This system implements the APM framework's Ad-Hoc debugging workflow as described in the [Agentic Project Management repository](https://github.com/sdi2200262/agentic-project-management). It provides structured debugging capabilities with delegation prompts and session management.

## 🔄 Workflow

The system follows the APM debugging workflow:

1. **Issue Complexity Assessment** - Automatically determines if an issue is simple or complex
2. **Local Debugging** - Attempts local resolution for simple issues (up to 2 attempts)
3. **Delegation System** - Generates structured prompts for complex or persistent issues
4. **Ad-Hoc Session Management** - Manages specialized debugging sessions
5. **Solution Integration** - Integrates findings back into the main workflow

## 📁 System Components

### Core Modules

- **`debug_system.py`** - Main debugging workflow implementation
- **`session_manager.py`** - Session persistence and management
- **`debug_cli.py`** - Command-line interface for debugging operations
- **`test_debug_system.py`** - Comprehensive test suite

### Integration

- **`app.py`** - Integrated with main application for automatic error handling
- **`main.py`** - Added debug command to main CLI

## 🚀 Quick Start

### Basic Usage

```bash
# Create a debug session
python main.py debug create --issue-id "bug_001" --title "Import Error" --description "Module not found"

# List all sessions
python main.py debug list

# Show session details
python main.py debug show --session-id "session_123"

# Handle an issue through the workflow
python main.py debug handle --issue-id "bug_002" --title "Build Failure" --description "LaTeX compilation failed"

# Complete a workflow with solution
python main.py debug complete --session-id "session_123" --solution "Fixed import path"

# Show statistics
python main.py debug stats

# Export session report
python main.py debug export --session-id "session_123" --format markdown --output report.md
```

### Integration with Main Application

The debug system is automatically integrated with the main LaTeX generator. When errors occur, they are automatically processed through the APM workflow:

```python
# Errors are automatically handled
app = Application()
success = app.generate_latex_file("input.csv", "output.tex", "My Timeline")
# If an error occurs, it's automatically processed through the debug system
```

## 🔧 Configuration

### Session Storage

Debug sessions are stored in:
- **Database**: `debug_sessions/sessions.db` (SQLite)
- **JSON Files**: `debug_sessions/*.json` (detailed session data)

### Settings

- **Max Local Attempts**: 2 (configurable)
- **Session Retention**: 30 days (configurable)
- **Cleanup Interval**: 1 hour (configurable)

## 📊 Features

### Issue Complexity Assessment

The system automatically assesses issue complexity based on:
- **Simple Indicators**: syntax errors, import errors, typos, missing files
- **Complex Indicators**: race conditions, memory leaks, performance issues, integration errors

### Local Debugging

For simple issues, the system attempts local resolution with:
- Structured attempt tracking
- Success/failure logging
- Duration measurement
- Error message capture

### Delegation Prompts

For complex issues, the system generates structured prompts including:
- Issue context and history
- Specific instructions for the Ad-Hoc agent
- Expected deliverables
- Constraints and success criteria

### Session Management

Comprehensive session tracking with:
- Persistent storage (database + JSON)
- Search and filtering capabilities
- Statistics and reporting
- Automatic cleanup of old sessions

## 🧪 Testing

Run the comprehensive test suite:

```bash
python test_debug_system.py
```

The test suite covers:
1. Simple issue handling
2. Complex issue delegation
3. Workflow integration
4. Session management
5. Solution integration
6. CLI interface

## 📋 CLI Commands

### Session Management

```bash
# Create a new session
python main.py debug create --issue-id "ID" --title "Title" --description "Description"

# List sessions with filtering
python main.py debug list --status resolved --limit 10

# Show detailed session information
python main.py debug show --session-id "SESSION_ID"

# Search sessions
python main.py debug search --query "import error"
```

### Workflow Operations

```bash
# Handle an issue through the complete workflow
python main.py debug handle --issue-id "ID" --title "Title" --description "Description"

# Complete a workflow with solution findings
python main.py debug complete --session-id "SESSION_ID" --solution "Solution details"
```

### Reporting and Analysis

```bash
# Show debug statistics
python main.py debug stats

# Export session data
python main.py debug export --session-id "SESSION_ID" --format markdown --output report.md

# Clean up old sessions
python main.py debug cleanup
```

## 🔍 Example Workflows

### Simple Issue Resolution

```bash
# 1. Create a simple issue
python main.py debug create --issue-id "simple_001" --title "File Not Found" --description "Configuration file missing"

# 2. Handle through workflow (will attempt local resolution)
python main.py debug handle --issue-id "simple_001" --title "File Not Found" --description "Configuration file missing"

# 3. If resolved locally, session will show as resolved
python main.py debug show --session-id "SESSION_ID"
```

### Complex Issue Delegation

```bash
# 1. Create a complex issue
python main.py debug create --issue-id "complex_001" --title "Performance Issue" --description "LaTeX generation is slow with large datasets"

# 2. Handle through workflow (will generate delegation prompt)
python main.py debug handle --issue-id "complex_001" --title "Performance Issue" --description "LaTeX generation is slow with large datasets" --save-prompt

# 3. Use the generated delegation prompt in an Ad-Hoc session
# 4. Complete the workflow with solution findings
python main.py debug complete --session-id "SESSION_ID" --solution "Optimized template generation algorithm"
```

## 📈 Statistics and Reporting

The system provides comprehensive statistics:

- **Total Sessions**: Number of debug sessions created
- **Resolution Rate**: Percentage of successfully resolved issues
- **Average Attempts**: Average number of attempts per session
- **Status Breakdown**: Distribution of session statuses
- **Complexity Breakdown**: Distribution of issue complexities

## 🔧 Advanced Usage

### Custom Debug Functions

You can provide custom debug functions for local resolution:

```python
def custom_debug_function():
    # Your custom debugging logic
    return "Solution found"

workflow = DebugWorkflow()
success, message, prompt = workflow.handle_issue(
    issue_id="custom_001",
    title="Custom Issue",
    description="Issue requiring custom debugging",
    debug_function=custom_debug_function
)
```

### Session Integration

Access session data programmatically:

```python
from src.session_manager import get_session_manager

session_manager = get_session_manager()
session = session_manager.load_session("session_id")
stats = session_manager.get_session_statistics()
```

### Error Context

Provide rich error context for better issue assessment:

```python
error_context = {
    "error_type": "ValueError",
    "error_message": "Invalid input",
    "file": "main.py",
    "line": 42,
    "stack_trace": "Full stack trace...",
    "user_input": "problematic input"
}
```

## 🛠️ Development

### Adding New Features

1. **New Issue Types**: Extend the complexity assessment logic
2. **Custom Recovery**: Add new recovery strategies for simple issues
3. **Export Formats**: Add new export formats for session data
4. **Integration Points**: Add integration with other systems

### Database Schema

The system uses SQLite with the following tables:
- **sessions**: Main session metadata
- **attempts**: Individual debug attempts
- **delegation_prompts**: Generated delegation prompts

### Configuration

Key configuration options:
- `max_local_attempts`: Maximum local debugging attempts
- `session_retention_days`: How long to keep completed sessions
- `cleanup_interval`: How often to run cleanup tasks

## 📚 API Reference

### DebugWorkflow

Main workflow orchestrator:

```python
workflow = DebugWorkflow()
success, message, delegation_prompt = workflow.handle_issue(
    issue_id, title, description, error_context, debug_function
)
```

### AdHocDebugger

Core debugging functionality:

```python
debugger = AdHocDebugger()
session = debugger.create_debug_session(issue_id, title, description, error_context)
attempt = debugger.attempt_local_debug(session_id, debug_function)
prompt = debugger.generate_delegation_prompt(session_id)
```

### SessionManager

Session persistence and management:

```python
session_manager = SessionManager()
session_manager.save_session(session)
session = session_manager.load_session(session_id)
stats = session_manager.get_session_statistics()
```

## 🤝 Contributing

Contributions are welcome! Areas for improvement:

- **New Recovery Strategies**: Add more automatic recovery options
- **Enhanced Analytics**: Better reporting and analysis features
- **Integration**: Connect with external debugging tools
- **UI**: Web-based interface for session management
- **Testing**: Additional test cases and scenarios

## 📄 License

This implementation follows the same license as the main project.

## 🔗 References

- [Agentic Project Management Framework](https://github.com/sdi2200262/agentic-project-management)
- [APM Documentation](https://github.com/sdi2200262/agentic-project-management/tree/main/docs)
- [APM Workflow Overview](https://github.com/sdi2200262/agentic-project-management/tree/main/docs/Workflow_Overview.md)

---

**Perfect for**: Complex debugging scenarios, systematic issue resolution, and maintaining debugging context across extended development sessions.
