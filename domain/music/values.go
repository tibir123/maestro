package music

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TrackID is a value object representing a unique track identifier.
// It provides type safety and validation for track identification.
type TrackID struct {
	value string
}

// NewTrackID creates a new TrackID with validation.
func NewTrackID(value string) TrackID {
	return TrackID{value: strings.TrimSpace(value)}
}

// Value returns the underlying string value of the TrackID.
func (t TrackID) Value() string {
	return t.value
}

// IsEmpty returns true if the TrackID has no value.
func (t TrackID) IsEmpty() bool {
	return t.value == ""
}

// Equals checks if two TrackIDs are the same.
func (t TrackID) Equals(other TrackID) bool {
	return t.value == other.value
}

// String returns the string representation of the TrackID.
func (t TrackID) String() string {
	return t.value
}

// PlaylistID is a value object representing a unique playlist identifier.
// It provides type safety and validation for playlist identification.
type PlaylistID struct {
	value string
}

// NewPlaylistID creates a new PlaylistID with validation.
func NewPlaylistID(value string) PlaylistID {
	return PlaylistID{value: strings.TrimSpace(value)}
}

// Value returns the underlying string value of the PlaylistID.
func (p PlaylistID) Value() string {
	return p.value
}

// IsEmpty returns true if the PlaylistID has no value.
func (p PlaylistID) IsEmpty() bool {
	return p.value == ""
}

// Equals checks if two PlaylistIDs are the same.
func (p PlaylistID) Equals(other PlaylistID) bool {
	return p.value == other.value
}

// String returns the string representation of the PlaylistID.
func (p PlaylistID) String() string {
	return p.value
}

// Duration is a value object representing a time duration in seconds.
// It provides validation and utility methods for working with track durations and positions.
type Duration struct {
	seconds int
}

// NewDuration creates a new Duration with validation.
func NewDuration(seconds int) Duration {
	if seconds < 0 {
		seconds = 0
	}
	return Duration{seconds: seconds}
}

// NewDurationFromTime creates a Duration from a time.Duration.
func NewDurationFromTime(d time.Duration) Duration {
	return NewDuration(int(d.Seconds()))
}

// Seconds returns the duration in seconds.
func (d Duration) Seconds() int {
	return d.seconds
}

// Minutes returns the duration in minutes (rounded down).
func (d Duration) Minutes() int {
	return d.seconds / 60
}

// Hours returns the duration in hours (rounded down).
func (d Duration) Hours() int {
	return d.seconds / 3600
}

// IsValid returns true if the duration is non-negative.
func (d Duration) IsValid() bool {
	return d.seconds >= 0
}

// IsZero returns true if the duration is zero.
func (d Duration) IsZero() bool {
	return d.seconds == 0
}

// Add adds another duration to this one.
func (d Duration) Add(other Duration) Duration {
	return NewDuration(d.seconds + other.seconds)
}

// Subtract subtracts another duration from this one (minimum 0).
func (d Duration) Subtract(other Duration) Duration {
	result := d.seconds - other.seconds
	if result < 0 {
		result = 0
	}
	return NewDuration(result)
}

