package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/cluster/domain"
	"time"
)

type Coordinator struct {
	store *domain.Store
	node  string
	ttl   time.Duration
}

func NewCoordinator(store *domain.Store, node string, ttl time.Duration) *Coordinator {
	return &Coordinator{store: store, node: node, ttl: ttl}
}
func (c *Coordinator) Claim(ctx context.Context, room string) (domain.Lease, error) {
	select {
	case <-ctx.Done():
		return domain.Lease{}, ctx.Err()
	default:
	}
	l, e := c.store.Acquire(room, c.node, c.ttl)
	if e != nil {
		return l, fmt.Errorf("claim room: %w", e)
	}
	return l, nil
}
func (c *Coordinator) Renew(ctx context.Context, l domain.Lease) (domain.Lease, error) {
	select {
	case <-ctx.Done():
		return l, ctx.Err()
	default:
	}
	n, e := c.store.Renew(l, c.ttl)
	if e != nil {
		return l, fmt.Errorf("renew room: %w", e)
	}
	return n, nil
}
