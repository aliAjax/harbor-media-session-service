package application

import (
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"sync"
)

// Remapper translates publisher SSRCs into per-track SSRCs assigned by the
// SFU. A single Remapper is shared by the packet ingest path and the cleanup
// path of one media session, so every accessor must hold the lock.
type Remapper struct {
	mu   sync.Mutex
	next uint32
	ssrc map[uint32]uint32
}

func NewRemapper() *Remapper { return &Remapper{next: 1000, ssrc: map[uint32]uint32{}} }

func (r *Remapper) Map(p domain.Packet) domain.Packet {
	r.mu.Lock()
	defer r.mu.Unlock()
	v := r.ssrc[p.SSRC]
	if v == 0 {
		r.next++
		v = r.next
		r.ssrc[p.SSRC] = v
	}
	p.SSRC = v
	return p
}

func (r *Remapper) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ssrc = map[uint32]uint32{}
	r.next = 1000
}
