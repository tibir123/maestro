package cli

import (
	"context"
	"strconv"

	"github.com/madstone-tech/maestro/domain/music"
	"github.com/spf13/cobra"
)

// CommandContext holds shared dependencies for all CLI commands
type CommandContext struct {
	Context         context.Context
	PlayerRepo      music.PlayerRepository
	OutputFormatter *OutputFormatter
}

// NewPlayCommand creates the play command
func NewPlayCommand(ctx *CommandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "play",
		Short: "Start or resume playback",
		Long:  "Start or resume music playback. If music is paused, it will resume. If stopped, it will start playing.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.OutputFormatter.Debug("Executing play command")

			// Check current state to determine if we should play or resume
			player, err := ctx.PlayerRepo.GetCurrentState(ctx.Context)
			if err != nil {
				ctx.OutputFormatter.Error(err)
				return err
			}

			if player.IsPaused() {
				err = ctx.PlayerRepo.Resume(ctx.Context)
				if err != nil {
					ctx.OutputFormatter.Error(err)
					return err
				}
				ctx.OutputFormatter.Success("Playback resumed")
			} else {
				err = ctx.PlayerRepo.Resume(ctx.Context) // In AppleScript, both play and resume use "play"
				if err != nil {
					ctx.OutputFormatter.Error(err)
					return err
				}
				ctx.OutputFormatter.Success("Playback started")
			}

			return nil
		},
	}
}

// NewPauseCommand creates the pause command
func NewPauseCommand(ctx *CommandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "pause",
		Short: "Pause playback",
		Long:  "Pause the currently playing music. Playback can be resumed with the play command.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.OutputFormatter.Debug("Executing pause command")

			err := ctx.PlayerRepo.Pause(ctx.Context)
			if err != nil {
				ctx.OutputFormatter.Error(err)
				return err
			}

			ctx.OutputFormatter.Success("Playback paused")
			return nil
		},
	}
}

// NewStopCommand creates the stop command
func NewStopCommand(ctx *CommandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop playback",
		Long:  "Stop music playback completely. This will clear the current track and reset the position.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.OutputFormatter.Debug("Executing stop command")

			err := ctx.PlayerRepo.Stop(ctx.Context)
			if err != nil {
				ctx.OutputFormatter.Error(err)
				return err
			}

			ctx.OutputFormatter.Success("Playback stopped")
			return nil
		},
	}
}

// NewResumeCommand creates the resume command
func NewResumeCommand(ctx *CommandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "Resume paused playback",
		Long:  "Resume playback if it is currently paused. This is an alias for the play command.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.OutputFormatter.Debug("Executing resume command")

			err := ctx.PlayerRepo.Resume(ctx.Context)
			if err != nil {
				ctx.OutputFormatter.Error(err)
				return err
			}

			ctx.OutputFormatter.Success("Playback resumed")
			return nil
		},
	}
}

// NewNextCommand creates the next command
func NewNextCommand(ctx *CommandContext) *cobra.Command {
	return &cobra.Command{
		Use:   "next",
		Short: "Skip to next track",
		Long:  "Skip to the next track in the current playlist or queue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.OutputFormatter.Debug("Executing next command")

			err := ctx.PlayerRepo.Next(ctx.Context)
			if err != nil {
				ctx.OutputFormatter.Error(err)
				return err
			}

			ctx.OutputFormatter.Success("Skipped to next track")
			return nil
		},
	}
}

// NewPreviousCommand creates the previous command
func NewPreviousCommand(ctx *CommandContext) *cobra.Command {
	return &cobra.Command{
		Use:     "previous",
		Aliases: []string{"prev"},
		Short:   "Skip to previous track",
		Long:    "Skip to the previous track in the current playlist or queue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.OutputFormatter.Debug("Executing previous command")

			err := ctx.PlayerRepo.Previous(ctx.Context)
			if err != nil {
				ctx.OutputFormatter.Error(err)
				return err
			}

			ctx.OutputFormatter.Success("Skipped to previous track")
			return nil
		},
	}
}

// NewVolumeCommand creates the volume command
func NewVolumeCommand(ctx *CommandContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume [level]",
		Short: "Get or set volume",
		Long: `Get the current volume level or set it to a specific value.
Volume should be a number between 0 and 100.

Examples:
  maestro volume        # Show current volume
  maestro volume 50     # Set volume to 50%
  maestro volume 0      # Mute
  maestro volume 100    # Maximum volume`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.OutputFormatter.Debug("Executing volume command")

			// If no arguments, show current volume
			if len(args) == 0 {
				player, err := ctx.PlayerRepo.GetCurrentState(ctx.Context)
				if err != nil {
					ctx.OutputFormatter.Error(err)
					return err
				}
				ctx.OutputFormatter.PrintVolume(player.Volume)
				return nil
			}

			// Parse volume level
			levelStr := args[0]
			level, err := strconv.Atoi(levelStr)
			if err != nil {
				ctx.OutputFormatter.Error(music.NewDomainError(music.ErrInvalidVolume, "volume must be a number between 0 and 100"))
				return err
			}

			// Validate and create volume
			volume := music.NewVolume(level)
			if !volume.IsValid() {
				ctx.OutputFormatter.Error(music.NewDomainError(music.ErrInvalidVolume, "volume must be between 0 and 100"))
				return music.NewDomainError(music.ErrInvalidVolume, "volume must be between 0 and 100")
			}

			// Set volume
			err = ctx.PlayerRepo.SetVolume(ctx.Context, volume)
			if err != nil {
				ctx.OutputFormatter.Error(err)
				return err
			}

			ctx.OutputFormatter.Success("Volume set to " + volume.String())
			return nil
		},
	}

	return cmd
}

// NewStatusCommand creates the status command
func NewStatusCommand(ctx *CommandContext) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Aliases: []string{"stat"},
		Short:   "Show current player status",
		Long:    "Display the current player status including state, volume, current track, and playback settings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.OutputFormatter.Debug("Executing status command")

			// Get current player state
			player, err := ctx.PlayerRepo.GetCurrentState(ctx.Context)
			if err != nil {
				ctx.OutputFormatter.Error(err)
				return err
			}

			// Get current track if there is one
			var track *music.Track
			if player.HasCurrentTrack() {
				track, err = ctx.PlayerRepo.GetCurrentTrack(ctx.Context)
				if err != nil {
					// Don't fail the command if we can't get track info, just log it
					ctx.OutputFormatter.Debug("Could not get current track: " + err.Error())
				}
			}

			ctx.OutputFormatter.PrintPlayerStatus(player, track)
			return nil
		},
	}
}
