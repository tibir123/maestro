-- Health check to ensure Music.app is accessible
tell application "Music"
	try
		get name
		return "ok"
	on error errMsg
		error "Music app not accessible: " & errMsg
	end try
end tell