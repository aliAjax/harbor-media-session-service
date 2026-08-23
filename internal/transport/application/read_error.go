package application

import "fmt"

func wrapUDPReadError(err error) error { return fmt.Errorf("udp read: %w", err) }
