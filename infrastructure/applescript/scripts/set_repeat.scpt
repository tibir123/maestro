-- Set repeat mode
-- Parameters: {{repeat_mode}} (off, all, one)
tell application "Music"
	set song repeat to {{repeat_mode}}
	return "ok"
end tell