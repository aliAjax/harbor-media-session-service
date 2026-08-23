package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/signaling/domain"
	"time"
)

type Dispatcher struct {
	hub     *Hub
	timeout time.Duration
}

func NewDispatcher(h *Hub, d time.Duration) *Dispatcher { return &Dispatcher{hub: h, timeout: d} }
func (d *Dispatcher) Dispatch(ctx context.Context, s *Session, m domain.Message) error {
	if !s.State.Accept(m) {
		return wrapDispatchError(ErrStaleMessage)
	}
	m.RoomID = s.State.RoomID
	select {
	case s.Send <- m:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d.timeout):
		return fmt.Errorf("send timeout")
	}
}
func (d *Dispatcher) Join(s *Session, room, pid string) {
	s.State.RoomID = room
	s.State.ParticipantID = pid
	s.State.Connected = true
}
