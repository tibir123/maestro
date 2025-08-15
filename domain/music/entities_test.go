package music

import (
	"errors"
	"testing"
)

func TestNewTrack(t *testing.T) {
	tests := []struct {
		name          string
		id            TrackID
		title         string
		artist        string
		album         string
		duration      Duration
		expectedError bool
		errorType     error
	}{
		{
			name:          "valid track",
			id:            NewTrackID("track-123"),
			title:         "Test Song",
			artist:        "Test Artist",
			album:         "Test Album",
			duration:      NewDuration(180),
			expectedError: false,
		},
		{
			name:          "empty track ID",
			id:            NewTrackID(""),
			title:         "Test Song",
			artist:        "Test Artist",
			album:         "Test Album",
			duration:      NewDuration(180),
			expectedError: true,
			errorType:     ErrInvalidTrackID,
		},
		{
			name:          "empty title",
			id:            NewTrackID("track-123"),
			title:         "",
			artist:        "Test Artist",
			album:         "Test Album",
			duration:      NewDuration(180),
			expectedError: true,
			errorType:     ErrInvalidTrack,
		},
		{
			name:          "empty artist",
			id:            NewTrackID("track-123"),
			title:         "Test Song",
			artist:        "",
			album:         "Test Album",
			duration:      NewDuration(180),
			expectedError: true,
			errorType:     ErrInvalidTrack,
		},
		{
			name:          "invalid duration",
			id:            NewTrackID("track-123"),
			title:         "Test Song",
			artist:        "Test Artist",
			album:         "Test Album",
			duration:      Duration{seconds: -1}, // Create invalid duration directly
			expectedError: true,
			errorType:     ErrInvalidTrack,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track, err := NewTrack(tt.id, tt.title, tt.artist, tt.album, tt.duration)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}

				var domainErr *DomainError
				if !IsType(err, &domainErr) {
					t.Errorf("expected DomainError, got %T", err)
					return
				}

				if !IsError(err, tt.errorType) {
					t.Errorf("expected error type %v, got %v", tt.errorType, domainErr.Code)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if track.ID != tt.id {
				t.Errorf("expected ID %v, got %v", tt.id, track.ID)
			}
			if track.Title != tt.title {
				t.Errorf("expected title %s, got %s", tt.title, track.Title)
			}
			if track.Artist != tt.artist {
				t.Errorf("expected artist %s, got %s", tt.artist, track.Artist)
			}
			if track.Album != tt.album {
				t.Errorf("expected album %s, got %s", tt.album, track.Album)
			}
			if track.Duration != tt.duration {
				t.Errorf("expected duration %v, got %v", tt.duration, track.Duration)
			}
		})
	}
}

func TestTrackEquals(t *testing.T) {
	id1 := NewTrackID("track-1")
	id2 := NewTrackID("track-2")

	track1, _ := NewTrack(id1, "Song 1", "Artist 1", "Album 1", NewDuration(180))
	track2, _ := NewTrack(id1, "Different Song", "Different Artist", "Different Album", NewDuration(240))
	track3, _ := NewTrack(id2, "Song 1", "Artist 1", "Album 1", NewDuration(180))

	if !track1.Equals(track2) {
		t.Error("tracks with same ID should be equal")
	}

	if track1.Equals(track3) {
		t.Error("tracks with different IDs should not be equal")
	}

	if track1.Equals(nil) {
		t.Error("track should not equal nil")
	}
}

func TestTrackString(t *testing.T) {
	track, _ := NewTrack(NewTrackID("track-1"), "Test Song", "Test Artist", "Test Album", NewDuration(180))
	expected := "Test Artist - Test Song"

	if track.String() != expected {
		t.Errorf("expected string %s, got %s", expected, track.String())
	}
}

func TestNewPlaylist(t *testing.T) {
	tests := []struct {
		name          string
		id            PlaylistID
		playlistName  string
		playlistType  PlaylistType
		readOnly      bool
		expectedError bool
		errorType     error
	}{
		{
			name:          "valid playlist",
			id:            NewPlaylistID("playlist-123"),
			playlistName:  "My Playlist",
			playlistType:  PlaylistTypeUser,
			readOnly:      false,
			expectedError: false,
		},
		{
			name:          "empty playlist ID",
			id:            NewPlaylistID(""),
			playlistName:  "My Playlist",
			playlistType:  PlaylistTypeUser,
			readOnly:      false,
			expectedError: true,
			errorType:     ErrInvalidPlaylistID,
		},
		{
			name:          "empty playlist name",
			id:            NewPlaylistID("playlist-123"),
			playlistName:  "",
			playlistType:  PlaylistTypeUser,
			readOnly:      false,
			expectedError: true,
			errorType:     ErrInvalidPlaylist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playlist, err := NewPlaylist(tt.id, tt.playlistName, tt.playlistType, tt.readOnly)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}

				var domainErr *DomainError
				if !IsType(err, &domainErr) {
					t.Errorf("expected DomainError, got %T", err)
					return
				}

				if !IsError(err, tt.errorType) {
					t.Errorf("expected error type %v, got %v", tt.errorType, domainErr.Code)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if playlist.ID != tt.id {
				t.Errorf("expected ID %v, got %v", tt.id, playlist.ID)
			}
			if playlist.Name != tt.playlistName {
				t.Errorf("expected name %s, got %s", tt.playlistName, playlist.Name)
			}
			if playlist.Type != tt.playlistType {
				t.Errorf("expected type %v, got %v", tt.playlistType, playlist.Type)
			}
			if playlist.ReadOnly != tt.readOnly {
				t.Errorf("expected readOnly %v, got %v", tt.readOnly, playlist.ReadOnly)
			}
			if len(playlist.Tracks) != 0 {
				t.Errorf("expected empty tracks list, got %d tracks", len(playlist.Tracks))
			}
		})
	}
}

