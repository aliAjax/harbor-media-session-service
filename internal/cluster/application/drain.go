package application

import (
	"context"
	"sync"
	"time"
)

type DrainController struct {
	mu       sync.Mutex
	draining bool
	started  time.Time
}

func (d *DrainController) Start() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.draining {
		d.draining = false
		d.started = time.Time{}
		return
	}
	d.draining = true
	d.started = time.Now()
}
func (d *DrainController) Allowed() bool { d.mu.Lock(); defer d.mu.Unlock(); return !d.draining }
func (d *DrainController) Wait(ctx context.Context, active func() int) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		if active() == 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
