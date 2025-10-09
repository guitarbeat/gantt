# ⚙️ Configuration Guide

## 🎯 Overview

Your project now has a comprehensive configuration system that organizes all settings, tools, and development environment configurations in one place!

## 📁 Configuration Structure

```
.config/                          # All configurations organized here
├── git/                          # Git settings and aliases
│   ├── config                   # Git configuration with project aliases
│   ├── attributes               # File handling rules
│   └── ignore                   # Project-specific gitignore
├── ide/                         # IDE and editor settings
│   └── vscode/                  # VSCode configuration
│       ├── settings.json        # Editor settings for Go/LaTeX/Markdown
│       ├── extensions.json      # Recommended extensions
│       └── launch.json          # Debug configurations
├── build/                       # Build system configuration
│   ├── makefile.conf           # Makefile settings and overrides
│   └── docker.conf             # Docker configuration
├── dev/                         # Development environment
│   └── environment.conf         # Development variables and settings
├── scripts/                     # Configuration management
│   └── config-manager.sh        # Setup and management script
└── README.md                    # Detailed configuration documentation
```

## 🚀 Quick Start

### 1. Initialize All Configurations
```bash
# Setup everything at once
./.config/scripts/config-manager.sh all
```

### 2. Check What's Configured
```bash
# See configuration status
./.config/scripts/config-manager.sh status
```

### 3. Use Project Aliases
```bash
# Git aliases are now available
git build          # Quick build
git release        # Create release
git organize       # Clean up files
git status-full    # Show project status
```

## 🎯 Key Features

### 🔧 Git Configuration
- **Project-specific aliases** for common tasks
- **File handling rules** for LaTeX, Go, PDFs
- **Comprehensive gitignore** patterns
- **Color-coded output** for better readability

### 💻 IDE Configuration
- **VSCode settings** optimized for Go + LaTeX + Markdown
- **Recommended extensions** for the project
- **Debug configurations** for different scenarios
- **File nesting** for related files

### 🔨 Build Configuration
- **Makefile overrides** with project-specific settings
- **Docker configuration** for containerized builds
- **Environment variables** management
- **Build targets** and optimization flags

### 🛠️ Development Environment
- **Development variables** for debugging
- **Hot reload** configuration
- **Tool management** and installation
- **Development script** for easy setup

## 📋 Configuration Commands

### Setup Commands
```bash
# Setup all configurations
./.config/scripts/config-manager.sh all

# Setup individual categories
./.config/scripts/config-manager.sh git
./.config/scripts/config-manager.sh ide
./.config/scripts/config-manager.sh build
./.config/scripts/config-manager.sh dev
```

### Management Commands
```bash
# Check configuration status
./.config/scripts/config-manager.sh status

# Clean all configurations
./.config/scripts/config-manager.sh clean

# Show help
./.config/scripts/config-manager.sh help
```

## 🎉 Benefits

### 🧠 **Reduced Cognitive Load**
- All configurations in one place
- Clear organization by category
- Easy to find and modify settings

### 🔄 **Easy Management**
- Single script to manage all configs
- Version control for configuration changes
- Easy to share and replicate setups

### 🛠️ **Development Ready**
- Pre-configured for optimal development
- IDE settings for Go, LaTeX, and Markdown
- Debug configurations for different scenarios

### 📚 **Well Documented**
- Clear documentation for each configuration
- Examples and usage patterns
- Troubleshooting guides

## 🔧 Usage Examples

### Git Workflow
```bash
# Use project aliases
git build          # Run quick build
git release        # Create timestamped release
git organize       # Clean up project files
git status-full    # Show detailed project status

# Standard git commands work as usual
git add .
git commit -m "Update configuration"
git push
```

### VSCode Development
```bash
# Open project in VSCode
code .

# Use debug configurations
# - Launch Planner (default)
# - Launch Planner (Silent)
# - Launch Planner (Validate)
# - Launch Planner (Preview)
# - Launch Planner (Custom Config)
```

### Development Workflow
```bash
# Use development script for environment setup
./dev.sh go run ./cmd/planner
./dev.sh make test
./dev.sh make build

# Or use standard commands with dev environment
source .env.dev
go run ./cmd/planner
```

## 📝 Configuration Files

### Git Configuration
- **Location**: `.config/git/`
- **Applied to**: `.git/config`, `.gitattributes`, `.gitignore`
- **Features**: Project aliases, file handling, color output

### IDE Configuration
- **Location**: `.config/ide/`
- **Applied to**: `.vscode/` directory
- **Features**: Editor settings, extensions, debug configs

### Build Configuration
- **Location**: `.config/build/`
- **Applied to**: `Makefile`, Docker files
- **Features**: Build settings, environment variables

### Development Environment
- **Location**: `.config/dev/`
- **Applied to**: `.env.dev`, `dev.sh`
- **Features**: Dev variables, tool management

## 🔄 Maintenance

### Updating Configurations
1. Edit files in `.config/` directory
2. Run `./.config/scripts/config-manager.sh [category]` to apply
3. Verify with `./.config/scripts/config-manager.sh status`

### Adding New Configurations
1. Create new config file in appropriate `.config/` subdirectory
2. Update `config-manager.sh` to handle new config
3. Test with status command

### Cleaning Configurations
```bash
# Remove all copied configuration files
./.config/scripts/config-manager.sh clean
```

## 🎯 Next Steps

1. **Explore configurations** - Check out the files in `.config/` directory
2. **Customize settings** - Modify configurations to your preferences
3. **Use project aliases** - Try the new git aliases
4. **Set up IDE** - Use VSCode with the provided settings
5. **Development workflow** - Use the development script for environment setup

---

**💡 Pro Tip**: The configuration system is designed to grow with your project. Add new configurations as needed, and they'll be automatically managed by the config-manager script!

**🎉 Your project is now fully configured and ready for development!**
