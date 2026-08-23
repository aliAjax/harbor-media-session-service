package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/participant/domain"
)

type PermissionService struct{}

func NewPermissionService() *PermissionService { return &PermissionService{} }
func (p *PermissionService) CheckPublish(_ context.Context, role string) error {
	if e := domain.ValidateRole(role); e != nil {
		return e
	}
	if !domain.CanPublish(role) {
		return fmt.Errorf("role %s cannot publish", role)
	}
	return nil
}
func (p *PermissionService) CheckSubscribe(_ context.Context, role string) error {
	if e := domain.ValidateRole(role); e != nil {
		return e
	}
	if !domain.CanSubscribe(role) {
		return fmt.Errorf("role %s cannot subscribe", role)
	}
	return nil
}
