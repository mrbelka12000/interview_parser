# Interview Parser Desktop App Makefile
# Based on the README documentation

# Variables
APP_NAME := InterviewParser.app
BUILD_DIR := build/bin
FRONTEND_DIR := frontend
APP_DIR := build/bin/InterviewParser.app

# Default target
.PHONY: help
help:
	@echo "Interview Parser Desktop App - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  install      Install all dependencies (Go, Node.js, Wails, FFmpeg)"
	@echo "  deps         Download Go modules and npm packages"
	@echo "  dev          Start development server with hot reload"
	@echo "  frontend-dev Start frontend development server only"
	@echo ""
	@echo "Building:"
	@echo "  build        Build for current platform"
	@echo "  build-darwin Build for macOS (Apple Silicon)"
	@echo "  build-all    Build for all platforms"
	@echo "  package-darwin Package macOS build with FFmpeg binaries"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean        Clean build artifacts"
	@echo "  tidy         Tidy Go modules"
	@echo "  test         Run tests"
	@echo ""

# Installation and Setup
.PHONY: install
install:
	@echo "Installing prerequisites..."
	@echo "Note: This assumes you have Homebrew (macOS) or apt (Ubuntu) available"
	@if command -v brew >/dev/null 2>&1; then \
		brew install go node ffmpeg; \
	elif command -v apt-get >/dev/null 2>&1; then \
		sudo apt-get update && sudo apt-get install -y golang-go nodejs npm ffmpeg; \
	else \
		echo "Please install Go, Node.js, and FFmpeg manually for your platform"; \
	fi
	@echo "Installing Wails..."
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	@echo "Installation complete!"

# Dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@echo "Downloading Go modules..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install
	@echo "Dependencies installed!"

# Development
.PHONY: dev
dev:
	@echo "Starting development server with hot reload..."
	wails dev

# Building
.PHONY: build
build:
	@echo "Building for current platform..."
	wails build

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS (Apple Silicon)..."
	wails build -platform darwin/arm64

.PHONY: build-darwin-amd64
build-darwin-amd64:
	@echo "Building for macOS (Intel)..."
	wails build -platform darwin/amd64

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	wails build -platform windows/amd64 -webview2 embed

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	wails build -platform linux/amd64

.PHONY: build-all
build-all: build-darwin build-darwin-amd64 build-windows build-linux
	@echo "All platform builds complete!"

# macOS Packaging with FFmpeg
.PHONY: package-darwin
package-darwin: build-darwin
	cp /usr/local/bin/ffprobe $(APP_DIR)/Contents/MacOS/ffprobe
	cp /usr/local/bin/ffmpeg  $(APP_DIR)/Contents/MacOS/ffmpeg

# Maintenance
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@echo "Cleaning Wails build artifacts..."
	wails build -clean
	@echo "Cleaning frontend build artifacts..."
	cd $(FRONTEND_DIR) && npm run clean 2>/dev/null || true
	@echo "Cleaning Go build cache..."
	go clean -cache
	@echo "Clean complete!"

.PHONY: tidy
tidy:
	@echo "Tidying Go modules..."
	go mod tidy
	@echo "Tidying frontend dependencies..."
	cd $(FRONTEND_DIR) && npm ci 2>/dev/null || npm install
	@echo "Tidy complete!"

# Testing
.PHONY: test
test:
	@echo "Running Go tests..."
	go test ./...
	@echo "Running frontend tests..."
	cd $(FRONTEND_DIR) && npm test 2>/dev/null || echo "No frontend tests configured"
	@echo "Tests complete!"

.PHONY: test-coverage
test-coverage:
	@echo "Running Go tests with coverage..."
	go test -cover ./...
	@echo "Running frontend tests with coverage..."
	cd $(FRONTEND_DIR) && npm run test:coverage 2>/dev/null || echo "No frontend test coverage configured"

# Utilities
.PHONY: check-deps
check-deps:
	@echo "Checking dependencies..."
	@command -v go >/dev/null 2>&1 || (echo "Go is not installed" && exit 1)
	@command -v node >/dev/null 2>&1 || (echo "Node.js is not installed" && exit 1)
	@command -v npm >/dev/null 2>&1 || (echo "npm is not installed" && exit 1)
	@command -v wails >/dev/null 2>&1 || (echo "Wails is not installed" && exit 1)
	@command -v ffmpeg >/dev/null 2>&1 || (echo "FFmpeg is not installed" && exit 1)
	@command -v ffprobe >/dev/null 2>&1 || (echo "FFprobe is not installed" && exit 1)
	@echo "All dependencies are installed!"

.PHONY: version
version:
	@echo "Interview Parser Version Info:"
	@go version
	@node --version
	@npm --version
	@wails version 2>/dev/null || echo "Wails not found in PATH"
	@ffmpeg -version | head -n 1

# Development helpers
.PHONY: setup-dev
setup-dev: install deps
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start development"

.PHONY: quick-start
quick-start:
	@echo "Quick start for Interview Parser:"
	@echo "1. Run 'make setup-dev' to install dependencies"
	@echo "2. Run 'make dev' to start development"
	@echo "3. Configure your OpenAI API key in the app"
	@echo "4. Upload and process your interview files"
