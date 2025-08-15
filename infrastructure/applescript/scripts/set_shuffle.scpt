-- Set shuffle mode
-- Parameters: {{shuffle_enabled}}
tell application "Music"
	set shuffle enabled to {{shuffle_enabled}}
	return "ok"
end tell