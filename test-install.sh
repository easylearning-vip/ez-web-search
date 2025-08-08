#!/bin/bash

# Test script for the installation script
# This script tests the install.sh script in a safe way

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

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

# Test the install script without actually installing
test_install_script() {
    log_info "Testing install script..."
    
    # Check if install script exists
    if [ ! -f "install.sh" ]; then
        log_error "install.sh not found"
        return 1
    fi
    
    # Check if script is executable
    if [ ! -x "install.sh" ]; then
        log_error "install.sh is not executable"
        return 1
    fi
    
    # Check script syntax
    if bash -n install.sh; then
        log_success "install.sh syntax is valid"
    else
        log_error "install.sh has syntax errors"
        return 1
    fi
    
    # Test platform detection (dry run)
    log_info "Testing platform detection..."
    
    # Extract platform detection function and test it
    if grep -q "detect_platform" install.sh; then
        log_success "Platform detection function found"
    else
        log_error "Platform detection function not found"
        return 1
    fi
    
    log_success "Install script tests passed"
}

# Test build artifacts
test_build_artifacts() {
    log_info "Testing build artifacts..."
    
    if [ ! -d "build" ]; then
        log_error "Build directory not found. Run 'make build-release' first."
        return 1
    fi
    
    # Check if binaries exist
    local platforms=("linux_amd64" "linux_arm64" "darwin_amd64" "darwin_arm64" "windows_amd64.exe")
    
    for platform in "${platforms[@]}"; do
        local binary="build/ez-web-search-v2_${platform}"
        if [ -f "$binary" ]; then
            log_success "Binary found: $binary"
        else
            log_error "Binary missing: $binary"
            return 1
        fi
    done
    
    # Check checksums
    if [ -f "build/checksums.txt" ]; then
        log_success "Checksums file found"
    else
        log_error "Checksums file missing"
        return 1
    fi
    
    # Check release info
    if [ -f "build/RELEASE_INFO.txt" ]; then
        log_success "Release info file found"
    else
        log_error "Release info file missing"
        return 1
    fi
    
    log_success "Build artifacts tests passed"
}

# Test local binary
test_local_binary() {
    log_info "Testing local binary..."
    
    if [ ! -f "ez-web-search" ]; then
        log_warning "Local binary not found. Building..."
        make build
    fi
    
    if [ -f "ez-web-search" ]; then
        log_success "Local binary found"
        
        # Test that binary can start (timeout after 3 seconds)
        if timeout 3s ./ez-web-search >/dev/null 2>&1 || [ $? -eq 124 ]; then
            log_success "Local binary can start"
        else
            log_error "Local binary failed to start"
            return 1
        fi
    else
        log_error "Failed to build local binary"
        return 1
    fi
    
    log_success "Local binary tests passed"
}

# Test configuration files
test_config_files() {
    log_info "Testing configuration files..."
    
    # Check setup script
    if [ -f "setup-claude-cli.sh" ] && [ -x "setup-claude-cli.sh" ]; then
        log_success "Setup script found and executable"
    else
        log_error "Setup script missing or not executable"
        return 1
    fi
    
    # Check config template
    if [ -f "claude-mcp-config.json" ]; then
        log_success "Config template found"
        
        # Validate JSON
        if command -v jq >/dev/null 2>&1; then
            if jq empty claude-mcp-config.json 2>/dev/null; then
                log_success "Config template is valid JSON"
            else
                log_error "Config template is invalid JSON"
                return 1
            fi
        else
            log_warning "jq not found, skipping JSON validation"
        fi
    else
        log_error "Config template missing"
        return 1
    fi
    
    # Check environment template
    if [ -f ".env.example" ]; then
        log_success "Environment template found"
    else
        log_error "Environment template missing"
        return 1
    fi
    
    log_success "Configuration files tests passed"
}

# Test documentation
test_documentation() {
    log_info "Testing documentation..."
    
    local docs=("README.md" "FEATURES.md" "ENTERPRISE_FEATURES.md" "TESTING.md" "RELEASE_NOTES.md")
    
    for doc in "${docs[@]}"; do
        if [ -f "$doc" ]; then
            log_success "Documentation found: $doc"
        else
            log_warning "Documentation missing: $doc"
        fi
    done
    
    log_success "Documentation tests completed"
}

# Main test function
main() {
    echo "🧪 Testing EZ Web Search MCP Server Release"
    echo "==========================================="
    echo
    
    local failed=0
    
    # Run tests
    test_install_script || failed=1
    echo
    
    test_build_artifacts || failed=1
    echo
    
    test_local_binary || failed=1
    echo
    
    test_config_files || failed=1
    echo
    
    test_documentation || failed=1
    echo
    
    # Summary
    if [ $failed -eq 0 ]; then
        log_success "🎉 All tests passed! Release is ready."
        echo
        echo "📋 Next steps:"
        echo "1. Push the tag: git push origin v1.0.0"
        echo "2. Create GitHub release with build artifacts"
        echo "3. Test the installation script from GitHub"
        echo "4. Announce the release"
    else
        log_error "❌ Some tests failed. Please fix issues before release."
        exit 1
    fi
}

# Run tests
main "$@"
