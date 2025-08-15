package applescript

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/madstone-tech/maestro/domain/music"
)

// PlayerRepository implements the music.PlayerRepository interface using AppleScript
// to control Music.app on macOS.
type PlayerRepository struct {
	executor *Executor
}

// NewPlayerRepository creates a new AppleScript-based player repository.
func NewPlayerRepository(executor *Executor) *PlayerRepository {
	if executor == nil {
		executor = NewExecutor(nil)
	}

	return &PlayerRepository{
		executor: executor,
	}
}

// Play starts playback of the specified track.
func (p *PlayerRepository) Play(ctx context.Context, trackID music.TrackID) error {
	if trackID.IsEmpty() {
		return music.NewDomainError(music.ErrInvalidTrackID, "track ID cannot be empty")
	}

	script := fmt.Sprintf(`
		tell application "Music"
			try
				set theTrack to track id %s
				play theTrack
			on error errMsg
				error "Failed to play track: " & errMsg
			end try
		end tell
	`, trackID.Value())

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(
			music.ErrOperationFailed,
			fmt.Sprintf("failed to play track %s", trackID.Value()),
			result.Error,
		)
	}

	return nil
}

// Pause pauses the current playback.
func (p *PlayerRepository) Pause(ctx context.Context) error {
	script := `
		tell application "Music"
			pause
		end tell
	`

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to pause playback", result.Error)
	}

	return nil
}

// Stop stops playback and clears the current track.
func (p *PlayerRepository) Stop(ctx context.Context) error {
	script := `
		tell application "Music"
			stop
		end tell
	`

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to stop playback", result.Error)
	}

	return nil
}

// Resume resumes paused playback.
func (p *PlayerRepository) Resume(ctx context.Context) error {
	script := `
		tell application "Music"
			play
		end tell
	`

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to resume playback", result.Error)
	}

	return nil
}

// Next advances to the next track in the current context.
func (p *PlayerRepository) Next(ctx context.Context) error {
	script := `
		tell application "Music"
			next track
		end tell
	`

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to skip to next track", result.Error)
	}

	return nil
}

// Previous goes back to the previous track in the current context.
func (p *PlayerRepository) Previous(ctx context.Context) error {
	script := `
		tell application "Music"
			previous track
		end tell
	`

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to skip to previous track", result.Error)
	}

	return nil
}

// Seek changes the playback position within the current track.
func (p *PlayerRepository) Seek(ctx context.Context, position music.Duration) error {
	if !position.IsValid() {
		return music.WrapInvalidPosition(position, nil)
	}

	script := fmt.Sprintf(`
		tell application "Music"
			try
				set player position to %d
			on error errMsg
				error "Failed to seek: " & errMsg
			end try
		end tell
	`, position.Seconds())

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(
			music.ErrOperationFailed,
			fmt.Sprintf("failed to seek to position %s", position.String()),
			result.Error,
		)
	}

	return nil
}

// SetVolume changes the playback volume.
func (p *PlayerRepository) SetVolume(ctx context.Context, volume music.Volume) error {
	if !volume.IsValid() {
		return music.WrapInvalidVolume(volume.Level(), nil)
	}

	script := fmt.Sprintf(`
		tell application "Music"
			set sound volume to %d
		end tell
	`, volume.Level())

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(
			music.ErrOperationFailed,
			fmt.Sprintf("failed to set volume to %s", volume.String()),
			result.Error,
		)
	}

	return nil
}

// SetShuffle enables or disables shuffle mode.
func (p *PlayerRepository) SetShuffle(ctx context.Context, enabled bool) error {
	script := fmt.Sprintf(`
		tell application "Music"
			set shuffle enabled to %t
		end tell
	`, enabled)

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(
			music.ErrOperationFailed,
			fmt.Sprintf("failed to set shuffle mode to %t", enabled),
			result.Error,
		)
	}

	return nil
}

// SetRepeat changes the repeat mode.
func (p *PlayerRepository) SetRepeat(ctx context.Context, mode music.RepeatMode) error {
	if !mode.IsValid() {
		return music.NewDomainError(music.ErrInvalidRepeatMode, "invalid repeat mode")
	}

	var repeatSetting string
	switch mode {
	case music.RepeatModeOff:
		repeatSetting = "off"
	case music.RepeatModeAll:
		repeatSetting = "all"
	case music.RepeatModeOne:
		repeatSetting = "one"
	default:
		return music.NewDomainError(music.ErrInvalidRepeatMode, "unsupported repeat mode")
	}

	script := fmt.Sprintf(`
		tell application "Music"
			set song repeat to %s
		end tell
	`, repeatSetting)

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(
			music.ErrOperationFailed,
			fmt.Sprintf("failed to set repeat mode to %s", mode.String()),
			result.Error,
		)
	}

	return nil
}

