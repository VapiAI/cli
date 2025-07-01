#!/bin/bash

# NPM Publishing Script for Vapi CLI
set -e

echo "🚀 Publishing Vapi CLI to npm..."

# Check if we're in the right directory
if [[ ! -f "npm-wrapper/package.json" ]]; then
    echo "❌ Error: Run this script from the root of the repository"
    exit 1
fi

# Check if user is logged in to npm
if ! npm whoami > /dev/null 2>&1; then
    echo "❌ Error: Not logged in to npm"
    echo "Run: npm login"
    exit 1
fi

# Get current version from VERSION file
if [[ -f "VERSION" ]]; then
    VERSION=$(cat VERSION)
    echo "📦 Using version: $VERSION"
else
    echo "❌ Error: VERSION file not found"
    exit 1
fi

# Update package.json version
cd npm-wrapper
echo "🔄 Updating package.json version to $VERSION..."
npm version "$VERSION" --no-git-tag-version

# Publish to npm
echo "📤 Publishing to npm..."
npm publish --access public

echo "✅ Successfully published @vapi-ai/cli@$VERSION to npm!"
echo ""
echo "Users can now install with:"
echo "  npm install -g @vapi-ai/cli"
echo ""
echo "Verify the publish at: https://www.npmjs.com/package/@vapi-ai/cli" 