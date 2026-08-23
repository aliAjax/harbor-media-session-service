package domain

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("room resource not found")
var ErrConflict = errors.New("room resource conflict")
var ErrClosed = errors.New("room is closed")

type Status string

const (
	Open     Status = "open"
	Draining Status = "draining"
	Closed   Status = "closed"
)

type Room struct {
	ID           string                  `json:"id"`
	TenantID     string                  `json:"tenant_id"`
	Name         string                  `json:"name"`
	Status       Status                  `json:"status"`
	NodeID       string                  `json:"node_id"`
	FencingToken uint64                  `json:"fencing_token"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	Participants map[string]*Participant `json:"participants"`
	Tracks       map[string]*Track       `json:"tracks"`
	mu           *sync.RWMutex
}
type Participant struct {
	ID        string    `json:"id"`
	Identity  string    `json:"identity"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
	SessionID string    `json:"session_id"`
	Connected bool      `json:"connected"`
}
type Track struct {
	ID          string    `json:"id"`
	Kind        string    `json:"kind"`
	Codec       string    `json:"codec"`
	PublisherID string    `json:"publisher_id"`
	SSRC        uint32    `json:"ssrc"`
	Layers      []Layer   `json:"layers"`
	CreatedAt   time.Time `json:"created_at"`
}
type Layer struct {
	RID     string `json:"rid"`
	Bitrate int64  `json:"bitrate"`
	Active  bool   `json:"active"`
}

func (r *Room) AddParticipant(p *Participant) error {
	if r.mu == nil {
		r.mu = &sync.RWMutex{}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Status != Open {
		return ErrClosed
	}
	if _, ok := r.Participants[p.ID]; ok {
		return ErrConflict
	}
	r.Participants[p.ID] = p
	r.UpdatedAt = time.Now()
	return nil
}
func (r *Room) RemoveParticipant(id string) {
	if r.mu == nil {
		r.mu = &sync.RWMutex{}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Participants, id)
	r.UpdatedAt = time.Now()
}
func (r *Room) AddTrack(t *Track) error {
	if r.mu == nil {
		r.mu = &sync.RWMutex{}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Status == Closed {
		return ErrClosed
	}
	if _, ok := r.Tracks[t.ID]; ok {
		return ErrConflict
	}
	r.Tracks[t.ID] = t
	r.UpdatedAt = time.Now()
	return nil
}
func (r *Room) Snapshot() Room {
	if r.mu == nil {
		r.mu = &sync.RWMutex{}
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := Room{ID: r.ID, TenantID: r.TenantID, Name: r.Name, Status: r.Status, NodeID: r.NodeID, FencingToken: r.FencingToken, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
	out.Participants = map[string]*Participant{}
	for k, v := range r.Participants {
		cp := *v
		out.Participants[k] = &cp
	}
	out.Tracks = map[string]*Track{}
	for k, v := range r.Tracks {
		cp := *v
		out.Tracks[k] = &cp
	}
	return out
}
