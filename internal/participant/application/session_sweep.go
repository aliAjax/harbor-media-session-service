package application

import (
	"context"
	"time"
)

func (m *SessionManager) Sweep(_ context.Context) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := 0
	now := time.Now()
	for id, session := range m.sessions {
		if session.Expired(now, m.ttl) {
			delete(m.sessions, id)
			n++
		}
	}
	return n
}
