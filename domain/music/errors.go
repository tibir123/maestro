package music

import (
	"errors"
	"fmt"
)

// Domain error types - these define the categories of errors that can occur
// in the music domain. They are used as base errors that can be wrapped
// with additional context.
var (
	// Track-related errors
	ErrTrackNotFound  = errors.New("track not found")
	ErrInvalidTrackID = errors.New("invalid track ID")
	ErrInvalidTrack   = errors.New("invalid track")

	// Playlist-related errors
	ErrPlaylistNotFound       = errors.New("playlist not found")
	ErrInvalidPlaylistID      = errors.New("invalid playlist ID")
	ErrInvalidPlaylist        = errors.New("invalid playlist")
	ErrPlaylistReadOnly       = errors.New("playlist is read-only")
	ErrTrackAlreadyInPlaylist = errors.New("track already in playlist")

	// Player-related errors
	ErrPlayerNotAvailable = errors.New("player not available")
	ErrInvalidPlayerState = errors.New("invalid player state")
	ErrInvalidVolume      = errors.New("invalid volume")
	ErrInvalidPosition    = errors.New("invalid position")
	ErrInvalidRepeatMode  = errors.New("invalid repeat mode")

	// Queue-related errors
	ErrQueueEmpty           = errors.New("queue is empty")
	ErrInvalidQueuePosition = errors.New("invalid queue position")

	// Library-related errors
	ErrLibraryNotAvailable = errors.New("library not available")
	ErrSearchFailed        = errors.New("search failed")
	ErrInvalidSearchQuery  = errors.New("invalid search query")

	// General errors
	ErrOperationFailed  = errors.New("operation failed")
	ErrTimeout          = errors.New("operation timed out")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidOperation = errors.New("invalid operation")
)

// DomainError represents an error that occurred in the music domain.
// It wraps underlying errors with additional context and provides
// methods for error classification and handling.
type DomainError struct {
	// Code is the specific error type (one of the Err* variables above)
	Code error

	// Message provides human-readable context about the error
	Message string

	// Cause is the underlying error that caused this domain error (optional)
	Cause error

	// Context provides additional key-value pairs for debugging (optional)
	Context map[string]interface{}
}

// NewDomainError creates a new DomainError with the specified code and message.
func NewDomainError(code error, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// NewDomainErrorWithCause creates a new DomainError wrapping an underlying error.
func NewDomainErrorWithCause(code error, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// Error implements the error interface.
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code.Error(), e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code.Error(), e.Message)
}

// Unwrap returns the underlying cause error for error unwrapping.
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Is checks if this error matches the target error.
func (e *DomainError) Is(target error) bool {
	return errors.Is(e.Code, target)
}

// WithContext adds context information to the error.
func (e *DomainError) WithContext(key string, value interface{}) *DomainError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// GetContext retrieves context information from the error.
func (e *DomainError) GetContext(key string) (interface{}, bool) {
	if e.Context == nil {
		return nil, false
	}
	value, exists := e.Context[key]
	return value, exists
}

// IsTrackError returns true if this is a track-related error.
func (e *DomainError) IsTrackError() bool {
	return errors.Is(e.Code, ErrTrackNotFound) ||
		errors.Is(e.Code, ErrInvalidTrackID) ||
		errors.Is(e.Code, ErrInvalidTrack)
}

// IsPlaylistError returns true if this is a playlist-related error.
func (e *DomainError) IsPlaylistError() bool {
	return errors.Is(e.Code, ErrPlaylistNotFound) ||
		errors.Is(e.Code, ErrInvalidPlaylistID) ||
		errors.Is(e.Code, ErrInvalidPlaylist) ||
		errors.Is(e.Code, ErrPlaylistReadOnly) ||
		errors.Is(e.Code, ErrTrackAlreadyInPlaylist)
}

// IsPlayerError returns true if this is a player-related error.
func (e *DomainError) IsPlayerError() bool {
	return errors.Is(e.Code, ErrPlayerNotAvailable) ||
		errors.Is(e.Code, ErrInvalidPlayerState) ||
		errors.Is(e.Code, ErrInvalidVolume) ||
		errors.Is(e.Code, ErrInvalidPosition) ||
		errors.Is(e.Code, ErrInvalidRepeatMode)
}

// IsQueueError returns true if this is a queue-related error.
func (e *DomainError) IsQueueError() bool {
	return errors.Is(e.Code, ErrQueueEmpty) ||
		errors.Is(e.Code, ErrInvalidQueuePosition)
}

