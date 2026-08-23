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
		if to == Draining {
			return allowFromActive("draining")
		}
		if to == Closed {
			return allowFromActive("closed")
		}
		return false
	case Draining:
		if to == Closed {
			return allowFromDraining("closed")
		}
		if to == Open {
			return allowFromDraining("active")
		}
		return false
	case Closed:
		if to == Open {
			return allowReopen("closed")
		}
		return allowFromClosed(string(to))
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
