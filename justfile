# Justfile for Roll - Virtual Dice Rolling Application

# Default recipe to display available commands
default:
    @just --list

# Build the application
build version="":
    @if [ -n "{{ version }}" ]; then \
        ./scripts/build.sh "{{ version }}"; \
    else \
        ./scripts/build.sh; \
    fi

# Build with automatic version detection from git
build-release:
    ./scripts/build.sh

# Run tests
test:
    go test -v ./...

# Run tests with race detection
test-race:
    go test -race -v ./internal/...

# Run tests with coverage
test-coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

# Run the application
run:
    go run .

# Format code
fmt:
    go fmt ./...

# Run linter
lint:
    golangci-lint run

# Tidy dependencies
tidy:
    go mod tidy

# Clean build artifacts
clean:
    rm -f roll
    rm -f coverage.out coverage.html
    rm -rf bin/

# Install dependencies for development
install-deps:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Cross-compile for different platforms (Note: may fail locally due to CGO dependencies)
# Use the GitHub release workflow for proper cross-compilation
build-all version="":
    #!/bin/bash
    VERSION="{{ version }}"
    if [ -z "$VERSION" ]; then
        if git describe --tags --exact-match HEAD 2>/dev/null; then
            VERSION=$(git describe --tags --exact-match HEAD | sed 's/^v//')
        elif git describe --tags 2>/dev/null; then
            VERSION=$(git describe --tags | sed 's/^v//')
        else
            VERSION="dev"
        fi
    fi
    echo "Building all platforms for version: $VERSION"
    echo "Note: Cross-compilation may fail locally due to CGO dependencies."
    echo "Use the GitHub release workflow for proper cross-platform builds."
    mkdir -p bin
    echo "Building Linux binary..."
    GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/sfkleach/roll/internal/info.Version=$VERSION" -o bin/roll-linux-amd64 . || echo "Linux build failed (expected if not on Linux)"
    echo "Building Windows binary..."
    GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/sfkleach/roll/internal/info.Version=$VERSION" -o bin/roll-windows-amd64.exe . || echo "Windows build failed (expected due to CGO)"
    echo "Building macOS AMD64 binary..."
    GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/sfkleach/roll/internal/info.Version=$VERSION" -o bin/roll-darwin-amd64 . || echo "macOS AMD64 build failed (expected due to CGO)"
    echo "Building macOS ARM64 binary..."
    GOOS=darwin GOARCH=arm64 go build -ldflags "-X github.com/sfkleach/roll/internal/info.Version=$VERSION" -o bin/roll-darwin-arm64 . || echo "macOS ARM64 build failed (expected due to CGO)"

# Check for security vulnerabilities
security:
    go list -json -deps ./... | nancy sleuth

# Full CI check (run all quality checks)
ci: fmt lint test-race
    @echo "All checks passed!"
