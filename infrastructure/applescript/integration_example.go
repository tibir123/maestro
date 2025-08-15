package applescript

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/madstone-tech/maestro/domain/music"
)

// IntegrationExample demonstrates how to use the AppleScript infrastructure.
// This file serves as both documentation and a basic integration test.
func IntegrationExample() {
	fmt.Println("=== Maestro AppleScript Infrastructure Integration Example ===")

	// Create executor with default configuration
	executor := NewExecutor(nil)

	// Check if maestro-exec is available
	fmt.Println("1. Checking if maestro-exec is available...")
	if err := executor.IsExecutable(); err != nil {
		log.Printf("ERROR: maestro-exec not available: %v", err)
		return
	}
	fmt.Println("   ✓ maestro-exec is available")

	// Create player repository
	playerRepo := NewPlayerRepository(executor)

	// Create context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check Music.app health
	fmt.Println("\n2. Checking Music.app health...")
	if err := playerRepo.HealthCheck(ctx); err != nil {
		log.Printf("ERROR: Music.app health check failed: %v", err)
		return
	}
	fmt.Println("   ✓ Music.app is accessible")

	// Get current player state
	fmt.Println("\n3. Getting current player state...")
	state, err := playerRepo.GetCurrentState(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get player state: %v", err)
		return
	}
	fmt.Printf("   Player State: %s\n", state.State.String())
	fmt.Printf("   Volume: %s\n", state.Volume.String())
	fmt.Printf("   Position: %s\n", state.Position.String())
	fmt.Printf("   Shuffle: %t\n", state.Shuffle)
	fmt.Printf("   Repeat: %s\n", state.Repeat.String())

	// Get current track if playing
	fmt.Println("\n4. Getting current track...")
	track, err := playerRepo.GetCurrentTrack(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get current track: %v", err)
		return
	}

	if track != nil {
		fmt.Printf("   Current Track: %s\n", track.String())
		fmt.Printf("   Album: %s\n", track.Album)
		fmt.Printf("   Duration: %s\n", track.Duration.String())
	} else {
		fmt.Println("   No track currently playing")
	}

	// Test volume control
	fmt.Println("\n5. Testing volume control...")
	originalVolume := state.Volume
	testVolume := music.NewVolume(75)

	if err := playerRepo.SetVolume(ctx, testVolume); err != nil {
		log.Printf("ERROR: Failed to set volume: %v", err)
	} else {
		fmt.Printf("   ✓ Set volume to %s\n", testVolume.String())

		// Restore original volume
		time.Sleep(1 * time.Second)
		if err := playerRepo.SetVolume(ctx, originalVolume); err != nil {
			log.Printf("WARNING: Failed to restore original volume: %v", err)
		} else {
			fmt.Printf("   ✓ Restored volume to %s\n", originalVolume.String())
		}
	}

	// Test playback controls (only if we have a track)
	if state.State != music.PlayerStateStopped {
		fmt.Println("\n6. Testing playback controls...")

		// Test pause
		if state.State == music.PlayerStatePlaying {
			if err := playerRepo.Pause(ctx); err != nil {
				log.Printf("ERROR: Failed to pause: %v", err)
			} else {
				fmt.Println("   ✓ Paused playback")
				time.Sleep(1 * time.Second)

				// Test resume
				if err := playerRepo.Resume(ctx); err != nil {
					log.Printf("ERROR: Failed to resume: %v", err)
				} else {
					fmt.Println("   ✓ Resumed playback")
				}
			}
		} else {
			// Test play
			if err := playerRepo.Resume(ctx); err != nil {
				log.Printf("ERROR: Failed to start playback: %v", err)
			} else {
				fmt.Println("   ✓ Started playback")
			}
		}
	} else {
		fmt.Println("\n6. Skipping playback controls (no track loaded)")
	}

	// Test shuffle and repeat
	fmt.Println("\n7. Testing shuffle and repeat modes...")
	originalShuffle := state.Shuffle
	originalRepeat := state.Repeat

	// Toggle shuffle
	if err := playerRepo.SetShuffle(ctx, !originalShuffle); err != nil {
		log.Printf("ERROR: Failed to toggle shuffle: %v", err)
	} else {
		fmt.Printf("   ✓ Set shuffle to %t\n", !originalShuffle)

		// Restore original
		time.Sleep(500 * time.Millisecond)
		if err := playerRepo.SetShuffle(ctx, originalShuffle); err != nil {
			log.Printf("WARNING: Failed to restore shuffle setting: %v", err)
		} else {
			fmt.Printf("   ✓ Restored shuffle to %t\n", originalShuffle)
		}
	}

	// Test repeat mode
	testRepeat := music.RepeatModeOne
	if originalRepeat == music.RepeatModeOne {
		testRepeat = music.RepeatModeAll
	}

	if err := playerRepo.SetRepeat(ctx, testRepeat); err != nil {
		log.Printf("ERROR: Failed to set repeat mode: %v", err)
	} else {
		fmt.Printf("   ✓ Set repeat to %s\n", testRepeat.String())

		// Restore original
		time.Sleep(500 * time.Millisecond)
		if err := playerRepo.SetRepeat(ctx, originalRepeat); err != nil {
			log.Printf("WARNING: Failed to restore repeat mode: %v", err)
		} else {
			fmt.Printf("   ✓ Restored repeat to %s\n", originalRepeat.String())
		}
	}

	// Test template execution
	fmt.Println("\n8. Testing template execution...")
	result := executor.ExecuteTemplate(ctx, "health_check", nil)
	if result.Error != nil {
		log.Printf("ERROR: Template execution failed: %v", result.Error)
	} else {
		fmt.Printf("   ✓ Template executed successfully (output: %s)\n", result.Output)
		fmt.Printf("   Duration: %v, Retries: %d\n", result.Duration, result.RetryCount)
	}

	fmt.Println("\n=== Integration Example Complete ===")
	fmt.Println("All tests completed. Check the output above for any errors.")
	fmt.Println("The AppleScript infrastructure is ready for use!")
}

// DemoUsage shows typical usage patterns for the AppleScript infrastructure.
func DemoUsage() {
	// This function demonstrates the typical patterns developers will use
	// when integrating with the AppleScript infrastructure layer.

	// 1. Create an executor (usually done once at application startup)
	config := DefaultExecutorConfig()
	config.ExecPath = "/path/to/maestro-exec" // Optional: specify custom path
	config.DefaultTimeout = 15 * time.Second  // Optional: customize timeout
	executor := NewExecutor(config)

	// 2. Create repository implementations
	playerRepo := NewPlayerRepository(executor)

	// 3. Use the repositories through the domain interfaces
	ctx := context.Background()

	// Basic playback control
	_ = playerRepo.Pause(ctx)
	_ = playerRepo.Resume(ctx)
	_ = playerRepo.Next(ctx)
	_ = playerRepo.Previous(ctx)

	// Volume control
	volume := music.NewVolume(80)
	_ = playerRepo.SetVolume(ctx, volume)

	// Player state
	state, _ := playerRepo.GetCurrentState(ctx)
	if state != nil {
		fmt.Printf("Player is %s at volume %s\n", state.State.String(), state.Volume.String())
	}

	// Track information
	track, _ := playerRepo.GetCurrentTrack(ctx)
	if track != nil {
		fmt.Printf("Now playing: %s\n", track.String())
	}
}
