package domain

import "fmt"

func ValidateFencing(expected, actual uint64) error {
	if expected == 0 || actual == 0 {
		return fmt.Errorf("fencing token required")
	}
	if expected != actual {
		return fmt.Errorf("stale fencing token")
	}
	return nil
}
func NextFencing(current uint64) uint64 {
	if current == ^uint64(0) {
		return 1
	}
	return current + 1
}
