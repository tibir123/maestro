# Maestro - Complete Project Specification v1.0.0

## Executive Summary

**Maestro** is a sophisticated music controller for macOS that enables command-line and terminal-based control of Music.app. Built with Go following Domain-Driven Design principles and Unix philosophy, Maestro provides a client-server architecture supporting CLI, TUI, and MCP interfaces for local music library control.

**Project Name**: Maestro  
**Command**: `maestro` / `mctl` (short form)  
**Version**: 1.0.0  
**Date**: January 2025  
**Platform**: macOS (latest version only)  
**Language**: Go (latest stable)  
**License**: TBD

## Table of Contents

1. [Core Principles](#core-principles)
2. [Architecture Overview](#architecture-overview)
3. [Technical Decisions](#technical-decisions)
4. [Project Structure](#project-structure)
5. [Component Specifications](#component-specifications)
6. [Configuration](#configuration)
7. [API Specifications](#api-specifications)
8. [Implementation Guidelines](#implementation-guidelines)
9. [Development Phases](#development-phases)
10. [Installation & Usage](#installation--usage)

## Core Principles

- **Unix Philosophy**: Small, composable tools that do one thing well
- **Domain-Driven Design**: Clean separation of business logic from infrastructure
- **Minimal Dependencies**: Pure Go + OS built-ins (AppleScript)
- **Local Library Focus**: Control Music.app's library (no streaming services)
- **Rolling Release Model**: Support latest macOS only, version tags for compatibility
- **User-First Design**: Cider-inspired UX with modern, intuitive controls

## Architecture Overview

### System Architecture

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
                                                  🔊 Audio Out
```

### Component Overview

| Component | Purpose | Binary Name |
|-----------|---------|-------------|
| **Daemon** | Core service managing state and Music.app control | `maestrod` |
| **CLI** | Command-line interface for scripting and quick commands | `maestro` / `mctl` |
| **TUI** | Terminal UI with rich visual interface | `maestro-tui` |
| **MCP** | Model Context Protocol server for AI integration | `maestro-mcp` |
| **Executor** | AppleScript executor for process isolation | `maestro-exec` |

## Technical Decisions

### Core Technology Stack

| Category | Decision | Rationale |
|----------|----------|-----------|
| **Music Control** | AppleScript via `osascript` | OS built-in, no external dependencies |
| **Configuration** | TOML with Viper | Go ecosystem standard, less error-prone than YAML |
| **Communication** | gRPC or WebSocket (benchmark to decide) | Reliable, bidirectional, streaming support |
| **Authentication** | mTLS with certificates | Secure without password management |
| **TUI Framework** | Bubble Tea + Charm tools | Modern Go TUI with great UX |
| **Logging** | Structured JSON | Machine-readable, queryable |

### Session Management Policy

| Client Type | Session Duration | Takeover | Heartbeat | Rate Limit |
|------------|------------------|----------|-----------|------------|
| **Human** (CLI/TUI) | 5 minutes | Yes (with notification) | Required | None |
| **MCP** (AI) | Stateless | No | N/A | 20/minute |
| **Admin** | 10 minutes | Priority | Required | None |

**Session Rules**:
- Timeout after 5 minutes of inactivity (configurable)
- Heartbeat required to maintain session
- Music pauses when session times out
- Takeover requires 5-second waiting period for current owner response
- External control (physical, Siri, etc.) releases session

### Performance Requirements

| Metric | Target | Hard Limit | Degradation Strategy |
|--------|--------|------------|---------------------|
| **Library Size** | 10,000 tracks | 50,000 tracks | Warn if search >1s |
| **Search Response** | <500ms | 1 second | Show partial results |
| **Command Response** | <100ms | 500ms | Queue up to 3 commands |
| **Memory Usage** | ~100MB | 1GB | Drop cache at 80% |
| **CPU Idle** | <0.1% | 1% | Reduce polling frequency |
| **CPU Active** | <10% | 25% | Throttle operations |
| **Concurrent Clients** | 10 | 10 | Reject new connections |

### Error Handling Strategy

```toml
[error_strategy]
max_retries = 5
backoff = "exponential"  # 100ms, 200ms, 400ms, 800ms, 1600ms
silent_retries = 2        # First 2 retries are silent
user_notification = "friendly"  # No technical jargon
auto_recovery = true      # Try to fix known issues
music_app_restart = "once"  # One restart attempt if frozen
command_queue_max = 3     # Queue up to 3 commands during recovery
```

### Platform Support Strategy

- **Model**: Latest macOS only (currently Sonoma 14.x)
- **Updates**: Rolling release, no backward compatibility
- **Versioning**: Semantic versioning with Git tags
- **Compatibility**: Each release tagged with supported macOS version
- **Migration**: Users on older macOS use older tagged releases

## Project Structure

```
maestro/
├── cmd/                        # Entry points (minimal logic)
│   ├── maestro/               # CLI client
│   ├── maestrod/              # Daemon service
│   ├── maestro-tui/           # Terminal UI
│   ├── maestro-mcp/           # MCP server
│   └── maestro-exec/          # AppleScript executor
│
├── domain/                    # Pure business logic (no external deps)
│   └── music/
│       ├── entities.go        # Track, Playlist, Player
│       ├── values.go          # TrackID, Duration, Volume
│       ├── services.go        # MusicService (use cases)
│       ├── repositories.go    # Port interfaces
│       └── errors.go          # Domain-specific errors
│
├── application/               # Application services (orchestration)
│   ├── commands/             # Command handlers (PlayTrack, etc.)
│   ├── queries/              # Query handlers (GetStatus, etc.)
│   └── session/              # Session management logic
│
├── infrastructure/            # External adapters
│   ├── applescript/          # Music.app control implementation
│   │   ├── executor.go       # Script execution
│   │   ├── player.go         # Player adapter
│   │   ├── library.go        # Library adapter
│   │   └── scripts/          # AppleScript templates
│   ├── grpc/                 # gRPC server/client
│   ├── websocket/            # WebSocket alternative
│   ├── cache/                # In-memory caching
│   ├── auth/                 # Certificate management
│   └── mcp/                  # MCP protocol adapter
│
├── presentation/              # User interfaces
│   ├── cli/                  # Command-line interface
│   │   ├── commands.go       # Command definitions
│   │   ├── repl.go          # Interactive mode
│   │   └── output.go        # JSON/text formatting
│   └── tui/                  # Terminal UI
│       ├── models/          # State management
│       ├── views/           # View components
│       ├── components/      # Reusable widgets
│       └── keys.go          # Keyboard handling
│
├── pkg/                       # Shared packages
│   ├── protocol/             # Protocol definitions
│   ├── models/               # Shared data models
│   ├── logger/               # Logging utilities
│   └── health/               # Health check server
│
├── configs/                   # Configuration files
│   ├── maestrod.toml         # Daemon config template
│   ├── maestro.toml          # Client config template
│   └── examples/             # Example configurations
│
├── certs/                     # Certificate management
│   ├── generate.sh           # Certificate generation
│   └── README.md             # Certificate setup guide
│
├── scripts/                   # Operational scripts
│   ├── install.sh            # Installation script
│   ├── maestrod.plist        # launchd service definition
│   └── release.sh            # Release tagging
│
├── docs/                      # Documentation
│   ├── API.md                # API reference
│   ├── ARCHITECTURE.md       # Architecture details
│   └── TROUBLESHOOTING.md    # Common issues
│
├── go.mod                     # Go module definition
├── go.sum                     # Go dependencies
├── Makefile                   # Build automation
├── README.md                  # Project overview
└── PROJECT_SPEC.md           # This document
```

[Rest of the content continues as in the original document...]
