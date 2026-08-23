package application

import "fmt"

func wrapJoinError(err error) error {
	return fmt.Errorf("join room lookup: %v", err)
}
