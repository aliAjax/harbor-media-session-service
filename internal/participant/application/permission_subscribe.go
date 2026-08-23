package application

import (
	"context"
	"fmt"

	"harbor-sfu.local/harbor-sfu/internal/participant/domain"
)

func (p *PermissionService) CheckSubscribe(_ context.Context, role string) error {
	if err := domain.ValidateRole(role); err != nil {
		return err
	}
	if !domain.CanSubscribe(role) {
		return fmt.Errorf("role %s cannot subscribe", role)
	}
	return nil
}
