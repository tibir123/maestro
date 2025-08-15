package music

import (
	"errors"
	"testing"
)

func TestNewDomainError(t *testing.T) {
	err := NewDomainError(ErrTrackNotFound, "track not found in library")

	if err.Code != ErrTrackNotFound {
		t.Errorf("expected code %v, got %v", ErrTrackNotFound, err.Code)
	}

	if err.Message != "track not found in library" {
		t.Errorf("expected message 'track not found in library', got '%s'", err.Message)
	}

	if err.Cause != nil {
		t.Errorf("expected no cause, got %v", err.Cause)
	}

	if err.Context == nil {
		t.Error("expected context to be initialized")
	}
}

func TestNewDomainErrorWithCause(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewDomainErrorWithCause(ErrTrackNotFound, "track not found", cause)

	if err.Code != ErrTrackNotFound {
		t.Errorf("expected code %v, got %v", ErrTrackNotFound, err.Code)
	}

	if err.Message != "track not found" {
		t.Errorf("expected message 'track not found', got '%s'", err.Message)
	}

	if err.Cause != cause {
		t.Errorf("expected cause %v, got %v", cause, err.Cause)
	}
}

func TestDomainErrorError(t *testing.T) {
	// Test without cause
	err1 := NewDomainError(ErrTrackNotFound, "track not found")
	expected1 := "track not found: track not found"
	if err1.Error() != expected1 {
		t.Errorf("expected error string '%s', got '%s'", expected1, err1.Error())
	}

	// Test with cause
	cause := errors.New("underlying error")
	err2 := NewDomainErrorWithCause(ErrTrackNotFound, "track not found", cause)
	expected2 := "track not found: track not found (caused by: underlying error)"
	if err2.Error() != expected2 {
		t.Errorf("expected error string '%s', got '%s'", expected2, err2.Error())
	}
}

func TestDomainErrorUnwrap(t *testing.T) {
	// Test without cause
	err1 := NewDomainError(ErrTrackNotFound, "track not found")
	if err1.Unwrap() != nil {
		t.Error("expected Unwrap() to return nil when no cause")
	}

	// Test with cause
	cause := errors.New("underlying error")
	err2 := NewDomainErrorWithCause(ErrTrackNotFound, "track not found", cause)
	if err2.Unwrap() != cause {
		t.Errorf("expected Unwrap() to return %v, got %v", cause, err2.Unwrap())
	}
}

func TestDomainErrorIs(t *testing.T) {
	err := NewDomainError(ErrTrackNotFound, "track not found")

	if !err.Is(ErrTrackNotFound) {
		t.Error("expected Is() to return true for matching error")
	}

	if err.Is(ErrPlaylistNotFound) {
		t.Error("expected Is() to return false for non-matching error")
	}
}

func TestDomainErrorContext(t *testing.T) {
	err := NewDomainError(ErrTrackNotFound, "track not found")

	// Test adding context
	err = err.WithContext("track_id", "track-123")
	err = err.WithContext("user_id", "user-456")

	// Test retrieving context
	trackID, exists := err.GetContext("track_id")
	if !exists {
		t.Error("expected track_id to exist in context")
	}
	if trackID != "track-123" {
		t.Errorf("expected track_id 'track-123', got '%v'", trackID)
	}

	userID, exists := err.GetContext("user_id")
	if !exists {
		t.Error("expected user_id to exist in context")
	}
	if userID != "user-456" {
		t.Errorf("expected user_id 'user-456', got '%v'", userID)
	}

	// Test non-existent key
	_, exists = err.GetContext("non_existent")
	if exists {
		t.Error("expected non_existent key to not exist")
	}
}

