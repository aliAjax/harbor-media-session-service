package application

import (
	"context"
	"sync"
	"time"
)

type AuditEvent struct {
	ID      string
	RoomID  string
	Action  string
	Actor   string
	At      time.Time
	Details map[string]string
}
type AuditLog struct {
	mu    sync.RWMutex
	items []AuditEvent
}

func NewAuditLog() *AuditLog { return &AuditLog{items: make([]AuditEvent, 0, 128)} }
func (a *AuditLog) Append(_ context.Context, e AuditEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.items = append(a.items, e)
}
func (a *AuditLog) List(_ context.Context, room string, limit int) []AuditEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if limit <= 0 || limit > len(a.items) {
		limit = len(a.items)
	}
	out := make([]AuditEvent, 0, limit)
	for i := len(a.items) - 1; i >= 0 && len(out) < limit; i-- {
		if room == "" || a.items[i].RoomID == room {
			out = append(out, a.items[i])
		}
	}
	return out
}
