// Package music contains the core domain entities for Maestro music controller.
// This package follows Domain-Driven Design principles and contains pure business logic
// with no external dependencies.
package music

import (
	"time"
)

// Track represents a music track in the domain.
// It encapsulates all the essential information about a track including metadata
// and provides the core business logic for track operations.
type Track struct {
	// ID uniquely identifies the track within the music library
	ID TrackID `json:"id"`

	// Title is the name of the track
	Title string `json:"title"`

	// Artist is the name of the track's artist or performer
	Artist string `json:"artist"`

	// Album is the name of the album this track belongs to
	Album string `json:"album"`

	// Duration is the length of the track
	Duration Duration `json:"duration"`
}

// NewTrack creates a new Track entity with validation.
// It ensures all required fields are provided and valid.
func NewTrack(id TrackID, title, artist, album string, duration Duration) (*Track, error) {
	if id.IsEmpty() {
		return nil, NewDomainError(ErrInvalidTrackID, "track ID cannot be empty")
	}

	if title == "" {
		return nil, NewDomainError(ErrInvalidTrack, "track title cannot be empty")
	}

	if artist == "" {
		return nil, NewDomainError(ErrInvalidTrack, "track artist cannot be empty")
	}

	if !duration.IsValid() {
		return nil, NewDomainError(ErrInvalidTrack, "track duration must be valid")
	}

	return &Track{
		ID:       id,
		Title:    title,
		Artist:   artist,
		Album:    album,
		Duration: duration,
	}, nil
}

// Equals checks if two tracks are the same based on their ID.
func (t *Track) Equals(other *Track) bool {
	if other == nil {
		return false
	}
	return t.ID.Equals(other.ID)
}

// String returns a human-readable representation of the track.
func (t *Track) String() string {
	return t.Artist + " - " + t.Title
}

// Playlist represents a collection of tracks in a specific order.
// It encapsulates playlist metadata and the ordered list of tracks.
type Playlist struct {
	// ID uniquely identifies the playlist
	ID PlaylistID `json:"id"`

	// Name is the display name of the playlist
	Name string `json:"name"`

	// Type indicates what kind of playlist this is (e.g., user-created, smart, library)
	Type PlaylistType `json:"type"`

	// ReadOnly indicates whether the playlist can be modified
	ReadOnly bool `json:"read_only"`

	// Tracks is the ordered list of tracks in this playlist
	Tracks []TrackID `json:"tracks"`

	// CreatedAt is when the playlist was created
	CreatedAt time.Time `json:"created_at"`

	// ModifiedAt is when the playlist was last modified
	ModifiedAt time.Time `json:"modified_at"`
}

// NewPlaylist creates a new Playlist entity with validation.
func NewPlaylist(id PlaylistID, name string, playlistType PlaylistType, readOnly bool) (*Playlist, error) {
	if id.IsEmpty() {
		return nil, NewDomainError(ErrInvalidPlaylistID, "playlist ID cannot be empty")
	}

	if name == "" {
		return nil, NewDomainError(ErrInvalidPlaylist, "playlist name cannot be empty")
	}

	now := time.Now()

	return &Playlist{
		ID:         id,
		Name:       name,
		Type:       playlistType,
		ReadOnly:   readOnly,
		Tracks:     make([]TrackID, 0),
		CreatedAt:  now,
		ModifiedAt: now,
	}, nil
}

// AddTrack adds a track to the playlist if it's not read-only.
func (p *Playlist) AddTrack(trackID TrackID) error {
	if p.ReadOnly {
		return NewDomainError(ErrPlaylistReadOnly, "cannot modify read-only playlist")
	}

	if trackID.IsEmpty() {
		return NewDomainError(ErrInvalidTrackID, "cannot add empty track ID to playlist")
	}

	// Check if track is already in playlist
	for _, existing := range p.Tracks {
		if existing.Equals(trackID) {
			return NewDomainError(ErrTrackAlreadyInPlaylist, "track is already in playlist")
		}
	}

	p.Tracks = append(p.Tracks, trackID)
	p.ModifiedAt = time.Now()

	return nil
}

// RemoveTrack removes a track from the playlist if it's not read-only.
func (p *Playlist) RemoveTrack(trackID TrackID) error {
	if p.ReadOnly {
		return NewDomainError(ErrPlaylistReadOnly, "cannot modify read-only playlist")
	}

	for i, existing := range p.Tracks {
		if existing.Equals(trackID) {
			// Remove track by slicing
			p.Tracks = append(p.Tracks[:i], p.Tracks[i+1:]...)
			p.ModifiedAt = time.Now()
			return nil
		}
	}

	return NewDomainError(ErrTrackNotFound, "track not found in playlist")
}

