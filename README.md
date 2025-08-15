# Maestro

[![CI](https://github.com/madstone-tech/maestro/actions/workflows/ci.yml/badge.svg)](https://github.com/madstone-tech/maestro/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/madstone-tech/maestro)](https://goreportcard.com/report/github.com/madstone-tech/maestro)
[![codecov](https://codecov.io/gh/madstone-tech/maestro/branch/main/graph/badge.svg)](https://codecov.io/gh/madstone-tech/maestro/branch/main)

Maestro is a sophisticated music controller for macOS that provides seamless control over Apple Music through multiple interfaces. Built with Go and designed for performance, security, and extensibility.

## Features

- 🎵 **Complete Music Control**: Play, pause, skip, volume, shuffle, repeat
- 🖥️ **Multiple Interfaces**: CLI, TUI, gRPC API, MCP server
- 🔒 **Enterprise Security**: mTLS authentication, secure session management
- 🚀 **High Performance**: Concurrent operations, efficient AppleScript integration
- 📊 **Comprehensive Logging**: Structured JSON logging with context tracing
- 🧪 **Well Tested**: 95%+ test coverage with integration tests

## Quick Start

### Prerequisites

- macOS 10.15+ (Catalina or later)
- Apple Music.app installed
- Go 1.25+ (for building from source)

### Installation

#### From Source

```bash
git clone https://github.com/madstone-tech/maestro.git
cd maestro
make build
sudo make install
```

#### Using Homebrew (Coming Soon)

```bash
brew install madstone-tech/tap/maestro
```

### Basic Usage

```bash
# Control music playback
maestro play
maestro pause
maestro next
maestro volume 75

# Get player status
maestro status
maestro status --json

# Interactive TUI (Coming Soon)
maestro-tui

# Start daemon for API access (Coming Soon)
maestrod
```

## Architecture

Maestro follows Domain-Driven Design principles with clean architecture:

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Presentation  │  │   Application   │  │  Infrastructure │
│                 │  │                 │  │                 │
│ • CLI (Cobra)   │◄─┤ • Use Cases     │◄─┤ • AppleScript   │
│ • TUI (Bubble)  │  │ • Services      │  │ • gRPC Server   │
│ • gRPC API      │  │ • Session Mgmt  │  │ • Logging       │
│ • MCP Server    │  │                 │  │ • Authentication│
└─────────────────┘  └─────────────────┘  └─────────────────┘
                              │
                    ┌─────────────────┐
                    │     Domain      │
                    │                 │
                    │ • Entities      │
                    │ • Value Objects │
                    │ • Repositories  │
                    │ • Domain Logic  │
                    └─────────────────┘
```

### Key Components

- **Domain Layer**: Core business logic, entities, and value objects
- **Application Layer**: Use cases, services, and application logic
- **Infrastructure Layer**: External integrations (AppleScript, gRPC, logging)
- **Presentation Layer**: User interfaces (CLI, TUI, API)

## Current Status (Phase 1 Complete)

✅ **Foundation Layer**
- Complete domain model with 95.3% test coverage
- AppleScript infrastructure with timeout handling
- Basic CLI with essential commands
- Structured JSON logging system
- CI/CD pipeline with quality gates

🚧 **In Development** (See [GitHub Issues](https://github.com/madstone-tech/maestro/issues))
- Phase 2: TUI interface with Bubble Tea
- Phase 3: gRPC API server with mTLS
- Phase 4: MCP (Model Context Protocol) server
- Phase 5: Advanced features and optimizations

## Commands

### Available Commands (Phase 1)

| Command | Description | Example |
|---------|-------------|---------|
| `play` | Start playback | `maestro play` |
| `pause` | Pause playback | `maestro pause` |
| `stop` | Stop playback | `maestro stop` |
| `next` | Next track | `maestro next` |
| `previous` | Previous track | `maestro previous` |
| `volume` | Set volume (0-100) | `maestro volume 75` |
| `status` | Get player status | `maestro status --json` |

### Global Flags

- `--json`: Output in JSON format
- `--log-level`: Set log level (debug, info, warn, error)
- `--timeout`: Set operation timeout (default: 5s)

## Development

### Prerequisites

- Go 1.25+
- macOS with Apple Music.app
- Make

### Building

```bash
# Install dependencies
make deps

# Build all binaries
make build

# Run tests
make test

# Run linters
make lint

# Run all pre-commit checks
make pre-commit
```

### Testing

```bash
# Run all tests with coverage
make test

# Run only unit tests
make test-unit

# Run integration tests (requires Music.app)
make test-integration

# Generate coverage report
make test-coverage

# Check coverage meets threshold (80%)
make test-coverage-check
```

### Code Quality

The project maintains high code quality standards:

- **80% minimum test coverage** enforced by CI
- **golangci-lint** with comprehensive rules
- **Security scanning** with govulncheck
- **Automated formatting** with gofmt
- **Pre-commit hooks** for quality checks

### Project Structure

```
maestro/
├── cmd/                    # Application entry points
│   ├── maestro/           # Main CLI application
│   ├── maestro-exec/      # AppleScript executor subprocess
│   ├── maestrod/          # Daemon server (Phase 3)
│   ├── maestro-tui/       # Terminal UI (Phase 2)
│   └── maestro-mcp/       # MCP server (Phase 4)
├── domain/music/          # Domain layer
│   ├── entities.go        # Core entities (Track, Playlist, Player)
│   ├── values.go          # Value objects (Duration, Volume, etc.)
│   ├── repositories.go    # Repository interfaces
│   └── errors.go          # Domain errors
├── infrastructure/        # Infrastructure layer
│   └── applescript/       # AppleScript integration
├── presentation/          # Presentation layer
│   └── cli/              # Command-line interface
├── pkg/                   # Shared packages
│   └── logger/           # Structured logging
├── .github/              # GitHub workflows and templates
└── docs/                 # Documentation
```

## Configuration

### Environment Variables

```bash
# Logging configuration
export MAESTRO_LOG_LEVEL=info
export MAESTRO_LOG_FORMAT=json
export MAESTRO_LOG_OUTPUT=stdout

# Server configuration (Phase 3)
export MAESTRO_SERVER_PORT=8443
export MAESTRO_TLS_CERT_PATH=/path/to/cert.pem
export MAESTRO_TLS_KEY_PATH=/path/to/key.pem

# Performance tuning
export MAESTRO_SCRIPT_TIMEOUT=5s
export MAESTRO_MAX_CONCURRENT_OPERATIONS=10
```

### Configuration File (Coming Soon)

```yaml
# ~/.maestro/config.yaml
server:
  port: 8443
  tls:
    cert_path: ~/.maestro/certs/server.crt
    key_path: ~/.maestro/certs/server.key
    
logging:
  level: info
  format: json
  
performance:
  script_timeout: 5s
  max_concurrent_ops: 10
```

## API Reference (Phase 3)

### gRPC Service (Coming Soon)

```protobuf
service MaestroService {
  rpc Play(PlayRequest) returns (PlayResponse);
  rpc Pause(PauseRequest) returns (PauseResponse);
  rpc GetStatus(StatusRequest) returns (StatusResponse);
  // ... more methods
}
```

### MCP Integration (Phase 4)

Maestro will provide Model Context Protocol support for AI assistants:

```json
{
  "name": "maestro",
  "description": "Music control for AI assistants",
  "tools": [
    {
      "name": "play_music",
      "description": "Control music playback"
    }
  ]
}
```

## Security

### Authentication (Phase 3)

- **mTLS authentication** for API access
- **Certificate-based client authentication**
- **Session management** with configurable timeouts
- **Rate limiting** and request validation

### Permissions

Maestro requires the following macOS permissions:
- **Automation permission** for Apple Music.app control
- **Accessibility permission** (if using advanced features)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make pre-commit` to ensure quality
5. Submit a pull request

### Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

## Roadmap

See our [detailed roadmap](ROADMAP.md) and [GitHub milestones](https://github.com/madstone-tech/maestro/milestones) for planned features.

### Upcoming Features

- **Phase 2**: Interactive TUI with real-time updates
- **Phase 3**: gRPC API server with streaming
- **Phase 4**: MCP server for AI integration
- **Phase 5**: Performance optimizations and advanced features

## Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/madstone-tech/maestro/issues)
- **Discussions**: [GitHub Discussions](https://github.com/madstone-tech/maestro/discussions)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Logging powered by [Logrus](https://github.com/sirupsen/logrus)
- Future TUI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Inspired by the Unix philosophy of composable tools

---

**Made with ❤️ by the Madstone Technologies team**