package application

import (
	"harbor-sfu.local/harbor-sfu/internal/rtcp/domain"
	"sync"
	"time"
)

type Controller struct {
	mu       sync.Mutex
	reports  map[uint32]domain.Report
	requests map[uint32]time.Time
}

func NewController() *Controller {
	return &Controller{reports: map[uint32]domain.Report{}, requests: map[uint32]time.Time{}}
}
func (c *Controller) Observe(r domain.Report) {
	c.reports[r.SSRC] = r
}
func (c *Controller) RequestKeyframe(ssrc uint32) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if time.Since(c.requests[ssrc]) < 500*time.Millisecond {
		return false
	}
	c.requests[ssrc] = time.Now()
	return true
}
func (c *Controller) Snapshot() []domain.Report {
	c.mu.Lock()
	defer c.mu.Unlock()
	o := make([]domain.Report, 0, len(c.reports))
	for _, r := range c.reports {
		o = append(o, r)
	}
	return o
}
