package domain

import "fmt"

func ValidateRole(role string) error {
	switch role {
	case "publisher", "subscriber", "moderator", "observer":
		return nil
	default:
		return fmt.Errorf("unsupported participant role %q", role)
	}
}
func CanPublish(role string) bool   { return role == "publisher" || role == "moderator" }
func CanSubscribe(role string) bool { return role != "observer" }
