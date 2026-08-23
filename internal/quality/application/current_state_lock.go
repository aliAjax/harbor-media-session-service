package application

import "sync"

type CurrentState struct {
	mu      sync.RWMutex
	current string
}

func (s *CurrentState) loadCurrentState() string {
	s.mu.RLock()
	v := s.current
	s.mu.RUnlock()
	return v
}
func (s *CurrentState) Load() string   { return s.loadCurrentState() }
func (s *CurrentState) Store(v string) { s.mu.Lock(); s.current = v; s.mu.Unlock() }
