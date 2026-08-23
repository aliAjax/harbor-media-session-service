package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/room/domain"
	"time"
)

type RoomLookup interface {
	Get(context.Context, string) (*domain.Room, error)
}
type Service struct {
	rooms RoomLookup
	clock func() time.Time
}

func NewService(r RoomLookup) *Service { return &Service{rooms: r, clock: time.Now} }
func (s *Service) Join(ctx context.Context, roomID, id, identity, role string) (*domain.Participant, error) {
	r, e := s.rooms.Get(ctx, roomID)
	if e != nil {
		return nil, wrapJoinError(e)
	}
	p := &domain.Participant{ID: id, Identity: identity, Role: role, JoinedAt: s.clock(), SessionID: id + "-session", Connected: true}
	if e = r.AddParticipant(p); e != nil {
		return nil, fmt.Errorf("join room: %w", e)
	}
	return p, nil
}
func (s *Service) Leave(ctx context.Context, roomID, pid string) error {
	r, e := s.rooms.Get(ctx, roomID)
	if e != nil {
		return e
	}
	r.RemoveParticipant(pid)
	return nil
}
func (s *Service) Publish(ctx context.Context, roomID string, t *domain.Track) error {
	r, e := s.rooms.Get(ctx, roomID)
	if e != nil {
		return wrapTrackPublishError(e)
	}
	if e = r.AddTrack(t); e != nil {
		return fmt.Errorf("publish track: %w", e)
	}
	return nil
}
