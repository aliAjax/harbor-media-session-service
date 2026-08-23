package application

import (
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"net"
	"sync"
	"time"
)

type Session struct {
	ID        string
	Remote    *net.UDPAddr
	CreatedAt time.Time
	Packets   uint64
	mu        sync.Mutex
	closed    bool
}

func (s *Session) Accept(p domain.Packet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("transport session closed")
	}
	s.Packets++
	return nil
}
func (s *Session) Close() { s.mu.Lock(); defer s.mu.Unlock(); s.closed = true }
