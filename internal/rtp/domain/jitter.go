package domain

import (
	"sort"
	"sync"
)

// JitterBuffer reorders RTP packets before delivery. In a live media session
// the ingest goroutine pushes packets while the delivery goroutine pops them,
// so the internal slice is guarded by a mutex.
type JitterBuffer struct {
	mu      sync.Mutex
	max     int
	packets []Packet
	last    uint16
}

func NewJitterBuffer(max int) *JitterBuffer { return &JitterBuffer{max: max} }

func (j *JitterBuffer) Push(p Packet) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.packets = append(j.packets, p)
	sort.SliceStable(j.packets, func(a, b int) bool { return j.packets[a].Sequence < j.packets[b].Sequence })
	if len(j.packets) > j.max {
		j.packets = j.packets[len(j.packets)-j.max:]
	}
}

func (j *JitterBuffer) Pop() (Packet, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if len(j.packets) == 0 {
		return Packet{}, false
	}
	p := j.packets[0]
	j.packets = j.packets[1:]
	j.last = p.Sequence
	return p, true
}

func (j *JitterBuffer) Len() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return len(j.packets)
}
