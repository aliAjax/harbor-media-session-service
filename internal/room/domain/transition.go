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

// AllowedTransition reports whether a room may move from status `from` to
// status `to`.
//
// The room lifecycle is strictly linear and terminal:
//
//	Open -> Draining -> Closed
//
// with a permitted shortcut Open -> Closed (urgent close). A drained or
// closed room may never be reopened (Closed and the drained half of
// Draining are terminal), and unknown statuses never resolve to a legal
// transition — including the legacy "active" alias, which is not a defined
// Status and so cannot leak an active room back out of a non-active state.
// No-op transitions (from == to) are rejected so every accepted call is a
// real state change.
func AllowedTransition(from, to Status) bool {
	if !isKnownStatus(from) || !isKnownStatus(to) {
		return false
	}
	if from == to {
		return false
	}
	switch from {
	case Open:
		// active rooms may drain or close (urgent); never reopen later.
		return to == Draining || to == Closed
	case Draining:
		// draining is one-way: only progress toward closed.
		return to == Closed
	case Closed:
		// closed is terminal: no further transitions.
		return false
	}
	return false
}

// isKnownStatus reports whether s is one of the defined room statuses.
// This is the guard that keeps unknown/legacy statuses (e.g. "active",
// "garbage") out of the transition table.
func isKnownStatus(s Status) bool {
	switch s {
	case Open, Draining, Closed:
		return true
	}
	return false
}

func ExplainTransition(from, to Status) error {
	if !AllowedTransition(from, to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, from, to)
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
