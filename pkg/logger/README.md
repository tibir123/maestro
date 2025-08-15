# Maestro Logging Package

A structured logging package for Maestro that provides JSON logging with logrus, designed for observability and production use.

## Features

- **Structured JSON Logging**: Machine-readable logs using logrus
- **Component-Specific Loggers**: Each system component can have its own logger
- **Context-Aware Logging**: Request tracing with context propagation
- **Configurable Output**: Support for stdout, stderr, or file output
- **Multiple Log Levels**: Debug, Info, Warn, Error, Fatal, Panic
- **Environment Configuration**: Configure via environment variables
- **Performance Optimized**: Minimal overhead for production use
- **Thread Safe**: Safe for concurrent use across goroutines

## Quick Start

### Basic Usage

```go
package main

import (
    "github.com/mdstn/maestro/pkg/logger"
)

func main() {
    // Initialize with default configuration
    err := logger.InitializeDefault()
    if err != nil {
        panic(err)
    }
    
    // Basic logging
    logger.Info("Application started")
    logger.Debug("Debug information", 
        logger.String("version", "1.0.0"),
        logger.Bool("debug_mode", true),
    )
    
    // Component-specific logging
    cliLogger := logger.Component("cli")
    cliLogger.Info("CLI initialized")
}
```

### Configuration

#### Programmatic Configuration

```go
// Development configuration (text format, debug level)
err := logger.InitializeDevelopment()

// Production configuration (JSON format, info level)
err := logger.InitializeProduction()

// Custom configuration
config := &logger.Config{
    Level:           "debug",
    Format:          "json",
    Output:          "/var/log/maestro/app.log",
    Component:       "maestro",
    EnableCaller:    true,
    EnableTimestamp: true,
}
err := logger.Initialize(config)
```

#### Environment Variables

```bash
export MAESTRO_LOG_LEVEL=debug
export MAESTRO_LOG_FORMAT=json
export MAESTRO_LOG_OUTPUT=stdout
export MAESTRO_LOG_COMPONENT=maestro
export MAESTRO_LOG_CALLER=true
```

```go
config := logger.LoadFromEnv()
err := logger.Initialize(config)
```

## Usage Patterns

### Component-Specific Logging

```go
// Create component-specific loggers
cliLogger := logger.Component("cli")
playerLogger := logger.Component("player")
sessionLogger := logger.Component("session")

cliLogger.Info("Processing command", logger.String("cmd", "play"))
playerLogger.Debug("Executing AppleScript", logger.String("script", "play.scpt"))
```

### Context-Aware Logging

```go
// Add context information for request tracing
ctx := context.Background()
ctx = context.WithValue(ctx, "request_id", "req-123")
ctx = context.WithValue(ctx, "session_id", "sess-456")

contextLogger := logger.WithContext(ctx).WithComponent("music_service")
contextLogger.Info("Processing request", logger.String("action", "play"))

// All subsequent logs will include request_id and session_id
```

### Operation Tracking

```go
logger := logger.Component("applescript").WithOperation("get_current_track")

start := time.Now()
logger.Debug("Starting operation")

// ... do work ...

logger.Info("Operation completed",
    logger.Duration("duration", time.Since(start)),
    logger.Bool("success", true),
)
```

### Error Logging

```go
logger := logger.Component("player")

err := executeCommand()
if err != nil {
    logger.Error("Command execution failed",
        logger.Error(err),
        logger.String("command", "play"),
        logger.Int("retry_count", 3),
    )
    return err
}
```

### Performance Monitoring

```go
perfLogger := logger.Component("performance")

start := time.Now()
result, err := searchLibrary(query)
duration := time.Since(start)

perfLogger.Info("Library search completed",
    logger.String("query", query),
    logger.Int("results", len(result)),
    logger.Duration("duration", duration),
    logger.Bool("success", err == nil),
)

if duration > 500*time.Millisecond {
    perfLogger.Warn("Slow operation detected",
        logger.Duration("duration", duration),
        logger.Duration("threshold", 500*time.Millisecond),
    )
}
```

