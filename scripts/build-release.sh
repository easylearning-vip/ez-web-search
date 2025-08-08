#!/bin/bash

# Build script for creating release binaries for multiple platforms

set -e

# Configuration
BINARY_NAME="ez-web-search-v2"
BUILD_DIR="build"
CMD_PATH="./cmd/server"

# Platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Clean build directory
clean_build() {
    log_info "Cleaning build directory..."
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    log_success "Build directory cleaned"
}

# Build for a specific platform
build_platform() {
    local platform=$1
    local os=$(echo $platform | cut -d'/' -f1)
    local arch=$(echo $platform | cut -d'/' -f2)
    
    log_info "Building for $os/$arch..."
    
    local output_name="${BINARY_NAME}_${os}_${arch}"
    if [ "$os" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local output_path="$BUILD_DIR/$output_name"
    
    # Set environment variables and build
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty)" \
        -o "$output_path" \
        "$CMD_PATH"
    
    # Verify the binary was created
    if [ -f "$output_path" ]; then
        local size=$(du -h "$output_path" | cut -f1)
        log_success "Built $output_name ($size)"
    else
        echo "❌ Failed to build $output_name"
        exit 1
    fi
}

# Create checksums
create_checksums() {
    log_info "Creating checksums..."
    
    cd "$BUILD_DIR"
    
    # Create SHA256 checksums
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum * > checksums.txt
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 * > checksums.txt
    else
        echo "⚠️  No checksum utility found, skipping checksums"
        cd ..
        return
    fi
    
    cd ..
    log_success "Checksums created"
}

# Create release archive
create_archive() {
    log_info "Creating release archive..."
    
    # Create a release info file
    cat > "$BUILD_DIR/RELEASE_INFO.txt" << EOF
EZ Web Search MCP Server Release
================================

Version: $(git describe --tags --always --dirty)
Build Date: $(date -u +"%Y-%m-%d %H:%M:%S UTC")
Git Commit: $(git rev-parse HEAD)

Included Binaries:
$(ls -la $BUILD_DIR/${BINARY_NAME}_* | awk '{print "- " $9 " (" $5 " bytes)"}')

Installation:
1. Download the appropriate binary for your platform
2. Make it executable: chmod +x ez-web-search-v2
3. Run the setup: ./setup-claude-cli.sh

Or use the one-click installer:
curl -fsSL https://raw.githubusercontent.com/easylearning-vip/ez-web-search/main/install.sh | bash

Documentation: https://github.com/easylearning-vip/ez-web-search
EOF
    
    log_success "Release archive prepared"
}

# Main build function
main() {
    echo "🔨 Building EZ Web Search MCP Server for multiple platforms"
    echo "==========================================================="
    echo
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo "❌ Not in a git repository"
        exit 1
    fi
    
    # Check if go is installed
    if ! command -v go >/dev/null 2>&1; then
        echo "❌ Go is not installed"
        exit 1
    fi
    
    # Clean build directory
    clean_build
    
    # Build for each platform
    for platform in "${PLATFORMS[@]}"; do
        build_platform "$platform"
    done
    
    # Create checksums
    create_checksums
    
    # Create release info
    create_archive
    
    echo
    log_success "🎉 Build completed successfully!"
    echo
    echo "📦 Built binaries:"
    ls -la "$BUILD_DIR"/${BINARY_NAME}_*
    echo
    echo "📋 Next steps:"
    echo "1. Test the binaries on target platforms"
    echo "2. Create a GitHub release"
    echo "3. Upload the binaries to the release"
    echo "4. Update the installation script if needed"
    echo
    echo "💡 Tip: Use 'gh release create' to create a GitHub release"
}

# Run main function
main "$@"
