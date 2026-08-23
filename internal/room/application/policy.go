package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/room/domain"
)

type PolicyStore interface {
	Get(context.Context, string) (domain.Policy, error)
	Set(context.Context, string, domain.Policy) error
}
type MemoryPolicies struct{ values map[string]domain.Policy }

func NewMemoryPolicies() *MemoryPolicies { return &MemoryPolicies{values: map[string]domain.Policy{}} }
func (m *MemoryPolicies) Get(_ context.Context, t string) (domain.Policy, error) {
	p, ok := m.values[t]
	if !ok {
		return domain.Policy{AllowPublish: true, AllowSubscribe: true, MaxTracks: 32}, nil
	}
	return p, nil
}
func (m *MemoryPolicies) Set(_ context.Context, t string, p domain.Policy) error {
	if p.MaxTracks < 0 {
		return fmt.Errorf("negative track limit")
	}
	m.values[t] = p
	return nil
}