func TestDomainErrorClassification(t *testing.T) {
	tests := []struct {
		name     string
		code     error
		isTrack  bool
		isPlayer bool
		isList   bool
		isQueue  bool
		isLib    bool
	}{
		{"track error", ErrTrackNotFound, true, false, false, false, false},
		{"invalid track ID", ErrInvalidTrackID, true, false, false, false, false},
		{"playlist error", ErrPlaylistNotFound, false, false, true, false, false},
		{"playlist read-only", ErrPlaylistReadOnly, false, false, true, false, false},
		{"player error", ErrPlayerNotAvailable, false, true, false, false, false},
		{"invalid volume", ErrInvalidVolume, false, true, false, false, false},
		{"queue error", ErrQueueEmpty, false, false, false, true, false},
		{"library error", ErrLibraryNotAvailable, false, false, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewDomainError(tt.code, "test message")

			if err.IsTrackError() != tt.isTrack {
				t.Errorf("expected IsTrackError() %v, got %v", tt.isTrack, err.IsTrackError())
			}

			if err.IsPlayerError() != tt.isPlayer {
				t.Errorf("expected IsPlayerError() %v, got %v", tt.isPlayer, err.IsPlayerError())
			}

			if err.IsPlaylistError() != tt.isList {
				t.Errorf("expected IsPlaylistError() %v, got %v", tt.isList, err.IsPlaylistError())
			}

			if err.IsQueueError() != tt.isQueue {
				t.Errorf("expected IsQueueError() %v, got %v", tt.isQueue, err.IsQueueError())
			}

			if err.IsLibraryError() != tt.isLib {
				t.Errorf("expected IsLibraryError() %v, got %v", tt.isLib, err.IsLibraryError())
			}
		})
	}
}

func TestDomainErrorRetryable(t *testing.T) {
	tests := []struct {
		name      string
		code      error
		retryable bool
		permanent bool
	}{
		{"timeout", ErrTimeout, true, false},
		{"player not available", ErrPlayerNotAvailable, true, false},
		{"library not available", ErrLibraryNotAvailable, true, false},
		{"operation failed", ErrOperationFailed, true, false},
		{"permission denied", ErrPermissionDenied, false, true},
		{"invalid track ID", ErrInvalidTrackID, false, true},
		{"invalid volume", ErrInvalidVolume, false, true},
		{"playlist read-only", ErrPlaylistReadOnly, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewDomainError(tt.code, "test message")

			if err.IsRetryable() != tt.retryable {
				t.Errorf("expected IsRetryable() %v, got %v", tt.retryable, err.IsRetryable())
			}

			if err.IsPermanent() != tt.permanent {
				t.Errorf("expected IsPermanent() %v, got %v", tt.permanent, err.IsPermanent())
			}
		})
	}
}

func TestWrapTrackNotFound(t *testing.T) {
	trackID := NewTrackID("track-123")
	cause := errors.New("underlying error")

	err := WrapTrackNotFound(trackID, cause)

	if !err.Is(ErrTrackNotFound) {
		t.Error("expected wrapped error to be ErrTrackNotFound")
	}

	if err.Cause != cause {
		t.Errorf("expected cause %v, got %v", cause, err.Cause)
	}

	contextTrackID, exists := err.GetContext("track_id")
	if !exists {
		t.Error("expected track_id in context")
	}
	if contextTrackID != trackID.Value() {
		t.Errorf("expected track_id '%s', got '%v'", trackID.Value(), contextTrackID)
	}
}

func TestWrapPlaylistNotFound(t *testing.T) {
	playlistID := NewPlaylistID("playlist-123")
	cause := errors.New("underlying error")

	err := WrapPlaylistNotFound(playlistID, cause)

	if !err.Is(ErrPlaylistNotFound) {
		t.Error("expected wrapped error to be ErrPlaylistNotFound")
	}

	if err.Cause != cause {
		t.Errorf("expected cause %v, got %v", cause, err.Cause)
	}

	contextPlaylistID, exists := err.GetContext("playlist_id")
	if !exists {
		t.Error("expected playlist_id in context")
	}
	if contextPlaylistID != playlistID.Value() {
		t.Errorf("expected playlist_id '%s', got '%v'", playlistID.Value(), contextPlaylistID)
	}
}

func TestWrapInvalidVolume(t *testing.T) {
	volume := 150
	cause := errors.New("underlying error")

	err := WrapInvalidVolume(volume, cause)

	if !err.Is(ErrInvalidVolume) {
		t.Error("expected wrapped error to be ErrInvalidVolume")
	}

	if err.Cause != cause {
		t.Errorf("expected cause %v, got %v", cause, err.Cause)
	}

	contextVolume, exists := err.GetContext("volume")
	if !exists {
		t.Error("expected volume in context")
	}
	if contextVolume != volume {
		t.Errorf("expected volume %d, got %v", volume, contextVolume)
	}
}

