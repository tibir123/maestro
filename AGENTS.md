# AGENTS.md - Coding Agent Guidelines for Maestro

## Build & Test Commands
```bash
make build              # Build all binaries
make test              # Run all tests  
go test -v -run TestName ./path/to/package  # Run single test
make fmt               # Format code (gofmt + go fmt)
make lint              # Run golangci-lint
make test-unit         # Unit tests only (-short flag)
make test-integration  # Integration tests only
```

## Code Style & Conventions
- **Architecture**: Follow DDD with clean separation - domain/ (pure logic), application/ (use cases), infrastructure/ (adapters), presentation/ (UI)
- **Imports**: Group as stdlib, external deps, internal packages with blank lines between groups
- **Error Handling**: Return errors, don't panic; wrap with context; use domain-specific errors
- **Naming**: Use Go conventions - exported CamelCase, unexported camelCase, interfaces end with -er suffix
- **Testing**: Use testify/assert for assertions; mock interfaces for unit tests; _test.go suffix
- **Dependencies**: Check go.mod before adding libraries; use existing patterns from codebase
- **Performance**: Command response <100ms target, 500ms hard limit; cache when appropriate
- **Logging**: Use logrus with structured JSON; avoid fmt.Print
- **No Comments**: Don't add code comments unless explicitly requested