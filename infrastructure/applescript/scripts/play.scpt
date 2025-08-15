-- Play a specific track by ID
-- Parameters: {{track_id}}
tell application "Music"
	try
		set theTrack to track id {{track_id}}
		play theTrack
		return "ok"
	on error errMsg
		error "Failed to play track: " & errMsg
	end try
end tell