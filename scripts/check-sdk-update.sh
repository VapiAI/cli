#!/bin/bash

# Check for Go SDK updates script
# This script replicates the logic used in the GitHub Actions workflow
# Usage: ./scripts/check-sdk-update.sh [--update]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -f "VERSION" ]; then
    print_error "This script must be run from the project root directory"
    exit 1
fi

print_info "Checking for Vapi Go SDK updates..."

# Get current SDK version from go.mod
current_version=$(grep "github.com/VapiAI/server-sdk-go" go.mod | awk '{print $2}')
print_info "Current SDK version: $current_version"

# Get latest SDK version from GitHub API
print_info "Fetching latest SDK version from GitHub..."
latest_version=$(curl -s https://api.github.com/repos/VapiAI/server-sdk-go/releases/latest | jq -r '.tag_name')

if [ "$latest_version" = "null" ] || [ -z "$latest_version" ]; then
    print_error "Failed to fetch latest SDK version from GitHub API"
    exit 1
fi

print_info "Latest SDK version: $latest_version"

# Compare versions (remove 'v' prefix for comparison)
current_clean=$(echo $current_version | sed 's/^v//')
latest_clean=$(echo $latest_version | sed 's/^v//')

if [ "$current_clean" = "$latest_clean" ]; then
    print_success "CLI is already using the latest SDK version ($current_version)"
    
    if [ "$1" != "--update" ]; then
        print_info "Use --update flag to force update anyway"
    fi
    
    if [ "$1" != "--update" ]; then
        exit 0
    fi
fi

if [ "$current_clean" != "$latest_clean" ]; then
    print_warning "SDK update available: $current_version -> $latest_version"
fi

# If --update flag is provided, perform the update
if [ "$1" = "--update" ]; then
    print_info "Updating Go SDK dependency..."
    
    # Update the dependency
    go get github.com/VapiAI/server-sdk-go@$latest_version
    go mod tidy
    
    print_success "Go SDK updated to $latest_version"
    
    # Run tests
    print_info "Running tests to verify compatibility..."
    if make test; then
        print_success "Tests passed"
    else
        print_error "Tests failed - SDK update may have introduced breaking changes"
        exit 1
    fi
    
    # Run linter
    print_info "Running linter..."
    if make lint; then
        print_success "Linting passed"
    else
        print_error "Linting failed - please fix the issues"
        exit 1
    fi
    
    # Get current CLI version
    current_cli_version=$(cat VERSION | tr -d '\n')
    print_info "Current CLI version: $current_cli_version"
    
    # Parse version components and bump patch version
    IFS='.' read -ra ADDR <<< "$current_cli_version"
    major=${ADDR[0]}
    minor=${ADDR[1]}
    patch=${ADDR[2]}
    
    new_patch=$((patch + 1))
    new_cli_version="$major.$minor.$new_patch"
    
    print_info "Bumping CLI version to: $new_cli_version"
    echo "$new_cli_version" > VERSION
    
    # Update version references in code
    if [ -f "cmd/version.go" ]; then
        sed -i "s/version = \"[^\"]*\"/version = \"$new_cli_version\"/" cmd/version.go
        print_info "Updated version reference in cmd/version.go"
    fi
    
    print_success "Local update completed!"
    print_info "Changes made:"
    print_info "  - Go SDK: $current_version -> $latest_version"
    print_info "  - CLI version: $current_cli_version -> $new_cli_version"
    print_info ""
    print_warning "Don't forget to commit and push these changes:"
    print_info "  git add ."
    print_info "  git commit -m \"chore: update Go SDK to $latest_version\""
    print_info "  git push"
    
else
    print_info ""
    print_info "To update manually, run:"
    print_info "  ./scripts/check-sdk-update.sh --update"
    print_info ""
    print_info "Or wait for the daily automated check at 2 AM UTC"
fi 