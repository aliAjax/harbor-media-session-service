package infrastructure

import (
	"context"
	"harbor-sfu.local/harbor-sfu/internal/room/domain"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	rooms map[string]*domain.Room
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{rooms: map[string]*domain.Room{}}
}
func (m *MemoryRepository) Create(_ context.Context, r *domain.Room) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.rooms[r.ID]; ok {
		return domain.ErrConflict
	}
	m.rooms[r.ID] = r
	return nil
}
func (m *MemoryRepository) Get(_ context.Context, id string) (*domain.Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rooms[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return r, nil
}
func (m *MemoryRepository) List(_ context.Context, t string) ([]domain.Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o := []domain.Room{}
	for _, r := range m.rooms {
		if t == "" || r.TenantID == t {
			o = append(o, r.Snapshot())
		}
	}
	return o, nil
}
func (m *MemoryRepository) Save(_ context.Context, r *domain.Room) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.rooms[r.ID]; !ok {
		return domain.ErrNotFound
	}
	m.rooms[r.ID] = r
	return nil
}
