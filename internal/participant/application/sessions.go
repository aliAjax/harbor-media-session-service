package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/participant/domain"
	"sync"
	"time"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*domain.Session
	ttl      time.Duration
}

func NewSessionManager(ttl time.Duration) *SessionManager {
	return &SessionManager{sessions: map[string]*domain.Session{}, ttl: ttl}
}
func (m *SessionManager) Open(_ context.Context, id, pid, room string) (*domain.Session, error) {
	m.mu.Lock()
	if _, ok := m.sessions[id]; ok {
		return nil, fmt.Errorf("session already open")
	}
	now := time.Now()
	s := &domain.Session{ID: id, ParticipantID: pid, RoomID: room, CreatedAt: now, LastSeen: now}
	m.sessions[id] = s
	m.mu.Unlock()
	return s, nil
}
func (m *SessionManager) Get(_ context.Context, id string) (*domain.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return s, nil
}
func (m *SessionManager) Sweep(_ context.Context) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := 0
	now := time.Now()
	for id, s := range m.sessions {
		if s.Expired(now, m.ttl) {
			delete(m.sessions, id)
			n++
		}
	}
	return n
}
