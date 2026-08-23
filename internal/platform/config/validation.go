package config

import (
	"fmt"
	"net"
)

func Validate(c Config) error {
	if c.HTTPAddr == "" || c.UDPAddr == "" {
		return fmt.Errorf("listen addresses are required")
	}
	if _, e := net.ResolveTCPAddr("tcp", c.HTTPAddr); e != nil {
		return fmt.Errorf("http address: %w", e)
	}
	if _, e := net.ResolveUDPAddr("udp", c.UDPAddr); e != nil {
		return fmt.Errorf("udp address: %w", e)
	}
	if c.MaxRooms < 1 || c.MaxParticipants < 1 {
		return fmt.Errorf("limits must be positive")
	}
	return nil
}
