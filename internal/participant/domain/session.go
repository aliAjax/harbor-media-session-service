package domain

import (
	"errors"
	"sync"
	"time"
)

var ErrSessionClosed = errors.New("participant session closed")

type Session struct {
	ID            string
	ParticipantID string
	RoomID        string
	CreatedAt     time.Time
	LastSeen      time.Time
	Closed        bool
	mu            sync.Mutex
}

func (s *Session) Touch() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Closed {
		return ErrSessionClosed
	}
	s.LastSeen = time.Now()
	return nil
}
func (s *Session) Close() { s.mu.Lock(); defer s.mu.Unlock(); s.Closed = true; s.LastSeen = time.Now() }
func (s *Session) Expired(now time.Time, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Closed || now.Sub(s.LastSeen) > ttl
}
