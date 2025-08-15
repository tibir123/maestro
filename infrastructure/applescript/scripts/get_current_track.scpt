-- Get current track information
tell application "Music"
	try
		if player state is stopped then
			return ""
		end if
		
		set theTrack to current track
		set trackID to (database ID of theTrack) as string
		set trackName to name of theTrack
		set trackArtist to artist of theTrack
		set trackAlbum to album of theTrack
		set trackDuration to duration of theTrack
		
		return trackID & "|" & trackName & "|" & trackArtist & "|" & trackAlbum & "|" & trackDuration
	on error errMsg
		error "Failed to get current track: " & errMsg
	end try
end tell