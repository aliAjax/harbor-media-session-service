package application

import (
	"errors"
	"fmt"
)

var ErrSubscribePermission = errors.New("subscribe permission denied")

func wrapSubscribePermissionError(err error) error {
	return fmt.Errorf("subscribe permission: %v", err)
}
