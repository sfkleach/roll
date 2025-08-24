# Justfile for Roll - Virtual Dice Rolling Application

# Default recipe to display available commands
default:
    @just --list

# Build the application
build:
    go build -v .

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

# Install dependencies for development
install-deps:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Cross-compile for different platforms
build-all:
    GOOS=linux GOARCH=amd64 go build -o bin/roll-linux-amd64 .
    GOOS=windows GOARCH=amd64 go build -o bin/roll-windows-amd64.exe .
    GOOS=darwin GOARCH=amd64 go build -o bin/roll-darwin-amd64 .
    GOOS=darwin GOARCH=arm64 go build -o bin/roll-darwin-arm64 .

# Check for security vulnerabilities
security:
    go list -json -deps ./... | nancy sleuth

# Full CI check (run all quality checks)
ci: fmt lint test-race
    @echo "All checks passed!"