func TestPlaylistAddTrack(t *testing.T) {
	playlist, _ := NewPlaylist(NewPlaylistID("playlist-1"), "Test Playlist", PlaylistTypeUser, false)
	readOnlyPlaylist, _ := NewPlaylist(NewPlaylistID("playlist-2"), "Read Only", PlaylistTypeLibrary, true)

	trackID1 := NewTrackID("track-1")
	trackID2 := NewTrackID("track-2")
	emptyTrackID := NewTrackID("")

	// Test adding to regular playlist
	err := playlist.AddTrack(trackID1)
	if err != nil {
		t.Errorf("unexpected error adding track: %v", err)
	}

	if playlist.TrackCount() != 1 {
		t.Errorf("expected 1 track, got %d", playlist.TrackCount())
	}

	if !playlist.ContainsTrack(trackID1) {
		t.Error("playlist should contain added track")
	}

	// Test adding duplicate track
	err = playlist.AddTrack(trackID1)
	if err == nil {
		t.Error("expected error when adding duplicate track")
	}
	if !IsError(err, ErrTrackAlreadyInPlaylist) {
		t.Errorf("expected ErrTrackAlreadyInPlaylist, got %v", err)
	}

	// Test adding to read-only playlist
	err = readOnlyPlaylist.AddTrack(trackID1)
	if err == nil {
		t.Error("expected error when adding to read-only playlist")
	}
	if !IsError(err, ErrPlaylistReadOnly) {
		t.Errorf("expected ErrPlaylistReadOnly, got %v", err)
	}

	// Test adding empty track ID
	err = playlist.AddTrack(emptyTrackID)
	if err == nil {
		t.Error("expected error when adding empty track ID")
	}
	if !IsError(err, ErrInvalidTrackID) {
		t.Errorf("expected ErrInvalidTrackID, got %v", err)
	}

	// Test adding another valid track
	err = playlist.AddTrack(trackID2)
	if err != nil {
		t.Errorf("unexpected error adding second track: %v", err)
	}

	if playlist.TrackCount() != 2 {
		t.Errorf("expected 2 tracks, got %d", playlist.TrackCount())
	}
}

func TestPlaylistRemoveTrack(t *testing.T) {
	playlist, _ := NewPlaylist(NewPlaylistID("playlist-1"), "Test Playlist", PlaylistTypeUser, false)
	readOnlyPlaylist, _ := NewPlaylist(NewPlaylistID("playlist-2"), "Read Only", PlaylistTypeLibrary, true)

	trackID1 := NewTrackID("track-1")
	trackID2 := NewTrackID("track-2")
	trackID3 := NewTrackID("track-3")

	// Add some tracks
	_ = playlist.AddTrack(trackID1)
	_ = playlist.AddTrack(trackID2)

	// Test removing existing track
	err := playlist.RemoveTrack(trackID1)
	if err != nil {
		t.Errorf("unexpected error removing track: %v", err)
	}

	if playlist.TrackCount() != 1 {
		t.Errorf("expected 1 track after removal, got %d", playlist.TrackCount())
	}

	if playlist.ContainsTrack(trackID1) {
		t.Error("playlist should not contain removed track")
	}

	// Test removing non-existent track
	err = playlist.RemoveTrack(trackID3)
	if err == nil {
		t.Error("expected error when removing non-existent track")
	}
	if !IsError(err, ErrTrackNotFound) {
		t.Errorf("expected ErrTrackNotFound, got %v", err)
	}

	// Test removing from read-only playlist
	err = readOnlyPlaylist.RemoveTrack(trackID1)
	if err == nil {
		t.Error("expected error when removing from read-only playlist")
	}
	if !IsError(err, ErrPlaylistReadOnly) {
		t.Errorf("expected ErrPlaylistReadOnly, got %v", err)
	}
}

