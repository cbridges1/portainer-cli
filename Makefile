# Portainer CLI Makefile

.PHONY: build test test-verbose test-coverage clean lint fmt help

# Default target
help:
	@echo "Available targets:"
	@echo "  build         - Build the CLI binary"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  lint          - Run linter (requires golangci-lint)"
	@echo "  fmt           - Format code"
	@echo "  help          - Show this help"

# Build the CLI binary
build:
	@echo "Building portainer-cli..."
	go build -o portainer-cli

# Run all tests
test:
	@echo "Running tests..."
	go test ./... -short

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test ./... -v -short

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -short -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific package tests
test-storage:
	@echo "Running storage package tests..."
	go test ./pkg/storage/... -v

test-client:
	@echo "Running client package tests..."
	go test ./pkg/client/... -v

test-cmd:
	@echo "Running cmd package tests..."
	go test ./cmd/... -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f portainer-cli
	rm -f coverage.out
	rm -f coverage.html

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with:"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.50.1"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run integration tests (requires running Portainer instance)
test-integration:
	@echo "Running integration tests..."
	@echo "Note: These tests require a running Portainer instance"
	go test ./test/... -tags=integration -v

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o dist/portainer-cli-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o dist/portainer-cli-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/portainer-cli-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o dist/portainer-cli-windows-amd64.exe

# Create dist directory
dist:
	mkdir -p dist

# Release build (optimized)
release: dist
	@echo "Building release binary..."
	go build -ldflags="-s -w" -o dist/portainer-cli

# Run security scan
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install it with:"; \
		echo "  go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Run all quality checks
quality: fmt lint test security
	@echo "All quality checks completed"