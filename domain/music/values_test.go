package music

import (
	"testing"
	"time"
)

func TestTrackID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal ID", "track-123", "track-123"},
		{"with spaces", "  track-456  ", "track-456"},
		{"empty string", "", ""},
		{"only spaces", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := NewTrackID(tt.input)
			if id.Value() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, id.Value())
			}
			if id.String() != tt.expected {
				t.Errorf("expected string %s, got %s", tt.expected, id.String())
			}
		})
	}
}

func TestTrackIDIsEmpty(t *testing.T) {
	emptyID := NewTrackID("")
	nonEmptyID := NewTrackID("track-123")

	if !emptyID.IsEmpty() {
		t.Error("empty ID should return true for IsEmpty()")
	}

	if nonEmptyID.IsEmpty() {
		t.Error("non-empty ID should return false for IsEmpty()")
	}
}

func TestTrackIDEquals(t *testing.T) {
	id1 := NewTrackID("track-123")
	id2 := NewTrackID("track-123")
	id3 := NewTrackID("track-456")

	if !id1.Equals(id2) {
		t.Error("identical IDs should be equal")
	}

	if id1.Equals(id3) {
		t.Error("different IDs should not be equal")
	}
}

func TestPlaylistID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal ID", "playlist-123", "playlist-123"},
		{"with spaces", "  playlist-456  ", "playlist-456"},
		{"empty string", "", ""},
		{"only spaces", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := NewPlaylistID(tt.input)
			if id.Value() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, id.Value())
			}
			if id.String() != tt.expected {
				t.Errorf("expected string %s, got %s", tt.expected, id.String())
			}
		})
	}
}

func TestPlaylistIDIsEmpty(t *testing.T) {
	emptyID := NewPlaylistID("")
	nonEmptyID := NewPlaylistID("playlist-123")

	if !emptyID.IsEmpty() {
		t.Error("empty ID should return true for IsEmpty()")
	}

	if nonEmptyID.IsEmpty() {
		t.Error("non-empty ID should return false for IsEmpty()")
	}
}

func TestPlaylistIDEquals(t *testing.T) {
	id1 := NewPlaylistID("playlist-123")
	id2 := NewPlaylistID("playlist-123")
	id3 := NewPlaylistID("playlist-456")

	if !id1.Equals(id2) {
		t.Error("identical IDs should be equal")
	}

	if id1.Equals(id3) {
		t.Error("different IDs should not be equal")
	}
}

func TestNewDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive duration", 180, 180},
		{"zero duration", 0, 0},
		{"negative duration", -60, 0}, // should be clamped to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := NewDuration(tt.input)
			if duration.Seconds() != tt.expected {
				t.Errorf("expected %d seconds, got %d", tt.expected, duration.Seconds())
			}
		})
	}
}

func TestNewDurationFromTime(t *testing.T) {
	timeDuration := 3*time.Minute + 30*time.Second // 210 seconds
	duration := NewDurationFromTime(timeDuration)

	if duration.Seconds() != 210 {
		t.Errorf("expected 210 seconds, got %d", duration.Seconds())
	}
}

func TestDurationMethods(t *testing.T) {
	duration := NewDuration(3661) // 1 hour, 1 minute, 1 second

	if duration.Minutes() != 61 {
		t.Errorf("expected 61 minutes, got %d", duration.Minutes())
	}

	if duration.Hours() != 1 {
		t.Errorf("expected 1 hour, got %d", duration.Hours())
	}

	if !duration.IsValid() {
		t.Error("positive duration should be valid")
	}

	if duration.IsZero() {
		t.Error("non-zero duration should not be zero")
	}
}

func TestDurationArithmetic(t *testing.T) {
	duration1 := NewDuration(60)  // 1 minute
	duration2 := NewDuration(120) // 2 minutes

	// Test Add
	sum := duration1.Add(duration2)
	if sum.Seconds() != 180 {
		t.Errorf("expected sum to be 180 seconds, got %d", sum.Seconds())
	}

	// Test Subtract
	diff := duration2.Subtract(duration1)
	if diff.Seconds() != 60 {
		t.Errorf("expected difference to be 60 seconds, got %d", diff.Seconds())
	}

	// Test Subtract with negative result (should clamp to 0)
	diff2 := duration1.Subtract(duration2)
	if diff2.Seconds() != 0 {
		t.Errorf("expected negative difference to be clamped to 0, got %d", diff2.Seconds())
	}
}

func TestDurationString(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"under a minute", 45, "0:45"},
		{"exactly one minute", 60, "1:00"},
		{"over a minute", 90, "1:30"},
		{"under an hour", 3599, "59:59"},
		{"exactly one hour", 3600, "1:00:00"},
		{"over an hour", 3661, "1:01:01"},
		{"zero duration", 0, "0:00"},
		{"negative duration", -1, "0:00"}, // NewDuration clamps to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := NewDuration(tt.seconds)
			if duration.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, duration.String())
			}
		})
	}
}

func TestDurationToTime(t *testing.T) {
	duration := NewDuration(180)
	timeDuration := duration.ToTime()

	if timeDuration != 180*time.Second {
		t.Errorf("expected 180 seconds, got %v", timeDuration)
	}
}

func TestNewVolume(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"valid volume", 50, 50},
		{"minimum volume", 0, 0},
		{"maximum volume", 100, 100},
		{"below minimum", -10, 0},   // should be clamped to 0
		{"above maximum", 150, 100}, // should be clamped to 100
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volume := NewVolume(tt.input)
			if volume.Level() != tt.expected {
				t.Errorf("expected level %d, got %d", tt.expected, volume.Level())
			}
		})
	}
}

