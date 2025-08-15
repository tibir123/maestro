// Package applescript provides the AppleScript infrastructure layer for Maestro.
//
// This package implements the domain repository interfaces using AppleScript
// to control Music.app on macOS. It serves as the bridge between the pure
// domain layer and the actual Music.app functionality.
//
// # Architecture
//
// The package is organized around several key components:
//
//   - Executor: Handles AppleScript execution with timeout and retry logic
//   - PlayerRepository: Implements music.PlayerRepository for playback control
//   - Script Templates: Reusable AppleScript files for common operations
//
// # Usage
//
// Basic usage involves creating an executor and using it to create repository
// implementations:
//
//	executor := applescript.NewExecutor(nil) // Use default config
//	playerRepo := applescript.NewPlayerRepository(executor)
//
//	// Use through domain interfaces
//	err := playerRepo.Play(ctx, trackID)
//	state, err := playerRepo.GetCurrentState(ctx)
//
// # Error Handling
//
// All methods return domain errors from the music package. These errors
// provide structured information about failures and support retry logic:
//
//	err := playerRepo.SetVolume(ctx, volume)
//	if music.IsRetryable(err) {
//		// Retry the operation
//	}
//
// # Script Templates
//
// The package includes pre-built AppleScript templates in the scripts/
// directory. These can be executed directly or used as templates with
// parameter substitution:
//
//	result := executor.ExecuteTemplate(ctx, "play", map[string]interface{}{
//		"track_id": "12345",
//	})
//
// # Configuration
//
// The executor can be configured with custom timeouts, retry counts,
// and the path to the maestro-exec binary:
//
//	config := &applescript.ExecutorConfig{
//		ExecPath:       "/custom/path/to/maestro-exec",
//		DefaultTimeout: 30 * time.Second,
//		MaxRetries:     5,
//		RetryDelay:     1 * time.Second,
//	}
//	executor := applescript.NewExecutor(config)
//
// # Requirements
//
// This package requires:
//   - macOS with Music.app installed
//   - The maestro-exec binary in PATH or specified path
//   - Proper permissions for AppleScript to control Music.app
//
// # Thread Safety
//
// All types in this package are thread-safe and can be used concurrently.
// However, Music.app itself has limitations on concurrent operations, so
// callers should consider serializing operations when necessary.
package applescript
