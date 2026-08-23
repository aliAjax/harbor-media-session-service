package domain

import "time"

func stampRecordingTransition(j *Job, _ Status) { j.UpdatedAt = time.Now() }
