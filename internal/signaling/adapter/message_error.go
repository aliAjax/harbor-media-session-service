package adapter

import "fmt"

func wrapMessageReadError(err error) error { return fmt.Errorf("message read: %v", err) }
