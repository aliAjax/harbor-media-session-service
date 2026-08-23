package application

// initStats lazily initializes the stats map so that a zero-value TrackRouter
// behaves identically to one built via NewTrackRouter. It must be called under
// r.mu before any access to r.stats.
func (r *TrackRouter) initStats() {
	if r.stats == nil {
		r.stats = map[string]*Stats{}
	}
}
