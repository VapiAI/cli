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

# Default target
all: test build

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Build the binary
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
	rm -f coverage.out coverage.html

# Run go mod tidy
tidy:
	@echo "Tidying modules..."
	$(GOMOD) tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

# Install the binary locally
install: build
	@echo "Installing $(BINARY_NAME)..."
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY_PATH) $(HOME)/.local/bin/
	@echo "Installed to $(HOME)/.local/bin/$(BINARY_NAME)"

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
	@echo "Available targets:"
	@echo "  build              Build the CLI binary"
	@echo "  install            Install the CLI to ~/.local/bin"
	@echo "  test               Run all tests"
	@echo "  lint               Run linters"
	@echo "  clean              Clean build artifacts"
	@echo ""
	@echo "Version management:"
	@echo "  version            Show current version"
	@echo "  version-set        Set version (requires VERSION=x.y.z)"
	@echo "  version-bump-major Bump major version (1.2.3 -> 2.0.0)"
	@echo "  version-bump-minor Bump minor version (1.2.3 -> 1.3.0)"
	@echo "  version-bump-patch Bump patch version (1.2.3 -> 1.2.4)"
	@echo ""
	@echo "Examples:"
	@echo "  make version-set VERSION=1.2.3"
	@echo "  make version-bump-patch"

.PHONY: all build build-all test test-coverage clean tidy deps lint install run help 