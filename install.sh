#!/bin/bash

# EZ Web Search MCP Server - One-Click Installation Script
# This script downloads, installs, and configures the EZ Web Search MCP Server

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="easylearning-vip/ez-web-search"
BINARY_NAME="ez-web-search"
INSTALL_DIR="$HOME/.local/bin"
CONFIG_DIR="$HOME/.claude"

# Functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $os in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case $arch in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    PLATFORM="${OS}_${ARCH}"
    log_info "Detected platform: $PLATFORM"
}

# Get latest release version
get_latest_version() {
    log_info "Fetching latest release information..."
    
    if command -v curl >/dev/null 2>&1; then
        LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        LATEST_VERSION=$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        log_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    if [ -z "$LATEST_VERSION" ]; then
        log_error "Failed to fetch latest version"
        exit 1
    fi
    
    log_success "Latest version: $LATEST_VERSION"
}

# Download and install binary
install_binary() {
    log_info "Downloading EZ Web Search MCP Server..."
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    # Download URL
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/${BINARY_NAME}_${PLATFORM}"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$INSTALL_DIR/$BINARY_NAME" "$DOWNLOAD_URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$INSTALL_DIR/$BINARY_NAME" "$DOWNLOAD_URL"
    fi
    
    # Make executable
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    log_success "Binary installed to $INSTALL_DIR/$BINARY_NAME"
}

# Add to PATH if needed
setup_path() {
    # Check if install directory is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        log_warning "$INSTALL_DIR is not in your PATH"
        
        # Add to shell profile
        SHELL_PROFILE=""
        if [ -f "$HOME/.bashrc" ]; then
            SHELL_PROFILE="$HOME/.bashrc"
        elif [ -f "$HOME/.zshrc" ]; then
            SHELL_PROFILE="$HOME/.zshrc"
        elif [ -f "$HOME/.profile" ]; then
            SHELL_PROFILE="$HOME/.profile"
        fi
        
        if [ -n "$SHELL_PROFILE" ]; then
            echo "" >> "$SHELL_PROFILE"
            echo "# Added by EZ Web Search MCP Server installer" >> "$SHELL_PROFILE"
            echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$SHELL_PROFILE"
            log_success "Added $INSTALL_DIR to PATH in $SHELL_PROFILE"
            log_warning "Please restart your shell or run: source $SHELL_PROFILE"
        else
            log_warning "Please manually add $INSTALL_DIR to your PATH"
        fi
    else
        log_success "$INSTALL_DIR is already in your PATH"
    fi
}

# Configure Claude Code CLI
configure_claude_cli() {
    log_info "Configuring Claude Code CLI..."
    
    # Check if Claude CLI is installed
    if ! command -v claude >/dev/null 2>&1; then
        log_warning "Claude Code CLI not found. Please install it from: https://claude.ai/cli"
        log_info "You can configure the MCP server manually later"
        return
    fi
    
    # Create Claude config directory
    mkdir -p "$CONFIG_DIR"
    
    # MCP configuration file
    MCP_CONFIG="$CONFIG_DIR/mcp_settings.json"
    
    # Backup existing config
    if [ -f "$MCP_CONFIG" ]; then
        log_warning "Existing MCP configuration found. Creating backup..."
        cp "$MCP_CONFIG" "$MCP_CONFIG.backup.$(date +%Y%m%d_%H%M%S)"
    fi
    
    # Prompt for BigModel API token
    echo
    log_info "BigModel API Token Configuration"
    echo "Please enter your BigModel API token:"
    echo "(Get one from: https://open.bigmodel.cn/)"
    echo "Press Enter to use the default test token"
    echo
    read -p "BigModel API Token: " BIGMODEL_TOKEN
    
    if [ -z "$BIGMODEL_TOKEN" ]; then
        BIGMODEL_TOKEN="0f405f7a11b946298b154f042a70f12b.s6VO3ITALpa3bhDo"
        log_warning "Using default test token"
    fi
    
    # Create MCP configuration
    cat > "$MCP_CONFIG" << EOF
{
  "mcpServers": {
    "ez-web-search": {
      "command": "$INSTALL_DIR/$BINARY_NAME",
      "env": {
        "BIGMODEL_TOKEN": "$BIGMODEL_TOKEN",
        "WEBFETCH_USER_AGENT_ROTATE": "true",
        "WEBFETCH_DELAY_MIN": "1s",
        "WEBFETCH_DELAY_MAX": "3s",
        "WEBFETCH_MAX_CONTENT_SIZE": "5000",
        "WEBFETCH_MAX_LINKS": "50",
        "WEBFETCH_MAX_IMAGES": "20",
        "PATH": "/usr/local/bin:/usr/bin:/bin:$INSTALL_DIR"
      }
    }
  }
}
EOF
    
    log_success "Claude Code CLI configured: $MCP_CONFIG"
}

# Test installation
test_installation() {
    log_info "Testing installation..."
    
    # Test binary
    if [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
        log_success "Binary is executable"
    else
        log_error "Binary is not executable"
        return 1
    fi
    
    # Test MCP configuration
    if [ -f "$CONFIG_DIR/mcp_settings.json" ]; then
        if command -v jq >/dev/null 2>&1; then
            if jq empty "$CONFIG_DIR/mcp_settings.json" 2>/dev/null; then
                log_success "MCP configuration is valid JSON"
            else
                log_error "MCP configuration is invalid JSON"
                return 1
            fi
        else
            log_success "MCP configuration file created"
        fi
    fi
    
    log_success "Installation test passed!"
}

# Main installation function
main() {
    echo "🚀 EZ Web Search MCP Server - One-Click Installer"
    echo "================================================="
    echo
    
    # Check prerequisites
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        log_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    # Detect platform
    detect_platform
    
    # Get latest version
    get_latest_version
    
    # Install binary
    install_binary
    
    # Setup PATH
    setup_path
    
    # Configure Claude CLI
    configure_claude_cli
    
    # Test installation
    test_installation
    
    echo
    log_success "🎉 Installation completed successfully!"
    echo
    echo "📋 Next steps:"
    echo "1. Restart your shell or run: source ~/.bashrc (or ~/.zshrc)"
    echo "2. Start Claude Code CLI: claude"
    echo "3. Test web search: Search for \"Go programming tutorials\""
    echo "4. Test web fetch: Fetch content from https://example.com"
    echo
    echo "📚 Documentation: https://github.com/$REPO"
    echo "🔧 Binary location: $INSTALL_DIR/$BINARY_NAME"
    echo "🔧 Config location: $CONFIG_DIR/mcp_settings.json"
    echo
    echo "💡 Tip: You can edit the config file to customize settings"
}

# Run main function
main "$@"
