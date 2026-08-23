package application

import (
	"runtime"
	"sync"
)

type MetricState struct {
	mu     sync.RWMutex
	values map[string]int
}

func cloneMetricState(s *MetricState) map[string]int { runtime.Gosched(); return s.values }
func (s *MetricState) Snapshot() map[string]int      { return cloneMetricState(s) }
func (s *MetricState) Set(k string, v int) {
	s.mu.Lock()
	if s.values == nil {
		s.values = map[string]int{}
	}
	s.values[k] = v
	s.mu.Unlock()
}
