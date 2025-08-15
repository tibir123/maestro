package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/madstone-tech/maestro/domain/music"
)

// OutputFormatter handles formatting and displaying command output
type OutputFormatter struct {
	jsonMode bool
	verbose  bool
	writer   io.Writer
}

// NewOutputFormatter creates a new output formatter
func NewOutputFormatter(jsonMode, verbose bool) *OutputFormatter {
	return &OutputFormatter{
		jsonMode: jsonMode,
		verbose:  verbose,
		writer:   os.Stdout,
	}
}

// SetWriter sets the output writer (useful for testing)
func (f *OutputFormatter) SetWriter(w io.Writer) {
	f.writer = w
}

// Success prints a success message
func (f *OutputFormatter) Success(message string) {
	if f.jsonMode {
		f.printJSON(map[string]interface{}{
			"success": true,
			"message": message,
		})
	} else {
		fmt.Fprintf(f.writer, "%s\n", message)
	}
}

// Error prints an error message
func (f *OutputFormatter) Error(err error) {
	if f.jsonMode {
		f.printJSON(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	} else {
		fmt.Fprintf(f.writer, "Error: %s\n", err.Error())
	}
}

// PrintPlayerStatus prints the current player status
func (f *OutputFormatter) PrintPlayerStatus(player *music.Player, track *music.Track) {
	if f.jsonMode {
		f.printPlayerStatusJSON(player, track)
	} else {
		f.printPlayerStatusText(player, track)
	}
}

// PrintTrack prints track information
func (f *OutputFormatter) PrintTrack(track *music.Track) {
	if f.jsonMode {
		f.printJSON(track)
	} else {
		if track != nil {
			fmt.Fprintf(f.writer, "%s - %s\n", track.Artist, track.Title)
			if track.Album != "" {
				fmt.Fprintf(f.writer, "Album: %s\n", track.Album)
			}
			fmt.Fprintf(f.writer, "Duration: %s\n", track.Duration.String())
		} else {
			fmt.Fprintf(f.writer, "No current track\n")
		}
	}
}

// PrintVolume prints volume information
func (f *OutputFormatter) PrintVolume(volume music.Volume) {
	if f.jsonMode {
		f.printJSON(map[string]interface{}{
			"volume": volume.Level(),
		})
	} else {
		fmt.Fprintf(f.writer, "Volume: %s\n", volume.String())
	}
}

// printPlayerStatusJSON prints player status in JSON format
func (f *OutputFormatter) printPlayerStatusJSON(player *music.Player, track *music.Track) {
	status := map[string]interface{}{
		"state":    player.State.String(),
		"volume":   player.Volume.Level(),
		"position": player.Position.String(),
		"shuffle":  player.Shuffle,
		"repeat":   player.Repeat.String(),
	}

	if track != nil {
		status["current_track"] = map[string]interface{}{
			"id":       track.ID.Value(),
			"title":    track.Title,
			"artist":   track.Artist,
			"album":    track.Album,
			"duration": track.Duration.String(),
		}
	} else {
		status["current_track"] = nil
	}

	f.printJSON(status)
}

// printPlayerStatusText prints player status in human-readable format
func (f *OutputFormatter) printPlayerStatusText(player *music.Player, track *music.Track) {
	fmt.Fprintf(f.writer, "Status: %s\n", player.State.String())
	fmt.Fprintf(f.writer, "Volume: %s\n", player.Volume.String())

	if track != nil {
		fmt.Fprintf(f.writer, "Now Playing: %s - %s\n", track.Artist, track.Title)
		if track.Album != "" {
			fmt.Fprintf(f.writer, "Album: %s\n", track.Album)
		}
		fmt.Fprintf(f.writer, "Position: %s / %s\n", player.Position.String(), track.Duration.String())
	} else {
		fmt.Fprintf(f.writer, "No current track\n")
	}

	fmt.Fprintf(f.writer, "Shuffle: %t\n", player.Shuffle)
	fmt.Fprintf(f.writer, "Repeat: %s\n", player.Repeat.String())
}

// printJSON prints data in JSON format
func (f *OutputFormatter) printJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(f.writer, `{"error": "Failed to format JSON output"}%s`, "\n")
		return
	}
	fmt.Fprintf(f.writer, "%s\n", jsonData)
}

// Debug prints debug information if verbose mode is enabled
func (f *OutputFormatter) Debug(message string) {
	if f.verbose {
		if f.jsonMode {
			f.printJSON(map[string]interface{}{
				"debug": message,
			})
		} else {
			fmt.Fprintf(f.writer, "DEBUG: %s\n", message)
		}
	}
}

// Info prints informational messages
func (f *OutputFormatter) Info(message string) {
	if f.jsonMode {
		f.printJSON(map[string]interface{}{
			"info": message,
		})
	} else {
		fmt.Fprintf(f.writer, "%s\n", message)
	}
}
