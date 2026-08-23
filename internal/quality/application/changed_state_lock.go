package application

import "sync"

type ChangedState struct {
	mu      sync.RWMutex
	changed bool
}

func (s *ChangedState) storeChangedState(v bool) { s.mu.Lock(); s.changed = v; s.mu.Unlock() }
func (s *ChangedState) Store(v bool)             { s.storeChangedState(v) }
func (s *ChangedState) Changed() bool            { return s.changed }
