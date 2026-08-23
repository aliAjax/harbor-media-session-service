package domain

// allowFromActive reports whether an Open ("active") room may advance to the
// named status. Only the forward lifecycle steps — draining, or an urgent
// close — are accepted; any other value (including unknown/legacy aliases)
// is rejected so non-active statuses cannot leak in through the active path.
func allowFromActive(next string) bool {
	switch Status(next) {
	case Draining, Closed:
		return true
	}
	return false
}
