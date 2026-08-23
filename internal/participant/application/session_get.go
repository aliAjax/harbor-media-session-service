package application

import (
	"context"
	"fmt"

	"harbor-sfu.local/harbor-sfu/internal/participant/domain"
)

func (m *SessionManager) Get(_ context.Context, id string) (*domain.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return s, nil
}
