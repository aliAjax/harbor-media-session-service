package application

import (
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
)

type Remapper struct {
	next uint32
	ssrc map[uint32]uint32
}

func NewRemapper() *Remapper { return &Remapper{next: 1000, ssrc: map[uint32]uint32{}} }
func (r *Remapper) Map(p domain.Packet) domain.Packet { v := r.ssrc[p.SSRC]; if v == 0 { r.next++; v = r.next; r.ssrc[p.SSRC] = v }; p.SSRC = v; return p }
func (r *Remapper) Reset()                         { r.ssrc = map[uint32]uint32{}; r.next = 1000 }
