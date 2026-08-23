package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/recording/domain"
	"sync"
	"time"
)

type Service struct {
	mu   sync.RWMutex
	jobs map[string]*domain.Job
}

func NewService() *Service { return &Service{jobs: map[string]*domain.Job{}} }
func (s *Service) Start(_ context.Context, id, room string) (*domain.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[id]; ok {
		return nil, fmt.Errorf("recording exists")
	}
	j := &domain.Job{ID: id, RoomID: room, Status: domain.Starting, StartedAt: time.Now(), UpdatedAt: time.Now()}
	s.jobs[id] = j
	if e := j.Transition(domain.Running); e != nil {
		return nil, e
	}
	return j, nil
}
func (s *Service) Get(_ context.Context, id string) (*domain.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	if !ok {
		return nil, fmt.Errorf("recording not found")
	}
	cp := *j
	return &cp, nil
}
func (s *Service) Stop(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok {
		return fmt.Errorf("recording not found")
	}
	if e := j.Transition(domain.Stopping); e != nil {
		return e
	}
	return j.Transition(domain.Completed)
}
