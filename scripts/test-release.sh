#!/usr/bin/env bash

# Test GoReleaser configuration locally
# This creates a snapshot release without pushing to GitHub

set -e

echo "Testing GoReleaser configuration..."

# Determine where go install puts binaries
GOBIN=$(go env GOPATH)/bin
if [ -z "$GOBIN" ]; then
    GOBIN="$HOME/go/bin"
fi

# Check if goreleaser is installed
if ! command -v goreleaser &> /dev/null && ! command -v "$GOBIN/goreleaser" &> /dev/null; then
    echo "GoReleaser is not installed. Installing..."
    go install github.com/goreleaser/goreleaser/v2@latest
    
    # Make sure it's accessible
    if [ -f "$GOBIN/goreleaser" ]; then
        echo "GoReleaser installed to $GOBIN"
        export PATH="$GOBIN:$PATH"
    else
        echo "Failed to install GoReleaser"
        exit 1
    fi
fi

# Use goreleaser from PATH or GOBIN
if command -v goreleaser &> /dev/null; then
    GORELEASER="goreleaser"
else
    GORELEASER="$GOBIN/goreleaser"
fi

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf dist/

# Run goreleaser in snapshot mode
echo "Running GoReleaser in snapshot mode..."
$GORELEASER release --snapshot --clean

echo ""
echo "✅ GoReleaser test completed successfully!"
echo ""
echo "Build artifacts are in the 'dist/' directory:"
ls -la dist/

echo ""
echo "To test a specific binary:"
echo "  ./dist/vapi_$(uname -s)_$(uname -m)/vapi --version" 