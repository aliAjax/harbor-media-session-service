package application

import (
	"context"
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
func Run(ctx context.Context, conn *net.UDPConn, accept func(domain.Packet)) error {
	b := make([]byte, 4096)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		n, _, e := conn.ReadFromUDP(b)
		if e != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if ne, ok := e.(net.Error); ok && ne.Timeout() {
				continue
			}
			return e
		}
		if p, e := domain.Parse(b[:n]); e == nil {
			accept(p)
		}
	}
}
