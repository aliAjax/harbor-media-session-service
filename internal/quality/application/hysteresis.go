package application

import "time"

type Hysteresis struct {
	last     string
	changed  time.Time
	cooldown time.Duration
}

func NewHysteresis(d time.Duration) *Hysteresis { return &Hysteresis{cooldown: d} }
func (h *Hysteresis) Apply(next string) string {
	if h.last == "" {
		h.last = next
		h.changed = time.Now()
		return next
	}
	if next != h.last && time.Since(h.changed) >= h.cooldown {
		h.last = next
		h.changed = time.Now()
	}
	return h.last
}
