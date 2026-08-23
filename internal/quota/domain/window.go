package domain

import (
	"sync"
	"time"
)

type Window struct {
	mu       sync.Mutex
	start    time.Time
	bytes    int64
	limit    int64
	duration time.Duration
}

func NewWindow(limit int64, d time.Duration) *Window {
	return &Window{start: time.Now(), limit: limit, duration: d}
}
func (w *Window) Allow(n int64) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	if now.Sub(w.start) >= w.duration {
		w.start = now
		w.bytes = 0
	}
	if w.limit > 0 && w.bytes+n > w.limit {
		return false
	}
	w.bytes += n
	return true
}
func (w *Window) Usage() int64 { return w.bytes }
