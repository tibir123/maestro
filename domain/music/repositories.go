package music

import (
	"context"
)

// PlayerRepository defines the interface for controlling music playback.
// This is a port that will be implemented by the infrastructure layer
// to provide actual Music.app control via AppleScript.
type PlayerRepository interface {
	// Play starts playback of the specified track
	Play(ctx context.Context, trackID TrackID) error

	// Pause pauses the current playback
	Pause(ctx context.Context) error

	// Stop stops playback and clears the current track
	Stop(ctx context.Context) error

	// Resume resumes paused playback
	Resume(ctx context.Context) error

	// Next advances to the next track in the current context
	Next(ctx context.Context) error

	// Previous goes back to the previous track in the current context
	Previous(ctx context.Context) error

	// Seek changes the playback position within the current track
	Seek(ctx context.Context, position Duration) error

	// SetVolume changes the playback volume
	SetVolume(ctx context.Context, volume Volume) error

	// SetShuffle enables or disables shuffle mode
	SetShuffle(ctx context.Context, enabled bool) error

	// SetRepeat changes the repeat mode
	SetRepeat(ctx context.Context, mode RepeatMode) error

	// GetCurrentState returns the current player state
	GetCurrentState(ctx context.Context) (*Player, error)

	// GetCurrentTrack returns the currently playing track, if any
	GetCurrentTrack(ctx context.Context) (*Track, error)
}

// LibrarySearchOptions contains options for searching the music library.
type LibrarySearchOptions struct {
	// Query is the search term to look for
	Query string

	// Artist filters results to tracks by specific artists
	Artist string

	// Album filters results to tracks from specific albums
	Album string

	// Limit is the maximum number of results to return (0 = no limit)
	Limit int

	// Offset is the number of results to skip (for pagination)
	Offset int
}

// LibraryRepository defines the interface for accessing the music library.
// This provides read-only access to tracks and playlists in Music.app.
type LibraryRepository interface {
	// Search finds tracks in the library based on the provided criteria
	Search(ctx context.Context, options LibrarySearchOptions) ([]*Track, error)

	// GetTrack retrieves a specific track by its ID
	GetTrack(ctx context.Context, trackID TrackID) (*Track, error)

	// GetTracks retrieves multiple tracks by their IDs
	GetTracks(ctx context.Context, trackIDs []TrackID) ([]*Track, error)

	// GetAllTracks returns all tracks in the library with pagination
	GetAllTracks(ctx context.Context, limit, offset int) ([]*Track, error)

	// GetTrackCount returns the total number of tracks in the library
	GetTrackCount(ctx context.Context) (int, error)

	// GetPlaylists returns all playlists accessible to the user
	GetPlaylists(ctx context.Context) ([]*Playlist, error)

	// GetPlaylist retrieves a specific playlist by its ID
	GetPlaylist(ctx context.Context, playlistID PlaylistID) (*Playlist, error)

	// GetPlaylistTracks returns all tracks in a specific playlist
	GetPlaylistTracks(ctx context.Context, playlistID PlaylistID) ([]*Track, error)

	// GetArtists returns a list of all artists in the library
	GetArtists(ctx context.Context) ([]string, error)

	// GetAlbums returns a list of all albums in the library
	GetAlbums(ctx context.Context) ([]string, error)

	// GetAlbumsByArtist returns albums by a specific artist
	GetAlbumsByArtist(ctx context.Context, artist string) ([]string, error)

	// GetTracksByArtist returns tracks by a specific artist
	GetTracksByArtist(ctx context.Context, artist string) ([]*Track, error)

	// GetTracksByAlbum returns tracks from a specific album
	GetTracksByAlbum(ctx context.Context, album string) ([]*Track, error)
}

