package domain

func allowFromClosed(next string) bool {
	return next == "active" || next == "draining" || next == "closed"
}
