package domain

import (
	"fmt"
	"time"
)

type Transition struct {
	From   Status
	To     Status
	At     time.Time
	Actor  string
	Reason string
}

func AllowedTransition(from, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case Open:
		return to == Draining || to == Closed
	case Draining:
		return to == Closed
	case Closed:
		return false
	}
	return false
}
func ExplainTransition(from, to Status) error {
	if !AllowedTransition(from, to) {
		return fmt.Errorf("transition %s -> %s is not allowed", from, to)
	}
	return nil
}
func NextStatus(s Status) Status {
	switch s {
	case Open:
		return Draining
	case Draining:
		return Closed
	default:
		return Closed
	}
}
