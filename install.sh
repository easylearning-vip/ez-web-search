#!/bin/bash

# EZ Web Search - Binary Installation Script
# This script downloads and installs the EZ Web Search binary

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
        log_info "Adding $INSTALL_DIR to PATH..."
        
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
            echo "# Added by EZ Web Search installer" >> "$SHELL_PROFILE"
            echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$SHELL_PROFILE"
            log_success "Added $INSTALL_DIR to PATH in $SHELL_PROFILE"
            
            # Source the profile to update current session
            source "$SHELL_PROFILE"
            log_success "PATH updated for current session"
        else
            log_warning "Please manually add $INSTALL_DIR to your PATH"
        fi
    else
        log_success "$INSTALL_DIR is already in your PATH"
    fi
}


# Main installation function
main() {
    echo "🚀 EZ Web Search - Binary Installer"
    echo "==================================="
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
    
    echo
    log_success "🎉 Installation completed successfully!"
    echo
    echo "📋 Ready to use:"
    echo "   $BINARY_NAME"
    echo
    echo "📚 Documentation: https://github.com/$REPO"
    echo "🔧 Binary location: $INSTALL_DIR/$BINARY_NAME"
}

# Run main function
main "$@"
