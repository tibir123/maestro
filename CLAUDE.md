# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Maestro is a sophisticated music controller for macOS that provides command-line and terminal-based control of Music.app. The project follows Domain-Driven Design (DDD) principles with clean architecture separation and uses Go as the primary language.

## Build Commands

```bash
# Build all components
make build

# Build individual binaries
make build-maestro        # CLI client
make build-maestrod       # Daemon service  
make build-maestro-tui    # Terminal UI
make build-maestro-mcp    # MCP server
make build-maestro-exec   # AppleScript executor

# Development helpers
make run-daemon          # Build and run daemon
make run-cli ARGS="..."  # Build and run CLI with arguments
make run-tui            # Build and run TUI

# Testing
make test               # Run all tests
make test-unit         # Unit tests only
make test-integration  # Integration tests only

# Code quality
make fmt               # Format code
make lint              # Run golangci-lint
make security          # Run govulncheck for vulnerabilities

# Other utilities
make deps              # Download dependencies
make clean             # Clean build artifacts
make install           # Install to /usr/local/bin
make certs             # Generate development certificates
```

## Architecture

### Layer Structure
- **Domain** (`domain/`): Pure business logic with no external dependencies
  - Entities: Track, Playlist, Player
  - Value objects: TrackID, Duration, Volume
  - Repository interfaces
  - Domain-specific errors

- **Application** (`application/`): Use cases and orchestration
  - Commands: Command handlers (PlayTrack, etc.)
  - Queries: Query handlers (GetStatus, etc.)
  - Session: Session management logic

- **Infrastructure** (`infrastructure/`): External adapters
  - AppleScript: Music.app control via osascript
  - gRPC/WebSocket: Client-server communication with mTLS
  - Cache: In-memory caching for performance
  - Auth: Certificate management

- **Presentation** (`presentation/`): User interfaces
  - CLI: Command-line interface with Cobra
  - TUI: Terminal UI using Bubble Tea

### Component Architecture
```
┌─────────────┐         gRPC/WebSocket        ┌──────────────┐
│   Clients   │ ◄───────────────────────────► │   maestrod   │
│ CLI/TUI/MCP │         with mTLS              │   (macOS)    │
└─────────────┘                                └──────┬───────┘
                                                       │
                                                  AppleScript
                                                       │
                                                       ▼
                                                  Music.app
```

## Key Implementation Guidelines

### Go Development Best Practices
- Always use the go-docs-coder agent when coding in Go (per global CLAUDE.md)
- Follow DDD principles: keep domain layer pure, no external dependencies
- Use dependency injection and interfaces for testability
- Maintain clean separation between layers

### Session Management
- Human clients (CLI/TUI): 5-minute timeout with takeover support
- MCP clients: Stateless with 20/minute rate limit  
- Admin clients: 10-minute timeout with priority access
- Heartbeat required to maintain sessions

### Error Handling Strategy
- Max 5 retries with exponential backoff
- Silent retries for first 2 attempts
- One Music.app restart attempt if frozen
- Queue up to 3 commands during recovery

### Performance Targets
- Command response: <100ms (hard limit 500ms)
- Search response: <500ms for 10K tracks
- Memory usage: ~100MB target, 1GB limit
- CPU idle: <0.1%, active <10%

## Development Workflow

1. **Start with Domain Layer**: Define entities and interfaces first
2. **AppleScript Integration**: Use `maestro-exec` for isolated script execution
3. **Build Infrastructure**: Implement repository interfaces
4. **Add Application Layer**: Command and query handlers
5. **Create Presentation**: CLI/TUI interfaces

## Certificate Management

Development certificates are required for mTLS communication:
```bash
./certs/generate.sh  # Generate dev certificates
```

## Testing Requirements

- Unit tests for domain layer (target 100% coverage)
- Integration tests for AppleScript interaction
- Session management scenarios
- Multi-client behavior testing
- Error recovery testing

## Configuration

- Uses TOML configuration files with Viper
- Templates in `configs/` directory
- Environment variable support
- Structured JSON logging with logrus

## Current Implementation Status

### ✅ Phase 1 Complete (v0.1.0)
**Core Foundation implemented and tested:**

- **Domain Layer**: Complete with entities, value objects, repository interfaces, and domain errors (95.3% test coverage)
- **AppleScript Infrastructure**: Full Music.app integration with process isolation via maestro-exec
- **CLI Interface**: All basic commands (play, pause, next, previous, volume, status) with JSON output support
- **Logging Infrastructure**: Production-ready structured logging with component-specific loggers
- **Build System**: Makefile + GoReleaser for local and release builds

**Immediately Usable:**
```bash
make build
./bin/maestro play
./bin/maestro status --json
./bin/maestro volume 75
```

### 🚧 Next Phase (v0.2.0) - Daemon & Communication
**Ready to implement:**
- Basic maestrod daemon structure
- gRPC/WebSocket communication layer
- Session management with mTLS authentication
- Client-server protocol implementation

### 🎯 Implementation Notes for Future Sessions

1. **Start with Application Layer**: Command and query handlers in `application/`
2. **Then Infrastructure**: gRPC server implementation in `infrastructure/grpc/`
3. **Session Management**: Multi-client coordination in `application/session/`
4. **Certificate System**: mTLS setup in `infrastructure/auth/`

All foundation components are solid and tested. The domain interfaces are ready for daemon implementation without any changes needed.