// QueueRepository defines the interface for managing the playback queue.
// This provides control over what tracks will play next.
type QueueRepository interface {
	// GetQueue returns the current playback queue
	GetQueue(ctx context.Context) (*Playlist, error)

	// AddToQueue adds a track to the end of the queue
	AddToQueue(ctx context.Context, trackID TrackID) error

	// AddTracksToQueue adds multiple tracks to the end of the queue
	AddTracksToQueue(ctx context.Context, trackIDs []TrackID) error

	// PlayNext adds a track to play immediately after the current track
	PlayNext(ctx context.Context, trackID TrackID) error

	// PlayLater adds a track to the end of the queue (same as AddToQueue)
	PlayLater(ctx context.Context, trackID TrackID) error

	// RemoveFromQueue removes a track from the queue by position
	RemoveFromQueue(ctx context.Context, position int) error

	// ClearQueue removes all tracks from the queue
	ClearQueue(ctx context.Context) error

	// ShuffleQueue randomizes the order of tracks in the queue
	ShuffleQueue(ctx context.Context) error

	// GetQueuePosition returns the current position in the queue (0-based)
	GetQueuePosition(ctx context.Context) (int, error)

	// SetQueuePosition moves to a specific position in the queue
	SetQueuePosition(ctx context.Context, position int) error

	// GetUpNext returns the next few tracks that will play
	GetUpNext(ctx context.Context, count int) ([]*Track, error)
}

// PlaylistRepository defines the interface for managing playlists.
// This provides create, update, and delete operations for user playlists.
type PlaylistRepository interface {
	// CreatePlaylist creates a new user playlist
	CreatePlaylist(ctx context.Context, name string) (*Playlist, error)

	// UpdatePlaylist updates playlist metadata (name, etc.)
	UpdatePlaylist(ctx context.Context, playlist *Playlist) error

	// DeletePlaylist removes a playlist (only user-created playlists)
	DeletePlaylist(ctx context.Context, playlistID PlaylistID) error

	// AddTrackToPlaylist adds a track to a playlist
	AddTrackToPlaylist(ctx context.Context, playlistID PlaylistID, trackID TrackID) error

	// RemoveTrackFromPlaylist removes a track from a playlist
	RemoveTrackFromPlaylist(ctx context.Context, playlistID PlaylistID, trackID TrackID) error

	// ReorderPlaylistTracks changes the order of tracks in a playlist
	ReorderPlaylistTracks(ctx context.Context, playlistID PlaylistID, trackIDs []TrackID) error

	// DuplicatePlaylist creates a copy of an existing playlist
	DuplicatePlaylist(ctx context.Context, playlistID PlaylistID, newName string) (*Playlist, error)
}

// RepositoryManager aggregates all repository interfaces for convenience.
// This can be implemented by a single struct that composes all the individual repositories,
// or used as a service locator pattern.
type RepositoryManager interface {
	PlayerRepository
	LibraryRepository
	QueueRepository
	PlaylistRepository
}

// SearchResult represents a single search result with relevance scoring.
type SearchResult struct {
	// Track is the matching track
	Track *Track

	// Score indicates the relevance of this result (0.0 - 1.0)
	Score float64

	// MatchedFields indicates which fields matched the search query
	MatchedFields []string
}

// LibraryStats provides statistics about the music library.
type LibraryStats struct {
	// TotalTracks is the total number of tracks in the library
	TotalTracks int

	// TotalPlaylists is the total number of playlists
	TotalPlaylists int

	// TotalArtists is the total number of unique artists
	TotalArtists int

	// TotalAlbums is the total number of unique albums
	TotalAlbums int

	// TotalDuration is the total duration of all tracks
	TotalDuration Duration

	// LastUpdated indicates when these stats were calculated
	LastUpdated int64
}

// LibraryStatsRepository provides access to library statistics and metadata.
type LibraryStatsRepository interface {
	// GetStats returns current library statistics
	GetStats(ctx context.Context) (*LibraryStats, error)

	// RefreshStats recalculates library statistics
	RefreshStats(ctx context.Context) (*LibraryStats, error)
}