func TestPlaylistEquals(t *testing.T) {
	id1 := NewPlaylistID("playlist-1")
	id2 := NewPlaylistID("playlist-2")

	playlist1, _ := NewPlaylist(id1, "Playlist 1", PlaylistTypeUser, false)
	playlist2, _ := NewPlaylist(id1, "Different Name", PlaylistTypeSmart, true)
	playlist3, _ := NewPlaylist(id2, "Playlist 1", PlaylistTypeUser, false)

	if !playlist1.Equals(playlist2) {
		t.Error("playlists with same ID should be equal")
	}

	if playlist1.Equals(playlist3) {
		t.Error("playlists with different IDs should not be equal")
	}

	if playlist1.Equals(nil) {
		t.Error("playlist should not equal nil")
	}
}

func TestNewPlayer(t *testing.T) {
	player := NewPlayer()

	if player.State != PlayerStateStopped {
		t.Errorf("expected initial state to be stopped, got %v", player.State)
	}

	if player.CurrentTrack != nil {
		t.Error("expected no current track initially")
	}

	if !player.Position.IsZero() {
		t.Error("expected initial position to be zero")
	}

	if player.Volume.Level() != 50 {
		t.Errorf("expected initial volume to be 50, got %d", player.Volume.Level())
	}

	if player.Shuffle {
		t.Error("expected shuffle to be disabled initially")
	}

	if player.Repeat != RepeatModeOff {
		t.Errorf("expected repeat to be off initially, got %v", player.Repeat)
	}
}

func TestPlayerSetVolume(t *testing.T) {
	player := NewPlayer()

	// Test valid volume
	validVolume := NewVolume(75)
	err := player.SetVolume(validVolume)
	if err != nil {
		t.Errorf("unexpected error setting valid volume: %v", err)
	}

	if player.Volume.Level() != 75 {
		t.Errorf("expected volume 75, got %d", player.Volume.Level())
	}

	// Test invalid volume (this shouldn't happen with NewVolume, but test anyway)
	invalidVolume := Volume{level: 150}
	err = player.SetVolume(invalidVolume)
	if err == nil {
		t.Error("expected error setting invalid volume")
	}
}

func TestPlayerPlayback(t *testing.T) {
	player := NewPlayer()
	trackID := NewTrackID("track-1")

	// Test play
	player.Play(&trackID)
	if !player.IsPlaying() {
		t.Error("player should be playing after Play()")
	}
	if player.CurrentTrack == nil || !player.CurrentTrack.Equals(trackID) {
		t.Error("current track should be set after Play()")
	}

	// Test pause
	player.Pause()
	if !player.IsPaused() {
		t.Error("player should be paused after Pause()")
	}
	if player.CurrentTrack == nil || !player.CurrentTrack.Equals(trackID) {
		t.Error("current track should remain set after Pause()")
	}

	// Test resume
	player.Resume()
	if !player.IsPlaying() {
		t.Error("player should be playing after Resume()")
	}

	// Test stop
	player.Stop()
	if !player.IsStopped() {
		t.Error("player should be stopped after Stop()")
	}
	if player.CurrentTrack != nil {
		t.Error("current track should be cleared after Stop()")
	}
	if !player.Position.IsZero() {
		t.Error("position should be reset after Stop()")
	}
}

func TestPlayerState(t *testing.T) {
	player := NewPlayer()

	// Test initial state
	if !player.IsStopped() {
		t.Error("player should initially be stopped")
	}
	if player.IsPlaying() || player.IsPaused() {
		t.Error("player should not be playing or paused initially")
	}

	// Test has current track
	if player.HasCurrentTrack() {
		t.Error("player should not have current track initially")
	}

	trackID := NewTrackID("track-1")
	player.Play(&trackID)

	if !player.HasCurrentTrack() {
		t.Error("player should have current track after play")
	}
}

func TestPlayerSettings(t *testing.T) {
	player := NewPlayer()

	// Test shuffle
	player.SetShuffle(true)
	if !player.Shuffle {
		t.Error("shuffle should be enabled")
	}

	player.SetShuffle(false)
	if player.Shuffle {
		t.Error("shuffle should be disabled")
	}

	// Test repeat
	player.SetRepeat(RepeatModeAll)
	if player.Repeat != RepeatModeAll {
		t.Errorf("expected repeat mode all, got %v", player.Repeat)
	}

	player.SetRepeat(RepeatModeOne)
	if player.Repeat != RepeatModeOne {
		t.Errorf("expected repeat mode one, got %v", player.Repeat)
	}
}

// Helper functions for testing
func IsType(err error, target interface{}) bool {
	switch target.(type) {
	case **DomainError:
		var domainErr *DomainError
		return errors.As(err, &domainErr)
	default:
		return false
	}
}

func IsError(err error, target error) bool {
	return errors.Is(err, target)
}
