-- Set the playback volume
-- Parameters: {{volume_level}}
tell application "Music"
	set sound volume to {{volume_level}}
	return "ok"
end tell