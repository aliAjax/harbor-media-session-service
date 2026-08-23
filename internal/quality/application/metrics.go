package application

import (
	"harbor-sfu.local/harbor-sfu/internal/quality/domain"
	"sync"
)

type MetricStore struct {
	mu    sync.RWMutex
	items map[string]domain.Metrics
}

func NewMetricStore() *MetricStore { return &MetricStore{items: map[string]domain.Metrics{}} }
func (m *MetricStore) Observe(id string, delta domain.Metrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.items[id]
	v.Packets += delta.Packets
	v.Bytes += delta.Bytes
	v.Lost += delta.Lost
	if delta.Jitter > 0 {
		v.Jitter = delta.Jitter
	}
	if delta.RTTMillis > 0 {
		v.RTTMillis = delta.RTTMillis
	}
	m.items[id] = v
}
func (m *MetricStore) Snapshot() map[string]domain.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o := map[string]domain.Metrics{}
	for k, v := range m.items {
		o[k] = v
	}
	return o
}
