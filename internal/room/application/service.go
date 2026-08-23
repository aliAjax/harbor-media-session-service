package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/room/domain"
	"time"
)

type Repository interface {
	Create(context.Context, *domain.Room) error
	Get(context.Context, string) (*domain.Room, error)
	List(context.Context, string) ([]domain.Room, error)
	Save(context.Context, *domain.Room) error
}
type Service struct {
	repo  Repository
	clock func() time.Time
}

func NewService(r Repository) *Service { return &Service{repo: r, clock: time.Now} }
func (s *Service) Create(ctx context.Context, tenant, id, name, node string) (*domain.Room, error) {
	now := s.clock()
	r := &domain.Room{ID: id, TenantID: tenant, Name: name, Status: domain.Open, NodeID: node, CreatedAt: now, UpdatedAt: now, Participants: map[string]*domain.Participant{}, Tracks: map[string]*domain.Track{}}
	if e := s.repo.Create(ctx, r); e != nil {
		return nil, fmt.Errorf("create room: %w", e)
	}
	return r, nil
}
func (s *Service) Get(ctx context.Context, id string) (*domain.Room, error) {
	r, e := s.repo.Get(ctx, id)
	if e != nil {
		return nil, fmt.Errorf("get room: %w", e)
	}
	return r, nil
}
func (s *Service) List(ctx context.Context, t string) ([]domain.Room, error) {
	r, e := s.repo.List(ctx, t)
	if e != nil {
		return nil, fmt.Errorf("list rooms: %w", e)
	}
	return r, nil
}
func (s *Service) Transition(ctx context.Context, id string, status domain.Status) error {
	r, e := s.Get(ctx, id)
	if e != nil {
		return e
	}
	// Reject anything that is not a defined room status before consulting the
	// transition table — unknown / legacy aliases (e.g. "active") cannot leak
	// in through this entry point. The table itself then enforces the linear,
	// terminal lifecycle: Open -> Draining -> Closed (plus Open -> Closed),
	// with no reopen from a drained or closed room.
	if !domain.AllowedTransition(r.Status, status) {
		return fmt.Errorf("%w: %s -> %s", domain.ErrInvalidTransition, r.Status, status)
	}
	r.Status = status
	r.UpdatedAt = s.clock()
	if e = s.repo.Save(ctx, r); e != nil {
		return fmt.Errorf("save room: %w", e)
	}
	return nil
}
