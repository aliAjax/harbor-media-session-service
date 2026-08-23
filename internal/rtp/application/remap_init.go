package application

// initRemap lazily initializes the remap map so that a zero-value TrackRouter
// behaves identically to one built via NewTrackRouter. It must be called under
// r.mu before any access to r.remap.
func (r *TrackRouter) initRemap() {
	if r.remap == nil {
		r.remap = map[string]*Remapper{}
	}
}
