# Vapi CLI Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Node.js parameters for MCP server
NPM=npm
MCP_DIR=mcp-docs-server
MCP_DIST=$(MCP_DIR)/dist

# Binary details
BINARY_NAME=vapi
BUILD_DIR=build
BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)

# Version information
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags - note the lowercase variable names to match main.go
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT) -X main.date=$(BUILD_TIME) -X main.builtBy=make"

# Default target - build both CLI and MCP server
all: test build build-mcp

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Build the CLI binary
build: $(BUILD_DIR)
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) -v

# Build MCP server
build-mcp:
	@echo "Building MCP docs server..."
	@if [ ! -d "$(MCP_DIR)/node_modules" ]; then \
		echo "Installing MCP server dependencies..."; \
		cd $(MCP_DIR) && $(NPM) install; \
	fi
	@cd $(MCP_DIR) && $(NPM) run build
	@echo "✅ MCP server built successfully"

# Install MCP server dependencies
mcp-deps:
	@echo "Installing MCP server dependencies..."
	@cd $(MCP_DIR) && $(NPM) install

# Clean MCP server
clean-mcp:
	@echo "Cleaning MCP server..."
	@rm -rf $(MCP_DIR)/dist
	@rm -rf $(MCP_DIR)/node_modules

# Test MCP server
test-mcp:
	@echo "Testing MCP server..."
	@cd $(MCP_DIR) && $(NPM) test

# Lint MCP server
lint-mcp:
	@echo "Linting MCP server..."
	@cd $(MCP_DIR) && $(NPM) run lint

# Publish MCP server to npm
publish-mcp:
	@echo "Publishing MCP server to npm..."
	@cd $(MCP_DIR) && $(NPM) publish

# Build for all platforms
build-all: $(BUILD_DIR) build-mcp
	@echo "Building for all platforms..."
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run all tests (CLI + MCP server)
test-all: test test-mcp

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf man
	rm -f coverage.out coverage.html

# Clean everything (CLI + MCP server)
clean-all: clean clean-mcp

# Run go mod tidy
tidy:
	@echo "Tidying modules..."
	$(GOMOD) tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Install all dependencies (CLI + MCP server)
deps-all: deps mcp-deps

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

# Run all linters (CLI + MCP server)
lint-all: lint lint-mcp

# Generate manual pages
man-pages: build
	@echo "Generating manual pages..."
	@mkdir -p man
	@$(BINARY_PATH) manual --output ./man
	@echo "✅ Manual pages generated in ./man/"

# Install the binary and manual pages locally
install: build man-pages
	@echo "Installing $(BINARY_NAME)..."
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY_PATH) $(HOME)/.local/bin/
	@echo "Installing manual pages..."
	@mkdir -p $(HOME)/.local/share/man/man1
	@cp ./man/*.1 $(HOME)/.local/share/man/man1/ 2>/dev/null || true
	@command -v mandb >/dev/null 2>&1 && mandb -q $(HOME)/.local/share/man || true
	@echo "Installed to $(HOME)/.local/bin/$(BINARY_NAME)"

# Install MCP server globally
install-mcp: build-mcp
	@echo "Installing MCP server globally..."
	@cd $(MCP_DIR) && $(NPM) install -g .
	@echo "✅ MCP server installed globally"

# Install everything
install-all: install install-mcp

# Run the binary
run: build
	$(BINARY_PATH)

# Version management
.PHONY: version version-get version-set version-bump-major version-bump-minor version-bump-patch
version: version-get

version-get:
	@./scripts/version.sh get

version-set:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required. Usage: make version-set VERSION=1.2.3"; \
		exit 1; \
	fi
	@./scripts/version.sh set $(VERSION)

version-bump-major:
	@./scripts/version.sh bump major

version-bump-minor:
	@./scripts/version.sh bump minor

version-bump-patch:
	@./scripts/version.sh bump patch

# Help target to show available commands
help:
	@echo "🚀 Vapi CLI & MCP Server Build System"
	@echo ""
	@echo "📦 CLI Commands:"
	@echo "  build              Build the CLI binary"
	@echo "  man-pages          Generate Unix manual pages"
	@echo "  install            Install the CLI and manual pages to ~/.local/"
	@echo "  test               Run CLI tests"
	@echo "  lint               Run CLI linters"
	@echo "  clean              Clean CLI build artifacts"
	@echo ""
	@echo "🔧 MCP Server Commands:"
	@echo "  build-mcp          Build the MCP docs server"
	@echo "  install-mcp        Install MCP server globally"
	@echo "  test-mcp           Run MCP server tests"
	@echo "  lint-mcp           Run MCP server linters" 
	@echo "  clean-mcp          Clean MCP server artifacts"
	@echo "  publish-mcp        Publish MCP server to npm"
	@echo ""
	@echo "🎯 Combined Commands:"
	@echo "  all                Build both CLI and MCP server"
	@echo "  build-all          Build CLI for all platforms + MCP server"
	@echo "  test-all           Run all tests (CLI + MCP server)"
	@echo "  lint-all           Run all linters (CLI + MCP server)"
	@echo "  clean-all          Clean everything"
	@echo "  deps-all           Install all dependencies"
	@echo "  install-all        Install CLI and MCP server"
	@echo ""
	@echo "📋 Version management:"
	@echo "  version            Show current version"
	@echo "  version-set        Set version (requires VERSION=x.y.z)"
	@echo "  version-bump-major Bump major version (1.2.3 -> 2.0.0)"
	@echo "  version-bump-minor Bump minor version (1.2.3 -> 1.3.0)"
	@echo "  version-bump-patch Bump patch version (1.2.3 -> 1.2.4)"
	@echo ""
	@echo "Examples:"
	@echo "  make all                    # Build everything"
	@echo "  make install-all            # Install CLI + MCP server"
	@echo "  make version-set VERSION=1.2.3"
	@echo "  make publish-mcp            # Publish MCP server to npm"

.PHONY: all build build-mcp build-all test test-mcp test-all test-coverage clean clean-mcp clean-all tidy deps mcp-deps deps-all lint lint-mcp lint-all man-pages install install-mcp install-all run publish-mcp help 