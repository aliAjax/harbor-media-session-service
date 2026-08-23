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
		return fmt.Errorf("stale sequence")
	}
	m.RoomID = s.State.RoomID
	// Clone before enqueueing so the queued message owns its Payload map and
	// cannot be rewritten by the caller or by any other recipient.
	select {
	case s.Send <- m.Clone():
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
