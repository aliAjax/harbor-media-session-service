package application

import (
	"context"
	"fmt"

	"harbor-sfu.local/harbor-sfu/internal/cluster/domain"
)

func (c *Coordinator) Renew(_ context.Context, lease domain.Lease) (domain.Lease, error) {
	renewed, err := c.store.Renew(lease, c.ttl)
	if err != nil {
		return lease, fmt.Errorf("renew room: %w", err)
	}
	return renewed, nil
}
