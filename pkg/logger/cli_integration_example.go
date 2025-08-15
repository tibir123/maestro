package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/madstone-tech/maestro/domain/music"
)

// This file demonstrates how to integrate the logging package into the existing CLI

// UpdatedCommandContext shows how to integrate logging into the CLI CommandContext
type UpdatedCommandContext struct {
	Context         context.Context
	PlayerRepo      music.PlayerRepository
	OutputFormatter *UpdatedOutputFormatter
	Logger          Logger // Add logger to context
}

// UpdatedOutputFormatter demonstrates how to enhance the existing OutputFormatter with logging
type UpdatedOutputFormatter struct {
	jsonMode bool
	verbose  bool
	logger   Logger // Add logger for internal operations
}

// NewUpdatedOutputFormatter creates an enhanced output formatter with logging
func NewUpdatedOutputFormatter(jsonMode, verbose bool, logger Logger) *UpdatedOutputFormatter {
	return &UpdatedOutputFormatter{
		jsonMode: jsonMode,
		verbose:  verbose,
		logger:   logger.WithComponent("output_formatter"),
	}
}

// Example showing how to update the existing play command with logging
func ExampleUpdatedPlayCommand(ctx *UpdatedCommandContext) error {
	// Create operation-specific logger
	logger := ctx.Logger.WithOperation("play_command")

	start := time.Now()
	logger.Debug("Executing play command")

	// Check current state to determine if we should play or resume
	player, err := ctx.PlayerRepo.GetCurrentState(ctx.Context)
	if err != nil {
		logger.Error("Failed to get current player state", Error(err))
		ctx.OutputFormatter.Error(err)
		return err
	}

	logger.Debug("Retrieved player state",
		String("current_state", player.State.String()),
		Bool("has_track", player.HasCurrentTrack()),
	)

	var action string
	if player.IsPaused() {
		action = "resume"
		err = ctx.PlayerRepo.Resume(ctx.Context)
		if err != nil {
			logger.Error("Failed to resume playback",
				Error(err),
				Duration("operation_time", time.Since(start)),
			)
			ctx.OutputFormatter.Error(err)
			return err
		}
		ctx.OutputFormatter.Success("Playback resumed")
		logger.Info("Playback resumed successfully")
	} else {
		action = "play"
		err = ctx.PlayerRepo.Resume(ctx.Context) // In AppleScript, both play and resume use "play"
		if err != nil {
			logger.Error("Failed to start playback",
				Error(err),
				Duration("operation_time", time.Since(start)),
			)
			ctx.OutputFormatter.Error(err)
			return err
		}
		ctx.OutputFormatter.Success("Playback started")
		logger.Info("Playback started successfully")
	}

	duration := time.Since(start)
	logger.Info("Play command completed",
		String("action", action),
		Duration("total_time", duration),
		Bool("success", true),
	)

	return nil
}

// Example showing how to enhance the status command with comprehensive logging
func ExampleUpdatedStatusCommand(ctx *UpdatedCommandContext) error {
	logger := ctx.Logger.WithOperation("status_command")

	start := time.Now()
	logger.Debug("Executing status command")

	// Get current player state
	player, err := ctx.PlayerRepo.GetCurrentState(ctx.Context)
	if err != nil {
		logger.Error("Failed to get player state",
			Error(err),
			Duration("operation_time", time.Since(start)),
		)
		ctx.OutputFormatter.Error(err)
		return err
	}

	logger.Debug("Retrieved player state",
		String("state", player.State.String()),
		Int("volume", player.Volume.Level()),
		Bool("shuffle", player.Shuffle),
		String("repeat", player.Repeat.String()),
	)

	// Get current track if there is one
	var track *music.Track
	if player.HasCurrentTrack() {
		trackStart := time.Now()
		track, err = ctx.PlayerRepo.GetCurrentTrack(ctx.Context)
		if err != nil {
			// Don't fail the command if we can't get track info
			logger.Warn("Could not get current track",
				Error(err),
				Duration("track_lookup_time", time.Since(trackStart)),
			)
		} else {
			logger.Debug("Retrieved current track",
				String("track_id", track.ID.Value()),
				String("title", track.Title),
				String("artist", track.Artist),
				Duration("track_lookup_time", time.Since(trackStart)),
			)
		}
	} else {
		logger.Debug("No current track playing")
	}

	ctx.OutputFormatter.PrintPlayerStatus(player, track)

	duration := time.Since(start)
	logger.Info("Status command completed",
		Duration("total_time", duration),
		Bool("has_track", track != nil),
		Bool("success", true),
	)

	return nil
}

