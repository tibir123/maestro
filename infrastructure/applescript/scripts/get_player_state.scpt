-- Get comprehensive player state information
tell application "Music"
	try
		set playerState to player state as string
		set playerVolume to sound volume
		set playerPosition to player position
		set shuffleState to shuffle enabled
		set repeatState to song repeat as string
		
		set currentTrackID to ""
		if player state is not stopped then
			try
				set currentTrackID to (database ID of current track) as string
			end try
		end if
		
		return playerState & "|" & playerVolume & "|" & playerPosition & "|" & shuffleState & "|" & repeatState & "|" & currentTrackID
	on error errMsg
		error "Failed to get player state: " & errMsg
	end try
end tell