// IsLibraryError returns true if this is a library-related error.
func (e *DomainError) IsLibraryError() bool {
	return errors.Is(e.Code, ErrLibraryNotAvailable) ||
		errors.Is(e.Code, ErrSearchFailed) ||
		errors.Is(e.Code, ErrInvalidSearchQuery)
}

// IsRetryable returns true if this error might succeed if retried.
func (e *DomainError) IsRetryable() bool {
	return errors.Is(e.Code, ErrTimeout) ||
		errors.Is(e.Code, ErrPlayerNotAvailable) ||
		errors.Is(e.Code, ErrLibraryNotAvailable) ||
		errors.Is(e.Code, ErrOperationFailed)
}

// IsPermanent returns true if this error is unlikely to succeed if retried.
func (e *DomainError) IsPermanent() bool {
	return errors.Is(e.Code, ErrPermissionDenied) ||
		errors.Is(e.Code, ErrInvalidTrackID) ||
		errors.Is(e.Code, ErrInvalidPlaylistID) ||
		errors.Is(e.Code, ErrInvalidVolume) ||
		errors.Is(e.Code, ErrInvalidPosition) ||
		errors.Is(e.Code, ErrInvalidSearchQuery) ||
		errors.Is(e.Code, ErrPlaylistReadOnly)
}

// Error wrapping utilities

// WrapTrackNotFound wraps a track not found error with additional context.
func WrapTrackNotFound(trackID TrackID, cause error) *DomainError {
	err := NewDomainErrorWithCause(
		ErrTrackNotFound,
		fmt.Sprintf("track with ID '%s' was not found", trackID.Value()),
		cause,
	)
	return err.WithContext("track_id", trackID.Value())
}

// WrapPlaylistNotFound wraps a playlist not found error with additional context.
func WrapPlaylistNotFound(playlistID PlaylistID, cause error) *DomainError {
	err := NewDomainErrorWithCause(
		ErrPlaylistNotFound,
		fmt.Sprintf("playlist with ID '%s' was not found", playlistID.Value()),
		cause,
	)
	return err.WithContext("playlist_id", playlistID.Value())
}

// WrapInvalidVolume wraps an invalid volume error with the attempted value.
func WrapInvalidVolume(volume int, cause error) *DomainError {
	err := NewDomainErrorWithCause(
		ErrInvalidVolume,
		fmt.Sprintf("volume %d is invalid (must be 0-100)", volume),
		cause,
	)
	return err.WithContext("volume", volume)
}

// WrapInvalidPosition wraps an invalid position error with the attempted value.
func WrapInvalidPosition(position Duration, cause error) *DomainError {
	err := NewDomainErrorWithCause(
		ErrInvalidPosition,
		fmt.Sprintf("position %s is invalid", position.String()),
		cause,
	)
	return err.WithContext("position_seconds", position.Seconds())
}

// WrapOperationTimeout wraps a timeout error with operation context.
func WrapOperationTimeout(operation string, timeout Duration, cause error) *DomainError {
	err := NewDomainErrorWithCause(
		ErrTimeout,
		fmt.Sprintf("operation '%s' timed out after %s", operation, timeout.String()),
		cause,
	)
	return err.WithContext("operation", operation).WithContext("timeout_seconds", timeout.Seconds())
}

// Error checking utilities

// IsTrackNotFound checks if an error is a track not found error.
func IsTrackNotFound(err error) bool {
	return errors.Is(err, ErrTrackNotFound)
}

// IsPlaylistNotFound checks if an error is a playlist not found error.
func IsPlaylistNotFound(err error) bool {
	return errors.Is(err, ErrPlaylistNotFound)
}

// IsInvalidVolume checks if an error is an invalid volume error.
func IsInvalidVolume(err error) bool {
	return errors.Is(err, ErrInvalidVolume)
}

// IsTimeout checks if an error is a timeout error.
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}

// IsRetryable checks if an error might succeed if retried.
func IsRetryable(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.IsRetryable()
	}

	// Default retry logic for non-domain errors
	return errors.Is(err, ErrTimeout) ||
		errors.Is(err, ErrPlayerNotAvailable) ||
		errors.Is(err, ErrLibraryNotAvailable)
}

// IsPermanent checks if an error is unlikely to succeed if retried.
func IsPermanent(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.IsPermanent()
	}

	// Default permanent error logic for non-domain errors
	return errors.Is(err, ErrPermissionDenied) ||
		errors.Is(err, ErrInvalidTrackID) ||
		errors.Is(err, ErrInvalidPlaylistID)
}
