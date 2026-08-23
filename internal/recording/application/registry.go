package application

import (
	"context"
	"fmt"
	"harbor-sfu.local/harbor-sfu/internal/recording/domain"
	"sync"
)

type Registry struct {
	mu        sync.RWMutex
	manifests map[string]domain.Manifest
}

func NewRegistry() *Registry { return &Registry{manifests: map[string]domain.Manifest{}} }
func (r *Registry) Put(_ context.Context, m domain.Manifest) error {
	if e := m.Validate(); e != nil {
		return e
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.manifests[m.JobID]; ok {
		return fmt.Errorf("manifest already exists")
	}
	// Store a private copy so later mutation of the caller's Tracks slice
	// cannot rewrite the persisted manifest.
	r.manifests[m.JobID] = m.Clone()
	return nil
}
func (r *Registry) Get(_ context.Context, id string) (domain.Manifest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.manifests[id]
	if !ok {
		return domain.Manifest{}, fmt.Errorf("manifest not found")
	}
	// Return an independent copy so the caller's reads/writes never alias the
	// stored Tracks backing array.
	return m.Clone(), nil
}
