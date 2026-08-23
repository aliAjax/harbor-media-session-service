package application

import (
	"harbor-sfu.local/harbor-sfu/internal/rtp/domain"
	"sync"
)

type TrackRouter struct {
	mu         sync.RWMutex
	forwarders map[string]*Forwarder
	remap      map[string]*Remapper
	stats      map[string]*Stats
}

func NewTrackRouter() *TrackRouter {
	return &TrackRouter{forwarders: map[string]*Forwarder{}, remap: map[string]*Remapper{}, stats: map[string]*Stats{}}
}
func (r *TrackRouter) Register(id string, f *Forwarder) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.initForwarders()
	r.initRemap()
	r.initStats()
	r.forwarders[id] = f
	r.remap[id] = NewRemapper()
	r.stats[id] = &Stats{}
}
func (r *TrackRouter) Unregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.initForwarders()
	r.initRemap()
	r.initStats()
	if f := r.forwarders[id]; f != nil {
		f.Close()
	}
	delete(r.forwarders, id)
	delete(r.remap, id)
	delete(r.stats, id)
}
func (r *TrackRouter) Publish(id string, p domain.Packet) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Reading a nil map yields the zero value, so a zero-value TrackRouter
	// without a prior Register returns false here without panic. Map mutation
	// is confined to the write-locked Register/Unregister paths.
	f := r.forwarders[id]
	if f == nil {
		return false
	}
	mapped := r.remap[id].Map(p)
	r.stats[id].Observe(mapped)
	f.Publish(mapped)
	return true
}
func (r *TrackRouter) Stats(id string) map[string]uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s := r.stats[id]; s != nil {
		return s.Snapshot()
	}
	return map[string]uint64{}
}
