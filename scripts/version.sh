#!/bin/bash

# Version management script for Vapi CLI
# Usage: 
#   ./scripts/version.sh get          - Show current version
#   ./scripts/version.sh set 1.2.3    - Set new version
#   ./scripts/version.sh bump major   - Bump major version (1.2.3 -> 2.0.0)
#   ./scripts/version.sh bump minor   - Bump minor version (1.2.3 -> 1.3.0)
#   ./scripts/version.sh bump patch   - Bump patch version (1.2.3 -> 1.2.4)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
VERSION_FILE="$PROJECT_ROOT/VERSION"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
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

# Get current version
get_version() {
    if [[ -f "$VERSION_FILE" ]]; then
        cat "$VERSION_FILE"
    else
        echo "0.0.1"
    fi
}

# Set version
set_version() {
    local new_version="$1"
    
    # Validate version format (basic semver check)
    if [[ ! $new_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        log_error "Invalid version format. Use semantic versioning (e.g., 1.2.3)"
        exit 1
    fi
    
    echo "$new_version" > "$VERSION_FILE"
    log_success "Version updated to $new_version"
}

# Bump version
bump_version() {
    local bump_type="$1"
    local current_version
    current_version=$(get_version)
    
    # Parse current version
    IFS='.' read -r major minor patch <<< "$current_version"
    
    case "$bump_type" in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            log_error "Invalid bump type. Use: major, minor, or patch"
            exit 1
            ;;
    esac
    
    local new_version="$major.$minor.$patch"
    set_version "$new_version"
    log_info "Bumped $bump_type version: $current_version → $new_version"
}

# Show usage
show_usage() {
    echo "Usage: $0 <command> [arguments]"
    echo ""
    echo "Commands:"
    echo "  get                 Show current version"
    echo "  set <version>       Set new version (e.g., 1.2.3)"
    echo "  bump major          Bump major version (1.2.3 -> 2.0.0)"
    echo "  bump minor          Bump minor version (1.2.3 -> 1.3.0)"
    echo "  bump patch          Bump patch version (1.2.3 -> 1.2.4)"
    echo ""
    echo "Examples:"
    echo "  $0 get"
    echo "  $0 set 1.2.3"
    echo "  $0 bump patch"
}

# Main script logic
main() {
    case "${1:-}" in
        get)
            echo "Current version: $(get_version)"
            ;;
        set)
            if [[ -z "${2:-}" ]]; then
                log_error "Version number required"
                show_usage
                exit 1
            fi
            set_version "$2"
            ;;
        bump)
            if [[ -z "${2:-}" ]]; then
                log_error "Bump type required (major, minor, or patch)"
                show_usage
                exit 1
            fi
            bump_version "$2"
            ;;
        help|--help|-h)
            show_usage
            ;;
        "")
            log_warning "No command specified"
            show_usage
            exit 1
            ;;
        *)
            log_error "Unknown command: $1"
            show_usage
            exit 1
            ;;
    esac
}

main "$@" 