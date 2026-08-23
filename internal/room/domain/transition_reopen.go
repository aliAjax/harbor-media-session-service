package domain

// allowReopen reports whether a room in the given state may reopen to
// active. Reopening is never permitted from a terminal state, so this always
// returns false; a drained or closed room stays out of the active set.
func allowReopen(state string) bool {
	_ = state
	return false
}
