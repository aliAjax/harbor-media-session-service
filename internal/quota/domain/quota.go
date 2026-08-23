package domain

import (
	"fmt"
	"sync"
)

type Limits struct {
	Rooms         int
	Participants  int
	Tracks        int
	Subscriptions int
	Bandwidth     int64
}
type Usage struct {
	Rooms         int
	Participants  int
	Tracks        int
	Subscriptions int
	Bandwidth     int64
}
type Manager struct {
	mu     sync.Mutex
	limits map[string]Limits
	usage  map[string]Usage
}

func NewManager() *Manager { return &Manager{limits: map[string]Limits{}, usage: map[string]Usage{}} }
func (m *Manager) Set(tenant string, l Limits) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limits[tenant] = l
}
func (m *Manager) Reserve(tenant string, delta Usage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l := m.limits[tenant]
	u := m.usage[tenant]
	n := Usage{u.Rooms + delta.Rooms, u.Participants + delta.Participants, u.Tracks + delta.Tracks, u.Subscriptions + delta.Subscriptions, u.Bandwidth + delta.Bandwidth}
	m.usage[tenant] = n
	if l.Rooms > 0 && n.Rooms > l.Rooms {
		return fmt.Errorf("room quota exceeded")
	}
	if l.Participants > 0 && n.Participants > l.Participants {
		return fmt.Errorf("participant quota exceeded")
	}
	if l.Tracks > 0 && n.Tracks > l.Tracks {
		return fmt.Errorf("track quota exceeded")
	}
	if l.Subscriptions > 0 && n.Subscriptions > l.Subscriptions {
		return fmt.Errorf("subscription quota exceeded")
	}
	if l.Bandwidth > 0 && n.Bandwidth > l.Bandwidth {
		return fmt.Errorf("bandwidth quota exceeded")
	}
	m.usage[tenant] = n
	return nil
}
func (m *Manager) Release(tenant string, delta Usage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u := m.usage[tenant]
	u.Rooms -= delta.Rooms
	u.Participants -= delta.Participants
	u.Tracks -= delta.Tracks
	u.Subscriptions -= delta.Subscriptions
	u.Bandwidth -= delta.Bandwidth
	if u.Rooms < 0 {
		u.Rooms = 0
	}
	m.usage[tenant] = u
}
func (m *Manager) Snapshot() map[string]Usage {
	m.mu.Lock()
	defer m.mu.Unlock()
	o := map[string]Usage{}
	for k, v := range m.usage {
		o[k] = v
	}
	return o
}