func TestVolumeValidation(t *testing.T) {
	validVolume := NewVolume(50)
	if !validVolume.IsValid() {
		t.Error("valid volume should return true for IsValid()")
	}

	// Test boundaries
	minVolume := NewVolume(0)
	maxVolume := NewVolume(100)

	if !minVolume.IsValid() {
		t.Error("minimum volume should be valid")
	}

	if !maxVolume.IsValid() {
		t.Error("maximum volume should be valid")
	}
}

func TestVolumeState(t *testing.T) {
	mutedVolume := NewVolume(0)
	if !mutedVolume.IsMuted() {
		t.Error("volume 0 should be muted")
	}

	maxVolume := NewVolume(100)
	if !maxVolume.IsMax() {
		t.Error("volume 100 should be max")
	}

	normalVolume := NewVolume(50)
	if normalVolume.IsMuted() || normalVolume.IsMax() {
		t.Error("volume 50 should not be muted or max")
	}
}

func TestVolumeAdjustment(t *testing.T) {
	volume := NewVolume(50)

	// Test increase
	increased := volume.Increase(25)
	if increased.Level() != 75 {
		t.Errorf("expected level 75 after increase, got %d", increased.Level())
	}

	// Test increase beyond maximum
	maxed := volume.Increase(75)
	if maxed.Level() != 100 {
		t.Errorf("expected level 100 after large increase, got %d", maxed.Level())
	}

	// Test decrease
	decreased := volume.Decrease(25)
	if decreased.Level() != 25 {
		t.Errorf("expected level 25 after decrease, got %d", decreased.Level())
	}

	// Test decrease below minimum
	muted := volume.Decrease(75)
	if muted.Level() != 0 {
		t.Errorf("expected level 0 after large decrease, got %d", muted.Level())
	}
}

func TestVolumeString(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{0, "0%"},
		{50, "50%"},
		{100, "100%"},
	}

	for _, tt := range tests {
		volume := NewVolume(tt.level)
		if volume.String() != tt.expected {
			t.Errorf("expected string %s, got %s", tt.expected, volume.String())
		}
	}
}

func TestVolumePercentage(t *testing.T) {
	tests := []struct {
		level    int
		expected float64
	}{
		{0, 0.0},
		{50, 0.5},
		{100, 1.0},
	}

	for _, tt := range tests {
		volume := NewVolume(tt.level)
		if volume.Percentage() != tt.expected {
			t.Errorf("expected percentage %f, got %f", tt.expected, volume.Percentage())
		}
	}
}

func TestPlayerStateEnum(t *testing.T) {
	tests := []struct {
		state    PlayerState
		expected string
		valid    bool
	}{
		{PlayerStateStopped, "stopped", true},
		{PlayerStatePlaying, "playing", true},
		{PlayerStatePaused, "paused", true},
		{PlayerStateBuffering, "buffering", true},
		{PlayerState(999), "unknown", false},
	}

	for _, tt := range tests {
		if tt.state.String() != tt.expected {
			t.Errorf("expected string %s for state %v, got %s", tt.expected, tt.state, tt.state.String())
		}

		if tt.state.IsValid() != tt.valid {
			t.Errorf("expected IsValid() %v for state %v, got %v", tt.valid, tt.state, tt.state.IsValid())
		}
	}
}

func TestRepeatModeEnum(t *testing.T) {
	tests := []struct {
		mode     RepeatMode
		expected string
		valid    bool
	}{
		{RepeatModeOff, "off", true},
		{RepeatModeAll, "all", true},
		{RepeatModeOne, "one", true},
		{RepeatMode(999), "unknown", false},
	}

	for _, tt := range tests {
		if tt.mode.String() != tt.expected {
			t.Errorf("expected string %s for mode %v, got %s", tt.expected, tt.mode, tt.mode.String())
		}

		if tt.mode.IsValid() != tt.valid {
			t.Errorf("expected IsValid() %v for mode %v, got %v", tt.valid, tt.mode, tt.mode.IsValid())
		}
	}
}

func TestPlaylistTypeEnum(t *testing.T) {
	tests := []struct {
		pType    PlaylistType
		expected string
		valid    bool
		readOnly bool
	}{
		{PlaylistTypeUser, "user", true, false},
		{PlaylistTypeSmart, "smart", true, false},
		{PlaylistTypeLibrary, "library", true, true},
		{PlaylistTypeQueue, "queue", true, true},
		{PlaylistTypeRecentlyPlayed, "recently_played", true, true},
		{PlaylistTypeRecentlyAdded, "recently_added", true, true},
		{PlaylistType(999), "unknown", false, false},
	}

	for _, tt := range tests {
		if tt.pType.String() != tt.expected {
			t.Errorf("expected string %s for type %v, got %s", tt.expected, tt.pType, tt.pType.String())
		}

		if tt.pType.IsValid() != tt.valid {
			t.Errorf("expected IsValid() %v for type %v, got %v", tt.valid, tt.pType, tt.pType.IsValid())
		}

		if tt.pType.IsReadOnly() != tt.readOnly {
			t.Errorf("expected IsReadOnly() %v for type %v, got %v", tt.readOnly, tt.pType, tt.pType.IsReadOnly())
		}
	}
}
