-- Seek to a specific position in the current track
-- Parameters: {{position_seconds}}
tell application "Music"
	try
		set player position to {{position_seconds}}
		return "ok"
	on error errMsg
		error "Failed to seek: " & errMsg
	end try
end tell