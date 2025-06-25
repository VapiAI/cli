#!/usr/bin/env bash
# Vapi CLI installation script
# Usage: curl -sSL https://vapi.ai/install.sh | bash

set -e

# Configuration
REPO="VapiAI/cli"
BINARY_NAME="vapi"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        darwin) OS="Darwin" ;;
        linux) OS="Linux" ;;
        mingw*|msys*|cygwin*) OS="Windows" ;;
        *) error "Unsupported operating system: $OS" ;;
    esac
    
    case "$ARCH" in
        x86_64) ARCH="x86_64" ;;
        amd64) ARCH="x86_64" ;;
        aarch64) ARCH="arm64" ;;
        arm64) ARCH="arm64" ;;
        armv7l) ARCH="armv7" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac
    
    PLATFORM="${OS}_${ARCH}"
    log "Detected platform: $PLATFORM"
}

# Get latest release version
get_latest_version() {
    log "Fetching latest version..."
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$VERSION" ]; then
        error "Failed to fetch latest version"
    fi
    
    log "Latest version: $VERSION"
}

# Download and install
install_vapi() {
    local url="https://github.com/$REPO/releases/download/$VERSION/${BINARY_NAME}_${PLATFORM}.tar.gz"
    local tmp_dir=$(mktemp -d)
    
    log "Downloading Vapi CLI..."
    log "URL: $url"
    
    if ! curl -sL "$url" -o "$tmp_dir/vapi.tar.gz"; then
        error "Failed to download Vapi CLI"
    fi
    
    log "Extracting..."
    tar -xzf "$tmp_dir/vapi.tar.gz" -C "$tmp_dir"
    
    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/"
    else
        log "Installing to $INSTALL_DIR (requires sudo)..."
        sudo mv "$tmp_dir/$BINARY_NAME" "$INSTALL_DIR/"
    fi
    
    # Make executable
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    # Cleanup
    rm -rf "$tmp_dir"
    
    log "Vapi CLI installed successfully!"
}

# Verify installation
verify_installation() {
    if command -v vapi &> /dev/null; then
        log "Verification: $(vapi --version)"
        echo ""
        log "Installation complete! 🎉"
        echo ""
        echo "Get started with:"
        echo "  vapi login"
        echo "  vapi --help"
    else
        warn "Vapi CLI was installed but not found in PATH"
        warn "You may need to add $INSTALL_DIR to your PATH"
        warn "Or restart your terminal"
    fi
}

# Main installation flow
main() {
    echo "==================================="
    echo "    Vapi CLI Installer"
    echo "==================================="
    echo ""
    
    detect_platform
    get_latest_version
    install_vapi
    verify_installation
}

# Run main function
main "$@" 