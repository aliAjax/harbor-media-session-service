package application

import (
	"context"
	"fmt"

	"harbor-sfu.local/harbor-sfu/internal/cluster/domain"
)

func (c *Coordinator) Renew(ctx context.Context, lease domain.Lease) (domain.Lease, error) {
	// Honor the control request's lifecycle: a renewal whose context has
	// already been cancelled must not extend the lease. A live context proceeds
	// normally.
	if err := ctx.Err(); err != nil {
		return lease, fmt.Errorf("renew room: %w", err)
	}
	renewed, err := c.store.Renew(lease, c.ttl)
	if err != nil {
		return lease, fmt.Errorf("renew room: %w", err)
	}
	return renewed, nil
}
