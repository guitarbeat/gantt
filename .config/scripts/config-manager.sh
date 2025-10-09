#!/bin/bash

# Configuration Manager for PhD Dissertation Planner
# This script manages all project configurations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CONFIG_DIR="$PROJECT_ROOT/.config"

# Function to show help
show_help() {
    echo -e "${BLUE}Configuration Manager for PhD Dissertation Planner${NC}"
    echo "======================================================"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  init          - Initialize configuration files"
    echo "  git           - Setup git configuration"
    echo "  ide           - Setup IDE configuration"
    echo "  build         - Setup build configuration"
    echo "  dev           - Setup development environment"
    echo "  all           - Setup all configurations"
    echo "  status        - Show configuration status"
    echo "  clean         - Clean configuration files"
    echo "  help          - Show this help message"
    echo ""
    echo "Options:"
    echo "  --force       - Force overwrite existing files"
    echo "  --backup      - Create backup before overwriting"
    echo "  --dry-run     - Show what would be done without making changes"
    echo ""
}

# Function to initialize git configuration
setup_git() {
    echo -e "${YELLOW}🔧 Setting up Git configuration...${NC}"
    
    # Copy git config to project root
    if [ -f "$CONFIG_DIR/git/config" ]; then
        cp "$CONFIG_DIR/git/config" "$PROJECT_ROOT/.git/config"
        echo -e "${GREEN}✅ Git config copied${NC}"
    fi
    
    # Copy git attributes
    if [ -f "$CONFIG_DIR/git/attributes" ]; then
        cp "$CONFIG_DIR/git/attributes" "$PROJECT_ROOT/.gitattributes"
        echo -e "${GREEN}✅ Git attributes copied${NC}"
    fi
    
    # Copy git ignore
    if [ -f "$CONFIG_DIR/git/ignore" ]; then
        cp "$CONFIG_DIR/git/ignore" "$PROJECT_ROOT/.gitignore"
        echo -e "${GREEN}✅ Git ignore copied${NC}"
    fi
}

# Function to setup IDE configuration
setup_ide() {
    echo -e "${YELLOW}🔧 Setting up IDE configuration...${NC}"
    
    # Create VSCode directory
    mkdir -p "$PROJECT_ROOT/.vscode"
    
    # Copy VSCode settings
    if [ -f "$CONFIG_DIR/ide/vscode/settings.json" ]; then
        cp "$CONFIG_DIR/ide/vscode/settings.json" "$PROJECT_ROOT/.vscode/settings.json"
        echo -e "${GREEN}✅ VSCode settings copied${NC}"
    fi
    
    # Copy VSCode extensions
    if [ -f "$CONFIG_DIR/ide/vscode/extensions.json" ]; then
        cp "$CONFIG_DIR/ide/vscode/extensions.json" "$PROJECT_ROOT/.vscode/extensions.json"
        echo -e "${GREEN}✅ VSCode extensions copied${NC}"
    fi
    
    # Copy VSCode launch config
    if [ -f "$CONFIG_DIR/ide/vscode/launch.json" ]; then
        cp "$CONFIG_DIR/ide/vscode/launch.json" "$PROJECT_ROOT/.vscode/launch.json"
        echo -e "${GREEN}✅ VSCode launch config copied${NC}"
    fi
}

# Function to setup build configuration
setup_build() {
    echo -e "${YELLOW}🔧 Setting up build configuration...${NC}"
    
    # Include build config in main Makefile
    if [ -f "$CONFIG_DIR/build/makefile.conf" ]; then
        echo "" >> "$PROJECT_ROOT/Makefile"
        echo "# Include build configuration" >> "$PROJECT_ROOT/Makefile"
        echo "-include .config/build/makefile.conf" >> "$PROJECT_ROOT/Makefile"
        echo -e "${GREEN}✅ Build configuration included in Makefile${NC}"
    fi
}

