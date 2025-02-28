# Go API Makefile

# Variables
BINARY_NAME=api-server
GO=go
MAIN_PATH=main.go
BUILD_DIR=build
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
VERSION?=0.1.0
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Default target executed when no arguments are given to make.
default: build

# Build the application
build:
	@echo "Building ${BINARY_NAME}..."
	@${GO} build -o ${BINARY_NAME} ${LDFLAGS}

# Run the application
run: build
	@echo "Running ${BINARY_NAME}..."
	@./${BINARY_NAME}

# Run the application in development mode
dev:
	@echo "Running in development mode..."
	@${GO} run ${MAIN_PATH}

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME}
	@rm -rf ${BUILD_DIR}

# Test the application
test:
	@echo "Running tests..."
	@${GO} test ./... -v

# Create a clean build for release
release: clean
	@echo "Building release version ${VERSION}..."
	@mkdir -p ${BUILD_DIR}
	@${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}

# Cross-compile for different platforms
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p ${BUILD_DIR}

	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64

	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe

	@echo "Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64

	@echo "Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64

# Install development tools
tools:
	@echo "Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
lint: tools
	@echo "Running linter..."
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	@${GO} fmt ./...

# Install the application
install: build
	@echo "Installing ${BINARY_NAME}..."
	@mv ${BINARY_NAME} /usr/local/bin/

# Docker related commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t ${BINARY_NAME}:${VERSION} .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8082:8082 ${BINARY_NAME}:${VERSION}

# Help message
help:
	@echo "Go API Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build      Build the application"
	@echo "  run        Build and run the application"
	@echo "  dev        Run in development mode with go run"
	@echo "  clean      Remove build artifacts"
	@echo "  test       Run tests"
	@echo "  release    Create a clean build for release"
	@echo "  build-all  Cross-compile for different platforms"
	@echo "  tools      Install development tools"
	@echo "  lint       Run linter"
	@echo "  fmt        Format code"
	@echo "  install    Install the application to /usr/local/bin"
	@echo "  docker-build Build Docker image"
	@echo "  docker-run   Run Docker container"
	@echo "  help       Show this help message"

.PHONY: default build run dev clean test release build-all tools lint fmt install docker-build docker-run help 