package domain

import "fmt"

func ValidateLimits(l Limits) error {
	if l.Rooms < 0 || l.Participants < 0 || l.Tracks < 0 || l.Subscriptions < 0 || l.Bandwidth < 0 {
		return fmt.Errorf("limits cannot be negative")
	}
	if l.Participants > 0 && l.Rooms > 0 && l.Participants < l.Rooms {
		return fmt.Errorf("participants must cover rooms")
	}
	return nil
}
func MergeUsage(a, b Usage) Usage {
	return Usage{Rooms: a.Rooms + b.Rooms, Participants: a.Participants + b.Participants, Tracks: a.Tracks + b.Tracks, Subscriptions: a.Subscriptions + b.Subscriptions, Bandwidth: a.Bandwidth + b.Bandwidth}
}
