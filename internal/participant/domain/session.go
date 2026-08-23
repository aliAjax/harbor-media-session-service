package domain

import (
	"errors"
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
}

func (s *Session) Touch() error { if s.Closed { return ErrSessionClosed }; s.LastSeen = time.Now(); return nil }
func (s *Session) Close() { s.Closed = true; s.LastSeen = time.Now() }
func (s *Session) Expired(now time.Time, ttl time.Duration) bool { return s.Closed || now.Sub(s.LastSeen) > ttl }