// Enhanced OutputFormatter methods with logging
func (f *UpdatedOutputFormatter) Success(message string) {
	f.logger.Debug("Displaying success message", String("message", message))

	if f.jsonMode {
		f.printJSON(map[string]interface{}{
			"success": true,
			"message": message,
		})
	} else {
		// Original text output
		// fmt.Fprintf(f.writer, "%s\n", message)
	}
}

func (f *UpdatedOutputFormatter) Error(err error) {
	f.logger.Debug("Displaying error message", Error(err))

	if f.jsonMode {
		f.printJSON(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	} else {
		// Original text output
		// fmt.Fprintf(f.writer, "Error: %s\n", err.Error())
	}
}

func (f *UpdatedOutputFormatter) PrintPlayerStatus(player *music.Player, track *music.Track) {
	f.logger.Debug("Displaying player status",
		String("state", player.State.String()),
		Int("volume", player.Volume.Level()),
		Bool("has_track", track != nil),
	)

	if f.jsonMode {
		f.printPlayerStatusJSON(player, track)
	} else {
		f.printPlayerStatusText(player, track)
	}
}

// Placeholder methods (implementation would be same as original)
func (f *UpdatedOutputFormatter) printJSON(data interface{}) {
	// Same as original implementation
	f.logger.Debug("Outputting JSON response")
}

func (f *UpdatedOutputFormatter) printPlayerStatusJSON(player *music.Player, track *music.Track) {
	// Same as original implementation
	f.logger.Debug("Outputting player status as JSON")
}

func (f *UpdatedOutputFormatter) printPlayerStatusText(player *music.Player, track *music.Track) {
	// Same as original implementation
	f.logger.Debug("Outputting player status as text")
}

// Example showing how to initialize logging in the main CLI entry point
func ExampleCLIMain() {
	// Initialize logging based on environment or config
	var config *Config

	// Check if we're in development mode
	if os.Getenv("MAESTRO_ENV") == "development" {
		config = DevelopmentConfig()
	} else {
		config = LoadFromEnv()
		if config.Component == "" {
			config.Component = "cli"
		}
	}

	// Initialize the global logger
	err := Initialize(config)
	if err != nil {
		// Fallback to stderr if logging setup fails
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize logging: %v\n", err)
		// Continue with basic logging
		_ = InitializeDefault()
	}

	// Create CLI logger
	cliLogger := Component("cli")
	cliLogger.Info("Maestro CLI starting",
		String("version", "1.0.0"),
		String("log_level", config.Level),
		String("log_format", config.Format),
	)

	// Create enhanced command context
	_ = &UpdatedCommandContext{
		Context: context.Background(),
		Logger:  cliLogger,
		// PlayerRepo and OutputFormatter would be initialized here
	}

	// Example of handling CLI commands with logging
	cliLogger.Debug("CLI initialization completed")
}

// Example showing how to add request ID tracking for better traceability
func ExampleRequestTracking() {
	logger := Component("cli")

	// Generate a unique request ID for this CLI invocation
	requestID := generateRequestID()
	ctx := context.WithValue(context.Background(), "request_id", requestID)

	// Create request-aware logger
	requestLogger := logger.WithContext(ctx)

	requestLogger.Info("CLI request started",
		String("command", "play"),
		String("user", getCurrentUser()),
	)

	// All subsequent operations in this request will include the request_id
	// This makes it easy to trace a single command execution through logs
}

// Helper functions (would be implemented elsewhere)
func generateRequestID() string {
	// Implementation would generate a unique ID
	return "req-" + time.Now().Format("20060102-150405") + "-abc123"
}

func getCurrentUser() string {
	// Implementation would get current user
	return os.Getenv("USER")
}
