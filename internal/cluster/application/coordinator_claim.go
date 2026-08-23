package application

import (
	"context"
	"fmt"

	"harbor-sfu.local/harbor-sfu/internal/cluster/domain"
)

func (c *Coordinator) Claim(ctx context.Context, room string) (domain.Lease, error) {
	// Honor the control request's lifecycle: a claim whose context has already
	// been cancelled must not take the lease, otherwise a room is seized for a
	// request the caller has abandoned. A live context is unaffected.
	if err := ctx.Err(); err != nil {
		return domain.Lease{}, fmt.Errorf("claim room: %w", err)
	}
	lease, err := c.store.Acquire(room, c.node, c.ttl)
	if err != nil {
		return lease, fmt.Errorf("claim room: %w", err)
	}
	return lease, nil
}
