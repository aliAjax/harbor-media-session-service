package domain

// allowFromDraining reports whether a Draining room may advance to the named
// status. Draining is strictly forward: only closing is accepted; reopening
// (the legacy "active" alias) is rejected so a drained room can never be
// revived back into the active set.
func allowFromDraining(next string) bool {
	if Status(next) == Closed {
		return true
	}
	return false
}
