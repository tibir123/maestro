# Build configuration
BINARY_DIR := bin
CMD_DIR := cmd
GOOS := darwin
GOARCH := arm64

# Binary names
MAESTRO := maestro
MAESTROD := maestrod
MAESTRO_TUI := maestro-tui
MAESTRO_MCP := maestro-mcp
MAESTRO_EXEC := maestro-exec

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -ldflags "-s -w"

.PHONY: all build clean test deps install uninstall

all: deps build

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build all binaries (Phase 1: only CLI and executor)
build: build-maestro build-maestro-exec

# Build all binaries (future - when daemon is implemented)
build-all: build-maestro build-maestrod build-maestro-tui build-maestro-mcp build-maestro-exec

build-maestro:
	@echo "Building maestro CLI..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(MAESTRO) ./$(CMD_DIR)/$(MAESTRO)

build-maestrod:
	@echo "Building maestrod daemon..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(MAESTROD) ./$(CMD_DIR)/$(MAESTROD)

build-maestro-tui:
	@echo "Building maestro-tui..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(MAESTRO_TUI) ./$(CMD_DIR)/$(MAESTRO_TUI)

build-maestro-mcp:
	@echo "Building maestro-mcp..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(MAESTRO_MCP) ./$(CMD_DIR)/$(MAESTRO_MCP)

build-maestro-exec:
	@echo "Building maestro-exec..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(MAESTRO_EXEC) ./$(CMD_DIR)/$(MAESTRO_EXEC)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...

test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v -race -coverprofile=coverage-unit.out -covermode=atomic -short ./...

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -race -coverprofile=coverage-integration.out -covermode=atomic -run Integration ./...

# Test coverage analysis
test-coverage: test
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-check: test
	@echo "Checking test coverage..."
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Coverage: $${coverage}%"; \
	if [ $$(echo "$${coverage} < 80.0" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $${coverage}% is below the required 80%"; \
		exit 1; \
	else \
		echo "✅ Coverage $${coverage}% meets the requirement"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

# Install binaries to /usr/local/bin
install: build
	@echo "Installing binaries..."
	@sudo cp $(BINARY_DIR)/$(MAESTRO) /usr/local/bin/
	@sudo cp $(BINARY_DIR)/$(MAESTROD) /usr/local/bin/
	@sudo cp $(BINARY_DIR)/$(MAESTRO_TUI) /usr/local/bin/
	@sudo cp $(BINARY_DIR)/$(MAESTRO_MCP) /usr/local/bin/
	@sudo cp $(BINARY_DIR)/$(MAESTRO_EXEC) /usr/local/bin/
	@echo "Creating mctl symlink..."
	@sudo ln -sf /usr/local/bin/$(MAESTRO) /usr/local/bin/mctl
	@echo "Installation complete!"

# Uninstall binaries
uninstall:
	@echo "Uninstalling binaries..."
	@sudo rm -f /usr/local/bin/$(MAESTRO)
	@sudo rm -f /usr/local/bin/$(MAESTROD)
	@sudo rm -f /usr/local/bin/$(MAESTRO_TUI)
	@sudo rm -f /usr/local/bin/$(MAESTRO_MCP)
	@sudo rm -f /usr/local/bin/$(MAESTRO_EXEC)
	@sudo rm -f /usr/local/bin/mctl
	@echo "Uninstall complete!"

# Development helpers
run-daemon:
	@$(GOBUILD) -o $(BINARY_DIR)/$(MAESTROD) ./$(CMD_DIR)/$(MAESTROD)
	@$(BINARY_DIR)/$(MAESTROD)

run-cli:
	@$(GOBUILD) -o $(BINARY_DIR)/$(MAESTRO) ./$(CMD_DIR)/$(MAESTRO)
	@$(BINARY_DIR)/$(MAESTRO) $(ARGS)

run-tui:
	@$(GOBUILD) -o $(BINARY_DIR)/$(MAESTRO_TUI) ./$(CMD_DIR)/$(MAESTRO_TUI)
	@$(BINARY_DIR)/$(MAESTRO_TUI)

# Generate certificates for development
certs:
	@echo "Generating development certificates..."
	@./certs/generate.sh

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .
	@go fmt ./...

# Run linters
lint:
	@echo "Running linters..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	fi
	@golangci-lint run --timeout=5m
	@echo "Running go vet..."
	@go vet ./...

# Check code formatting
fmt-check:
	@echo "Checking code formatting..."
	@gofmt -l . | tee /tmp/gofmt-output
	@test ! -s /tmp/gofmt-output || (echo "❌ Code is not formatted. Run 'make fmt' to fix." && exit 1)
	@echo "✅ Code formatting is correct"

# Check for security vulnerabilities
security:
	@echo "Checking for vulnerabilities..."
	@if ! command -v govulncheck &> /dev/null; then \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	@govulncheck ./...

# Pre-commit checks - run before committing
pre-commit: fmt-check lint test-coverage-check security
	@echo "✅ All pre-commit checks passed!"

# CI checks - what CI runs
ci: deps lint test-coverage-check security build
	@echo "✅ All CI checks passed!"

# Quick checks for development
quick-check: fmt-check test-unit
	@echo "✅ Quick checks passed!"

# Generate documentation
docs:
	@echo "Generating documentation..."
	@godoc -http=:6060

# Show help
help:
	@echo "Maestro - Build Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Build Targets:"
	@echo "  all              Build everything (default)"
	@echo "  deps             Download dependencies"
	@echo "  build            Build all binaries"
	@echo "  clean            Clean build artifacts"
	@echo "  install          Install binaries to /usr/local/bin"
	@echo "  uninstall        Remove installed binaries"
	@echo ""
	@echo "Testing Targets:"
	@echo "  test             Run all tests with coverage"
	@echo "  test-unit        Run unit tests only"
	@echo "  test-integration Run integration tests only"
	@echo "  test-coverage    Generate HTML coverage report"
	@echo "  test-coverage-check  Check coverage meets 80% threshold"
	@echo ""
	@echo "Quality Targets:"
	@echo "  lint             Run linters and static analysis"
	@echo "  fmt              Format code"
	@echo "  fmt-check        Check code formatting"
	@echo "  security         Check for security vulnerabilities"
	@echo "  pre-commit       Run all pre-commit checks"
	@echo "  ci               Run all CI checks"
	@echo "  quick-check      Run quick development checks"
	@echo ""
	@echo "Development Targets:"
	@echo "  run-daemon       Build and run daemon"
	@echo "  run-cli          Build and run CLI (use ARGS=...)"
	@echo "  run-tui          Build and run TUI"
	@echo "  certs            Generate development certificates"
	@echo "  docs             Start documentation server"
	@echo ""
	@echo "Examples:"
	@echo "  make pre-commit     # Run before committing"
	@echo "  make quick-check    # Quick dev feedback"
	@echo "  make test-coverage  # Generate coverage report"
	@echo "  make run-cli ARGS='status --json'  # Run CLI with args"
