# ⚙️ Configuration Directory

This directory contains all configuration files for the PhD Dissertation Planner project, organized by category for easy management.

## 📁 Directory Structure

```
.config/
├── git/                    # Git configuration
│   ├── config             # Git settings and aliases
│   ├── attributes         # Git file attributes
│   └── ignore             # Project-specific gitignore
├── ide/                   # IDE configuration
│   └── vscode/           # VSCode settings
│       ├── settings.json  # Editor settings
│       ├── extensions.json # Recommended extensions
│       └── launch.json    # Debug configurations
├── build/                 # Build configuration
│   ├── makefile.conf     # Makefile settings
│   └── docker.conf       # Docker settings
├── dev/                   # Development environment
│   └── environment.conf   # Dev environment variables
├── scripts/               # Configuration management
│   └── config-manager.sh  # Configuration setup script
└── README.md             # This file
```

## 🚀 Quick Setup

### Initialize All Configurations
```bash
# Setup all configurations at once
./.config/scripts/config-manager.sh all

# Or setup individually
./.config/scripts/config-manager.sh git
./.config/scripts/config-manager.sh ide
./.config/scripts/config-manager.sh build
./.config/scripts/config-manager.sh dev
```

### Check Configuration Status
```bash
# See what's configured
./.config/scripts/config-manager.sh status
```

## 📋 Configuration Details

### 🔧 Git Configuration (`git/`)

- **`config`** - Git settings, aliases, and project-specific options
- **`attributes`** - File handling rules for different file types
- **`ignore`** - Project-specific gitignore patterns

**Key Features:**
- Project-specific aliases (`git build`, `git release`, `git organize`)
- LaTeX file handling
- Go file formatting
- PDF binary handling
- Comprehensive ignore patterns

### 💻 IDE Configuration (`ide/`)

- **VSCode Settings** - Editor configuration for Go, LaTeX, and Markdown
- **Extensions** - Recommended extensions for the project
- **Launch Configs** - Debug configurations for different scenarios

**Key Features:**
- Go development setup with language server
- LaTeX Workshop integration
- Markdown preview settings
- File nesting for related files
- Search exclusions for generated files

### 🔨 Build Configuration (`build/`)

- **Makefile Settings** - Build targets and Go/LaTeX settings
- **Docker Settings** - Container configuration

**Key Features:**
- Go build flags and test settings
- LaTeX engine configuration
- File patterns and directories
- Environment variable management

### 🛠️ Development Environment (`dev/`)

- **Environment Variables** - Development-specific settings
- **Hot Reload** - File watching configuration
- **Development Tools** - Required tools and installation

**Key Features:**
- Development mode flags
- Debug and verbose logging
- File watching for hot reload
- Development tool management

## 🎯 Usage Examples

### Git Aliases
```bash
# Use project-specific aliases
git build          # Run quick build
git release        # Create release
git organize       # Clean up files
git status-full    # Show project status
```

### VSCode Development
```bash
# Launch with different configurations
# - Launch Planner (default)
# - Launch Planner (Silent)
# - Launch Planner (Validate)
# - Launch Planner (Preview)
# - Launch Planner (Custom Config)
```

### Development Workflow
```bash
# Use development script
./dev.sh go run ./cmd/planner
./dev.sh make test
./dev.sh make build
```

## 🔄 Configuration Management

### Adding New Configurations
1. Create configuration file in appropriate directory
2. Update `config-manager.sh` to handle new config
3. Test with `./.config/scripts/config-manager.sh status`

### Updating Existing Configurations
1. Edit configuration file in `.config/` directory
2. Run `./.config/scripts/config-manager.sh [category]` to apply
3. Verify with `./.config/scripts/config-manager.sh status`

### Cleaning Configurations
```bash
# Remove all copied configuration files
./.config/scripts/config-manager.sh clean
```

## 📝 Configuration Files

### Git Configuration
- **Purpose**: Project-specific Git settings and aliases
- **Location**: `.config/git/`
- **Applied to**: `.git/config`, `.gitattributes`, `.gitignore`

### IDE Configuration
- **Purpose**: Editor settings and development environment
- **Location**: `.config/ide/`
- **Applied to**: `.vscode/` directory

### Build Configuration
- **Purpose**: Build system settings and Docker configuration
- **Location**: `.config/build/`
- **Applied to**: `Makefile`, Docker files

### Development Environment
- **Purpose**: Development-specific environment variables
- **Location**: `.config/dev/`
- **Applied to**: `.env.dev`, `dev.sh`

## 🎉 Benefits

- **🎯 Centralized Management** - All configurations in one place
- **🔄 Easy Updates** - Single script to manage all configs
- **📋 Version Control** - Track configuration changes
- **🛠️ Development Ready** - Pre-configured for optimal development
- **🧹 Clean Separation** - Configs separate from source code
- **📚 Well Documented** - Clear documentation for each config

## 🔧 Troubleshooting

### Configuration Not Applied
```bash
# Check status
./.config/scripts/config-manager.sh status

# Re-apply configuration
./.config/scripts/config-manager.sh [category]
```

### Git Configuration Issues
```bash
# Check git config
git config --list --local

# Reset git config
git config --unset-all [key]
```

### IDE Configuration Issues
```bash
# Check VSCode settings
cat .vscode/settings.json

# Reload VSCode window
# Cmd+Shift+P -> "Developer: Reload Window"
```

---

**💡 Pro Tip**: Use `./.config/scripts/config-manager.sh status` to quickly see what's configured and what needs attention!