# Function to setup development environment
setup_dev() {
    echo -e "${YELLOW}🔧 Setting up development environment...${NC}"
    
    # Create development environment file
    if [ -f "$CONFIG_DIR/dev/environment.conf" ]; then
        cp "$CONFIG_DIR/dev/environment.conf" "$PROJECT_ROOT/.env.dev"
        echo -e "${GREEN}✅ Development environment file created${NC}"
    fi
    
    # Create development script
    cat > "$PROJECT_ROOT/dev.sh" << 'EOF'
#!/bin/bash
# Development script for PhD Dissertation Planner

# Load development environment
if [ -f ".env.dev" ]; then
    source .env.dev
fi

# Set development environment variables
export PLANNER_CSV_FILE="${DEV_PLANNER_CSV_FILE:-input_data/research_timeline_v5.1_comprehensive.csv}"
export PLANNER_CONFIG_FILE="${DEV_PLANNER_CONFIG_FILE:-configs/base.yaml}"
export PLANNER_OUTPUT_DIR="${DEV_PLANNER_OUTPUT_DIR:-generated}"
export PLANNER_SILENT="${DEV_PLANNER_SILENT:-0}"
export PLANNER_DEBUG="${DEV_PLANNER_DEBUG:-1}"
export PLANNER_VERBOSE="${DEV_PLANNER_VERBOSE:-1}"

# Run the command
exec "$@"
EOF
    
    chmod +x "$PROJECT_ROOT/dev.sh"
    echo -e "${GREEN}✅ Development script created${NC}"
}

# Function to show configuration status
show_status() {
    echo -e "${BLUE}📊 Configuration Status${NC}"
    echo "======================"
    
    echo -e "\n${YELLOW}Git Configuration:${NC}"
    if [ -f "$PROJECT_ROOT/.git/config" ]; then
        echo "✅ Git config: $PROJECT_ROOT/.git/config"
    else
        echo "❌ Git config: Not found"
    fi
    
    if [ -f "$PROJECT_ROOT/.gitattributes" ]; then
        echo "✅ Git attributes: $PROJECT_ROOT/.gitattributes"
    else
        echo "❌ Git attributes: Not found"
    fi
    
    if [ -f "$PROJECT_ROOT/.gitignore" ]; then
        echo "✅ Git ignore: $PROJECT_ROOT/.gitignore"
    else
        echo "❌ Git ignore: Not found"
    fi
    
    echo -e "\n${YELLOW}IDE Configuration:${NC}"
    if [ -d "$PROJECT_ROOT/.vscode" ]; then
        echo "✅ VSCode directory: $PROJECT_ROOT/.vscode"
        if [ -f "$PROJECT_ROOT/.vscode/settings.json" ]; then
            echo "  ✅ Settings: $PROJECT_ROOT/.vscode/settings.json"
        fi
        if [ -f "$PROJECT_ROOT/.vscode/extensions.json" ]; then
            echo "  ✅ Extensions: $PROJECT_ROOT/.vscode/extensions.json"
        fi
        if [ -f "$PROJECT_ROOT/.vscode/launch.json" ]; then
            echo "  ✅ Launch: $PROJECT_ROOT/.vscode/launch.json"
        fi
    else
        echo "❌ VSCode directory: Not found"
    fi
    
    echo -e "\n${YELLOW}Build Configuration:${NC}"
    if [ -f "$CONFIG_DIR/build/makefile.conf" ]; then
        echo "✅ Build config: $CONFIG_DIR/build/makefile.conf"
    else
        echo "❌ Build config: Not found"
    fi
    
    echo -e "\n${YELLOW}Development Environment:${NC}"
    if [ -f "$PROJECT_ROOT/.env.dev" ]; then
        echo "✅ Dev environment: $PROJECT_ROOT/.env.dev"
    else
        echo "❌ Dev environment: Not found"
    fi
    
    if [ -f "$PROJECT_ROOT/dev.sh" ]; then
        echo "✅ Dev script: $PROJECT_ROOT/dev.sh"
    else
        echo "❌ Dev script: Not found"
    fi
}

# Function to clean configuration files
clean_config() {
    echo -e "${YELLOW}🧹 Cleaning configuration files...${NC}"
    
    # Remove copied files
    rm -f "$PROJECT_ROOT/.git/config"
    rm -f "$PROJECT_ROOT/.gitattributes"
    rm -f "$PROJECT_ROOT/.gitignore"
    rm -rf "$PROJECT_ROOT/.vscode"
    rm -f "$PROJECT_ROOT/.env.dev"
    rm -f "$PROJECT_ROOT/dev.sh"
    
    echo -e "${GREEN}✅ Configuration files cleaned${NC}"
}

# Main execution
main() {
    case "${1:-help}" in
        init)
            setup_git
            setup_ide
            setup_build
            setup_dev
            echo -e "${GREEN}🎉 Configuration initialized!${NC}"
            ;;
        git)
            setup_git
            ;;
        ide)
            setup_ide
            ;;
        build)
            setup_build
            ;;
        dev)
            setup_dev
            ;;
        all)
            setup_git
            setup_ide
            setup_build
            setup_dev
            echo -e "${GREEN}🎉 All configurations setup!${NC}"
            ;;
        status)
            show_status
            ;;
        clean)
            clean_config
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            echo -e "${RED}Error: Unknown command '$1'${NC}"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
