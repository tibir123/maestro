# Maestro Development Roadmap

## Project Overview
Building a sophisticated music controller for macOS with Go, following DDD principles and Unix philosophy.

## Phase 1: Foundation (Week 1-2)
**Goal**: Establish core infrastructure and basic functionality

### Week 1: Core Domain & Infrastructure Setup

#### Day 1-2: Domain Layer
- [ ] **Domain Entities** (`domain/music/entities.go`)
  - [ ] Define `Track` entity with ID, Title, Artist, Album, Duration
  - [ ] Define `Playlist` entity with ID, Name, Type, ReadOnly, Tracks
  - [ ] Define `Player` entity with State, Track, Position, Volume, Shuffle, Repeat
  - [ ] Write unit tests for entities

- [ ] **Value Objects** (`domain/music/values.go`)
  - [ ] Implement `TrackID`, `PlaylistID` types
  - [ ] Implement `Duration` with validation (seconds)
  - [ ] Implement `Volume` with validation (0-100)
  - [ ] Define enums: `PlayerState`, `RepeatMode`, `PlaylistType`
  - [ ] Write unit tests for value objects

- [ ] **Repository Interfaces** (`domain/music/repositories.go`)
  - [ ] Define `PlayerRepository` interface (Play, Pause, Stop, Next, Previous, etc.)
  - [ ] Define `LibraryRepository` interface (Search, GetTracks, GetPlaylists)
  - [ ] Define `QueueRepository` interface (GetQueue, AddToQueue, ClearQueue)

- [ ] **Domain Errors** (`domain/music/errors.go`)
  - [ ] Define domain-specific errors (TrackNotFound, InvalidVolume, etc.)
  - [ ] Implement error wrapping utilities

#### Day 3-4: AppleScript Executor & Basic Infrastructure
- [ ] **maestro-exec Implementation** (`cmd/maestro-exec/main.go`)
  - [ ] Simple stdin/stdout AppleScript executor (30 lines)
  - [ ] Error handling on stderr
  - [ ] Basic integration test with Music.app

- [ ] **AppleScript Infrastructure** (`infrastructure/applescript/`)
  - [ ] `executor.go`: Script execution wrapper with timeout
  - [ ] `player.go`: Implement PlayerRepository interface
  - [ ] Basic play/pause/next commands
  - [ ] Create script templates in `scripts/` directory

- [ ] **Logging Setup** (`pkg/logger/`)
  - [ ] Structured JSON logging with logrus
  - [ ] Log levels configuration
  - [ ] Component-specific loggers

#### Day 5-6: Basic Daemon & CLI
- [ ] **maestrod Skeleton** (`cmd/maestrod/main.go`)
  - [ ] Basic service structure
  - [ ] Configuration loading with Viper
  - [ ] Health check endpoint setup
  - [ ] Graceful shutdown handling

- [ ] **Simple CLI** (`cmd/maestro/main.go`)
  - [ ] Cobra command structure
  - [ ] Basic commands: play, pause, next, stop
  - [ ] Direct AppleScript execution (no daemon yet)
  - [ ] Text and JSON output formats

- [ ] **Configuration Files**
  - [ ] Create `configs/maestrod.toml.example`
  - [ ] Create `configs/maestro.toml.example`
  - [ ] Configuration validation

### Week 2: Communication & Session Management

#### Day 7-8: gRPC/WebSocket Setup
- [ ] **Protocol Definitions** (`pkg/protocol/`)
  - [ ] Define protobuf messages for commands and responses
  - [ ] Generate Go code from protobuf
  - [ ] Create shared models package

- [ ] **gRPC Implementation** (`infrastructure/grpc/`)
  - [ ] Server implementation with mTLS
  - [ ] Client implementation
  - [ ] Connection pooling
  - [ ] Benchmark gRPC vs WebSocket performance

- [ ] **Certificate Management** (`certs/`)
  - [ ] Write `generate.sh` script for development certs
  - [ ] CA, server, and client certificate generation
  - [ ] Certificate README with setup instructions

#### Day 9-10: Session Management
- [ ] **Session Manager** (`application/session/`)
  - [ ] Session lifecycle (acquire, release, heartbeat)
  - [ ] Timeout handling (5 minutes default)
  - [ ] Client type differentiation (Human, MCP, Admin)
  - [ ] Takeover mechanism with notification

- [ ] **Integration with Daemon**
  - [ ] Wire up session management in maestrod
  - [ ] Add session commands to CLI
  - [ ] Session status endpoint
  - [ ] Heartbeat implementation

#### Day 11-12: Error Handling & Retry Logic
- [ ] **Retry Mechanism** (`infrastructure/applescript/retry.go`)
  - [ ] Exponential backoff implementation
  - [ ] Silent retries for first 2 attempts
  - [ ] Max 5 retries configuration
  - [ ] Error classification (retryable vs permanent)

