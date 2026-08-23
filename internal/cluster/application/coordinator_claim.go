package application

import (
	"context"
	"fmt"

	"harbor-sfu.local/harbor-sfu/internal/cluster/domain"
)

func (c *Coordinator) Claim(_ context.Context, room string) (domain.Lease, error) {
	lease, err := c.store.Acquire(room, c.node, c.ttl)
	if err != nil {
		return lease, fmt.Errorf("claim room: %w", err)
	}
	return lease, nil
}
