package application

import (
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/signaling/domain"
	"sync"
)

type Session struct {
	ID    string
	State domain.State
	Send  chan domain.Message
	Done  chan struct{}
}
type Hub struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewHub() *Hub { return &Hub{sessions: map[string]*Session{}} }
func (h *Hub) Add(s *Session) error {
	if _, ok := h.sessions[s.ID]; ok {
		return fmt.Errorf("session exists")
	}
	h.sessions[s.ID] = s
	return nil
}
func (h *Hub) Remove(id string) {
	if s, ok := h.sessions[id]; ok {
		close(s.Done)
		delete(h.sessions, id)
	}
}
func (h *Hub) Broadcast(room string, m domain.Message, except string) {
	for id, s := range h.sessions {
		if id != except && s.State.RoomID == room {
			select {
			case s.Send <- m:
			default:
			}
		}
	}
}
func (h *Hub) Count() int { return len(h.sessions) }
