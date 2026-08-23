package domain

import (
	"errors"
	"time"
)

type Status string

const (
	Starting  Status = "starting"
	Running   Status = "running"
	Stopping  Status = "stopping"
	Completed Status = "completed"
	Failed    Status = "failed"
)

var ErrInvalid = errors.New("invalid recording transition")

type Job struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	Status    Status    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Error     string    `json:"error,omitempty"`
}

func (j *Job) Transition(s Status) error {
	valid := map[Status][]Status{Starting: {Running, Failed}, Running: {Stopping, Failed}, Stopping: {Completed, Failed}}
	for _, x := range valid[j.Status] {
		if x == s {
			j.Status = s
			j.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrInvalid
}
