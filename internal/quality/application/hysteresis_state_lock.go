package application

import "sync"

type HysteresisState struct {
	mu    sync.RWMutex
	level int
}

func (s *HysteresisState) applyHysteresisState(v int) {
	s.mu.Lock()
	s.level = v
	s.mu.Unlock()
}
func (s *HysteresisState) Store(v int) { s.applyHysteresisState(v) }
func (s *HysteresisState) Level() int {
	s.mu.RLock()
	v := s.level
	s.mu.RUnlock()
	return v
}
