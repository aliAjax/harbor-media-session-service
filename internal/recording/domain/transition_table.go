package domain

func canAdvanceRecording(from, to Status) bool {
	valid := map[Status][]Status{Starting: {Running, Failed}, Running: {Stopping, Failed}, Stopping: {Completed, Failed, Running}, Completed: {Running}}
	for _, candidate := range valid[from] {
		if candidate == to {
			return true
		}
	}
	return false
}
