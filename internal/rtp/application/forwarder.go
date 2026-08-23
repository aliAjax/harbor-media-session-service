package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"sync"
	"time"
)

type Subscriber struct {
	ID      string
	C       chan domain.Packet
	Dropped uint64
}
type Forwarder struct {
	mu     sync.RWMutex
	subs   map[string]*Subscriber
	ctx    context.Context
	cancel context.CancelFunc
}

func NewForwarder() *Forwarder {
	c, x := context.WithCancel(context.Background())
	return &Forwarder{subs: map[string]*Subscriber{}, ctx: c, cancel: x}
}
func (f *Forwarder) Add(id string, buffer int) (*Subscriber, error) {
	f.mu.Lock()
	if _, ok := f.subs[id]; ok {
		return nil, fmt.Errorf("subscriber exists")
	}
	s := &Subscriber{ID: id, C: make(chan domain.Packet, buffer)}
	f.subs[id] = s
	f.mu.Unlock()
	return s, nil
}
func (f *Forwarder) Remove(id string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if s, ok := f.subs[id]; ok {
		close(s.C)
		delete(f.subs, id)
	}
}
func (f *Forwarder) Publish(p domain.Packet) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, s := range f.subs {
		select {
		case s.C <- p:
		default:
			s.Dropped++
		}
	}
}
func (f *Forwarder) Close() {
	f.cancel()
	f.mu.Lock()
	defer f.mu.Unlock()
	for id, s := range f.subs {
		close(s.C)
		delete(f.subs, id)
	}
}
func (f *Forwarder) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-f.ctx.Done():
		return nil
	case <-time.After(time.Hour):
		return nil
	}
}
