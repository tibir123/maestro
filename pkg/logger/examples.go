package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// Example integration examples for the logging package

// ExampleCLIIntegration demonstrates how to integrate logging into CLI commands
func ExampleCLIIntegration() {
	// Initialize logging for CLI
	config := DefaultConfig()
	config.Component = "cli"
	config.Level = "info"
	config.Format = "text" // CLI typically uses text format for user readability

	// Set output based on environment
	if os.Getenv("MAESTRO_LOG_OUTPUT") != "" {
		config.Output = os.Getenv("MAESTRO_LOG_OUTPUT")
	}

	err := Initialize(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}

	// Create component-specific logger
	cliLogger := Component("cli")

	// Example usage in a CLI command
	rootCmd := &cobra.Command{
		Use:   "maestro",
		Short: "Music controller for macOS",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Log command execution
			cliLogger.Info("Command execution started",
				String("command", cmd.Use),
				String("args", fmt.Sprintf("%v", args)),
			)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// Log command completion
			cliLogger.Info("Command execution completed",
				String("command", cmd.Use),
			)
		},
	}

	// Add play command with logging
	playCmd := &cobra.Command{
		Use:   "play",
		Short: "Start or resume playback",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := cliLogger.WithOperation("play")

			start := time.Now()
			logger.Debug("Starting play operation")

			// Simulate command execution
			ctx := context.Background()
			err := executePlayCommand(ctx, logger)

			duration := time.Since(start)
			if err != nil {
				logger.Error("Play command failed",
					Error(err),
					Duration("duration", duration),
				)
				return err
			}

			logger.Info("Play command succeeded",
				Duration("duration", duration),
			)
			return nil
		},
	}

	rootCmd.AddCommand(playCmd)
}

// ExampleDaemonIntegration demonstrates how to integrate logging into the daemon
func ExampleDaemonIntegration() {
	// Initialize structured JSON logging for daemon
	config := ProductionConfig()
	config.Component = "maestrod"
	config.Output = "/var/log/maestro/maestrod.log"

	// Enable file rotation for production
	config.FileRotation = &FileRotationConfig{
		MaxSize:    100, // 100MB
		MaxAge:     30,  // 30 days
		MaxBackups: 10,  // 10 backup files
		Compress:   true,
	}

	err := Initialize(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize daemon logging: %v\n", err)
		os.Exit(1)
	}

	// Create component-specific loggers
	daemonLogger := Component("daemon")
	sessionLogger := Component("session")
	playerLogger := Component("player")

	// Log daemon startup
	daemonLogger.Info("Maestro daemon starting",
		String("version", "1.0.0"),
		String("config_file", "/etc/maestro/maestrod.toml"),
	)

	// Example of session management logging
	sessionID := "session-123"
	ctx := context.WithValue(context.Background(), "session_id", sessionID)

	sessionLoggerWithCtx := sessionLogger.WithContext(ctx)
	sessionLoggerWithCtx.Info("New session created",
		String("client_type", "cli"),
		String("client_id", "client-456"),
	)

	// Example of player operation logging
	playerLogger.Debug("Executing AppleScript command",
		String("script", "play.scpt"),
		String("session_id", sessionID),
	)
}

// ExampleInfrastructureIntegration demonstrates logging in infrastructure layers
func ExampleInfrastructureIntegration() {
	// Initialize logging for infrastructure components
	err := InitializeDevelopment()
	if err != nil {
		panic(err)
	}

	// AppleScript executor logging
	scriptLogger := Component("applescript")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-789")

	// Log script execution
	scriptLogger.WithContext(ctx).Info("Executing AppleScript",
		String("script_name", "get_current_track.scpt"),
		String("timeout", "5s"),
	)

	// gRPC server logging
	grpcLogger := Component("grpc")

	grpcLogger.Info("gRPC server starting",
		String("address", ":9090"),
		Bool("tls_enabled", true),
		String("cert_file", "/etc/maestro/certs/server.crt"),
	)

	// Cache operations logging
	cacheLogger := Component("cache")

	cacheLogger.Debug("Cache operation",
		String("operation", "get"),
		String("key", "track:123"),
		Bool("hit", true),
	)

	// Health check logging
	healthLogger := Component("health")

	healthLogger.Info("Health check completed",
		Bool("music_app_running", true),
		Bool("grpc_server_healthy", true),
		Duration("check_duration", 150*time.Millisecond),
	)
}

// ExampleContextAwareLogging demonstrates context-aware logging patterns
func ExampleContextAwareLogging() {
	err := InitializeDevelopment()
	if err != nil {
		panic(err)
	}

	// Create a context with request tracing information
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-abc123")
	ctx = context.WithValue(ctx, "session_id", "sess-def456")
	ctx = context.WithValue(ctx, "user_id", "user-ghi789")

	// Get a logger with context
	logger := WithContext(ctx).WithComponent("music_service")

	// All subsequent logs will include the context information
	logger.Info("Processing music command",
		String("command", "play"),
		String("track_id", "track-123"),
	)

	// Nested operation logging
	operationLogger := logger.WithOperation("track_lookup")
	operationLogger.Debug("Looking up track in library",
		String("query", "artist:Beatles title:Yesterday"),
	)

	operationLogger.Info("Track found",
		String("track_id", "track-456"),
		String("title", "Yesterday"),
		String("artist", "The Beatles"),
	)
}

// ExamplePerformanceLogging demonstrates performance monitoring
func ExamplePerformanceLogging() {
	err := InitializeDevelopment()
	if err != nil {
		panic(err)
	}

	logger := Component("performance")

	// Measure operation duration
	start := time.Now()

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	duration := time.Since(start)

	// Log performance metrics
	logger.Info("Operation performance",
		String("operation", "search_library"),
		Duration("duration", duration),
		Int("results_count", 42),
		Bool("cached", false),
	)

	// Log slow operations
	if duration > 500*time.Millisecond {
		logger.Warn("Slow operation detected",
			String("operation", "search_library"),
			Duration("duration", duration),
			String("threshold", "500ms"),
		)
	}
}

// ExampleErrorLogging demonstrates comprehensive error logging
func ExampleErrorLogging() {
	err := InitializeDevelopment()
	if err != nil {
		panic(err)
	}

	logger := Component("error_handling")

	// Simulate an error scenario
	ctx := context.WithValue(context.Background(), "session_id", "sess-error-123")
	contextLogger := logger.WithContext(ctx)

	// Log different types of errors with context
	contextLogger.Error("Music.app is not responding",
		String("operation", "get_current_track"),
		String("script", "get_current_track.scpt"),
		Int("retry_count", 3),
		Bool("will_restart_app", true),
	)

	contextLogger.Warn("Command timeout exceeded",
		String("command", "search"),
		Duration("timeout", 10*time.Second),
		String("query", "artist:long-search-term"),
	)

	// Fatal error logging (this would exit the program)
	// contextLogger.Fatal("Unable to connect to Music.app after restart",
	//     String("last_error", "osascript failed"),
	//     Int("restart_attempts", 3),
	// )
}

// Helper function to simulate command execution
func executePlayCommand(ctx context.Context, logger Logger) error {
	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Simulate random success/failure
	if time.Now().UnixNano()%2 == 0 {
		return fmt.Errorf("music app not responding")
	}

	return nil
}
