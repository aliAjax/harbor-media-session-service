package application

import (
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"sync/atomic"
)

type Stats struct {
	Packets      atomic.Uint64
	Bytes        atomic.Uint64
	Lost         atomic.Uint64
	LastSequence atomic.Uint32
}

func (s *Stats) Observe(p domain.Packet) {
	s.Packets.Add(1)
	s.Bytes.Add(uint64(len(p.Payload) + 12))
	s.LastSequence.Store(uint32(p.Sequence))
}
func (s *Stats) MarkLost(n uint64) { s.Lost.Add(n) }
func (s *Stats) Snapshot() map[string]uint64 {
	return map[string]uint64{"packets": s.Packets.Load(), "bytes": s.Bytes.Load(), "lost": s.Lost.Load(), "last_sequence": uint64(s.LastSequence.Load())}
}
