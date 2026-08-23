package adapter

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"io"
	"sync"
)

type Sink interface {
	WritePacket(context.Context, domain.Packet) error
	Close() error
}
type MemorySink struct {
	mu      sync.Mutex
	packets []domain.Packet
	closed  bool
}

func NewMemorySink() *MemorySink { return &MemorySink{} }
func (s *MemorySink) WritePacket(ctx context.Context, p domain.Packet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("sink closed")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.packets = append(s.packets, p)
	return nil
}
func (s *MemorySink) Close() error { s.mu.Lock(); defer s.mu.Unlock(); s.closed = true; return nil }
func (s *MemorySink) Count() int   { s.mu.Lock(); defer s.mu.Unlock(); return len(s.packets) }
func Copy(ctx context.Context, r io.Reader, s Sink) error {
	b := make([]byte, 2048)
	for {
		n, e := r.Read(b)
		if n > 0 {
			p, e2 := domain.Parse(b[:n])
			if e2 == nil {
				if e2 = s.WritePacket(ctx, p); e2 != nil {
					return e2
				}
			}
		}
		if e == io.EOF {
			return nil
		}
		if e != nil {
			return e
		}
	}
}
