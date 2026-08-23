package application

import (
	"context"
	"harbor-sfu.local/harbor-sfu/internal/participant/domain"
)

type PermissionService struct{}

func NewPermissionService() *PermissionService { return &PermissionService{} }
func (p *PermissionService) CheckPublish(_ context.Context, role string) error {
	if e := domain.ValidateRole(role); e != nil {
		return e
	}
	if !domain.CanPublish(role) {
		return wrapPublishPermissionError(ErrPublishPermission)
	}
	return nil
}
func (p *PermissionService) CheckSubscribe(_ context.Context, role string) error {
	if e := domain.ValidateRole(role); e != nil {
		return e
	}
	if !domain.CanSubscribe(role) {
		return wrapSubscribePermissionError(ErrSubscribePermission)
	}
	return nil
}
