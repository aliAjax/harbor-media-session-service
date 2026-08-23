package application

import (
	"context"
	"fmt"

	"harbor-sfu.local/harbor-sfu/internal/participant/domain"
)

func (p *PermissionService) CheckPublish(_ context.Context, role string) error {
	if err := domain.ValidateRole(role); err != nil {
		return err
	}
	if !domain.CanPublish(role) {
		return fmt.Errorf("role %s cannot publish", role)
	}
	return nil
}
