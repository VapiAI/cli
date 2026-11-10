# Vapi CLI Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

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

# Default target - test and build the CLI
all: test build

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Build the CLI binary
build: $(BUILD_DIR)
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) -v

# Build for all platforms
build-all: $(BUILD_DIR)
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

# Alias for all CLI tests
test-all: test

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

# Alias for cleaning CLI artifacts
clean-all: clean

# Run go mod tidy
tidy:
	@echo "Tidying modules..."
	$(GOMOD) tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Alias for CLI dependencies
deps-all: deps

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

# Alias for CLI linting
lint-all: lint

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

# Alias for installing the CLI and manual pages
install-all: install

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
	@echo "🚀 Vapi CLI Build System"
	@echo ""
	@echo "📦 CLI Commands:"
	@echo "  build              Build the CLI binary"
	@echo "  man-pages          Generate Unix manual pages"
	@echo "  install            Install the CLI and manual pages to ~/.local/"
	@echo "  test               Run CLI tests"
	@echo "  lint               Run CLI linters"
	@echo "  clean              Clean CLI build artifacts"
	@echo ""
	@echo "📋 Version management:"
	@echo "  version            Show current version"
	@echo "  version-set        Set version (requires VERSION=x.y.z)"
	@echo "  version-bump-major Bump major version (1.2.3 -> 2.0.0)"
	@echo "  version-bump-minor Bump minor version (1.2.3 -> 1.3.0)"
	@echo "  version-bump-patch Bump patch version (1.2.3 -> 1.2.4)"
	@echo ""
	@echo "Examples:"
	@echo "  make all                    # Test and build the CLI"
	@echo "  make install-all            # Install CLI and manual pages"
	@echo "  make version-set VERSION=1.2.3"

.PHONY: all build build-all test test-all test-coverage clean clean-all tidy deps deps-all lint lint-all man-pages install install-all run help
