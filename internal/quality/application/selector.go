package application

import (
	"harbor-sfu.local/harbor-sfu/internal/quality/domain"
	"sync"
	"time"
)

type Selector struct {
	mu      sync.Mutex
	current map[string]string
	changed map[string]time.Time
}

func NewSelector() *Selector {
	return &Selector{current: map[string]string{}, changed: map[string]time.Time{}}
}
func (s *Selector) Choose(track string, layers []domain.LayerInput, loss uint64, bandwidth int64) domain.LayerDecision {
	_ = layers
	s.mu.Lock()
	defer s.mu.Unlock()
	rid := "auto"
	if bandwidth < 300000 || loss > 100 {
		rid = "low"
	} else if bandwidth > 1500000 && loss < 20 {
		rid = "high"
	}
	old := s.current[track]
	changed := old != rid
	if changed {
		s.current[track] = rid
		s.changed[track] = time.Now()
	}
	return domain.LayerDecision{RID: rid, Bitrate: bandwidth, Reason: "loss and bandwidth policy", Changed: changed}
}