- [ ] **Music.app Recovery**
  - [ ] Detect Music.app not responding
  - [ ] Implement restart mechanism (once per session)
  - [ ] Command queueing during recovery

## Phase 2: Core Features (Week 3-4)
**Goal**: Implement essential music control features

### Week 3: Library & Queue Operations

#### Day 13-14: Library Operations
- [ ] **Library Repository** (`infrastructure/applescript/library.go`)
  - [ ] GetAllTracks with pagination
  - [ ] Search implementation with AppleScript
  - [ ] GetPlaylists and GetPlaylistTracks
  - [ ] Performance optimization for large libraries

- [ ] **Caching Layer** (`infrastructure/cache/`)
  - [ ] LRU cache implementation (100MB limit)
  - [ ] Page-based caching for search results
  - [ ] TTL management (5 minutes default)
  - [ ] Cache invalidation strategy

#### Day 15-16: Queue Management
- [ ] **Queue Repository** (`infrastructure/applescript/queue.go`)
  - [ ] GetQueue implementation
  - [ ] AddToQueue, PlayNext, PlayLater
  - [ ] ClearQueue functionality
  - [ ] Queue state synchronization

- [ ] **Application Commands** (`application/commands/`)
  - [ ] PlayTrackCommand handler
  - [ ] QueueManagementCommand handlers
  - [ ] SearchCommand with result formatting

#### Day 17-18: Advanced Player Controls
- [ ] **Extended Player Features**
  - [ ] Volume control with validation
  - [ ] Seek functionality
  - [ ] Shuffle and repeat modes
  - [ ] Get current track position

- [ ] **State Polling** (`infrastructure/applescript/poller.go`)
  - [ ] Adaptive polling intervals (250ms playing, 2s paused, 5s stopped)
  - [ ] State change notifications
  - [ ] Efficient batch status queries

### Week 4: CLI Enhancement & Testing

#### Day 19-20: CLI Improvements
- [ ] **Enhanced CLI** (`presentation/cli/`)
  - [ ] All playback commands
  - [ ] Search with formatted results
  - [ ] Queue management commands
  - [ ] Session management commands
  - [ ] Daemon control (start, stop, status)

- [ ] **REPL Mode** (`presentation/cli/repl.go`)
  - [ ] Interactive command mode
  - [ ] Command history
  - [ ] Tab completion for commands
  - [ ] Contextual help

#### Day 21-22: Testing & Documentation
- [ ] **Unit Tests**
  - [ ] Domain layer tests (100% coverage)
  - [ ] Session management tests
  - [ ] Cache behavior tests
  - [ ] Command handler tests

- [ ] **Integration Tests**
  - [ ] End-to-end command flow
  - [ ] Session timeout scenarios
  - [ ] Error recovery testing
  - [ ] Multi-client scenarios

- [ ] **Documentation**
  - [ ] API.md with all commands
  - [ ] ARCHITECTURE.md details
  - [ ] Basic troubleshooting guide

## Implementation Order for Claude Code

When working with Claude Code, follow this implementation sequence:

### Priority 1: Core Foundation (Start Here)
1. Domain entities and value objects
2. maestro-exec implementation
3. Basic AppleScript player control
4. Simple CLI with direct commands

### Priority 2: Daemon & Communication
1. Basic maestrod structure
2. gRPC server setup
3. CLI to daemon communication
4. Session management

### Priority 3: Essential Features
1. Library search
2. Queue management
3. Caching layer
4. Error handling & retry

### Priority 4: User Interface
1. TUI foundation with Bubble Tea
2. Now Playing view
3. Queue view
4. Search functionality

### Priority 5: Production Features
1. MCP server
2. Installation scripts
3. Performance optimization
4. Documentation

## Success Criteria

### Technical Metrics
- [ ] All commands respond < 100ms
- [ ] Search completes < 500ms for 10K songs
- [ ] Memory usage < 100MB typical
- [ ] CPU usage < 0.1% idle
- [ ] 99% command success rate after retries

### Feature Completeness
- [ ] All Music.app controls working
- [ ] Session management functional
- [ ] TUI all views implemented
- [ ] MCP server operational
- [ ] Installation automated

### Quality Metrics
- [ ] Unit test coverage > 80%
- [ ] Integration tests passing
- [ ] No critical bugs
- [ ] Documentation complete
- [ ] Error handling comprehensive

## Version Milestones

- **v0.1.0**: Basic CLI with play/pause/next
- **v0.2.0**: Daemon with session management
- **v0.3.0**: Search and queue functionality
- **v0.4.0**: Complete TUI
- **v0.5.0**: MCP support
- **v1.0.0**: Production release

---

**Last Updated**: January 2025  
**Version**: 1.0.0  
**Status**: Ready for Implementation
