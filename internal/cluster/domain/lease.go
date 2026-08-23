package domain

import (
	"errors"
	"sync"
	"time"
)

var ErrLeaseLost = errors.New("cluster lease lost")

type Lease struct {
	RoomID    string
	Owner     string
	Token     uint64
	ExpiresAt time.Time
}
type Store struct {
	mu     sync.Mutex
	leases map[string]Lease
	seq    uint64
}

func NewStore() *Store { return &Store{leases: map[string]Lease{}} }
func (s *Store) Acquire(room, owner string, ttl time.Duration) (Lease, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if old, ok := s.leases[room]; ok && old.ExpiresAt.After(now) && old.Owner != owner {
		return Lease{}, ErrLeaseLost
	}
	s.seq++
	l := Lease{RoomID: room, Owner: owner, Token: s.seq, ExpiresAt: now.Add(ttl)}
	s.leases[room] = l
	return l, nil
}
func (s *Store) Renew(l Lease, ttl time.Duration) (Lease, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.leases[l.RoomID]
	if !ok || old.Owner != l.Owner || old.Token != l.Token {
		return Lease{}, ErrLeaseLost
	}
	l.ExpiresAt = time.Now().Add(ttl)
	s.leases[l.RoomID] = l
	return l, nil
}
func (s *Store) Release(l Lease) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if old, ok := s.leases[l.RoomID]; ok && old.Token == l.Token {
		delete(s.leases, l.RoomID)
	}
}