// ContainsTrack checks if the playlist contains the specified track.
func (p *Playlist) ContainsTrack(trackID TrackID) bool {
	for _, existing := range p.Tracks {
		if existing.Equals(trackID) {
			return true
		}
	}
	return false
}

// TrackCount returns the number of tracks in the playlist.
func (p *Playlist) TrackCount() int {
	return len(p.Tracks)
}

// IsEmpty returns true if the playlist has no tracks.
func (p *Playlist) IsEmpty() bool {
	return len(p.Tracks) == 0
}

// Equals checks if two playlists are the same based on their ID.
func (p *Playlist) Equals(other *Playlist) bool {
	if other == nil {
		return false
	}
	return p.ID.Equals(other.ID)
}

// Player represents the state and configuration of the music player.
// It encapsulates the current playback state, position, and player settings.
type Player struct {
	// State indicates the current playback state (playing, paused, stopped)
	State PlayerState `json:"state"`

	// CurrentTrack is the track currently loaded in the player (may be nil)
	CurrentTrack *TrackID `json:"current_track,omitempty"`

	// Position is the current playback position within the track
	Position Duration `json:"position"`

	// Volume is the current playback volume (0-100)
	Volume Volume `json:"volume"`

	// Shuffle indicates whether shuffle mode is enabled
	Shuffle bool `json:"shuffle"`

	// Repeat indicates the current repeat mode
	Repeat RepeatMode `json:"repeat"`

	// LastUpdated is when the player state was last updated
	LastUpdated time.Time `json:"last_updated"`
}

// NewPlayer creates a new Player entity with default settings.
func NewPlayer() *Player {
	return &Player{
		State:        PlayerStateStopped,
		CurrentTrack: nil,
		Position:     NewDuration(0),
		Volume:       NewVolume(50), // Default to 50% volume
		Shuffle:      false,
		Repeat:       RepeatModeOff,
		LastUpdated:  time.Now(),
	}
}

// SetVolume changes the player volume with validation.
func (p *Player) SetVolume(volume Volume) error {
	if !volume.IsValid() {
		return NewDomainError(ErrInvalidVolume, "volume must be between 0 and 100")
	}

	p.Volume = volume
	p.LastUpdated = time.Now()

	return nil
}

// SetPosition changes the playback position with validation.
func (p *Player) SetPosition(position Duration) error {
	if !position.IsValid() {
		return NewDomainError(ErrInvalidPosition, "position must be valid")
	}

	p.Position = position
	p.LastUpdated = time.Now()

	return nil
}

// Play starts playback with the specified track.
func (p *Player) Play(trackID *TrackID) {
	p.State = PlayerStatePlaying
	p.CurrentTrack = trackID
	p.LastUpdated = time.Now()
}

// Pause pauses playback without changing the current track.
func (p *Player) Pause() {
	if p.State == PlayerStatePlaying {
		p.State = PlayerStatePaused
		p.LastUpdated = time.Now()
	}
}

// Stop stops playback and clears the current track.
func (p *Player) Stop() {
	p.State = PlayerStateStopped
	p.CurrentTrack = nil
	p.Position = NewDuration(0)
	p.LastUpdated = time.Now()
}

// Resume resumes playback if currently paused.
func (p *Player) Resume() {
	if p.State == PlayerStatePaused {
		p.State = PlayerStatePlaying
		p.LastUpdated = time.Now()
	}
}

// SetShuffle enables or disables shuffle mode.
func (p *Player) SetShuffle(enabled bool) {
	p.Shuffle = enabled
	p.LastUpdated = time.Now()
}

// SetRepeat changes the repeat mode.
func (p *Player) SetRepeat(mode RepeatMode) {
	p.Repeat = mode
	p.LastUpdated = time.Now()
}

// IsPlaying returns true if the player is currently playing.
func (p *Player) IsPlaying() bool {
	return p.State == PlayerStatePlaying
}

// IsPaused returns true if the player is currently paused.
func (p *Player) IsPaused() bool {
	return p.State == PlayerStatePaused
}

// IsStopped returns true if the player is currently stopped.
func (p *Player) IsStopped() bool {
	return p.State == PlayerStateStopped
}

// HasCurrentTrack returns true if there is a track currently loaded.
func (p *Player) HasCurrentTrack() bool {
	return p.CurrentTrack != nil && !p.CurrentTrack.IsEmpty()
}
