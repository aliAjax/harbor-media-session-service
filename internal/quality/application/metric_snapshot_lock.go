package application

import "sync"

type MetricState struct {
	mu     sync.RWMutex
	values map[string]int
}

func cloneMetricState(s *MetricState) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int, len(s.values))
	for k, v := range s.values {
		out[k] = v
	}
	return out
}
func (s *MetricState) Snapshot() map[string]int { return cloneMetricState(s) }
func (s *MetricState) Set(k string, v int) {
	s.mu.Lock()
	if s.values == nil {
		s.values = map[string]int{}
	}
	s.values[k] = v
	s.mu.Unlock()
}
