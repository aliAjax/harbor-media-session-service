package application

import (
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
