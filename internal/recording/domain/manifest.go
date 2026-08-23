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
	return nil
}
func (m Manifest) JSON() []byte { b, _ := json.Marshal(m); return b }

// Clone returns a deep copy of the manifest so the caller and the registry
// never share the Tracks backing array. A nil Tracks slice stays nil.
func (m Manifest) Clone() Manifest {
	cp := m
	if m.Tracks != nil {
		tracks := make([]string, len(m.Tracks))
		copy(tracks, m.Tracks)
		cp.Tracks = tracks
	}
	return cp
}
