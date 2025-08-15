package logger

import (
	"context"
	"fmt"
	"time"
)

// Demo demonstrates the logging package functionality
// This file can be used to manually test the logging implementation
func Demo() {
	fmt.Println("=== Maestro Logging Package Demo ===")
	fmt.Println()

	// 1. Basic Configuration Demo
	fmt.Println("1. Testing different configurations:")

	fmt.Println("\n--- Development Configuration (Text Format) ---")
	err := InitializeDevelopment()
	if err != nil {
		panic(err)
	}

	Info("Development mode initialized")
	Debug("This is a debug message in development mode")
	WithFields(String("feature", "logging"), Bool("enabled", true)).Info("Feature status")

	fmt.Println("\n--- Production Configuration (JSON Format) ---")
	err = InitializeProduction()
	if err != nil {
		panic(err)
	}

	Info("Production mode initialized")
	Debug("This debug message should not appear in production")
	WithFields(String("service", "maestro"), String("version", "1.0.0")).Info("Service information")

	// 2. Component-specific logging
	fmt.Println("\n2. Component-specific logging:")

	cliLogger := Component("cli")
	cliLogger.Info("CLI component initialized")

	playerLogger := Component("player")
	playerLogger.Info("Player component initialized")

	sessionLogger := Component("session")
	sessionLogger.Info("Session component initialized")

	// 3. Context-aware logging
	fmt.Println("\n3. Context-aware logging:")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-demo-123")
	ctx = context.WithValue(ctx, "session_id", "sess-demo-456")
	ctx = context.WithValue(ctx, "user_id", "user-demo-789")

	contextLogger := WithContext(ctx).WithComponent("demo")
	contextLogger.Info("Processing user request", String("action", "play"))
	contextLogger.Info("Request completed successfully", Duration("processing_time", 150*time.Millisecond))

	// 4. Operation tracking
	fmt.Println("\n4. Operation tracking:")

	operationLogger := contextLogger.WithOperation("search_music")
	operationLogger.Debug("Starting music search", String("query", "artist:Beatles"))
	operationLogger.Info("Search completed", Int("results", 42), Duration("search_time", 230*time.Millisecond))

	// 5. Error logging
	fmt.Println("\n5. Error logging:")

	errorLogger := Component("applescript")
	errorLogger.Error("Failed to execute script",
		Error(fmt.Errorf("osascript: execution error")),
		String("script", "get_current_track.scpt"),
		Int("retry_count", 3),
	)

	errorLogger.Warn("Slow operation detected",
		String("operation", "library_search"),
		Duration("duration", 1200*time.Millisecond),
		Duration("threshold", 500*time.Millisecond),
	)

	// 6. Performance monitoring
	fmt.Println("\n6. Performance monitoring:")

	perfLogger := Component("performance")
	start := time.Now()

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	perfLogger.Info("Operation completed",
		String("operation", "get_player_state"),
		Duration("duration", time.Since(start)),
		Bool("success", true),
		String("cache_status", "hit"),
	)

	// 7. Structured logging benefits
	fmt.Println("\n7. Structured logging benefits:")

	structuredLogger := Component("metrics")
	structuredLogger.Info("System metrics",
		Int("active_sessions", 3),
		Int("memory_usage_mb", 85),
		Float64("cpu_usage_percent", 2.5),
		Bool("music_app_healthy", true),
		String("last_command", "play"),
		Duration("uptime", 2*time.Hour+15*time.Minute),
	)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("The logging package provides:")
	fmt.Println("✓ Structured JSON logging with logrus")
	fmt.Println("✓ Component-specific loggers")
	fmt.Println("✓ Context-aware logging for request tracing")
	fmt.Println("✓ Configurable log levels and formats")
	fmt.Println("✓ Performance monitoring capabilities")
	fmt.Println("✓ Production-ready error handling")
}

// DemoEnvironmentConfiguration shows how to configure logging via environment variables
func DemoEnvironmentConfiguration() {
	fmt.Println("=== Environment Configuration Demo ===")
	fmt.Println()

	fmt.Println("Set these environment variables to configure logging:")
	fmt.Println("export MAESTRO_LOG_LEVEL=debug")
	fmt.Println("export MAESTRO_LOG_FORMAT=json")
	fmt.Println("export MAESTRO_LOG_OUTPUT=stdout")
	fmt.Println("export MAESTRO_LOG_COMPONENT=maestro")
	fmt.Println("export MAESTRO_LOG_CALLER=true")

	config := LoadFromEnv()
	fmt.Printf("\nCurrent configuration from environment:\n")
	fmt.Printf("  Level: %s\n", config.Level)
	fmt.Printf("  Format: %s\n", config.Format)
	fmt.Printf("  Output: %s\n", config.Output)
	fmt.Printf("  Component: %s\n", config.Component)
	fmt.Printf("  EnableCaller: %t\n", config.EnableCaller)
}

// DemoIntegrationPatterns shows common integration patterns
func DemoIntegrationPatterns() {
	fmt.Println("=== Integration Patterns Demo ===")
	fmt.Println()

	// Initialize logging
	err := InitializeDevelopment()
	if err != nil {
		panic(err)
	}

	// Pattern 1: CLI Command Logging
	fmt.Println("1. CLI Command Pattern:")
	cliLogger := Component("cli")
	start := time.Now()

	cliLogger.Info("Command execution started",
		String("command", "play"),
		String("args", "[]"),
	)

	// Simulate command execution
	time.Sleep(25 * time.Millisecond)

	cliLogger.Info("Command execution completed",
		String("command", "play"),
		Duration("execution_time", time.Since(start)),
		Bool("success", true),
	)

	// Pattern 2: Infrastructure Service Logging
	fmt.Println("\n2. Infrastructure Service Pattern:")
	serviceLogger := Component("applescript_executor")

	serviceLogger.Debug("Executing AppleScript",
		String("script", "get_current_track.scpt"),
		Duration("timeout", 5*time.Second),
	)

	serviceLogger.Info("AppleScript execution completed",
		String("script", "get_current_track.scpt"),
		Duration("execution_time", 120*time.Millisecond),
		Int("output_length", 256),
	)

	// Pattern 3: Session Management Logging
	fmt.Println("\n3. Session Management Pattern:")
	sessionLogger := Component("session_manager")

	sessionID := "sess_" + fmt.Sprintf("%d", time.Now().UnixNano())
	ctx := context.WithValue(context.Background(), "session_id", sessionID)

	sessionLoggerWithCtx := sessionLogger.WithContext(ctx)
	sessionLoggerWithCtx.Info("Session created",
		String("client_type", "cli"),
		String("client_id", "client-123"),
		Duration("timeout", 5*time.Minute),
	)

	sessionLoggerWithCtx.Debug("Session heartbeat received")

	sessionLoggerWithCtx.Info("Session terminated",
		String("reason", "timeout"),
		Duration("session_duration", 4*time.Minute+30*time.Second),
	)

	fmt.Println("\n=== Integration Patterns Complete ===")
}