func TestWrapInvalidPosition(t *testing.T) {
	position := NewDuration(300)
	cause := errors.New("underlying error")

	err := WrapInvalidPosition(position, cause)

	if !err.Is(ErrInvalidPosition) {
		t.Error("expected wrapped error to be ErrInvalidPosition")
	}

	if err.Cause != cause {
		t.Errorf("expected cause %v, got %v", cause, err.Cause)
	}

	contextPosition, exists := err.GetContext("position_seconds")
	if !exists {
		t.Error("expected position_seconds in context")
	}
	if contextPosition != position.Seconds() {
		t.Errorf("expected position_seconds %d, got %v", position.Seconds(), contextPosition)
	}
}

func TestWrapOperationTimeout(t *testing.T) {
	operation := "play_track"
	timeout := NewDuration(5)
	cause := errors.New("underlying error")

	err := WrapOperationTimeout(operation, timeout, cause)

	if !err.Is(ErrTimeout) {
		t.Error("expected wrapped error to be ErrTimeout")
	}

	if err.Cause != cause {
		t.Errorf("expected cause %v, got %v", cause, err.Cause)
	}

	contextOperation, exists := err.GetContext("operation")
	if !exists {
		t.Error("expected operation in context")
	}
	if contextOperation != operation {
		t.Errorf("expected operation '%s', got '%v'", operation, contextOperation)
	}

	contextTimeout, exists := err.GetContext("timeout_seconds")
	if !exists {
		t.Error("expected timeout_seconds in context")
	}
	if contextTimeout != timeout.Seconds() {
		t.Errorf("expected timeout_seconds %d, got %v", timeout.Seconds(), contextTimeout)
	}
}

func TestErrorCheckingUtilities(t *testing.T) {
	// Test IsTrackNotFound
	trackErr := NewDomainError(ErrTrackNotFound, "track not found")
	if !IsTrackNotFound(trackErr) {
		t.Error("expected IsTrackNotFound to return true")
	}

	otherErr := NewDomainError(ErrPlaylistNotFound, "playlist not found")
	if IsTrackNotFound(otherErr) {
		t.Error("expected IsTrackNotFound to return false for non-track error")
	}

	// Test IsPlaylistNotFound
	if !IsPlaylistNotFound(otherErr) {
		t.Error("expected IsPlaylistNotFound to return true")
	}

	if IsPlaylistNotFound(trackErr) {
		t.Error("expected IsPlaylistNotFound to return false for non-playlist error")
	}

	// Test IsInvalidVolume
	volumeErr := NewDomainError(ErrInvalidVolume, "invalid volume")
	if !IsInvalidVolume(volumeErr) {
		t.Error("expected IsInvalidVolume to return true")
	}

	if IsInvalidVolume(trackErr) {
		t.Error("expected IsInvalidVolume to return false for non-volume error")
	}

	// Test IsTimeout
	timeoutErr := NewDomainError(ErrTimeout, "operation timed out")
	if !IsTimeout(timeoutErr) {
		t.Error("expected IsTimeout to return true")
	}

	if IsTimeout(trackErr) {
		t.Error("expected IsTimeout to return false for non-timeout error")
	}
}

func TestUtilityRetryableAndPermanent(t *testing.T) {
	// Test with DomainError
	retryableErr := NewDomainError(ErrTimeout, "timeout")
	if !IsRetryable(retryableErr) {
		t.Error("expected timeout error to be retryable")
	}

	permanentErr := NewDomainError(ErrInvalidTrackID, "invalid ID")
	if !IsPermanent(permanentErr) {
		t.Error("expected invalid ID error to be permanent")
	}

	// Test with non-domain errors
	timeoutErr := ErrTimeout
	if !IsRetryable(timeoutErr) {
		t.Error("expected base timeout error to be retryable")
	}

	invalidIDErr := ErrInvalidTrackID
	if !IsPermanent(invalidIDErr) {
		t.Error("expected base invalid ID error to be permanent")
	}

	// Test with generic error
	genericErr := errors.New("generic error")
	if IsRetryable(genericErr) {
		t.Error("expected generic error to not be retryable by default")
	}

	if IsPermanent(genericErr) {
		t.Error("expected generic error to not be permanent by default")
	}
}
