package main

import (
	"context"
	"fmt"
	"os"

	"github.com/madstone-tech/maestro/infrastructure/applescript"
	"github.com/madstone-tech/maestro/presentation/cli"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	jsonOutput bool
	verbose    bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "maestro",
		Short: "Maestro - Control your music from the command line",
		Long: `Maestro is a command-line interface for controlling music playback.
It provides simple commands to play, pause, skip tracks, and manage volume
using your system's music player.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add global flags
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Initialize infrastructure
	executor := applescript.NewExecutor(nil)
	playerRepo := applescript.NewPlayerRepository(executor)

	// Create command context (OutputFormatter will be set in PersistentPreRun)
	ctx := context.Background()
	cmdCtx := &cli.CommandContext{
		Context:    ctx,
		PlayerRepo: playerRepo,
	}

	// Set up PersistentPreRun to initialize OutputFormatter after flags are parsed
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cmdCtx.OutputFormatter = cli.NewOutputFormatter(jsonOutput, verbose)
	}

	// Add all commands
	rootCmd.AddCommand(cli.NewPlayCommand(cmdCtx))
	rootCmd.AddCommand(cli.NewPauseCommand(cmdCtx))
	rootCmd.AddCommand(cli.NewStopCommand(cmdCtx))
	rootCmd.AddCommand(cli.NewResumeCommand(cmdCtx))
	rootCmd.AddCommand(cli.NewNextCommand(cmdCtx))
	rootCmd.AddCommand(cli.NewPreviousCommand(cmdCtx))
	rootCmd.AddCommand(cli.NewVolumeCommand(cmdCtx))
	rootCmd.AddCommand(cli.NewStatusCommand(cmdCtx))

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		if jsonOutput {
			fmt.Fprintf(os.Stderr, `{"error": "%s"}%s`, err.Error(), "\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
		os.Exit(1)
	}
}
