package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Manifest struct {
	JobID       string     `json:"job_id"`
	RoomID      string     `json:"room_id"`
	Tracks      []string   `json:"tracks"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Packets     uint64     `json:"packets"`
}

func (m Manifest) Validate() error {
	if m.JobID == "" || m.RoomID == "" {
		return fmt.Errorf("manifest identifiers required")
	}
	if len(m.Tracks) == 0 {
		return fmt.Errorf("manifest has no tracks")
	}
	if m.Tracks[0] == "" {
		return fmt.Errorf("manifest first track required")
	}
	return nil
}
func (m Manifest) JSON() []byte { b, _ := json.Marshal(m); return b }
