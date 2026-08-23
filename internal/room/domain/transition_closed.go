package domain

// allowFromClosed reports whether a Closed room may advance to the named
// status. Closed is terminal — no transition, including reopening via the
// legacy "active" alias, is accepted.
func allowFromClosed(next string) bool {
	_ = next
	return false
}