// String returns a human-readable duration string in format "MM:SS" or "H:MM:SS".
func (d Duration) String() string {
	if d.seconds < 0 {
		return "0:00"
	}

	hours := d.seconds / 3600
	minutes := (d.seconds % 3600) / 60
	seconds := d.seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// ToTime converts the Duration to a time.Duration.
func (d Duration) ToTime() time.Duration {
	return time.Duration(d.seconds) * time.Second
}

// Volume is a value object representing audio volume from 0 to 100.
// It provides validation and ensures volume stays within valid bounds.
type Volume struct {
	level int
}

// NewVolume creates a new Volume with validation (0-100).
func NewVolume(level int) Volume {
	if level < 0 {
		level = 0
	} else if level > 100 {
		level = 100
	}
	return Volume{level: level}
}

// Level returns the volume level (0-100).
func (v Volume) Level() int {
	return v.level
}

// IsValid returns true if the volume is between 0 and 100 inclusive.
func (v Volume) IsValid() bool {
	return v.level >= 0 && v.level <= 100
}

// IsMuted returns true if the volume is 0.
func (v Volume) IsMuted() bool {
	return v.level == 0
}

// IsMax returns true if the volume is at maximum (100).
func (v Volume) IsMax() bool {
	return v.level == 100
}

// Increase increases the volume by the specified amount (capped at 100).
func (v Volume) Increase(amount int) Volume {
	return NewVolume(v.level + amount)
}

// Decrease decreases the volume by the specified amount (minimum 0).
func (v Volume) Decrease(amount int) Volume {
	return NewVolume(v.level - amount)
}

// String returns the volume as a percentage string.
func (v Volume) String() string {
	return strconv.Itoa(v.level) + "%"
}

// Percentage returns the volume as a float between 0.0 and 1.0.
func (v Volume) Percentage() float64 {
	return float64(v.level) / 100.0
}

// PlayerState represents the current state of the music player.
type PlayerState int

const (
	// PlayerStateStopped indicates the player is stopped with no track loaded
	PlayerStateStopped PlayerState = iota

	// PlayerStatePlaying indicates the player is actively playing a track
	PlayerStatePlaying

	// PlayerStatePaused indicates the player is paused with a track loaded
	PlayerStatePaused

	// PlayerStateBuffering indicates the player is loading/buffering content
	PlayerStateBuffering
)

// String returns the string representation of the PlayerState.
func (ps PlayerState) String() string {
	switch ps {
	case PlayerStateStopped:
		return "stopped"
	case PlayerStatePlaying:
		return "playing"
	case PlayerStatePaused:
		return "paused"
	case PlayerStateBuffering:
		return "buffering"
	default:
		return "unknown"
	}
}

// IsValid returns true if the PlayerState is a valid value.
func (ps PlayerState) IsValid() bool {
	return ps >= PlayerStateStopped && ps <= PlayerStateBuffering
}

// RepeatMode represents the repeat behavior of the player.
type RepeatMode int

const (
	// RepeatModeOff indicates no repeat - play through once and stop
	RepeatModeOff RepeatMode = iota

	// RepeatModeAll indicates repeat all tracks in the current context
	RepeatModeAll

	// RepeatModeOne indicates repeat the current track indefinitely
	RepeatModeOne
)

// String returns the string representation of the RepeatMode.
func (rm RepeatMode) String() string {
	switch rm {
	case RepeatModeOff:
		return "off"
	case RepeatModeAll:
		return "all"
	case RepeatModeOne:
		return "one"
	default:
		return "unknown"
	}
}

// IsValid returns true if the RepeatMode is a valid value.
func (rm RepeatMode) IsValid() bool {
	return rm >= RepeatModeOff && rm <= RepeatModeOne
}

// PlaylistType represents the type of playlist.
type PlaylistType int

const (
	// PlaylistTypeUser indicates a user-created playlist
	PlaylistTypeUser PlaylistType = iota

	// PlaylistTypeSmart indicates a smart playlist with dynamic rules
	PlaylistTypeSmart

	// PlaylistTypeLibrary indicates the main music library
	PlaylistTypeLibrary

	// PlaylistTypeQueue indicates the current play queue
	PlaylistTypeQueue

	// PlaylistTypeRecentlyPlayed indicates recently played tracks
	PlaylistTypeRecentlyPlayed

	// PlaylistTypeRecentlyAdded indicates recently added tracks
	PlaylistTypeRecentlyAdded
)

// String returns the string representation of the PlaylistType.
func (pt PlaylistType) String() string {
	switch pt {
	case PlaylistTypeUser:
		return "user"
	case PlaylistTypeSmart:
		return "smart"
	case PlaylistTypeLibrary:
		return "library"
	case PlaylistTypeQueue:
		return "queue"
	case PlaylistTypeRecentlyPlayed:
		return "recently_played"
	case PlaylistTypeRecentlyAdded:
		return "recently_added"
	default:
		return "unknown"
	}
}

// IsValid returns true if the PlaylistType is a valid value.
func (pt PlaylistType) IsValid() bool {
	return pt >= PlaylistTypeUser && pt <= PlaylistTypeRecentlyAdded
}

// IsReadOnly returns true if this playlist type should be read-only.
func (pt PlaylistType) IsReadOnly() bool {
	switch pt {
	case PlaylistTypeLibrary, PlaylistTypeQueue, PlaylistTypeRecentlyPlayed, PlaylistTypeRecentlyAdded:
		return true
	default:
		return false
	}
}
