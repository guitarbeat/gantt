# Cursor CLI Integration Guide

This document provides a comprehensive guide to the Cursor CLI integration in the PhD Dissertation Planner project.

## 🎯 Overview

The project now includes extensive Cursor CLI integration for AI-powered development workflows, including pre-commit hooks, test enhancement, code analysis, and automated improvements.

## 🚀 Quick Start

### 1. Install Cursor CLI
```bash
npm install -g @cursor/cli
```

### 2. Install Cursor CLI Hooks
```bash
make install-cursor-hooks
```

### 3. Test the Integration
```bash
make cursor-stats
make cursor-structure
```

## 📋 Available Commands

### Pre-commit Hooks
- `make install-cursor-hooks` - Install Cursor CLI pre-commit hooks
- `make test-cursor-hooks` - Test Cursor CLI hooks without committing
- `make uninstall-cursor-hooks` - Remove Cursor CLI hooks and restore previous
- `make cursor-precommit` - Run Cursor CLI pre-commit checks manually

### Test Enhancement
- `make cursor-test-enhance` - Run AI-powered test enhancement

### Development Tools
- `make cursor-dev-tools` - Run all Cursor CLI development tools
- `make cursor-dev-refactor` - Refactor code with AI
- `make cursor-dev-review` - AI-powered code review
- `make cursor-dev-optimize` - Optimize code performance with AI
- `make cursor-dev-docs` - Generate documentation with AI
- `make cursor-dev-fix` - Fix code issues with AI
- `make cursor-dev-complexity` - Analyze code complexity with AI
- `make cursor-dev-security` - Security analysis with AI
- `make cursor-dev-api-docs` - Generate API documentation with AI

### Simple Integration
- `make cursor-open` - Open entire project in Cursor
- `make cursor-file FILE=path` - Open specific file in Cursor
- `make cursor-structure` - Show project structure
- `make cursor-stats` - Show project statistics

## 🔧 Configuration

### Cursor CLI Configuration
The project includes a comprehensive Cursor CLI configuration file at `.cursor/config.yaml` that defines:
- Project context and preferences
- Code generation settings
- Testing preferences
- AI behavior configuration
- Security and performance preferences

### Pre-commit Hook Configuration
The pre-commit hooks are configured in `scripts/dev/cursor-precommit.sh` and include:
- Go formatting and linting
- YAML validation
- Large file detection
- Merge conflict detection
- AI-powered code analysis

## 🧪 Testing Integration

### Test Enhancement Tools
The `scripts/dev/cursor-test-enhancer.sh` script provides:
- AI-powered test coverage analysis
- Missing test generation
- Test failure analysis
- Performance benchmarking
- Test documentation generation

### Development Tools
The `scripts/dev/cursor-dev-tools.sh` script provides:
- Code refactoring with AI
- AI-powered code review
- Performance optimization
- Documentation generation
- Security analysis
- Complexity analysis

## 🤖 CI/CD Integration

### GitHub Actions Workflow
The `.github/workflows/cursor-ai-enhancement.yml` workflow provides:
- Automated AI code analysis
- AI-powered code optimization
- Test enhancement
- Documentation generation
- Automated PR creation with AI improvements

### Usage
1. Go to GitHub Actions
2. Select "Cursor AI Enhancement" workflow
3. Click "Run workflow"
4. Configure options:
   - Run AI analysis
   - Run AI optimization
   - Run AI testing
   - Create PR with changes

## 📁 File Structure

```
.cursor/
├── config.yaml                    # Cursor CLI configuration

scripts/dev/
├── cursor-precommit.sh            # Pre-commit hook script
├── cursor-test-enhancer.sh        # Test enhancement tools
├── cursor-dev-tools.sh            # Development tools
├── cursor-simple.sh               # Simple integration tools
├── install-cursor-hooks.sh        # Installation script
└── git-hook-pre-commit            # Git hook wrapper

.github/workflows/
└── cursor-ai-enhancement.yml      # AI enhancement workflow
```

## 🔍 How It Works

### Pre-commit Hooks
1. **File Analysis**: Scans staged files for various issues
2. **AI Processing**: Uses Cursor CLI for intelligent analysis
3. **Automatic Fixes**: Applies AI-suggested improvements
4. **Validation**: Ensures all checks pass before commit

### Test Enhancement
1. **Coverage Analysis**: Generates and analyzes test coverage
2. **Missing Tests**: Identifies and generates missing tests
3. **Test Quality**: Enhances existing tests with AI
4. **Documentation**: Generates test documentation

### Development Tools
1. **Code Analysis**: Comprehensive code review with AI
2. **Refactoring**: AI-powered code improvements
3. **Optimization**: Performance and memory optimization
4. **Documentation**: Automated documentation generation

## 🛠️ Customization

### Adding New Commands
1. Add new functions to the appropriate script
2. Add Makefile targets
3. Update help documentation
4. Test the integration

### Modifying AI Behavior
1. Edit `.cursor/config.yaml`
2. Adjust AI preferences and rules
3. Test with sample code
4. Commit changes

### Extending Workflows
1. Add new steps to GitHub Actions workflow
2. Create new scripts for specific tasks
3. Update documentation
4. Test in CI/CD environment

## 🐛 Troubleshooting

### Common Issues

**Cursor CLI not found:**
```bash
npm install -g @cursor/cli
```

**Hooks not working:**
```bash
make uninstall-cursor-hooks
make install-cursor-hooks
```

**Permission errors:**
```bash
chmod +x scripts/dev/*.sh
```

**YAML syntax errors:**
```bash
# Check YAML files
python3 -c "import yaml; yaml.safe_load(open('file.yaml'))"
```

### Debug Mode
Enable verbose logging by setting environment variables:
```bash
export CURSOR_DEBUG=1
export CURSOR_VERBOSE=1
```

## 📚 Examples

### Basic Usage
```bash
# Install and test
make install-cursor-hooks
make cursor-stats

# Development workflow
make cursor-dev-review
make cursor-test-enhance
git add .
git commit -m "feat: new feature"
```

### Advanced Usage
```bash
# Specific file analysis
make cursor-file FILE=internal/calendar/task_stacker.go

# Comprehensive analysis
make cursor-dev-tools
make cursor-test-enhance

# CI/CD integration
# Trigger GitHub Actions workflow manually
```

## 🎉 Benefits

### For Developers
- **Faster Development**: AI-powered code suggestions and fixes
- **Better Code Quality**: Automated analysis and improvements
- **Comprehensive Testing**: AI-generated and enhanced tests
- **Rich Documentation**: Automated documentation generation

### For the Project
- **Consistent Quality**: Standardized AI-powered checks
- **Automated Improvements**: Continuous code enhancement
- **Better Testing**: Higher test coverage and quality
- **Comprehensive Documentation**: Up-to-date project docs

## 🔮 Future Enhancements

- **Custom AI Models**: Project-specific AI training
- **Advanced Analytics**: Code quality metrics and trends
- **Team Collaboration**: Shared AI configurations
- **Integration Expansion**: More development tools and workflows

---

*This integration provides a solid foundation for AI-enhanced development workflows while maintaining the flexibility to adapt to changing project needs.*
