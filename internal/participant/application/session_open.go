package application

import (
	"context"
	"fmt"
	"time"

	"harbor-sfu.local/harbor-sfu/internal/participant/domain"
)

func (m *SessionManager) Open(_ context.Context, id, pid, room string) (*domain.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.sessions[id]; ok {
		return nil, fmt.Errorf("session already open")
	}
	now := time.Now()
	s := &domain.Session{ID: id, ParticipantID: pid, RoomID: room, CreatedAt: now, LastSeen: now}
	m.sessions[id] = s
	return s, nil
}
