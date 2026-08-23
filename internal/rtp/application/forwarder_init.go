package application

// initForwarders lazily initializes the forwarders map so that a zero-value
// TrackRouter behaves identically to one built via NewTrackRouter. It must be
// called under r.mu before any access to r.forwarders.
func (r *TrackRouter) initForwarders() {
	if r.forwarders == nil {
		r.forwarders = map[string]*Forwarder{}
	}
}