// GetCurrentState returns the current player state.
func (p *PlayerRepository) GetCurrentState(ctx context.Context) (*music.Player, error) {
	script := `
		tell application "Music"
			try
				set playerState to player state as string
				set playerVolume to sound volume
				set playerPosition to player position
				set shuffleState to shuffle enabled
				set repeatState to song repeat as string
				
				set currentTrackID to ""
				if player state is not stopped then
					try
						set currentTrackID to (database ID of current track) as string
					end try
				end if
				
				return playerState & "|" & playerVolume & "|" & playerPosition & "|" & shuffleState & "|" & repeatState & "|" & currentTrackID
			on error errMsg
				error "Failed to get player state: " & errMsg
			end try
		end tell
	`

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return nil, music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to get current player state", result.Error)
	}

	return p.parsePlayerState(result.Output)
}

// GetCurrentTrack returns the currently playing track, if any.
func (p *PlayerRepository) GetCurrentTrack(ctx context.Context) (*music.Track, error) {
	script := `
		tell application "Music"
			try
				if player state is stopped then
					return ""
				end if
				
				set theTrack to current track
				set trackID to (database ID of theTrack) as string
				set trackName to name of theTrack
				set trackArtist to artist of theTrack
				set trackAlbum to album of theTrack
				set trackDuration to duration of theTrack
				
				return trackID & "|" & trackName & "|" & trackArtist & "|" & trackAlbum & "|" & trackDuration
			on error errMsg
				error "Failed to get current track: " & errMsg
			end try
		end tell
	`

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return nil, music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to get current track", result.Error)
	}

	if strings.TrimSpace(result.Output) == "" {
		return nil, nil // No current track
	}

	return p.parseTrack(result.Output)
}

// parsePlayerState parses the player state output from AppleScript.
func (p *PlayerRepository) parsePlayerState(output string) (*music.Player, error) {
	parts := strings.Split(strings.TrimSpace(output), "|")
	if len(parts) != 6 {
		return nil, music.NewDomainError(music.ErrOperationFailed, "invalid player state format")
	}

	// Parse player state
	var state music.PlayerState
	switch strings.ToLower(parts[0]) {
	case "stopped":
		state = music.PlayerStateStopped
	case "playing":
		state = music.PlayerStatePlaying
	case "paused":
		state = music.PlayerStatePaused
	default:
		state = music.PlayerStateStopped
	}

	// Parse volume
	volumeLevel, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, music.NewDomainError(music.ErrOperationFailed, "invalid volume format")
	}
	volume := music.NewVolume(volumeLevel)

	// Parse position
	var position music.Duration
	if parts[2] == "missing value" || parts[2] == "" {
		// When stopped or no current track, position is 0
		position = music.NewDuration(0)
	} else {
		positionSeconds, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, music.NewDomainError(music.ErrOperationFailed, "invalid position format")
		}
		position = music.NewDuration(int(positionSeconds))
	}

	// Parse shuffle
	shuffle := strings.ToLower(parts[3]) == "true"

	// Parse repeat mode
	var repeat music.RepeatMode
	switch strings.ToLower(parts[4]) {
	case "off":
		repeat = music.RepeatModeOff
	case "all":
		repeat = music.RepeatModeAll
	case "one":
		repeat = music.RepeatModeOne
	default:
		repeat = music.RepeatModeOff
	}

	// Parse current track ID
	var currentTrack *music.TrackID
	if parts[5] != "" {
		trackID := music.NewTrackID(parts[5])
		currentTrack = &trackID
	}

	player := music.NewPlayer()
	player.State = state
	player.Volume = volume
	player.Position = position
	player.Shuffle = shuffle
	player.Repeat = repeat
	player.CurrentTrack = currentTrack

	return player, nil
}

// parseTrack parses track information from AppleScript output.
func (p *PlayerRepository) parseTrack(output string) (*music.Track, error) {
	parts := strings.Split(strings.TrimSpace(output), "|")
	if len(parts) != 5 {
		return nil, music.NewDomainError(music.ErrOperationFailed, "invalid track format")
	}

	trackID := music.NewTrackID(parts[0])
	title := parts[1]
	artist := parts[2]
	album := parts[3]

	// Parse duration (in seconds)
	durationSeconds, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, music.NewDomainError(music.ErrOperationFailed, "invalid track duration format")
	}
	duration := music.NewDuration(int(durationSeconds))

	return music.NewTrack(trackID, title, artist, album, duration)
}

// HealthCheck performs a basic health check to ensure Music.app is accessible.
func (p *PlayerRepository) HealthCheck(ctx context.Context) error {
	script := `
		tell application "Music"
			try
				get name
				return "ok"
			on error errMsg
				error "Music app not accessible: " & errMsg
			end try
		end tell
	`

	result := p.executor.Execute(ctx, script)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(music.ErrPlayerNotAvailable, "Music.app health check failed", result.Error)
	}

	if result.Output != "ok" {
		return music.NewDomainError(music.ErrPlayerNotAvailable, "Music.app returned unexpected response")
	}

	return nil
}