## Field Types

The package provides convenient field helper functions:

```go
logger.String("key", "value")           // String field
logger.Int("count", 42)                 // Integer field
logger.Int64("timestamp", 1609459200)   // Int64 field
logger.Float64("ratio", 3.14)           // Float64 field
logger.Bool("enabled", true)            // Boolean field
logger.Duration("elapsed", time.Second) // Duration field
logger.Error(err)                       // Error field (key="error")
logger.Any("data", complexStruct)       // Any type field
```

## Integration Examples

### CLI Integration

```go
func NewPlayCommand(ctx *CommandContext) *cobra.Command {
    return &cobra.Command{
        Use: "play",
        RunE: func(cmd *cobra.Command, args []string) error {
            logger := ctx.Logger.WithOperation("play_command")
            
            start := time.Now()
            logger.Debug("Executing play command")
            
            err := ctx.PlayerRepo.Play(ctx.Context)
            if err != nil {
                logger.Error("Play command failed",
                    logger.Error(err),
                    logger.Duration("duration", time.Since(start)),
                )
                return err
            }
            
            logger.Info("Play command succeeded",
                logger.Duration("duration", time.Since(start)),
            )
            return nil
        },
    }
}
```

### Infrastructure Integration

```go
type AppleScriptExecutor struct {
    logger logger.Logger
}

func (e *AppleScriptExecutor) ExecuteScript(ctx context.Context, script string) error {
    logger := e.logger.WithContext(ctx).WithOperation("execute_script")
    
    start := time.Now()
    logger.Debug("Starting script execution",
        logger.String("script", script),
    )
    
    // ... execute script ...
    
    logger.Info("Script execution completed",
        logger.String("script", script),
        logger.Duration("execution_time", time.Since(start)),
    )
    
    return nil
}
```

## Configuration Options

| Option | Default | Description |
|--------|---------|-------------|
| `Level` | `"info"` | Minimum log level (debug, info, warn, error) |
| `Format` | `"json"` | Log format (json, text) |
| `Output` | `"stdout"` | Output destination (stdout, stderr, file path) |
| `Component` | `"maestro"` | Default component name |
| `EnableCaller` | `false` | Include caller information in logs |
| `EnableTimestamp` | `true` | Include timestamp in logs |
| `TimestampFormat` | `RFC3339` | Timestamp format |

## Output Examples

### JSON Format (Production)

```json
{
  "component": "cli",
  "level": "info",
  "msg": "Command execution started",
  "time": "2025-01-15T10:30:45Z",
  "command": "play",
  "request_id": "req-123"
}
```

### Text Format (Development)

```
time="2025-01-15T10:30:45Z" level=info msg="Command execution started" component=cli command=play request_id=req-123
```

## Best Practices

1. **Use Component Loggers**: Create component-specific loggers for better organization
2. **Include Context**: Use context-aware logging for request tracing
3. **Structured Fields**: Use structured fields instead of string interpolation
4. **Performance Logging**: Log operation durations for performance monitoring
5. **Error Context**: Include relevant context when logging errors
6. **Log Levels**: Use appropriate log levels (debug for development, info for normal operations)
7. **Avoid Secrets**: Never log sensitive information like passwords or tokens

## Environment Configuration

For production deployments, configure logging via environment variables:

```bash
# Production settings
export MAESTRO_LOG_LEVEL=info
export MAESTRO_LOG_FORMAT=json
export MAESTRO_LOG_OUTPUT=/var/log/maestro/maestro.log

# Development settings
export MAESTRO_LOG_LEVEL=debug
export MAESTRO_LOG_FORMAT=text
export MAESTRO_LOG_OUTPUT=stdout
export MAESTRO_LOG_CALLER=true
```

## Testing

Run the test suite:

```bash
go test ./pkg/logger -v
```

## Dependencies

- [logrus](https://github.com/sirupsen/logrus) - Structured logging library

## Thread Safety

The logging package is thread-safe and can be used safely across multiple goroutines. Each logger instance maintains its own state and logrus handles concurrent writes safely.